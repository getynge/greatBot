package commands

import (
	"github.com/bwmarrin/discordgo"
	"os"
)

var ratecmd = os.Getenv("DISCORD_RATE_CMD")
var rategood = os.Getenv("DISCORD_RATE_GOOD")
var ratebad = os.Getenv("DISCORD_RATE_BAD")

type CommandHandler func(session *discordgo.Session, message *discordgo.MessageCreate)

type CommandSpec struct {
	Handler CommandHandler
	Syntax  map[string]string
}

// A list of commands to be routed by the command router
// Each item in the list is a function, and a spec describing the arguments acceptable
// The format for argument definitions is as such:
// nil denotes that any set of arguments is acceptable, note that in this case the map passed to the function will be nil
// an empty map denotes that no arguments are acceptable
// Any types described in parser can be used in the type spec
// If a single variable can be multiple types, separate them by pipes
// If a variable is optional, end the type list with a ?
// If a variable can occur any number of times, it must be named "*"
// If an infinitely occurring variable is marked as optional, then 0 or more arguments can be passed, otherwise
// 1 or more arguments must be passed
// There are no guarantees about the key names of infinitely occurring variables
var CommandList = map[string]CommandSpec{
	ratecmd: {
		RateUser(rategood, ratebad),
		map[string]string{
			"user": "Mention?",
		},
	},
}
