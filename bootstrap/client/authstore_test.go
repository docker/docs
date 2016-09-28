package client

import (
	"testing"
)

func TestSanitizeHostnameRegular(t *testing.T) {
	h := "foo.bar.com"
	expected := "foo_bar_com"

	res := sanitizeHostname(h)

	if res != expected {
		t.Fatalf("expected %s; received %s", expected, res)
	}
}

func TestSanitizeHostnameDashes(t *testing.T) {
	h := "foo.bar.internal-ec2.com"
	expected := "foo_bar_internal_ec2_com"

	res := sanitizeHostname(h)

	if res != expected {
		t.Fatalf("expected %s; received %s", expected, res)
	}
}
