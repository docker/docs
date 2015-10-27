package api_test

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/docker/notary/cryptoservice"
	"github.com/docker/notary/signer"
	"github.com/docker/notary/signer/api"
	"github.com/docker/notary/trustmanager"
	"github.com/docker/notary/tuf/data"
	"github.com/miekg/pkcs11"
	"github.com/stretchr/testify/assert"

	pb "github.com/docker/notary/proto"
)

var (
	server              *httptest.Server
	reader              io.Reader
	deleteKeyBaseURL    string
	createKeyBaseURL    string
	keyInfoBaseURL      string
	signBaseURL         string
	passphraseRetriever = func(string, string, bool, int) (string, bool, error) { return "passphrase", false, nil }
)

func SetupHSMEnv(t *testing.T) (*pkcs11.Ctx, pkcs11.SessionHandle) {
	var libPath = "/usr/local/lib/softhsm/libsofthsm2.so"
	if _, err := os.Stat(libPath); err != nil {
		t.Skipf("Skipping test. Library path: %s does not exist", libPath)
	}

	p := pkcs11.New(libPath)

	if p == nil {
		t.Fatalf("Failed to init library")
	}

	if err := p.Initialize(); err != nil {
		t.Fatalf("Initialize error %s\n", err.Error())
	}

	slots, err := p.GetSlotList(true)
	if err != nil {
		t.Fatalf("Failed to list HSM slots %s", err)
	}

	session, err := p.OpenSession(slots[0], pkcs11.CKF_SERIAL_SESSION|pkcs11.CKF_RW_SESSION)
	if err != nil {
		t.Fatalf("Failed to Start Session with HSM %s", err)
	}

	if err = p.Login(session, pkcs11.CKU_USER, "1234"); err != nil {
		t.Fatalf("User PIN %s\n", err.Error())
	}

	return p, session
}

func setup(cryptoServices signer.CryptoServiceIndex) {
	server = httptest.NewServer(api.Handlers(cryptoServices))
	deleteKeyBaseURL = fmt.Sprintf("%s/delete", server.URL)
	createKeyBaseURL = fmt.Sprintf("%s/new", server.URL)
	keyInfoBaseURL = fmt.Sprintf("%s", server.URL)
	signBaseURL = fmt.Sprintf("%s/sign", server.URL)
}

func TestDeleteKeyHandlerReturns404WithNonexistentKey(t *testing.T) {
	keyStore := trustmanager.NewKeyMemoryStore(passphraseRetriever)
	cryptoService := cryptoservice.NewCryptoService("", keyStore)
	setup(signer.CryptoServiceIndex{data.ED25519Key: cryptoService, data.RSAKey: cryptoService, data.ECDSAKey: cryptoService})

	fakeID := "c62e6d68851cef1f7e55a9d56e3b0c05f3359f16838cad43600f0554e7d3b54d"

	keyID := &pb.KeyID{ID: fakeID}
	requestJson, _ := json.Marshal(keyID)
	reader = strings.NewReader(string(requestJson))

	request, err := http.NewRequest("POST", deleteKeyBaseURL, reader)
	assert.Nil(t, err)

	res, err := http.DefaultClient.Do(request)
	assert.Nil(t, err)

	assert.Equal(t, 404, res.StatusCode)
}

func TestDeleteKeyHandler(t *testing.T) {
	keyStore := trustmanager.NewKeyMemoryStore(passphraseRetriever)
	cryptoService := cryptoservice.NewCryptoService("", keyStore)
	setup(signer.CryptoServiceIndex{data.ED25519Key: cryptoService, data.RSAKey: cryptoService, data.ECDSAKey: cryptoService})

	tufKey, _ := cryptoService.Create("", data.ED25519Key)
	assert.NotNil(t, tufKey)

	requestJson, _ := json.Marshal(&pb.KeyID{ID: tufKey.ID()})
	reader = strings.NewReader(string(requestJson))

	request, err := http.NewRequest("POST", deleteKeyBaseURL, reader)
	assert.Nil(t, err)

	res, err := http.DefaultClient.Do(request)
	assert.Nil(t, err)

	assert.Equal(t, 200, res.StatusCode)
}

