package events

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/getynge/greatBot/commands"
	"os"
	"strings"
)

var ratecmd = os.Getenv("DISCORD_RATE_CMD")
var rategood = os.Getenv("DISCORD_RATE_GOOD")
var ratebad = os.Getenv("DISCORD_RATE_BAD")

func validateCommand(command *commands.Command) (err error) {
	spec, has := commands.CommandList[command.Name]

	if !has {
		err = errors.New(fmt.Sprintf("command %s does not exist", command.Name))
		return
	}

	if _, has := spec.Syntax["*"]; len(command.Args) > len(spec.Syntax) && !has {
		err = errors.New("too many arguments")
		return
	}

	index := 0

	for _, typeRestrictions := range spec.Syntax {
		optional := strings.HasSuffix(typeRestrictions, "?")

		if index >= len(command.Args) && optional {
			return
		} else if index >= len(command.Args) {
			err = errors.New("not enough arguments")
			return
		}

		listed := strings.Split(typeRestrictions, "|")
		current := command.Args[index]
		matched := false

		for _, rest := range listed {
			restriction := strings.TrimSuffix(rest, "?")
			matched = (current.Username != "" && restriction == "Username") ||
				(current.Channel != "" && restriction == "Channel") ||
				(current.Mention != "" && restriction == "Mention") ||
				(current.PlainText != "" && restriction == "PlainText") ||
				(current.Quote != "" && restriction == "Quote")
			if matched {
				break
			}
		}

		if !matched && strings.HasSuffix(typeRestrictions, "?") && index != len(command.Args)-1 {
			continue
		} else if !matched {
			errstr := fmt.Sprintf("command syntax is incorrect, correct syntax is %v", spec.Syntax)
			err = errors.New(errstr)
			return
		}
		index++
	}

	return
}

func (id *EventDispatcher) routeEvent(session *discordgo.Session, message *discordgo.MessageCreate) {
	if !strings.HasPrefix(message.Content, id.prefix) {
		return
	}
	command := strings.TrimPrefix(message.Content, id.prefix)
	result, err := id.parser.ParseCommand(command)

	if err != nil {
		session.ChannelMessageSend(message.ChannelID, err.Error())
		return
	}

	if cmd := result.Command; cmd != nil {
		err := validateCommand(result.Command)
		if err != nil {
			session.ChannelMessageSend(message.ChannelID, err.Error())
			return
		}
		commands.CommandList[cmd.Name].Handler(session, message)
	} else if cmds := result.Commands; cmds != nil {
		for _, cmd := range cmds {
			err := validateCommand(cmd)
			if err != nil {
				session.ChannelMessageSend(message.ChannelID, err.Error())
				continue
			}
		}
	}
}
