package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/dhe-deploy/ha-integration/ha_utils"
	"github.com/docker/dhe-deploy/manager/versions"
	"github.com/docker/dhe-deploy/shared/dtrutil"
	"github.com/docker/distribution/registry/client/auth"
	"github.com/docker/distribution/registry/client/transport"
	"github.com/docker/notary"
	"github.com/docker/notary/client"
	"github.com/docker/notary/passphrase"
	"github.com/docker/notary/trustpinning"
	"github.com/docker/notary/tuf/data"
	"github.com/stretchr/testify/require"
)

type credentialStore struct {
	username      string
	password      string
	refreshTokens map[string]string
}

func (tcs *credentialStore) Basic(url *url.URL) (string, string) {
	return tcs.username, tcs.password
}

// refresh tokens are the long lived tokens that can be used instead of a password
func (tcs *credentialStore) RefreshToken(u *url.URL, service string) string {
	return tcs.refreshTokens[service]
}

func (tcs *credentialStore) SetRefreshToken(u *url.URL, service string, token string) {
	if tcs.refreshTokens != nil {
		tcs.refreshTokens[service] = token
	}
}

func GetGarantAuthTransport(server string, scopes []string, username, password string) (http.RoundTripper, error) {
	insecureTransport, err := dtrutil.HTTPTransport(true, nil, "", "")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v2/", server), nil)
	if err != nil {
		return nil, err
	}
	pingClient := &http.Client{
		Transport: insecureTransport,
		Timeout:   5 * time.Second,
	}
	resp, err := pingClient.Do(req)
	if err != nil {
		return nil, err
	}
	challengeManager := auth.NewSimpleChallengeManager()
	if err := challengeManager.AddResponse(resp); err != nil {
		return nil, err
	}

	creds := credentialStore{
		username:      username,
		password:      password,
		refreshTokens: make(map[string]string),
	}

	var scopeObjs []auth.Scope
	for _, scopeName := range scopes {
		scopeObjs = append(scopeObjs, auth.RepositoryScope{
			Repository: scopeName,
			Actions:    []string{"push", "pull"},
		})
	}

	// allow setting multiple scopes so we don't have to reauth
	tokenHandler := auth.NewTokenHandlerWithOptions(auth.TokenHandlerOptions{
		Transport:   insecureTransport,
		Credentials: &creds,
		Scopes:      scopeObjs,
	})

	authedTransport := transport.NewTransport(insecureTransport, auth.NewAuthorizer(challengeManager, tokenHandler))
	return authedTransport, nil
}

func requireNotaryPushPullSuccess(confDir, gun, replicaServer, pushTargetName string, additionalTargetsExpected []client.Target, t require.TestingT) client.Target {
	rt, err := GetGarantAuthTransport(
		replicaServer, []string{gun}, ha_utils.GetAdminUser(), ha_utils.GetAdminPassword())
	require.NoError(t, err)

	repo, err := client.NewNotaryRepository(
		confDir, gun, replicaServer, rt, passphrase.ConstantRetriever("password"), trustpinning.TrustPinConfig{})
	require.NoError(t, err)

	fileMeta, err := data.NewFileMeta(strings.NewReader(pushTargetName), notary.SHA256)
	require.NoError(t, err)
	pushTarget := client.Target{
		Name:   pushTargetName,
		Hashes: fileMeta.Hashes,
		Length: fileMeta.Length,
	}

	err = repo.AddTarget(&pushTarget)
	require.NoError(t, err, "Could not add a target to the notary repo")

	err = repo.Publish()
	if _, ok := err.(client.ErrRepoNotInitialized); ok {
		// initialize the repo and try publishing again
		func() {
			rootPublicKey, err := repo.CryptoService.Create(data.CanonicalRootRole, "", data.ECDSAKey)
			require.NoError(t, err)
			err = repo.Initialize([]string{rootPublicKey.ID()})
			require.NoError(t, err)
		}()
		err = repo.Publish()
	}
	require.NoError(t, err, "Could not publish notary repo")

	targets, err := repo.ListTargets()
	require.NoError(t, err, "Could not list notary targets")

	gotTargets := make(map[string]client.Target)
	for _, targetWithRole := range targets {
		require.Equal(t, data.CanonicalTargetsRole, targetWithRole.Role)
		gotTargets[targetWithRole.Target.Name] = targetWithRole.Target
	}

	for _, expectedTarget := range append(additionalTargetsExpected, pushTarget) {
		gotten, ok := gotTargets[expectedTarget.Name]
		require.True(t, ok, "Could not find %s in trust data", expectedTarget.Name)
		require.Equal(t, expectedTarget, gotten, "targets are different")
	}

	return pushTarget
}

func NotaryTest(machines []ha_utils.Machine, imageArgs DTRImageArgs, t require.TestingT) {
	minVersionSupported := v2_1_0
	currVersion, err := versions.TagToSemver(imageArgs.DTRTag)
	require.NoError(t, err)
	// skip if we are on version < 2.1
	if currVersion.LT(minVersionSupported) {
		return
	}
	repoName := "notary-no-images"
	// note, garant doesn't seem to issue tokens for DefaultInstallation.API.GetHost():DefaultInstallation.API.GetApiClientPort()/adminuser/repoName and
	// notary operations 401
	gun := path.Join(ha_utils.GetAdminUser(), repoName)

	// login and create the repo in DTR, else we won't be able to get a garant token for the repo
	require.NoError(t, DefaultInstallation.API.Login(ha_utils.GetAdminUser(), ha_utils.GetAdminPassword()))
	_, err = DefaultInstallation.API.CreateRepository(ha_utils.GetAdminUser(), repoName, "meh", "meh", "public")
	require.NoError(t, err)
	defer DefaultInstallation.API.DeleteRepository(ha_utils.GetAdminUser(), repoName)

	repos, err := DefaultInstallation.API.ListRepositories(ha_utils.GetAdminUser())
	require.NoError(t, err)
	var found bool
	for _, repo := range repos {
		if repo.Name == repoName {
			found = true
		}
	}
	require.True(t, found, "Repo %s was never created", repoName)

	// set up the notary scratch space
	notaryConfDir, err := ioutil.TempDir("", "notary-only-dtr-ha-test")
	require.NoError(t, err)
	defer os.RemoveAll(notaryConfDir)

	targets := make([]client.Target, 0)

	// push from every machine, then assert that pull gets the data from the
	// latest pull and all the other pushes on other machines
	for _, machine := range machines {
		if _, ok := DefaultInstallation.Replicas[machine.GetName()]; ok {
			machineIP, err := machine.GetIP()
			require.NoError(t, err)

			replicaServer := fmt.Sprintf("https://%s", machineIP)
			if DefaultInstallation.Args.ReplicaHTTPSPort != "" {
				replicaServer = replicaServer + ":" + DefaultInstallation.Args.ReplicaHTTPSPort
			}
			pushTarget := requireNotaryPushPullSuccess(notaryConfDir, gun, replicaServer, machine.GetName(), targets, t)
			targets = append(targets, pushTarget)

			// clear out the tuf metadata, but we want to keep the keys so we can keep signing.
			os.RemoveAll(filepath.Join(notaryConfDir, "tuf"))
		}
	}
}
