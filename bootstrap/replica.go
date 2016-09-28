package bootstrap

import (
	"fmt"

	"github.com/docker/dhe-deploy/manager/versions"
)

// Replica represents the metadata for a DTR replica to differentiate
// from the versioning of the bootstrap container being ran.
type Replica struct {
	ReplicaID string
	Version   string
}

func NewReplica(replicaID string, version string) *Replica {
	return &Replica{
		replicaID,
		version,
	}
}

// CheckEqualVersion returns an error if version or deploy.Version cannot
// be parsed into a semantic version or if its semantic versions of
// the image tags are not equal between the replica and bootstrap container
func (r *Replica) CheckEqualVersion(version string) error {
	replicaVersion, err := versions.TagToSemver(r.Version)
	if err != nil {
		return err
	}

	semversion, err := versions.TagToSemver(version)
	if err != nil {
		return err
	}

	if !replicaVersion.Equals(semversion) {
		return fmt.Errorf("Failed version match. Replica version %s doesn't match bootstrap version %s.", r.Version, version)
	}
	return nil
}

// CheckUpgradeableVersion returns an error if version or deploy.Version cannot
// be parsed into a semantic version of if its semantic versions of
// the image tags differ by a major version or the replica has a higher
// minor version than the bootstrap container or the bootstrap container
// is more than one minor version ahead of the replica
func (r *Replica) CheckUpgradeableVersion(version string) error {
	replicaVersion, err := versions.TagToSemver(r.Version)
	if err != nil {
		return err
	}

	semversion, err := versions.TagToSemver(version)
	if err != nil {
		return err
	}

	if replicaVersion.Major != semversion.Major || replicaVersion.Minor > semversion.Minor || replicaVersion.Minor+1 < semversion.Minor {
		return fmt.Errorf("Couldn't upgrade from %s to %s. Bootstrap only supports upgrading from versions: %d.%d.x - %d.%d.x", r.Version, version, semversion.Major, semversion.Minor-1, semversion.Major, semversion.Minor)
	}
	return nil
}

// RequiresMigration returns an error if version or deploy.Version cannot
// be parsed into a semantic version or if it isn't an upgradeable version.
// It also returns a boolean that indicates whether the upgrade should run
// schema and data migrations due to a minor version bump.
func (r *Replica) RequiresMigration(version string) (bool, error) {
	if err := r.CheckUpgradeableVersion(version); err != nil {
		return false, err
	}

	replicaVersion, err := versions.TagToSemver(r.Version)
	if err != nil {
		return false, err
	}

	semversion, err := versions.TagToSemver(version)
	if err != nil {
		return false, err
	}

	return semversion.Minor > replicaVersion.Minor, nil
}

type Replicas []*Replica

// ReplicaIDs returns a list of replica IDs as an array of strings
func (r Replicas) ReplicaIDs() []string {
	var replicaIDs []string
	for _, replica := range r {
		replicaIDs = append(replicaIDs, replica.ReplicaID)
	}
	return replicaIDs
}

// GetReplica returns the Replica instance referenced by its replicaID
func (r Replicas) GetReplica(replicaID string) *Replica {
	for _, replica := range r {
		if replica.ReplicaID == replicaID {
			return replica
		}
	}
	return nil
}
