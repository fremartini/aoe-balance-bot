package main

import (
	"aoe-bot/internal/api"
	"aoe-bot/internal/bot"
	"aoe-bot/internal/commands/elo"
	"aoe-bot/internal/discord"
	"aoe-bot/internal/logger"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

const prefix = "!"

func New(
	api *api.Api,
	playerMapping map[string]string,
	session *discordgo.Session,
	logger *logger.Logger) map[string]bot.Command {
	return map[string]bot.Command{
		withPrefix("team"): {
			Handle: func(context *bot.Context, args []string) error {
				logger.Info("Not implemented")

				return nil
			},
			Hint: "Create two teams consisting of players in your current channel",
		},
		withPrefix("elo"): {
			Handle: func(context *bot.Context, args []string) error {
				messageProvider := discord.New(session)

				handler := elo.New(api, messageProvider, playerMapping, logger)

				return handler.Handle(context)
			},
			Hint: "Get your current 1v1 ELO",
		},
	}
}

func withPrefix(cmd string) string {
	return fmt.Sprintf("%s%s", prefix, cmd)
}
