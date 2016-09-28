package jose

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"io"
	"math/big"
)

// SignatureAlgorithm holds info about a JWS signing algorithm.
type SignatureAlgorithm interface {
	// Name returns the "alg" (Algorithm) header parameter value for a JSON
	// Web Signature as defined in
	// https://tools.ietf.org/html/rfc7518#section-3.1
	Name() string
	// KeyType returns the "kty" (Key Type) parameter value for a JSON Web
	// Key as defined in https://tools.ietf.org/html/rfc7518#section-6.1
	KeyType() string
	// Hash returns a cryptographic hash function identifier that is used
	// with the signature algorithm as detailed in
	// https://tools.ietf.org/html/rfc7518#section-3.1
	Hash() crypto.Hash
}

type signatureAlg struct {
	name    string
	keyType string
	hash    crypto.Hash
}

func (sigAlg signatureAlg) Name() string {
	return sigAlg.name
}

func (sigAlg signatureAlg) KeyType() string {
	return sigAlg.keyType
}

func (sigAlg signatureAlg) Hash() crypto.Hash {
	return sigAlg.hash
}

// SignatureAlgs used by this package.
var (
	rs256 = signatureAlg{"RS256", "RSA", crypto.SHA256}
	rs384 = signatureAlg{"RS384", "RSA", crypto.SHA384}
	rs512 = signatureAlg{"RS512", "RSA", crypto.SHA512}
	es256 = signatureAlg{"ES256", "EC", crypto.SHA256}
	es384 = signatureAlg{"ES384", "EC", crypto.SHA384}
	es512 = signatureAlg{"ES512", "EC", crypto.SHA512}
)

// RS256 returns the JWA Digital Signature Algorithm "RS256" defined in
// https://tools.ietf.org/html/rfc7518#section-3.3
func RS256() SignatureAlgorithm {
	return rs256
}

// RS384 returns the JWA Digital Signature Algorithm "RS384" defined in
// https://tools.ietf.org/html/rfc7518#section-3.3
func RS384() SignatureAlgorithm {
	return rs384
}

// RS512 returns the JWA Digital Signature Algorithm "RS512" defined in
// https://tools.ietf.org/html/rfc7518#section-3.3
func RS512() SignatureAlgorithm {
	return rs512
}

// ES256 returns the JWA Digital Signature Algorithm "ES256" defined in
// https://tools.ietf.org/html/rfc7518#section-3.4
func ES256() SignatureAlgorithm {
	return es256
}

// ES384 returns the JWA Digital Signature Algorithm "ES384" defined in
// https://tools.ietf.org/html/rfc7518#section-3.4
func ES384() SignatureAlgorithm {
	return es384
}

// ES512 returns the JWA Digital Signature Algorithm "ES512" defined in
// https://tools.ietf.org/html/rfc7518#section-3.4
func ES512() SignatureAlgorithm {
	return es512
}

// Signer is the interface used for generating a JWS signature.
type Signer interface {
	// SignatureAlgorithm is the signature algorithm used with this signer
	// as detailed in https://tools.ietf.org/html/rfc7518#section-3
	SignatureAlgorithm
	// Sign reads the given data and generates a signature of the content
	// using the specified SignatureAlgorithm as detailed in
	// https://tools.ietf.org/html/rfc7515#section-5.1
	Sign(data io.Reader) (signature string, err error)
}

type rsaSigner struct {
	key *rsa.PrivateKey
	SignatureAlgorithm
}

func (signer *rsaSigner) Sign(data io.Reader) (signature string, err error) {
	hasher := signer.Hash().New()
	if _, err = io.Copy(hasher, data); err != nil {
		return "", fmt.Errorf("error reading data to sign: %s", err)
	}
	hashed := hasher.Sum(nil)

	sigBytes, err := rsa.SignPKCS1v15(rand.Reader, signer.key, signer.Hash(), hashed)
	if err != nil {
		return "", fmt.Errorf("unable to sign hashed data: %s", err)
	}

	return Base64URLEncode(sigBytes), nil
}

type ecdsaSigner struct {
	key *ecdsa.PrivateKey
	SignatureAlgorithm
}

func (signer *ecdsaSigner) Sign(data io.Reader) (signature string, err error) {
	hasher := signer.Hash().New()
	if _, err = io.Copy(hasher, data); err != nil {
		return "", fmt.Errorf("error reading data to sign: %s", err)
	}
	hashed := hasher.Sum(nil)

	r, s, err := ecdsa.Sign(rand.Reader, signer.key, hashed)
	if err != nil {
		return "", fmt.Errorf("unable to sign hashed data: %s", err)
	}

	rBytes := r.Bytes()
	sBytes := s.Bytes()
	octetLength := (signer.key.PublicKey.Params().BitSize + 7) >> 3
	// We MUST include leading zeros in the output.
	rBytes = append(make([]byte, octetLength-len(rBytes), octetLength), rBytes...)
	sBytes = append(make([]byte, octetLength-len(sBytes), octetLength), sBytes...)

	return Base64URLEncode(append(rBytes, sBytes...)), nil
}

// Verifier is the interface used for verifying a JWS signature.
type Verifier interface {
	// SignatureAlgorithm is the signature algorithm used with this
	// verifier as detailed in
	// https://tools.ietf.org/html/rfc7518#section-3
	SignatureAlgorithm
	// Verify reads the given data and checks it against the given
	// base64url-encoded signature using the specified SignatureAlgorithm
	// as detailed in https://tools.ietf.org/html/rfc7515#section-5.2
	Verify(data io.Reader, signature string) error
}

type rsaVerifier struct {
	key *rsa.PublicKey
	SignatureAlgorithm
}

func (verifier *rsaVerifier) Verify(data io.Reader, signature string) error {
	sigBytes, err := Base64URLDecode(signature)
	if err != nil {
		return fmt.Errorf("unable to decode signature: %s", err)
	}

	hasher := verifier.Hash().New()
	if _, err = io.Copy(hasher, data); err != nil {
		return fmt.Errorf("error reading data to verify: %s", err)
	}
	hashed := hasher.Sum(nil)

	if err := rsa.VerifyPKCS1v15(verifier.key, verifier.Hash(), hashed, sigBytes); err != nil {
		return fmt.Errorf("invalid %s signature: %s", verifier.Name(), err)
	}

	return nil
}

type ecdsaVerifier struct {
	key *ecdsa.PublicKey
	SignatureAlgorithm
}

func (verifier *ecdsaVerifier) Verify(data io.Reader, signature string) error {
	// signature is the concatenation of (r, s), base64Url encoded.
	sigBytes, err := Base64URLDecode(signature)
	if err != nil {
		return fmt.Errorf("unable to decode signature: %s", err)
	}

	sigLength := len(sigBytes)
	expectedOctetLength := 2 * ((verifier.key.Params().BitSize + 7) >> 3)
	if sigLength != expectedOctetLength {
		return fmt.Errorf("%s signature length is %d octets long, should be %d", verifier.Name(), sigLength, expectedOctetLength)
	}

	rBytes := sigBytes[:sigLength/2]
	sBytes := sigBytes[sigLength/2:]
	r := new(big.Int).SetBytes(rBytes)
	s := new(big.Int).SetBytes(sBytes)

	hasher := verifier.Hash().New()
	if _, err = io.Copy(hasher, data); err != nil {
		return fmt.Errorf("error reading data to verify: %s", err)
	}
	hashed := hasher.Sum(nil)

	if !ecdsa.Verify(verifier.key, hashed, r, s) {
		return fmt.Errorf("invalid %s signature", verifier.Name())
	}

	return nil
}
