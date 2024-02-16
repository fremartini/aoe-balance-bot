package librematch

import (
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
}

const base_url = "https://aoe-api.reliclink.com"

func New(logger *logger.Logger) *api {
	return &api{
		logger: logger,
	}
}

type lobbyResponse struct {
	Matches []*Lobby `json:"matches"`
}

func (a *api) GetLobbies() ([]*Lobby, error) {
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

	return response.Matches, nil
}

type playerResponse struct {
	StatGroups       []*StatGroup        `json:"statGroups"`
	LeaderboardStats []*LeaderboardStats `json:"leaderboardStats"`
}

func (a *api) GetPlayer(playerId uint) (*Player, error) {
	url := fmt.Sprintf("https://aoe-api.reliclink.com/community/leaderboard/GetPersonalStat?title=age2&profile_ids=[\"%d\"]", playerId)

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

	firstStatGroup := response.StatGroups[0]

	var rating uint = 1000

	if len(response.LeaderboardStats) > 0 {
		// player has played a ranked 1v1 match
		_1v1stats, _, exists := list.FirstWhere(response.LeaderboardStats, func(leaderboardStats *LeaderboardStats) bool {
			return leaderboardStats.LeaderboardId == 3
		})

		if exists {
			rating = (*_1v1stats).Rating
		} else {
			rating = response.LeaderboardStats[0].Rating
		}
	}

	player := &Player{
		ProfileId: firstStatGroup.Members[0].ProfileId,
		Name:      firstStatGroup.Members[0].Name,
		Alias:     firstStatGroup.Members[0].Alias,
		Rating:    rating,
	}

	return player, nil
}
