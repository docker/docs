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
	"github.com/docker/notary/utils"
	"github.com/spf13/cobra"
)

var cmdTufList = &cobra.Command{
	Use:   "list [ GUN ]",
	Short: "Lists targets for a trusted collection.",
	Long:  "Lists all targets for a trusted collection identified by the Globally Unique Name.",
	Run:   tufList,
}

var cmdTufAdd = &cobra.Command{
	Use:   "add [ GUN ] <target> <file>",
	Short: "adds the file as a target to the trusted collection.",
	Long:  "adds the file as a target to the local trusted collection identified by the Globally Unique Name.",
	Run:   tufAdd,
}

var cmdTufRemove = &cobra.Command{
	Use:   "remove [ GUN ] <target>",
	Short: "Removes a target from a trusted collection.",
	Long:  "removes a target from the local trusted collection identified by the Globally Unique Name.",
	Run:   tufRemove,
}

var cmdTufInit = &cobra.Command{
	Use:   "init [ GUN ]",
	Short: "initializes a local trusted collection.",
	Long:  "initializes a local trusted collection identified by the Globally Unique Name.",
	Run:   tufInit,
}

var cmdTufLookup = &cobra.Command{
	Use:   "lookup [ GUN ] <target>",
	Short: "Looks up a specific target in a trusted collection.",
	Long:  "looks up a specific target in a trusted collection identified by the Globally Unique Name.",
	Run:   tufLookup,
}

var cmdTufPublish = &cobra.Command{
	Use:   "publish [ GUN ]",
	Short: "publishes the local trusted collection.",
	Long:  "publishes the local trusted collection identified by the Globally Unique Name, sending the local changes to a remote trusted server.",
	Run:   tufPublish,
}

var cmdTufStatus = &cobra.Command{
	Use:   "status [ GUN ]",
	Short: "displays status of unpublished changes to the local trusted collection.",
	Long:  "displays status of unpublished changes to the local trusted collection identified by the Globally Unique Name.",
	Run:   tufStatus,
}

var cmdVerify = &cobra.Command{
	Use:   "verify [ GUN ] <target>",
	Short: "verifies if the content is included in the trusted collection",
	Long:  "verifies if the data passed in STDIN is included in the trusted collection identified by the Global Unique Name.",
	Run:   verify,
}

func tufAdd(cmd *cobra.Command, args []string) {
	if len(args) < 3 {
		cmd.Usage()
		fatalf("must specify a GUN, target, and path to target data")
	}

	gun := args[0]
	targetName := args[1]
	targetPath := args[2]

	parseConfig()
	// no online operations are performed by add so the transport argument
	// should be nil
	nRepo, err := notaryclient.NewNotaryRepository(trustDir, gun, getRemoteTrustServer(), nil, retriever)
	if err != nil {
		fatalf(err.Error())
	}

	target, err := notaryclient.NewTarget(targetName, targetPath)
	if err != nil {
		fatalf(err.Error())
	}
	err = nRepo.AddTarget(target)
	if err != nil {
		fatalf(err.Error())
	}
	fmt.Printf("Addition of %s to %s staged for next publish.\n", targetName, gun)
}

func tufInit(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		fatalf("Must specify a GUN")
	}

	gun := args[0]
	parseConfig()

	nRepo, err := notaryclient.NewNotaryRepository(trustDir, gun, getRemoteTrustServer(), getTransport(gun, false), retriever)
	if err != nil {
		fatalf(err.Error())
	}

	keysMap := nRepo.KeyStoreManager.RootKeyStore().ListKeys()

	var rootKeyID string
	if len(keysMap) < 1 {
		fmt.Println("No root keys found. Generating a new root key...")
		rootKeyID, err = nRepo.KeyStoreManager.GenRootKey("ECDSA")
		if err != nil {
			fatalf(err.Error())
		}
	} else {
		// TODO(diogo): ask which root key to use
		for keyID := range keysMap {
			rootKeyID = keyID
		}

		fmt.Printf("Root key found, using: %s\n", rootKeyID)
	}

	rootCryptoService, err := nRepo.KeyStoreManager.GetRootCryptoService(rootKeyID)
	if err != nil {
		fatalf(err.Error())
	}

	err = nRepo.Initialize(rootCryptoService)
	if err != nil {
		fatalf(err.Error())
	}
}

func tufList(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		fatalf("must specify a GUN")
	}
	gun := args[0]
	parseConfig()

	nRepo, err := notaryclient.NewNotaryRepository(trustDir, gun, getRemoteTrustServer(), getTransport(gun, true), retriever)
	if err != nil {
		fatalf(err.Error())
	}

	// Retreive the remote list of signed targets
	targetList, err := nRepo.ListTargets()
	if err != nil {
		fatalf(err.Error())
	}

	// Print all the available targets
	for _, t := range targetList {
		fmt.Printf("%s %x %d\n", t.Name, t.Hashes["sha256"], t.Length)
	}
}

