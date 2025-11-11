package parser

import (
	"fmt"
	"iter"
	"reflect"
	"regexp"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/runningwild/javelin/opcode"
)

type Instruction struct {
	Mnemonic        string           `@ArchIdent`
	Operands        []Operand        `( @@ ("," @@)*   )?`
	OptionalOperand *OptionalOperand `( "{" "," @@ "}" )?`
}

func (inst Instruction) Make() {
	fmt.Printf("%s ", inst.Mnemonic)
	for i, op := range inst.Operands {
		if i > 0 {
			fmt.Printf(", ")
		}
		fmt.Printf("%s", op.Make())
	}
	fmt.Printf("\n")
}

// func (inst Instruction) ParseableStruct() string {
// ADD <Wd>, <Wn>, <Wm>{, <shift> #<amount>}
/*

	type AddInstruction struct {
		Wd Register                `"ADD" "<" "W" @RegisterNumber`
		Wn Register                `  "," "<" "W" @RegisterNumber`
		Wm Register                `  "," "<" "W" @RegisterNumber`
		Opt *AddOptionalParameter  ` ("," @@)?`
	}

	type AddOptionalParameter struct {
		Shift string  `@Shift`
		Amount uint32 `@Integer`
	}

*/

// ADD <Wd|WSP>, <Wn|WSP>, #<imm>{, <shift>}
// ADD <Wd|WSP>, <Wn|WSP>, <Wm>{, <extend> {#<amount>}}
// ADD <Xd|SP>, <Xn|SP>, #<imm>{, <shift>}
// ADD <Xd|SP>, <Xn|SP>, <R><m>{, <extend> {#<amount>}}
// ADD <Xd>, <Xn>, <Xm>{, <shift> #<amount>}
//}

type OptionalOperand struct {
	Option            string     `"<" @Variable ">"`
	Parameter         *Parameter `( @@ |`
	OptionalParameter *Parameter `  ("{" @@ "}") )? `
}

type Operand struct {
	Parameter *Parameter `@@ |`
	Address   *Address   `@@`
}

func (op *Operand) Make() string {
	if op.Parameter != nil {
		return op.Parameter.Make()
	}
	return op.Address.Make()
}

type Address struct {
	Register       Parameter  `"[" @@`
	Offset         *Parameter `( (    "," @@    ) |`
	OptionalOffset *Parameter `  ("{" "," @@ "}") )?`
	Bang           *string    `"]" (@"!")?`
}

func (a *Address) Make() string {
	return "address"
}

type Parameter struct {
	Imm   *string             `("#" "<" @Variable      ">") |`
	F     *FixedWidthRegister `(    "<" @@             ">") |`
	V     *Expression         `(    "<" "R" ">" "<" @@ ">") |`
	Label *string             `     "<" @Variable      ">"`
}

func (p *Parameter) Make() string {
	return "param"
}

type Register struct {
	// Possible forms:
	// <Xn>, <Xm>, <Wn>, <Wn>, <R><m>, <X(n+1)>, <Xn|SP>
	Fixed    *FixedWidthRegister    `@@ |`
	Variable *VariableWidthRegister `@@`
}

type FixedWidthRegister struct {
	Width RegisterWidth `@@`
	Index Expression    `@@`
	Alt   *string       `("|" @ArchIdent)?`
}

type RegisterWidth struct {
	Width string `@("X"|"W")`
}

type VariableWidthRegister struct {
	Width string     `"R"`
	Index Expression `@@`
}

type Value struct {
	Variable *string `@Variable |`
	Number   *string `@Number`
}

type Expression struct {
	Value  *Value            `(@@`
	Suffix *ExpressionSuffix `@@?) |`
	Paren  *Expression       `"(" @@ ")"`
}
type ExpressionSuffix struct {
	Operator   string     `@("-" | "+")`
	Expression Expression `@@`
}

func ParseInstruction(s string) (*Instruction, error) {
	return Parser.ParseString("", s)
}

var Parser *participle.Parser[Instruction]
var asmDef *lexer.StatefulDefinition

func Lexer() *lexer.StatefulDefinition {
	return asmDef
}

func init() {
	asmDef = lexer.MustSimple([]lexer.SimpleRule{
		{"ArchIdent", `[A-Z]+`},
		{"Variable", `[a-z][a-z0-9]*`},
		{"Comment", `(?i)rem[^\n]*`},
		{"String", `"(\\"|[^"])*"`},
		{"Punct", `[-[!@#$%^&*()+_={}\|:;"'<,>.?/]|]`},
		{"Number", `[-+]?(\d*\.)?\d+`},
		{"whitespace", `[ \t]+`},
	})
	Parser = participle.MustBuild[Instruction](
		participle.Lexer(asmDef),
	)
}

func Tokenize(asm string) ([]Token, error) {
	mnemonicRE := regexp.MustCompile(`^\s*(\S+)\s*`)
	m := mnemonicRE.FindStringSubmatch(asm)
	if m == nil {
		return nil, fmt.Errorf("no leading mnemonic")
	}
	asm = asm[len(m[0]):]
	var toks []Token
	toks = append(toks, Token{"mnemonic", m[1]})

	res := []namedRe{
		{"space", regexp.MustCompile(`^\s+`)},
		{"ident", regexp.MustCompilePOSIX(`^[a-zA-Z][a-zA-Z0-9_]*`)},
		{"number", regexp.MustCompilePOSIX(`^(0x)?[0-9]+`)},
		{"stuff", regexp.MustCompile(`^[!|#,<>{}()[\]]`)},
		{"math", regexp.MustCompile(`^[+-]`)},
	}
tokenLoop:
	for len(asm) > 0 {
		for _, nr := range res {
			m := nr.re.FindStringSubmatch(asm)
			if m == nil {
				continue
			}
			toks = append(toks, Token{nr.name, m[0]})
			asm = asm[len(m[0]):]
			continue tokenLoop
		}
		return nil, fmt.Errorf("failed to tokenize at %q", asm)
	}
	return toks, nil
}

type namedRe struct {
	name string
	re   *regexp.Regexp
}

type Token struct {
	kind string
	Str  string
}

// Scan takes a pattern and parses the assembly according to the type asmParams, if the parsing is
// successful it encodes the information into a variable of type asmParams and passes that into the
// txl function and returns the result.
func Scanner[asmParams any](pattern string, txl func(asmParams) (opcode.Instruction, error)) (func(string) (asmParams, error), error) {
	return nil, nil
}

func printFields(typ reflect.Type, yield func(string) bool) {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if !yield(fmt.Sprintf("Field: %s, Type: %s\n", field.Name, field.Type)) {
			return
		}
		printFields(field.Type, yield)
	}
}

func X[asmParams any](f func(asmParams) bool) iter.Seq[string] {
	var params asmParams
	typ := reflect.TypeOf(params)
	return func(yield func(string) bool) {
		printFields(typ, yield)
	}
}
