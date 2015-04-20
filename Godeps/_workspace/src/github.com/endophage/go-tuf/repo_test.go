package tuf

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/agl/ed25519"
	. "gopkg.in/check.v1"
	"github.com/endophage/go-tuf/data"
	"github.com/endophage/go-tuf/store"
	//	"github.com/endophage/go-tuf/encrypted"
	tuferr "github.com/endophage/go-tuf/errors"
	"github.com/endophage/go-tuf/signed"
	"github.com/endophage/go-tuf/util"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type RepoSuite struct{}

var _ = Suite(&RepoSuite{})

func (RepoSuite) TestNewRepo(c *C) {
	trust := signed.NewEd25519()

	meta := map[string]json.RawMessage{
		"root.json": []byte(`{
		  "signed": {
		    "_type": "root",
		    "version": 1,
		    "expires": "2015-12-26T03:26:55.821520874Z",
		    "keys": {},
		    "roles": {}
		  },
		  "signatures": []
		}`),
		"targets.json": []byte(`{
		  "signed": {
		    "_type": "targets",
		    "version": 1,
		    "expires": "2015-03-26T03:26:55.82155686Z",
		    "targets": {}
		  },
		  "signatures": []
		}`),
		"snapshot.json": []byte(`{
		  "signed": {
		    "_type": "snapshot",
		    "version": 1,
		    "expires": "2015-01-02T03:26:55.821585981Z",
		    "meta": {}
		  },
		  "signatures": []
		}`),
		"timestamp.json": []byte(`{
		  "signed": {
		    "_type": "timestamp",
		    "version": 1,
		    "expires": "2014-12-27T03:26:55.821599702Z",
		    "meta": {}
		  },
		  "signatures": []
		}`),
	}
	db := util.GetSqliteDB()
	defer util.FlushDB(db)
	local := store.DBStore(db, "")

	for k, v := range meta {
		local.SetMeta(k, v)
	}

	r, err := NewRepo(trust, local, "sha256")
	c.Assert(err, IsNil)

	root, err := r.root()
	c.Assert(err, IsNil)
	c.Assert(root.Type, Equals, "root")
	c.Assert(root.Version, Equals, 1)
	c.Assert(root.Keys, NotNil)
	c.Assert(root.Keys, HasLen, 0)

	targets, err := r.targets()
	c.Assert(err, IsNil)
	c.Assert(targets.Type, Equals, "targets")
	c.Assert(targets.Version, Equals, 1)
	c.Assert(targets.Targets, NotNil)
	c.Assert(targets.Targets, HasLen, 0)

	snapshot, err := r.snapshot()
	c.Assert(err, IsNil)
	c.Assert(snapshot.Type, Equals, "snapshot")
	c.Assert(snapshot.Version, Equals, 1)
	c.Assert(snapshot.Meta, NotNil)
	c.Assert(snapshot.Meta, HasLen, 0)

	timestamp, err := r.timestamp()
	c.Assert(err, IsNil)
	c.Assert(timestamp.Type, Equals, "timestamp")
	c.Assert(timestamp.Version, Equals, 1)
	c.Assert(timestamp.Meta, NotNil)
	c.Assert(timestamp.Meta, HasLen, 0)
}

func (RepoSuite) TestInit(c *C) {
	trust := signed.NewEd25519()

	db := util.GetSqliteDB()
	defer util.FlushDB(db)
	local := store.DBStore(
		db,
		"",
		//map[string][]byte{"/foo.txt": []byte("foo")},
	)
	local.AddBlob("/foo.txt", util.SampleMeta())

	r, err := NewRepo(trust, local, "sha256")
	c.Assert(err, IsNil)

	// Init() sets root.ConsistentSnapshot
	for _, v := range []bool{true, false} {
		c.Assert(r.Init(v), IsNil)
		root, err := r.root()
		c.Assert(err, IsNil)
		c.Assert(root.ConsistentSnapshot, Equals, v)
	}

	// Init() fails if targets have been added
	c.Assert(r.AddTarget("foo.txt", nil), IsNil)
	c.Assert(r.Init(true), Equals, tuferr.ErrInitNotAllowed)
}

func genKey(c *C, r *Repo, role string) string {
	id, err := r.GenKey(role)
	c.Assert(err, IsNil)
	return id
}

