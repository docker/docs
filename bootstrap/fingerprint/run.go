// Package dump_certs display the public certs for this Orca
package fingerprint

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/orca/bootstrap/certs"
	"github.com/docker/orca/bootstrap/client"
	"github.com/docker/orca/bootstrap/config"
)

func fingerprint(c *cli.Context) (int, error) {
	config.HandleGlobalArgs(c)

	ec, err := client.New()
	if err != nil {
		return 1, err
	}
	if ec.IsTty() {
		return 1, fmt.Errorf("Fingerprints will be corrupted with TTY mode enabled.  Please re-run without the '-t' flag.")
	}

	if !config.InPhase2 {
		config.OrcaInstanceID = "fingerprint" // Dummy name since it doesn't matter
		// TODO - Might want to make this have affinity to orca-controller so it can be run against orca itself
		config.Phase2VolumeMounts = append(config.Phase2VolumeMounts,
			fmt.Sprintf("%s:%s", config.OrcaServerCertVolumeName, config.OrcaServerCertVolumeMount))
		return ec.StartPhase2(os.Args[1:], false)

	} else {

		if fingerprint, err := certs.GetFingerprint(filepath.Join(config.OrcaServerCertVolumeMount, "cert.pem")); err != nil {
			return 1, fmt.Errorf("Unable to locate UCP cert on this host.  Is the UCP server installed here? - %s", err)
		} else {
			fmt.Println(fingerprint)
		}
	}
	return 0, nil
}

// Run the fingerprint flow
func Run(c *cli.Context) {
	if code, err := fingerprint(c); err != nil {
		log.Fatal(err)
	} else {
		os.Exit(code)
	}
}
