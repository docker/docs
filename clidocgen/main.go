package main

import (
	"fmt"
	"os"
	"text/template"

	"github.com/docker/dhe-deploy/bootstrap/app"

	"github.com/codegangsta/cli"
)

const (
	overviewTemplateFile = "/templates/index.md.tmpl"
	commandTemplateFile  = "/templates/command.md.tmpl"
	outputDir            = "/output/"
)

func main() {
	dtrApp := app.NewApp()

	generateOverviewDocs(dtrApp, overviewTemplateFile, outputDir)
	generateCommandDocs(dtrApp.Commands, commandTemplateFile, outputDir)
}

func generateOverviewDocs(app *cli.App, templateFile string, outputDir string) {
	template, err := template.ParseFiles(templateFile)
	handleErr(err)

	file, err := os.Create(fmt.Sprintf("%s%s.md", outputDir, "index"))
	handleErr(err)
	defer file.Close()

	err = template.Execute(file, app)
	handleErr(err)
}

func generateCommandDocs(commands []cli.Command, templateFile string, outputDir string) {
	template, err := template.ParseFiles(templateFile)
	handleErr(err)

	for _, command := range commands {
		file, err := os.Create(fmt.Sprintf("%s%s.md", outputDir, command.Name))
		handleErr(err)
		defer file.Close()

		err = template.Execute(file, command)
		handleErr(err)
	}
}

func handleErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
