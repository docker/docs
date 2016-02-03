package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/docker/notary/passphrase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// the default location for the config file is in ~/.notary/config.json - even if it doesn't exist.
func TestNotaryConfigFileDefault(t *testing.T) {
	commander := &notaryCommander{
		getRetriever: func() passphrase.Retriever { return passphrase.ConstantRetriever("pass") },
	}

	config, err := commander.parseConfig()
	assert.NoError(t, err)
	configFileUsed := config.ConfigFileUsed()
	assert.True(t, strings.HasSuffix(configFileUsed,
		filepath.Join(".notary", "config.json")), "Unknown config file: %s", configFileUsed)
}

// the default server address is notary-server
func TestRemoteServerDefault(t *testing.T) {
	tempDir := tempDirWithConfig(t, "{}")
	defer os.RemoveAll(tempDir)
	configFile := filepath.Join(tempDir, "config.json")

	commander := &notaryCommander{
		getRetriever: func() passphrase.Retriever { return passphrase.ConstantRetriever("pass") },
	}

	// set a blank config file, so it doesn't check ~/.notary/config.json by default
	// and execute a random command so that the flags are parsed
	cmd := commander.GetCommand()
	cmd.SetArgs([]string{"-c", configFile, "list"})
	cmd.SetOutput(new(bytes.Buffer)) // eat the output
	cmd.Execute()

	config, err := commander.parseConfig()
	assert.NoError(t, err)
	assert.Equal(t, "https://notary-server:4443", getRemoteTrustServer(config))
}

// providing a config file uses the config file's server url instead
func TestRemoteServerUsesConfigFile(t *testing.T) {
	tempDir := tempDirWithConfig(t, `{"remote_server": {"url": "https://myserver"}}`)
	defer os.RemoveAll(tempDir)
	configFile := filepath.Join(tempDir, "config.json")

	commander := &notaryCommander{
		getRetriever: func() passphrase.Retriever { return passphrase.ConstantRetriever("pass") },
	}

	// set a config file, so it doesn't check ~/.notary/config.json by default,
	// and execute a random command so that the flags are parsed
	cmd := commander.GetCommand()
	cmd.SetArgs([]string{"-c", configFile, "list"})
	cmd.SetOutput(new(bytes.Buffer)) // eat the output
	cmd.Execute()

	config, err := commander.parseConfig()
	assert.NoError(t, err)
	assert.Equal(t, "https://myserver", getRemoteTrustServer(config))
}

// a command line flag overrides the config file's server url
func TestRemoteServerCommandLineFlagOverridesConfig(t *testing.T) {
	tempDir := tempDirWithConfig(t, `{"remote_server": {"url": "https://myserver"}}`)
	defer os.RemoveAll(tempDir)
	configFile := filepath.Join(tempDir, "config.json")

	commander := &notaryCommander{
		getRetriever: func() passphrase.Retriever { return passphrase.ConstantRetriever("pass") },
	}

	// set a config file, so it doesn't check ~/.notary/config.json by default,
	// and execute a random command so that the flags are parsed
	cmd := commander.GetCommand()
	cmd.SetArgs([]string{"-c", configFile, "-s", "http://overridden", "list"})
	cmd.SetOutput(new(bytes.Buffer)) // eat the output
	cmd.Execute()

	config, err := commander.parseConfig()
	assert.NoError(t, err)
	assert.Equal(t, "http://overridden", getRemoteTrustServer(config))
}

var exampleValidCommands = []string{
	"init repo",
	"list repo",
	"status repo",
	"publish repo",
	"add repo v1 somefile",
	"verify repo v1",
	"key list",
	"key rotate repo",
	"key generate rsa",
	"key backup tempfile.zip",
	"key export e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 backup.pem",
	"key restore tempfile.zip",
	"key import backup.pem",
	"key remove e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
	"key passwd e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
	"cert list",
	"cert remove e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
	"delegation list repo",
	"delegation add repo targets/releases path/to/pem/file.pem",
	"delegation remove repo targets/releases",
}

// config parsing bugs are propagated in all commands
func TestConfigParsingErrorsPropagatedByCommands(t *testing.T) {
	tempdir, err := ioutil.TempDir("", "empty-dir")
	require.NoError(t, err)
	defer os.RemoveAll(tempdir)

	for _, args := range exampleValidCommands {
		b := new(bytes.Buffer)
		cmd := NewNotaryCommand()
		cmd.SetOutput(b)

		cmd.SetArgs(append(
			[]string{"-c", filepath.Join(tempdir, "idonotexist.json"), "-d", tempdir},
			strings.Fields(args)...))
		err = cmd.Execute()

		require.Error(t, err, "expected error when running %s", args)
		require.Contains(t, err.Error(), "error opening config file", "running %s", args)
		require.NotContains(t, b.String(), "Usage:")
	}
}

// insufficient arguments produce an error before any parsing of configs happens
func TestInsufficientArgumentsReturnsErrorAndPrintsUsage(t *testing.T) {
	tempdir, err := ioutil.TempDir("", "empty-dir")
	require.NoError(t, err)
	defer os.RemoveAll(tempdir)

	for _, args := range exampleValidCommands {
		b := new(bytes.Buffer)
		cmd := NewNotaryCommand()
		cmd.SetOutput(b)

		arglist := strings.Fields(args)
		if args == "key list" || args == "cert list" || args == "key generate rsa" {
			// in these case, "key" or "cert" or "key generate" are valid commands, so add an arg to them instead
			arglist = append(arglist, "extraArg")
		} else {
			arglist = arglist[:len(arglist)-1]
		}

		invalid := strings.Join(arglist, " ")

		cmd.SetArgs(append(
			[]string{"-c", filepath.Join(tempdir, "idonotexist.json"), "-d", tempdir}, arglist...))
		err = cmd.Execute()

		require.NotContains(t, err.Error(), "error opening config file", "running %s", invalid)
		// it's a usage error, so the usage is printed
		require.Contains(t, b.String(), "Usage:", "expected usage when running %s", invalid)
	}
}
