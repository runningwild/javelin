package instruction_test

import (
	"testing"

	instruction "github.com/runningwild/javelin/parser/instructions"
)

func TestAdd_6(t *testing.T) {
	if _, err := instruction.Parse("ADD X1, X2, #4"); err != nil {
		t.Errorf("%v", err)
	}
	if _, err := instruction.Parse("ADD X1, X2, X3"); err != nil {
		t.Errorf("%v", err)
	}
	if _, err := instruction.Parse("ADD X1, X2, X3"); err != nil {
		t.Errorf("%v", err)
	}
}
