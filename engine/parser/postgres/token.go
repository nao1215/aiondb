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

// match checks whether the argument str matches the SQL token specified in the argument.
// The argument str can be entered in either uppercase or lowercase.
func (l *Lexer) match(str []byte, token core.TokenID) bool {
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

// matchSingleChar checks whether it matches the single character token.
func (l *Lexer) matchSingleChar(char byte, token core.TokenID) bool {
	if l.Position() > l.InstructionLength() {
		return false
	}

	if l.lex.Instruction.Content[l.Position()] != char {
		return false
	}
	l.appendToken(core.Token{ID: token, Lexeme: core.Lexeme(char)})
	return true
}

// matchStringToken checks whether it matches the string token.
func (l *Lexer) matchStringToken() bool {
	i := l.Position()
	for i < l.InstructionLength() &&
		(unicode.IsLetter(rune(l.lex.Instruction.Content[i])) ||
			unicode.IsDigit(rune(l.lex.Instruction.Content[i])) ||
			l.lex.Instruction.Content[i] == '_' ||
			l.lex.Instruction.Content[i] == '@' /* || l.instruction[i] == '.'*/) {
		i++
	}

	if i == l.Position() {
		return false
	}

	l.appendToken(core.Token{
		ID:     core.TokenIDString,
		Lexeme: core.Lexeme(l.lex.Instruction.Content[l.Position():i]),
	})
	l.lex.Position.Current = i
	return true
}

// matchNumberToken checks whether it matches the number token.
func (l *Lexer) matchNumberToken() bool {
	i := l.Position()
	for i < l.InstructionLength() && unicode.IsDigit(rune(l.lex.Instruction.Content[i])) {
		i++
	}

	if i < l.InstructionLength() && l.lex.Instruction.Content[i] == '.' {
		i++
		for i < l.InstructionLength() && unicode.IsDigit(rune(l.lex.Instruction.Content[i])) {
			i++
		}
	}

	if i == l.Position() {
		return false
	}

	l.appendToken(core.Token{
		ID:     core.TokenIDNumber,
		Lexeme: core.Lexeme(l.lex.Instruction.Content[l.Position():i]),
	})
	l.lex.Position.Current = i
	return true
}

// 2015-09-10 14:03:09.444695269 +0200 CEST);
func (l *Lexer) matchDateToken() bool {
	i := l.Position()
	for i < l.InstructionLength() &&
		l.lex.Instruction.Content[i] != ',' &&
		l.lex.Instruction.Content[i] != ')' {
		i++
	}

	data := string(l.lex.Instruction.Content[l.Position():i])
	_, err := core.ParseDate(data)
	if err != nil {
		return false
	}
	l.appendToken(core.Token{ID: core.TokenIDDate, Lexeme: core.Lexeme(data)})
	l.lex.Position.Current = i

	return true
}

// matchSingleQuoteToken checks whether it matches the single quote token.
func (l *Lexer) matchSingleQuoteToken() bool {
	if l.lex.Instruction.Content[l.Position()] != '\'' {
		return false
	}

	l.appendToken(core.Token{ID: core.TokenIDSingleQuote, Lexeme: core.Lexeme("'")})

	if l.matchSingleQuotedStringToken() {
		l.appendToken(core.Token{ID: core.TokenIDSingleQuote, Lexeme: core.Lexeme("'")})
		return true
	}
	return true
}

// matchSingleQuotedStringToken checks whether it matches the single quoted string token.
func (l *Lexer) matchSingleQuotedStringToken() bool {
	i := l.Position()
	for i < l.InstructionLength() && l.lex.Instruction.Content[i] != '\'' {
		i++
	}

	l.appendToken(core.Token{
		ID:     core.TokenIDString,
		Lexeme: core.Lexeme(l.lex.Instruction.Content[l.Position():i]),
	})
	l.lex.Position.Current = i
	return true
}

// matchDoubleQuoteToken checks whether it matches the double quote token.
func (l *Lexer) matchDoubleQuoteToken() bool {
	if l.lex.Instruction.Content[l.Position()] != '"' {
		return false
	}
	l.appendToken(core.Token{ID: core.TokenIDDoubleQuote, Lexeme: "\""})

	if l.matchDoubleQuotedStringToken() {
		l.appendToken(core.Token{ID: core.TokenIDDoubleQuote, Lexeme: "\""})
		return true
	}
	return true
}

// matchDoubleQuotedStringToken checks whether it matches the double quoted string token.
func (l *Lexer) matchDoubleQuotedStringToken() bool {
	i := l.Position()
	for i < l.InstructionLength() && l.lex.Instruction.Content[i] != '"' {
		i++
	}
	l.appendToken(core.Token{
		ID:     core.TokenIDString,
		Lexeme: core.Lexeme(l.lex.Instruction.Content[l.Position():i]),
	})
	l.lex.Position.Current = i

	return true
}

// matchEscapedStringToken checks whether it matches the escaped string token.
func (l *Lexer) matchEscapedStringToken() bool {
	i := l.Position()
	if l.lex.Instruction.Content[i] != '$' || l.lex.Instruction.Content[i+1] != '$' {
		return false
	}
	i += 2

	for i+1 < l.InstructionLength() &&
		!(l.lex.Instruction.Content[i] == '$' &&
			l.lex.Instruction.Content[i+1] == '$') {
		i++
	}
	i++

	if i == l.InstructionLength() {
		return false
	}

	tokenID := core.TokenIDNumber
	escaped := l.lex.Instruction.Content[l.Position()+2 : i-1]

	for _, r := range escaped {
		if unicode.IsDigit(rune(r)) == false {
			tokenID = core.TokenIDString
		}
	}

	_, err := core.ParseDate(string(escaped))
	if err == nil {
		tokenID = core.TokenIDDate
	}

	l.appendToken(core.Token{ID: tokenID, Lexeme: core.Lexeme(escaped)})
	l.lex.Position.Current = i + 1

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
	return l.match([]byte("now()"), core.TokenIDNow)
}

// matchUniqueToken checks whether it matches the unique token.
func (l *Lexer) matchUniqueToken() bool {
	return l.match([]byte("unique"), core.TokenIDUnique)
}

// matchLocalTimestampToken checks whether it matches the localtimestamp token.
func (l *Lexer) matchLocalTimestampToken() bool {
	return l.match([]byte("localtimestamp"), core.TokenIDLocalTimestamp)
}

// matchDefaultToken checks whether it matches the default token.
func (l *Lexer) matchDefaultToken() bool {
	return l.match([]byte("default"), core.TokenIDDefault)
}

// matchTrueToken checks whether it matches the true token.
func (l *Lexer) matchTrueToken() bool {
	return l.match([]byte("true"), core.TokenIDTrue)
}

// matchFlaseToken checks whether it matches the false token.
func (l *Lexer) matchFalseToken() bool {
	return l.match([]byte("false"), core.TokenIDFalse)
}

// matchAscToken checks whether it matches the asc token.
func (l *Lexer) matchAscToken() bool {
	return l.match([]byte("asc"), core.TokenIDAsc)
}

// matchDescToken checks whether it matches the desc token.
func (l *Lexer) matchDescToken() bool {
	return l.match([]byte("desc"), core.TokenIDDesc)
}

// matchAndToken checks whether it matches the and token.
func (l *Lexer) matchAndToken() bool {
	return l.match([]byte("and"), core.TokenIDAnd)
}

// matchOrToken checks whether it matches the or token.
func (l *Lexer) matchOrToken() bool {
	return l.match([]byte("or"), core.TokenIDOr)
}

// matchInToken checks whether it matches the in token.
func (l *Lexer) matchInToken() bool {
	return l.match([]byte("in"), core.TokenIDIn)
}

// matchReturningToken checks whether it matches the returning token.
func (l *Lexer) matchReturningToken() bool {
	return l.match([]byte("returning"), core.TokenIDReturning)
}

// matchTruncateToken checks whether it matches the truncate token.
func (l *Lexer) matchTruncateToken() bool {
	return l.match([]byte("truncate"), core.TokenIDTruncate)
}

// matchDropToken checks whether it matches the drop token.
func (l *Lexer) matchDropToken() bool {
	return l.match([]byte("drop"), core.TokenIDDrop)
}

// matchGrantToken checks whether it matches the grant token.
func (l *Lexer) matchGrantToken() bool {
	return l.match([]byte("grant"), core.TokenIDGrant)
}

// matchWithToken checks whether it matches the with token.
func (l *Lexer) matchWithToken() bool {
	return l.match([]byte("with"), core.TokenIDWith)
}

// matchTimeToken checks whether it matches the time token.
func (l *Lexer) matchTimeToken() bool {
	return l.match([]byte("time"), core.TokenIDTime)
}

// matchZoneToken checks whether it matches the zone token.
func (l *Lexer) matchZoneToken() bool {
	return l.match([]byte("zone"), core.TokenIDZone)
}

// matchIsToken checks whether it matches the is token.
func (l *Lexer) matchIsToken() bool {
	return l.match([]byte("is"), core.TokenIDIs)
}

// matchForToken checks whether it matches the for token.
func (l *Lexer) matchForToken() bool {
	return l.match([]byte("for"), core.TokenIDFor)
}

// matchLimitToken checks whether it matches the limit token.
func (l *Lexer) matchLimitToken() bool {
	return l.match([]byte("limit"), core.TokenIDLimit)
}

// matchOrderToken checks whether it matches the order token.
func (l *Lexer) matchOrderToken() bool {
	return l.match([]byte("order"), core.TokenIDOrder)
}

// matchByToken checks whether it matches the by token.
func (l *Lexer) matchByToken() bool {
	return l.match([]byte("by"), core.TokenIDBy)
}

// matchSetToken checks whether it matches the set token.
func (l *Lexer) matchSetToken() bool {
	return l.match([]byte("set"), core.TokenIDSet)
}

// matchUpdateToken checks whether it matches the update token.
func (l *Lexer) matchUpdateToken() bool {
	return l.match([]byte("update"), core.TokenIDUpdate)
}

// matchCreateToken checks whether it matches the create token.
func (l *Lexer) matchCreateToken() bool {
	return l.match([]byte("create"), core.TokenIDCreate)
}

// matchSelectToken checks whether it matches the select token.
func (l *Lexer) matchSelectToken() bool {
	return l.match([]byte("select"), core.TokenIDSelect)
}

// matchDistinctToken checks whether it matches the distinct token.
func (l *Lexer) matchDistinctToken() bool {
	return l.match([]byte("distinct"), core.TokenIDDistinct)
}

// matchInsertToken	checks whether it matches the insert token.
func (l *Lexer) matchInsertToken() bool {
	return l.match([]byte("insert"), core.TokenIDInsert)
}

// matchWhereToken checks whether it matches the where token.
func (l *Lexer) matchWhereToken() bool {
	return l.match([]byte("where"), core.TokenIDWhere)
}

func (l *Lexer) matchFromToken() bool {
	return l.match([]byte("from"), core.TokenIDFrom)
}

// matchTableToken checks whether it matches the table token.
func (l *Lexer) matchTableToken() bool {
	return l.match([]byte("table"), core.TokenIDTable)
}

// matchNullToken checks whether it matches the null token.
func (l *Lexer) matchNullToken() bool {
	return l.match([]byte("null"), core.TokenIDNull)
}

// matchIfToken checks whether it matches the if token.
func (l *Lexer) matchIfToken() bool {
	return l.match([]byte("if"), core.TokenIDIf)
}

// matchNotToken checks whether it matches the not token.
func (l *Lexer) matchNotToken() bool {
	return l.match([]byte("not"), core.TokenIDNot)
}

// matchExistsToken checks whether it matches the exists token.
func (l *Lexer) matchExistsToken() bool {
	return l.match([]byte("exists"), core.TokenIDExists)
}

// matchCountToken checks whether it matches the count token.
func (l *Lexer) matchCountToken() bool {
	return l.match([]byte("count"), core.TokenIDCount)
}

// matchDeleteToken checks whether it matches the delete token.
func (l *Lexer) matchDeleteToken() bool {
	return l.match([]byte("delete"), core.TokenIDDelete)
}

// matchAutoIncrementToken checks whether it matches the auto_increment token.
func (l *Lexer) matchAutoIncrementToken() bool {
	if l.match([]byte("auto_increment"), core.TokenIDAutoincrement) {
		return true
	}
	return l.match([]byte("autoincrement"), core.TokenIDAutoincrement)
}

// matchPrimaryToken checks whether it matches the primary token.
func (l *Lexer) matchPrimaryToken() bool {
	return l.match([]byte("primary"), core.TokenIDPrimary)
}

// matchKeyToken checks whether it matches the key token.
func (l *Lexer) matchKeyToken() bool {
	return l.match([]byte("key"), core.TokenIDKey)
}

// matchIntoToken checks whether it matches the into token.
func (l *Lexer) matchIntoToken() bool {
	return l.match([]byte("into"), core.TokenIDInto)
}

// matchValuesToken checks whether it matches the values token.
func (l *Lexer) matchValuesToken() bool {
	return l.match([]byte("values"), core.TokenIDValues)
}

// matchJoinToken checks whether it matches the join token.
func (l *Lexer) matchJoinToken() bool {
	return l.match([]byte("join"), core.TokenIDJoin)
}

// matchOnToken checks whether it matches the on token.
func (l *Lexer) matchOnToken() bool {
	return l.match([]byte("on"), core.TokenIDOn)
}

// matchOffsetToken checks whether it matches the offset token.
func (l *Lexer) matchOffsetToken() bool {
	return l.match([]byte("offset"), core.TokenIDOffset)
}

// matchIndexToken checks whether it matches the index token.
func (l *Lexer) matchIndexToken() bool {
	return l.match([]byte("index"), core.TokenIDIndex)
}

// matchCollateToken checks whether it matches the collate token.
func (l *Lexer) matchCollateToken() bool {
	return l.match([]byte("collate"), core.TokenIDCollate)
}

// matchNocaseToken	checks whether it matches the nocase token.
func (l *Lexer) matchNocaseToken() bool {
	return l.match([]byte("nocase"), core.TokenIDNocase)
}

// matchSemicolonToken checks whether it matches the semicolon token.
func (l *Lexer) matchSemicolonToken() bool {
	return l.matchSingleChar(';', core.TokenIDSemicolon)
}

// matchPeriodToken checks whether it matches the period token.
func (l *Lexer) matchPeriodToken() bool {
	return l.matchSingleChar('.', core.TokenIDPeriod)
}

// matchBracketOpeningToken checks whether it matches the bracket opening token.
func (l *Lexer) matchBracketOpeningToken() bool {
	return l.matchSingleChar('(', core.TokenIDBracketOpening)
}

// matchBracketClosingToken checks whether it matches the bracket closing token.
func (l *Lexer) matchBracketClosingToken() bool {
	return l.matchSingleChar(')', core.TokenIDBracketClosing)
}

// matchCommaToken checks whether it matches the comma token.
func (l *Lexer) matchCommaToken() bool {
	return l.matchSingleChar(',', core.TokenIDComma)
}

// matchStarToken checks whether it matches the star token.
func (l *Lexer) matchStarToken() bool {
	return l.matchSingleChar('*', core.TokenIDStar)
}

// matchEqualityToken checks whether it matches the equality token.
func (l *Lexer) matchEqualityToken() bool {
	return l.matchSingleChar('=', core.TokenIDEquality)
}

// matchDistinctnessToken checks whether it matches the distinctness token.
func (l *Lexer) matchDistinctnessToken() bool {
	return l.match([]byte("<>"), core.TokenIDDistinctness)
}

// matchLeftDipleToken checks whether it matches the left diple token.
func (l *Lexer) matchLeftDipleToken() bool {
	return l.matchSingleChar('<', core.TokenIDLeftDiple)
}

// matchRightDipleToken checks whether it matches the right diple token.
func (l *Lexer) matchRightDipleToken() bool {
	return l.matchSingleChar('>', core.TokenIDRightDiple)
}

// matchLessOrEqualToken checks whether it matches the less or equal token.
func (l *Lexer) matchLessOrEqualToken() bool {
	return l.match([]byte("<="), core.TokenIDLessOrEqual)
}

// matchGreaterOrEqualToken checks whether it matches the greater or equal token.
func (l *Lexer) matchGreaterOrEqualToken() bool {
	return l.match([]byte(">="), core.TokenIDGreaterOrEqual)
}

// matchBacktickToken checks whether it matches the backtick token.
func (l *Lexer) matchBacktickToken() bool {
	return l.matchSingleChar('`', core.TokenIDBacktick)
}
