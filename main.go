package main

import (
	"aoe-bot/internal/bot"
	"aoe-bot/internal/cache"
	"aoe-bot/internal/config"
	"aoe-bot/internal/domain"
	"aoe-bot/internal/logger"
	"crypto/tls"
	"net/http"
)

const (
	Prefix = "!"
)

func main() {
	config, err := config.Read()

	if err != nil {
		panic(err)
	}

	logger := logger.New(config.LogLevel)

	logger.Infof("Log level %d, Cache expiry %d, Cache size %d", config.LogLevel, config.Cache.ExpiryHours, config.Cache.MaxSize)

	b, err := bot.New(logger, Prefix, config.Token)

	if err != nil {
		panic(err)
	}

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	playerCache := cache.New[uint, *domain.Player](config.Cache.ExpiryHours, config.Cache.MaxSize, logger)

	commands := New(b.Session, logger, playerCache, Prefix)

	b.Run(commands, config.Port)
}
