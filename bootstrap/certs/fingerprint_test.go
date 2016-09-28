package certs

import (
	"fmt"
	"github.com/docker/orca/bootstrap/config"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetFingerprint(t *testing.T) {
	config.CertDir = fmt.Sprintf("/tmp/%d_test_dir", os.Getpid())

	if err := os.Mkdir(config.CertDir, 0700); err != nil {
		t.Errorf("Failed to setup tmp dir %s %s", config.CertDir, err)
	}
	defer os.RemoveAll(config.CertDir)

	// Predefined/hardcoded values
	data := `-----BEGIN CERTIFICATE-----
MIIFEjCCAvygAwIBAgIIOaFXTz/3NbswCwYJKoZIhvcNAQENMBcxFTATBgNVBAMT
DE9yY2EgUm9vdCBDQTAeFw0xNTA4MTIyMDQyMDBaFw0yMDA4MTAyMDQyMDBaMBcx
FTATBgNVBAMTDE9yY2EgUm9vdCBDQTCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCC
AgoCggIBAMMy+ooz1u+ECAvm63ZkVWZJePfSUjMGLEsRZZAKRhJE8Xe4aos2iu6J
zbU57rvbVs274a4KdYq1CG+EslD8JUPILmHE9EmdD3I1X7ccQ8kIGm5yRNWz79A+
KUWC+/OC8CRYyXzw2L2+/zsCSNPAR2CB+G/xEI+cv8ePE6UFGzPQvs3+EO0Xn/U5
G4+vcLC9JTaILqswZxtv9eAgRUT842AUtRKWD63p5DSh9d75sv4blux5/hrZsDp5
50aDCJTDksZdxCmxJJTpRSDAQs4ueIrx7sIOYbMt6pVHAx8VBpq53DZbrLV69UHz
6gfHgbrvgg8FJVGuiLbRapnM6YX5E4jZRrUm1OF5Gp90vtcoiXsf0xaw7I1euZxT
JlDS1Fs4MmLlGDjn3SI8za+Rl3/UrSi0m2hCvaCwzATQKXGyEpozRGLpazo3BfaV
sXwvN4o7IezQAvYJMlCf6lOdzkrX8B5ly+OhRBhbZJgzp9Efkf00TwxR2WHh/RNB
jNt/LRy4oVxF89/sHf0NsnM0WwQUfL18sDZ/uhUfY2mmtR6Dxm9nCTJPLiNCZ5om
nhlgXMgQigOZpzjtbxUQMdiQTyBjz5NizpSesxvHX+VhOUw+0ehiTv5LzReOci1g
KXNoHCZBjb7P+Juqrvv3vs2uIbwwppYNDZB0touDmalNfPV7dX0/AgMBAAGjZjBk
MA4GA1UdDwEB/wQEAwIABjASBgNVHRMBAf8ECDAGAQH/AgECMB0GA1UdDgQWBBQS
6sp1GI6mrrK7Ez4D5zgCBKtlCTAfBgNVHSMEGDAWgBQS6sp1GI6mrrK7Ez4D5zgC
BKtlCTALBgkqhkiG9w0BAQ0DggIBAGp72If6sJfG0RduGc2LE7FYhXW9u5Y7KT3s
fZtsybv1PxL34vwmUUdwi1q2rN+2qmY4oc68fwWvW4ZSDSmG7Ig5VG0b7/x9XQwd
/mEcRS6YWN71QLbnLj3aF7axixsnnzfnIyg57H52kL4eeS00MBvUkti4YIxt66LD
BnJ7LkrsL95Bm9eYRPOzF8hMDi5jU/JzK0LEBG4pZY4vTkfPkpLmgmOZhPmeaOSO
SlCs0VqDH48PJZcWtsEus72xFfeoP41tz2jZda9olBEPwTyk95nCjnQWwM9Ao3Yp
0LOrhHYr+uPY4cpmprTiqL4GfuP/3fogP1Lak4plnlLy+FHmXr58Ou01bPZ0AO1P
WAHcx/zoCabD0gUW7YbogM89Q6ZCwQcU9SGWuhroZncQYoQ41Q3ZVjH1UooRov00
KK/Cp/ve4OsVBTHhb3uZqp775B7CFnI7zz9b9LU2TCVVPuyocShiEUwCVYTdEc8T
yKSWRnwlXOKnhhVfn1ges8sjwb794yd15DuHa0j9asWhVxOP+8Q+v3iiUUCIXKxr
RA5z60YIrYQ/QPSGmqDLoMw31X2VXeL3Nl4kebVKaKzwF5cmKNIqPBi6oNn3lui5
cNh+xk+gyKxBoirws/9HeYJr0E+SBgVX2RepniOveCFZjLhuwM3vKEAtkFTSNCRx
ZoP56+dK
-----END CERTIFICATE-----
`
	expected := "5C:9D:0C:44:7C:FE:F0:D2:97:80:97:BA:9B:B6:F8:37:53:53:82:D1:33:29:17:F4:47:BD:F3:EA:2D:93:0B:5B"
	filename := filepath.Join(config.CertDir, "fingerprint.dat")

	if err := ioutil.WriteFile(filename, []byte(data), 0600); err != nil {
		t.Errorf("Unexpected error %s", err)
	}

	output, err := GetFingerprint(filename)
	if err != nil {
		t.Errorf("Unexpected error %s", err)
	}

	if !strings.Contains(output, expected) {
		t.Errorf("Expected to contain %s, got %s", expected, output)
	}
}

func TestGetFingerprintBadFile(t *testing.T) {
	_, err := GetFingerprint("/bogus/filename/that/doesnt/exist")
	if err == nil {
		t.Error("Expected error, but it didn't fail")
	}
}

func TestGetFingerprintMalformed(t *testing.T) {
	config.CertDir = fmt.Sprintf("/tmp/%d_test_dir", os.Getpid())

	if err := os.Mkdir(config.CertDir, 0700); err != nil {
		t.Errorf("Failed to setup tmp dir %s %s", config.CertDir, err)
	}
	defer os.RemoveAll(config.CertDir)

	// Predefined/hardcoded values
	data := `Not a cert`
	filename := filepath.Join(config.CertDir, "fingerprint.dat")

	if err := ioutil.WriteFile(filename, []byte(data), 0600); err != nil {
		t.Errorf("Unexpected error %s", err)
	}

	_, err := GetFingerprint(filename)
	if err != ErrMalformedCert {
		t.Errorf("expected malformed cert error, got %s", err)
	}
}
