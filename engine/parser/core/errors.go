// Package core implements the common functions used in the parsing
// of SQL queries. SQL queries have dialects specific to each RDBMS.
// AION DB has been implemented to accommodate these dialects. This
// package aims to minimize code duplication.
package core

import (
	"errors"
	"fmt"
)

var ()

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
