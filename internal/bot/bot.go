package bot

import (
	"aoe-bot/internal/list"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type bot struct {
	commands            map[*regexp.Regexp]Command
	Session             *discordgo.Session
	prefix              string
	whitelistedChannels []string
}

func New(
	prefix,
	token string,
	whitelistedChannels []string,
) (*bot, error) {
	discord, err := discordgo.New("Bot " + token)

	if err != nil {
		return nil, err
	}

	return &bot{
		Session:             discord,
		prefix:              prefix,
		whitelistedChannels: whitelistedChannels,
	}, nil
}

func (b *bot) onInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionMessageComponent {
		return
	}

	data := i.MessageComponentData()

	// '|' is used as the delimiter. If this is present the command carries additional data
	if strings.Contains(data.CustomID, "|") {
		args := strings.Split(data.CustomID, "|")
		command := args[0]
		rest := args[1:]

		commandWithPrefix := fmt.Sprintf("%s%s", b.prefix, command)
		newArgs := append([]string{commandWithPrefix}, rest...)

		b.tryCommand(commandWithPrefix, i.Message, newArgs)
	}

	// fallback - ignore the "This interaction failed"
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})
}

func (b *bot) Run(commands map[*regexp.Regexp]Command, port *uint) {
	b.commands = commands

	b.Session.Identify.Presence.Game.Name = fmt.Sprintf("%shelp", b.prefix)

	b.Session.AddHandler(b.onMessage)

	b.Session.AddHandler(b.onInteraction)

	log.Print("Starting bot")

	err := b.Session.Open()

	if err != nil {
		panic(err)
	}

	defer b.Session.Close()

	if port != nil {
		m := http.NewServeMux()

		log.Printf("Starting web server on port %d", *port)

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
		log.Print("No port provided. Web server not starting")
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

	b.tryCommand(action, message.Message, args)
}

func (b *bot) tryCommand(action string, message *discordgo.Message, args []string) {
	for k, v := range b.commands {
		if !k.MatchString(action) {
			continue
		}

		log.Printf("Handling action: %s %s", action, args)

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
