// Copyright 2013 The Go-SQLite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sqlite3

/*
#include "sqlite3.h"
*/
import "C"

// Fundamental SQLite data types. These are returned by Stmt.DataTypes method.
// [http://www.sqlite.org/c3ref/c_blob.html]
const (
	INTEGER = C.SQLITE_INTEGER // 1
	FLOAT   = C.SQLITE_FLOAT   // 2
	TEXT    = C.SQLITE_TEXT    // 3
	BLOB    = C.SQLITE_BLOB    // 4
	NULL    = C.SQLITE_NULL    // 5
)

// General result codes returned by the SQLite API. When converted to an error,
// OK and ROW become nil, and DONE becomes either nil or io.EOF, depending on
// the context in which the statement is executed. All other codes are returned
// via the Error struct.
// [http://www.sqlite.org/c3ref/c_abort.html]
const (
	OK         = C.SQLITE_OK         // 0   = Successful result
	ERROR      = C.SQLITE_ERROR      // 1   = SQL error or missing database
	INTERNAL   = C.SQLITE_INTERNAL   // 2   = Internal logic error in SQLite
	PERM       = C.SQLITE_PERM       // 3   = Access permission denied
	ABORT      = C.SQLITE_ABORT      // 4   = Callback routine requested an abort
	BUSY       = C.SQLITE_BUSY       // 5   = The database file is locked
	LOCKED     = C.SQLITE_LOCKED     // 6   = A table in the database is locked
	NOMEM      = C.SQLITE_NOMEM      // 7   = A malloc() failed
	READONLY   = C.SQLITE_READONLY   // 8   = Attempt to write a readonly database
	INTERRUPT  = C.SQLITE_INTERRUPT  // 9   = Operation terminated by sqlite3_interrupt()
	IOERR      = C.SQLITE_IOERR      // 10  = Some kind of disk I/O error occurred
	CORRUPT    = C.SQLITE_CORRUPT    // 11  = The database disk image is malformed
	NOTFOUND   = C.SQLITE_NOTFOUND   // 12  = Unknown opcode in sqlite3_file_control()
	FULL       = C.SQLITE_FULL       // 13  = Insertion failed because database is full
	CANTOPEN   = C.SQLITE_CANTOPEN   // 14  = Unable to open the database file
	PROTOCOL   = C.SQLITE_PROTOCOL   // 15  = Database lock protocol error
	EMPTY      = C.SQLITE_EMPTY      // 16  = Database is empty
	SCHEMA     = C.SQLITE_SCHEMA     // 17  = The database schema changed
	TOOBIG     = C.SQLITE_TOOBIG     // 18  = String or BLOB exceeds size limit
	CONSTRAINT = C.SQLITE_CONSTRAINT // 19  = Abort due to constraint violation
	MISMATCH   = C.SQLITE_MISMATCH   // 20  = Data type mismatch
	MISUSE     = C.SQLITE_MISUSE     // 21  = Library used incorrectly
	NOLFS      = C.SQLITE_NOLFS      // 22  = Uses OS features not supported on host
	AUTH       = C.SQLITE_AUTH       // 23  = Authorization denied
	FORMAT     = C.SQLITE_FORMAT     // 24  = Auxiliary database format error
	RANGE      = C.SQLITE_RANGE      // 25  = 2nd parameter to sqlite3_bind out of range
	NOTADB     = C.SQLITE_NOTADB     // 26  = File opened that is not a database file
	NOTICE     = C.SQLITE_NOTICE     // 27  = Notifications from sqlite3_log()
	WARNING    = C.SQLITE_WARNING    // 28  = Warnings from sqlite3_log()
	ROW        = C.SQLITE_ROW        // 100 = sqlite3_step() has another row ready
	DONE       = C.SQLITE_DONE       // 101 = sqlite3_step() has finished executing
)

