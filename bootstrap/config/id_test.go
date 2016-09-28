package config

import (
	"strings"
	"testing"
)

func TestGetNewID(t *testing.T) {
	OrcaInstanceID = ""
	OrcaInstanceKey = ""
	expected := "-----BEGIN"

	if err := GetNewID(); err != nil {
		t.Errorf("Failed to get new ID %s", err)
	}
	// Now make sure things are gone
	if OrcaInstanceID == "" {
		t.Error("Instance ID failed to generate")
	}
	if OrcaInstanceKey == "" {
		t.Error("Key failed to generate")
	}
	if !strings.Contains(OrcaInstanceKey, expected) {
		t.Errorf("Key should contain %s, was %s", expected, OrcaInstanceKey)
	}
}
