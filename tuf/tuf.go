// Package tuf defines the core TUF logic around manipulating a repo.
package tuf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/signed"
	"github.com/docker/notary/tuf/utils"
)

// ErrSigVerifyFail - signature verification failed
type ErrSigVerifyFail struct{}

func (e ErrSigVerifyFail) Error() string {
	return "Error: Signature verification failed"
}

// ErrMetaExpired - metadata file has expired
type ErrMetaExpired struct{}

func (e ErrMetaExpired) Error() string {
	return "Error: Metadata has expired"
}

// ErrLocalRootExpired - the local root file is out of date
type ErrLocalRootExpired struct{}

func (e ErrLocalRootExpired) Error() string {
	return "Error: Local Root Has Expired"
}

// ErrNotLoaded - attempted to access data that has not been loaded into
// the repo. This means specifically that the relevant JSON file has not
// been loaded.
type ErrNotLoaded struct {
	Role string
}

func (err ErrNotLoaded) Error() string {
	return fmt.Sprintf("%s role has not been loaded", err.Role)
}

// ErrStopWalk - used by visitor functions to signal WalkTargets to stop walking
type ErrStopWalk struct{}

func (e ErrStopWalk) Error() string {
	return "Error: signalled to stop the walk"
}

// ErrContinueWalk - used by visitor functions to signal WalkTargets to continue walking
type ErrContinueWalk struct{}

func (e ErrContinueWalk) Error() string {
	return "Error: signalled to continue the walk"
}

// Repo is an in memory representation of the TUF Repo.
// It operates at the data.Signed level, accepting and producing
// data.Signed objects. Users of a Repo are responsible for
// fetching raw JSON and using the Set* functions to populate
// the Repo instance.
type Repo struct {
	Root          *data.SignedRoot
	Targets       map[string]*data.SignedTargets
	Snapshot      *data.SignedSnapshot
	Timestamp     *data.SignedTimestamp
	cryptoService signed.CryptoService
}

// NewRepo initializes a Repo instance with a CryptoService.
// If the Repo will only be used for reading, the CryptoService
// can be nil.
func NewRepo(cryptoService signed.CryptoService) *Repo {
	repo := &Repo{
		Targets:       make(map[string]*data.SignedTargets),
		cryptoService: cryptoService,
	}
	return repo
}

// AddBaseKeys is used to add keys to the role in root.json
func (tr *Repo) AddBaseKeys(role string, keys ...data.PublicKey) error {
	if tr.Root == nil {
		return ErrNotLoaded{Role: data.CanonicalRootRole}
	}
	ids := []string{}
	for _, k := range keys {
		// Store only the public portion
		tr.Root.Signed.Keys[k.ID()] = k
		tr.Root.Signed.Roles[role].KeyIDs = append(tr.Root.Signed.Roles[role].KeyIDs, k.ID())
		ids = append(ids, k.ID())
	}
	tr.Root.Dirty = true

	// also, whichever role was switched out needs to be re-signed
	// root has already been marked dirty
	switch role {
	case data.CanonicalSnapshotRole:
		if tr.Snapshot != nil {
			tr.Snapshot.Dirty = true
		}
	case data.CanonicalTargetsRole:
		if target, ok := tr.Targets[data.CanonicalTargetsRole]; ok {
			target.Dirty = true
		}
	case data.CanonicalTimestampRole:
		if tr.Timestamp != nil {
			tr.Timestamp.Dirty = true
		}
	}
	return nil
}

// ReplaceBaseKeys is used to replace all keys for the given role with the new keys
func (tr *Repo) ReplaceBaseKeys(role string, keys ...data.PublicKey) error {
	r, err := tr.GetBaseRole(role)
	if err != nil {
		return err
	}
	err = tr.RemoveBaseKeys(role, r.ListKeyIDs()...)
	if err != nil {
		return err
	}
	return tr.AddBaseKeys(role, keys...)
}

