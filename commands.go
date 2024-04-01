package main

import (
	"aoe-bot/internal/bot"
	"aoe-bot/internal/cache"
	"aoe-bot/internal/commands/balance"
	"aoe-bot/internal/commands/balance/strategies"
	"aoe-bot/internal/discord"
	"aoe-bot/internal/domain"
	"aoe-bot/internal/librematch"
	"aoe-bot/internal/logger"
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func New(
	session *discordgo.Session,
	logger *logger.Logger,
	playerCache *cache.Cache[uint, *domain.Player],
	prefix string,
) map[*regexp.Regexp]bot.Command {
	return map[*regexp.Regexp]bot.Command{
		regexp.MustCompile(`aoe2de:\/\/0/\d*`): {
			Handle: func(context *bot.Context, args []string) {
				discordAPI := discord.New(session)

				lobbyId := parseAoeLobbyId(args)

				librematchApi := librematch.New(logger, playerCache)

				teamStrategy := strategies.NewBruteForce()

				handler := balance.New(librematchApi, discordAPI, teamStrategy, logger)

				handler.Handle(context, lobbyId)
			},
			Hint:   "Create two teams of players in a lobby",
			Hidden: true,
		},

		withPrefix(prefix, "balance"): {
			Handle: func(context *bot.Context, args []string) {
				// discard command
				args = args[1:]

				discordAPI := discord.New(session)

				if len(args) == 0 {
					discordAPI.ChannelMessageSendReply(context.ChannelId, "Missing game id", context.MessageId, context.GuildId)
					return
				}

				lobbyId := parseAoeLobbyId(args)

				librematchApi := librematch.New(logger, playerCache)

				teamStrategy := strategies.NewBruteForce()

				handler := balance.New(librematchApi, discordAPI, teamStrategy, logger)

				handler.Handle(context, lobbyId)
			},
			Hint:   "Create two teams of players in a lobby",
			Hidden: false,
		},
	}
}

func parseAoeLobbyId(args []string) string {
	fullLobbyId := strings.Split(args[0], "/")
	lobbyId := fullLobbyId[len(fullLobbyId)-1]

	return lobbyId
}

func withPrefix(prefix, cmd string) *regexp.Regexp {
	s := fmt.Sprintf("%s%s", prefix, cmd)
	return regexp.MustCompile(s)
}
