// Copyright 2013 The Go-SQLite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sqlite3_test

import (
	"bytes"
	"database/sql"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"
	"unsafe"

	. "code.google.com/p/go-sqlite/go1/sqlite3"
	_ "code.google.com/p/go-sqlite/go1/sqlite3/codec"
)

var key = flag.String("key", "", "codec key to use for all tests")

// skip, when set to true, causes all remaining tests to be skipped.
var skip = false

type T struct{ *testing.T }

func begin(t *testing.T) T {
	if skip {
		t.SkipNow()
	}
	return T{t}
}

func (t T) skipRestIfFailed() {
	skip = skip || t.Failed()
}

func (t T) open(name string) *Conn {
	codec := false
	if *key != "" && name == ":memory:" {
		name = t.tmpFile()
		codec = true
	}
	c, err := Open(name)
	if c == nil || err != nil {
		t.Fatalf(cl("Open(%q) unexpected error: %v"), name, err)
	}
	if codec {
		if err = c.Key("main", []byte(*key)); err != nil {
			t.Fatalf(cl("Key() unexpected error: %v"), err)
		}
	}
	return c
}

func (t T) close(c io.Closer) {
	if c != nil {
		if db, _ := c.(*Conn); db != nil {
			if path := db.Path("main"); path != "" {
				defer os.Remove(path)
			}
		}
		if err := c.Close(); err != nil {
			if !t.Failed() {
				t.Fatalf(cl("(%T).Close() unexpected error: %v"), c, err)
			}
			t.FailNow()
		}
	}
}

func (t T) prepare(c *Conn, sql string) *Stmt {
	s, err := c.Prepare(sql)
	if s == nil || err != nil {
		t.Fatalf(cl("c.Prepare(%q) unexpected error: %v"), sql, err)
	}
	return s
}

func (t T) query(cs interface{}, args ...interface{}) (s *Stmt) {
	var sql string
	var err error
	if c, ok := cs.(*Conn); ok {
		sql = args[0].(string)
		s, err = c.Query(sql, args[1:]...)
	} else {
		s = cs.(*Stmt)
		sql = s.String()
		err = s.Query(args...)
	}
	if s == nil || err != nil {
		t.Fatalf(cl("(%T).Query(%q) unexpected error: %v"), cs, sql, err)
	}
	return
}

func (t T) exec(cs interface{}, args ...interface{}) {
	var sql string
	var err error
	if c, ok := cs.(*Conn); ok {
		sql = args[0].(string)
		err = c.Exec(sql, args[1:]...)
	} else {
		s := cs.(*Stmt)
		sql = s.String()
		err = s.Exec(args...)
	}
	if err != nil {
		t.Fatalf(cl("(%T).Exec(%q) unexpected error: %v"), cs, sql, err)
	}
}

func (t T) scan(s *Stmt, dst ...interface{}) {
	if err := s.Scan(dst...); err != nil {
		t.Fatalf(cl("s.Scan() unexpected error: %v"), err)
	}
}

func (t T) next(s *Stmt, want error) {
	if have := s.Next(); have != want {
		if want == nil {
			t.Fatalf(cl("s.Next() unexpected error: %v"), have)
		} else {
			t.Fatalf(cl("s.Next() expected %v; got %v"), want, have)
		}
	}
}

func (t T) tmpFile() string {
	f, err := ioutil.TempFile("", "go-sqlite.db.")
	if err != nil {
		t.Fatalf(cl("tmpFile() unexpected error: %v"), err)
	}
	defer f.Close()
	return f.Name()
}

func (t T) errCode(have error, want int) {
	if e, ok := have.(*Error); !ok || e.Code() != want {
		t.Fatalf(cl("errCode() expected error code [%d]; got %v"), want, have)
	}
}

func cl(s string) string {
	_, thisFile, _, _ := runtime.Caller(1)
	_, testFile, line, ok := runtime.Caller(2)
	if ok && thisFile == testFile {
		return fmt.Sprintf("%d: %s", line, s)
	}
	return s
}

func sHdr(s string) *reflect.StringHeader {
	return (*reflect.StringHeader)(unsafe.Pointer(&s))
}

func bHdr(b []byte) *reflect.SliceHeader {
	return (*reflect.SliceHeader)(unsafe.Pointer(&b))
}

func TestLib(T *testing.T) {
	t := begin(T)
	defer t.skipRestIfFailed()

	if v, min := VersionNum(), 3007017; v < min {
		t.Errorf("VersionNum() expected >= %d; got %d", min, v)
	}
	if SingleThread() {
		t.Errorf("SingleThread() expected false")
	}

	sql := "CREATE TABLE x(a)"
	if Complete(sql) {
		t.Errorf("Complete(%q) expected false", sql)
	}
	if sql += ";"; !Complete(sql) {
		t.Errorf("Complete(%q) expected true", sql)
	}
}

