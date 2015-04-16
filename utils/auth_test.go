package utils

import (
	"testing"
)

func TestInsecureAuthorization(t *testing.T) {
	auther := InsecureAuthorizer{}
	ctx := Context{}
	err := auther.Authorize(&ctx, SSNoAuth)
	if err != nil {
		t.Fatalf("Failed to authorize with InsecureAuthorizer")
	}
	if !ctx.Authorization().HasScope(SSCreate) {
		t.Fatalf("InsecureAuthorization failed to approve a scope")
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