// Extended result codes returned by the SQLite API. Extended result codes are
// enabled by default for all new Conn objects. Use Error.Code()&0xFF to convert
// an extended code to a general one.
// [http://www.sqlite.org/c3ref/c_abort_rollback.html]
const (
	IOERR_READ              = C.SQLITE_IOERR_READ              // (SQLITE_IOERR | (1<<8))
	IOERR_SHORT_READ        = C.SQLITE_IOERR_SHORT_READ        // (SQLITE_IOERR | (2<<8))
	IOERR_WRITE             = C.SQLITE_IOERR_WRITE             // (SQLITE_IOERR | (3<<8))
	IOERR_FSYNC             = C.SQLITE_IOERR_FSYNC             // (SQLITE_IOERR | (4<<8))
	IOERR_DIR_FSYNC         = C.SQLITE_IOERR_DIR_FSYNC         // (SQLITE_IOERR | (5<<8))
	IOERR_TRUNCATE          = C.SQLITE_IOERR_TRUNCATE          // (SQLITE_IOERR | (6<<8))
	IOERR_FSTAT             = C.SQLITE_IOERR_FSTAT             // (SQLITE_IOERR | (7<<8))
	IOERR_UNLOCK            = C.SQLITE_IOERR_UNLOCK            // (SQLITE_IOERR | (8<<8))
	IOERR_RDLOCK            = C.SQLITE_IOERR_RDLOCK            // (SQLITE_IOERR | (9<<8))
	IOERR_DELETE            = C.SQLITE_IOERR_DELETE            // (SQLITE_IOERR | (10<<8))
	IOERR_BLOCKED           = C.SQLITE_IOERR_BLOCKED           // (SQLITE_IOERR | (11<<8))
	IOERR_NOMEM             = C.SQLITE_IOERR_NOMEM             // (SQLITE_IOERR | (12<<8))
	IOERR_ACCESS            = C.SQLITE_IOERR_ACCESS            // (SQLITE_IOERR | (13<<8))
	IOERR_CHECKRESERVEDLOCK = C.SQLITE_IOERR_CHECKRESERVEDLOCK // (SQLITE_IOERR | (14<<8))
	IOERR_LOCK              = C.SQLITE_IOERR_LOCK              // (SQLITE_IOERR | (15<<8))
	IOERR_CLOSE             = C.SQLITE_IOERR_CLOSE             // (SQLITE_IOERR | (16<<8))
	IOERR_DIR_CLOSE         = C.SQLITE_IOERR_DIR_CLOSE         // (SQLITE_IOERR | (17<<8))
	IOERR_SHMOPEN           = C.SQLITE_IOERR_SHMOPEN           // (SQLITE_IOERR | (18<<8))
	IOERR_SHMSIZE           = C.SQLITE_IOERR_SHMSIZE           // (SQLITE_IOERR | (19<<8))
	IOERR_SHMLOCK           = C.SQLITE_IOERR_SHMLOCK           // (SQLITE_IOERR | (20<<8))
	IOERR_SHMMAP            = C.SQLITE_IOERR_SHMMAP            // (SQLITE_IOERR | (21<<8))
	IOERR_SEEK              = C.SQLITE_IOERR_SEEK              // (SQLITE_IOERR | (22<<8))
	IOERR_DELETE_NOENT      = C.SQLITE_IOERR_DELETE_NOENT      // (SQLITE_IOERR | (23<<8))
	IOERR_MMAP              = C.SQLITE_IOERR_MMAP              // (SQLITE_IOERR | (24<<8))
	IOERR_GETTEMPPATH       = C.SQLITE_IOERR_GETTEMPPATH       // (SQLITE_IOERR | (25<<8))
	LOCKED_SHAREDCACHE      = C.SQLITE_LOCKED_SHAREDCACHE      // (SQLITE_LOCKED |  (1<<8))
	BUSY_RECOVERY           = C.SQLITE_BUSY_RECOVERY           // (SQLITE_BUSY   |  (1<<8))
	BUSY_SNAPSHOT           = C.SQLITE_BUSY_SNAPSHOT           // (SQLITE_BUSY   |  (2<<8))
	CANTOPEN_NOTEMPDIR      = C.SQLITE_CANTOPEN_NOTEMPDIR      // (SQLITE_CANTOPEN | (1<<8))
	CANTOPEN_ISDIR          = C.SQLITE_CANTOPEN_ISDIR          // (SQLITE_CANTOPEN | (2<<8))
	CANTOPEN_FULLPATH       = C.SQLITE_CANTOPEN_FULLPATH       // (SQLITE_CANTOPEN | (3<<8))
	CORRUPT_VTAB            = C.SQLITE_CORRUPT_VTAB            // (SQLITE_CORRUPT | (1<<8))
	READONLY_RECOVERY       = C.SQLITE_READONLY_RECOVERY       // (SQLITE_READONLY | (1<<8))
	READONLY_CANTLOCK       = C.SQLITE_READONLY_CANTLOCK       // (SQLITE_READONLY | (2<<8))
	READONLY_ROLLBACK       = C.SQLITE_READONLY_ROLLBACK       // (SQLITE_READONLY | (3<<8))
	ABORT_ROLLBACK          = C.SQLITE_ABORT_ROLLBACK          // (SQLITE_ABORT | (2<<8))
	CONSTRAINT_CHECK        = C.SQLITE_CONSTRAINT_CHECK        // (SQLITE_CONSTRAINT | (1<<8))
	CONSTRAINT_COMMITHOOK   = C.SQLITE_CONSTRAINT_COMMITHOOK   // (SQLITE_CONSTRAINT | (2<<8))
	CONSTRAINT_FOREIGNKEY   = C.SQLITE_CONSTRAINT_FOREIGNKEY   // (SQLITE_CONSTRAINT | (3<<8))
	CONSTRAINT_FUNCTION     = C.SQLITE_CONSTRAINT_FUNCTION     // (SQLITE_CONSTRAINT | (4<<8))
	CONSTRAINT_NOTNULL      = C.SQLITE_CONSTRAINT_NOTNULL      // (SQLITE_CONSTRAINT | (5<<8))
	CONSTRAINT_PRIMARYKEY   = C.SQLITE_CONSTRAINT_PRIMARYKEY   // (SQLITE_CONSTRAINT | (6<<8))
	CONSTRAINT_TRIGGER      = C.SQLITE_CONSTRAINT_TRIGGER      // (SQLITE_CONSTRAINT | (7<<8))
	CONSTRAINT_UNIQUE       = C.SQLITE_CONSTRAINT_UNIQUE       // (SQLITE_CONSTRAINT | (8<<8))
	CONSTRAINT_VTAB         = C.SQLITE_CONSTRAINT_VTAB         // (SQLITE_CONSTRAINT | (9<<8))
	NOTICE_RECOVER_WAL      = C.SQLITE_NOTICE_RECOVER_WAL      // (SQLITE_NOTICE | (1<<8))
	NOTICE_RECOVER_ROLLBACK = C.SQLITE_NOTICE_RECOVER_ROLLBACK // (SQLITE_NOTICE | (2<<8))
	WARNING_AUTOINDEX       = C.SQLITE_WARNING_AUTOINDEX       // (SQLITE_WARNING | (1<<8))
)

