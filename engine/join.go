package engine

// virtualRow is the resultset after FROM and JOIN transformations
// The key of the map is the lexeme (table.attribute) of the value (i.e: user.name)
type virtualRow map[string]Value
