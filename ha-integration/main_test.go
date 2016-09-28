package ha_integration

import (
	"os"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/dhe-deploy/ha-integration/ha_utils"
	"github.com/docker/dhe-deploy/ha-integration/suite"
	"github.com/docker/dhe-deploy/ha-integration/util"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"
)

func TestSuite(t *testing.T) {
	if config.DefaultReporterConfig.Verbose || testing.Verbose() {
		log.SetLevel(log.DebugLevel)
	}
	ha_utils.LogFile = GinkgoWriter
	log.SetOutput(GinkgoWriter)
	log.Debug("logging in debug mode")

	RegisterFailHandler(Fail)
	RunSpecs(t, "Everything Suite")
}

var machines []ha_utils.Machine

var _ = BeforeSuite(func() {
	util.DefaultInstallation.Replicas = make(map[string]string)
	imageArgs, err := util.GenerateImageArgs(os.Getenv("DTR_REPO"), os.Getenv("DTR_TAG"))
	if err != nil {
		panic(err)
	}

	util.DefaultInstallation.Args = util.DTRArgs{
		UCPArgs: util.UCPArgs{
			AdminUsername:  "admin",
			AdminPassword:  "orca",
			UCPInsecureTLS: true,
		},
		DTRImageArgs: imageArgs,
		HTTPArgs: util.HTTPArgs{
			HTTPProxy:        "",
			HTTPSProxy:       "",
			NoProxy:          false,
			ReplicaHTTPPort:  "",
			ReplicaHTTPSPort: "",
			DTRHost:          "",
		},
		EnziArgs: util.EnziArgs{
			EnziInsecureTLS: true,
			EnziCA:          "",
		},
	}

	util.DefaultInstallation.PrimaryReplicaID = ""

	log.Info("about to prepare")
	machines = suite.Prepare()

	// at this point the deployment step is done and we're about to start some tests
	if testing.Short() {
		Skip("skipping HA integration tests in short mode.")
	}

	dtrImageArgs := util.DTRImageArgs{}
	err = util.PushDTRImages(machines, dtrImageArgs)
	if err != nil {
		log.Error(err)
	}
})

var _ = Describe("in a prepared cluster", func() {
	BeforeEach(func() {
		util.PurgeDTR(machines)
	})

	installTests()
	upgradeTests()
	removeTests()
})
