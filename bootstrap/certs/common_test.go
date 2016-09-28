package certs

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	cfssllog "github.com/cloudflare/cfssl/log"
)

func TestGenerateCSR(t *testing.T) {
	cfssllog.Level = cfssllog.LevelWarning
	cn, ou := "myCN", "myOU"
	hosts := []string{"hostname1", "hostname2"}
	csrExpected := "-----BEGIN CERTIFICATE REQUEST-----"
	keyExpected := "-----BEGIN RSA PRIVATE KEY-----"

	csr, key, err := GenerateCSR(cn, ou, hosts)
	if err != nil {
		t.Errorf("Unexpected CSR failure %s", err)
	}
	if !strings.Contains(string(csr), csrExpected) {
		t.Errorf("Was expecting to contain %s, but found %s", csrExpected, csr)
	}
	if !strings.Contains(string(key), keyExpected) {
		t.Errorf("Was expecting to contain %s, but found %s", keyExpected, key)
	}
}

func TestInitLocalNodePos(t *testing.T) {
	cfssllog.Level = cfssllog.LevelWarning

	caDir := fmt.Sprintf("/tmp/%d_test_dir_ca", os.Getpid())
	certDir := fmt.Sprintf("/tmp/%d_test_dir_cert", os.Getpid())

	if err := os.Mkdir(certDir, 0700); err != nil {
		t.Errorf("Failed to setup tmp dir %s %s", certDir, err)
	}
	defer os.RemoveAll(certDir)
	if err := os.Mkdir(caDir, 0700); err != nil {
		t.Errorf("Failed to setup tmp dir %s %s", caDir, err)
	}
	defer os.RemoveAll(caDir)

	certPrefix := "certprefix"
	nodeName := "nodename"
	hostnames := []string{"hostname1", "hostname2"}
	if err := writeCAFiles(caDir); err != nil {
		t.Errorf("Unexpected failure %s", err)
	}

	err := InitLocalNode(caDir, certDir, certPrefix, "", nodeName, hostnames, os.Geteuid(), os.Getgid())
	if err != nil {
		t.Errorf("Unexpected failure %s", err)
	}
	// TODO - verify the files look plausible
}

// Write out canned versions of some files the signer needs
func writeCAFiles(dirname string) error {
	caCert := `-----BEGIN CERTIFICATE-----
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
	caKey := `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----
`

	if err := ioutil.WriteFile(filepath.Join(dirname, "cert.pem"), []byte(caCert), 0600); err != nil {
		return err
	}
	if err := ioutil.WriteFile(filepath.Join(dirname, "key.pem"), []byte(caKey), 0600); err != nil {
		return err
	}
	return nil
}
