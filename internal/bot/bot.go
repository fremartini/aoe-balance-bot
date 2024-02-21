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
	logger           *logger.Logger
	commands         map[string]Command
	Session          *discordgo.Session
	aoe2lobbyIdRegex *regexp.Regexp
}

func New(
	logger *logger.Logger,
	token string,
) (*bot, error) {
	discord, err := discordgo.New("Bot " + token)

	if err != nil {
		return nil, err
	}

	r := regexp.MustCompile(`aoe2de:\/\/0/\d*`)

	return &bot{
		logger:           logger,
		Session:          discord,
		aoe2lobbyIdRegex: r,
	}, nil
}

func (b *bot) Run(commands map[string]Command) {
	b.commands = commands

	b.Session.Identify.Presence.Game.Name = "!help"

	b.Session.AddHandler(b.onMessage)

	b.logger.Info("Starting bot")

	b.Session.Open()
	defer b.Session.Close()

	m := http.NewServeMux()
	s := http.Server{Addr: ":8080", Handler: m}
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		b.logger.Info("Received ping")
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
			s.WriteString(fmt.Sprintf("%s\t\t\t\t%s\n", k, c.Hint))
		}

		session.ChannelMessageSend(message.ChannelID, s.String())

		return
	}

	if b.aoe2lobbyIdRegex.MatchString(action) {
		// user pasted an aoe2 lobby id into the chat. Treat it as !balance command
		groups := b.aoe2lobbyIdRegex.FindStringSubmatch(action)

		action = "!balance"

		args = []string{args[0], groups[0]}
	}

	command, ok := b.commands[action]

	if !ok {
		return
	}

	b.logger.Infof("Handling action: %s %s", action, args[1:])

	context := &Context{
		UserId:    message.Author.ID,
		ChannelId: message.ChannelID,
		ServerId:  message.GuildID,
	}

	if err := command.Handle(context, args[1:]); err != nil {
		b.logger.Fatal(err.Error())
	}
}
