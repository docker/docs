package main

import (
	"io"
	"io/ioutil"
	"os"

	tuf "github.com/endophage/gotuf/client"
	"github.com/endophage/gotuf/utils"
	"github.com/flynn/go-docopt"
)

func init() {
	register("get", cmdGet, `
usage: tuf-client get [-s|--store=<path>] <url> <target>

Options:
  -s <path>    The path to the local file store [default: tuf.db]

Get a target from the repository.
  `)
}

type tmpFile struct {
	*os.File
}

func (t *tmpFile) Delete() error {
	t.Close()
	return os.Remove(t.Name())
}

func cmdGet(args *docopt.Args, client *tuf.Client) error {
	if _, err := client.Update(); err != nil && !tuf.IsLatestSnapshot(err) {
		return err
	}
	target := utils.NormalizeTarget(args.String["<target>"])
	file, err := ioutil.TempFile("", "gotuf")
	if err != nil {
		return err
	}
	tmp := tmpFile{file}
	if err := client.Download(target, &tmp); err != nil {
		return err
	}
	defer tmp.Delete()
	if _, err := tmp.Seek(0, os.SEEK_SET); err != nil {
		return err
	}
	_, err = io.Copy(os.Stdout, file)
	return err
}
