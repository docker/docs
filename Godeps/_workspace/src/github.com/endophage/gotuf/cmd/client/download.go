package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/endophage/gotuf/keys"

	"github.com/endophage/gotuf"
	"github.com/endophage/gotuf/client"
	"github.com/endophage/gotuf/data"
	"github.com/endophage/gotuf/store"
)

var commandDownload = cli.Command{
	Name:   "download",
	Usage:  "provide the path to a target you wish to download.",
	Action: download,
}

func init() {
	commandDownload.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "root, r",
			Value:  "",
			Usage:  "The local file path to the current, or immediately previous, root.json for the TUF repo.",
			EnvVar: "TUF_ROOT",
		},
		cli.StringFlag{
			Name:   "host, h",
			Value:  "",
			Usage:  "The scheme and hostname of the TUF repository, e.g. http://example.com/.",
			EnvVar: "TUF_HOST",
		},
		cli.StringFlag{
			Name:   "meta, m",
			Value:  "",
			Usage:  "The path prefix for TUF metadata files.",
			EnvVar: "TUF_META_PREFIX",
		},
		cli.StringFlag{
			Name:   "ext, e",
			Value:  "json",
			Usage:  "The file extension for TUF metadata files.",
			EnvVar: "TUF_META_EXT",
		},
		cli.StringFlag{
			Name:   "targets, t",
			Value:  "",
			Usage:  "The path prefix for target files.",
			EnvVar: "TUF_TARGETS_PREFIX",
		},
	}
}

func download(ctx *cli.Context) {
	if len(ctx.Args()) < 1 {
		fmt.Println("At least one target name must be provided.")
		return
	}
	var root []byte
	r := &data.Signed{}
	err := json.Unmarshal(root, r)
	if err != nil {
		fmt.Println("Could not read initial root.json")
		return
	}
	kdb := keys.NewDB()
	repo := tuf.NewTufRepo(kdb, nil)
	repo.SetRoot(r)
	remote, err := store.NewHTTPStore(
		ctx.String("host"),
		ctx.String("meta"),
		ctx.String("ext"),
		ctx.String("targets"),
	)
	cached := store.NewFileCacheStore(remote, "/tmp/tuf")
	if err != nil {
		fmt.Println(err)
		return
	}

	client := client.NewClient(repo, cached, kdb)

	err = client.Update()
	if err != nil {
		fmt.Println(err)
		return
	}
	filename := filepath.Base(ctx.Args()[0])
	f, err := os.OpenFile(filename, os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	m := client.TargetMeta(ctx.Args()[0])
	if m == nil {
		fmt.Println("Requested package not found.")
		return
	}
	err = client.DownloadTarget(f, ctx.Args()[0], m)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Requested pacakge downloaded.")
}
