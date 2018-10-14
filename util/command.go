package util

import "strings"

type Command []string

func ExtractCommand(prefix string, body string) (out Command) {
	if strings.HasPrefix(body, prefix) {
		trimmed := strings.Trim(body, prefix)
		out = strings.Split(trimmed, " ")
	}
	return
}