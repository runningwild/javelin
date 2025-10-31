package main

import (
	"fmt"
	"strconv"
	"strings"

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
	Add *Add `@@`
}

type Add struct {
	AddImmediate        *AddImmediate        `"add" (@@ |`
	AddShiftedRegister  *AddShiftedRegister  `       @@ |`
	AddExtendedRegister *AddExtendedRegister `       @@ |`
	AddVector           *AddVector           `       @@ )`
}

type AddImmediate struct {
	Rd  int    `@RegisterGeneral ","`
	Rn  int    `@RegisterGeneral ","`
	Imm string `"#" @Integer`
}

func (i *AddImmediate) Validate() ([]OpcodeInstruction, error) {
	immStr := i.Imm
	base := 10
	if strings.HasPrefix(immStr, "0x") {
		immStr = immStr[2:]
		base = 16
	}
	imm, err := strconv.ParseInt(immStr, base, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse immediate: %w", err)
	}
	if imm&(^0xfff) != 0 {
		return nil, fmt.Errorf("immediate %d overflows 12 bits", imm)
	}
	return nil, nil
}

type AddShiftedRegister struct {
	Rd  int     `@RegisterGeneral  ","`
	Rn  int     `@RegisterGeneral  ","`
	Rm  int     `@RegisterGeneral (","`
	Dir *string `  @Shift`
	Amt *int    `  @Integer)?`
}

type AddExtendedRegister struct {
	Rd     int     `@RegisterGeneral  ","`
	Rn     int     `@RegisterGeneral  ","`
	Rm     int     `@RegisterGeneral (","`
	Extend *string `  @(Extend|"lsl")`
	Amt    *int    `  @Integer?)?`
}

type AddVector struct {
	Vd RegisterNeon `@RegisterNeon ","`
	Vn RegisterNeon `@RegisterNeon ","`
	Vm RegisterNeon `@RegisterNeon`
}

type RegisterNeon struct {
	V int    `@RegisterNeon`
	T string `@TypeSpecifier`
}

var ( // global parser variables
	asmLexer = lexer.MustSimple([]lexer.SimpleRule{
		{"whitespace", `\s+`},
		{"Mnemonic", `add`},
		{"RegisterGeneral", `r([12]?[0-9]|30|31)\b`},
		{"RegisterNeon", `v([12]?[0-9]|30|31)\b`},
		{"TypeSpecifier", `[.](16b)`},
		{"Integer", `[0-9]+|0x[a-fA-F0-9]`},
		{"Shift", `lsl|lsr|asr`},
		{"Extend", `(ux|sx)(t[bhwx])`},
		{"Punct", `[,#]`},
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
		case i.Add != nil:
			switch {
			case i.Add.AddImmediate != nil:
				p.Instructions = append(p.Instructions, i.Add.AddImmediate)
			case i.Add.AddExtendedRegister != nil:
				p.Instructions = append(p.Instructions, i.Add.AddExtendedRegister)
			case i.Add.AddShiftedRegister != nil:
				p.Instructions = append(p.Instructions, i.Add.AddShiftedRegister)
			case i.Add.AddVector != nil:
				p.Instructions = append(p.Instructions, i.Add.AddVector)
			}
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
