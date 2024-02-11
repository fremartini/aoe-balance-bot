package team_test

import (
	"aoe-bot/internal/bot"
	"aoe-bot/internal/commands/team"
	"aoe-bot/internal/discord"
	"aoe-bot/internal/logger"
	"errors"
	"slices"
	"testing"
)

func TestHandle_ReturnsTeams(t *testing.T) {
	// arrange
	user1 := discord.User{
		Username: "user1",
		Id:       "id1",
	}

	user2 := discord.User{
		Username: "user2",
		Id:       "id2",
	}

	user3 := discord.User{
		Username: "user3",
		Id:       "id3",
	}

	user4 := discord.User{
		Username: "user4",
		Id:       "id4",
	}

	steamId1 := "steamId1"
	steamId2 := "steamId2"
	steamId3 := "steamId3"
	steamId4 := "steamId4"

	mock := &mock{
		FakeFindUserVoiceChannel: func(guildId, userid string) (string, error) {
			return "voiceChannel", nil
		},
		FakeFindUsersInVoiceChannel: func(serverId, channelId string) ([]*discord.User, error) {
			return []*discord.User{
				&user1, &user2, &user3, &user4,
			}, nil
		},
		FakeGetPlayer: func(s string) (int, error) {
			if s == steamId1 {
				return 1800, nil
			}

			if s == steamId2 {
				return 1300, nil
			}

			if s == steamId3 {
				return 800, nil
			}

			if s == steamId4 {
				return 2200, nil
			}

			return 0, errors.New("invalid player")
		},
	}

	players := map[string]string{
		user1.Id: steamId1,
		user2.Id: steamId2,
		user3.Id: steamId3,
		user4.Id: steamId4,
	}

	logger := logger.New(0)

	handler := team.New(mock, mock, players, logger)

	context := &bot.Context{}

	// act
	teams, unknowns, err := handler.Handle(context)

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

	if len(unknowns) != 0 {
		t.Error("unknowns should be empty")
	}
}

type mock struct {
	FakeGetPlayer               func(string) (int, error)
	FakeFindUserVoiceChannel    func(guildId, userid string) (string, error)
	FakeFindUsersInVoiceChannel func(serverId, channelId string) ([]*discord.User, error)
}

func (m *mock) GetPlayer(steamId string) (int, error) {
	return m.FakeGetPlayer(steamId)
}

func (m *mock) FindUserVoiceChannel(guildId, userid string) (string, error) {
	return m.FakeFindUserVoiceChannel(guildId, userid)
}

func (m *mock) FindUsersInVoiceChannel(serverId, channelId string) ([]*discord.User, error) {
	return m.FakeFindUsersInVoiceChannel(serverId, channelId)
}
