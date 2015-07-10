package utils

import (
	"encoding/hex"
	"testing"

	"github.com/endophage/gotuf/data"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type UtilSuite struct{}

var _ = Suite(&UtilSuite{})

func (UtilSuite) TestFileMetaEqual(c *C) {
	type test struct {
		name string
		b    data.FileMeta
		a    data.FileMeta
		err  func(test) error
	}
	fileMeta := func(length int64, hashes map[string]string) data.FileMeta {
		m := data.FileMeta{Length: length, Hashes: make(map[string]data.HexBytes, len(hashes))}
		for typ, hash := range hashes {
			v, err := hex.DecodeString(hash)
			c.Assert(err, IsNil)
			m.Hashes[typ] = v
		}
		return m
	}
	tests := []test{
		{
			name: "wrong length",
			a:    data.FileMeta{Length: 1},
			b:    data.FileMeta{Length: 2},
			err:  func(test) error { return ErrWrongLength },
		},
		{
			name: "wrong sha512 hash",
			a:    fileMeta(10, map[string]string{"sha512": "111111"}),
			b:    fileMeta(10, map[string]string{"sha512": "222222"}),
			err:  func(t test) error { return ErrWrongHash{"sha512", t.b.Hashes["sha512"], t.a.Hashes["sha512"]} },
		},
		{
			name: "intersecting hashes",
			a:    fileMeta(10, map[string]string{"sha512": "111111", "md5": "222222"}),
			b:    fileMeta(10, map[string]string{"sha512": "111111", "sha256": "333333"}),
			err:  func(test) error { return nil },
		},
		{
			name: "no common hashes",
			a:    fileMeta(10, map[string]string{"sha512": "111111"}),
			b:    fileMeta(10, map[string]string{"sha256": "222222", "md5": "333333"}),
			err:  func(t test) error { return ErrNoCommonHash{t.b.Hashes, t.a.Hashes} },
		},
	}
	for _, t := range tests {
		c.Assert(FileMetaEqual(t.a, t.b), DeepEquals, t.err(t), Commentf("name = %s", t.name))
	}
}

func (UtilSuite) TestNormalizeTarget(c *C) {
	for before, after := range map[string]string{
		"":                    "/",
		"foo.txt":             "/foo.txt",
		"/bar.txt":            "/bar.txt",
		"foo//bar.txt":        "/foo/bar.txt",
		"/with/./a/dot":       "/with/a/dot",
		"/with/double/../dot": "/with/dot",
	} {
		c.Assert(NormalizeTarget(before), Equals, after)
	}
}

func (UtilSuite) TestHashedPaths(c *C) {
	hexBytes := func(s string) data.HexBytes {
		v, err := hex.DecodeString(s)
		c.Assert(err, IsNil)
		return v
	}
	hashes := data.Hashes{
		"sha512": hexBytes("abc123"),
		"sha256": hexBytes("def456"),
	}
	paths := HashedPaths("foo/bar.txt", hashes)
	// cannot use DeepEquals as the returned order is non-deterministic
	c.Assert(paths, HasLen, 2)
	expected := map[string]struct{}{"foo/abc123.bar.txt": {}, "foo/def456.bar.txt": {}}
	for _, path := range paths {
		if _, ok := expected[path]; !ok {
			c.Fatalf("unexpected path: %s", path)
		}
		delete(expected, path)
	}
}
