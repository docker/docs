package notary

// application wide constants
const (
	// Require a minimum of one threshold for roles, currently we do not support a higher threshold
	MinThreshold    = 1
	PrivKeyPerms    = 0700
	PubCertPerms    = 0755
	TrustedCertsDir = "trusted_certificates"
)
