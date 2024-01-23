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
	logger *logger.Logger
}

func New(logger *logger.Logger) *bot {
	return &bot{
		logger: logger,
	}
}

func (b *bot) Run(token string) error {
	discord, err := discordgo.New("Bot " + token)

	if err != nil {
		return err
	}

	discord.AddHandler(b.newMessage)

	discord.Open()
	defer discord.Close()

	b.logger.Info("Starting bot")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	return nil
}

func (b *bot) newMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {
	// prevent responding to own messages
	if message.Author.ID == discord.State.User.ID {
		return
	}

	split := strings.Split(message.Content, " ")

	action := split[0]

	b.logger.Infof("Handling action: %s", action)

	if action == "!help" {
		s := strings.Builder{}

		for k, c := range actions {
			s.WriteString(fmt.Sprintf("%s - %s\n", k, c.hint))
		}

		discord.ChannelMessageSend(message.ChannelID, s.String())

		return
	}

	command, ok := actions[action]

	if !ok {
		discord.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Unrecognized command '%s'", action))
		return
	}

	if err := command.handle(split[1:]); err != nil {
		b.logger.Fatal(err.Error())
	}
}
