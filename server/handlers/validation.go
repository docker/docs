package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/docker/notary/tuf"
	"github.com/docker/notary/tuf/data"

	"github.com/Sirupsen/logrus"
	"github.com/docker/notary/server/storage"
	"github.com/docker/notary/tuf/keys"
	"github.com/docker/notary/tuf/signed"
	"github.com/docker/notary/tuf/utils"
)

// ErrValidation represents a general validation error
type ErrValidation struct {
	msg string
}

func (err ErrValidation) Error() string {
	return fmt.Sprintf("An error occurred during validation: %s", err.msg)
}

// ErrBadHierarchy represents a missing snapshot at this current time.
// When delegations are implemented it will also represent a missing
// delegation parent
type ErrBadHierarchy struct {
	msg string
}

func (err ErrBadHierarchy) Error() string {
	return fmt.Sprintf("Hierarchy of updates in incorrect: %s", err.msg)
}

// ErrBadRoot represents a failure validating the root
type ErrBadRoot struct {
	msg string
}

func (err ErrBadRoot) Error() string {
	return fmt.Sprintf("The root being updated is invalid: %s", err.msg)
}

// ErrBadTargets represents a failure to validate a targets (incl delegations)
type ErrBadTargets struct {
	msg string
}

func (err ErrBadTargets) Error() string {
	return fmt.Sprintf("The targets being updated is invalid: %s", err.msg)
}

// ErrBadSnapshot represents a failure to validate the snapshot
type ErrBadSnapshot struct {
	msg string
}

func (err ErrBadSnapshot) Error() string {
	return fmt.Sprintf("The snapshot being updated is invalid: %s", err.msg)
}

// validateUpload checks that the updates being pushed
// are semantically correct and the signatures are correct
func validateUpdate(gun string, updates []storage.MetaUpdate, store storage.MetaStore) error {
	kdb := keys.NewDB()
	repo := tuf.NewRepo(kdb, nil)
	rootRole := data.RoleName(data.CanonicalRootRole)
	targetsRole := data.RoleName(data.CanonicalTargetsRole)
	snapshotRole := data.RoleName(data.CanonicalSnapshotRole)

	// check that the necessary roles are present:
	roles := make(map[string]storage.MetaUpdate)
	for _, v := range updates {
		roles[v.Role] = v
	}
	if err := hierarchyOK(roles); err != nil {
		logrus.Error("ErrBadHierarchy: ", err.Error())
		return ErrBadHierarchy{msg: err.Error()}
	}
	logrus.Debug("Successfully validated hierarchy")

	var root *data.SignedRoot
	oldRootJSON, err := store.GetCurrent(gun, rootRole)
	if _, ok := err.(*storage.ErrNotFound); err != nil && !ok {
		// problem with storage. No expectation we can
		// write if we can't read so bail.
		logrus.Error("error reading previous root: ", err.Error())
		return err
	}
	if rootUpdate, ok := roles[rootRole]; ok {
		// if root is present, validate its integrity, possibly
		// against a previous root
		if root, err = validateRoot(gun, oldRootJSON, rootUpdate.Data); err != nil {
			logrus.Error("ErrBadRoot: ", err.Error())
			return ErrBadRoot{msg: err.Error()}
		}
		// setting root will update keys db
		if err = repo.SetRoot(root); err != nil {
			logrus.Error("ErrValidation: ", err.Error())
			return ErrValidation{msg: err.Error()}
		}
		logrus.Debug("Successfully validated root")
	} else {
		if oldRootJSON == nil {
			return ErrValidation{msg: "no pre-existing root and no root provided in update."}
		}
		parsedOldRoot := &data.SignedRoot{}
		if err := json.Unmarshal(oldRootJSON, parsedOldRoot); err != nil {
			return ErrValidation{msg: "pre-existing root is corrupted and no root provided in update."}
		}
		if err = repo.SetRoot(parsedOldRoot); err != nil {
			logrus.Error("ErrValidation: ", err.Error())
			return ErrValidation{msg: err.Error()}
		}
	}

	// TODO: validate delegated targets roles.
	var t *data.SignedTargets
	if _, ok := roles[targetsRole]; ok {
		if t, err = validateTargets(targetsRole, roles, kdb); err != nil {
			logrus.Error("ErrBadTargets: ", err.Error())
			return ErrBadTargets{msg: err.Error()}
		}
		repo.SetTargets(targetsRole, t)
	}
	logrus.Debug("Successfully validated targets")

	var oldSnap *data.SignedSnapshot
	oldSnapJSON, err := store.GetCurrent(gun, snapshotRole)
	if _, ok := err.(*storage.ErrNotFound); err != nil && !ok {
		// problem with storage. No expectation we can
		// write if we can't read so bail.
		logrus.Error("error reading previous snapshot: ", err.Error())
		return err
	} else if err == nil {
		oldSnap = &data.SignedSnapshot{}
		if err := json.Unmarshal(oldSnapJSON, oldSnap); err != nil {
			oldSnap = nil
		}
	}

	if err := validateSnapshot(snapshotRole, oldSnap, roles[snapshotRole], roles, kdb); err != nil {
		logrus.Error("ErrBadSnapshot: ", err.Error())
		return ErrBadSnapshot{msg: err.Error()}
	}
	logrus.Debug("Successfully validated snapshot")
	return nil
}

