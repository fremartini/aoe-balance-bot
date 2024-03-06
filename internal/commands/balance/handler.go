package balance

import (
	"aoe-bot/internal/bot"
	"aoe-bot/internal/discord"
	"aoe-bot/internal/domain"
	internalErrors "aoe-bot/internal/errors"
	"aoe-bot/internal/list"
	"aoe-bot/internal/logger"
	"errors"
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
	ChannelMessageSendReply(channelID, content, messageId, guildId string)
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

func (h *handler) Handle(context *bot.Context, lobbyId string) {
	h.logger.Infof("Trying to find lobby with id %s", lobbyId)

	lobbies, err := h.dataProvider.GetLobbies()

	if err != nil {
		h.handleError(err, context)
		return
	}

	lobby, found := list.FirstWhere(lobbies, func(lobby *domain.Lobby) bool {
		s := fmt.Sprintf("%d", lobby.Id)
		return s == lobbyId
	})

	if !found {
		e := fmt.Sprintf("Public lobby %s not found", lobbyId)
		internalError := internalErrors.NewApplicationError(e)

		h.handleError(internalError, context)
		return
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
}

func (h *handler) printOutput(context bot.Context, teams []*Team, lobbyId string) {
	var sb strings.Builder

	gameIdStr := fmt.Sprintf(`Game id **%s**`, lobbyId)
	sb.WriteString(gameIdStr)
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

		diffStr := fmt.Sprintf("ELO difference: **%d** in favor of **Team %d**\n\n", diff, highestEloTeam)
		sb.WriteString(diffStr)
	}

	joinStr := fmt.Sprintf(`[Click here to join](https://aoe2lobby.com/j/%s)`, lobbyId)
	sb.WriteString(joinStr)

	h.messageProvider.ChannelMessageSendReply(context.ChannelId, sb.String(), context.MessageId, context.GuildId)
}

func (h *handler) handleError(err error, context *bot.Context) {
	var serverErr *internalErrors.ServerError
	if errors.As(err, &serverErr) {
		h.messageProvider.ChannelMessageSendReply(context.ChannelId, "Server error", context.MessageId, context.GuildId)
		return
	}

	var notFoundErr *internalErrors.NotFoundError
	if errors.As(err, &notFoundErr) {
		h.messageProvider.ChannelMessageSendReply(context.ChannelId, "Unknown player", context.MessageId, context.GuildId)
		return
	}

	var applicationErr *internalErrors.ApplicationError
	if errors.As(err, &applicationErr) {
		h.messageProvider.ChannelMessageSendReply(context.ChannelId, applicationErr.Message, context.MessageId, context.GuildId)
		return
	}

	h.logger.Warnf("Unhandlded error %v", err)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
