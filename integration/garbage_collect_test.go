package integration

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/docker/dhe-deploy/adminserver/api/common/forms"
	"github.com/docker/dhe-deploy/integration/framework"
	"github.com/docker/dhe-deploy/integration/util"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type GCTestSuite struct {
	suite.Suite
	*framework.IntegrationFramework
	u *util.Util

	restoreSettings func()
}

// SetupSuite sets up the test suite

func (suite *GCTestSuite) SetupSuite() {
	suite.IntegrationFramework, suite.u = setupFramework(suite)

	suite.restoreSettings = suite.u.SwitchAuth()

	// we use the admin user for everything
	suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)

	// ensure that migration has completed before we embark upon our tests;
	// this is a prequisite as GC doesn't work without tagmigration completed
	// (it requires tagstore data)
	job, err := suite.API.RunJobByAction("tagmigration")
	require.NotEmpty(suite.T(), job.ID, fmt.Sprintf("tagmiration not created: %#v", job))
	require.Nil(suite.T(), err)
	err = suite.u.WaitForJob(job.ID)
	require.Nil(suite.T(), err, "Failed to get complete tagmigration task in time: %v", err)
}

func (suite *GCTestSuite) TearDownSuite() {
	defer suite.restoreSettings()
	if err := suite.u.PollAvailable(); err != nil {
		suite.T().Fatalf("DTR failed to return to normal after GC tests")
	}
}

func (suite *GCTestSuite) SetupTest() {
	util.WipeDTRIgnorableLoggedErrors()
	util.WipeDockerIgnorableLoggedErrors()
}

func (suite *GCTestSuite) TearDownTest() {
	suite.u.TestLogs()
}

func (suite *GCTestSuite) GetDiskUsage() int64 {
	// we look at the storage of all replicas on the machine
	stdout, stderr, err := suite.SSH.RunRemoteCommand("sudo docker volume ls | grep dtr-registry | awk '{print $2}' | xargs -I+ sudo docker run --rm -v +:/storage busybox du -s /storage | awk '{print $1}'")
	if err != nil {
		suite.T().Fatalf("Failed to run df command: %v, stderr: %s", err, stderr)
	}
	cleaned := strings.TrimSpace(stdout)
	lines := strings.Split(cleaned, "\n")
	totalSize := int64(0)
	for _, line := range lines {
		size, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			suite.T().Fatalf("Failed to parse size: %v, stderr: %s", err, stderr)
		}
		totalSize += size
	}
	return totalSize
}

// Test that GC doesn't delete things that should not be deleted
func (suite *GCTestSuite) TestRunWithNoDeletionGC() {
	beforeSize := suite.GetDiskUsage()
	suite.runGC()
	// either there is the same amount of disk usage or more (due to logging etc)
	afterGCSize := suite.GetDiskUsage()
	if afterGCSize < beforeSize {
		suite.T().Fatalf("After \"GC with nothing to collect\" space used not greater or equal. before: %v after: %v", beforeSize, afterGCSize)
	}
	util.AppendDockerIgnorableLoggedErrors([]string{
		"No such container",
	})
}

// Test that GC deletes manifests that aren't referenced by any tags when
// run in MarkByTag mode
func (suite *GCTestSuite) TestRunGCInMarkByTagMode() {
	// set the GC mode to MarkByTag
	suite.setGCMode("tag")

	repoName := "mytrue"
	tag := "rick_sanchez_tag" + strconv.FormatInt(util.Prng.Int63(), 16)
	beforeSize := suite.GetDiskUsage()

	defer suite.u.LoadPackedImage()()
	defer suite.u.CreateRepoWithChecks(suite.Config.AdminUsername, repoName, "test", "long test", "public")()
	defer suite.u.TagImageWithChecks(suite.Config.AdminNamespace(), repoName, tag, suite.u.PackedImageName)()

	func() {
		// this will delete the image tag, but not the manifest, from the
		// tagstore and blobstore
		defer suite.u.PushImageWithChecks(suite.Config.AdminNamespace(), repoName, tag)()

		afterPushSize := suite.GetDiskUsage()
		if beforeSize >= afterPushSize {
			suite.T().Fatalf("After push size not larger: %v >= %v", beforeSize, afterPushSize)
		}
	}()

	// at this point the tag is deleted but the manifest should still exist
	mfsts, err := suite.API.GetRepositoryManifests(suite.Config.AdminUsername, repoName)
	require.Nil(suite.T(), err, "%s", err)
	require.Equal(suite.T(), 1, len(mfsts))

	afterTagDeleteSize := suite.GetDiskUsage()
	suite.runGC()
	afterGCSize := suite.GetDiskUsage()
	if afterGCSize >= afterTagDeleteSize {
		suite.T().Fatalf("After GC size not smaller: %v >= %v.", afterGCSize, afterTagDeleteSize)
	}

	// Ensure that running in tag mode we delete the untagged maniefst
	mfsts, err = suite.API.GetRepositoryManifests(suite.Config.AdminUsername, repoName)
	require.Nil(suite.T(), err, "%s", err)
	require.Equal(suite.T(), 0, len(mfsts))

	util.AppendDockerIgnorableLoggedErrors([]string{
		"Error streaming logs: unexpected EOF",
		"No such image",
		"No such container",
	})
}

