package tuf

import (
	"bytes"
	"encoding/json"
	"fmt"
	//"io"
	"path"
	"strings"
	"time"

	cjson "github.com/tent/canonical-json-go"

	"github.com/endophage/go-tuf/data"
	"github.com/endophage/go-tuf/errors"
	"github.com/endophage/go-tuf/keys"
	"github.com/endophage/go-tuf/signed"
	"github.com/endophage/go-tuf/store"
	"github.com/endophage/go-tuf/util"
)

type CompressionType uint8

const (
	CompressionTypeNone CompressionType = iota
	CompressionTypeGzip
)

// topLevelManifests determines the order signatures are verified when committing.
var topLevelManifests = []string{
	"root.json",
	"targets.json",
	"snapshot.json",
	"timestamp.json",
}

// snapshotManifests is the list of default filenames that should be included in the
// snapshots.json. If using delegated targets, additional, dynamic files should also
// be included in snapshots.
var snapshotManifests = []string{
	"root.json",
	"targets.json",
}

// Repo represents an instance of a TUF repo
type Repo struct {
	trust          *signed.Signer
	local          store.LocalStore
	hashAlgorithms []string
	meta           map[string]json.RawMessage
}

// NewRepo is a factory function for instantiating new TUF repos objects.
// If the local store is already populated, local.GetMeta() will initialise
// the Repo with the appropriate state.
func NewRepo(trust signed.TrustService, local store.LocalStore, hashAlgorithms ...string) (*Repo, error) {
	r := &Repo{trust: signed.NewSigner(trust), local: local, hashAlgorithms: hashAlgorithms}

	var err error
	r.meta, err = local.GetMeta()
	if err != nil {
		return nil, err
	}
	return r, nil
}

// Init attempts to initialize a brand new TUF repo. It will fail if
// an existing targets file is detected.
func (r *Repo) Init(consistentSnapshot bool) error {
	t, err := r.targets()
	if err != nil {
		return err
	}
	if len(t.Targets) > 0 {
		return errors.ErrInitNotAllowed
	}
	root := data.NewRoot()
	root.ConsistentSnapshot = consistentSnapshot
	return r.setMeta("root.json", root)
}

func (r *Repo) db() (*keys.DB, error) {
	db := keys.NewDB()
	root, err := r.root()
	if err != nil {
		return nil, err
	}
	for _, k := range root.Keys {
		if err := db.AddKey(&keys.PublicKey{*k, k.ID()}); err != nil {
			return nil, err
		}
	}
	for name, role := range root.Roles {
		if err := db.AddRole(name, role); err != nil {
			return nil, err
		}
	}
	return db, nil
}

func (r *Repo) root() (*data.Root, error) {
	rootJSON, ok := r.meta["root.json"]
	if !ok {
		return data.NewRoot(), nil
	}
	s := &data.Signed{}
	if err := json.Unmarshal(rootJSON, s); err != nil {
		return nil, err
	}
	root := data.NewRoot()
	if err := json.Unmarshal(s.Signed, root); err != nil {
		return nil, err
	}
	return root, nil
}

func (r *Repo) snapshot() (*data.Snapshot, error) {
	snapshotJSON, ok := r.meta["snapshot.json"]
	if !ok {
		return data.NewSnapshot(), nil
	}
	s := &data.Signed{}
	if err := json.Unmarshal(snapshotJSON, s); err != nil {
		return nil, err
	}
	snapshot := data.NewSnapshot()
	if err := json.Unmarshal(s.Signed, snapshot); err != nil {
		return nil, err
	}
	return snapshot, nil
}

func (r *Repo) targets() (*data.Targets, error) {
	targetsJSON, ok := r.meta["targets.json"]
	if !ok {
		return data.NewTargets(), nil
	}
	s := &data.Signed{}
	if err := json.Unmarshal(targetsJSON, s); err != nil {
		return nil, err
	}
	targets := data.NewTargets()
	if err := json.Unmarshal(s.Signed, targets); err != nil {
		return nil, err
	}
	return targets, nil
}

func (r *Repo) timestamp() (*data.Timestamp, error) {
	timestampJSON, ok := r.meta["timestamp.json"]
	if !ok {
		return data.NewTimestamp(), nil
	}
	s := &data.Signed{}
	if err := json.Unmarshal(timestampJSON, s); err != nil {
		return nil, err
	}
	timestamp := data.NewTimestamp()
	if err := json.Unmarshal(s.Signed, timestamp); err != nil {
		return nil, err
	}
	return timestamp, nil
}

