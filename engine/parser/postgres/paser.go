// Package postgres parses SQL queries with postgreSQL syntax.
package postgres

import (
	"fmt"

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

	for p.hasNext() {
		// Found a new instruction
		if tokens[p.index].ID == core.TokenIDSemicolon {
			p.index++
			continue
		}

		// Ignore space token, not needed anymore
		if tokens[p.index].ID == core.TokenIDSpace {
			p.index++
			continue
		}
		// Now,
		// Create a logical tree of all tokens
		// We start with first order query
		// CREATE, SELECT, INSERT, UPDATE, DELETE, TRUNCATE, DROP, EXPLAIN
		switch tokens[p.index].ID {
		case core.TokenIDCreate:
			s, err := p.parseCreate(tokens)
			if err != nil {
				return nil, err
			}
			p.stmt = append(p.stmt, *s)
			break
			// TODO: implement other statements
		}
	}
	return p.stmt, nil
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
	return p.index+1 < len(p.tokens)
}

// mustHaveNext returns the next token if it exists.
func (p *Parser) mustHaveNext(tokenTypes ...core.TokenID) (t core.Token, err error) {
	if !p.hasNext() {
		return t, p.syntaxError()
	}
	if err = p.next(); err != nil {
		return t, err
	}

	for _, tokenType := range tokenTypes {
		if p.is(tokenType) {
			return p.tokens[p.index], nil
		}
	}
	return t, p.syntaxError()
}

// is returns true if the current token is one of the specified tokens.
func (p *Parser) is(tokenTypes ...core.TokenID) bool {
	for _, tokenType := range tokenTypes {
		if p.current().ID == tokenType {
			return true
		}
	}
	return false
}

// isNot returns true if the current token is not one of the specified tokens.
func (p *Parser) isNot(tokenTypes ...core.TokenID) bool {
	return !p.is(tokenTypes...)
}

// isNext returns true if the next token is one of the specified tokens.
func (p *Parser) isNext(tokenTypes ...core.TokenID) (t core.Token, err error) {
	if !p.hasNext() {
		return t, p.syntaxError()
	}
	for _, tokenType := range tokenTypes {
		if p.tokens[p.index+1].ID == tokenType {
			return p.tokens[p.index+1], nil
		}
	}
	return t, p.syntaxError()
}

// current returns the current token.
func (p *Parser) current() core.Token {
	return p.tokens[p.index]
}

// consumeToken consumes the current token and returns it.
func (p *Parser) consumeToken(tokenTypes ...core.TokenID) (*core.Decl, error) {
	if !p.is(tokenTypes...) {
		return nil, p.syntaxError()
	}
	decl := core.NewDecl(p.tokens[p.index])
	p.next()
	return decl, nil
}

// syntaxError returns a syntax error.
func (p *Parser) syntaxError() error {
	if p.index == 0 {
		return fmt.Errorf("Syntax error near %s %s",
			p.tokens[p.index].Lexeme, p.tokens[p.index+1].Lexeme)
	} else if !p.hasNext() {
		return fmt.Errorf("Syntax error near %s %s",
			p.tokens[p.index-1].Lexeme, p.tokens[p.index].Lexeme)
	}
	return fmt.Errorf("Syntax error near %s %s %s",
		p.tokens[p.index-1].Lexeme, p.tokens[p.index].Lexeme, p.tokens[p.index+1].Lexeme)
}

// parseAttribute parse an attribute of the form
// table.foo
// table.*
// "table".foo
// "table"."foo"
// foo
func (p *Parser) parseAttribute() (*core.Decl, error) {
	quoted := false
	quoteToken := core.TokenIDDoubleQuote

	if p.is(core.TokenIDDoubleQuote) || p.is(core.TokenIDBacktick) {
		quoteToken = p.current().ID
		quoted = true
		if err := p.next(); err != nil {
			return nil, err
		}
	}

	// should be a StringToken here
	// If there is a point after, it's a table name,
	// if not, it's the attribute
	if !p.is(core.TokenIDString, core.TokenIDStar) {
		return nil, p.syntaxError()
	}
	decl := core.NewDecl(p.current())

	if quoted {
		// Check there is a closing quote
		if _, err := p.mustHaveNext(quoteToken); err != nil {
			return nil, err
		}
	}
	quoted = false

	// If no next token, and not quoted, then is was the attribute name
	if err := p.next(); err != nil {
		return decl, nil
	}

	if !p.is(core.TokenIDPeriod) {
		return decl, nil
	}
	if _, err := p.consumeToken(core.TokenIDPeriod); err != nil {
		return nil, err
	}

	// mayby attribute is quoted as well
	if p.is(core.TokenIDDoubleQuote) || p.is(core.TokenIDBacktick) {
		quoteToken = p.current().ID
		quoted = true
		if err := p.next(); err != nil {
			return nil, err
		}
	}

	// if so, next must be the attribute name or a star
	attributeDecl, err := p.consumeToken(core.TokenIDString, core.TokenIDStar)
	if err != nil {
		return nil, err
	}
	attributeDecl.Append(decl)

	if quoted {
		// Check there is a closing quote
		if _, err := p.consumeToken(quoteToken); err != nil {
			return nil, fmt.Errorf("expected closing quote: %s", err)
		}
	}
	return attributeDecl, nil
}

// parseQuotedToken parse a token of the form
// table
// "table"
func (p *Parser) parseQuotedToken() (*core.Decl, error) {
	quoted := false
	quoteToken := core.TokenIDDoubleQuote

	if p.is(core.TokenIDDoubleQuote) || p.is(core.TokenIDBacktick) {
		quoted = true
		quoteToken = p.current().ID
		if err := p.next(); err != nil {
			return nil, err
		}
	}

	// shoud be a StringToken here
	if !p.is(core.TokenIDString) {
		return nil, p.syntaxError()
	}
	decl := core.NewDecl(p.current())

	if quoted {
		// Check there is a closing quote
		if _, err := p.mustHaveNext(quoteToken); err != nil {
			return nil, err
		}
	}
	p.next()
	return decl, nil
}

// parseType parse a type of the form.
func (p *Parser) parseType() (*core.Decl, error) {
	typeDecl, err := p.consumeToken(core.TokenIDString)
	if err != nil {
		return nil, err
	}

	// Maybe a complex type
	if !p.is(core.TokenIDBracketOpening) {
		return typeDecl, nil
	}

	if _, err = p.consumeToken(core.TokenIDBracketOpening); err != nil {
		return nil, err
	}

	sizeDecl, err := p.consumeToken(core.TokenIDNumber)
	if err != nil {
		return nil, err
	}
	typeDecl.Append(sizeDecl)

	if _, err = p.consumeToken(core.TokenIDBracketClosing); err != nil {
		return nil, err
	}
	return typeDecl, nil
}

func (p *Parser) parseStringLiteral() (*core.Decl, error) {
	singleQuoted := p.is(core.TokenIDSingleQuote)
	_, err := p.consumeToken(core.TokenIDSingleQuote, core.TokenIDDoubleQuote)
	if err != nil {
		return nil, err
	}

	valueDecl, err := p.consumeToken(core.TokenIDString)
	if err != nil {
		return nil, err
	}

	if (singleQuoted && p.is(core.TokenIDDoubleQuote)) || (!singleQuoted && p.is(core.TokenIDSingleQuote)) {
		return nil, fmt.Errorf("quotation marks do not match")
	}
	if _, err = p.consumeToken(core.TokenIDSingleQuote, core.TokenIDDoubleQuote); err != nil {
		return nil, err
	}
	return valueDecl, nil
}
