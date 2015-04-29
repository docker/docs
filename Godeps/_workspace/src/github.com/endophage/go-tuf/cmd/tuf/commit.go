package main

import (
	"github.com/endophage/go-tuf"
	"github.com/flynn/go-docopt"
)

func init() {
	register("commit", cmdCommit, `
usage: tuf commit

Commit staged files to the repository.
`)
}

func cmdCommit(args *docopt.Args, repo *tuf.Repo) error {
	return repo.Commit()
}
