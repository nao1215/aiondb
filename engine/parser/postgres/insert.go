package postgres

import "github.com/nao1215/aiondb/engine/parser/core"

// parseInsert parses an INSERT statement.
//
// The generated AST is as follows:
//
//	|-> "INSERT" (InsertToken)
//	    |-> "INTO" (IntoToken)
//	        |-> table name
//	            |-> column name
//	            |-> (...)
//	    |-> "VALUES" (ValuesToken)
//	        |-> "(" (BracketOpeningToken)
//	            |-> value
//	            |-> (...)
//	        |-> (...)
//	    |-> "RETURNING" (ReturningToken) (optional)
//	        |-> column name
func (p *Parser) parseInsert() (*core.Statement, error) {
	stmt := &core.Statement{}

	// Set INSERT decl
	insertDecl, err := p.consumeToken(core.TokenIDInsert)
	if err != nil {
		return nil, err
	}
	stmt.Decls = append(stmt.Decls, insertDecl)

	// should be INTO
	intoDecl, err := p.consumeToken(core.TokenIDInto)
	if err != nil {
		return nil, err
	}
	insertDecl.Append(intoDecl)

	// should be table Name
	tableDecl, err := p.parseQuotedToken()
	if err != nil {
		return nil, err
	}
	intoDecl.Append(tableDecl)

	_, err = p.consumeToken(core.TokenIDBracketOpening)
	if err != nil {
		return nil, err
	}

	// concerned attribute
	for {
		decl, err := p.parseListElement()
		if err != nil {
			return nil, err
		}
		tableDecl.Append(decl)

		if p.is(core.TokenIDBracketClosing) {
			if _, err = p.consumeToken(core.TokenIDBracketClosing); err != nil {
				return nil, err
			}

			break
		}

		_, err = p.consumeToken(core.TokenIDComma)
		if err != nil {
			return nil, err
		}
	}

	// should be VALUES
	valuesDecl, err := p.consumeToken(core.TokenIDValues)
	if err != nil {
		return nil, err
	}
	insertDecl.Append(valuesDecl)

	for {
		openingBracketDecl, err := p.consumeToken(core.TokenIDBracketOpening)
		if err != nil {
			return nil, err
		}
		valuesDecl.Append(openingBracketDecl)

		// should be a list of values for specified attributes
		for {
			decl, err := p.parseListElement()
			if err != nil {
				return nil, err
			}
			openingBracketDecl.Append(decl)

			if p.is(core.TokenIDBracketClosing) {
				if _, err := p.consumeToken(core.TokenIDBracketClosing); err != nil {
					return nil, err
				}
				break
			}

			_, err = p.consumeToken(core.TokenIDComma)
			if err != nil {
				return nil, err
			}
		}

		if p.is(core.TokenIDComma) {
			if _, err := p.consumeToken(core.TokenIDComma); err != nil {
				return nil, err
			}
			continue
		}
		break
	}

	// we may have `returning "something"` here
	if retDecl, err := p.consumeToken(core.TokenIDReturning); err == nil {
		insertDecl.Append(retDecl)

		// returned attribute
		attrDecl, err := p.parseAttribute()
		if err != nil {
			return nil, err
		}
		retDecl.Append(attrDecl)
	}
	return stmt, nil
}
