package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/orca/controller/commands"
	"github.com/docker/orca/version"
)

const (
	STORE_KEY = "orca"
)

func main() {
	app := cli.NewApp()
	app.Name = "orca"
	app.Usage = "docker orchestration"
	app.Version = version.FullVersion()
	app.Author = ""
	app.Email = ""

	// Always log with json format to make it easier for external tools to process
	log.SetFormatter(&log.JSONFormatter{})
	app.Before = func(c *cli.Context) error {
		if c.Bool("debug") {
			log.SetLevel(log.DebugLevel)
			commands.Debug = true
		}
		return nil
	}
	app.Commands = []cli.Command{
		commands.CmdServer,
	}
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, D",
			Usage: "enable debug",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

	// Server should run forever, if we get to here, something went wrong, so exit with failure
	os.Exit(1)
}
