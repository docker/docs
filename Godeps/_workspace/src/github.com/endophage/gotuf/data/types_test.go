package data

import (
	"bytes"

	. "gopkg.in/check.v1"
)

type TypesSuite struct{}

var _ = Suite(&TypesSuite{})

func (TypesSuite) TestGenerateFileMetaDefault(c *C) {
	// default is sha512
	r := bytes.NewReader([]byte("foo"))
	meta, err := NewFileMeta(r, "sha512")
	c.Assert(err, IsNil)
	c.Assert(meta.Length, Equals, int64(3))
	hashes := meta.Hashes
	c.Assert(hashes, HasLen, 1)
	hash, ok := hashes["sha512"]
	if !ok {
		c.Fatal("missing sha512 hash")
	}
	c.Assert(hash.String(), DeepEquals, "f7fbba6e0636f890e56fbbf3283e524c6fa3204ae298382d624741d0dc6638326e282c41be5e4254d8820772c5518a2c5a8c0c7f7eda19594a7eb539453e1ed7")
}

func (TypesSuite) TestGenerateFileMetaExplicit(c *C) {
	r := bytes.NewReader([]byte("foo"))
	meta, err := NewFileMeta(r, "sha256", "sha512")
	c.Assert(err, IsNil)
	c.Assert(meta.Length, Equals, int64(3))
	hashes := meta.Hashes
	c.Assert(hashes, HasLen, 2)
	for name, val := range map[string]string{
		"sha256": "2c26b46b68ffc68ff99b453c1d30413413422d706483bfa0f98a5e886266e7ae",
		"sha512": "f7fbba6e0636f890e56fbbf3283e524c6fa3204ae298382d624741d0dc6638326e282c41be5e4254d8820772c5518a2c5a8c0c7f7eda19594a7eb539453e1ed7",
	} {
		hash, ok := hashes[name]
		if !ok {
			c.Fatalf("missing %s hash", name)
		}
		c.Assert(hash.String(), DeepEquals, val)
	}
}
