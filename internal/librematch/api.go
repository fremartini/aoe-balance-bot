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
)

type api struct {
	logger *logger.Logger
	cache  *cache.Cache[uint, *domain.Player]
}

const (
	leaderboard_1v1      = 3
	leaderboard_teamgame = 4
	base_url             = "https://aoe-api.reliclink.com"
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
	url := fmt.Sprintf("%s/community/advertisement/findAdvertisements?title=age2", base_url)

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

func (a *api) GetPlayer(playerId uint) (*domain.Player, error) {

	// cache lookup
	p, exists := a.cache.Contains(playerId)

	if exists {
		a.logger.Infof("Found %d in cache (%s)", playerId, (*p).Name)
		return *p, nil
	}

	url := fmt.Sprintf("%s/community/leaderboard/GetPersonalStat?title=age2&profile_ids=[\"%d\"]", base_url, playerId)

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

	var rating_1v1 *uint = nil
	var rating_tg *uint = nil

	if len(response.LeaderboardStats) > 0 {
		stats1v1, exists := list.FirstWhere(response.LeaderboardStats, func(leaderboardStats *LeaderboardStats) bool {
			return leaderboardStats.LeaderboardId == leaderboard_1v1
		})

		if exists {
			rating_1v1 = &(*stats1v1).Rating
		}

		statstg, exists := list.FirstWhere(response.LeaderboardStats, func(leaderboardStats *LeaderboardStats) bool {
			return leaderboardStats.LeaderboardId == leaderboard_teamgame
		})

		if exists {
			rating_tg = &(*statstg).Rating
		}
	}

	firstStatGroup := response.StatGroups[0]

	player := &domain.Player{
		PlayerId:   firstStatGroup.Members[0].ProfileId,
		SteamId:    firstStatGroup.Members[0].Name,
		Name:       firstStatGroup.Members[0].Alias,
		Rating_1v1: rating_1v1,
		Rating_TG:  rating_tg,
	}

	a.logger.Infof("Inserted %d into cache (%s)", playerId, player.Name)
	a.cache.Insert(playerId, player)

	return player, nil
}
