package engine

import (
	"fmt"

	"github.com/nao1215/aiondb/engine/parser/core"
)

// fromExecutor returns a slice of tables from a FROM declaration
func fromExecutor(fromDecl *core.Decl) []*Table {
	var tables []*Table
	for _, t := range fromDecl.DeclList {
		tables = append(tables, NewTable(t.Lexeme.String()))
	}
	return tables
}

// whereExecutor returns a slice of predicates from a WHERE declaration.
func whereExecutor(whereDecl *core.Decl, fromTableName string) ([]Predicate, error) {
	var predicates []Predicate
	var err error
	whereDecl.String(0)

	for i := range whereDecl.DeclList {
		var p Predicate
		tableName := fromTableName
		cond := whereDecl.DeclList[i]

		// 1 PREDICATE
		if cond.Lexeme == "1" {
			predicates = append(predicates, TruePredicate)
			continue
		}

		if len(cond.DeclList) == 0 {
			// log.Debug("whereExecutor: you must be AND or OR: %v", cond)
			continue
		}

		switch cond.DeclList[0].TokenID {
		case core.TokenIDEquality, core.TokenIDDistinctness, core.TokenIDLeftDiple, core.TokenIDRightDiple, core.TokenIDLessOrEqual, core.TokenIDGreaterOrEqual:
			// log.Debug("whereExecutor: it's = <> < > <= >=\n")
		case core.TokenIDIn:
			// log.Debug("whereExecutor: it's IN\n")
		case core.TokenIDNot:
			// log.Debug("whereExecutor: it's NOT\n")
		case core.TokenIDIs:
			// log.Debug("whereExecutor: it's IS token\n")
			// log.Debug("whereExecutor: %+v\n", cond.DeclList[0])
		default:
			// log.Debug("it's the table name ! -> %s", cond.DeclList[0].Lexeme)
			tableName = cond.DeclList[0].Lexeme.String()
			cond.DeclList = cond.DeclList[1:]
		}

		p.LeftValue.lexeme = whereDecl.DeclList[i].Lexeme.String()
		// Handle IN keyword
		if cond.DeclList[0].TokenID == core.TokenIDIn {
			err := inExecutor(cond.DeclList[0], &p)
			if err != nil {
				return nil, err
			}
			p.LeftValue.table = tableName
			predicates = append(predicates, p)
			continue
		}

		// Handle NOT IN keywords
		if cond.DeclList[0].TokenID == core.TokenIDNot && cond.DeclList[0].DeclList[0].TokenID == core.TokenIDIn {
			err := notInExecutor(cond.DeclList[0].DeclList[0], &p)
			if err != nil {
				return nil, err
			}
			p.LeftValue.table = tableName
			predicates = append(predicates, p)
			continue
		}

		// Handle IS NULL and IS NOT NULL
		if cond.DeclList[0].TokenID == core.TokenIDIs {
			err := isExecutor(cond.DeclList[0], &p)
			if err != nil {
				return nil, err
			}
			p.LeftValue.table = tableName
			predicates = append(predicates, p)
			continue
		}

		if len(cond.DeclList) < 2 {
			return nil, fmt.Errorf("malformed predicate \"%s\"", cond.Lexeme)
		}

		// The first element of the list is then the relation of the attribute
		op := cond.DeclList[0]
		val := cond.DeclList[1]

		p.Operator, err = NewOperator(op.TokenID, op.Lexeme.String())
		if err != nil {
			return nil, err
		}
		p.RightValue.lexeme = val.Lexeme.String()
		p.RightValue.valid = true
		p.LeftValue.table = tableName
		predicates = append(predicates, p)
	}

	if len(predicates) == 0 {
		return nil, fmt.Errorf("no predicates provided")
	}
	return predicates, nil
}

// inExecutor handles the IN operator
func inExecutor(inDecl *core.Decl, p *Predicate) error {
	inDecl.String(0)
	p.Operator = inOperator

	// Put everything in a []string
	var values []string
	for i := range inDecl.DeclList {
		values = append(values, inDecl.DeclList[i].Lexeme.String())
	}
	p.RightValue.v = values

	return nil
}

// notInExecutor handles the NOT IN operator
func notInExecutor(inDecl *core.Decl, p *Predicate) error {
	inDecl.String(0)
	p.Operator = notInOperator

	// Put everything in a []string
	var values []string
	for i := range inDecl.DeclList {
		values = append(values, inDecl.DeclList[i].Lexeme.String())
	}
	p.RightValue.v = values

	return nil
}

// isExecutor handles the IS operator
func isExecutor(isDecl *core.Decl, p *Predicate) error {
	isDecl.String(0)

	if isDecl.DeclList[0].TokenID == core.TokenIDNull {
		p.Operator = isNullOperator
	} else {
		p.Operator = isNotNullOperator
	}
	return nil
}
