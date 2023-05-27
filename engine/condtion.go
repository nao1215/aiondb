package engine

import (
	"fmt"

	"github.com/nao1215/aiondb/engine/parser/core"
	"github.com/nao1215/aiondb/engine/protocol"
)

// ifExecutor executes if statement
func ifExecutor(e *Engine, ifDecl *core.Decl, conn protocol.EngineConn) error {
	if len(ifDecl.DeclList) == 0 {
		return fmt.Errorf("malformed condition")
	}
	if e.opsExecutors[ifDecl.DeclList[0].TokenID] != nil {
		return e.opsExecutors[ifDecl.DeclList[0].TokenID](e, ifDecl.DeclList[0], conn)
	}
	return fmt.Errorf("error near %s, unknown keyword", ifDecl.DeclList[0].Lexeme.String())
}

// notExecutor executes not statement
// TODO: implement
func notExecutor(_ *Engine, _ *core.Decl, _ protocol.EngineConn) error {
	return nil
}

// existsExecutor executes exists statement
// TODO: implement
func existsExecutor(_ *Engine, _ *core.Decl, _ protocol.EngineConn) error {
	return nil
}
