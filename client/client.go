package client

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/docker/notary/client/changelist"
	"github.com/docker/notary/trustmanager"
	"github.com/endophage/gotuf"
	tufclient "github.com/endophage/gotuf/client"
	"github.com/endophage/gotuf/data"
	"github.com/endophage/gotuf/keys"
	"github.com/endophage/gotuf/signed"
	"github.com/endophage/gotuf/store"
)

type ErrRepoNotInitialized struct{}

func (err *ErrRepoNotInitialized) Error() string {
	return "Repository has not been initialized"
}

// Default paths should end with a '/' so directory creation works correctly
const (
	trustDir    string = "/trusted_certificates/"
	privDir     string = "/private/"
	tufDir      string = "/tuf/"
	rootKeysDir string = privDir + "/root_keys/"
)
const rsaKeySize int = 2048

// ErrRepositoryNotExist gets returned when trying to make an action over a repository
/// that doesn't exist
var ErrRepositoryNotExist = errors.New("repository does not exist")

type UnlockedSigner struct {
	privKey *data.PrivateKey
	signer  *signed.Signer
}

type NotaryClient struct {
	baseDir          string
	caStore          trustmanager.X509Store
	certificateStore trustmanager.X509Store
	rootKeyStore     *trustmanager.KeyFileStore
}

