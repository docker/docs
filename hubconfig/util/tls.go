package util

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/dhe-deploy"
	"github.com/docker/dhe-deploy/hubconfig"
)

func GenTLSCert(domainName string) (*tls.Certificate, error) {
	now := time.Now()
	if domainName == "" {
		domainName = deploy.UnconfiguredCertSentinelCN
	}
	// if we don't put the IP in the SAN, browsers will not consider the cert valid
	ipAddresses := []net.IP{}
	if ip := net.ParseIP(domainName); ip != nil {
		ipAddresses = append(ipAddresses, ip)
	}
	// generate random serial number
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, err
	}
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}
	// RFC-compliant mandatory subjectKeyId for CAs
	derBytes, err := x509.MarshalPKIXPublicKey(key.Public())
	if err != nil {
		return nil, err
	}
	// subjectKeyId is sha1 of public key bytes
	hash := sha1.New()
	_, err = hash.Write(derBytes)
	if err != nil {
		return nil, err
	}
	subjectKeyId := hash.Sum(nil)
	certTemplate := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Country:            []string{"US"},
			Organization:       []string{"Docker"},
			OrganizationalUnit: []string{"Docker"},
			Locality:           []string{"San Francisco"},
			CommonName:         domainName,
		},
		NotBefore:             now,
		NotAfter:              now.Add(365 * 24 * time.Hour),
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{},
		SubjectKeyId:          subjectKeyId,
		BasicConstraintsValid: true,
		IPAddresses:           ipAddresses,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, &certTemplate, &certTemplate, &key.PublicKey, key)
	if err != nil {
		return nil, fmt.Errorf("Failed to generate TLS certificate: %s", err)
	}

	keyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}
	certPEM := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	}
	cert, err := tls.X509KeyPair(pem.EncodeToMemory(certPEM), pem.EncodeToMemory(keyPEM))
	if err != nil {
		return nil, fmt.Errorf("Failed to generate TLS key/certificate pair: %s", err)
	}
	log.WithField("domain", domainName).Info("Generated TLS certificate.")
	return &cert, nil
}

func CertificateToPEM(certificate *tls.Certificate) ([]byte, error) {
	if certificate == nil {
		return []byte{}, nil
	}
	if certificate.Certificate == nil {
		return []byte{}, nil
	}
	tlsBuffer := new(bytes.Buffer)
	for _, cert := range certificate.Certificate {
		if err := pem.Encode(tlsBuffer, &pem.Block{Type: "CERTIFICATE", Bytes: cert}); err != nil {
			return nil, err
		}
	}
	switch k := certificate.PrivateKey.(type) {
	case *rsa.PrivateKey:
		if err := pem.Encode(tlsBuffer, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}); err != nil {
			return nil, err
		}
	case *ecdsa.PrivateKey:
		if b, err := x509.MarshalECPrivateKey(k); err != nil {
			return nil, fmt.Errorf("Unable to marshal ECDSA private key: %v", err)
		} else if err := pem.Encode(tlsBuffer, &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("Unsupported TLS private key type: %T", k)
	}
	return tlsBuffer.Bytes(), nil
}

func CertificateToCertPEM(certificate *tls.Certificate) ([]byte, error) {
	if certificate == nil {
		return []byte{}, nil
	}
	if certificate.Certificate == nil {
		return []byte{}, nil
	}
	tlsBuffer := new(bytes.Buffer)
	for _, cert := range certificate.Certificate {
		if err := pem.Encode(tlsBuffer, &pem.Block{Type: "CERTIFICATE", Bytes: cert}); err != nil {
			return nil, err
		}
	}
	return tlsBuffer.Bytes(), nil
}

func CertificateToKeyPEM(certificate *tls.Certificate) ([]byte, error) {
	if certificate == nil {
		return []byte{}, nil
	}
	tlsBuffer := new(bytes.Buffer)
	switch k := certificate.PrivateKey.(type) {
	case *rsa.PrivateKey:
		if err := pem.Encode(tlsBuffer, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}); err != nil {
			return nil, err
		}
	case *ecdsa.PrivateKey:
		if b, err := x509.MarshalECPrivateKey(k); err != nil {
			return nil, fmt.Errorf("Unable to marshal ECDSA private key: %v", err)
		} else if err := pem.Encode(tlsBuffer, &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("Unsupported TLS private key type: %T", k)
	}
	return tlsBuffer.Bytes(), nil
}

