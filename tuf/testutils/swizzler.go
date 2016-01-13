package testutils

import (
	"bytes"
	"path"
	"time"

	"github.com/docker/notary/cryptoservice"
	"github.com/docker/notary/passphrase"
	"github.com/docker/notary/trustmanager"
	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/signed"
	"github.com/docker/notary/tuf/store"
	"github.com/jfrazelle/go/canonical/json"
)

const (
	maxSize = 5 << 20
)

// ErrNoKeyForRole returns an error when the cryptoservice provided to
// MetadataSwizzler has no key for a particular role
type ErrNoKeyForRole struct {
	Role string
}

func (e ErrNoKeyForRole) Error() string {
	return "Swizzler's cryptoservice has no key for role " + e.Role
}

// MetadataSwizzler fuzzes the metadata in a MetadataStore
type MetadataSwizzler struct {
	gun           string
	MetadataCache store.MetadataStore
	CryptoService signed.CryptoService
	Roles         []string // list of Roles in the metadataStore
}

// signs the new metadata, replacing whatever signature was there
func serializeMetadata(cs signed.CryptoService, s *data.Signed, role string,
	pubKeys ...data.PublicKey) ([]byte, error) {

	// delete the existing signatures
	s.Signatures = []data.Signature{}

	if len(pubKeys) > 0 {
		if err := signed.Sign(cs, s, pubKeys...); err != nil {
			if _, ok := err.(signed.ErrNoKeys); ok {
				return nil, ErrNoKeyForRole{Role: role}
			}
			return nil, err
		}
	} else if role == data.CanonicalRootRole {
		// if this is root metadata, we have to get the keys from the root because they
		// are certs
		root := &data.Root{}
		if err := json.Unmarshal(s.Signed, root); err != nil {
			return nil, err
		}
		for _, pubKeyID := range root.Roles[data.CanonicalRootRole].KeyIDs {
			if err := signed.Sign(cs, s, root.Keys[pubKeyID]); err != nil {
				if _, ok := err.(signed.ErrNoKeys); ok {
					return nil, ErrNoKeyForRole{Role: role}
				}
				return nil, err
			}
		}
	} else {
		pubKeyIDs := cs.ListKeys(role)
		if len(pubKeyIDs) < 1 {
			return nil, ErrNoKeyForRole{role}
		}
		for _, pubKeyID := range pubKeyIDs {
			if err := signed.Sign(cs, s, cs.GetKey(pubKeyID)); err != nil {
				return nil, err
			}
		}
	}

	metaBytes, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	return metaBytes, nil
}

// gets a Signed from the metadata store
func signedFromStore(cache store.MetadataStore, role string) (*data.Signed, error) {
	b, err := cache.GetMeta(role, maxSize)
	if err != nil {
		return nil, err
	}

	signed := &data.Signed{}
	if err := json.Unmarshal(b, signed); err != nil {
		return nil, err
	}

	return signed, nil
}

var delegatedRoles = []string{"targets/a", "targets/a/b"}

// DefaultRoles are the defualt roles that NewMetadataSwizzler creates
var DefaultRoles = append(data.BaseRoles, delegatedRoles...)

// NewMetadataSwizzler creates a new tuf.Repo and generates metadata to fuzz.
func NewMetadataSwizzler(gun string) (*MetadataSwizzler, error) {
	_, tufRepo, cs, err := EmptyRepo(gun)
	if err != nil {
		return nil, err
	}

	for _, delgName := range delegatedRoles {
		// create a delegations key and a delegation in the tuf repo
		delgKey, err := cs.Create(delgName, data.ED25519Key)
		if err != nil {
			return nil, err
		}
		role, err := data.NewRole(delgName, 1, []string{}, []string{"/"}, []string{})
		if err != nil {
			return nil, err
		}
		if err := tufRepo.UpdateDelegations(role, []data.PublicKey{delgKey}); err != nil {
			return nil, err
		}

		// create the targets metadata
		if _, ok := tufRepo.Targets[delgName]; !ok {
			_, err := tufRepo.InitTargets(delgName)
			if err != nil {
				return nil, err
			}
		}
	}

	meta := make(map[string][]byte)

	// now we need to create a signed target/serialize, add it to the snapshot,
	// and save the signed target metadata to the store - this must be done
	// in a separate loop because "targets/a" can't be serialized until "targets/a/b"
	// has been added
	for _, delgName := range delegatedRoles {
		signedThing, err := tufRepo.SignTargets(delgName, data.DefaultExpires("targets"))
		if err != nil {
			return nil, err
		}
		metaBytes, err := json.MarshalCanonical(signedThing)
		if err != nil {
			return nil, err
		}
		meta[delgName] = metaBytes
	}

	// these need to be generated after the delegations are created and signed so
	// the snapshot will have the delegation metadata
	rs, tgs, ss, ts, err := Sign(tufRepo)
	if err != nil {
		return nil, err
	}
	rf, tgf, sf, tf, err := Serialize(rs, tgs, ss, ts)
	if err != nil {
		return nil, err
	}

	meta[data.CanonicalRootRole] = rf
	meta[data.CanonicalSnapshotRole] = sf
	meta[data.CanonicalTargetsRole] = tgf
	meta[data.CanonicalTimestampRole] = tf

	swizzler := MetadataSwizzler{
		gun:           gun,
		MetadataCache: store.NewMemoryStore(meta, nil),
		CryptoService: cs,
		Roles:         DefaultRoles,
	}

	return &swizzler, nil
}

