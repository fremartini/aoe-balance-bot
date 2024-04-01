package main

import (
	"aoe-bot/internal/bot"
	"aoe-bot/internal/cache"
	"aoe-bot/internal/config"
	"aoe-bot/internal/domain"
	"aoe-bot/internal/logger"
	"encoding/json"
	"fmt"
)

func main() {
	config, err := config.Read()

	if err != nil {
		panic(err)
	}

	logger := logger.New(config.LogLevel)

	logger.Infof("Log level %d, Cache expiry %d, Cache size %d", config.LogLevel, config.Cache.ExpiryHours, config.Cache.MaxSize)

	b, err := bot.New(logger, config.Token)

	if err != nil {
		panic(err)
	}

	playerCache := cache.New[uint, *domain.Player](config.Cache.ExpiryHours, config.Cache.MaxSize, logger)

	commands := New(b.Session, logger, playerCache)

	fmt.Println(prettyPrint(config))

	b.Run(commands)
}

func prettyPrint(a any) string {
	b, _ := json.Marshal(a)

	return string(b)
}
