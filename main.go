package main

import (
	"aoe-bot/internal/bot"
	"aoe-bot/internal/cache"
	"aoe-bot/internal/config"
	"aoe-bot/internal/domain"
	"aoe-bot/internal/list"
	"crypto/tls"
	"fmt"
	"log"
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

	channelStr := list.Map(config.WhitelistedChannels, func(c uint) string {
		return fmt.Sprintf("%v", c)
	})

	portStr := "nil"
	if config.Port != nil {
		portStr = fmt.Sprintf("%d", *config.Port)
	}

	log.Printf(
		"Cache expiry %d, Cache size %d, Trust insecure certificates %t, port %s, Whitelisted channels [%s]",
		config.Cache.ExpiryHours,
		config.Cache.MaxSize,
		config.TrustInsecureCertificates,
		portStr,
		strings.Join(channelStr, ","))

	b, err := bot.New(Prefix, config.Token, channelStr)

	if err != nil {
		panic(err)
	}

	if config.TrustInsecureCertificates {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	playerCache := cache.New[uint, *domain.Player](config.Cache.ExpiryHours, config.Cache.MaxSize)

	commands := New(b.Session, playerCache, Prefix)

	b.Run(commands, config.Port)
}
