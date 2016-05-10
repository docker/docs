package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/docker/distribution/registry/client/auth"
	"github.com/docker/distribution/registry/client/transport"
	"github.com/docker/docker/pkg/term"
	"github.com/docker/go-connections/tlsconfig"
	"github.com/docker/notary"
	notaryclient "github.com/docker/notary/client"
	"github.com/docker/notary/passphrase"
	"github.com/docker/notary/trustpinning"
	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmdTufListTemplate = usageTemplate{
	Use:   "list [ GUN ]",
	Short: "Lists targets for a remote trusted collection.",
	Long:  "Lists all targets for a remote trusted collection identified by the Globally Unique Name. This is an online operation.",
}

var cmdTufAddTemplate = usageTemplate{
	Use:   "add [ GUN ] <target> <file>",
	Short: "Adds the file as a target to the trusted collection.",
	Long:  "Adds the file as a target to the local trusted collection identified by the Globally Unique Name. This is an offline operation.  Please then use `publish` to push the changes to the remote trusted collection.",
}

var cmdTufAddHashTemplate = usageTemplate{
	Use:   "addhash [ GUN ] <target> <byte size> <hashes>",
	Short: "Adds the byte size and hash(es) as a target to the trusted collection.",
	Long:  "Adds the specified byte size and hash(es) as a target to the local trusted collection identified by the Globally Unique Name. This is an offline operation.  Please then use `publish` to push the changes to the remote trusted collection.",
}

var cmdTufRemoveTemplate = usageTemplate{
	Use:   "remove [ GUN ] <target>",
	Short: "Removes a target from a trusted collection.",
	Long:  "Removes a target from the local trusted collection identified by the Globally Unique Name. This is an offline operation.  Please then use `publish` to push the changes to the remote trusted collection.",
}

var cmdTufInitTemplate = usageTemplate{
	Use:   "init [ GUN ]",
	Short: "Initializes a local trusted collection.",
	Long:  "Initializes a local trusted collection identified by the Globally Unique Name. This is an online operation.",
}

var cmdTufLookupTemplate = usageTemplate{
	Use:   "lookup [ GUN ] <target>",
	Short: "Looks up a specific target in a remote trusted collection.",
	Long:  "Looks up a specific target in a remote trusted collection identified by the Globally Unique Name.",
}

var cmdTufPublishTemplate = usageTemplate{
	Use:   "publish [ GUN ]",
	Short: "Publishes the local trusted collection.",
	Long:  "Publishes the local trusted collection identified by the Globally Unique Name, sending the local changes to a remote trusted server.",
}

var cmdTufStatusTemplate = usageTemplate{
	Use:   "status [ GUN ]",
	Short: "Displays status of unpublished changes to the local trusted collection.",
	Long:  "Displays status of unpublished changes to the local trusted collection identified by the Globally Unique Name.",
}

var cmdTufVerifyTemplate = usageTemplate{
	Use:   "verify [ GUN ] <target>",
	Short: "Verifies if the content is included in the remote trusted collection",
	Long:  "Verifies if the data passed in STDIN is included in the remote trusted collection identified by the Global Unique Name.",
}

type tufCommander struct {
	// these need to be set
	configGetter func() (*viper.Viper, error)
	retriever    passphrase.Retriever

	// these are for command line parsing - no need to set
	roles  []string
	sha256 string
	sha512 string

	input  string
	output string
	quiet  bool
}

