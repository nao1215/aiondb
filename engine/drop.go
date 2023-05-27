package engine

import (
	"fmt"

	"github.com/nao1215/aiondb/engine/parser/core"
	"github.com/nao1215/aiondb/engine/protocol"
)

// dropExecutor executes drop statement
func dropExecutor(e *Engine, dropDecl *core.Decl, conn protocol.EngineConn) error {
	// Should have table token
	if dropDecl.DeclList == nil ||
		len(dropDecl.DeclList) != 1 ||
		dropDecl.DeclList[0].TokenID != core.TokenIDTable ||
		len(dropDecl.DeclList[0].DeclList) != 1 {
		return fmt.Errorf("unexpected drop arguments")
	}

	table := dropDecl.DeclList[0].DeclList[0].Lexeme.String()
	if r := e.relation(table); r == nil {
		return fmt.Errorf("relation '%s' not found", table)
	}
	e.drop(table)

	return conn.WriteResult(0, 1)
}
