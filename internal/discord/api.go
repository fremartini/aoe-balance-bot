package discord

import (
	"aoe-bot/internal/errors"

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

func (p *api) FindUserVoiceChannel(serverId, userId string) (string, error) {
	guild, err := p.session.State.Guild(serverId)

	if err != nil {
		return "", err
	}

	for _, vs := range guild.VoiceStates {
		if vs.UserID != userId {
			continue
		}

		return vs.ChannelID, nil
	}

	return "", errors.NewApplicationError("You are not in a voice channel")
}

func (p *api) FindUsersInVoiceChannel(serverId, channelId string) ([]*User, error) {
	guild, err := p.session.State.Guild(serverId)

	if err != nil {
		return []*User{}, err
	}

	users := []*User{}

	for _, vs := range guild.VoiceStates {
		if vs.ChannelID != channelId {
			continue
		}

		discordUser := vs.Member.User

		user := User{
			Username: discordUser.Username,
			Id:       discordUser.ID,
		}

		users = append(users, &user)
	}

	return users, nil
}
