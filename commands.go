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

	"github.com/bwmarrin/discordgo"
)

const (
	aoe2LobbyRegex = `aoe2de:\/\/0/\d*`
)

func New(
	session *discordgo.Session,
	logger *logger.Logger,
	playerCache *cache.Cache[uint, *domain.Player],
	prefix string,
) map[*regexp.Regexp]bot.Command {
	return map[*regexp.Regexp]bot.Command{
		regexp.MustCompile(aoe2LobbyRegex): {
			Handle: func(context *bot.Context, args []string) {
				discordAPI := discord.New(session, logger)

				librematchApi := librematch.New(logger, playerCache)

				teamBalanceStrategy := strategies.NewBruteForce()

				handler := balance.New(librematchApi, discordAPI, teamBalanceStrategy, logger)

				handler.Handle(context, args)
			},
			Hint:   "Create two teams of players in a lobby",
			Hidden: true,
		},
		withPrefix(prefix, "balance"): {
			Handle: func(context *bot.Context, args []string) {
				// discard command name
				args = args[1:]

				discordAPI := discord.New(session, logger)

				if len(args) == 0 {
					discordAPI.ChannelMessageSendReply(context.ChannelId, "Missing game id", context.MessageId, context.GuildId)
					return
				}

				librematchApi := librematch.New(logger, playerCache)

				teamBalanceStrategy := strategies.NewBruteForce()

				handler := balance.New(librematchApi, discordAPI, teamBalanceStrategy, logger)

				handler.Handle(context, args)
			},
			Hint:   "Create two teams of players in a lobby",
			Hidden: false,
		},
	}
}

func withPrefix(prefix, cmd string) *regexp.Regexp {
	s := fmt.Sprintf("%s%s", prefix, cmd)
	return regexp.MustCompile(s)
}
