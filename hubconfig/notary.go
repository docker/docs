package hubconfig

// TODO: these should be in notary

// NotaryTLSServerOptions represents the server certificates required to
// set up a TLS-enabled server, and the CA used to validate the any client certs
// for mutual TLS
type NotaryTLSServerOptions struct {
	ServerCert string `json:"tls_cert_file",`
	ServerKey  string `json:"tls_key_file"`
	ClientCA   string `json:"client_ca_file,omitempty"`
}

// NotaryTLSClientOptions represents the client certificates required to connect
// to a server, and the CA used to validate the server's TLS certificate
type NotaryTLSClientOptions struct {
	ClientCert string `json:"tls_client_cert,omitempty"`
	ClientKey  string `json:"tls_client_key,omitempty"`
	ServerCA   string `json:"tls_ca_file,omitempty"`
}

// NotaryStorage represents the storage section of the config for Notary Server or Signer
type NotaryStorage struct {
	Backend    string `json:"backend"`
	URL        string `json:"db_url,omitempty"`
	DB         string `json:"database,omitempty"`
	ClientCert string `json:"client_cert_file,omitempty"`
	ClientKey  string `json:"client_key_file,omitempty"`
	ServerCA   string `json:"tls_ca_file,omitempty"`
	Username   string `json:"username,omitempty"`
	Password   string `json:"password,omitempty"`
}

// NotaryLogging represents logging options in the config for Notary Server or Signer
type NotaryLogging struct {
	Level string `json:"level,omitempty"`
}

// Note, skip Notary reporting entirely for now

// NotaryListeningServer returns how to configure a notary service to listen on ports
type NotaryListeningServer struct {
	HTTPAddr string `json:"http_addr,omitempty"`
	GRPCAddr string `json:"grpc_addr,omitempty"`
	NotaryTLSServerOptions
}

// NotaryTrustService tells a Notary Server how to connect to a Notary Signer
type NotaryTrustService struct {
	Type         string `json:"type"`
	Hostname     string `json:"hostname,omitempty"`
	Port         string `json:"port,omitempty"`
	KeyAlgorithm string `json:"key_algorithm,omitempty"`
	NotaryTLSClientOptions
}

// JSONGarantOptions are a struct that represents garant authentication options
type JSONGarantOptions struct {
	Realm      string `json:"realm"`
	Service    string `json:"service"`
	Issuer     string `json:"issuer"`
	CertBundle string `json:"rootcertbundle"`
}

// NotaryAuth tells a Notary Server what garant server to use
type NotaryAuth struct {
	Type    string            `json:"type"`
	Options JSONGarantOptions `json:"options,omitempty"`
}

// NotaryServerConfig is the config file format for running a Notary Server
type NotaryServerConfig struct {
	Server       NotaryListeningServer `json:"server"`
	TrustService NotaryTrustService    `json:"trust_service"`
	Auth         NotaryAuth            `json:"auth,omitempty"`
	Storage      NotaryStorage         `json:"storage"`
	Logging      NotaryLogging         `json:"logging,omitempty"`
}

// NotarySignerConfig is the config file format for running a Notary Signer
type NotarySignerConfig struct {
	Server NotaryListeningServer `json:"server"`

	Storage NotaryStorage `json:"storage"`
	Logging NotaryLogging `json:"logging,omitempty"`
}
