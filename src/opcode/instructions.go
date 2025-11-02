// Package mnemonics contains everything necessary for parsing assembly.  This does the conversion
// into mnemonics, which are the instructions written by the programmer, which then know how to
// convert themselves into opcodes, which are the instructions actually executed.
package opcode

import (
	"github.com/runningwild/javelin/machine"
)

type Instruction interface {
	Encode() uint32
	Execute(m *machine.Machine)
}

// C6.2.5 ADD (immediate)
type AddImmedite struct {
	Sf byte // 1 bit
	// Op = b0
	// S  = b0
	// OpCode = b100010
	Sh  byte   // 1 bit
	Imm uint16 // 12 bits
	Rn  uint16 // 5 bits
	Rd  uint16 // 5 bits
}

func (op *AddImmedite) Encode() uint32 {
	return uint32(op.Sf)<<31 |
		0b0<<30 |
		0b0<<29 |
		0b100010<<23 |
		uint32(op.Sh)<<22 |
		uint32(op.Imm)<<10 |
		uint32(op.Rn)<<5 |
		uint32(op.Rd)
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

func (op *AddShiftedRegister) Encode() uint32 {
	return uint32(op.Sf)<<31 |
		(0b0 << 30) |
		(0b0 << 29) |
		(0b01011 << 24) |
		uint32(op.Shift)<<22 |
		(0b0 << 21) |
		uint32(op.Rm)<<16 |
		uint32(op.Imm)<<10 |
		uint32(op.Rn)<<5 |
		uint32(op.Rd)
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
	Sf byte // 1 bit
	// Op = b0
	// S  = b0
	// OpCode = b01011
	// _ = b1
	Opt byte   // 2 bits
	Imm uint16 // 3 bits
	// _ = b00
	Rm  uint16 // 5 bits
	Rn  uint16 // 5 bits
	Rd  uint16 // 5 bits
}

func (op *AddExtendedRegister) Encode() uint32 {
	return uint32(op.Sf)<<31 |
		(0b0 << 30) |
		(0b0 << 29) |
		(0b01011 << 24) |
		(0b1 << 21) |
		uint32(op.Opt)<<19 |
		uint32(op.Imm)<<16 |
		(0b00 << 14) |
		uint32(op.Rm)<<10 |
		uint32(op.Rn)<<5 |
		uint32(op.Rd)
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
	Q    byte   // 1 bit
	Size byte   // 2 bits
	Rm   uint16 // 5 bits
	Rn   uint16 // 5 bits
	Rd   uint16 // 5 bits
}

func (op *AddVector) Encode() uint32 {
	return (0b0 << 31) |
		uint32(op.Q)<<30 |
		(0b0 << 29) |
		(0b01110 << 24) |
		uint32(op.Size)<<22 |
		(0b1 << 21) |
		uint32(op.Rm)<<16 |
		(0b0010 << 12) |
		(0b10 << 10) |
		uint32(op.Rn)<<5 |
		uint32(op.Rd)
}

func (op *AddVector) Execute(m *machine.Machine) {
	var esize int
	var elements int
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

	lanes := 128 / esize
	if op.Q == 1 {
		lanes = 256 / esize
	}

	for i := 0; i < lanes; i++ {
		op1 := m.V[op.Rn].Get(i, esize)
		op2 := m.V[op.Rm].Get(i, esize)
		result := op1 + op2
		m.V[op.Rd].Set(i, esize, result)
	}
}
