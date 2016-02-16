package data

import (
	"errors"
	"fmt"
	"path"

	"github.com/docker/go/canonical/json"
)

// SignedTargets is a fully unpacked targets.json, or target delegation
// json file
type SignedTargets struct {
	Signatures []Signature
	Signed     Targets
	Dirty      bool
}

// Targets is the Signed components of a targets.json or delegation json file
type Targets struct {
	SignedCommon
	Targets     Files       `json:"targets"`
	Delegations Delegations `json:"delegations,omitempty"`
}

// isValidTargetsStructure returns an error, or nil, depending on whether the content of the struct
// is valid for targets metadata.  This does not check signatures or expiry, just that
// the metadata content is valid.
func isValidTargetsStructure(t Targets, roleName string) error {
	if roleName != CanonicalTargetsRole && !IsDelegation(roleName) {
		return ErrInvalidRole{Role: roleName}
	}

	// even if it's a delegated role, the metadata type is "Targets"
	expectedType := TUFTypes[CanonicalTargetsRole]
	if t.Type != expectedType {
		return ErrInvalidMeta{
			Role: roleName, Msg: fmt.Sprintf("expected type %s, not %s", expectedType, t.Type)}
	}

	for _, roleObj := range t.Delegations.Roles {
		if !IsDelegation(roleObj.Name) || path.Dir(roleObj.Name) != roleName {
			return ErrInvalidMeta{
				Role: roleName, Msg: fmt.Sprintf("delegation role %s invalid", roleObj.Name)}
		}
		if roleObj.Threshold < 1 {
			return ErrInvalidMeta{
				Role: roleName, Msg: fmt.Sprintf("invalid threshold for %s: %v ", roleObj.Name, roleObj.Threshold)}
		}
		for _, keyID := range roleObj.KeyIDs {
			if _, ok := t.Delegations.Keys[keyID]; !ok {
				return ErrInvalidMeta{
					Role: roleName,
					Msg:  fmt.Sprintf("%s role specifies key ID %s without corresponding key", roleObj.Name, keyID),
				}
			}
		}
	}
	return nil
}

// NewTargets intiializes a new empty SignedTargets object
func NewTargets() *SignedTargets {
	return &SignedTargets{
		Signatures: make([]Signature, 0),
		Signed: Targets{
			SignedCommon: SignedCommon{
				Type:    TUFTypes["targets"],
				Version: 0,
				Expires: DefaultExpires("targets"),
			},
			Targets:     make(Files),
			Delegations: *NewDelegations(),
		},
		Dirty: true,
	}
}

// GetMeta attempts to find the targets entry for the path. It
// will return nil in the case of the target not being found.
func (t SignedTargets) GetMeta(path string) *FileMeta {
	for p, meta := range t.Signed.Targets {
		if p == path {
			return &meta
		}
	}
	return nil
}

// GetDelegations filters the roles and associated keys that may be
// the signers for the given target path. If no appropriate roles
// can be found, it will simply return nil for the return values.
// The returned slice of Role will have order maintained relative
// to the role slice on Delegations per TUF spec proposal on using
// order to determine priority.
func (t SignedTargets) GetDelegations(path string) []*Role {
	var roles []*Role
	for _, r := range t.Signed.Delegations.Roles {
		if r.CheckPaths(path) {
			roles = append(roles, r)
			continue
		}
	}
	return roles
}

// AddTarget adds or updates the meta for the given path
func (t *SignedTargets) AddTarget(path string, meta FileMeta) {
	t.Signed.Targets[path] = meta
	t.Dirty = true
}

// AddDelegation will add a new delegated role with the given keys,
// ensuring the keys either already exist, or are added to the map
// of delegation keys
func (t *SignedTargets) AddDelegation(role *Role, keys []*PublicKey) error {
	return errors.New("Not Implemented")
}

// ToSigned partially serializes a SignedTargets for further signing
func (t *SignedTargets) ToSigned() (*Signed, error) {
	s, err := defaultSerializer.MarshalCanonical(t.Signed)
	if err != nil {
		return nil, err
	}
	signed := json.RawMessage{}
	err = signed.UnmarshalJSON(s)
	if err != nil {
		return nil, err
	}
	sigs := make([]Signature, len(t.Signatures))
	copy(sigs, t.Signatures)
	return &Signed{
		Signatures: sigs,
		Signed:     signed,
	}, nil
}

// MarshalJSON returns the serialized form of SignedTargets as bytes
func (t *SignedTargets) MarshalJSON() ([]byte, error) {
	signed, err := t.ToSigned()
	if err != nil {
		return nil, err
	}
	return defaultSerializer.Marshal(signed)
}

// TargetsFromSigned fully unpacks a Signed object into a SignedTargets, given
// a role name (so it can validate the SignedTargets object)
func TargetsFromSigned(s *Signed, roleName string) (*SignedTargets, error) {
	t := Targets{}
	if err := defaultSerializer.Unmarshal(s.Signed, &t); err != nil {
		return nil, err
	}
	if err := isValidTargetsStructure(t, roleName); err != nil {
		return nil, err
	}
	sigs := make([]Signature, len(s.Signatures))
	copy(sigs, s.Signatures)
	return &SignedTargets{
		Signatures: sigs,
		Signed:     t,
	}, nil
}
