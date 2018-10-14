package parser

import (
	"strconv"
	"strings"
)

type Command []string

func LexCommand(prefix string, body string) (out Command) {

	var arguments []string

	if strings.HasPrefix(body, prefix) {
		trimmed := strings.Trim(body, prefix)
		arguments = strings.Split(trimmed, " ")
	} else {
		return
	}
	// TODO: change this to tokens once I finish adding parser support
	out = arguments

	results := make([]interface{}, len(arguments))
	results[0] = arguments[0]

	for i, argument := range arguments[1:] {
		if strings.HasPrefix(argument, "<@") && strings.HasSuffix(argument, ">") && len(argument) == 22 {
			results[i] = userReference{argument}
		} else if strings.HasPrefix(argument, "#") {
			results[i] = channelReference{argument}
		} else if strings.Contains(argument, "#") {
			parts := strings.Split(argument, "#")
			if len(parts) != 2 || len(parts[1]) != 4 {
				results[i] = argument
				continue
			}
			parsed, err := strconv.Atoi(parts[1])

			if err != nil {
				results[i] = argument
				continue
			}
			results[i] = userName{parts[0], parsed}
		} else {
			results[i] = argument
		}
	}

	return
}