// RemoveBaseKeys is used to remove keys from the roles in root.json
func (tr *Repo) RemoveBaseKeys(role string, keyIDs ...string) error {
	if tr.Root == nil {
		return ErrNotLoaded{Role: data.CanonicalRootRole}
	}
	var keep []string
	toDelete := make(map[string]struct{})
	// remove keys from specified role
	for _, k := range keyIDs {
		toDelete[k] = struct{}{}
		for _, rk := range tr.Root.Signed.Roles[role].KeyIDs {
			if k != rk {
				keep = append(keep, rk)
			}
		}
	}
	tr.Root.Signed.Roles[role].KeyIDs = keep

	// determine which keys are no longer in use by any roles
	for roleName, r := range tr.Root.Signed.Roles {
		if roleName == role {
			continue
		}
		for _, rk := range r.KeyIDs {
			if _, ok := toDelete[rk]; ok {
				delete(toDelete, rk)
			}
		}
	}

	// remove keys no longer in use by any roles
	for k := range toDelete {
		delete(tr.Root.Signed.Keys, k)
		// remove the signing key from the cryptoservice if it
		// isn't a root key. Root keys must be kept for rotation
		// signing
		if role != data.CanonicalRootRole {
			tr.cryptoService.RemoveKey(k)
		}
	}
	tr.Root.Dirty = true
	return nil
}

// GetBaseRole gets a base role from this repo's metadata
func (tr *Repo) GetBaseRole(name string) (data.BaseRole, error) {
	if !data.ValidRole(name) {
		return data.BaseRole{}, data.ErrInvalidRole{Role: name, Reason: "invalid base role name"}
	}
	if tr.Root == nil {
		return data.BaseRole{}, ErrNotLoaded{data.CanonicalRootRole}
	}
	roleData, ok := tr.Root.Signed.Roles[name]
	if !ok {
		return data.BaseRole{}, data.ErrInvalidRole{Role: name, Reason: "role not found in root file"}
	}
	// Get all public keys for the base role from TUF metadata
	keyIDs := roleData.KeyIDs
	pubKeys := make(map[string]data.PublicKey)
	for _, keyID := range keyIDs {
		pubKey, ok := tr.Root.Signed.Keys[keyID]
		if !ok {
			return data.BaseRole{}, data.ErrInvalidRole{
				Role:   name,
				Reason: fmt.Sprintf("key with ID %s was not found in root metadata", keyID),
			}
		}
		pubKeys[keyID] = pubKey
	}

	return data.BaseRole{
		Name:      name,
		Keys:      pubKeys,
		Threshold: roleData.Threshold,
	}, nil
}

// GetDelegationRole gets a delegation role from this repo's metadata, walking from the targets role down to the delegation itself
func (tr *Repo) GetDelegationRole(name string) (data.DelegationRole, error) {
	if !data.IsDelegation(name) {
		return data.DelegationRole{}, data.ErrInvalidRole{Role: name, Reason: "invalid delegation name"}
	}
	if tr.Root == nil {
		return data.DelegationRole{}, ErrNotLoaded{data.CanonicalRootRole}
	}
	_, ok := tr.Root.Signed.Roles[data.CanonicalTargetsRole]
	if !ok {
		return data.DelegationRole{}, ErrNotLoaded{data.CanonicalTargetsRole}
	}
	// Traverse target metadata, down to delegation itself
	// Get all public keys for the base role from TUF metadata
	_, ok = tr.Targets[data.CanonicalTargetsRole]
	if !ok {
		return data.DelegationRole{}, ErrNotLoaded{data.CanonicalTargetsRole}
	}

	// Start with top level roles in targets. Walk the chain of ancestors
	// until finding the desired role, or we run out of targets files to search.
	var foundRole *data.DelegationRole
	buildDelegationRoleVisitor := func(tgt *data.SignedTargets, roleName string) error {
		if tgt == nil {
			return ErrContinueWalk{}
		}

		// Try to find the delegation and build a DelegationRole structure
		for _, role := range tgt.Signed.Delegations.Roles {
			if role.Name == name {
				pubKeys := make(map[string]data.PublicKey)
				for _, keyID := range role.KeyIDs {
					pubKey, ok := tgt.Signed.Delegations.Keys[keyID]
					if !ok {
						// Couldn't retrieve all keys, so stop walking
						return ErrStopWalk{}
					}
					pubKeys[keyID] = pubKey
				}
				foundRole = &data.DelegationRole{
					BaseRole: data.BaseRole{
						Name:      role.Name,
						Keys:      pubKeys,
						Threshold: role.Threshold,
					},
					Paths: role.Paths,
				}
				return ErrStopWalk{}
			}
		}
		return ErrContinueWalk{}
	}

	// Walk to the parent of this delegation, since that is where its role metadata exists
	tr.WalkTargets("", path.Dir(name), buildDelegationRoleVisitor)

	// We never found the delegation. In the context of this repo it is considered
	// invalid. N.B. it may be that it existed at one point but an ancestor has since
	// been modified/removed.
	if foundRole == nil {
		return data.DelegationRole{}, data.ErrInvalidRole{Role: name, Reason: "delegation does not exist"}
	}

	return *foundRole, nil
}

