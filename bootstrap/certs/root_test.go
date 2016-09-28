package certs

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	cfssllog "github.com/cloudflare/cfssl/log"
	"github.com/docker/orca/bootstrap/config"
)

func TestInitCAPos(t *testing.T) {
	cfssllog.Level = cfssllog.LevelWarning
	certDir := fmt.Sprintf("/tmp/%d_test_dir", os.Getpid())
	expected := "----BEGIN"

	if err := os.Mkdir(certDir, 0700); err != nil {
		t.Errorf("Failed to setup tmp dir %s %s", config.CertDir, err)
	}
	defer os.RemoveAll(certDir)
	if err := InitCA("foo", certDir); err != nil {
		t.Errorf("Unexpected error %s", err)
	}

	// Check that the expected files got created
	for _, filename := range []string{filepath.Join(certDir, "cert.pem"), filepath.Join(certDir, "key.pem")} {

		if _, err := os.Stat(filename); os.IsNotExist(err) {
			t.Errorf("Expected file %s didn't exist", filename)
		} else {
			if output, err := ioutil.ReadFile(filename); err != nil {
				t.Errorf("Failed to read in file %s: %s", filename, err)
			} else if !strings.Contains(string(output), expected) {
				t.Errorf("File %s didn't contain %s", filename, expected)
			}
		}
	}
}
