package ha_integration

import (
	"fmt"
	"os"
	"strconv"

	"github.com/docker/dhe-deploy/ha-integration/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func installTests() {
	Context("install tests", func() {
		It("installs regularly with constraints", func() {
			dtrConstraints := os.Getenv("DTR_CONSTRAINTS")
			dtrEnablePProf := os.Getenv("DTR_ENABLE_PPROF") != ""
			args := util.InstallArgs{
				ConstraintArgs: util.ConstraintArgs{Constraints: dtrConstraints},
				PProfArgs:      util.PProfArgs{EnablePProf: dtrEnablePProf},
			}

			strReplicas := os.Getenv("DTR_REPLICAS")
			if strReplicas == "" {
				strReplicas = "0"
			}

			dtrReplicas, err := strconv.ParseInt(strReplicas, 10, 64)
			Expect(err).To(BeNil())

			util.DeployDTR(machines, args, int(dtrReplicas))
			// we can't do a smoke test because it doesn't handle constraints properly and
			// has no idea what to expect
		})
		It("installs regularly on 3 machines and removes forcefully", func() {
			// TODO: convert to ginkgo style tests
			// TODO: get rid of the t
			t := GinkgoT()
			if len(machines) < 3 {
				Skip(fmt.Sprintf("Needed 3 machines to run install test, have %d", len(machines)))
			}
			args := util.InstallArgs{}
			imageArgs := util.DefaultInstallation.Args.DTRImageArgs
			util.DeployDTR(machines, args, 3)
			util.FullSmokeTest(machines, imageArgs, t)
			util.NotaryTest(machines, imageArgs, t)
			util.PushPullTest(machines, imageArgs, t)
			util.RemoveDTR(machines, imageArgs)
			util.FullSmokeTest(machines, imageArgs, t)
		})
		It("installs regularly on 1 machine and can push/pull", func() {
			// TODO: once running HA tests works with the filesystem driver or another driver, remove this test as well
			// as the check in `PushPullTest` for exactly one DTR replica machine
			t := GinkgoT()
			if len(machines) < 1 || len(machines) > 2 {
				Skip(fmt.Sprintf("Needed 1 machine to run install test, have %d", len(machines)))
			}
			args := util.InstallArgs{}
			imageArgs := util.DefaultInstallation.Args.DTRImageArgs
			util.DeployDTR(machines, args, 1)
			util.PushPullTest(machines, imageArgs, t)
			util.RemoveDTR(machines, imageArgs)
		})
	})
}
