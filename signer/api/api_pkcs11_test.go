// +build pkcs11

package api_test

import (
	"os"
	"testing"

	"github.com/miekg/pkcs11"
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
