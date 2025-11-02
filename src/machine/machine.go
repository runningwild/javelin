package machine

import (
	"encoding/binary"
	"fmt"
)

// VectorRegister represents a 128-bit vector register.
type VectorRegister [16]byte

// Get returns the value of a lane in the vector register.
func (v *VectorRegister) Get(lane, esize int) uint64 {
	offset := lane * (esize / 8)
	switch esize {
	case 8:
		return uint64(v[offset])
	case 16:
		return uint64(binary.LittleEndian.Uint16(v[offset:]))
	case 32:
		return uint64(binary.LittleEndian.Uint32(v[offset:]))
	case 64:
		return binary.LittleEndian.Uint64(v[offset:])
	}
	return 0
}

// Set sets the value of a lane in the vector register.
func (v *VectorRegister) Set(lane, esize int, val uint64) {
	offset := lane * (esize / 8)
	switch esize {
	case 8:
		v[offset] = byte(val)
	case 16:
		binary.LittleEndian.PutUint16(v[offset:], uint16(val))
	case 32:
		binary.LittleEndian.PutUint32(v[offset:], uint32(val))
	case 64:
		binary.LittleEndian.PutUint64(v[offset:], val)
	}
}

// Machine represents the state of the ARMv8-A machine.
type Machine struct {
	// General-purpose registers x0-x30.
	R [31]uint64
	// Vector registers v0-v31.
	V [32]VectorRegister
	// Program Counter.
	PC uint64
	// Stack Pointer.
	SP uint64
	// Current Program Status Register.
	CPSR uint32
	// Memory. A simple byte slice for simulation.
	Memory []byte
}

// New creates a new Machine with initialized memory.
func New(memorySize int) *Machine {
	return &Machine{
		Memory: make([]byte, memorySize),
	}
}

// PrintState prints the current state of the machine's registers and PC.
func (m *Machine) PrintState() {
	fmt.Println("Registers:")
	for i := 0; i < 31; i++ {
		fmt.Printf("  x%d: 0x%x\n", i, m.R[i])
	}
	fmt.Printf("  PC: 0x%x\n", m.PC)
	fmt.Printf("  SP: 0x%x\n", m.SP)
	fmt.Printf("CPSR: 0x%x\n", m.CPSR)
}
