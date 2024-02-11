package playermapper

import (
	"aoe-bot/internal/logger"
	"encoding/json"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

type messageProvider interface {
	ChannelMessages(channelID string, limit int, beforeID string, afterID string, aroundID string, options ...discordgo.RequestOption) (st []*discordgo.Message, err error)
}

func New(logger *logger.Logger, provider messageProvider, channelId string) *PlayerMapper {
	return &PlayerMapper{
		logger:    logger,
		provider:  provider,
		mapping:   map[string]string{},
		channelId: channelId,
	}
}

type PlayerMapper struct {
	logger    *logger.Logger
	provider  messageProvider
	mapping   map[string]string
	channelId string
}

func (m *PlayerMapper) BuildPlayerMapping() error {
	messages, err := m.provider.ChannelMessages(m.channelId, 100, "", "", "")

	if err != nil {
		return err
	}

	m.logger.Info("Building player map")

	playerMapping := map[string]string{}

	for _, msg := range messages {
		if _, err := strconv.Atoi(msg.Content); err == nil {
			playerMapping[msg.Author.ID] = msg.Content
		}
	}

	v, err := json.MarshalIndent(playerMapping, "", " ")

	if err != nil {
		return err
	}

	m.logger.Infof("Built player map: %s", v)

	m.mapping = playerMapping

	return nil
}

func (m *PlayerMapper) GetPlayers() map[string]string {
	return m.mapping
}

func (m *PlayerMapper) AddPlayer(discordId, steamId string) {
	m.mapping[discordId] = steamId
	m.logger.Infof("Added entry to player mapping: %s --> %s", discordId, steamId)
}
