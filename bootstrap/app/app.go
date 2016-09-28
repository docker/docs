package app

import (
	"fmt"

	"github.com/docker/dhe-deploy"
	"github.com/docker/dhe-deploy/bootstrap/commands"

	"github.com/codegangsta/cli"
)

const (
	// Changes the codegangsta default help to explain what DTR is
	AppHelpTemplate = `Docker Trusted Registry

{{if .Usage}}{{.Usage}}
{{end}}USAGE:
   {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} {{if .Flags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}
   {{if .Version}}{{if not .HideVersion}}
VERSION:
   {{.Version}}
   {{end}}{{end}}{{if len .Authors}}
AUTHOR(S):
   {{range .Authors}}{{ . }}{{end}}
   {{end}}{{if .Commands}}
COMMANDS:{{range .Categories}}{{if .Name}}
  {{.Name}}{{ ":" }}{{end}}{{range .Commands}}
    {{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}{{end}}
{{end}}{{end}}{{if .Flags}}
GLOBAL OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}{{end}}{{if .Copyright }}
COPYRIGHT:
   {{.Copyright}}
   {{end}}`

	// Description for how to run the app
	AppDescrition = `This tool has commands to install, configure, and backup Docker
Trusted Registry (DTR). It also allows uninstalling DTR.
By default the tool runs in interactive mode. It prompts you for
the values needed.

Additional help is available for each command with the '--help' option.`

	AppUsageExample = `docker run -it --rm docker/dtr \
    command [command options]`
)

func NewApp() *cli.App {
	app := cli.NewApp()
	app.Name = "dtr"
	app.HelpName = "dtr"
	app.Usage = AppDescrition
	app.Version = fmt.Sprintf("%s (%s)", deploy.ShortVersion, deploy.GitSHA)
	app.Author = "Docker Inc"
	app.UsageText = AppUsageExample
	app.HideVersion = true
	cli.AppHelpTemplate = AppHelpTemplate
	cli.CommandHelpTemplate = commands.CommandHelpTemplate

	app.Commands = commands.Commands

	return app
}
