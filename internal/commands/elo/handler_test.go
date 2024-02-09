package elo_test

import (
	"aoe-bot/internal/bot"
	"aoe-bot/internal/commands/elo"
	internalErrors "aoe-bot/internal/errors"
	"aoe-bot/internal/logger"
	"errors"
	"testing"
)

func TestHandle_UnknownPlayer_ReturnsError(t *testing.T) {
	// arrange
	mock := &mock{}

	logger := logger.New(0)

	mapping := map[string]string{}

	handler := elo.New(mock, mapping, logger)

	context := &bot.Context{
		UserId:    "authorId",
		ChannelId: "channelId",
	}

	// act
	_, err := handler.Handle(context)

	// assert
	if !errors.Is(err, &internalErrors.NotFoundError{}) {
		t.Errorf("actual was not expected error")
	}
}

func TestHandle_KnownPlayer_ReturnsRating(t *testing.T) {
	// arrange
	expected := 1000

	mock := &mock{
		FakeGetPlayer: func(s string) (int, error) {
			return expected, nil
		},
	}

	logger := logger.New(0)

	mapping := map[string]string{
		"authorId": "steamId",
	}

	handler := elo.New(mock, mapping, logger)

	context := &bot.Context{
		UserId:    "authorId",
		ChannelId: "channelId",
	}

	// act
	rating, err := handler.Handle(context)

	// assert
	if err != nil {
		t.Error("error was not nil")
	}

	if rating != expected {
		t.Errorf("expected %v got %v", rating, expected)
	}
}

type mock struct {
	FakeGetPlayer func(string) (int, error)
}

func (m *mock) GetPlayer(steamId string) (int, error) {
	return m.FakeGetPlayer(steamId)
}
