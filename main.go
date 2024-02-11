package main

import (
	"aoe-bot/internal/bot"
	"aoe-bot/internal/config"
	"aoe-bot/internal/logger"
	"encoding/json"
	"strconv"

	"github.com/bwmarrin/discordgo"
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

	playerMapping, err := getPlayerMapping(logger, b.Session, config.SteamIdChannel)

	if err != nil {
		panic(err)
	}

	commands := New(playerMapping, b.Session, logger)

	b.Run(commands)
}

func getPlayerMapping(logger *logger.Logger, session *discordgo.Session, channelId string) (map[string]string, error) {
	messages, err := session.ChannelMessages(channelId, 100, "", "", "")

	if err != nil {
		return nil, err
	}

	logger.Info("Building player map")

	playerMapping := map[string]string{}

	for _, msg := range messages {
		if _, err := strconv.Atoi(msg.Content); err == nil {
			playerMapping[msg.Author.ID] = msg.Content
		}
	}

	v, err := json.MarshalIndent(playerMapping, "", " ")

	if err != nil {
		return nil, err
	}

	logger.Infof("Built player map: %s", v)

	return playerMapping, nil
}
