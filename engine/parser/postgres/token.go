package postgres

import (
	"unicode"

	"github.com/nao1215/aiondb/engine/parser/core"
)

// appendToken appends a token to the lexer.
// Now that the verification of the current position (character) is complete,
// the next position will check.
func (l *Lexer) appendToken(t core.Token) {
	l.lex.Tokens = append(l.lex.Tokens, t)
	l.lex.Position.Current += uint64(len(t.Lexeme))
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

// matchAscToken checks whether it matches the asc token.
func (l *Lexer) matchAscToken() bool {
	return l.Match([]byte("asc"), core.TokenIDAsc)
}

// matchDescToken checks whether it matches the desc token.
func (l *Lexer) matchDescToken() bool {
	return l.Match([]byte("desc"), core.TokenIDDesc)
}

// matchAndToken checks whether it matches the and token.
func (l *Lexer) matchAndToken() bool {
	return l.Match([]byte("and"), core.TokenIDAnd)
}

// matchOrToken checks whether it matches the or token.
func (l *Lexer) matchOrToken() bool {
	return l.Match([]byte("or"), core.TokenIDOr)
}

// matchInToken checks whether it matches the in token.
func (l *Lexer) matchInToken() bool {
	return l.Match([]byte("in"), core.TokenIDIn)
}

// matchReturningToken checks whether it matches the returning token.
func (l *Lexer) matchReturningToken() bool {
	return l.Match([]byte("returning"), core.TokenIDReturning)
}

// matchTruncateToken checks whether it matches the truncate token.
func (l *Lexer) matchTruncateToken() bool {
	return l.Match([]byte("truncate"), core.TokenIDTruncate)
}

// matchDropToken checks whether it matches the drop token.
func (l *Lexer) matchDropToken() bool {
	return l.Match([]byte("drop"), core.TokenIDDrop)
}

// matchGrantToken checks whether it matches the grant token.
func (l *Lexer) matchGrantToken() bool {
	return l.Match([]byte("grant"), core.TokenIDGrant)
}

// matchWithToken checks whether it matches the with token.
func (l *Lexer) matchWithToken() bool {
	return l.Match([]byte("with"), core.TokenIDWith)
}

// matchTimeToken checks whether it matches the time token.
func (l *Lexer) matchTimeToken() bool {
	return l.Match([]byte("time"), core.TokenIDTime)
}

// matchZoneToken checks whether it matches the zone token.
func (l *Lexer) matchZoneToken() bool {
	return l.Match([]byte("zone"), core.TokenIDZone)
}

// matchIsToken checks whether it matches the is token.
func (l *Lexer) matchIsToken() bool {
	return l.Match([]byte("is"), core.TokenIDIs)
}

// matchForToken checks whether it matches the for token.
func (l *Lexer) matchForToken() bool {
	return l.Match([]byte("for"), core.TokenIDFor)
}

// matchLimitToken checks whether it matches the limit token.
func (l *Lexer) matchLimitToken() bool {
	return l.Match([]byte("limit"), core.TokenIDLimit)
}

// matchOrderToken checks whether it matches the order token.
func (l *Lexer) matchOrderToken() bool {
	return l.Match([]byte("order"), core.TokenIDOrder)
}

// matchByToken checks whether it matches the by token.
func (l *Lexer) matchByToken() bool {
	return l.Match([]byte("by"), core.TokenIDBy)
}

// matchSetToken checks whether it matches the set token.
func (l *Lexer) matchSetToken() bool {
	return l.Match([]byte("set"), core.TokenIDSet)
}

// matchUpdateToken checks whether it matches the update token.
func (l *Lexer) matchUpdateToken() bool {
	return l.Match([]byte("update"), core.TokenIDUpdate)
}

// matchCreateToken checks whether it matches the create token.
func (l *Lexer) matchCreateToken() bool {
	return l.Match([]byte("create"), core.TokenIDCreate)
}

// matchSelectToken checks whether it matches the select token.
func (l *Lexer) matchSelectToken() bool {
	return l.Match([]byte("select"), core.TokenIDSelect)
}

// matchDistinctToken checks whether it matches the distinct token.
func (l *Lexer) matchDistinctToken() bool {
	return l.Match([]byte("distinct"), core.TokenIDDistinct)
}

// matchInsertToken	checks whether it matches the insert token.
func (l *Lexer) matchInsertToken() bool {
	return l.Match([]byte("insert"), core.TokenIDInsert)
}

// matchWhereToken checks whether it matches the where token.
func (l *Lexer) matchWhereToken() bool {
	return l.Match([]byte("where"), core.TokenIDWhere)
}

func (l *Lexer) matchFromToken() bool {
	return l.Match([]byte("from"), core.TokenIDFrom)
}

// matchTableToken checks whether it matches the table token.
func (l *Lexer) matchTableToken() bool {
	return l.Match([]byte("table"), core.TokenIDTable)
}

// matchNullToken checks whether it matches the null token.
func (l *Lexer) matchNullToken() bool {
	return l.Match([]byte("null"), core.TokenIDNull)
}

// matchIfToken checks whether it matches the if token.
func (l *Lexer) matchIfToken() bool {
	return l.Match([]byte("if"), core.TokenIDIf)
}

// matchNotToken checks whether it matches the not token.
func (l *Lexer) matchNotToken() bool {
	return l.Match([]byte("not"), core.TokenIDNot)
}

// matchExistsToken checks whether it matches the exists token.
func (l *Lexer) matchExistsToken() bool {
	return l.Match([]byte("exists"), core.TokenIDExists)
}

// matchCountToken checks whether it matches the count token.
func (l *Lexer) matchCountToken() bool {
	return l.Match([]byte("count"), core.TokenIDCount)
}

// matchDeleteToken checks whether it matches the delete token.
func (l *Lexer) matchDeleteToken() bool {
	return l.Match([]byte("delete"), core.TokenIDDelete)
}

// matchAutoIncrementToken checks whether it matches the auto_increment token.
func (l *Lexer) matchAutoIncrementToken() bool {
	if l.Match([]byte("auto_increment"), core.TokenIDAutoincrement) {
		return true
	}
	return l.Match([]byte("autoincrement"), core.TokenIDAutoincrement)
}
