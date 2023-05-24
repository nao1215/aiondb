package postgres

import "github.com/nao1215/aiondb/engine/parser/core"

// parseTruncate parses a TRUNCATE statement.
func (p *Parser) parseTruncate() (*core.Statement, error) {
	stmt := &core.Statement{}

	// Set TRUNCATE decl
	trDecl, err := p.consumeToken(core.TokenIDTruncate)
	if err != nil {
		return nil, err
	}
	stmt.Decls = append(stmt.Decls, trDecl)

	// Should be a table name
	nameDecl, err := p.parseQuotedToken()
	if err != nil {
		return nil, err
	}
	trDecl.Append(nameDecl)

	return stmt, nil
}
