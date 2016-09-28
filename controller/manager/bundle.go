package manager

import (
	"archive/zip"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"strings"
	"text/template"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/orca/auth"
	"github.com/docker/orca/pkg/pki"
	"github.com/docker/orca/utils"
)

var (
	configTemplates = map[string]string{
		"env.sh": `export DOCKER_TLS_VERIFY=1
export DOCKER_CERT_PATH="$(pwd)"
export DOCKER_HOST={{.DockerHost}}
#
# Bundle for user {{.Username}}
# UCP Instance ID {{.UcpID}}
#
{{.AdminBlock}}# Run this command from within this directory to configure your shell:
# eval $(<env.sh)
`,
		"env.ps1": `$Env:DOCKER_TLS_VERIFY = "1"
$Env:DOCKER_CERT_PATH = $(Split-Path $script:MyInvocation.MyCommand.Path)
$Env:DOCKER_HOST = "{{.DockerHost}}"
#
# Bundle for user {{.Username}}
# UCP Instance ID {{.UcpID}}
#
{{.AdminBlock}}# Run this command from within this directory to configure your shell:
# Import-Module .\env.ps1
`,
		"env.cmd": `@echo off
set DOCKER_TLS_VERIFY=1
set DOCKER_CERT_PATH=%~dp0
set DOCKER_HOST={{.DockerHost}}
REM
REM Bundle for user {{.Username}}
REM UCP Instance ID {{.UcpID}}
REM
{{.CMDAdminBlock}}REM Run this command from within this directory to configure your shell:
REM .\env.cmd
`,
	}
)

func generatePrivateKey() (*rsa.PrivateKey, error) {
	// private key
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	return privKey, nil
}

func generateCSR(cn, id string, privateKey *rsa.PrivateKey) (*pki.CertificateSigningRequest, error) {
	tmpl := x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName:         cn,
			Country:            []string{""},
			Province:           []string{""},
			Locality:           []string{""},
			Organization:       []string{"Orca: " + id},
			OrganizationalUnit: []string{"Client"},
		},
		EmailAddresses: []string{""},
	}

	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &tmpl, privateKey)
	if err != nil {
		return nil, err
	}

	block := pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csrBytes,
	}

	csr := pem.EncodeToMemory(&block)

	pkiCSR := &pki.CertificateSigningRequest{
		CertificateRequest: string(csr),
		Profile:            "client",
	}

	return pkiCSR, nil
}

func getPublicKey(certificate string) (string, error) {
	block, _ := pem.Decode([]byte(certificate))
	if block == nil {
		return "", fmt.Errorf("failed to parse certificate")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return "", err
	}

	der, err := x509.MarshalPKIXPublicKey(cert.PublicKey)
	if err != nil {
		return "", err
	}

	var keyBytes bytes.Buffer
	pem.Encode(&keyBytes, &pem.Block{Type: "PUBLIC KEY", Bytes: der})

	key := strings.TrimSpace(keyBytes.String())

	return key, nil
}

func (m DefaultManager) GenerateClientBundle(ctx *auth.Context, host, publicKeyLabel string) ([]byte, error) {
	privateKey, err := generatePrivateKey()
	if err != nil {
		return nil, err
	}

	privateKeyBytes := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})

	privateKeyEncoded := string(privateKeyBytes)

	account := ctx.User
	username := account.Username

	csr, err := generateCSR(username, m.ID(), privateKey)
	if err != nil {
		return nil, err
	}

	log.Debugf("signing client certificate request: cn=%s", username)

	// Figure out if this is an admin or regular user
	isAdmin := account.Admin

	clientCertificate := ""
	chain := ""
	if isAdmin {
		// Use swarm for admins so they can access swarm if they need
		if certResponse, err := m.ClusterSignCSR(csr); err != nil {
			log.Debugf("error requesting to sign csr: %s", err)
			return nil, err
		} else {
			clientCertificate = utils.JoinCerts(certResponse.Certificate, m.clusterCAChain)
			chain = certResponse.CertificateChain
		}
	} else {
		// Non-admins get orca signed certs
		// TODO - Check for external CA config, and give a better error
		if certResponse, err := m.ClientSignCSR(csr); err != nil {
			log.Debugf("error requesting to sign csr: %s", err)
			return nil, err
		} else {
			clientCertificate = utils.JoinCerts(certResponse.Certificate, m.clientCAChain)
			chain = certResponse.CertificateChain
		}
	}

	pubKey, err := getPublicKey(clientCertificate)
	if err != nil {
		return nil, err
	}

	adminBlock := `# This admin cert will also work directly against Swarm and the individual
# engine proxies for troubleshooting.  After sourcing this env file, use
# "docker info" to discover the location of Swarm managers and engines.
# and use the --host option to override $DOCKER_HOST
#
`
	cmdAdminBlock := `REM This admin cert will also work directly against Swarm and the individual
REM engine proxies for troubleshooting.  After sourcing this env file, use
REM "docker info" to discover the location of Swarm managers and engines.
REM and use the --host option to override $DOCKER_HOST
REM
`
	if !isAdmin {
		adminBlock = ""
		cmdAdminBlock = ""
	}

	config := struct {
		DockerHost    string
		UcpID         string
		Username      string
		AdminBlock    string
		CMDAdminBlock string
	}{
		DockerHost:    fmt.Sprintf("tcp://%s", host),
		UcpID:         m.ID(),
		Username:      username,
		AdminBlock:    adminBlock,
		CMDAdminBlock: cmdAdminBlock,
	}

	// generate bundle
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	var files = []struct {
		Name, Data string
	}{
		{"ca.pem", utils.JoinCerts(chain, m.clusterCAChain, m.clientCAChain)},
		{"cert.pem", clientCertificate},
		{"key.pem", privateKeyEncoded},
		{"cert.pub", pubKey},
	}
	for filename, tmpl := range configTemplates {
		t, err := template.New("config").Parse(tmpl)
		if err != nil {
			return nil, err
		}

		var tBuf bytes.Buffer
		if err := t.Execute(&tBuf, config); err != nil {
			return nil, err
		}

		envScript := tBuf.String()
		files = append(files, struct {
			Name, Data string
		}{filename, envScript})
	}

	for _, file := range files {
		h := &zip.FileHeader{Name: file.Name}
		h.SetModTime(time.Now())
		if file.Name == "key.pem" {
			h.SetMode(0600)
		} else {
			h.SetMode(0644)
		}

		f, err := w.CreateHeader(h)
		if err != nil {
			return nil, err
		}
		_, err = f.Write([]byte(file.Data))
		if err != nil {
			return nil, err
		}
	}

	if err := w.Close(); err != nil {
		return nil, err
	}

	// update user account with pub key
	if err := m.GetAuthenticator().AddUserPublicKey(account, publicKeyLabel, privateKey.Public()); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Add a CA (and/or intermediate cert) to our list of trusted user certs
func (m *DefaultManager) AddTrustedCert(caCert string) {
	m.clientCAChain = utils.JoinCerts(m.clientCAChain, caCert)
}
