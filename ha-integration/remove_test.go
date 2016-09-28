package ha_integration

import (
	"github.com/docker/dhe-deploy/ha-integration/ha_utils"
	"github.com/docker/dhe-deploy/ha-integration/util"
	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
)

func removeTests() {
	Context("install tests", func() {
		It("removes from healthy clusters", func() {
			testHealthyRemoves(machines)
		})
		It("force removes always", func() {
			testForceRemoves(machines)
		})
	})
}

// TODO: convert to ginkgo style tests
func testHealthyRemoves(machines []ha_utils.Machine) {
	// TODO: get rid of the t
	t := GinkgoT()
	installArgs := util.InstallArgs{}
	installMachine := util.GetAvailableMachine(machines)
	assert.NotNil(t, installMachine)

	imageArgs := util.DefaultInstallation.Args.DTRImageArgs

	util.InstallDTRNode(installMachine, installArgs)
	util.FullSmokeTest(machines, imageArgs, t)

	// We're not sure if we have a second one in the pool, so we'll not do all the tests based on availability
	joinMachine := util.GetAvailableMachine(machines)

	// Join/Remove tests from a 2 replica cluster
	if joinMachine != nil {
		joinArgs := util.JoinArgs{}

		// First a regular remove on a joined node
		util.JoinDTRNode(joinMachine, joinArgs)
		util.FullSmokeTest(machines, imageArgs, t)

		removeArgs := util.RemoveArgs{}
		removeArgs.ExistingReplicaID = util.DefaultInstallation.Replicas[installMachine.GetName()]
		removeArgs.ReplicaID = util.DefaultInstallation.Replicas[joinMachine.GetName()]

		util.RemoveDTRNode(joinMachine, removeArgs)

		assert.Equal(t, len(util.DefaultInstallation.Replicas), 1)
		util.FullSmokeTest(machines, imageArgs, t)

		// Now we will join a node again and remove the first node we installed
		util.JoinDTRNode(joinMachine, joinArgs)
		util.FullSmokeTest(machines, imageArgs, t)

		removeArgs.ExistingReplicaID = util.DefaultInstallation.Replicas[joinMachine.GetName()]
		removeArgs.ReplicaID = util.DefaultInstallation.Replicas[installMachine.GetName()]

		util.RemoveDTRNode(installMachine, removeArgs)

		// Confirm that remove succeeded
		assert.Equal(t, len(util.DefaultInstallation.Replicas), 1)
		util.FullSmokeTest(machines, imageArgs, t)

		removeArgs.ExistingReplicaID = util.DefaultInstallation.Replicas[joinMachine.GetName()]
		removeArgs.ReplicaID = util.DefaultInstallation.Replicas[joinMachine.GetName()]

		util.RemoveDTRNode(joinMachine, removeArgs)

		// Assert that remove succeeded
		assert.Equal(t, len(util.DefaultInstallation.Replicas), 0)
		util.FullSmokeTest(machines, imageArgs, t)
	} else {
		removeArgs := util.RemoveArgs{}
		removeArgs.ExistingReplicaID = util.DefaultInstallation.Replicas[installMachine.GetName()]
		removeArgs.ReplicaID = util.DefaultInstallation.Replicas[installMachine.GetName()]

		util.RemoveDTRNode(installMachine, removeArgs)

		// Assert that remove succeeded
		assert.Equal(t, len(util.DefaultInstallation.Replicas), 0)
		util.FullSmokeTest(machines, imageArgs, t)
	}
}

// TODO: convert to ginkgo style tests
func testForceRemoves(machines []ha_utils.Machine) {
	// TODO: get rid of the t
	t := GinkgoT()
	installArgs := util.InstallArgs{}
	installMachine := util.GetAvailableMachine(machines)
	assert.NotNil(t, installMachine)

	imageArgs := util.DefaultInstallation.Args.DTRImageArgs

	util.InstallDTRNode(installMachine, installArgs)
	util.FullSmokeTest(machines, imageArgs, t)

	// We're not sure how many nodes we have
	firstJoinMachine := util.GetAvailableMachine(machines)
	if firstJoinMachine != nil {
		joinArgs := util.JoinArgs{}
		removeArgs := util.RemoveArgs{}

		util.JoinDTRNode(firstJoinMachine, joinArgs)
		util.FullSmokeTest(machines, imageArgs, t)

		secondJoinMachine := util.GetAvailableMachine(machines)
		if secondJoinMachine != nil {

			util.JoinDTRNode(secondJoinMachine, joinArgs)
			util.FullSmokeTest(machines, imageArgs, t)

			removeArgs.ReplicaID = util.DefaultInstallation.Replicas[secondJoinMachine.GetName()]
			removeArgs.ExistingReplicaID = util.DefaultInstallation.Replicas[secondJoinMachine.GetName()]
			removeArgs.Force = true

			util.RemoveDTRNode(secondJoinMachine, removeArgs)

			// Assert that remove failed
			assert.Equal(t, len(util.DefaultInstallation.Replicas), 2)
			// Confirm we didn't break the cluster
			util.FullSmokeTest(machines, imageArgs, t)
		}
	}
}
