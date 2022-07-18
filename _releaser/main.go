package main

import (
	"log"

	"github.com/alecthomas/kong"
)

var (
	version = "dev"
	cli     struct {
		Version kong.VersionFlag
		Aws     AwsCmd     `kong:"cmd,name=aws"`
		Netlify NetlifyCmd `kong:"cmd,name=netlify"`
	}
)

func main() {
	log.SetFlags(0)
	ctx := kong.Parse(&cli,
		kong.Name("releaser"),
		kong.UsageOnError(),
		kong.Vars{
			"version": version,
		},
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: true,
		}))
	ctx.FatalIfErrorf(ctx.Run())
}
