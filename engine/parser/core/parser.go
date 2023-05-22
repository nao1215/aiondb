package core

import (
	"fmt"
	"os"
)

// Decl structure is the node to statement declaration tree
type Decl struct {
	// TokenID is token id
	TokenID TokenID
	// Lexeme is token lexeme
	Lexeme Lexeme
	// DeclList is the list of declaration
	DeclList []*Decl
}

// NewDecl initialize a Decl struct from a given token
func NewDecl(t Token) *Decl {
	return &Decl{
		TokenID: t.ID,
		Lexeme:  t.Lexeme,
	}
}

// Append creates a new leaf with given Decl
func (d *Decl) Append(subDecl *Decl) {
	d.DeclList = append(d.DeclList, subDecl)
}

// String prints the declaration tree in console
func (d *Decl) String(depth int) {
	indent := ""
	for i := 0; i < depth; i++ {
		indent = fmt.Sprintf("%s    ", indent)
	}

	fmt.Fprintf(os.Stdout, "%s|-> %s\n", indent, d.Lexeme)
	for _, v := range d.DeclList {
		v.String(depth + 1)
	}
}

// Statement define a valid SQL statement
type Statement struct {
	// Decls is the list of declaration
	Decls []*Decl
}

// PrettyPrint prints statement's declarations on console with indentation
func (s *Statement) PrettyPrint() {
	for _, d := range s.Decls {
		d.String(0)
	}
}
