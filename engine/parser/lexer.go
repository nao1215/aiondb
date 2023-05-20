package parser

import (
	"github.com/nao1215/aiondb/engine/parser/core"
	"github.com/nao1215/aiondb/engine/parser/postgres"
)

// Lexer is an lex interface.
type Lexer interface {
	Lex() ([]core.Token, error)
}

// NewLexer returns a new Lexer.
func NewLexer(input string, mode core.SQLSyntaxMode) Lexer {
	switch mode {
	case core.SQLSyntaxModePostgreSQL:
		return postgres.NewLexer(input)
	}
	return postgres.NewLexer(input)
}
