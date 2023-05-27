package engine

import (
	"sync"
)

// Relation is a table with column and rows
// AKA File
type Relation struct {
	// mu is the mutex used to protect the rows slice.
	sync.RWMutex
	// table is the table of the relation
	table *Table
	// rows is the slice of rows of the relation
	rows []*Tuple
}

// NewRelation initializes a new Relation struct
func NewRelation(t *Table) *Relation {
	r := &Relation{
		table: t,
	}
	return r
}

// Insert a tuple in relation
func (r *Relation) Insert(t *Tuple) error {
	// Maybe do somthing like lock read/write here
	// Maybe index
	r.rows = append(r.rows, t)
	return nil
}