// GetAllLoadedRoles returns a list of all role entries loaded in this TUF repo, could be empty
func (tr *Repo) GetAllLoadedRoles() []*data.Role {
	var res []*data.Role
	if tr.Root == nil {
		// if root isn't loaded, we should consider we have no loaded roles because we can't
		// trust any other state that might be present
		return res
	}
	for name, rr := range tr.Root.Signed.Roles {
		res = append(res, &data.Role{
			RootRole: *rr,
			Name:     name,
		})
	}
	for _, delegate := range tr.Targets {
		for _, r := range delegate.Signed.Delegations.Roles {
			res = append(res, r)
		}
	}
	return res
}

// GetDelegation finds the role entry representing the provided
// role name along with its associated public keys, or ErrInvalidRole
func (tr *Repo) GetDelegation(role string) (*data.Role, data.Keys, error) {
	if !data.IsDelegation(role) {
		return nil, nil, data.ErrInvalidRole{Role: role, Reason: "not a valid delegated role"}
	}

	parent := path.Dir(role)

	// check the parent role
	if _, err := tr.GetDelegationRole(parent); parent != data.CanonicalTargetsRole && err != nil {
		return nil, nil, data.ErrInvalidRole{Role: role, Reason: "parent role not found"}
	}

	// check the parent role's metadata
	p, ok := tr.Targets[parent]
	if !ok { // the parent targetfile may not exist yet, so it can't be in the list
		return nil, nil, data.ErrNoSuchRole{Role: role}
	}

	foundAt := utils.FindRoleIndex(p.Signed.Delegations.Roles, role)
	if foundAt < 0 {
		return nil, nil, data.ErrNoSuchRole{Role: role}
	}
	delegationRole := p.Signed.Delegations.Roles[foundAt]
	keys := make(data.Keys)
	for _, keyID := range delegationRole.KeyIDs {
		keys[keyID] = p.Signed.Delegations.Keys[keyID]
	}
	return delegationRole, keys, nil
}

// UpdateDelegations updates the appropriate delegations, either adding
// a new delegation or updating an existing one. If keys are
// provided, the IDs will be added to the role (if they do not exist
// there already), and the keys will be added to the targets file.
func (tr *Repo) UpdateDelegations(role *data.Role, keys []data.PublicKey) error {
	if !data.IsDelegation(role.Name) {
		return data.ErrInvalidRole{Role: role.Name, Reason: "not a valid delegated role"}
	}
	parent := path.Dir(role.Name)

	if err := tr.VerifyCanSign(parent); err != nil {
		return err
	}

	// check the parent role's metadata
	p, ok := tr.Targets[parent]
	if !ok { // the parent targetfile may not exist yet - if not, then create it
		var err error
		p, err = tr.InitTargets(parent)
		if err != nil {
			return err
		}
	}

	for _, k := range keys {
		if !utils.StrSliceContains(role.KeyIDs, k.ID()) {
			role.KeyIDs = append(role.KeyIDs, k.ID())
		}
		p.Signed.Delegations.Keys[k.ID()] = k
	}

	// if the role has fewer keys than the threshold, it
	// will never be able to create a valid targets file
	// and should be considered invalid.
	if len(role.KeyIDs) < role.Threshold {
		return data.ErrInvalidRole{Role: role.Name, Reason: "insufficient keys to meet threshold"}
	}

	foundAt := utils.FindRoleIndex(p.Signed.Delegations.Roles, role.Name)

	if foundAt >= 0 {
		p.Signed.Delegations.Roles[foundAt] = role
	} else {
		p.Signed.Delegations.Roles = append(p.Signed.Delegations.Roles, role)
	}
	// We've made a change to parent. Set it to dirty
	p.Dirty = true

	// We don't actually want to create the new delegation metadata yet.
	// When we add a delegation, it may only be signable by a key we don't have
	// (hence we are delegating signing).

	utils.RemoveUnusedKeys(p)

	return nil
}

