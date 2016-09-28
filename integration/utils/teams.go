package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	neturl "net/url"
	"path"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/orca"
	"github.com/docker/orca/auth"
	enziclient "github.com/docker/orca/enzi/api/client"
	"github.com/stretchr/testify/require"
)

// TeamLDAPOpts hold optional arguments to the AddTeam function.
type TeamLDAPOpts struct {
	GroupDN    string
	MemberAttr string
}

// Returns the team ID
func AddTeam(client *http.Client, serverURL, adminUser, adminPassword, teamname string, ldapOpts ...TeamLDAPOpts) (string, error) {
	orcaURL, err := neturl.Parse(serverURL)
	if err != nil {
		return "", err
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
		return "", err
	}

	orcaURL.Path = "/api/teams"

	reqData := map[string]string{
		"name": teamname,
	}

	if len(ldapOpts) > 0 {
		reqData["ldapdn"] = ldapOpts[0].GroupDN
		reqData["ldap_member_attr"] = ldapOpts[0].MemberAttr
	}

	reqJSON, err := json.Marshal(reqData)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", orcaURL.String(), bytes.NewBuffer(reqJSON))
	if err != nil {
		// Should never fail
		return "", err
	}
	req.Header.Set(GetTokenHeader(token))
	resp, err := client.Do(req)
	if err != nil {
		log.Debug("Failed to make request")
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		return "", fmt.Errorf(string(body))
	}
	log.Debugf("Succesfully added team %s", teamname)
	return path.Base(resp.Header.Get("Location")), nil
}

func FindTeam(client *http.Client, serverURL, adminUser, adminPassword, teamname string) (string, error) {
	orcaURL, err := neturl.Parse(serverURL)
	if err != nil {
		return "", err
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
		return "", err
	}

	orcaURL.Path = "/api/teams"

	req, err := http.NewRequest("GET", orcaURL.String(), nil)
	if err != nil {
		// Should never fail
		return "", err
	}
	req.Header.Set(GetTokenHeader(token))
	resp, err := client.Do(req)
	if err != nil {
		log.Debug("Failed to make request")
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		return "", fmt.Errorf(string(body))
	}

	var teams []auth.Team
	if err := json.Unmarshal(body, &teams); err != nil {
		return "", err
	}
	for _, team := range teams {
		if team.Name == teamname {
			return team.Id, nil
		}
	}
	return "", fmt.Errorf("Failed to locate team %s", teamname)
}

func AddAccessList(client *http.Client, serverURL, adminUser, adminPassword, teamID, label string, role auth.Role) error {
	log.Debugf("Attempting to setup access : %s->%s", teamID, label)
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

	orcaURL.Path = "/api/accesslists"

	data := []byte(
		fmt.Sprintf(`{
    "role":%d,
    "teamId":"%s",
    "label":"%s"
}`,
			role,
			teamID,
			label,
		))

	// Now create the access lists
	req, err := http.NewRequest("POST", orcaURL.String(), bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set(GetTokenHeader(token))
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		return fmt.Errorf(string(body))
	}
	log.Debugf("Succesfully added access list with label %s", label)
	return nil
}

func AddTeamMember(client *http.Client, serverURL, adminUser, adminPassword, teamID, username string) error {
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

	orcaURL.Path = fmt.Sprintf("/api/teams/%s/members/add/%s", teamID, username)

	req, err := http.NewRequest("PUT", orcaURL.String(), nil)
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
		return fmt.Errorf(string(body))
	}
	log.Debugf("Succesfully added team %s member %s", teamID, username)
	return nil
}

func GetTeamMember(client *http.Client, serverURL, adminUser, adminPassword, teamID string) ([]string, error) {
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

	orcaURL.Path = fmt.Sprintf("/api/teams/%s", teamID)

	req, err := http.NewRequest("GET", orcaURL.String(), nil)
	if err != nil {
		// Should never fail
		return nil, err
	}
	req.Header.Set(GetTokenHeader(token))
	resp, err := client.Do(req)
	if err != nil {
		log.Debugf("Failed to make request: %s", req.URL.String())
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		log.Debugf("Failed to make request: %s", req.URL.String())
		return nil, fmt.Errorf(string(body))
	}
	var team auth.Team
	if err := json.Unmarshal(body, &team); err != nil {
		log.Debugf("Failed to unmarshal request: %s: %s", req.URL.String(), string(body))
		return nil, err
	}

	// Now double check we get the same results from the other endpoint
	orcaURL.Path = "/api/accounts"
	orcaURL.RawQuery = "teamId=" + teamID
	req, err = http.NewRequest("GET", orcaURL.String(), nil)
	if err != nil {
		// Should never fail
		log.Debugf("Failed to make request: %s", req.URL.String())
		return nil, err
	}
	req.Header.Set(GetTokenHeader(token))
	resp, err = client.Do(req)
	if err != nil {
		log.Debugf("Failed to make request: %s", req.URL.String())
		return nil, err
	}

	body, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		log.Debugf("Failed to make request: %s", req.URL.String())
		return nil, fmt.Errorf(string(body))
	}
	var members []auth.Account
	if err := json.Unmarshal(body, &members); err != nil {
		log.Debugf("Failed to unmarshal request: %s: %s", req.URL.String(), string(body))
		return nil, err
	}

	if len(team.ManagedMembers)+len(team.DiscoveredMembers) != len(members) {
		return nil, fmt.Errorf("Inconsistent membership: %#v does not match %#v", team, members)
	}

	// Not exhaustive, but probably good enough to catch any bugs...
	fullMembers := append(team.ManagedMembers, team.DiscoveredMembers...)
	for _, member := range members {
		found := false
		for _, comp := range fullMembers {
			if comp == member.Username {
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("Inconsistent membership: %#v does not match %#v", team, members)
		}
	}
	return fullMembers, nil
}

