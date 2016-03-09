package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"sort"
	"strings"

	"github.com/Sirupsen/logrus"

	"github.com/docker/notary/server/snapshot"
	"github.com/docker/notary/server/storage"
	"github.com/docker/notary/server/timestamp"
	"github.com/docker/notary/tuf"
	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/signed"
	"github.com/docker/notary/tuf/utils"
	"github.com/docker/notary/tuf/validation"
)

// validateUpload checks that the updates being pushed
// are semantically correct and the signatures are correct
// A list of possibly modified updates are returned if all
// validation was successful. This allows the snapshot to be
// created and added if snapshotting has been delegated to the
// server
func validateUpdate(cs signed.CryptoService, gun string, updates []storage.MetaUpdate, store storage.MetaStore) ([]storage.MetaUpdate, error) {
	repo := tuf.NewRepo(cs)
	rootRole := data.CanonicalRootRole
	snapshotRole := data.CanonicalSnapshotRole

	// some delegated targets role may be invalid based on other updates
	// that have been made by other clients. We'll rebuild the slice of
	// updates with only the things we should actually update
	updatesToApply := make([]storage.MetaUpdate, 0, len(updates))

	roles := make(map[string]storage.MetaUpdate)
	for _, v := range updates {
		roles[v.Role] = v
	}

	var root *data.SignedRoot
	_, oldRootJSON, err := store.GetCurrent(gun, rootRole)
	if _, ok := err.(storage.ErrNotFound); err != nil && !ok {
		// problem with storage. No expectation we can
		// write if we can't read so bail.
		logrus.Error("error reading previous root: ", err.Error())
		return nil, err
	}
	if rootUpdate, ok := roles[rootRole]; ok {
		// if root is present, validate its integrity, possibly
		// against a previous root
		if root, err = validateRoot(gun, oldRootJSON, rootUpdate.Data); err != nil {
			logrus.Error("ErrBadRoot: ", err.Error())
			return nil, validation.ErrBadRoot{Msg: err.Error()}
		}

		// setting root will update keys db
		if err = repo.SetRoot(root); err != nil {
			logrus.Error("ErrValidation: ", err.Error())
			return nil, validation.ErrValidation{Msg: err.Error()}
		}
		logrus.Debug("Successfully validated root")
		updatesToApply = append(updatesToApply, rootUpdate)
	} else {
		if oldRootJSON == nil {
			return nil, validation.ErrValidation{Msg: "no pre-existing root and no root provided in update."}
		}
		parsedOldRoot := &data.SignedRoot{}
		if err := json.Unmarshal(oldRootJSON, parsedOldRoot); err != nil {
			return nil, validation.ErrValidation{Msg: "pre-existing root is corrupted and no root provided in update."}
		}
		if err = repo.SetRoot(parsedOldRoot); err != nil {
			logrus.Error("ErrValidation: ", err.Error())
			return nil, validation.ErrValidation{Msg: err.Error()}
		}
	}

	targetsToUpdate, err := loadAndValidateTargets(gun, repo, roles, store)
	if err != nil {
		return nil, err
	}
	updatesToApply = append(updatesToApply, targetsToUpdate...)

	// there's no need to load files from the database if no targets etc...
	// were uploaded because that means they haven't been updated and
	// the snapshot will already contain the correct hashes and sizes for
	// those targets (incl. delegated targets)
	logrus.Debug("Successfully validated targets")

	// At this point, root and targets must have been loaded into the repo
	if _, ok := roles[snapshotRole]; ok {
		var oldSnap *data.SignedSnapshot
		_, oldSnapJSON, err := store.GetCurrent(gun, snapshotRole)
		if _, ok := err.(storage.ErrNotFound); err != nil && !ok {
			// problem with storage. No expectation we can
			// write if we can't read so bail.
			logrus.Error("error reading previous snapshot: ", err.Error())
			return nil, err
		} else if err == nil {
			oldSnap = &data.SignedSnapshot{}
			if err := json.Unmarshal(oldSnapJSON, oldSnap); err != nil {
				oldSnap = nil
			}
		}

		if err := loadAndValidateSnapshot(snapshotRole, oldSnap, roles[snapshotRole], roles, repo); err != nil {
			logrus.Error("ErrBadSnapshot: ", err.Error())
			return nil, validation.ErrBadSnapshot{Msg: err.Error()}
		}
		logrus.Debug("Successfully validated snapshot")
		updatesToApply = append(updatesToApply, roles[snapshotRole])
	} else {
		// Check:
		//   - we have a snapshot key
		//   - it matches a snapshot key signed into the root.json
		// Then:
		//   - generate a new snapshot
		//   - add it to the updates
		update, err := generateSnapshot(gun, repo, store)
		if err != nil {
			return nil, err
		}
		updatesToApply = append(updatesToApply, *update)
	}

	// generate a timestamp immediately
	update, err := generateTimestamp(gun, repo, store)
	if err != nil {
		return nil, err
	}

	return append(updatesToApply, *update), nil
}

