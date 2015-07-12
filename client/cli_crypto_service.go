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

	"github.com/Sirupsen/logrus"
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
	err = ccs.keyStore.AddKey(filepath.Join(ccs.gun, privKey.ID()), privKey)
	if err != nil {
		return nil, fmt.Errorf("failed to add key to filestore: %v", err)
	}

	logrus.Debugf("generated new RSA key for role: %s and keyID: %s", role, privKey.ID())

	return data.PublicKeyFromPrivate(*privKey), nil
}

// Sign returns the signatures for the payload with a set of keyIDs. It ignores
// errors to sign and expects the called to validate if the number of returned
// signatures is adequate.
func (ccs *RSACryptoService) Sign(keyIDs []string, payload []byte) ([]data.Signature, error) {
	// Create hasher and hash data
	hash := crypto.SHA256
	hashed := sha256.Sum256(payload)

	signatures := make([]data.Signature, 0, len(keyIDs))
	for _, keyid := range keyIDs {
		// ccs.gun will be empty if this is the root key
		keyName := filepath.Join(ccs.gun, keyid)

		var privKey *data.PrivateKey
		var err error
		var method string

		// Read PrivateKey from file.
		// TODO(diogo): This assumes both that only root keys are encrypted and
		// that encrypted keys are always X509 encoded
		if ccs.passphrase != "" {
			// This is a root key
			privKey, err = ccs.keyStore.GetDecryptedKey(keyName, ccs.passphrase)
			// This method is added to the signature so verifiers can decode
			// the X509 certificate before trying to parse the public materials
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
			logrus.Debugf("ignoring error attempting to retrieve RSA key ID: %s, %v", privKey.ID(), err)
			continue
		}

		sig, err := rsaSign(privKey, hash, hashed[:])
		if err != nil {
			// If the rsaSign method got called with a non RSA private key,
			// we ignore this call.
			// This might happen when root is ECDSA, targets and snapshots RSA, and
			// gotuf still attempts to sign root with this cryptoserver
			// return nil, err
			logrus.Debugf("ignoring error attempting to RSA sign with keyID: %s, %v", keyid, err)
			continue
		}

		logrus.Debugf("appending RSA signature with Key ID: %s and method %s", privKey.ID(), method)

		// Append signatures to result array
		signatures = append(signatures, data.Signature{
			KeyID:     keyid,
			Method:    method,
			Signature: sig[:],
		})
	}

	return signatures, nil
}

func rsaSign(privKey *data.PrivateKey, hash crypto.Hash, hashed []byte) ([]byte, error) {
	if strings.ToUpper(privKey.Cipher()) != "RSA" {
		return nil, fmt.Errorf("private key type not supported: %s", privKey.Cipher())
	}

	// Create an rsa.PrivateKey out of the private key bytes
	rsaPrivKey, err := x509.ParsePKCS1PrivateKey(privKey.Private())
	if err != nil {
		return nil, err
	}

	// Use the RSA key to RSASSA-PSS sign the data
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
	err = ccs.keyStore.AddKey(filepath.Join(ccs.gun, privKey.ID()), privKey)
	if err != nil {
		return nil, fmt.Errorf("failed to add key to filestore: %v", err)
	}

	logrus.Debugf("generated new ECDSA key for role: %s with keyID: %s", role, privKey.ID())

	return data.PublicKeyFromPrivate(*privKey), nil
}

// Sign returns the signatures for the payload with a set of keyIDs. It ignores
// errors to sign and expects the called to validate if the number of returned
// signatures is adequate.
func (ccs *ECDSACryptoService) Sign(keyIDs []string, payload []byte) ([]data.Signature, error) {
	// Create hasher and hash data
	hash := crypto.SHA256
	hashed := sha256.Sum256(payload)

	signatures := make([]data.Signature, 0, len(keyIDs))
	for _, keyid := range keyIDs {
		// ccs.gun will be empty if this is the root key
		keyName := filepath.Join(ccs.gun, keyid)

		var privKey *data.PrivateKey
		var err error
		var method string

		// Read PrivateKey from file
		// TODO(diogo): This assumes both that only root keys are encrypted and
		// that encrypted keys are always X509 encoded
		if ccs.passphrase != "" {
			// This is a root key
			privKey, err = ccs.keyStore.GetDecryptedKey(keyName, ccs.passphrase)
			// This method is added to the signature so verifiers can decode
			// the X509 certificate before trying to parse the public materials
			method = "ECDSA-X509"
		} else {
			privKey, err = ccs.keyStore.GetKey(keyName)
			method = "ECDSA"
		}
		if err != nil {
			// Note that GetDecryptedKey always fails on InitRepo.
			// InitRepo gets a signer that doesn't have access to
			// the root keys. Continuing here is safe because we
			// end up not returning any signatures.
			// TODO(diogo): figure out if there are any specific error types to
			// check. We're swallowing all errors.
			logrus.Debugf("Ignoring error attempting to retrieve ECDSA key ID: %s, %v", keyid, err)
			continue
		}
		if err != nil {
			fmt.Println("ERROR: ", err.Error())
		}

		sig, err := ecdsaSign(privKey, hash, hashed[:])
		if err != nil {
			// If the ecdsaSign method got called with a non ECDSA private key,
			// we ignore this call.
			// This might happen when root is RSA, targets and snapshots ECDSA, and
			// gotuf still attempts to sign root with this cryptoserver
			// return nil, err
			logrus.Debugf("ignoring error attempting to ECDSA sign with keyID: %s, %v", privKey.ID(), err)
			continue
		}

		logrus.Debugf("appending ECDSA signature with Key ID: %s and method: %s", privKey.ID(), method)

		// Append signatures to result array
		signatures = append(signatures, data.Signature{
			KeyID:     keyid,
			Method:    method,
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
