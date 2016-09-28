package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	neturl "net/url"
	"os"
	"strings"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/orca/auth"
	"github.com/samalba/dockerclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func GetAdminPassword() string {
	password := os.Getenv("UCP_ADMIN_PASSWORD")
	if password == "" {
		password = "orca"
	}
	return password
}

func GetAdminUser() string {
	username := os.Getenv("UCP_ADMIN")
	if username == "" {
		username = "admin"
	}
	return username
}

func CreateNewUser(client *http.Client, serverURL, adminUser, adminPassword, newUsername, newPassword string, isAdmin bool, role auth.Role) error {
	account := &auth.Account{
		FirstName: newUsername,
		LastName:  newUsername,
		Username:  newUsername,
		Password:  newPassword,
		Admin:     isAdmin,
		Role:      role,
	}
	return SaveUser(client, serverURL, adminUser, adminPassword, account)
}

// TODO belongs elsewhere
func GetVersionWithRetries(client *dockerclient.DockerClient, retries int) (*dockerclient.Version, error) {
	var lastError error
	for i := 0; i < retries; i++ {
		version, err := client.Version()
		if err == nil {
			return version, err
		}
		log.Debugf("Failure for get version, retrying... %s", err)
		lastError = err
		time.Sleep(1 * time.Second)
	}
	return nil, lastError
}

func DisableUser(client *http.Client, serverURL, adminUser, adminPassword, username string) error {
	account, err := GetUser(client, serverURL, adminUser, adminPassword, username)
	if err != nil {
		return err
	}
	account.Disabled = true
	return SaveUser(client, serverURL, adminUser, adminPassword, account)
}

func SetUserRole(client *http.Client, serverURL, adminUser, adminPassword, username string, newRole auth.Role) error {
	account, err := GetUser(client, serverURL, adminUser, adminPassword, username)
	if err != nil {
		return err
	}
	account.Role = newRole
	return SaveUser(client, serverURL, adminUser, adminPassword, account)
}

