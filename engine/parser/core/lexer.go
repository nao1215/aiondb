package core

// Instruction is a single Instruction in the query.
type Instruction struct {
	// Content is the Content of the instruction.
	Content []byte
	// Length is the Length of the instruction.
	Length uint64
}

// newInstruction creates a new instruction.
func newInstruction(input string) *Instruction {
	return &Instruction{
		Content: []byte(input),
		Length:  uint64(len(input)),
	}
}

// Position is the Position of the instruction.
type Position struct {
	// Current is the Current position of the instruction.
	Current uint64
	// Security is the positon that provides information to detect syntax errors.
	Security uint64
}

// IsSyntaxErr is whether a syntax error occurred during lexical analysis.
// Call after reading all matcher processes.
func (p *Position) IsSyntaxErr() bool {
	if p.Current == 0 && p.Security == 0 {
		return false
	}
	return p.Current == p.Security
}

// newPosition creates a new position.
func newPosition() *Position {
	return &Position{
		Current:  0,
		Security: 0,
	}
}

// Lex is information used during lexical analysis.
type Lex struct {
	// Tokens is the Tokens of the query.
	Tokens []Token
	// instructions is the instructions of the query.
	Instruction *Instruction
	// Position is the Position of the instruction.
	Position *Position
}

// NewLex creates a new lexer.
func NewLex(input string) *Lex {
	return &Lex{
		Tokens:      []Token{},
		Instruction: newInstruction(input),
		Position:    newPosition(),
	}
}

// Matcher tries to match given string to an SQL token
type Matcher func() bool

// Matchers is a list of matchers
type Matchers []Matcher
