package utils

import (
	"fmt"
	"io/ioutil"
	neturl "net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

type ComposeHandle struct {
	DirPath    string
	DockerHost string
}

func ComposePresent() bool {
	_, err := exec.LookPath("docker-machine")
	if err != nil {
		return false
	}
	return true
}

func (h *ComposeHandle) Cleanup() {
	// Consider preserving if PRESERVE_TEST_MACHINE is set
	if err := os.RemoveAll(h.DirPath); err != nil {
		log.Errorf("Failed to remove %s %s", h.DirPath, err)
	}
}

func (h *ComposeHandle) RunComposeCommand(args []string) (string, error) {
	log.Debugf("Running docker-compose %v", args)
	cmd := exec.Command("docker-compose", args...)
	cmd.Dir = h.DirPath
	cmd.Env = append(cmd.Env, fmt.Sprintf("DOCKER_CERT_PATH=%s", h.DirPath))
	cmd.Env = append(cmd.Env, fmt.Sprintf("DOCKER_HOST=tcp://%s:443", h.DockerHost))
	cmd.Env = append(cmd.Env, "DOCKER_TLS_VERIFY=1")
	log.Debugf("Running %#v", cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Error(string(out))
		return "", err
	}
	return string(out), nil
}

func NewComposeHandle(serverURL, username, password string) (*ComposeHandle, error) {
	zr, err := GetBundle(serverURL, username, password)
	if err != nil {
		return nil, err
	}

	// Make the temp dir and expand the cert bundle there
	dirname, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, err
	}
	u, err := neturl.Parse(serverURL)
	if err != nil {
		return nil, err
	}
	h := &ComposeHandle{
		DirPath:    dirname,
		DockerHost: strings.Split(u.Host, ":")[0],
	}

	log.Debugf("Writing bundle out to %s", h.DirPath)
	for _, zf := range zr.File {

		src, err := zf.Open()
		if err != nil {
			return nil, err
		}
		defer src.Close()
		var data []byte
		data, err = ioutil.ReadAll(src)
		if err != nil {
			return nil, err
		}
		if err := ioutil.WriteFile(filepath.Join(h.DirPath, zf.Name), data, 0600); err != nil {
			return nil, err
		}
	}

	return h, err
}

// The following routines are test routines that can be used in different contexts

func TestSimpleCompose(t *testing.T, serverURL, username, password string) {
	// Might want to soften and skip if this is a problem...
	require.True(t, ComposePresent())

	ch, err := NewComposeHandle(serverURL, username, password)
	require.Nil(t, err)
	defer ch.Cleanup()

	// Write out a sample compose file with a very simple app
	err = ioutil.WriteFile(filepath.Join(ch.DirPath, "docker-compose.yml"),
		[]byte(`
app:
    image: busybox
    command: sleep 24h
`), 0600)
	require.Nil(t, err)
	log.Debug("running compose")
	output, err := ch.RunComposeCommand([]string{"up", "-d"})
	require.Nil(t, err)
	log.Debug(output)

	// TODO - Might be nice to check that it's running, but it seems compose exits with error pretty well

	/* Might want to make this tunable for true idempotency
	output, err = ch.RunComposeCommand([]string{"stop"})
	require.Nil(t, err)
	log.Debug(output)

	output, err = ch.RunComposeCommand([]string{"rm", "-f"})
	require.Nil(t, err)
	log.Debug(output)
	*/
}

// TODO - would be nice to make a slightly more complex compose app..., maybe with a build too