func TestCreate(T *testing.T) {
	t := begin(T)
	defer t.skipRestIfFailed()

	checkPath := func(c *Conn, name, want string) {
		if have := c.Path(name); have != want {
			t.Fatalf(cl("c.Path() expected %q; got %q"), want, have)
		}
	}
	sql := "CREATE TABLE x(a); INSERT INTO x VALUES(1);"
	tmp := t.tmpFile()

	// File
	os.Remove(tmp)
	c := t.open(tmp)
	defer t.close(c)
	checkPath(c, "main", tmp)
	t.exec(c, sql)
	if err := c.Close(); err != nil {
		t.Fatalf("c.Close() unexpected error: %v", err)
	}
	if err := c.Exec(sql); err != ErrBadConn {
		t.Fatalf("c.Exec() expected %v; got %v", ErrBadConn, err)
	}

	// URI (existing)
	uri := strings.NewReplacer("?", "%3f", "#", "%23").Replace(tmp)
	if runtime.GOOS == "windows" {
		uri = "/" + strings.Replace(uri, "\\", "/", -1)
	}
	c = t.open("file:" + uri)
	defer t.close(c)
	checkPath(c, "main", tmp)
	t.exec(c, "INSERT INTO x VALUES(2)")

	// Temporary (in-memory)
	if *key == "" {
		c = t.open(":memory:")
		defer t.close(c)
		checkPath(c, "main", "")
		t.exec(c, sql)
	}

	// Temporary (file)
	c = t.open("")
	defer t.close(c)
	checkPath(c, "main", "")
	t.exec(c, sql)
}

func TestQuery(T *testing.T) {
	t := begin(T)
	defer t.skipRestIfFailed()

	sql := `
		CREATE TABLE x(a, b, c, d, e);
		INSERT INTO x VALUES(NULL, 123, 1.23, 'TEXT', x'424C4F42');
	`
	type row struct {
		a interface{}
		b int
		c float64
		d string
		e []byte
	}
	want := &row{nil, 123, 1.23, "TEXT", []byte("BLOB")}
	have := &row{}

	c := t.open(":memory:")
	defer t.close(c)
	t.exec(c, sql)

	s := t.query(c, "SELECT * FROM x")
	defer t.close(s)
	t.scan(s, &have.a, &have.b, &have.c, &have.d, &have.e)
	if !reflect.DeepEqual(have, want) {
		t.Errorf("s.Scan() expected %v; got %v", want, have)
	}

	t.next(s, io.EOF)
	t.close(s)
	if err := s.Query(); err != ErrBadStmt {
		t.Errorf("s.Query() expected %v; got %v", ErrBadStmt, err)
	}
}

func TestScan(T *testing.T) {
	t := begin(T)
	defer t.skipRestIfFailed()

	type types struct {
		v         interface{}
		int       int
		int64     int64
		float64   float64
		bool      bool
		string    string
		bytes     []byte
		Time      time.Time
		RawString RawString
		RawBytes  RawBytes
		Writer    io.Writer
	}
	scan := func(s *Stmt, dst ...interface{}) {
		t.query(s) // Re-query to avoid interference from type conversion
		t.scan(s, dst...)
	}
	skipCols := make([]interface{}, 0, 16)
	scanNext := func(s *Stmt, have, want *types) {
		scan(s, append(skipCols, &have.v)...)
		scan(s, append(skipCols, &have.int)...)
		scan(s, append(skipCols, &have.int64)...)
		scan(s, append(skipCols, &have.float64)...)
		scan(s, append(skipCols, &have.bool)...)
		scan(s, append(skipCols, &have.string)...)
		scan(s, append(skipCols, &have.bytes)...)
		scan(s, append(skipCols, &have.Time)...)
		scan(s, append(skipCols, have.Writer)...)

		// RawString must be copied, RawBytes (last access) can be used directly
		scan(s, append(skipCols, &have.RawString)...)
		have.RawString = RawString(have.RawString.Copy())
		scan(s, append(skipCols, &have.RawBytes)...)

		if !reflect.DeepEqual(have, want) {
			t.Fatalf(cl("scanNext() expected\n%#v; got\n%#v"), want, have)
		}
		skipCols = append(skipCols, nil)
	}

	c := t.open(":memory:")
	defer t.close(c)
	t.exec(c, `
		CREATE TABLE x(a, b, c, d, e, f, g, h, i);
		INSERT INTO x VALUES(NULL, '', x'', 0, 0.0, 4.2, 42, '42', x'3432');
	`)
	s := t.query(c, "SELECT * FROM x")
	defer t.close(s)

	// Verify data types
	wantT := []uint8{NULL, TEXT, BLOB, INTEGER, FLOAT, FLOAT, INTEGER, TEXT, BLOB}
	if haveT := s.DataTypes(); !reflect.DeepEqual(haveT, wantT) {
		t.Fatalf(cl("s.DataTypes() expected %v; got %v"), wantT, haveT)
	}

	// NULL
	have := &types{Writer: new(bytes.Buffer)}
	want := &types{Writer: new(bytes.Buffer)}
	scanNext(s, have, want)

	// ''
	want.v = ""
	want.Time = time.Unix(0, 0)
	scanNext(s, have, want)

	// x''
	want.v = []byte(nil)
	scanNext(s, have, want)

	// 0
	want.v = int64(0)
	want.string = "0"
	want.bytes = []byte("0")
	want.RawString = RawString("0")
	want.RawBytes = RawBytes("0")
	want.Writer.Write([]byte("0"))
	scanNext(s, have, want)

	// 0.0
	want.v = 0.0
	want.string = "0.0"
	want.bytes = []byte("0.0")
	want.RawString = RawString("0.0")
	want.RawBytes = RawBytes("0.0")
	want.Writer.Write([]byte("0.0"))
	scanNext(s, have, want)

	// 4.2
	want.v = 4.2
	want.int = 4
	want.int64 = 4
	want.float64 = 4.2
	want.bool = true
	want.string = "4.2"
	want.bytes = []byte("4.2")
	want.Time = time.Unix(4, 0)
	want.RawString = RawString("4.2")
	want.RawBytes = RawBytes("4.2")
	want.Writer.Write([]byte("4.2"))
	scanNext(s, have, want)

	// 42
	want.v = int64(42)
	want.int = 42
	want.int64 = 42
	want.float64 = 42
	want.string = "42"
	want.bytes = []byte("42")
	want.Time = time.Unix(42, 0)
	want.RawString = RawString("42")
	want.RawBytes = RawBytes("42")
	want.Writer.Write([]byte("42"))
	scanNext(s, have, want)

	// '42'
	want.v = "42"
	want.Writer.Write([]byte("42"))
	scanNext(s, have, want)

	// x'3432'
	want.v = []byte("42")
	want.Writer.Write([]byte("42"))
	scanNext(s, have, want)

	// Zero destinations
	t.scan(s)

	// Unsupported type
	var f32 float32
	t.errCode(s.Scan(&f32), MISUSE)

	// EOF
	t.next(s, io.EOF)
	if err := s.Scan(); err != io.EOF {
		t.Fatalf("s.Scan() expected EOF; got %v", err)
	}
}

