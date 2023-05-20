package core

// instruction is a single instruction in the query.
type instruction struct {
	// content is the content of the instruction.
	content []byte
	// length is the length of the instruction.
	length uint64
}

// newInstruction creates a new instruction.
func newInstruction(input string) instruction {
	return instruction{
		content: []byte(input),
		length:  uint64(len(input)),
	}
}

// position is the position of the instruction.
type position uint64

// Lexer is a lexer for the sql query.
type Lexer struct {
	// tokens is the tokens of the query.
	tokens []Token
	// instructions is the instructions of the query.
	instruction instruction
	// position is the position of the instruction.
	position position
}

// NewLexer creates a new lexer.
func NewLexer(input string) *Lexer {
	return &Lexer{
		instruction: newInstruction(input),
	}
}
