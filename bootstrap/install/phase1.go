package install

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/dhe-deploy/bootstrap"
	"github.com/docker/dhe-deploy/bootstrap/flags"
)

func Phase1(c *cli.Context) (string, error) {
	var err error
	log.Info("Beginning Docker Trusted Registry installation")

	bootstrap.PromptIfNotSet(c, flags.DTRHostFlag)

	dropperOpts, err := bootstrap.ParseAndPromptDropperOpts(c)
	if err != nil {
		return "", err
	}

	bs, err := bootstrap.GetBootstrapClient(c, dropperOpts)
	if err != nil {
		return "", err
	}
	log.Debugf("bootstrap = %q", bs)

	container := dropperOpts.MakeContainerConfig([]string{"install"})
	bootstrap.SetEnvFromFlags(c, container, flags.InstallFlags...)

	// Look to see if http/https ports are already in use, and exclude those nodes
	if !flags.NoUCP {
		nodes, err := bootstrap.SetupNodePortConstraints(c, bs)
		if err != nil {
			return "", err
		}
		container.ExcludeNodes = *nodes
	}

	replicaIDs, err := bs.ListReplicaIDs(true)
	if err != nil {
		return "", err
	}
	if len(replicaIDs) > 0 {
		return "", fmt.Errorf("We don't currently support running multiple DTRs on the same UCP cluster. You have %d existing DTR replicas. Please remove them before trying to install again.", len(replicaIDs))
	}

	replicaID, err := bs.ValidateOrGetNewReplicaID(flags.ReplicaID)
	if err != nil {
		return "", err
	}
	bs.SetReplicaID(replicaID)

	err = bootstrap.Phase2Execute(container, bs)
	if err != nil {
		bootstrap.CheckConstraintsError(err)
		return "", err
	}
	return replicaID, nil
}
