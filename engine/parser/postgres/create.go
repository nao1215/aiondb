package postgres

import (
	"fmt"
	"strings"

	"github.com/nao1215/aiondb/engine/parser/core"
)

// parseCreate parses the CREATE statement.
func (p *Parser) parseCreate(tokens []core.Token) (*core.Statement, error) {
	stmt := &core.Statement{}

	// Set CREATE decl
	createDecl := core.NewDecl(tokens[p.index])
	stmt.Decls = append(stmt.Decls, createDecl)

	// After create token, should be either TABLE, INDEX, ...
	if !p.hasNext() {
		return nil, core.ErrParseAfterCreateToken
	}
	p.index++

	switch tokens[p.index].ID {
	case core.TokenIDTable:
		d, err := p.parseTable(tokens)
		if err != nil {
			return nil, err
		}
		createDecl.Append(d)
		break
	case core.TokenIDIndex:
		d, err := p.parseIndex(tokens)
		if err != nil {
			return nil, err
		}
		createDecl.Append(d)
		break
	case core.TokenIDUnique:
		u, err := p.consumeToken(core.TokenIDUnique)
		if err != nil {
			return nil, err
		}
		// should have index after unique here
		if !p.hasNext() || tokens[p.index].ID != core.TokenIDIndex {
			return nil, fmt.Errorf("expected INDEX after UNIQUE")
		}
		d, err := p.parseIndex(tokens)
		if err != nil {
			return nil, err
		}
		d.Append(u)
		createDecl.Append(d)
		break
	default:
		return nil, fmt.Errorf("Parsing error near <%s>", tokens[p.index].Lexeme)
	}
	return stmt, nil
}

// parseTable parses the CREATE TABLE statement.
func (p *Parser) parseTable(tokens []core.Token) (*core.Decl, error) {
	var err error
	tableDecl := core.NewDecl(tokens[p.index])
	p.index++

	if tableDecl, err = p.parseIf(tableDecl); err != nil {
		return nil, err
	}

	// Now we should found table name
	nameTable, err := p.parseAttribute()
	if err != nil {
		return nil, p.syntaxError()
	}
	tableDecl.Append(nameTable)

	// Now we should found brackets
	if !p.hasNext() || tokens[p.index].ID != core.TokenIDBracketOpening {
		return nil, fmt.Errorf("Table name token must be followed by table definition")
	}
	p.index++

	for p.index < len(tokens) {
		switch p.current().ID {
		case core.TokenIDPrimary:
			_, err := p.parsePrimaryKey()
			if err != nil {
				return nil, err
			}
			continue
		default:
		}

		// Closing bracket ?
		if tokens[p.index].ID == core.TokenIDBracketClosing {
			p.consumeToken(core.TokenIDBracketClosing)
			break
		}

		// New attribute name
		newAttribute, err := p.parseQuotedToken()
		if err != nil {
			return nil, err
		}
		tableDecl.Append(newAttribute)

		newAttributeType, err := p.parseType()
		if err != nil {
			return nil, err
		}
		newAttribute.Append(newAttributeType)

		// All the following tokens until bracket or comma are column constraints.
		// Column constraints can be listed in any order.
		for p.isNot(core.TokenIDBracketClosing, core.TokenIDComma) {
			switch p.current().ID {
			case core.TokenIDUnique:
				uniqueDecl, err := p.consumeToken(core.TokenIDUnique)
				if err != nil {
					return nil, err
				}
				newAttribute.Append(uniqueDecl)
			case core.TokenIDNot:
				if _, err = p.isNext(core.TokenIDNull); err == nil {
					notDecl, err := p.consumeToken(core.TokenIDNot)
					if err != nil {
						return nil, err
					}
					newAttribute.Append(notDecl)

					nullDecl, err := p.consumeToken(core.TokenIDNull)
					if err != nil {
						return nil, err
					}
					notDecl.Append(nullDecl)
				}
			case core.TokenIDPrimary:
				if _, err = p.isNext(core.TokenIDKey); err == nil {
					newPrimary := core.NewDecl(tokens[p.index])
					newAttribute.Append(newPrimary)

					if err = p.next(); err != nil {
						return nil, fmt.Errorf("Unexpected end")
					}

					newKey := core.NewDecl(tokens[p.index])
					newPrimary.Append(newKey)

					if err = p.next(); err != nil {
						return nil, fmt.Errorf("Unexpected end")
					}
				}
			case core.TokenIDAutoincrement:
				autoincDecl, err := p.consumeToken(core.TokenIDAutoincrement)
				if err != nil {
					return nil, err
				}
				newAttribute.Append(autoincDecl)
			case core.TokenIDWith:
				if strings.ToLower(newAttributeType.Lexeme.String()) == "timestamp" {
					withDecl, err := p.consumeToken(core.TokenIDWith)
					if err != nil {
						return nil, err
					}

					timeDecl, err := p.consumeToken(core.TokenIDTime)
					if err != nil {
						return nil, err
					}

					zoneDecl, err := p.consumeToken(core.TokenIDZone)
					if err != nil {
						return nil, err
					}
					newAttributeType.Append(withDecl)
					withDecl.Append(timeDecl)
					timeDecl.Append(zoneDecl)
				}
			case core.TokenIDDefault:
				dDecl, err := p.parseDefaultClause()
				if err != nil {
					return nil, err
				}
				newAttribute.Append(dDecl)
			default: // Unknown column constraint
				return nil, p.syntaxError()
			}
		}

		// The current token is either closing bracked or comma.
		// Closing bracket means table parsing stops.
		if tokens[p.index].ID == core.TokenIDBracketClosing {
			p.index++
			break
		}
		// Comma means continue on next table column.
		p.index++
	}
	return tableDecl, nil
}

