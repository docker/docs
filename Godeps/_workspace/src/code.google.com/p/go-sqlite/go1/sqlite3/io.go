// Copyright 2013 The Go-SQLite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sqlite3

/*
#include "sqlite3.h"
*/
import "C"

import (
	"io"
	"runtime"
)

// ErrBlobFull is returned by BlobIO.Write when there isn't enough space left to
// write the provided bytes.
var ErrBlobFull = &Error{ERROR, "incremental write failed, no space left"}

// BlobIO is a handle to a single BLOB (binary large object) or TEXT value
// opened for incremental I/O. This allows the value to be treated as a file for
// reading and writing. The value length cannot be changed using this API; use
// an UPDATE statement for that. The recommended way of allocating space for a
// BLOB is to use the ZeroBlob type or the zeroblob() SQL function.
// [http://www.sqlite.org/c3ref/blob.html]
type BlobIO struct {
	conn *Conn
	blob *C.sqlite3_blob

	row int64 // ROWID of the row containing the BLOB/TEXT value
	len int   // Value length in bytes
	off int   // Current read/write offset
}

// newBlobIO initializes an incremental I/O operation.
func newBlobIO(c *Conn, db, tbl, col string, row int64, rw bool) (*BlobIO, error) {
	db += "\x00"
	tbl += "\x00"
	col += "\x00"

	var blob *C.sqlite3_blob
	rc := C.sqlite3_blob_open(c.db, cStr(db), cStr(tbl), cStr(col),
		C.sqlite3_int64(row), cBool(rw), &blob)
	if rc != OK {
		return nil, libErr(rc, c.db)
	}

	b := &BlobIO{
		conn: c,
		blob: blob,
		row:  row,
		len:  int(C.sqlite3_blob_bytes(blob)),
	}
	runtime.SetFinalizer(b, (*BlobIO).Close)
	return b, nil
}

// Close releases all resources associated with the incremental I/O operation.
// It is important to check the error returned by this method, since disk I/O
// and other types of errors may not be reported until the changes are actually
// committed to the database.
// [http://www.sqlite.org/c3ref/blob_close.html]
func (b *BlobIO) Close() error {
	if blob := b.blob; blob != nil {
		b.blob = nil
		b.len = 0
		b.off = 0
		runtime.SetFinalizer(b, nil)
		if rc := C.sqlite3_blob_close(blob); rc != OK {
			return libErr(rc, b.conn.db)
		}
	}
	return nil
}

// Conn returns the connection that that created this incremental I/O operation.
func (b *BlobIO) Conn() *Conn {
	return b.conn
}

// Row returns the ROWID of the row containing the BLOB/TEXT value.
func (b *BlobIO) Row() int64 {
	return b.row
}

// Len returns the length of the BLOB/TEXT value in bytes. It is not possible to
// read/write/seek beyond this length. The length changes to 0 if the I/O handle
// expires due to an update of any column in the same row. This condition is
// indicated by an ABORT error code returned from Read or Write. An expired
// handle is closed automatically and cannot be reopened. Any writes that
// occurred before the abort are not rolled back.
// [http://www.sqlite.org/c3ref/blob_bytes.html]
func (b *BlobIO) Len() int {
	return b.len
}

// Read implements the io.Reader interface.
// [http://www.sqlite.org/c3ref/blob_read.html]
func (b *BlobIO) Read(p []byte) (n int, err error) {
	if b.blob == nil {
		return 0, ErrBadIO
	}
	if b.off >= b.len {
		return 0, io.EOF
	}
	if n = b.len - b.off; len(p) < n {
		n = len(p)
	}
	rc := C.sqlite3_blob_read(b.blob, cBytes(p), C.int(n), C.int(b.off))
	return b.io(rc, n)
}

// Write implements the io.Writer interface. The number of bytes written is
// always either 0 or len(p). ErrBlobFull is returned if there isn't enough
// space left to write all of p.
// [http://www.sqlite.org/c3ref/blob_write.html]
func (b *BlobIO) Write(p []byte) (n int, err error) {
	if b.blob == nil {
		return 0, ErrBadIO
	}
	if n = len(p); b.off+n > b.len {
		// Doesn't make sense to do a partial write. Better to return quickly
		// and let the caller reallocate the BLOB.
		return 0, ErrBlobFull
	}
	rc := C.sqlite3_blob_write(b.blob, cBytes(p), C.int(n), C.int(b.off))
	return b.io(rc, n)
}

// Seek implements the io.Seeker interface.
func (b *BlobIO) Seek(offset int64, whence int) (ret int64, err error) {
	if b.blob == nil {
		return 0, ErrBadIO
	}
	switch whence {
	case 0:
	case 1:
		offset += int64(b.off)
	case 2:
		offset += int64(b.len)
	default:
		return 0, pkgErr(MISUSE, "invalid whence for BlobIO.Seek (%d)", whence)
	}
	if offset < 0 || offset > int64(b.len) {
		return 0, pkgErr(MISUSE, "invalid offset for BlobIO.Seek (%d)", offset)
	}
	b.off = int(offset)
	return offset, nil
}

// Reopen closes the current value and opens another one in the same column,
// specified by its ROWID. If an error is encountered, the I/O handle becomes
// unusable and is automatically closed.
// [http://www.sqlite.org/c3ref/blob_reopen.html]
func (b *BlobIO) Reopen(row int64) error {
	if b.blob == nil {
		return ErrBadIO
	}
	if rc := C.sqlite3_blob_reopen(b.blob, C.sqlite3_int64(row)); rc != OK {
		err := libErr(rc, b.conn.db)
		b.Close()
		return err
	}
	b.row = row
	b.len = int(C.sqlite3_blob_bytes(b.blob))
	b.off = 0
	return nil
}

// io handles the completion of a single Read/Write call.
func (b *BlobIO) io(rc C.int, n int) (int, error) {
	if rc == OK {
		b.off += n
		return n, nil
	}
	err := libErr(rc, b.conn.db)
	if rc == ABORT {
		b.Close()
	}
	return 0, err
}
