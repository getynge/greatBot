package commands

import (
	"fmt"
	"github.com/alecthomas/participle"
	"os"
)

type Argument struct {
	Channel      *int64 `'#'@Int`
	Mention      *int64 `| '<''@'@Int'>'`
	Username     string `| @Ident`
	Discriminant *int64 `@Int`
}

type Command struct {
	Name string      `@Ident`
	Args []*Argument `[ { @@ } ]`
}

type Commands struct {
	Command  *Command   `@@`
	Commands []*Command `| { @@';' }`
}

func ParseCommand(text string) (parsed *Commands) {
	parser, err := participle.Build(&Commands{})
	if err != nil {
		fmt.Println("Parser syntax malformed, exiting...")
		fmt.Printf("Error %s\n", err.Error())
		os.Exit(1)
	}

	parsed = &Commands{}
	err = parser.ParseString(text, parsed)
	if err != nil {
		parsed = nil
		fmt.Printf("Bad command syntax %s\n", err.Error())
	}

	return
}
