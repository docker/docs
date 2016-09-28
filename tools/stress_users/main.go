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
		log.Fatal("You must specify the number of users to create on the command line")
	}
	count, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	client, err := utils.GetClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	serverURL := utils.GetServerURLFromClient(client)

	admin := os.Getenv("ORCA_ADMIN_USER")
	if admin == "" {
		admin = "admin"
	}
	password := os.Getenv("ORCA_ADMIN_PASSWORD")
	if password == "" {
		password = "orca"
	}

	// Get down to work...
	total, err := utils.AddUsers(client.HTTPClient, serverURL, count, "", "", admin, password, false, true)

	log.Infof("Created %d users", total)
	if err != nil {
		log.Fatalf("Failure: %s", err)
	}
}
