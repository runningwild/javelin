package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/runningwild/javelin/parser"
)

//go:embed instruction.tmpl
var instTmpl string
var t = template.Must(template.New("instructionTmpl").Parse(instTmpl))

//go:embed base.tmpl
var baseTmpl string
var b = template.Must(template.New("baseTmpl").Parse(instTmpl))

var (
	input  = flag.String("input", "", "Input file which should be text version of instruction list")
	outdir = flag.String("output-dir", "instructions", "output directory to write generated files to")
)

type txp struct {
	Section     string
	Descriptor  string
	Instruction *parser.Instruction
}

func (txp) Upperize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[0:1]) + s[1:]
}

func main() {
	flag.Parse()
	if err := doit(context.Background(), *input, *outdir); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}

func doit(ctx context.Context, input, outdir string) error {
	if err := os.MkdirAll(outdir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory %q: %w", outdir, err)
	}
	f, err := os.Open(input)
	if err != nil {
		return fmt.Errorf("failed to open input file %q: %w", input, err)
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("failed to read input file %q: %w", input, err)
	}

	instHeadingRE := regexp.MustCompile(`C6[.]2[.](\d+)\s+(([A-Z][A-Z0-9]*(, )?)*)+\s*(\S.*\S)?\s*$`)
	var section string
	var opcodes []string
	var descriptor string
	var instDefRE *regexp.Regexp

	N := 20
	for _, line := range strings.Split(string(data), "\n") {
		m := instHeadingRE.FindStringSubmatch(line)
		if m != nil {
			section = m[1]
			opcodes = strings.Split(m[2], ", ")
			descriptor = m[5]
			fmt.Printf("%s %v %s\n", section, opcodes, descriptor)
			var err error
			instDefRE, err = regexp.Compile(fmt.Sprintf(`^\s*(%s)\s+(.*)\s*$`, strings.Join(opcodes, "|")))
			if err != nil {
				return fmt.Errorf("failed to compile regexp for %v: %w", opcodes, err)
			}
			continue
		}
		if instDefRE == nil {
			continue
		}
		m = instDefRE.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		fmt.Printf("%s %s\n", m[1], m[2])
		inst, err := parser.ParseInstruction(line)
		if err != nil {
			err = fmt.Errorf("failed to parse %q: %w", line, err)
			fmt.Printf("%v\n", err)
			continue
		}
		fmt.Printf("%v\n", inst)
		outpath := filepath.Join(outdir, fmt.Sprintf("%s_%s.go", strings.ToLower(inst.Mnemonic), section))
		outfile, err := os.Create(outpath)
		if err != nil {
			return fmt.Errorf("failed to create output file %q: %w", outpath, err)
		}
		defer outfile.Close()
		target := txp{
			Instruction: inst,
			Descriptor:  descriptor,
			Section:     section,
		}
		if err := t.Execute(outfile, target); err != nil {
			return fmt.Errorf("failed to execute template on %q: %w", inst.Mnemonic, err)
		}
		N--
		if N < 0 {
			break
		}
	}

	os.WriteFile(filepath.Join(outdir, "base.go"), []byte(baseTmpl), 0644)
	if err != nil {
		return fmt.Errorf("failed to create base.go in output dir: %w", err)
	}

	return nil
}
