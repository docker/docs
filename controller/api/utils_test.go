package api

import (
	"testing"
)

func TestParamBoolValueTrue(t *testing.T) {
	v := "1"
	if !paramBoolValue(v) {
		t.Errorf("expected true for %s", v)
	}

	v = "True"
	if !paramBoolValue(v) {
		t.Errorf("expected true for %s", v)
	}

	v = "true"
	if !paramBoolValue(v) {
		t.Errorf("expected true for %s", v)
	}

	v = "true   "
	if !paramBoolValue(v) {
		t.Errorf("expected true for %s", v)
	}
}

func TestParamBoolValueFalse(t *testing.T) {
	v := "0"
	if paramBoolValue(v) {
		t.Errorf("expected false for %s", v)
	}

	v = "False"
	if paramBoolValue(v) {
		t.Errorf("expected false for %s", v)
	}

	v = "false"
	if paramBoolValue(v) {
		t.Errorf("expected false for %s", v)
	}

	v = "false   "
	if paramBoolValue(v) {
		t.Errorf("expected false for %s", v)
	}
}