func TestScanDynamic(T *testing.T) {
	t := begin(T)
	defer t.skipRestIfFailed()

	type row struct {
		a, b, c, d interface{}
		m          RowMap
	}
	scanNext := func(s *Stmt, have, want *row) {
		switch len(want.m) {
		case 0:
			t.scan(s, &have.a, &have.b, &have.c, &have.d)
		case 1:
			t.scan(s, &have.a, &have.b, &have.c, have.m)
		case 2:
			t.scan(s, &have.a, &have.b, have.m)
		case 3:
			t.scan(s, &have.a, have.m)
		case 4:
			t.scan(s, have.m)
		}
		if !reflect.DeepEqual(have, want) {
			t.Fatalf(cl("scanNext() expected\n%#v; got\n%#v"), want, have)
		}
		if err := s.Next(); err == io.EOF {
			t.query(s)
		} else if err != nil {
			t.Fatalf(cl("s.Next() unexpected error: %v"), err)
		}
	}

	c := t.open(":memory:")
	defer t.close(c)
	t.exec(c, `
		-- Affinity: NONE, NUMERIC, NUMERIC, NUMERIC
		CREATE TABLE x(a, b DATE, c time, d BOOLEAN);
		INSERT INTO x VALUES(NULL, NULL, NULL, NULL);
		INSERT INTO x VALUES('', '', '', '');
		INSERT INTO x VALUES(x'', x'', x'', x'');
		INSERT INTO x VALUES(0, 0, 0, 0);
		INSERT INTO x VALUES(0.0, 0.0, 0.0, 0.0);
		INSERT INTO x VALUES(4.2, 4.2, 4.2, 4.2);
		INSERT INTO x VALUES(42, 42, 42, 42);
		INSERT INTO x VALUES('42', '42', '42', '42');
		INSERT INTO x VALUES(x'3432', x'3432', x'3432', x'3432');
	`)
	s := t.query(c, "SELECT * FROM x ORDER BY rowid")
	defer t.close(s)

	// NULL
	have, want := &row{}, &row{}
	scanNext(s, have, want)

	// ''
	want = &row{"", "", "", "", nil}
	scanNext(s, have, want)

	// x''
	want = &row{[]byte(nil), []byte(nil), []byte(nil), []byte(nil), nil}
	scanNext(s, have, want)

	// 0
	t0 := time.Unix(0, 0)
	want = &row{int64(0), t0, t0, false, nil}
	scanNext(s, have, want)

	// 0.0
	want = &row{0.0, t0, t0, false, nil}
	scanNext(s, have, want)

	// 4.2
	want = &row{4.2, 4.2, 4.2, 4.2, nil}
	scanNext(s, have, want)

	// 42
	t42 := time.Unix(42, 0)
	want = &row{int64(42), t42, t42, true, nil}
	scanNext(s, have, want)

	// '42'
	want = &row{"42", t42, t42, true, nil}
	scanNext(s, have, want)

	// x'3432'
	want = &row{[]byte("42"), []byte("42"), []byte("42"), []byte("42"), nil}
	scanNext(s, have, want)

	// NULL (reset)
	have = &row{m: RowMap{}}
	want = &row{m: RowMap{"a": nil, "b": nil, "c": nil, "d": nil}}
	scanNext(s, have, want)

	// ''
	have = &row{m: RowMap{}}
	want = &row{m: RowMap{"a": "", "b": "", "c": "", "d": ""}}
	scanNext(s, have, want)

	// x''
	t.next(s, nil)

	// 0
	have = &row{m: RowMap{}}
	want = &row{m: RowMap{"a": int64(0), "b": t0, "c": t0, "d": false}}
	scanNext(s, have, want)

	// 0.0
	have = &row{m: RowMap{}}
	want = &row{a: 0.0, m: RowMap{"b": t0, "c": t0, "d": false}}
	scanNext(s, have, want)

	// 4.2
	have = &row{m: RowMap{}}
	want = &row{a: 4.2, b: 4.2, m: RowMap{"c": 4.2, "d": 4.2}}
	scanNext(s, have, want)

	// 42
	have = &row{m: RowMap{}}
	want = &row{a: int64(42), b: t42, c: t42, m: RowMap{"d": true}}
	scanNext(s, have, want)

	// Too many destinations
	t.errCode(s.Scan(&have.a, &have.b, &have.c, &have.d, &have.m), MISUSE)
}

