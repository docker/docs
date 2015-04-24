package auth

import (
	"net/http"
	"testing"
)

func TestInsecureAuthorization(t *testing.T) {
	auther := InsecureAuthorizer{}
	_, err := auther.Authorize(&http.Request{}, SSNoAuth)
	if err != nil {
		t.Fatalf("Failed to authorize with InsecureAuthorizer")
	}

}

func TestNoAuthorization(t *testing.T) {
	auth := NoAuthorization{}
	if auth.HasScope(SSCreate) {
		t.Fatalf("NoAuthorization should not have any scopes")
	}
}

func TestSimpleScope(t *testing.T) {
	scope1 := SimpleScope("Test")
	scope2 := SimpleScope("Test")
	if !scope1.Compare(scope2) {
		t.Fatalf("Expected scope1 and scope2 to match")
	}

	scope3 := SimpleScope("Test")
	scope4 := SimpleScope("Don't Match")
	if scope3.Compare(scope4) {
		t.Fatalf("Expected scope3 and scope4 not to match")
	}
}
