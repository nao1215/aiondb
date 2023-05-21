package core

import (
	"errors"
	"fmt"
)

var (
	// ErrLexerSyntax means "syntax error in user input"
	ErrLexerSyntax = errors.New("failed to lex instruction. syntax error")
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
