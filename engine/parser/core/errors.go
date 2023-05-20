package core

import (
	"errors"
	"fmt"
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
