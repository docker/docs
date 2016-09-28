package hubconfig

type UserHubConfig struct {
	DTRHost      string `yaml:"dtr_host" json:"dtr_host"`
	AuthBypassCA string `yaml:"auth_bypass_ca" json:"auth_bypass_ca"`
	AuthBypassOU string `yaml:"auth_bypass_ou" json:"auth_bypass_ou"`
	WebTLSCert   string `yaml:"web_tls_cert"`
	WebTLSKey    string `yaml:"web_tls_key"`
	WebTLSCA     string `yaml:"web_tls_ca"`

	// DisableUpgrades prevents handlers wrapped with the `upgradeHandlerMiddleware` from being
	// called if set to true.
	DisableUpgrades bool `yaml:"disable_upgrades" json:"disable_upgrades"`
	// ReleaseChannel overrides the release channel baked into the product, specifies where to check for
	// upgrades if not empty.
	ReleaseChannel     string `yaml:"release_channel" json:"release_channel"`
	ReportAnalytics    bool   `yaml:"report_analytics" json:"report_analytics"`
	AnonymizeAnalytics bool   `yaml:"include_license_id" json:"include_license_id"`

	GCMode string `yaml:"gc_mode" json:"gc_mode"`
}

type ReplicaConfig struct {
	HTTPPort  uint16
	HTTPSPort uint16
	Node      string
	Version   string
}

// XXX: should we add json names?
type HAConfig struct {
	ReplicaConfig    map[string]ReplicaConfig
	HTTPProxy        string
	HTTPSProxy       string
	NoProxy          string
	LogProtocol      string
	LogHost          string
	LogLevel         string
	LogTLSCACert     string
	LogTLSCert       string
	LogTLSKey        string
	LogTLSSkipVerify bool
	// Should these be in userhubconfig?
	// If you screw up their values evereything fails, but they don't have to be read-only
	UCPHost               string
	UCPCA                 string
	UCPVerifyCert         bool
	EnziHost              string
	EnziCA                string
	EnziVerifyCert        bool
	EnablePProf           bool
	EtcdHeartbeatInterval int
	EtcdElectionTimeout   int
	EtcdSnapshotCount     int
}

type LicenseConfig struct {
	KeyID         string `yaml:"key_id" json:"key_id"`
	PrivateKey    string `yaml:"private_key" json:"private_key,omitempty"`
	AutoRefresh   bool   `yaml:"auto_refresh" json:"auto_refresh,omitempty"`
	Authorization string `yaml:"authorization" json:"authorization,omitempty"`
}
