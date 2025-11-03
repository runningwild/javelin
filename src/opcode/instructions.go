// Package mnemonics contains everything necessary for parsing assembly.  This does the conversion
// into mnemonics, which are the instructions written by the programmer, which then know how to
// convert themselves into opcodes, which are the instructions actually executed.
package opcode

import (
	"fmt"

	"github.com/runningwild/javelin/machine"
)

type Instruction interface {
	Encode() uint32
	Execute(m *machine.Machine)
}

// C6.2.5 ADD (immediate)
type AddImmedite struct {
	Sf  uint32 // 1 bit
	Sh  uint32 // 1 bit
	Imm uint32 // 12 bits
	Rn  uint32 // 5 bits
	Rd  uint32 // 5 bits
}

type bits struct {
	v uint32
	b int
}

func buildUint32(sbs ...bits) uint32 {
	var V uint32
	t := 0
	for _, sb := range sbs {
		V = V << sb.b
		V = V | (sb.v & uint32((1<<sb.b)-1))
		t += sb.b
	}
	if t != 32 {
		panic(fmt.Sprintf("tried to construct uint32 with %d bits", t))
	}
	return V
}

func (op *AddImmedite) Encode() uint32 {
	return buildUint32([]bits{
		{op.Sf, 1},
		{0, 1}, // op
		{0, 1}, // S
		{0b100010, 6},
		{op.Sh, 1},
		{op.Imm, 12},
		{op.Rn, 5},
		{op.Rd, 5},
	}...)
}

func (op *AddImmedite) Execute(m *machine.Machine) {
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
type AddShiftedRegister struct {
	Sf    uint32 // 1 bit
	Shift uint32 // 2 bits
	Rm    uint32 // 5 bits
	Imm   uint32 // 6 bits
	Rn    uint32 // 5 bits
	Rd    uint32 // 5 bits
}

func (op *AddShiftedRegister) Encode() uint32 {
	return buildUint32([]bits{
		{op.Sf, 1},
		{0, 1}, // op
		{0, 1}, // S
		{0b01011, 5},
		{op.Shift, 2},
		{0, 1},
		{op.Rm, 5},
		{op.Imm, 6},
		{op.Rn, 5},
		{op.Rd, 5},
	}...)
}

func (op *AddShiftedRegister) Execute(m *machine.Machine) {
	var datasize uint = 32
	if op.Sf&0x01 == 1 {
		datasize = 64
	}

	var op1 uint64
	if op.Rn&0b11111 == 31 {
		op1 = m.SP
	} else {
		op1 = m.R[op.Rn&0b11111]
	}

	var op2 uint64
	op2 = m.R[op.Rm&0b11111]

	shift_type := op.Shift
	shift_amount := uint64(op.Imm)

	var shifted_op2 uint64
	switch shift_type {
	case 0b00: // LSL
		shifted_op2 = op2 << shift_amount
	case 0b01: // LSR
		shifted_op2 = op2 >> shift_amount
	case 0b10: // ASR
		if datasize == 32 {
			shifted_op2 = uint64(int32(op2) >> shift_amount)
		} else {
			shifted_op2 = uint64(int64(op2) >> shift_amount)
		}
	}

	var result uint64
	if datasize == 32 {
		result = uint64(uint32(op1) + uint32(shifted_op2))
	} else {
		result = op1 + shifted_op2
	}

	if op.Rd&0b11111 == 31 {
		m.SP = result
	} else {
		m.R[op.Rd&0b11111] = result
	}
}

// C6.2.4 ADD (extended register)
type AddExtendedRegister struct {
	Sf  uint32 // 1 bit
	Opt byte   // 2 bits
	Imm uint32 // 3 bits
	Rm  uint32 // 5 bits
	Rn  uint32 // 5 bits
	Rd  uint32 // 5 bits
}

func (op *AddExtendedRegister) Encode() uint32 {
}

func (op *AddExtendedRegister) Execute(m *machine.Machine) {
	var datasize uint = 32
	if op.Sf&0x01 == 1 {
		datasize = 64
	}

	var op1 uint64
	if op.Rn&0b11111 == 31 {
		op1 = m.SP
	} else {
		op1 = m.R[op.Rn&0b11111]
	}

	op2 := m.R[op.Rm&0b11111]

	var extended_op2 uint64
	extend_type := op.Opt
	shift := op.Imm

	switch extend_type {
	case 0b000: // UXTB
		extended_op2 = uint64(uint8(op2))
	case 0b001: // UXTH
		extended_op2 = uint64(uint16(op2))
	case 0b010: // UXTW
		extended_op2 = uint64(uint32(op2))
	case 0b011: // UXTX
		extended_op2 = op2
	case 0b100: // SXTB
		extended_op2 = uint64(int8(op2))
	case 0b101: // SXTH
		extended_op2 = uint64(int16(op2))
	case 0b110: // SXTW
		extended_op2 = uint64(int32(op2))
	case 0b111: // SXTX
		extended_op2 = op2
	}

	shifted_op2 := extended_op2 << shift

	var result uint64
	if datasize == 32 {
		result = uint64(uint32(op1) + uint32(shifted_op2))
	} else {
		result = op1 + shifted_op2
	}

	if op.Rd&0b11111 == 31 {
		m.SP = result
	} else {
		m.R[op.Rd&0b11111] = result
	}
}

// ADD (vector)
type AddVector struct {
	Q    uint32 // 1 bit
	Size uint32 // 2 bits
	Rm   uint32 // 5 bits
	Rn   uint32 // 5 bits
	Rd   uint32 // 5 bits
}

func (op *AddVector) Encode() uint32 {
}

func (op *AddVector) Execute(m *machine.Machine) {
	var esize int
	switch op.Size {
	case 0b00:
		esize = 8
	case 0b01:
		esize = 16
	case 0b10:
		esize = 32
	case 0b11:
		esize = 64
	}

	datasize := 64
	if op.Q == 1 {
		datasize = 128
	}
	lanes := datasize / esize

	for i := 0; i < lanes; i++ {
		op1 := m.V[op.Rn].Get(i, esize)
		op2 := m.V[op.Rm].Get(i, esize)
		result := op1 + op2
		m.V[op.Rd].Set(i, esize, result)
	}
}
