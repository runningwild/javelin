// Package mnemonics contains everything necessary for parsing assembly.  This does the conversion
// into mnemonics, which are the instructions written by the programmer, which then know how to
// convert themselves into opcodes, which are the instructions actually executed.
package mnemonics

import (
	"fmt"

	"github.com/runningwild/javelin/machine"
)

type MnemonicInstruction interface {
	Validate() ([]OpcodeInstruction, error)
}

type OpcodeInstruction interface {
	// Execute(m *Machine)
}

type Instruction struct {
	AddGeneral *AddGeneral `@@ |`
	AddNeon    *AddNeon    `@@`
}

type AddGeneral struct {
	Rd int `"add" @RegisterGeneral`
	Rn int `  "," @RegisterGeneral`
	Rm int `  "," @RegisterGeneral`
}

type AddNeon struct {
	Vd RegisterNeon `"add" @RegisterNeon`
	Vn RegisterNeon `  "," @RegisterNeon`
	Vm RegisterNeon `  "," @RegisterNeon`
}

type RegisterNeon struct {
	V int    `@RegisterNeon`
	T string `@TypeSpecifier`
}

func (i *AddGeneral) Validate() ([]OpcodeInstruction, error) {
	return nil, nil
}

func (i *AddNeon) Validate() ([]OpcodeInstruction, error) {
	if i.Vd.T != i.Vn.T || i.Vd.T != i.Vm.T {
		return nil, fmt.Errorf("type specifiers do not match: (%s, %s, %s))", i.Vd.T, i.Vn.T, i.Vm.T)
	}
	return nil, nil
}

type Machine struct{}

// C6.2.5 ADD (immediate)
type OpcodeAddImmedite struct {
	Sf byte // 1 bit
	// Op = b0
	// S  = b0
	// OpCode = b100010
	Sh  byte   // 1 bit
	Imm uint16 // 12 bits
	Rn  uint16 // 5 bits
	Rd  uint16 // 5 bits
}

func (op *OpcodeAddImmedite) Encode() uint32 {
	return uint32(op.Sf)<<31 |
		0b0<<30 |
		0b0<<29 |
		0b100010<<23 |
		uint32(op.Sh)<<22 |
		uint32(op.Imm)<<10 |
		uint32(op.Rn)<<5 |
		uint32(op.Rd)
}

func (op *OpcodeAddImmedite) Execute(m *machine.Machine) {
	var datamask uint64 = 0xffffffff
	if op.Sf&0x01 == 1 {
		datamask = 0xffffffffffffffff
	}
	imm := uint32(op.Imm)
	if op.Sh&0x01 == 1 {
		imm = imm << 12
	}
	var op1 uint64
	if op.Rn&0b11111 == 31 {
		op1 = m.SP
	} else {
		op1 = m.R[op.Rn&0b11111]
	}
	op1 = op1 & datamask
	var op2 uint64
	op2 = uint64(imm)
	op2 = op2 & datamask
	if op.Rn&0b11111 == 31 {
		m.SP = op1 + op2
	} else {
		m.R[op.Rd&0b11111] = op1 + op2
	}
}

// C6.2.6 ADD (shifted register)
type OpcodeAddShiftedRegister struct {
	Sf byte // 1 bit
	// Op = b0
	// S  = b0
	// OpCode = b01011
	Shift byte // 2 bits
	// _ = b0
	Rm  uint16 // 5 bits
	Imm uint16 // 6 bits
	Rn  uint16 // 5 bits
	Rd  uint16 // 5 bits
}

func (op *OpcodeAddShiftedRegister) Encode() uint32 {
}

func (op *OpcodeAddShiftedRegister) Execute(m *Machine) {
}
