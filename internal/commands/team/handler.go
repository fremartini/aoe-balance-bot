package team

import (
	"aoe-bot/internal/bot"
	"aoe-bot/internal/logger"
	"cmp"
	"fmt"
	"slices"
)

type dataProvider interface {
	GetPlayer(steamId string) (int, error)
}

type discordInfoProvider interface {
	FindUserVoiceChannel(guildId, userid string) (string, error)
	FindUsersInVoiceChannel(serverId, channelId string) ([]*string, error)
}

type handler struct {
	dataProvider        dataProvider
	discordInfoProvider discordInfoProvider
	logger              *logger.Logger
	playerMapping       map[string]string
}

func New(
	dataProvider dataProvider,
	discordInfoProvider discordInfoProvider,
	playerMapping map[string]string,
	logger *logger.Logger) *handler {
	return &handler{
		dataProvider:        dataProvider,
		discordInfoProvider: discordInfoProvider,
		playerMapping:       playerMapping,
		logger:              logger,
	}
}

func (h *handler) Handle(context *bot.Context) ([]*Team, error) {
	// find users channel
	channelId, err := h.discordInfoProvider.FindUserVoiceChannel(context.ServerId, context.UserId)

	if err != nil {
		return nil, err
	}

	// find all usernames in channel
	usernames, err := h.discordInfoProvider.FindUsersInVoiceChannel(context.ServerId, channelId)

	if err != nil {
		return nil, err
	}

	// map discord ids to steam ids
	players := []*Player{}
	unknowns := []string{}

	for _, username := range usernames {
		steamId, ok := h.playerMapping[*username]

		if !ok {
			unknowns = append(unknowns, *username)
			continue
		}

		rating, err := h.dataProvider.GetPlayer(steamId)

		if err != nil {
			h.logger.Warnf("could not get ranking for %s: %s", steamId, err)
			continue
		}

		player := &Player{
			DiscordName: *username,
			SteamId:     steamId,
			Rating:      rating,
		}

		players = append(players, player)
	}

	fmt.Println(unknowns)

	// create teams
	team1, team2 := createTeams(players)

	return []*Team{team1, team2}, nil
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
