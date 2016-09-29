package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

// read bytes from private key PEM storage
const privateKey string = `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----
`

func main() {
	var udid string
	// Get UDID from CLI
	if len(os.Args) > 1 {
		udid = os.Args[1]
	} else {
		udid = detectUDID()
	}

	logrus.Debugf("UDID found: ==>%s<===", udid)
	signatureBase64, err := signUdid(udid)
	if err != nil {
		logrus.Error(err)
	}
	logrus.Debugf(signatureBase64)
	logrus.Infof("Writing signature to %s/private-invites-signature.txt", appContainerPath())
	if err := writeSignatureInFile(signatureBase64); err != nil {
		logrus.Error(err)
	}
}

// signUdid signs the user udid according to the server private key
func signUdid(udid string) (string, error) {

	// create a signer object from the private key bytes
	signer, err := ssh.ParsePrivateKey([]byte(privateKey))
	// signer, err := ssh.ParsePrivateKey(privateKeyBytes)
	if err != nil {
		logrus.Error("unable to parse private key", err.Error())
	}
	signature, err := signer.Sign(nil, []byte(udid))
	if err != nil {
		logrus.Error("ERROR:", err.Error())
		return "", err
	}
	// encode signature into base64
	signatureBase64 := base64.StdEncoding.EncodeToString(signature.Blob)
	return signatureBase64, nil
}

func appContainerPath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(usr.HomeDir, "Library", "Containers", "com.docker.docker", "Data")
}

// encode signature into base64 and write it in a file
func writeSignatureInFile(signatureBase64 string) error {
	// get path to group container
	containerPath := appContainerPath()
	filepath := path.Join(containerPath, "private-invites-signature.txt")
	err := ioutil.WriteFile(filepath, []byte(signatureBase64), os.ModePerm)
	return err
}

var udidPrefix = []byte("Hardware UUID:")

func detectUDID() string {
	cmd := exec.Command("/usr/sbin/system_profiler", "SPHardwareDataType")
	lines := getLines(cmd)
	for _, line := range lines {
		line = bytes.TrimSpace(line)
		if bytes.HasPrefix(line, udidPrefix) {
			return string(bytes.TrimSpace(bytes.TrimPrefix(line, udidPrefix)))
		}
	}
	logrus.Fatal("udid not found; is this a Mac OSX?")
	return ""
}

func getLines(cmd *exec.Cmd) [][]byte {
	out := combinedOutput(cmd)
	return bytes.Split(out, []byte("\n"))
}
func combinedOutput(cmd *exec.Cmd) []byte {
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(strings.Join(cmd.Args, "\n"))
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return out
}
