package postgres

import (
	"github.com/nao1215/aiondb/engine/parser/core"
)

func (p *Parser) parseDrop() (*core.Statement, error) {
	stmt := &core.Statement{}

	trDecl, err := p.consumeToken(core.TokenIDDrop)
	if err != nil {
		return nil, err
	}
	stmt.Decls = append(stmt.Decls, trDecl)

	tableDecl, err := p.consumeToken(core.TokenIDTable)
	if err != nil {
		return nil, err
	}
	trDecl.Append(tableDecl)

	// Should be a table name
	nameDecl, err := p.parseQuotedToken()
	if err != nil {
		return nil, err
	}
	tableDecl.Append(nameDecl)

	return stmt, nil
}