func (r *Repo) GenKey(role string) (string, error) {
	return r.GenKeyWithExpires(role, data.DefaultExpires("root"))
}

func (r *Repo) GenKeyWithExpires(keyRole string, expires time.Time) (string, error) {
	if !keys.ValidRole(keyRole) {
		return "", errors.ErrInvalidRole{keyRole}
	}

	if !validExpires(expires) {
		return "", errors.ErrInvalidExpires{expires}
	}

	root, err := r.root()
	if err != nil {
		return "", err
	}

	key, err := r.trust.Create()
	if err != nil {
		return "", err
	}
	if err := r.local.SaveKey(keyRole, &key.Key); err != nil {
		return "", err
	}

	role, ok := root.Roles[keyRole]
	if !ok {
		role = &data.Role{KeyIDs: []string{}, Threshold: 1}
		root.Roles[keyRole] = role
	}
	role.KeyIDs = append(role.KeyIDs, key.ID)

	root.Keys[key.ID] = &key.Key
	root.Expires = expires.Round(time.Second)
	root.Version++

	return key.ID, r.setMeta("root.json", root)
}

func validExpires(expires time.Time) bool {
	return expires.Sub(time.Now()) > 0
}

func (r *Repo) RootKeys() ([]*data.Key, error) {
	root, err := r.root()
	if err != nil {
		return nil, err
	}
	role, ok := root.Roles["root"]
	if !ok {
		return nil, nil
	}
	rootKeys := make([]*data.Key, len(role.KeyIDs))
	for i, id := range role.KeyIDs {
		key, ok := root.Keys[id]
		if !ok {
			return nil, fmt.Errorf("tuf: invalid root metadata")
		}
		rootKeys[i] = key
	}
	return rootKeys, nil
}

func (r *Repo) RevokeKey(role, id string) error {
	return r.RevokeKeyWithExpires(role, id, data.DefaultExpires("root"))
}

func (r *Repo) RevokeKeyWithExpires(keyRole, id string, expires time.Time) error {
	if !keys.ValidRole(keyRole) {
		return errors.ErrInvalidRole{keyRole}
	}

	if !validExpires(expires) {
		return errors.ErrInvalidExpires{expires}
	}

	root, err := r.root()
	if err != nil {
		return err
	}

	if _, ok := root.Keys[id]; !ok {
		return errors.ErrKeyNotFound{keyRole, id}
	}

	role, ok := root.Roles[keyRole]
	if !ok {
		return errors.ErrKeyNotFound{keyRole, id}
	}

	keyIDs := make([]string, 0, len(role.KeyIDs))
	for _, keyID := range role.KeyIDs {
		if keyID == id {
			continue
		}
		keyIDs = append(keyIDs, keyID)
	}
	if len(keyIDs) == len(role.KeyIDs) {
		return errors.ErrKeyNotFound{keyRole, id}
	}
	role.KeyIDs = keyIDs

	delete(root.Keys, id)
	root.Roles[keyRole] = role
	root.Expires = expires.Round(time.Second)
	root.Version++

	return r.setMeta("root.json", root)
}

func (r *Repo) setMeta(name string, meta interface{}) error {
	keys, err := r.getKeys(strings.TrimSuffix(name, ".json"))
	if err != nil {
		return err
	}
	b, err := cjson.Marshal(meta)
	if err != nil {
		return err
	}
	s := &data.Signed{Signed: b}
	err = r.trust.Sign(s, keys...)
	if err != nil {
		return err
	}
	b, err = json.Marshal(s)
	if err != nil {
		return err
	}
	r.meta[name] = b
	return r.local.SetMeta(name, b)
}

func (r *Repo) Sign(name string) error {
	role := strings.TrimSuffix(name, ".json")
	if !keys.ValidRole(role) {
		return errors.ErrInvalidRole{role}
	}

	s, err := r.signedMeta(name)
	if err != nil {
		return err
	}

	keys, err := r.getKeys(role)
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return errors.ErrInsufficientKeys{name}
	}

	r.trust.Sign(s, keys...)

	b, err := json.Marshal(s)
	if err != nil {
		return err
	}
	r.meta[name] = b
	return r.local.SetMeta(name, b)
}

