// Package parser parses the string (e.g., SQL query) received by the aion shell
// and returns the result. If an SQL query is received as input, it returns a declaration tree.
// This package include lexer function, too.
package parser

import (
	"github.com/nao1215/aiondb/engine/parser/core"
	"github.com/nao1215/aiondb/engine/parser/postgres"
)

// Parser is an interface introduced to comprehensively
// parse the SQL syntax of common RDBMS (e.g. MySQL, SQLite, PostgreSQL, Oracle).
type Parser interface {
	Parse(input string) (string, error)
}

// NewParser returns a new Parser.
func NewParser(mode core.SQLSyntaxMode) Parser {
	switch mode {
	case core.SQLSyntaxModePostgreSQL:
		return postgres.NewParser()
	}
	return postgres.NewParser()
}
