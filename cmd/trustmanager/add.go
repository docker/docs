package main

import (
	"net/url"
	"os"

	"github.com/codegangsta/cli"
)

var (
	commandAdd = cli.Command{
		Name:        "add",
		Usage:       "Add an entry to the trusted certificate authority list.",
		Description: "Add an entry to the trusted certificate authority list.",
		Action:      add,
	}
)

func add(ctx *cli.Context) {
	args := []string(ctx.Args())

	if len(args) < 1 {
		cli.ShowCommandHelp(ctx, ctx.Command.Name)
		errorf("must specify a URL or file.")
	}

	// Verify if argument is a valid URL
	url, err := url.Parse(args[0])
	if err == nil && url.Scheme != "" {
		err = caStore.AddCertFromURL(args[0])
		if err != nil {
			errorf("error adding certificate to CA Store: %v", err)
		}
		// Verify is argument is a valid file
	} else if _, err := os.Stat(args[0]); err == nil {
		if err := caStore.AddCertFromFile(args[0]); err != nil {
			errorf("error adding certificate from file: %v", err)
		}
	} else {
		errorf("please provide a file location or URL for CA certificate.")
	}
}