// DeleteDelegation removes a delegated targets role from its parent
// targets object. It also deletes the delegation from the snapshot.
// DeleteDelegation will only make use of the role Name field.
func (tr *Repo) DeleteDelegation(role data.Role) error {
	if !data.IsDelegation(role.Name) {
		return data.ErrInvalidRole{Role: role.Name, Reason: "not a valid delegated role"}
	}
	// the role variable must not be used past this assignment for safety
	name := role.Name

	parent := path.Dir(name)
	if err := tr.VerifyCanSign(parent); err != nil {
		return err
	}

	// delete delegated data from Targets map and Snapshot - if they don't
	// exist, these are no-op
	delete(tr.Targets, name)
	tr.Snapshot.DeleteMeta(name)

	p, ok := tr.Targets[parent]
	if !ok {
		// if there is no parent metadata (the role exists though), then this
		// is as good as done.
		return nil
	}

	foundAt := utils.FindRoleIndex(p.Signed.Delegations.Roles, name)

	if foundAt >= 0 {
		var roles []*data.Role
		// slice out deleted role
		roles = append(roles, p.Signed.Delegations.Roles[:foundAt]...)
		if foundAt+1 < len(p.Signed.Delegations.Roles) {
			roles = append(roles, p.Signed.Delegations.Roles[foundAt+1:]...)
		}
		p.Signed.Delegations.Roles = roles

		utils.RemoveUnusedKeys(p)

		p.Dirty = true
	} // if the role wasn't found, it's a good as deleted

	return nil
}

// InitRoot initializes an empty root file with the 4 core roles passed to the
// method, and the consistent flag.
func (tr *Repo) InitRoot(root, timestamp, snapshot, targets data.BaseRole, consistent bool) error {
	rootRoles := make(map[string]*data.RootRole)
	rootKeys := make(map[string]data.PublicKey)

	for _, r := range []data.BaseRole{root, timestamp, snapshot, targets} {
		rootRoles[r.Name] = &data.RootRole{
			Threshold: r.Threshold,
			KeyIDs:    r.ListKeyIDs(),
		}
		for kid, k := range r.Keys {
			rootKeys[kid] = k
		}
	}
	r, err := data.NewRoot(rootKeys, rootRoles, consistent)
	if err != nil {
		return err
	}
	tr.Root = r
	return nil
}

// InitTargets initializes an empty targets, and returns the new empty target
func (tr *Repo) InitTargets(role string) (*data.SignedTargets, error) {
	if !data.IsDelegation(role) && role != data.CanonicalTargetsRole {
		return nil, data.ErrInvalidRole{
			Role:   role,
			Reason: fmt.Sprintf("role is not a valid targets role name: %s", role),
		}
	}
	targets := data.NewTargets()
	tr.Targets[role] = targets
	return targets, nil
}

// InitSnapshot initializes a snapshot based on the current root and targets
func (tr *Repo) InitSnapshot() error {
	if tr.Root == nil {
		return ErrNotLoaded{Role: data.CanonicalRootRole}
	}
	root, err := tr.Root.ToSigned()
	if err != nil {
		return err
	}

	if _, ok := tr.Targets[data.CanonicalTargetsRole]; !ok {
		return ErrNotLoaded{Role: data.CanonicalTargetsRole}
	}
	targets, err := tr.Targets[data.CanonicalTargetsRole].ToSigned()
	if err != nil {
		return err
	}
	snapshot, err := data.NewSnapshot(root, targets)
	if err != nil {
		return err
	}
	tr.Snapshot = snapshot
	return nil
}

// InitTimestamp initializes a timestamp based on the current snapshot
func (tr *Repo) InitTimestamp() error {
	snap, err := tr.Snapshot.ToSigned()
	if err != nil {
		return err
	}
	timestamp, err := data.NewTimestamp(snap)
	if err != nil {
		return err
	}

	tr.Timestamp = timestamp
	return nil
}

// SetRoot sets the Repo.Root field to the SignedRoot object.
func (tr *Repo) SetRoot(s *data.SignedRoot) error {
	tr.Root = s
	return nil
}

