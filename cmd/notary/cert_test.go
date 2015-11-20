package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/docker/notary/trustmanager"
	"github.com/stretchr/testify/assert"
)

func generateCertificate(t *testing.T, gun string, expireInHours int64) *x509.Certificate {
	template, err := trustmanager.NewCertificate(gun)
	assert.NoError(t, err)
	template.NotAfter = template.NotBefore.Add(
		time.Hour * time.Duration(expireInHours))

	ecdsaPrivKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.NoError(t, err)

	certBytes, err := x509.CreateCertificate(rand.Reader, template, template,
		ecdsaPrivKey.Public(), ecdsaPrivKey)
	assert.NoError(t, err)

	cert, err := x509.ParseCertificate(certBytes)
	assert.NoError(t, err)
	return cert
}

// If there are no certs in the cert store store, a message that there are no
// certs should be displayed.
func TestPrettyPrintZeroCerts(t *testing.T) {
	var b bytes.Buffer
	prettyPrintCerts([]*x509.Certificate{}, &b)
	text, err := ioutil.ReadAll(&b)
	assert.NoError(t, err)

	lines := strings.Split(strings.TrimSpace(string(text)), "\n")
	assert.Len(t, lines, 1)
	assert.Equal(t, "No trusted root certificates present.", lines[0])
}

// Certificates are pretty-printed in table form sorted by gun and then expiry
func TestPrettySortedCerts(t *testing.T) {
	unsorted := []*x509.Certificate{
		generateCertificate(t, "xylitol", 77),    // 3 days 5 hours
		generateCertificate(t, "xylitol", 12),    // less than 1 day
		generateCertificate(t, "cheesecake", 25), // a little more than 1 day
		generateCertificate(t, "baklava", 239),   // almost 10 days
	}

	var b bytes.Buffer
	prettyPrintCerts(unsorted, &b)
	text, err := ioutil.ReadAll(&b)
	assert.NoError(t, err)

	expected := [][]string{
		{"baklava", "9 days"},
		{"cheesecake", "1 day"},
		{"xylitol", "< 1 day"},
		{"xylitol", "3 days"},
	}

	lines := strings.Split(strings.TrimSpace(string(text)), "\n")
	assert.Len(t, lines, len(expected)+2)

	// starts with headers
	assert.True(t, reflect.DeepEqual(strings.Fields(lines[0]), strings.Fields(
		"GUN     FINGERPRINT OF TRUSTED ROOT CERTIFICATE      EXPIRES IN")))
	assert.Equal(t, "----", lines[1][:4])

	for i, line := range lines[2:] {
		splitted := strings.Fields(line)
		assert.True(t, len(splitted) >= 3)
		assert.Equal(t, expected[i][0], splitted[0])
		assert.Equal(t, expected[i][1], strings.Join(splitted[2:], " "))
	}
}