func TestState(T *testing.T) {
	t := begin(T)
	defer t.skipRestIfFailed()

	type connState struct {
		AutoCommit        bool
		LastInsertId      int64
		RowsAffected      int
		TotalRowsAffected int
	}
	checkConn := func(c *Conn, want *connState) {
		have := &connState{
			c.AutoCommit(),
			c.LastInsertId(),
			c.RowsAffected(),
			c.TotalRowsAffected(),
		}
		if !reflect.DeepEqual(have, want) {
			t.Fatalf(cl("checkConn() expected\n%#v; got\n%#v"), want, have)
		}
	}

	c := t.open(":memory:")
	defer t.close(c)
	checkConn(c, &connState{true, 0, 0, 0})

	t.exec(c, "BEGIN")
	checkConn(c, &connState{false, 0, 0, 0})

	t.exec(c, "ROLLBACK")
	checkConn(c, &connState{true, 0, 0, 0})

	t.exec(c, "CREATE TABLE x(a, b INTEGER, c text)")
	checkConn(c, &connState{true, 0, 0, 0})

	t.exec(c, "INSERT INTO x VALUES(NULL, 42, 42)")
	checkConn(c, &connState{true, 1, 1, 1})

	t.exec(c, "INSERT INTO x VALUES(-42, 4.2, x'42')")
	checkConn(c, &connState{true, 2, 1, 2})

	type stmtState struct {
		Conn       *Conn
		Valid      bool
		Busy       bool
		ReadOnly   bool
		NumParams  int
		NumColumns int
		Params     []string
		Columns    []string
		DeclTypes  []string
		DataTypes  []uint8
	}
	checkStmt := func(s *Stmt, want *stmtState) {
		have := &stmtState{
			s.Conn(),
			s.Valid(),
			s.Busy(),
			s.ReadOnly(),
			s.NumParams(),
			s.NumColumns(),
			s.Params(),
			s.Columns(),
			s.DeclTypes(),
			s.DataTypes(),
		}
		if !reflect.DeepEqual(have, want) {
			t.Fatalf(cl("checkStmt() expected\n%#v; got\n%#v"), want, have)
		}
	}
	closed := &stmtState{Conn: c, ReadOnly: true}

	// SELECT
	s := t.prepare(c, "SELECT * FROM x ORDER BY rowid")
	defer t.close(s)
	want := &stmtState{
		Conn:       c,
		Valid:      true,
		ReadOnly:   true,
		NumColumns: 3,
		Columns:    []string{"a", "b", "c"},
		DeclTypes:  []string{"", "INTEGER", "TEXT"},
	}
	checkStmt(s, want)

	t.query(s)
	want.Busy = true
	want.DataTypes = []uint8{NULL, INTEGER, TEXT}
	checkStmt(s, want)

	t.query(s)
	var _a, _b, _c []byte
	t.scan(s, &_a, &_b, &_c) // This scan must not change DataTypes
	checkStmt(s, want)

	t.next(s, nil)
	want.DataTypes = []uint8{INTEGER, FLOAT, BLOB}
	checkStmt(s, want)

	t.next(s, io.EOF)
	want.Busy = false
	want.DataTypes = nil
	checkStmt(s, want)

	t.close(s)
	checkStmt(s, closed)

	// INSERT (unnamed parameters)
	s = t.prepare(c, "INSERT INTO x VALUES(?, ?, ?)")
	defer t.close(s)
	want = &stmtState{
		Conn:      c,
		Valid:     true,
		NumParams: 3,
	}
	checkStmt(s, want)

	t.exec(s, nil, nil, nil)
	checkStmt(s, want)

	t.close(s)
	checkStmt(s, closed)

	// INSERT (named parameters)
	s = t.prepare(c, "INSERT INTO x VALUES(:a, @B, $c)")
	defer t.close(s)
	want = &stmtState{
		Conn:      c,
		Valid:     true,
		NumParams: 3,
		Params:    []string{":a", "@B", "$c"},
	}
	checkStmt(s, want)

	// Comment
	s = t.prepare(c, "-- This is a comment")
	checkStmt(s, closed)
}

func TestTail(T *testing.T) {
	t := begin(T)
	defer t.skipRestIfFailed()

	c := t.open(":memory:")
	defer t.close(c)

	check := func(sql, tail string) {
		s := t.prepare(c, sql)
		defer t.close(s)
		head := sql[:len(sql)-len(tail)]
		tail = sql[len(head):]

		// s.String() must be a prefix of sql
		if s.String() != head {
			t.Errorf(cl("s.String() expected %q; got %q"), head, s.String())
		} else if sHdr(s.String()).Data != sHdr(sql).Data {
			t.Errorf(cl("s.String() isn't a pointer into sql"))
		}

		// s.Tail must be a suffix of sql or ""
		if s.Tail != tail {
			t.Errorf(cl("s.Tail expected %q; got %q"), tail, s.Tail)
		} else if tail == "" && sHdr(s.Tail).Data == sHdr(tail).Data {
			t.Errorf(cl("s.Tail is a pointer into sql"))
		} else if tail != "" && sHdr(s.Tail).Data != sHdr(tail).Data {
			t.Errorf(cl("s.Tail isn't a pointer into sql"))
		}
	}
	head := "CREATE TABLE x(a);"
	tail := " -- comment"

	check("", "")
	check(head, "")
	check(head+tail, tail)
	check(tail, "")
}

