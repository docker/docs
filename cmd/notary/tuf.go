package main

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"crypto/subtle"

	"github.com/Sirupsen/logrus"
	"github.com/docker/distribution/registry/client/auth"
	"github.com/docker/distribution/registry/client/transport"
	"github.com/docker/docker/pkg/term"
	notaryclient "github.com/docker/notary/client"
	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmdTufList = &cobra.Command{
	Use:   "list [ GUN ]",
	Short: "Lists targets for a remote trusted collection.",
	Long:  "Lists all targets for a remote trusted collection identified by the Globally Unique Name. This is an online operation.",
	Run:   tufList,
}

var cmdTufAdd = &cobra.Command{
	Use:   "add [ GUN ] <target> <file>",
	Short: "Adds the file as a target to the trusted collection.",
	Long:  "Adds the file as a target to the local trusted collection identified by the Globally Unique Name. This is an offline operation.  Please then use `publish` to push the changes to the remote trusted collection.",
	Run:   tufAdd,
}

var cmdTufRemove = &cobra.Command{
	Use:   "remove [ GUN ] <target>",
	Short: "Removes a target from a trusted collection.",
	Long:  "Removes a target from the local trusted collection identified by the Globally Unique Name. This is an offline operation.  Please then use `publish` to push the changes to the remote trusted collection.",
	Run:   tufRemove,
}

var cmdTufInit = &cobra.Command{
	Use:   "init [ GUN ]",
	Short: "Initializes a local trusted collection.",
	Long:  "Initializes a local trusted collection identified by the Globally Unique Name. This is an online operation.",
	Run:   tufInit,
}

var cmdTufLookup = &cobra.Command{
	Use:   "lookup [ GUN ] <target>",
	Short: "Looks up a specific target in a remote trusted collection.",
	Long:  "Looks up a specific target in a remote trusted collection identified by the Globally Unique Name.",
	Run:   tufLookup,
}

var cmdTufPublish = &cobra.Command{
	Use:   "publish [ GUN ]",
	Short: "Publishes the local trusted collection.",
	Long:  "Publishes the local trusted collection identified by the Globally Unique Name, sending the local changes to a remote trusted server.",
	Run:   tufPublish,
}

var cmdTufStatus = &cobra.Command{
	Use:   "status [ GUN ]",
	Short: "Displays status of unpublished changes to the local trusted collection.",
	Long:  "Displays status of unpublished changes to the local trusted collection identified by the Globally Unique Name.",
	Run:   tufStatus,
}

var cmdVerify = &cobra.Command{
	Use:   "verify [ GUN ] <target>",
	Short: "Verifies if the content is included in the remote trusted collection",
	Long:  "Verifies if the data passed in STDIN is included in the remote trusted collection identified by the Global Unique Name.",
	Run:   verify,
}

func tufAdd(cmd *cobra.Command, args []string) {
	if len(args) < 3 {
		cmd.Usage()
		fatalf("Must specify a GUN, target, and path to target data")
	}
	parseConfig()

	gun := args[0]
	targetName := args[1]
	targetPath := args[2]

	// no online operations are performed by add so the transport argument
	// should be nil
	nRepo, err := notaryclient.NewNotaryRepository(mainViper.GetString("trust_dir"), gun, getRemoteTrustServer(mainViper), nil, retriever)
	if err != nil {
		fatalf(err.Error())
	}

	target, err := notaryclient.NewTarget(targetName, targetPath)
	if err != nil {
		fatalf(err.Error())
	}
	// If roles is empty, we default to adding to targets
	err = nRepo.AddTarget(target, roles...)
	if err != nil {
		fatalf(err.Error())
	}
	cmd.Printf(
		"Addition of target \"%s\" to repository \"%s\" staged for next publish.\n",
		targetName, gun)
}

func tufInit(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		fatalf("Must specify a GUN")
	}

	parseConfig()
	gun := args[0]

	nRepo, err := notaryclient.NewNotaryRepository(mainViper.GetString("trust_dir"), gun, getRemoteTrustServer(mainViper), getTransport(mainViper, gun, false), retriever)
	if err != nil {
		fatalf(err.Error())
	}

	rootKeyList := nRepo.CryptoService.ListKeys(data.CanonicalRootRole)

	var rootKeyID string
	if len(rootKeyList) < 1 {
		cmd.Println("No root keys found. Generating a new root key...")
		rootPublicKey, err := nRepo.CryptoService.Create(data.CanonicalRootRole, data.ECDSAKey)
		rootKeyID = rootPublicKey.ID()
		if err != nil {
			fatalf(err.Error())
		}
	} else {
		// Choses the first root key available, which is initialization specific
		// but should return the HW one first.
		rootKeyID = rootKeyList[0]
		cmd.Printf("Root key found, using: %s\n", rootKeyID)
	}

	err = nRepo.Initialize(rootKeyID)
	if err != nil {
		fatalf(err.Error())
	}
}

