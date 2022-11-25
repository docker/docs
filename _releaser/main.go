package main

import (
	"log"
	"os"
	"path/filepath"

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

// getEnvOrSecret retrieves secret's value from secret file or env
func getEnvOrSecret(name string) string {
	if v, ok := os.LookupEnv(name); ok {
		return v
	}
	b, err := os.ReadFile(filepath.Join("/run/secrets", name))
	if err != nil {
		return ""
	}
	return string(b)
}
