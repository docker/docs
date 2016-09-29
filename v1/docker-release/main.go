package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/docker/pinata/v1/docker-release/commands"
)

func main() {
	if c, err := commands.RootCmd.ExecuteC(); err != nil {
		logrus.Errorf("%s\n", err)
		c.Println(c.UsageString())
		// os.Exit(-1)
	}
}
