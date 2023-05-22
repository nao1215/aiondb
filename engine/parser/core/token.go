package core

// TokenID is the type of token ID.
type TokenID uint64

// Lexeme is minimal meaningful unit of language.
type Lexeme string

// String returns string representation of Lexeme.
func (l Lexeme) String() string {
	return string(l)
}

const (
	//=======================
	// Punctuation token
	//=======================

	// TokenIDSpace is the token ID for space.
	TokenIDSpace TokenID = 0
	// TokenIDSemicolon is the token ID for semicolon.
	TokenIDSemicolon TokenID = 1
	// TokenIDComma is the token ID for comma.
	TokenIDComma TokenID = 2
	// TokenIDBracketOpening is the token ID for bracket opening.
	TokenIDBracketOpening TokenID = 3
	// TokenIDBracketClosing is the token ID for bracket closing.
	TokenIDBracketClosing TokenID = 4
	// TokenIDLeftDiple is the token ID for left diple.
	TokenIDLeftDiple TokenID = 5
	// TokenIDRightDiple is the token ID for right diple.
	TokenIDRightDiple TokenID = 6
	// TokenIDLessOrEqual is the token ID for less or equal.
	TokenIDLessOrEqual TokenID = 7
	// TokenIDGreaterOrEqual is the token ID for greater or equal.
	TokenIDGreaterOrEqual TokenID = 8
	// TokenIDBacktick is the token ID for backtick.
	TokenIDBacktick TokenID = 9

	//=======================
	// Quote token
	//=======================

	// TokenIDDoubleQuote is the token ID for double quote.
	TokenIDDoubleQuote TokenID = 100
	// TokenIDSingleQuote is the token ID for single quote.
	TokenIDSingleQuote TokenID = 101
	// TokenIDStar is the token ID for star.
	TokenIDStar TokenID = 102
	// TokenIDEquality is the token ID for equal.
	TokenIDEquality TokenID = 103
	// TokenIDDistinctness is the token ID for distinctness.
	TokenIDDistinctness TokenID = 104
	// TokenIDPeriod is the token ID for period.
	TokenIDPeriod TokenID = 105

	//=======================
	// First order token
	//=======================

	// TokenIDCreate is the token ID for create.
	TokenIDCreate TokenID = 200
	// TokenIDSelect is the token ID for select.
	TokenIDSelect TokenID = 201
	// TokenIDInsert is the token ID for insert.
	TokenIDInsert TokenID = 202
	// TokenIDUpdate is the token ID for update.
	TokenIDUpdate TokenID = 203
	// TokenIDDelete is the token ID for delete.
	TokenIDDelete TokenID = 204
	// TokenIDExplain is the token ID for explain.
	TokenIDExplain TokenID = 205
	// TokenIDTruncate is the token ID for truncate.
	TokenIDTruncate TokenID = 206
	// TokenIDDrop is the token ID for drop.
	TokenIDDrop TokenID = 207
	// TokenIDGrant is the token ID for grant.
	TokenIDGrant TokenID = 208
	// TokenIDDistinct is the token ID for distinct.
	TokenIDDistinct TokenID = 209

	//=======================
	//  Second order token
	//=======================

	// TokenIDFrom is the token ID for all.
	TokenIDFrom TokenID = 300
	// TokenIDWhere is the token ID for where.
	TokenIDWhere TokenID = 301
	// TokenIDTable is the token ID for table.
	TokenIDTable TokenID = 302
	// TokenIDInto is the token ID for into.
	TokenIDInto TokenID = 303
	// TokenIDValues is the token ID for values.
	TokenIDValues TokenID = 304
	// TokenIDJoin is the token ID for join.
	TokenIDJoin TokenID = 305
	// TokenIDOn is the token ID for on.
	TokenIDOn TokenID = 306
	// TokenIDIf is the token ID for if.
	TokenIDIf TokenID = 307
	// TokenIDNot is the token ID for not.
	TokenIDNot TokenID = 308
	// TokenIDExists is the token ID for exists.
	TokenIDExists TokenID = 309
	// TokenIDNull is the token ID for null.
	TokenIDNull TokenID = 310
	// TokenIDAutoincrement is the token ID for autoincrement.
	TokenIDAutoincrement TokenID = 311
	// TokenIDCount is the token ID for count.
	TokenIDCount TokenID = 312
	// TokenIDSet is the token ID for set.
	TokenIDSet TokenID = 313
	// TokenIDOrder is the token ID for order.
	TokenIDOrder TokenID = 314
	// TokenIDBy is the token ID for by.
	TokenIDBy TokenID = 315
	// TokenIDWith is the token ID for with.
	TokenIDWith TokenID = 316
	// TokenIDTime is the token ID for time.
	TokenIDTime TokenID = 317
	// TokenIDZone is the token ID for zone.
	TokenIDZone TokenID = 318
	// TokenIDReturning is the token ID for returning.
	TokenIDReturning TokenID = 319
	// TokenIDIn is the token ID for in.
	TokenIDIn TokenID = 320
	// TokenIDAnd is the token ID for and.
	TokenIDAnd TokenID = 321
	// TokenIDOr is the token ID for or.
	TokenIDOr TokenID = 322
	// TokenIDAsc is the token ID for asc.
	TokenIDAsc TokenID = 323
	// TokenIDDesc is the token ID for desc.
	TokenIDDesc TokenID = 324
	// TokenIDLimit is the token ID for limit.
	TokenIDLimit TokenID = 325
	// TokenIDIs is the token ID for is.
	TokenIDIs TokenID = 326
	// TokenIDFor is the token ID for for.
	TokenIDFor TokenID = 327
	// TokenIDDefault is the token ID for default.
	TokenIDDefault TokenID = 328
	// TokenIDLocalTimestamp is the token ID for local timestamp.
	TokenIDLocalTimestamp TokenID = 329
	// TokenIDTrue is the token ID for true.
	TokenIDTrue TokenID = 330
	// TokenIDFalse is the token ID for false.
	TokenIDFalse TokenID = 331
	// TokenIDUnique is the token ID for unique.
	TokenIDUnique TokenID = 332
	// TokenIDNow is the token ID for now.
	TokenIDNow TokenID = 333
	// TokenIDOffset is the token ID for offset.
	TokenIDOffset TokenID = 334
	// TokenIDIndex is the token ID for index.
	TokenIDIndex TokenID = 335
	// TokenIDCollate is the token ID for collate.
	TokenIDCollate TokenID = 336
	// TokenIDNocase is the token ID for nocase.
	TokenIDNocase TokenID = 337

	//=======================
	//  Type token
	//=======================

	// TokenIDText is the token ID for text.
	TokenIDText TokenID = 400
	// TokenIDInt is the token ID for int.
	TokenIDInt TokenID = 401
	// TokenIDPrimary is the token ID for primary.
	TokenIDPrimary TokenID = 402
	// TokenIDKey is the token ID for key.
	TokenIDKey TokenID = 403
	// TokenIDString is the token ID for string.
	TokenIDString TokenID = 404
	// TokenIDNumber is the token ID for number.
	TokenIDNumber TokenID = 405
	// TokenIDDate is the token ID for date.
	TokenIDDate TokenID = 406
)

// Token in lexical analysis is the smallest unit
// of meaning in a language.
type Token struct {
	// ID is the token ID.
	ID TokenID
	// Lexeme is the token lexeme.
	Lexeme Lexeme
}

// StripSpaces strips spaces from tokens.
func StripSpaces(t []Token) (ret []Token) {
	for i := range t {
		if t[i].ID != TokenIDSpace {
			ret = append(ret, t[i])
		}
	}
	return ret
}
