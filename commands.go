package main

import (
	"aoe-bot/internal/aoe2"
	"aoe-bot/internal/bot"
	"aoe-bot/internal/commands/elo"
	"aoe-bot/internal/commands/team"
	"aoe-bot/internal/discord"
	internalErrors "aoe-bot/internal/errors"
	"aoe-bot/internal/logger"
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

const prefix = "!"

func New(
	playerMapping map[string]string,
	session *discordgo.Session,
	logger *logger.Logger) map[string]bot.Command {
	return map[string]bot.Command{
		withPrefix("team"): {
			Handle: func(context *bot.Context, args []string) error {
				aoe2NetApi := aoe2.New(logger)

				discordAPI := discord.New(session)

				handler := team.New(aoe2NetApi, discordAPI, playerMapping, logger)

				teams, err := handler.Handle(context)

				if err == nil {
					fmt.Println(teams)

					return nil
				}

				return handleError(err, context.ChannelId, discordAPI)
			},
			Hint: "Create two teams consisting of players in your current channel",
		},
		withPrefix("elo"): {
			Handle: func(context *bot.Context, args []string) error {
				aoe2NetApi := aoe2.New(logger)

				discordAPI := discord.New(session)

				handler := elo.New(aoe2NetApi, playerMapping, logger)

				rating, err := handler.Handle(context)

				if err == nil {
					discordAPI.ChannelMessageSend(context.ChannelId, fmt.Sprintf("Your rating is %v", rating))

					return nil
				}

				return handleError(err, context.ChannelId, discordAPI)

			},
			Hint: "Get your current 1v1 ELO",
		},
	}
}

func withPrefix(cmd string) string {
	return fmt.Sprintf("%s%s", prefix, cmd)
}

type MessageSender interface {
	ChannelMessageSend(channelID, content string)
}

func handleError(err error, channelId string, api MessageSender) error {
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
