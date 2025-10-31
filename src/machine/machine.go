package main

import (
	"fmt"
)

// Machine represents the state of the ARMv8-A machine.
type Machine struct {
	// General-purpose registers x0-x30.
	R [31]uint64
	// Program Counter.
	PC uint64
	// Stack Pointer.
	SP uint64
	// Current Program Status Register.
	CPSR uint32
	// Memory. A simple byte slice for simulation.
	Memory []byte
}

// NewMachine creates a new Machine with initialized memory.
func NewMachine(memorySize int) *Machine {
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