func TestKeyInfoHandler(t *testing.T) {
	keyStore := trustmanager.NewKeyMemoryStore(passphraseRetriever)
	cryptoService := cryptoservice.NewCryptoService("", keyStore)
	setup(signer.CryptoServiceIndex{data.ED25519Key: cryptoService, data.RSAKey: cryptoService, data.ECDSAKey: cryptoService})

	tufKey, _ := cryptoService.Create("", data.ED25519Key)
	assert.NotNil(t, tufKey)

	keyInfoURL := fmt.Sprintf("%s/%s", keyInfoBaseURL, tufKey.ID())

	request, err := http.NewRequest("GET", keyInfoURL, nil)
	assert.Nil(t, err)

	res, err := http.DefaultClient.Do(request)
	assert.Nil(t, err)

	jsonBlob, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)

	var pubKey *pb.PublicKey
	err = json.Unmarshal(jsonBlob, &pubKey)
	assert.Nil(t, err)

	assert.Equal(t, tufKey.ID(), pubKey.KeyInfo.KeyID.ID)
	assert.Equal(t, 200, res.StatusCode)
}

func TestKeyInfoHandlerReturns404WithNonexistentKey(t *testing.T) {
	// We associate both key types with this signing service to bypass the
	// ID -> keyType logic in the tests
	keyStore := trustmanager.NewKeyMemoryStore(passphraseRetriever)
	cryptoService := cryptoservice.NewCryptoService("", keyStore)
	setup(signer.CryptoServiceIndex{data.ED25519Key: cryptoService, data.RSAKey: cryptoService, data.ECDSAKey: cryptoService})

	fakeID := "c62e6d68851cef1f7e55a9d56e3b0c05f3359f16838cad43600f0554e7d3b54d"
	keyInfoURL := fmt.Sprintf("%s/%s", keyInfoBaseURL, fakeID)

	request, err := http.NewRequest("GET", keyInfoURL, nil)
	assert.Nil(t, err)

	res, err := http.DefaultClient.Do(request)
	assert.Nil(t, err)

	assert.Equal(t, 404, res.StatusCode)
}

func TestHSMCreateKeyHandler(t *testing.T) {
	ctx, session := SetupHSMEnv(t)
	defer ctx.Destroy()
	defer ctx.Finalize()
	defer ctx.CloseSession(session)
	defer ctx.Logout(session)

	cryptoService := api.NewRSAHardwareCryptoService(ctx, session)
	setup(signer.CryptoServiceIndex{data.RSAKey: cryptoService})

	createKeyURL := fmt.Sprintf("%s/%s", createKeyBaseURL, data.RSAKey)

	request, err := http.NewRequest("POST", createKeyURL, nil)
	assert.Nil(t, err)

	res, err := http.DefaultClient.Do(request)
	assert.Nil(t, err)

	jsonBlob, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)

	var keyInfo *pb.PublicKey
	err = json.Unmarshal(jsonBlob, &keyInfo)
	assert.Nil(t, err)

	assert.Equal(t, 200, res.StatusCode)
}

func TestSoftwareCreateKeyHandler(t *testing.T) {
	keyStore := trustmanager.NewKeyMemoryStore(passphraseRetriever)
	cryptoService := cryptoservice.NewCryptoService("", keyStore)
	setup(signer.CryptoServiceIndex{data.ED25519Key: cryptoService, data.RSAKey: cryptoService, data.ECDSAKey: cryptoService})

	createKeyURL := fmt.Sprintf("%s/%s", createKeyBaseURL, data.ED25519Key)

	request, err := http.NewRequest("POST", createKeyURL, nil)
	assert.Nil(t, err)

	res, err := http.DefaultClient.Do(request)
	assert.Nil(t, err)

	assert.Equal(t, 200, res.StatusCode)

	jsonBlob, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)

	var keyInfo *pb.PublicKey
	err = json.Unmarshal(jsonBlob, &keyInfo)
	assert.Nil(t, err)
}

