package engine

import (
	"github.com/nao1215/aiondb/engine/protocol"
)

// distinct is a wrapper around a connection that removes duplicate rows
type distinct struct {
	// realConn is the real connection
	realConn protocol.EngineConn
	// seen is the map of all seen rows
	seen seen
	// len is the number of columns to compare
	len int
}

// distinctedConn returns a connection that removes duplicate rows
func distinctedConn(conn protocol.EngineConn, len int) protocol.EngineConn {
	return &distinct{
		realConn: conn,
		len:      len,
		seen:     make(seen),
	}
}

// ReadStatement returns the next statement. It is not needed for this wrapper
func (l *distinct) ReadStatement() (string, error) {
	return "", nil
}

// WriteResult writes the result of the statement. It is not needed for this wrapper
func (l *distinct) WriteResult(_ int64, _ int64) error {
	return nil
}

// WriteError writes an error.
func (l *distinct) WriteError(err error) error {
	return l.realConn.WriteError(err)
}

// WriteRowHeader writes the header of the row.
func (l *distinct) WriteRowHeader(header []string) error {
	if l.len > 0 {
		// Postgres returns only columns outside of DISTINCT ON
		return l.realConn.WriteRowHeader(header[l.len:])
	}
	return l.realConn.WriteRowHeader(header)
}

// WriteRow writes a row.
func (l *distinct) WriteRow(row []string) error {
	if l.len > 0 {
		if l.seen.exists(row[:l.len]) {
			return nil
		}
		// Postgres returns only columns outside of DISTINCT ON
		return l.realConn.WriteRow(row[l.len:])
	}
	if l.seen.exists(row) {
		return nil
	}
	return l.realConn.WriteRow(row)
}

// WriteRowEnd writes the end of the row.
func (l *distinct) WriteRowEnd() error {
	return l.realConn.WriteRowEnd()
}

// equalRows returns true if the two rows are equal
func (l *distinct) equalRows(a, b []string) bool {
	if l.len > 0 {
		if len(a) < l.len || len(b) < l.len {
			return false
		}

		for idx := 0; idx < l.len; idx++ {
			if a[idx] != b[idx] {
				return false
			}
		}
		return true
	}

	if len(a) != len(b) {
		return false
	}
	for idx := range a {
		if a[idx] != b[idx] {
			return false
		}
	}
	return true
}

// seen is a map of seen rows
type seen map[string]seen

// exists returns true if the row exists in the map
func (s seen) exists(r []string) bool {
	if c, ok := s[r[0]]; ok {
		if len(r) == 1 {
			return true
		}
		return c.exists(r[1:])
	}

	s[r[0]] = make(seen)
	if len(r) == 1 {
		return false
	}
	// does not exists, but we want to populate the tree fully
	return s[r[0]].exists(r[1:])
}
