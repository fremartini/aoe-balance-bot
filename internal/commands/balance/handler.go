package balance

import (
	"aoe-bot/internal/bot"
	"aoe-bot/internal/domain"
	internalErrors "aoe-bot/internal/errors"
	"aoe-bot/internal/list"
	"aoe-bot/internal/logger"
	"fmt"
	"strings"
)

type gameDataProvider interface {
	GetLobbies() ([]*domain.Lobby, error)
	GetPlayers(playerIds []uint) ([]*domain.Player, error)
}

type messageProvider interface {
	ChannelMessageSendReply(channelID, content, messageId, guildId string) error
	ChannelMessageDelete(channelID string, messageID string) error
}

type teamProvider interface {
	CreateTeams(players []*Player) (*Team, *Team)
}

type handler struct {
	gameDataProvider gameDataProvider
	messageProvider  messageProvider
	teamProvider     teamProvider
	logger           *logger.Logger
}

func New(
	gameDataProvider gameDataProvider,
	messageProvider messageProvider,
	teamProvider teamProvider,
	logger *logger.Logger,
) *handler {
	return &handler{
		gameDataProvider: gameDataProvider,
		messageProvider:  messageProvider,
		teamProvider:     teamProvider,
		logger:           logger,
	}
}

func (h *handler) Handle(context *bot.Context, args []string) {
	lobbyId := parseAoeLobbyId(args)

	h.logger.Infof("Trying to find lobby with id %s", lobbyId)

	lobbies, err := h.gameDataProvider.GetLobbies()

	if err != nil {
		h.handleError(err, context)
		return
	}

	lobby, ok := list.FirstWhere(lobbies, func(lobby *domain.Lobby) bool {
		lobbyIdStr := fmt.Sprintf("%d", lobby.Id)
		return lobbyIdStr == lobbyId
	})

	if !ok {
		e := fmt.Sprintf("Public lobby %s not found", lobbyId)
		h.logger.Info(e)

		h.printLobbyNotFound(context, lobbyId)

		return
	}

	memberIds := (**lobby).Members

	h.logger.Infof("Found lobby %s (%s). It has %d players", (*lobby).Title, lobbyId, len(memberIds))

	players, err := h.gameDataProvider.GetPlayers(memberIds)

	if err != nil {
		h.handleError(err, context)
		return
	}

	playersWithELO := list.Map(players, func(p *domain.Player) *Player {
		if p.Rating_1v1 == nil {
			var defaultElo uint = 1000
			p.Rating_1v1 = &defaultElo
		}

		return &Player{
			Rating: *p.Rating_1v1,
			Name:   p.Name,
		}
	})

	team1, team2 := h.teamProvider.CreateTeams(playersWithELO)

	h.printLobbyOutput(context, []*Team{team1, team2}, *lobby)
}

func (h *handler) printLobbyNotFound(context *bot.Context, lobbyId string) {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("<@%s> used %s", context.UserId, context.Command))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf(`Lobby **%s** could not be found`, lobbyId))

	sb.WriteString("\n\nPossible reasons:\n")
	sb.WriteString("* The lobby is private or does not exist\n")
	sb.WriteString("* The external lobby provider is not up to date\n\n")

	sb.WriteString(fmt.Sprintf(`If this is an error, you can [click here to join](https://aoe2lobby.com/j/%s)`, lobbyId))

	h.messageProvider.ChannelMessageSendReply(context.ChannelId, sb.String(), context.MessageId, context.GuildId)

	h.messageProvider.ChannelMessageDelete(context.ChannelId, context.MessageId)
}

func (h *handler) printLobbyOutput(context *bot.Context, teams []*Team, lobby *domain.Lobby) {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("<@%s> used %s", context.UserId, context.Command))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf(`Lobby **%s** (%d)`, lobby.Title, lobby.Id))
	sb.WriteString("\n\n")

	t1 := teams[0]
	t2 := teams[1]

	totalLobbyMembers := len(t1.Players) + len(t2.Players)

	if totalLobbyMembers > 1 {
		for teamNumber, team := range teams {
			players := team.Players

			sb.WriteString(fmt.Sprintf("**Team %d:**\n", teamNumber+1))
			for _, player := range players {
				s := fmt.Sprintf("%s (%d)\n", player.Name, player.Rating)
				sb.WriteString(s)
			}
			sb.WriteString("\n")
		}

		highestEloTeam := 1
		if t2.ELO > t1.ELO {
			highestEloTeam = 2
		}

		diff := abs(int(t1.ELO) - int(t2.ELO))

		sb.WriteString(fmt.Sprintf("ELO difference: **%d** in favor of **Team %d**\n\n", diff, highestEloTeam))
	}

	sb.WriteString(fmt.Sprintf(`[Click here to join](https://aoe2lobby.com/j/%d)`, lobby.Id))

	h.messageProvider.ChannelMessageSendReply(context.ChannelId, sb.String(), context.MessageId, context.GuildId)

	h.messageProvider.ChannelMessageDelete(context.ChannelId, context.MessageId)
}

func (h *handler) handleError(err error, context *bot.Context) {
	switch e := err.(type) {
	default:
		h.logger.Warnf("Unhandlded error %v", err)
	case *internalErrors.ServerError:
		h.logger.Warn(e.Message)
		msg := fmt.Sprintf("**Server error** \n%s", e.Message)
		h.messageProvider.ChannelMessageSendReply(context.ChannelId, msg, context.MessageId, context.GuildId)
	case *internalErrors.NotFoundError:
		h.messageProvider.ChannelMessageSendReply(context.ChannelId, "Unknown player", context.MessageId, context.GuildId)
	case *internalErrors.ApplicationError:
		h.messageProvider.ChannelMessageSendReply(context.ChannelId, e.Message, context.MessageId, context.GuildId)
	}
}

func parseAoeLobbyId(args []string) string {
	fullLobbyId := strings.Split(args[0], "/")
	lobbyId := fullLobbyId[len(fullLobbyId)-1]

	return lobbyId
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
