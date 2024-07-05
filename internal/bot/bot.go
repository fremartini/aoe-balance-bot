package bot

import (
	"aoe-bot/internal/list"
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
	logger              *logger.Logger
	commands            map[*regexp.Regexp]Command
	Session             *discordgo.Session
	prefix              string
	whitelistedChannels []string
}

func New(
	logger *logger.Logger,
	prefix,
	token string,
	whitelistedChannels []string,
) (*bot, error) {
	discord, err := discordgo.New("Bot " + token)

	if err != nil {
		return nil, err
	}

	return &bot{
		logger:              logger,
		Session:             discord,
		prefix:              prefix,
		whitelistedChannels: whitelistedChannels,
	}, nil
}

func (b *bot) Run(commands map[*regexp.Regexp]Command, port *uint) {
	b.commands = commands

	b.Session.Identify.Presence.Game.Name = fmt.Sprintf("%shelp", b.prefix)

	b.Session.AddHandler(b.onMessage)

	b.logger.Info("Starting bot")

	err := b.Session.Open()

	if err != nil {
		panic(err)
	}

	defer b.Session.Close()

	if port != nil {
		m := http.NewServeMux()

		b.logger.Infof("Starting web server on port %d", *port)

		server := http.Server{
			Addr:    fmt.Sprintf(":%d", *port),
			Handler: m,
		}
		m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("I'm alive"))
		})

		if err := server.ListenAndServe(); err != nil {
			panic(err)
		}
	} else {
		b.logger.Info("No port provided. Web server not starting")
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

	// message was sent in a channel that was not whitelisted.
	// if there are no entries in the list, no whitelisting should be applied
	if len(b.whitelistedChannels) > 0 && !list.Contains(b.whitelistedChannels, message.ChannelID) {
		return
	}

	args := strings.Split(message.Content, " ")

	action := args[0]

	if action == fmt.Sprintf("%shelp", b.prefix) {
		b.printHelp(session, message)

		return
	}

	for k, v := range b.commands {
		if !k.MatchString(action) {
			continue
		}

		b.logger.Infof("Handling action: %s %s", action, args)

		context := &Context{
			UserId:    message.Author.ID,
			ChannelId: message.ChannelID,
			GuildId:   message.GuildID,
			MessageId: message.ID,
			Command:   action,
		}

		v.Handle(context, args)

		break
	}
}

func (b *bot) printHelp(session *discordgo.Session, message *discordgo.MessageCreate) {
	builder := strings.Builder{}

	for k, c := range b.commands {
		if c.Hidden {
			continue
		}

		builder.WriteString(fmt.Sprintf("%s\t\t\t\t%s\n", k, c.Hint))
	}

	session.ChannelMessageSend(message.ChannelID, builder.String())
}
