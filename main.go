package main

import (
	"aoe-bot/internal/bot"
	"aoe-bot/internal/config"
	"aoe-bot/internal/logger"
	"fmt"
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

	commands := New(b.Session, logger)

	fmt.Println(commands)

	//b.Run(commands)
}