type NotaryRepository struct {
	Gun              string
	baseURL          string
	tufRepoPath      string
	transport        http.RoundTripper
	signer           *signed.Signer
	tufRepo          *tuf.TufRepo
	fileStore        store.MetadataStore
	privKeyStore     *trustmanager.KeyFileStore
	caStore          trustmanager.X509Store
	certificateStore trustmanager.X509Store
	rootSigner       *UnlockedSigner
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
func NewClient(baseDir string) (*NotaryClient, error) {
	trustDir := filepath.Join(baseDir, trustDir)
	rootKeysDir := filepath.Join(baseDir, rootKeysDir)

	nClient := &NotaryClient{baseDir: baseDir}

	if err := nClient.loadKeys(trustDir, rootKeysDir); err != nil {
		return nil, err
	}

	return nClient, nil
}

// Initialize creates a new repository by using rootKey as the root Key for the
// TUF repository.
func (r *NotaryRepository) Initialize() error {
	remote, err := getRemoteStore(r.Gun)
	rawTSKey, err := remote.GetKey("timestamp")
	if err != nil {
		return err
	}

	parsedKey := &data.TUFKey{}
	err = json.Unmarshal(rawTSKey, parsedKey)
	if err != nil {
		return err
	}

	timestampKey := data.NewPublicKey(parsedKey.Cipher(), parsedKey.Public())
	rootKey := r.rootSigner.PublicKey()

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

	if err := kdb.AddRole(rootRole); err != nil {
		return err
	}
	if err := kdb.AddRole(targetsRole); err != nil {
		return err
	}
	if err := kdb.AddRole(snapshotRole); err != nil {
		return err
	}
	if err := kdb.AddRole(timestampRole); err != nil {
		return err
	}

	r.tufRepo = tuf.NewTufRepo(kdb, r.signer)

	r.fileStore, err = store.NewFilesystemStore(
		r.tufRepoPath,
		"metadata",
		"json",
		"targets",
	)
	if err != nil {
		return err
	}

	if err := r.tufRepo.InitRepo(false); err != nil {
		return err
	}

	if err := r.saveMetadata(uSigner.signer); err != nil {
		return err
	}

	// Creates an empty snapshot
	return r.snapshot()
}

// AddTarget adds a new target to the repository, forcing a timestamps check from TUF
func (r *NotaryRepository) AddTarget(target *Target) error {
	cl, err := changelist.NewFileChangelist(filepath.Join(r.tufRepoPath, "changelist"))
	if err != nil {
		return err
	}
	fmt.Printf("Adding target \"%s\" with sha256 \"%s\" and size %d bytes.\n", target.Name, target.Hashes["sha256"], target.Length)

	meta := data.FileMeta{Length: target.Length, Hashes: target.Hashes}
	metaJSON, err := json.Marshal(meta)
	if err != nil {
		return err
	}

	c := changelist.NewTufChange(changelist.ActionCreate, "targets", "target", target.Name, metaJSON)
	err = cl.Add(c)
	if err != nil {
		return err
	}
	return cl.Close()
}

// ListTargets lists all targets for the current repository
func (r *NotaryRepository) ListTargets() ([]*Target, error) {
	//r.bootstrapRepo()

	c, err := r.bootstrapClient()
	if err != nil {
		return nil, err
	}

	err = c.Update()
	if err != nil {
		return nil, err
	}

	targetList := make([]*Target, 0)
	for name, meta := range r.tufRepo.Targets["targets"].Signed.Targets {
		target := &Target{Name: name, Hashes: meta.Hashes, Length: meta.Length}
		targetList = append(targetList, target)
	}

	return targetList, nil
}

// GetTargetByName returns a target given a name
func (r *NotaryRepository) GetTargetByName(name string) (*Target, error) {
	//r.bootstrapRepo()

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
	c, err := r.bootstrapClient() // just need the repo to be initialized from remote
	if err != nil {
		if _, ok := err.(*store.ErrMetaNotFound); ok {
			// attempt to load locally to see if it's already init'ed
			err := r.bootstrapRepo()
			if err != nil {
				logrus.Debug("Repository not initialized during Publish")
				return &ErrRepoNotInitialized{} // caller must init
			}
		} else {
			logrus.Error("Could not publish Repository: ", err.Error())
			return err
		}
	}
	err = c.Update()
	if err != nil {
		return err
	}

	cl, err := changelist.NewFileChangelist(filepath.Join(r.tufRepoPath, "changelist"))
	if err != nil {
		logrus.Debug("Error initializing changelist")
		return err
	}
	err = applyChangelist(r.tufRepo, cl)
	if err != nil {
		logrus.Debug("Error applying changelist")
		return err
	}

	root, err := r.tufRepo.SignRoot(data.DefaultExpires("root"), r.rootSigner)
	if err != nil {
		return err
	}
	targets, err := r.tufRepo.SignTargets("targets", data.DefaultExpires("targets"), nil)
	if err != nil {
		return err
	}
	snapshot, err := r.tufRepo.SignSnapshot(data.DefaultExpires("snapshot"), nil)
	if err != nil {
		return err
	}

	rootJSON, err := json.Marshal(root)
	if err != nil {
		return err
	}
	targetsJSON, err := json.Marshal(targets)
	if err != nil {
		return err
	}
	snapshotJSON, err := json.Marshal(snapshot)
	if err != nil {
		return err
	}

	remote, err := getRemoteStore(r.Gun)
	if err != nil {
		return err
	}
	err = remote.SetMeta("root", rootJSON)
	if err != nil {
		return err
	}
	err = remote.SetMeta("targets", targetsJSON)
	if err != nil {
		return err
	}
	err = remote.SetMeta("snapshot", snapshotJSON)
	if err != nil {
		return err
	}

	return nil
}

func (r *NotaryRepository) bootstrapRepo() error {
	fileStore, err := store.NewFilesystemStore(
		r.tufRepoPath,
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

func (r *NotaryRepository) saveMetadata(rootSigner *signed.Signer) error {
	signedRoot, err := r.tufRepo.SignRoot(data.DefaultExpires("root"), rootSigner)
	if err != nil {
		return err
	}

	rootJSON, _ := json.Marshal(signedRoot)
	return r.fileStore.SetMeta("root", rootJSON)
}

func (r *NotaryRepository) snapshot() error {
	fmt.Println("Saving changes to Trusted Collection.")

	for t, _ := range r.tufRepo.Targets {
		signedTargets, err := r.tufRepo.SignTargets(t, data.DefaultExpires("targets"), nil)
		if err != nil {
			return err
		}
		targetsJSON, _ := json.Marshal(signedTargets)
		parentDir := filepath.Dir(t)
		os.MkdirAll(parentDir, 0755)
		r.fileStore.SetMeta(t, targetsJSON)
	}

	signedSnapshot, err := r.tufRepo.SignSnapshot(data.DefaultExpires("snapshot"), nil)
	if err != nil {
		return err
	}
	snapshotJSON, _ := json.Marshal(signedSnapshot)

	return r.fileStore.SetMeta("snapshot", snapshotJSON)
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
		k, _ := pem.Decode([]byte(rootSigned.Keys[fingerprint].Public()))
		logrus.Debug("Root PEM: ", k)
		logrus.Debug("Root ID: ", fingerprint)
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

	if len(certs) < 1 {
		return errors.New("could not validate the path to a trusted root")
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
	if err != nil {
		return nil, err
	}
	root := &data.Signed{}
	err = json.Unmarshal(rootJSON, root)
	if err != nil {
		return nil, err
	}

	err = r.ValidateRoot(root)
	if err != nil {
		return nil, err
	}

	kdb := keys.NewDB()
	r.tufRepo = tuf.NewTufRepo(kdb, r.signer)

	err = r.tufRepo.SetRoot(root)
	if err != nil {
		return nil, err
	}

	// TODO(dlaw): Where does this keyDB come in

	return tufclient.NewClient(
		r.tufRepo,
		remote,
		kdb,
	), nil
}

// ListPrivateKeys lists all availables private keys. Does not include private key
// material
func (c *NotaryClient) ListPrivateKeys() []string {
	// TODO(diogo): Make this work
	for _, k := range c.rootKeyStore.ListAll() {
		fmt.Println(k)
	}
	return nil
}

// GenRootKey generates a new root key protected by a given passphrase
func (c *NotaryClient) GenRootKey(passphrase string) (string, error) {
	// TODO(diogo): Refactor TUF Key creation. We should never see crypto.privatekeys
	// Generates a new RSA key
	rsaPrivKey, err := rsa.GenerateKey(rand.Reader, rsaKeySize)
	if err != nil {
		return "", fmt.Errorf("could not generate private key: %v", err)
	}

	// Encode the private key in PEM format since that is the final storage format
	pemPrivKey, err := trustmanager.KeyToPEM(rsaPrivKey)
	if err != nil {
		return "", fmt.Errorf("failed to encode the private key: %v", err)
	}

	tufPrivKey, err := trustmanager.RSAToPrivateKey(rsaPrivKey)
	if err != nil {
		return "", fmt.Errorf("failed to convert private key: ", err)
	}

	c.rootKeyStore.AddEncrypted(tufPrivKey.ID(), pemPrivKey, passphrase)

	return tufPrivKey.ID(), nil
}

// GetRootSigner retreives a root key that includes the ID and a signer
func (c *NotaryClient) GetRootSigner(rootKeyID, passphrase string) (*UnlockedSigner, error) {
	pemPrivKey, err := c.rootKeyStore.GetDecrypted(rootKeyID, passphrase)
	if err != nil {
		return nil, fmt.Errorf("could not get decrypted root key: %v", err)
	}

	tufPrivKey, err := trustmanager.TufParsePEMPrivateKey(pemPrivKey)
	if err != nil {
		return nil, fmt.Errorf("could not get parse root key: %v", err)
	}

	signer := signed.NewSigner(NewRootCryptoService(c.rootKeyStore, passphrase))

	return &UnlockedSigner{
		privKey: tufPrivKey,
		signer:  signer}, nil
}

// GetRepository returns a new repository
func (c *NotaryClient) GetRepository(gun string, baseURL string, transport http.RoundTripper, uSigner *UnlockedSigner) (*NotaryRepository, error) {
	privKeyStore, err := trustmanager.NewKeyFileStore(filepath.Join(c.baseDir, privDir))
	if err != nil {
		return nil, err
	}

	signer := signed.NewSigner(NewCryptoService(gun, privKeyStore))

	return &NotaryRepository{Gun: gun,
		baseURL:          baseURL,
		tufRepoPath:      filepath.Join(c.baseDir, tufDir, gun),
		transport:        transport,
		signer:           signer,
		privKeyStore:     privKeyStore,
		caStore:          c.caStore,
		certificateStore: c.certificateStore,
		rootSigner:       uSigner,
	}, nil
}

func (c *NotaryClient) InitRepository(gun string, baseURL string, transport http.RoundTripper, uSigner *UnlockedSigner) (*NotaryRepository, error) {
	// Creates and saves a trusted certificate for this store, with this root key
	rootCert, err := uSigner.GenerateCertificate(gun)
	if err != nil {
		return nil, err
	}
	c.certificateStore.AddCert(rootCert)
	rootKey := data.NewPublicKey("RSA", trustmanager.CertToPEM(rootCert))
	err = c.rootKeyStore.Link(uSigner.ID(), rootKey.ID())
	if err != nil {
		return nil, err
	}

	privKeyStore, err := trustmanager.NewKeyFileStore(filepath.Join(c.baseDir, privDir))
	if err != nil {
		return nil, err
	}

	signer := signed.NewSigner(NewCryptoService(gun, privKeyStore))

	nRepo := &NotaryRepository{Gun: gun,
		baseURL:          baseURL,
		tufRepoPath:      filepath.Join(c.baseDir, tufDir, gun),
		transport:        transport,
		signer:           signer,
		privKeyStore:     privKeyStore,
		caStore:          c.caStore,
		certificateStore: c.certificateStore,
		rootSigner:       uSigner,
	}

	err = nRepo.Initialize()
	if err != nil {
		return nil, err
	}

	return nRepo, nil
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

// ID gets a consistent ID based on the PrivateKey bytes and cipher type
func (uk *UnlockedSigner) ID() string {
	return uk.PublicKey().ID()
}

// PublicKey Returns the public key associated with the Root Key
func (uk *UnlockedSigner) PublicKey() *data.PublicKey {
	return data.PublicKeyFromPrivate(*uk.privKey)
}

// GenerateCertificate
func (uk *UnlockedSigner) GenerateCertificate(gun string) (*x509.Certificate, error) {
	privKey, err := x509.ParsePKCS1PrivateKey(uk.privKey.Private())
	if err != nil {
		return nil, fmt.Errorf("failed to parse root key: %v (%s)", gun, err.Error())
	}

	//TODO (diogo): We're hardcoding the Organization to be the GUN. Probably want to change it
	template := trustmanager.NewCertificate(gun, gun)
	derBytes, err := x509.CreateCertificate(rand.Reader, template, template, privKey.Public(), privKey)
	if err != nil {
		return nil, fmt.Errorf("failed to generate the certificate for: %v (%s)", gun, err.Error())
	}

	// Encode the new certificate into PEM
	cert, err := x509.ParseCertificate(derBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the certificate for key: %v (%s)", gun, err.Error())
	}

	return cert, nil
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

func applyChangelist(repo *tuf.TufRepo, cl changelist.Changelist) error {
	changes := cl.List()
	var err error
	for _, c := range changes {
		if c.Scope() == "targets" {
			applyTargetsChange(repo, c)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func applyTargetsChange(repo *tuf.TufRepo, c changelist.Change) error {
	var err error
	meta := &data.FileMeta{}
	err = json.Unmarshal(c.Content(), meta)
	if err != nil {
		return nil
	}
	if c.Action() == changelist.ActionCreate {
		files := data.Files{c.Path(): *meta}
		_, err = repo.AddTargets("targets", files)
	} else if c.Action() == changelist.ActionDelete {
		err = repo.RemoveTargets("targets", c.Path())
	}
	if err != nil {
		// TODO(endophage): print out rem entries as files that couldn't
		//                  be added.
		return err
	}
	return nil
}
