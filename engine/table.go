package engine

// Table is defined by a name and attributes
// A table with data is called a Relation
type Table struct {
	name       string
	attributes []Attribute
}

// NewTable initializes a new Table
func NewTable(name string) *Table {
	t := &Table{
		name: name,
	}
	return t
}