func HubConfigTLSDomainConsistent(userHubConfig *hubconfig.UserHubConfig) error {
	if userHubConfig == nil {
		// this shouldn't happen
		return fmt.Errorf("no userHubConfig provided for validation")
	}

	// check if the cert parses and the private key matches
	cert, err := tls.X509KeyPair([]byte(userHubConfig.WebTLSCert), []byte(userHubConfig.WebTLSKey))
	if err != nil {
		return fmt.Errorf("Failed to parse TLS key/certificate pair: %s", err)
	}

	// check if the cert matches the domain
	domain := strings.Split(userHubConfig.DTRHost, ":")[0]
	if !TLSCertMatchesDomain(&cert, domain) {
		return fmt.Errorf("Domain in certificate does not match configured domain: %s", domain)
	}

	// check if the CA is valid for the cert / chain / domain name
	certPool := x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM([]byte(userHubConfig.WebTLSCert))
	if !ok {
		// this is technically impossible because we checked above that the cert is valid and has a valid private key
		return fmt.Errorf("Failed to parse cert")
	}
	caCertPool := x509.NewCertPool()
	ok = caCertPool.AppendCertsFromPEM([]byte(userHubConfig.WebTLSCA))
	if !ok {
		return fmt.Errorf("Failed to parse CA")
	}
	opts := x509.VerifyOptions{
		DNSName:       domain,
		Intermediates: certPool,
		Roots:         caCertPool,
	}
	// this is guaranteed to have a certificate because tls.X509KeyPair above didn't fail
	x509Cert, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return err
	}
	_, err = x509Cert.Verify(opts)
	if err != nil {
		return err
	}

	return nil
}

// if no domain name is set yet, assume that it matches the cert, then if later it stops matching, we'll regenerate it; if no domain name and no cert are set, we generate a cert with a dummy name
func TLSCertMatchesDomain(tlsCertificate *tls.Certificate, userHubConfigDomainName string) bool {
	if tlsCertificate == nil {
		log.WithField("domainName", userHubConfigDomainName).Info("No TLS certificate found.")
		return false
	} else if userHubConfigDomainName != "" {
		x509cert, err := x509.ParseCertificate(tlsCertificate.Certificate[0])
		if err != nil {
			log.WithField("error", err).Error("Unable to parse TLS certificate as x509 certificate.")
			return false
		} else if x509cert.VerifyHostname(userHubConfigDomainName) != nil {
			log.WithField("domainName", userHubConfigDomainName).Info("TLS certificate does not match domain name")
			return false
		}
	}
	return true
}

func SetTLSCertificateInHubConfig(userHubConfig *hubconfig.UserHubConfig, cert *tls.Certificate, ca *tls.Certificate) error {
	certPEM, err := CertificateToCertPEM(cert)
	if err != nil {
		return err
	}
	keyPEM, err := CertificateToKeyPEM(cert)
	if err != nil {
		return err
	}
	caPEM, err := CertificateToCertPEM(ca)
	if err != nil {
		return err
	}
	userHubConfig.WebTLSCert = string(certPEM)
	userHubConfig.WebTLSKey = string(keyPEM)
	userHubConfig.WebTLSCA = string(caPEM)
	return nil
}

func DefaultUserHubConfig(host string) (*hubconfig.UserHubConfig, error) {
	cert, err := GenTLSCert(strings.Split(host, ":")[0])
	if err != nil {
		return nil, fmt.Errorf("failed to generate default UserHubConfig: %s", err)
	}

	var userHubConfig hubconfig.UserHubConfig
	userHubConfig.DTRHost = host

	err = SetTLSCertificateInHubConfig(&userHubConfig, cert, cert)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize default UserHubConfig: %s", err)
	}

	// turn on analytics reporting by default
	userHubConfig.ReportAnalytics = true
	userHubConfig.AnonymizeAnalytics = false
	return &userHubConfig, nil
}
