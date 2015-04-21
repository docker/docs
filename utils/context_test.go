package utils

import (
	"net/http"
	"testing"

	"github.com/endophage/go-tuf/signed"
)

func TestNewContex(t *testing.T) {
	r, err := http.NewRequest("GET", "http://localhost/test/url", nil)
	if err != nil {
		t.Fatalf("Error creating request: %s", err.Error())
	}
	ctx := NewContext(r, &signed.Ed25519{})

	if ctx.Resource() != "/test/url" {
		t.Fatalf("Context has incorrect resource")
	}
}

func TestContextTrust(t *testing.T) {
	ctx := context{}

	if ctx.Trust() != nil {
		t.Fatalf("Update this test now that Trust has been implemented")
	}
}