func TestParams(T *testing.T) {
	t := begin(T)
	defer t.skipRestIfFailed()

	c := t.open(":memory:")
	defer t.close(c)
	t.exec(c, "CREATE TABLE x(a, b, c, d)")

	dt := func(v interface{}) uint8 {
		switch v.(type) {
		case int64:
			return INTEGER
		case float64:
			return FLOAT
		case string:
			return TEXT
		case []byte:
			return BLOB
		}
		return NULL
	}
	verify := func(_a, _b, _c, _d interface{}) {
		s := t.query(c, "SELECT * FROM x ORDER BY rowid LIMIT 1")
		defer t.close(s)

		wantT := []uint8{dt(_a), dt(_b), dt(_c), dt(_d)}
		if haveT := s.DataTypes(); !reflect.DeepEqual(haveT, wantT) {
			t.Fatalf(cl("s.DataTypes() expected %v; got %v"), wantT, haveT)
		}

		type row struct{ a, b, c, d interface{} }
		want, have := &row{_a, _b, _c, _d}, &row{}
		t.scan(s, &have.a, &have.b, &have.c, &have.d)
		if !reflect.DeepEqual(have, want) {
			t.Fatalf(cl("verify() expected\n%#v; got\n%#v"), want, have)
		}
		t.exec(c, "DELETE FROM x WHERE rowid=(SELECT min(rowid) FROM x)")
	}

	// Unnamed
	sql := "INSERT INTO x VALUES(?, ?, ?, ?)"
	s := t.prepare(c, sql)
	defer t.close(s)

	t.exec(s, nil, nil, nil, nil)
	verify(nil, nil, nil, nil)

	t.exec(s, int(0), int(1), int64(math.MinInt64), int64(math.MaxInt64))
	verify(int64(0), int64(1), int64(math.MinInt64), int64(math.MaxInt64))

	t.exec(s, 0.0, 1.0, math.SmallestNonzeroFloat64, math.MaxFloat64)
	verify(0.0, 1.0, math.SmallestNonzeroFloat64, math.MaxFloat64)

	t.exec(s, false, true, "", "x\x00y")
	verify(int64(0), int64(1), "", "x\x00y")

	t.exec(s, []byte(nil), []byte{}, []byte{0}, []byte("1"))
	verify([]byte(nil), []byte(nil), []byte{0}, []byte("1"))

	t.exec(s, time.Unix(0, 0), time.Unix(1, 0), RawString(""), RawString("x"))
	verify(int64(0), int64(1), "", "x")

	t.exec(s, RawBytes(""), RawBytes("x"), ZeroBlob(0), ZeroBlob(2))
	verify([]byte(nil), []byte("x"), []byte(nil), []byte{0, 0})

	// Issue 1: string/[]byte zero values are not NULLs
	var s1, s2 string
	var b1, b2 []byte
	*sHdr(s1) = reflect.StringHeader{}
	*sHdr(s2) = reflect.StringHeader{Data: 1}
	*bHdr(b1) = reflect.SliceHeader{}
	*bHdr(b2) = reflect.SliceHeader{Data: 1}
	t.exec(s, s1, s2, b1, b2)
	verify("", "", []byte(nil), []byte(nil))

	// Invalid
	t.errCode(s.Exec(), MISUSE)
	t.errCode(s.Exec(0, 0, 0), MISUSE)
	t.errCode(s.Exec(0, 0, 0, 0, 0), MISUSE)
	t.errCode(s.Exec(NamedArgs{}), MISUSE)

	// Named
	s = t.prepare(c, "INSERT INTO x VALUES(:a, @B, :a, $d)")
	defer t.close(s)

	t.exec(s, NamedArgs{})
	verify(nil, nil, nil, nil)

	t.exec(s, 0, 1, 2)
	verify(int64(0), int64(1), int64(0), int64(2))

	t.exec(s, NamedArgs{":a": "a", "@B": "b", "$d": "d", "$c": nil})
	verify("a", "b", "a", "d")

	t.exec(s, NamedArgs{"@B": RawString("hello"), "$d": RawBytes("world")})
	verify(nil, "hello", nil, []byte("world"))

	// Invalid
	t.errCode(s.Exec(), MISUSE)
	t.errCode(s.Exec(0, 0), MISUSE)
	t.errCode(s.Exec(0, 0, 0, 0), MISUSE)
	t.errCode(s.Exec(NamedArgs(nil)), MISUSE)
	t.errCode(s.Exec(0, NamedArgs{}), MISUSE)
	t.errCode(s.Exec(0, 0, NamedArgs{}), MISUSE)

	// Conn.Query
	if s, err := c.Query(sql, 1, 2, 3, 4); s != nil || err != io.EOF {
		t.Fatalf("c.Query(%q) expected <nil>, EOF; got %v, %v", sql, s, err)
	}
	verify(int64(1), int64(2), int64(3), int64(4))

	// Conn.Exec
	t.exec(c, `
		INSERT INTO x VALUES(?, ?, NULL, NULL);
		INSERT INTO x VALUES(NULL, NULL, ?, ?);
	`, 1, 2, 3, 4)
	verify(int64(1), int64(2), nil, nil)
	verify(nil, nil, int64(3), int64(4))

	t.exec(c, `
		INSERT INTO x VALUES($a, $b, NULL, NULL);
		INSERT INTO x VALUES($a, $a, $c, $d);
	`, NamedArgs{"$a": "a", "$b": "b", "$c": "c", "$d": "d"})
	verify("a", "b", nil, nil)
	verify("a", "a", "c", "d")

	t.errCode(c.Exec(sql, 0, 0, 0), MISUSE)
	t.errCode(c.Exec(sql, 0, 0, 0, 0, 0), MISUSE)
}

