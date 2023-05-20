// Package parser parses the string (e.g., SQL query) received by the aion shell
// and returns the result. If an SQL query is received as input, it returns a declaration tree.
package parser

// Parser XXX
type Parser struct {
}

// NewParser returns a new Parser.
func NewParser() *Parser {
	return &Parser{}
}

// Parse parses the string (e.g., SQL query).
func (p *Parser) Parse(input string) (string, error) {
	return "", nil
}
