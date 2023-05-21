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
			return nil, core.Wrap(core.ErrLexerSyntax,
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
		l.matchAscToken,
		l.matchDescToken,
		l.matchAndToken,
		l.matchOrToken,
		l.matchInToken,
		l.matchReturningToken,
		l.matchTruncateToken,
		l.matchDropToken,
		l.matchGrantToken,
		l.matchWithToken,
		l.matchTimeToken,
		l.matchZoneToken,
		l.matchIsToken,
		l.matchForToken,
		l.matchLimitToken,
		l.matchOrderToken,
		l.matchByToken,
		l.matchSetToken,
		l.matchUpdateToken,
		l.matchCreateToken,
		l.matchSelectToken,
		l.matchDistinctToken,
		l.matchInsertToken,
		l.matchWhereToken,
		l.matchFromToken,
		l.matchTableToken,
		l.matchNullToken,
		l.matchIfToken,
		l.matchNotToken,
		l.matchExistsToken,
		l.matchCountToken,
		l.matchDeleteToken,
		l.matchAutoIncrementToken,
		l.matchPrimaryToken,
		l.matchKeyToken,
		l.matchIntoToken,
		l.matchValuesToken,
		l.matchJoinToken,
		l.matchOnToken,
		l.matchOffsetToken,
		l.matchIndexToken,
		l.matchCollateToken,
		l.matchNocaseToken,
		l.matchSingleQuoteToken,
		l.matchDoubleQuoteToken,
		l.matchDateToken,
		l.matchEscapedStringToken,
		l.matchStringToken,
		l.matchNumberToken,
		l.matchSemicolonToken,
		l.matchPeriodToken,
		l.matchBracketOpeningToken,
		l.matchBracketClosingToken,
		l.matchCommaToken,
		l.matchStarToken,
		l.matchEqualityToken,
		l.matchDistinctnessToken,
		l.matchLeftDipleToken,
		l.matchRightDipleToken,
		l.matchLessOrEqualToken,
		l.matchGreaterOrEqualToken,
		l.matchBacktickToken,
	}
}
