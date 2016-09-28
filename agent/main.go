package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"

	"github.com/docker/orca/agent/agent"
	"github.com/docker/orca/agent/agent/reconcile"
	"github.com/docker/orca/agent/proxy"
	"github.com/docker/orca/agent/testserver"
	"github.com/docker/orca/version"
)

// The `agent` command launches the UCP agent which performs state inspection and launches the `ucp-reconcile` container
var cmdAgent = cli.Command{
	Name:  "agent",
	Usage: "Inspect the current node state and bootstrap UCP",
	Description: `
	`,
	Action: agent.StartAgent,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "d",
			Usage: "Path to the Docker socket",
			Value: "/var/run/docker.sock",
		},
		cli.StringFlag{
			Name:   "image_version",
			Usage:  "The version of the docker image and tag to use for the current node",
			EnvVar: "IMAGE_VERSION",
			Hidden: true,
		},
		cli.StringFlag{
			Name:   "secret",
			Usage:  "The secret to be used as part of the CSR handshake when joining an existing UCP cluster",
			EnvVar: "SECRET",
			Hidden: true,
		},
		cli.StringFlag{
			Name:   "swarm_port",
			Usage:  "The port on the host machine where the swarm-join container will listen on",
			EnvVar: "SWARM_PORT",
		},
		cli.StringFlag{
			Name:   "ucp_instance_id",
			Usage:  "The instance ID of the UCP node to be bootstrapped",
			EnvVar: "UCP_INSTANCE_ID",
		},
		cli.StringFlag{
			Name:   "controller_port",
			Usage:  "The port on the host machine where the ucp-controller container will listen on, if the node is a manager",
			EnvVar: "CONTROLLER_PORT",
		},
		cli.StringSliceFlag{
			Name:   "dns",
			Usage:  "Set custom DNS servers for the UCP infrastructure containers",
			EnvVar: "DNS",
			Value:  &cli.StringSlice{},
		},
		cli.StringSliceFlag{
			Name:   "dns_opt",
			Usage:  "Set DNS options for the UCP infrastructure containers",
			EnvVar: "DNS_OPT",
			Value:  &cli.StringSlice{},
		},
		cli.StringSliceFlag{
			Name:   "dns_search",
			Usage:  "Set custom DNS search domains for the UCP infrastructure containers",
			EnvVar: "DNS_SEARCH",
			Value:  &cli.StringSlice{},
		},
		cli.BoolFlag{
			Name:   "debug, D",
			Usage:  "Enable Debug output",
			EnvVar: "DEBUG",
		},
	},
}

// The `proxy` command proxies all requests to the docker engine
var cmdProxy = cli.Command{
	Name:  "proxy",
	Usage: "Run a Docker Engine proxy with mTLS",
	Description: `
		To enable TLS, pass the SSL_CA, SSL_CERT and SSL_KEY environment
		variables. To disable SSL verification, set the SSL_SKIP_VERIFY 
		environemnt variable. CLI flags can be used in place of any of these
		options as well.
	`,
	Action: proxy.ProxyServer,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "d",
			Usage: "Path to the Docker socket",
			Value: "/var/run/docker.sock",
		},
		cli.StringFlag{
			Name:  "listen-address, l",
			Usage: "Listen Address",
			Value: ":2376",
		},
		cli.StringFlag{
			Name:   "ca",
			Usage:  "Path to CA certificate",
			EnvVar: "SSL_CA",
		},
		cli.StringFlag{
			Name:   "cert",
			Usage:  "Path to server certificate",
			EnvVar: "SSL_CERT",
		},
		cli.StringFlag{
			Name:   "key",
			Usage:  "Path to certificate key",
			EnvVar: "SSL_KEY",
		},
		cli.StringFlag{
			Name:   "insecure, i",
			Usage:  "Disable SSL verifcation",
			EnvVar: "SSL_SKIP_VERIFY",
		},
		cli.BoolFlag{
			Name:  "debug, D",
			Usage: "Enable Debug output",
		},
	},
}

// The `test-server` command returns a simple 200 OK at the / endpoint, meant to debug port connectivity
var cmdTestServer = cli.Command{
	Name:  "test-server",
	Usage: "Returns 200 OK at /",
	Description: `
	`,
	Action: testserver.TestServer,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "listen-address, l",
			Usage: "Listen Address",
			Value: ":2376",
		},
	},
}

// The `reconcile` command starts phase2 of the agent, trigerring a current vs expected state reconciliation
var cmdReconcile = cli.Command{
	Name:  "reconcile",
	Usage: "Launches the Phase 2 Reconciliation of the UCP Agent",
	Description: `
	`,
	Action: reconcile.Reconcile,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "d",
			Usage: "Path to the Docker socket",
			Value: "/var/run/docker.sock",
		},
		cli.StringFlag{
			Name:  "payload",
			Usage: "The payload passed from the agent to the reconcile container",
		},
	},
}

// Driver function
func main() {
	app := cli.NewApp()
	app.Name = "UCP Agent"
	app.Usage = "UCP Node Agent - Proxy, Bootstrapper and Reconcile"
	app.Version = version.FullVersion()
	app.Commands = []cli.Command{
		cmdProxy,
		cmdAgent,
		cmdTestServer,
		cmdReconcile,
	}
	log.SetFormatter(&log.JSONFormatter{})

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

	os.Exit(1)
}
