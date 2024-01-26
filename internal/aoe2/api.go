package aoe2

import (
	"aoe-bot/internal/logger"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type api struct {
	logger *logger.Logger
}

func New(logger *logger.Logger) *api {
	return &api{
		logger: logger,
	}
}

const base_url = "https://aoe2.net/api"

type playerResponse struct {
	Rating int
}

func (a *api) GetPlayer(steamId string) (int, error) {
	uri := fmt.Sprintf("%s/player/ratinghistory?game=aoe2de&leaderboard_id=3&count=1&steam_id=%s", base_url, steamId)

	a.logger.Infof("Requesting ELO for Steam user %s", steamId)

	resp, err := http.Get(uri)

	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return 0, err
	}

	players := &[]playerResponse{}
	json.Unmarshal(body, players)

	if len(*players) == 0 {
		return 0, fmt.Errorf("no user with id %s exists", steamId)
	}

	return (*players)[0].Rating, nil
}