func tufList(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		fatalf("Must specify a GUN")
	}
	parseConfig()
	gun := args[0]

	nRepo, err := notaryclient.NewNotaryRepository(mainViper.GetString("trust_dir"), gun, getRemoteTrustServer(mainViper), getTransport(mainViper, gun, true), retriever)
	if err != nil {
		fatalf(err.Error())
	}

	// Retrieve the remote list of signed targets, prioritizing the passed-in list over targets
	roles = append(roles, data.CanonicalTargetsRole)
	targetList, err := nRepo.ListTargets(roles...)
	if err != nil {
		fatalf(err.Error())
	}

	prettyPrintTargets(targetList, cmd.Out())
}

func tufLookup(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		cmd.Usage()
		fatalf("Must specify a GUN and target")
	}
	parseConfig()

	gun := args[0]
	targetName := args[1]

	nRepo, err := notaryclient.NewNotaryRepository(mainViper.GetString("trust_dir"), gun, getRemoteTrustServer(mainViper), getTransport(mainViper, gun, true), retriever)
	if err != nil {
		fatalf(err.Error())
	}

	target, err := nRepo.GetTargetByName(targetName)
	if err != nil {
		fatalf(err.Error())
	}

	cmd.Println(target.Name, fmt.Sprintf("sha256:%x", target.Hashes["sha256"]), target.Length)
}

func tufStatus(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		fatalf("Must specify a GUN")
	}

	parseConfig()
	gun := args[0]

	nRepo, err := notaryclient.NewNotaryRepository(mainViper.GetString("trust_dir"), gun, getRemoteTrustServer(mainViper), nil, retriever)
	if err != nil {
		fatalf(err.Error())
	}

	cl, err := nRepo.GetChangelist()
	if err != nil {
		fatalf(err.Error())
	}

	if len(cl.List()) == 0 {
		cmd.Printf("No unpublished changes for %s\n", gun)
		return
	}

	cmd.Printf("Unpublished changes for %s:\n\n", gun)
	cmd.Printf("%-10s%-10s%-12s%s\n", "action", "scope", "type", "path")
	cmd.Println("----------------------------------------------------")
	for _, ch := range cl.List() {
		cmd.Printf("%-10s%-10s%-12s%s\n", ch.Action(), ch.Scope(), ch.Type(), ch.Path())
	}
}

func tufPublish(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		fatalf("Must specify a GUN")
	}

	parseConfig()
	gun := args[0]

	cmd.Println("Pushing changes to", gun)

	nRepo, err := notaryclient.NewNotaryRepository(mainViper.GetString("trust_dir"), gun, getRemoteTrustServer(mainViper), getTransport(mainViper, gun, false), retriever)
	if err != nil {
		fatalf(err.Error())
	}

	err = nRepo.Publish()
	if err != nil {
		fatalf(err.Error())
	}
}

func tufRemove(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		cmd.Usage()
		fatalf("Must specify a GUN and target")
	}
	parseConfig()

	gun := args[0]
	targetName := args[1]

	// no online operation are performed by remove so the transport argument
	// should be nil.
	repo, err := notaryclient.NewNotaryRepository(mainViper.GetString("trust_dir"), gun, getRemoteTrustServer(mainViper), nil, retriever)
	if err != nil {
		fatalf(err.Error())
	}
	// If roles is empty, we default to removing from targets
	err = repo.RemoveTarget(targetName, roles...)
	if err != nil {
		fatalf(err.Error())
	}

	cmd.Printf("Removal of %s from %s staged for next publish.\n", targetName, gun)
}

func verify(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		cmd.Usage()
		fatalf("Must specify a GUN and target")
	}

	parseConfig()

	// Reads all of the data on STDIN
	payload, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fatalf("Error reading content from STDIN: %v", err)
	}

	gun := args[0]
	targetName := args[1]
	nRepo, err := notaryclient.NewNotaryRepository(mainViper.GetString("trust_dir"), gun, getRemoteTrustServer(mainViper), getTransport(mainViper, gun, true), retriever)
	if err != nil {
		fatalf(err.Error())
	}

	target, err := nRepo.GetTargetByName(targetName)
	if err != nil {
		logrus.Error("notary: data not present in the trusted collection.")
		os.Exit(-11)
	}

	// Create hasher and hash data
	stdinHash := sha256.Sum256(payload)
	serverHash := target.Hashes["sha256"]

	if subtle.ConstantTimeCompare(stdinHash[:], serverHash) == 0 {
		logrus.Error("notary: data not present in the trusted collection.")
		os.Exit(1)
	} else {
		_, _ = os.Stdout.Write(payload)
	}
	return
}

type passwordStore struct {
	anonymous bool
}

