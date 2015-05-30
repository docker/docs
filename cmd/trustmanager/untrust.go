package main

import (
	"fmt"

	"github.com/codegangsta/cli"
)

var (
	commandUntrust = cli.Command{
		Name:        "untrust",
		Usage:       "remove trust from a specifice certificate authority",
		Description: "remove trust from a specifice certificate authority.",
		Action:      untrust,
	}
)

func untrust(ctx *cli.Context) {
	args := []string(ctx.Args())

	if len(args) < 1 {
		cli.ShowCommandHelp(ctx, ctx.Command.Name)
		errorf("must specify a SHA256 SubjectKeyID of the certificate")
	}

	cert, err := caStore.GetCertificateBySKID(args[0])
	if err != nil {
		errorf("certificate not found")
	}

	fmt.Printf("Removing: ")
	print_cert(cert)

	err = caStore.RemoveCert(cert)
	if err != nil {
		errorf("failed to remove certificate for Key Store")
	}
}
