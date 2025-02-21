package discord

import (
	"aoe-bot/internal/list"
	"aoe-bot/internal/ui"

	"github.com/bwmarrin/discordgo"
)

type api struct {
	session *discordgo.Session
}

func New(
	session *discordgo.Session,
) *api {
	return &api{
		session: session,
	}
}

func (a *api) ChannelMessageSend(channelID, content string) error {
	_, err := a.session.ChannelMessageSend(channelID, content)

	return err
}

func (a *api) ChannelMessageSendReply(channelID, content, messageId, guildId string) error {
	_, err := a.session.ChannelMessageSendReply(channelID, content, &discordgo.MessageReference{
		MessageID: messageId,
		ChannelID: channelID,
		GuildID:   guildId,
	})

	return err
}

func (a *api) ChannelMessageDelete(channelID, messageID string) error {
	return a.session.ChannelMessageDelete(channelID, messageID)
}

func (a *api) ChannelMessageSendContentWithButton(channelId, content string, buttons []*ui.Button) error {
	goButtons := list.Map(buttons, func(b *ui.Button) discordgo.MessageComponent {
		return &discordgo.Button{
			Label:    b.Label,
			Style:    discordgo.ButtonStyle(b.Style),
			CustomID: b.Id,
			URL:      b.Url,
		}
	})
	_, err := a.session.ChannelMessageSendComplex(channelId, &discordgo.MessageSend{
		Content: content,
		Components: []discordgo.MessageComponent{
			&discordgo.ActionsRow{
				Components: goButtons,
			},
		},
	})

	return err
}