func loadAndValidateTargets(gun string, repo *tuf.Repo, roles map[string]storage.MetaUpdate, store storage.MetaStore) ([]storage.MetaUpdate, error) {
	targetsRoles := make(utils.RoleList, 0)
	for role := range roles {
		if role == data.CanonicalTargetsRole || data.IsDelegation(role) {
			targetsRoles = append(targetsRoles, role)
		}
	}

	// N.B. RoleList sorts paths with fewer segments first.
	// By sorting, we'll always process shallower targets updates before deeper
	// ones (i.e. we'll load and validate targets before targets/foo). This
	// helps ensure we only load from storage when necessary in a cleaner way.
	sort.Sort(targetsRoles)

	updatesToApply := make([]storage.MetaUpdate, 0, len(targetsRoles))
	for _, role := range targetsRoles {
		// don't load parent if current role is "targets",
		// we must load all ancestor roles for delegations to validate the full parent chain
		ancestorRole := role
		for ancestorRole != data.CanonicalTargetsRole {
			ancestorRole = path.Dir(ancestorRole)
			if _, ok := repo.Targets[ancestorRole]; !ok {
				err := loadTargetsFromStore(gun, ancestorRole, repo, store)
				if err != nil {
					return nil, err
				}
			}
		}
		var (
			t   *data.SignedTargets
			err error
		)
		if t, err = validateTargets(role, roles, repo); err != nil {
			if _, ok := err.(data.ErrInvalidRole); ok {
				// role wasn't found in its parent. It has been removed
				// or never existed. Drop this role from the update
				// (by not adding it to updatesToApply)
				continue
			}
			logrus.Error("ErrBadTargets: ", err.Error())
			return nil, validation.ErrBadTargets{Msg: err.Error()}
		}
		// this will load keys and roles into the kdb
		err = repo.SetTargets(role, t)
		if err != nil {
			return nil, err
		}
		updatesToApply = append(updatesToApply, roles[role])
	}
	return updatesToApply, nil
}

func loadTargetsFromStore(gun, role string, repo *tuf.Repo, store storage.MetaStore) error {
	_, tgtJSON, err := store.GetCurrent(gun, role)
	if err != nil {
		return err
	}
	t := &data.SignedTargets{}
	err = json.Unmarshal(tgtJSON, t)
	if err != nil {
		return err
	}
	return repo.SetTargets(role, t)
}

// generateSnapshot generates a new snapshot from the previous one in the store - this assumes all
// the other roles except timestamp have already been set on the repo, and will set the generated
// snapshot on the repo as well
func generateSnapshot(gun string, repo *tuf.Repo, store storage.MetaStore) (*storage.MetaUpdate, error) {
	var prev *data.SignedSnapshot
	_, currentJSON, err := store.GetCurrent(gun, data.CanonicalSnapshotRole)
	if err == nil {
		prev = new(data.SignedSnapshot)
		if err = json.Unmarshal(currentJSON, prev); err != nil {
			logrus.Error("Failed to unmarshal existing snapshot for GUN ", gun)
			return nil, err
		}
	}

	if _, ok := err.(storage.ErrNotFound); !ok && err != nil {
		return nil, err
	}

	metaUpdate, err := snapshot.NewSnapshotUpdate(prev, repo)
	switch err.(type) {
	case signed.ErrInsufficientSignatures, signed.ErrNoKeys:
		// If we cannot sign the snapshot, then we don't have keys for the snapshot,
		// and the client should have submitted a snapshot
		return nil, validation.ErrBadHierarchy{
			Missing: data.CanonicalSnapshotRole,
			Msg:     "no snapshot was included in update and server does not hold current snapshot key for repository"}
	case nil:
		return metaUpdate, nil

	default:
		return nil, validation.ErrValidation{Msg: err.Error()}
	}
}

