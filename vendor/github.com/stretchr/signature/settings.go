package signature

var (
	// PrivateKeyKey is the key (URL field) for the private key.
	PrivateKeyKey string = "private"

	// BodyHashKey is the key (URL field) for the body hash used for signing requests.
	BodyHashKey string = "bodyhash"

	// SignatureKey is the key (URL field) for the signature of requests.
	SignatureKey string = "sign"
)
