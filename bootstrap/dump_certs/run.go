// Package dump_certs display the public certs for this Orca
package dump_certs

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/orca/bootstrap/client"
	"github.com/docker/orca/bootstrap/config"
)

func dump(c *cli.Context) (int, error) {
	config.HandleGlobalArgs(c)

	ec, err := client.New()
	if err != nil {
		return 1, err
	}
	if ec.IsTty() {
		return 1, fmt.Errorf("Certificates will be corrupted with TTY mode enabled.  Please re-run without the '-t' flag.")
	}

	if !config.InPhase2 {
		config.OrcaInstanceID = "dump-certs" // Dummy name since it doesn't matter
		// TODO - Might want to make this have affinity to orca-controller so it can be run against orca itself
		var volume string
		var mount string

		if c.Bool("cluster") {
			volume = config.SwarmNodeCertVolumeName
			mount = config.SwarmNodeCertVolumeMount
		} else {
			volume = config.OrcaServerCertVolumeName
			mount = config.OrcaServerCertVolumeMount
		}
		config.Phase2VolumeMounts = append(config.Phase2VolumeMounts, []string{
			fmt.Sprintf("%s:%s:ro", volume, mount),
		}...)
		return ec.StartPhase2(os.Args[1:], false)

	} else {
		var mount string
		if c.Bool("cluster") {
			mount = config.SwarmNodeCertVolumeMount
		} else {
			mount = config.OrcaServerCertVolumeMount
		}
		files := []string{filepath.Join(mount, "ca.pem")}
		if !c.Bool("ca") {
			files = append(files, filepath.Join(mount, "cert.pem"))
		}
		for _, filename := range files {
			if data, err := ioutil.ReadFile(filename); err != nil {
				return 1, fmt.Errorf("Missing expected cert file on this host %s: %s", filename, err)
			} else {
				os.Stdout.Write(data)
			}
		}
	}
	return 0, nil
}

// Run the dump cert flow
func Run(c *cli.Context) {
	if code, err := dump(c); err != nil {
		log.Fatal(err)
	} else {
		os.Exit(code)
	}
}
