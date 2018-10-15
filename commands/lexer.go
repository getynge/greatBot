package commands

import (
	"strconv"
	"strings"
)

func LexCommand(prefix string, body string) (results []interface{}) {

	var arguments []string

	if strings.HasPrefix(body, prefix) {
		trimmed := strings.Trim(body, prefix)
		arguments = strings.Split(trimmed, " ")
	} else {
		return
	}

	results = make([]interface{}, len(arguments))
	results[0] = arguments[0]

	for i, argument := range arguments[1:] {
		if strings.HasPrefix(argument, "<@") && strings.HasSuffix(argument, ">") && len(argument) == 22 {
			results[i+1] = userReference(argument)
		} else if strings.HasPrefix(argument, "#") {
			results[i+1] = channelReference(argument)
		} else if strings.Contains(argument, "#") {
			parts := strings.Split(argument, "#")

			// should probably be using regex for this
			if len(parts) != 2 || len(parts[1]) != 4 {
				results[i+1] = argument
				continue
			}
			_, err := strconv.Atoi(parts[1])

			if err != nil {
				results[i+1] = argument
				continue
			}
			results[i+1] = userName(argument)
		} else {
			results[i+1] = argument
		}
	}

	return
}
