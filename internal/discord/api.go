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

func (p *api) ChannelMessageSend(channelID, content string) {
	p.session.ChannelMessageSend(channelID, content)
}

func (p *api) ChannelMessageSendReply(channelID, content, messageId, guildId string) {
	p.session.ChannelMessageSendReply(channelID, content, &discordgo.MessageReference{
		MessageID: messageId,
		ChannelID: channelID,
		GuildID:   guildId,
	})
}

func (p *api) ChannelMessageEdit(channelID, content, messageId string) {
	p.session.ChannelMessageEdit(channelID, messageId, content)
}
