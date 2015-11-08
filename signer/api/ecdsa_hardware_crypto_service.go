// +build pkcs11

package api

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"math/big"

	"github.com/Sirupsen/logrus"
	"github.com/docker/notary/passphrase"
	"github.com/docker/notary/trustmanager"
	"github.com/docker/notary/tuf/data"
	"github.com/miekg/pkcs11"
)

const (
	USER_PIN    = "123456"
	SO_USER_PIN = "010203040506070801020304050607080102030405060708"
	numSlots    = 4 // number of slots in the yubikey

	KeymodeNone      = 0
	KeymodeTouch     = 1 // touch enabled
	KeymodePinOnce   = 2 // require pin entry once
	KeymodePinAlways = 4 // require pin entry all the time
)

// what key mode to use when generating keys
var yubikeyKeymode = KeymodeTouch | KeymodePinOnce

// SetYubikeyKeyMode - sets the mode when generating yubikey keys.
// This is to be used for testing.  It does nothing if not building with tag
// pkcs11.
func SetYubikeyKeyMode(keyMode int) error {
	// technically 7 (1 | 2 | 4) is valid, but KeymodePinOnce +
	// KeymdoePinAlways don't really make sense together
	if keyMode < 0 || keyMode > 5 {
		return errors.New("Invalid key mode")
	}
	yubikeyKeymode = keyMode
	return nil
}

// Hardcoded yubikey PKCS11 ID
var YUBIKEY_ROOT_KEY_ID = []byte{2}

type ErrBackupFailed struct {
	err string
}

func (err ErrBackupFailed) Error() string {
	return fmt.Sprintf("Failed to backup private key to: %s", err.err)
}

type yubiSlot struct {
	role   string
	slotID []byte
}

// YubiPrivateKey represents a private key inside of a yubikey
type YubiPrivateKey struct {
	data.ECDSAPublicKey
	passRetriever passphrase.Retriever
	slot          []byte
}

type YubikeySigner struct {
	YubiPrivateKey
}

func NewYubiPrivateKey(slot []byte, pubKey data.ECDSAPublicKey, passRetriever passphrase.Retriever) *YubiPrivateKey {
	return &YubiPrivateKey{
		ECDSAPublicKey: pubKey,
		passRetriever:  passRetriever,
		slot:           slot,
	}
}

func (ys *YubikeySigner) Public() crypto.PublicKey {
	publicKey, err := x509.ParsePKIXPublicKey(ys.YubiPrivateKey.Public())
	if err != nil {
		return nil
	}

	return publicKey
}

// CryptoSigner returns a crypto.Signer tha wraps the YubiPrivateKey. Needed for
// Certificate generation only
func (y *YubiPrivateKey) CryptoSigner() crypto.Signer {
	return &YubikeySigner{YubiPrivateKey: *y}
}

// Private is not implemented in hardware  keys
func (y *YubiPrivateKey) Private() []byte {
	// We cannot return the private material from a Yubikey
	// TODO(david): We probably want to return an error here
	return nil
}

func (y YubiPrivateKey) SignatureAlgorithm() data.SigAlgorithm {
	return data.ECDSASignature
}

func (y *YubiPrivateKey) Sign(rand io.Reader, msg []byte, opts crypto.SignerOpts) ([]byte, error) {
	ctx, session, err := SetupHSMEnv(pkcs11Lib)
	if err != nil {
		return nil, err
	}
	defer cleanup(ctx, session)

	sig, err := sign(ctx, session, y.slot, y.passRetriever, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to sign using Yubikey: %v", err)
	}

	return sig, nil
}

