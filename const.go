package notary

// application wide constants
const (
	// MaxDownloadSize is the maximum size we'll download for metadata if no limit is given
	MaxDownloadSize int64 = 100 << 20
	// MaxTimestampSize is the maximum size of timestamp metadata - 1MiB.
	MaxTimestampSize int64 = 1 << 20
	// MinRSABitSize is the minimum bit size for RSA keys allowed in notary
	MinRSABitSize = 2048
	// MinThreshold requires a minimum of one threshold for roles; currently we do not support a higher threshold
	MinThreshold = 1
	// PrivKeyPerms are the file permissions to use when writing private keys to disk
	PrivKeyPerms = 0700
	// PubCertPerms are the file permissions to use when writing public certificates to disk
	PubCertPerms = 0755
	// Sha256HexSize is how big a Sha256 hex is in number of characters
	Sha256HexSize = 64
	// TrustedCertsDir is the directory, under the notary repo base directory, where trusted certs are stored
	TrustedCertsDir = "trusted_certificates"
)
