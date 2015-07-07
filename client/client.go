package client

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/docker/notary/trustmanager"
	"github.com/endophage/gotuf"
	tufclient "github.com/endophage/gotuf/client"
	"github.com/endophage/gotuf/data"
	"github.com/endophage/gotuf/keys"
	"github.com/endophage/gotuf/signed"
	"github.com/endophage/gotuf/store"

	"github.com/spf13/viper"
)

// Default paths should end with a '/' so directory creation works correctly
const trustDir string = "/trusted_certificates/"
const privDir string = "/private/"
const tufDir string = "/tuf/"
const rootKeysDir string = "/root_keys/"

// ErrRepositoryNotExist gets returned when trying to make an action over a repository
/// that doesn't exist
var ErrRepositoryNotExist = errors.New("repository does not exist")

// Client is the interface that defines the Notary Client type
type Client interface {
	ListPrivateKeys() []data.PrivateKey
	GenRootKey(passphrase string) (*data.PublicKey, error)
	GetRepository(gun string, baseURL string, transport http.RoundTripper) (Repository, error)
}

// Repository is the interface that represents a Notary Repository
type Repository interface {
	Update() error
	Initialize(key *data.PublicKey) error

	AddTarget(target *Target) error
	ListTargets() ([]*Target, error)
	GetTargetByName(name string) (*Target, error)

	Publish() error
}

type NotaryClient struct {
	caStore          trustmanager.X509Store
	certificateStore trustmanager.X509Store
	rootKeyStore     trustmanager.EncryptedFileStore
}

type NotaryRepository struct {
	Gun              string
	baseURL          string
	transport        http.RoundTripper
	signer           *signed.Signer
	tufRepo          *tuf.TufRepo
	fileStore        store.MetadataStore
	privKeyStore     trustmanager.FileStore
	caStore          trustmanager.X509Store
	certificateStore trustmanager.X509Store
}

// Target represents a simplified version of the data TUF operates on.
type Target struct {
	Name   string
	Hashes data.Hashes
	Length int64
}

