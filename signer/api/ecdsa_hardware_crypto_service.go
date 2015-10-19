package api

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"log"

	"github.com/docker/notary/trustmanager"
	"github.com/endophage/gotuf/data"
	"github.com/miekg/pkcs11"
)

// YubikeyECDSACryptoService is an implementation of SigningService
type YubikeyECDSACryptoService struct {
	context *pkcs11.Ctx
	session pkcs11.SessionHandle
}

// Create creates a key and returns its public components
func (s *YubikeyECDSACryptoService) Create(role string, algo data.KeyAlgorithm) (data.PublicKey, error) {
	if role != "root" {
		return nil, errors.New("can only generate root keys inside of Yubikey")
	}
	// TODO(diogo): Generate ECDSA key
	// TODO(diogo): Generate Certificate from ECDSA public key
	// TODO(diogo): Import Private key and certificate into yubikey
	// TODO(diogo): Return TUF Key
	return trustmanager.GenerateECDSAKey(rand.Reader)
}

// RemoveKey removes a key from the key database
func (s *YubikeyECDSACryptoService) RemoveKey(keyID string) error {
	return errors.New("cannot delete root key from Yubikey")
}

// GetKey returns the public components of a particular key
func (s *YubikeyECDSACryptoService) GetKey(keyID string) data.PublicKey {
	// TODO(diogo): Get the certificate out of the Yubikey and return the ID from it
	class = pkcs11.CKO_PRIVATE_KEY
	privateKeyTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, class),
		pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_ECDSA),
		pkcs11.NewAttribute(pkcs11.CKA_ID, []byte{2}),
	}

	if err := s.context.FindObjectsInit(session, privateKeyTemplate); err != nil {
		return nil
	}
	obj, b, err := s.context.FindObjects(session, 1)
	if err != nil {
		return nil
	}
	if err = s.context.FindObjectsFinal(session); e != nil {
		return nil
	}
	if len(obj) != 1 {
		return nil
	}

	return data.NewPublicKey(data.ECDSAKey, []byte{0})
}

// Sign returns a signature for a given signature request
func (s *YubikeyECDSACryptoService) Sign(keyIDs []string, payload []byte) ([]data.Signature, error) {
	signatures := make([]data.Signature, 0, len(keyIDs))
	for _, keyid := range keyIDs {
		key := s.GetKey(keyID)
		if key == nil {
			// We skip keys that aren't found
			continue
		}

		// Define the ECDSA Private key template
		class = pkcs11.CKO_PRIVATE_KEY
		privateKeyTemplate := []*pkcs11.Attribute{
			pkcs11.NewAttribute(pkcs11.CKA_CLASS, class),
			pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_ECDSA),
			pkcs11.NewAttribute(pkcs11.CKA_ID, []byte{2}),
		}

		if err := s.context.FindObjectsInit(session, privateKeyTemplate); err != nil {
			return nil, err
		}
		obj, b, err := s.context.FindObjects(session, 1)
		if err != nil {
			return nil, err
		}
		if err = s.context.FindObjectsFinal(session); e != nil {
			return nil, err
		}
		if len(obj) != 1 {
			return nil, errors.New("length of objects found not 1")
		}

		var sig []byte
		s.context.SignInit(s.session, []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_ECDSA, nil)}, obj[0])

		// Get the SHA256 of the payload
		digest := sha256.Sum256(payload)

		sig, err := s.context.Sign(s.session, digest)
		if err != nil {
			log.Printf("Error while signing: %s", err)
			return nil, err
		}

		if sig == nil {
			return nil, errors.New("Failed to create signature")
		}

		signatures = append(signatures, data.Signature{
			KeyID:     keyid,
			Method:    data.ECDSASignature,
			Signature: sig[:],
		})
	}
	return signatures, nil
}

// NewYubikeyECDSACryptoService returns an instance of RSAHardwareCryptoService
func NewYubikeyECDSACryptoService(ctx *pkcs11.Ctx, session pkcs11.SessionHandle) *YubikeyECDSACryptoService {
	return &YubikeyECDSACryptoService{
		context: ctx,
		session: session,
	}
}
