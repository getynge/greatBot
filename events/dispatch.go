package events

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/getynge/greatBot/commands"
	"strings"
)

type EventDispatcher struct {
	botID  string
	prefix string
	parser commands.Parser
}

func NewDispatcher(botID string, prefix string) EventDispatcher {
	parser, err := commands.InitParser()

	if err != nil {
		panic(err) // parser compilation errors are fatal, so might as well panic
	}

	return EventDispatcher{
		botID,
		prefix,
		parser,
	}
}

func (*EventDispatcher) ReadyHandler(session *discordgo.Session, ready *discordgo.Ready) {
	err := session.UpdateStatus(0, "with itself")
	if err != nil {
		fmt.Println("Could not set status, continuing...")
	}

	servers := session.State.Guilds
	fmt.Printf("Bot listening on %d servers\n", len(servers))
}

func (id *EventDispatcher) ServerJoinHandler(session *discordgo.Session, guilds *discordgo.GuildCreate) {
	if guilds.Guild.Unavailable {
		return
	}

	for _, channel := range guilds.Guild.Channels {
		if channel.ID == guilds.Guild.ID {
			_, _ = session.ChannelMessageSend(channel.ID, fmt.Sprintf("GreatBot is ready! Use %shelp to get a list of commands, whenever support for that rolls around", id.prefix))
			return
		}
	}
}

func (id *EventDispatcher) CommandHandler(session *discordgo.Session, message *discordgo.MessageCreate) {
	user := message.Author
	if user.ID == id.botID {
		return
	}

	if strings.HasPrefix(message.Content, id.prefix) && !strings.HasPrefix(message.Content, fmt.Sprintf("%s ", id.prefix)) {
		id.routeEvent(session, message)
	}
}
