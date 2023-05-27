// Package engine contains the main engine of AION DB.
// It is responsible for listening for new connections, parsing the SQL statements and executing them.
// It also contains the relations and the operations executors.
package engine

import (
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/nao1215/aiondb/engine/parser"
	"github.com/nao1215/aiondb/engine/parser/core"
	"github.com/nao1215/aiondb/engine/protocol"
)

// executor is a function that executes a statement
type executor func(*Engine, *core.Decl, protocol.EngineConn) error

// Engine is the root struct of AION DB server.
// It contains the endpoint to accept connections from drivers,
// the relations and the operations executors.
type Engine struct {
	// endpoint is the entrypoint of the engine.
	endpoint protocol.EngineEndpoint
	// relations is the map of all relations.
	relations map[string]*Relation
	// opsExecutors is the map of all operations executors.
	opsExecutors map[core.TokenID]executor
	// stop is the channel used to stop the listening loop.
	// Any value send to this channel (through Engine.stop), will stop the listening loop
	stop chan bool
	// parser is the parser used to parse the SQL statements.
	parser parser.Parser
	// mu is the mutex used to protect the relations map.
	sync.Mutex
}

// New initialize a new AION DB server.
func New(endpoint protocol.EngineEndpoint) (e *Engine, err error) {
	e = &Engine{
		endpoint: endpoint,
	}

	e.stop = make(chan bool)
	e.opsExecutors = map[core.TokenID]executor{
		core.TokenIDCreate: createExecutor,
		// core.TokenIDTable:    createTableExecutor,
		// core.TokenIDSelect:   selectExecutor,
		// core.TokenIDInsert:   insertIntoTableExecutor,
		// core.TokenIDDelete:   deleteExecutor,
		// core.TokenIDUpdate:   updateExecutor,
		core.TokenIDIf:       ifExecutor,
		core.TokenIDNot:      notExecutor,
		core.TokenIDExists:   existsExecutor,
		core.TokenIDTruncate: truncateExecutor,
		core.TokenIDDrop:     dropExecutor,
		core.TokenIDGrant:    grantExecutor,
	}
	e.relations = make(map[string]*Relation)
	e.parser = parser.NewParser(core.SQLSyntaxModePostgreSQL)

	e.start()
	return
}

// start starts the listening loop.
func (e *Engine) start() {
	go e.listen()
}

// listen listens for new connections.
func (e *Engine) listen() {
	newConnectionChannel := make(chan protocol.EngineConn)

	// Accept new connections
	go func() {
		for {
			conn, err := e.endpoint.Accept()
			if err != nil {
				e.Stop()
				return
			}
			newConnectionChannel <- conn
		}
	}()

	// Handle new connections
	for {
		select {
		case conn := <-newConnectionChannel:
			go e.handleConnection(conn)
		case <-e.stop:
			e.endpoint.Close()
			return
		}
	}
}

// Stop stops the listening loop.
func (e *Engine) Stop() {
	if e.stop == nil {
		// already stopped
		return
	}

	go func() {
		e.stop <- true
		close(e.stop)
		e.stop = nil
	}()
}

// handleConnection handles a new connection.
func (e *Engine) handleConnection(conn protocol.EngineConn) {
	for {
		stmt, err := conn.ReadStatement()
		if errors.Is(err, io.EOF) {
			// TODO: close engine if there is no conn left
			return
		}
		if err != nil {
			return
		}

		stmtList, err := e.parser.Parse(stmt)
		if err != nil {
			// TODO: handle error
			conn.WriteError(err) //nolint
			continue
		}

		err = e.executeQueries(stmtList, conn)
		if err != nil {
			// TODO: handle error
			conn.WriteError(err) //nolint
			continue
		}
	}
}

// executeQuery executes a single query.
func (e *Engine) executeQueries(stmts []core.Statement, conn protocol.EngineConn) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("fatal error: %s", r)
			return
		}
	}()

	for _, v := range stmts {
		err = e.executeQuery(v, conn)
		if err != nil {
			return err
		}
	}
	return nil
}

// executeQuery executes a single query.
func (e *Engine) executeQuery(stmt core.Statement, conn protocol.EngineConn) error {
	if e.opsExecutors[stmt.Decls[0].TokenID] != nil {
		return e.opsExecutors[stmt.Decls[0].TokenID](e, stmt.Decls[0], conn)
	}
	return errors.New("not implemented")
}

// relation returns the relation with the given name.
func (e *Engine) relation(name string) *Relation {
	e.Lock()
	r := e.relations[name]
	e.Unlock()
	return r
}

func (e *Engine) drop(name string) {
	e.Lock()
	delete(e.relations, name)
	e.Unlock()
}

// createExecutor executes a CREATE statement.
func createExecutor(e *Engine, createDecl *core.Decl, conn protocol.EngineConn) error {
	if len(createDecl.DeclList) == 0 {
		return errors.New("parsing failed, no declaration after CREATE")
	}

	if e.opsExecutors[createDecl.DeclList[0].TokenID] != nil {
		return e.opsExecutors[createDecl.DeclList[0].TokenID](e, createDecl.DeclList[0], conn)
	}
	return errors.New("parsing failed, unknown token " + createDecl.DeclList[0].Lexeme.String())
}

// grantExecutor executes a GRANT statement.
func grantExecutor(_ *Engine, _ *core.Decl, conn protocol.EngineConn) error {
	return conn.WriteResult(0, 0)
}
