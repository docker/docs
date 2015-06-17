package main

import (
	"github.com/endophage/gotuf"
	"github.com/flynn/go-docopt"
)

func init() {
	register("clean", cmdClean, `
usage: tuf clean

Remove all staged manifests.
  `)
}

func cmdClean(args *docopt.Args, repo *tuf.Repo) error {
	return repo.Clean()
}
