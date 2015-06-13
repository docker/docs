package main

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"fmt"
	"math"
	"time"

	"github.com/codegangsta/cli"
)

var (
	commandList = cli.Command{
		Name:        "list",
		Usage:       `List the currently trusted certificate authorities.`,
		Description: `List the currently trusted certificate authorities.`,
		Action:      list,
	}
)

func list(ctx *cli.Context) {
	// Load all the certificates
	trustedCAs := caStore.GetCertificates()

	fmt.Println("CAs Loaded:")
	for _, c := range trustedCAs {
		print_cert(c)
	}

}

func print_cert(cert *x509.Certificate) {
	timeDifference := cert.NotAfter.Sub(time.Now())
	fmt.Printf("Certificate: %s ; Expires in: %v days; SKID: %s\n", printPkix(cert.Subject), math.Floor(timeDifference.Hours()/24), hex.EncodeToString(cert.SubjectKeyId[:]))
}

func printPkix(pkixName pkix.Name) string {
	return fmt.Sprintf("%s - %s", pkixName.CommonName, pkixName.Organization)
}
