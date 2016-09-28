package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/orca/auth"
	enziclient "github.com/docker/orca/enzi/api/client"
	"github.com/docker/orca/enzi/api/forms"
	enziconfig "github.com/docker/orca/enzi/config"
	enzischema "github.com/docker/orca/enzi/schema"
	enziworker "github.com/docker/orca/enzi/worker"
	"github.com/samalba/dockerclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// These all assume the standard test rig setup - might want to make configurable

// This verifies that we have the environment variables set for LDAP
func LDAPTestsEnabled(t *testing.T) bool {
	if os.Getenv("LDAP_URL") == "" {
		t.Skip(`Skipping LDAP test - you must set the following
export LDAP_URL="ldap://192.168.1.2"
export LDAP_SEARCH_DN="ucpsearcher"
export LDAP_SEARCH_PASSWORD="secret"
export LDAP_ADMIN="ucpadmin"
export LDAP_ADMIN_PASSWORD="secret"
export LDAP_USER="jdoe"
export LDAP_USER_PASSWORD="secret"
export LDAP_BASE_DN="dc=ad,dc=dckr,dc=org"
export LDAP_USER_SEARCH="(memberOf=CN=dheusers,OU=Groups,DC=ad,DC=dckr,DC=org)"
export LDAP_TEAM_DN="CN=dtradmins,OU=Groups,DC=ad,DC=dckr,DC=org"
export LDAP_TEAM_MEMBER_ATTR="member"
`)
		return false
	}
	return true
}

