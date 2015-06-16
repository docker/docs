package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/docker/libtrust"
)

func main() {
	keyFile, err := filepath.Abs("/home/derek/.docker/key.json")
	if err != nil {
		log.Fatalf("Error getting path: %s", err)
	}
	pk, err := libtrust.LoadKeyFile(keyFile)
	if err != nil {
		log.Fatalf("Error loading key file: %s", err)
	}

	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("Error reading stdin: %s", err)
	}

	sig, err := libtrust.NewJSONSignature(input)
	if err != nil {
		log.Fatalf("Error creating JSON signature: %s", err)
	}

	err = sig.Sign(pk)
	if err != nil {
		log.Fatalf("Error signing: %s", err)
	}
	//log.Printf("Private key (%s): %s", pk.KeyType(), pk.KeyID())
	jws, err := sig.JWS()
	if err != nil {
		log.Fatalf("Error getting JWS: %s", err)
	}

	os.Stdout.Write(jws)

}