func TestHSMSignHandler(t *testing.T) {
	ctx, session := SetupHSMEnv(t)
	defer ctx.Destroy()
	defer ctx.Finalize()
	defer ctx.CloseSession(session)
	defer ctx.Logout(session)

	cryptoService := api.NewRSAHardwareCryptoService(ctx, session)
	setup(signer.CryptoServiceIndex{data.RSAKey: cryptoService})

	tufKey, _ := cryptoService.Create("", data.RSAKey)

	sigRequest := &pb.SignatureRequest{KeyID: &pb.KeyID{ID: tufKey.ID()}, Content: make([]byte, 10)}
	requestJson, _ := json.Marshal(sigRequest)

	reader = strings.NewReader(string(requestJson))

	request, err := http.NewRequest("POST", signBaseURL, reader)

	assert.Nil(t, err)

	res, err := http.DefaultClient.Do(request)
	assert.Nil(t, err)

	jsonBlob, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)

	var sig *pb.Signature
	err = json.Unmarshal(jsonBlob, &sig)
	assert.Nil(t, err)

	assert.Equal(t, tufKey.ID, sig.KeyInfo.KeyID.ID)
	assert.Equal(t, 200, res.StatusCode)
}

func TestSoftwareSignHandler(t *testing.T) {
	keyStore := trustmanager.NewKeyMemoryStore(passphraseRetriever)
	cryptoService := cryptoservice.NewCryptoService("", keyStore)
	setup(signer.CryptoServiceIndex{data.ED25519Key: cryptoService, data.RSAKey: cryptoService, data.ECDSAKey: cryptoService})

	tufKey, err := cryptoService.Create("", data.ED25519Key)
	assert.Nil(t, err)

	sigRequest := &pb.SignatureRequest{KeyID: &pb.KeyID{ID: tufKey.ID()}, Content: make([]byte, 10)}
	requestJson, _ := json.Marshal(sigRequest)

	reader = strings.NewReader(string(requestJson))

	request, err := http.NewRequest("POST", signBaseURL, reader)

	assert.Nil(t, err)

	res, err := http.DefaultClient.Do(request)
	assert.Nil(t, err)

	assert.Equal(t, 200, res.StatusCode)

	jsonBlob, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)

	var sig *pb.Signature
	err = json.Unmarshal(jsonBlob, &sig)
	assert.Nil(t, err)

	assert.Equal(t, tufKey.ID(), sig.KeyInfo.KeyID.ID)
}

func TestSoftwareSignWithInvalidRequestHandler(t *testing.T) {
	keyStore := trustmanager.NewKeyMemoryStore(passphraseRetriever)
	cryptoService := cryptoservice.NewCryptoService("", keyStore)
	setup(signer.CryptoServiceIndex{data.ED25519Key: cryptoService, data.RSAKey: cryptoService, data.ECDSAKey: cryptoService})

	requestJson := "{\"blob\":\"7d16f1d0b95310a7bc557747fc4f20fcd41c1c5095ae42f189df0717e7d7f4a0a2b55debce630f43c4ac099769c612965e3fda3cd4c0078ee6a460f14fa19307\"}"
	reader = strings.NewReader(requestJson)

	request, err := http.NewRequest("POST", signBaseURL, reader)

	assert.Nil(t, err)

	res, err := http.DefaultClient.Do(request)
	assert.Nil(t, err)

	jsonBlob, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)

	var sig *pb.Signature
	err = json.Unmarshal(jsonBlob, &sig)

	assert.Equal(t, 400, res.StatusCode)
}

func TestSignHandlerReturns404WithNonexistentKey(t *testing.T) {
	keyStore := trustmanager.NewKeyMemoryStore(passphraseRetriever)
	cryptoService := cryptoservice.NewCryptoService("", keyStore)
	setup(signer.CryptoServiceIndex{data.ED25519Key: cryptoService, data.RSAKey: cryptoService, data.ECDSAKey: cryptoService})

	fakeID := "c62e6d68851cef1f7e55a9d56e3b0c05f3359f16838cad43600f0554e7d3b54d"

	cryptoService.Create("", data.ED25519Key)

	sigRequest := &pb.SignatureRequest{KeyID: &pb.KeyID{ID: fakeID}, Content: make([]byte, 10)}
	requestJson, _ := json.Marshal(sigRequest)

	reader = strings.NewReader(string(requestJson))

	request, err := http.NewRequest("POST", signBaseURL, reader)
	assert.Nil(t, err)

	res, err := http.DefaultClient.Do(request)
	assert.Nil(t, err)

	assert.Equal(t, 404, res.StatusCode)
}
