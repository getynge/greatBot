package commands

import (
	"errors"
	"fmt"
	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
)

type Argument struct {
	Channel   string `@Channel`
	Mention   string `| @Mention`
	Username  string `| @Username`
	Quote     string `| @String`
	PlainText string `| @Identifier | @Any`
}

type Command struct {
	Name string      `@Identifier`
	Args []*Argument `[ { @@ } ]`
}

type Commands struct {
	Command  *Command   `@@`
	Commands []*Command `| { @@';' }`
}

type Parser struct {
	parser *participle.Parser
}

func InitParser() (parser Parser, err error) {
	parser = Parser{}

	// The regexp lexer stores input in memory, but this should be fine since chat messages are relatively small
	cmdLexer, err := lexer.Regexp(
		`(?m)` +
			`(\s+)` +
			`|(^//.*$)` +
			`|(/\*.*\*/)` + // discord bot supports comments, because why not
			`|(?P<String>"(?:\\.|[^"])*")` +
			`|(?P<Mention><@[0-9]{18}>)` +
			`|(?P<Channel>#[a-zA-Z]+)` +
			`|(?P<Username>.*#[0-9]{4})` +
			`|(?P<Identifier>[a-zA-Z]+)` +
			`|(?P<Any>\S+)`, // used for nicknames, which may not conform to the identifier group
	)

	if err != nil {
		return
	}
	fmt.Println(cmdLexer.Symbols())
	parser.parser, err = participle.Build(&Commands{}, participle.Lexer(cmdLexer))

	return
}

func (parserHolder *Parser) ParseCommand(text string) (parsed *Commands, err error) {
	parser := parserHolder.parser

	parsed = &Commands{}
	err = parser.ParseString(text, parsed)
	if err != nil {
		parsed = nil
		errtxt := fmt.Sprintf("Bad command syntax %s\n", err.Error())
		err = errors.New(errtxt)
	}

	return
}
