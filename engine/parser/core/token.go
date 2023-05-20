package core

// TokenID is the type of token ID.
type TokenID uint64

// Lexeme is minimal meaningful unit of language.
type Lexeme string

// Token in lexical analysis is the smallest unit
// of meaning in a language.
type Token struct {
	// ID is the token ID.
	ID TokenID
	// Lexeme is the token lexeme.
	Lexeme Lexeme
}
