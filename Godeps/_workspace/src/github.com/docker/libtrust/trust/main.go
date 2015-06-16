package main

import (
	"fmt"
	"os"
	"path"
	"runtime"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

func main() {
	dir := path.Join(getHomeDir(), ".docker")
	app := cli.NewApp()
	app.Name = "trust"
	app.Usage = "manage keys and grants"
	app.Commands = []cli.Command{
		cli.Command{
			Name:   "tag",
			Usage:  "create or update a tag",
			Action: actionTag,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "commit",
					Usage: "Whether to immediately commit tag",
				},
			},
		},
		cli.Command{
			Name:   "untag",
			Usage:  "delete a tag",
			Action: actionUntag,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "ref,r",
					Value: "",
					Usage: "Tag or Hash reference",
				},
			},
		},
		cli.Command{
			Name:   "commit",
			Usage:  "commit target changes",
			Action: actionCommit,
		},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "dir,d",
			Value: dir,
			Usage: "Directory for repository",
		},
	}

	app.Run(os.Args)
}

func getHomeDir() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("USERPROFILE")
	}
	return os.Getenv("HOME")
}

func updateRepository() {
	log.Infof("Getting updated targets from repository")
	log.Infof("Verifying local keys match repository")
}

func actionTag(c *cli.Context) {
	hash := c.Args().First()
	if len(hash) == 0 {
		cli.ShowCommandHelp(c, "tag")
		return
	}

	tag := c.Args().Get(1)
	if len(tag) == 0 {
		cli.ShowCommandHelp(c, "tag")
		return
	}
	// TODO parse tag (Get last index of ':')

	updateRepository()
	log.Infof("Checking hash exists for name")
	log.Infof("Tagging %s as %s", hash, tag)
	fmt.Printf("%s+ %s %s%s\n", Green, tag, hash, Clear)
}

func actionUntag(c *cli.Context) {
	tagOrHash := c.Args().First()
	if len(tagOrHash) == 0 {
		cli.ShowCommandHelp(c, "untag")
		return
	}
	updateRepository()
	fmt.Printf("%s+ %s%s\n", Red, tagOrHash, Clear)
}

func actionCommit(c *cli.Context) {
}

const Black = "\x1b[30;1m"
const Red = "\x1b[31;1m"
const Green = "\x1b[32;1m"
const Yellow = "\x1b[33;1m"
const Blue = "\x1b[34;1m"
const Magenta = "\x1b[35;1m"
const Cyan = "\x1b[36;1m"
const White = "\x1b[37;1m"
const Clear = "\x1b[0m"
