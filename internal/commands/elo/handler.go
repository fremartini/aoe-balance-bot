package elo

import (
	"aoe-bot/internal/bot"
	"aoe-bot/internal/logger"
	"errors"
	"fmt"
)

type dataProvider interface {
	GetPlayer(steamId string) (int, error)
}

type handler struct {
	provider      dataProvider
	logger        *logger.Logger
	playerMapping map[string]string
}

func New(provider dataProvider,
	playerMapping map[string]string,
	logger *logger.Logger) *handler {
	return &handler{
		provider:      provider,
		playerMapping: playerMapping,
		logger:        logger,
	}
}

func (h *handler) Handle(context *bot.Context) error {
	h.logger.Infof("Getting ELO info for Discord user %s", context.AuthorId)

	steamId, ok := h.playerMapping[context.AuthorId]

	if !ok {
		return errors.New("unknown player")
	}

	rating, err := h.provider.GetPlayer(steamId)

	if err != nil {
		return err
	}

	fmt.Println(rating)

	return nil
}
