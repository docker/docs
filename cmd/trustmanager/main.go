package main

import (
	"crypto/x509"
	"fmt"
	"os"
	"os/user"
	"path"

	"github.com/codegangsta/cli"
	"github.com/docker/vetinari/trustmanager"
)

const caDir string = ".docker/trust/certificate_authorities/"
const repoDir string = ".docker/trust/repositories/"

var caStore trustmanager.X509Store
var repoStore trustmanager.X509Store

func init() {

	// Retrieve current user to get home directory
	usr, err := user.Current()
	if err != nil {
		errorf("cannot get current user: %v", err)
	}

	// Get home directory for current user
	homeDir := usr.HomeDir
	if homeDir == "" {
		errorf("cannot get current user home directory")
	}

	// Ensure the existence of the CAs directory
	fullCaDir, fullRepoDir := setupDefaultDirectories(homeDir)

	// TODO(diogo): inspect permissions of the directories/files. Warn.
	caStore = trustmanager.NewX509FilteredFileStore(fullCaDir, func(cert *x509.Certificate) bool {
		return cert.IsCA
	})
	repoStore = trustmanager.NewX509FileStore(fullRepoDir)
}

func main() {
	app := cli.NewApp()
	app.Name = "keymanager"
	app.Usage = "trust keymanager"

	app.Commands = []cli.Command{
		commandTrust,
		commandList,
		commandUntrust,
	}

	app.RunAndExitOnError()
}

func errorf(format string, args ...interface{}) {
	fmt.Printf("* fatal: "+format+"\n", args...)
	os.Exit(1)
}

func setupDefaultDirectories(homeDir string) (string, string) {
	fullCaDir := path.Join(homeDir, path.Dir(caDir))
	if err := os.MkdirAll(fullCaDir, 0700); err != nil {
		errorf("cannot create directory: %v", err)
	}

	// Ensure the existence of the repositories directory
	fullRepoDir := path.Join(homeDir, path.Dir(repoDir))
	if err := os.MkdirAll(fullRepoDir, 0700); err != nil {
		errorf("cannot create directory: %v", err)
	}

	return fullCaDir, fullRepoDir
}