func (t *tufCommander) AddToCommand(cmd *cobra.Command) {
	cmd.AddCommand(cmdTufInitTemplate.ToCommand(t.tufInit))
	cmd.AddCommand(cmdTufStatusTemplate.ToCommand(t.tufStatus))
	cmd.AddCommand(cmdTufPublishTemplate.ToCommand(t.tufPublish))
	cmd.AddCommand(cmdTufLookupTemplate.ToCommand(t.tufLookup))

	cmdTufList := cmdTufListTemplate.ToCommand(t.tufList)
	cmdTufList.Flags().StringSliceVarP(
		&t.roles, "roles", "r", nil, "Delegation roles to list targets for (will shadow targets role)")
	cmd.AddCommand(cmdTufList)

	cmdTufAdd := cmdTufAddTemplate.ToCommand(t.tufAdd)
	cmdTufAdd.Flags().StringSliceVarP(&t.roles, "roles", "r", nil, "Delegation roles to add this target to")
	cmd.AddCommand(cmdTufAdd)

	cmdTufRemove := cmdTufRemoveTemplate.ToCommand(t.tufRemove)
	cmdTufRemove.Flags().StringSliceVarP(&t.roles, "roles", "r", nil, "Delegation roles to remove this target from")
	cmd.AddCommand(cmdTufRemove)

	cmdTufAddHash := cmdTufAddHashTemplate.ToCommand(t.tufAddByHash)
	cmdTufAddHash.Flags().StringSliceVarP(&t.roles, "roles", "r", nil, "Delegation roles to add this target to")
	cmdTufAddHash.Flags().StringVar(&t.sha256, notary.SHA256, "", "hex encoded sha256 of the target to add")
	cmdTufAddHash.Flags().StringVar(&t.sha512, notary.SHA512, "", "hex encoded sha512 of the target to add")
	cmd.AddCommand(cmdTufAddHash)

	cmdTufVerify := cmdTufVerifyTemplate.ToCommand(t.tufVerify)
	cmdTufVerify.Flags().StringVarP(&t.input, "input", "i", "", "Read from a file, instead of STDIN")
	cmdTufVerify.Flags().StringVarP(&t.output, "output", "o", "", "Write to a file, instead of STDOUT")
	cmdTufVerify.Flags().BoolVarP(&t.quiet, "quiet", "q", false, "No output except for errors")
	cmd.AddCommand(cmdTufVerify)
}

func (t *tufCommander) tufAddByHash(cmd *cobra.Command, args []string) error {
	if len(args) < 3 || t.sha256 == "" && t.sha512 == "" {
		cmd.Usage()
		return fmt.Errorf("Must specify a GUN, target, byte size of target data, and at least one hash")
	}
	config, err := t.configGetter()
	if err != nil {
		return err
	}

	gun := args[0]
	targetName := args[1]
	targetSize := args[2]

	targetInt64Len, err := strconv.ParseInt(targetSize, 0, 64)
	if err != nil {
		return err
	}

	trustPin, err := getTrustPinning(config)
	if err != nil {
		return err
	}

	// no online operations are performed by add so the transport argument
	// should be nil
	nRepo, err := notaryclient.NewNotaryRepository(
		config.GetString("trust_dir"), gun, getRemoteTrustServer(config), nil, t.retriever, trustPin)
	if err != nil {
		return err
	}

	targetHash := data.Hashes{}
	if t.sha256 != "" {
		if len(t.sha256) != notary.Sha256HexSize {
			return fmt.Errorf("invalid sha256 hex contents provided")
		}
		sha256Hash, err := hex.DecodeString(t.sha256)
		if err != nil {
			return err
		}
		targetHash[notary.SHA256] = sha256Hash
	}
	if t.sha512 != "" {
		if len(t.sha512) != notary.Sha512HexSize {
			return fmt.Errorf("invalid sha512 hex contents provided")
		}
		sha512Hash, err := hex.DecodeString(t.sha512)
		if err != nil {
			return err
		}
		targetHash[notary.SHA512] = sha512Hash
	}

	// Manually construct the target with the given byte size and hashes
	target := &notaryclient.Target{Name: targetName, Hashes: targetHash, Length: targetInt64Len}

	// If roles is empty, we default to adding to targets
	if err = nRepo.AddTarget(target, t.roles...); err != nil {
		return err
	}
	// Include the hash algorithms we're using for pretty printing
	hashesUsed := []string{}
	for hashName := range targetHash {
		hashesUsed = append(hashesUsed, hashName)
	}
	cmd.Printf(
		"Addition of target \"%s\" by %s hash to repository \"%s\" staged for next publish.\n",
		targetName, strings.Join(hashesUsed, ", "), gun)
	return nil
}

func (t *tufCommander) tufAdd(cmd *cobra.Command, args []string) error {
	if len(args) < 3 {
		cmd.Usage()
		return fmt.Errorf("Must specify a GUN, target, and path to target data")
	}
	config, err := t.configGetter()
	if err != nil {
		return err
	}

	gun := args[0]
	targetName := args[1]
	targetPath := args[2]

	trustPin, err := getTrustPinning(config)
	if err != nil {
		return err
	}

	// no online operations are performed by add so the transport argument
	// should be nil
	nRepo, err := notaryclient.NewNotaryRepository(
		config.GetString("trust_dir"), gun, getRemoteTrustServer(config), nil, t.retriever, trustPin)
	if err != nil {
		return err
	}

	target, err := notaryclient.NewTarget(targetName, targetPath)
	if err != nil {
		return err
	}
	// If roles is empty, we default to adding to targets
	if err = nRepo.AddTarget(target, t.roles...); err != nil {
		return err
	}
	cmd.Printf(
		"Addition of target \"%s\" to repository \"%s\" staged for next publish.\n",
		targetName, gun)
	return nil
}

