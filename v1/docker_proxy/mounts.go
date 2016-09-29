package proxy

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/docker/engine-api/types/container"
)

// MountRewriter rewrites mount paths
type MountRewriter interface {
	RewriteMounts(hc *container.HostConfig) error
}

// WindowsMountRewriter rewrites windows mount paths
type WindowsMountRewriter struct{}

// RewriteMounts rewrites Windows mount paths
func (w *WindowsMountRewriter) RewriteMounts(hc *container.HostConfig) error {
	if hc != nil {
		for i, spec := range hc.Binds {
			adjusted, err := adjustMount(spec, hc.VolumeDriver)
			if err != nil {
				return err
			}

			hc.Binds[i] = adjusted
		}
	}

	return nil
}

// adjustMount translates the mount paths on windows
func adjustMount(originalSpec string, volumeDriver string) (string, error) {
	spec := originalSpec

	// Replace C: with /C
	if len(spec) >= 2 && spec[1] == ':' {
		spec = fmt.Sprintf("/%c%s", spec[0], spec[2:])
	}

	// Extract source
	parts := strings.SplitN(spec, ":", 2)
	source := parts[0]
	if len(source) == 0 {
		return spec, nil
	}

	// Turn into a Linux path
	source = strings.Replace(source, `\`, `/`, -1)

	// Remove any .. links here
	source = path.Clean(source)

	// Rewrite source:dest spec
	spec = source
	if len(parts) == 2 {
		spec += ":" + parts[1]
	}

	if err := verifyDriveIsShared(source); err != nil {
		return "", err
	}

	log.Printf("Rewrote mount %s (volumeDriver=%s) to %s\n", originalSpec, volumeDriver, spec)

	return spec, nil
}

// verifyDriveIsShared checks that the source path can be accessed from the backend
func verifyDriveIsShared(source string) error {
	drive := extractDrive(source)
	if drive == "" {
		return nil
	}

	cliExe := findBackendCli()
	if cliExe == "" {
		log.Printf("Docker CLI not found. Cannot verify if the drive is shared. Let's say it is...")
		return nil
	}

	output, err := exec.Command(cliExe, "-SharedDrives").Output()
	if err != nil {
		return err
	}

	sharedDrives := strings.Split(strings.TrimSpace(string(output)), ",")
	for _, sharedDrive := range sharedDrives {
		log.Println(sharedDrive)
		if drive == sharedDrive {
			return nil
		}
	}

	return fmt.Errorf("%s: drive is not shared. Please share it in Docker for Windows Settings", drive)
}

// extractDrive extracts the (upper case) drive letter from a mount source path.
// Returns empty string if it doesn't look like a path to a shared drive.
func extractDrive(source string) string {
	if len(source) < 2 {
		return ""
	}

	if source[0] != '/' {
		return ""
	}

	if len(source) > 2 && source[2] != '/' {
		return ""
	}

	return strings.ToUpper(source[1:2])
}

// findBackendCli finds DockerCli.exe in the path. It first searches in the development path.
func findBackendCli() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}

	devCli := filepath.Join(dir, "..", "..", "build", "win", "DockerCli.exe")
	if _, err := os.Stat(devCli); err == nil {
		return devCli
	}

	prodCli := filepath.Join(dir, "..", "DockerCli.exe")
	if _, err := os.Stat(prodCli); err == nil {
		return prodCli
	}

	return ""
}
