package testdb

import (
	"database/sql/driver"
	"errors"
)

type conn struct {
	queries      map[string]query
	queryFunc    func(query string, args []driver.Value) (driver.Rows, error)
	execFunc     func(query string, args []driver.Value) (driver.Result, error)
	beginFunc    func() (driver.Tx, error)
	commitFunc   func() error
	rollbackFunc func() error
}

func newConn() *conn {
	return &conn{
		queries: make(map[string]query),
	}
}

func (c *conn) Prepare(query string) (driver.Stmt, error) {
	s := new(stmt)

	if c.queryFunc != nil {
		s.queryFunc = func(args []driver.Value) (driver.Rows, error) {
			return c.queryFunc(query, args)
		}
	}

	if c.execFunc != nil {
		s.execFunc = func(args []driver.Value) (driver.Result, error) {
			return c.execFunc(query, args)
		}
	}

	if q, ok := d.conn.queries[getQueryHash(query)]; ok {
		if s.queryFunc == nil && q.rows != nil {
			s.queryFunc = func(args []driver.Value) (driver.Rows, error) {
				if q.rows != nil {
					if rows, ok := q.rows.(*rows); ok {
						return rows.clone(), nil
					}
					return q.rows, nil
				}
				return nil, q.err
			}
		}

		if s.execFunc == nil && q.result != nil {
			s.execFunc = func(args []driver.Value) (driver.Result, error) {
				if q.result != nil {
					return q.result, nil
				}
				return nil, q.err
			}
		}
	}

	if s.queryFunc == nil && s.execFunc == nil {
		return new(stmt), errors.New("Query not stubbed: " + query)
	}

	return s, nil
}

func (*conn) Close() error {
	return nil
}

func (c *conn) Begin() (driver.Tx, error) {
	if c.beginFunc != nil {
		return c.beginFunc()
	}

	t := &Tx{}
	if c.commitFunc != nil {
		t.SetCommitFunc(c.commitFunc)
	}
	if c.rollbackFunc != nil {
		t.SetRollbackFunc(c.rollbackFunc)
	}

	return t, nil
}

func (c *conn) Query(query string, args []driver.Value) (driver.Rows, error) {
	if c.queryFunc != nil {
		return c.queryFunc(query, args)
	}
	if q, ok := d.conn.queries[getQueryHash(query)]; ok {
		if rows, ok := q.rows.(*rows); ok {
			return rows.clone(), q.err
		}
		return q.rows, q.err
	}
	return nil, errors.New("Query not stubbed: " + query)
}

func (c *conn) Exec(query string, args []driver.Value) (driver.Result, error) {
	if c.execFunc != nil {
		return c.execFunc(query, args)
	}

	if q, ok := d.conn.queries[getQueryHash(query)]; ok {
		if q.result != nil {
			return q.result, nil
		} else if q.err != nil {
			return nil, q.err
		}
	}

	return nil, errors.New("Exec call not stubbed: " + query)
}
