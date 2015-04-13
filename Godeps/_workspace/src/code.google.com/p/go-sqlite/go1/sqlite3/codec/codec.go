// Copyright 2013 The Go-SQLite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package codec

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"hash"
	"io"
	"strings"

	. "code.google.com/p/go-sqlite/go1/sqlite3"
)

func init() {
	RegisterCodec("aes-hmac", newAesHmac)
	RegisterCodec("hexdump", newHexDump)
}

// Errors returned by codec implementations.
var (
	codecErr = NewError(ERROR, "unspecified codec error")
	prngErr  = NewError(ERROR, "csprng not available")
	keyErr   = NewError(MISUSE, "invalid codec key format")
)

// parseKey extracts the codec name, options, and anything left over from a key
// in the format "<name>:<options>:<tail...>".
func parseKey(key []byte) (name string, opts map[string]string, tail []byte) {
	k := bytes.SplitN(key, []byte{':'}, 3)
	name = string(k[0])
	opts = make(map[string]string)
	if len(k) > 1 && len(k[1]) > 0 {
		for _, opt := range strings.Split(string(k[1]), ",") {
			if i := strings.Index(opt, "="); i > 0 {
				opts[opt[:i]] = opt[i+1:]
			} else {
				opts[opt] = ""
			}
		}
	}
	if len(k) > 2 && len(k[2]) > 0 {
		tail = k[2]
	}
	return
}

// hkdf implements the HMAC-based Key Derivation Function, as described in RFC
// 5869. The extract step is skipped if salt == nil. It is the caller's
// responsibility to set salt "to a string of HashLen zeros," if such behavior
// is desired. It returns the function that performs the expand step using the
// provided info value, which must be appendable. The derived key is valid until
// the next expansion.
func hkdf(ikm, salt []byte, dkLen int, h func() hash.Hash) func(info []byte) []byte {
	if salt != nil {
		prf := hmac.New(h, salt)
		prf.Write(ikm)
		ikm = prf.Sum(nil)
	}
	prf := hmac.New(h, ikm)
	hLen := prf.Size()
	n := (dkLen + hLen - 1) / hLen
	dk := make([]byte, dkLen, n*hLen)

	return func(info []byte) []byte {
		info = append(info, 0)
		ctr := &info[len(info)-1]
		for i, t := 1, dk[:0]; i <= n; i++ {
			*ctr = byte(i)
			prf.Reset()
			prf.Write(t)
			prf.Write(info)
			t = prf.Sum(t[len(t):])
		}
		return dk
	}
}

// rnd fills b with bytes from a CSPRNG.
func rnd(b []byte) bool {
	_, err := io.ReadFull(rand.Reader, b)
	return err == nil
}

// wipe overwrites b with zeros.
func wipe(b []byte) {
	for i := range b {
		b[i] = 0
	}
}

// suiteId constructs a canonical cipher suite identifier.
type suiteId struct {
	Cipher  string
	KeySize string
	Mode    string
	MAC     string
	Hash    string
	Trunc   string
}

func (s *suiteId) Id() []byte {
	id := make([]byte, 0, 64)
	section := func(parts ...string) {
		for i, p := range parts {
			if p != "" {
				parts = parts[i:]
				goto write
			}
		}
		return
	write:
		if len(id) > 0 {
			id = append(id, ',')
		}
		id = append(id, parts[0]...)
		for _, p := range parts[1:] {
			if p != "" {
				id = append(id, '-')
				id = append(id, p...)
			}
		}
	}
	section(s.Cipher, s.KeySize, s.Mode)
	section(s.MAC, s.Hash, s.Trunc)
	return id
}
