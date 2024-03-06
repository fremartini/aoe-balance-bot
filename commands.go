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
	"strings"

	"github.com/bwmarrin/discordgo"
)

const prefix = "!"

func New(
	session *discordgo.Session,
	logger *logger.Logger,
	playerCache *cache.Cache[uint, *domain.Player]) map[string]bot.Command {
	return map[string]bot.Command{
		withPrefix("balance"): {
			Handle: func(context *bot.Context, args []string) {
				discordAPI := discord.New(session)

				if len(args) == 0 {
					discordAPI.ChannelMessageSendReply(context.ChannelId, "Missing game id", context.MessageId, context.GuildId)
					return
				}

				fullLobbyId := strings.Split(args[0], "/")
				lobbyId := fullLobbyId[len(fullLobbyId)-1]

				librematchApi := librematch.New(logger, playerCache)

				handler := balance.New(librematchApi, discordAPI, discordAPI, logger)

				handler.Handle(context, lobbyId)
			},
			Hint: "Create two teams of players in a lobby",
		},
	}
}

func withPrefix(cmd string) string {
	return fmt.Sprintf("%s%s", prefix, cmd)
}
