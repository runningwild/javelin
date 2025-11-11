package parser

import (
	"testing"

	"github.com/alecthomas/participle/v2"
)

func TestParseParameter(t *testing.T) {
	p := participle.MustBuild[Parameter](
		participle.Lexer(asmDef),
	)
	for _, reg := range []string{
		"<Xm>",
		"<Wn>",
		"<Wbeans>",
		"<R><n>",
		"<X(s+1)>",
	} {
		_, err := p.ParseString("", reg)
		if err != nil {
			t.Errorf("failed to parse %q: %v", reg, err)
		}
	}
}

func TestParseExpression(t *testing.T) {
	p := participle.MustBuild[Expression](
		participle.Lexer(asmDef),
	)
	for _, reg := range []string{
		"s+1",
	} {
		_, err := p.ParseString("", reg)
		if err != nil {
			t.Errorf("failed to parse %q: %v", reg, err)
		}
	}
}

func TestParseAddress(t *testing.T) {
	p := participle.MustBuild[Address](
		participle.Lexer(asmDef),
	)
	for _, reg := range []string{
		"[<Xn>]",
		"[<Xn|SP>]",
		"[<Xn>, #<imm>]",
		"[<Xn|SP>, #<imm>]",
	} {
		_, err := p.ParseString("", reg)
		if err != nil {
			t.Errorf("failed to parse %q: %v", reg, err)
		}
	}
}

func TestParseOptionalOperand(t *testing.T) {
	p := participle.MustBuild[OptionalOperand](
		participle.Lexer(asmDef),
	)
	for _, reg := range []string{
		"<extend>",
		"<extend> {#<amount>}",
	} {
		_, err := p.ParseString("", reg)
		if err != nil {
			t.Errorf("failed to parse %q: %v", reg, err)
		}
	}
}
