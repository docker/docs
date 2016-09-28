package main

import (
	"fmt"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/orca/bootstrap/backup"
	"github.com/docker/orca/bootstrap/config"
	"github.com/docker/orca/bootstrap/dump_certs"
	"github.com/docker/orca/bootstrap/fingerprint"
	"github.com/docker/orca/bootstrap/id"
	"github.com/docker/orca/bootstrap/images"
	"github.com/docker/orca/bootstrap/install"
	"github.com/docker/orca/bootstrap/regen_certs"
	"github.com/docker/orca/bootstrap/restart"
	"github.com/docker/orca/bootstrap/restore"
	"github.com/docker/orca/bootstrap/stop"
	"github.com/docker/orca/bootstrap/support"
	"github.com/docker/orca/bootstrap/uninstall"
	"github.com/docker/orca/bootstrap/upgrade"
	orcaconfig "github.com/docker/orca/config"
	"github.com/docker/orca/version"
)

func buildFlags(flags ...[]cli.Flag) []cli.Flag {
	res := []cli.Flag{}
	for _, flagSet := range flags {
		res = append(res, flagSet...)
	}
	return res
}

func main() {
	commonFlags := []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, D",
			Usage: "Enable debug mode",
		},
		cli.BoolFlag{
			Name:  "jsonlog",
			Usage: "Produce json formatted output for easier parsing",
		},
	}
	commonInteractiveFlags := []cli.Flag{
		cli.BoolFlag{
			Name:  "interactive, i",
			Usage: "Enable interactive mode.  You will be prompted to enter all required information",
		},
	}
	commonCredentialFlags := []cli.Flag{
		cli.StringFlag{
			Name:   "admin-username",
			Usage:  "Specify the UCP admin username",
			EnvVar: "UCP_ADMIN_USER",
		},
		cli.StringFlag{
			Name:   "admin-password",
			Usage:  "Specify the UCP admin password",
			EnvVar: "UCP_ADMIN_PASSWORD",
		},
	}
	commonInstallFlags := []cli.Flag{
		// XXX Might want a "redeploy" option, to replace containers but persist state
		cli.BoolFlag{
			Name:  "fresh-install",
			Usage: "Destroy any existing state and start fresh",
		},
		cli.StringSliceFlag{
			Name:  "san",
			Usage: "Additional Subject Alternative Names for certs (e.g. --san foo1.bar.com --san foo2.bar.com)",
			Value: &cli.StringSlice{},
		},
		cli.StringFlag{
			Name:   "host-address",
			Usage:  "Specify the visible IP/hostname for this node (override automatic detection)",
			EnvVar: "UCP_HOST_ADDRESS",
		},
		cli.IntFlag{
			Name:  "swarm-port",
			Usage: "Select what port to run the local Swarm manager on",
			Value: 2376,
		},
		cli.IntFlag{
			Name:  "controller-port",
			Usage: "Select what port to run the local UCP Controller on",
			Value: 443,
		},
		cli.IntFlag{
			Name:  "swarm-grpc-port",
			Usage: "Select what port to run Swarm GRPC on",
			Value: 2377,
		},
		cli.StringSliceFlag{
			Name:   "dns",
			Usage:  "Set custom DNS servers for the UCP infrastructure containers",
			EnvVar: "DNS",
			Value:  &cli.StringSlice{},
		},
		cli.StringSliceFlag{
			Name:   "dns-opt",
			Usage:  "Set DNS options for the UCP infrastructure containers",
			EnvVar: "DNS_OPT",
			Value:  &cli.StringSlice{},
		},
		cli.StringSliceFlag{
			Name:   "dns-search",
			Usage:  "Set custom DNS search domains for the UCP infrastructure containers",
			EnvVar: "DNS_SEARCH",
			Value:  &cli.StringSlice{},
		},
		cli.IntFlag{
			Name:  "kv-timeout",
			Usage: "Timeout in milliseconds for the KV store (set higher for a multi-datacenter cluster)",
			Value: 1000,
		},
	}
	commonImageFlags := []cli.Flag{
		cli.StringFlag{
			Name:   "image-version",
			Value:  strings.Split(version.FullVersion(), " ")[0],
			Hidden: true,
		},
		cli.StringFlag{
			Name:  "pull",
			Usage: "Specify image pull behavior ('always', when 'missing', or 'never')",
			Value: "missing",
		},
		cli.StringFlag{
			Name:   "registry-username",
			Usage:  "Specify the username to pull required images with",
			EnvVar: "REGISTRY_USERNAME",
		},
		cli.StringFlag{
			Name:   "registry-password",
			Usage:  "Specify the password to pull required images with",
			EnvVar: "REGISTRY_PASSWORD",
		},
	}
	cli.AppHelpTemplate = `Docker Universal Control Plane Tool

This tool has commands to 'install' and 'uninstall' UCP.
This tool must run as a container with a well-known name and with the
docker.sock volume mounted, which you can cut-and-paste from the usage
example below.

This tool will generate TLS certificates and will attempt to determine
your hostname and primary IP addresses.  This may be overridden with the
'--host-address' option.  The tool may not be able to discover your
externally visible fully qualified hostname.  For proper certificate
verification, you should pass one or more Subject Alternative Names with
'--san' during 'install' that matches the fully qualified hostname you
intend to use to access the given system.

Many settings can be passed as flags or environment variables. When passing as
an environment variable, use the 'docker run -e VARIABLE_NAME ...' syntax to
pass the value from your shell, or 'docker run -e VARIABLE_NAME=value ...' to
specify the value explicitly on the command line.

Additional help is available for each command with the '--help' option.

USAGE:
   {{.Name}} command [arguments...]

VERSION:
   {{.Version}}

COMMANDS:
   {{range .Commands}}{{if and (ne .Name "engine-discovery") (ne .Name "join")}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
   {{end}}{{end}}{{if .Flags}}
GLOBAL OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}
`
	cli.CommandHelpTemplate = `Docker Universal Control Plane Tool

   {{.Name}} - {{.Usage}}

USAGE: {{if or (or (eq .Name "dump-certs") (eq .Name "fingerprint")) (eq .Name "id") }}
   docker run --rm \
        --name ucp \
        -v /var/run/docker.sock:/var/run/docker.sock \
        docker/ucp \
        {{.Name}}{{if .Flags}} [command options]{{end}}
{{else if eq .Name "backup" }}
   docker run --rm -i \
        --name ucp \
        -v /var/run/docker.sock:/var/run/docker.sock \
        docker/ucp \
        {{.Name}}{{if .Flags}} [command options]{{end}} > backup.tar
{{else if eq .Name "restore" }}
   docker run --rm -i \
        --name ucp \
        -v /var/run/docker.sock:/var/run/docker.sock \
        docker/ucp \
        {{.Name}}{{if .Flags}} [command options]{{end}} < backup.tar
{{else if eq .Name "support" }}
   docker run --rm \
        --name ucp \
        -v /var/run/docker.sock:/var/run/docker.sock \
        docker/ucp \
        {{.Name}}{{if .Flags}} [command options]{{end}} > docker-support.tgz
{{else}}
   docker run --rm -it \
        --name ucp \
        -v /var/run/docker.sock:/var/run/docker.sock \
        docker/ucp \
        {{.Name}}{{if .Flags}} [command options]{{end}}
{{end}} {{if .Description}}
DESCRIPTION:
   {{.Description}}{{end}}{{if .Flags}}

OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{ end }}
`
	app := cli.NewApp()
	app.Name = `docker run --rm -it \
        --name ucp \
        -v /var/run/docker.sock:/var/run/docker.sock \
        docker/ucp \
        `
	app.Version = version.FullVersion()

	// Skip the first 2 ports (controller and swarm proxy).
	fixedPorts := make([]string, len(orcaconfig.RequiredPorts)-2)
	for i, port := range orcaconfig.RequiredPorts[2:] {
		fixedPorts[i] = fmt.Sprintf("%d", *port)
	}
	fixedPortsList := strings.Join(fixedPorts, ", ")

	app.Commands = []cli.Command{
		{
			Name:  "install",
			Usage: "Install UCP on this engine",
			Description: fmt.Sprintf(`
The 'install' command will install the UCP controller on the
local engine. If you intend to install a multi-node cluster,
you must open firewall ports between the engines for the
following ports:

    %d or the '--orca-port'
    %d or the '--swarm-port'
    %s
    4789(udp) and 7946(tcp/udp) for overlay networking

You can optionally use an externally generated and signed certificate
for the UCP controller by specifying '--external-server-cert'.  Create a storage
volume named 'ucp-controller-server-certs' with ca.pem, cert.pem, and key.pem
in the root directory before running the install.

A license file can optionally be injected during install by volume
mounting the file at '/docker_subscription.lic' in the tool.  E.g.,
-v /path/to/my%s:%s

`, orcaconfig.OrcaPort, orcaconfig.SwarmPort, fixedPortsList,
				config.LicenseFile, config.LicenseFile),
			Action: install.Run,
			Flags: buildFlags(commonFlags, commonInteractiveFlags, commonCredentialFlags, commonInstallFlags, commonImageFlags,
				[]cli.Flag{
					cli.BoolFlag{
						Name:  "swarm-experimental",
						Usage: "Enable experimental swarm features",
					},
					cli.BoolFlag{
						Name:  "disable-tracking",
						Usage: "Disable anonymous tracking and analytics",
					},
					cli.BoolFlag{
						Name:  "disable-usage",
						Usage: "Disable anonymous usage reporting",
					},
					cli.BoolFlag{ // Deprecated flag - remove eventually
						Name:   "external-ucp-ca",
						Hidden: true,
					},
					cli.BoolFlag{
						Name:  "external-server-cert",
						Usage: "Set up UCP Controller with an externally signed server certificate",
					},
					cli.BoolFlag{
						Name:  "preserve-certs",
						Usage: "Don't (re)generate certs on the host if existing ones are found",
					},
					cli.BoolFlag{
						Name:  "binpack",
						Usage: "Set Swarm scheduler to binpack mode (default spread)",
					},
					cli.BoolFlag{
						Name:  "random",
						Usage: "Set Swarm scheduler to random mode (default spread)",
					},
				}),
		},
		{
			Name:        "join",
			Description: "The join command is no longer used.  To join a node to UCP, simply run `docker swarm join ...`",
			Action: func(c *cli.Context) {
				log.Fatal("The join command is no longer used.  To join a node to UCP, simply run `docker swarm join ...`")
			},
			Flags: buildFlags(commonFlags, commonInteractiveFlags, commonCredentialFlags, commonInstallFlags, commonImageFlags,
				[]cli.Flag{
					cli.StringFlag{
						Name:   "url",
						Usage:  "The connection URL for the remote UCP controller",
						EnvVar: "UCP_URL",
						Hidden: true,
					},
					cli.BoolFlag{
						Name:   "insecure-fingerprint",
						Usage:  "Do not verify fingerprint of the remote UCP controller. Insecure on non-private networks.",
						Hidden: true,
					},
					cli.StringFlag{
						Name:   "fingerprint",
						Usage:  "The fingerprint of the UCP controller you trust",
						EnvVar: "UCP_FINGERPRINT",
						Hidden: true,
					},
					cli.BoolFlag{
						Name:   "replica",
						Usage:  "Configure this node to be a UCP controller replica",
						Hidden: true,
					},
					cli.BoolFlag{ // Deprecated flag - remove eventually
						Name:   "external-ucp-ca",
						Hidden: true,
					},
					cli.BoolFlag{
						Name:   "external-server-cert",
						Usage:  "(Replica only) Set up UCP Controller with an externally signed server certificate",
						Hidden: true,
					},
					cli.StringFlag{
						Name:   "passphrase",
						Usage:  "Decrypt the Root CA tar file with the provided passphrase",
						EnvVar: "UCP_PASSPHRASE",
						Hidden: true,
					},
				}),
		},
		{
			Name:   "restart",
			Usage:  "Start or restart UCP components on this engine",
			Action: restart.Run,
			Flags:  commonFlags,
		},
		{
			Name:   "stop",
			Usage:  "Stop UCP components running on this engine",
			Action: stop.Run,
			Flags:  commonFlags,
		},
		{
			Name:   "upgrade",
			Usage:  "Upgrade the UCP components on this engine",
			Action: upgrade.Run,
			Description: `
When upgrading UCP, you must run the 'upgrade' command against every
engine in your cluster.  You should upgrade your controller and replica
nodes first, followed by your compute nodes.  If you plan to upgrade your
engine as well, upgrade the engine first, before upgrading UCP on the node.

After upgrading each node, confirm the node is present in the UCP console
before proceeding to the next node.
`,
			Flags: buildFlags(commonFlags, commonInteractiveFlags, commonCredentialFlags, commonImageFlags,
				[]cli.Flag{
					cli.StringFlag{
						Name:  "id",
						Usage: "The ID of the UCP instance to upgrade",
					},
				}),
		},
		{
			Name:   "images",
			Usage:  "Verify the UCP images on this engine",
			Action: images.Run,
			Description: `
This command will verify all the required images used by UCP on the current engine.
By default, this will pull any missing images. Use the '--pull' argument to change
behavior.
`,
			Flags: buildFlags(commonFlags, commonImageFlags,
				[]cli.Flag{
					cli.BoolFlag{
						Name:  "list",
						Usage: "Don't do anything, just list the images used by UCP",
					},
				}),
		},
		{
			Name:  "uninstall",
			Usage: "Uninstall UCP components from this engine",
			Description: `
When uninstalling UCP, you must run the 'uninstall' command against every
engine in your cluster.
`,
			Action: uninstall.Run,
			Flags: buildFlags(commonFlags, commonInteractiveFlags, commonImageFlags,
				[]cli.Flag{
					cli.StringFlag{
						Name:  "id",
						Usage: "The ID of the UCP instance to uninstall",
					},
					cli.BoolFlag{
						Name:  "preserve-certs",
						Usage: "Don't delete the certs on the host",
					},
					cli.BoolFlag{
						Name:  "preserve-images",
						Usage: "Don't delete images on the host",
					},
				}),
		},
		{
			Name:  "dump-certs",
			Usage: "Dump out the public certs for this UCP controller",
			Description: `
This utility will dump out the public certificates for the UCP controller
running on the local engine.  This can then be used to populate local
certificate trust stores as desired.

When connecting UCP to DTR, use the output of '--cluster --ca' to
configure DTR.
`,
			Action: dump_certs.Run,
			Flags: buildFlags(commonFlags,
				[]cli.Flag{
					cli.BoolFlag{
						Name:  "ca",
						Usage: "Dump only the contents of the ca.pem file (default is to dump both ca and cert)",
					},
					cli.BoolFlag{
						Name:  "cluster",
						Usage: "Dump the internal UCP Cluster Root CA and cert instead of the public server cert",
					},
				}),
		},
		{
			Name:  "fingerprint",
			Usage: "Dump out the TLS fingerprint for the UCP controller running on this engine",
			Description: `
This utility will display the certificate fingerprint of the UCP controller
running on the local engine.  This can be used when scripting 'join'
operations for the '--fingerprint' flag.
`,
			Action: fingerprint.Run,
		},
		{
			Name:  "support",
			Usage: "Generate a support dump for this engine",
			Description: `
This utility will produce a support dump file on stdout for this local node.
`,
			Action: support.Run,
		},
		{
			Name:  "id",
			Usage: "Dump out the ID of the UCP components running on this engine",
			Description: `
This utility will display the ID of the local UCP components running
on this node.  This ID matches what you see when you run 'docker info'
pointed to the UCP controller(s) and is required by various commands
in this tool as confirmation.
`,
			Action: id.Run,
		},
		{
			Name:        "engine-discovery",
			Description: "The engine-discovery command is no longer used.  Overlay networking is enabled automatically via swarm-mode",
			Action: func(c *cli.Context) {
				log.Fatal("The engine-discovery command is no longer used.  Overlay networking is enabled automatically via swarm-mode")
			},
			// Keep the old flags around in case someone has a script - we want them to get the above message, not a "invalid flag" error
			Flags: buildFlags(commonFlags, []cli.Flag{
				cli.BoolFlag{
					Name:   "update",
					Usage:  "Apply engine discovery configuration changes",
					Hidden: true,
				},
				cli.StringSliceFlag{
					Name:   "controller",
					Usage:  "Update discovery with the external IP address or hostname of the controller(s)",
					Value:  &cli.StringSlice{},
					Hidden: true,
				},
				cli.StringFlag{
					Name:   "host-address",
					Usage:  "Update the external IP address this node advertises itself as",
					EnvVar: "UCP_HOST_ADDRESS",
					Hidden: true,
				},
				// Used to get the right version in the integration tests
				cli.StringFlag{
					Name:   "image-version",
					Value:  strings.Split(version.FullVersion(), " ")[0],
					Hidden: true,
				},
			}),
		},
		{
			Name:  "backup",
			Usage: "Stream a tar file to stdout containing all UCP data volumes",
			Description: `
This utility will dump out a tar file containing all the contents of the
volumes used by UCP on this controller.  This can be used to make periodic
backups suitable for use in the 'restore' command.  Only UCP infrastructure
containers are backed up by this tool.

When backing up an HA cluster, take backups of all controllers, one at a
time, in quick succession, and keep track of the exact time and sequence
when you performed each backup.  You will need this timestamp/sequence
information if you restore more than one controller together.

WARNING: During the backup, all UCP infrastructure containers will be
temporarily stopped on this controller to prevent data corruption.  No user
containers will be stopped during the backup.

WARNING: This backup will contain private keys and other sensitive information
and should be stored securely.  You may use the '--passphrase' flag to enable
built-in PGP compatible encryption.
`,
			Action: backup.Run,
			Flags: buildFlags(commonFlags, commonInteractiveFlags, []cli.Flag{
				cli.StringFlag{
					Name:  "id",
					Usage: "The ID of the UCP instance to backup",
				},
				cli.BoolFlag{
					Name:  "root-ca-only",
					Usage: "Backup only the root CA certificates and keys from this controller node",
				},
				cli.StringFlag{
					Name:   "passphrase",
					Usage:  "Encrypt the tar file with the provided passphrase",
					EnvVar: "UCP_PASSPHRASE",
				},
			}),
		},
		{
			Name:  "restore",
			Usage: "Restore a UCP cluster from a backup tar file.",
			Description: `
This utility will restore the state of this controller based on a tar file 
generated by the 'backup' command.  Any UCP containers that are running on this
host will be stopped prior to restoring the data. After the data is replaced, 
the containers will be restarted.

The full restore operation can be only performed on the same node where a backup 
was taken. Other controllers of the cluster need to be stopped before attempting 
a restore operation.

The restore operation can be also used to replicate Root CA material across 
controllers in the same cluster, by using the --root-ca-only flag. In this case,
no controllers need to be stopped and the same backup file can be used to 
replicate Root CA material across multiple controllers of the same cluster.

WARNING: Existing state on this node will be lost and replaced by the 
contents of the backup. 

When the restore operation is run using --interactive mode, the backup file 
needs to be mounted under /backup.tar within this container. If the 
--interactive flag is not set, the backup file  will be read from stdin.
`,
			Action: restore.Run,
			Flags: buildFlags(commonFlags, commonInteractiveFlags, []cli.Flag{
				cli.StringFlag{
					Name:  "id",
					Usage: "The ID of the UCP instance to backup",
				},
				cli.BoolFlag{
					Name:  "root-ca-only",
					Usage: "Restore only the Root CA certificates and keys on this controller node (leaving all other data intact)",
				},
				cli.StringFlag{
					Name:   "passphrase",
					Usage:  "Decrypt the tar file with the provided passphrase",
					EnvVar: "UCP_PASSPHRASE",
				},
			}),
		},
		{
			Name:  "regen-certs",
			Usage: "Regenerate keys and certificates for a UCP controller",
			Description: `
This utility will generate new private keys and certs for UCP controllers.

By default it will leave the Root CA keys and certs intact and only
regenerate server and client certs on the controller.  This can be used
to change the list of SANs within the certs after install and refresh
the expiration of the certificates.

You may regenerate the Root CAs with this tool using "--root-ca-only"
then follow a multi-step procedure to regenerate all certs in the cluster.

WARNING: REGENERATING THE ROOT CAs IS A DISRUPTIVE OPERATION!

First run "regen-certs --root-ca-only" on one controller.  If this is an
HA cluster, then perform a "backup --root-ca-only" on this controller,
and "restore --root-ca-only" on all other controllers.  Then on all of
the controllers run "regen-certs" during which the cluster will become
unavailable until 1/2+1 of the controllers are running with new certs.
Once all controllers have new certs, restart all the docker daemons on
the controller nodes.  Once the cluster controllers have recovered, run
"join --fresh-install" on all non-controller nodes to re-join them to
the cluster.  After completing the process, all user bundles will be
invalid and new bundles must be downloaded.
`,
			Action: regen_certs.Run,
			Flags: buildFlags(commonFlags, commonInteractiveFlags, []cli.Flag{
				// TODO - add commonCredentialFlags
				cli.BoolFlag{
					Name:  "root-ca-only",
					Usage: "Regenerate the Root CAs on this node (Do only once in an HA cluster!)",
				},
				cli.StringFlag{
					Name:  "id",
					Usage: "The ID of the UCP instance to regenerate certificates for",
				},
				cli.StringSliceFlag{
					Name:  "san",
					Usage: "Additional Subject Alternative Names for certs (e.g. --san foo1.bar.com --san foo2.bar.com)",
					Value: &cli.StringSlice{},
				},
				cli.BoolFlag{
					Name:  "external-server-cert",
					Usage: "Omit regenerating the UCP Controller web server certificate signed with an external CA",
				},
				// We have to run an address test, and need to be able to pass image version to use the right image on the system
				cli.StringFlag{
					Name:   "image-version",
					Value:  strings.Split(version.FullVersion(), " ")[0],
					Hidden: true,
				},
			}),
		},
	}
	// app.Flags = []cli.Flag{}

	defer func() {
		if r := recover(); r != nil {
			log.Error("Error: ", r)
			os.Exit(1)
		}
	}()

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}
