// Package parser parses the string (e.g., SQL query) received by the aion shell
// and returns the result. If an SQL query is received as input, it returns a declaration tree.
// This package include lexer function, too.
package parser

import (
	"github.com/nao1215/aiondb/engine/parser/postgres"
)

// SQLSyntaxMode is the SQL syntax mode.
type SQLSyntaxMode int

const (
	// SQLSyntaxModeDefault is the default SQL syntax mode. It is the same as MySQL.
	SQLSyntaxModeDefault SQLSyntaxMode = 0
	// SQLSyntaxModeMySQL is the MySQL SQL syntax mode.
	SQLSyntaxModeMySQL SQLSyntaxMode = 1
	// SQLSyntaxModePostgreSQL is the PostgreSQL SQL syntax mode.
	SQLSyntaxModePostgreSQL SQLSyntaxMode = 2
	// SQLSyntaxModeOracle is the Oracle SQL syntax mode.
	SQLSyntaxModeOracle SQLSyntaxMode = 3
	// SQLSyntaxModeSQLite is the SQLite SQL syntax mode.
	SQLSyntaxModeSQLite SQLSyntaxMode = 4
)

// Parser is an interface introduced to comprehensively
// parse the SQL syntax of common RDBMS (e.g. MySQL, SQLite, PostgreSQL, Oracle).
type Parser interface {
	Parse(input string) (string, error)
}

// NewParser returns a new Parser.
func NewParser(mode SQLSyntaxMode) Parser {
	switch mode {
	case SQLSyntaxModePostgreSQL:
		return postgres.NewParser()
	}
	return postgres.NewParser()
}
