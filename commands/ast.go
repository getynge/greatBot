package commands

import (
	"errors"
	"fmt"
	"math"
)

func cmdError(message string, command interface{}) error {
	return errors.New(fmt.Sprintf(message, command))
}

func ParseCommand(prefix string, command string) (out *Node, errs []error) {
	var selectedCommand *Command
	lexed := LexCommand(prefix, command)
	errs = make([]error, 0)
	switch lexed[0].(type) {
	case string:
		command, exists := Commands[lexed[0].(string)]
		if !exists {
			errs = append(errs, cmdError("%+v is not a valid command", lexed[0]))
			return
		}
		out = &Node{
			CommandName,
			lexed[0].(string),
			nil,
		}
		selectedCommand = &command
	default:
		errs = append(errs, cmdError("%+v is not a valid command", lexed[0]))
		return
	}

	var args []interface{}

	if len(lexed) == 1 {
		args = make([]interface{}, 0)
	} else {
		args = lexed[1:]
	}

	matchedIndex := -1
	minsize := math.MaxInt8
	maxsize := 0
	for i, expectedGroup := range selectedCommand.Expects {
		adjustedSize := len(expectedGroup)
		if adjustedSize != 0 && (expectedGroup[len(expectedGroup)-1] == RepeatAll || expectedGroup[len(expectedGroup)-1] == RepeatLast) {
			adjustedSize = len(expectedGroup) - 1
			if len(args)%adjustedSize == 0 {
				matchedIndex = i
				break
			}
			continue
		}
		if adjustedSize > maxsize {
			maxsize = adjustedSize
		}
		if adjustedSize < minsize {
			minsize = adjustedSize
		}
		if len(args) == len(expectedGroup) {
			matchedIndex = i
			break
		}
	}

	if matchedIndex == -1 {
		if minsize == maxsize {
			errs = append(errs, errors.New(fmt.Sprintf("expected %d arguments but %d were supplied", maxsize, len(args))))
			return
		}
		errs = append(errs, errors.New(fmt.Sprintf("expected between %d and %d arguments but %d were supplied", minsize, maxsize, len(args))))
		return
	}

	link := out
	matched := true
	expectedGroup := selectedCommand.Expects[matchedIndex]

	for i := range args {
		var expected ExpressionType

		if expectedGroup[len(expectedGroup)-1] == RepeatLast && i > len(expectedGroup)-2 {
			expected = expectedGroup[len(expectedGroup)-2]
		} else if expectedGroup[len(expectedGroup)-1] == RepeatAll {
			expected = expectedGroup[i%(len(expectedGroup)-1)]
		} else {
			expected = expectedGroup[i]
		}

		switch expected {
		case UserName:
			_, matches := args[i].(userName)
			if !matches {
				matched = false
				errs = append(errs, cmdError("Expected a username (name#number), got %+v", args[i]))
				continue
			}
			link.Next = &Node{
				UserName,
				string(args[i].(userName)),
				nil,
			}
		case UserMention:
			_, matches := args[i].(userReference)
			if !matches {
				matched = false
				errs = append(errs, cmdError("Expected a mention, got %+v", args[i]))
				continue
			}
			link.Next = &Node{
				UserMention,
				string(args[i].(userReference)),
				nil,
			}
		case ChannelName:
			_, matches := args[i].(channelReference)
			if !matches {
				matched = false
				errs = append(errs, cmdError("Expected a channel, got %+v", args[i]))
				continue
			}
			link.Next = &Node{
				ChannelName,
				string(args[i].(channelReference)),
				nil,
			}
		case Literal:
			_, matches := args[i].(string)
			if !matches {
				matched = false
				errs = append(errs, cmdError("Expected plaintext, got %+v", args[i]))
				continue
			}
			fallthrough
		case Any:
			link.Next = &Node{
				ChannelName,
				args[i].(string),
				nil,
			}
		}
	}

	if matched {
		errs = nil
	}
	return
}
