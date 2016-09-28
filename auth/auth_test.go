package auth

import (
	"testing"
)

const (
	testPass  = "FOOPASS.+&^"
	testUser  = "admin"
	testToken = "12345"
)

func TestHash(t *testing.T) {
	h, err := Hash(testPass)
	if err != nil {
		t.Error(err)
	}

	if len(h) == 0 {
		t.Errorf("expected a hashed password; go a zero length string")
	}
}
