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
	"strings"

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

				teams, unknowns, err := handler.Handle(context)

				if err == nil {
					var sb strings.Builder

					for teamNumber, team := range teams {
						players := team.Players

						sb.WriteString(fmt.Sprintf("Team %d:\n", teamNumber+1))
						for _, player := range players {
							s := fmt.Sprintf("%s (%d)\n", player.DiscordName, player.Rating)
							sb.WriteString(s)
						}
						sb.WriteString("\n")
					}

					if len(unknowns) > 0 {
						sb.WriteString("Unknown players (missing steam id)\n")

						for _, username := range unknowns {
							s := fmt.Sprintf("%s\n", username)
							sb.WriteString(s)
						}

						sb.WriteString("\n")
					}

					discordAPI.ChannelMessageSend(context.ChannelId, sb.String())

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
