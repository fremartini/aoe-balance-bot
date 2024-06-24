package discord

import (
	"github.com/bwmarrin/discordgo"
)

type api struct {
	session *discordgo.Session
}

func New(session *discordgo.Session) *api {
	return &api{
		session: session,
	}
}

func (p *api) ChannelMessageSend(channelID, content string) error {
	_, err := p.session.ChannelMessageSend(channelID, content)

	return err
}

func (p *api) ChannelMessageSendReply(channelID, content, messageId, guildId string) error {
	_, err := p.session.ChannelMessageSendReply(channelID, content, &discordgo.MessageReference{
		MessageID: messageId,
		ChannelID: channelID,
		GuildID:   guildId,
	})

	return err
}

func (p *api) ChannelMessageDelete(channelID string, messageID string) error {
	return p.session.ChannelMessageDelete(channelID, messageID)
}
