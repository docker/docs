// Package id displays the UCP instance ID
package id

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/orca/bootstrap/client"
	"github.com/docker/orca/bootstrap/config"
	orcaconfig "github.com/docker/orca/config"
)

func id(c *cli.Context) (int, error) {
	config.HandleGlobalArgs(c)

	ec, err := client.New()
	if err != nil {
		return 1, err
	}
	/* TODO After HA CA change merges
	   if ec.IsTty() {
	       return 1, fmt.Errorf("ID will be corrupted with TTY mode enabled.  Please re-run without the '-t' flag.")
	   }
	*/

	containers := ec.FindContainers(orcaconfig.RuntimeContainerNames)
	ids := client.GetInstanceIDs(containers)
	if len(ids) == 0 {
		return 1, fmt.Errorf("No running UCP instances detected on this engine")
	} else if len(ids) > 1 {
		return 1, fmt.Errorf("Multiple UCP instances detected: %v", ids)
	}
	fmt.Println(ids[0])
	return 0, nil
}

// Run the id flow
func Run(c *cli.Context) {
	if code, err := id(c); err != nil {
		log.Fatal(err)
	} else {
		os.Exit(code)
	}
}
