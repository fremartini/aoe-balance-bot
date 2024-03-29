package main

import (
	"aoe-bot/internal/bot"
	"aoe-bot/internal/cache"
	"aoe-bot/internal/commands/balance"
	"aoe-bot/internal/discord"
	"aoe-bot/internal/domain"
	"aoe-bot/internal/librematch"
	"aoe-bot/internal/logger"
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const prefix = "!"

func New(
	session *discordgo.Session,
	logger *logger.Logger,
	playerCache *cache.Cache[uint, *domain.Player]) map[*regexp.Regexp]bot.Command {
	return map[*regexp.Regexp]bot.Command{
		regexp.MustCompile(`aoe2de:\/\/0/\d*`): {
			Handle: func(context *bot.Context, args []string) {
				discordAPI := discord.New(session)

				lobbyId := parseAoeLobbyId(args)

				librematchApi := librematch.New(logger, playerCache)

				handler := balance.New(librematchApi, discordAPI, discordAPI, logger)

				handler.Handle(context, lobbyId)
			},
			Hint:   "Create two teams of players in a lobby",
			Hidden: true,
		},

		withPrefix("balance"): {
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

				handler := balance.New(librematchApi, discordAPI, discordAPI, logger)

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

func withPrefix(cmd string) *regexp.Regexp {
	s := fmt.Sprintf("%s%s", prefix, cmd)
	return regexp.MustCompile(s)
}
