package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/runningwild/javelin/machine"
	"github.com/runningwild/javelin/opcode"
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
	Validate() ([]opcode.Instruction, error)
}

type OpcodeInstruction interface {
	Execute(m *machine.Machine)
}

// AST structures for participle

type AsmProgram struct {
	Instructions []MnemonicInstruction
}

type asmInstructions struct {
	Instructions []AsmInstruction `(@@)*`
}

type AsmInstruction struct {
	Add *Add `@@`
}

type Add struct {
	General *AddGeneralSuffix `"add" ( @@ |`
	Vector  *AddVectorSuffix  `        @@ )`
}

type AddGeneralSuffix struct {
	Rd        int                 `@RegisterGeneral ","`
	Rn        int                 `@RegisterGeneral ","`
	Immediate *AddImmediateSuffix `( @@ |`
	Shifted   *AddShiftedSuffix   `  @@ |`
	Extended  *AddExtendedSuffix  `  @@ )`
}

// NEXT: Shifted and Extended cannot share a prefix so we need to split something out
type AddImmediateSuffix struct {
	Imm string `"#" @Integer`
}

type AddShiftedOrExtendedSuffix struct {
	Rm       int                `@RegisterGeneral`
	Shifted  *AddShiftedSuffix  `( @@ |`
	Extended *AddExtendedSuffix `  @@ )`
}

type AddShiftedSuffix struct {
	Dir      *string `(@Shift`
	ShiftAmt *int    `"#" @Integer)?`
}

type AddExtendedSuffix struct {
	Extend   *string `(@(Extend|"lsl")`
	ShiftAmt *int    `("#" @Integer)?)?`
}

type AddVectorSuffix struct {
	Vd RegisterNeon `@@ ","`
	Vn RegisterNeon `@@ ","`
	Vm RegisterNeon `@@`
}

func (a *Add) Validate() ([]opcode.Instruction, error) {
	switch {
	case a.General != nil:
	case a.Vector != nil:
		// TODO return opcode.AddVector
	}
	return nil, nil
}
func parseImmediate(immstr string) (int64, error) {
	base := 10
	if strings.HasPrefix(immstr, "0x") {
		immstr = immstr[2:]
		base = 16
	}
	return strconv.ParseInt(immstr, base, 64)
}

type Addx struct {
	AddShiftedRegister  *AddShiftedRegister  `"add" (@@ |`
	AddImmediate        *AddImmediate        `       @@ |`
	AddExtendedRegister *AddExtendedRegister `       @@ |`
	AddVector           *AddVector           `       @@ )`
}

type AddImmediate struct {
	Rd  int    `@RegisterGeneral ","`
	Rn  int    `@RegisterGeneral ","`
	Imm string `"#" @Integer`
}

func (i *AddImmediate) Validate() ([]opcode.Instruction, error) {
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
	var op opcode.AddImmedite
	op.Rd = uint32(i.Rd)
	op.Rn = uint32(i.Rn)

	if imm&0xfff == imm {
		op.Imm = uint32(imm)
		op.Sh = 0
		return []opcode.Instruction{&op}, nil
	}
	if (imm>>12)&0xfff == imm {
		op.Imm = uint32(imm >> 12)
		op.Sh = 1
		return []opcode.Instruction{&op}, nil
	}

	return nil, fmt.Errorf("immediate %d cannot be represented as a 12-bit value with an optional 12-bit left shift", imm)
}

type AddShiftedRegister struct {
	Rd  int     `@RegisterGeneral  ","`
	Rn  int     `@RegisterGeneral  ","`
	Rm  int     `@RegisterGeneral (","`
	Dir *string `  @Shift`
	Amt *int    `  @Integer)?`
}

func (i *AddShiftedRegister) Validate() ([]opcode.Instruction, error) {
	return nil, nil
}

type AddExtendedRegister struct {
	Rd     int     `@RegisterGeneral  ","`
	Rn     int     `@RegisterGeneral  ","`
	Rm     int     `@RegisterGeneral (","`
	Extend *string `  @(Extend|"lsl")`
	Amt    *int    `  @Integer?)?`
}

func (i *AddExtendedRegister) Validate() ([]opcode.Instruction, error) {
	return nil, nil
}

type AddVector struct {
	Vd RegisterNeon `@@ ","`
	Vn RegisterNeon `@@ ","`
	Vm RegisterNeon `@@`
}

func (i *AddVector) Validate() ([]opcode.Instruction, error) {
	return nil, nil
}

type RegisterNeon struct {
	V int    `@RegisterNeon`
	T string `@TypeSpecifier`
}
