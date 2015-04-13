// Copyright 2013 The Go-SQLite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sqlite3

import "C"

import (
	"database/sql"
	"database/sql/driver"
	"io"
	"reflect"
	"time"
	"unsafe"
)

// Driver implements the interface required by database/sql.
type Driver string

func register(name string) {
	defer func() { recover() }()
	sql.Register(name, Driver(name))
}

func (Driver) Open(name string) (driver.Conn, error) {
	c, err := Open(name)
	if err != nil {
		return nil, err
	}
	c.BusyTimeout(5 * time.Second)
	return &conn{c}, nil
}

// conn implements driver.Conn.
type conn struct {
	*Conn
}

func (c *conn) Prepare(query string) (driver.Stmt, error) {
	if c.Conn.db == nil {
		return nil, driver.ErrBadConn
	}
	s, err := c.Conn.Prepare(query)
	if err != nil {
		return nil, err
	}
	return &stmt{s, false}, nil
}

func (c *conn) Begin() (driver.Tx, error) {
	if c.Conn.db == nil {
		return nil, driver.ErrBadConn
	}
	if err := c.Conn.Begin(); err != nil {
		return nil, err
	}
	return c.Conn, nil
}

func (c *conn) Exec(query string, args []driver.Value) (driver.Result, error) {
	if c.Conn.db == nil {
		return nil, driver.ErrBadConn
	}
	if err := c.Conn.Exec(query, vtoi(args)...); err != nil {
		return nil, err
	}
	// TODO: Do the driver.Result values need to be cached?
	return result{c.Conn}, nil
}

// stmt implements driver.Stmt.
type stmt struct {
	*Stmt
	closed bool
}

func (s *stmt) Close() error {
	if !s.closed {
		s.closed = true
		if !s.Stmt.Busy() {
			return s.Stmt.Close()
		}
	}
	return nil
}

func (s *stmt) NumInput() int {
	return s.Stmt.NumParams()
}

func (s *stmt) Exec(args []driver.Value) (driver.Result, error) {
	if err := s.Stmt.Exec(vtoi(args)...); err != nil {
		return nil, err
	}
	return result{s.Stmt.Conn()}, nil
}

func (s *stmt) Query(args []driver.Value) (driver.Rows, error) {
	if err := s.Stmt.Query(vtoi(args)...); err != nil && err != io.EOF {
		return nil, err
	}
	return &rows{s, true}, nil
}

// result implements driver.Result.
type result struct {
	*Conn
}

func (r result) LastInsertId() (int64, error) {
	return int64(r.Conn.LastInsertId()), nil
}

func (r result) RowsAffected() (int64, error) {
	return int64(r.Conn.RowsAffected()), nil
}

// rows implements driver.Rows.
type rows struct {
	*stmt
	first bool
}

func (r *rows) Close() error {
	if r.stmt.closed {
		return r.stmt.Stmt.Close()
	}
	r.stmt.Stmt.Reset()
	return nil
}

func (r *rows) Next(dest []driver.Value) error {
	if r.first {
		r.first = false
		if !r.stmt.Stmt.Busy() {
			return io.EOF
		}
	} else if err := r.stmt.Stmt.Next(); err != nil {
		return err
	}
	for i := range dest {
		v := (*interface{})(&dest[i])
		err := r.stmt.Stmt.scanDynamic(C.int(i), v, true)
		if err != nil {
			return err
		}
	}
	return nil
}

// vtoi converts []driver.Value to []interface{} without copying the contents.
func vtoi(v []driver.Value) (i []interface{}) {
	if len(v) > 0 {
		h := (*reflect.SliceHeader)(unsafe.Pointer(&i))
		h.Data = uintptr(unsafe.Pointer(&v[0]))
		h.Len = len(v)
		h.Cap = cap(v)
	}
	return
}
