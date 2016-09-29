package commands

import (
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "upload new release artifacts",
	Run:   upload,
	PreRun: func(cmd *cobra.Command, args []string) {
		err := checkArgs(args)
		if err != nil {
			logrus.Fatal(err)
		}
	},
}

func upload(cmd *cobra.Command, args []string) {
	var symbols, installerFile string
	releasePath := getReleasePath(true)

	logrus.Debugf("Creating Path: %v", releasePath)
	pushData("", releasePath, true)
	for _, item := range args {
		s3Object := releasePath + path.Base(item)
		if strings.Contains(strings.ToUpper(item), "NOTE") {
			s3Object = releasePath + "NOTES"
		}
		pushData(item, s3Object, true)
		if strings.Contains(item, "Docker") {
			if strings.Contains(item, ".dmg") || strings.Contains(item, ".msi") {
				installerFile = item
			}
			if strings.Contains(strings.ToLower(item), ".dsym") {
				symbols = item
			}
			shasum, err := genShasum(item)
			if err != nil {
				logrus.Error(err)
			}
			s3Object = releasePath + path.Base(item) + ".sha256sum"
			pushData(shasum, s3Object, true)
		}
	}
	uploadMetadata(releasePath, installerFile, symbols)
}

func uploadMetadata(releasePath string, installerFile string, symbols string) {
	// Create MetaData for uploaded items
	version, build, err := checkVersion()
	if err != nil {
		logrus.Fatal(err)
	}
	tmpPath, errDir := ioutil.TempDir("", "xmlbuild")
	if errDir != nil {
		logrus.Error(errDir)
	}
	tmpMetadata := tmpPath + "/metadata.json"

	defer os.RemoveAll(tmpPath) // clean up

	installerInfo := getFileStat(installerFile)

	releaseMeta := metaData{
		Version:      version,
		Build:        build,
		Sha1:         sha1Flag,
		HumanVersion: version + " (build: " + build + ")",
		LastUpdate:   installerInfo.ModTime().Format(time.RFC3339),
		Size:         strconv.FormatInt(installerInfo.Size(), 10),
		Main:         s3URL + releasePath + path.Base(installerFile),
		Assets: []asset{
			{Signature: s3URL + releasePath + path.Base(installerFile) + ".sha256sum"},
			{Notes: s3URL + releasePath + "NOTES"},
		},
		Channel: channelFlag,
		Arch:    archFlag,
	}
	if symbols != "" {
		releaseMeta.Assets = append(releaseMeta.Assets, asset{Symbols: s3URL + getReleasePath(true) + path.Base(symbols)})
	}
	err = genMetaData(releaseMeta, tmpMetadata)
	if err != nil {
		logrus.Fatal(err)
	}
	pushData(tmpMetadata, releasePath+"metadata.json", true)
}
