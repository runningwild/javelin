package main

import (
	"testing"
)

func TestADD(t *testing.T) {
	m := NewMachine()
	m.R[1] = 5
	// m.R[2] = 3

	// // ADD R3, R1, R2 (R3 = R1 + R2)
	// (&ADDInstruction{
	// 	Rd: 3,
	// 	Rn: 1,
	// 	Rm: 2,
	// }).Execute(m)

	// if m.R[3] != 8 {
	// 	t.Errorf("ADD failed: expected R3 to be 8, got %d", m.R[3])
	// }
}