func validateSnapshot(role string, oldSnap *data.SignedSnapshot, snapUpdate storage.MetaUpdate, roles map[string]storage.MetaUpdate, kdb *keys.KeyDB) error {
	s := &data.Signed{}
	err := json.Unmarshal(snapUpdate.Data, s)
	if err != nil {
		return errors.New("could not parse snapshot")
	}
	// version specifically gets validated when writing to store to
	// better handle race conditions there.
	if err := signed.Verify(s, role, 0, kdb); err != nil {
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
	snapshotRole := data.RoleName(data.CanonicalSnapshotRole)
	timestampRole := data.RoleName(data.CanonicalTimestampRole) // just in case
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

func validateTargets(role string, roles map[string]storage.MetaUpdate, kdb *keys.KeyDB) (*data.SignedTargets, error) {
	// TODO: when delegations are being validated, validate parent
	//       role exists for any delegation
	s := &data.Signed{}
	err := json.Unmarshal(roles[role].Data, s)
	if err != nil {
		return nil, fmt.Errorf("could not parse %s", role)
	}
	// version specifically gets validated when writing to store to
	// better handle race conditions there.
	if err := signed.Verify(s, role, 0, kdb); err != nil {
		return nil, err
	}
	t, err := data.TargetsFromSigned(s)
	if err != nil {
		return nil, err
	}
	if !data.ValidTUFType(t.Signed.Type, data.CanonicalTargetsRole) {
		return nil, fmt.Errorf("%s has wrong type", role)
	}
	return t, nil
}

// check the snapshot is present. If it is, the hierarchy
// of the update is OK. This seems like a simplistic check
// but is completely sufficient for all possible use cases:
// 1. the user is updating only the snapshot.
// 2. the user is updating a targets (incl. delegations) or
//    root metadata. This requires they also provide a new
//    snapshot.
// N.B. users should never be updating timestamps. The server
//      always handles timestamping. If the user does send a
//      timestamp, the server will replace it on next
//      GET timestamp.jsonshould it detect the current
//      snapshot has a different hash to the one in the timestamp.
func hierarchyOK(roles map[string]storage.MetaUpdate) error {
	snapshotRole := data.RoleName(data.CanonicalSnapshotRole)
	if _, ok := roles[snapshotRole]; !ok {
		return errors.New("snapshot missing from update")
	}
	return nil
}

func validateRoot(gun string, oldRoot, newRoot []byte) (*data.SignedRoot, error) {
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
	if err := checkRoot(parsedOldRoot, parsedNewRoot); err != nil {
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

// checkRoot returns true if no rotation, or a valid
// rotation has taken place, and the threshold number of signatures
// are valid.
func checkRoot(oldRoot, newRoot *data.SignedRoot) error {
	rootRole := data.RoleName(data.CanonicalRootRole)
	targetsRole := data.RoleName(data.CanonicalTargetsRole)
	snapshotRole := data.RoleName(data.CanonicalSnapshotRole)
	timestampRole := data.RoleName(data.CanonicalTimestampRole)

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
	// at a minimum, check the 4 required roles are present
	for _, r := range []string{rootRole, targetsRole, snapshotRole, timestampRole} {
		role, ok := root.Signed.Roles[r]
		if !ok {
			return fmt.Errorf("missing required %s role from root", r)
		}
		if role.Threshold < 1 {
			return fmt.Errorf("%s role has invalid threshold", r)
		}
		if len(role.KeyIDs) < role.Threshold {
			return fmt.Errorf("%s role has insufficient number of keys", r)
		}
	}
	return nil
}
