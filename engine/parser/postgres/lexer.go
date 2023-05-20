package postgres

import (
	"fmt"
	"os"

	"github.com/nao1215/aiondb/engine/parser/core"
)

// Lexer is a lexer for PostgreSQL.
type Lexer struct {
	// Lex is information used during lexical analysis.
	lex      *core.Lex
	matchers *core.Matchers
}

// NewLexer returns a new Lexer.
func NewLexer(input string) *Lexer {
	return &Lexer{
		lex:      core.NewLex(input),
		matchers: newMatchers(),
	}
}

// Lex performs lexical analysis.
func (l *Lexer) Lex() ([]core.Token, error) {
	isMatch := false
	for l.lex.Position.Current < l.lex.Instruction.Length {
		for _, m := range *l.matchers {
			isMatch = false // reset. not delete it
			if isMatch = m(); isMatch == true {
				l.lex.Position.Security = l.lex.Position.Current
				break
			}
		}

		if isMatch {
			continue
		}

		if l.lex.Position.IsSyntaxErr() {
			fmt.Fprintf(os.Stderr, "Cannot lex <%s>, stuck at pos %d -> [%c]",
				l.lex.Instruction.Content, l.lex.Position.Current, l.lex.Instruction.Content[l.lex.Position.Current])
			return nil, core.Wrap(core.ErrLexerSyntaxErr,
				fmt.Sprintf("near %s", l.lex.Instruction.Content[l.lex.Position.Current:]))
		}
		l.lex.Position.Security = l.lex.Position.Current
	}
	return l.lex.Tokens, nil
}

// newMatchers returns a new Matchers.
func newMatchers() *core.Matchers {
	return &core.Matchers{}
}
