package engine

import (
	"fmt"
	"strings"
	"time"

	"github.com/nao1215/aiondb/engine/parser/core"
)

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

// NewAttribute initialize a new Attribute struct
func NewAttribute(name string, typeName string, autoIncrement bool) Attribute {
	a := Attribute{
		name:          name,
		typeName:      typeName,
		autoIncrement: autoIncrement,
	}
	return a
}

// parseAttribute parses a declaration and returns an Attribute
func parseAttribute(decl *core.Decl) (Attribute, error) {
	attr := Attribute{}

	// Attribute name
	if decl.TokenID != core.TokenIDString {
		return attr, fmt.Errorf("engine: expected attribute name, got %v", decl.TokenID)
	}
	attr.name = decl.Lexeme.String()

	// Attribute type
	if len(decl.DeclList) < 1 {
		return attr, fmt.Errorf("attribute %s has no type", decl.Lexeme.String())
	}
	if decl.DeclList[0].TokenID != core.TokenIDString {
		return attr, fmt.Errorf("engine: expected attribute type, got %v:%v", decl.DeclList[0].TokenID, decl.DeclList[0].Lexeme)
	}
	attr.typeName = decl.DeclList[0].Lexeme.String()

	// Maybe domain and special thing like primary key
	typeDecl := decl.DeclList[1:]
	for i := range typeDecl {
		if typeDecl[i].TokenID == core.TokenIDAutoincrement {
			attr.autoIncrement = true
		}

		if typeDecl[i].TokenID == core.TokenIDDefault {
			switch typeDecl[i].DeclList[0].TokenID {
			case core.TokenIDLocalTimestamp, core.TokenIDNow:
				attr.defaultValue = func() interface{} { return time.Now().Format(core.DateLongFormat) }
			default:
				attr.defaultValue = typeDecl[i].DeclList[0].Lexeme
			}
		}
		// Check if attribute is unique
		if typeDecl[i].TokenID == core.TokenIDUnique {
			attr.unique = true
		}
	}

	if strings.ToLower(attr.typeName) == "bigserial" {
		attr.autoIncrement = true
	}
	return attr, nil
}
