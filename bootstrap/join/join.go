package join

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/dhe-deploy/bootstrap"
)

func join(c *cli.Context) (int, error) {
	if !bootstrap.IsPhase2() {
		return phase1(c)
	} else {
		return phase2(c)
	}

	return 0, nil
}

func Run(c *cli.Context) {
	bootstrap.ConfigureLogging()
	if rc, err := join(c); err != nil {
		log.Fatal(err)
	} else {
		os.Exit(rc)
	}
}
