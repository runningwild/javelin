package main

import (
	"fmt"
	"os"

	"github.com/runningwild/javelin/machine"
	"github.com/runningwild/javelin/parser"
)

func main() {
	addAsm := `
add x2, x3, x5
add w2, w3, #5
`
	p := parser.New(addAsm)
	insts, err := p.Parse()
	if err != nil {
		fmt.Printf("Failed to parse ADD: %v", err)
		os.Exit(1)
	}

	m := machine.New(1024 * 1024)
	m.R[3] = 10
	m.R[5] = 20

	for _, inst := range insts {
		inst.Execute(m)
	}

	fmt.Printf("R[2]: %d\n", m.R[2])
}