// SetInvalidJSON corrupts metadata into something that is no longer valid JSON
func (m *MetadataSwizzler) SetInvalidJSON(role string) error {
	metaBytes, err := m.MetadataCache.GetMeta(role, maxSize)
	if err != nil {
		return err
	}
	return m.MetadataCache.SetMeta(role, metaBytes[5:])
}

// SetInvalidSigned corrupts the metadata into something that is valid JSON,
// but not unmarshallable into signed JSON
func (m *MetadataSwizzler) SetInvalidSigned(role string) error {
	signedThing, err := signedFromStore(m.MetadataCache, role)
	if err != nil {
		return err
	}
	metaBytes, err := json.MarshalCanonical(map[string]interface{}{
		"signed":     signedThing.Signed,
		"signatures": "not list",
	})
	if err != nil {
		return err
	}
	return m.MetadataCache.SetMeta(role, metaBytes)
}

// SetInvalidSignedMeta corrupts the metadata into something that is unmarshallable
// as a Signed object, but not unmarshallable into a SignedMeta object
func (m *MetadataSwizzler) SetInvalidSignedMeta(role string) error {
	signedThing, err := signedFromStore(m.MetadataCache, role)
	if err != nil {
		return err
	}

	var unmarshalled map[string]interface{}
	if err := json.Unmarshal(signedThing.Signed, &unmarshalled); err != nil {
		return err
	}

	unmarshalled["_type"] = []string{"not a string"}
	unmarshalled["version"] = "string not int"
	unmarshalled["expires"] = "cannot be parsed as time"

	metaBytes, err := json.MarshalCanonical(unmarshalled)
	if err != nil {
		return err
	}
	signedThing.Signed = json.RawMessage(metaBytes)

	metaBytes, err = serializeMetadata(m.CryptoService, signedThing, role)
	if err != nil {
		return err
	}
	return m.MetadataCache.SetMeta(role, metaBytes)
}

// TODO: corrupt metadata in such a way that it can be unmarshalled as a
// SignedMeta, but not as a SignedRoot or SignedTarget, etc. (Signed*)

// SetInvalidMetadataType unmarshallable, but has the wrong metadata type (not
// actually a metadata type)
func (m *MetadataSwizzler) SetInvalidMetadataType(role string) error {
	signedThing, err := signedFromStore(m.MetadataCache, role)
	if err != nil {
		return err
	}

	var unmarshalled map[string]interface{}
	if err := json.Unmarshal(signedThing.Signed, &unmarshalled); err != nil {
		return err
	}

	unmarshalled["_type"] = "not_real"

	metaBytes, err := json.MarshalCanonical(unmarshalled)
	if err != nil {
		return err
	}
	signedThing.Signed = json.RawMessage(metaBytes)

	metaBytes, err = serializeMetadata(m.CryptoService, signedThing, role)
	if err != nil {
		return err
	}
	return m.MetadataCache.SetMeta(role, metaBytes)
}

// InvalidateMetadataSignatures signs with the right key(s) but wrong hash
func (m *MetadataSwizzler) InvalidateMetadataSignatures(role string) error {
	signedThing, err := signedFromStore(m.MetadataCache, role)
	if err != nil {
		return err
	}
	sigs := make([]data.Signature, len(signedThing.Signatures))
	for i, origSig := range signedThing.Signatures {
		sigs[i] = data.Signature{
			KeyID:     origSig.KeyID,
			Signature: []byte("invalid signature"),
			Method:    origSig.Method,
		}
	}
	signedThing.Signatures = sigs

	metaBytes, err := json.Marshal(signedThing)
	if err != nil {
		return err
	}
	return m.MetadataCache.SetMeta(role, metaBytes)
}

