package events

import (
	"github.com/bwmarrin/discordgo"
	"github.com/getynge/greatBot/commands"
	"github.com/getynge/greatBot/util"
	"os"
)

var ratecmd = os.Getenv("DISCORD_RATE_CMD")
var rategood = os.Getenv("DISCORD_RATE_GOOD")
var ratebad = os.Getenv("DISCORD_RATE_BAD")

func (id *EventDispatcher) routeEvent(session *discordgo.Session, message *discordgo.MessageCreate) {
	command := util.ExtractCommand(id.prefix, message.Content)
	if command == nil {
		return
	}

	switch command[0] {
	case ratecmd:
		commands.RateUser(session, message, rategood, ratebad)
	}
}