// getKeys returns signing keys from local storage.
//
// Only keys contained in the keys db are returned (i.e. local keys which have
// been revoked are omitted), except for the root role in which case all local
// keys are returned (revoked root keys still need to sign new root metadata so
// clients can verify the new root.json and update their keys db accordingly).
func (r *Repo) getKeys(name string) ([]*keys.PublicKey, error) {
	localKeys, err := r.local.GetKeys(name)
	if err != nil {
		return nil, err
	}
	if name == "root" {
		rootkeys := make([]*keys.PublicKey, 0, len(localKeys))
		for _, key := range localKeys {
			rootkeys = append(rootkeys, &keys.PublicKey{*key, key.ID()})
		}
		return rootkeys, nil
	}
	db, err := r.db()
	if err != nil {
		return nil, err
	}
	role := db.GetRole(name)
	if role == nil {
		return nil, nil
	}
	if len(role.KeyIDs) == 0 {
		return nil, nil
	}
	rolekeys := make([]*keys.PublicKey, 0, len(role.KeyIDs))
	for _, key := range localKeys {
		if role.ValidKey(key.ID()) {
			rolekeys = append(rolekeys, &keys.PublicKey{*key, key.ID()})
		}
	}
	return rolekeys, nil
}

func (r *Repo) signedMeta(name string) (*data.Signed, error) {
	b, ok := r.meta[name]
	if !ok {
		return nil, errors.ErrMissingMetadata{name}
	}
	s := &data.Signed{}
	if err := json.Unmarshal(b, s); err != nil {
		return nil, err
	}
	return s, nil
}

func validManifest(name string) bool {
	for _, m := range topLevelManifests {
		if m == name {
			return true
		}
	}
	return false
}

func (r *Repo) AddTarget(path string, custom json.RawMessage) error {
	return r.AddTargets([]string{path}, custom)
}

func (r *Repo) AddTargets(paths []string, custom json.RawMessage) error {
	return r.AddTargetsWithExpires(paths, custom, data.DefaultExpires("targets"))
}

func (r *Repo) AddTargetWithExpires(path string, custom json.RawMessage, expires time.Time) error {
	return r.AddTargetsWithExpires([]string{path}, custom, expires)
}

func (r *Repo) AddTargetsWithExpires(paths []string, custom json.RawMessage, expires time.Time) error {
	if !validExpires(expires) {
		return errors.ErrInvalidExpires{expires}
	}

	t, err := r.targets()
	if err != nil {
		return err
	}
	normalizedPaths := make([]string, len(paths))
	for i, path := range paths {
		normalizedPaths[i] = util.NormalizeTarget(path)
	}
	if err := r.local.WalkStagedTargets(normalizedPaths, func(path string, meta data.FileMeta) (err error) {
		t.Targets[util.NormalizeTarget(path)] = meta
		return nil
	}); err != nil {
		return err
	}
	t.Expires = expires.Round(time.Second)
	t.Version++
	return r.setMeta("targets.json", t)
}

func (r *Repo) RemoveTarget(path string) error {
	return r.RemoveTargets([]string{path})
}

func (r *Repo) RemoveTargets(paths []string) error {
	return r.RemoveTargetsWithExpires(paths, data.DefaultExpires("targets"))
}

func (r *Repo) RemoveTargetWithExpires(path string, expires time.Time) error {
	return r.RemoveTargetsWithExpires([]string{path}, expires)
}

// If paths is empty, all targets will be removed.
func (r *Repo) RemoveTargetsWithExpires(paths []string, expires time.Time) error {
	if !validExpires(expires) {
		return errors.ErrInvalidExpires{expires}
	}

	t, err := r.targets()
	if err != nil {
		return err
	}
	if len(paths) == 0 {
		t.Targets = make(data.Files)
	} else {
		removed := false
		for _, path := range paths {
			path = util.NormalizeTarget(path)
			if _, ok := t.Targets[path]; !ok {
				continue
			}
			removed = true
			delete(t.Targets, path)
		}
		if !removed {
			return nil
		}
	}
	t.Expires = expires.Round(time.Second)
	t.Version++
	return r.setMeta("targets.json", t)
}

func (r *Repo) Snapshot(t CompressionType) error {
	return r.SnapshotWithExpires(t, data.DefaultExpires("snapshot"))
}

func (r *Repo) SnapshotWithExpires(t CompressionType, expires time.Time) error {
	if !validExpires(expires) {
		return errors.ErrInvalidExpires{expires}
	}

	snapshot, err := r.snapshot()
	if err != nil {
		return err
	}
	db, err := r.db()
	if err != nil {
		return err
	}
	// TODO: generate compressed manifests
	for _, name := range snapshotManifests {
		if err := r.verifySignature(name, db); err != nil {
			return err
		}
		var err error
		snapshot.Meta[name], err = r.fileMeta(name)
		if err != nil {
			return err
		}
	}
	snapshot.Expires = expires.Round(time.Second)
	snapshot.Version++
	return r.setMeta("snapshot.json", snapshot)
}