func TestTx(T *testing.T) {
	t := begin(T)
	defer t.skipRestIfFailed()

	c := t.open(":memory:")
	defer t.close(c)
	t.exec(c, "CREATE TABLE x(a)")

	// Begin/Commit
	if err := c.Begin(); err != nil {
		t.Fatalf("c.Begin() unexpected error: %v", err)
	}
	t.exec(c, "INSERT INTO x VALUES(1)")
	t.exec(c, "INSERT INTO x VALUES(2)")
	if err := c.Commit(); err != nil {
		t.Fatalf("c.Commit() unexpected error: %v", err)
	}

	// Begin/Rollback
	if err := c.Begin(); err != nil {
		t.Fatalf("c.Begin() unexpected error: %v", err)
	}
	t.exec(c, "INSERT INTO x VALUES(3)")
	t.exec(c, "INSERT INTO x VALUES(4)")
	if err := c.Rollback(); err != nil {
		t.Fatalf("c.Rollback() unexpected error: %v", err)
	}

	// Verify
	s := t.query(c, "SELECT * FROM x ORDER BY rowid")
	defer t.close(s)
	var i int
	if t.scan(s, &i); i != 1 {
		t.Fatalf("s.Scan() expected 1; got %d", i)
	}
	t.next(s, nil)
	if t.scan(s, &i); i != 2 {
		t.Fatalf("s.Scan() expected 2; got %d", i)
	}
	t.next(s, io.EOF)
}

func TestIO(T *testing.T) {
	t := begin(T)

	c := t.open(":memory:")
	defer t.close(c)
	t.exec(c, "CREATE TABLE x(a)")
	t.exec(c, "INSERT INTO x VALUES(?)", ZeroBlob(8))
	t.exec(c, "INSERT INTO x VALUES(?)", "hello, world")

	// Open
	b, err := c.BlobIO("main", "x", "a", 1, true)
	if b == nil || err != nil {
		t.Fatalf("c.BlobIO() unexpected error: %v", err)
	}
	defer t.close(b)

	// State
	if b.Conn() != c {
		t.Fatalf("b.Conn() expected %p; got %p", c, b.Conn())
	}
	if b.Row() != 1 {
		t.Fatalf("b.Row() expected 1; got %d", b.Row())
	}
	if b.Len() != 8 {
		t.Fatalf("b.Len() expected 8; got %d", b.Len())
	}

	// Write
	in := []byte("1234567")
	if n, err := b.Write(in); n != 7 || err != nil {
		t.Fatalf("b.Write(%q) expected 7, <nil>; got %d, %v", in, n, err)
	}
	in = []byte("89")
	if n, err := b.Write(in); n != 0 || err != ErrBlobFull {
		t.Fatalf("b.Write(%q) expected 0, ErrBlobFull; got %d, %v", in, n, err)
	}

	// Reopen
	if err := b.Reopen(2); err != nil {
		t.Fatalf("b.Reopen(2) unexpected error: %v", err)
	}
	if b.Row() != 2 {
		t.Fatalf("b.Row() expected 2; got %d", b.Row())
	}
	if b.Len() != 12 {
		t.Fatalf("b.Len() expected 12; got %d", b.Len())
	}

	// Read
	for i := 0; i < 2; i++ {
		out := make([]byte, 13)
		if n, err := b.Read(out); n != 12 || err != nil {
			t.Fatalf("b.Read() #%d expected 12, <nil>; got %d, %v", i, n, err)
		}
		have := string(out)
		if want := "hello, world\x00"; have != want {
			t.Fatalf("b.Read() #%d expected %q; got %q", i, have, want)
		}
		if p, err := b.Seek(0, 0); p != 0 || err != nil {
			t.Fatalf("b.Seek() #%d expected 0, <nil>; got %d, %v", i, p, err)
		}
	}

	// Close
	t.close(b)
	if err := b.Reopen(1); err != ErrBadIO {
		t.Fatalf("b.Reopen(1) expected %v; got %v", ErrBadIO, err)
	}

	// Verify
	s := t.query(c, "SELECT * FROM x ORDER BY rowid")
	defer t.close(s)
	var have string
	t.scan(s, &have)
	if want := "1234567\x00"; have != want {
		t.Fatalf("s.Scan() expected %q; got %q", want, have)
	}
	t.next(s, nil)
	t.scan(s, &have)
	if want := "hello, world"; have != want {
		t.Fatalf("s.Scan() expected %q; got %q", want, have)
	}
	t.next(s, io.EOF)
}

func TestBackup(T *testing.T) {
	t := begin(T)

	c1, c2 := t.open(":memory:"), t.open(":memory:")
	defer t.close(c1)
	defer t.close(c2)
	t.exec(c1, "CREATE TABLE x(a)")
	t.exec(c1, "INSERT INTO x VALUES(?)", "1234567\x00")
	t.exec(c1, "INSERT INTO x VALUES(?)", "hello, world")

	// Backup
	b, err := c1.Backup("main", c2, "main")
	if b == nil || err != nil {
		t.Fatalf("b.Backup() unexpected error: %v", err)
	}
	defer t.close(b)
	if pr, pt := b.Progress(); pr != 0 || pt != 0 {
		t.Fatalf("b.Progress() expected 0, 0; got %d, %d", pr, pt)
	}
	if err = b.Step(1); err != nil {
		t.Fatalf("b.Step(1) expected <nil>; got %v", err)
	}
	if pr, pt := b.Progress(); pr != 1 || pt != 2 {
		t.Fatalf("b.Progress() expected 1, 2; got %d, %d", pr, pt)
	}
	if err = b.Step(-1); err != io.EOF {
		t.Fatalf("b.Step(-1) expected EOF; got %v", err)
	}

	// Close
	t.close(b)
	if err = b.Step(-1); err != ErrBadBackup {
		t.Fatalf("b.Step(-1) expected %v; got %v", ErrBadBackup, err)
	}

	// Verify
	s := t.query(c2, "SELECT * FROM x ORDER BY rowid")
	defer t.close(s)
	var have string
	t.scan(s, &have)
	if want := "1234567\x00"; have != want {
		t.Fatalf("s.Scan() expected %q; got %q", want, have)
	}
	t.next(s, nil)
	t.scan(s, &have)
	if want := "hello, world"; have != want {
		t.Fatalf("s.Scan() expected %q; got %q", want, have)
	}
	t.next(s, io.EOF)
}