func SetupLDAPAuth(t *testing.T, client *dockerclient.DockerClient) {
	orcaURL := *client.URL
	orcaURL.Path = ""

	// Create a secondary non-admin user for testing after switching to LDAP
	newUsername := "newuser3"
	newPassword := "supersecret"
	log.Debug("Creating user")
	require.Nil(t, CreateNewUser(nil, orcaURL.String(), GetAdminUser(), GetAdminPassword(), newUsername, newPassword, false, auth.RestrictedControl))

	log.Debugf("Attempting to setup LDAP")
	ldapURL := os.Getenv("LDAP_URL")
	require.True(t, ldapURL != "")

	// Make an eNZi Session.
	enziCreds := &enziclient.BasicAuthenticator{
		Username: GetAdminUser(),
		Password: GetAdminPassword(),
	}
	enziSession := enziclient.New(client.HTTPClient, orcaURL.Host, "enzi", enziCreds)

	// Prepare the LDAP Settings form.
	recoveryAdminPassword := os.Getenv("LDAP_ADMIN_PASSWORD")
	usernameAttr := os.Getenv("LDAP_USERNAME_ATTR")
	if usernameAttr == "" {
		// Fallback to the default for Active Directory.
		usernameAttr = "sAMAccountName"
	}
	ldapConfigForm := forms.LDAPSettings{
		RecoveryAdminUsername: os.Getenv("LDAP_ADMIN"),
		RecoveryAdminPassword: &recoveryAdminPassword,
		ServerURL:             ldapURL,
		StartTLS:              true,
		TLSSkipVerify:         true,
		ReaderDN:              os.Getenv("LDAP_SEARCH_DN"),
		ReaderPassword:        os.Getenv("LDAP_SEARCH_PASSWORD"),
		UserSearchConfigs: []forms.UserSearchOpts{{
			BaseDN: os.Getenv("LDAP_BASE_DN"),
			// WholeSubtree handles when the search base DN is more
			// than one link away from the users. We're not sure
			// what is expected given the environment variables.
			ScopeSubtree: true,
			UsernameAttr: usernameAttr,
			Filter:       os.Getenv("LDAP_USER_SEARCH"),
		}},
		AdminSyncOpts: forms.MemberSyncOpts{
			// Disable admin syncing. This is the default, but we
			// want to be explicit in this test.
			EnableSync: false,
		},
		SyncSchedule: "@hourly",
	}

	// Make the request to set LDAP settings in eNZi.
	_, err := enziSession.SetLDAPSettings(ldapConfigForm)
	require.Nil(t, err)

	// Make the request to set the user default role in UCP. We need to
	// update it from the current enziConfig.
	orcaURL.Path = "/api/config/auth2"
	resp, err := client.HTTPClient.Get(orcaURL.String())
	require.Nil(t, err)

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		require.FailNow(t, string(body))
	}
	require.Nil(t, err)

	var ucpAuthConfig auth.AuthenticatorConfiguration
	require.Nil(t, json.Unmarshal(body, &ucpAuthConfig))

	// Set the new user default role to "Restricted Control".
	ucpAuthConfig.EnziConfig.UserDefaultRole = auth.RestrictedControl
	ucpAuthConfigJSON, err := json.Marshal(ucpAuthConfig)
	require.Nil(t, err)

	resp, err = client.HTTPClient.Post(orcaURL.String(), "application/json", bytes.NewBuffer(ucpAuthConfigJSON))
	require.Nil(t, err)

	body, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		require.FailNow(t, string(body))
	}
	require.Nil(t, err)

	// That doesn't actually switch eNZi to use LDAP though, there's a
	// separate endpoint to actually perform the switch.
	_, err = enziSession.SetAuthConfig(forms.AuthConfig{Backend: enziconfig.AuthBackendLDAP})
	require.Nil(t, err)

	log.Debugf("Succesfully configured LDAP") // Hooray!

	// Next, make a single team which is to be synced with an LDAP group.
	ldapOpts := TeamLDAPOpts{
		GroupDN:    os.Getenv("LDAP_TEAM_DN"),
		MemberAttr: os.Getenv("LDAP_TEAM_MEMBER_ATTR"),
	}
	if ldapOpts.MemberAttr == "" {
		// Fallback to the default value.
		ldapOpts.MemberAttr = "member"
	}

	teamID, err := AddTeam(nil, orcaURL.String(), os.Getenv("LDAP_ADMIN"), os.Getenv("LDAP_ADMIN_PASSWORD"), "testldapgroup", ldapOpts)
	require.Nil(t, err)

	// The sync is scheduled to run hourly, but we want to trigger a sync
	// now.
	runSync(t, client.HTTPClient, orcaURL.Host)

	orcaURL.Path = ""

	log.Debugf("Verifying LDAP recovery admin can login")
	client, err = GetUserDockerClient(orcaURL.String(), os.Getenv("LDAP_ADMIN"), os.Getenv("LDAP_ADMIN_PASSWORD"))
	require.Nil(t, err)

	// The recovery admin account should work with certs.
	log.Debugf("Verifying recovery admin client cert connection works")
	version, err := client.Version()
	require.Nil(t, err)
	require.True(t, strings.Contains(version.Version, "ucp"))

	// And double check to make sure we can't login with an invalid password
	log.Debugf("Verifying builtin auth admin can't login with bad password ")
	_, err = GetUserDockerClient(orcaURL.String(), GetAdminUser(), "invalid"+GetAdminPassword())
	require.NotNil(t, err)
	// An empty password is a very important case as well.
	_, err = GetUserDockerClient(orcaURL.String(), GetAdminUser(), "")
	require.NotNil(t, err)

	// Verify non-admin account can't login
	log.Debugf("Verifying builtin non-admin can no longer login")
	_, err = GetUserDockerClient(orcaURL.String(), newUsername, newPassword)
	require.NotNil(t, err)

	// And double check to make sure we can't login using the regular admin
	log.Debugf("Verifying builtin auth admin can't login ")
	orcaURL.Path = ""
	_, err = GetUserDockerClient(orcaURL.String(), GetAdminUser(), GetAdminPassword())
	require.NotNil(t, err)

	// Try adding an access list for the LDAP-synced team.
	require.Nil(t, AddAccessList(nil, orcaURL.String(), os.Getenv("LDAP_ADMIN"), os.Getenv("LDAP_ADMIN_PASSWORD"), teamID, "busybox", 2))

	// Finally, make sure that the LDAP user shows up since the sync has
	// finished.
	log.Infof("Ensuring user %s appears from the background sync", os.Getenv("LDAP_USER"))
	userURL := orcaURL
	userURL.Path = "/api/accounts/" + os.Getenv("LDAP_USER")

	resp, err = client.HTTPClient.Get(userURL.String())
	require.Nil(t, err)

	body, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		require.FailNow(t, string(body))
	}
	require.Nil(t, err)

	log.Info("Detected the account")

	// TODO Might want to validate the user account one last time, but we have other tests that
	// will pop if it's missing, so just fall through at this point.

	// Make sure that the user is a member of the synced team.
	members, err := GetTeamMember(nil, orcaURL.String(), os.Getenv("LDAP_ADMIN"), os.Getenv("LDAP_ADMIN_PASSWORD"), teamID)
	require.Nil(t, err)

	foundUser := false
	for _, username := range members {
		if username == os.Getenv("LDAP_USER") {
			foundUser = true
			break
		}
	}

	require.True(t, foundUser)
}

