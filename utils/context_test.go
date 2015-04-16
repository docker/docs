package utils

import (
	"net/http"
	"testing"
)

func TestContextFactory(t *testing.T) {
	r, err := http.NewRequest("GET", "http://localhost/test/url", nil)
	if err != nil {
		t.Fatalf("Error creating request: %s", err.Error())
	}
	ctx := ContextFactory(r)

	if ctx.Resource() != "/test/url" {
		t.Fatalf("Context has incorrect resource")
	}
}

func TestContext(t *testing.T) {
	ctx := Context{}

	if ctx.Signer() != nil {
		t.Fatalf("Update this test now that Signer has been implemented")
	}
}
