// Package parser parses the string (e.g., SQL query) received by the aion shell
// and returns the result. If an SQL query is received as input, it returns a declaration tree.
package parser

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