// Codes used by SQLite to indicate the operation type when invoking authorizer
// and row update callbacks.
// [http://www.sqlite.org/c3ref/c_alter_table.html]
const (
	CREATE_INDEX        = C.SQLITE_CREATE_INDEX        // 1
	CREATE_TABLE        = C.SQLITE_CREATE_TABLE        // 2
	CREATE_TEMP_INDEX   = C.SQLITE_CREATE_TEMP_INDEX   // 3
	CREATE_TEMP_TABLE   = C.SQLITE_CREATE_TEMP_TABLE   // 4
	CREATE_TEMP_TRIGGER = C.SQLITE_CREATE_TEMP_TRIGGER // 5
	CREATE_TEMP_VIEW    = C.SQLITE_CREATE_TEMP_VIEW    // 6
	CREATE_TRIGGER      = C.SQLITE_CREATE_TRIGGER      // 7
	CREATE_VIEW         = C.SQLITE_CREATE_VIEW         // 8
	DELETE              = C.SQLITE_DELETE              // 9
	DROP_INDEX          = C.SQLITE_DROP_INDEX          // 10
	DROP_TABLE          = C.SQLITE_DROP_TABLE          // 11
	DROP_TEMP_INDEX     = C.SQLITE_DROP_TEMP_INDEX     // 12
	DROP_TEMP_TABLE     = C.SQLITE_DROP_TEMP_TABLE     // 13
	DROP_TEMP_TRIGGER   = C.SQLITE_DROP_TEMP_TRIGGER   // 14
	DROP_TEMP_VIEW      = C.SQLITE_DROP_TEMP_VIEW      // 15
	DROP_TRIGGER        = C.SQLITE_DROP_TRIGGER        // 16
	DROP_VIEW           = C.SQLITE_DROP_VIEW           // 17
	INSERT              = C.SQLITE_INSERT              // 18
	PRAGMA              = C.SQLITE_PRAGMA              // 19
	READ                = C.SQLITE_READ                // 20
	SELECT              = C.SQLITE_SELECT              // 21
	TRANSACTION         = C.SQLITE_TRANSACTION         // 22
	UPDATE              = C.SQLITE_UPDATE              // 23
	ATTACH              = C.SQLITE_ATTACH              // 24
	DETACH              = C.SQLITE_DETACH              // 25
	ALTER_TABLE         = C.SQLITE_ALTER_TABLE         // 26
	REINDEX             = C.SQLITE_REINDEX             // 27
	ANALYZE             = C.SQLITE_ANALYZE             // 28
	CREATE_VTABLE       = C.SQLITE_CREATE_VTABLE       // 29
	DROP_VTABLE         = C.SQLITE_DROP_VTABLE         // 30
	FUNCTION            = C.SQLITE_FUNCTION            // 31
	SAVEPOINT           = C.SQLITE_SAVEPOINT           // 32
)

