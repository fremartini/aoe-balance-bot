package main

import (
	"aoe-bot/internal/bot"
	"aoe-bot/internal/commands/balance"
	"aoe-bot/internal/discord"
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
	logger *logger.Logger) map[string]bot.Command {
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

				librematchApi := librematch.New(logger)

				handler := balance.New(librematchApi, discordAPI, logger)

				teams, err := handler.Handle(context, lobbyId)

				if err != nil {
					return handleError(err, context.ChannelId, discordAPI)
				}

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

				diff := abs(int(teams[0].ELO) - int(teams[1].ELO))
				diffStr := fmt.Sprintf("ELO difference: %d\n", diff)
				sb.WriteString(diffStr)

				discordAPI.ChannelMessageSend(context.ChannelId, sb.String())

				return nil
			},
			Hint: "Create two teams of players in a lobby",
		},
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
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