// generateTimestamp generates a new timestamp from the previous one in the store - this assumes all
// the other roles have already been set on the repo, and will set the generated timestamp on the repo as well
func generateTimestamp(gun string, repo *tuf.Repo, store storage.MetaStore) (*storage.MetaUpdate, error) {
	var prev *data.SignedTimestamp
	_, currentJSON, err := store.GetCurrent(gun, data.CanonicalTimestampRole)
	if err == nil {
		prev = new(data.SignedTimestamp)
		if err = json.Unmarshal(currentJSON, prev); err != nil {
			logrus.Error("Failed to unmarshal existing timestamp for GUN ", gun)
			return nil, err
		}
	}

	if _, ok := err.(storage.ErrNotFound); !ok && err != nil {
		return nil, err
	}

	metaUpdate, err := timestamp.NewTimestampUpdate(prev, repo)
	if err != nil {
		_, noSigs := err.(signed.ErrInsufficientSignatures)
		_, noKeys := err.(signed.ErrNoKeys)
		// If we cannot sign the timestamp, then we don't have keys for the timestamp,
		// and the client screwed up their root
		if noSigs || noKeys {
			return nil, validation.ErrBadRoot{
				Msg: fmt.Sprintf("none of the following timestamp keys exist on the server: %s",
					strings.Join(repo.Root.Signed.Roles[data.CanonicalTimestampRole].KeyIDs, ", ")),
			}
		}

		return nil, validation.ErrValidation{Msg: err.Error()}
	}
	return metaUpdate, nil
}

// loadAndValidateSnapshot validates that the given snapshot update is valid.  It also sets the new snapshot
// on the TUF repo, if it is valid
func loadAndValidateSnapshot(role string, oldSnap *data.SignedSnapshot, snapUpdate storage.MetaUpdate, roles map[string]storage.MetaUpdate, repo *tuf.Repo) error {
	s := &data.Signed{}
	err := json.Unmarshal(snapUpdate.Data, s)
	if err != nil {
		return errors.New("could not parse snapshot")
	}
	// version specifically gets validated when writing to store to
	// better handle race conditions there.
	snapshotRole, err := repo.GetBaseRole(role)
	if err != nil {
		return err
	}
	if err := signed.Verify(s, snapshotRole, 0); err != nil {
		return err
	}

	snap, err := data.SnapshotFromSigned(s)
	if err != nil {
		return errors.New("could not parse snapshot")
	}
	if !data.ValidTUFType(snap.Signed.Type, data.CanonicalSnapshotRole) {
		return errors.New("snapshot has wrong type")
	}
	err = checkSnapshotEntries(role, oldSnap, snap, roles)
	if err != nil {
		return err
	}
	repo.SetSnapshot(snap)
	return nil
}

func checkSnapshotEntries(role string, oldSnap, snap *data.SignedSnapshot, roles map[string]storage.MetaUpdate) error {
	snapshotRole := data.CanonicalSnapshotRole
	timestampRole := data.CanonicalTimestampRole
	for r, update := range roles {
		if r == snapshotRole || r == timestampRole {
			continue
		}
		m, ok := snap.Signed.Meta[r]
		if !ok {
			return fmt.Errorf("snapshot missing metadata for %s", r)
		}
		if int64(len(update.Data)) != m.Length {
			return fmt.Errorf("snapshot has incorrect length for %s", r)
		}

		if !checkHashes(m, update.Data) {
			return fmt.Errorf("snapshot has incorrect hashes for %s", r)
		}
	}
	return nil
}

func checkHashes(meta data.FileMeta, update []byte) bool {
	for alg, digest := range meta.Hashes {
		d := utils.DoHash(alg, update)
		if !bytes.Equal(digest, d) {
			return false
		}
	}
	return true
}

