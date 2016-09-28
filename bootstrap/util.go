package bootstrap

import (
	"bufio"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/coreos/etcd/client"
	"github.com/docker/dhe-deploy"
	"github.com/docker/dhe-deploy/bootstrap/flags"
	"github.com/docker/dhe-deploy/bootstrap/ucpclient"
	"github.com/docker/dhe-deploy/manager/schema"
	jschema "github.com/docker/dhe-deploy/pkg/jobrunner-framework/schema"
	"github.com/docker/dhe-deploy/shared/containers"
	"github.com/docker/dhe-deploy/shared/dtrutil"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/network"
	"gopkg.in/dancannon/gorethink.v2"

	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/net/context"
)

func IsPhase2() bool {
	phase2, _ := strconv.ParseBool(os.Getenv(deploy.Phase2EnvVar))
	return phase2
}

// SetEnvFromFlags takes a list of flags and passes them on from the current command to
// the next container through environment variables. This makes it easy to pass parameters
// from phase 1 to phase 2 for various commands
func SetEnvFromFlags(c *cli.Context, container *containers.ContainerConfig, theFlags ...cli.Flag) {
	if container.Environment == nil {
		container.Environment = map[string]string{}
	}
	for _, flag := range theFlags {
		if flags.IsSet(c, flag) {
			container.Environment[flags.EnvFor(flag)] = flags.String(flag)
		}
	}
}

func PromptIfNotSet(c *cli.Context, cliFlags ...cli.StringFlag) {
	for _, flag := range cliFlags {
		if flags.IsSet(c, flag) {
			continue
		}

		if flag == flags.PasswordFlag {
			*flag.Destination = PromptPassword(fmt.Sprintf("%s: ", flags.PasswordFlag.Name))
		} else {
			PromptFlag(flag)
		}
	}
}

func PromptFlag(flag cli.Flag) {
	for {
		var err error
		result := PromptString(fmt.Sprintf("%s (%s): ", flag.GetName(), flags.UsageFor(flag)), "")

		switch matchedFlag := flag.(type) {
		case cli.StringFlag:
			*matchedFlag.Destination = result
			err = nil
		case cli.BoolFlag:
			*matchedFlag.Destination, err = strconv.ParseBool(result)
		case cli.IntFlag:
			var i int64
			i, err = strconv.ParseInt(result, 10, 32)
			*matchedFlag.Destination = int(i)
		default:
			log.Warnf("Could not set destination for a flag of unknown type: %v", flag)
			os.Exit(1)
		}
		if err == nil {
			break
		}
		log.Printf("Failed to parse input, try again: %s\n", err.Error())
	}
}

func PromptString(prompt, defaultStr string) string {
	var line string
	var err error
	bio := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(log.StandardLogger().Out, prompt)
		line, err = bio.ReadString('\n')
		if err != nil {
			log.Errorf("unable to read input: %s", err)
			os.Exit(1)
		}
		line = strings.TrimSpace(line)
		if line != "" {
			break
		} else if defaultStr != "" {
			return defaultStr
		}
	}
	return line
}

func PromptStringRepeated(prompt, defaultStr string, answers []string) string {
	stringInSlice := func(s string, sl []string) bool {
		for _, s2 := range sl {
			if s == strings.ToUpper(s2) {
				return true
			}
		}
		return false
	}
	resp := strings.ToUpper(PromptString(prompt, defaultStr))
	for !stringInSlice(resp, answers) {
		resp = strings.ToUpper(PromptString(prompt, defaultStr))
	}
	return resp
}

func PromptPassword(prompt string) string {
	var password string
	for {
		logWriter := log.StandardLogger().Out
		fmt.Fprint(logWriter, prompt)
		passwordBytes, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			if strings.Contains(err.Error(), "inappropriate ioctl for device") {
				log.Error("cannot prompt for a password because the container was not started in tty mode (-t)")
				os.Exit(1)
			}
			log.Errorf("unable to scan password: %s", err)
			os.Exit(1)
		}
		fmt.Fprintln(logWriter)

		password = string(passwordBytes)
		if password != "" {
			break
		}
	}
	return password
}

func CreateManager(replicaID string) schema.Manager {
	schemaManager := schema.NewManager(func() (*gorethink.Session, error) {
		return dtrutil.GetRethinkSession(replicaID)
	})

	return schemaManager
}

func SetupReplicaAndNode(c *cli.Context, bs Bootstrap) error {
	if bs.GetReplicaID() == "" {
		bs.SetReplicaID(flags.ReplicaID)
	}
	return SetupNode(c, bs)
}

func SetupNode(c *cli.Context, bs Bootstrap) error {
	return SetupNodeWithName(c, bs, deploy.BootstrapPhase2ContainerName)
}

