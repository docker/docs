package main

import "github.com/codegangsta/cli"

const (
	osFlag             = "os"
	osVersionFlag      = "os-version"
	driverFlag         = "driver"
	esxHostFlag        = "esx-host"
	esxUserFlag        = "esx-user"
	esxPassFlag        = "esx-pass"
	esxDatastoreFlag   = "esx-datastore"
	tokenFlag          = "token"
	buildFlag          = "build"
	registryProxyFlag  = "registry-proxy"
	skipRevertFlag     = "skip-revert"
	skipInstallFlag    = "skip-install"
	skipTestFlag       = "skip-test"
	runTestCommandFlag = "run-test-command"
	grabDirFlag        = "grab-dir"

	driverESX    = "esx"
	driverFusion = "fusion"
)

var flags = []string{
	osFlag,
	osVersionFlag,
	driverFlag,
	esxHostFlag,
	esxUserFlag,
	esxPassFlag,
	esxDatastoreFlag,
	tokenFlag,
	buildFlag,
	registryProxyFlag,
	skipRevertFlag,
}

var cliFlags = []cli.Flag{
	cli.StringFlag{
		Name:   driverFlag + ", d",
		Value:  "esx",
		Usage:  "Driver to use esx|fusion",
		EnvVar: "PITFALL_DRIVER",
	},
	cli.StringFlag{
		Name:   osFlag + ", o",
		Value:  "osx",
		Usage:  "OS under test",
		EnvVar: "PITFALL_OS",
	},
	cli.StringFlag{
		Name:   osVersionFlag + ", x",
		Value:  "10.11",
		Usage:  "OS Version under test",
		EnvVar: "PITFALL_OS_VERSION",
	},
	cli.StringFlag{
		Name:   esxHostFlag + ", e",
		Value:  "172.16.1.10",
		Usage:  "ESXi Host",
		EnvVar: "PITFALL_ESX_HOST",
	},
	cli.StringFlag{
		Name:   esxUserFlag + ", u",
		Value:  "root",
		Usage:  "ESXi User",
		EnvVar: "PITFALL_ESX_USER",
	},
	cli.StringFlag{
		Name:   esxPassFlag + ", p",
		Value:  "slartibartfast",
		Usage:  "ESXi Password",
		EnvVar: "PITFALL_ESX_PASS",
	},
	cli.StringFlag{
		Name:   esxDatastoreFlag + " , s",
		Value:  "datastore2",
		Usage:  "ESXi Datastore",
		EnvVar: "PITFALL_ESX_DATASTORE",
	},
	cli.StringFlag{
		Name:   buildFlag + ", b",
		Usage:  "Build Number",
		EnvVar: "PITFALL_BUILD_ID",
	},
	cli.StringFlag{
		Name:   registryProxyFlag,
		Usage:  "URL to registry mirror/proxy",
		EnvVar: "PITFALL_REGISTRY_PROXY",
	},
	cli.StringFlag{
		Name:  runTestCommandFlag,
		Usage: "Command to run tests",
		Value: "./rt-local -l release,nostart -v -x run",
	},
	cli.StringFlag{
		Name:  grabDirFlag,
		Usage: "Directory to grab after testing, relative to test dir",
		Value: "_results/",
	},
	cli.BoolFlag{
		Name:  skipRevertFlag,
		Usage: "Skip revert step and connect to an already running VM",
	},
	cli.BoolFlag{
		Name:  skipInstallFlag,
		Usage: "Skip install step and use already installed D4x",
	},
	cli.BoolFlag{
		Name:  skipTestFlag,
		Usage: "Skip test run (mainly for debugging pitfall)",
	},
}

// ParseConfig parses the config in to a map
func ParseConfig(c *cli.Context) map[string]string {
	config := make(map[string]string)
	for _, flag := range flags {
		v := c.String(flag)
		config[flag] = v
	}
	return config
}
