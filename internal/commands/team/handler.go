package team

import (
	"aoe-bot/internal/bot"
	"aoe-bot/internal/logger"
	"fmt"
)

type dataProvider interface {
	GetPlayer(steamId string) (int, error)
}

type discordInfoProvider interface {
	ChannelMessageSend(channelID, content string)
	FindUserVoiceChannel(guildId, userid string) (string, error)
	FindUsersInVoiceChannel(serverId, channelId string) ([]*string, error)
}

type handler struct {
	dataProvider        dataProvider
	discordInfoProvider discordInfoProvider
	logger              *logger.Logger
	playerMapping       map[string]string
}

func New(
	dataProvider dataProvider,
	discordInfoProvider discordInfoProvider,
	playerMapping map[string]string,
	logger *logger.Logger) *handler {
	return &handler{
		dataProvider:        dataProvider,
		discordInfoProvider: discordInfoProvider,
		playerMapping:       playerMapping,
		logger:              logger,
	}
}

func (h *handler) Handle(context *bot.Context) error {
	// find users channel
	channelId, err := h.discordInfoProvider.FindUserVoiceChannel(context.ServerId, context.UserId)

	if err != nil {
		return err
	}

	// find all users in channel
	users, err := h.discordInfoProvider.FindUsersInVoiceChannel(context.ServerId, channelId)

	if err != nil {
		return err
	}

	// map discord ids to steam ids
	steamIds := []string{}
	unknowns := []string{}

	for _, user := range users {
		steamId, ok := h.playerMapping[*user]

		if !ok {
			unknowns = append(unknowns, *user)
		}

		steamIds = append(steamIds, steamId)
	}

	fmt.Println(steamIds)
	fmt.Println(unknowns)

	// lookup steam ids for all users
	ratings := map[string]int{}

	for _, steamId := range steamIds {
		rating, err := h.dataProvider.GetPlayer(steamId)

		if err != nil {
			h.logger.Warnf("could not get ranking for %s: %s", steamId, err)
			continue
		}

		ratings[steamId] = rating
	}

	// create teams
	for k, v := range ratings {
		fmt.Println(k, v)
	}

	return nil
}
