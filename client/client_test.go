package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/docker/distribution/registry/api/v2"
	"github.com/docker/notary/client/changelist"
	"github.com/docker/notary/cryptoservice"
	"github.com/docker/notary/server/handlers"
	"github.com/docker/notary/server/storage"
	"github.com/docker/notary/trustmanager"
	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/utils"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

const timestampKeyJSON = `{"keytype":"rsa","keyval":{"public":"MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAyyvBtTg2xzYS+MTTIBqSpI4V78tt8Yzqi7Jki/Z6NqjiDvcnbgcTqNR2t6B2W5NjGdp/hSaT2jyHM+kdmEGaPxg/zIuHbL3NIp4e0qwovWiEgACPIaELdn8O/kt5swsSKl1KMvLCH1sM86qMibNMAZ/hXOwd90TcHXCgZ91wHEAmsdjDC3dB0TT+FBgOac8RM01Y196QrZoOaDMTWh0EQfw7YbXAElhFVDFxBzDdYWbcIHSIogXQmq0CP+zaL/1WgcZZIClt2M6WCaxxF1S34wNn45gCvVZiZQ/iKWHerSr/2dGQeGo+7ezMSutRzvJ+01fInD86RS/CEtBCFZ1VyQIDAQAB","private":"MIIEpAIBAAKCAQEAyyvBtTg2xzYS+MTTIBqSpI4V78tt8Yzqi7Jki/Z6NqjiDvcnbgcTqNR2t6B2W5NjGdp/hSaT2jyHM+kdmEGaPxg/zIuHbL3NIp4e0qwovWiEgACPIaELdn8O/kt5swsSKl1KMvLCH1sM86qMibNMAZ/hXOwd90TcHXCgZ91wHEAmsdjDC3dB0TT+FBgOac8RM01Y196QrZoOaDMTWh0EQfw7YbXAElhFVDFxBzDdYWbcIHSIogXQmq0CP+zaL/1WgcZZIClt2M6WCaxxF1S34wNn45gCvVZiZQ/iKWHerSr/2dGQeGo+7ezMSutRzvJ+01fInD86RS/CEtBCFZ1VyQIDAQABAoIBAHar8FFxrE1gAGTeUpOF8fG8LIQMRwO4U6eVY7V9GpWiv6gOJTHXYFxU/aL0Ty3eQRxwy9tyVRo8EJz5pRex+e6ws1M+jLOviYqW4VocxQ8dZYd+zBvQfWmRfah7XXJ/HPUx2I05zrmR7VbGX6Bu4g5w3KnyIO61gfyQNKF2bm2Q3yblfupx3URvX0bl180R/+QN2Aslr4zxULFE6b+qJqBydrztq+AAP3WmskRxGa6irFnKxkspJqUpQN1mFselj6iQrzAcwkRPoCw0RwCCMq1/OOYvQtgxTJcO4zDVlbw54PvnxPZtcCWw7fO8oZ2Fvo2SDo75CDOATOGaT4Y9iqECgYEAzWZSpFbN9ZHmvq1lJQg//jFAyjsXRNn/nSvyLQILXltz6EHatImnXo3v+SivG91tfzBI1GfDvGUGaJpvKHoomB+qmhd8KIQhO5MBdAKZMf9fZqZofOPTD9xRXECCwdi+XqHBmL+l1OWz+O9Bh+Qobs2as/hQVgHaoXhQpE0NkTcCgYEA/Tjf6JBGl1+WxQDoGZDJrXoejzG9OFW19RjMdmPrg3t4fnbDtqTpZtCzXxPTCSeMrvplKbqAqZglWyq227ksKw4p7O6YfyhdtvC58oJmivlLr6sFaTsER7mDcYce8sQpqm+XQ8IPbnOk0Z1l6g56euTwTnew49uy25M6U1xL0P8CgYEAxEXv2Kw+OVhHV5PX4BBHHj6we88FiDyMfwM8cvfOJ0datekf9X7ImZkmZEAVPJpWBMD+B0J0jzU2b4SLjfFVkzBHVOH2Ob0xCH2MWPAWtekin7OKizUlPbW5ZV8b0+Kq30DQ/4a7D3rEhK8UPqeuX1tHZox1MAqrgbq3zJj4yvcCgYEAktYPKPm4pYCdmgFrlZ+bA0iEPf7Wvbsd91F5BtHsOOM5PQQ7e0bnvWIaEXEad/2CG9lBHlBy2WVLjDEZthILpa/h6e11ao8KwNGY0iKBuebT17rxOVMqqTjPGt8CuD2994IcEgOPFTpkAdUmyvG4XlkxbB8F6St17NPUB5DGuhsCgYA//Lfytk0FflXEeRQ16LT1YXgV7pcR2jsha4+4O5pxSFw/kTsOfJaYHg8StmROoyFnyE3sg76dCgLn0LENRCe5BvDhJnp5bMpQldG3XwcAxH8FGFNY4LtV/2ZKnJhxcONkfmzQPOmTyedOzrKQ+bNURsqLukCypP7/by6afBY4dA=="}}`
const timestampECDSAKeyJSON = `
{"keytype":"ecdsa","keyval":{"public":"MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEgl3rzMPMEKhS1k/AX16MM4PdidpjJr+z4pj0Td+30QnpbOIARgpyR1PiFztU8BZlqG3cUazvFclr2q/xHvfrqw==","private":"MHcCAQEEIDqtcdzU7H3AbIPSQaxHl9+xYECt7NpK7B1+6ep5cv9CoAoGCCqGSM49AwEHoUQDQgAEgl3rzMPMEKhS1k/AX16MM4PdidpjJr+z4pj0Td+30QnpbOIARgpyR1PiFztU8BZlqG3cUazvFclr2q/xHvfrqw=="}}`