func (RepoSuite) TestGenKey(c *C) {
	trust := signed.NewEd25519()

	sqldb := util.GetSqliteDB()
	defer util.FlushDB(sqldb)
	local := store.DBStore(sqldb, "")
	r, err := NewRepo(trust, local, "sha256")
	c.Assert(err, IsNil)

	// generate a key for an unknown role
	_, err = r.GenKey("foo")
	c.Assert(err, Equals, tuferr.ErrInvalidRole{"foo"})

	// generate a root key
	id := genKey(c, r, "root")

	// check root metadata is correct
	root, err := r.root()
	c.Assert(err, IsNil)
	c.Assert(root.Roles, NotNil)
	c.Assert(root.Roles, HasLen, 1)
	c.Assert(root.Keys, NotNil)
	c.Assert(root.Keys, HasLen, 1)
	rootRole, ok := root.Roles["root"]
	if !ok {
		c.Fatal("missing root role")
	}
	c.Assert(rootRole.KeyIDs, HasLen, 1)
	keyID := rootRole.KeyIDs[0]
	c.Assert(keyID, Equals, id)
	k, ok := root.Keys[keyID]
	if !ok {
		c.Fatal("missing key")
	}
	c.Assert(k.ID(), Equals, keyID)
	c.Assert(k.Value.Public, HasLen, ed25519.PublicKeySize)
	//c.Assert(k.Value.Private, IsNil)

	// check root key + role are in db
	db, err := r.db()
	c.Assert(err, IsNil)
	rootKey := db.GetKey(keyID)
	c.Assert(rootKey, NotNil)
	c.Assert(rootKey.ID, Equals, keyID)
	role := db.GetRole("root")
	c.Assert(role.KeyIDs, DeepEquals, []string{keyID})

	// check the key was saved correctly
	localKeys, err := local.GetKeys("root")
	c.Assert(err, IsNil)
	c.Assert(localKeys, HasLen, 1)
	c.Assert(localKeys[0].ID(), Equals, keyID)

	// check RootKeys() is correct
	rootKeys, err := r.RootKeys()
	c.Assert(err, IsNil)
	c.Assert(rootKeys, HasLen, 1)
	c.Assert(rootKeys[0].ID(), Equals, rootKey.ID)
	c.Assert(rootKeys[0].Value.Public, DeepEquals, rootKey.Key.Value.Public)
	//c.Assert(rootKeys[0].Value.Private, IsNil)

	// generate two targets keys
	genKey(c, r, "targets")
	genKey(c, r, "targets")

	// check root metadata is correct
	root, err = r.root()
	c.Assert(err, IsNil)
	c.Assert(root.Roles, HasLen, 2)
	c.Assert(root.Keys, HasLen, 3)
	targetsRole, ok := root.Roles["targets"]
	if !ok {
		c.Fatal("missing targets role")
	}
	c.Assert(targetsRole.KeyIDs, HasLen, 2)
	targetKeyIDs := make([]string, 0, 2)
	db, err = r.db()
	c.Assert(err, IsNil)
	for _, id := range targetsRole.KeyIDs {
		targetKeyIDs = append(targetKeyIDs, id)
		_, ok = root.Keys[id]
		if !ok {
			c.Fatal("missing key")
		}
		key := db.GetKey(id)
		c.Assert(key, NotNil)
		c.Assert(key.ID, Equals, id)
	}
	role = db.GetRole("targets")
	c.Assert(role.KeyIDs, DeepEquals, targetKeyIDs)

	// check RootKeys() is unchanged
	rootKeys, err = r.RootKeys()
	c.Assert(err, IsNil)
	c.Assert(rootKeys, HasLen, 1)
	c.Assert(rootKeys[0].ID(), Equals, rootKey.ID)

	// check the keys were saved correctly
	localKeys, err = local.GetKeys("targets")
	c.Assert(err, IsNil)
	c.Assert(localKeys, HasLen, 2)
	for _, key := range localKeys {
		found := false
		for _, id := range targetsRole.KeyIDs {
			if id == key.ID() {
				found = true
			}
		}
		if !found {
			c.Fatal("missing key")
		}
	}

	// check root.json got staged
	meta, err := local.GetMeta()
	c.Assert(err, IsNil)
	rootJSON, ok := meta["root.json"]
	if !ok {
		c.Fatal("missing root metadata")
	}
	s := &data.Signed{}
	c.Assert(json.Unmarshal(rootJSON, s), IsNil)
	stagedRoot := &data.Root{}
	c.Assert(json.Unmarshal(s.Signed, stagedRoot), IsNil)
	c.Assert(stagedRoot.Type, Equals, root.Type)
	c.Assert(stagedRoot.Version, Equals, root.Version)
	c.Assert(stagedRoot.Expires.UnixNano(), Equals, root.Expires.UnixNano())
	c.Assert(stagedRoot.Keys, DeepEquals, root.Keys)
	c.Assert(stagedRoot.Roles, DeepEquals, root.Roles)
}

