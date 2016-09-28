package reconcile

import (
	"fmt"
	"os"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/engine-api/client"

	"github.com/docker/orca/agent/agent/components"
	"github.com/docker/orca/agent/agent/utils"
	"github.com/docker/orca/bootstrap/config"
	butils "github.com/docker/orca/bootstrap/utils"
	orcaconfig "github.com/docker/orca/config"
)

// Reconcile is invoked for the `reconcile` command
func Reconcile(c *cli.Context) error {
	// Create an engine-api client to the local engine
	dockerSocket := c.String("d")
	if _, err := os.Stat(dockerSocket); os.IsNotExist(err) {
		return fmt.Errorf("Unable to locate the Docker Socket at %s", dockerSocket)
	}
	dclient, err := client.NewClient(fmt.Sprintf("unix://%s", dockerSocket), "", nil, nil)
	if err != nil {
		return fmt.Errorf("could not create a docker engine-api client: %s", err)
	}
	log.Infof("Initialized engine-api client, version %s", dclient.ClientVersion())

	expected, current, err := utils.DeserializeReconcileArgs(c.String("payload"))
	if err != nil {
		return err
	}

	// TODO: figure out log levels
	log.SetLevel(log.DebugLevel)

	// Set up global package configuration goop - TODO: remove
	config.OrcaInstanceID = expected.UCPInstanceID
	orcaconfig.OrcaPort, err = strconv.Atoi(expected.ControllerPort)
	if err != nil {
		return err
	}
	orcaconfig.SwarmPort, err = strconv.Atoi(expected.SwarmPort)
	if err != nil {
		return err
	}
	orcaconfig.ImageVersion = expected.ImageVersion
	config.DNS = expected.DNS
	config.DNSOpt = expected.DNSOpt
	config.DNSSearch = expected.DNSSearch
	config.OrcaHostAddress = expected.HostAddress
	config.OrcaLocalName = config.OrcaHostAddress
	// TODO: gather a set of addresses from various interfaces
	config.OrcaHostnames = []string{"127.0.0.1", "localhost", config.OrcaHostAddress}
	log.Info("Configuring node as agent with the following SANs: ", config.OrcaHostnames)

	// TODO: exit fast on blocked ports
	/* XXX WRONG!  Need some way to turn this off when we're reconciling a node that's already set up
	// Detect fresh install files
	if !utils.IsFreshInstall() {
		// Skip the port check in a fresh install, as it's done by the bootstrapper earlier
		ec, err := engineClient.NewBareClient()
		if err != nil {
			return err
		}
		err = ec.CheckPorts(orcaconfig.RequiredPorts)
		if err != nil {
			return err
		}
	}
	*/

	// XXX Also wrong - but need to figure out the right place to wire this up
	if err := butils.VerifyPermissions(true); err != nil {
		return err
	}

	// Begin reconciliation
	for _, component := range components.ComponentList {
		log.Debugf("Reconciling component %#v", component)
		err = component.Reconcile(dclient, expected, current)
		if err != nil {
			log.Warnf("Component %#v failed to reconcile with %s", component, err)
			return err
		}
	}

	return nil
}