func createTestServer(t *testing.T) (*httptest.Server, *http.ServeMux) {
	mux := http.NewServeMux()
	// TUF will request /v2/docker.com/notary/_trust/tuf/timestamp.key
	// Return a canned timestamp.key
	mux.HandleFunc("/v2/docker.com/notary/_trust/tuf/timestamp.key", func(w http.ResponseWriter, r *http.Request) {
		// Also contains the private key, but for the purpose of this
		// test, we don't care
		fmt.Fprint(w, timestampECDSAKeyJSON)
	})

	ts := httptest.NewServer(mux)

	return ts, mux
}

// TestInitRepo runs through the process of initializing a repository and makes
// sure the repository looks correct on disk.
// We test this with both an RSA and ECDSA root key
func TestInitRepo(t *testing.T) {
	testInitRepo(t, data.ECDSAKey)
	if !testing.Short() {
		testInitRepo(t, data.RSAKey)
	}
}

func testInitRepo(t *testing.T, rootType data.KeyAlgorithm) {
	gun := "docker.com/notary"
	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	defer os.RemoveAll(tempBaseDir)

	assert.NoError(t, err, "failed to create a temporary directory: %s", err)

	ts, _ := createTestServer(t)
	defer ts.Close()

	repo, err := NewNotaryRepository(tempBaseDir, gun, ts.URL, http.DefaultTransport, passphraseRetriever)
	assert.NoError(t, err, "error creating repo: %s", err)

	rootKeyID, err := repo.KeyStoreManager.GenRootKey(rootType.String())
	assert.NoError(t, err, "error generating root key: %s", err)

	rootCryptoService, err := repo.KeyStoreManager.GetRootCryptoService(rootKeyID)
	assert.NoError(t, err, "error retrieving root key: %s", err)

	err = repo.Initialize(rootCryptoService)
	assert.NoError(t, err, "error creating repository: %s", err)

	// Inspect contents of the temporary directory
	expectedDirs := []string{
		"private",
		filepath.Join("private", "tuf_keys", filepath.FromSlash(gun)),
		filepath.Join("private", "root_keys"),
		"trusted_certificates",
		filepath.Join("trusted_certificates", filepath.FromSlash(gun)),
		"tuf",
		filepath.Join("tuf", filepath.FromSlash(gun), "metadata"),
	}
	for _, dir := range expectedDirs {
		fi, err := os.Stat(filepath.Join(tempBaseDir, dir))
		assert.NoError(t, err, "missing directory in base directory: %s", dir)
		assert.True(t, fi.Mode().IsDir(), "%s is not a directory", dir)
	}

	// Look for keys in private. The filenames should match the key IDs
	// in the private key store.
	privKeyList := repo.KeyStoreManager.NonRootKeyStore().ListFiles()
	for _, privKeyName := range privKeyList {
		privKeyFileName := filepath.Join(repo.KeyStoreManager.NonRootKeyStore().BaseDir(), privKeyName)
		_, err := os.Stat(privKeyFileName)
		assert.NoError(t, err, "missing private key: %s", privKeyName)
	}

	// Look for keys in root_keys
	// There should be a file named after the key ID of the root key we
	// passed in.
	rootKeyFilename := rootCryptoService.ID() + "_root.key"
	_, err = os.Stat(filepath.Join(tempBaseDir, "private", "root_keys", rootKeyFilename))
	assert.NoError(t, err, "missing root key")

	certificates := repo.KeyStoreManager.TrustedCertificateStore().GetCertificates()
	assert.Len(t, certificates, 1, "unexpected number of certificates")

	certID, err := trustmanager.FingerprintCert(certificates[0])
	assert.NoError(t, err, "unable to fingerprint the certificate")

	// There should be a trusted certificate
	_, err = os.Stat(filepath.Join(tempBaseDir, "trusted_certificates", filepath.FromSlash(gun), certID+".crt"))
	assert.NoError(t, err, "missing trusted certificate")

	// Sanity check the TUF metadata files. Verify that they exist, the JSON is
	// well-formed, and the signatures exist. For the root.json file, also check
	// that the root, snapshot, and targets key IDs are present.
	expectedTUFMetadataFiles := []string{
		filepath.Join("tuf", filepath.FromSlash(gun), "metadata", "root.json"),
		filepath.Join("tuf", filepath.FromSlash(gun), "metadata", "snapshot.json"),
		filepath.Join("tuf", filepath.FromSlash(gun), "metadata", "targets.json"),
	}
	for _, filename := range expectedTUFMetadataFiles {
		fullPath := filepath.Join(tempBaseDir, filename)
		_, err := os.Stat(fullPath)
		assert.NoError(t, err, "missing TUF metadata file: %s", filename)

		jsonBytes, err := ioutil.ReadFile(fullPath)
		assert.NoError(t, err, "error reading TUF metadata file %s: %s", filename, err)

		var decoded data.Signed
		err = json.Unmarshal(jsonBytes, &decoded)
		assert.NoError(t, err, "error parsing TUF metadata file %s: %s", filename, err)

		assert.Len(t, decoded.Signatures, 1, "incorrect number of signatures in TUF metadata file %s", filename)

		assert.NotEmpty(t, decoded.Signatures[0].KeyID, "empty key ID field in TUF metadata file %s", filename)
		assert.NotEmpty(t, decoded.Signatures[0].Method, "empty method field in TUF metadata file %s", filename)
		assert.NotEmpty(t, decoded.Signatures[0].Signature, "empty signature in TUF metadata file %s", filename)

		// Special case for root.json: also check that the signed
		// content for keys and roles
		if strings.HasSuffix(filename, "root.json") {
			var decodedRoot data.Root
			err := json.Unmarshal(decoded.Signed, &decodedRoot)
			assert.NoError(t, err, "error parsing root.json signed section: %s", err)

			assert.Equal(t, "Root", decodedRoot.Type, "_type mismatch in root.json")

			// Expect 4 keys in the Keys map: root, targets, snapshot, timestamp
			assert.Len(t, decodedRoot.Keys, 4, "wrong number of keys in root.json")

			roleCount := 0
			for role := range decodedRoot.Roles {
				roleCount++
				if role != "root" && role != "snapshot" && role != "targets" && role != "timestamp" {
					t.Fatalf("unexpected role %s in root.json", role)
				}
			}
			assert.Equal(t, 4, roleCount, "wrong number of roles (%d) in root.json", roleCount)
		}
	}
}

