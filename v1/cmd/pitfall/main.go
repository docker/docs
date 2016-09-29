package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"golang.org/x/crypto/ssh"
)

const (
	osxSnapshot = "SSH Enabled"
	winSnapshot = "Hyper-V"
	// OSX is a streing constant for Apple OSX
	OSX = "osx"
	// Win is a string constant for Windows
	Win = "win"
)

var (
	supportedOS = []string{OSX, Win}
	osxVMs      = map[string]string{
		"10.10": "OS X 10.10",
		"10.11": "OS X 10.11",
	}
	winVMs = map[string]string{
		"10": "Microsoft Windows 10 (64-bit)",
	}
)

func fatal(err error) {
	log.Fatalf("😿: %s\n", err.Error())
	os.Exit(1)
}

func main() {
	app := cli.NewApp()
	app.Name = "pitfall"
	app.Usage = "Pinata Integration Tests For All!"
	app.UsageText = "pitfall [global options] path_to_tests"
	app.Version = "1.0.0"
	app.Action = appMain
	app.Flags = cliFlags
	app.Run(os.Args)
}

func appMain(c *cli.Context) {
	var testPath string
	if c.NArg() > 0 {
		testPath = c.Args()[0]
	} else {
		fatal(fmt.Errorf("Please provide the path to the 'tests' directory"))
	}

	config := ParseConfig(c)

	var driver Driver
	var err error
	switch config[driverFlag] {
	case driverESX:
		driver, err = NewESXDriver(config)
	case driverFusion:
		driver, err = NewFusionDriver(config)
	default:
		fatal(fmt.Errorf("Driver %s is not supported", config[driverFlag]))
	}
	if err != nil {
		fatal(err)
	}

	skipInstall := c.Bool(skipInstallFlag)
	build := c.String("build")
	if build == "" {
		if skipInstall {
			build = "unknown"
		} else {
			fatal(fmt.Errorf("Build ID is reqired"))
		}
	}

	skipRevert := c.Bool(skipRevertFlag)
	var ip string
	// Revert and start VM
	if !skipRevert {
		ip, err = driver.RevertVMToSnapshot()
		if err != nil {
			fatal(err)
		}
	} else {
		ip, err = driver.GetIP()
		if err != nil {
			fatal(err)
		}
	}

	// Initiate SSH
	sshConfig := &ssh.ClientConfig{
		User: "docker",
		Auth: []ssh.AuthMethod{
			ssh.Password("containyourself"),
		},
	}

	remoteHost := fmt.Sprintf("[%s]:22", ip)
	var client *ssh.Client
	for i := 0; i < 12; i++ {
		client, err = ssh.Dial("tcp", remoteHost, sshConfig)
		if err != nil {
			log.Warning("Failed to dial ", remoteHost, ": ", err, " (retrying)")
			time.Sleep(10 * time.Second)
		} else {
			break
		}
	}
	if err != nil {
		fatal(fmt.Errorf("Failed to dial %s, no more retries", remoteHost))
	} else {
		log.Info("Connected to ", remoteHost)
	}

	var path string
	if filepath.IsAbs(testPath) {
		path = testPath
	} else {
		path, err = filepath.Abs(testPath)
		if err != nil {
			fatal(err)
		}
	}

	log.Info("Uploading Test Dir to server")
	var errc <-chan error
	errc, err = UploadDir(client, path, "pinata")
	if err != nil {
		fatal(err)
	}
	err = <-errc
	if err != nil {
		fatal(err)
	}

	regProxy := c.String("registry-proxy")

	if !skipInstall {
		log.Info("Attempting to install Docker for Mac")

		if regProxy != "" {
			log.Info("Asking client to use registry proxy at ", regProxy)
		}

		err = ExecCmd(client, fmt.Sprintf("cd /Users/docker/pinata && ./mac-autoinst.sh %s", regProxy), false)
		if err != nil {
			fatal(err)
		}

	} else {
		log.Info("Skipping install step")
		if regProxy != "" {
			log.Info("(Registry proxy flag ignored, using existing installation)")
		}
	}

	log.Info("Setting bashrc...")
	err = ExecCmd(client, "echo 'export PATH=/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin' > .bashrc", false)
	if err != nil {
		fatal(err)
	}

	skipTest := c.Bool(skipTestFlag)
	var testsFailed bool
	if !skipTest {
		log.Info("Running the tests...")
		runTestCommand := c.String(runTestCommandFlag)
		log.Info("Command: ", runTestCommand)
		err = ExecCmd(client, fmt.Sprintf("cd /Users/docker/pinata && %s", runTestCommand), true)
		testsFailed = false
		if err != nil {
			testsFailed = true
		}
	} else {
		testsFailed = false
	}

	grabDir := c.String(grabDirFlag)
	log.Info("Grabbing results from ", grabDir)
	err = ExecCmd(client, fmt.Sprintf("cd /Users/docker/pinata && tar -zcvf results.tar.gz %s", grabDir), false)
	if err != nil {
		fatal(err)
	}

	var file string
	for i := 0; i < 9999999; i++ {
		file = fmt.Sprintf("b-%s-%s-%s-results-%d.tar.gz", build, c.String("os"), c.String("os-version"), i)
		if _, err := os.Stat(file); err != nil {
			break
		}
	}

	log.Info("Storing logs in ", file)

	err = DownloadResults(client, "/Users/docker/pinata/results.tar.gz", file)
	if err != nil {
		fatal(err)
	}

	if testsFailed {
		fatal(fmt.Errorf("The tests failed. See the logs in %s", file))
	}
	log.Infof("SUCCESS! 🐳 Logs available in %s", file)
}
