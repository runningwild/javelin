package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/runningwild/javelin/parser"
)

func main() {
	if err := doit(); err != nil {
		fmt.Printf("%v\n", err)
		return
	}
}

func doit() error {
	f, err := os.Open("../../docs/encodings.txt")
	if err != nil {
		return err
	}
	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	fmt.Printf("Fails:\n")
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		_, err := parser.Parser.Parse("", strings.NewReader(line))
		if err != nil {
			fmt.Printf("%s: %v\n", line, err)
		}
	}
	fmt.Printf("Successes:\n")
	for _, line := range lines {
		out, err := parser.Parser.Parse("", strings.NewReader(line))
		if err == nil {
			fmt.Printf("%v\n", out)
		}
	}
	return nil
}
