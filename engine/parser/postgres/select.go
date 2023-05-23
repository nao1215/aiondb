package postgres

import (
	"errors"

	"github.com/nao1215/aiondb/engine/parser/core"
)

// parseSelect parses the SELECT statement.
func (p *Parser) parseSelect(tokens []core.Token) (*core.Statement, error) {
	stmt := &core.Statement{}
	var err error

	// Create select decl
	selectDecl := core.NewDecl(tokens[p.index])
	stmt.Decls = append(stmt.Decls, selectDecl)

	// After select token, should be either
	// a StarToken
	// a list of table names + (StarToken Or Attribute)
	// a builtin func (COUNT, MAX, ...)
	if err = p.next(); err != nil {
		return nil, core.ErrParseAfterSelectToken
	}

	distinctDecl, distinctOpen, err := p.parseDistinct(selectDecl)
	if err != nil {
		return nil, err
	}

	if err = p.parseColumnBeforeFromToken(selectDecl, distinctDecl, distinctOpen); err != nil {
		return nil, err
	}

	// Must be from now
	if tokens[p.index].ID != core.TokenIDFrom {
		return nil, errors.New("syntax error near " + tokens[p.index].Lexeme.String())
	}
	fromDecl := core.NewDecl(tokens[p.index])
	selectDecl.Append(fromDecl)

	// Now must be a list of table
	for {
		// string
		if err = p.next(); err != nil {
			return nil, errors.New("unexpected end. Syntax error near " + tokens[p.index].Lexeme.String())
		}
		tableNameDecl, err := p.parseAttribute()
		if err != nil {
			return nil, err
		}
		fromDecl.Append(tableNameDecl)

		// If no next, then it's implicit where
		if !p.hasNext() {
			appendImplicitWhereAll(selectDecl)
			return stmt, nil
		}
		// if not comma, break
		if tokens[p.index].ID != core.TokenIDComma {
			break // No more table
		}
	}

	// JOIN OR ...?
	for p.is(core.TokenIDJoin) {
		joinDecl, err := p.parseJoin()
		if err != nil {
			return nil, err
		}
		selectDecl.Append(joinDecl)
	}

	hazWhereClause := false
	for {
		switch p.current().ID {
		case core.TokenIDWhere:
			err := p.parseWhere(selectDecl)
			if err != nil {
				return nil, err
			}
			hazWhereClause = true
		case core.TokenIDOrder:
			if hazWhereClause == false {
				// WHERE clause is implicit
				appendImplicitWhereAll(selectDecl)
			}
			err := p.parseOrderBy(selectDecl)
			if err != nil {
				return nil, err
			}
		case core.TokenIDLimit:
			limitDecl, err := p.consumeToken(core.TokenIDLimit)
			if err != nil {
				return nil, err
			}
			selectDecl.Append(limitDecl)

			numDecl, err := p.consumeToken(core.TokenIDNumber)
			if err != nil {
				return nil, err
			}
			limitDecl.Append(numDecl)
		case core.TokenIDOffset:
			offsetDecl, err := p.consumeToken(core.TokenIDOffset)
			if err != nil {
				return nil, err
			}
			selectDecl.Append(offsetDecl)

			offsetValue, err := p.consumeToken(core.TokenIDNumber)
			if err != nil {
				return nil, err
			}
			offsetDecl.Append(offsetValue)
		case core.TokenIDFor:
			err := p.parseForUpdate(selectDecl)
			if err != nil {
				return nil, err
			}
		default:
			return stmt, nil
		}
	}
}

// parseDistinct parse 'distinct' clause.
func (p *Parser) parseDistinct(selectDecl *core.Decl) (*core.Decl, bool, error) {
	if !p.is(core.TokenIDDistinct) {
		return nil, false, nil
	}

	var distinctDecl *core.Decl
	distinctDecl, err := p.consumeToken(core.TokenIDDistinct)
	if err != nil {
		return distinctDecl, false, err
	}

	distinctOpen := false
	if p.is(core.TokenIDOn) {
		if err := p.next(); err != nil {
			return distinctDecl, false, err
		}
		if !p.is(core.TokenIDBracketOpening) {
			return distinctDecl, false, errors.New("syntax error. opening bracket expected")
		}
		if err := p.next(); err != nil {
			return distinctDecl, false, err
		}
		distinctOpen = true
	}
	selectDecl.Append(distinctDecl)

	return distinctDecl, distinctOpen, nil
}

// parseColumnBeforeFromToken parses the column before FROM token.
func (p *Parser) parseColumnBeforeFromToken(selectDecl, distinctDecl *core.Decl, distinctOpen bool) error {
	for {
		switch {
		case p.is(core.TokenIDCount):
			attrDecl, err := p.parseBuiltinFunc()
			if err != nil {
				return err
			}
			selectDecl.Append(attrDecl)
		default:
			attrDecl, err := p.parseAttribute()
			if err != nil {
				return err
			}
			if distinctOpen {
				distinctDecl.Append(attrDecl)
			}
			selectDecl.Append(attrDecl)
		}

		switch {
		case distinctOpen && p.is(core.TokenIDBracketClosing):
			if err := p.next(); err != nil {
				return err
			}
			distinctOpen = false
			continue
		case p.is(core.TokenIDComma):
			if err := p.next(); err != nil {
				return err
			}
			continue
		}
		break
	}
	return nil
}

// appendImplicitWhereAll appends implicit where clause.
func appendImplicitWhereAll(decl *core.Decl) {
	whereDecl := core.NewDecl(core.Token{
		ID:     core.TokenIDWhere,
		Lexeme: "where",
	})

	whereDecl.Append(core.NewDecl(core.Token{
		ID:     core.TokenIDNumber,
		Lexeme: "1",
	}))
	decl.Append(whereDecl)
}

// parseOrderBy parses 'order by' clause.
func (p *Parser) parseOrderBy(selectDecl *core.Decl) error {
	orderDecl, err := p.consumeToken(core.TokenIDOrder)
	if err != nil {
		return err
	}
	selectDecl.Append(orderDecl)

	_, err = p.consumeToken(core.TokenIDBy)
	if err != nil {
		return err
	}

	for {
		// parse attribute now
		attrDecl, err := p.parseAttribute()
		if err != nil {
			return err
		}
		orderDecl.Append(attrDecl)

		if p.is(core.TokenIDAsc, core.TokenIDDesc) {
			decl, err := p.consumeToken(core.TokenIDAsc, core.TokenIDDesc)
			if err != nil {
				return err
			}
			attrDecl.Append(decl)
		}

		if !p.is(core.TokenIDComma) {
			break
		}

		if _, err = p.consumeToken(core.TokenIDComma); err != nil {
			return nil
		}
	}
	return nil
}

func (p *Parser) parseForUpdate(decl *core.Decl) error {
	// Optionnal
	if !p.is(core.TokenIDFor) {
		return nil
	}
	d, err := p.consumeToken(core.TokenIDFor)
	if err != nil {
		return err
	}

	u, err := p.consumeToken(core.TokenIDUpdate)
	if err != nil {
		return err
	}

	d.Append(u)
	decl.Append(d)
	return nil
}