// TODO: AddExtraSignedInfo - add an extra field to Signed that doesn't get
//	unmarshalled, and the whole thing is correctly signed, so shouldn't cause
//  problems there.  Should this fail a canonical JSON check?

// RemoveMetadata deletes the metadata entirely
func (m *MetadataSwizzler) RemoveMetadata(role string) error {
	return m.MetadataCache.RemoveMeta(role)
}

// SignMetadataWithInvalidKey signs the metadata with the wrong key
func (m *MetadataSwizzler) SignMetadataWithInvalidKey(role string) error {
	signedThing, err := signedFromStore(m.MetadataCache, role)
	if err != nil {
		return err
	}

	// create an invalid key, but not in the existing CryptoService
	cs := cryptoservice.NewCryptoService(
		m.gun, trustmanager.NewKeyMemoryStore(passphrase.ConstantRetriever("")))
	key, err := createKey(cs, m.gun, role)
	if err != nil {
		return err
	}

	metaBytes, err := serializeMetadata(cs, signedThing, "root", key)
	if err != nil {
		return err
	}
	return m.MetadataCache.SetMeta(role, metaBytes)
}

// OffsetMetadataVersion updates the metadata version
func (m *MetadataSwizzler) OffsetMetadataVersion(role string, offset int) error {
	signedThing, err := signedFromStore(m.MetadataCache, role)
	if err != nil {
		return err
	}

	var unmarshalled map[string]interface{}
	if err := json.Unmarshal(signedThing.Signed, &unmarshalled); err != nil {
		return err
	}

	oldVersion, ok := unmarshalled["version"].(float64)
	if !ok {
		oldVersion = float64(0) // just ignore the error and set it to 0
	}
	unmarshalled["version"] = int(oldVersion) + offset

	metaBytes, err := json.MarshalCanonical(unmarshalled)
	if err != nil {
		return err
	}
	signedThing.Signed = json.RawMessage(metaBytes)

	metaBytes, err = serializeMetadata(m.CryptoService, signedThing, role)
	if err != nil {
		return err
	}
	return m.MetadataCache.SetMeta(role, metaBytes)
}

// ExpireMetadata expires the metadata, which would make it invalid - don't do anything if
// we don't have the timestamp key
func (m *MetadataSwizzler) ExpireMetadata(role string) error {
	signedThing, err := signedFromStore(m.MetadataCache, role)
	if err != nil {
		return err
	}

	var unmarshalled map[string]interface{}
	if err := json.Unmarshal(signedThing.Signed, &unmarshalled); err != nil {
		return err
	}

	unmarshalled["expires"] = time.Now().AddDate(-1, -1, -1)

	metaBytes, err := json.MarshalCanonical(unmarshalled)
	if err != nil {
		return err
	}
	signedThing.Signed = json.RawMessage(metaBytes)

	metaBytes, err = serializeMetadata(m.CryptoService, signedThing, role)
	if err != nil {
		return err
	}
	return m.MetadataCache.SetMeta(role, metaBytes)
}

// SetThreshold sets a threshold for a metadata role - can invalidate metadata for which
// the threshold is increased, if there aren't enough signatures or can be invalid because
// the threshold is 0
func (m *MetadataSwizzler) SetThreshold(role string, newThreshold int) error {
	roleSpecifier := data.CanonicalRootRole
	if data.IsDelegation(role) {
		roleSpecifier = path.Dir(role)
	}

	b, err := m.MetadataCache.GetMeta(roleSpecifier, maxSize)
	if err != nil {
		return err
	}

	signedThing := &data.Signed{}
	if err := json.Unmarshal(b, signedThing); err != nil {
		return err
	}

	if roleSpecifier == data.CanonicalRootRole {
		signedRoot, err := data.RootFromSigned(signedThing)
		if err != nil {
			return err
		}
		signedRoot.Signed.Roles[role].Threshold = newThreshold
		if signedThing, err = signedRoot.ToSigned(); err != nil {
			return err
		}
	} else {
		signedTargets, err := data.TargetsFromSigned(signedThing)
		if err != nil {
			return err
		}
		for _, roleObject := range signedTargets.Signed.Delegations.Roles {
			if roleObject.Name == role {
				roleObject.Threshold = newThreshold
				break
			}
		}
		if signedThing, err = signedTargets.ToSigned(); err != nil {
			return err
		}
	}

	metaBytes, err := serializeMetadata(m.CryptoService, signedThing, roleSpecifier)
	if err != nil {
		return err
	}
	return m.MetadataCache.SetMeta(roleSpecifier, metaBytes)
}