func (RepoSuite) TestRevokeKey(c *C) {
	trust := signed.NewEd25519()

	db := util.GetSqliteDB()
	defer util.FlushDB(db)
	local := store.DBStore(db, "")
	r, err := NewRepo(trust, local, "sha256")
	c.Assert(err, IsNil)

	// revoking a key for an unknown role returns ErrInvalidRole
	c.Assert(r.RevokeKey("foo", ""), DeepEquals, tuferr.ErrInvalidRole{"foo"})

	// revoking a key which doesn't exist returns ErrKeyNotFound
	c.Assert(r.RevokeKey("root", "nonexistent"), DeepEquals, tuferr.ErrKeyNotFound{"root", "nonexistent"})

	// generate keys
	genKey(c, r, "root")
	genKey(c, r, "targets")
	genKey(c, r, "targets")
	genKey(c, r, "snapshot")
	genKey(c, r, "timestamp")
	root, err := r.root()
	c.Assert(err, IsNil)
	c.Assert(root.Roles, NotNil)
	c.Assert(root.Roles, HasLen, 4)
	c.Assert(root.Keys, NotNil)
	c.Assert(root.Keys, HasLen, 5)

	// revoke a key
	targetsRole, ok := root.Roles["targets"]
	if !ok {
		c.Fatal("missing targets role")
	}
	c.Assert(targetsRole.KeyIDs, HasLen, 2)
	id := targetsRole.KeyIDs[0]
	c.Assert(r.RevokeKey("targets", id), IsNil)

	// check root was updated
	root, err = r.root()
	c.Assert(err, IsNil)
	c.Assert(root.Roles, NotNil)
	c.Assert(root.Roles, HasLen, 4)
	c.Assert(root.Keys, NotNil)
	c.Assert(root.Keys, HasLen, 4)
	targetsRole, ok = root.Roles["targets"]
	if !ok {
		c.Fatal("missing targets role")
	}
	c.Assert(targetsRole.KeyIDs, HasLen, 1)
	c.Assert(targetsRole.KeyIDs[0], Not(Equals), id)
}

func (RepoSuite) TestSign(c *C) {
	trust := signed.NewEd25519()

	baseMeta := map[string]json.RawMessage{"root.json": []byte(`{"signed":{},"signatures":[]}`)}
	db := util.GetSqliteDB()
	defer util.FlushDB(db)
	local := store.DBStore(db, "")
	local.SetMeta("root.json", baseMeta["root.json"])
	r, err := NewRepo(trust, local, "sha256")
	c.Assert(err, IsNil)

	// signing with no keys returns ErrInsufficientKeys
	c.Assert(r.Sign("root.json"), Equals, tuferr.ErrInsufficientKeys{"root.json"})

	checkSigIDs := func(keyIDs ...string) {
		meta, err := local.GetMeta()
		if err != nil {
			c.Fatal("failed to retrieve meta")
		}
		rootJSON, ok := meta["root.json"]
		if !ok {
			c.Fatal("missing root.json")
		}
		s := &data.Signed{}
		c.Assert(json.Unmarshal(rootJSON, s), IsNil)
		c.Assert(s.Signatures, HasLen, len(keyIDs))
		for i, id := range keyIDs {
			c.Assert(s.Signatures[i].KeyID, Equals, id)
		}
	}

	// signing with an available key generates a signature
	//key, err := signer.Create()
	kID, err := r.GenKey("root")
	c.Assert(err, IsNil)
	//c.Assert(local.SaveKey("root", key.SerializePrivate()), IsNil)
	c.Assert(r.Sign("root.json"), IsNil)
	checkSigIDs(kID)

	// signing again does not generate a duplicate signature
	c.Assert(r.Sign("root.json"), IsNil)
	checkSigIDs(kID)

	// signing with a new available key generates another signature
	//newKey, err := signer.Create()
	newkID, err := r.GenKey("root")
	c.Assert(err, IsNil)
	//c.Assert(local.SaveKey("root", newKey.SerializePrivate()), IsNil)
	c.Assert(r.Sign("root.json"), IsNil)
	checkSigIDs(kID, newkID)
}