func GetNodeWithContainer(bs Bootstrap, containerName string) (string, error) {
	nodeName, err := bs.ContainerNode(containerName)
	if err != nil {
		log.Errorf("Couldn't figure out node for the second phase: %s", err)
		return "", err
	}
	return nodeName, nil
}

func SetupNodeWithName(c *cli.Context, bs Bootstrap, containerName string) error {
	// Get node name for bootstrap.  Docker Swarm sets the node information,
	// however, Docker Engine does not.  If we're not using UCP, don't bother
	// setting the node information.
	if !flags.NoUCP {
		nodeName, err := bs.ContainerNode(containerName)
		if err != nil {
			log.Errorf("Couldn't figure out node for the second phase: %s", err)
			return err
		}
		bs.SetNodeName(nodeName)
		log.Debugf("Node name is '%s'", nodeName)
	}
	return nil
}

func GetBridgeNetworkName(c *cli.Context, bs Bootstrap) string {
	if flags.NoUCP {
		return deploy.BridgeNetworkName
	}
	return fmt.Sprintf("%s/%s", bs.GetNodeName(), deploy.BridgeNetworkName)
}

func SetupNodePortConstraints(c *cli.Context, bs Bootstrap) (*[]string, error) {
	pubPorts := []int{}
	var err error

	log.Debugf("replica http port = %d", flags.ReplicaHTTPPort)
	log.Debugf("replica https port = %d", flags.ReplicaHTTPSPort)

	if flags.ReplicaHTTPPort > 0 {
		pubPorts = append(pubPorts, flags.ReplicaHTTPPort)
	} else {
		pubPorts = append(pubPorts, deploy.AdminPort)
	}

	if flags.ReplicaHTTPSPort > 0 {
		pubPorts = append(pubPorts, flags.ReplicaHTTPSPort)
	} else {
		pubPorts = append(pubPorts, deploy.AdminTlsPort)
	}

	// need to find a node that has a port open
	nodes, err := bs.CreateNodeConstraintsFromPort(pubPorts)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

func CheckConstraintsError(err error) {
	if strings.Contains(err.Error(), "Unable to find a node that satisfies the following conditions") {
		port := flags.ReplicaHTTPPort
		if port == 0 {
			port = deploy.AdminPort
		}
		port2 := flags.ReplicaHTTPSPort
		if port2 == 0 {
			port2 = deploy.AdminTlsPort
		}
		log.Errorf("Make sure that there is a node in your UCP cluster where port %d and port %d are open. Also confirm that all UCP node have the DTR images or can obtain them from Docker Hub.", port, port2)
	}
	return
}

func Phase2NetworkConnect(bs Bootstrap, name string) error {
	log.Infof("Connecting to network: %s", name)
	// XXX: begin workaround for docker 1.11-rc2 bug with connecting to networks
	err := dtrutil.Poll(time.Millisecond*100, 10, func() error {
		log.Info("Waiting for phase2 container to be known to the Docker daemon")
		state, err := bs.ContainerInspect(deploy.BootstrapPhase2ContainerName)
		if err != nil {
			return err
		}
		if !state.State.Running {
			log.Debug("phase2 not in running state")
			return fmt.Errorf("phase2 not in running state")
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Failed to wait for phase2 to enter running state: %s", err)
	}
	log.Debug("phase2 is running now")
	// XXX: end workaround for docker 1.11 bug with connecting to networks during early start

	endpointCfg := &network.EndpointSettings{}
	netId, err := bs.NetworkID(name)
	if err != nil {
		return err
	}
	err = bs.NetworkConnect(netId, deploy.BootstrapPhase2ContainerName, endpointCfg)
	if err != nil {
		log.Errorf("Couldn't attach phase2 container to %s: %s", netId, err)
		return err
	}
	log.Debugf("Connected phase2 container to network %s aka %s", netId, name)
	return err
}

func GetEtcdTransport() (client.CancelableTransport, error) {
	cert, err := containers.EtcdCertStore.LoadKeyPair()
	if err != nil {
		log.Errorf("Couldn't load cert pair: %s", err)
		return nil, err
	}

	caCertPool, err := containers.EtcdCACertStore.LoadCACertPool()
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}

	var etcdTransport client.CancelableTransport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSClientConfig:     tlsConfig,
		TLSHandshakeTimeout: 10 * time.Second,
	}
	return etcdTransport, nil
}

func GetEtcdConn(clientURL string) (client.Client, error) {
	etcdTransport, err := GetEtcdTransport()
	if err != nil {
		log.Errorf("Couldn't create etcd client: %s", err)
		return nil, err
	}

	cfg := client.Config{
		Endpoints: []string{clientURL},
		Transport: etcdTransport,
	}

	etcdClient, err := client.New(cfg)
	if err != nil {
		log.Errorf("Couldn't create etcd client: %s", err)
		return nil, err
	}

	kApi := client.NewKeysAPI(etcdClient)

	// DNS can take a little bit to come up on the overlay network so
	// test to make certain everything has come up
	err = dtrutil.Poll(500*time.Millisecond, 60, func() error {
		var err error
		log.Infof("Connecting to etcd...")
		_, err = kApi.Get(context.Background(), "/", nil)
		if err != nil {
			log.Debugf("Can't connect to etcd yet: %s", err)
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	log.Debugf("Connected to etcd")
	return etcdClient, nil
}

// GetNotaryTransport returns an HTTP transport with Notary CA and the correct client keys set.
func GetNotaryTransport() (*http.Transport, error) {
	serverCertBytes, err := ioutil.ReadFile(containers.NotaryCACertStore.CertPath())
	if err != nil {
		return nil, err
	}
	return dtrutil.HTTPTransport(false, []string{string(serverCertBytes)},
		containers.NotaryCertStore.CertPath(), containers.NotaryCertStore.KeyPath())
}

func ConfigureLogging() {
	if flags.Debug {
		log.SetLevel(log.DebugLevel)
	}
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
}

func MakeRegistryAuth(username, password string) string {
	// should be base64 url encoded
	authCfg := types.AuthConfig{
		Username: username,
		Password: password,
	}

	buf, _ := json.Marshal(authCfg)

	return base64.URLEncoding.EncodeToString(buf)
}

// UCPConnTest is used to double-check that ucp is accessible in case connecting to the overlay network broke the
// networking within the container. This can happen on some kernels. Docker tries to prevent this, but can fail to detect it.
func UCPConnTest() error {
	log.Info("Starting UCP connectivity test")
	client, err := dtrutil.HTTPClient(flags.UCPInsecureTLS, flags.UCPCA)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s", flags.UCPHost()), nil)
	if err != nil {
		return err
	}
	res, err := dtrutil.DoRequestWithClient(req, client)
	if err != nil {
		return fmt.Errorf("Failed to connect to UCP after attaching networks. Make sure that your overlay networking works correctly. Error: %s", err)
	}
	res.Body.Close()
	log.Info("UCP connectivity test passed")
	return nil
}

// UCPCertTest is used to make sure UCP has the correct SAN in its cert for DTR to be able to use it
// with the provided domain name
func UCPCertTest(host string, insecure bool, ca string) error {
	log.Info("Validating UCP cert")
	httpClient, err := dtrutil.HTTPClient(insecure, ca)
	if err != nil {
		return fmt.Errorf("Failed to create http client to test UCP cert: %s", err)
	}
	ucp := ucpclient.New(host, httpClient)
	ca, err = ucp.GetCA()
	if err != nil {
		return fmt.Errorf("Failed to get UCP CA: %s", err)
	}
	log.Debug("UCP CA: %s", ca)

	secureHTTPClient, err := dtrutil.HTTPClient(false, ca)
	ucp = ucpclient.New(host, secureHTTPClient)
	err = ucp.Ping()
	if err != nil {
		return fmt.Errorf("Failed to connect to UCP; make sure that you are using a domain listed in UCP's TLS certificate's subject alternate names: %s", err)
	}

	log.Info("UCP cert validation successful")
	return nil
}

func IsNoSuchImageErr(err string) bool {
	return strings.Contains(err, "No such container") || strings.Contains(err, "not found") || strings.Contains(err, "unable to detect a container with the given ID or name")
}

func IdempotentMsg(command string) string {
	return fmt.Sprintf("%s has failed. Try running it again.", command)
}

func MigrateDatabase(replicaID string, numReplicas uint) (schema.Manager, error) {
	// first handle the jobrunner tables
	dbSession, err := dtrutil.GetRethinkSession(replicaID)
	if err != nil {
		return nil, fmt.Errorf("Failed to create db session: %s", err)
	}

	jobSchemaManager := jschema.NewJobrunnerManager(deploy.JobrunnerDBName, dbSession)
	err = jobSchemaManager.SetupDB(numReplicas)
	if err != nil {
		return nil, fmt.Errorf("Failed to set up jobrunner tables: %s", err)
	}

	// then continue to do the dtr tables the old way
	schemaManager := CreateManager(replicaID)
	if err = schemaManager.Initialize(); err != nil {
		return nil, fmt.Errorf("Couldn't initialize schema manager: %s", err)
	}
	if err = schemaManager.Migrate(); err != nil {
		return nil, fmt.Errorf("Couldn't migrate schema manager: %s", err)
	}
	return schemaManager, nil
}

func RequireNoTTY() error {
	if log.IsTerminal() {
		return fmt.Errorf("Terminal mode detected. This command does not work in terminal mode. Please remove -t from the docker run command's arguments.")
	}
	return nil
}
