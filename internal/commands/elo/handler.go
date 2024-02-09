package elo

import (
	"aoe-bot/internal/bot"
	"aoe-bot/internal/errors"
	"aoe-bot/internal/logger"
)

type dataProvider interface {
	GetPlayer(steamId string) (int, error)
}

type handler struct {
	dataProvider  dataProvider
	playerMapping map[string]string
	logger        *logger.Logger
}

func New(
	dataProvider dataProvider,
	playerMapping map[string]string,
	logger *logger.Logger) *handler {
	return &handler{
		dataProvider:  dataProvider,
		playerMapping: playerMapping,
		logger:        logger,
	}
}

func (h *handler) Handle(context *bot.Context) (int, error) {
	h.logger.Infof("Getting ELO info for Discord user %s", context.UserId)

	steamId, ok := h.playerMapping[context.UserId]

	if !ok {
		return 0, errors.NewNotFoundError()
	}

	return h.dataProvider.GetPlayer(steamId)
}