func tufLookup(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		cmd.Usage()
		fatalf("must specify a GUN and target")
	}
	gun := args[0]
	targetName := args[1]
	parseConfig()

	nRepo, err := notaryclient.NewNotaryRepository(trustDir, gun, getRemoteTrustServer(), getTransport(gun, true), retriever)
	if err != nil {
		fatalf(err.Error())
	}

	target, err := nRepo.GetTargetByName(targetName)
	if err != nil {
		fatalf(err.Error())
	}

	fmt.Println(target.Name, fmt.Sprintf("sha256:%x", target.Hashes["sha256"]), target.Length)
}

func tufStatus(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		fatalf("Must specify a GUN")
	}

	gun := args[0]
	parseConfig()

	nRepo, err := notaryclient.NewNotaryRepository(trustDir, gun, getRemoteTrustServer(), nil, retriever)
	if err != nil {
		fatalf(err.Error())
	}

	cl, err := nRepo.GetChangelist()
	if err != nil {
		fatalf(err.Error())
	}

	if len(cl.List()) == 0 {
		fmt.Printf("No unpublished changes for %s\n", gun)
		return
	}

	fmt.Printf("Unpublished changes for %s:\n\n", gun)
	fmt.Printf("%-10s%-10s%-12s%s\n", "action", "scope", "type", "path")
	fmt.Println("----------------------------------------------------")
	for _, ch := range cl.List() {
		fmt.Printf("%-10s%-10s%-12s%s\n", ch.Action(), ch.Scope(), ch.Type(), ch.Path())
	}
}

func tufPublish(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		cmd.Usage()
		fatalf("Must specify a GUN")
	}

	gun := args[0]
	parseConfig()

	fmt.Println("Pushing changes to ", gun, ".")

	nRepo, err := notaryclient.NewNotaryRepository(trustDir, gun, getRemoteTrustServer(), getTransport(gun, false), retriever)
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
		fatalf("must specify a GUN and target")
	}
	gun := args[0]
	targetName := args[1]
	parseConfig()

	// no online operation are performed by remove so the transport argument
	// should be nil.
	repo, err := notaryclient.NewNotaryRepository(trustDir, gun, getRemoteTrustServer(), nil, retriever)
	if err != nil {
		fatalf(err.Error())
	}
	err = repo.RemoveTarget(targetName)
	if err != nil {
		fatalf(err.Error())
	}

	fmt.Printf("Removal of %s from %s staged for next publish.\n", targetName, gun)
}

func verify(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		cmd.Usage()
		fatalf("must specify a GUN and target")
	}
	parseConfig()

	// Reads all of the data on STDIN
	payload, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fatalf("error reading content from STDIN: %v", err)
	}

	gun := args[0]
	targetName := args[1]
	nRepo, err := notaryclient.NewNotaryRepository(trustDir, gun, getRemoteTrustServer(), getTransport(gun, true), retriever)
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

func getTransport(gun string, readOnly bool) http.RoundTripper {
	// Attempt to get a root CA from the config file. Nil is the host defaults.
	rootCAFile := mainViper.GetString("remote_server.root_ca")
	if rootCAFile != "" {
		// If we haven't been given an Absolute path, we assume it's relative
		// from the configuration directory (~/.notary by default)
		if !filepath.IsAbs(rootCAFile) {
			rootCAFile = filepath.Join(configPath, rootCAFile)
		}
	}

	insecureSkipVerify := false
	if mainViper.IsSet("remote_server.skipTLSVerify") {
		insecureSkipVerify = mainViper.GetBool("remote_server.skipTLSVerify")
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

	return tokenAuth(base, gun, readOnly)
}

func tokenAuth(baseTransport *http.Transport, gun string, readOnly bool) http.RoundTripper {
	// TODO(dmcgowan): add notary specific headers
	authTransport := transport.NewTransport(baseTransport)
	pingClient := &http.Client{
		Transport: authTransport,
		Timeout:   5 * time.Second,
	}
	trustServerURL := getRemoteTrustServer()
	endpoint, err := url.Parse(trustServerURL)
	if err != nil {
		fatalf("could not parse remote trust server url (%s): %s", trustServerURL, err.Error())
	}
	if endpoint.Scheme == "" {
		fatalf("trust server url has to be in the form of http(s)://URL:PORT. Got: %s", trustServerURL)
	}
	subPath, err := url.Parse("v2/")
	if err != nil {
		fatalf("failed to parse v2 subpath. This error should not have been reached. Please report it as an issue at https://github.com/docker/notary/issues: %s", err.Error())
	}
	endpoint = endpoint.ResolveReference(subPath)
	req, err := http.NewRequest("GET", endpoint.String(), nil)
	if err != nil {
		fatalf(err.Error())
	}
	resp, err := pingClient.Do(req)
	if err != nil {
		fatalf(err.Error())
	}
	defer resp.Body.Close()

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

func getRemoteTrustServer() string {
	if remoteTrustServer == "" {
		configRemote := mainViper.GetString("remote_server.url")
		if configRemote != "" {
			remoteTrustServer = configRemote
		} else {
			remoteTrustServer = defaultServerURL
		}
	}
	return remoteTrustServer
}
