package main

import (
	"aoe-bot/internal/bot"
	"aoe-bot/internal/cache"
	"aoe-bot/internal/commands/balance"
	"aoe-bot/internal/discord"
	"aoe-bot/internal/domain"
	internalErrors "aoe-bot/internal/errors"
	"aoe-bot/internal/librematch"
	"aoe-bot/internal/logger"
	"errors"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const prefix = "!"

func New(
	session *discordgo.Session,
	logger *logger.Logger,
	playerCache *cache.Cache[uint, *domain.Player]) map[string]bot.Command {
	return map[string]bot.Command{
		withPrefix("balance"): {
			Handle: func(context *bot.Context, args []string) error {
				discordAPI := discord.New(session)

				if len(args) == 0 {
					discordAPI.ChannelMessageSend(context.ChannelId, "Missing game id")
					return nil
				}

				fullLobbyId := strings.Split(args[0], "/")
				lobbyId := fullLobbyId[len(fullLobbyId)-1]

				librematchApi := librematch.New(logger, playerCache)

				handler := balance.New(librematchApi, discordAPI, discordAPI, logger)

				err := handler.Handle(context, lobbyId)

				if err != nil {
					return handleError(err, context.ChannelId, discordAPI)
				}

				return nil
			},
			Hint: "Create two teams of players in a lobby",
		},
	}
}

func withPrefix(cmd string) string {
	return fmt.Sprintf("%s%s", prefix, cmd)
}

type messageSender interface {
	ChannelMessageSend(channelID, content string)
}

func handleError(err error, channelId string, api messageSender) error {
	var serverErr *internalErrors.ServerError
	if errors.As(err, &serverErr) {
		api.ChannelMessageSend(channelId, "Server error")

		return nil
	}

	var notFoundErr *internalErrors.NotFoundError
	if errors.As(err, &notFoundErr) {
		api.ChannelMessageSend(channelId, "Unknown player")

		return nil
	}

	var applicationErr *internalErrors.ApplicationError
	if errors.As(err, &applicationErr) {
		api.ChannelMessageSend(channelId, applicationErr.Message)

		return nil
	}

	return err
}
