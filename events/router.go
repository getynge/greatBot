package events

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/getynge/greatBot/commands"
	"os"
)

var ratecmd = os.Getenv("DISCORD_RATE_CMD")
var rategood = os.Getenv("DISCORD_RATE_GOOD")
var ratebad = os.Getenv("DISCORD_RATE_BAD")

func (id *EventDispatcher) routeEvent(session *discordgo.Session, message *discordgo.MessageCreate) {
	tree, errs := commands.ParseCommand(id.prefix, message.Content)
	if errs != nil {
		for _, err := range errs {
			session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Error: %s", err.Error()))
		}
		return
	}

	commands.Commands[tree.Literal].CMD(session, message, tree)
}
