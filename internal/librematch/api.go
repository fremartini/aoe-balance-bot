package librematch

import (
	"aoe-bot/internal/cache"
	"aoe-bot/internal/domain"
	"aoe-bot/internal/errors"
	"aoe-bot/internal/list"
	"aoe-bot/internal/logger"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type api struct {
	logger *logger.Logger
	cache  *cache.Cache[uint, *domain.Player]
}

const (
	LEADERBOARD_1V1      = 3
	LEADERBOARD_TEAMGAME = 4
	BASE_URL             = "https://aoe-api.reliclink.com"
)

func New(logger *logger.Logger, cache *cache.Cache[uint, *domain.Player]) *api {
	return &api{
		logger: logger,
		cache:  cache,
	}
}

type lobbyResponse struct {
	Matches []*Lobby `json:"matches"`
}

func (a *api) GetLobbies() ([]*domain.Lobby, error) {
	url := fmt.Sprintf("%s/community/advertisement/findAdvertisements?title=age2", BASE_URL)

	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		errorMsg := fmt.Sprintf("(%s) %d", string(body), resp.StatusCode)

		return nil, errors.NewServerError(errorMsg)
	}

	response := &lobbyResponse{}
	json.Unmarshal(body, response)

	lobbies := list.Map(response.Matches, func(lobby *Lobby) *domain.Lobby {
		memberIds := list.Map(lobby.MatchMembers, func(member *Member) uint {
			return member.ProfileId
		})

		return &domain.Lobby{
			Id:      lobby.Id,
			Members: memberIds,
		}
	})

	return lobbies, nil
}

type playerResponse struct {
	StatGroups       []*StatGroup        `json:"statGroups"`
	LeaderboardStats []*LeaderboardStats `json:"leaderboardStats"`
}

func (a *api) GetPlayers(playerIds []uint) ([]*domain.Player, error) {
	playersAlreadyFetched, playerIdsToFetch := a.findKnownAndUnknownPlayers(playerIds)

	// all players found in cache
	if len(playerIdsToFetch) == 0 {
		return playersAlreadyFetched, nil
	}

	idString := buildIdString(playerIdsToFetch)

	url := fmt.Sprintf("%s/community/leaderboard/GetPersonalStat?title=age2&profile_ids=[%s]", BASE_URL, idString)

	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		errorMsg := fmt.Sprintf("(%s) %d", string(body), resp.StatusCode)

		return nil, errors.NewServerError(errorMsg)
	}

	response := &playerResponse{}
	json.Unmarshal(body, response)

	players := list.Map(response.StatGroups, func(statGroup *StatGroup) *domain.Player {
		statGroupMember := statGroup.Members[0]

		player := parsePlayer(statGroup.Id, response.LeaderboardStats, statGroupMember)

		a.cache.Insert(player.PlayerId, player)

		return player
	})

	players = append(players, playersAlreadyFetched...)

	return players, nil
}

func (a *api) findKnownAndUnknownPlayers(playerIds []uint) ([]*domain.Player, []uint) {
	playersAlreadyFetched := []*domain.Player{}
	playerIdsToFetch := []uint{}

	for _, id := range playerIds {
		player, ok := a.cache.Contains(id)

		if ok {
			playersAlreadyFetched = append(playersAlreadyFetched, *player)
			continue
		}

		playerIdsToFetch = append(playerIdsToFetch, id)
	}

	return playersAlreadyFetched, playerIdsToFetch
}

func buildIdString(playerIds []uint) string {
	str := list.Map(playerIds, func(id uint) string {
		return fmt.Sprintf("%d", id)
	})

	return strings.Join(str, ",")
}

func parsePlayer(id uint, leaderboardStats []*LeaderboardStats, statGroupMember *StatGroupMember) *domain.Player {
	var rating_1v1 *uint = nil
	var rating_tg *uint = nil

	stats := list.Where(leaderboardStats, func(stat *LeaderboardStats) bool {
		return stat.StatGroupId == id
	})

	stats1v1, exists := list.FirstWhere(stats, func(leaderboardStats *LeaderboardStats) bool {
		return leaderboardStats.LeaderboardId == LEADERBOARD_1V1
	})

	if exists {
		rating_1v1 = &(*stats1v1).Rating
	}

	statstg, exists := list.FirstWhere(stats, func(leaderboardStats *LeaderboardStats) bool {
		return leaderboardStats.LeaderboardId == LEADERBOARD_TEAMGAME
	})

	if exists {
		rating_tg = &(*statstg).Rating
	}

	player := &domain.Player{
		PlayerId:   statGroupMember.ProfileId,
		SteamId:    statGroupMember.Name,
		Name:       statGroupMember.Alias,
		Rating_1v1: rating_1v1,
		Rating_TG:  rating_tg,
	}

	return player
}
