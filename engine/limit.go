package engine

import (
	"github.com/nao1215/aiondb/engine/protocol"
)

// limit is a wrapper around protocol.EngineConn that limits the number of rows.
type limit struct {
	// realConn is the underlying connection
	realConn protocol.EngineConn
	// limit is the maximum number of rows to return
	limit int
	// current is the number of rows returned so far
	current int
}

// limitedConn returns a connection that limits the number of rows returned.
func limitedConn(conn protocol.EngineConn, l int) protocol.EngineConn {
	c := &limit{
		realConn: conn,
		limit:    l,
		current:  0,
	}
	return c
}

// ReadStatement reads a statement from the underlying connection.
// NOTE: This should not be used.
func (l *limit) ReadStatement() (string, error) {
	return "", nil
}

// WriteResult writes a result to the underlying connection.
// NOTE: This should not be used.
func (l *limit) WriteResult(_, _ int64) error {
	return nil
}

// WriteError writes an error to the underlying connection.
func (l *limit) WriteError(err error) error {
	return l.realConn.WriteError(err)
}

// WriteRowHeader writes a row header to the underlying connection.
func (l *limit) WriteRowHeader(header []string) error {
	return l.realConn.WriteRowHeader(header)
}

// WriteRow writes a row to the underlying connection.
func (l *limit) WriteRow(row []string) error {
	if l.current == l.limit {
		// We are done here
		return nil
	}
	l.current++
	return l.realConn.WriteRow(row)
}

// WriteRowEnd writes a row end to the underlying connection.
func (l *limit) WriteRowEnd() error {
	return l.realConn.WriteRowEnd()
}

// offset is a wrapper around protocol.EngineConn that skips the first N rows.
type offset struct {
	// realConn is the underlying connection
	realConn protocol.EngineConn
	// offset is the number of rows to skip
	offset int
	// current is the number of rows skipped so far
	current int
}

// offsetedConn returns a connection that skips the first N rows.
func offsetedConn(conn protocol.EngineConn, o int) protocol.EngineConn {
	c := &offset{
		realConn: conn,
		offset:   o,
	}
	return c
}

// ReadStatement reads a statement from the underlying connection.
// NOTE: This should not be used.
func (l *offset) ReadStatement() (string, error) {
	return "", nil
}

// WriteResult writes a result to the underlying connection.
// NOTE: This should not be used.
func (l *offset) WriteResult(_, _ int64) error {
	return nil
}

// WriteError writes an error to the underlying connection.
func (l *offset) WriteError(err error) error {
	return l.realConn.WriteError(err)
}

// WriteRowHeader writes a row header to the underlying connection.
func (l *offset) WriteRowHeader(header []string) error {
	return l.realConn.WriteRowHeader(header)
}

// WriteRow writes a row to the underlying connection.
func (l *offset) WriteRow(row []string) error {
	if l.current < l.offset {
		// skip this line
		l.current++
		return nil
	}
	return l.realConn.WriteRow(row)
}

// WriteRowEnd writes a row end to the underlying connection.
func (l *offset) WriteRowEnd() error {
	return l.realConn.WriteRowEnd()
}
