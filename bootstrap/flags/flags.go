package flags

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/dhe-deploy"
	"github.com/docker/dhe-deploy/hubconfig/defaultconfigs"
)

var (
	DropperAndConfigFlags = []cli.Flag{
		UCPInsecureTLSFlag,
		UCPCAFlag,
	}
	DropperFlags = []cli.Flag{
		UCPHostFlag,
		UsernameFlag,
		PasswordFlag,
		DebugFlag,
		HubUsernameFlag,
		HubPasswordFlag,
	}
	// XXX: Proxy settings might have to be both dropper options and config options, so should we have two separate flags? same with insecure/CAs?
	ConfigFlags = []cli.Flag{
		HTTPProxyFlag,
		HTTPSProxyFlag,
		NoProxyFlag,
		ReplicaHTTPPortFlag,
		ReplicaHTTPSPortFlag,
		LogProtocolFlag,
		LogHostFlag,
		LogLevelFlag,
		LogTLSCACertFlag,
		LogTLSCertFlag,
		LogTLSKeyFlag,
		LogTLSSkipVerifyFlag,
		DTRHostFlag,
		EnablePProfFlag,
		EtcdHeartbeatIntervalFlag,
		EtcdElectionTimeoutFlag,
		EtcdSnapshotCountFlag,
	}
	InstallFlags     = []cli.Flag{}
	JoinFlags        = []cli.Flag{}
	ReconfigureFlags = []cli.Flag{}
	RestoreFlags     = []cli.Flag{}
	BackupFlags      = []cli.Flag{}
	MigrateFlags     = []cli.Flag{}
	DumpCertsFlags   = []cli.Flag{}
	RemoveFlags      = []cli.Flag{}
	UpgradeFlags     = []cli.Flag{}
)

func init() {
	// In production we don't allow certain flags to be passed
	// Specifically, we refuse to allow our customers to run a standalone enzi
	// or deploy DTR without using UCP
	// Also, since there are no problems with the way we use jwts for docker commands
	// we don't need to expose the UseBundle flag
	if !deploy.IsProduction() {
		DropperAndConfigFlags = append(DropperAndConfigFlags,
			EnziInsecureTLSFlag,
			EnziCAFlag,
		)
		DropperFlags = append(DropperFlags,
			NoUCPFlag,
			UseBundleFlag,
		)
		ConfigFlags = append(ConfigFlags, EnziHostFlag)
		MigrateFlags = append(MigrateFlags, EnziHostFlag)
	}
	InstallFlags = append(append(append(DropperFlags, ConfigFlags...), DropperAndConfigFlags...),
		UCPNodeFlag,
		ReplicaIDFlag,
		UnsafeFlag,
		ExtraEnvsFlag,
	)
	JoinFlags = append(append(DropperFlags, DropperAndConfigFlags...),
		UCPNodeFlag,
		ReplicaIDFlag,
		UnsafeFlag,
		ExistingReplicaIDFlag,
		ReplicaHTTPPortFlag,
		ReplicaHTTPSPortFlag,
		SkipNetworkTestFlag,
		ExtraEnvsFlag,
	)
	RemoveFlags = append(append(DropperFlags, DropperAndConfigFlags...),
		ForceRemoveFlag,
		ReplicaIDFlag,
		ExistingReplicaIDFlag,
	)
	ReconfigureFlags = append(append(append(DropperFlags, ConfigFlags...), DropperAndConfigFlags...),
		ExistingReplicaIDFlag,
	)
	RestoreFlags = append(append(append(DropperFlags, ConfigFlags...), DropperAndConfigFlags...),
		UCPNodeFlag,
		ReplicaIDFlag,
		ConfigOnlyFlag,
	)
	BackupFlags = append(append(DropperFlags, DropperAndConfigFlags...),
		ExistingReplicaIDFlag,
		ConfigOnlyFlag,
	)
	MigrateFlags = append(append(append(DropperFlags, MigrateFlags...), DropperAndConfigFlags...),
		RunFullMigrationFlag,
		DTRHostFlag,
		DTRInsecureTLSFlag,
		DTRCAFlag,
		HTTPProxyFlag,
		HTTPSProxyFlag,
		NoProxyFlag,
	)
	DumpCertsFlags = append(append(DropperFlags, DropperAndConfigFlags...),
		ExistingReplicaIDFlag,
	)
	UpgradeFlags = append(append(append(DropperFlags, ConfigFlags...), DropperAndConfigFlags...),
		ExistingReplicaIDFlag,
	)
}

