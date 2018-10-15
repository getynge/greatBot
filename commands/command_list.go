package commands

import (
	"github.com/bwmarrin/discordgo"
	"os"
)

type Node struct {
	ExprType ExpressionType
	Literal  string
	Next     *Node
}

type ExpressionType int

// the "tree" being passed to a CommandFunc is the parsed AST with the first leaf removed
type CommandFunc func(session *discordgo.Session, message *discordgo.MessageCreate, tree *Node)

const (
	CommandName ExpressionType = iota // commands, in the form of a string
	UserName                          // a user name, in the form of username#usernumber
	UserMention                       // a user mention, in the form of <@0000000000000000000>
	ChannelName                       // a channel name, in the form of #name
	Literal                           // any string

	// meta expression types below, these types cannot land in the ast but change the behavior of the parser

	Any        // bypass all checks and assume type is a literal
	RepeatLast // command has an infinite number of arguments, all of which are of the last type on the command list. needs to be last in list
	RepeatAll  // command arguments can be repeated, in order, indefinitely. needs to be last in list
)

type Command struct {
	Expects [][]ExpressionType // every permutation of acceptable argument types
	CMD     CommandFunc
}

var Commands = map[string]Command{
	os.Getenv("DISCORD_RATE_CMD"): {
		[][]ExpressionType{
			{
				UserMention,
			},
			{
				UserName,
			},
			{},
		},
		RateUser(os.Getenv("DISCORD_RATE_GOOD"), os.Getenv("DISCORD_RATE_BAD")),
	},
}
