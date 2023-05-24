package postgres

import "github.com/nao1215/aiondb/engine/parser/core"

// parseUpdate parses an UPDATE statement.
func (p *Parser) parseUpdate() (*core.Statement, error) {
	stmt := &core.Statement{}

	// Set DELETE decl
	updateDecl, err := p.consumeToken(core.TokenIDUpdate)
	if err != nil {
		return nil, err
	}
	stmt.Decls = append(stmt.Decls, updateDecl)

	// should be table name
	nameDecl, err := p.parseQuotedToken()
	if err != nil {
		return nil, err
	}
	updateDecl.Append(nameDecl)

	// should be SET
	setDecl, err := p.consumeToken(core.TokenIDSet)
	if err != nil {
		return nil, err
	}
	updateDecl.Append(setDecl)

	// should be a list of equality
	gotClause := false
	for p.tokens[p.index].ID != core.TokenIDWhere {
		if !p.hasNext() && gotClause {
			break
		}

		attributeDecl, err := p.parseAttribution()
		if err != nil {
			return nil, err
		}
		setDecl.Append(attributeDecl)
		if _, err := p.consumeToken(core.TokenIDComma); err != nil {
			return nil, err
		}

		// Got at least one clause
		gotClause = true
	}

	err = p.parseWhere(updateDecl)
	if err != nil {
		return nil, err
	}
	return stmt, nil
}