func validateTargets(role string, roles map[string]storage.MetaUpdate, repo *tuf.Repo) (*data.SignedTargets, error) {
	// TODO: when delegations are being validated, validate parent
	//       role exists for any delegation
	s := &data.Signed{}
	err := json.Unmarshal(roles[role].Data, s)
	if err != nil {
		return nil, fmt.Errorf("could not parse %s", role)
	}
	// version specifically gets validated when writing to store to
	// better handle race conditions there.
	var targetOrDelgRole data.BaseRole
	if role == data.CanonicalTargetsRole {
		targetOrDelgRole, err = repo.GetBaseRole(role)
		if err != nil {
			logrus.Debugf("no %s role loaded", role)
			return nil, err
		}
	} else {
		delgRole, err := repo.GetDelegationRole(role)
		if err != nil {
			logrus.Debugf("no %s delegation role loaded", role)
			return nil, err
		}
		targetOrDelgRole = delgRole.BaseRole
	}
	if err := signed.Verify(s, targetOrDelgRole, 0); err != nil {
		return nil, err
	}
	t, err := data.TargetsFromSigned(s, role)
	if err != nil {
		return nil, err
	}
	if !data.ValidTUFType(t.Signed.Type, data.CanonicalTargetsRole) {
		return nil, fmt.Errorf("%s has wrong type", role)
	}
	return t, nil
}

// validateRoot returns the parsed data.SignedRoot object if the new root:
// - is a valid root metadata object
// - has the correct number of timestamp keys
// - validates against the previous root's signatures (if there was a rotation)
// - is valid against itself (signature-wise)
func validateRoot(gun string, oldRoot, newRoot []byte) (
	*data.SignedRoot, error) {

	parsedNewSigned := &data.Signed{}
	err := json.Unmarshal(newRoot, parsedNewSigned)
	if err != nil {
		return nil, err
	}

	// validates the structure of the root metadata
	parsedNewRoot, err := data.RootFromSigned(parsedNewSigned)
	if err != nil {
		return nil, err
	}

	newRootRole, _ := parsedNewRoot.BuildBaseRole(data.CanonicalRootRole)
	if err != nil { // should never happen, since the root metadata has been validated
		return nil, err
	}

	newTimestampRole, err := parsedNewRoot.BuildBaseRole(data.CanonicalTimestampRole)
	if err != nil { // should never happen, since the root metadata has been validated
		return nil, err
	}
	// According to the TUF spec, any role may have more than one signing
	// key and require a threshold signature.  However, notary-server
	// creates the timestamp, and there is only ever one, so a threshold
	// greater than one would just always fail validation
	if newTimestampRole.Threshold != 1 {
		return nil, fmt.Errorf("timestamp role has invalid threshold")
	}

	if oldRoot != nil {
		if err := checkAgainstOldRoot(oldRoot, newRootRole, parsedNewSigned); err != nil {
			return nil, err
		}
	}

	if err := signed.VerifyRoot(parsedNewSigned, newRootRole.Threshold, newRootRole.Keys); err != nil {
		return nil, err
	}

	return parsedNewRoot, nil
}

// checkAgainstOldRoot errors if an invalid root rotation has taken place
func checkAgainstOldRoot(oldRoot []byte, newRootRole data.BaseRole, newSigned *data.Signed) error {
	parsedOldRoot := &data.SignedRoot{}
	err := json.Unmarshal(oldRoot, parsedOldRoot)
	if err != nil {
		logrus.Warn("Old root could not be parsed, and cannot be used to check the new root.")
		return nil
	}

	oldRootRole, err := parsedOldRoot.BuildBaseRole(data.CanonicalRootRole)
	if err != nil {
		logrus.Warn("Old root does not have a valid root role, and cannot be used to check the new root.")
		return nil
	}

	// if the set of keys has changed between the old root and new root, then a root
	// rotation may have taken place
	rotation := len(oldRootRole.Keys) != len(newRootRole.Keys)
	if !rotation { // if the number of keys is the same, we need to check every key
		for kid := range newRootRole.Keys {
			if _, ok := oldRootRole.Keys[kid]; !ok {
				// if there is any difference in keys, a key rotation may have
				// occurred.
				rotation = true
			}
		}
	}

	if rotation {
		if err := signed.VerifyRoot(newSigned, oldRootRole.Threshold, oldRootRole.Keys); err != nil {
			return fmt.Errorf("rotation detected and new root was not signed with at least %d old keys",
				oldRootRole.Threshold)
		}
	}

	return nil
}
