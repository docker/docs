// Copyright 2013 The Go-SQLite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package codec provides authenticated encryption and other codecs for the sqlite3
package.

This package has no public interface and should be imported with the blank
identifier to register the codecs. Use Conn.Key after opening a new connection,
or the KEY clause in an ATTACH statement, to use one of the codecs:

	c, _ := sqlite3.Open("file1.db")
	c.Key("main", []byte("aes-hmac::secretkey1"))
	c.Exec("ATTACH DATABASE 'file2.db' AS two KEY 'aes-hmac::secretkey2'")

If the KEY clause is omitted, SQLite uses the key from the main database, which
may no longer be valid depending on how the codec is implemented (e.g. aes-hmac
destroys the master key after initialization). Specify an empty string as the
key to disable this behavior.

Codec Operation

Each SQLite database and journal file consists of one or more pages of identical
size. Each page may have extra space reserved at the end, which SQLite will not
use in any way. The exact number of bytes reserved is stored at offset 20 of the
database file header, so the value is between 0 and 255. SQLite requires each
page to have at least 480 usable bytes, so the value cannot exceed 32 bytes with
a page size of 512. This extra space may be used by a codec to store per-page
Initialization Vectors (IVs), Message Authentication Codes (MACs), or any other
information.

CodecFunc is called to initialize a registered codec when a key with a matching
prefix is provided. If it returns a non-nil Codec implementation, Codec.Reserve
is called to determine how much space this codec needs reserved in each page for
correct operation. Codec.Resize is called to provide the current page size and
reserve values, and for all subsequent changes. The page size may be changed
before the database file is created. Once the first CREATE TABLE statement is
executed, the page size and reserve values are fixed.

Codec.Encode is called when a page is about to be written to the disk.
Codec.Decode is called when a page was just read from the disk. This happens for
both the main database file and the journal/WAL, so the pages are always encoded
on the disk and decoded in memory. Codec.Free is called to free all codec
resources when the database is detached.

AES-HMAC

The aes-hmac codec provides authenticated encryption using the Advanced
Encryption Standard (AES) cipher and the Hash-based Message Authentication Code
(HMAC) in Encrypt-then-MAC mode. Each page has an independent, pseudorandom IV,
which is regenerated every time the page is encrypted, and an authentication
tag, which is verified before the page is decrypted. The codec requires 32 bytes
per page to store this information.

The key format is "aes-hmac:<options>:<master-key>", where <options> is a
comma-separated list of codec options described below, and <master-key> is the
key from which separate encryption and authentication keys are derived.

SECURITY WARNING: The master key is called a "key" and not a "password" for a
reason. It is not passed through pbkdf2, bcrypt, scrypt, or any other key
stretching function. The application is expected to ensure that this key is
sufficiently resistant to brute-force attacks. Ideally, it should be obtained
from a cryptographically secure pseudorandom number generator (CSPRNG), such as
the one provided by the crypto/rand package.

The encryption and authentication keys are derived from the master key using the
HMAC-based Key Derivation Function (HKDF), as described in RFC 5869. The salt is
the codec name ("aes-hmac") extended with NULLs to HashLen bytes, and info is
the codec configuration string (e.g. "aes-128-ctr,hmac-sha1-128"). This is done
to obtain two keys of the required lengths, which are also bound to the codec
configuration.

The default configuration is AES-128-CTR cipher and HMAC-SHA1-128 authentication
(HMAC output is truncated to 128 bits). The following options may be used to
change the defaults:

	192
		AES-192 block cipher.
	256
		AES-256 block cipher.
	ofb
		Output feedback mode of operation.
	sha256
		SHA-256 hash function used by HKDF and HMAC.

For example, "aes-hmac:256,ofb,sha256:<master-key>" will use the AES-256-OFB
cipher and HMAC-SHA256-128 authentication.

HEXDUMP

The hexdump codec logs all method calls and dumps the page content for each
encode/decode operation to a file. It is intended to be used as an aid when
writing your own codecs.

The key format is "hexdump:<options>:<file>", where <options> is a
comma-separated list of codec options described below, and <file> is the output
destination. The default destination is stderr. Dash ("-") means stdout. For
obvious reasons, this codec cannot be used with an encrypted database except to
see the first Codec.Decode call for page 1.

The following options are supported:

	quiet
		Do not output a hex dump of each page.
	reserve=N
		Reserve N bytes in each page. The default is -1, which means don't
		change the current reserve value.
*/
package codec
