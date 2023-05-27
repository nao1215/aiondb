package engine

import (
	"fmt"

	"github.com/nao1215/aiondb/engine/parser/core"
	"github.com/nao1215/aiondb/engine/protocol"
)

// truncateExecutor executes truncate statement
func truncateExecutor(e *Engine, trDecl *core.Decl, conn protocol.EngineConn) error {
	// get tables to be deleted
	table := NewTable(trDecl.DeclList[0].Lexeme.String())
	return truncateTable(e, table, conn)
}

// truncateTable truncates table
func truncateTable(e *Engine, table *Table, conn protocol.EngineConn) error {
	var rowsDeleted int64

	// get relations and write lock them
	r := e.relation(table.name)
	if r == nil {
		return fmt.Errorf("table %v not found", table.name)
	}
	r.Lock()
	defer r.Unlock()

	if r.rows != nil {
		rowsDeleted = int64(len(r.rows))
	}
	r.rows = make([]*Tuple, 0)

	return conn.WriteResult(0, rowsDeleted)
}
