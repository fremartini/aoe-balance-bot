package bot

import (
	"aoe-bot/internal/logger"
	playermapper "aoe-bot/internal/player_mapper"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type bot struct {
	logger           *logger.Logger
	commands         map[string]Command
	Session          *discordgo.Session
	steamIdChannelId string
	mapper           *playermapper.PlayerMapper
}

func New(
	logger *logger.Logger,
	token string,
	steamIdChannelId string,
) (*bot, error) {
	discord, err := discordgo.New("Bot " + token)

	if err != nil {
		return nil, err
	}

	return &bot{
		logger:           logger,
		Session:          discord,
		steamIdChannelId: steamIdChannelId,
	}, nil
}

func (b *bot) Run(commands map[string]Command, mapper *playermapper.PlayerMapper) {
	b.commands = commands
	b.mapper = mapper

	b.Session.Identify.Presence.Game.Name = "!help"

	b.Session.AddHandler(b.onMessage)

	b.logger.Info("Starting bot")

	b.Session.Open()
	defer b.Session.Close()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func (b *bot) onMessage(session *discordgo.Session, message *discordgo.MessageCreate) {
	// prevent responding to own messages
	if message.Author.ID == session.State.User.ID {
		return
	}

	// user posted id in the special 'id' channel
	if message.ChannelID == b.steamIdChannelId {
		// content is assumed to be the users steam id
		steamId := message.Content

		b.mapper.AddPlayer(message.Author.ID, steamId)
		return
	}

	split := strings.Split(message.Content, " ")

	action := split[0]

	b.logger.Infof("Handling action: %s", action)

	if action == "!help" {
		s := strings.Builder{}

		for k, c := range b.commands {
			s.WriteString(fmt.Sprintf("%s\t\t\t\t%s\n", k, c.Hint))
		}

		session.ChannelMessageSend(message.ChannelID, s.String())

		return
	}

	command, ok := b.commands[action]

	if !ok {
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Unrecognized command '%s'", action))
		return
	}

	context := &Context{
		UserId:    message.Author.ID,
		ChannelId: message.ChannelID,
		ServerId:  message.GuildID,
	}

	if err := command.Handle(context, split[1:]); err != nil {
		b.logger.Fatal(err.Error())
	}
}
