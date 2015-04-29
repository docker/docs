package main

import (
	"log"

	"github.com/endophage/go-tuf"
	"github.com/flynn/go-docopt"
)

func init() {
	register("regenerate", cmdRegenerate, `
usage: tuf regenerate [--consistent-snapshot=false]

Recreate the targets manifest.
  `)
}

func cmdRegenerate(args *docopt.Args, repo *tuf.Repo) error {
	// TODO: implement this
	log.Println("not implemented")
	return nil
}
