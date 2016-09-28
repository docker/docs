package manager

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"testing"
)

const (
	testCert = `-----BEGIN CERTIFICATE-----
MIIC6TCCAdOgAwIBAgIQIPyjWjpPHLHum3S6ztlT2TALBgkqhkiG9w0BAQswEjEQ
MA4GA1UEChMHdW5rbm93bjAeFw0xNTA4MTkxNzI4MDBaFw0xODA4MDMxNzI4MDBa
MBIxEDAOBgNVBAoTB3Vua25vd24wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEK
AoIBAQC4lYYA6cvD/5kc/YSFcHatdVePPVh0+tbiO9Gbvy81NFcMzIHhMlc5HPDb
Ln1AkcTVjYMdKZo0gMHONS4qPuCY4V4i0h17OkkzAKDEHtiI1v8mcTzPiFm95X1Y
G2eIoN4gnji90KoHd55LOM/4Xs3IrltCwIWeJD/EeVWl9sSbi2o67q80GgtHWMB1
NDebw1G/lFR2XJywfu1sjXLyfT2dq2+v9nsMZfjxpU7EEBacOhG0MMszvNxQmC2a
STZn79jE+8PtS+GXiN4W/GN4lVVXkJfFVj2hFlOyxHovPuA64+XH56lsmd0IslUm
4UfYwG4G+TrZ7+9ynP1IK2mJQ7dPAgMBAAGjPzA9MA4GA1UdDwEB/wQEAwIAqDAd
BgNVHSUEFjAUBggrBgEFBQcDAgYIKwYBBQUHAwEwDAYDVR0TAQH/BAIwADALBgkq
hkiG9w0BAQsDggEBAJsRshSa/IUmJsFgSyoBypCV7VGR/lKonoLtId4OZEolVrBM
nPbOwoTKbF9HFbpC5hWEiOmr/OvldqjgEkwVuG2qc+xgX5uLFUrvdZpcb7PC7iDP
/qIRSR3nG2sFdByDtEL8clxRLhCheNBbZFZRoqLiR8rEMdlG137Cis8dUJhF9Jzw
SnEHJUTgAwf3W5KTFTAND/BfBcS8Bdy+n1dJDXQQXQy2YAqwDjyBJl6t/Gf1G9Qp
z9VtRW5xddjjK5G9fD53Ym6qKq6/wDuxGxq6sciq+XLCb9i3ROvfLFIU04ZDe2Ww
QaZOMRvewW3l/fOR1HJTkeBT3HpAdWV12kyxDSY=
-----END CERTIFICATE-----
`
)

func TestGetPublicKey(t *testing.T) {
	expectedPubKey := `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAuJWGAOnLw/+ZHP2EhXB2
rXVXjz1YdPrW4jvRm78vNTRXDMyB4TJXORzw2y59QJHE1Y2DHSmaNIDBzjUuKj7g
mOFeItIdezpJMwCgxB7YiNb/JnE8z4hZveV9WBtniKDeIJ44vdCqB3eeSzjP+F7N
yK5bQsCFniQ/xHlVpfbEm4tqOu6vNBoLR1jAdTQ3m8NRv5RUdlycsH7tbI1y8n09
natvr/Z7DGX48aVOxBAWnDoRtDDLM7zcUJgtmkk2Z+/YxPvD7Uvhl4jeFvxjeJVV
V5CXxVY9oRZTssR6Lz7gOuPlx+epbJndCLJVJuFH2MBuBvk62e/vcpz9SCtpiUO3
TwIDAQAB
-----END PUBLIC KEY-----`
	pubKey, err := getPublicKey(testCert)
	if err != nil {
		t.Fatalf("error getting public key: %s", err)
	}

	if pubKey != expectedPubKey {
		t.Fatalf("public key does not match:\nExpected:\n%s\n\nReceived:\n%s\n", expectedPubKey, pubKey)
	}
}

func TestGenerateCSR(t *testing.T) {
	commonName := "test"
	pk, err := generatePrivateKey()
	if err != nil {
		t.Fatal(err)
	}

	csr, err := generateCSR(commonName, "ID", pk)
	if err != nil {
		t.Fatal(err)
	}

	d := []byte(csr.CertificateRequest)
	data := bytes.TrimSpace(d)
	block, _ := pem.Decode(data)
	if block == nil {
		t.Fatalf("no PEM block found")
	}

	r, err := x509.ParseCertificateRequest(block.Bytes)
	if err != nil {
		t.Fatal(err)
	}

	if r.Subject.CommonName != commonName {
		t.Fatalf("common name does not match: expected %s; received %s", commonName, r.Subject.CommonName)
	}
}
