package postgres

import "github.com/nao1215/aiondb/engine/parser/core"

// parseDelete parses a DELETE statement.
func (p *Parser) parseDelete() (*core.Statement, error) {
	stmt := &core.Statement{}

	// Set DELETE decl
	deleteDecl, err := p.consumeToken(core.TokenIDDelete)
	if err != nil {
		return nil, err
	}
	stmt.Decls = append(stmt.Decls, deleteDecl)

	// should be From
	fromDecl, err := p.consumeToken(core.TokenIDFrom)
	if err != nil {
		return nil, err
	}
	deleteDecl.Append(fromDecl)

	// Should be a table name
	nameDecl, err := p.parseQuotedToken()
	if err != nil {
		return nil, err
	}
	fromDecl.Append(nameDecl)

	// MAY be WHERE  here
	if !p.hasNext() {
		return stmt, nil
	}

	err = p.parseWhere(deleteDecl)
	if err != nil {
		return nil, err
	}
	return stmt, nil
}