func (ps passwordStore) Basic(u *url.URL) (string, string) {
	if ps.anonymous {
		return "", ""
	}

	stdin := bufio.NewReader(os.Stdin)
	fmt.Fprintf(os.Stdout, "Enter username: ")

	userIn, err := stdin.ReadBytes('\n')
	if err != nil {
		logrus.Errorf("error processing username input: %s", err)
		return "", ""
	}

	username := strings.TrimSpace(string(userIn))

	state, err := term.SaveState(0)
	if err != nil {
		logrus.Errorf("error saving terminal state, cannot retrieve password: %s", err)
		return "", ""
	}
	term.DisableEcho(0, state)
	defer term.RestoreTerminal(0, state)

	fmt.Fprintf(os.Stdout, "Enter password: ")

	userIn, err = stdin.ReadBytes('\n')
	fmt.Fprintln(os.Stdout)
	if err != nil {
		logrus.Errorf("error processing password input: %s", err)
		return "", ""
	}
	password := strings.TrimSpace(string(userIn))

	return username, password
}

// getTransport returns an http.RoundTripper to be used for all http requests.
// It correctly handles the auth challenge/credentials required to interact
// with a notary server over both HTTP Basic Auth and the JWT auth implemented
// in the notary-server
// The readOnly flag indicates if the operation should be performed as an
// anonymous read only operation. If the command entered requires write
// permissions on the server, readOnly must be false
func getTransport(config *viper.Viper, gun string, readOnly bool) http.RoundTripper {
	// Attempt to get a root CA from the config file. Nil is the host defaults.
	rootCAFile := config.GetString("remote_server.root_ca")
	if rootCAFile != "" {
		// If we haven't been given an Absolute path, we assume it's relative
		// from the configuration directory (~/.notary by default)
		if !filepath.IsAbs(rootCAFile) {
			rootCAFile = filepath.Join(configPath, rootCAFile)
		}
	}

	insecureSkipVerify := false
	if config.IsSet("remote_server.skipTLSVerify") {
		insecureSkipVerify = config.GetBool("remote_server.skipTLSVerify")
	}
	tlsConfig, err := utils.ConfigureClientTLS(&utils.ClientTLSOpts{
		RootCAFile:         rootCAFile,
		InsecureSkipVerify: insecureSkipVerify,
	})
	if err != nil {
		logrus.Fatal("Unable to configure TLS: ", err.Error())
	}

	base := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig:     tlsConfig,
		DisableKeepAlives:   true,
	}
	trustServerURL := getRemoteTrustServer(config)
	return tokenAuth(trustServerURL, base, gun, readOnly)
}

func tokenAuth(trustServerURL string, baseTransport *http.Transport, gun string,
	readOnly bool) http.RoundTripper {

	// TODO(dmcgowan): add notary specific headers
	authTransport := transport.NewTransport(baseTransport)
	pingClient := &http.Client{
		Transport: authTransport,
		Timeout:   5 * time.Second,
	}
	endpoint, err := url.Parse(trustServerURL)
	if err != nil {
		fatalf("Could not parse remote trust server url (%s): %s", trustServerURL, err.Error())
	}
	if endpoint.Scheme == "" {
		fatalf("Trust server url has to be in the form of http(s)://URL:PORT. Got: %s", trustServerURL)
	}
	subPath, err := url.Parse("v2/")
	if err != nil {
		fatalf("Failed to parse v2 subpath. This error should not have been reached. Please report it as an issue at https://github.com/docker/notary/issues: %s", err.Error())
	}
	endpoint = endpoint.ResolveReference(subPath)
	req, err := http.NewRequest("GET", endpoint.String(), nil)
	if err != nil {
		fatalf(err.Error())
	}
	resp, err := pingClient.Do(req)
	if err != nil {
		logrus.Errorf("could not reach %s: %s", trustServerURL, err.Error())
		logrus.Info("continuing in offline mode")
		return nil
	}
	// non-nil err means we must close body
	defer resp.Body.Close()
	if (resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices) &&
		resp.StatusCode != http.StatusUnauthorized {
		// If we didn't get a 2XX range or 401 status code, we're not talking to a notary server.
		// The http client should be configured to handle redirects so at this point, 3XX is
		// not a valid status code.
		logrus.Errorf("could not reach %s: %d", trustServerURL, resp.StatusCode)
		logrus.Info("continuing in offline mode")
		return nil
	}

	challengeManager := auth.NewSimpleChallengeManager()
	if err := challengeManager.AddResponse(resp); err != nil {
		fatalf(err.Error())
	}

	ps := passwordStore{anonymous: readOnly}
	tokenHandler := auth.NewTokenHandler(authTransport, ps, gun, "push", "pull")
	basicHandler := auth.NewBasicHandler(ps)
	modifier := transport.RequestModifier(auth.NewAuthorizer(challengeManager, tokenHandler, basicHandler))
	return transport.NewTransport(baseTransport, modifier)
}

func getRemoteTrustServer(config *viper.Viper) string {
	if remoteTrustServer == "" {
		configRemote := config.GetString("remote_server.url")
		if configRemote != "" {
			remoteTrustServer = configRemote
		} else {
			remoteTrustServer = defaultServerURL
		}
	}
	return remoteTrustServer
}
