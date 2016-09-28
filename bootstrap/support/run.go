package support

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/orca/bootstrap/client"
	"github.com/docker/orca/bootstrap/config"
)

func support(c *cli.Context) (int, error) {
	config.HandleGlobalArgs(c)

	ec, err := client.New()
	if err != nil {
		return 1, err
	}
	if ec.IsTty() {
		return 1, fmt.Errorf("The support dump will be corrupted with TTY mode enabled.  Please re-run without the '-t' flag.")
	}

	return ec.SupportDump()
}

// Run the fingerprint flow
func Run(c *cli.Context) {
	if code, err := support(c); err != nil {
		log.Fatal(err)
	} else {
		os.Exit(code)
	}
}