func (t *tufCommander) tufInit(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		cmd.Usage()
		return fmt.Errorf("Must specify a GUN")
	}

	config, err := t.configGetter()
	if err != nil {
		return err
	}
	gun := args[0]

	rt, err := getTransport(config, gun, false)
	if err != nil {
		return err
	}

	trustPin, err := getTrustPinning(config)
	if err != nil {
		return err
	}

	nRepo, err := notaryclient.NewNotaryRepository(
		config.GetString("trust_dir"), gun, getRemoteTrustServer(config), rt, t.retriever, trustPin)
	if err != nil {
		return err
	}

	rootKeyList := nRepo.CryptoService.ListKeys(data.CanonicalRootRole)

	var rootKeyID string
	if len(rootKeyList) < 1 {
		cmd.Println("No root keys found. Generating a new root key...")
		rootPublicKey, err := nRepo.CryptoService.Create(data.CanonicalRootRole, "", data.ECDSAKey)
		rootKeyID = rootPublicKey.ID()
		if err != nil {
			return err
		}
	} else {
		// Choses the first root key available, which is initialization specific
		// but should return the HW one first.
		rootKeyID = rootKeyList[0]
		cmd.Printf("Root key found, using: %s\n", rootKeyID)
	}

	if err = nRepo.Initialize(rootKeyID); err != nil {
		return err
	}
	return nil
}

func (t *tufCommander) tufList(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		cmd.Usage()
		return fmt.Errorf("Must specify a GUN")
	}
	config, err := t.configGetter()
	if err != nil {
		return err
	}
	gun := args[0]

	rt, err := getTransport(config, gun, true)
	if err != nil {
		return err
	}

	trustPin, err := getTrustPinning(config)
	if err != nil {
		return err
	}

	nRepo, err := notaryclient.NewNotaryRepository(
		config.GetString("trust_dir"), gun, getRemoteTrustServer(config), rt, t.retriever, trustPin)
	if err != nil {
		return err
	}

	// Retrieve the remote list of signed targets, prioritizing the passed-in list over targets
	roles := append(t.roles, data.CanonicalTargetsRole)
	targetList, err := nRepo.ListTargets(roles...)
	if err != nil {
		return err
	}

	prettyPrintTargets(targetList, cmd.Out())
	return nil
}

func (t *tufCommander) tufLookup(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		cmd.Usage()
		return fmt.Errorf("Must specify a GUN and target")
	}
	config, err := t.configGetter()
	if err != nil {
		return err
	}

	gun := args[0]
	targetName := args[1]

	rt, err := getTransport(config, gun, true)
	if err != nil {
		return err
	}

	trustPin, err := getTrustPinning(config)
	if err != nil {
		return err
	}

	nRepo, err := notaryclient.NewNotaryRepository(
		config.GetString("trust_dir"), gun, getRemoteTrustServer(config), rt, t.retriever, trustPin)
	if err != nil {
		return err
	}

	target, err := nRepo.GetTargetByName(targetName)
	if err != nil {
		return err
	}

	cmd.Println(target.Name, fmt.Sprintf("sha256:%x", target.Hashes["sha256"]), target.Length)
	return nil
}

func (t *tufCommander) tufStatus(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		cmd.Usage()
		return fmt.Errorf("Must specify a GUN")
	}

	config, err := t.configGetter()
	if err != nil {
		return err
	}
	gun := args[0]

	trustPin, err := getTrustPinning(config)
	if err != nil {
		return err
	}

	nRepo, err := notaryclient.NewNotaryRepository(
		config.GetString("trust_dir"), gun, getRemoteTrustServer(config), nil, t.retriever, trustPin)
	if err != nil {
		return err
	}

	cl, err := nRepo.GetChangelist()
	if err != nil {
		return err
	}

	if len(cl.List()) == 0 {
		cmd.Printf("No unpublished changes for %s\n", gun)
		return nil
	}

	cmd.Printf("Unpublished changes for %s:\n\n", gun)
	cmd.Printf("%-10s%-10s%-12s%s\n", "action", "scope", "type", "path")
	cmd.Println("----------------------------------------------------")
	for _, ch := range cl.List() {
		cmd.Printf("%-10s%-10s%-12s%s\n", ch.Action(), ch.Scope(), ch.Type(), ch.Path())
	}
	return nil
}

