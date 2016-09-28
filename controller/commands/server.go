package commands

import (
	"io/ioutil"
	"net"
	"net/url"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/orca/controller/api"
	"github.com/docker/orca/controller/manager"
	"github.com/docker/orca/utils"
	"github.com/docker/orca/version"
)

var (
	controllerManager *manager.Manager
	Debug             bool
)

var CmdServer = cli.Command{
	Name:   "server",
	Usage:  "run controller",
	Action: cmdServer,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "listen, l",
			Usage: "listen address",
			Value: ":8080",
		},
		cli.StringFlag{
			Name:  "discovery, k",
			Usage: "discovery address",
			Value: "",
		},
		cli.StringFlag{
			Name:  "discovery-tls-ca-cert",
			Value: "",
			Usage: "discovery tls ca certificate",
		},
		cli.StringFlag{
			Name:  "discovery-tls-cert",
			Value: "",
			Usage: "discovery tls certificate",
		},
		cli.StringFlag{
			Name:  "discovery-tls-key",
			Value: "",
			Usage: "discovery tls key",
		},
		cli.BoolFlag{
			Name:  "disable-usage-info",
			Usage: "disable anonymous usage reporting",
		},
		cli.BoolFlag{
			Name:  "disable-tracking",
			Usage: "disable anonymous tracking",
		},
		cli.StringFlag{
			Name:  "swarm-url",
			Usage: "docker swarm v1 address",
		},
		cli.StringFlag{
			Name:  "docker-url",
			Usage: "docker engine address",
		},
		cli.StringFlag{
			Name:  "tls-ca-cert",
			Value: "",
			Usage: "tls ca certificate",
		},
		cli.StringFlag{
			Name:  "tls-cert",
			Value: "",
			Usage: "tls certificate",
		},
		cli.StringFlag{
			Name:  "tls-key",
			Value: "",
			Usage: "tls key",
		},
		cli.StringFlag{
			Name:  "orca-tls-ca-cert",
			Usage: "Orca TLS CA Cert",
			Value: "",
		},
		cli.StringFlag{
			Name:  "orca-tls-cert",
			Usage: "Orca TLS Cert",
			Value: "",
		},
		cli.StringFlag{
			Name:  "orca-tls-key",
			Usage: "Orca TLS Key",
			Value: "",
		},
		cli.BoolFlag{
			Name:  "allow-insecure",
			Usage: "enable insecure tls communication",
		},
		cli.StringSliceFlag{
			Name:  "auth-whitelist-cidr",
			Usage: "whitelist CIDR to bypass auth",
			Value: &cli.StringSlice{},
		},
		cli.IntFlag{
			Name:  "support-timeout",
			Usage: "Timeout in seconds for individual nodes when getting support logs (default: 60)",
			Value: 60,
		},
		cli.StringFlag{
			Name:  "support-image",
			Usage: "Image to use when generating support dumps (default: docker/ucp-dsinfo:" + version.TagVersion() + ")",
			Value: "docker/ucp-dsinfo:" + version.TagVersion(),
		},
	},
}

