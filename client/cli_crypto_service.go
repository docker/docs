package client

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/docker/notary/trustmanager"
	"github.com/endophage/gotuf/data"
)

type genericCryptoService struct {
	gun        string
	passphrase string
	keyStore   *trustmanager.KeyFileStore
}

// RSACryptoService implements Sign and Create, holding a specific GUN and keystore to
// operate on
type RSACryptoService struct {
	genericCryptoService
}

// ECDSACryptoService implements Sign and Create, holding a specific GUN and keystore to
// operate on
type ECDSACryptoService struct {
	genericCryptoService
}

// NewRSACryptoService returns an instance of CryptoService
func NewRSACryptoService(gun string, keyStore *trustmanager.KeyFileStore, passphrase string) *RSACryptoService {
	return &RSACryptoService{genericCryptoService{gun: gun, keyStore: keyStore, passphrase: passphrase}}
}

// Create is used to generate keys for targets, snapshots and timestamps
func (ccs *RSACryptoService) Create(role string) (*data.PublicKey, error) {
	privKey, err := trustmanager.GenerateRSAKey(rand.Reader, rsaKeySize)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA key: %v", err)
	}

	// Store the private key into our keystore with the name being: /GUN/ID.key
	ccs.keyStore.AddKey(filepath.Join(ccs.gun, privKey.ID()), privKey)

	return data.PublicKeyFromPrivate(*privKey), nil
}

// Sign returns the signatures for data with the given root Key ID, falling back
// if not rootKeyID is found
func (ccs *RSACryptoService) Sign(keyIDs []string, payload []byte) ([]data.Signature, error) {
	// Create hasher and hash data
	hash := crypto.SHA256
	hashed := sha256.Sum256(payload)

	signatures := make([]data.Signature, 0, len(keyIDs))
	for _, fingerprint := range keyIDs {
		// ccs.gun will be empty if this is the root key
		keyName := filepath.Join(ccs.gun, fingerprint)

		var privKey *data.PrivateKey
		var err error
		var method string

		// Read PrivateKey from file.
		// This assumes root keys always have to have a non-empty passphrase.
		if ccs.passphrase != "" {
			// This is a root key
			privKey, err = ccs.keyStore.GetDecryptedKey(keyName, ccs.passphrase)
			method = "RSASSA-PSS-X509"
		} else {
			privKey, err = ccs.keyStore.GetKey(keyName)
			method = "RSASSA-PSS"
		}
		if err != nil {
			// Note that GetDecryptedKey always fails on InitRepo.
			// InitRepo gets a signer that doesn't have access to
			// the root keys. Continuing here is safe because we
			// end up not returning any signatures.
			continue
		}

		sig, err := rsaSign(privKey, hash, hashed[:])
		if err != nil {
			return nil, err
		}

		// Append signatures to result array
		signatures = append(signatures, data.Signature{
			KeyID:     fingerprint,
			Method:    method,
			Signature: sig[:],
		})
	}

	return signatures, nil
}

func rsaSign(privKey *data.PrivateKey, hash crypto.Hash, hashed []byte) ([]byte, error) {
	if strings.ToLower(privKey.Cipher()) != "rsa" {
		return nil, fmt.Errorf("private key type not supported: %s", privKey.Cipher())
	}

	// Create an rsa.PrivateKey out of the private key bytes
	rsaPrivKey, err := x509.ParsePKCS1PrivateKey(privKey.Private())
	if err != nil {
		return nil, err
	}

	// Use the RSA key to sign the data
	sig, err := rsa.SignPSS(rand.Reader, rsaPrivKey, hash, hashed[:], &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash})
	if err != nil {
		return nil, err
	}

	return sig, nil
}

// NewECDSACryptoService returns an instance of CryptoService
func NewECDSACryptoService(gun string, keyStore *trustmanager.KeyFileStore, passphrase string) *ECDSACryptoService {
	return &ECDSACryptoService{genericCryptoService{gun: gun, keyStore: keyStore, passphrase: passphrase}}
}

// Create is used to generate keys for targets, snapshots and timestamps
func (ccs *ECDSACryptoService) Create(role string) (*data.PublicKey, error) {
	privKey, err := trustmanager.GenerateECDSAKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate EC key: %v", err)
	}

	// Store the private key into our keystore with the name being: /GUN/ID.key
	ccs.keyStore.AddKey(filepath.Join(ccs.gun, privKey.ID()), privKey)

	return data.PublicKeyFromPrivate(*privKey), nil
}

// Sign returns the signatures for data with the given root Key ID, falling back
// if not rootKeyID is found
func (ccs *ECDSACryptoService) Sign(keyIDs []string, payload []byte) ([]data.Signature, error) {
	// Create hasher and hash data
	hash := crypto.SHA256
	hashed := sha256.Sum256(payload)

	signatures := make([]data.Signature, 0, len(keyIDs))
	for _, fingerprint := range keyIDs {
		// ccs.gun will be empty if this is the root key
		keyName := filepath.Join(ccs.gun, fingerprint)

		var privKey *data.PrivateKey
		var err error
		// var method string

		// Read PrivateKey from file
		if ccs.passphrase != "" {
			// This is a root key
			privKey, err = ccs.keyStore.GetDecryptedKey(keyName, ccs.passphrase)
		} else {
			privKey, err = ccs.keyStore.GetKey(keyName)
		}
		if err != nil {
			// Note that GetDecryptedKey always fails on InitRepo.
			// InitRepo gets a signer that doesn't have access to
			// the root keys. Continuing here is safe because we
			// end up not returning any signatures.
			continue
		}

		sig, err := ecdsaSign(privKey, hash, hashed[:])
		if err != nil {
			return nil, err
		}

		// Append signatures to result array
		signatures = append(signatures, data.Signature{
			KeyID:     fingerprint,
			Method:    "ECDSA",
			Signature: sig[:],
		})
	}

	return signatures, nil
}

func ecdsaSign(privKey *data.PrivateKey, hash crypto.Hash, hashed []byte) ([]byte, error) {
	if strings.ToLower(privKey.Cipher()) != "ecdsa" {
		return nil, fmt.Errorf("private key type not supported: %s", privKey.Cipher())
	}

	// Create an ecdsa.PrivateKey out of the private key bytes
	ecdsaPrivKey, err := x509.ParseECPrivateKey(privKey.Private())
	if err != nil {
		return nil, err
	}

	// Use the ECDSA key to sign the data
	r, s, err := ecdsa.Sign(rand.Reader, ecdsaPrivKey, hashed[:])
	if err != nil {
		return nil, err
	}

	rBytes, sBytes := r.Bytes(), s.Bytes()
	octetLength := (ecdsaPrivKey.Params().BitSize + 7) >> 3
	// MUST include leading zeros in the output
	rBuf := make([]byte, octetLength-len(rBytes), octetLength)
	sBuf := make([]byte, octetLength-len(sBytes), octetLength)

	rBuf = append(rBuf, rBytes...)
	sBuf = append(sBuf, sBytes...)

	return append(rBuf, sBuf...), nil
}
