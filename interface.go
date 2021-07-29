package lexer

type Lexer interface {
	NextToken() (Token, bool)
	Start()
	Emit(t TokenType)
	Next() (r rune)
	Peek() (r rune)
	Ignore()
	Backup()
	Accept(valid string) bool
	AcceptRun(valid string)
	Errorf(format string, args ...interface{}) StateFn
}

type Token interface {
	Type() TokenType
	At() int
	Value() string
}
