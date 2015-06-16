package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/docker/libtrust"
)

var ca string
var chain string
var cert string
var signKey string
var validKeys string

func init() {
	flag.StringVar(&signKey, "k", "", "Private key to use for signing (pem or JWS file)")
}

func main() {
	flag.Parse()
	if signKey == "" {
		log.Fatalf("Missing key")
	}

	content, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("Error reading content from stdin: %s", err)
	}

	sig, err := libtrust.NewJSONSignature(content)
	if err != nil {
		log.Fatalf("Error creating JSON signature: %s", err)
	}

	privKey, err := libtrust.LoadKeyFile(signKey)
	if err != nil {
		log.Fatalf("Error loading priv key: %s", err)
	}
	sig.Sign(privKey)

	// Output signed content to stdout
	out, err := sig.PrettySignature("signatures")
	if err != nil {
		log.Fatalf("Error formatting output: %s", err)
	}
	_, err = os.Stdout.Write(out)
	if err != nil {
		log.Fatalf("Error writing output: %s", err)
	}
}
