package main

import (
	"testing"
)

func TestParseInstruction(t *testing.T) {
	// Test ADD
	addAsm := "add r7, r2, r4"
	addInst, err := ParseInstruction(addAsm)
	if err != nil {
		t.Fatalf("Failed to parse ADD: %v", err)
	}
	if expected, ok := addInst.(*ADDInstruction); !ok {
		t.Errorf("ADD parsed as a %T, want %T", addInst, expected)
	}

	// Test VADD.I32
	vaddAsm := "VADD.I32 D4, D0, D2"
	vaddInst, err := ParseInstruction(vaddAsm)
	if err != nil {
		t.Fatalf("Failed to parse VADD.I32: %v", err)
	}
	if expected, ok := addInst.(*ADDInstruction); !ok {
		t.Errorf("VADD parsed as a %T, want %T", vaddInst, expected)
	}

	// Test invalid instruction
	invalidAsm := "MOV R0, #1"
	_, err = ParseInstruction(invalidAsm)
	if err == nil {
		t.Errorf("Expected error for invalid instruction, got nil")
	}

	// Test invalid register
	invalidRegAsm := "ADD R33, R1, R2"
	_, err = ParseInstruction(invalidRegAsm)
	if err == nil {
		t.Errorf("Expected error for invalid register, got nil")
	}
}