func (t *tufCommander) tufPublish(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		cmd.Usage()
		return fmt.Errorf("Must specify a GUN")
	}

	config, err := t.configGetter()
	if err != nil {
		return err
	}
	gun := args[0]

	cmd.Println("Pushing changes to", gun)

	rt, err := getTransport(config, gun, false)
	if err != nil {
		return err
	}

	trustPin, err := getTrustPinning(config)
	if err != nil {
		return err
	}

	nRepo, err := notaryclient.NewNotaryRepository(
		config.GetString("trust_dir"), gun, getRemoteTrustServer(config), rt, t.retriever, trustPin)
	if err != nil {
		return err
	}

	if err = nRepo.Publish(); err != nil {
		return err
	}
	return nil
}

func (t *tufCommander) tufRemove(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("Must specify a GUN and target")
	}
	config, err := t.configGetter()
	if err != nil {
		return err
	}

	gun := args[0]
	targetName := args[1]

	trustPin, err := getTrustPinning(config)
	if err != nil {
		return err
	}

	// no online operation are performed by remove so the transport argument
	// should be nil.
	repo, err := notaryclient.NewNotaryRepository(
		config.GetString("trust_dir"), gun, getRemoteTrustServer(config), nil, t.retriever, trustPin)
	if err != nil {
		return err
	}
	// If roles is empty, we default to removing from targets
	if err = repo.RemoveTarget(targetName, t.roles...); err != nil {
		return err
	}

	cmd.Printf("Removal of %s from %s staged for next publish.\n", targetName, gun)
	return nil
}