// SetTimestamp parses the Signed object into a SignedTimestamp object
// and sets the Repo.Timestamp field.
func (tr *Repo) SetTimestamp(s *data.SignedTimestamp) error {
	tr.Timestamp = s
	return nil
}

// SetSnapshot parses the Signed object into a SignedSnapshots object
// and sets the Repo.Snapshot field.
func (tr *Repo) SetSnapshot(s *data.SignedSnapshot) error {
	tr.Snapshot = s
	return nil
}

// SetTargets sets the SignedTargets object agaist the role in the
// Repo.Targets map.
func (tr *Repo) SetTargets(role string, s *data.SignedTargets) error {
	tr.Targets[role] = s
	return nil
}

// TargetMeta returns the FileMeta entry for the given path in the
// targets file associated with the given role. This may be nil if
// the target isn't found in the targets file.
func (tr Repo) TargetMeta(role, path string) *data.FileMeta {
	if t, ok := tr.Targets[role]; ok {
		if m, ok := t.Signed.Targets[path]; ok {
			return &m
		}
	}
	return nil
}

// TargetDelegations returns a slice of Roles that are valid publishers
// for the target path provided.
func (tr Repo) TargetDelegations(role, path string) []*data.Role {
	var roles []*data.Role
	if t, ok := tr.Targets[role]; ok {
		for _, r := range t.Signed.Delegations.Roles {
			if r.CheckPaths(path) {
				roles = append(roles, r)
			}
		}
	}
	return roles
}

// VerifyCanSign returns nil if the role exists and we have at least one
// signing key for the role, false otherwise.  This does not check that we have
// enough signing keys to meet the threshold, since we want to support the use
// case of multiple signers for a role.  It returns an error if the role doesn't
// exist or if there are no signing keys.
func (tr *Repo) VerifyCanSign(roleName string) error {
	var (
		role data.BaseRole
		err  error
	)
	// we only need the BaseRole part of a delegation because we're just
	// checking KeyIDs
	if data.IsDelegation(roleName) {
		r, err := tr.GetDelegationRole(roleName)
		if err != nil {
			return err
		}
		role = r.BaseRole
	} else {
		role, err = tr.GetBaseRole(roleName)
	}
	if err != nil {
		return data.ErrInvalidRole{Role: roleName, Reason: "does not exist"}
	}

	for keyID, k := range role.Keys {
		check := []string{keyID}
		if canonicalID, err := utils.CanonicalKeyID(k); err == nil {
			check = append(check, canonicalID)
		}
		for _, id := range check {
			p, _, err := tr.cryptoService.GetPrivateKey(id)
			if err == nil && p != nil {
				return nil
			}
		}
	}
	return signed.ErrNoKeys{KeyIDs: role.ListKeyIDs()}
}

// used for walking the targets/delegations tree, potentially modifying the underlying SignedTargets for the repo
type walkVisitorFunc func(*data.SignedTargets, string) error

// WalkTargets will apply the specified visitor function to iteratively walk the targets/delegation metadata tree,
// until receiving a ErrStopWalk.  The walk starts from the base "targets" role, and searches for the correct targetPath and/or rolePath
// to call the visitor function on.
func (tr *Repo) WalkTargets(targetPath, rolePath string, visitTarget walkVisitorFunc) error {
	// Start with the base targets role, which implicitly has the "" targets path
	targetsRole, ok := tr.Root.Signed.Roles[data.CanonicalTargetsRole]
	if !ok {
		return data.ErrInvalidRole{Role: data.CanonicalTargetsRole, Reason: "role not found in root file"}
	}
	// Make the targets role have the empty path, when we treat it as a delegation role
	roles := []*data.Role{
		{
			RootRole: *targetsRole,
			Name:     data.CanonicalTargetsRole,
			Paths:    []string{""},
		},
	}

	for len(roles) > 0 {
		role := roles[0]
		roles = roles[1:]

		// Check the role metadata
		signedTgt, ok := tr.Targets[role.Name]
		if !ok {
			// The role meta doesn't exist in the repo so continue onward
			continue
		}

		// We're at a prefix of the desired role subtree, so add its delegation role children and continue walking
		if strings.HasPrefix(rolePath, role.Name+"/") {
			roles = append(roles, signedTgt.Signed.Delegations.Roles...)
			continue
		}

		// Determine whether to visit this role or not:
		// If the paths validate against the specified targetPath and the rolePath is empty or is in the subtree
		if isValidPath(targetPath, role) && isAncestorRole(role.Name, rolePath) {
			// If we had matching path or role name, visit this target and determine whether or not to keep walking
			err := visitTarget(signedTgt, role.Name)
			switch err.(type) {
			case ErrStopWalk:
				// If the visitor function signalled a stop, return nil to finish the walk
				return nil
			case ErrContinueWalk:
				// If the visitor function signalled to continue, add this role's delegation to the walk
				roles = append(roles, signedTgt.Signed.Delegations.Roles...)
			default:
				// Return out if we got a different error or nil
				return err
			}

		}
	}
	return nil
}

