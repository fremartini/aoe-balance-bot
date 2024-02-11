package main

import (
	"aoe-bot/internal/bot"
	"aoe-bot/internal/config"
	"aoe-bot/internal/logger"
	playermapper "aoe-bot/internal/player_mapper"
)

func main() {
	config, err := config.Read()

	if err != nil {
		panic(err)
	}

	logger := logger.New(config.LogLevel)

	b, err := bot.New(logger, config.Token, config.SteamIdChannel)

	if err != nil {
		panic(err)
	}

	playerMapper := playermapper.New(logger, b.Session, config.SteamIdChannel)

	err = playerMapper.BuildPlayerMapping()

	if err != nil {
		panic(err)
	}

	commands := New(playerMapper, b.Session, logger)

	b.Run(commands, playerMapper)
}