func (RepoSuite) TestCommit(c *C) {
	trust := signed.NewEd25519()

	//files := map[string][]byte{"/foo.txt": []byte("foo"), "/bar.txt": []byte("bar")}
	db := util.GetSqliteDB()
	defer util.FlushDB(db)
	local := store.DBStore(db, "")
	r, err := NewRepo(trust, local, "sha256")
	c.Assert(err, IsNil)

	// commit without root.json
	c.Assert(r.Commit(), DeepEquals, tuferr.ErrMissingMetadata{"root.json"})

	// commit without targets.json
	genKey(c, r, "root")
	c.Assert(r.Commit(), DeepEquals, tuferr.ErrMissingMetadata{"targets.json"})

	// commit without snapshot.json
	genKey(c, r, "targets")
	local.AddBlob("/foo.txt", util.SampleMeta())
	c.Assert(r.AddTarget("foo.txt", nil), IsNil)
	c.Assert(r.Commit(), DeepEquals, tuferr.ErrMissingMetadata{"snapshot.json"})

	// commit without timestamp.json
	genKey(c, r, "snapshot")
	c.Assert(r.Snapshot(CompressionTypeNone), IsNil)
	c.Assert(r.Commit(), DeepEquals, tuferr.ErrMissingMetadata{"timestamp.json"})

	// commit with timestamp.json but no timestamp key
	c.Assert(r.Timestamp(), IsNil)
	c.Assert(r.Commit(), DeepEquals, tuferr.ErrInsufficientSignatures{"timestamp.json", signed.ErrNoSignatures})

	// commit success
	genKey(c, r, "timestamp")
	c.Assert(r.Snapshot(CompressionTypeNone), IsNil)
	c.Assert(r.Timestamp(), IsNil)
	c.Assert(r.Commit(), IsNil)

	// commit with an invalid root hash in snapshot.json due to new key creation
	genKey(c, r, "targets")
	c.Assert(r.Sign("targets.json"), IsNil)
	c.Assert(r.Commit(), DeepEquals, errors.New("tuf: invalid root.json in snapshot.json: wrong length"))

	// commit with an invalid targets hash in snapshot.json
	c.Assert(r.Snapshot(CompressionTypeNone), IsNil)
	local.AddBlob("/bar.txt", util.SampleMeta())
	c.Assert(r.AddTarget("bar.txt", nil), IsNil)
	c.Assert(r.Commit(), DeepEquals, errors.New("tuf: invalid targets.json in snapshot.json: wrong length"))

	// commit with an invalid timestamp
	c.Assert(r.Snapshot(CompressionTypeNone), IsNil)
	// TODO: Change this test once Snapshot() supports compression and we
	//       can guarantee the error will end in "wrong length" by
	//       compressing a file and thus changing the size of snapshot.json
	err = r.Commit()
	c.Assert(err, NotNil)
	c.Assert(err.Error()[0:44], Equals, "tuf: invalid snapshot.json in timestamp.json")

	// commit with a role's threshold greater than number of keys
	root, err := r.root()
	c.Assert(err, IsNil)
	role, ok := root.Roles["timestamp"]
	if !ok {
		c.Fatal("missing timestamp role")
	}
	c.Assert(role.KeyIDs, HasLen, 1)
	c.Assert(role.Threshold, Equals, 1)
	c.Assert(r.RevokeKey("timestamp", role.KeyIDs[0]), IsNil)
	c.Assert(r.Snapshot(CompressionTypeNone), IsNil)
	c.Assert(r.Timestamp(), IsNil)
	c.Assert(r.Commit(), DeepEquals, tuferr.ErrNotEnoughKeys{"timestamp", 0, 1})
}

type tmpDir struct {
	path string
	c    *C
}

func newTmpDir(c *C) *tmpDir {
	return &tmpDir{path: c.MkDir(), c: c}
}

func (t *tmpDir) assertExists(path string) {
	if _, err := os.Stat(filepath.Join(t.path, path)); os.IsNotExist(err) {
		t.c.Fatalf("expected path to exist but it doesn't: %s", path)
	}
}

func (t *tmpDir) assertNotExist(path string) {
	if _, err := os.Stat(filepath.Join(t.path, path)); !os.IsNotExist(err) {
		t.c.Fatalf("expected path to not exist but it does: %s", path)
	}
}

func (t *tmpDir) assertHashedFilesExist(path string, hashes data.Hashes) {
	t.c.Assert(len(hashes) > 0, Equals, true)
	for _, path := range util.HashedPaths(path, hashes) {
		t.assertExists(path)
	}
}

func (t *tmpDir) assertHashedFilesNotExist(path string, hashes data.Hashes) {
	t.c.Assert(len(hashes) > 0, Equals, true)
	for _, path := range util.HashedPaths(path, hashes) {
		t.assertNotExist(path)
	}
}

func (t *tmpDir) assertEmpty(dir string) {
	path := filepath.Join(t.path, dir)
	f, err := os.Stat(path)
	if os.IsNotExist(err) {
		t.c.Fatalf("expected dir to exist but it doesn't: %s", dir)
	}
	t.c.Assert(err, IsNil)
	t.c.Assert(f.IsDir(), Equals, true)
	entries, err := ioutil.ReadDir(path)
	t.c.Assert(err, IsNil)
	// check that all (if any) entries are also empty
	for _, e := range entries {
		t.assertEmpty(filepath.Join(dir, e.Name()))
	}
}

