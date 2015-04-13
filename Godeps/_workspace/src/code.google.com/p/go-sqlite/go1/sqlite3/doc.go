// Copyright 2013 The Go-SQLite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package sqlite3 provides an interface to SQLite version 3 databases.

Database connections are created either by using this package directly or with
the "sqlite3" database/sql driver. The direct interface, which is described
below, exposes SQLite-specific features, such as incremental I/O and online
backups. The driver is recommended when your application has to support multiple
database engines.

Installation

Minimum requirements are Go 1.1+ with CGO enabled and GCC/MinGW C compiler. The
SQLite amalgamation version 3.8.0.2 (2013-09-03) is compiled as part of the
package (see http://www.sqlite.org/amalgamation.html). Compilation options are
defined at the top of sqlite3.go (#cgo CFLAGS). Dynamic linking with a shared
SQLite library is not supported.

Windows users should install mingw-w64 (http://mingw-w64.sourceforge.net/),
TDM64-GCC (http://tdm-gcc.tdragon.net/), or another MinGW distribution, and make
sure that gcc.exe is available from the %PATH%. MSYS is not required.

Run 'go get code.google.com/p/go-sqlite/go1/sqlite3' to download, build, and
install the package.

Concurrency

A single connection instance and all of its derived objects (prepared
statements, backup operations, etc.) may NOT be used concurrently from multiple
goroutines without external synchronization. The only exception is
Conn.Interrupt(), which may be called from another goroutine to abort a
long-running operation. It is safe to use separate connection instances
concurrently, even if they are accessing the same database file. For example:

	// ERROR (without any extra synchronization)
	c, _ := sqlite3.Open("sqlite.db")
	go use(c)
	go use(c)

	// OK
	c1, _ := sqlite3.Open("sqlite.db")
	c2, _ := sqlite3.Open("sqlite.db")
	go use(c1)
	go use(c2)

Maps

Use NamedArgs map to bind values to named statement parameters (see
http://www.sqlite.org/lang_expr.html#varparam). Use RowMap to retrieve the
current row as a map of column/value pairs. Here is a short example with the
error-handling code omitted for brevity:

	c, _ := sqlite3.Open(":memory:")
	c.Exec("CREATE TABLE x(a, b, c)")

	args := sqlite3.NamedArgs{"$a": 1, "$b": "demo"}
	c.Exec("INSERT INTO x VALUES($a, $b, $c)", args) // $c will be NULL

	sql := "SELECT rowid, * FROM x"
	row := make(sqlite3.RowMap)
	for s, err := c.Query(sql); err == nil; err = s.Next() {
		var rowid int64
		s.Scan(&rowid, row)     // Assigns 1st column to rowid, the rest to row
		fmt.Println(rowid, row) // Prints "1 map[a:1 b:demo c:<nil>]"
	}

Data Types

See http://www.sqlite.org/datatype3.html for a description of the SQLite data
type system. The following Go data types are supported as arguments to prepared
statements and may be used in NamedArgs:

	Go Type    SQLite Type  Notes
	---------  -----------  ----------------------------------------------------
	<nil>      NULL         Unbound parameters are NULL by default.
	int        INTEGER
	int64      INTEGER
	float64    FLOAT
	bool       INTEGER      Converted as false = 0, true = 1.
	string     TEXT         SQLite makes a private copy when the value is bound.
	[]byte     BLOB         SQLite makes a private copy when the value is bound.
	time.Time  INTEGER      Converted by calling Unix().
	RawString  TEXT         SQLite uses the value directly without copying. The
	                        caller must keep a reference to the value for the
	                        duration of the query to prevent garbage collection.
	RawBytes   BLOB         Same as RawString. The value must not be modified
	                        for the duration of the query.
	ZeroBlob   BLOB         Allocates a zero-filled BLOB of the specified length
	                        (e.g. ZeroBlob(4096) allocates 4KB).

Note that the table above describes how the value is bound to the statement. The
final storage class is determined according to the column affinity rules.

See http://www.sqlite.org/c3ref/column_blob.html for a description of how column
values are retrieved from the results of a query. The following static Go data
types are supported for retrieving column values:

	Go Type     Req. Type  Notes
	----------  ---------  ---------------------------------------------------
	*int        INTEGER
	*int64      INTEGER
	*float64    FLOAT
	*bool       INTEGER    Converted as 0 = false, otherwise true.
	*string     TEXT       The caller receives a copy of the value.
	*[]byte     BLOB       The caller receives a copy of the value.
	*time.Time  INTEGER    Converted by calling time.Unix(). Text values are not
	                       supported, but the conversion can be performed with
	                       the date and time SQL functions.
	*RawString  TEXT       The value is used directly without copying and
	                       remains valid until the next Stmt method call.
	*RawBytes   BLOB       Same as *RawString. The value must not be modified.
	                       Re-slicing is ok, but be careful with append().
	io.Writer   BLOB       The value is written out directly into the writer.

For *interface{} and RowMap arguments, the Go data type is dynamically selected
based on the SQLite storage class and column declaration prefix:

	SQLite Type  Col. Decl.  Go Type    Notes
	-----------  ----------  ---------  ----------------------------------------
	NULL                     <nil>
	INTEGER      "DATE..."   time.Time  Converted by calling time.Unix().
	INTEGER      "TIME..."   time.Time  Converted by calling time.Unix().
	INTEGER      "BOOL..."   bool       Converted as 0 = false, otherwise true.
	INTEGER                  int64
	FLOAT                    float64
	TEXT                     string
	BLOB                     []byte

Database Names

Methods that require a database name as one of the arguments (e.g. Conn.Path())
expect the symbolic name by which the database is known to the connection, not a
path to a file. Valid database names are "main", "temp", or a name specified
after the AS clause in an ATTACH statement.

Callbacks

SQLite can execute callbacks for various internal events. The package provides
types and methods for registering callback handlers. Unless stated otherwise in
SQLite documentation, callback handlers are not reentrant and must not do
anything to modify the associated database connection. This includes
preparing/running any other SQL statements. The safest bet is to avoid all
interactions with Conn, Stmt, and other related objects within the handler.

Codecs and Encryption

SQLite has an undocumented codec API, which operates between the pager and VFS
layers, and is used by the SQLite Encryption Extension (SEE) to encrypt database
and journal contents. Consider purchasing a SEE license if you require
production-quality encryption support (http://www.hwaci.com/sw/sqlite/see.html).

This package has an experimental API (read: unstable, may eat your data) for
writing codecs in Go. The "codec" subpackage provides additional documentation
and several existing codec implementations.

Codecs are registered via the RegisterCodec function for a specific key prefix.
For example, the "aes-hmac" codec is initialized when a key in the format
"aes-hmac:<...>" is provided to an attached database. The key format after the
first colon is codec-specific. See CodecFunc for more information.

The codec API has several limitations. Codecs cannot be used for in-memory or
temporary databases. Once a database is created, the page size and the amount of
reserved space at the end of each page cannot be changed (i.e. "PRAGMA
page_size=N; VACUUM;" will not work). Online backups will fail unless the
destination database has the same page size and reserve values. Bytes 16 through
23 of page 1 (the database header, see http://www.sqlite.org/fileformat2.html)
cannot be altered, so it is always possible to identify encrypted SQLite
databases.

The rekey function is currently not implemented. The key can only be changed via
the backup API or by dumping and restoring the database contents.
*/
package sqlite3