func TestBusyHandler(T *testing.T) {
	t := begin(T)

	tmp := t.tmpFile()
	c1 := t.open(tmp)
	defer t.close(c1)
	c2 := t.open(tmp)
	defer t.close(c2)
	t.exec(c1, "CREATE TABLE x(a); BEGIN; INSERT INTO x VALUES(1);")

	try := func(sql string, want, terr time.Duration) {
		start := time.Now()
		err := c2.Exec(sql)
		have := time.Since(start)
		if have < want-terr || want+terr < have {
			t.Fatalf(cl("c2.Exec(%q) timeout expected %v; got %v"), sql, want, have)
		}
		t.errCode(err, BUSY)
	}
	want := 100 * time.Millisecond
	terr := 50 * time.Millisecond

	// Default
	try("INSERT INTO x VALUES(2)", 0, terr/2)

	// Built-in
	if prev := c2.BusyTimeout(want); prev != nil {
		t.Fatalf("c2.BusyTimeout() expected <nil>; got %v", prev)
	}
	try("INSERT INTO x VALUES(3)", want, terr)

	// Custom
	calls := 0
	handler := func(count int) (retry bool) {
		calls++
		time.Sleep(10 * time.Millisecond)
		return calls == count+1 && calls < 10
	}
	if prev := c2.BusyFunc(handler); prev != nil {
		t.Fatalf("c2.BusyFunc() expected <nil>; got %v", prev)
	}
	try("INSERT INTO x VALUES(4)", want, terr)
	if calls != 10 {
		t.Fatalf("calls expected 10; got %d", calls)
	}

	// Disable
	if prev := c2.BusyTimeout(0); prev == nil {
		t.Fatalf("c2.BusyTimeout() expected %v; got %v", handler, prev)
	}
	try("INSERT INTO x VALUES(5)", 0, terr/2)
}

func TestTxHandler(T *testing.T) {
	t := begin(T)
	defer t.skipRestIfFailed()

	c := t.open(":memory:")
	defer t.close(c)
	t.exec(c, "CREATE TABLE x(a)")

	commit := 0
	rollback := 0
	c.CommitFunc(func() (abort bool) { commit++; return commit >= 2 })
	c.RollbackFunc(func() { rollback++ })

	// Allow
	c.Begin()
	t.exec(c, "INSERT INTO x VALUES(1)")
	t.exec(c, "INSERT INTO x VALUES(2)")
	if err := c.Commit(); err != nil {
		t.Fatalf("c.Commit() unexpected error: %v", err)
	}

	// Deny
	c.Begin()
	t.exec(c, "INSERT INTO x VALUES(3)")
	t.exec(c, "INSERT INTO x VALUES(4)")
	t.errCode(c.Commit(), CONSTRAINT_COMMITHOOK)

	// Verify
	if commit != 2 || rollback != 1 {
		t.Fatalf("commit/rollback expected 2/1; got %d/%d", commit, rollback)
	}
	s := t.query(c, "SELECT * FROM x ORDER BY rowid")
	defer t.close(s)
	var i int
	if t.scan(s, &i); i != 1 {
		t.Fatalf("s.Scan() expected 1; got %d", i)
	}
	t.next(s, nil)
	if t.scan(s, &i); i != 2 {
		t.Fatalf("s.Scan() expected 2; got %d", i)
	}
	t.next(s, io.EOF)
}

func TestUpdateHandler(T *testing.T) {
	t := begin(T)
	defer t.skipRestIfFailed()

	c := t.open(":memory:")
	defer t.close(c)
	t.exec(c, "CREATE TABLE x(a)")

	type update struct {
		op      int
		db, tbl string
		row     int64
	}
	var have *update
	verify := func(want *update) {
		if !reflect.DeepEqual(have, want) {
			t.Fatalf(cl("verify() expected %v; got %v"), want, have)
		}
	}
	c.UpdateFunc(func(op int, db, tbl RawString, row int64) {
		have = &update{op, db.Copy(), tbl.Copy(), row}
	})

	t.exec(c, "INSERT INTO x VALUES(1)")
	verify(&update{INSERT, "main", "x", 1})

	t.exec(c, "INSERT INTO x VALUES(2)")
	verify(&update{INSERT, "main", "x", 2})

	t.exec(c, "UPDATE x SET a=3 WHERE rowid=1")
	verify(&update{UPDATE, "main", "x", 1})

	t.exec(c, "DELETE FROM x WHERE rowid=2")
	verify(&update{DELETE, "main", "x", 2})
}