func (t *tmpDir) assertFileContent(path, content string) {
	actual := t.readFile(path)
	t.c.Assert(string(actual), Equals, content)
}

func (t *tmpDir) stagedTargetPath(path string) string {
	return filepath.Join(t.path, "staged", "targets", path)
}

func (t *tmpDir) writeStagedTarget(path, data string) {
	path = t.stagedTargetPath(path)
	t.c.Assert(os.MkdirAll(filepath.Dir(path), 0755), IsNil)
	t.c.Assert(ioutil.WriteFile(path, []byte(data), 0644), IsNil)
}

func (t *tmpDir) readFile(path string) []byte {
	t.assertExists(path)
	data, err := ioutil.ReadFile(filepath.Join(t.path, path))
	t.c.Assert(err, IsNil)
	return data
}

func (RepoSuite) TestCommitFileSystem(c *C) {
	trust := signed.NewEd25519()
	tmp := newTmpDir(c)
	local := store.FileSystemStore(tmp.path, nil)
	r, err := NewRepo(trust, local, "sha256")
	c.Assert(err, IsNil)

	// don't use consistent snapshots to make the checks simpler
	c.Assert(r.Init(false), IsNil)

	// generating keys should stage root.json and create repo dirs
	genKey(c, r, "root")
	genKey(c, r, "targets")
	genKey(c, r, "snapshot")
	genKey(c, r, "timestamp")
	tmp.assertExists("staged/root.json")
	tmp.assertEmpty("repository")
	tmp.assertEmpty("staged/targets")

	// adding a non-existent file fails
	c.Assert(r.AddTarget("foo.txt", nil), Equals, tuferr.ErrFileNotFound{tmp.stagedTargetPath("foo.txt")})
	tmp.assertEmpty("repository")

	// adding a file stages targets.json
	tmp.writeStagedTarget("foo.txt", "foo")
	c.Assert(r.AddTarget("foo.txt", nil), IsNil)
	tmp.assertExists("staged/targets.json")
	tmp.assertEmpty("repository")
	t, err := r.targets()
	c.Assert(err, IsNil)
	c.Assert(t.Targets, HasLen, 1)
	if _, ok := t.Targets["/foo.txt"]; !ok {
		c.Fatal("missing target file: /foo.txt")
	}

	// Snapshot() stages snapshot.json
	c.Assert(r.Snapshot(CompressionTypeNone), IsNil)
	tmp.assertExists("staged/snapshot.json")
	tmp.assertEmpty("repository")

	// Timestamp() stages timestamp.json
	c.Assert(r.Timestamp(), IsNil)
	tmp.assertExists("staged/timestamp.json")
	tmp.assertEmpty("repository")

	// committing moves files from staged -> repository
	c.Assert(r.Commit(), IsNil)
	tmp.assertExists("repository/root.json")
	tmp.assertExists("repository/targets.json")
	tmp.assertExists("repository/snapshot.json")
	tmp.assertExists("repository/timestamp.json")
	tmp.assertFileContent("repository/targets/foo.txt", "foo")
	tmp.assertEmpty("staged/targets")
	tmp.assertEmpty("staged")

	// adding and committing another file moves it into repository/targets
	tmp.writeStagedTarget("path/to/bar.txt", "bar")
	c.Assert(r.AddTarget("path/to/bar.txt", nil), IsNil)
	tmp.assertExists("staged/targets.json")
	c.Assert(r.Snapshot(CompressionTypeNone), IsNil)
	c.Assert(r.Timestamp(), IsNil)
	c.Assert(r.Commit(), IsNil)
	tmp.assertFileContent("repository/targets/foo.txt", "foo")
	tmp.assertFileContent("repository/targets/path/to/bar.txt", "bar")
	tmp.assertEmpty("staged/targets")
	tmp.assertEmpty("staged")

	// removing and committing a file removes it from repository/targets
	c.Assert(r.RemoveTarget("foo.txt"), IsNil)
	tmp.assertExists("staged/targets.json")
	c.Assert(r.Snapshot(CompressionTypeNone), IsNil)
	c.Assert(r.Timestamp(), IsNil)
	c.Assert(r.Commit(), IsNil)
	tmp.assertNotExist("repository/targets/foo.txt")
	tmp.assertFileContent("repository/targets/path/to/bar.txt", "bar")
	tmp.assertEmpty("staged/targets")
	tmp.assertEmpty("staged")
}

