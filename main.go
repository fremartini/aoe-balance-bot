package main

import (
	"aoe-bot/internal/bot"
	"aoe-bot/internal/cache"
	"aoe-bot/internal/config"
	"aoe-bot/internal/librematch"
	"aoe-bot/internal/logger"
)

func main() {
	config, err := config.Read()

	if err != nil {
		panic(err)
	}

	logger := logger.New(config.LogLevel)

	b, err := bot.New(logger, config.Token)

	if err != nil {
		panic(err)
	}

	playerCache := cache.New[uint, *librematch.Player]()

	commands := New(b.Session, logger, playerCache)

	b.Run(commands)
}
