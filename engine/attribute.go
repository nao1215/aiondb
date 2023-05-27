package engine

// Domain is the set of allowable values for an Attribute.
type Domain struct{}

// Attribute is a named column of a relation
// AKA Field
// AKA Column
type Attribute struct {
	// name is the name of the attribute
	name string
	// typeName is the name of the type of the attribute
	typeName string
	// typeInstance is the instance of the type of the attribute
	typeInstance interface{}
	// defaultValue is the default value of the attribute
	defaultValue interface{}
	// domain is the set of allowable values for the attribute
	domain Domain
	// autoIncrement is true if the attribute is auto-incremented
	autoIncrement bool
	// unique is true if the attribute is unique
	unique bool
}
