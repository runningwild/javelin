package main

import (
	"fmt"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

// InstructionType defines the type of ARM instruction
type InstructionType int

const (
	INST_LDR InstructionType = iota
	INST_STR
	INST_ADD
	INST_VLD1
	INST_VST1
	INST_VADD
	INST_UNKNOWN
)

// MnemonicInstruction represents a parsed ARM instruction
type MnemonicInstruction interface {
	Validate() ([]OpcodeInstruction, error)
}

type OpcodeInstruction interface {
	Execute(m *Machine)
}

// AST structures for participle

type AsmProgram struct {
	Instructions []MnemonicInstruction
}

type asmInstructions struct {
	Instructions []AsmInstruction `@@*`
}

type AsmInstruction struct {
	ADD *ADDInstruction `@@`
}

type ADDInstruction struct {
	R []int          `"add" ((@RegisterGeneral "," @RegisterGeneral "," @RegisterGeneral) |`
	V []RegisterNeon `       (@@ "," @@ "," @@))`
}

type RegisterNeon struct {
	V int    `@RegisterNeon`
	T string `@TypeSpecifier`
}

func (i *ADDInstruction) Execute(m *Machine) {
	if len(i.R) > 0 {
		m.R[i.R[0]] = m.R[i.R[1]] + m.R[i.R[2]]
	} else {
		// do neon
	}
}

func (i *ADDInstruction) Validate() error {
	if (i.R != nil) == (i.V != nil) {
		return fmt.Errorf("exactly one of general or neon registers should be specified")
	}
	if i.R != nil {
		if got, want := len(i.R), 3; got != want {
			return fmt.Errorf("general registers specified but got %d instead of the expected %d", got, want)
		}
	} else {
		if got, want := len(i.V), 3; got != want {
			return fmt.Errorf("neon registers specified but got %d instead of the expected %d", got, want)
		}
		if i.V[0].T != i.V[1].T || i.V[0].T != i.V[2].T {
			return fmt.Errorf("type specifiers do not match: (%s, %s, %s))", i.V[0].T, i.V[1].T, i.V[2].T)
		}
	}
	return nil
}

var ( // global parser variables
	asmLexer = lexer.MustSimple([]lexer.SimpleRule{
		{"whitespace", `\s+`},
		{"Mnemonic", `add`},
		{"RegisterGeneral", `r([12]?[0-9]|30|31)\b`},
		{"RegisterNeon", `v([12]?[0-9]|30|31)\b`},
		{"TypeSpecifier", `[.](16b)`},
		{"Punct", `,`},
	},
	)
	parser = participle.MustBuild[asmInstructions](
		participle.Lexer(asmLexer),
		participle.Elide("whitespace"),
		participle.Map(func(t lexer.Token) (lexer.Token, error) {
			t.Value = t.Value[1:]
			return t, nil
		}, "RegisterGeneral", "RegisterNeon", "TypeSpecifier"),
	)
)

// ParseProgram parses a single line of ARM assembly into an Instruction struct
func ParseProgram(line string) (*AsmProgram, error) {
	program, err := parser.ParseString("", line)
	if err != nil {
		return nil, fmt.Errorf("parsing error: %w", err)
	}
	var p AsmProgram
	for _, i := range program.Instructions {
		switch {
		case i.ADD != nil:
			p.Instructions = append(p.Instructions, i.ADD)
		}
	}
	fmt.Printf("root: %v\n", program)
	for _, i := range p.Instructions {
		if err := i.Validate(); err != nil {
			return nil, err
		}
	}
	return &p, nil
}
