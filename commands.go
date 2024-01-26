package main

import (
	"aoe-bot/internal/aoe2"
	"aoe-bot/internal/bot"
	"aoe-bot/internal/commands/elo"
	"aoe-bot/internal/discord"
	"aoe-bot/internal/logger"
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
				logger.Info("Not implemented")

				return nil
			},
			Hint: "Create two teams consisting of players in your current channel",
		},
		withPrefix("elo"): {
			Handle: func(context *bot.Context, args []string) error {
				aoe2NetApi := aoe2.New(logger)

				discordAPI := discord.New(session)

				handler := elo.New(aoe2NetApi, discordAPI, playerMapping, logger)

				return handler.Handle(context)
			},
			Hint: "Get your current 1v1 ELO",
		},
	}
}

func withPrefix(cmd string) string {
	return fmt.Sprintf("%s%s", prefix, cmd)
}
