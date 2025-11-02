package parser

import (
	"testing"
)

func TestLexer(t *testing.T) {
	input := "add x1, x2, #123\n"
	_, tokens := lex(input)

	expectedTokens := []token{
		{tokenIdentifier, "add"},
		{tokenRegister, "x1"},
		{tokenComma, ","},
		{tokenRegister, "x2"},
		{tokenComma, ","},
		{tokenNumber, "#123"},
		{tokenNewline, "\n"},
		{tokenEOF, ""},
	}

	for i, expected := range expectedTokens {
		got := <-tokens
		if got.typ != expected.typ || got.val != expected.val {
			t.Errorf("Token %d: expected %v, got %v", i, expected, got)
		}
	}
}
