package util

import (
	"crypto/x509"
	"io/ioutil"
)

// this file is pretty much all copied from go's crypto/x509 internal code. :(

// Possible certificate files; stop after finding one.
var certFiles = []string{
	"/etc/ssl/certs/ca-certificates.crt", // Debian/Ubuntu/Gentoo etc.
	"/etc/pki/tls/certs/ca-bundle.crt",   // Fedora/RHEL
	"/etc/ssl/ca-bundle.pem",             // OpenSUSE
	"/etc/pki/tls/cacert.pem",            // OpenELEC
}

// Possible directories with certificate files; stop after successfully
// reading at least one file from a directory.
var certDirectories = []string{
	"/etc/ssl/certs",               // SLES10/SLES11, https://golang.org/issue/12139
	"/system/etc/security/cacerts", // Android
}

func systemRootsPool() *x509.CertPool {
	roots := x509.NewCertPool()
	for _, file := range certFiles {
		data, err := ioutil.ReadFile(file)
		if err == nil {
			roots.AppendCertsFromPEM(data)
			return roots
		}
	}

	for _, directory := range certDirectories {
		fis, err := ioutil.ReadDir(directory)
		if err != nil {
			continue
		}
		rootsAdded := false
		for _, fi := range fis {
			data, err := ioutil.ReadFile(directory + "/" + fi.Name())
			if err == nil && roots.AppendCertsFromPEM(data) {
				rootsAdded = true
			}
		}
		if rootsAdded {
			return roots
		}
	}

	panic("unable to find certs")

	// All of the files failed to load. systemRoots will be nil which will
	// trigger a specific error at verification time.
	return roots
}
