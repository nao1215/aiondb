package core

// SQLSyntaxMode is the SQL syntax mode.
type SQLSyntaxMode uint64

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