func cmdServer(c *cli.Context) {
	listenAddr := c.String("listen")
	authWhitelist := c.StringSlice("auth-whitelist-cidr")
	supportTimeout := c.Int("support-timeout")
	supportImage := c.String("support-image")
	discoveryAddr := c.String("discovery")
	discoveryTlsCaCert := c.String("discovery-tls-ca-cert")
	discoveryTlsCert := c.String("discovery-tls-cert")
	discoveryTlsKey := c.String("discovery-tls-key")

	log.Infof("orca version %s", version.FullVersion())

	if len(authWhitelist) > 0 {
		log.Infof("whitelisting the following subnets: %v", authWhitelist)
	}

	swarmClassicURL, err := url.Parse(c.String("swarm-url"))
	if err != nil {
		log.Fatal(err)
	}
	engineProxyURL, err := url.Parse(c.String("docker-url"))
	if err != nil {
		log.Fatal(err)
	}
	discoveryURL, err := url.Parse(discoveryAddr)
	if err != nil {
		log.Fatalf("unable to parse discovery URL: %s", err)
	}
	hostAddr, _, err := net.SplitHostPort(discoveryURL.Host)
	if err != nil {
		log.Fatalf("unable to split host:port from discovery URL %s: %s", discoveryURL, err)
	}

	tlsCaFilePath := c.String("tls-ca-cert")
	tlsCertFilePath := c.String("tls-cert")
	tlsKeyFilepath := c.String("tls-key")
	allowInsecure := c.Bool("allow-insecure")

	// Read the swarm key, cert, and ca file here because the contents are
	// needed in multiple places.
	swarmKeyPEM, err := ioutil.ReadFile(tlsKeyFilepath)
	if err != nil {
		log.Fatalf("unable to read swarm cluster TLS Key file: %s", err)
	}
	swarmCertPEM, err := ioutil.ReadFile(tlsCertFilePath)
	if err != nil {
		log.Fatalf("unable to read swarm cluster TLS Cert file: %s", err)
	}
	swarmCaPEM, err := ioutil.ReadFile(tlsCaFilePath)
	if err != nil {
		log.Fatalf("unable to read swarm cluster TLS CA file: %s", err)
	}

	log.Debugf("discovery: %s", discoveryAddr)

	client, swarmTransport, err := utils.GetClient(swarmClassicURL.String(), swarmCaPEM, swarmCertPEM, swarmKeyPEM, allowInsecure)
	if err != nil {
		log.Fatal(err)
	}

	// Create an engine-api client with the Engine Proxy
	proxyClient, proxyTransport, err := utils.GetClient(engineProxyURL.String(), swarmCaPEM, swarmCertPEM, swarmKeyPEM, allowInsecure)
	if err != nil {
		log.Fatal(err)
	}

	orcaTLSCertFile := c.String("orca-tls-cert")
	orcaTLSKeyFile := c.String("orca-tls-key")
	orcaTLSCACertFile := c.String("orca-tls-ca-cert")

	// Also read the controller key, cert, and ca file here because the
	// contents are needed in multiple places.
	controllerKeyPEM, err := ioutil.ReadFile(orcaTLSKeyFile)
	if err != nil {
		log.Fatalf("unable to read controller TLS Key file: %s", err)
	}
	controllerCertPEM, err := ioutil.ReadFile(orcaTLSCertFile)
	if err != nil {
		log.Fatalf("unable to read controller TLS Cert file: %s", err)
	}
	controllerCaPEM, err := ioutil.ReadFile(orcaTLSCACertFile)
	if err != nil {
		log.Fatalf("unable to read controller TLS CA file: %s", err)
	}

	dtrUrl := c.String("dtr-url")
	dtrInsecure := c.Bool("dtr-insecure")
	dtrAdmin := c.String("dtr-admin")

	config := &manager.DefaultManagerConfig{
		SwarmClassicURL:        swarmClassicURL,
		EngineProxyURL:         engineProxyURL,
		DiscoveryAddr:          discoveryAddr,
		HostAddr:               hostAddr,
		DiscoveryTLSCaCertPath: discoveryTlsCaCert,
		DiscoveryTLSCertPath:   discoveryTlsCert,
		DiscoveryTLSKeyPath:    discoveryTlsKey,
		Client:                 client,
		ClientTransport:        swarmTransport,
		ProxyClient:            proxyClient,
		ProxyTransport:         proxyTransport,
		SupportTimeout:         supportTimeout,
		SupportImage:           supportImage,
		ControllerCAPEM:        controllerCaPEM,
		ControllerCertPEM:      controllerCertPEM,
		ControllerKeyPEM:       controllerKeyPEM,
		DtrUrl:                 dtrUrl,
		DtrInsecure:            dtrInsecure,
		DtrAdmin:               dtrAdmin,
		SwarmCAPEM:             swarmCaPEM,
		SwarmCertPEM:           swarmCertPEM,
		SwarmKeyPEM:            swarmKeyPEM,
	}

	controllerManager, err := manager.NewDefaultManager(config)
	if err != nil {
		log.Fatal(err)
	}

	log.Debugf("connected to classic swarm: url=%s", swarmClassicURL.String())
	log.Debugf("connected to engine proxy: url=%s", engineProxyURL.String())

	apiConfig := api.ApiConfig{
		ListenAddr:         listenAddr,
		Manager:            controllerManager,
		AuthWhiteListCIDRs: authWhitelist,
		AllowInsecure:      allowInsecure,
		SwarmCAPEM:         swarmCaPEM,
		ControllerCAPEM:    controllerCaPEM,
		ControllerCertPEM:  controllerCertPEM,
		ControllerKeyPEM:   controllerKeyPEM,
		EnableProfiling:    Debug,
	}

	orcaApi, err := api.NewApi(apiConfig)
	if err != nil {
		log.Fatal(err)
	}

	if err := orcaApi.Run(); err != nil {
		log.Fatal(err)
	}
}
