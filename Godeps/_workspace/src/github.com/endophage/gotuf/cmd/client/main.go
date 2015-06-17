package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "tufc"
	app.Usage = "tuf download <package name>"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config",
			Value: "config.json",
			Usage: "Set the path to a json configuration file for the TUF repo you want to interact with.",
		},
	}

	app.Commands = []cli.Command{
		commandDownload,
	}

	app.Run(os.Args)
}
