package engine

// Tuple is a row in a relation
type Tuple struct {
	// Values is the list of values of the tuple
	Values []interface{}
}

// NewTuple should check that value are for the right Attribute and match domain.
func NewTuple(values ...interface{}) *Tuple {
	t := &Tuple{}
	t.Values = append(t.Values, values...)
	return t
}

// Append add a value to the tuple
func (t *Tuple) Append(value interface{}) {
	t.Values = append(t.Values, value)
}