// TestAddListTarget adds a target to the repo and confirms that the changelist
// is updated correctly. Then it calls ListTargets and checks the return value.
// Using ListTargets involves serving signed metadata files over the test's
// internal HTTP server.
// We test this with both an RSA and ECDSA root key
func TestAddListTarget(t *testing.T) {
	testAddListTarget(t, data.ECDSAKey)
	if !testing.Short() {
		testAddListTarget(t, data.RSAKey)
	}
}

func testAddListTarget(t *testing.T, rootType data.KeyAlgorithm) {
	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	defer os.RemoveAll(tempBaseDir)

	assert.NoError(t, err, "failed to create a temporary directory: %s", err)

	gun := "docker.com/notary"

	ts, mux := createTestServer(t)
	defer ts.Close()

	repo, err := NewNotaryRepository(tempBaseDir, gun, ts.URL, http.DefaultTransport, passphraseRetriever)
	assert.NoError(t, err, "error creating repository: %s", err)

	rootKeyID, err := repo.KeyStoreManager.GenRootKey(rootType.String())
	assert.NoError(t, err, "error generating root key: %s", err)

	rootCryptoService, err := repo.KeyStoreManager.GetRootCryptoService(rootKeyID)
	assert.NoError(t, err, "error retreiving root key: %s", err)

	err = repo.Initialize(rootCryptoService)
	assert.NoError(t, err, "error creating repository: %s", err)

	// tests need to manually boostrap timestamp as client doesn't generate it
	err = repo.tufRepo.InitTimestamp()
	assert.NoError(t, err, "error creating repository: %s", err)

	// Add fixtures/intermediate-ca.crt as a target. There's no particular reason
	// for using this file except that it happens to be available as
	// a fixture.
	latestTarget, err := NewTarget("latest", "../fixtures/intermediate-ca.crt")
	assert.NoError(t, err, "error creating target")
	err = repo.AddTarget(latestTarget)
	assert.NoError(t, err, "error adding target")

	// Look for the changelist file
	changelistDirPath := filepath.Join(tempBaseDir, "tuf", filepath.FromSlash(gun), "changelist")

	changelistDir, err := os.Open(changelistDirPath)
	assert.NoError(t, err, "could not open changelist directory")

	fileInfos, err := changelistDir.Readdir(0)
	assert.NoError(t, err, "could not read changelist directory")

	// Should only be one file in the directory
	assert.Len(t, fileInfos, 1, "wrong number of changelist files found")

	clName := fileInfos[0].Name()
	raw, err := ioutil.ReadFile(filepath.Join(changelistDirPath, clName))
	assert.NoError(t, err, "could not read changelist file %s", clName)

	c := &changelist.TufChange{}
	err = json.Unmarshal(raw, c)
	assert.NoError(t, err, "could not unmarshal changelist file %s", clName)

	assert.EqualValues(t, changelist.ActionCreate, c.Actn)
	assert.Equal(t, "targets", c.Role)
	assert.Equal(t, "target", c.ChangeType)
	assert.Equal(t, "latest", c.ChangePath)
	assert.NotEmpty(t, c.Data)

	changelistDir.Close()

	// Create a second target
	currentTarget, err := NewTarget("current", "../fixtures/intermediate-ca.crt")
	assert.NoError(t, err, "error creating target")
	err = repo.AddTarget(currentTarget)
	assert.NoError(t, err, "error adding target")

	changelistDir, err = os.Open(changelistDirPath)
	assert.NoError(t, err, "could not open changelist directory")

	// There should now be a second file in the directory
	fileInfos, err = changelistDir.Readdir(0)
	assert.NoError(t, err, "could not read changelist directory")

	assert.Len(t, fileInfos, 2, "wrong number of changelist files found")

	newFileFound := false
	for _, fileInfo := range fileInfos {
		if fileInfo.Name() != clName {
			clName2 := fileInfo.Name()
			raw, err := ioutil.ReadFile(filepath.Join(changelistDirPath, clName2))
			assert.NoError(t, err, "could not read changelist file %s", clName2)

			c := &changelist.TufChange{}
			err = json.Unmarshal(raw, c)
			assert.NoError(t, err, "could not unmarshal changelist file %s", clName2)

			assert.EqualValues(t, changelist.ActionCreate, c.Actn)
			assert.Equal(t, "targets", c.Role)
			assert.Equal(t, "target", c.ChangeType)
			assert.Equal(t, "current", c.ChangePath)
			assert.NotEmpty(t, c.Data)

			newFileFound = true
			break
		}
	}

	assert.True(t, newFileFound, "second changelist file not found")

	changelistDir.Close()

	// Now test ListTargets. In preparation, we need to expose some signed
	// metadata files on the internal HTTP server.

	// Apply the changelist. Normally, this would be done by Publish

	// load the changelist for this repo
	cl, err := changelist.NewFileChangelist(filepath.Join(tempBaseDir, "tuf", filepath.FromSlash(gun), "changelist"))
	assert.NoError(t, err, "could not open changelist")

	// apply the changelist to the repo
	err = applyChangelist(repo.tufRepo, cl)
	assert.NoError(t, err, "could not apply changelist")

	var tempKey data.TUFKey
	json.Unmarshal([]byte(timestampECDSAKeyJSON), &tempKey)

	repo.KeyStoreManager.NonRootKeyStore().AddKey(filepath.Join(filepath.FromSlash(gun), tempKey.ID()), "nonroot", &tempKey)

	// Because ListTargets will clear this
	savedTUFRepo := repo.tufRepo

	rootJSONFile := filepath.Join(tempBaseDir, "tuf", filepath.FromSlash(gun), "metadata", "root.json")
	rootFileBytes, err := ioutil.ReadFile(rootJSONFile)

	signedTargets, err := savedTUFRepo.SignTargets("targets", data.DefaultExpires("targets"), nil)
	assert.NoError(t, err)

	signedSnapshot, err := savedTUFRepo.SignSnapshot(data.DefaultExpires("snapshot"), nil)
	assert.NoError(t, err)

	signedTimestamp, err := savedTUFRepo.SignTimestamp(data.DefaultExpires("timestamp"), nil)
	assert.NoError(t, err)

	mux.HandleFunc("/v2/docker.com/notary/_trust/tuf/root.json", func(w http.ResponseWriter, r *http.Request) {
		assert.NoError(t, err)
		fmt.Fprint(w, string(rootFileBytes))
	})

	mux.HandleFunc("/v2/docker.com/notary/_trust/tuf/timestamp.json", func(w http.ResponseWriter, r *http.Request) {
		timestampJSON, _ := json.Marshal(signedTimestamp)
		fmt.Fprint(w, string(timestampJSON))
	})

	mux.HandleFunc("/v2/docker.com/notary/_trust/tuf/snapshot.json", func(w http.ResponseWriter, r *http.Request) {
		snapshotJSON, _ := json.Marshal(signedSnapshot)
		fmt.Fprint(w, string(snapshotJSON))
	})

	mux.HandleFunc("/v2/docker.com/notary/_trust/tuf/targets.json", func(w http.ResponseWriter, r *http.Request) {
		targetsJSON, _ := json.Marshal(signedTargets)
		fmt.Fprint(w, string(targetsJSON))
	})

	targets, err := repo.ListTargets()
	assert.NoError(t, err)

	// Should be two targets
	assert.Len(t, targets, 2, "unexpected number of targets returned by ListTargets")

	if targets[0].Name == "latest" {
		assert.Equal(t, latestTarget, targets[0], "latest target does not match")
		assert.Equal(t, currentTarget, targets[1], "current target does not match")
	} else if targets[0].Name == "current" {
		assert.Equal(t, currentTarget, targets[0], "current target does not match")
		assert.Equal(t, latestTarget, targets[1], "latest target does not match")
	} else {
		t.Fatalf("unexpected target name: %s", targets[0].Name)
	}

	// Also test GetTargetByName
	newLatestTarget, err := repo.GetTargetByName("latest")
	assert.NoError(t, err)
	assert.Equal(t, latestTarget, newLatestTarget, "latest target does not match")

	newCurrentTarget, err := repo.GetTargetByName("current")
	assert.NoError(t, err)
	assert.Equal(t, currentTarget, newCurrentTarget, "current target does not match")
}

