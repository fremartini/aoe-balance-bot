package balance

import (
	"aoe-bot/internal/bot"
	"aoe-bot/internal/discord"
	"aoe-bot/internal/domain"
	"aoe-bot/internal/errors"
	"aoe-bot/internal/list"
	"aoe-bot/internal/logger"
	"cmp"
	"fmt"
	"slices"
	"strings"
)

type gameDataProvider interface {
	GetLobbies() ([]*domain.Lobby, error)
	GetPlayer(playerId uint) (*domain.Player, error)
}

type userDataProvider interface {
	FindUserVoiceChannel(guildId, userid string) (string, error)
	FindUsersInVoiceChannel(serverId, channelId string) ([]*discord.User, error)
}

type messageProvider interface {
	ChannelMessageSend(channelID, content string)
}

type handler struct {
	dataProvider        gameDataProvider
	discordInfoProvider userDataProvider
	messageProvider     messageProvider
	logger              *logger.Logger
}

func New(
	gameDataProvider gameDataProvider,
	userDataProvider userDataProvider,
	messageProvider messageProvider,
	logger *logger.Logger) *handler {
	return &handler{
		dataProvider:        gameDataProvider,
		discordInfoProvider: userDataProvider,
		messageProvider:     messageProvider,
		logger:              logger,
	}
}

func (h *handler) Handle(context *bot.Context, lobbyId string) error {
	h.logger.Infof("Trying to find lobby with id %s", lobbyId)

	lobbies, err := h.dataProvider.GetLobbies()

	if err != nil {
		return err
	}

	lobby, found := list.FirstWhere(lobbies, func(lobby *domain.Lobby) bool {
		s := fmt.Sprintf("%d", lobby.Id)
		return s == lobbyId
	})

	if !found {
		e := fmt.Sprintf("Lobby id %s not found", lobbyId)
		return errors.NewApplicationError(e)
	}

	memberIds := (**lobby).Members

	h.logger.Infof("Found lobby with id %s. It has %d players", lobbyId, len(memberIds))

	players := []*Player{}
	for _, memberId := range memberIds {
		p, err := h.dataProvider.GetPlayer(memberId)

		if err != nil {
			h.logger.Warnf("Error while fetching player: %s", err)
		}

		h.logger.Infof("Got data for profile id %d. Name %s", memberId, p.Name)

		if p.Rating_1v1 == nil {
			var defaultElo uint = 1000
			p.Rating_1v1 = &defaultElo
		}

		players = append(players, &Player{
			Rating: *p.Rating_1v1,
			Name:   p.Name,
		})
	}

	t1, t2 := createTeams(players)

	h.printTeams(*context, []*Team{t1, t2})

	return nil
}

func createTeams(players []*Player) (*Team, *Team) {
	t1 := &Team{}
	t2 := &Team{}

	var t1Rating uint = 0
	var t2Rating uint = 0

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

	t1.ELO = t1Rating
	t2.ELO = t2Rating

	return t1, t2
}

func (h *handler) printTeams(context bot.Context, teams []*Team) {
	var sb strings.Builder
	for teamNumber, team := range teams {
		players := team.Players

		sb.WriteString(fmt.Sprintf("Team %d:\n", teamNumber+1))
		for _, player := range players {
			s := fmt.Sprintf("%s (%d)\n", player.Name, player.Rating)
			sb.WriteString(s)
		}
		sb.WriteString("\n")
	}

	diff := abs(int(teams[0].ELO) - int(teams[1].ELO))
	diffStr := fmt.Sprintf("ELO difference: %d\n", diff)
	sb.WriteString(diffStr)

	h.messageProvider.ChannelMessageSend(context.ChannelId, sb.String())
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
