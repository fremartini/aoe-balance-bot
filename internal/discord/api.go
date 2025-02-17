package discord

import (
	"fmt"

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
	err := a.session.ChannelMessageDelete(channelID, messageID)

	return err
}

func (a *api) ChannelMessageSendContentWithButton(channelID, buttonLabel, payload, content string) error {
	customIdWithPayload := fmt.Sprintf("%s|%s", "balance", payload)

	button := &discordgo.Button{
		Label:    buttonLabel,
		Style:    discordgo.PrimaryButton,
		CustomID: customIdWithPayload,
	}

	_, err := a.session.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Content: content,
		Components: []discordgo.MessageComponent{
			&discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{button},
			},
		},
	})

	return err
}