// Core SQLite performance counters that can be queried with Status.
// [http://www.sqlite.org/c3ref/c_status_malloc_count.html]
const (
	STATUS_MEMORY_USED        = C.SQLITE_STATUS_MEMORY_USED        // 0
	STATUS_PAGECACHE_USED     = C.SQLITE_STATUS_PAGECACHE_USED     // 1
	STATUS_PAGECACHE_OVERFLOW = C.SQLITE_STATUS_PAGECACHE_OVERFLOW // 2
	STATUS_SCRATCH_USED       = C.SQLITE_STATUS_SCRATCH_USED       // 3
	STATUS_SCRATCH_OVERFLOW   = C.SQLITE_STATUS_SCRATCH_OVERFLOW   // 4
	STATUS_MALLOC_SIZE        = C.SQLITE_STATUS_MALLOC_SIZE        // 5
	STATUS_PARSER_STACK       = C.SQLITE_STATUS_PARSER_STACK       // 6
	STATUS_PAGECACHE_SIZE     = C.SQLITE_STATUS_PAGECACHE_SIZE     // 7
	STATUS_SCRATCH_SIZE       = C.SQLITE_STATUS_SCRATCH_SIZE       // 8
	STATUS_MALLOC_COUNT       = C.SQLITE_STATUS_MALLOC_COUNT       // 9
)

// Connection performance counters that can be queried with Conn.Status.
// [http://www.sqlite.org/c3ref/c_dbstatus_options.html]
const (
	DBSTATUS_LOOKASIDE_USED      = C.SQLITE_DBSTATUS_LOOKASIDE_USED      // 0
	DBSTATUS_CACHE_USED          = C.SQLITE_DBSTATUS_CACHE_USED          // 1
	DBSTATUS_SCHEMA_USED         = C.SQLITE_DBSTATUS_SCHEMA_USED         // 2
	DBSTATUS_STMT_USED           = C.SQLITE_DBSTATUS_STMT_USED           // 3
	DBSTATUS_LOOKASIDE_HIT       = C.SQLITE_DBSTATUS_LOOKASIDE_HIT       // 4
	DBSTATUS_LOOKASIDE_MISS_SIZE = C.SQLITE_DBSTATUS_LOOKASIDE_MISS_SIZE // 5
	DBSTATUS_LOOKASIDE_MISS_FULL = C.SQLITE_DBSTATUS_LOOKASIDE_MISS_FULL // 6
	DBSTATUS_CACHE_HIT           = C.SQLITE_DBSTATUS_CACHE_HIT           // 7
	DBSTATUS_CACHE_MISS          = C.SQLITE_DBSTATUS_CACHE_MISS          // 8
	DBSTATUS_CACHE_WRITE         = C.SQLITE_DBSTATUS_CACHE_WRITE         // 9
	DBSTATUS_DEFERRED_FKS        = C.SQLITE_DBSTATUS_DEFERRED_FKS        // 10
)

// Statement performance counters that can be queried with Stmt.Status.
// [http://www.sqlite.org/c3ref/c_stmtstatus_counter.html]
const (
	STMTSTATUS_FULLSCAN_STEP = C.SQLITE_STMTSTATUS_FULLSCAN_STEP // 1
	STMTSTATUS_SORT          = C.SQLITE_STMTSTATUS_SORT          // 2
	STMTSTATUS_AUTOINDEX     = C.SQLITE_STMTSTATUS_AUTOINDEX     // 3
	STMTSTATUS_VM_STEP       = C.SQLITE_STMTSTATUS_VM_STEP       // 4
)

// Per-connection limits that can be queried and changed with Conn.Limit.
// [http://www.sqlite.org/c3ref/c_limit_attached.html]
const (
	LIMIT_LENGTH              = C.SQLITE_LIMIT_LENGTH              // 0
	LIMIT_SQL_LENGTH          = C.SQLITE_LIMIT_SQL_LENGTH          // 1
	LIMIT_COLUMN              = C.SQLITE_LIMIT_COLUMN              // 2
	LIMIT_EXPR_DEPTH          = C.SQLITE_LIMIT_EXPR_DEPTH          // 3
	LIMIT_COMPOUND_SELECT     = C.SQLITE_LIMIT_COMPOUND_SELECT     // 4
	LIMIT_VDBE_OP             = C.SQLITE_LIMIT_VDBE_OP             // 5
	LIMIT_FUNCTION_ARG        = C.SQLITE_LIMIT_FUNCTION_ARG        // 6
	LIMIT_ATTACHED            = C.SQLITE_LIMIT_ATTACHED            // 7
	LIMIT_LIKE_PATTERN_LENGTH = C.SQLITE_LIMIT_LIKE_PATTERN_LENGTH // 8
	LIMIT_VARIABLE_NUMBER     = C.SQLITE_LIMIT_VARIABLE_NUMBER     // 9
	LIMIT_TRIGGER_DEPTH       = C.SQLITE_LIMIT_TRIGGER_DEPTH       // 10
)
