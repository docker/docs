package utils

import (
	"encoding/hex"
	"testing"

	"github.com/docker/notary/tuf/data"
	"github.com/stretchr/testify/require"
)

func TestFileMetaEqual(t *testing.T) {
	type test struct {
		name string
		b    data.FileMeta
		a    data.FileMeta
		err  func(test) error
	}
	fileMeta := func(length int64, hashes map[string]string) data.FileMeta {
		m := data.FileMeta{Length: length, Hashes: make(map[string][]byte, len(hashes))}
		for typ, hash := range hashes {
			v, err := hex.DecodeString(hash)
			require.NoError(t, err, "hash not in hex")
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
	for _, run := range tests {
		require.Equal(t, FileMetaEqual(run.a, run.b), run.err(run), "Files not equivalent")
	}
}

func TestNormalizeTarget(t *testing.T) {
	for before, after := range map[string]string{
		"":                    "/",
		"foo.txt":             "/foo.txt",
		"/bar.txt":            "/bar.txt",
		"foo//bar.txt":        "/foo/bar.txt",
		"/with/./a/dot":       "/with/a/dot",
		"/with/double/../dot": "/with/dot",
	} {
		require.Equal(t, NormalizeTarget(before), after, "Path normalization did not output expected.")
	}
}

func TestHashedPaths(t *testing.T) {
	hexBytes := func(s string) []byte {
		v, err := hex.DecodeString(s)
		require.NoError(t, err, "String was not hex")
		return v
	}
	hashes := data.Hashes{
		"sha512": hexBytes("abc123"),
		"sha256": hexBytes("def456"),
	}
	paths := HashedPaths("foo/bar.txt", hashes)
	// cannot use DeepEquals as the returned order is non-deterministic
	require.Len(t, paths, 2, "Expected 2 paths")
	expected := map[string]struct{}{"foo/abc123.bar.txt": {}, "foo/def456.bar.txt": {}}
	for _, path := range paths {
		if _, ok := expected[path]; !ok {
			t.Fatalf("unexpected path: %s", path)
		}
		delete(expected, path)
	}
}
