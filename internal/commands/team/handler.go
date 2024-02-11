package team

import (
	"aoe-bot/internal/bot"
	"aoe-bot/internal/discord"
	"aoe-bot/internal/logger"
	"cmp"
	"slices"
)

type dataProvider interface {
	GetPlayer(steamId string) (int, error)
}

type discordInfoProvider interface {
	FindUserVoiceChannel(guildId, userid string) (string, error)
	FindUsersInVoiceChannel(serverId, channelId string) ([]*discord.User, error)
}

type playerProvider interface {
	GetPlayers() map[string]string
}

type handler struct {
	dataProvider        dataProvider
	discordInfoProvider discordInfoProvider
	playerProvider      playerProvider
	logger              *logger.Logger
}

func New(
	dataProvider dataProvider,
	discordInfoProvider discordInfoProvider,
	playerProvider playerProvider,
	logger *logger.Logger) *handler {
	return &handler{
		dataProvider:        dataProvider,
		discordInfoProvider: discordInfoProvider,
		playerProvider:      playerProvider,
		logger:              logger,
	}
}

func (h *handler) Handle(context *bot.Context) ([]*Team, []string, error) {
	channelId, err := h.discordInfoProvider.FindUserVoiceChannel(context.ServerId, context.UserId)

	if err != nil {
		return nil, nil, err
	}

	users, err := h.discordInfoProvider.FindUsersInVoiceChannel(context.ServerId, channelId)

	if err != nil {
		return nil, nil, err
	}

	players := []*Player{}
	unknowns := []string{}

	for _, user := range users {
		steamId, ok := h.playerProvider.GetPlayers()[user.Id]

		if !ok {
			unknowns = append(unknowns, user.Username)
			continue
		}

		rating, err := h.dataProvider.GetPlayer(steamId)

		if err != nil {
			h.logger.Warnf("could not get ranking for %s: %s", steamId, err)
			continue
		}

		player := &Player{
			DiscordName: user.Username,
			SteamId:     steamId,
			Rating:      rating,
		}

		players = append(players, player)
	}

	team1, team2 := createTeams(players)

	return []*Team{team1, team2}, unknowns, nil
}

func createTeams(players []*Player) (*Team, *Team) {
	t1 := &Team{}
	t2 := &Team{}

	t1Rating := 0
	t2Rating := 0

	slices.SortFunc(players, func(a, b *Player) int {
		return cmp.Compare(b.Rating, a.Rating)
	})

	for _, player := range players {
		if t1Rating < t2Rating {
			t1.Players = append(t1.Players, player)
			t1Rating += player.Rating
		} else {
			t2.Players = append(t2.Players, player)
			t2Rating += player.Rating
		}
	}

	return t1, t2
}
