package main

import (
	"os"

	"github.com/docker/dhe-deploy/moshpit/dtr"
	"github.com/docker/dhe-deploy/pkg/moshpit-framework/builder"
)

func main() {
	app := builder.BuildMoshpit(dtr.Setup, dtr.ClientRun)
	app.Run(os.Args)
}
