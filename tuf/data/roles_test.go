package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanonicalRole(t *testing.T) {

	testRoles := map[string]string{
		CanonicalRootRole:      "testRoot",
		CanonicalTargetsRole:   "testTargets",
		CanonicalSnapshotRole:  "testSnapshot",
		CanonicalTimestampRole: "testTimestamp",
		"garbageRole":          "testGarbageRole",
	}

	SetValidRoles(testRoles)

	// make sure roles were set correctly
	assert.Equal(t, "testRoot", ValidRoles[CanonicalRootRole])
	assert.Equal(t, "testTargets", ValidRoles[CanonicalTargetsRole])
	assert.Equal(t, "testSnapshot", ValidRoles[CanonicalSnapshotRole])
	assert.Equal(t, "testTimestamp", ValidRoles[CanonicalTimestampRole])
	// check SetValidRoles doesn't allow non-valid roles in
	assert.Equal(t, "", ValidRoles["garbageRole"])

	// check when looking up CanonicalRole from configured role
	assert.Equal(t, CanonicalRootRole, CanonicalRole("testRoot"))
	assert.Equal(t, CanonicalTargetsRole, CanonicalRole("testTargets"))
	assert.Equal(t, CanonicalSnapshotRole, CanonicalRole("testSnapshot"))
	assert.Equal(t, CanonicalTimestampRole, CanonicalRole("testTimestamp"))
	assert.Equal(t, "", CanonicalRole("testGarbageRole"))

	// check when looking up CanonicalRole with canonical role
	assert.Equal(t, CanonicalRootRole, CanonicalRole(CanonicalRootRole))
	assert.Equal(t, CanonicalTargetsRole, CanonicalRole(CanonicalTargetsRole))
	assert.Equal(t, CanonicalSnapshotRole, CanonicalRole(CanonicalSnapshotRole))
	assert.Equal(t, CanonicalTimestampRole, CanonicalRole(CanonicalTimestampRole))
	assert.Equal(t, "", CanonicalRole("garbageRole"))

	assert.Equal(t, "", CanonicalRole("not found"))

	// reset ValidRoles so other tests aren't messed up
	ValidRoles = map[string]string{
		CanonicalRootRole:      CanonicalRootRole,
		CanonicalTargetsRole:   CanonicalTargetsRole,
		CanonicalSnapshotRole:  CanonicalSnapshotRole,
		CanonicalTimestampRole: CanonicalTimestampRole,
	}
}