// helper function that returns whether the candidateChild role name is an ancestor or equal to the candidateAncestor role name
// Will return true if given an empty candidateAncestor role name
func isAncestorRole(candidateChild, candidateAncestor string) bool {
	return candidateAncestor == "" || candidateAncestor == candidateChild || strings.HasPrefix(candidateChild, candidateAncestor+"/")
}

// helper function that returns whether the delegation Role is valid against the given path
// Will return true if given an empty candidatePath
func isValidPath(candidatePath string, delgRole *data.Role) bool {
	return candidatePath == "" || delgRole.CheckPaths(candidatePath)
}

// AddTargets will attempt to add the given targets specifically to
// the directed role. If the metadata for the role doesn't exist yet,
// AddTargets will create one.
func (tr *Repo) AddTargets(role string, targets data.Files) (data.Files, error) {

	err := tr.VerifyCanSign(role)
	if err != nil {
		return nil, err
	}

	// check the role's metadata
	t, ok := tr.Targets[role]
	if !ok { // the targetfile may not exist yet - if not, then create it
		var err error
		t, err = tr.InitTargets(role)
		if err != nil {
			return nil, err
		}
	}

	var r data.DelegationRole
	if role != data.CanonicalTargetsRole {
		// we only call r.CheckPaths if the role is not "targets"
		// so r being nil is fine in the case role == "targets"
		r, err = tr.GetDelegationRole(role)
		if err != nil {
			return nil, err
		}
	}

	invalid := make(data.Files)
	for path, target := range targets {
		if role == data.CanonicalTargetsRole || r.CheckPaths(path) {
			t.Signed.Targets[path] = target
		} else {
			invalid[path] = target
		}
	}
	t.Dirty = true
	if len(invalid) > 0 {
		return invalid, fmt.Errorf("Could not add all targets")
	}
	return nil, nil
}

// RemoveTargets removes the given target (paths) from the given target role (delegation)
func (tr *Repo) RemoveTargets(role string, targets ...string) error {
	if err := tr.VerifyCanSign(role); err != nil {
		return err
	}

	// if the role exists but metadata does not yet, then our work is done
	t, ok := tr.Targets[role]
	if ok {
		for _, path := range targets {
			delete(t.Signed.Targets, path)
		}
		t.Dirty = true
	}

	return nil
}

// UpdateSnapshot updates the FileMeta for the given role based on the Signed object
func (tr *Repo) UpdateSnapshot(role string, s *data.Signed) error {
	jsonData, err := json.Marshal(s)
	if err != nil {
		return err
	}
	meta, err := data.NewFileMeta(bytes.NewReader(jsonData), "sha256")
	if err != nil {
		return err
	}
	tr.Snapshot.Signed.Meta[role] = meta
	tr.Snapshot.Dirty = true
	return nil
}

// UpdateTimestamp updates the snapshot meta in the timestamp based on the Signed object
func (tr *Repo) UpdateTimestamp(s *data.Signed) error {
	jsonData, err := json.Marshal(s)
	if err != nil {
		return err
	}
	meta, err := data.NewFileMeta(bytes.NewReader(jsonData), "sha256")
	if err != nil {
		return err
	}
	tr.Timestamp.Signed.Meta["snapshot"] = meta
	tr.Timestamp.Dirty = true
	return nil
}

