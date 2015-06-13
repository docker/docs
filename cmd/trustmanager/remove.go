package main

import (
	"fmt"

	"github.com/codegangsta/cli"
)

var (
	commandRemove = cli.Command{
		Name:        "remove",
		Usage:       "remove trust from a specific certificate authority",
		Description: "remove trust from a specific certificate authority.",
		Action:      remove,
	}
)

func remove(ctx *cli.Context) {
	args := []string(ctx.Args())

	if len(args) < 1 {
		cli.ShowCommandHelp(ctx, ctx.Command.Name)
		fatalf("must specify a SHA256 SubjectKeyID of the certificate")
	}

	cert, err := caStore.GetCertificateBySKID(args[0])
	if err != nil {
		fatalf("certificate not found")
	}

	fmt.Printf("Removing: ")
	print_cert(cert)

	err = caStore.RemoveCert(cert)
	if err != nil {
		fatalf("failed to remove certificate for Key Store")
	}
}
