package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/docker/orca/enzi/cmd"
)

func main() {
	cmd.Worker.Subcommands = []cli.Command{
		cmd.TickJob,
		cmd.LDAPSyncJob,
		cmd.CleanupDBJob,
	}

	cmd.Root.Commands = []cli.Command{
		cmd.SyncDB,
		cmd.DrainDBServer,
		cmd.WaitForDB,
		cmd.CreateAdmin,
		cmd.Passwd,
		cmd.LDAPSearch,
		cmd.APIServer,
		cmd.Worker,
	}

	cmd.Root.Run(os.Args)
}
