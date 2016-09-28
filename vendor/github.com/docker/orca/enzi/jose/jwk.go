package jose

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"
)

// JWK is used to represent a JSON Web Keys with fields for either RSA or ECDSA
// public keys.
type JWK struct {
	// Key ID value is arbitrary but should uniquely identify the public
	// key, usually via a hash.
	ID string `json:"kid"`
	// Key Type is either "RSA" or "EC"
	KeyType string `json:"kty"`

	// Fields for RSA public keys.
	Modulus  string `json:"n,omitempty"`
	Exponent string `json:"e,omitempty"`

	// Fields for ECDSA public keys. Curve is one of "P-256", "P-384", and
	// "P-521".
	Curve       string `json:"crv,omitempty"`
	XCoordinate string `json:"x,omitempty"`
	YCoordinate string `json:"y,omitempty"`
}

func (k *JWK) toRSA() (*rsa.PublicKey, []SignatureAlgorithm, error) {
	// First, convert the Modulus from base64url to *big.Int
	modulusBytes, err := Base64URLDecode(k.Modulus)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid RSA Modulus: %s", err)
	}

	modulus := new(big.Int).SetBytes(modulusBytes)

	// Next, convert the public exponent from base64url to int
	exponentBytes, err := Base64URLDecode(k.Exponent)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid RSA Exponent: %s", err)
	}
	// Only the minimum number of bytes were used to represent the
	// exponent, but binary.BigEndian.Uint32 expects at least 4 bytes, so
	// we need to add zero padding if necassary.
	exponentBytes = append(make([]byte, 4-len(exponentBytes), 4), exponentBytes...)

	exponent := int(binary.BigEndian.Uint32(exponentBytes))

	return &rsa.PublicKey{
		N: modulus,
		E: exponent,
	}, []SignatureAlgorithm{rs256, rs384, rs512}, nil
}

func (k *JWK) toECDSA() (*ecdsa.PublicKey, []SignatureAlgorithm, error) {
	// First, determine the Curve.
	var (
		curve  elliptic.Curve
		sigAlg SignatureAlgorithm
	)
	switch k.Curve {
	case "P-256":
		curve = elliptic.P256()
		sigAlg = es256
	case "P-384":
		curve = elliptic.P384()
		sigAlg = es384
	case "P-521":
		curve = elliptic.P521()
		sigAlg = es512
	default:
		return nil, nil, fmt.Errorf("unknown elliptic curve: %q", k.Curve)
	}

	// Determine the curve coordinate byte length to validate the X and Y
	// coordinate values.
	curveByteLen := (curve.Params().BitSize + 7) >> 3

	// Next, convert the X Coordinate from base64 to *big.Int
	xBytes, err := Base64URLDecode(k.XCoordinate)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid ECDSA X Coordinate: %s", err)
	}
	if len(xBytes) != curveByteLen {
		return nil, nil, fmt.Errorf("invalid ECDSA X Coordinate byte length: got %d, expected %d", len(xBytes), curveByteLen)
	}

	x := new(big.Int).SetBytes(xBytes)

	// Next, convert the Y Coordinate from base64 to *big.Int
	yBytes, err := Base64URLDecode(k.YCoordinate)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid ECDSA Y Coordinate: %s", err)
	}
	if len(yBytes) != curveByteLen {
		return nil, nil, fmt.Errorf("invalid ECDSA Y Coordinate byte length: got %d, expected %d", len(yBytes), curveByteLen)
	}

	y := new(big.Int).SetBytes(yBytes)

	return &ecdsa.PublicKey{
		Curve: curve,
		X:     x,
		Y:     y,
	}, []SignatureAlgorithm{sigAlg}, nil
}

// PublicKey represents an RSA or ECDSA JWK which can be used to verify JSON
// Web Signatures.
type PublicKey struct {
	JWK

	publicKey crypto.PublicKey
	// There is always at least one supported signing algorithm.
	sigAlgs []SignatureAlgorithm
}

// JWKSet represents a set of Public Keys.
type JWKSet struct {
	Keys []*PublicKey `json:"keys"`
}

// MarshalJSON marshals this PublicKey into JSON as a JWK.
func (k *PublicKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(k.JWK)
}

// NewPublicKeyJWK converts the given JWK into a PublicKey which can be used to
// verify JWS signatures.
func NewPublicKeyJWK(jwk JWK) (*PublicKey, error) {
	pubKey := &PublicKey{
		JWK: jwk,
	}

	var err error
	switch jwk.KeyType {
	case "RSA":
		pubKey.publicKey, pubKey.sigAlgs, err = jwk.toRSA()
	case "EC":
		pubKey.publicKey, pubKey.sigAlgs, err = jwk.toECDSA()
	default:
		err = fmt.Errorf("unknown JWK key type: %q", jwk.KeyType)
	}

	if err != nil {
		return nil, fmt.Errorf("unable to decode JWK public key: %s", err)
	}

	return pubKey, nil
}