// SignRoot signs the root
func (tr *Repo) SignRoot(expires time.Time) (*data.Signed, error) {
	logrus.Debug("signing root...")
	tr.Root.Signed.Expires = expires
	tr.Root.Signed.Version++
	root, err := tr.GetBaseRole(data.CanonicalRootRole)
	if err != nil {
		return nil, err
	}
	signed, err := tr.Root.ToSigned()
	if err != nil {
		return nil, err
	}
	signed, err = tr.sign(signed, root)
	if err != nil {
		return nil, err
	}
	tr.Root.Signatures = signed.Signatures
	return signed, nil
}

// SignTargets signs the targets file for the given top level or delegated targets role
func (tr *Repo) SignTargets(role string, expires time.Time) (*data.Signed, error) {
	logrus.Debugf("sign targets called for role %s", role)
	if _, ok := tr.Targets[role]; !ok {
		return nil, data.ErrInvalidRole{
			Role:   role,
			Reason: "SignTargets called with non-existant targets role",
		}
	}
	tr.Targets[role].Signed.Expires = expires
	tr.Targets[role].Signed.Version++
	signed, err := tr.Targets[role].ToSigned()
	if err != nil {
		logrus.Debug("errored getting targets data.Signed object")
		return nil, err
	}

	var targets data.BaseRole
	if role == data.CanonicalTargetsRole {
		targets, err = tr.GetBaseRole(role)
	} else {
		tr, err := tr.GetDelegationRole(role)
		if err != nil {
			return nil, err
		}
		targets = tr.BaseRole
	}
	if err != nil {
		return nil, err
	}

	signed, err = tr.sign(signed, targets)
	if err != nil {
		logrus.Debug("errored signing ", role)
		return nil, err
	}
	tr.Targets[role].Signatures = signed.Signatures
	return signed, nil
}

// SignSnapshot updates the snapshot based on the current targets and root then signs it
func (tr *Repo) SignSnapshot(expires time.Time) (*data.Signed, error) {
	logrus.Debug("signing snapshot...")
	signedRoot, err := tr.Root.ToSigned()
	if err != nil {
		return nil, err
	}
	err = tr.UpdateSnapshot("root", signedRoot)
	if err != nil {
		return nil, err
	}
	tr.Root.Dirty = false // root dirty until changes captures in snapshot
	for role, targets := range tr.Targets {
		signedTargets, err := targets.ToSigned()
		if err != nil {
			return nil, err
		}
		err = tr.UpdateSnapshot(role, signedTargets)
		if err != nil {
			return nil, err
		}
		targets.Dirty = false
	}
	tr.Snapshot.Signed.Expires = expires
	tr.Snapshot.Signed.Version++
	signed, err := tr.Snapshot.ToSigned()
	if err != nil {
		return nil, err
	}
	snapshot, err := tr.GetBaseRole(data.CanonicalSnapshotRole)
	if err != nil {
		return nil, err
	}
	signed, err = tr.sign(signed, snapshot)
	if err != nil {
		return nil, err
	}
	tr.Snapshot.Signatures = signed.Signatures
	return signed, nil
}

// SignTimestamp updates the timestamp based on the current snapshot then signs it
func (tr *Repo) SignTimestamp(expires time.Time) (*data.Signed, error) {
	logrus.Debug("SignTimestamp")
	signedSnapshot, err := tr.Snapshot.ToSigned()
	if err != nil {
		return nil, err
	}
	err = tr.UpdateTimestamp(signedSnapshot)
	if err != nil {
		return nil, err
	}
	tr.Timestamp.Signed.Expires = expires
	tr.Timestamp.Signed.Version++
	signed, err := tr.Timestamp.ToSigned()
	if err != nil {
		return nil, err
	}
	timestamp, err := tr.GetBaseRole(data.CanonicalTimestampRole)
	if err != nil {
		return nil, err
	}
	signed, err = tr.sign(signed, timestamp)
	if err != nil {
		return nil, err
	}
	tr.Timestamp.Signatures = signed.Signatures
	tr.Snapshot.Dirty = false // snapshot is dirty until changes have been captured in timestamp
	return signed, nil
}

func (tr Repo) sign(signedData *data.Signed, role data.BaseRole) (*data.Signed, error) {
	ks := role.ListKeys()
	if len(ks) < 1 {
		return nil, signed.ErrNoKeys{}
	}
	err := signed.Sign(tr.cryptoService, signedData, ks...)
	if err != nil {
		return nil, err
	}
	return signedData, nil
}
