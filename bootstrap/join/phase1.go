package join

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/dhe-deploy"
	"github.com/docker/dhe-deploy/bootstrap"
	"github.com/docker/dhe-deploy/bootstrap/flags"
)

func phase1(c *cli.Context) (int, error) {
	var ourReplicaID string

	log.Infof("Beginning Docker Trusted Registry replica join")
	dropperOpts, err := bootstrap.ParseAndPromptDropperOpts(c)
	if err != nil {
		return 1, err
	}

	bs, err := bootstrap.GetBootstrapClient(c, dropperOpts)
	if err != nil {
		return 1, err
	}

	replica, err := bs.ExistingReplicaFlagPicker("Choose a replica to join to", true)
	if err != nil {
		return 1, err
	}

	if err := replica.CheckEqualVersion(deploy.Version); err != nil {
		return 1, err
	}

	ourReplicaID, err = bs.ValidateOrGetNewReplicaID(flags.ReplicaID)
	if err != nil {
		return 1, err
	}
	bs.SetReplicaID(ourReplicaID)

	container := dropperOpts.MakeContainerConfig([]string{"join"})
	bootstrap.SetEnvFromFlags(c, container, flags.JoinFlags...)

	// Look to see if http/https ports are already in use, and exclude those nodes
	if !flags.NoUCP {
		nodes, err := bootstrap.SetupNodePortConstraints(c, bs)
		if err != nil {
			return 1, err
		}
		container.ExcludeNodes = *nodes
	}

	err = bootstrap.Phase2Execute(container, bs)
	if err != nil {
		bootstrap.CheckConstraintsError(err)
		return 1, err
	}
	return 0, nil
}
