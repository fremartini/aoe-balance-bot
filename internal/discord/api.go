package discord

import (
	"errors"

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

	return "", errors.New("user is not in a voice channel")
}

func (p *api) FindUsersInVoiceChannel(serverId, channelId string) ([]*string, error) {
	guild, err := p.session.State.Guild(serverId)

	if err != nil {
		return []*string{}, err
	}

	discordIds := []*string{}

	for _, vs := range guild.VoiceStates {
		if vs.ChannelID != channelId {
			continue
		}

		user := vs.Member.User

		discordIds = append(discordIds, &user.ID)
	}

	return discordIds, nil
}