func (RepoSuite) TestConsistentSnapshot(c *C) {
	trust := signed.NewEd25519()
	tmp := newTmpDir(c)
	local := store.FileSystemStore(tmp.path, nil)
	r, err := NewRepo(trust, local, "sha512", "sha256")
	c.Assert(err, IsNil)

	genKey(c, r, "root")
	genKey(c, r, "targets")
	genKey(c, r, "snapshot")
	genKey(c, r, "timestamp")
	tmp.writeStagedTarget("foo.txt", "foo")
	c.Assert(r.AddTarget("foo.txt", nil), IsNil)
	tmp.writeStagedTarget("dir/bar.txt", "bar")
	c.Assert(r.AddTarget("dir/bar.txt", nil), IsNil)
	c.Assert(r.Snapshot(CompressionTypeNone), IsNil)
	c.Assert(r.Timestamp(), IsNil)
	c.Assert(r.Commit(), IsNil)

	hashes, err := r.fileHashes()
	c.Assert(err, IsNil)

	// root.json, targets.json and snapshot.json should exist at both hashed and unhashed paths
	for _, path := range []string{"root.json", "targets.json", "snapshot.json"} {
		repoPath := filepath.Join("repository", path)
		tmp.assertHashedFilesExist(repoPath, hashes[path])
		tmp.assertExists(repoPath)
	}

	// target files should exist at hashed but not unhashed paths
	for _, path := range []string{"targets/foo.txt", "targets/dir/bar.txt"} {
		repoPath := filepath.Join("repository", path)
		tmp.assertHashedFilesExist(repoPath, hashes[path])
		tmp.assertNotExist(repoPath)
	}

	// timestamp.json should exist at an unhashed path (it doesn't have a hash)
	tmp.assertExists("repository/timestamp.json")

	// removing a file should remove the hashed files
	c.Assert(r.RemoveTarget("foo.txt"), IsNil)
	c.Assert(r.Snapshot(CompressionTypeNone), IsNil)
	c.Assert(r.Timestamp(), IsNil)
	c.Assert(r.Commit(), IsNil)
	tmp.assertHashedFilesNotExist("repository/targets/foo.txt", hashes["targets/foo.txt"])
	tmp.assertNotExist("repository/targets/foo.txt")

	// targets should be returned by new repo
	newRepo, err := NewRepo(trust, local, "sha512", "sha256")
	c.Assert(err, IsNil)
	t, err := newRepo.targets()
	c.Assert(err, IsNil)
	c.Assert(t.Targets, HasLen, 1)
	if _, ok := t.Targets["/dir/bar.txt"]; !ok {
		c.Fatal("missing targets file: dir/bar.txt")
	}
}

func (RepoSuite) TestExpiresAndVersion(c *C) {
	trust := signed.NewEd25519()

	//files := map[string][]byte{"/foo.txt": []byte("foo")}
	db := util.GetSqliteDB()
	defer util.FlushDB(db)
	local := store.DBStore(db, "")
	r, err := NewRepo(trust, local, "sha256")
	c.Assert(err, IsNil)

	past := time.Now().Add(-1 * time.Second)
	_, genKeyErr := r.GenKeyWithExpires("root", past)
	for _, err := range []error{
		genKeyErr,
		r.AddTargetWithExpires("foo.txt", nil, past),
		r.RemoveTargetWithExpires("foo.txt", past),
		r.SnapshotWithExpires(CompressionTypeNone, past),
		r.TimestampWithExpires(past),
	} {
		c.Assert(err, Equals, tuferr.ErrInvalidExpires{past})
	}

	genKey(c, r, "root")
	root, err := r.root()
	c.Assert(err, IsNil)
	c.Assert(root.Version, Equals, 1)

	expires := time.Now().Add(24 * time.Hour)
	_, err = r.GenKeyWithExpires("root", expires)
	c.Assert(err, IsNil)
	root, err = r.root()
	c.Assert(err, IsNil)
	c.Assert(root.Expires.Unix(), DeepEquals, expires.Round(time.Second).Unix())
	c.Assert(root.Version, Equals, 2)

	expires = time.Now().Add(12 * time.Hour)
	role, ok := root.Roles["root"]
	if !ok {
		c.Fatal("missing root role")
	}
	c.Assert(role.KeyIDs, HasLen, 2)
	c.Assert(r.RevokeKeyWithExpires("root", role.KeyIDs[0], expires), IsNil)
	root, err = r.root()
	c.Assert(err, IsNil)
	c.Assert(root.Expires.Unix(), DeepEquals, expires.Round(time.Second).Unix())
	c.Assert(root.Version, Equals, 3)

	expires = time.Now().Add(6 * time.Hour)
	genKey(c, r, "targets")
	local.AddBlob("/foo.txt", util.SampleMeta())
	c.Assert(r.AddTargetWithExpires("foo.txt", nil, expires), IsNil)
	targets, err := r.targets()
	c.Assert(err, IsNil)
	c.Assert(targets.Expires.Unix(), Equals, expires.Round(time.Second).Unix())
	c.Assert(targets.Version, Equals, 1)

	expires = time.Now().Add(2 * time.Hour)
	c.Assert(r.RemoveTargetWithExpires("foo.txt", expires), IsNil)
	targets, err = r.targets()
	c.Assert(err, IsNil)
	c.Assert(targets.Expires.Unix(), Equals, expires.Round(time.Second).Unix())
	c.Assert(targets.Version, Equals, 2)

	expires = time.Now().Add(time.Hour)
	genKey(c, r, "snapshot")
	c.Assert(r.SnapshotWithExpires(CompressionTypeNone, expires), IsNil)
	snapshot, err := r.snapshot()
	c.Assert(err, IsNil)
	c.Assert(snapshot.Expires.Unix(), Equals, expires.Round(time.Second).Unix())
	c.Assert(snapshot.Version, Equals, 1)

	c.Assert(r.Snapshot(CompressionTypeNone), IsNil)
	snapshot, err = r.snapshot()
	c.Assert(err, IsNil)
	c.Assert(snapshot.Version, Equals, 2)

	expires = time.Now().Add(10 * time.Minute)
	genKey(c, r, "timestamp")
	c.Assert(r.TimestampWithExpires(expires), IsNil)
	timestamp, err := r.timestamp()
	c.Assert(err, IsNil)
	c.Assert(timestamp.Expires.Unix(), Equals, expires.Round(time.Second).Unix())
	c.Assert(timestamp.Version, Equals, 1)

	c.Assert(r.Timestamp(), IsNil)
	timestamp, err = r.timestamp()
	c.Assert(err, IsNil)
	c.Assert(timestamp.Version, Equals, 2)
}

