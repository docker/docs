package api

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"errors"
	"log"
	"math/big"

	"github.com/docker/notary/signer/keys"
	"github.com/docker/notary/tuf/data"
	"github.com/miekg/pkcs11"
)

// RSAHardwareCryptoService is an implementation of SigningService
type RSAHardwareCryptoService struct {
	keys    map[string]*keys.HSMRSAKey
	context *pkcs11.Ctx
	session pkcs11.SessionHandle
}

// Create creates a key and returns its public components
func (s *RSAHardwareCryptoService) Create(role string, algo data.KeyAlgorithm) (data.PublicKey, error) {
	// For now generate random labels for keys
	// (diogo): add link between keyID and label in database so we can support multiple keys
	randomLabel := make([]byte, 32)
	_, err := rand.Read(randomLabel)
	if err != nil {
		return nil, errors.New("Could not generate a random key label.")
	}

	// Set the public key template
	// CKA_TOKEN: Guarantees key persistence in hardware
	// CKA_LABEL: Identifies this specific key inside of the HSM
	publicKeyTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_TOKEN, true),
		pkcs11.NewAttribute(pkcs11.CKA_PUBLIC_EXPONENT, []byte{3}),
		pkcs11.NewAttribute(pkcs11.CKA_MODULUS_BITS, 2048),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, string(randomLabel)),
	}
	privateKeyTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_TOKEN, true),
		pkcs11.NewAttribute(pkcs11.CKA_PRIVATE, true),
		pkcs11.NewAttribute(pkcs11.CKA_SIGN, true),
		pkcs11.NewAttribute(pkcs11.CKA_LABEL, string(randomLabel)),
	}

	// Generate a new RSA private/public keypair inside of the HSM
	pub, priv, err := s.context.GenerateKeyPair(s.session,
		[]*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_RSA_PKCS_KEY_PAIR_GEN, nil)},
		publicKeyTemplate, privateKeyTemplate)
	if err != nil {
		return nil, errors.New("Could not generate a new key inside of the HSM.")
	}

	// (diogo): This template is used for the GetAttribute
	template := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_PUBLIC_EXPONENT, nil),
		pkcs11.NewAttribute(pkcs11.CKA_MODULUS_BITS, nil),
		pkcs11.NewAttribute(pkcs11.CKA_MODULUS, nil),
	}

	// Retrieve the public-key material to be able to create a new HSMRSAKey
	attr, err := s.context.GetAttributeValue(s.session, pub, template)
	if err != nil {
		return nil, errors.New("Failed to get Attribute value.")
	}

	// We're going to store the elements of the RSA Public key, exponent and Modulus inside of exp and mod
	var exp int
	mod := big.NewInt(0)

	// Iterate through all the attributes of this key and saves CKA_PUBLIC_EXPONENT and CKA_MODULUS. Removes ordering specific issues.
	for _, a := range attr {
		if a.Type == pkcs11.CKA_PUBLIC_EXPONENT {
			exp, _ = readInt(a.Value)
		}

		if a.Type == pkcs11.CKA_MODULUS {
			mod.SetBytes(a.Value)
		}
	}

	rsaPublicKey := rsa.PublicKey{N: mod, E: exp}
	// Using x509 to Marshal the Public key into der encoding
	pubBytes, err := x509.MarshalPKIXPublicKey(&rsaPublicKey)
	if err != nil {
		return nil, errors.New("Failed to Marshal public key.")
	}

	// (diogo): Ideally I would like to return base64 PEM encoded public keys to the client
	k := keys.NewHSMRSAKey(pubBytes, priv)

	keyID := k.ID()

	s.keys[keyID] = k

	return k, nil
}

// RemoveKey removes a key from the key database
func (s *RSAHardwareCryptoService) RemoveKey(keyID string) error {
	if _, ok := s.keys[keyID]; !ok {
		return keys.ErrInvalidKeyID
	}

	delete(s.keys, keyID)
	return nil
}

// GetKey returns the public components of a particular key
func (s *RSAHardwareCryptoService) GetKey(keyID string) data.PublicKey {
	key, ok := s.keys[keyID]
	if !ok {
		return nil
	}
	return key
}

// Sign returns a signature for a given signature request
func (s *RSAHardwareCryptoService) Sign(keyIDs []string, payload []byte) ([]data.Signature, error) {
	signatures := make([]data.Signature, 0, len(keyIDs))
	for _, keyid := range keyIDs {
		privateKey, present := s.keys[keyid]
		if !present {
			// We skip keys that aren't found
			continue
		}

		priv := privateKey.PKCS11ObjectHandle()
		var sig []byte
		var err error
		for i := 0; i < 3; i++ {
			s.context.SignInit(s.session, []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_SHA256_RSA_PKCS, nil)}, priv)

			sig, err = s.context.Sign(s.session, payload)
			if err != nil {
				log.Printf("Error while signing: %s", err)
				continue
			}

			digest := sha256.Sum256(payload)
			pub, err := x509.ParsePKIXPublicKey(privateKey.Public())
			if err != nil {
				log.Printf("Failed to parse public key: %s\n", err)
				return nil, err
			}

			rsaPub, ok := pub.(*rsa.PublicKey)
			if !ok {
				log.Printf("Value returned from ParsePKIXPublicKey was not an RSA public key")
				return nil, err
			}

			err = rsa.VerifyPKCS1v15(rsaPub, crypto.SHA256, digest[:], sig)
			if err != nil {
				log.Printf("Failed verification. Retrying: %s", err)
				continue
			}
			break
		}

		if sig == nil {
			return nil, errors.New("Failed to create signature")
		}

		signatures = append(signatures, data.Signature{
			KeyID:     keyid,
			Method:    data.RSAPKCS1v15Signature,
			Signature: sig[:],
		})
	}

	return signatures, nil
}

// NewRSAHardwareCryptoService returns an instance of RSAHardwareCryptoService
func NewRSAHardwareCryptoService(ctx *pkcs11.Ctx, session pkcs11.SessionHandle) *RSAHardwareCryptoService {
	return &RSAHardwareCryptoService{
		keys:    make(map[string]*keys.HSMRSAKey),
		context: ctx,
		session: session,
	}
}

// readInt converts a []byte into an int. It is used to convert the RSA Public key exponent into an int to create a crypto.PublicKey
func readInt(data []byte) (int, error) {
	var ret int
	if len(data) > 4 {
		return 0, errors.New("Cannot convert byte array due to size")
	}

	for i, a := range data {
		ret |= (int(a) << uint(i*8))
	}
	return ret, nil
}
