package main

import (
	"fmt"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type Register struct {
	Number int `@Register`
}

func main() {
	lex := lexer.MustSimple([]lexer.SimpleRule{
		{Name: "Register", Pattern: `r\d+`},
		{Name: "whitespace", Pattern: `\s+`},
	})
	mapper := participle.Map(func(t lexer.Token) (lexer.Token, error) {
		t.Value = t.Value[1:]
		return t, nil
	}, "Register")
	parser := participle.MustBuild[Register](participle.Lexer(lex), mapper)

	reg, err := parser.ParseString("", "r5")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Parsed register number: %d\n", reg.Number)

	reg, err = parser.ParseString("", "r16")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Parsed register number: %d\n", reg.Number)
}
