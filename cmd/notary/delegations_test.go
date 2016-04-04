package main

import (
	"crypto/rand"
	"crypto/x509"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/docker/notary/cryptoservice"
	"github.com/docker/notary/trustmanager"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

var testTrustDir = "trust_dir"

func setup() *delegationCommander {
	return &delegationCommander{
		configGetter: func() (*viper.Viper, error) {
			mainViper := viper.New()
			mainViper.Set("trust_dir", testTrustDir)
			return mainViper, nil
		},
		retriever: nil,
	}
}

func TestAddInvalidDelegationName(t *testing.T) {
	// Cleanup after test
	defer os.RemoveAll(testTrustDir)

	// Setup certificate
	tempFile, err := ioutil.TempFile("/tmp", "pemfile")
	require.NoError(t, err)
	cert, _, err := generateValidTestCert()
	_, err = tempFile.Write(trustmanager.CertToPEM(cert))
	require.NoError(t, err)
	tempFile.Close()
	defer os.Remove(tempFile.Name())

	// Setup commander
	commander := setup()

	// Should error due to invalid delegation name (should be prefixed by "targets/")
	err = commander.delegationAdd(commander.GetCommand(), []string{"gun", "INVALID_NAME", tempFile.Name()})
	require.Error(t, err)
}

func TestAddInvalidDelegationCert(t *testing.T) {
	// Cleanup after test
	defer os.RemoveAll(testTrustDir)

	// Setup certificate
	tempFile, err := ioutil.TempFile("/tmp", "pemfile")
	require.NoError(t, err)
	cert, _, err := generateExpiredTestCert()
	_, err = tempFile.Write(trustmanager.CertToPEM(cert))
	require.NoError(t, err)
	tempFile.Close()
	defer os.Remove(tempFile.Name())

	// Setup commander
	commander := setup()

	// Should error due to expired cert
	err = commander.delegationAdd(commander.GetCommand(), []string{"gun", "targets/delegation", tempFile.Name(), "--paths", "path"})
	require.Error(t, err)
}

func TestAddInvalidShortPubkeyCert(t *testing.T) {
	// Cleanup after test
	defer os.RemoveAll(testTrustDir)

	// Setup certificate
	tempFile, err := ioutil.TempFile("/tmp", "pemfile")
	require.NoError(t, err)
	cert, _, err := generateShortRSAKeyTestCert()
	_, err = tempFile.Write(trustmanager.CertToPEM(cert))
	require.NoError(t, err)
	tempFile.Close()
	defer os.Remove(tempFile.Name())

	// Setup commander
	commander := setup()

	// Should error due to short RSA key
	err = commander.delegationAdd(commander.GetCommand(), []string{"gun", "targets/delegation", tempFile.Name(), "--paths", "path"})
	require.Error(t, err)
}

func TestRemoveInvalidDelegationName(t *testing.T) {
	// Cleanup after test
	defer os.RemoveAll(testTrustDir)

	// Setup commander
	commander := setup()

	// Should error due to invalid delegation name (should be prefixed by "targets/")
	err := commander.delegationRemove(commander.GetCommand(), []string{"gun", "INVALID_NAME", "fake_key_id1", "fake_key_id2"})
	require.Error(t, err)
}

func TestRemoveAllInvalidDelegationName(t *testing.T) {
	// Cleanup after test
	defer os.RemoveAll(testTrustDir)

	// Setup commander
	commander := setup()

	// Should error due to invalid delegation name (should be prefixed by "targets/")
	err := commander.delegationRemove(commander.GetCommand(), []string{"gun", "INVALID_NAME"})
	require.Error(t, err)
}

func TestAddInvalidNumArgs(t *testing.T) {
	// Setup commander
	commander := setup()

	// Should error due to invalid number of args (2 instead of 3)
	err := commander.delegationAdd(commander.GetCommand(), []string{"not", "enough"})
	require.Error(t, err)
}

func TestListInvalidNumArgs(t *testing.T) {
	// Setup commander
	commander := setup()

	// Should error due to invalid number of args (0 instead of 1)
	err := commander.delegationsList(commander.GetCommand(), []string{})
	require.Error(t, err)
}

func TestRemoveInvalidNumArgs(t *testing.T) {
	// Setup commander
	commander := setup()

	// Should error due to invalid number of args (1 instead of 2)
	err := commander.delegationRemove(commander.GetCommand(), []string{"notenough"})
	require.Error(t, err)
}

func generateValidTestCert() (*x509.Certificate, string, error) {
	privKey, err := trustmanager.GenerateECDSAKey(rand.Reader)
	if err != nil {
		return nil, "", err
	}
	keyID := privKey.ID()
	startTime := time.Now()
	endTime := startTime.AddDate(10, 0, 0)
	cert, err := cryptoservice.GenerateCertificate(privKey, "gun", startTime, endTime)
	if err != nil {
		return nil, "", err
	}
	return cert, keyID, nil
}

func generateExpiredTestCert() (*x509.Certificate, string, error) {
	privKey, err := trustmanager.GenerateECDSAKey(rand.Reader)
	if err != nil {
		return nil, "", err
	}
	keyID := privKey.ID()
	// Set to Unix time 0 start time, valid for one more day
	startTime := time.Unix(0, 0)
	endTime := startTime.AddDate(0, 0, 1)
	cert, err := cryptoservice.GenerateCertificate(privKey, "gun", startTime, endTime)
	if err != nil {
		return nil, "", err
	}
	return cert, keyID, nil
}

func generateShortRSAKeyTestCert() (*x509.Certificate, string, error) {
	// 1024 bits is too short
	privKey, err := trustmanager.GenerateRSAKey(rand.Reader, 1024)
	if err != nil {
		return nil, "", err
	}
	keyID := privKey.ID()
	startTime := time.Now()
	endTime := startTime.AddDate(10, 0, 0)
	cert, err := cryptoservice.GenerateCertificate(privKey, "gun", startTime, endTime)
	if err != nil {
		return nil, "", err
	}
	return cert, keyID, nil
}
