package main

import (
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

	b := bot.New(logger)

	if err := b.Run(config.Token); err != nil {
		panic(err)
	}
}