// NewTarget  is a helper method that returns a Target
func NewTarget(targetName string, targetPath string) (*Target, error) {
	b, err := ioutil.ReadFile(targetPath)
	if err != nil {
		return nil, err
	}

	meta, err := data.NewFileMeta(bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	return &Target{Name: targetName, Hashes: meta.Hashes, Length: meta.Length}, nil
}

// NewClient is a helper method that returns a new notary Client, given a config
// file. It makes the assumption that the base directory for the config file will
// be the place where trust information is being cached locally.
func NewClient(trustDir, rootKeysDir string) (*NotaryClient, error) {
	nClient := &NotaryClient{}

	err := nClient.loadKeys(trustDir, rootKeysDir)
	if err != nil {
		return nil, err
	}

	return nClient, nil
}

// Update forces TUF to download the remote timestamps and verify if there are
// any remote changes.
func (r *NotaryRepository) Update() error {
	return nil
}

// Initialize creates a new repository by using rootKey as the root Key for the
// TUF repository.
func (r *NotaryRepository) Initialize(rootKey *data.PublicKey) error {
	remote, err := getRemoteStore(r.Gun)
	rawTSKey, err := remote.GetKey("timestamp")
	if err != nil {
		return err
	}
	fmt.Println("RawKey: ", string(rawTSKey))
	parsedKey := &data.TUFKey{}
	err = json.Unmarshal(rawTSKey, parsedKey)
	if err != nil {
		return err
	}
	timestampKey := data.NewPublicKey(parsedKey.Cipher(), parsedKey.Public())

	targetsKey, err := r.signer.Create("targets")
	if err != nil {
		return err
	}
	snapshotKey, err := r.signer.Create("snapshot")
	if err != nil {
		return err
	}

	kdb := keys.NewDB()

	kdb.AddKey(rootKey)
	kdb.AddKey(targetsKey)
	kdb.AddKey(snapshotKey)
	kdb.AddKey(timestampKey)

	rootRole, err := data.NewRole("root", 1, []string{rootKey.ID()}, nil, nil)
	if err != nil {
		return err
	}
	targetsRole, err := data.NewRole("targets", 1, []string{targetsKey.ID()}, nil, nil)
	if err != nil {
		return err
	}
	snapshotRole, err := data.NewRole("snapshot", 1, []string{snapshotKey.ID()}, nil, nil)
	if err != nil {
		return err
	}
	timestampRole, err := data.NewRole("timestamp", 1, []string{timestampKey.ID()}, nil, nil)
	if err != nil {
		return err
	}

	err = kdb.AddRole(rootRole)
	if err != nil {
		return err
	}
	err = kdb.AddRole(targetsRole)
	if err != nil {
		return err
	}
	err = kdb.AddRole(snapshotRole)
	if err != nil {
		return err
	}
	err = kdb.AddRole(timestampRole)
	if err != nil {
		return err
	}

	r.tufRepo = tuf.NewTufRepo(kdb, r.signer)

	r.fileStore, err = store.NewFilesystemStore(
		path.Join(viper.GetString("tufDir")),
		"metadata",
		"json",
		"targets",
	)
	if err != nil {
		return err
	}

	err = r.tufRepo.InitRepo(false)
	if err != nil {
		return err
	}

	r.saveRepo()
	return nil
}

// AddTarget adds a new target to the repository, forcing a timestamps check from TUF
func (r *NotaryRepository) AddTarget(target *Target) error {
	r.bootstrapRepo()

	fmt.Printf("Adding target \"%s\" with sha256 \"%s\" and size %d bytes.\n", target.Name, target.Hashes["sha256"], target.Length)

	meta := data.FileMeta{Length: target.Length, Hashes: target.Hashes}
	_, err := r.tufRepo.AddTargets("targets", data.Files{target.Name: meta})
	if err != nil {
		return err
	}

	r.saveRepo()

	return nil
}

// ListTargets lists all targets for the current repository
func (r *NotaryRepository) ListTargets() ([]*Target, error) {
	r.bootstrapRepo()

	c, err := r.bootstrapClient()
	if err != nil {
		return nil, err
	}

	err = c.Update()
	if err != nil {
		return nil, err
	}

	// TODO(diogo): return hashes
	for name, meta := range r.tufRepo.Targets["targets"].Signed.Targets {
		fmt.Println(name, " ", meta.Hashes["sha256"], " ", meta.Length)
	}

	return nil, nil
}

// GetTargetByName returns a target given a name
func (r *NotaryRepository) GetTargetByName(name string) (*Target, error) {
	r.bootstrapRepo()

	c, err := r.bootstrapClient()
	if err != nil {
		return nil, err
	}

	err = c.Update()
	if err != nil {
		return nil, err
	}

	meta := c.TargetMeta(name)
	if meta == nil {
		return nil, errors.New("Meta is nil for target")
	}

	return &Target{Name: name, Hashes: meta.Hashes, Length: meta.Length}, nil
}

// Publish pushes the local changes in signed material to the remote notary-server
func (r *NotaryRepository) Publish() error {
	r.bootstrapRepo()

	remote, err := getRemoteStore(r.Gun)

	root, err := r.fileStore.GetMeta("root", 0)
	if err != nil {
		return err
	}
	targets, err := r.fileStore.GetMeta("targets", 0)
	if err != nil {
		return err
	}
	snapshot, err := r.fileStore.GetMeta("snapshot", 0)
	if err != nil {
		return err
	}

	err = remote.SetMeta("root", root)
	if err != nil {
		return err
	}
	err = remote.SetMeta("targets", targets)
	if err != nil {
		return err
	}
	err = remote.SetMeta("snapshot", snapshot)
	if err != nil {
		return err
	}

	return nil
}

func (r *NotaryRepository) bootstrapRepo() error {
	fileStore, err := store.NewFilesystemStore(
		path.Join(viper.GetString("tufDir")),
		"metadata",
		"json",
		"targets",
	)
	if err != nil {
		return err
	}

	kdb := keys.NewDB()
	tufRepo := tuf.NewTufRepo(kdb, r.signer)

	fmt.Println("Loading trusted collection.")
	rootJSON, err := fileStore.GetMeta("root", 0)
	if err != nil {
		return err
	}
	root := &data.Signed{}
	err = json.Unmarshal(rootJSON, root)
	if err != nil {
		return err
	}
	tufRepo.SetRoot(root)
	targetsJSON, err := fileStore.GetMeta("targets", 0)
	if err != nil {
		return err
	}
	targets := &data.Signed{}
	err = json.Unmarshal(targetsJSON, targets)
	if err != nil {
		return err
	}
	tufRepo.SetTargets("targets", targets)
	snapshotJSON, err := fileStore.GetMeta("snapshot", 0)
	if err != nil {
		return err
	}
	snapshot := &data.Signed{}
	err = json.Unmarshal(snapshotJSON, snapshot)
	if err != nil {
		return err
	}
	tufRepo.SetSnapshot(snapshot)

	r.tufRepo = tufRepo
	r.fileStore = fileStore

	return nil
}

func (r *NotaryRepository) saveRepo() error {
	signedRoot, err := r.tufRepo.SignRoot(data.DefaultExpires("root"))
	if err != nil {
		return err
	}
	rootJSON, _ := json.Marshal(signedRoot)
	r.fileStore.SetMeta("root", rootJSON)

	fmt.Println("Saving changes to Trusted Collection.")

	for t, _ := range r.tufRepo.Targets {
		signedTargets, err := r.tufRepo.SignTargets(t, data.DefaultExpires("targets"))
		if err != nil {
			return err
		}
		targetsJSON, _ := json.Marshal(signedTargets)
		parentDir := filepath.Dir(t)
		os.MkdirAll(parentDir, 0755)
		r.fileStore.SetMeta(t, targetsJSON)
	}

	signedSnapshot, err := r.tufRepo.SignSnapshot(data.DefaultExpires("snapshot"))
	if err != nil {
		return err
	}
	snapshotJSON, _ := json.Marshal(signedSnapshot)
	r.fileStore.SetMeta("snapshot", snapshotJSON)

	return nil
}

/*
validateRoot iterates over every root key included in the TUF data and attempts
to validate the certificate by first checking for an exact match on the certificate
store, and subsequently trying to find a valid chain on the caStore.

Example TUF Content for root role:
"roles" : {
  "root" : {
    "threshold" : 1,
      "keyids" : [
        "e6da5c303d572712a086e669ecd4df7b785adfc844e0c9a7b1f21a7dfc477a38"
      ]
  },
 ...
}

Example TUF Content for root key:
"e6da5c303d572712a086e669ecd4df7b785adfc844e0c9a7b1f21a7dfc477a38" : {
	"keytype" : "RSA",
	"keyval" : {
	  "private" : "",
	  "public" : "Base64-encoded, PEM encoded x509 Certificate"
	}
}
*/
func (r *NotaryRepository) ValidateRoot(root *data.Signed) error {
	rootSigned := &data.Root{}
	err := json.Unmarshal(root.Signed, rootSigned)
	if err != nil {
		return err
	}
	certs := make(map[string]*data.PublicKey)
	for _, fingerprint := range rootSigned.Roles["root"].KeyIDs {
		// TODO(dlaw): currently assuming only one cert contained in
		// public key entry. Need to fix when we want to pass in chains.
		k, _ := pem.Decode([]byte(rootSigned.Keys["kid"].Public()))

		decodedCerts, err := x509.ParseCertificates(k.Bytes)
		if err != nil {
			continue
		}

		// TODO(diogo): Assuming that first certificate is the leaf-cert. Need to
		// iterate over all decodedCerts and find a non-CA one (should be the last).
		leafCert := decodedCerts[0]
		leafID := trustmanager.FingerprintCert(leafCert)

		// Check to see if there is an exact match of this certificate.
		// Checking the CommonName is not required since ID is calculated over
		// Cert.Raw. It's included to prevent breaking logic with changes of how the
		// ID gets computed.
		_, err = r.certificateStore.GetCertificateByFingerprint(leafID)
		if err == nil && leafCert.Subject.CommonName == r.Gun {
			certs[fingerprint] = rootSigned.Keys[fingerprint]
		}

		// Check to see if this leafCertificate has a chain to one of the Root CAs
		// of our CA Store.
		certList := []*x509.Certificate{leafCert}
		err = trustmanager.Verify(r.caStore, r.Gun, certList)
		if err == nil {
			certs[fingerprint] = rootSigned.Keys[fingerprint]
		}
	}
	_, err = signed.VerifyRoot(root, 0, certs, 1)

	return err
}

func (r *NotaryRepository) bootstrapClient() (*tufclient.Client, error) {
	remote, err := getRemoteStore(r.Gun)
	if err != nil {
		return nil, err
	}

	rootJSON, err := remote.GetMeta("root", 5<<20)
	root := &data.Signed{}
	err = json.Unmarshal(rootJSON, root)
	if err != nil {
		return nil, err
	}

	err = r.ValidateRoot(root)
	if err != nil {
		return nil, err
	}
	err = r.tufRepo.SetRoot(root)
	if err != nil {
		return nil, err
	}

	// TODO(dlaw): Where does this keyDB come in
	kdb := keys.NewDB()

	return tufclient.NewClient(
		r.tufRepo,
		remote,
		kdb,
	), nil
}

// ListPrivateKeys lists all availables private keys. Does not include private key
// material
func (c *NotaryClient) ListPrivateKeys() []data.PrivateKey {
	// TODO(diogo): Make this work
	for _, k := range c.rootKeyStore.ListAll() {
		fmt.Println(k)
	}
	return nil
}

// GenRootKey generates a new root key protected by a given passphrase
func (c *NotaryClient) GenRootKey(passphrase string) (*data.PublicKey, error) {
	// When generating a root key, passing in a GUN  to put into the cert
	// doesn't make sense since this key can be used for multiple distinct
	// repositories.
	pemKey, _, err := GenerateKeyAndCert("TUF root key")
	if err != nil {
		return nil, fmt.Errorf("could not generate private key: %v", err)
	}

	c.rootKeyStore.AddEncrypted("root", pemKey, passphrase)

	return data.NewPublicKey("RSA", pemKey), nil
}

// GetRepository returns a new repository
func (c *NotaryClient) GetRepository(gun string, baseURL string, transport http.RoundTripper) (*NotaryRepository, error) {
	privKeyStore, err := trustmanager.NewKeyFileStore(viper.GetString("privDir"))
	if err != nil {
		return nil, err
	}

	signer := signed.NewSigner(NewCryptoService(gun, privKeyStore))

	return &NotaryRepository{Gun: gun,
		baseURL:          baseURL,
		transport:        transport,
		signer:           signer,
		caStore:          c.caStore,
		certificateStore: c.certificateStore}, nil
}

func (c *NotaryClient) loadKeys(trustDir, rootKeysDir string) error {
	// Load all CAs that aren't expired and don't use SHA1
	caStore, err := trustmanager.NewX509FilteredFileStore(trustDir, func(cert *x509.Certificate) bool {
		return cert.IsCA && cert.BasicConstraintsValid && cert.SubjectKeyId != nil &&
			time.Now().Before(cert.NotAfter) &&
			cert.SignatureAlgorithm != x509.SHA1WithRSA &&
			cert.SignatureAlgorithm != x509.DSAWithSHA1 &&
			cert.SignatureAlgorithm != x509.ECDSAWithSHA1
	})
	if err != nil {
		return err
	}

	// Load all individual (non-CA) certificates that aren't expired and don't use SHA1
	certificateStore, err := trustmanager.NewX509FilteredFileStore(trustDir, func(cert *x509.Certificate) bool {
		return !cert.IsCA &&
			time.Now().Before(cert.NotAfter) &&
			cert.SignatureAlgorithm != x509.SHA1WithRSA &&
			cert.SignatureAlgorithm != x509.DSAWithSHA1 &&
			cert.SignatureAlgorithm != x509.ECDSAWithSHA1
	})
	if err != nil {
		return err
	}

	rootKeyStore, err := trustmanager.NewKeyFileStore(rootKeysDir)
	if err != nil {
		return err
	}

	c.caStore = caStore
	c.certificateStore = certificateStore
	c.rootKeyStore = rootKeyStore

	return nil
}

// Use this to initialize remote HTTPStores from the config settings
func getRemoteStore(gun string) (store.RemoteStore, error) {
	return store.NewHTTPStore(
		"https://notary:4443/v2/"+gun+"/_trust/tuf/",
		"",
		"json",
		"",
		"key",
	)
}
