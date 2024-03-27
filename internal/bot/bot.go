package bot

import (
	"aoe-bot/internal/logger"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type bot struct {
	logger   *logger.Logger
	commands map[*regexp.Regexp]Command
	Session  *discordgo.Session
}

func New(logger *logger.Logger, token string) (*bot, error) {
	discord, err := discordgo.New("Bot " + token)

	if err != nil {
		return nil, err
	}

	return &bot{
		logger:  logger,
		Session: discord,
	}, nil
}

func (b *bot) Run(commands map[*regexp.Regexp]Command) {
	b.commands = commands

	b.Session.Identify.Presence.Game.Name = "!help"

	b.Session.AddHandler(b.onMessage)

	b.logger.Info("Starting bot")

	b.Session.Open()
	defer b.Session.Close()

	m := http.NewServeMux()
	s := http.Server{Addr: ":8080", Handler: m}
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("I'm alive"))
	})

	if err := s.ListenAndServe(); err != nil {
		panic(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func (b *bot) onMessage(session *discordgo.Session, message *discordgo.MessageCreate) {
	// prevent responding to own messages
	if message.Author.ID == session.State.User.ID {
		return
	}

	args := strings.Split(message.Content, " ")

	action := args[0]

	if action == "!help" {
		s := strings.Builder{}

		for k, c := range b.commands {
			if c.Hidden {
				continue
			}

			s.WriteString(fmt.Sprintf("%s\t\t\t\t%s\n", k, c.Hint))
		}

		session.ChannelMessageSend(message.ChannelID, s.String())

		return
	}

	context := &Context{
		UserId:    message.Author.ID,
		ChannelId: message.ChannelID,
		GuildId:   message.GuildID,
		MessageId: message.ID,
	}

	for k, v := range b.commands {
		if k.MatchString(action) {
			b.logger.Infof("Handling action: %s %s", action, args)

			v.Handle(context, args)
		}
	}
}