// parseIndex parses 'index' tokens.
// INDEX index_name ON table_name (col1, col2)
func (p *Parser) parseIndex(tokens []core.Token) (*core.Decl, error) {
	var err error
	indexDecl := core.NewDecl(tokens[p.index])
	p.index++

	// Maybe have "IF NOT EXISTS" here
	if indexDecl, err = p.parseIf(indexDecl); err != nil {
		return nil, err
	}

	// Now we should found index name
	nameIndex, err := p.parseAttribute()
	if err != nil {
		return nil, p.syntaxError()
	}
	indexDecl.Append(nameIndex)

	// ON
	if !p.hasNext() || tokens[p.index].ID != core.TokenIDOn {
		return nil, fmt.Errorf("Expected ON")
	}
	p.index++

	// Now we should found table name
	nameTable, err := p.parseAttribute()
	if err != nil {
		return nil, p.syntaxError()
	}
	indexDecl.Append(nameTable)

	// Now we should found brackets
	if !p.hasNext() || tokens[p.index].ID != core.TokenIDBracketOpening {
		return nil, fmt.Errorf("Table name token must be followed by table definition")
	}
	p.index++

	for p.index < len(tokens) {
		// New attribute name
		newAttribute, err := p.parseQuotedToken()
		if err != nil {
			return nil, err
		}
		indexDecl.Append(newAttribute)

		// Closing bracket ?
		if tokens[p.index].ID == core.TokenIDBracketClosing {
			p.consumeToken(core.TokenIDBracketClosing)
			break
		}

		// All the following tokens until bracket or comma are column constraints.
		// Column constraints can be listed in any order.
		for p.isNot(core.TokenIDBracketClosing, core.TokenIDComma) {
			switch p.current().ID {
			case core.TokenIDCollate:
				collateDecl, err := p.consumeToken(core.TokenIDCollate)
				if err != nil {
					return nil, p.syntaxError()
				}
				newAttribute.Append(collateDecl)

				n, err := p.consumeToken(core.TokenIDNocase)
				if err != nil {
					return nil, p.syntaxError()
				}
				collateDecl.Append(n)
			default:
				// Unknown column constraint
				return nil, p.syntaxError()
			}
		}
		// The current token is either closing bracked or comma.
		// Closing bracket means table parsing stops.
		if tokens[p.index].ID == core.TokenIDBracketClosing {
			p.index++
			break
		}
		// Comma means continue on next table column.
		p.index++
	}
	return indexDecl, nil
}

// parseIf parses 'if' tokens.
func (p *Parser) parseIf(decl *core.Decl) (*core.Decl, error) {
	if !p.is(core.TokenIDIf) {
		return nil, nil
	}

	ifDecl, err := p.consumeToken(core.TokenIDIf)
	if err != nil {
		return nil, err
	}
	decl.Append(ifDecl)

	if !p.is(core.TokenIDNot) {
		return nil, nil
	}

	notDecl, err := p.consumeToken(core.TokenIDNot)
	if err != nil {
		return nil, err
	}
	ifDecl.Append(notDecl)

	if !p.is(core.TokenIDExists) {
		return nil, p.syntaxError()
	}

	existsDecl, err := p.consumeToken(core.TokenIDExists)
	if err != nil {
		return nil, err
	}
	notDecl.Append(existsDecl)

	return decl, nil
}

// parsePrimaryKey parses 'primary key' tokens.
func (p *Parser) parsePrimaryKey() (*core.Decl, error) {
	primaryDecl, err := p.consumeToken(core.TokenIDPrimary)
	if err != nil {
		return nil, err
	}

	keyDecl, err := p.consumeToken(core.TokenIDKey)
	if err != nil {
		return nil, err
	}
	primaryDecl.Append(keyDecl)

	if _, err = p.consumeToken(core.TokenIDBracketOpening); err != nil {
		return nil, err
	}

	for {
		d, err := p.parseQuotedToken()
		if err != nil {
			return nil, err
		}

		if d, err = p.consumeToken(core.TokenIDComma, core.TokenIDBracketClosing); err != nil {
			return nil, err
		}
		if d.TokenID == core.TokenIDBracketClosing {
			break
		}
	}
	return primaryDecl, nil
}

// parseDefaultClause parses 'default' tokens.
func (p *Parser) parseDefaultClause() (*core.Decl, error) {
	dDecl, err := p.consumeToken(core.TokenIDDefault)
	if err != nil {
		return nil, err
	}

	var vDecl *core.Decl
	if p.is(core.TokenIDSingleQuote) || p.is(core.TokenIDDoubleQuote) {
		vDecl, err = p.parseStringLiteral()
	} else {
		vDecl, err = p.consumeToken(core.TokenIDFalse, core.TokenIDNumber, core.TokenIDLocalTimestamp, core.TokenIDNow)
	}
	if err != nil {
		return nil, err
	}

	dDecl.Append(vDecl)
	return dDecl, nil
}
