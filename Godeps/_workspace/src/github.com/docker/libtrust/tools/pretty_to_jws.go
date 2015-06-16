package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/docker/libtrust"
)

func main() {
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("Error reading stdin: %s", err)
	}

	sig, err := libtrust.ParsePrettySignature(input, "signatures")
	if err != nil {
		log.Fatalf("Error parsing pretty signature: %s", err)
	}

	jws, err := sig.JWS()
	if err != nil {
		log.Fatalf("Error creating JWS: %s", err)
	}

	os.Stdout.Write(jws)
}
