package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/docker/libtrust"
)

var ca string
var chain string
var cert string
var signKey string
var validKeys string

func init() {
	flag.StringVar(&ca, "ca", "", "Certificate authorities (pem file)")
	flag.StringVar(&chain, "chain", "", "Certificate chain to include (pem file)")
	flag.StringVar(&cert, "cert", "", "Certificate used to sign")
	flag.StringVar(&signKey, "k", "", "Private key to use for signing (pem or JWS file)")
	flag.StringVar(&validKeys, "keys", "", "File containing list of valid keys and namespaces")
}

func LoadValidKeys(filename string) (map[string][]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	validKeys := make(map[string][]string)
	r := bufio.NewScanner(f)
	for r.Scan() {
		parts := strings.Split(r.Text(), " ")
		if len(parts) < 2 {
			return nil, errors.New("Invalid line input: expecting <KeyId> <namespace> ...")
		}
		validKeys[parts[0]] = parts[1:]
	}
	return validKeys, nil
}

func main() {
	flag.Parse()
	if ca == "" {
		log.Fatal("Missing ca")
	}
	if chain == "" {
		log.Fatalf("Missing chain")
	}
	if cert == "" {
		log.Fatalf("Missing certificate")
	}
	if signKey == "" {
		log.Fatalf("Missing key")
	}
	if validKeys == "" {
		log.Fatalf("Missing valid keys")
	}

	caPool, err := libtrust.LoadCertificatePool(ca)
	if err != nil {
		log.Fatalf("Error loading ca certs: %s", err)
	}

	chainCerts, err := libtrust.LoadCertificateBundle(chain)
	if err != nil {
		log.Fatalf("Error loading chain certificates; %s", err)
	}
	chainPool := x509.NewCertPool()
	for _, cert := range chainCerts {
		chainPool.AddCert(cert)
	}

	signCert, err := tls.LoadX509KeyPair(cert, signKey)
	if err != nil {
		log.Fatalf("Error loading key: %s", err)
	}
	if signCert.Certificate == nil {
		log.Fatalf("Signed Cert is empty")
	}

	validKeyMap, err := LoadValidKeys(validKeys)
	if err != nil {
		log.Fatalf("Error loading valid keys: %s", err)
	}

	verifyOptions := x509.VerifyOptions{
		Intermediates: chainPool,
		Roots:         caPool,
	}

	parsedCert, err := x509.ParseCertificate(signCert.Certificate[0])
	if err != nil {
		log.Fatalf("Error parsing certificate: %s", err)
	}

	chains, err := parsedCert.Verify(verifyOptions)
	if err != nil {
		log.Fatalf("Error verifying certificate: %s", err)
	}
	if len(chains) == 0 {
		log.Fatalf("No verified chains")
	}

	content, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("Error reading content from stdin: %s", err)
	}

	sig, err := libtrust.ParsePrettySignature(content, "signatures")

	buildKeys, err := sig.Verify()
	if err != nil {
		log.Fatalf("Error verifying signatures: %s", err)
	}

	type manifest struct {
		Name string `json:"name"`
		Tag  string `json:"tag"`
	}
	var buildManifest manifest
	payload, err := sig.Payload()
	if err != nil {
		log.Fatalf("Error retrieving payload: %s", err)
	}
	err = json.Unmarshal(payload, &buildManifest)
	if err != nil {
		log.Fatalf("Error unmarshalling build manifest: %s", err)
	}

	log.Printf("Build keys: %#v", buildKeys)
	// Check keys against list of valid keys
	var foundKey bool
	for _, key := range buildKeys {
		keyID := key.KeyID()
		log.Printf("Checking key id: %s", keyID)
		namespaces, ok := validKeyMap[keyID]
		if ok {
			for _, namespace := range namespaces {
				if namespace == "*" || strings.HasPrefix(buildManifest.Name, namespace) {
					foundKey = true
				}
			}
		}

	}

	if !foundKey {
		log.Fatalf("No valid key found for build")
	}

	verifiedSig, err := libtrust.NewJSONSignature(content)
	if err != nil {
		log.Fatalf("Error creating JSON signature: %s", err)
	}

	privKey, err := libtrust.FromCryptoPrivateKey(signCert.PrivateKey)
	if err != nil {
		log.Fatalf("Error converting priv key: %s", err)
	}
	signChain := make([]*x509.Certificate, 1, len(chainCerts)+1)
	signChain[0] = parsedCert
	err = verifiedSig.SignWithChain(privKey, append(signChain, chainCerts...))
	if err != nil {
		log.Fatalf("Error signing with chain: %s", err)
	}

	// Output signed content to stdout
	out, err := verifiedSig.PrettySignature("verifySignatures")
	if err != nil {
		log.Fatalf("Error formatting output: %s", err)
	}
	_, err = os.Stdout.Write(out)
	if err != nil {
		log.Fatalf("Error writing output: %s", err)
	}
}