// TestValidateRootKey verifies that the public data in root.json for the root
// key is a valid x509 certificate.
func TestValidateRootKey(t *testing.T) {
	testValidateRootKey(t, data.ECDSAKey)
	if !testing.Short() {
		testValidateRootKey(t, data.RSAKey)
	}
}

func testValidateRootKey(t *testing.T, rootType data.KeyAlgorithm) {
	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	defer os.RemoveAll(tempBaseDir)

	assert.NoError(t, err, "failed to create a temporary directory: %s", err)

	gun := "docker.com/notary"

	ts, _ := createTestServer(t)
	defer ts.Close()

	repo, err := NewNotaryRepository(tempBaseDir, gun, ts.URL, http.DefaultTransport, passphraseRetriever)
	assert.NoError(t, err, "error creating repository: %s", err)

	rootKeyID, err := repo.KeyStoreManager.GenRootKey(rootType.String())
	assert.NoError(t, err, "error generating root key: %s", err)

	rootCryptoService, err := repo.KeyStoreManager.GetRootCryptoService(rootKeyID)
	assert.NoError(t, err, "error retreiving root key: %s", err)

	err = repo.Initialize(rootCryptoService)
	assert.NoError(t, err, "error creating repository: %s", err)

	rootJSONFile := filepath.Join(tempBaseDir, "tuf", filepath.FromSlash(gun), "metadata", "root.json")

	jsonBytes, err := ioutil.ReadFile(rootJSONFile)
	assert.NoError(t, err, "error reading TUF metadata file %s: %s", rootJSONFile, err)

	var decoded data.Signed
	err = json.Unmarshal(jsonBytes, &decoded)
	assert.NoError(t, err, "error parsing TUF metadata file %s: %s", rootJSONFile, err)

	var decodedRoot data.Root
	err = json.Unmarshal(decoded.Signed, &decodedRoot)
	assert.NoError(t, err, "error parsing root.json signed section: %s", err)

	keyids := []string{}
	for role, roleData := range decodedRoot.Roles {
		if role == "root" {
			keyids = append(keyids, roleData.KeyIDs...)
		}
	}
	assert.NotEmpty(t, keyids)

	for _, keyid := range keyids {
		if key, ok := decodedRoot.Keys[keyid]; !ok {
			t.Fatal("key id not found in keys")
		} else {
			_, err := trustmanager.LoadCertFromPEM(key.Public())
			assert.NoError(t, err, "key is not a valid cert")
		}
	}
}

