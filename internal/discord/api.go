package discord

import "github.com/bwmarrin/discordgo"

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
