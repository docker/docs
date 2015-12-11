package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/docker/notary/client/changelist"
	tuf "github.com/docker/notary/tuf"
	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/keys"
	"github.com/docker/notary/tuf/store"
)

// Use this to initialize remote HTTPStores from the config settings
func getRemoteStore(baseURL, gun string, rt http.RoundTripper) (store.RemoteStore, error) {
	return store.NewHTTPStore(
		baseURL+"/v2/"+gun+"/_trust/tuf/",
		"",
		"json",
		"",
		"key",
		rt,
	)
}

func applyChangelist(repo *tuf.Repo, cl changelist.Changelist) error {
	it, err := cl.NewIterator()
	if err != nil {
		return err
	}
	index := 0
	for it.HasNext() {
		c, err := it.Next()
		if err != nil {
			return err
		}
		switch c.Scope() {
		case changelist.ScopeTargets:
			err = applyTargetsChange(repo, c)
		case changelist.ScopeRoot:
			err = applyRootChange(repo, c)
		default:
			logrus.Debug("scope not supported: ", c.Scope())
		}
		index++
		if err != nil {
			return err
		}
	}
	logrus.Debugf("applied %d change(s)", index)
	return nil
}

func applyTargetsChange(repo *tuf.Repo, c changelist.Change) error {
	var err error
	switch c.Action() {
	case changelist.ActionCreate:
		logrus.Debug("changelist add: ", c.Path())
		meta := &data.FileMeta{}
		err = json.Unmarshal(c.Content(), meta)
		if err != nil {
			return err
		}
		files := data.Files{c.Path(): *meta}
		_, err = repo.AddTargets(c.Scope(), files)
	case changelist.ActionDelete:
		logrus.Debug("changelist remove: ", c.Path())
		err = repo.RemoveTargets(c.Scope(), c.Path())
	default:
		logrus.Debug("action not yet supported: ", c.Action())
	}
	if err != nil {
		return err
	}
	return nil
}

func applyRootChange(repo *tuf.Repo, c changelist.Change) error {
	var err error
	switch c.Type() {
	case changelist.TypeRootRole:
		err = applyRootRoleChange(repo, c)
	default:
		logrus.Debug("type of root change not yet supported: ", c.Type())
	}
	return err // might be nil
}

func applyRootRoleChange(repo *tuf.Repo, c changelist.Change) error {
	switch c.Action() {
	case changelist.ActionCreate:
		// replaces all keys for a role
		d := &changelist.TufRootData{}
		err := json.Unmarshal(c.Content(), d)
		if err != nil {
			return err
		}
		err = repo.ReplaceBaseKeys(d.RoleName, d.Keys...)
		if err != nil {
			return err
		}
	default:
		logrus.Debug("action not yet supported for root: ", c.Action())
	}
	return nil
}

func nearExpiry(r *data.SignedRoot) bool {
	plus6mo := time.Now().AddDate(0, 6, 0)
	return r.Signed.Expires.Before(plus6mo)
}

// Fetches a public key from a remote store, given a gun and role
func getRemoteKey(url, gun, role string, rt http.RoundTripper) (data.PublicKey, error) {
	remote, err := getRemoteStore(url, gun, rt)
	if err != nil {
		return nil, err
	}
	rawPubKey, err := remote.GetKey(role)
	if err != nil {
		return nil, err
	}

	pubKey, err := data.UnmarshalPublicKey(rawPubKey)
	if err != nil {
		return nil, err
	}

	return pubKey, nil
}

// add a key to a KeyDB, and create a role for the key and add it.
func addKeyForRole(kdb *keys.KeyDB, role string, key data.PublicKey) error {
	theRole, err := data.NewRole(role, 1, []string{key.ID()}, nil, nil)
	if err != nil {
		return err
	}
	kdb.AddKey(key)
	if err := kdb.AddRole(theRole); err != nil {
		return err
	}
	return nil
}

// signs and serializes the metadata for a canonical role in a tuf repo to JSON
func serializeCanonicalRole(tufRepo *tuf.Repo, role string) (out []byte, err error) {
	var s *data.Signed
	switch role {
	case data.CanonicalRootRole:
		s, err = tufRepo.SignRoot(data.DefaultExpires(role))
	case data.CanonicalSnapshotRole:
		s, err = tufRepo.SignSnapshot(data.DefaultExpires(role))
	case data.CanonicalTargetsRole:
		s, err = tufRepo.SignTargets(role, data.DefaultExpires(role))
	default:
		err = fmt.Errorf("%s not supported role to sign on the client", role)
	}

	if err != nil {
		return
	}

	return json.Marshal(s)
}
