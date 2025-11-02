package parser

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

type tokenType int

const (
	tokenError tokenType = iota
	tokenEOF
	tokenIdentifier
	tokenNumber
	tokenRegister
	tokenComma
	tokenNewline
)

type token struct {
	typ tokenType
	val string
}

func (t token) String() string {
	switch t.typ {
	case tokenEOF:
		return "EOF"
	case tokenError:
		return t.val
	}
	return fmt.Sprintf("%q", t.val)
}

type lexer struct {
	input  string
	start  int
	pos    int
	width  int
	tokens chan token
}

func lex(input string) (*lexer, chan token) {
	l := &lexer{
		input:  input,
		tokens: make(chan token),
	}
	go l.run()
	return l, l.tokens
}

func (l *lexer) run() {
	for state := lexCode; state != nil; {
		state = state(l)
	}
	close(l.tokens)
}

func (l *lexer) emit(t tokenType) {
	l.tokens <- token{t, l.input[l.start:l.pos]}
	l.start = l.pos
}

func (l *lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return 0
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width
	return r
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func lexCode(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case r == 0:
			l.emit(tokenEOF)
			return nil
		case r == '\n':
			l.emit(tokenNewline)
		case unicode.IsSpace(r):
			l.ignore()
		case r == ',':
			l.emit(tokenComma)
		case r == '#':
			return lexNumber
		case r == 'X' || r == 'x' || r == 'W' || r == 'w':
			l.backup()
			return lexRegister
		case '0' <= r && r <= '9':
			l.backup()
			return lexNumber
		case unicode.IsLetter(r):
			return lexIdentifier
		default:
			return l.errorf("unrecognized character: %q", r)
		}
	}
}

func lexIdentifier(l *lexer) stateFn {
	for {
		r := l.next()
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			l.backup()
			l.emit(tokenIdentifier)
			return lexCode
		}
	}
}

func lexNumber(l *lexer) stateFn {
	// Optional '#' prefix
	if l.peek() == '#' {
		l.next()
	}
	// TODO: Handle hex
	for {
		r := l.next()
		if !unicode.IsDigit(r) {
			l.backup()
			l.emit(tokenNumber)
			return lexCode
		}
	}
}

func lexRegister(l *lexer) stateFn {
	l.next() // X,x,W,w
	r := l.next()
	if !unicode.IsDigit(r) {
		return l.errorf("expected a digit after register prefix")
	}
	for {
		r := l.next()
		if !unicode.IsDigit(r) {
			l.backup()
			break
		}
	}
	l.emit(tokenRegister)
	return lexCode
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.tokens <- token{
		tokenError,
		fmt.Sprintf(format, args...),
	}
	return nil
}

type stateFn func(*lexer) stateFn
