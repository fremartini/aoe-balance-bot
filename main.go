package main

import (
	"aoe-bot/internal/api"
	"aoe-bot/internal/bot"
	"aoe-bot/internal/config"
	"aoe-bot/internal/logger"
)

func main() {
	config, err := config.Read()

	if err != nil {
		panic(err)
	}

	logger := logger.New(config.LogLevel)

	api := api.New(logger)

	playerMapping := map[string]string{
		//TODO: fetch mapping somewhere
	}

	commands := New(api, playerMapping, logger)

	b := bot.New(logger, commands)

	if err := b.Run(config.Token); err != nil {
		panic(err)
	}
}
