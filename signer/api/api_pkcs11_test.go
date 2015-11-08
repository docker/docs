// +build pkcs11

package api_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/docker/notary/signer"
	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/tuf/signed"
	"github.com/miekg/pkcs11"
	"github.com/stretchr/testify/assert"

	pb "github.com/docker/notary/proto"
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

func TestHSMCreateKeyHandler(t *testing.T) {
	ctx, session := SetupHSMEnv(t)
	defer ctx.Destroy()
	defer ctx.Finalize()
	defer ctx.CloseSession(session)
	defer ctx.Logout(session)

	cryptoService := signed.NewEd25519()
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

func TestHSMSignHandler(t *testing.T) {
	ctx, session := SetupHSMEnv(t)
	defer ctx.Destroy()
	defer ctx.Finalize()
	defer ctx.CloseSession(session)
	defer ctx.Logout(session)

	cryptoService := signed.NewEd25519()
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