// Pre 1.1 setup routine (used for upgrade scenarios)
func LegacySetupLDAPAuth(serverURL string) error {
	log.Debug("In LegacySetupLDAPAuth")
	tlsConfig, err := GetUserTLSConfig(serverURL, GetAdminUser(), GetAdminPassword())
	if err != nil {
		log.Debug("GetUserTLSConfig failed")
		return err
	}
	tlsConfig.InsecureSkipVerify = true

	tr := &http.Transport{
		TLSClientConfig:     tlsConfig,
		DisableKeepAlives:   true,
		MaxIdleConnsPerHost: 0,
	}

	client := &http.Client{
		Transport: tr,
	}

	orcaURL, err := url.Parse(serverURL)
	if err != nil {
		return err
	}
	log.Debugf("Attempting to setup LDAP")
	ldapURL := os.Getenv("LDAP_URL")
	if ldapURL == "" {
		return fmt.Errorf("You must set LDAP_URL and friends to enable this testcase")
	}

	orcaURL.Path = "/api/config/auth"

	// WARNING - this has assumptions in it - might want to genralize if we test against multiple AD/LDAP rigs...
	ldapUsernameAttr := os.Getenv("LDAP_USERNAME_ATTR")
	if ldapUsernameAttr == "" {
		ldapUsernameAttr = "sAMAccountName"
	}
	authConfig := auth.AuthenticatorConfiguration{
		AuthenticatorType: auth.AuthenticatorLDAP,
		LDAPConfig: auth.LDAPSettings{
			ServerURL:         ldapURL,
			AdminUsername:     os.Getenv("LDAP_ADMIN"),
			AdminPassword:     os.Getenv("LDAP_ADMIN_PASSWORD"),
			StartTLS:          true,
			TLSSkipVerify:     true,
			ReaderDN:          os.Getenv("LDAP_SEARCH_DN"),
			ReaderPassword:    os.Getenv("LDAP_SEARCH_PASSWORD"),
			UserBaseDN:        os.Getenv("LDAP_BASE_DN"),
			UserLoginAttrName: ldapUsernameAttr,
			UserSearchFilter:  os.Getenv("LDAP_USER_SEARCH"),
			SyncInterval:      60,
			UserDefaultRole:   2,
		},
	}
	data, err := json.Marshal(authConfig)
	if err != nil {
		return err
	}

	//Periodically dump out the sync message in case the test is having problems
	done := false
	go func() {
		linesShown := 0
		statusURL := orcaURL
		statusURL.Path = "/api/authsync"
		for !done {
			time.Sleep(10 * time.Second)
			resp, err := client.Get(statusURL.String())
			if err != nil {
				log.Errorf("Failed to get status: %s", err)
				continue
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Errorf("Failed to read body: %s", err)
				continue
			}
			lines := strings.Split(string(body), "\n")
			for ; linesShown < len(lines); linesShown++ {
				// XXX Warning - fragile - we need a more programatic way to do this!
				if strings.Contains(lines[linesShown], "Sync completed") {
					done = true
				}
				log.Infof("Sync status: %d %s", linesShown, lines[linesShown])
			}
		}
		log.Info("Exiting sync monitor")
	}()
	defer func() { /*Give up showing sync messages if the test exits */ done = true }()

	// If you're having problems, uncomment this to get a dump of the settings..
	// log.Infof("LDAP Settings: %s", string(data))

	req, err := http.NewRequest("POST", orcaURL.String(), bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Debugf("LDAP config request was failed")
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	resp.Body.Close()
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		log.Debugf("LDAP config request was rejected by the server")
		return fmt.Errorf(string(body))
	}
	log.Debugf("Succesfully configured LDAP")

	// Give the background sync a little time to get going, and possibly finish if the user list is short
	time.Sleep(5 * time.Second)

	log.Debugf("Verifying LDAP admin can login ")
	_, err = GetUserDockerClient(orcaURL.String(), os.Getenv("LDAP_ADMIN"), os.Getenv("LDAP_ADMIN_PASSWORD"))
	if err != nil {
		return err
	}

	// Now spin for a while until the LDAP user shows up or the sync finishes
	log.Infof("Spinning for a while until user %s appears from the background sync", os.Getenv("LDAP_USER"))
	userURL := orcaURL
	userURL.Path = "/api/accounts/" + os.Getenv("LDAP_USER")
	for !done {
		time.Sleep(10 * time.Second)
		resp, err := client.Get(userURL.String())
		if err != nil {
			//log.Debugf("XXX Failed to get user: %s", err)
			continue
		}
		if resp.StatusCode > 299 || resp.StatusCode < 200 {
			//log.Debugf("XXX status: %d", resp.StatusCode)
			continue
		}
		// Got a 2xx response, break out
		log.Info("Detected the account")
		break
	}
	log.Info("LDAP was configured")
	return nil
}

