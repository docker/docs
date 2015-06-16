package main

import (
	"code.google.com/p/go.crypto/ssh/agent"
	"log"
	"net"
	"os"
)

func main() {
	c, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	agent := agent.NewClient(c)

	log.Println("Listing agent keys")

	log.Println("SSH")
	keys, err := agent.List()
	if err != nil {
		log.Fatalf("Error listing keys: %s", err)
	}
	for _, k := range keys {
		log.Printf("Key: %s", k.String())
		log.Printf("Type: %s", k.Type())
	}

}