// TestPublish creates a repo, instantiates a notary server, and publishes
// the repo to the server.
// We test this with both an RSA and ECDSA root key
func TestPublish(t *testing.T) {
	testPublish(t, data.ECDSAKey)
	if !testing.Short() {
		testPublish(t, data.RSAKey)
	}
}

func testPublish(t *testing.T, rootType data.KeyAlgorithm) {
	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	defer os.RemoveAll(tempBaseDir)

	assert.NoError(t, err, "failed to create a temporary directory: %s", err)

	gun := "docker.com/notary"

	// Set up server
	ctx := context.WithValue(context.Background(), "metaStore", storage.NewMemStorage())

	// Do not pass one of the const KeyAlgorithms here as the value! Passing a
	// string is in itself good test that we are handling it correctly as we will
	// be receiving a string from the configuration.
	ctx = context.WithValue(ctx, "keyAlgorithm", "ecdsa")

	hand := utils.RootHandlerFactory(nil, ctx,
		cryptoservice.NewCryptoService("", trustmanager.NewKeyMemoryStore(passphraseRetriever)))

	r := mux.NewRouter()
	r.Methods("POST").Path("/v2/{imageName:" + v2.RepositoryNameRegexp.String() + "}/_trust/tuf/").Handler(hand(handlers.AtomicUpdateHandler, "push", "pull"))
	r.Methods("GET").Path("/v2/{imageName:" + v2.RepositoryNameRegexp.String() + "}/_trust/tuf/{tufRole:(root|targets|snapshot)}.json").Handler(hand(handlers.GetHandler, "pull"))
	r.Methods("GET").Path("/v2/{imageName:" + v2.RepositoryNameRegexp.String() + "}/_trust/tuf/timestamp.json").Handler(hand(handlers.GetTimestampHandler, "pull"))
	r.Methods("GET").Path("/v2/{imageName:" + v2.RepositoryNameRegexp.String() + "}/_trust/tuf/timestamp.key").Handler(hand(handlers.GetTimestampKeyHandler, "push", "pull"))
	//r.Methods("POST").Path("/v2/{imageName:" + server.RepositoryNameRegexp + "}/_trust/tuf/{tufRole:(root|targets|timestamp|snapshot)}.json").Handler(hand(handlers.UpdateHandler, "push", "pull"))
	r.Methods("DELETE").Path("/v2/{imageName:" + v2.RepositoryNameRegexp.String() + "}/_trust/tuf/").Handler(hand(handlers.DeleteHandler, "push", "pull"))

	ts := httptest.NewServer(r)

	repo, err := NewNotaryRepository(tempBaseDir, gun, ts.URL, http.DefaultTransport, passphraseRetriever)
	assert.NoError(t, err, "error creating repository: %s", err)

	rootKeyID, err := repo.KeyStoreManager.GenRootKey(rootType.String())
	assert.NoError(t, err, "error generating root key: %s", err)

	rootCryptoService, err := repo.KeyStoreManager.GetRootCryptoService(rootKeyID)
	assert.NoError(t, err, "error retreiving root key: %s", err)

	err = repo.Initialize(rootCryptoService)
	assert.NoError(t, err, "error creating repository: %s", err)

	// Add fixtures/intermediate-ca.crt as a target. There's no particular reason
	// for using this file except that it happens to be available as
	// a fixture.
	latestTarget, err := NewTarget("latest", "../fixtures/intermediate-ca.crt")
	assert.NoError(t, err, "error creating target")
	err = repo.AddTarget(latestTarget)
	assert.NoError(t, err, "error adding target")

	// Look for the changelist file
	changelistDirPath := filepath.Join(tempBaseDir, "tuf", filepath.FromSlash(gun), "changelist")

	changelistDir, err := os.Open(changelistDirPath)
	assert.NoError(t, err, "could not open changelist directory")

	fileInfos, err := changelistDir.Readdir(0)
	assert.NoError(t, err, "could not read changelist directory")

	// Should only be one file in the directory
	assert.Len(t, fileInfos, 1, "wrong number of changelist files found")

	clName := fileInfos[0].Name()
	raw, err := ioutil.ReadFile(filepath.Join(changelistDirPath, clName))
	assert.NoError(t, err, "could not read changelist file %s", clName)

	c := &changelist.TufChange{}
	err = json.Unmarshal(raw, c)
	assert.NoError(t, err, "could not unmarshal changelist file %s", clName)

	assert.EqualValues(t, changelist.ActionCreate, c.Actn)
	assert.Equal(t, "targets", c.Role)
	assert.Equal(t, "target", c.ChangeType)
	assert.Equal(t, "latest", c.ChangePath)
	assert.NotEmpty(t, c.Data)

	changelistDir.Close()

	// Create a second target
	currentTarget, err := NewTarget("current", "../fixtures/intermediate-ca.crt")
	assert.NoError(t, err, "error creating target")
	err = repo.AddTarget(currentTarget)
	assert.NoError(t, err, "error adding target")

	changelistDir, err = os.Open(changelistDirPath)
	assert.NoError(t, err, "could not open changelist directory")

	// There should now be a second file in the directory
	fileInfos, err = changelistDir.Readdir(0)
	assert.NoError(t, err, "could not read changelist directory")

	assert.Len(t, fileInfos, 2, "wrong number of changelist files found")

	newFileFound := false
	for _, fileInfo := range fileInfos {
		if fileInfo.Name() != clName {
			clName2 := fileInfo.Name()
			raw, err := ioutil.ReadFile(filepath.Join(changelistDirPath, clName2))
			assert.NoError(t, err, "could not read changelist file %s", clName2)

			c := &changelist.TufChange{}
			err = json.Unmarshal(raw, c)
			assert.NoError(t, err, "could not unmarshal changelist file %s", clName2)

			assert.EqualValues(t, changelist.ActionCreate, c.Actn)
			assert.Equal(t, "targets", c.Role)
			assert.Equal(t, "target", c.ChangeType)
			assert.Equal(t, "current", c.ChangePath)
			assert.NotEmpty(t, c.Data)

			newFileFound = true
			break
		}
	}

	assert.True(t, newFileFound, "second changelist file not found")

	changelistDir.Close()

	// Test loading changelist
	changes := make(map[string]changelist.Change)
	cl, err := repo.GetChangelist()
	assert.NoError(t, err, "could not get changelist for repo")

	assert.Len(t, cl.List(), 2, "Wrong number of changes returned from changelist")
	for _, ch := range cl.List() {
		changes[ch.Path()] = ch
	}

	currentChange := changes["current"]
	assert.NotNil(t, currentChange, "Expected changelist to contain a change for path 'current'")
	assert.EqualValues(t, changelist.ActionCreate, currentChange.Action())
	assert.Equal(t, "targets", currentChange.Scope())
	assert.Equal(t, "target", currentChange.Type())
	assert.Equal(t, "current", currentChange.Path())

	latestChange := changes["latest"]
	assert.NotNil(t, latestChange, "Expected changelist to contain a change for path 'latest'")
	assert.EqualValues(t, changelist.ActionCreate, latestChange.Action())
	assert.Equal(t, "targets", latestChange.Scope())
	assert.Equal(t, "target", latestChange.Type())
	assert.Equal(t, "latest", latestChange.Path())

	// Now test Publish
	err = repo.Publish()
	assert.NoError(t, err)

	changelistDir, err = os.Open(changelistDirPath)
	assert.NoError(t, err, "could not open changelist directory")
	fileInfos, err = changelistDir.Readdir(0)
	assert.NoError(t, err, "could not read changelist directory")
	// Should only be one file in the directory
	assert.Len(t, fileInfos, 0, "wrong number of changelist files found")

	// Create a new repo and pull from the server
	tempBaseDir2, err := ioutil.TempDir("", "notary-test-")
	defer os.RemoveAll(tempBaseDir2)

	assert.NoError(t, err, "failed to create a temporary directory: %s", err)

	repo2, err := NewNotaryRepository(tempBaseDir, gun, ts.URL, http.DefaultTransport, passphraseRetriever)
	assert.NoError(t, err, "error creating repository: %s", err)

	targets, err := repo2.ListTargets()
	assert.NoError(t, err)

	// Should be two targets
	assert.Len(t, targets, 2, "unexpected number of targets returned by ListTargets")

	if targets[0].Name == "latest" {
		assert.Equal(t, latestTarget, targets[0], "latest target does not match")
		assert.Equal(t, currentTarget, targets[1], "current target does not match")
	} else if targets[0].Name == "current" {
		assert.Equal(t, currentTarget, targets[0], "current target does not match")
		assert.Equal(t, latestTarget, targets[1], "latest target does not match")
	} else {
		t.Fatalf("unexpected target name: %s", targets[0].Name)
	}

	// Also test GetTargetByName
	newLatestTarget, err := repo2.GetTargetByName("latest")
	assert.NoError(t, err)
	assert.Equal(t, latestTarget, newLatestTarget, "latest target does not match")

	newCurrentTarget, err := repo2.GetTargetByName("current")
	assert.NoError(t, err)
	assert.Equal(t, currentTarget, newCurrentTarget, "current target does not match")
}

