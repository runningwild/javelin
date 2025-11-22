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
var instData string
var instTmpl = template.Must(template.New("instructionTmpl").Parse(instData))

//go:embed stubs.tmpl
var stubsData string
var stubsTmpl = template.Must(template.New("stubsTmpl").Parse(stubsData))

//go:embed base.tmpl
var baseData string
var baseTmpl = template.Must(template.New("baseTmpl").Parse(baseData))

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
		instpath := filepath.Join(outdir, fmt.Sprintf("%s_%s.go", strings.ToLower(inst.Mnemonic), section))
		instfile, err := os.Create(instpath)
		if err != nil {
			return fmt.Errorf("failed to create output instruction file %q: %w", instpath, err)
		}
		defer instfile.Close()
		stubpath := filepath.Join(outdir, fmt.Sprintf("%s_%s_stubs.go", strings.ToLower(inst.Mnemonic), section))
		stubfile, err := os.Create(stubpath)
		if err != nil {
			return fmt.Errorf("failed to create output stub file %q: %w", stubpath, err)
		}
		defer stubfile.Close()
		target := txp{
			Instruction: inst,
			Descriptor:  descriptor,
			Section:     section,
		}
		if err := instTmpl.Execute(instfile, target); err != nil {
			return fmt.Errorf("failed to execute instruction template on %q: %w", inst.Mnemonic, err)
		}
		if err := stubsTmpl.Execute(stubfile, target); err != nil {
			return fmt.Errorf("failed to execute stubs template on %q: %w", inst.Mnemonic, err)
		}
		N--
		if N < 0 {
			break
		}
	}

	os.WriteFile(filepath.Join(outdir, "base.go"), []byte(baseData), 0644)
	if err != nil {
		return fmt.Errorf("failed to create base.go in output dir: %w", err)
	}

	return nil
}