func TestSchema(T *testing.T) {
	t := begin(T)

	c := t.open(":memory:")
	defer t.close(c)
	t.exec(c, "CREATE TABLE x(a int)")
	t.exec(c, "INSERT INTO x VALUES(1)")
	t.exec(c, "INSERT INTO x VALUES(2)")

	checkCols := func(s *Stmt, want ...string) {
		if have := s.NumColumns(); have != len(want) {
			t.Fatalf(cl("s.NumColumns() expected %d; got %d"), len(want), have)
		}
		if have := s.Columns(); !reflect.DeepEqual(have, want) {
			t.Fatalf(cl("s.Columns() expected %v; got %v"), want, have)
		}
	}
	checkDecls := func(s *Stmt, want ...string) {
		if have := s.DeclTypes(); !reflect.DeepEqual(have, want) {
			t.Fatalf(cl("s.DeclTypes() expected %v; got %v"), want, have)
		}
	}
	checkRow := func(s *Stmt, want RowMap) {
		have := RowMap{}
		if t.scan(s, have); !reflect.DeepEqual(have, want) {
			t.Fatalf(cl("s.Scan() expected %v; got %v"), want, have)
		}
	}
	s := t.query(c, "SELECT * FROM x ORDER BY rowid")
	defer t.close(s)

	checkCols(s, "a")
	checkDecls(s, "INT")
	checkRow(s, RowMap{"a": int64(1)})

	// Schema changes do not affect running statements
	t.exec(c, "ALTER TABLE x ADD b text")
	t.next(s, nil)

	checkCols(s, "a")
	checkDecls(s, "INT")
	checkRow(s, RowMap{"a": int64(2)})
	t.next(s, io.EOF)

	checkCols(s, "a")
	checkDecls(s, "INT")
	t.query(s)

	checkCols(s, "a", "b")
	checkDecls(s, "INT", "TEXT")
	checkRow(s, RowMap{"a": int64(1), "b": nil})
	t.next(s, nil)

	checkCols(s, "a", "b")
	checkDecls(s, "INT", "TEXT")
	checkRow(s, RowMap{"a": int64(2), "b": nil})
	t.next(s, io.EOF)
}

func TestDriver(T *testing.T) {
	t := begin(T)

	c, err := sql.Open("sqlite3", ":memory:")
	if c == nil || err != nil {
		t.Fatalf("sql.Open() unexpected error: %v", err)
	}
	defer t.close(c)

	// Setup
	sql := "CREATE TABLE x(a, b, c)"
	r, err := c.Exec(sql)
	if r == nil || err != nil {
		t.Fatalf("c.Exec(%q) unexpected error: %v", sql, err)
	}
	if id, err := r.LastInsertId(); id != 0 || err != nil {
		t.Fatalf("r.LastInsertId() expected 0, <nil>; got %d, %v", id, err)
	}
	if n, err := r.RowsAffected(); n != 0 || err != nil {
		t.Fatalf("r.RowsAffected() expected 0, <nil>; got %d, %v", n, err)
	}

	// Prepare
	sql = "INSERT INTO x VALUES(?, ?, ?)"
	s, err := c.Prepare(sql)
	if err != nil {
		t.Fatalf("c.Prepare(%q) unexpected error: %v", sql, err)
	}
	defer t.close(s)

	// Multiple inserts
	r, err = s.Exec(1, 2.2, "test")
	if err != nil {
		t.Fatalf("s.Exec(%q) unexpected error: %v", sql, err)
	}
	if id, err := r.LastInsertId(); id != 1 || err != nil {
		t.Fatalf("r.LastInsertId() expected 1, <nil>; got %d, %v", id, err)
	}
	if n, err := r.RowsAffected(); n != 1 || err != nil {
		t.Fatalf("r.RowsAffected() expected 1, <nil>; got %d, %v", n, err)
	}

	r, err = s.Exec(3, []byte{4}, nil)
	if err != nil {
		t.Fatalf("s.Exec(%q) unexpected error: %v", sql, err)
	}
	if id, err := r.LastInsertId(); id != 2 || err != nil {
		t.Fatalf("r.LastInsertId() expected 1, <nil>; got %d, %v", id, err)
	}
	if n, err := r.RowsAffected(); n != 1 || err != nil {
		t.Fatalf("r.RowsAffected() expected 1, <nil>; got %d, %v", n, err)
	}

	// Select all rows
	sql = "SELECT rowid, * FROM x ORDER BY rowid"
	rows, err := c.Query(sql)
	if rows == nil || err != nil {
		t.Fatalf("c.Query() unexpected error: %v", err)
	}
	defer t.close(rows)

	// Row information
	want := []string{"rowid", "a", "b", "c"}
	if have, err := rows.Columns(); !reflect.DeepEqual(have, want) {
		t.Fatalf("rows.Columns() expected %v, <nil>; got %v, %v", want, have, err)
	}

	// Verify
	table := [][]interface{}{
		{int64(1), int64(1), float64(2.2), []byte("test")},
		{int64(2), int64(3), []byte{4}, nil},
	}
	for i, want := range table {
		if !rows.Next() {
			t.Fatalf("rows.Next(%d) expected true", i)
		}
		have := make([]interface{}, 4)
		if err := rows.Scan(&have[0], &have[1], &have[2], &have[3]); err != nil {
			t.Fatalf("rows.Scan() unexpected error: %v", err)
		}
		if !reflect.DeepEqual(have, want) {
			t.Fatalf("rows.Scan() expected %v; got %v", want, have)
		}
	}
	if rows.Next() {
		t.Fatalf("rows.Next() expected false")
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("rows.Err() unexpected error: %v", err)
	}
}
