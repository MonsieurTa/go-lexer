package lexer

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const (
	EOFRune    rune      = -1
	ErrorToken TokenType = -1
	EmptyToken TokenType = 0
)

type TokenType int

type StateFn func(Lexer) StateFn

type Token struct {
	Type  TokenType
	At    int
	Value string
}

// lexer holds the state of the scanner
type lexer struct {
	name       string     // used for error reports
	input      string     // the string being scanned
	start      int        // start position of current item
	pos        int        // current position in the input
	width      int        // width of last rune read
	tokens     chan Token // channel of scanned tokens
	startState StateFn
}

func New(name, s string, startState StateFn) Lexer {
	return &lexer{
		name:       name,
		input:      s,
		startState: startState,
	}
}

func (l *lexer) Start() chan Token {
	buffSize := len(l.input) / 2
	if buffSize == 0 {
		buffSize = 1
	}
	l.tokens = make(chan Token, buffSize)

	go l.run()
	return l.tokens
}

func (l *lexer) run() {
	for state := l.startState; state != nil; {
		state = state(l)
	}
	close(l.tokens)
}

func (l *lexer) Emit(t TokenType) {
	l.tokens <- Token{t, l.start, l.input[l.start:l.pos]}
	l.start = l.pos
}

func (l *lexer) Next() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return EOFRune
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return
}

func (l *lexer) Peek() (r rune) {
	r = l.Next()
	l.Backup()
	return
}

func (l *lexer) Ignore() {
	l.start = l.pos
}

func (l *lexer) Backup() {
	l.pos -= l.width
}

func (l *lexer) Accept(valid string) bool {
	if strings.ContainsRune(valid, l.Next()) {
		return true
	}
	l.Backup()
	return false
}

func (l *lexer) AcceptRun(valid string) {
	for strings.ContainsRune(valid, l.Next()) {
	}
	l.Backup()
}

func (l *lexer) Errorf(format string, args ...interface{}) StateFn {
	l.tokens <- Token{
		ErrorToken,
		l.start,
		fmt.Sprintf(format, args...),
	}
	return nil
}
