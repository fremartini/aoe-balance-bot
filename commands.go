package main

import (
	"aoe-bot/internal/api"
	"aoe-bot/internal/bot"
	"aoe-bot/internal/elo"
	"aoe-bot/internal/logger"
	"fmt"
)

const prefix = "!"

func New(
	api *api.Api, 
	playerMapping map[string]string,
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
				handler := elo.New(api, playerMapping, logger)

				return handler.Handle(context)
			},
			Hint: "Get your current 1v1 ELO",
		},
	}
}

func withPrefix(cmd string) string {
	return fmt.Sprintf("%s%s", prefix, cmd)
}
