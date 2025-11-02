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
	rd, rdType, err := p.parseRegister()
	if err != nil {
		return nil, err
	}
	if p.cur.typ != tokenComma {
		return nil, fmt.Errorf("expected comma after destination register")
	}
	p.next()
	rn, _, err := p.parseRegister()
	if err != nil {
		return nil, err
	}
	if p.cur.typ != tokenComma {
		return nil, fmt.Errorf("expected comma after first source register")
	}
	p.next()

	switch p.cur.typ {
	case tokenRegister:
		rm, _, err := p.parseRegister()
		if err != nil {
			return nil, err
		}
		var sf byte
		if rdType == "x" {
			sf = 1
		}
		if p.cur.typ == tokenEOF || p.cur.typ == tokenNewline {
			return []opcode.Instruction{&opcode.AddShiftedRegister{Sf: sf, Rd: uint16(rd), Rn: uint16(rn), Rm: uint16(rm)}}, nil
		}
		if p.cur.typ != tokenComma {
			return nil, fmt.Errorf("expected comma after second source register")
		}
		p.next()

		shiftType, err := p.parseShiftType()
		if err != nil {
			return nil, err
		}

		if p.cur.typ != tokenNumber {
			return nil, fmt.Errorf("expected shift amount")
		}
		immString := p.cur.val
		if immString[0] == '#' {
			immString = immString[1:]
		}
		imm, err := strconv.Atoi(immString)
		if err != nil {
			return nil, fmt.Errorf("invalid immediate value: %v", err)
		}

		return []opcode.Instruction{&opcode.AddShiftedRegister{Sf: sf, Rd: uint16(rd), Rn: uint16(rn), Rm: uint16(rm), Shift: shiftType, Imm: uint16(imm)}}, nil

	case tokenNumber:
		immString := p.cur.val
		if immString[0] == '#' {
			immString = immString[1:]
		}
		imm, err := strconv.Atoi(immString)
		if err != nil {
			return nil, fmt.Errorf("invalid immediate value: %v", err)
		}
		var sf byte
		if rdType == "x" {
			sf = 1
		}
		var sh byte = 0
		return []opcode.Instruction{&opcode.AddImmedite{Sf: sf, Sh: sh, Rd: uint16(rd), Rn: uint16(rn), Imm: uint16(imm)}}, nil
	default:
		return nil, fmt.Errorf("unexpected token in add instruction: %s", p.cur)
	}
}

func (p *Parser) parseShiftType() (byte, error) {
	if p.cur.typ != tokenIdentifier {
		return 0, fmt.Errorf("expected shift type, got %s", p.cur)
	}
	shiftType := p.cur.val
	p.next()
	switch strings.ToLower(shiftType) {
	case "lsl":
		return 0, nil
	case "lsr":
		return 1, nil
	case "asr":
		return 2, nil
	default:
		return 0, fmt.Errorf("unknown shift type: %s", shiftType)
	}
}

func (p *Parser) parseRegister() (int, string, error) {
	if p.cur.typ != tokenRegister {
		return 0, "", fmt.Errorf("expected register, got %s", p.cur)
	}
	regType := string(p.cur.val[0])
	reg, err := strconv.Atoi(p.cur.val[1:])
	if err != nil {
		return 0, "", fmt.Errorf("invalid register number: %v", err)
	}
	p.next()
	return reg, regType, nil
}
