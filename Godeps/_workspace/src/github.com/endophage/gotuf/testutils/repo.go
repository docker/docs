package testutils

import (
	"encoding/json"
	"math/rand"
	"time"

	"github.com/endophage/gotuf/data"
	"github.com/endophage/gotuf/utils"
	fuzz "github.com/google/gofuzz"

	tuf "github.com/endophage/gotuf"
	"github.com/endophage/gotuf/keys"
	"github.com/endophage/gotuf/signed"
)

// EmptyRepo creates an in memory key database, crypto service
// and initializes a repo with no targets or delegations.
func EmptyRepo() (*keys.KeyDB, *tuf.Repo, signed.CryptoService) {
	c := signed.NewEd25519()
	kdb := keys.NewDB()
	r := tuf.NewRepo(kdb, c)

	for _, role := range []string{"root", "targets", "snapshot", "timestamp"} {
		key, _ := c.Create(role, data.ED25519Key)
		role, _ := data.NewRole(role, 1, []string{key.ID()}, nil, nil)
		kdb.AddKey(key)
		kdb.AddRole(role)
	}

	r.InitRepo(false)
	return kdb, r, c
}

// AddTarget generates a fake target and adds it to a repo.
func AddTarget(role string, r *tuf.Repo) (name string, meta data.FileMeta, content []byte, err error) {
	randness := fuzz.Continue{}
	content = RandomByteSlice(1024)
	name = randness.RandString()
	t := data.FileMeta{
		Length: int64(len(content)),
		Hashes: data.Hashes{
			"sha256": utils.DoHash("sha256", content),
			"sha512": utils.DoHash("sha512", content),
		},
	}
	files := data.Files{name: t}
	_, err = r.AddTargets(role, files)
	return
}

// RandomByteSlice generates some random data to be used for testing only
func RandomByteSlice(maxSize int) []byte {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	contentSize := r.Intn(maxSize)
	content := make([]byte, contentSize)
	for i := range content {
		content[i] = byte(r.Int63() & 0xff)
	}
	return content
}

// Sign signs all top level roles in a repo in the appropriate order
func Sign(repo *tuf.Repo) (root, targets, snapshot, timestamp *data.Signed, err error) {
	root, err = repo.SignRoot(data.DefaultExpires("root"), nil)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	targets, err = repo.SignTargets("targets", data.DefaultExpires("targets"), nil)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	snapshot, err = repo.SignSnapshot(data.DefaultExpires("snapshot"), nil)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	timestamp, err = repo.SignTimestamp(data.DefaultExpires("timestamp"), nil)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	return
}

// Serialize takes the Signed objects for the 4 top level roles and serializes them all to JSON
func Serialize(sRoot, sTargets, sSnapshot, sTimestamp *data.Signed) (root, targets, snapshot, timestamp []byte, err error) {
	root, err = json.Marshal(sRoot)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	targets, err = json.Marshal(sTargets)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	snapshot, err = json.Marshal(sSnapshot)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	timestamp, err = json.Marshal(sTimestamp)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	return
}
