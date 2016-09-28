package pki

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"testing"
)

func mockSignHandler(w http.ResponseWriter, r *http.Request) {

	response := `{
    "success": true,
    "result": {
        "certificate": "-----BEGIN CERTIFICATE-----\n12345\n-----END CERTIFICATE-----\n"
    },
    "errors": [],
    "messages": []
}
    `
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func TestApiSignCSR(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(mockSignHandler))
	defer ts.Close()

	tlsConfig := &tls.Config{}

	c, err := NewDefaultClient(ts.URL, tlsConfig)
	if err != nil {
		t.Fatal(err)
	}

	csr := &CertificateSigningRequest{
		CertificateRequest: "-----BEGIN CERTIFICATE REQUEST-----\n12345\n-----END CERTIFICATE REQUEST-----\n",
	}

	cert, err := c.SignCSR(csr)
	if err != nil {
		t.Fatal(err)
	}

	if cert == nil {
		t.Fatalf("expected certificate; received nil")
	}

	expectedCert := "-----BEGIN CERTIFICATE-----\n12345\n-----END CERTIFICATE-----\n"
	if cert.Certificate != expectedCert {
		t.Fatal("certificate received does not match expected")
	}
}
