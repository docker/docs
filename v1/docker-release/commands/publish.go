package commands

import (
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var version, build string

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "publish a new release",
	Run:   publish,
	PreRun: func(cmd *cobra.Command, args []string) {
		err := checkArgs(args)
		if err != nil {
			logrus.Fatal(err)
		}
		version, build, err = checkVersion()
		if err != nil {
			logrus.Fatal(err)
		}
		if len(args) > 0 {
			upload(cmd, args)
		}
	},
}

func publish(cmd *cobra.Command, args []string) {

	releasePath := getReleasePath(true)
	description, _ := readNotes(s3DomainURL)
	var artifact string
	switch archFlag {
	case "win":
		artifact = "InstallDocker.msi"
	case "mac":
		artifact = "Docker.dmg"
	}
	s3Artifact := getObject(releasePath + artifact)
	tmpPath, errDir := ioutil.TempDir("", "xmlbuild")
	if errDir != nil {
		logrus.Error(errDir)
	}
	tmpRelease := tmpPath + "/release.xml"

	defer os.RemoveAll(tmpPath) // clean up

	release := appXML{
		Version: "2.0",
		Sparkle: "http://www.andymatuschak.org/xml-namespaces/sparkle",
		DC:      "http://purl.org/dc/elements/1.1/",
		Channel: channelXML{
			Title: "Docker for " + archFlag,
			Link:  s3URL + getReleasePath(true) + artifact,
			Item: itemXML{
				Title: "Version " + version + " (" + build + ")", //Version 1.11.1-beta11 (build: 6974)
				Date:  s3Artifact.LastModified.Format(time.RFC3339),
				Enclosure: enclosureXML{
					URL:          s3URL + getReleasePath(true) + artifact,
					Version:      build,
					ShortVersion: version,
					Length:       strconv.FormatInt(*s3Artifact.ContentLength, 10),
					Type:         "application/octet-stream",
				},
			},
		},
	}
	err := genXML(release, tmpRelease, description)
	if err != nil {
		logrus.Fatal(err)
	}

	pushData(tmpRelease, getReleasePath(false)+"appcast.xml", true)
	createLatest(getReleasePath(true)+artifact, getReleasePath(false)+artifact)
	createLatest(getReleasePath(true)+artifact+".sha256sum", getReleasePath(false)+artifact+".sha256sum")
}
