package parser

import (
	"reflect"
	"testing"

	"github.com/runningwild/javelin/opcode"
)

func TestParser(t *testing.T) {
	input := "add x1, x2, #123\n"
	p := New(input)
	insts, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse() failed: %v", err)
	}
	if len(insts) != 1 {
		t.Fatalf("Expected 1 instruction, got %d", len(insts))
	}

	expected := &opcode.AddImmedite{
		Sf:  1,
		Sh:  0,
		Rd:  1,
		Rn:  2,
		Imm: 123,
	}

	if !reflect.DeepEqual(insts[0], expected) {
		t.Errorf("Expected instruction %v, got %v", expected, insts[0])
	}
}

func TestParserAddShiftedRegister(t *testing.T) {
	input := "add x1, x2, x3, lsl #4\n"
	p := New(input)
	insts, err := p.Parse()
	if err != nil {
		t.Fatalf("Parse() failed: %v", err)
	}
	if len(insts) != 1 {
		t.Fatalf("Expected 1 instruction, got %d", len(insts))
	}

	expected := &opcode.AddShiftedRegister{
		Sf:    1,
		Shift: 0,
		Rd:    1,
		Rn:    2,
		Rm:    3,
		Imm:   4,
	}

	if !reflect.DeepEqual(insts[0], expected) {
		t.Errorf("Expected instruction %v, got %v", expected, insts[0])
	}
}