func TestRotate(t *testing.T) {
	// Temporary directory where test files will be created
	tempBaseDir, err := ioutil.TempDir("", "notary-test-")
	defer os.RemoveAll(tempBaseDir)

	assert.NoError(t, err, "failed to create a temporary directory: %s", err)

	gun := "docker.com/notary"

	// Set up server
	ctx := context.WithValue(context.Background(), "metaStore", storage.NewMemStorage())

	// Do not pass one of the const KeyAlgorithms here as the value! Passing a
	// string is in itself good test that we are handling it correctly as we will
	// be receiving a string from the configuration.
	ctx = context.WithValue(ctx, "keyAlgorithm", "ecdsa")

	hand := utils.RootHandlerFactory(nil, ctx,
		cryptoservice.NewCryptoService("", trustmanager.NewKeyMemoryStore(passphraseRetriever)))

	r := mux.NewRouter()
	r.Methods("POST").Path("/v2/{imageName:" + v2.RepositoryNameRegexp.String() + "}/_trust/tuf/").Handler(hand(handlers.AtomicUpdateHandler, "push", "pull"))
	r.Methods("GET").Path("/v2/{imageName:" + v2.RepositoryNameRegexp.String() + "}/_trust/tuf/{tufRole:(root|targets|snapshot)}.json").Handler(hand(handlers.GetHandler, "pull"))
	r.Methods("GET").Path("/v2/{imageName:" + v2.RepositoryNameRegexp.String() + "}/_trust/tuf/timestamp.json").Handler(hand(handlers.GetTimestampHandler, "pull"))
	r.Methods("GET").Path("/v2/{imageName:" + v2.RepositoryNameRegexp.String() + "}/_trust/tuf/timestamp.key").Handler(hand(handlers.GetTimestampKeyHandler, "push", "pull"))
	//r.Methods("POST").Path("/v2/{imageName:" + server.RepositoryNameRegexp + "}/_trust/tuf/{tufRole:(root|targets|timestamp|snapshot)}.json").Handler(hand(handlers.UpdateHandler, "push", "pull"))
	r.Methods("DELETE").Path("/v2/{imageName:" + v2.RepositoryNameRegexp.String() + "}/_trust/tuf/").Handler(hand(handlers.DeleteHandler, "push", "pull"))

	ts := httptest.NewServer(r)

	repo, err := NewNotaryRepository(tempBaseDir, gun, ts.URL, http.DefaultTransport, passphraseRetriever)
	assert.NoError(t, err, "error creating repository: %s", err)

	rootKeyID, err := repo.KeyStoreManager.GenRootKey(data.ECDSAKey.String())
	assert.NoError(t, err, "error generating root key: %s", err)

	rootCryptoService, err := repo.KeyStoreManager.GetRootCryptoService(rootKeyID)
	assert.NoError(t, err, "error retreiving root key: %s", err)

	err = repo.Initialize(rootCryptoService)
	assert.NoError(t, err, "error creating repository: %s", err)

	// Add fixtures/intermediate-ca.crt as a target. There's no particular reason
	// for using this file except that it happens to be available as
	// a fixture.
	// Adding a target will allow us to confirm the repository is still valid after
	// rotating the keys.
	latestTarget, err := NewTarget("latest", "../fixtures/intermediate-ca.crt")
	assert.NoError(t, err, "error creating target")
	err = repo.AddTarget(latestTarget)
	assert.NoError(t, err, "error adding target")

	// Publish
	err = repo.Publish()
	assert.NoError(t, err)

	// Get root.json and capture targets + snapshot key IDs
	repo.GetTargetByName("latest") // force a pull
	targetsKeyIDs := repo.tufRepo.Root.Signed.Roles["targets"].KeyIDs
	snapshotKeyIDs := repo.tufRepo.Root.Signed.Roles["snapshot"].KeyIDs
	assert.Len(t, targetsKeyIDs, 1)
	assert.Len(t, snapshotKeyIDs, 1)

	// Do rotation
	repo.RotateKeys()

	// Publish
	err = repo.Publish()
	assert.NoError(t, err)

	// Get root.json. Check targets + snapshot keys have changed
	// and that they match those found in the changelist.
	_, err = repo.GetTargetByName("latest") // force a pull
	assert.NoError(t, err)
	newTargetsKeyIDs := repo.tufRepo.Root.Signed.Roles["targets"].KeyIDs
	newSnapshotKeyIDs := repo.tufRepo.Root.Signed.Roles["snapshot"].KeyIDs
	assert.Len(t, newTargetsKeyIDs, 1)
	assert.Len(t, newSnapshotKeyIDs, 1)
	assert.NotEqual(t, targetsKeyIDs[0], newTargetsKeyIDs[0])
	assert.NotEqual(t, snapshotKeyIDs[0], newSnapshotKeyIDs[0])

	// Confirm changelist dir empty after publishing changes
	// Look for the changelist file
	changelistDirPath := filepath.Join(tempBaseDir, "tuf", filepath.FromSlash(gun), "changelist")
	changelistDir, err := os.Open(changelistDirPath)
	assert.NoError(t, err, "could not open changelist directory")
	fileInfos, err := changelistDir.Readdir(0)
	assert.NoError(t, err, "could not read changelist directory")
	// Should only be one file in the directory
	assert.Len(t, fileInfos, 0, "wrong number of changelist files found")
}
