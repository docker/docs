package signature

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"github.com/stretchr/stew/strings"
	"github.com/stretchr/tracer"
	"io"
)

var HashWithKeysSeparator string = ":"

// HashFunc represents funcs that can hash a string.
type HashFunc func(s string) string

// Hash hashes a string using the current HashFunc.
//
// To tell Signature to use a different hashing algorithm, you
// just need to assign a different HashFunc to the Hash variable.
//
// To use the MD5 hash:
//
//     signature.Hash = signature.MD5Hash
//
// Or you can write your own hashing function:
//
//     signature.Hash = func(s string) string {
//	     // TODO: do your own hashing here
//     }
var Hash HashFunc = SHA1Hash

// SHA1Hash hashes a string using the SHA-1 hash algorithm as defined in RFC 3174.
var SHA1Hash HashFunc = func(s string) string {
	hash := sha1.New()
	hash.Write([]byte(s))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// MD5Hash hashes a string using the MD5 hash algorithm as defined in RFC 1321.
var MD5Hash HashFunc = func(s string) string {
	md5 := md5.New()
	io.WriteString(md5, s)
	return fmt.Sprintf("%x", md5.Sum(nil))
}

// HashWithKeys generates a hash of the specified bytes by first merging them with
// the specified private key.
//
// For format is:
//
//     BODY:PUBLIC:PRIVATE
//
// The keywords should be replaced with the actual bytes, but the colons are literals as the
// HashWithKeysSeparator value is used to separate the values.
//
// Useful for hashing non URLs (such as response bodies etc.)
func HashWithKeys(body, publicKey, privateKey []byte) string {
	return HashWithKeysWithTrace(body, publicKey, privateKey, nil)
}

// HashWithKey does the same as HashWithKey, but uses a single key.
func HashWithKey(body, key []byte) string {
	return HashWithKeyWithTrace(body, key, nil)
}

func HashWithKeysWithTrace(body, publicKey, privateKey []byte, t *tracer.Tracer) string {

	if t.Should(tracer.LevelDebug) {
		t.Trace(tracer.LevelDebug, "HashWithKeys: body=", body)
		t.Trace(tracer.LevelDebug, "HashWithKeys: publicKey=", publicKey)
		t.Trace(tracer.LevelDebug, "HashWithKeys: privateKey=", privateKey)
	}

	hash := Hash(string(strings.JoinBytes([]byte(HashWithKeysSeparator), body, publicKey, privateKey)))

	if t.Should(tracer.LevelDebug) {
		t.Trace(tracer.LevelDebug, "HashWithKeys: Output: %s", hash)
	}

	return hash

}

func HashWithKeyWithTrace(body, key []byte, t *tracer.Tracer) string {

	if t.Should(tracer.LevelDebug) {
		t.Trace(tracer.LevelDebug, "HashWithKeys: body=", body)
		t.Trace(tracer.LevelDebug, "HashWithKeys: key=", key)
	}

	hash := Hash(string(strings.JoinBytes([]byte(HashWithKeysSeparator), body, key)))

	if t.Should(tracer.LevelDebug) {
		t.Trace(tracer.LevelDebug, "HashWithKeys: Output: %s", hash)
	}

	return hash

}
