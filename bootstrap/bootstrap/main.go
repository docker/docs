package main

import (
	"os"

	"github.com/docker/dhe-deploy/bootstrap/app"
)

func main() {
	bootstapper := app.NewApp()
	bootstapper.Run(os.Args)
}
