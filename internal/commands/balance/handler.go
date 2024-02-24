package balance

import (
	"aoe-bot/internal/bot"
	"aoe-bot/internal/discord"
	"aoe-bot/internal/domain"
	"aoe-bot/internal/errors"
	"aoe-bot/internal/list"
	"aoe-bot/internal/logger"
	"fmt"
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

	t1, t2 := CreateTeamsBruteForce(players)

	h.printOutput(*context, []*Team{t1, t2}, lobbyId)

	return nil
}

func (h *handler) printOutput(context bot.Context, teams []*Team, lobbyId string) {
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

	highestEloTeam := 1
	if teams[1].ELO > teams[0].ELO {
		highestEloTeam = 2
	}

	diffStr := fmt.Sprintf("ELO difference: %d in favor of team %d\n\n", diff, highestEloTeam)
	sb.WriteString(diffStr)

	joinStr := fmt.Sprintf(`[Click here to join](https://aoe2lobby.com/j/%s)`, lobbyId)
	sb.WriteString(joinStr)

	h.messageProvider.ChannelMessageSend(context.ChannelId, sb.String())
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
