package storage

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/docker/go/canonical/json"
	"github.com/docker/notary"
	"github.com/docker/notary/tuf/data"
)

// ConsistentMetaStorage wraps a MetaStore in order to walk the TUF tree for GetCurrent in a consistent manner
type ConsistentMetaStorage struct {
	MetaStore
}

// GetCurrent gets a specific TUF record, by walking from the current Timestamp to other metadata by checksum
func (cms *ConsistentMetaStorage) GetCurrent(gun, tufRole string) (*time.Time, []byte, error) {
	timestampTime, timestampJSON, err := cms.MetaStore.GetCurrent(gun, data.CanonicalTimestampRole)
	if err != nil {
		return nil, nil, err
	}
	// If we wanted data for the timestamp role, we're done here
	if tufRole == data.CanonicalTimestampRole {
		return timestampTime, timestampJSON, nil
	}

	// If we want to lookup another role, walk to it via current timestamp --> snapshot by checksum --> desired role
	timestampMeta := &data.SignedTimestamp{}
	if err := json.Unmarshal(timestampJSON, timestampMeta); err != nil {
		return nil, nil, fmt.Errorf("could not parse current timestamp")
	}
	snapshotChecksums, err := timestampMeta.GetSnapshot()
	if err != nil || snapshotChecksums == nil {
		return nil, nil, fmt.Errorf("could not retrieve latest snapshot checksum")
	}
	snapshotSha256Bytes, ok := snapshotChecksums.Hashes[notary.SHA256]
	if !ok {
		return nil, nil, fmt.Errorf("could not retrieve latest snapshot sha256")
	}
	snapshotSha256Hex := hex.EncodeToString(snapshotSha256Bytes[:])

	// Get the snapshot by checksum
	snapshotTime, snapshotJSON, err := cms.GetChecksum(gun, data.CanonicalSnapshotRole, snapshotSha256Hex)
	if err != nil {
		return nil, nil, err
	}
	// If we wanted data for the snapshot role, we're done here
	if tufRole == data.CanonicalSnapshotRole {
		return snapshotTime, snapshotJSON, nil
	}

	// If it's a different role, we should have the checksum in snapshot metadata, and we can use it to GetChecksum()
	snapshotMeta := &data.SignedSnapshot{}
	if err := json.Unmarshal(snapshotJSON, snapshotMeta); err != nil {
		return nil, nil, fmt.Errorf("could not parse current snapshot")
	}
	roleMeta, err := snapshotMeta.GetMeta(tufRole)
	if err != nil {
		return nil, nil, err
	}
	roleSha256Bytes, ok := roleMeta.Hashes[notary.SHA256]
	if !ok {
		return nil, nil, fmt.Errorf("could not retrieve latest %s sha256", tufRole)
	}
	roleSha256Hex := hex.EncodeToString(roleSha256Bytes[:])
	return cms.MetaStore.GetChecksum(gun, tufRole, roleSha256Hex)
}