func (r *Repo) Timestamp() error {
	return r.TimestampWithExpires(data.DefaultExpires("timestamp"))
}

func (r *Repo) TimestampWithExpires(expires time.Time) error {
	if !validExpires(expires) {
		return errors.ErrInvalidExpires{expires}
	}

	db, err := r.db()
	if err != nil {
		return err
	}
	if err := r.verifySignature("snapshot.json", db); err != nil {
		return err
	}
	timestamp, err := r.timestamp()
	if err != nil {
		return err
	}
	timestamp.Meta["snapshot.json"], err = r.fileMeta("snapshot.json")
	if err != nil {
		return err
	}
	timestamp.Expires = expires.Round(time.Second)
	timestamp.Version++
	return r.setMeta("timestamp.json", timestamp)
}

func (r *Repo) fileHashes() (map[string]data.Hashes, error) {
	hashes := make(map[string]data.Hashes)
	addHashes := func(name string, meta data.Files) {
		if m, ok := meta[name]; ok {
			hashes[name] = m.Hashes
		}
	}
	timestamp, err := r.timestamp()
	if err != nil {
		return nil, err
	}
	snapshot, err := r.snapshot()
	if err != nil {
		return nil, err
	}
	addHashes("root.json", snapshot.Meta)
	addHashes("targets.json", snapshot.Meta)
	addHashes("snapshot.json", timestamp.Meta)
	t, err := r.targets()
	if err != nil {
		return nil, err
	}
	for name, meta := range t.Targets {
		hashes[path.Join("targets", name)] = meta.Hashes
	}
	return hashes, nil
}

func (r *Repo) Commit() error {
	// check we have all the metadata
	for _, name := range topLevelManifests {
		if _, ok := r.meta[name]; !ok {
			return errors.ErrMissingMetadata{name}
		}
	}

	// check roles are valid
	root, err := r.root()
	if err != nil {
		return err
	}
	for name, role := range root.Roles {
		if len(role.KeyIDs) < role.Threshold {
			return errors.ErrNotEnoughKeys{name, len(role.KeyIDs), role.Threshold}
		}
	}

	// verify hashes in snapshot.json are up to date
	snapshot, err := r.snapshot()
	if err != nil {
		return err
	}
	for _, name := range snapshotManifests {
		expected, ok := snapshot.Meta[name]
		if !ok {
			return fmt.Errorf("tuf: snapshot.json missing hash for %s", name)
		}
		actual, err := r.fileMeta(name)
		if err != nil {
			return err
		}
		if err := util.FileMetaEqual(actual, expected); err != nil {
			return fmt.Errorf("tuf: invalid %s in snapshot.json: %s", name, err)
		}
	}

	// verify hashes in timestamp.json are up to date
	timestamp, err := r.timestamp()
	if err != nil {
		return err
	}
	snapshotMeta, err := r.fileMeta("snapshot.json")
	if err != nil {
		return err
	}
	if err := util.FileMetaEqual(snapshotMeta, timestamp.Meta["snapshot.json"]); err != nil {
		return fmt.Errorf("tuf: invalid snapshot.json in timestamp.json: %s", err)
	}

	// verify all signatures are correct
	db, err := r.db()
	if err != nil {
		return err
	}
	for _, name := range topLevelManifests {
		if err := r.verifySignature(name, db); err != nil {
			return err
		}
	}

	hashes, err := r.fileHashes()
	if err != nil {
		return err
	}
	return r.local.Commit(r.meta, root.ConsistentSnapshot, hashes)
}

func (r *Repo) Clean() error {
	return r.local.Clean()
}

func (r *Repo) verifySignature(name string, db *keys.DB) error {
	s, err := r.signedMeta(name)
	if err != nil {
		return err
	}
	role := strings.TrimSuffix(name, ".json")
	if err := signed.Verify(s, role, 0, db); err != nil {
		return errors.ErrInsufficientSignatures{name, err}
	}
	return nil
}

func (r *Repo) fileMeta(name string) (data.FileMeta, error) {
	b, ok := r.meta[name]
	if !ok {
		return data.FileMeta{}, errors.ErrMissingMetadata{name}
	}
	return util.GenerateFileMeta(bytes.NewReader(b), r.hashAlgorithms...)
}
