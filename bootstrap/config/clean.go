package config

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
)

// Remove the certs and persistent files from the host
func CleanupState(mount string) error {
	files, err := filepath.Glob(fmt.Sprintf("%s/*", mount))
	if err != nil {
		return fmt.Errorf("Unable to remove old files at %s: %s", mount, err)
	}

	for _, name := range files {
		if err := os.RemoveAll(name); err != nil {
			log.Warnf("Failed to remove %s: %s", name, err)
			// XXX Should we bail?  If it's "extra cruft" we didn't put there, this
			//     might cause problems...
		}
	}
	return nil
}
