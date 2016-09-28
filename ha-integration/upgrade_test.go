package ha_integration

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/dhe-deploy/ha-integration/ha_utils"
	"github.com/docker/dhe-deploy/ha-integration/util"
	. "github.com/onsi/ginkgo"
)

func upgradeTests() {
	Context("upgrade tests", func() {
		It("upgrades regularly", func() {
			testUpgrade(machines)
		})
	})
}

// TODO: convert to ginkgo style tests
func testUpgrade(machines []ha_utils.Machine) {
	// TODO: get rid of the t
	t := GinkgoT()
	upgradeTag := os.Getenv("DTR_UPGRADE_TAG")
	if upgradeTag == "" {
		log.Info("Skipping upgrade test due to unspecified upgrade tag")
		return
	}
	if len(machines) < 3 {
		log.Warn(fmt.Sprintf("DTR will experience downtime with only %d machines. There must be at least 3 machines to have zero downtime upgrades.", len(machines)))
	}

	upgradeImageArgs, err := util.GenerateImageArgs(os.Getenv("DTR_REPO"), upgradeTag)
	if err != nil {
		log.Error(err)
	}

	upgradeArgs := util.UpgradeArgs{}
	upgradeArgs.DTRImageArgs = upgradeImageArgs

	err = util.PushDTRImages(machines, upgradeImageArgs)
	if err != nil {
		log.Error(err)
	}

	args := util.InstallArgs{}
	util.DeployDTR(machines, args, len(machines))
	util.FullSmokeTest(machines, util.DefaultInstallation.Args.DTRImageArgs, t)

	util.UpgradeDTR(machines, upgradeArgs)
	util.FullSmokeTest(machines, upgradeImageArgs, t)
	util.NotaryTest(machines, upgradeImageArgs, t)
	util.PushPullTest(machines, upgradeImageArgs, t)

	util.RemoveDTR(machines, upgradeImageArgs)
	util.FullSmokeTest(machines, upgradeImageArgs, t)
}