var (
	Debug     bool
	DebugFlag = cli.BoolFlag{
		Name:        "debug",
		Usage:       "Enable debug mode, provides additional logging",
		EnvVar:      "DEBUG",
		Destination: &Debug,
	}

	NoUCP     bool
	NoUCPFlag = cli.BoolFlag{
		Name:        "noucp",
		EnvVar:      "NO_UCP",
		Destination: &NoUCP,
	}
	UCPHostRaw  string
	UCPHostFlag = cli.StringFlag{
		Name:        "ucp-url",
		Usage:       "Specify the UCP controller URL including domain and port",
		EnvVar:      "UCP_URL",
		Destination: &UCPHostRaw,
	}
	UCPInsecureTLS     bool
	UCPInsecureTLSFlag = cli.BoolFlag{
		Name:        "ucp-insecure-tls",
		Usage:       "Disable TLS verification for UCP",
		EnvVar:      "UCP_INSECURE_TLS",
		Destination: &UCPInsecureTLS,
	}
	UCPCA     string
	UCPCAFlag = cli.StringFlag{
		Name:        "ucp-ca",
		Usage:       "Use a PEM-encoded TLS CA certificate for UCP",
		EnvVar:      "UCP_CA",
		Destination: &UCPCA,
	}
	UCPNode     string
	UCPNodeFlag = cli.StringFlag{
		Name:        "ucp-node",
		Usage:       "Specify the host to install Docker Trusted Registry",
		EnvVar:      "UCP_NODE",
		Destination: &UCPNode,
	}
	UseBundle     bool
	UseBundleFlag = cli.BoolFlag{
		Name:        "use-bundle",
		Usage:       "Alternative method for authenticating on UCP using a user bundle",
		EnvVar:      "USE_BUNDLE",
		Destination: &UseBundle,
	}
	Username     string
	UsernameFlag = cli.StringFlag{
		Name:        "ucp-username",
		Usage:       "Specify the UCP admin username",
		EnvVar:      "UCP_USERNAME",
		Destination: &Username,
	}
	Password     string
	PasswordFlag = cli.StringFlag{
		Name:        "ucp-password",
		Usage:       "Specify the UCP admin password",
		EnvVar:      "UCP_PASSWORD",
		Destination: &Password,
	}
	HubUsername     string
	HubUsernameFlag = cli.StringFlag{
		Name:        "hub-username",
		Usage:       "Specify the Docker Hub username for pulling images",
		EnvVar:      "HUB_USERNAME",
		Destination: &HubUsername,
	}
	HubPassword     string
	HubPasswordFlag = cli.StringFlag{
		Name:        "hub-password",
		Usage:       "Specify the Docker Hub password for pulling images",
		EnvVar:      "HUB_PASSWORD",
		Destination: &HubPassword,
	}

	ReplicaID     string
	ReplicaIDFlag = cli.StringFlag{
		Name:        "replica-id",
		Usage:       "Specify the replica ID. Must be unique per replica, leave blank for random",
		EnvVar:      deploy.ReplicaIDEnvVar,
		Destination: &ReplicaID,
	}
	ExistingReplicaID     string
	ExistingReplicaIDFlag = cli.StringFlag{
		Name:        "existing-replica-id",
		Usage:       "ID of an existing replica in a cluster",
		EnvVar:      "DTR_EXISTING_REPLICA_ID",
		Destination: &ExistingReplicaID,
	}
	SkipNetworkTest     bool
	SkipNetworkTestFlag = cli.BoolFlag{
		Name:        "skip-network-test",
		Usage:       "Enable this flag to skip the overlay networking test",
		EnvVar:      "DTR_SKIP_NETWORK_TEST",
		Destination: &SkipNetworkTest,
	}
	Unsafe     bool
	UnsafeFlag = cli.BoolFlag{
		Name:        "unsafe",
		Usage:       "Enable this flag to skip safety checks when installing or joining",
		EnvVar:      "DTR_UNSAFE",
		Destination: &Unsafe,
	}

	EnziHostRaw  string
	EnziHostFlag = cli.StringFlag{
		Name:        "enzi-host",
		Usage:       "Specify the Enzi host using the host[:port] format",
		EnvVar:      "ENZI_HOST",
		Destination: &EnziHostRaw,
	}
	EnziInsecureTLS     bool
	EnziInsecureTLSFlag = cli.BoolFlag{
		Name:        "enzi-insecure-tls",
		Usage:       "Disable TLS verification for Enzi",
		EnvVar:      "ENZI_TLS_INSECURE",
		Destination: &EnziInsecureTLS,
	}
	EnziCA     string
	EnziCAFlag = cli.StringFlag{
		Name:        "enzi-ca",
		Usage:       "Use a PEM-encoded TLS CA certificate for Enzi",
		EnvVar:      "ENZI_TLS_CA",
		Destination: &EnziCA,
	}

	DTRHostRaw  string
	DTRHostFlag = cli.StringFlag{
		Name:        "dtr-external-url",
		Usage:       "Specify the external domain name and port for DTR. If using a load balancer, use its external URL instead.",
		EnvVar:      "DTR_EXTERNAL_URL",
		Destination: &DTRHostRaw,
	}
	DTRInsecureTLS     bool
	DTRInsecureTLSFlag = cli.BoolFlag{
		Name:        "dtr-insecure-tls",
		Usage:       "Disable TLS verification for DTR",
		EnvVar:      "DTR_INSECURE_TLS",
		Destination: &DTRInsecureTLS,
	}
	DTRCA     string
	DTRCAFlag = cli.StringFlag{
		Name:        "dtr-ca",
		Usage:       "PEM-encoded TLS CA cert for DTR",
		EnvVar:      "DTR_CA",
		Destination: &DTRCA,
	}
	ReplicaHTTPPort     int
	ReplicaHTTPPortFlag = cli.IntFlag{
		Name:        "replica-http-port",
		Usage:       "Specify the public HTTP port for the DTR replica; 0 means unchanged/default",
		EnvVar:      "REPLICA_HTTP_PORT",
		Destination: &ReplicaHTTPPort,
	}
	ReplicaHTTPSPort     int
	ReplicaHTTPSPortFlag = cli.IntFlag{
		Name:        "replica-https-port",
		Usage:       "Specify the public HTTPS port for the DTR replica; 0 means unchanged/default",
		EnvVar:      "REPLICA_HTTPS_PORT",
		Destination: &ReplicaHTTPSPort,
	}
	ExtraEnvs     string
	ExtraEnvsFlag = cli.StringFlag{
		Name:        "extra-envs",
		Usage:       fmt.Sprintf("List of extra environment variables to use for deploying the DTR containers for the replica. This can be used to specify swarm constraints. Separate the environment variables with ampersands (&). You can escape actual ampersands with backslashes (\\). Can't be used in combination with --%s", UCPNodeFlag.Name),
		EnvVar:      "EXTRA_ENVS",
		Destination: &ExtraEnvs,
	}
	RunFullMigration     bool
	RunFullMigrationFlag = cli.BoolFlag{
		Name:        "run-full-migration",
		Usage:       "Run full migration procedure instead of dumping configurations",
		EnvVar:      "RUN_FULL_MIGRATION",
		Destination: &RunFullMigration,
	}
	LogProtocol     string
	LogProtocolFlag = cli.StringFlag{
		Name:        "log-protocol",
		Usage:       fmt.Sprintf("The protocol for sending container logs: tcp, tcp+tls, udp or internal. Default: %s", defaultconfigs.DefaultHAConfig.LogProtocol),
		EnvVar:      "LOG_PROTOCOL",
		Destination: &LogProtocol,
	}
	LogHostRaw  string
	LogHostFlag = cli.StringFlag{
		Name:        "log-host",
		Usage:       fmt.Sprintf("Endpoint to send logs to, required if --%s is tcp or udp", LogProtocolFlag.Name),
		EnvVar:      "LOG_HOST",
		Destination: &LogHostRaw,
	}
	LogLevel     string
	LogLevelFlag = cli.StringFlag{
		Name:        "log-level",
		Usage:       fmt.Sprintf("Log level for container logs. Default: %s", defaultconfigs.DefaultHAConfig.LogLevel),
		EnvVar:      "LOG_LEVEL",
		Destination: &LogLevel,
	}
	LogTLSCACert     string
	LogTLSCACertFlag = cli.StringFlag{
		Name:        "log-tls-ca-cert",
		Usage:       "PEM-encoded TLS CA cert for DTR logging driver. This option is ignored if the address protocol is not tcp+tls.",
		EnvVar:      "LOG_TLS_CA_CERT",
		Destination: &LogTLSCACert,
	}
	LogTLSCert     string
	LogTLSCertFlag = cli.StringFlag{
		Name:        "log-tls-cert",
		Usage:       "PEM-encoded TLS cert for DTR logging driver. This option is ignored if the address protocol is not tcp+tls.",
		EnvVar:      "LOG_TLS_CERT",
		Destination: &LogTLSCert,
	}
	LogTLSKey     string
	LogTLSKeyFlag = cli.StringFlag{
		Name:        "log-tls-key",
		Usage:       "PEM-encoded TLS key for DTR logging driver. This option is ignored if the address protocol is not tcp+tls.",
		EnvVar:      "LOG_TLS_KEY",
		Destination: &LogTLSKey,
	}
	LogTLSSkipVerify     bool
	LogTLSSkipVerifyFlag = cli.BoolFlag{
		Name:        "log-tls-skip-verify",
		Usage:       "Configures DTR logging driver's TLS verification. This verification is enabled by default, but it can be overrided by setting this option to true. This option is ignored if the address protocol is not tcp+tls.",
		EnvVar:      "LOG_TLS_SKIP_VERIFY",
		Destination: &LogTLSSkipVerify,
	}
	HTTPProxy     string
	HTTPProxyFlag = cli.StringFlag{
		Name:        "http-proxy",
		Usage:       "Set the HTTP proxy for outgoing requests",
		EnvVar:      "DTR_HTTP_PROXY",
		Destination: &HTTPProxy,
	}
	HTTPSProxy     string
	HTTPSProxyFlag = cli.StringFlag{
		Name:        "https-proxy",
		Usage:       "Set the HTTPS proxy for outgoing requests",
		EnvVar:      "DTR_HTTPS_PROXY",
		Destination: &HTTPSProxy,
	}
	NoProxy     string
	NoProxyFlag = cli.StringFlag{
		Name:        "no-proxy",
		Usage:       "Set the list of domains to not proxy to",
		EnvVar:      "DTR_NO_PROXY",
		Destination: &NoProxy,
	}
	ConfigOnly     bool
	ConfigOnlyFlag = cli.BoolFlag{
		Name:        "config-only",
		Usage:       "Backup/restore only the configurations of DTR and not the database",
		EnvVar:      "DTR_CONFIG_ONLY",
		Destination: &ConfigOnly,
	}
	ForceRemove     bool
	ForceRemoveFlag = cli.BoolFlag{
		Name:        "force-remove",
		Usage:       fmt.Sprintf("Force removal of replica even if it can break your cluster's state. Necessary only when --%s == --%s.", ExistingReplicaIDFlag.Name, ReplicaIDFlag.Name),
		EnvVar:      "DTR_CONFIG_ONLY",
		Destination: &ForceRemove,
	}
	EnablePProf     bool
	EnablePProfFlag = cli.BoolFlag{
		Name:        "enable-pprof",
		Usage:       "Enables pprof profiling of the server",
		EnvVar:      "DTR_PPROF",
		Destination: &EnablePProf,
	}
	EtcdHeartbeatInterval     int
	EtcdHeartbeatIntervalFlag = cli.IntFlag{
		Name:        "etcd-heartbeat-interval",
		Usage:       "Set etcd's frequency (ms) that its leader will notify followers that it is still the leader.",
		EnvVar:      "ETCD_HEARTBEAT_INTERVAL",
		Destination: &EtcdHeartbeatInterval,
	}
	EtcdElectionTimeout     int
	EtcdElectionTimeoutFlag = cli.IntFlag{
		Name:        "etcd-election-timeout",
		Usage:       "Set etcd's timeout (ms) for how long a follower node will go without hearing a heartbeat before attempting to become leader itself.",
		EnvVar:      "ETCD_ELECTION_TIMEOUT",
		Destination: &EtcdElectionTimeout,
	}
	EtcdSnapshotCount     int
	EtcdSnapshotCountFlag = cli.IntFlag{
		Name:        "etcd-snapshot-count",
		Usage:       "Set etcd's number of changes before creating a snapshot.",
		EnvVar:      "ETCD_SNAPSHOT_COUNT",
		Destination: &EtcdSnapshotCount,
	}
)

