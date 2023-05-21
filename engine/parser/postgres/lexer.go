package postgres

import (
	"fmt"
	"os"
	"unicode"

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
	lexer := &Lexer{
		lex: core.NewLex(input),
	}
	lexer.matchers = newMatchers(lexer)
	return lexer
}

// Lex performs lexical analysis.
// TODO: move to core
func (l *Lexer) Lex() ([]core.Token, error) {
	isMatch := false
	for l.Position() < l.InstructionLength() {
		for _, m := range *l.matchers {
			isMatch = false // reset. not delete it
			if isMatch = m(); isMatch == true {
				l.lex.Position.Security = l.Position()
				break
			}
		}

		if isMatch {
			continue
		}

		if l.lex.Position.IsSyntaxErr() {
			fmt.Fprintf(os.Stderr, "Cannot lex <%s>, stuck at pos %d -> [%c]",
				l.lex.Instruction.Content, l.Position(), l.lex.Instruction.Content[l.Position()])
			return nil, core.Wrap(core.ErrLexerSyntaxErr,
				fmt.Sprintf("near %s", l.lex.Instruction.Content[l.Position():]))
		}
		l.lex.Position.Security = l.Position()
	}
	return l.lex.Tokens, nil
}

// Position returns the current position of the lexer.
func (l *Lexer) Position() uint64 {
	return l.lex.Position.Current
}

// InstructionLength returns the length of the instruction.
func (l *Lexer) InstructionLength() uint64 {
	return l.lex.Instruction.Length
}

// newMatchers returns a new Matchers.
func newMatchers(l *Lexer) *core.Matchers {
	return &core.Matchers{
		l.matchSpaceToken,
		l.matchNowToken,
		l.matchUniqueToken,
		l.matchLocalTimestampToken,
		l.matchDefaultToken,
		l.matchTrueToken,
		l.matchFalseToken,
	}
}

// appendToken appends a token to the lexer.
// Now that the verification of the current position (character) is complete,
// the next position will check.
func (l *Lexer) appendToken(t core.Token) {
	l.lex.Tokens = append(l.lex.Tokens, t)
	l.lex.Position.Current++
}

// Match checks whether the argument str matches the SQL token specified in the argument.
// The argument str can be entered in either uppercase or lowercase.
func (l *Lexer) Match(str []byte, token core.TokenID) bool {
	if l.Position()+uint64(len(str))-1 > l.InstructionLength() {
		return false
	}

	for i := range str {
		if unicode.ToLower(rune(l.lex.Instruction.Content[int(l.Position())+i])) != unicode.ToLower(rune(str[i])) {
			return false
		}
	}

	// if next character is still a string, it means it doesn't match
	// ie: COUNT shoulnd match COUNTRY
	if l.InstructionLength() > l.Position()+uint64(len(str)) {
		if unicode.IsLetter(rune(l.lex.Instruction.Content[l.Position()+uint64(len(str))])) ||
			l.lex.Instruction.Content[l.Position()+uint64(len(str))] == '_' {
			return false
		}
	}
	l.appendToken(core.Token{ID: token, Lexeme: core.Lexeme(str)})
	return true
}

// matchSpaceToken checks whether it matches the space(e.g. " ") token.
func (l *Lexer) matchSpaceToken() bool {
	if !unicode.IsSpace(rune(l.lex.Instruction.Content[l.lex.Position.Current])) {
		return false
	}
	l.appendToken(core.Token{ID: core.TokenIDSpace, Lexeme: " "})
	return true
}

// matchNowToken checks whether it matches the now() token.
func (l *Lexer) matchNowToken() bool {
	return l.Match([]byte("now()"), core.TokenIDNow)
}

// matchUniqueToken checks whether it matches the unique token.
func (l *Lexer) matchUniqueToken() bool {
	return l.Match([]byte("unique"), core.TokenIDUnique)
}

// matchLocalTimestampToken checks whether it matches the localtimestamp token.
func (l *Lexer) matchLocalTimestampToken() bool {
	return l.Match([]byte("localtimestamp"), core.TokenIDLocalTimestamp)
}

// matchDefaultToken checks whether it matches the default token.
func (l *Lexer) matchDefaultToken() bool {
	return l.Match([]byte("default"), core.TokenIDDefault)
}

// matchTrueToken checks whether it matches the true token.
func (l *Lexer) matchTrueToken() bool {
	return l.Match([]byte("true"), core.TokenIDTrue)
}

// matchFlaseToken checks whether it matches the false token.
func (l *Lexer) matchFalseToken() bool {
	return l.Match([]byte("false"), core.TokenIDFalse)
}
