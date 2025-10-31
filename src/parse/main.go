package main

import (
	"fmt"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

var (
	customLexer = lexer.Must(lexer.NewSimple([]lexer.SimpleRule{
		{Name: "Foobar", Pattern: `foo|bar`},
		{Name: "whitespace", Pattern: `\s+`},
	}))
)

type Program struct {
	Tokens []string `@Foobar*`
}

type Token struct {
	Foobar string `@Foobar`
}

func main() {
	parser := participle.MustBuild[Program](
		participle.Lexer(customLexer),
	)
	program, err := parser.ParseString("", "foo bar foo")
	if err != nil {
		panic(err)
	}
	for _, token := range program.Tokens {
		fmt.Println(token.Foobar)
	}
}
