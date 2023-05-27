// Package protocol defines the communication interface between the driver and the engine.
// It is used by the driver to send queries to the engine and by the engine to send results to the driver.
// The engine is the server and the driver is the client. The driver is the only one to initiate a connection.
package protocol

// DriverConn is a networking helper hiding implementation
// either with channels or network sockets.
type DriverConn interface {
	WriteQuery(query string) error
	WriteExec(stmt string) error
	ReadResult() (lastInsertedID int64, rowsAffected int64, err error)
	ReadRows() (chan []string, error)
	Close()
}

// EngineConn is a networking helper hiding implementation
// either with channels or network sockets.
type EngineConn interface {
	ReadStatement() (string, error)
	WriteResult(lastInsertedID int64, rowsAffected int64) error
	WriteError(err error) error
	WriteRowHeader(header []string) error
	WriteRow(row []string) error
	WriteRowEnd() error
}

// EngineEndpoint is the query entrypoint of RamSQL engine.
type EngineEndpoint interface {
	Accept() (EngineConn, error)
	Close()
}

// DriverEndpoint instanciates either
// - an Engine and communication channels
// - a network socket to connect to an existing RamSQL engine
type DriverEndpoint interface {
	New(string) (DriverConn, error)
}