func (RepoSuite) TestHashAlgorithm(c *C) {
	trust := signed.NewEd25519()

	//files := map[string][]byte{"/foo.txt": []byte("foo")}
	db := util.GetSqliteDB()
	defer util.FlushDB(db)
	local := store.DBStore(db, "docker.io/testImage")
	type hashTest struct {
		args     []string
		expected []string
	}
	for _, test := range []hashTest{
		{args: []string{}, expected: []string{"sha512"}},
		{args: []string{"sha256"}},
		{args: []string{"sha512", "sha256"}},
	} {
		// generate metadata with specific hash functions
		r, err := NewRepo(trust, local, test.args...)
		c.Assert(err, IsNil)
		genKey(c, r, "root")
		genKey(c, r, "targets")
		genKey(c, r, "snapshot")
		local.AddBlob("/foo.txt", util.SampleMeta())
		c.Assert(r.AddTarget("foo.txt", nil), IsNil)
		c.Assert(r.Snapshot(CompressionTypeNone), IsNil)
		c.Assert(r.Timestamp(), IsNil)

		// check metadata has correct hash functions
		if test.expected == nil {
			test.expected = test.args
		}
		targets, err := r.targets()
		c.Assert(err, IsNil)
		snapshot, err := r.snapshot()
		c.Assert(err, IsNil)
		timestamp, err := r.timestamp()
		c.Assert(err, IsNil)
		for name, file := range map[string]data.FileMeta{
			"foo.txt":       targets.Targets["/foo.txt"],
			"root.json":     snapshot.Meta["root.json"],
			"targets.json":  snapshot.Meta["targets.json"],
			"snapshot.json": timestamp.Meta["snapshot.json"],
		} {
			for _, hashAlgorithm := range test.expected {
				if _, ok := file.Hashes[hashAlgorithm]; !ok {
					c.Fatalf("expected %s hash to contain hash func %s, got %s", name, hashAlgorithm, file.HashAlgorithms())
				}
			}
		}
	}
}

func testPassphraseFunc(p []byte) util.PassphraseFunc {
	return func(string, bool) ([]byte, error) { return p, nil }
}

