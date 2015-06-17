package store

import (
	//	"fmt"
	"testing"

	"github.com/endophage/gotuf/data"
	"github.com/endophage/gotuf/testutils"
)

// TestDBStore just ensures we can initialize an empty store.
// Nothing to test, just ensure no crashes :-)
func TestDBStore(t *testing.T) {
	db := testutils.GetSqliteDB()
	defer testutils.FlushDB(db)
	_ = DBStore(db, "")
}

func TestLoadFiles(t *testing.T) {
	db := testutils.GetSqliteDB()
	defer testutils.FlushDB(db)
	store := DBStore(db, "docker.io/testImage")
	testmeta := testutils.SampleMeta()
	store.AddBlob("/foo.txt", testmeta)

	called := false
	check := func(path string, meta data.FileMeta) error {
		if called {
			t.Fatal("Store only has one item but check called > once.")
		} else {
			called = true
		}

		if path != "/foo.txt" {
			t.Fatal("Path is incorrect", path)
		}

		if meta.Length != testmeta.Length {
			t.Fatal("Length is incorrect")
		}

		if len(meta.Hashes) != len(testmeta.Hashes) {
			t.Fatal("Hashes map has been modified")
		}

		return nil
	}
	store.WalkStagedTargets([]string{}, check)
	if !called {
		t.Fatal("Walk func never called")
	}
}

func TestAddBlob(t *testing.T) {
	db := testutils.GetSqliteDB()
	defer testutils.FlushDB(db)
	store := DBStore(db, "docker.io/testImage")
	testmeta := testutils.SampleMeta()
	store.AddBlob("/foo.txt", testmeta)

	called := false
	check := func(path string, meta data.FileMeta) error {
		if called {
			t.Fatal("Store only has one item but check called > once.")
		} else {
			called = true
		}

		if path != "/foo.txt" {
			t.Fatal("Path is incorrect")
		}

		if meta.Length != 1 {
			t.Fatal("Length is incorrect")
		}

		sha256, ok256 := meta.Hashes["sha256"]
		sha512, ok512 := meta.Hashes["sha512"]
		if len(meta.Hashes) != 2 || !ok256 || !ok512 {
			t.Fatal("Hashes map has been modified")
		}

		hash := data.HexBytes{0x01, 0x02}
		if sha256[0] != hash[0] || sha256[1] != hash[1] {
			t.Fatal("SHA256 has been modified")
		}
		hash = data.HexBytes{0x03, 0x04}
		if sha512[0] != hash[0] || sha512[1] != hash[1] {
			t.Fatal("SHA512 has been modified")
		}
		return nil
	}

	store.WalkStagedTargets([]string{}, check)

	if !called {
		t.Fatal("Walk func never called")
	}
}

func TestRemoveBlob(t *testing.T) {
	testPath := "/foo.txt"
	db := testutils.GetSqliteDB()
	defer testutils.FlushDB(db)
	store := DBStore(db, "docker.io/testImage")
	meta := testutils.SampleMeta()

	store.AddBlob(testPath, meta)

	called := false
	check := func(path string, meta data.FileMeta) error {
		called = true
		return nil
	}

	store.RemoveBlob(testPath)

	store.WalkStagedTargets([]string{}, check)

	if called {
		t.Fatal("Walk func called on empty db")
	}

}

func TestLoadFilesWithPath(t *testing.T) {
	db := testutils.GetSqliteDB()
	defer testutils.FlushDB(db)
	store := DBStore(db, "docker.io/testImage")
	meta := testutils.SampleMeta()

	store.AddBlob("/foo.txt", meta)
	store.AddBlob("/bar.txt", meta)

	called := false
	check := func(path string, meta data.FileMeta) error {
		if called {
			t.Fatal("Store only has one item but check called > once.")
		} else {
			called = true
		}

		if path != "/foo.txt" {
			t.Fatal("Path is incorrect")
		}

		return nil
	}

	store.WalkStagedTargets([]string{"/foo.txt"}, check)

	if !called {
		t.Fatal("Walk func never called")
	}
}
