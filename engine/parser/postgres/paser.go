// Package postgres parses SQL queries with postgreSQL syntax.
package postgres

import (
	"github.com/nao1215/aiondb/engine/parser/core"
)

// Parser parses SQL queries conforming to PostgreSQL.
// It satisfies the Parser interface of the parser package.
type Parser struct {
	stmt     []core.Statement
	index    int
	tokenLen int
	tokens   []core.Token
}

// NewParser returns a new Parser.
func NewParser() *Parser {
	return &Parser{}
}

// Parse parses the string (e.g., SQL query).
func (p *Parser) Parse(input string) ([]core.Statement, error) {
	lexer := NewLexer(input)
	tokens, err := lexer.Lex()
	if err != nil {
		return nil, err
	}

	stmt, err := p.parse(tokens)
	if err != nil {
		return nil, err
	}

	if len(stmt) == 0 {
		return nil, core.Wrap(core.ErrParserSyntax, input)
	}
	return stmt, nil
}

// parse parses the tokens. It is called by the Parse method.
func (p *Parser) parse(tokens []core.Token) ([]core.Statement, error) {
	tokens = core.StripSpaces(tokens)

	p.tokens = tokens
	p.tokenLen = len(tokens)
	p.index = 0

	return nil, nil
}

// next returns the next token.
func (p *Parser) next() error {
	if p.hasNext() {
		p.index++
		return nil
	}
	return core.ErrEndOfStatement
}

// hasNext returns true if the next token exists.
func (p *Parser) hasNext() bool {
	if p.index+1 < len(p.tokens) {
		return true
	}
	return false
}