// addECDSAKey adds a key to the yubikey
func addECDSAKey(
	ctx *pkcs11.Ctx,
	session pkcs11.SessionHandle,
	privKey data.PrivateKey,
	pkcs11KeyID []byte,
	passRetriever passphrase.Retriever,
	role string,
	backupStore trustmanager.KeyStore,
) error {
	logrus.Debugf("Got into add key with key: %s\n", privKey.ID())

	// TODO(diogo): Figure out CKU_SO with yubikey
	err := login(ctx, session, passRetriever, pkcs11.CKU_SO, SO_USER_PIN)
	if err != nil {
		return err
	}
	defer ctx.Logout(session)

	// Create an ecdsa.PrivateKey out of the private key bytes
	ecdsaPrivKey, err := x509.ParseECPrivateKey(privKey.Private())
	if err != nil {
		return err
	}

	ecdsaPrivKeyD := ecdsaPrivKey.D.Bytes()
	logrus.Debugf("Getting D bytes: %v\n", ecdsaPrivKeyD)

	template, err := trustmanager.NewCertificate(role)
	if err != nil {
		return fmt.Errorf("failed to create the certificate template: %v", err)
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, template, template, ecdsaPrivKey.Public(), ecdsaPrivKey)
	if err != nil {
		return fmt.Errorf("failed to create the certificate: %v", err)
	}

	certTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_CERTIFICATE),
		pkcs11.NewAttribute(pkcs11.CKA_VALUE, certBytes),
		pkcs11.NewAttribute(pkcs11.CKA_ID, pkcs11KeyID),
	}

	privateKeyTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PRIVATE_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_ECDSA),
		pkcs11.NewAttribute(pkcs11.CKA_ID, pkcs11KeyID),
		pkcs11.NewAttribute(pkcs11.CKA_EC_PARAMS, []byte{0x06, 0x08, 0x2a, 0x86, 0x48, 0xce, 0x3d, 0x03, 0x01, 0x07}),
		pkcs11.NewAttribute(pkcs11.CKA_VALUE, ecdsaPrivKeyD),
		pkcs11.NewAttribute(pkcs11.CKA_VENDOR_DEFINED, yubikeyKeymode),
	}

	_, err = ctx.CreateObject(session, certTemplate)
	if err != nil {
		return fmt.Errorf("error importing: %v", err)
	}

	_, err = ctx.CreateObject(session, privateKeyTemplate)
	if err != nil {
		return fmt.Errorf("error importing: %v", err)
	}

	err = backupStore.AddKey(privKey.ID(), role, privKey)
	if err != nil {
		return ErrBackupFailed{err: err.Error()}
	}

	return nil
}

func getECDSAKey(ctx *pkcs11.Ctx, session pkcs11.SessionHandle, pkcs11KeyID []byte) (*data.ECDSAPublicKey, string, error) {
	findTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_TOKEN, true),
		pkcs11.NewAttribute(pkcs11.CKA_ID, pkcs11KeyID),
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PUBLIC_KEY),
	}

	attrTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, []byte{0}),
		pkcs11.NewAttribute(pkcs11.CKA_EC_POINT, []byte{0}),
		pkcs11.NewAttribute(pkcs11.CKA_EC_PARAMS, []byte{0}),
	}

	if err := ctx.FindObjectsInit(session, findTemplate); err != nil {
		logrus.Debugf("Failed to init: %s\n", err.Error())
		return nil, "", err
	}
	obj, b, err := ctx.FindObjects(session, 1)
	if err != nil {
		logrus.Debugf("Failed to find: %s %v\n", err.Error(), b)
		return nil, "", err
	}
	if err := ctx.FindObjectsFinal(session); err != nil {
		logrus.Debugf("Failed to finalize: %s\n", err.Error())
		return nil, "", err
	}
	if len(obj) != 1 {
		logrus.Debugf("should have found one object")
		return nil, "", errors.New("no matching keys found inside of yubikey")
	}

	// Retrieve the public-key material to be able to create a new HSMRSAKey
	attr, err := ctx.GetAttributeValue(session, obj[0], attrTemplate)
	if err != nil {
		logrus.Debugf("Failed to get Attribute for: %v\n", obj[0])
		return nil, "", err
	}

	// Iterate through all the attributes of this key and saves CKA_PUBLIC_EXPONENT and CKA_MODULUS. Removes ordering specific issues.
	var rawPubKey []byte
	for _, a := range attr {
		if a.Type == pkcs11.CKA_EC_POINT {
			rawPubKey = a.Value
		}

	}

	ecdsaPubKey := ecdsa.PublicKey{Curve: elliptic.P256(), X: new(big.Int).SetBytes(rawPubKey[3:35]), Y: new(big.Int).SetBytes(rawPubKey[35:])}
	pubBytes, err := x509.MarshalPKIXPublicKey(&ecdsaPubKey)
	if err != nil {
		logrus.Debugf("Failed to Marshal public key")
		return nil, "", err
	}

	// TODO(diogo): Actually get the right alias from the certificate instead of
	// alwars returning data.CanonicalRootRole
	return data.NewECDSAPublicKey(pubBytes), data.CanonicalRootRole, nil
}