// UnmarshalJSON unmarshals a JWK from the given JSON-encoded data and sets the
// public key value and list of supported signature algorithms.
func (k *PublicKey) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &k.JWK); err != nil {
		return err
	}

	pubKey, err := NewPublicKeyJWK(k.JWK)
	if pubKey != nil {
		*k = *pubKey
	}

	return err
}

// NewPublicKey converts the given crypto public key (either RSA or ECDSA) into
// a PublicKey which can be used to verify JWS signatures.
func NewPublicKey(publicKey crypto.PublicKey) (*PublicKey, error) {
	switch publicKey := publicKey.(type) {
	case *rsa.PublicKey:
		return newRSAPublicKey(publicKey)
	case *ecdsa.PublicKey:
		return newECDSAPublicKey(publicKey)
	default:
		return nil, fmt.Errorf("unable to handle public key type: %T", publicKey)
	}
}

func newRSAPublicKey(publicKey *rsa.PublicKey) (*PublicKey, error) {
	// First, encode the public key Modulus to base64url
	modulus := Base64URLEncode(publicKey.N.Bytes())

	// Next, encode the public key Exponent to base64url.
	exponentBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(exponentBytes, uint32(publicKey.E))
	// We MUST use the minimum number of octets to represent the exponent
	// value. Note that if it happens to be zero, it should be 1 zero
	// octet. It is currently 4 octets, we just need to trim off up to 3
	// preceeding zero octets.
	var i int
	for i < 3 && exponentBytes[i] == 0 {
		i++
	}
	exponent := Base64URLEncode(exponentBytes[i:])

	// Finally, generate the keyID.
	keyID, err := makeKeyID(publicKey)
	if err != nil {
		return nil, fmt.Errorf("unable to generate key ID: %s", err)
	}

	return &PublicKey{
		JWK: JWK{
			ID:       keyID,
			KeyType:  "RSA",
			Modulus:  modulus,
			Exponent: exponent,
		},
		publicKey: publicKey,
		sigAlgs:   []SignatureAlgorithm{rs256, rs384, rs512},
	}, nil
}

func newECDSAPublicKey(publicKey *ecdsa.PublicKey) (*PublicKey, error) {
	// First, determine the named curve.
	var (
		curve  string
		sigAlg signatureAlg
	)

	switch publicKey.Curve.Params().Name {
	case elliptic.P256().Params().Name:
		curve = "P-256"
		sigAlg = es256
	case elliptic.P384().Params().Name:
		curve = "P-384"
		sigAlg = es384
	case elliptic.P521().Params().Name:
		curve = "P-521"
		sigAlg = es512
	default:
		return nil, fmt.Errorf("unknown curve name: %q", publicKey.Curve.Params().Name)
	}

	// Next, encode the X and Y Coordinates to base64url.
	xBytes := publicKey.X.Bytes()
	yBytes := publicKey.Y.Bytes()
	octetLength := (publicKey.Params().BitSize + 7) >> 3
	// We MUST include leading zeros in the output so that x, y are each
	// octetLength bytes long.
	xBytes = append(make([]byte, octetLength-len(xBytes), octetLength), xBytes...)
	yBytes = append(make([]byte, octetLength-len(yBytes), octetLength), yBytes...)

	xCoordinate := Base64URLEncode(xBytes)
	yCoordinate := Base64URLEncode(yBytes)

	// Finally, generate the keyID.
	keyID, err := makeKeyID(publicKey)
	if err != nil {
		return nil, fmt.Errorf("unable to generate key ID: %s", err)
	}

	return &PublicKey{
		JWK: JWK{
			ID:          keyID,
			KeyType:     "EC",
			Curve:       curve,
			XCoordinate: xCoordinate,
			YCoordinate: yCoordinate,
		},
		publicKey: publicKey,
		sigAlgs:   []SignatureAlgorithm{sigAlg},
	}, nil
}

// CryptoPublicKey returns the crypto.PublicKey value for this JWK.
func (k *PublicKey) CryptoPublicKey() crypto.PublicKey {
	return k.publicKey
}

// SignatureAlgs returns a slice of identifiers of JWS signing algorithms which
// can be verified by this JWK (or used to sign if used with the private key).
func (k *PublicKey) SignatureAlgs() []SignatureAlgorithm {
	return k.sigAlgs
}

