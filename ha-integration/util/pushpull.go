package util

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/blang/semver"
	"github.com/docker/dhe-deploy/ha-integration/ha_utils"
	integrationutil "github.com/docker/dhe-deploy/integration/util"
	"github.com/docker/dhe-deploy/manager/versions"
	"github.com/docker/dhe-deploy/shared/dtrutil"
	"github.com/docker/engine-api/client"
	"github.com/stretchr/testify/require"
)

type pushPullTestInfo struct {
	trustSupported bool
	imageName      string
	accountName    string
	password       string
	expectedTags   []string
}

func pushPullFromMachine(newTag string, binDocker *integrationutil.DockerTrustClient, testInfo pushPullTestInfo, t require.TestingT) {
	_, err := binDocker.Login(testInfo.accountName, testInfo.password)
	require.NoError(t, err)

	for _, tag := range testInfo.expectedTags {
		_, err := binDocker.Pull(testInfo.imageName+":"+tag, testInfo.trustSupported)
		require.NoError(t, err)
	}

	// get an image we can tag and push
	var tagImage string
	if len(testInfo.expectedTags) == 0 {
		tagImage = "tianon/true"
		_, err = binDocker.Pull(tagImage, false)
		require.NoError(t, err)
	} else {
		tagImage = testInfo.imageName + ":" + testInfo.expectedTags[len(testInfo.expectedTags)-1]
	}

	// tag and push the new image
	newImage := testInfo.imageName + ":" + newTag
	_, err = binDocker.Tag(tagImage, newImage)
	require.NoError(t, err)

	_, err = binDocker.Push(newImage, testInfo.trustSupported)
	require.NoError(t, err)

	_, err = binDocker.Pull(newImage, testInfo.trustSupported)
	require.NoError(t, err)

	logrus.Debugf("Successfully push/pulled %s", newImage)
}

func reconfigureDTRLB(machines []ha_utils.Machine, client *client.Client, t require.TestingT) {
	// reconfigure the load balancer to point to a single node - this should tear down the prevoius load balancer.
	// We also set the constraint that it should be on the same node as the previous DTR load balancer, otherwise
	// the DTR host, etc. will all have to change
	dtrLBNode := SetupDefaultLoadBalancer(
		machines, client, []string{fmt.Sprintf("constraint:node==%s", DefaultInstallation.DTRLBNode)},
		DTRLoadBalancerPort, DefaultDTRNodePort, DTRLoadBalancerContainerName)

	// Sanity check to make sure the load balancer will be installed on the same node.
	// If it isn't then this entire test is faulty, so fail.
	dtrIP, err := ha_utils.GetIP(dtrLBNode)
	require.NoError(t, err)
	require.Equal(t, DefaultInstallation.Args.DTRHost, fmt.Sprintf("%s:%d", dtrIP, DTRLoadBalancerPort))

	// wait until the LB is up
	attempts := 0
	err = dtrutil.Poll(time.Second, 300, func() error {
		attempts++
		c := MakeAPIClient(DefaultInstallation.Args.DTRHost)
		return c.Login(ha_utils.GetAdminUser(), ha_utils.GetAdminPassword())
	})
	require.NoError(t, err)

	var names []string
	for _, machine := range machines {
		names = append(names, machine.GetName())
	}
	logrus.Infof("push/pull test: LB pointing only %s came up after %d attempts", strings.Join(names, ", "), attempts)
}

// Reconfigures the load balancer to push/pull to each individual node
func PushPullTest(machines []ha_utils.Machine, imageArgs DTRImageArgs, t require.TestingT) {
	if len(DefaultInstallation.Replicas) == 0 {
		return
	}

	// TODO: once the image store is fixed for ha tests with replicas > 1, remove this part
	if len(DefaultInstallation.Replicas) > 1 {
		logrus.Warn("Cannot do the push/pull test in HA mode until we fix running HA tests with filesystem driver")
		return
	}

	currVersion, err := versions.TagToSemver(imageArgs.DTRTag)
	require.NoError(t, err)
	trustSupported := currVersion.GTE(semver.Version{Major: 2, Minor: 1, Patch: 0, Pre: []semver.PRVersion{{VersionStr: "alpha"}}})

	// we won't need the hub username and password to pull down any public images
	ucpClient, err := GetDockerClient(
		DefaultInstallation.Args.UCPHost, DefaultInstallation.Args.UCPArgs.AdminUsername,
		DefaultInstallation.Args.UCPArgs.AdminPassword, "", "")
	require.NoError(t, err)

	ca, err := DefaultInstallation.API.GetCA()
	require.NoError(t, err)

	repoName := integrationutil.RandStringBytes(20)

	// login and create the repo in DTR, else we won't be able to get a garant token for the repo
	require.NoError(t, DefaultInstallation.API.Login(ha_utils.GetAdminUser(), ha_utils.GetAdminPassword()))
	_, err = DefaultInstallation.API.CreateRepository(ha_utils.GetAdminUser(), repoName, "meh", "meh", "public")
	require.NoError(t, err)
	defer DefaultInstallation.API.DeleteRepository(ha_utils.GetAdminUser(), repoName)

	repos, err := DefaultInstallation.API.ListRepositories(ha_utils.GetAdminUser())
	require.NoError(t, err)
	var found bool
	for _, repo := range repos {
		if repo.Name == repoName {
			found = true
		}
	}
	require.True(t, found, "Repo %s was never created", repoName)

	binDocker, err := integrationutil.NewDockerClientWithTrust(ucpClient, DefaultInstallation.Args.DTRHost, "docker:dind", ca)
	require.NoError(t, err)
	defer binDocker.Cleanup()

	// push from every machine, then assert that pull gets the data from the
	// latest pull and all the other pushes on other machines
	testInfo := pushPullTestInfo{
		trustSupported: trustSupported,
		imageName:      filepath.Join(DefaultInstallation.Args.DTRHost, ha_utils.GetAdminUser(), repoName),
		accountName:    ha_utils.GetAdminUser(),
		password:       ha_utils.GetAdminPassword(),
	}

	for _, machine := range machines {
		if _, ok := DefaultInstallation.Replicas[machine.GetName()]; ok {
			reconfigureDTRLB([]ha_utils.Machine{machine}, ucpClient, t)
			newTag := machine.GetName()
			pushPullFromMachine(newTag, binDocker, testInfo, t)
			testInfo.expectedTags = append(testInfo.expectedTags, newTag)
		}
	}

	// put the LB back the way it was, and push/pull again
	reconfigureDTRLB(machines, ucpClient, t)
	pushPullFromMachine("originalAgain", binDocker, testInfo, t)
}