// Sign returns a signature for a given signature request
func sign(ctx *pkcs11.Ctx, session pkcs11.SessionHandle, pkcs11KeyID []byte, passRetriever passphrase.Retriever, payload []byte) ([]byte, error) {
	err := login(ctx, session, passRetriever, pkcs11.CKU_USER, USER_PIN)
	if err != nil {
		return nil, fmt.Errorf("error logging in: %v", err)
	}
	defer ctx.Logout(session)

	// Define the ECDSA Private key template
	class := pkcs11.CKO_PRIVATE_KEY
	privateKeyTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, class),
		pkcs11.NewAttribute(pkcs11.CKA_KEY_TYPE, pkcs11.CKK_ECDSA),
		pkcs11.NewAttribute(pkcs11.CKA_ID, pkcs11KeyID),
	}

	if err := ctx.FindObjectsInit(session, privateKeyTemplate); err != nil {
		return nil, err
	}
	obj, _, err := ctx.FindObjects(session, 1)
	if err != nil {
		return nil, err
	}
	if err = ctx.FindObjectsFinal(session); err != nil {
		return nil, err
	}
	if len(obj) != 1 {
		return nil, errors.New("length of objects found not 1")
	}

	var sig []byte
	ctx.SignInit(session, []*pkcs11.Mechanism{pkcs11.NewMechanism(pkcs11.CKM_ECDSA, nil)}, obj[0])

	// Get the SHA256 of the payload
	digest := sha256.Sum256(payload)

	fmt.Println("Please touch the attached Yubikey to perform signing.")
	sig, err = ctx.Sign(session, digest[:])
	if err != nil {
		logrus.Debugf("Error while signing: %s", err)
		return nil, err
	}

	if sig == nil {
		return nil, errors.New("Failed to create signature")
	}
	return sig[:], nil
}

func removeKey(ctx *pkcs11.Ctx, session pkcs11.SessionHandle, pkcs11KeyID []byte, passRetriever passphrase.Retriever, keyID string) error {
	err := login(ctx, session, passRetriever, pkcs11.CKU_SO, SO_USER_PIN)
	if err != nil {
		return err
	}
	defer ctx.Logout(session)

	template := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_TOKEN, true),
		pkcs11.NewAttribute(pkcs11.CKA_ID, pkcs11KeyID),
		//pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_PRIVATE_KEY),
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_CERTIFICATE),
	}

	if err := ctx.FindObjectsInit(session, template); err != nil {
		logrus.Printf("Failed to init: %s\n", err.Error())
		return err
	}
	obj, b, err := ctx.FindObjects(session, 1)
	if err != nil {
		logrus.Printf("Failed to find: %s %v\n", err.Error(), b)
		return err
	}
	if err := ctx.FindObjectsFinal(session); err != nil {
		logrus.Printf("Failed to finalize: %s\n", err.Error())
		return err
	}
	if len(obj) != 1 {
		logrus.Printf("should have found one object")
		return err
	}

	// Delete the certificate
	err = ctx.DestroyObject(session, obj[0])
	if err != nil {
		logrus.Printf("Failed to delete cert")
		return err
	}
	return nil
}