func GetUser(client *http.Client, serverURL, adminUser, adminPassword, username string) (*auth.Account, error) {
	orcaURL, err := neturl.Parse(serverURL)
	if err != nil {
		return nil, err
	}

	if client == nil {
		client = &http.Client{
			Transport: &http.Transport{
				// Sloppy for testing only - don't copy this into production code!
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
			Timeout: 10 * time.Second, // Be pretty aggressive on timeouts for testing
		}
	}

	// Login and get a token for the admin user
	token, err := GetOrcaToken(client, serverURL, adminUser, adminPassword)
	if err != nil {
		return nil, err
	}

	orcaURL.Path = "/api/accounts/" + username

	var account auth.Account

	req, err := http.NewRequest("GET", orcaURL.String(), nil)
	if err != nil {
		// Should never fail
		return nil, err
	}
	req.Header.Set(GetTokenHeader(token))
	resp, err := client.Do(req)
	if err != nil {
		log.Debug("Failed to make request")
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		if err == nil {
			return nil, fmt.Errorf("Failed to fetch user (%d): %s", resp.StatusCode, string(body))
		} else {
			return nil, fmt.Errorf("Failed to fetch user (%d): %s - %s", resp.StatusCode, string(body), err)
		}
	}
	if err := json.Unmarshal(body, &account); err != nil {
		log.Debug("Failed to unmarshal")
		return nil, err
	}
	return &account, nil
}

func DeleteUser(client *http.Client, serverURL, adminUser, adminPassword, username string) error {
	orcaURL, err := neturl.Parse(serverURL)
	if err != nil {
		return err
	}

	if client == nil {
		client = &http.Client{
			Transport: &http.Transport{
				// Sloppy for testing only - don't copy this into production code!
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
			Timeout: 10 * time.Second, // Be pretty aggressive on timeouts for testing
		}
	}

	// Login and get a token for the admin user
	token, err := GetOrcaToken(client, serverURL, adminUser, adminPassword)
	if err != nil {
		return err
	}

	orcaURL.Path = "/api/accounts/" + username

	req, err := http.NewRequest("DELETE", orcaURL.String(), nil)
	if err != nil {
		// Should never fail
		return err
	}
	req.Header.Set(GetTokenHeader(token))
	resp, err := client.Do(req)
	if err != nil {
		log.Debug("Failed to make request")
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		if err == nil {
			return fmt.Errorf("Failed to fetch user (%d): %s", resp.StatusCode, string(body))
		} else {
			return fmt.Errorf("Failed to fetch user (%d): %s - %s", resp.StatusCode, string(body), err)
		}
	}
	return nil
}

func GetAllUsers(client *http.Client, serverURL, adminUser, adminPassword string) ([]auth.Account, error) {
	orcaURL, err := neturl.Parse(serverURL)
	if err != nil {
		return nil, err
	}

	if client == nil {
		client = &http.Client{
			Transport: &http.Transport{
				// Sloppy for testing only - don't copy this into production code!
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
			Timeout: 10 * time.Second, // Be pretty aggressive on timeouts for testing
		}
	}

	// Login and get a token for the admin user
	token, err := GetOrcaToken(client, serverURL, adminUser, adminPassword)
	if err != nil {
		return nil, err
	}

	orcaURL.Path = "/api/accounts"

	var accounts []auth.Account

	req, err := http.NewRequest("GET", orcaURL.String(), nil)
	if err != nil {
		// Should never fail
		return nil, err
	}
	req.Header.Set(GetTokenHeader(token))
	resp, err := client.Do(req)
	if err != nil {
		log.Debug("Failed to make request")
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		if err == nil {
			return nil, fmt.Errorf("Failed to fetch user (%d): %s", resp.StatusCode, string(body))
		} else {
			return nil, fmt.Errorf("Failed to fetch user (%d): %s - %s", resp.StatusCode, string(body), err)
		}
	}
	if err := json.Unmarshal(body, &accounts); err != nil {
		log.Debug("Failed to unmarshal")
		return nil, err
	}
	return accounts, nil
}

func SaveUser(client *http.Client, serverURL, adminUser, adminPassword string, account *auth.Account) error {
	orcaURL, err := neturl.Parse(serverURL)
	if err != nil {
		return err
	}

	if client == nil {
		client = &http.Client{
			Transport: &http.Transport{
				// Sloppy for testing only - don't copy this into production code!
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
			Timeout: 10 * time.Second, // Be pretty aggressive on timeouts for testing
		}
	}

	// Login and get a token for the admin user
	token, err := GetOrcaToken(client, serverURL, adminUser, adminPassword)
	if err != nil {
		return err
	}

	orcaURL.Path = "/api/accounts"

	reqJson, err := json.Marshal(account)
	if err != nil {
		log.Debug("Failed to generate json")
		return err
	}

	req, err := http.NewRequest("POST", orcaURL.String(), bytes.NewBuffer(reqJson))
	if err != nil {
		// Should never fail
		return err
	}
	req.Header.Set(GetTokenHeader(token))
	resp, err := client.Do(req)
	if err != nil {
		log.Debug("Failed to make request")
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		if err == nil {
			return fmt.Errorf("Failed to update user (%d): %s", resp.StatusCode, string(body))
		} else {
			return err
		}
	}
	log.Debugf("Succesfully updated account %s", account.Username)
	return nil
}

// The following routines are test routines that can be used in different contexts

func TestAddUserEmptyPass(t *testing.T, serverURL string) {
	newUsername := "newemptyuser"
	newPassword := ""
	log.Debug("Creating user with empty password")
	require.NotNil(t, CreateNewUser(nil, serverURL, GetAdminUser(), GetAdminPassword(), newUsername, newPassword, false, auth.RestrictedControl))
}

func TestAddUser(t *testing.T, serverURL string) {
	newUsername := "newuser1"
	newPassword := "s"
	log.Debug("Creating user")
	require.Nil(t, CreateNewUser(nil, serverURL, GetAdminUser(), GetAdminPassword(), newUsername, newPassword, false, auth.RestrictedControl))

	// Now try to get a client as that user, and verify the system is accessible
	log.Debug("getting docker client as new user")
	client, err := GetUserDockerClient(serverURL, newUsername, newPassword)
	require.Nil(t, err)
	log.Debug("getting version")
	version, err := client.Version()
	require.Nil(t, err)
	require.True(t, strings.Contains(version.Version, "ucp"))
	log.Debug("PASS: User created, and can access the system")
}

func TestNonAdminUserNoSwarmForYou(t *testing.T, serverURL []string) {
	newUsername := "newrwuser1"
	newPassword := "supersecret"
	log.Debug("Creating user")
	require.Nil(t, CreateNewUser(nil, serverURL[0], GetAdminUser(), GetAdminPassword(), newUsername, newPassword, false, auth.RestrictedControl))

	// Now try to get a client as that user, and verify the system is accessible
	// Do the permutations against all the URLs
	for _, bundleServerURL := range serverURL {
		log.Debugf("getting docker client from %s as new user", bundleServerURL)
		client, err := GetUserDockerClient(bundleServerURL, newUsername, newPassword)
		require.Nil(t, err)
		log.Debug("getting version from ucp")
		version, err := GetVersionWithRetries(client, 5)
		require.Nil(t, err)
		require.True(t, strings.Contains(version.Version, "ucp"))

		for _, accessServerURL := range serverURL {
			swarmURL, err := neturl.Parse(accessServerURL)
			require.Nil(t, err)
			log.Debugf("Accessing swarm via %s", swarmURL.String())
			// Now try to point to Swarm with the same stuff (TODO - figure this out from the /info results)
			swarmURL.Host = strings.Split(swarmURL.Host, ":")[0] + ":3376" // WARNING - fragile assumption
			swarmClient, err := dockerclient.NewDockerClient(swarmURL.String(), client.TLSConfig)
			if err != nil {
				log.Debug("Failed to connect to swarm - probably OK")
				continue
			}
			_, err = swarmClient.Info()
			require.NotNil(t, err)
			log.Debugf("Swarm connection was rejected as expected: %s", err)
		}
	}
}

func TestNonAdminUserNoProxyForYou(t *testing.T, serverURL []string) {
	newUsername := "newnonadminuser1"
	newPassword := "supersecret"
	log.Debug("Creating user")
	require.Nil(t, CreateNewUser(nil, serverURL[0], GetAdminUser(), GetAdminPassword(), newUsername, newPassword, false, auth.RestrictedControl))

	// Now try to get a client as that user, and verify the system is accessible
	// Do the permutations against all the URLs
	for _, bundleServerURL := range serverURL {
		log.Debugf("getting docker client from %s as new user", bundleServerURL)
		client, err := GetUserDockerClient(bundleServerURL, newUsername, newPassword)
		require.Nil(t, err)
		log.Debug("getting version from ucp")
		version, err := GetVersionWithRetries(client, 5)
		require.Nil(t, err)
		require.True(t, strings.Contains(version.Version, "ucp"))
		for _, accessServerURL := range serverURL {
			proxyURL, err := neturl.Parse(accessServerURL)
			require.Nil(t, err)

			// Now try to point to proxy with the same stuff (TODO - figure this out from the /info results)
			proxyURL.Host = strings.Split(proxyURL.Host, ":")[0] + ":12376" // WARNING - fragile assumption
			log.Debugf("Accessing proxy via %s", proxyURL.String())
			proxyClient, err := dockerclient.NewDockerClient(proxyURL.String(), client.TLSConfig)
			if err != nil {
				log.Debug("Failed to connect to proxy - probably OK")
				continue
			}
			_, err = proxyClient.Info()
			require.NotNil(t, err)
			log.Debugf("Proxy connection was rejected as expected: %s", err)
		}
	}
}

func TestAdminHasSwarm(t *testing.T, serverURL []string) {
	newUsername := "newadminuser1"
	newPassword := "supersecret"
	log.Debug("Creating user")
	require.Nil(t, CreateNewUser(nil, serverURL[0], GetAdminUser(), GetAdminPassword(), newUsername, newPassword, true, auth.RestrictedControl))

	// Now try to get a client as that user, and verify the system is accessible
	// Do the permutations against all the URLs
	for _, bundleServerURL := range serverURL {
		log.Debugf("getting docker client from %s as new user", bundleServerURL)
		client, err := GetUserDockerClient(bundleServerURL, newUsername, newPassword)
		require.Nil(t, err)
		log.Debug("getting version")
		version, err := GetVersionWithRetries(client, 5)
		require.Nil(t, err)
		require.True(t, strings.Contains(version.Version, "ucp"))
		for _, accessServerURL := range serverURL {

			// Now try to point to Swarm with the same stuff (TODO - figure this out from the /info results)
			swarmURL, err := neturl.Parse(accessServerURL)
			require.Nil(t, err)
			swarmURL.Host = strings.Split(swarmURL.Host, ":")[0] + ":3376" // WARNING - fragile assumption
			log.Debugf("Accessing swarm via %s", swarmURL.String())
			swarmClient, err := dockerclient.NewDockerClient(swarmURL.String(), client.TLSConfig)
			require.Nil(t, err)
			version, err = swarmClient.Version()
			require.Nil(t, err)
			require.True(t, strings.Contains(version.Version, "swarm"))
		}
	}
}

func TestDisableUser(t *testing.T, serverURL string) {
	newUsername := "newuser2"
	newPassword := "supersecret"
	log.Debug("Creating user")
	require.Nil(t, CreateNewUser(nil, serverURL, GetAdminUser(), GetAdminPassword(), newUsername, newPassword, false, auth.RestrictedControl))

	// Now try to get a client as that user, and verify the system is accessible
	log.Debug("getting docker client as new user")
	client, err := GetUserDockerClient(serverURL, newUsername, newPassword)
	require.Nil(t, err)
	log.Debug("getting version")
	version, err := client.Version()
	require.Nil(t, err)
	require.True(t, strings.Contains(version.Version, "ucp"))

	log.Debug("Disabling user")
	// Now disable the user
	require.Nil(t, DisableUser(nil, serverURL, GetAdminUser(), GetAdminPassword(), newUsername))

	// Try to access it with the existing client
	_, err = client.Version()
	require.NotNil(t, err)

	// Make sure we can't get a client (does a login)
	_, err = GetUserDockerClient(serverURL, newUsername, newPassword)
	require.NotNil(t, err)
}

// Try various escalation attacks and make sure they're blocked
func TestEscalationsBlocked(t *testing.T, serverURL string) {
	newUsername := "newuser3"
	newPassword := "supersecret"
	log.Debug("Creating user")
	client := &http.Client{
		Transport: &http.Transport{
			// Sloppy for testing only - don't copy this into production code!
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 10 * time.Second, // Be pretty aggressive on timeouts for testing
	}

	// Start out non-admin
	require.Nil(t, CreateNewUser(client, serverURL, GetAdminUser(), GetAdminPassword(), newUsername, newPassword, false, auth.RestrictedControl))

	// Now try to get a client as that user, and verify the system is accessible
	log.Debug("getting docker client as new user")
	dclient, err := GetUserDockerClient(serverURL, newUsername, newPassword)
	require.Nil(t, err)
	log.Debug("getting version")
	version, err := dclient.Version()
	require.Nil(t, err)
	require.True(t, strings.Contains(version.Version, "ucp"))

	log.Debug("Trying to get admin for ourself without being admin")
	account, err := GetUser(client, serverURL, newUsername, newPassword, newUsername)
	require.Nil(t, err)
	require.Equal(t, account.Password, "")
	account.Admin = true
	err = SaveUser(client, serverURL, newUsername, newPassword, account)
	log.Debugf("Attempt to escalate got %s result", err) // We don't really care, as long as it didn't work
	// Verifying we didn't actually get admin
	account, err = GetUser(client, serverURL, newUsername, newPassword, newUsername)
	require.Nil(t, err)
	require.Equal(t, account.Password, "")
	require.False(t, account.Admin)
	log.Debug("Couldn't escalate ourself to admin - good!")

	log.Debug("Trying to increase our role")
	account.Role = auth.FullControl
	err = SaveUser(client, serverURL, newUsername, newPassword, account)
	log.Debugf("Attempt to escalate got %s result", err) // We don't really care, as long as it didn't work
	// Verifying we didn't actually get admin
	account, err = GetUser(client, serverURL, newUsername, newPassword, newUsername)
	require.Nil(t, err)
	require.Equal(t, account.Password, "")
	require.True(t, account.Role <= auth.RestrictedControl)
	log.Debug("Couldn't escalate ourself to FullControl role - good!")

	log.Debug("Trying to add ourself to a team")

	teamname := "testteam1"
	_, err = AddTeam(client, serverURL, GetAdminUser(), GetAdminPassword(), teamname)
	require.Nil(t, err)
	account.ManagedTeams = []string{teamname}
	err = SaveUser(client, serverURL, newUsername, newPassword, account)
	log.Debugf("Attempt to escalate by adding teams %s result", err) // We don't really care, as long as it didn't work
	// Verifying we didn't actually get admin
	account, err = GetUser(client, serverURL, newUsername, newPassword, newUsername)
	require.Nil(t, err)
	require.True(t, len(account.ManagedTeams) == 0)
	log.Debug("Couldn't escalate by adding teams")

	// Now do some other changes related to actually granting admin rights
	log.Debug("Marking new user as admin with real admin account")
	account.Admin = true
	err = SaveUser(client, serverURL, GetAdminUser(), GetAdminPassword(), account)
	require.Nil(t, err)

	// Does the old bundle still work?
	log.Debug("Trying existing bundle")
	version, err = dclient.Version()
	require.Nil(t, err)
	require.True(t, strings.Contains(version.Version, "ucp"))
	// We could verify no-swarm here, but we have other tests that cover that scenario

	// Get a new bundle, should be admin
	log.Debug("Getting new bundle")
	dclient, err = GetUserDockerClient(serverURL, newUsername, newPassword)
	require.Nil(t, err)
	log.Debug("getting version")
	version, err = dclient.Version()
	require.Nil(t, err)
	require.True(t, strings.Contains(version.Version, "ucp"))
	// And verify it works against swarm
	swarmURL := *dclient.URL
	swarmURL.Host = strings.Split(swarmURL.Host, ":")[0] + ":3376" // WARNING - fragile assumption
	swarmClient, err := dockerclient.NewDockerClient(swarmURL.String(), dclient.TLSConfig)
	require.Nil(t, err)
	version, err = swarmClient.Version()
	require.Nil(t, err)
	require.True(t, strings.Contains(version.Version, "swarm"))

	log.Debug("Trying to drop our own admin bit")
	account.Admin = false
	err = SaveUser(client, serverURL, newUsername, newPassword, account)
	// Verifying we were able to drop it
	account, err = GetUser(client, serverURL, newUsername, newPassword, newUsername)
	require.Nil(t, err)
	require.False(t, account.Admin)
	log.Debug("We dropped out own admin")

	// One last attempt to escalate just to be sure.
	account.Admin = true
	err = SaveUser(client, serverURL, newUsername, newPassword, account)
	log.Debugf("Attempt to escalate got %s result", err) // We don't really care, as long as it didn't work
	// Verifying we didn't actually get admin
	account, err = GetUser(client, serverURL, newUsername, newPassword, newUsername)
	require.Nil(t, err)
	require.Equal(t, account.Password, "")
	require.False(t, account.Admin)
	log.Debug("Couldn't escalate ourself back to admin - good!")

	// Make sure a non-admin can't delete an account
	newUsername2 := "newuser4"
	newPassword2 := "supersecret"
	log.Debug("Creating another user")
	require.Nil(t, CreateNewUser(client, serverURL, GetAdminUser(), GetAdminPassword(), newUsername2, newPassword2, false, auth.RestrictedControl))
	log.Debug("Attempting to delete another user with non-admin account")
	require.NotNil(t, DeleteUser(client, serverURL, newUsername, newPassword, newUsername2))
	log.Debug("Delete with an admin account")
	require.Nil(t, DeleteUser(client, serverURL, GetAdminUser(), GetAdminPassword(), newUsername2))

	log.Debug("Verifying admin can't see users hashed passwords")
	accounts, err := GetAllUsers(client, serverURL, GetAdminUser(), GetAdminPassword())
	require.Nil(t, err)
	for _, acct := range accounts {
		require.Equal(t, acct.Password, "")
	}

	log.Debug("Verifying non-admin can't get account list")
	accounts, err = GetAllUsers(client, serverURL, newUsername, newPassword)
	assert.NotNil(t, err)
	assert.Equal(t, len(accounts), 0)

}