func (suite *GCTestSuite) TestRunGCInMarkByManifestMode() {
	suite.setGCMode("manifest")

	repoName := "mytrue"
	tag := "rick_sanchez_manifest" + strconv.FormatInt(util.Prng.Int63(), 16)
	beforeSize := suite.GetDiskUsage()

	defer suite.u.LoadPackedImage()()
	defer suite.u.CreateRepoWithChecks(suite.Config.AdminUsername, repoName, "test", "long test", "public")()
	defer suite.u.TagImageWithChecks(suite.Config.AdminNamespace(), repoName, tag, suite.u.PackedImageName)()

	// ensure that deleting in mark by manifest doesn't remove any blobs
	// if we only delete the tag.
	func() {
		func() {
			// this will delete the image tag, but not the manifest, from the
			// tagstore and blobstore
			defer suite.u.PushImageWithChecks(suite.Config.AdminNamespace(), repoName, tag)()
			afterPushSize := suite.GetDiskUsage()
			if beforeSize >= afterPushSize {
				suite.T().Fatalf("After push size not larger: %v >= %v", beforeSize, afterPushSize)
			}
		}()

		// at this point the tag is deleted but the manifest should still exist
		mfsts, _ := suite.API.GetRepositoryManifests(suite.Config.AdminUsername, repoName)
		require.Equal(suite.T(), 1, len(mfsts))

		afterTagDeleteSize := suite.GetDiskUsage()
		suite.runGC()
		afterGCSize := suite.GetDiskUsage()
		if afterGCSize != afterTagDeleteSize {
			suite.T().Fatalf("gc in manifest mode should have same size with no deleted manifests: %v != %v.", afterGCSize, afterTagDeleteSize)
		}

		// make sure the manifest still exists
		mfsts, _ = suite.API.GetRepositoryManifests(suite.Config.AdminUsername, repoName)
		require.Equal(suite.T(), 1, len(mfsts))
	}()

	// ensure that deleting in mark by manifest cleans up if we remove manifests
	func() {
		func() {
			// this will delete the image tag, but not the manifest, from the
			// tagstore and blobstore
			defer suite.u.PushImageWithChecks(suite.Config.AdminNamespace(), repoName, tag)()
			afterPushSize := suite.GetDiskUsage()
			if beforeSize >= afterPushSize {
				suite.T().Fatalf("After push size not larger: %v >= %v", beforeSize, afterPushSize)
			}

			// Delete the manifest
			mfsts, _ := suite.API.GetRepositoryManifests(suite.Config.AdminUsername, repoName)
			require.Equal(suite.T(), 1, len(mfsts))
			suite.API.DeleteManifest(suite.Config.AdminUsername, repoName, mfsts[0].Digest)
			mfsts, _ = suite.API.GetRepositoryManifests(suite.Config.AdminUsername, repoName)
			require.Equal(suite.T(), 0, len(mfsts))
		}()

		afterDeleteSize := suite.GetDiskUsage()
		suite.runGC()
		afterGCSize := suite.GetDiskUsage()
		if afterGCSize >= afterDeleteSize {
			suite.T().Fatalf("after gc in manifest mode not smaller: %v >= %v.", afterGCSize, afterDeleteSize)
		}
	}()
}

// runGC is a helper function to run GC
func (suite *GCTestSuite) runGC() {
	job, err := suite.API.RunJobByAction("gc")
	require.NotEmpty(suite.T(), job.ID, fmt.Sprintf("job not created: %#v", job))
	require.Nil(suite.T(), err)
	err = suite.u.WaitForJob(job.ID)

	require.Nil(suite.T(), err, "Failed to get complete gc task in time: %v", err)
}

func (suite *GCTestSuite) setGCMode(mode string) {
	// set the GC mode to MarkByTag
	err := suite.API.SetHTTPSettings(&forms.Settings{
		GCMode: &mode,
	})
	require.Nil(suite.T(), err, "%s", err)
}

func TestGCSuite(t *testing.T) {
	suite.Run(t, new(GCTestSuite))
}
