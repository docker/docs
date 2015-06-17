package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/endophage/gotuf"
	"github.com/endophage/gotuf/data"
	"github.com/flynn/go-docopt"
)

func init() {
	register("meta", cmdMeta, `
usage: tuftools meta [<path>...]

Generate sample metadata for file(s) given by path.

`)
}

func cmdMeta(args *docopt.Args, repo *tuf.Repo) error {
	paths := args.All["<path>"].([]string)
	for _, file := range paths {
		reader, _ := os.Open(file)
		meta, _ := data.NewFileMeta(reader, "sha256")
		jsonBytes, err := json.Marshal(meta)
		if err != nil {
			return err
		}
		filename := fmt.Sprintf("%s.meta.json", file)
		err = ioutil.WriteFile(filename, jsonBytes, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}
