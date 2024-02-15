package bot

import (
	"aoe-bot/internal/logger"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type bot struct {
	logger   *logger.Logger
	commands map[string]Command
	Session  *discordgo.Session
}

func New(
	logger *logger.Logger,
	token string,
) (*bot, error) {
	discord, err := discordgo.New("Bot " + token)

	if err != nil {
		return nil, err
	}

	return &bot{
		logger:  logger,
		Session: discord,
	}, nil
}

func (b *bot) Run(commands map[string]Command) {
	b.commands = commands

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

	split := strings.Split(message.Content, " ")

	action := split[0]

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
		return
	}

	b.logger.Infof("Handling action: %s", action)

	context := &Context{
		UserId:    message.Author.ID,
		ChannelId: message.ChannelID,
		ServerId:  message.GuildID,
	}

	if err := command.Handle(context, split[1:]); err != nil {
		b.logger.Fatal(err.Error())
	}
}
