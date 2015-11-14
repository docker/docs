// Actually start up a notary server and run through basic TUF and key
// interactions via the client.

// Note - if using Yubikey, retrieving pins/touch doesn't seem to work right
// when running in the midst of all tests.

package main

import (
	"bytes"
	"crypto/rand"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/Sirupsen/logrus"
	ctxu "github.com/docker/distribution/context"
	"github.com/docker/notary/cryptoservice"
	"github.com/docker/notary/server"
	"github.com/docker/notary/server/storage"
	"github.com/docker/notary/trustmanager"
	"github.com/docker/notary/tuf/data"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

var testPassphrase = "passphrase"

// run a command and return the output as a string
func runCommand(t *testing.T, tempDir string, args ...string) (string, error) {
	// Using a new viper and Command so we don't have state between command invocations
	mainViper = viper.New()
	cmd := &cobra.Command{}
	setupCommand(cmd)

	b := new(bytes.Buffer)

	// Create an empty config file so we don't load the default on ~/.notary/config.json
	configFile := filepath.Join(tempDir, "config.json")

	cmd.SetArgs(append([]string{"-c", configFile, "-d", tempDir}, args...))
	cmd.SetOutput(b)
	retErr := cmd.Execute()
	output, err := ioutil.ReadAll(b)
	assert.NoError(t, err)

	return string(output), retErr
}

// makes a testing notary-server
func setupServer() *httptest.Server {
	// Set up server
	ctx := context.WithValue(
		context.Background(), "metaStore", storage.NewMemStorage())

	ctx = context.WithValue(ctx, "keyAlgorithm", data.ECDSAKey)

	// Eat the logs instead of spewing them out
	var b bytes.Buffer
	l := logrus.New()
	l.Out = &b
	ctx = ctxu.WithLogger(ctx, logrus.NewEntry(l))

	cryptoService := cryptoservice.NewCryptoService(
		"", trustmanager.NewKeyMemoryStore(retriever))
	return httptest.NewServer(server.RootHandler(nil, ctx, cryptoService))
}

// Initializes a repo, adds a target, publishes the target, lists the target,
// verifies the target, and then removes the target.
func TestClientTufInteraction(t *testing.T) {
	// -- setup --
	cleanup := setUp(t)
	defer cleanup()

	tempDir := tempDirWithConfig(t, "{}")
	defer os.RemoveAll(tempDir)

	server := setupServer()
	defer server.Close()

	tempFile, err := ioutil.TempFile("/tmp", "targetfile")
	assert.NoError(t, err)
	tempFile.Close()
	defer os.Remove(tempFile.Name())

	var (
		output string
		target = "sdgkadga"
	)
	// -- tests --

	// init repo
	_, err = runCommand(t, tempDir, "-s", server.URL, "init", "gun")
	assert.NoError(t, err)

	// add a target
	_, err = runCommand(t, tempDir, "add", "gun", target, tempFile.Name())
	assert.NoError(t, err)

	// check status - see target
	output, err = runCommand(t, tempDir, "status", "gun")
	assert.NoError(t, err)
	assert.True(t, strings.Contains(output, target))

	// publish repo
	_, err = runCommand(t, tempDir, "-s", server.URL, "publish", "gun")
	assert.NoError(t, err)

	// check status - no targets
	output, err = runCommand(t, tempDir, "status", "gun")
	assert.NoError(t, err)
	assert.False(t, strings.Contains(string(output), target))

	// list repo - see target
	output, err = runCommand(t, tempDir, "-s", server.URL, "list", "gun")
	assert.NoError(t, err)
	assert.True(t, strings.Contains(string(output), target))

	// lookup target and repo - see target
	output, err = runCommand(t, tempDir, "-s", server.URL, "lookup", "gun", target)
	assert.NoError(t, err)
	assert.True(t, strings.Contains(string(output), target))

	// verify repo - empty file
	output, err = runCommand(t, tempDir, "-s", server.URL, "verify", "gun", target)
	assert.NoError(t, err)

	// remove target
	_, err = runCommand(t, tempDir, "remove", "gun", target)
	assert.NoError(t, err)

	// publish repo
	_, err = runCommand(t, tempDir, "-s", server.URL, "publish", "gun")
	assert.NoError(t, err)

	// list repo - don't see target
	output, err = runCommand(t, tempDir, "-s", server.URL, "list", "gun")
	assert.NoError(t, err)
	assert.False(t, strings.Contains(string(output), target))
}

// Splits a string into lines, and returns any lines that are not empty (
// striped of whitespace)
func splitLines(chunk string) []string {
	splitted := strings.Split(strings.TrimSpace(chunk), "\n")
	var results []string

	for _, line := range splitted {
		line := strings.TrimSpace(line)
		if line != "" {
			results = append(results, line)
		}
	}

	return results
}

// List keys, parses the output, and returns the keys as an array of root key
// IDs and an array of signing key IDs
func GetKeys(t *testing.T, tempDir string) ([]string, []string) {
	output, err := runCommand(t, tempDir, "key", "list")
	assert.NoError(t, err)

	parts := strings.Split(output, "# Signing keys:")
	assert.Len(t, parts, 2)

	fixed := make([][]string, 2)
	for i, part := range parts {
		fixed[i] = splitLines(
			strings.TrimPrefix(strings.TrimSpace(part), "# Root keys:"))
		sort.Strings(fixed[i])
	}

	return fixed[0], fixed[1]
}

// List keys, parses the output, and asserts something about the number of root
// keys and number of signing keys, as well as returning them.
func assertNumKeys(t *testing.T, tempDir string, numRoot, numSigning int,
	rootOnDisk bool) ([]string, []string) {

	uniqueKeys := make(map[string]struct{})
	root, signing := GetKeys(t, tempDir)
	assert.Len(t, signing, numSigning)
	for i, rootKeyLine := range root {
		keyID := strings.Split(rootKeyLine, "-")[0]
		keyID = strings.TrimSpace(keyID)
		root[i] = keyID
		uniqueKeys[keyID] = struct{}{}
		_, err := os.Stat(filepath.Join(
			tempDir, "private", "root_keys", keyID+"_root.key"))
		// os.IsExist checks to see if the error is because a file already
		// exist, and hence doesn't actually the right funciton to use here
		assert.Equal(t, rootOnDisk, !os.IsNotExist(err))

		// this function is declared is in the build-tagged setup files
		verifyRootKeyOnHardware(t, keyID)
	}
	assert.Len(t, uniqueKeys, numRoot)
	return root, signing
}

// Adds the given target to the gun, publishes it, and lists it to ensure that
// it appears.  Returns the listing output.
func assertSuccessfullyPublish(
	t *testing.T, tempDir, url, gun, target, fname string) string {

	_, err := runCommand(t, tempDir, "add", gun, target, fname)
	assert.NoError(t, err)

	_, err = runCommand(t, tempDir, "-s", url, "publish", gun)
	assert.NoError(t, err)

	output, err := runCommand(t, tempDir, "-s", url, "list", gun)
	assert.NoError(t, err)
	assert.True(t, strings.Contains(string(output), target))

	return output
}

// Tests root key generation and key rotation
func TestClientKeyGenerationRotation(t *testing.T) {
	// -- setup --
	cleanup := setUp(t)
	defer cleanup()

	tempDir := tempDirWithConfig(t, "{}")
	defer os.RemoveAll(tempDir)

	tempfiles := make([]string, 2)
	for i := 0; i < 2; i++ {
		tempFile, err := ioutil.TempFile("/tmp", "targetfile")
		assert.NoError(t, err)
		tempFile.Close()
		tempfiles[i] = tempFile.Name()
		defer os.Remove(tempFile.Name())
	}

	server := setupServer()
	defer server.Close()

	var target = "sdgkadga"

	// -- tests --

	// starts out with no keys
	assertNumKeys(t, tempDir, 0, 0, true)

	// generate root key produces a single root key and no other keys
	_, err := runCommand(t, tempDir, "key", "generate", data.ECDSAKey)
	assert.NoError(t, err)
	assertNumKeys(t, tempDir, 1, 0, true)

	// initialize a repo, should have signing keys and no new root key
	_, err = runCommand(t, tempDir, "-s", server.URL, "init", "gun")
	assert.NoError(t, err)
	origRoot, origSign := assertNumKeys(t, tempDir, 1, 2, true)

	// publish using the original keys
	assertSuccessfullyPublish(t, tempDir, server.URL, "gun", target, tempfiles[0])

	// rotate the signing keys
	_, err = runCommand(t, tempDir, "key", "rotate", "gun")
	assert.NoError(t, err)
	root, sign := assertNumKeys(t, tempDir, 1, 4, true)
	assert.Equal(t, origRoot[0], root[0])
	// there should be the new keys and the old keys
	for _, origKey := range origSign {
		found := false
		for _, key := range sign {
			if key == origKey {
				found = true
			}
		}
		assert.True(t, found, "Old key not found in list of old and new keys")
	}

	// publish the key rotation
	_, err = runCommand(t, tempDir, "-s", server.URL, "publish", "gun")
	assert.NoError(t, err)
	root, sign = assertNumKeys(t, tempDir, 1, 2, true)
	assert.Equal(t, origRoot[0], root[0])
	// just do a cursory rotation check that the keys aren't equal anymore
	for _, origKey := range origSign {
		for _, key := range sign {
			assert.NotEqual(
				t, key, origKey, "One of the signing keys was not removed")
		}
	}

	// publish using the new keys
	output := assertSuccessfullyPublish(
		t, tempDir, server.URL, "gun", target+"2", tempfiles[1])
	// assert that the previous target is sitll there
	assert.True(t, strings.Contains(string(output), target))
}

// Tests backup/restore root+signing keys - repo with restored keys should be
// able to publish successfully
func TestClientKeyBackupAndRestore(t *testing.T) {
	// -- setup --
	cleanup := setUp(t)
	defer cleanup()

	dirs := make([]string, 3)
	for i := 0; i < 3; i++ {
		tempDir := tempDirWithConfig(t, "{}")
		defer os.RemoveAll(tempDir)
		dirs[i] = tempDir
	}

	tempfiles := make([]string, 2)
	for i := 0; i < 2; i++ {
		tempFile, err := ioutil.TempFile("/tmp", "tempfile")
		assert.NoError(t, err)
		tempFile.Close()
		tempfiles[i] = tempFile.Name()
		defer os.Remove(tempFile.Name())
	}

	server := setupServer()
	defer server.Close()

	var (
		target = "sdgkadga"
		err    error
	)

	// create two repos and publish a target
	for _, gun := range []string{"gun1", "gun2"} {
		_, err = runCommand(t, dirs[0], "-s", server.URL, "init", gun)
		assert.NoError(t, err)

		assertSuccessfullyPublish(
			t, dirs[0], server.URL, gun, target, tempfiles[0])
	}
	assertNumKeys(t, dirs[0], 1, 4, true)

	// -- tests --
	zipfile := tempfiles[0] + ".zip"
	defer os.Remove(zipfile)

	// backup then restore all keys
	_, err = runCommand(t, dirs[0], "key", "backup", zipfile)
	assert.NoError(t, err)

	_, err = runCommand(t, dirs[1], "key", "restore", zipfile)
	assert.NoError(t, err)
	assertNumKeys(t, dirs[1], 1, 4, !rootOnHardware()) // all keys should be there

	// can list and publish to both repos using restored keys
	for _, gun := range []string{"gun1", "gun2"} {
		output, err := runCommand(t, dirs[1], "-s", server.URL, "list", gun)
		assert.NoError(t, err)
		assert.True(t, strings.Contains(string(output), target))

		assertSuccessfullyPublish(
			t, dirs[1], server.URL, gun, target+"2", tempfiles[1])
	}

	// backup and restore keys for one gun
	_, err = runCommand(t, dirs[0], "key", "backup", zipfile, "-g", "gun1")
	assert.NoError(t, err)

	_, err = runCommand(t, dirs[2], "key", "restore", zipfile)
	assert.NoError(t, err)

	// this function is declared is in the build-tagged setup files
	if rootOnHardware() {
		// hardware root is still present, and the key will ONLY be on hardware
		// and not on disk
		assertNumKeys(t, dirs[2], 1, 2, false)
	} else {
		// only 2 signing keys should be there, and no root key
		assertNumKeys(t, dirs[2], 0, 2, true)
	}
}

// Generate a root key and export the root key only.  Return the key ID
// exported.
func exportRoot(t *testing.T, exportTo string) string {
	tempDir := tempDirWithConfig(t, "{}")
	defer os.RemoveAll(tempDir)

	// generate root key produces a single root key and no other keys
	_, err := runCommand(t, tempDir, "key", "generate", data.ECDSAKey)
	assert.NoError(t, err)
	oldRoot, _ := assertNumKeys(t, tempDir, 1, 0, true)

	// export does not require a password
	oldRetriever := retriever
	retriever = nil
	defer func() { // but import will, later
		retriever = oldRetriever
	}()

	_, err = runCommand(
		t, tempDir, "key", "export", oldRoot[0], exportTo)
	assert.NoError(t, err)

	return oldRoot[0]
}

// Tests import/export root key only
func TestClientKeyImportExportRootOnly(t *testing.T) {
	// -- setup --
	cleanup := setUp(t)
	defer cleanup()

	tempDir := tempDirWithConfig(t, "{}")
	defer os.RemoveAll(tempDir)

	server := setupServer()
	defer server.Close()

	var (
		target    = "sdgkadga"
		rootKeyID string
	)

	tempFile, err := ioutil.TempFile("/tmp", "pemfile")
	assert.NoError(t, err)
	// close later, because we might need to write to it
	defer os.Remove(tempFile.Name())

	// -- tests --

	if rootOnHardware() {
		t.Log("Cannot export a key from hardware. Will generate one to import.")

		privKey, err := trustmanager.GenerateECDSAKey(rand.Reader)
		assert.NoError(t, err)

		pemBytes, err := trustmanager.EncryptPrivateKey(privKey, testPassphrase)
		assert.NoError(t, err)

		nBytes, err := tempFile.Write(pemBytes)
		assert.NoError(t, err)
		tempFile.Close()
		assert.Equal(t, len(pemBytes), nBytes)
		rootKeyID = privKey.ID()
	} else {
		tempFile.Close()
		rootKeyID = exportRoot(t, tempFile.Name())
	}

	// import the key
	_, err = runCommand(t, tempDir, "key", "import", tempFile.Name())
	assert.NoError(t, err)

	// if there is hardware available, root will only be on hardware, and not
	// on disk
	newRoot, _ := assertNumKeys(t, tempDir, 1, 0, !rootOnHardware())
	assert.Equal(t, rootKeyID, newRoot[0])

	// Just to make sure, init a repo and publish
	_, err = runCommand(t, tempDir, "-s", server.URL, "init", "gun")
	assert.NoError(t, err)
	assertNumKeys(t, tempDir, 1, 2, !rootOnHardware())
	assertSuccessfullyPublish(
		t, tempDir, server.URL, "gun", target, tempFile.Name())
}

func assertNumCerts(t *testing.T, tempDir string, expectedNum int) []string {
	output, err := runCommand(t, tempDir, "cert", "list")
	assert.NoError(t, err)
	certs := splitLines(
		strings.TrimPrefix(strings.TrimSpace(output), "# Trusted Certificates:"))

	assert.Len(t, certs, expectedNum)
	return certs
}

// TestClientCertInteraction
func TestClientCertInteraction(t *testing.T) {
	// -- setup --
	cleanup := setUp(t)
	defer cleanup()

	tempDir := tempDirWithConfig(t, "{}")
	defer os.RemoveAll(tempDir)

	server := setupServer()
	defer server.Close()

	// -- tests --
	_, err := runCommand(t, tempDir, "-s", server.URL, "init", "gun1")
	assert.NoError(t, err)
	_, err = runCommand(t, tempDir, "-s", server.URL, "init", "gun2")
	assert.NoError(t, err)
	certs := assertNumCerts(t, tempDir, 2)

	// remove certs for one gun
	_, err = runCommand(t, tempDir, "cert", "remove", "-g", "gun1", "-y")
	assert.NoError(t, err)
	certs = assertNumCerts(t, tempDir, 1)

	// remove a single cert
	certID := strings.TrimSpace(strings.Split(certs[0], " ")[1])
	// passing an empty gun here because the string for the previous gun has
	// has already been stored (a drawback of running these commands without)
	// shelling out
	_, err = runCommand(t, tempDir, "cert", "remove", certID, "-y", "-g", "")
	assert.NoError(t, err)
	assertNumCerts(t, tempDir, 0)
}

// Tests default root key generation
func TestDefaultRootKeyGeneration(t *testing.T) {
	// -- setup --
	cleanup := setUp(t)
	defer cleanup()

	tempDir := tempDirWithConfig(t, "{}")
	defer os.RemoveAll(tempDir)

	// -- tests --

	// starts out with no keys
	assertNumKeys(t, tempDir, 0, 0, true)

	// generate root key with no algorithm produces a single ECDSA root key and no other keys
	_, err := runCommand(t, tempDir, "key", "generate")
	assert.NoError(t, err)
	assertNumKeys(t, tempDir, 1, 0, true)
}

func tempDirWithConfig(t *testing.T, config string) string {
	tempDir, err := ioutil.TempDir("/tmp", "repo")
	assert.NoError(t, err)
	err = ioutil.WriteFile(filepath.Join(tempDir, "config.json"), []byte(config), 0644)
	assert.NoError(t, err)
	return tempDir
}

func TestMain(m *testing.M) {
	if testing.Short() {
		// skip
		os.Exit(0)
	}
	os.Exit(m.Run())
}
