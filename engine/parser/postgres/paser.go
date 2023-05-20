// Package postgres parses SQL queries with postgreSQL syntax.
package postgres

// Parser parses SQL queries conforming to PostgreSQL.
// It satisfies the Parser interface of the parser package.
type Parser struct{}

// NewParser returns a new Parser.
func NewParser() *Parser {
	return &Parser{}
}

// Parse parses the string (e.g., SQL query).
func (p *Parser) Parse(input string) (string, error) {
	return "", nil
}
