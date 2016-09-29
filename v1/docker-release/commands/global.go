package commands

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"text/tabwriter"

	"github.com/Sirupsen/logrus"
)

// Print
func printData(templ string, data interface{}) {
	funcMap := template.FuncMap{
		"join": strings.Join,
	}

	w := tabwriter.NewWriter(os.Stdout, 1, 8, 2, ' ', 0)
	t := template.Must(template.New("help").Funcs(funcMap).Parse(templ))
	err := t.Execute(w, data)
	if err != nil {
		// If the writer is closed, t.Execute will fail, and there's nothing
		// we can do to recover.
		if os.Getenv("CLI_TEMPLATE_ERROR_DEBUG") != "" {
			logrus.Errorf("CLI TEMPLATE ERROR: %#v\n", err)
		}
		return
	}
	w.Flush()
}

func getReleasePath(addBuild bool) string {
	if addBuild == false {
		return archFlag + "/" + channelFlag + "/"
	}
	return archFlag + "/" + channelFlag + "/" + buildFlag + "/"
}

func genShasum(file string) (string, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	checksum := sha256.Sum256(data)

	// write shasum to tmp
	tmpFile, errDir := ioutil.TempFile("", "shasum")
	if errDir != nil {
		logrus.Error(errDir)
	}
	defer tmpFile.Close()
	shasum := hex.EncodeToString(checksum[:]) + "  " + path.Base(file) + "\n"
	logrus.Debug("Shasum 256: ", shasum)
	err = ioutil.WriteFile(tmpFile.Name(), []byte(shasum), 0644)
	if err != nil {
		logrus.Error("Could not write Shasum to tmp file")
		return "", errors.New("Could not write Shasum in tmp: " + tmpFile.Name())
	}
	return tmpFile.Name(), nil
}

func checkVersion() (version string, build string, err error) {
	re, _ := regexp.Compile(`^([0-9]+\.[0-9]+\.[0-9-a-z]+)\.([0-9]+)$`)
	result := re.FindStringSubmatch(buildFlag)
	if len(result) < 2 {
		return version, build, errors.New("Incorrect build flag")
	}
	version = result[1]
	build = result[2]
	if humanFlag != "" {
		version = humanFlag
	}
	logrus.Debugf("Version: %v - Build: %v", version, build)
	return version, build, nil
}

func checkArgs(args []string) error {
	if archFlag == "" || channelFlag == "" || buildFlag == "" {
		return errors.New("Release channel, archFlag or build number missing")
	}
	return nil
}

func readNotes(url string) ([]byte, error) {
	contents := []byte{}
	response, err := http.Get(url + "NOTES")
	if err != nil {
		logrus.Errorf("%s", err)
	} else {
		defer response.Body.Close()
		contents, err = ioutil.ReadAll(response.Body)
		if err != nil {
			logrus.Errorf("%s", err)
		}
		logrus.Debugf("%s\n", string(contents))
	}
	return contents, err
}

func getFileStat(filePath string) os.FileInfo {
	file, err := os.Open(filePath)
	if err != nil {
		logrus.Error("Failed to open file ", err)
		return nil
	}
	defer file.Close()
	fi, err := file.Stat()
	if err != nil {
		logrus.Error("Failed to stat opened file ", err)
		return nil
	}
	return fi
}
