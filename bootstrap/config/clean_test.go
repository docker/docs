package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestCleanupState(t *testing.T) {
	CertDir = fmt.Sprintf("/tmp/%d_test_dir", os.Getpid())

	if err := os.Mkdir(CertDir, 0700); err != nil {
		t.Errorf("Failed to setup tmp dir %s %s", CertDir, err)
	}
	defer os.RemoveAll(CertDir)

	file1 := filepath.Join(CertDir, CertFilename)
	if err := ioutil.WriteFile(file1, []byte("hello"), 0600); err != nil {
		t.Errorf("Failed to setup cert %s", err)
	}

	if err := CleanupState(CertDir); err != nil {
		t.Errorf("Failed to cleanup %s", err)
	}
	// Now make sure things are gone
	if _, err := os.Stat(file1); !os.IsNotExist(err) {
		t.Errorf("Expected file %s still exists: %s", file1, err)
	}
}