func listKeys(ctx *pkcs11.Ctx, session pkcs11.SessionHandle) (keys map[string]yubiSlot, err error) {
	keys = make(map[string]yubiSlot)
	findTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_TOKEN, true),
		//pkcs11.NewAttribute(pkcs11.CKA_ID, pkcs11KeyID),
		pkcs11.NewAttribute(pkcs11.CKA_CLASS, pkcs11.CKO_CERTIFICATE),
	}

	attrTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_ID, []byte{0}),
		pkcs11.NewAttribute(pkcs11.CKA_VALUE, []byte{0}),
	}

	if err = ctx.FindObjectsInit(session, findTemplate); err != nil {
		logrus.Debugf("Failed to init: %s\n", err.Error())
		return
	}
	objs, b, err := ctx.FindObjects(session, numSlots)
	for err == nil {
		var o []pkcs11.ObjectHandle
		o, b, err = ctx.FindObjects(session, numSlots)
		if err != nil {
			continue
		}
		if len(o) == 0 {
			break
		}
		objs = append(objs, o...)
	}
	if err != nil {
		logrus.Debugf("Failed to find: %s %v\n", err.Error(), b)
		if len(objs) == 0 {
			return nil, err
		}
	}
	if err = ctx.FindObjectsFinal(session); err != nil {
		logrus.Debugf("Failed to finalize: %s\n", err.Error())
		return
	}
	if len(objs) == 0 {
		return nil, errors.New("No keys found in yubikey.")
	}
	logrus.Debugf("Found %d objects matching list filters", len(objs))
	for _, obj := range objs {
		var (
			cert *x509.Certificate
			slot []byte
		)
		// Retrieve the public-key material to be able to create a new HSMRSAKey
		attr, err := ctx.GetAttributeValue(session, obj, attrTemplate)
		if err != nil {
			logrus.Debugf("Failed to get Attribute for: %v\n", obj)
			continue
		}

		// Iterate through all the attributes of this key and saves CKA_PUBLIC_EXPONENT and CKA_MODULUS. Removes ordering specific issues.
		for _, a := range attr {
			if a.Type == pkcs11.CKA_ID {
				slot = a.Value
			}
			if a.Type == pkcs11.CKA_VALUE {
				cert, err = x509.ParseCertificate(a.Value)
				if err != nil {
					continue
				}
				if !data.ValidRole(cert.Subject.CommonName) {
					continue
				}
			}
		}
		var ecdsaPubKey *ecdsa.PublicKey
		switch cert.PublicKeyAlgorithm {
		case x509.ECDSA:
			ecdsaPubKey = cert.PublicKey.(*ecdsa.PublicKey)
		default:
			logrus.Infof("Unsupported x509 PublicKeyAlgorithm: %d", cert.PublicKeyAlgorithm)
			continue
		}

		pubBytes, err := x509.MarshalPKIXPublicKey(ecdsaPubKey)
		if err != nil {
			logrus.Debugf("Failed to Marshal public key")
			continue
		}

		keys[data.NewECDSAPublicKey(pubBytes).ID()] = yubiSlot{
			role:   cert.Subject.CommonName,
			slotID: slot,
		}
	}
	return
}

