package postgres

import (
	"errors"

	"github.com/nao1215/aiondb/engine/parser/core"
)

// parseWhere parses the WHERE clause.
func (p *Parser) parseWhere(selectDecl *core.Decl) error {
	// May be WHERE  here
	// Can be ORDER BY if WHERE cause if implicit
	whereDecl, err := p.consumeToken(core.TokenIDWhere)
	if err != nil {
		return err
	}
	selectDecl.Append(whereDecl)

	// Now should be a list of: Attribute and Operator and Value
	gotClause := false
	for {
		if !p.hasNext() && gotClause {
			break
		}
		if p.is(core.TokenIDOrder, core.TokenIDLimit, core.TokenIDFor) {
			break
		}
		attributeDecl, err := p.parseCondition()
		if err != nil {
			return err
		}
		whereDecl.Append(attributeDecl)

		if p.is(core.TokenIDAnd, core.TokenIDAnd) {
			linkDecl, err := p.consumeToken(p.current().ID)
			if err != nil {
				return err
			}
			whereDecl.Append(linkDecl)
		}
		// Got at least one clause
		gotClause = true
	}
	return nil
}

// parseCondition
func (p *Parser) parseCondition() (*core.Decl, error) {
	// Optionnaly, brackets

	// We may have the WHERE 1 condition
	if t := p.current(); t.ID == core.TokenIDNumber && t.Lexeme == "1" {
		attributeDecl := core.NewDecl(t)
		if err := p.next(); err != nil {
			return nil, err
		}

		// in case of 1 = 1
		if p.current().ID == core.TokenIDEquality {
			t, err := p.isNext(core.TokenIDNumber)
			if err == nil && t.Lexeme == "1" {
				if _, err := p.consumeToken(core.TokenIDEquality); err != nil {
					return nil, err
				}
				if _, err := p.consumeToken(core.TokenIDNumber); err != nil {
					return nil, err
				}
			}
		}
		return attributeDecl, nil
	}

	// do we have brackets ?
	hasBracket := false
	if p.is(core.TokenIDBracketOpening) {
		if _, err := p.consumeToken(core.TokenIDBracketOpening); err != nil {
			return nil, err
		}
		hasBracket = true
	}

	// Attribute
	attributeDecl, err := p.parseAttribute()
	if err != nil {
		return nil, err
	}

	switch p.current().ID {
	case core.TokenIDEquality, core.TokenIDDistinctness, core.TokenIDLeftDiple, core.TokenIDRightDiple, core.TokenIDLessOrEqual, core.TokenIDGreaterOrEqual:
		decl, err := p.consumeToken(p.current().ID)
		if err != nil {
			return nil, err
		}
		attributeDecl.Append(decl)
	case core.TokenIDIn:
		inDecl, err := p.parseIn()
		if err != nil {
			return nil, err
		}
		attributeDecl.Append(inDecl)
		return attributeDecl, nil
	case core.TokenIDNot:
		notDecl, err := p.consumeToken(p.current().ID)
		if err != nil {
			return nil, err
		}

		if p.current().ID != core.TokenIDIn {
			return nil, errors.New("expected IN after NOT")
		}

		inDecl, err := p.parseIn()
		if err != nil {
			return nil, err
		}
		notDecl.Append(inDecl)
		attributeDecl.Append(notDecl)
		return attributeDecl, nil
	case core.TokenIDIs:
		decl, err := p.consumeToken(core.TokenIDIs)
		if err != nil {
			return nil, err
		}
		attributeDecl.Append(decl)
		if p.current().ID == core.TokenIDNot {
			notDecl, err := p.consumeToken(core.TokenIDNot)
			if err != nil {
				return nil, err
			}
			decl.Append(notDecl)
		}
		if p.current().ID == core.TokenIDNull {
			nullDecl, err := p.consumeToken(core.TokenIDNull)
			if err != nil {
				return nil, err
			}
			decl.Append(nullDecl)
		}
		return attributeDecl, nil
	default:
	}

	// Value
	valueDecl, err := p.parseValue()
	if err != nil {
		return nil, err
	}
	attributeDecl.Append(valueDecl)

	if hasBracket {
		if _, err = p.consumeToken(core.TokenIDBracketClosing); err != nil {
			return nil, err
		}
	}
	return attributeDecl, nil
}