func OrcaSuiteEnableLDAPLegacy(s *OrcaTestSuite) error {
	serverURLs := s.GetServerURLs()
	log.Info("Attempting to set up LDAP using legacy API")
	err := LegacySetupLDAPAuth(serverURLs[0])
	if err != nil {
		return err
	}

	// SUPER HACK!  LDAP in HA is broken in 1.0 and requires restarting the controllers
	if err := s.RestartAllControllers(); err != nil {
		return err
	}
	s.IsLDAP = true

	// This is a little hokey...
	// Overwrite the admin in the env so GetAdmin* will work
	os.Setenv("UCP_ADMIN", os.Getenv("LDAP_ADMIN"))
	os.Setenv("UCP_ADMIN_PASSWORD", os.Getenv("LDAP_ADMIN_PASSWORD"))

	log.Info("Make sure all the servers are happy with LDAP...")
	results := make(chan error, len(serverURLs))
	var wg sync.WaitGroup
	for i := 0; i < len(serverURLs); i++ {
		wg.Add(1)
		go func(serverURL string) {
			defer wg.Done()
			var err error
			for i := 0; i < 60; i++ {
				log.Debugf("Getting admin client from %s", serverURL)
				_, err = GetUserDockerClient(serverURL, GetAdminUser(), GetAdminPassword())
				if err == nil {
					log.Infof("Server %s looks properly configured", serverURL)
					break
				}
				log.Debugf("Error: %s %s - retrying...", serverURL, err)
				time.Sleep(1 * time.Second)
			}
			results <- err
		}(serverURLs[i])
	}
	wg.Wait()
	close(results)
	for err := range results {
		if err != nil {
			return err
		}
	}

	return err
}

func runSync(t *testing.T, httpClient *http.Client, authProviderAddr string) {
	// Make a new session using credentials of the LDAP recovery admin
	// user.
	enziCreds := &enziclient.BasicAuthenticator{
		Username: os.Getenv("LDAP_ADMIN"),
		Password: os.Getenv("LDAP_ADMIN_PASSWORD"),
	}
	enziSession := enziclient.New(httpClient, authProviderAddr, "enzi", enziCreds)

	log.Debug("Triggering LDAP Sync Job")
	jobResp, err := enziSession.CreateJob(forms.JobSubmission{Action: enziworker.ActionLdapSync})
	require.Nil(t, err)
	require.NotEmpty(t, jobResp.ID)

	// Wait for the sync job to finish.
	for {
		log.Debug("Waiting for LDAP Sync Job ...")
		// Give it a little more time to run. If the job is still in
		// "waiting" status, this should give it enough time to get
		// picked up by a worker.
		time.Sleep(5 * time.Second)

		jobResp, err = enziSession.GetJob(jobResp.ID)
		require.Nil(t, err)

		if jobResp.Status == enzischema.JobStatusDone {
			break
		}

		jobLogsReadCloser, err := enziSession.GetJobLogs(jobResp.ID)
		require.Nil(t, err)

		jobLogs, err := ioutil.ReadAll(jobLogsReadCloser)
		require.Nil(t, err)

		require.Equal(t, enzischema.JobStatusRunning, jobResp.Status, string(jobLogs))
		jobLogsReadCloser.Close()
	}

	log.Debug("LDAP Sync Job Completed")
}

func SetupBasicAuth(t *testing.T, client *dockerclient.DockerClient) {
	log.Debugf("Attempting to setup builtin auth")

	orcaURL := *client.URL

	// Make a new session using credentials of the LDAP recovery admin
	// user.
	enziCreds := &enziclient.BasicAuthenticator{
		Username: os.Getenv("LDAP_ADMIN"),
		Password: os.Getenv("LDAP_ADMIN_PASSWORD"),
	}
	enziSession := enziclient.New(client.HTTPClient, orcaURL.Host, "enzi", enziCreds)

	_, err := enziSession.SetAuthConfig(forms.AuthConfig{Backend: enziconfig.AuthBackendManaged})
	require.Nil(t, err)

	log.Debugf("Succesfully configured builtin auth")
}

// Reusable test cases

