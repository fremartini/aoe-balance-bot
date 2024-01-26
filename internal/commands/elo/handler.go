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

type messageProvider interface {
	ChannelMessageSend(channelID, content string)
}

type handler struct {
	dataProvider    dataProvider
	messageProvider messageProvider
	logger          *logger.Logger
	playerMapping   map[string]string
}

func New(
	dataProvider dataProvider,
	messageProvider messageProvider,
	playerMapping map[string]string,
	logger *logger.Logger) *handler {
	return &handler{
		dataProvider:    dataProvider,
		messageProvider: messageProvider,
		playerMapping:   playerMapping,
		logger:          logger,
	}
}

func (h *handler) Handle(context *bot.Context) error {
	h.logger.Infof("Getting ELO info for Discord user %s", context.AuthorId)

	steamId, ok := h.playerMapping[context.AuthorId]

	if !ok {
		h.messageProvider.ChannelMessageSend(context.ChannelId, "Unknown player")
		return errors.New("unknown player")
	}

	rating, err := h.dataProvider.GetPlayer(steamId)

	if err != nil {
		return err
	}

	h.messageProvider.ChannelMessageSend(context.ChannelId, fmt.Sprintf("Your rating is %v", rating))

	return nil
}
