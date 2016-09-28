package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/samalba/dockerclient"
	"github.com/stretchr/testify/require"
)

type FileHook struct {
	Writer *os.File
	Mutex  sync.Mutex
}

func (h *FileHook) Fire(entry *log.Entry) error {
	// Crude locking just to make sure we don't stomp on another parallel output
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	line, err := entry.String()
	if err != nil {
		// Should not happen, but don't log and create a loop
		return err
	}
	h.Writer.WriteString(line)
	// We could error on a short write, but meh...
	return nil
}

func (h *FileHook) Levels() []log.Level {
	return []log.Level{
		log.PanicLevel,
		log.FatalLevel,
		log.ErrorLevel,
		log.WarnLevel,
		log.InfoLevel,
		log.DebugLevel,
	}
}

// Wire up a log for all the tests
func init() {
	logDir := os.Getenv("LOG_DIR")
	if logDir == "" {
		logDir = os.Getenv("SRC_DIR")
	}
	if logDir == "" {
		logDir = "./"
	}
	logFile := os.Getenv("INTEGRATION_LOG_FILE")
	if logFile == "" {
		logFile = "integration.log"
	}
	logFilePath := filepath.Join(logDir, logFile)
	log.Infof("Setting up logging to also go to %s", logFilePath)
	fp, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Error("Failed to open log file! %s", err)
		return
	}
	h := FileHook{
		Writer: fp,
	}
	log.AddHook(&h)
}

// TODO - remove this after we ship 1.0
func ChangeLoggingLevelOld(t *testing.T, dclient *dockerclient.DockerClient, level string) {
	log.Debugf("Attempting to change logging level to %s", level)
	client := dclient.HTTPClient
	orcaURL := *dclient.URL
	orcaURL.Path = "/api/config/logging"
	data := []byte(fmt.Sprintf(`{"Level":"%s"}`, level)) // XXX This is the "old" config label
	req, err := http.NewRequest("POST", orcaURL.String(), bytes.NewBuffer(data))
	require.Nil(t, err)
	resp, err := client.Do(req)
	require.Nil(t, err)
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		require.FailNow(t, string(body))
	}
	log.Debugf("Succesfully set logging level")
}

func ChangeLoggingLevel(t *testing.T, dclient *dockerclient.DockerClient, level string) {
	log.Debugf("Attempting to change logging level to %s", level)
	client := dclient.HTTPClient
	orcaURL := *dclient.URL
	orcaURL.Path = "/api/config/logging"
	data := []byte(fmt.Sprintf(`{"level":"%s"}`, level))
	req, err := http.NewRequest("POST", orcaURL.String(), bytes.NewBuffer(data))
	require.Nil(t, err)
	resp, err := client.Do(req)
	require.Nil(t, err)
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		require.FailNow(t, string(body))
	}
	log.Debugf("Succesfully set logging level")
}

func GetLoggingConfig(t *testing.T, dclient *dockerclient.DockerClient) string {

	client := dclient.HTTPClient
	orcaURL := *dclient.URL
	orcaURL.Path = "/api/config/logging"
	req, err := http.NewRequest("GET", orcaURL.String(), nil)
	require.Nil(t, err)
	resp, err := client.Do(req)
	require.Nil(t, err)
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		require.FailNow(t, string(body))
	}
	return string(body)
}