func getNextEmptySlot(ctx *pkcs11.Ctx, session pkcs11.SessionHandle) ([]byte, error) {
	findTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_TOKEN, true),
	}
	attrTemplate := []*pkcs11.Attribute{
		pkcs11.NewAttribute(pkcs11.CKA_ID, []byte{0}),
	}

	if err := ctx.FindObjectsInit(session, findTemplate); err != nil {
		logrus.Debugf("Failed to init: %s\n", err.Error())
		return nil, err
	}
	objs, b, err := ctx.FindObjects(session, numSlots)
	for err == nil {
		var o []pkcs11.ObjectHandle
		o, b, err = ctx.FindObjects(session, numSlots)
		if err != nil {
			continue
		}
		if len(o) == 0 {
			break
		}
		objs = append(objs, o...)
	}
	taken := make(map[int]bool)
	if err != nil {
		logrus.Debugf("Failed to find: %s %v\n", err.Error(), b)
		return nil, err
	}
	for _, obj := range objs {
		// Retrieve the slot ID
		attr, err := ctx.GetAttributeValue(session, obj, attrTemplate)
		if err != nil {
			logrus.Debugf("Failed to get Attribute for: %v\n", obj)
			continue
		}

		// Iterate through attributes. If an ID attr was found, mark it as taken
		for _, a := range attr {
			if a.Type == pkcs11.CKA_ID {
				if len(a.Value) < 1 {
					continue
				}
				// a byte will always be capable of representing all slot IDs
				// for the Yubikeys
				slotNum := int(a.Value[0])
				if slotNum >= len(taken) {
					// defensive
					continue
				}
				taken[slotNum] = true
			}
		}
	}
	for i := 0; i < numSlots; i++ {
		if !taken[i] {
			return []byte{byte(i)}, nil
		}
	}
	return nil, errors.New("Yubikey has no available slots.")
}

type YubiKeyStore struct {
	passRetriever passphrase.Retriever
	keys          map[string]yubiSlot
	backupStore   trustmanager.KeyStore
}

func NewYubiKeyStore(backupStore trustmanager.KeyStore, passphraseRetriever passphrase.Retriever) (*YubiKeyStore, error) {
	s := &YubiKeyStore{
		passRetriever: passphraseRetriever,
		keys:          make(map[string]yubiSlot),
		backupStore:   backupStore,
	}
	s.ListKeys() // populate keys field
	return s, nil
}

func (s *YubiKeyStore) ListKeys() map[string]string {
	if len(s.keys) > 0 {
		return buildKeyMap(s.keys)
	}
	ctx, session, err := SetupHSMEnv(pkcs11Lib)
	if err != nil {
		return nil
	}
	defer cleanup(ctx, session)
	keys, err := listKeys(ctx, session)
	if err != nil {
		return nil
	}
	s.keys = keys
	return buildKeyMap(keys)
}

func (s *YubiKeyStore) AddKey(keyID, role string, privKey data.PrivateKey) error {
	// We only allow adding root keys for now
	if role != data.CanonicalRootRole {
		return fmt.Errorf("yubikey only supports storing root keys, got %s for key: %s\n", role, keyID)
	}

	ctx, session, err := SetupHSMEnv(pkcs11Lib)
	if err != nil {
		return err
	}
	defer cleanup(ctx, session)

	if k, ok := s.keys[keyID]; ok {
		if k.role == role {
			// already have the key and it's associated with the correct role
			return nil
		}
	}

	slot, err := getNextEmptySlot(ctx, session)
	if err != nil {
		return err
	}
	logrus.Debugf("Using yubikey slot %v", slot)
	err = addECDSAKey(ctx, session, privKey, slot, s.passRetriever, role, s.backupStore)
	if err == nil {
		s.keys[privKey.ID()] = yubiSlot{
			role:   role,
			slotID: slot,
		}
		return nil
	}
	return err
}

func (s *YubiKeyStore) GetKey(keyID string) (data.PrivateKey, string, error) {
	ctx, session, err := SetupHSMEnv(pkcs11Lib)
	if err != nil {
		return nil, "", err
	}
	defer cleanup(ctx, session)

	key, ok := s.keys[keyID]
	if !ok {
		return nil, "", errors.New("no matching keys found inside of yubikey")
	}

	pubKey, alias, err := getECDSAKey(ctx, session, key.slotID)
	if err != nil {
		return nil, "", err
	}
	// Check to see if we're returning the intended keyID
	if pubKey.ID() != keyID {
		return nil, "", fmt.Errorf("expected root key: %s, but found: %s\n", keyID, pubKey.ID())
	}
	privKey := NewYubiPrivateKey(key.slotID, *pubKey, s.passRetriever)
	if privKey == nil {
		return nil, "", errors.New("could not initialize new YubiPrivateKey")
	}

	return privKey, alias, err
}

