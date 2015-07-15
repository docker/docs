package cryptoservice

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"fmt"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/docker/notary/trustmanager"
	"github.com/endophage/gotuf/data"
)

const (
	rsaKeySize = 2048 // Used for snapshots and targets keys
)

// CryptoService implements Sign and Create, holding a specific GUN and keystore to
// operate on
type CryptoService struct {
	gun        string
	passphrase string
	keyStore   *trustmanager.KeyFileStore
}

// NewCryptoService returns an instance of CryptoService
func NewCryptoService(gun string, keyStore *trustmanager.KeyFileStore, passphrase string) *CryptoService {
	return &CryptoService{gun: gun, keyStore: keyStore, passphrase: passphrase}
}

// Create is used to generate keys for targets, snapshots and timestamps
func (ccs *CryptoService) Create(role string, algorithm data.KeyAlgorithm) (*data.PublicKey, error) {
	var privKey *data.PrivateKey
	var err error

	switch algorithm {
	case data.RSAKey:
		privKey, err = trustmanager.GenerateRSAKey(rand.Reader, rsaKeySize)
		if err != nil {
			return nil, fmt.Errorf("failed to generate RSA key: %v", err)
		}
	case data.ECDSAKey:
		privKey, err = trustmanager.GenerateECDSAKey(rand.Reader)
		if err != nil {
			return nil, fmt.Errorf("failed to generate EC key: %v", err)
		}
	default:
		return nil, fmt.Errorf("private key type not supported for key generation: %s", algorithm)
	}
	logrus.Debugf("generated new %s key for role: %s and keyID: %s", algorithm, role, privKey.ID())

	// Store the private key into our keystore with the name being: /GUN/ID.key
	err = ccs.keyStore.AddKey(filepath.Join(ccs.gun, privKey.ID()), privKey)
	if err != nil {
		return nil, fmt.Errorf("failed to add key to filestore: %v", err)
	}
	return data.PublicKeyFromPrivate(*privKey), nil
}

// Sign returns the signatures for the payload with a set of keyIDs. It ignores
// errors to sign and expects the called to validate if the number of returned
// signatures is adequate.
func (ccs *CryptoService) Sign(keyIDs []string, payload []byte) ([]data.Signature, error) {
	// Create hasher and hash data
	hash := crypto.SHA256
	hashed := sha256.Sum256(payload)

	signatures := make([]data.Signature, 0, len(keyIDs))
	for _, keyid := range keyIDs {
		// ccs.gun will be empty if this is the root key
		keyName := filepath.Join(ccs.gun, keyid)

		var privKey *data.PrivateKey
		var err error

		// Read PrivateKey from file and decrypt it if there is a passphrase.
		if ccs.passphrase != "" {
			privKey, err = ccs.keyStore.GetDecryptedKey(keyName, ccs.passphrase)
		} else {
			privKey, err = ccs.keyStore.GetKey(keyName)
		}
		if err != nil {
			// Note that GetDecryptedKey always fails on InitRepo.
			// InitRepo gets a signer that doesn't have access to
			// the root keys. Continuing here is safe because we
			// end up not returning any signatures.
			logrus.Debugf("ignoring error attempting to retrieve key ID: %s, %v", keyid, err)
			continue
		}

		algorithm := privKey.Algorithm()
		var sigAlgorithm data.SigAlgorithm
		var sig []byte

		switch algorithm {
		case data.RSAKey:
			sig, err = rsaSign(privKey, hash, hashed[:])
			sigAlgorithm = data.RSAPSSSignature
		case data.ECDSAKey:
			sig, err = ecdsaSign(privKey, hashed[:])
			sigAlgorithm = data.ECDSASignature
		}
		if err != nil {
			logrus.Debugf("ignoring error attempting to %s sign with keyID: %s, %v", algorithm, keyid, err)
			continue
		}

		logrus.Debugf("appending %s signature with Key ID: %s", algorithm, keyid)

		// Append signatures to result array
		signatures = append(signatures, data.Signature{
			KeyID:     keyid,
			Method:    sigAlgorithm,
			Signature: sig[:],
		})
	}

	return signatures, nil
}

func rsaSign(privKey *data.PrivateKey, hash crypto.Hash, hashed []byte) ([]byte, error) {
	if privKey.Algorithm() != data.RSAKey {
		return nil, fmt.Errorf("private key type not supported: %s", privKey.Algorithm())
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

func ecdsaSign(privKey *data.PrivateKey, hashed []byte) ([]byte, error) {
	if privKey.Algorithm() != data.ECDSAKey {
		return nil, fmt.Errorf("private key type not supported: %s", privKey.Algorithm())
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
