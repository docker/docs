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

// Backup is a handle to an online backup operation between two databases.
// [http://www.sqlite.org/c3ref/backup.html]
type Backup struct {
	src  *Conn
	dst  *Conn
	bkup *C.sqlite3_backup
}

// newBackup initializes an online backup operation from src.srcName to
// dst.dstName.
func newBackup(src *Conn, srcName string, dst *Conn, dstName string) (*Backup, error) {
	srcName += "\x00"
	dstName += "\x00"

	bkup := C.sqlite3_backup_init(dst.db, cStr(dstName), src.db, cStr(srcName))
	if bkup == nil {
		return nil, libErr(C.sqlite3_errcode(dst.db), dst.db)
	}

	b := &Backup{src, dst, bkup}
	runtime.SetFinalizer(b, (*Backup).Close)
	return b, nil
}

// Close releases all resources associated with the backup operation. It is safe
// to call this method prior to backup completion to abort the operation.
// [http://www.sqlite.org/c3ref/backup_finish.html#sqlite3backupfinish]
func (b *Backup) Close() error {
	if bkup := b.bkup; bkup != nil {
		b.bkup = nil
		runtime.SetFinalizer(b, nil)
		if rc := C.sqlite3_backup_finish(bkup); rc != OK {
			return libErr(rc, b.dst.db)
		}
	}
	return nil
}

// Conn returns the source and destination connections that are used by this
// backup operation. The destination connection must not be used until the
// backup operation is closed.
func (b *Backup) Conn() (src, dst *Conn) {
	return b.src, b.dst
}

// Step copies up to n pages to the destination database. If n is negative, all
// remaining pages are copied. io.EOF is returned upon successful backup
// completion.
// [http://www.sqlite.org/c3ref/backup_finish.html#sqlite3backupstep]
func (b *Backup) Step(n int) error {
	if b.bkup == nil {
		return ErrBadBackup
	}
	if rc := C.sqlite3_backup_step(b.bkup, C.int(n)); rc != OK {
		// Do not close automatically since that clears the progress info
		if rc == DONE {
			return io.EOF
		}
		return libErr(rc, b.dst.db)
	}
	return nil
}

// Progress returns the number of pages that still need to be backed up and the
// total number of pages in the source database. The values are updated after
// each call to Step and are reset to 0 after the backup is closed. The total
// number of pages may change if the source database is modified during the
// backup operation.
// [http://www.sqlite.org/c3ref/backup_finish.html#sqlite3backupremaining]
func (b *Backup) Progress() (remaining, total int) {
	if b.bkup != nil {
		remaining = int(C.sqlite3_backup_remaining(b.bkup))
		total = int(C.sqlite3_backup_pagecount(b.bkup))
	}
	return
}
