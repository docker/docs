package main

import (
	"bytes"
	"github.com/docker/libtrust"
	"log"
)

func main() {
	ids, err := libtrust.ListSSHAgentIDs()
	if err != nil {
		log.Fatalf("Error listing ssh agent ids: %s", err)
	}

	for i := range ids {
		var id libtrust.ID
		id = ids[i]
		log.Printf("ID: %#v", id.JSONWebKey())

		signed, err := id.Sign(bytes.NewReader([]byte("hello there")))
		if err != nil {
			log.Fatalf("Error signing: %s", err)
		}

		log.Printf("Signed\n%x", signed)
	}
}
