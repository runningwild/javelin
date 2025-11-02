package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/runningwild/javelin/opcode"
)

type Parser struct {
	lexer  *lexer
	tokens chan token
	cur    token
}

func New(input string) *Parser {
	l, c := lex(input)
	p := &Parser{
		lexer:  l,
		tokens: c,
	}
	p.next()
	return p
}

func (p *Parser) next() {
	p.cur = <-p.tokens
}

func (p *Parser) Parse() ([]opcode.Instruction, error) {
	var insts []opcode.Instruction
	for p.cur.typ != tokenEOF {
		if p.cur.typ == tokenNewline {
			p.next()
			continue
		}
		inst, err := p.parseInstruction()
		if err != nil {
			return nil, err
		}
		insts = append(insts, inst...)
		p.next()
	}
	return insts, nil
}

func (p *Parser) parseInstruction() ([]opcode.Instruction, error) {
	switch strings.ToLower(p.cur.val) {
	case "add":
		p.next()
		return p.parseAdd()
	default:
		return nil, fmt.Errorf("unknown instruction: %s", p.cur.val)
	}
}

func (p *Parser) parseAdd() ([]opcode.Instruction, error) {
	rd, err := p.parseRegister()
	if err != nil {
		return nil, err
	}
	if p.cur.typ != tokenComma {
		return nil, fmt.Errorf("expected comma after destination register")
	}
	p.next()
	rn, err := p.parseRegister()
	if err != nil {
		return nil, err
	}
	if p.cur.typ != tokenComma {
		return nil, fmt.Errorf("expected comma after first source register")
	}
	p.next()

	switch p.cur.typ {
	case tokenRegister:
		rm, err := p.parseRegister()
		if err != nil {
			return nil, err
		}
		// TODO: Handle shifting
		return []opcode.Instruction{&opcode.AddShiftedRegister{Rd: uint16(rd), Rn: uint16(rn), Rm: uint16(rm)}}, nil
	case tokenNumber:
		imm, err := strconv.Atoi(p.cur.val)
		if err != nil {
			return nil, fmt.Errorf("invalid immediate value: %v", err)
		}
		var sf byte = 1
		var sh byte = 0
		// TODO: This is not right
		return []opcode.Instruction{&opcode.AddImmedite{Sf: sf, Sh: sh, Rd: uint16(rd), Rn: uint16(rn), Imm: uint16(imm)}}, nil
	default:
		return nil, fmt.Errorf("unexpected token in add instruction: %s", p.cur)
	}
}

func (p *Parser) parseRegister() (int, error) {
	if p.cur.typ != tokenRegister {
		return 0, fmt.Errorf("expected register, got %s", p.cur)
	}
	reg, err := strconv.Atoi(p.cur.val[1:])
	if err != nil {
		return 0, fmt.Errorf("invalid register number: %v", err)
	}
	p.next()
	return reg, nil
}
