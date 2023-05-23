package core

import (
	"errors"
	"fmt"
)

var (
	// ErrLexerSyntax means "lex syntax error in user input"
	ErrLexerSyntax = errors.New("failed to lex instruction. syntax error")
	// ErrParserSyntax means "parse syntax error in user input"
	ErrParserSyntax = errors.New("failed to parse instruction. syntax error")
	// ErrEndOfStatement means "end of statement"
	ErrEndOfStatement = errors.New("end of statement")
	// ErrParseAfterCreateToken means "parse error after 'create' token"
	ErrParseAfterCreateToken = errors.New("'create' token must be followed by 'token', 'index'")
	// ErrParseAfterSelectToken means "parse error after 'select' token"
	ErrParseAfterSelectToken = errors.New("'select' token must be followed by attributes to select")
	// ErrNotDateFormat means input data is "not a date format"
	ErrNotDateFormat = errors.New("not a date format")
)

// Wrap return wrapping error with message.
// If e is nil, return new error with msg. If msg is empty string, return e.
func Wrap(e error, message string) error {
	if message == "" {
		return e
	}
	if e == nil {
		return errors.New(message)
	}
	return fmt.Errorf("%s: %w", message, e)
}
