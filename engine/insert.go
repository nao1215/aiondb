package engine

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/nao1215/aiondb/engine/parser/core"
	"github.com/nao1215/aiondb/engine/protocol"
)

// insertIntoTableExecutor is the executor for INSERT INTO statements.
func insertIntoTableExecutor(e *Engine, insertDecl *core.Decl, conn protocol.EngineConn) error {
	// Get table and concerned attributes and write lock it
	intoDecl := insertDecl.DeclList[0]
	r, attributes, err := getRelation(e, intoDecl)
	if err != nil {
		return err
	}
	r.Lock()
	defer r.Unlock()

	// Check for RETURNING clause
	var returnedID string
	if len(insertDecl.DeclList) > 2 {
		for i := range insertDecl.DeclList {
			if insertDecl.DeclList[i].TokenID == core.TokenIDReturning {
				returningDecl := insertDecl.DeclList[i]
				returnedID = returningDecl.Lexeme.String()
				break
			}
		}
	}

	// Create a new tuple with values
	ids := []int64{}
	valuesDecl := insertDecl.DeclList[1]
	for _, valueListDecl := range valuesDecl.DeclList {
		// TODO handle all inserts atomically
		id, err := insert(r, attributes, valueListDecl.DeclList, returnedID)
		if err != nil {
			return err
		}
		ids = append(ids, id)
	}

	// if RETURNING decl is not present
	if returnedID != "" {
		if err := conn.WriteRowHeader([]string{returnedID}); err != nil {
			return err
		}
		for _, id := range ids {
			if err := conn.WriteRow([]string{fmt.Sprintf("%v", id)}); err != nil {
				return err
			}
		}
		return conn.WriteRowEnd()
	}
	return conn.WriteResult(ids[len(ids)-1], (int64)(len(ids)))
}

// getRelation returns the relation and the attributes of the table
func getRelation(e *Engine, intoDecl *core.Decl) (*Relation, []*core.Decl, error) {
	// Decl[0] is the table name
	r := e.relation(intoDecl.DeclList[0].Lexeme.String())
	if r == nil {
		return nil, nil, errors.New("table " + intoDecl.DeclList[0].Lexeme.String() + " does not exist")
	}

	for i := range intoDecl.DeclList[0].DeclList {
		if err := attributeExistsInTable(e, intoDecl.DeclList[0].DeclList[i].Lexeme.String(), intoDecl.DeclList[0].Lexeme.String()); err != nil {
			return nil, nil, err
		}
	}
	return r, intoDecl.DeclList[0].DeclList, nil
}

// insert inserts a new tuple in the relation
func insert(r *Relation, attributes []*core.Decl, values []*core.Decl, returnedID string) (int64, error) {
	var assigned bool
	var id int64
	var valuesindex int

	// Create tuple
	t := NewTuple()

	for attrindex, attr := range r.table.attributes {
		assigned = false
		for x, decl := range attributes {
			if attr.name == decl.Lexeme.String() && !attr.autoIncrement {
				// Before adding value in tuple, check it's not a builtin func or arithmetic operation
				switch values[x].TokenID {
				case core.TokenIDNow:
					t.Append(time.Now().Format(core.DateLongFormat))
				default:
					switch strings.ToLower(attr.typeName) {
					case "int64", "int":
						val, err := strconv.ParseInt(values[x].Lexeme.String(), 10, 64)
						if err != nil {
							return 0, err
						}
						t.Append(val)
					case "numeric", "decimal":
						val, err := strconv.ParseFloat(values[x].Lexeme.String(), 64)
						if err != nil {
							return 0, err
						}
						t.Append(val)
					default:
						t.Append(values[x].Lexeme)
					}
				}
				valuesindex = x
				assigned = true
				if returnedID == attr.name {
					var err error
					id, err = strconv.ParseInt(values[x].Lexeme.String(), 10, 64)
					if err != nil {
						return 0, err
					}
				}
			}
		}

		// If attribute is AUTO INCREMENT, compute it and assign it
		if attr.autoIncrement {
			assigned = true
			id = int64(len(r.rows) + 1)
			t.Append(id)
		}

		// Do we have a UNIQUE attribute ? if so
		if attr.unique {
			for i := range r.rows { // check all value already in relation (yup, no index tree)
				val, ok := r.rows[i].Values[attrindex].(string)
				if !ok {
					return 0, fmt.Errorf("failed to type assertion")
				}
				if val == string(values[valuesindex].Lexeme) {
					return 0, fmt.Errorf("unique constraint violation")
				}
			}
		}

		// If values was not explicitly given, set default value
		if !assigned {
			switch val := attr.defaultValue.(type) {
			case func() interface{}:
				t.Append(val)
			default:
				t.Append(attr.defaultValue)
			}
		}
	}

	// Insert tuple
	if err := r.Insert(t); err != nil {
		return 0, err
	}
	return id, nil
}
