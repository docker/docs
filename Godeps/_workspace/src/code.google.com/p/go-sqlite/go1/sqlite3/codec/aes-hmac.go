// Copyright 2013 The Go-SQLite Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package codec

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"hash"

	. "code.google.com/p/go-sqlite/go1/sqlite3"
)

type aesHmac struct {
	key  []byte  // Key provided to newAesHmac with the master key removed
	buf  []byte  // Page encryption buffer
	hdr  [4]byte // Header included in each HMAC calculation (page number)
	tLen int     // Tag length in bytes (HMAC truncation)

	// Hash function and chaining mode constructors
	hash func() hash.Hash
	mode func(block cipher.Block, iv []byte) cipher.Stream

	// Block cipher and HMAC initialized from the master key
	block cipher.Block
	hmac  hash.Hash
}

func newAesHmac(ctx *CodecCtx, key []byte) (Codec, *Error) {
	name, opts, mk := parseKey(key)
	if len(mk) == 0 {
		return nil, keyErr
	}
	defer wipe(mk)

	// Configure the codec
	c := &aesHmac{
		key:  key[:len(key)-len(mk)],
		tLen: 16,
		hash: sha1.New,
		mode: cipher.NewCTR,
	}
	suite := suiteId{
		Cipher:  "aes",
		KeySize: "128",
		Mode:    "ctr",
		MAC:     "hmac",
		Hash:    "sha1",
		Trunc:   "128",
	}
	kLen := 16
	if err := c.config(opts, &suite, &kLen); err != nil {
		return nil, err
	}

	// Derive encryption and authentication keys
	hLen := c.hash().Size()
	salt := make([]byte, hLen)
	copy(salt, name)
	dk := hkdf(mk, salt, kLen+hLen, c.hash)(suite.Id())
	defer wipe(dk)

	// Initialize the block cipher and HMAC
	var err error
	if c.block, err = aes.NewCipher(dk[:kLen]); err != nil {
		return nil, NewError(MISUSE, err.Error())
	}
	c.hmac = hmac.New(c.hash, dk[kLen:])
	return c, nil
}

func (c *aesHmac) Reserve() int {
	return aes.BlockSize + c.tLen
}

func (c *aesHmac) Resize(pageSize, reserve int) {
	if reserve != c.Reserve() {
		panic("sqlite3: codec reserve value mismatch")
	}
	hLen := c.hash().Size()
	c.buf = make([]byte, pageSize, pageSize-c.tLen+hLen)
}

func (c *aesHmac) Encode(p []byte, n uint32, op int) ([]byte, *Error) {
	iv := c.pIV(c.buf)
	if !rnd(iv) {
		return nil, prngErr
	}
	c.mode(c.block, iv).XORKeyStream(c.buf, c.pText(p))
	if n == 1 {
		copy(c.buf[16:], p[16:24])
	}
	c.auth(c.buf, n, false)
	return c.buf, nil
}

func (c *aesHmac) Decode(p []byte, n uint32, op int) *Error {
	if !c.auth(p, n, true) {
		return codecErr
	}
	if n == 1 {
		copy(c.buf, p[16:24])
	}
	c.mode(c.block, c.pIV(p)).XORKeyStream(p, c.pText(p))
	if n == 1 {
		copy(p[16:24], c.buf)
	}
	return nil
}

func (c *aesHmac) Key() []byte {
	return c.key
}

func (c *aesHmac) Free() {
	c.buf = nil
	c.block = nil
	c.hmac = nil
}

// config applies the codec options that were provided in the key.
func (c *aesHmac) config(opts map[string]string, s *suiteId, kLen *int) *Error {
	for k := range opts {
		switch k {
		case "192":
			s.KeySize = k
			*kLen = 24
		case "256":
			s.KeySize = k
			*kLen = 32
		case "ofb":
			s.Mode = k
			c.mode = cipher.NewOFB
		case "sha256":
			s.Hash = k
			c.hash = sha256.New
		default:
			return NewError(MISUSE, "invalid codec option: "+k)
		}
	}
	return nil
}

// auth calculates and verifies the HMAC tag for page p. It returns true iff the
// tag is successfully verified.
func (c *aesHmac) auth(p []byte, n uint32, verify bool) bool {
	c.hdr[0] = byte(n >> 24)
	c.hdr[1] = byte(n >> 16)
	c.hdr[2] = byte(n >> 8)
	c.hdr[3] = byte(n)

	tag := c.pTag(c.buf)
	c.hmac.Reset()
	c.hmac.Write(c.hdr[:])
	c.hmac.Write(c.pAuth(p))
	c.hmac.Sum(tag[:0])

	return verify && hmac.Equal(tag, c.pTag(p))
}

// pAuth returns the page subslice that gets authenticated.
func (c *aesHmac) pAuth(p []byte) []byte {
	return p[:len(p)-c.tLen]
}

// pText returns the page subslice that gets encrypted.
func (c *aesHmac) pText(p []byte) []byte {
	return p[:len(p)-c.tLen-aes.BlockSize]
}

// pIV returns the page initialization vector.
func (c *aesHmac) pIV(p []byte) []byte {
	return p[len(p)-c.tLen-aes.BlockSize : len(p)-c.tLen]
}

// pTag returns the page authentication tag.
func (c *aesHmac) pTag(p []byte) []byte {
	return p[len(p)-c.tLen:]
}
