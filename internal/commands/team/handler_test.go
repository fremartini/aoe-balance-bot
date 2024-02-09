package team_test

import (
	"aoe-bot/internal/bot"
	"aoe-bot/internal/commands/team"
	"aoe-bot/internal/logger"
	"errors"
	"slices"
	"testing"
)

func TestHandle_ReturnsTeams(t *testing.T) {
	// arrange
	user1 := "user1"
	user2 := "user2"
	user3 := "user3"
	user4 := "user4"

	id1 := "id1"
	id2 := "id2"
	id3 := "id3"
	id4 := "id4"

	mock := &mock{
		FakeFindUserVoiceChannel: func(guildId, userid string) (string, error) {
			return "voiceChannel", nil
		},
		FakeFindUsersInVoiceChannel: func(serverId, channelId string) ([]*string, error) {
			return []*string{
				&user1, &user2, &user3, &user4,
			}, nil
		},
		FakeGetPlayer: func(s string) (int, error) {
			if s == id1 {
				return 1800, nil
			}

			if s == id2 {
				return 1300, nil
			}

			if s == id3 {
				return 800, nil
			}

			if s == id4 {
				return 2200, nil
			}

			return 0, errors.New("invalid player")
		},
	}

	players := map[string]string{
		user1: id1,
		user2: id2,
		user3: id3,
		user4: id4,
	}

	logger := logger.New(0)

	handler := team.New(mock, mock, players, logger)

	context := &bot.Context{}

	// act
	teams, err := handler.Handle(context)

	// assert
	if err != nil {
		t.Error(err)
	}

	if len(teams) != 2 {
		t.Error("teams should have length 2")
	}

	t1 := teams[0]
	t2 := teams[1]

	if !slices.ContainsFunc(t1.Players, func(p *team.Player) bool {
		return p.Rating == 1800
	}) {
		t.Error("missing ranking 1800")
	}

	if !slices.ContainsFunc(t1.Players, func(p *team.Player) bool {
		return p.Rating == 1300
	}) {
		t.Error("missing ranking 1300")
	}

	if !slices.ContainsFunc(t2.Players, func(p *team.Player) bool {
		return p.Rating == 2200
	}) {
		t.Error("missing ranking 2200")
	}

	if !slices.ContainsFunc(t2.Players, func(p *team.Player) bool {
		return p.Rating == 800
	}) {
		t.Error("missing ranking 800")
	}
}

type mock struct {
	FakeGetPlayer               func(string) (int, error)
	FakeFindUserVoiceChannel    func(guildId, userid string) (string, error)
	FakeFindUsersInVoiceChannel func(serverId, channelId string) ([]*string, error)
}

func (m *mock) GetPlayer(steamId string) (int, error) {
	return m.FakeGetPlayer(steamId)
}

func (m *mock) FindUserVoiceChannel(guildId, userid string) (string, error) {
	return m.FakeFindUserVoiceChannel(guildId, userid)
}

func (m *mock) FindUsersInVoiceChannel(serverId, channelId string) ([]*string, error) {
	return m.FakeFindUsersInVoiceChannel(serverId, channelId)
}
