package discord

import (
	"aoe-bot/internal/logger"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type api struct {
	session *discordgo.Session
	logger  *logger.Logger
}

func New(
	session *discordgo.Session,
	logger *logger.Logger,
) *api {
	return &api{
		session: session,
		logger:  logger,
	}
}

func (a *api) ChannelMessageSend(channelID, content string) error {
	_, err := a.session.ChannelMessageSend(channelID, content)

	if err != nil {
		msg := fmt.Sprintf("Could not send message: %s", err)
		a.logger.Warn(msg)
	}

	return err
}

func (a *api) ChannelMessageSendReply(channelID, content, messageId, guildId string) error {
	_, err := a.session.ChannelMessageSendReply(channelID, content, &discordgo.MessageReference{
		MessageID: messageId,
		ChannelID: channelID,
		GuildID:   guildId,
	})

	if err != nil {
		msg := fmt.Sprintf("Could not reply to message: %s", err)
		a.logger.Warn(msg)
	}

	return err
}

func (a *api) ChannelMessageDelete(channelID, messageID string) error {
	err := a.session.ChannelMessageDelete(channelID, messageID)

	if err != nil {
		msg := fmt.Sprintf("Could not delete message: %s", err)
		a.logger.Warn(msg)
	}

	return err
}
