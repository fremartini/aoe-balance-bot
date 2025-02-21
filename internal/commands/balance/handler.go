package balance

import (
	"aoe-bot/internal/bot"
	"aoe-bot/internal/domain"
	internalErrors "aoe-bot/internal/errors"
	"aoe-bot/internal/list"
	"aoe-bot/internal/ui"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type gameDataProvider interface {
	GetLobbies() ([]*domain.Lobby, error)
	GetPlayers(playerIds []uint) ([]*domain.Player, error)
}

type messageProvider interface {
	ChannelMessageSendReply(channelID, content, messageId, guildId string) error
	ChannelMessageDelete(channelID string, messageID string) error
	ChannelMessageSendContentWithButton(channelId, content string, buttons []*ui.Button) error
}

type teamProvider interface {
	CreateTeams(players []*Player) (*Team, *Team)
}

type handler struct {
	gameDataProvider gameDataProvider
	messageProvider  messageProvider
	teamProvider     teamProvider
}

func New(
	gameDataProvider gameDataProvider,
	messageProvider messageProvider,
	teamProvider teamProvider,
) *handler {
	return &handler{
		gameDataProvider: gameDataProvider,
		messageProvider:  messageProvider,
		teamProvider:     teamProvider,
	}
}

func (h *handler) Handle(context *bot.Context, args []string) error {
	lobbyId := parseAoeLobbyId(args)

	log.Printf("Trying to find lobby with id %s", lobbyId)

	lobbies, err := h.gameDataProvider.GetLobbies()

	if err != nil {
		h.handleError(err, context)
		return nil
	}

	lobby, ok := list.FirstWhere(lobbies, func(lobby *domain.Lobby) bool {
		lobbyIdStr := fmt.Sprintf("%d", lobby.Id)
		return lobbyIdStr == lobbyId
	})

	if !ok {
		e := fmt.Sprintf("Public lobby %s not found", lobbyId)
		log.Print(e)

		return h.printLobbyNotFound(context, lobbyId)
	}

	memberIds := (**lobby).Members

	log.Printf("Found lobby %s (%s). It has %d players", (*lobby).Title, lobbyId, len(memberIds))

	players, err := h.gameDataProvider.GetPlayers(memberIds)

	if err != nil {
		h.handleError(err, context)
		return nil
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

	return h.printLobbyOutput(context, []*Team{team1, team2}, *lobby)
}

func (h *handler) printLobbyNotFound(context *bot.Context, lobbyId string) error {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("<@%s> used %s", context.UserId, context.Command))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf(`Lobby **%s** could not be found`, lobbyId))

	sb.WriteString("\n\nPossible reasons:\n")
	sb.WriteString("* The lobby is private or does not exist\n")
	sb.WriteString("* The external lobby provider is not up to date\n\n")

	retryButton := &ui.Button{
		Label: "Retry",
		Style: uint(discordgo.PrimaryButton),
		Id:    fmt.Sprintf("%s|%s", "balance", lobbyId),
	}

	joinButton := &ui.Button{
		Label: "Join",
		Style: uint(discordgo.LinkButton),
		Url:   fmt.Sprintf(`https://aoe2lobby.com/j/%s`, lobbyId),
	}

	err := h.messageProvider.ChannelMessageSendContentWithButton(context.ChannelId, sb.String(), []*ui.Button{retryButton, joinButton})

	if err != nil {
		return err
	}

	return h.messageProvider.ChannelMessageDelete(context.ChannelId, context.MessageId)
}

func (h *handler) printLobbyOutput(context *bot.Context, teams []*Team, lobby *domain.Lobby) error {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("<@%s> used %s", context.UserId, context.Command))
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf(`Lobby **%s** (%d)`, lobby.Title, lobby.Id))
	sb.WriteString("\n\n")

	t1 := teams[0]
	t2 := teams[1]

	// sort teams highest ELO to lowest
	sort.Slice(t1.Players, func(i, j int) bool {
		return t1.Players[i].Rating > t1.Players[j].Rating
	})

	sort.Slice(t2.Players, func(i, j int) bool {
		return t2.Players[i].Rating > t2.Players[j].Rating
	})

	totalLobbyMembers := len(t1.Players) + len(t2.Players)

	if totalLobbyMembers > 1 {
		for teamNumber, team := range teams {
			sb.WriteString(fmt.Sprintf("**Team %d:**\n", teamNumber+1))
			for _, player := range team.Players {
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

	recalculateButton := &ui.Button{
		Label: "Recalculate",
		Style: uint(discordgo.PrimaryButton),
		Id:    fmt.Sprintf("%s|%v", "balance", lobby.Id),
	}

	joinButton := &ui.Button{
		Label: "Join",
		Style: uint(discordgo.LinkButton),
		Url:   fmt.Sprintf(`https://aoe2lobby.com/j/%v`, lobby.Id),
	}

	err := h.messageProvider.ChannelMessageSendContentWithButton(context.ChannelId, sb.String(), []*ui.Button{recalculateButton, joinButton})

	if err != nil {
		return err
	}

	return h.messageProvider.ChannelMessageDelete(context.ChannelId, context.MessageId)
}

func (h *handler) handleError(err error, context *bot.Context) {
	switch e := err.(type) {
	default:
		log.Printf("Unhandlded error %v", err)
	case *internalErrors.ServerError:
		log.Print(e.Message)
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
