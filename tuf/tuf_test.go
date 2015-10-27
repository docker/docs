package tuf

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/keys"
	"github.com/docker/notary/tuf/signed"
)

func initRepo(t *testing.T, cryptoService signed.CryptoService, keyDB *keys.KeyDB) *Repo {
	rootKey, err := cryptoService.Create("root", data.ED25519Key)
	if err != nil {
		t.Fatal(err)
	}
	targetsKey, err := cryptoService.Create("targets", data.ED25519Key)
	if err != nil {
		t.Fatal(err)
	}
	snapshotKey, err := cryptoService.Create("snapshot", data.ED25519Key)
	if err != nil {
		t.Fatal(err)
	}
	timestampKey, err := cryptoService.Create("timestamp", data.ED25519Key)
	if err != nil {
		t.Fatal(err)
	}

	keyDB.AddKey(rootKey)
	keyDB.AddKey(targetsKey)
	keyDB.AddKey(snapshotKey)
	keyDB.AddKey(timestampKey)

	rootRole := &data.Role{
		Name: "root",
		RootRole: data.RootRole{
			KeyIDs:    []string{rootKey.ID()},
			Threshold: 1,
		},
	}
	targetsRole := &data.Role{
		Name: "targets",
		RootRole: data.RootRole{
			KeyIDs:    []string{targetsKey.ID()},
			Threshold: 1,
		},
	}
	snapshotRole := &data.Role{
		Name: "snapshot",
		RootRole: data.RootRole{
			KeyIDs:    []string{snapshotKey.ID()},
			Threshold: 1,
		},
	}
	timestampRole := &data.Role{
		Name: "timestamp",
		RootRole: data.RootRole{
			KeyIDs:    []string{timestampKey.ID()},
			Threshold: 1,
		},
	}

	keyDB.AddRole(rootRole)
	keyDB.AddRole(targetsRole)
	keyDB.AddRole(snapshotRole)
	keyDB.AddRole(timestampRole)

	repo := NewRepo(keyDB, cryptoService)
	err = repo.InitRepo(false)
	if err != nil {
		t.Fatal(err)
	}
	return repo
}

func writeRepo(t *testing.T, dir string, repo *Repo) {
	//err := os.Remove(dir)
	//if err != nil {
	//	t.Fatal(err)
	//}
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		t.Fatal(err)
	}
	signedRoot, err := repo.SignRoot(data.DefaultExpires("root"), nil)
	if err != nil {
		t.Fatal(err)
	}
	rootJSON, _ := json.Marshal(signedRoot)
	ioutil.WriteFile(dir+"/root.json", rootJSON, 0755)

	for r := range repo.Targets {
		signedTargets, err := repo.SignTargets(r, data.DefaultExpires("targets"), nil)
		if err != nil {
			t.Fatal(err)
		}
		targetsJSON, _ := json.Marshal(signedTargets)
		p := path.Join(dir, r+".json")
		parentDir := filepath.Dir(p)
		os.MkdirAll(parentDir, 0755)
		ioutil.WriteFile(p, targetsJSON, 0755)
	}

	signedSnapshot, err := repo.SignSnapshot(data.DefaultExpires("snapshot"), nil)
	if err != nil {
		t.Fatal(err)
	}
	snapshotJSON, _ := json.Marshal(signedSnapshot)
	ioutil.WriteFile(dir+"/snapshot.json", snapshotJSON, 0755)

	signedTimestamp, err := repo.SignTimestamp(data.DefaultExpires("timestamp"), nil)
	if err != nil {
		t.Fatal(err)
	}
	timestampJSON, _ := json.Marshal(signedTimestamp)
	ioutil.WriteFile(dir+"/timestamp.json", timestampJSON, 0755)
}

func TestInitRepo(t *testing.T) {
	ed25519 := signed.NewEd25519()
	keyDB := keys.NewDB()
	repo := initRepo(t, ed25519, keyDB)
	writeRepo(t, "/tmp/tufrepo", repo)
}

func TestUpdateDelegations(t *testing.T) {
	ed25519 := signed.NewEd25519()
	keyDB := keys.NewDB()
	repo := initRepo(t, ed25519, keyDB)

	testKey, err := ed25519.Create("targets/test", data.ED25519Key)
	if err != nil {
		t.Fatal(err)
	}
	role, err := data.NewRole("targets/test", 1, []string{testKey.ID()}, []string{"test"}, []string{})
	if err != nil {
		t.Fatal(err)
	}

	err = repo.UpdateDelegations(role, []data.Key{testKey}, "")
	if err != nil {
		t.Fatal(err)
	}

	testDeepKey, err := ed25519.Create("targets/test/deep", data.ED25519Key)
	if err != nil {
		t.Fatal(err)
	}
	roleDeep, err := data.NewRole("targets/test/deep", 1, []string{testDeepKey.ID()}, []string{"test/deep"}, []string{})
	if err != nil {
		t.Fatal(err)
	}

	err = repo.UpdateDelegations(roleDeep, []data.Key{testDeepKey}, "")
	if err != nil {
		t.Fatal(err)
	}

	writeRepo(t, "/tmp/tufdelegation", repo)
}
