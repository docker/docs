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

	"github.com/docker/notary/server/storage"
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
	oldRootJSON, err := store.GetCurrent(gun, rootRole)
	if _, ok := err.(storage.ErrNotFound); err != nil && !ok {
		// problem with storage. No expectation we can
		// write if we can't read so bail.
		logrus.Error("error reading previous root: ", err.Error())
		return nil, err
	}
	if rootUpdate, ok := roles[rootRole]; ok {
		// if root is present, validate its integrity, possibly
		// against a previous root
		if root, err = validateRoot(gun, oldRootJSON, rootUpdate.Data, store); err != nil {
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
		oldSnapJSON, err := store.GetCurrent(gun, snapshotRole)
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

		if err := validateSnapshot(snapshotRole, oldSnap, roles[snapshotRole], roles, repo); err != nil {
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
	return updatesToApply, nil
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
	tgtJSON, err := store.GetCurrent(gun, role)
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

func generateSnapshot(gun string, repo *tuf.Repo, store storage.MetaStore) (*storage.MetaUpdate, error) {
	role, err := repo.GetBaseRole(data.CanonicalSnapshotRole)
	if err != nil {
		return nil, validation.ErrBadRoot{Msg: "root did not include snapshot role"}
	}

	algo, keyBytes, err := store.GetKey(gun, data.CanonicalSnapshotRole)
	if err != nil {
		return nil, validation.ErrBadHierarchy{Msg: "could not retrieve snapshot key. client must provide snapshot"}
	}
	foundK := data.NewPublicKey(algo, keyBytes)

	validKey := false
	for _, id := range role.ListKeyIDs() {
		if id == foundK.ID() {
			validKey = true
			break
		}
	}
	if !validKey {
		return nil, validation.ErrBadHierarchy{
			Missing: data.CanonicalSnapshotRole,
			Msg:     "no snapshot was included in update and server does not hold current snapshot key for repository"}
	}

	currentJSON, err := store.GetCurrent(gun, data.CanonicalSnapshotRole)
	if err != nil {
		if _, ok := err.(storage.ErrNotFound); !ok {
			return nil, validation.ErrValidation{Msg: err.Error()}
		}
	}
	var sn *data.SignedSnapshot
	if currentJSON != nil {
		sn = new(data.SignedSnapshot)
		err := json.Unmarshal(currentJSON, sn)
		if err != nil {
			return nil, validation.ErrValidation{Msg: err.Error()}
		}
		err = repo.SetSnapshot(sn)
		if err != nil {
			return nil, validation.ErrValidation{Msg: err.Error()}
		}
	} else {
		// this will only occurr if no snapshot has ever been created for the repository
		err := repo.InitSnapshot()
		if err != nil {
			return nil, validation.ErrBadSnapshot{Msg: err.Error()}
		}
	}
	sgnd, err := repo.SignSnapshot(data.DefaultExpires(data.CanonicalSnapshotRole))
	if err != nil {
		return nil, validation.ErrBadSnapshot{Msg: err.Error()}
	}
	sgndJSON, err := json.Marshal(sgnd)
	if err != nil {
		return nil, validation.ErrBadSnapshot{Msg: err.Error()}
	}
	return &storage.MetaUpdate{
		Role:    data.CanonicalSnapshotRole,
		Version: repo.Snapshot.Signed.Version,
		Data:    sgndJSON,
	}, nil
}

func validateSnapshot(role string, oldSnap *data.SignedSnapshot, snapUpdate storage.MetaUpdate, roles map[string]storage.MetaUpdate, repo *tuf.Repo) error {
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

func validateRoot(gun string, oldRoot, newRoot []byte, store storage.MetaStore) (
	*data.SignedRoot, error) {

	var parsedOldRoot *data.SignedRoot
	parsedNewRoot := &data.SignedRoot{}

	if oldRoot != nil {
		parsedOldRoot = &data.SignedRoot{}
		err := json.Unmarshal(oldRoot, parsedOldRoot)
		if err != nil {
			// TODO(david): if we can't read the old root should we continue
			//             here to check new root self referential integrity?
			//             This would permit recovery of a repo with a corrupted
			//             root.
			logrus.Warn("Old root could not be parsed.")
		}
	}
	err := json.Unmarshal(newRoot, parsedNewRoot)
	if err != nil {
		return nil, err
	}

	// Don't update if a timestamp key doesn't exist.
	algo, keyBytes, err := store.GetKey(gun, data.CanonicalTimestampRole)
	if err != nil || algo == "" || keyBytes == nil {
		return nil, fmt.Errorf("no timestamp key for %s", gun)
	}
	timestampKey := data.NewPublicKey(algo, keyBytes)

	if err := checkRoot(parsedOldRoot, parsedNewRoot, timestampKey); err != nil {
		// TODO(david): how strict do we want to be here about old signatures
		//              for rotations? Should the user have to provide a flag
		//              which gets transmitted to force a root update without
		//              correct old key signatures.
		return nil, err
	}

	if !data.ValidTUFType(parsedNewRoot.Signed.Type, data.CanonicalRootRole) {
		return nil, fmt.Errorf("root has wrong type")
	}
	return parsedNewRoot, nil
}

// checkRoot errors if an invalid rotation has taken place, if the
// threshold number of signatures is invalid, if there are an invalid
// number of roles and keys, or if the timestamp keys are invalid
func checkRoot(oldRoot, newRoot *data.SignedRoot, timestampKey data.PublicKey) error {
	rootRole := data.CanonicalRootRole
	targetsRole := data.CanonicalTargetsRole
	snapshotRole := data.CanonicalSnapshotRole
	timestampRole := data.CanonicalTimestampRole

	var oldRootRole *data.RootRole
	newRootRole, ok := newRoot.Signed.Roles[rootRole]
	if !ok {
		return errors.New("new root is missing role entry for root role")
	}

	oldThreshold := 1
	rotation := false
	oldKeys := map[string]data.PublicKey{}
	newKeys := map[string]data.PublicKey{}
	if oldRoot != nil {
		// check for matching root key IDs
		oldRootRole = oldRoot.Signed.Roles[rootRole]
		oldThreshold = oldRootRole.Threshold

		for _, kid := range oldRootRole.KeyIDs {
			k, ok := oldRoot.Signed.Keys[kid]
			if !ok {
				// if the key itself wasn't contained in the root
				// we're skipping it because it could never have
				// been used to validate this root.
				continue
			}
			oldKeys[kid] = data.NewPublicKey(k.Algorithm(), k.Public())
		}

		// super simple check for possible rotation
		rotation = len(oldKeys) != len(newRootRole.KeyIDs)
	}
	// if old and new had the same number of keys, iterate
	// to see if there's a difference.
	for _, kid := range newRootRole.KeyIDs {
		k, ok := newRoot.Signed.Keys[kid]
		if !ok {
			// if the key itself wasn't contained in the root
			// we're skipping it because it could never have
			// been used to validate this root.
			continue
		}
		newKeys[kid] = data.NewPublicKey(k.Algorithm(), k.Public())

		if oldRoot != nil {
			if _, ok := oldKeys[kid]; !ok {
				// if there is any difference in keys, a key rotation may have
				// occurred.
				rotation = true
			}
		}
	}
	newSigned, err := newRoot.ToSigned()
	if err != nil {
		return err
	}
	if rotation {
		err = signed.VerifyRoot(newSigned, oldThreshold, oldKeys)
		if err != nil {
			return fmt.Errorf("rotation detected and new root was not signed with at least %d old keys", oldThreshold)
		}
	}
	err = signed.VerifyRoot(newSigned, newRootRole.Threshold, newKeys)
	if err != nil {
		return err
	}
	root, err := data.RootFromSigned(newSigned)
	if err != nil {
		return err
	}

	var timestampKeyIDs []string

	// at a minimum, check the 4 required roles are present
	for _, r := range []string{rootRole, targetsRole, snapshotRole, timestampRole} {
		role, ok := root.Signed.Roles[r]
		if !ok {
			return fmt.Errorf("missing required %s role from root", r)
		}
		// According to the TUF spec, any role may have more than one signing
		// key and require a threshold signature.  However, notary-server
		// creates the timestamp, and there is only ever one, so a threshold
		// greater than one would just always fail validation
		if (r == timestampRole && role.Threshold != 1) || role.Threshold < 1 {
			return fmt.Errorf("%s role has invalid threshold", r)
		}
		if len(role.KeyIDs) < role.Threshold {
			return fmt.Errorf("%s role has insufficient number of keys", r)
		}

		if r == timestampRole {
			timestampKeyIDs = role.KeyIDs
		}
	}

	// ensure that at least one of the timestamp keys specified in the role
	// actually exists

	for _, keyID := range timestampKeyIDs {
		if timestampKey.ID() == keyID {
			return nil
		}
	}
	return fmt.Errorf("none of the following timestamp keys exist: %s",
		strings.Join(timestampKeyIDs, ", "))
}