//func (RepoSuite) TestKeyPersistence(c *C) {
//	tmp := newTmpDir(c)
//	passphrase := []byte("s3cr3t")
//	store := FileSystemStore(tmp.path, testPassphraseFunc(passphrase))
//
//	assertEqual := func(actual []*data.Key, expected []*keys.Key) {
//		c.Assert(actual, HasLen, len(expected))
//		for i, key := range expected {
//			c.Assert(actual[i].ID(), Equals, key.ID)
//			c.Assert(actual[i].Value.Public, DeepEquals, data.HexBytes(key.Public[:]))
//			c.Assert(actual[i].Value.Private, DeepEquals, data.HexBytes(key.Private[:]))
//		}
//	}
//
//	assertKeys := func(role string, enc bool, expected []*keys.Key) {
//		keysJSON := tmp.readFile("keys/" + role + ".json")
//		pk := &persistedKeys{}
//		c.Assert(json.Unmarshal(keysJSON, pk), IsNil)
//
//		// check the persisted keys are correct
//		var actual []*data.Key
//		if enc {
//			c.Assert(pk.Encrypted, Equals, true)
//			decrypted, err := encrypted.Decrypt(pk.Data, passphrase)
//			c.Assert(err, IsNil)
//			c.Assert(json.Unmarshal(decrypted, &actual), IsNil)
//		} else {
//			c.Assert(pk.Encrypted, Equals, false)
//			c.Assert(json.Unmarshal(pk.Data, &actual), IsNil)
//		}
//		assertEqual(actual, expected)
//
//		// check GetKeys is correct
//		actual, err := store.GetKeys(role)
//		c.Assert(err, IsNil)
//		assertEqual(actual, expected)
//	}
//
//	// save a key and check it gets encrypted
//	key, err := keys.NewKey()
//	c.Assert(err, IsNil)
//	c.Assert(store.SaveKey("root", key.SerializePrivate()), IsNil)
//	assertKeys("root", true, []*keys.Key{key})
//
//	// save another key and check it gets added to the existing keys
//	newKey, err := keys.NewKey()
//	c.Assert(err, IsNil)
//	c.Assert(store.SaveKey("root", newKey.SerializePrivate()), IsNil)
//	assertKeys("root", true, []*keys.Key{key, newKey})
//
//	// check saving a key to an encrypted file without a passphrase fails
//	insecureStore := FileSystemStore(tmp.path, nil)
//	key, err = keys.NewKey()
//	c.Assert(err, IsNil)
//	c.Assert(insecureStore.SaveKey("root", key.SerializePrivate()), Equals, ErrPassphraseRequired{"root"})
//
//	// save a key to an insecure store and check it is not encrypted
//	key, err = keys.NewKey()
//	c.Assert(err, IsNil)
//	c.Assert(insecureStore.SaveKey("targets", key.SerializePrivate()), IsNil)
//	assertKeys("targets", false, []*keys.Key{key})
//}

func (RepoSuite) TestManageMultipleTargets(c *C) {
	trust := signed.NewEd25519()
	tmp := newTmpDir(c)
	local := store.FileSystemStore(tmp.path, nil)
	r, err := NewRepo(trust, local)
	c.Assert(err, IsNil)
	// don't use consistent snapshots to make the checks simpler
	c.Assert(r.Init(false), IsNil)
	genKey(c, r, "root")
	genKey(c, r, "targets")
	genKey(c, r, "snapshot")
	genKey(c, r, "timestamp")

	assertRepoTargets := func(paths ...string) {
		t, err := r.targets()
		c.Assert(err, IsNil)
		for _, path := range paths {
			if _, ok := t.Targets[path]; !ok {
				c.Fatalf("missing target file: %s", path)
			}
		}
	}

	// adding and committing multiple files moves correct targets from staged -> repository
	tmp.writeStagedTarget("foo.txt", "foo")
	tmp.writeStagedTarget("bar.txt", "bar")
	c.Assert(r.AddTargets([]string{"foo.txt", "bar.txt"}, nil), IsNil)
	c.Assert(r.Snapshot(CompressionTypeNone), IsNil)
	c.Assert(r.Timestamp(), IsNil)
	c.Assert(r.Commit(), IsNil)
	assertRepoTargets("/foo.txt", "/bar.txt")
	tmp.assertExists("repository/targets/foo.txt")
	tmp.assertExists("repository/targets/bar.txt")

	// adding all targets moves them all from staged -> repository
	count := 10
	files := make([]string, count)
	for i := 0; i < count; i++ {
		files[i] = fmt.Sprintf("/file%d.txt", i)
		tmp.writeStagedTarget(files[i], "data")
	}
	c.Assert(r.AddTargets(nil, nil), IsNil)
	c.Assert(r.Snapshot(CompressionTypeNone), IsNil)
	c.Assert(r.Timestamp(), IsNil)
	c.Assert(r.Commit(), IsNil)
	tmp.assertExists("repository/targets/foo.txt")
	tmp.assertExists("repository/targets/bar.txt")
	assertRepoTargets(files...)
	for _, file := range files {
		tmp.assertExists("repository/targets/" + file)
	}
	tmp.assertEmpty("staged/targets")
	tmp.assertEmpty("staged")

	// removing all targets removes them from the repository and targets.json
	c.Assert(r.RemoveTargets(nil), IsNil)
	c.Assert(r.Snapshot(CompressionTypeNone), IsNil)
	c.Assert(r.Timestamp(), IsNil)
	c.Assert(r.Commit(), IsNil)
	tmp.assertEmpty("repository/targets")
	t, err := r.targets()
	c.Assert(err, IsNil)
	c.Assert(t.Targets, HasLen, 0)
}