// getSignatureAlg iterates through the given list of signing algorithm names
// and returns the first matching signing algorithm which is supported by this
// key. A non-nil error is returned if no alg is supported.
func (k *PublicKey) getSignatureAlg(algs ...string) (SignatureAlgorithm, error) {
	sigAlgSet := make(map[string]SignatureAlgorithm, len(k.sigAlgs))
	for _, sigAlg := range k.sigAlgs {
		sigAlgSet[sigAlg.Name()] = sigAlg
	}

	for _, alg := range algs {
		if sigAlg, ok := sigAlgSet[alg]; ok {
			return sigAlg, nil
		}
	}

	// none of the given algorithms are supported, generate a list of valid
	// algs to include in the error message.
	choices := make([]string, len(k.sigAlgs))
	for i, validAlg := range k.sigAlgs {
		choices[i] = validAlg.Name()
	}

	return nil, fmt.Errorf("signing algorithims %q not supported by this key - valid choices are: %s", algs, choices)
}

// Verifier creates a verifier using this Public Key JWK. The returned
// Verifier will use the hash function designated by given alg to verify a
// signature. The given alg value should correspond to the name of one of this
// JWK's supported SignatureAlgs otherwise an error is returned indicating so.
func (k *PublicKey) Verifier(alg string) (Verifier, error) {
	sigAlg, err := k.getSignatureAlg(alg)
	if err != nil {
		return nil, fmt.Errorf("unable to load signature algorithm: %s", err)
	}

	switch publicKey := k.publicKey.(type) {
	case *rsa.PublicKey:
		return &rsaVerifier{
			key:                publicKey,
			SignatureAlgorithm: sigAlg,
		}, nil
	case *ecdsa.PublicKey:
		return &ecdsaVerifier{
			key:                publicKey,
			SignatureAlgorithm: sigAlg,
		}, nil
	default:
		return nil, fmt.Errorf("unable to handle public key type: %T", publicKey)
	}
}

// PrivateKey represents an RSA or ECDSA JWK which can be used to generate JSON
// Web Signatures.
type PrivateKey struct {
	PublicKey

	privateKey crypto.PrivateKey
}

// NewPrivateKey generates a new JWK private key from the given private key
// which must be either of type (*rsa.PrivateKey) or (*ecdsa.PrivateKey). The
// resulting *PrivateKey can be used to generate JWS signatures.
func NewPrivateKey(privateKey crypto.PrivateKey) (*PrivateKey, error) {
	switch privateKey := privateKey.(type) {
	case *rsa.PrivateKey:
		return newRSAPrivateKey(privateKey)
	case *ecdsa.PrivateKey:
		return newECDSAPrivateKey(privateKey)
	default:
		return nil, fmt.Errorf("unable to handle private key type: %T", privateKey)
	}
}

func newRSAPrivateKey(privateKey *rsa.PrivateKey) (*PrivateKey, error) {
	publicKey, err := newRSAPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, err
	}

	return &PrivateKey{
		PublicKey:  *publicKey,
		privateKey: privateKey,
	}, nil
}

func newECDSAPrivateKey(privateKey *ecdsa.PrivateKey) (*PrivateKey, error) {
	publicKey, err := newECDSAPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, err
	}

	return &PrivateKey{
		PublicKey:  *publicKey,
		privateKey: privateKey,
	}, nil
}

// CryptoPrivateKey returns the crypto.PrivateKey value for this JWK.
func (k *PrivateKey) CryptoPrivateKey() crypto.PrivateKey {
	return k.privateKey
}

// Signer creates a Signer using this Private Key JWK. If no algorithm name is
// specified, a default signing algorithm is used: the first in this key's list
// of supported signing algorithms. If one or more algorithm names is
// specified, the first matching compatible signing algorithm is used. If no
// signing algorithms are available then an error is returned indicating so.
// The returned Signer will use the hash function designated by the selected
// algorithm to sign data.
func (k *PrivateKey) Signer(algs ...string) (Signer, error) {
	var (
		sigAlg = k.sigAlgs[0] // The default signing algorithm.
		err    error
	)

	// If any algs were specified, try to load one of those.
	if len(algs) > 0 {
		sigAlg, err = k.getSignatureAlg(algs...)
		if err != nil {
			return nil, fmt.Errorf("unable to load signature algorithm: %s", err)
		}
	}

	switch privateKey := k.privateKey.(type) {
	case *rsa.PrivateKey:
		return &rsaSigner{
			key:                privateKey,
			SignatureAlgorithm: sigAlg,
		}, nil
	case *ecdsa.PrivateKey:
		return &ecdsaSigner{
			key:                privateKey,
			SignatureAlgorithm: sigAlg,
		}, nil
	default:
		return nil, fmt.Errorf("unable to handle private key type: %T", privateKey)
	}
}
