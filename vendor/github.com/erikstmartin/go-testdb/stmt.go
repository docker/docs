package testdb

import (
	"database/sql/driver"
)

type stmt struct {
	queryFunc func(args []driver.Value) (driver.Rows, error)
	execFunc  func(args []driver.Value) (driver.Result, error)
}

func (s *stmt) Close() error {
	return nil
}

func (s *stmt) NumInput() int {
	// This prevents the sql package from validating the number of inputs
	return -1
}

func (s *stmt) Exec(args []driver.Value) (driver.Result, error) {
	return s.execFunc(args)
}

func (s *stmt) Query(args []driver.Value) (driver.Rows, error) {
	return s.queryFunc(args)
}
