package main

import (
	"aoe-bot/internal/bot"
	"aoe-bot/internal/cache"
	"aoe-bot/internal/config"
	"aoe-bot/internal/domain"
	"aoe-bot/internal/list"
	"aoe-bot/internal/logger"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
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

	channelStr := list.Map(config.WhitelistedChannels, func(c uint) string {
		return fmt.Sprintf("%v", c)
	})

	logger.Infof(
		"Log level %d, Cache expiry %d, Cache size %d, Trust insecure certificates %t, Whitelisted channels [%s]",
		config.LogLevel,
		config.Cache.ExpiryHours,
		config.Cache.MaxSize,
		config.TrustInsecureCertificates,
		strings.Join(channelStr, ","))

	b, err := bot.New(logger, Prefix, config.Token, channelStr)

	if err != nil {
		panic(err)
	}

	if config.TrustInsecureCertificates {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	playerCache := cache.New[uint, *domain.Player](config.Cache.ExpiryHours, config.Cache.MaxSize, logger)

	commands := New(b.Session, logger, playerCache, Prefix)

	b.Run(commands, config.Port)
}
