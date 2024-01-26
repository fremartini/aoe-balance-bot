package team

import (
	"aoe-bot/internal/bot"
	"aoe-bot/internal/logger"
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
	return nil
}
