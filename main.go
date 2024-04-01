package main

import (
	"aoe-bot/internal/bot"
	"aoe-bot/internal/cache"
	"aoe-bot/internal/config"
	"aoe-bot/internal/domain"
	"aoe-bot/internal/logger"
)

func main() {
	config, err := config.Read()

	if err != nil {
		panic(err)
	}

	logger := logger.New(config.LogLevel)

	logger.Infof("Log level %d", config.LogLevel)
	logger.Infof("Cache expiry %d", config.Cache.ExpiryHours)
	logger.Infof("Cache size %d", config.Cache.MaxSize)

	b, err := bot.New(logger, config.Token)

	if err != nil {
		panic(err)
	}

	playerCache := cache.New[uint, *domain.Player](config.Cache.ExpiryHours, config.Cache.MaxSize, logger)

	commands := New(b.Session, logger, playerCache)

	b.Run(commands, config.Port)
}