func TestManagedMembers(t *testing.T, serverURL string) {
	// Create a managed group, and verify membership looks sane
	log.Debug("Creating team")
	teamID, err := AddTeam(nil, serverURL, GetAdminUser(), GetAdminPassword(), "testmanagedgrouplocal")
	require.Nil(t, err)
	newUsername := "testnewuser"
	newPassword := "supersecret"
	log.Debug("Creating user")
	require.Nil(t, CreateNewUser(nil, serverURL, GetAdminUser(), GetAdminPassword(), newUsername, newPassword, false, auth.RestrictedControl))
	log.Debug("Adding user to team")
	require.Nil(t, AddTeamMember(nil, serverURL, GetAdminUser(), GetAdminPassword(), teamID, newUsername))
	members, err := GetTeamMember(nil, serverURL, GetAdminUser(), GetAdminPassword(), teamID)
	require.Nil(t, err)
	require.Equal(t, members, []string{newUsername})
	log.Debug("Managed membership looks good")
}

var upgradeUserNames = []string{
	"upgrade_user_1",
	"upgrade_user_2",
	"upgrade_user_3",
}

var upgradeTeamName = "upgrade_team"

// SetupUpgradeUsersAndTeam sets up a few users and a team with them as members
// to be used for upgrade tests. This function should be run to setup the team
// *before* the upgrade. After the upgrade, the suite should run the function
// TestUpgradeUsersAndTeam.
func SetupUpgradeUsersAndTeam(t *testing.T, serverURL string) {
	// Create a single team using the UCP API.
	teamID, err := AddTeam(nil, serverURL, GetAdminUser(), GetAdminPassword(), upgradeTeamName)
	require.Nil(t, err)

	// Create a few users and add them to the team using the UCP accounts
	// API.
	for _, upgradeUserName := range upgradeUserNames {
		require.Nil(t, CreateNewUser(nil, serverURL, GetAdminUser(), GetAdminPassword(), upgradeUserName, "maplesyrupbowflexfolgers", false, auth.None))
		require.Nil(t, AddTeamMember(nil, serverURL, GetAdminUser(), GetAdminPassword(), teamID, upgradeUserName))
	}
}

// TestUpgradeUsersAndTeam tests that the users and team created in the
// function SetupUpgradeUsersAndTeam were migrated successfully after the
// upgrade.
func TestUpgradeUsersAndTeam(t *testing.T, serverAddr string) {
	httpClient := &http.Client{
		Transport: &http.Transport{
			// Sloppy for testing only - don't copy this into production code!
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 10 * time.Second, // Be pretty aggressive on timeouts for testing
	}

	basicAuth := &enziclient.BasicAuthenticator{
		Username: GetAdminUser(),
		Password: GetAdminPassword(),
	}

	enziSession := enziclient.New(httpClient, serverAddr, "enzi", basicAuth)

	// First, make sure that the accounts exist in the auth service. They
	// should all be members of the default org. They should also be a
	// member of the upgrade test team.
	for _, upgradeUserName := range upgradeUserNames {
		user, err := enziSession.GetAccount(upgradeUserName)
		require.Nil(t, err, "unable to get user %s", upgradeUserName)
		require.False(t, user.IsOrg)

		orgMember, err := enziSession.GetOrganizationMember(orca.UCPDefaultOrg, upgradeUserName)
		require.Nil(t, err, "unable to get default org membership for user %s", upgradeUserName)
		require.NotNil(t, orgMember, "user is not a member of the org")

		teamMember, err := enziSession.GetTeamMember(orca.UCPDefaultOrg, upgradeTeamName, upgradeUserName)
		require.Nil(t, err, "unable to get team membership for user %s", upgradeUserName)
		require.NotNil(t, teamMember, "user is not a member of the team")
	}
}
