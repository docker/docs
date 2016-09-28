package main

import (
	"os"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/orca/integration/utils"
)

// Simple utility to stress out a system
func main() {
	log.SetLevel(log.DebugLevel)
	if len(os.Args) != 2 {
		log.Fatal("You must specify the number of containers to create on the command line")
	}
	count, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	client, err := utils.GetClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	// Get down to work...
	total, err := utils.CreateContainers(client, count, true)

	log.Infof("Created %d containers", total)
	if err != nil {
		log.Fatalf("Failure: %s", err)
	}
}
