package instruction

import (
	"fmt"
	"regexp"

	"github.com/alecthomas/participle/v2/lexer"
	"github.com/runningwild/javelin/machine"
)

type Instruction interface {
	Encode() uint32
	Execute(*machine.Machine)
}

var parsers map[string][]Parser = map[string][]Parser{}

// A return of nil, nil means that the mnemonic doesn't match.
type Parser func(string) (Instruction, error)

var asmDef = lexer.MustSimple([]lexer.SimpleRule{
	{"Register", `W|X|R`},
	{"HashTag", `#`},
	{"Comment", `//[^\n]*$`},
	{"Integer", `[-+]?\d+`},
	{"Ident", `[a-zA-Z][a-zA-Z0-9_]*`},
	{"Comma", `,`},
	{"whitespace", `[ \t]+`},
})

var instructionRE = regexp.MustCompile(`^\s*(\S+)[. \t]?.*$`)

func Parse(s string) (Instruction, error) {
	m := instructionRE.FindStringSubmatch(s)
	if m == nil {
		return nil, fmt.Errorf("failed to parse mnemonic")
	}
	ps := parsers[m[1]]
	if len(ps) == 0 {
		return nil, fmt.Errorf("no parsers for mnemonic %q", m[1])
	}
	for _, p := range ps {
		inst, err := p(s)
		if inst != nil && err == nil {
			return inst, nil
		}
		if err != nil {
			return nil, err
		}
	}
	return nil, fmt.Errorf("all parsers failed")
}