// ChangeRootKey swaps out the root key with a new key, and re-signs the metadata
// with the new key
func (m *MetadataSwizzler) ChangeRootKey() error {
	key, err := createKey(m.CryptoService, m.gun, data.CanonicalRootRole)
	if err != nil {
		return err
	}

	b, err := m.MetadataCache.GetMeta(data.CanonicalRootRole, maxSize)
	if err != nil {
		return err
	}

	signedRoot := &data.SignedRoot{}
	if err := json.Unmarshal(b, signedRoot); err != nil {
		return err
	}

	signedRoot.Signed.Keys[key.ID()] = key
	signedRoot.Signed.Roles[data.CanonicalRootRole].KeyIDs = []string{key.ID()}

	var signedThing *data.Signed
	if signedThing, err = signedRoot.ToSigned(); err != nil {
		return err
	}

	metaBytes, err := serializeMetadata(m.CryptoService, signedThing, data.CanonicalRootRole)
	if err != nil {
		return err
	}
	return m.MetadataCache.SetMeta(data.CanonicalRootRole, metaBytes)
}

// UpdateSnapshotHashes updates the snapshot to reflect the latest hash changes, to
// ensure that failure isn't because the snapshot has the wrong hash.
func (m *MetadataSwizzler) UpdateSnapshotHashes(roles ...string) error {
	var (
		metaBytes      []byte
		snapshotSigned *data.Signed
		err            error
	)
	if metaBytes, err = m.MetadataCache.GetMeta(data.CanonicalSnapshotRole, maxSize); err != nil {
		return err
	}

	snapshot := data.SignedSnapshot{}
	if err = json.Unmarshal(metaBytes, &snapshot); err != nil {
		return err
	}

	// just rebuild everything if roles is not specified
	if len(roles) == 0 {
		roles = m.Roles
	}

	for _, role := range roles {
		if role != data.CanonicalSnapshotRole && role != data.CanonicalTimestampRole {
			if metaBytes, err = m.MetadataCache.GetMeta(role, maxSize); err != nil {
				return err
			}

			meta, err := data.NewFileMeta(bytes.NewReader(metaBytes), "sha256")
			if err != nil {
				return err
			}

			snapshot.Signed.Meta[role] = meta
		}
	}

	if snapshotSigned, err = snapshot.ToSigned(); err != nil {
		return err
	}
	metaBytes, err = serializeMetadata(m.CryptoService, snapshotSigned, data.CanonicalSnapshotRole)
	if err != nil {
		return err
	}
	return m.MetadataCache.SetMeta(data.CanonicalSnapshotRole, metaBytes)
}

// UpdateTimestampHash updates the timestamp to reflect the latest snapshot changes, to
// ensure that failure isn't because the timestamp has the wrong hash.
func (m *MetadataSwizzler) UpdateTimestampHash() error {
	var (
		metaBytes       []byte
		timestamp       = &data.SignedTimestamp{}
		timestampSigned *data.Signed
		err             error
	)
	if metaBytes, err = m.MetadataCache.GetMeta(data.CanonicalTimestampRole, maxSize); err != nil {
		return err
	}
	// we can't just create a new timestamp, because then the expiry would be
	// different
	if err = json.Unmarshal(metaBytes, timestamp); err != nil {
		return err
	}

	if metaBytes, err = m.MetadataCache.GetMeta(data.CanonicalSnapshotRole, maxSize); err != nil {
		return err
	}

	snapshotMeta, err := data.NewFileMeta(bytes.NewReader(metaBytes), "sha256")
	if err != nil {
		return err
	}

	timestamp.Signed.Meta[data.CanonicalSnapshotRole] = snapshotMeta

	timestampSigned, err = timestamp.ToSigned()
	if err != nil {
		return err
	}
	metaBytes, err = serializeMetadata(m.CryptoService, timestampSigned, data.CanonicalTimestampRole)
	if err != nil {
		return err
	}
	return m.MetadataCache.SetMeta(data.CanonicalTimestampRole, metaBytes)
}