// This test configures ldap, verifies basic sanity, then reverts back to basic auth for idempotency
func TestLDAPBasic(t *testing.T, serverURL string) {
	if !LDAPTestsEnabled(t) {
		return
	}

	client, err := GetUserDockerClient(serverURL, GetAdminUser(), GetAdminPassword())
	require.Nil(t, err)
	log.Info("Enabling LDAP")
	SetupLDAPAuth(t, client)
	log.Info("Getting user client")
	userClient, err := GetUserDockerClient(serverURL, os.Getenv("LDAP_USER"), os.Getenv("LDAP_USER_PASSWORD"))
	require.Nil(t, err)
	version, err := userClient.Version()
	require.Nil(t, err)
	require.True(t, strings.Contains(version.Version, "ucp"))

	// Some negative tests
	log.Info("Verify bad password scenarios")
	_, err = GetUserDockerClient(serverURL, os.Getenv("LDAP_USER"), "invalidpassword")
	require.NotNil(t, err)
	_, err = GetUserDockerClient(serverURL, os.Getenv("LDAP_USER"), "")
	require.NotNil(t, err)

	// Create a managed group, and verify membership looks sane
	teamID, err := AddTeam(nil, serverURL, os.Getenv("LDAP_ADMIN"), os.Getenv("LDAP_ADMIN_PASSWORD"), "testmanagedgroup")
	require.Nil(t, err)
	require.Nil(t, AddTeamMember(nil, serverURL, os.Getenv("LDAP_ADMIN"), os.Getenv("LDAP_ADMIN_PASSWORD"), teamID, os.Getenv("LDAP_USER")))
	members, err := GetTeamMember(nil, serverURL, os.Getenv("LDAP_ADMIN"), os.Getenv("LDAP_ADMIN_PASSWORD"), teamID)
	assert.Nil(t, err)                                         // Let test keep running so we can unwind
	assert.Equal(t, members, []string{os.Getenv("LDAP_USER")}) // Let test keep running so we can unwind

	// At this point the builtin admin user is unusable
	client, err = GetUserDockerClient(serverURL, os.Getenv("LDAP_ADMIN"), os.Getenv("LDAP_ADMIN_PASSWORD"))
	require.Nil(t, err)
	log.Info("Reverting to basic auth")
	SetupBasicAuth(t, client)

	log.Info("Verifying LDAP account no longer works")
	userClient, err = GetUserDockerClient(serverURL, os.Getenv("LDAP_USER"), os.Getenv("LDAP_USER_PASSWORD"))
	require.NotNil(t, err)
	log.Info("Looks good!")

}

// Must have LDAP already configured for this - verify non-admin client can connect
func TestLDAPBasicClientConnect(t *testing.T, serverURL string) {
	userClient, err := GetUserDockerClient(serverURL, os.Getenv("LDAP_USER"), os.Getenv("LDAP_USER_PASSWORD"))
	require.Nil(t, err)
	version, err := userClient.Version()
	require.Nil(t, err)
	require.True(t, strings.Contains(version.Version, "ucp"))

	// Some negative tests
	log.Info("Verify bad password scenarios")
	_, err = GetUserDockerClient(serverURL, os.Getenv("LDAP_USER"), "invalidpassword")
	require.NotNil(t, err)
	_, err = GetUserDockerClient(serverURL, os.Getenv("LDAP_USER"), "")
	require.NotNil(t, err)
}

// Toggle a users role, sync, and make sure the role doesn't change
func TestUserRoleDoesntChangeOnSync(t *testing.T, serverURL string) {

	// Bump the user up to Full control
	require.Nil(t, SetUserRole(nil, serverURL, os.Getenv("LDAP_ADMIN"), os.Getenv("LDAP_ADMIN_PASSWORD"),
		os.Getenv("LDAP_USER"), auth.FullControl))

	client, err := GetUserDockerClient(serverURL, os.Getenv("LDAP_ADMIN"), os.Getenv("LDAP_ADMIN_PASSWORD"))
	require.Nil(t, err)
	orcaURL := *client.URL

	// Trigger a sync
	runSync(t, client.HTTPClient, orcaURL.Host)

	log.Infof("Ensuring user %s role did not change after sync", os.Getenv("LDAP_USER"))

	userURL := orcaURL
	userURL.Path = "/api/accounts/" + os.Getenv("LDAP_USER")

	account, err := GetUser(nil, serverURL, os.Getenv("LDAP_ADMIN"), os.Getenv("LDAP_ADMIN_PASSWORD"), os.Getenv("LDAP_USER"))
	require.Nil(t, err)
	require.Equal(t, auth.FullControl, account.Role) // If we ever drop full control, fail

	log.Info("User kept their increased access during the sync")
}

/*
TODO - test case ideas:

* Enable LDAP, delete admin account, re-enable builtin, verify admin can login with default password
* change admin password, Enable LDAP, re-enable builtin, verify admin can login with non-standard password

*/