func (t *tufCommander) tufVerify(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		cmd.Usage()
		return fmt.Errorf("Must specify a GUN and target")
	}

	config, err := t.configGetter()
	if err != nil {
		return err
	}

	payload, err := getPayload(t)
	if err != nil {
		return err
	}

	gun := args[0]
	targetName := args[1]

	rt, err := getTransport(config, gun, true)
	if err != nil {
		return err
	}

	trustPin, err := getTrustPinning(config)
	if err != nil {
		return err
	}

	nRepo, err := notaryclient.NewNotaryRepository(
		config.GetString("trust_dir"), gun, getRemoteTrustServer(config), rt, t.retriever, trustPin)
	if err != nil {
		return err
	}

	target, err := nRepo.GetTargetByName(targetName)
	if err != nil {
		return fmt.Errorf("error retrieving target by name:%s, error:%v", targetName, err)
	}

	if err := data.CheckHashes(payload, targetName, target.Hashes); err != nil {
		return fmt.Errorf("data not present in the trusted collection, %v", err)
	}

	return feedback(t, payload)
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

	if term.IsTerminal(0) {
		state, err := term.SaveState(0)
		if err != nil {
			logrus.Errorf("error saving terminal state, cannot retrieve password: %s", err)
			return "", ""
		}
		term.DisableEcho(0, state)
		defer term.RestoreTerminal(0, state)
	}

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
func getTransport(config *viper.Viper, gun string, readOnly bool) (http.RoundTripper, error) {
	// Attempt to get a root CA from the config file. Nil is the host defaults.
	rootCAFile := utils.GetPathRelativeToConfig(config, "remote_server.root_ca")
	clientCert := utils.GetPathRelativeToConfig(config, "remote_server.tls_client_cert")
	clientKey := utils.GetPathRelativeToConfig(config, "remote_server.tls_client_key")

	insecureSkipVerify := false
	if config.IsSet("remote_server.skipTLSVerify") {
		insecureSkipVerify = config.GetBool("remote_server.skipTLSVerify")
	}

	if clientCert == "" && clientKey != "" || clientCert != "" && clientKey == "" {
		return nil, fmt.Errorf("either pass both client key and cert, or neither")
	}

	tlsConfig, err := tlsconfig.Client(tlsconfig.Options{
		CAFile:             rootCAFile,
		InsecureSkipVerify: insecureSkipVerify,
		CertFile:           clientCert,
		KeyFile:            clientKey,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to configure TLS: %s", err.Error())
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
	readOnly bool) (http.RoundTripper, error) {

	// TODO(dmcgowan): add notary specific headers
	authTransport := transport.NewTransport(baseTransport)
	pingClient := &http.Client{
		Transport: authTransport,
		Timeout:   5 * time.Second,
	}
	endpoint, err := url.Parse(trustServerURL)
	if err != nil {
		return nil, fmt.Errorf("Could not parse remote trust server url (%s): %s", trustServerURL, err.Error())
	}
	if endpoint.Scheme == "" {
		return nil, fmt.Errorf("Trust server url has to be in the form of http(s)://URL:PORT. Got: %s", trustServerURL)
	}
	subPath, err := url.Parse("v2/")
	if err != nil {
		return nil, fmt.Errorf("Failed to parse v2 subpath. This error should not have been reached. Please report it as an issue at https://github.com/docker/notary/issues: %s", err.Error())
	}
	endpoint = endpoint.ResolveReference(subPath)
	req, err := http.NewRequest("GET", endpoint.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := pingClient.Do(req)
	if err != nil {
		logrus.Errorf("could not reach %s: %s", trustServerURL, err.Error())
		logrus.Info("continuing in offline mode")
		return nil, nil
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
		return nil, nil
	}

	challengeManager := auth.NewSimpleChallengeManager()
	if err := challengeManager.AddResponse(resp); err != nil {
		return nil, err
	}

	ps := passwordStore{anonymous: readOnly}

	var actions []string
	if readOnly {
		actions = []string{"pull"}
	} else {
		actions = []string{"push", "pull"}
	}
	tokenHandler := auth.NewTokenHandler(authTransport, ps, gun, actions...)
	basicHandler := auth.NewBasicHandler(ps)

	modifier := auth.NewAuthorizer(challengeManager, tokenHandler, basicHandler)

	if !readOnly {
		return newAuthRoundTripper(transport.NewTransport(baseTransport, modifier)), nil
	}

	// Try to authenticate read only repositories using basic username/password authentication
	return newAuthRoundTripper(transport.NewTransport(baseTransport, modifier),
		transport.NewTransport(baseTransport, auth.NewAuthorizer(challengeManager, auth.NewTokenHandler(authTransport, passwordStore{anonymous: false}, gun, actions...)))), nil
}

func getRemoteTrustServer(config *viper.Viper) string {
	if configRemote := config.GetString("remote_server.url"); configRemote != "" {
		return configRemote
	}
	return defaultServerURL
}

func getTrustPinning(config *viper.Viper) (trustpinning.TrustPinConfig, error) {
	var ok bool
	// Need to parse out Certs section from config
	certMap := config.GetStringMap("trust_pinning.certs")
	resultCertMap := make(map[string][]string)
	for gun, certSlice := range certMap {
		var castedCertSlice []interface{}
		if castedCertSlice, ok = certSlice.([]interface{}); !ok {
			return trustpinning.TrustPinConfig{}, fmt.Errorf("invalid format for trust_pinning.certs")
		}
		certsForGun := make([]string, len(castedCertSlice))
		for idx, certIDInterface := range castedCertSlice {
			if certID, ok := certIDInterface.(string); ok {
				certsForGun[idx] = certID
			} else {
				return trustpinning.TrustPinConfig{}, fmt.Errorf("invalid format for trust_pinning.certs")
			}
		}
		resultCertMap[gun] = certsForGun
	}
	return trustpinning.TrustPinConfig{
		DisableTOFU: config.GetBool("trust_pinning.disable_tofu"),
		CA:          config.GetStringMapString("trust_pinning.ca"),
		Certs:       resultCertMap,
	}, nil
}

// authRoundTripper tries to authenticate the requests via multiple HTTP transactions (until first succeed)
type authRoundTripper struct {
	trippers []http.RoundTripper
}

func newAuthRoundTripper(trippers ...http.RoundTripper) http.RoundTripper {
	return &authRoundTripper{trippers: trippers}
}

func (a *authRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {

	var resp *http.Response
	// Try all run all transactions
	for _, t := range a.trippers {
		var err error
		resp, err = t.RoundTrip(req)
		// Reject on error
		if err != nil {
			return resp, err
		}

		// Stop when request is authorized/unknown error
		if resp.StatusCode != http.StatusUnauthorized {
			return resp, nil
		}
	}

	// Return the last response
	return resp, nil
}