func TrimProto(url string) string {
	return strings.Split(strings.TrimPrefix(strings.TrimPrefix(url, "http://"), "https://"), "/")[0]
}

func UCPHost() string {
	return TrimProto(UCPHostRaw)
}
func DTRHost() string {
	return TrimProto(DTRHostRaw)
}
func EnziHost() string {
	return TrimProto(EnziHostRaw)
}
func LogHost() string {
	return TrimProto(LogHostRaw)
}

func UsageFor(flag cli.Flag) string {
	switch matchedFlag := flag.(type) {
	case cli.StringFlag:
		return matchedFlag.Usage
	case cli.BoolFlag:
		return matchedFlag.Usage
	case cli.IntFlag:
		return matchedFlag.Usage
	default:
		logrus.Warnf("Could not get usage for a flag of unknown type: %v", flag)
		os.Exit(1)
	}
	return ""
}
func EnvFor(flag cli.Flag) string {
	switch matchedFlag := flag.(type) {
	case cli.StringFlag:
		return matchedFlag.EnvVar
	case cli.BoolFlag:
		return matchedFlag.EnvVar
	case cli.IntFlag:
		return matchedFlag.EnvVar
	default:
		logrus.Warnf("Could not get environment variable for a flag of unknown type: %v", flag)
		os.Exit(1)
	}
	return ""
}
func String(flag cli.Flag) string {
	switch matchedFlag := flag.(type) {
	case cli.StringFlag:
		return *matchedFlag.Destination
	case cli.BoolFlag:
		return strconv.FormatBool(*matchedFlag.Destination)
	case cli.IntFlag:
		return strconv.FormatInt(int64(*matchedFlag.Destination), 10)
	default:
		logrus.Warnf("Could not get value for a flag of unknown type: %v", flag)
		os.Exit(1)
	}
	return ""
}

func IsSet(c *cli.Context, flag cli.Flag) bool {
	switch matchedFlag := flag.(type) {
	case cli.StringFlag:
		if *matchedFlag.Destination != "" {
			return true
		}
		_, found := os.LookupEnv(matchedFlag.EnvVar)
		return c.IsSet(matchedFlag.Name) || found
	case cli.BoolFlag:
		// XXX - is this correct?
		if *matchedFlag.Destination {
			return true
		}
		_, found := os.LookupEnv(matchedFlag.EnvVar)
		return c.IsSet(matchedFlag.Name) || found
	case cli.IntFlag:
		if *matchedFlag.Destination > 0 {
			return true
		}
		_, found := os.LookupEnv(matchedFlag.EnvVar)
		return c.IsSet(matchedFlag.Name) || found
	default:
		logrus.Warnf("Could not run IsSet on a flag of unknown type: %v", flag)
		os.Exit(1)
	}
	return false
}
