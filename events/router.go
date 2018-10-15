package events

import (
	"github.com/bwmarrin/discordgo"
	"github.com/getynge/greatBot/commands"
	"os"
	"strings"
)

var ratecmd = os.Getenv("DISCORD_RATE_CMD")
var rategood = os.Getenv("DISCORD_RATE_GOOD")
var ratebad = os.Getenv("DISCORD_RATE_BAD")

func (id *EventDispatcher) routeEvent(session *discordgo.Session, message *discordgo.MessageCreate) {
	if !strings.HasPrefix(message.Content, id.prefix) {
		return
	}
	command := strings.TrimPrefix(message.Content, id.prefix)
	result := commands.ParseCommand(command)
	if result != nil && result.Command != nil && result.Command.Name == ratecmd {
		commands.RateUser(rategood, ratebad)(session, message)
	}
}