func (s *YubiKeyStore) RemoveKey(keyID string) error {
	ctx, session, err := SetupHSMEnv(pkcs11Lib)
	if err != nil {
		return nil
	}
	defer cleanup(ctx, session)
	key, ok := s.keys[keyID]
	if !ok {
		return errors.New("Key not present in yubikey")
	}
	return removeKey(ctx, session, key.slotID, s.passRetriever, keyID)
}

func (s *YubiKeyStore) ExportKey(keyID string) ([]byte, error) {
	logrus.Debugf("Attempting to export: %s key inside of YubiKeyStore", keyID)
	return nil, errors.New("Keys cannot be exported from a Yubikey.")
}

func (s *YubiKeyStore) ImportKey(pemBytes []byte, keyID string) error {
	logrus.Debugf("Attempting to import: %s key inside of YubiKeyStore", keyID)
	privKey, _, err := trustmanager.GetPasswdDecryptBytes(s.passRetriever, pemBytes, "imported", "root")
	if err != nil {
		return err
	}
	return s.AddKey(privKey.ID(), "root", privKey)
}

func cleanup(ctx *pkcs11.Ctx, session pkcs11.SessionHandle) {
	ctx.CloseSession(session)
	ctx.Finalize()
	ctx.Destroy()
}

// SetupHSMEnv is a method that depends on the existences
func SetupHSMEnv(libraryPath string) (*pkcs11.Ctx, pkcs11.SessionHandle, error) {
	p := pkcs11.New(libraryPath)

	if p == nil {
		return nil, 0, errors.New("Failed to init library")
	}

	if err := p.Initialize(); err != nil {
		return nil, 0, fmt.Errorf("Initialize error %s\n", err.Error())
	}

	slots, err := p.GetSlotList(true)
	if err != nil {
		return nil, 0, fmt.Errorf("Failed to list HSM slots %s", err)
	}
	// Check to see if we got any slots from the HSM.
	if len(slots) < 1 {
		return nil, 0, fmt.Errorf("No HSM Slots found")
	}

	// CKF_SERIAL_SESSION: TRUE if cryptographic functions are performed in serial with the application; FALSE if the functions may be performed in parallel with the application.
	// CKF_RW_SESSION: TRUE if the session is read/write; FALSE if the session is read-only
	session, err := p.OpenSession(slots[0], pkcs11.CKF_SERIAL_SESSION|pkcs11.CKF_RW_SESSION)
	if err != nil {
		return nil, 0, fmt.Errorf("Failed to Start Session with HSM %s", err)
	}

	return p, session, nil
}

// YubikeyEnabled returns true if a Yubikey can be accessed
func YubikeyAccessible() bool {
	ctx, session, err := SetupHSMEnv(pkcs11Lib)
	if err != nil {
		return false
	}
	defer cleanup(ctx, session)
	return true
}

func login(ctx *pkcs11.Ctx, session pkcs11.SessionHandle, passRetriever passphrase.Retriever, userFlag uint, defaultPassw string) error {
	// try default password
	err := ctx.Login(session, userFlag, defaultPassw)
	if err == nil {
		return nil
	}

	// default failed, ask user for password
	for attempts := 0; ; attempts++ {
		var (
			giveup bool
			err    error
			user   string
		)
		if userFlag == pkcs11.CKU_SO {
			user = "SO Pin"
		} else {
			user = "User Pin"
		}
		passwd, giveup, err := passRetriever(user, "yubikey", false, attempts)
		// Check if the passphrase retriever got an error or if it is telling us to give up
		if giveup || err != nil {
			return trustmanager.ErrPasswordInvalid{}
		}
		if attempts > 2 {
			return trustmanager.ErrAttemptsExceeded{}
		}

		// Try to convert PEM encoded bytes back to a PrivateKey using the passphrase
		err = ctx.Login(session, userFlag, passwd)
		if err == nil {
			return nil
		}
	}
	return nil
}

func buildKeyMap(keys map[string]yubiSlot) map[string]string {
	res := make(map[string]string)
	for k, v := range keys {
		res[k] = v.role
	}
	return res
}
