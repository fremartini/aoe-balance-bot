package elo_test

import (
	"aoe-bot/internal/bot"
	"aoe-bot/internal/commands/elo"
	"aoe-bot/internal/logger"
	"fmt"
	"strings"
	"testing"
)

func TestHandle_UnknownPlayer_ReturnsError(t *testing.T) {
	// arrange
	mock := &mock{
		FakeChannelMessageSend: func(s1, s2 string) {},
	}

	logger := logger.New(0)

	mapping := map[string]string{}

	handler := elo.New(mock, mock, mapping, logger)

	context := &bot.Context{
		AuthorId:  "authorId",
		ChannelId: "channelId",
	}

	// act
	actual := handler.Handle(context)

	// assert
	if actual == nil {
		t.Errorf("actual was null")
	}
}

func TestHandle_KnownPlayer_SendsMessage(t *testing.T) {
	// arrange
	called := false
	rating := 1000
	authorId := "authorId"
	steamId := "steamId"

	mock := &mock{
		FakeChannelMessageSend: func(channelId, content string) {
			if strings.Contains(content, fmt.Sprint(rating)) {
				called = true
			}
		},
		FakeGetPlayer: func(id string) (int, error) {
			if id == steamId {
				return rating, nil
			}
			return 0, nil

		},
	}

	logger := logger.New(0)

	mapping := map[string]string{
		authorId: steamId,
	}

	handler := elo.New(mock, mock, mapping, logger)

	context := &bot.Context{
		AuthorId:  authorId,
		ChannelId: "channelId",
	}

	// act
	handler.Handle(context)

	// assert
	if !called {
		t.Errorf("message was not sent")
	}
}

type mock struct {
	FakeGetPlayer          func(string) (int, error)
	FakeChannelMessageSend func(string, string)
}

func (m *mock) GetPlayer(steamId string) (int, error) {
	return m.FakeGetPlayer(steamId)
}

func (m *mock) ChannelMessageSend(channelID, content string) {
	m.FakeChannelMessageSend(channelID, content)
}
