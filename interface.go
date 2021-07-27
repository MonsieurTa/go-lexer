package lexer

type Lexer interface {
	Start() chan Token
	Emit(t TokenType)
	Next() (r rune)
	Peek() (r rune)
	Ignore()
	Backup()
	Accept(valid string) bool
	AcceptRun(valid string)
	Errorf(format string, args ...interface{}) StateFn
}
