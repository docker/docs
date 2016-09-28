package api

import (
	"bytes"
	"crypto/tls"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"

	"github.com/docker/orca/auth"
	"github.com/docker/orca/integration/utils"
)

// A Test Team is used to hold all users of the authorization test
const TestTeamName = "auth-test-team"
const TestTeamLabel = "authzlabel"
const TestTeamLabelRole = auth.RestrictedControl

type API struct {
	Path   string
	Method string
	Body   string // JSON-formatted body string
	Status int
	Resp   string // JSON-formatted body string
}

type UserTest struct {
	Username  string
	Password  string
	Role      auth.Role
	Endpoints []API
}

const userPassChangeBody = "{\"new_password\":\"secret2\", \"old_password\":\"secret\"}"
const adminPassChangeBody = "{\"new_password\":\"orca\", \"old_password\":\"orca\"}"

const NoSuchContainer = "Error: No such container: 2222"

// OpenAuthEndpoints should be accessible by all authorized users in the same manner
var OpenAuthEndpoints = []API{
	{"/api/banner", "GET", "", http.StatusOK, ""},
	{"/api/catalog", "GET", "", http.StatusOK, ""},
	{"/api/nodes", "GET", "", http.StatusOK, ""},
	{"/api/clientbundle", "GET", "", http.StatusOK, ""},
	{"/api/clientbundle", "POST", "", http.StatusOK, ""},
	{"/api/applications", "GET", "", http.StatusOK, ""},
}

var TestMatrix = []UserTest{
	{
		Username: "usernoaccess",
		Password: "secret",
		Role:     auth.None,
		Endpoints: []API{
			{"/api/account", "GET", "", http.StatusOK, ""},
			{"/api/accounts", "GET", "", http.StatusForbidden, ""},
			{"/api/accounts", "POST", "", http.StatusInternalServerError, ""},
			{"/api/accounts/usernoaccess", "GET", "", http.StatusOK, ""},
			{"/api/accounts/usernoaccess", "DELETE", "", http.StatusForbidden, ""},
			{"/api/accesslists?username=usernoaccess", "GET", "", http.StatusOK, ""},
			{"/api/accesslists?username=admin", "GET", "", http.StatusForbidden, ""},
			{"/api/accesslists", "POST", "{}", http.StatusForbidden, ""},
			{"/api/accesslists", "GET", "", http.StatusForbidden, ""},
			{"/api/accesslists/5/2", "GET", "", http.StatusForbidden, ""},
			{"/api/accesslists/{teamid}/5", "GET", "", http.StatusForbidden, ""},
			{"/api/accesslists/{teamid}/5", "DELETE", "", http.StatusForbidden, ""},
			{"/api/teams", "POST", "", http.StatusForbidden, ""},
			{"/api/teams", "GET", "", http.StatusForbidden, ""},
			{"/api/teams/{teamid}", "GET", "", http.StatusOK, ""},
			{"/api/nodes/authorize", "POST", "{}", http.StatusForbidden, ""},
			{"/api/authsync", "GET", "", http.StatusForbidden, ""},
			{"/api/authsync", "POST", "{}", http.StatusForbidden, ""},
			{"/api/support", "POST", "", http.StatusForbidden, ""},
			{"/api/accounts/admin", "GET", "", http.StatusForbidden, ""},
			{"/api/webhookkeys", "POST", "", http.StatusForbidden, ""},
			{"/api/teams/5", "PUT", "", http.StatusForbidden, ""},
			{"/api/teams/5", "GET", "", http.StatusForbidden, ""},
			{"/api/teams/5", "DELETE", "", http.StatusForbidden, ""},
			{"/api/teams/{teamid}/members/add/usernoaccess", "PUT", "", http.StatusForbidden, ""},
			{"/api/teams/{teamid}/members/remove/usernoaccess", "DELETE", "", http.StatusForbidden, ""},
			{"/api/containers/2222/scale?n=3", "POST", "", http.StatusForbidden, ""},
			{"/api/containerlogs/2222", "POST", "", http.StatusForbidden, ""},
			{"/api/webhookkeys", "GET", "", http.StatusForbidden, ""},
			{"/api/webhookkeys/key", "GET", "", http.StatusForbidden, ""},
			{"/api/webhookkeys/key", "DELETE", "", http.StatusForbidden, ""},
			{"/api/config", "GET", "", http.StatusForbidden, ""},
			{"/api/config/registry", "GET", "", http.StatusForbidden, ""},
			{"/account/changepassword", "POST", "", http.StatusBadRequest, ""},
			{"/account/changepassword", "POST", userPassChangeBody, http.StatusOK, ""},
			{"/account/changepassword", "POST", adminPassChangeBody, http.StatusInternalServerError, ""},
			{"/api/teams/{teamid}", "DELETE", "", http.StatusForbidden, ""},
		},
	},
	{
		Username: "userviewonly",
		Password: "secret",
		Role:     auth.View,
		Endpoints: []API{
			{"/api/account", "GET", "", http.StatusOK, ""},
			{"/api/accounts", "GET", "", http.StatusForbidden, ""},
			{"/api/accounts", "POST", "", http.StatusInternalServerError, ""},
			{"/api/accounts/userviewonly", "GET", "", http.StatusOK, ""},
			{"/api/accounts/userviewonly", "DELETE", "", http.StatusForbidden, ""},
			{"/api/accounts/admin", "GET", "", http.StatusForbidden, ""},
			{"/api/accesslists", "POST", "{}", http.StatusForbidden, ""},
			{"/api/accesslists", "GET", "", http.StatusForbidden, ""},
			{"/api/accesslists?username=userviewonly", "GET", "", http.StatusOK, ""},
			{"/api/accesslists?username=admin", "GET", "", http.StatusForbidden, ""},
			{"/api/accesslists/{teamid}/5", "GET", "", http.StatusForbidden, ""},
			{"/api/accesslists/{teamid}/5", "DELETE", "", http.StatusForbidden, ""},
			{"/api/accesslists/5/2", "GET", "", http.StatusForbidden, ""},
			{"/api/teams", "POST", "", http.StatusForbidden, ""},
			{"/api/teams", "GET", "", http.StatusForbidden, ""},
			{"/api/nodes/authorize", "POST", "{}", http.StatusForbidden, ""},
			{"/api/authsync", "GET", "", http.StatusForbidden, ""},
			{"/api/authsync", "POST", "{}", http.StatusForbidden, ""},
			{"/api/support", "POST", "", http.StatusForbidden, ""},
			{"/api/webhookkeys", "POST", "", http.StatusForbidden, ""},
			{"/api/teams/5", "PUT", "", http.StatusForbidden, ""},
			{"/api/teams/5", "GET", "", http.StatusForbidden, ""},
			{"/api/teams/{teamid}", "GET", "", http.StatusOK, ""},
			{"/api/teams/5", "DELETE", "", http.StatusForbidden, ""},
			{"/api/teams/{teamid}/members/add/usernoaccess", "PUT", "", http.StatusForbidden, ""},
			{"/api/teams/{teamid}/members/remove/usernoaccess", "DELETE", "", http.StatusForbidden, ""},
			{"/api/containers/2222/scale?n=3", "POST", "", http.StatusForbidden, ""},
			{"/api/containerlogs/2222", "POST", "", http.StatusForbidden, ""},
			{"/api/webhookkeys", "GET", "", http.StatusForbidden, ""},
			{"/api/webhookkeys/key", "GET", "", http.StatusForbidden, ""},
			{"/api/webhookkeys/key", "DELETE", "", http.StatusForbidden, ""},
			{"/api/config", "GET", "", http.StatusForbidden, ""},
			{"/api/config/registry", "GET", "", http.StatusForbidden, ""},
			{"/account/changepassword", "POST", "", http.StatusBadRequest, ""},
			{"/account/changepassword", "POST", userPassChangeBody, http.StatusOK, ""},
			{"/account/changepassword", "POST", adminPassChangeBody, http.StatusInternalServerError, ""},
			{"/api/teams/{teamid}", "DELETE", "", http.StatusForbidden, ""},
		},
	},
	{
		Username: "userrestricted",
		Password: "secret",
		Role:     auth.RestrictedControl,
		Endpoints: []API{
			{"/api/account", "GET", "", http.StatusOK, ""},
			{"/api/accounts", "GET", "", http.StatusForbidden, ""},
			{"/api/accounts", "POST", "", http.StatusInternalServerError, ""},
			{"/api/accounts/userrestricted", "GET", "", http.StatusOK, ""},
			{"/api/accounts/userrestricted", "DELETE", "", http.StatusForbidden, ""},
			{"/api/accounts/admin", "GET", "", http.StatusForbidden, ""},
			{"/api/accesslists", "POST", "{}", http.StatusForbidden, ""},
			{"/api/accesslists", "GET", "", http.StatusForbidden, ""},
			{"/api/accesslists?username=userrestricted", "GET", "", http.StatusOK, ""},
			{"/api/accesslists?username=admin", "GET", "", http.StatusForbidden, ""},
			{"/api/accesslists/{teamid}/5", "GET", "", http.StatusForbidden, ""},
			{"/api/accesslists/{teamid}/5", "DELETE", "", http.StatusForbidden, ""},
			{"/api/accesslists/5/2", "GET", "", http.StatusForbidden, ""},
			{"/api/teams", "POST", "", http.StatusForbidden, ""},
			{"/api/teams", "GET", "", http.StatusForbidden, ""},
			{"/api/nodes/authorize", "POST", "{}", http.StatusForbidden, ""},
			{"/api/authsync", "GET", "", http.StatusForbidden, ""},
			{"/api/authsync", "POST", "{}", http.StatusForbidden, ""},
			{"/api/support", "POST", "", http.StatusForbidden, ""},
			{"/api/webhookkeys", "POST", "", http.StatusForbidden, ""},
			{"/api/teams/5", "PUT", "", http.StatusForbidden, ""},
			{"/api/teams/5", "GET", "", http.StatusForbidden, ""},
			{"/api/teams/{teamid}", "GET", "", http.StatusOK, ""},
			{"/api/teams/5", "DELETE", "", http.StatusForbidden, ""},
			{"/api/teams/{teamid}/members/add/usernoaccess", "PUT", "", http.StatusForbidden, ""},
			{"/api/teams/{teamid}/members/remove/usernoaccess", "DELETE", "", http.StatusForbidden, ""},
			{"/api/containers/2222/scale?n=3", "POST", "", http.StatusForbidden, ""},
			{"/api/containerlogs/2222", "POST", "", http.StatusForbidden, ""},
			{"/api/webhookkeys", "GET", "", http.StatusForbidden, ""},
			{"/api/webhookkeys/key", "GET", "", http.StatusForbidden, ""},
			{"/api/webhookkeys/key", "DELETE", "", http.StatusForbidden, ""},
			{"/api/config", "GET", "", http.StatusForbidden, ""},
			{"/api/config/registry", "GET", "", http.StatusForbidden, ""},
			{"/account/changepassword", "POST", "", http.StatusBadRequest, ""},
			{"/account/changepassword", "POST", userPassChangeBody, http.StatusOK, ""},
			{"/account/changepassword", "POST", adminPassChangeBody, http.StatusInternalServerError, ""},
			{"/api/teams/{teamid}", "DELETE", "", http.StatusForbidden, ""},
		},
	},
	{
		Username: "userfull",
		Password: "secret",
		Role:     auth.FullControl,
		Endpoints: []API{
			{"/api/account", "GET", "", http.StatusOK, ""},
			{"/api/accounts", "GET", "", http.StatusForbidden, ""},
			{"/api/accounts", "POST", "", http.StatusInternalServerError, ""},
			{"/api/accounts/userfull", "GET", "", http.StatusOK, ""},
			{"/api/accounts/userfull", "DELETE", "", http.StatusForbidden, ""},
			{"/api/accounts/admin", "GET", "", http.StatusForbidden, ""},
			{"/api/accesslists", "GET", "", http.StatusForbidden, ""},
			{"/api/accesslists", "POST", "{}", http.StatusForbidden, ""},
			{"/api/accesslists?username=userfull", "GET", "", http.StatusOK, ""},
			{"/api/accesslists?username=admin", "GET", "", http.StatusForbidden, ""},
			{"/api/accesslists/{teamid}/5", "GET", "", http.StatusForbidden, ""},
			{"/api/accesslists/{teamid}/5", "DELETE", "", http.StatusForbidden, ""},
			{"/api/accesslists/5/2", "GET", "", http.StatusForbidden, ""},
			{"/api/teams", "POST", "", http.StatusForbidden, ""},
			{"/api/teams", "GET", "", http.StatusForbidden, ""},
			{"/api/teams/{teamid}", "GET", "", http.StatusOK, ""},
			{"/api/nodes/authorize", "POST", "{}", http.StatusForbidden, ""},
			{"/api/authsync", "GET", "", http.StatusForbidden, ""},
			{"/api/authsync", "POST", "{}", http.StatusForbidden, ""},
			{"/api/support", "POST", "", http.StatusForbidden, ""},
			{"/api/webhookkeys", "POST", "", http.StatusForbidden, ""},
			{"/api/teams/5", "PUT", "", http.StatusForbidden, ""},
			{"/api/teams/5", "GET", "", http.StatusForbidden, ""},
			{"/api/teams/5", "DELETE", "", http.StatusForbidden, ""},
			{"/api/teams/{teamid}/members/add/usernoaccess", "PUT", "", http.StatusForbidden, ""},
			{"/api/teams/{teamid}/members/remove/usernoaccess", "DELETE", "", http.StatusForbidden, ""},
			{"/api/containers/2222/scale?n=3", "POST", "", http.StatusForbidden, ""},
			{"/api/containerlogs/2222", "POST", "", http.StatusForbidden, ""},
			{"/api/webhookkeys", "GET", "", http.StatusForbidden, ""},
			{"/api/webhookkeys/key", "GET", "", http.StatusForbidden, ""},
			{"/api/webhookkeys/key", "DELETE", "", http.StatusForbidden, ""},
			{"/api/config", "GET", "", http.StatusForbidden, ""},
			{"/api/config/registry", "GET", "", http.StatusForbidden, ""},
			{"/account/changepassword", "POST", "", http.StatusBadRequest, ""},
			{"/account/changepassword", "POST", userPassChangeBody, http.StatusOK, ""},
			{"/account/changepassword", "POST", adminPassChangeBody, http.StatusInternalServerError, ""},
			{"/api/teams/{teamid}", "DELETE", "", http.StatusForbidden, ""},
		},
	},
	{
		Username: "admin",
		Password: "orca",
		Role:     auth.Admin,
		Endpoints: []API{
			{"/api/account", "GET", "", http.StatusOK, ""},
			{"/api/accounts", "GET", "", http.StatusOK, ""},
			{"/api/accounts", "POST", "", http.StatusInternalServerError, ""},
			{"/api/accesslists?username=usernoaccess", "GET", "", http.StatusOK, ""},
			{"/api/accesslists?username=admin", "GET", "", http.StatusOK, ""},
			{"/api/accesslists", "POST", "{}", http.StatusCreated, ""},
			{"/api/accesslists", "GET", "", http.StatusOK, ""},
			{"/api/accesslists/{teamid}/5", "GET", "", http.StatusInternalServerError, ""},
			{"/api/accesslists/{teamid}/5", "DELETE", "", http.StatusInternalServerError, ""},
			{"/api/accesslists/5/2", "GET", "", http.StatusInternalServerError, ""},
			{"/api/teams", "GET", "", http.StatusOK, ""},
			{"/api/teams", "POST", "{}", http.StatusInternalServerError, ""},
			{"/api/teams/{teamid}", "GET", "", http.StatusOK, ""},
			{"/api/nodes/authorize", "POST", "{}", http.StatusInternalServerError, ""},
			{"/api/authsync", "GET", "", http.StatusOK, ""},
			{"/api/authsync", "POST", "{}", http.StatusOK, ""},
			{"/api/support", "POST", "", http.StatusOK, ""},
			{"/api/accounts/admin", "GET", "", http.StatusOK, ""},
			{"/api/webhookkeys", "POST", "", http.StatusInternalServerError, ""},
			{"/api/teams/5", "PUT", "", http.StatusInternalServerError, ""},
			{"/api/teams/5", "GET", "", http.StatusInternalServerError, ""},
			{"/api/teams/5", "DELETE", "", http.StatusInternalServerError, ""},
			{"/api/teams/{teamid}/members/add/usernoaccess", "PUT", "", http.StatusNoContent, ""},
			{"/api/teams/{teamid}/members/remove/usernoaccess", "DELETE", "", http.StatusNoContent, ""},
			{"/api/accounts/usernoaccess", "GET", "", http.StatusOK, ""},
			{"/api/accounts/usernoaccess", "DELETE", "", http.StatusNoContent, ""},
			{"/api/containers/2222/scale?n=3", "POST", "", http.StatusBadRequest, NoSuchContainer},
			{"/api/containerlogs/2222", "POST", "", http.StatusBadRequest, NoSuchContainer},
			{"/api/webhookkeys", "GET", "", http.StatusOK, ""},
			{"/api/webhookkeys/key", "GET", "", http.StatusInternalServerError, ""},
			{"/api/webhookkeys/key", "DELETE", "", http.StatusInternalServerError, ""},
			{"/api/config", "GET", "", http.StatusOK, ""},
			{"/api/config/registry", "GET", "", http.StatusOK, ""},
			{"/account/changepassword", "POST", "", http.StatusBadRequest, ""},
			{"/account/changepassword", "POST", userPassChangeBody, http.StatusInternalServerError, ""},
			{"/account/changepassword", "POST", adminPassChangeBody, http.StatusOK, ""},
			{"/api/teams/{teamid}", "DELETE", "", http.StatusNoContent, ""},
		},
	},
}

func getAdminHTTPClient(serverURL string) *http.Client {
	// Authenticate with admin user through TLS
	tlsConfig, _ := utils.GetUserTLSConfig(serverURL, utils.GetAdminUser(), utils.GetAdminPassword())
	tlsConfig.InsecureSkipVerify = true

	tr := &http.Transport{
		TLSClientConfig:     tlsConfig,
		DisableKeepAlives:   true,
		MaxIdleConnsPerHost: 0,
	}

	return &http.Client{
		Timeout:   TIMEOUT,
		Transport: tr,
	}
}

// TestNonAdmin endpoints manually pokes at all the API endpoints that
// should be accessible by an admin user but unauthorized by a non-admin
func (s *APITestSuite) TestAPIAuthorization() {
	require := require.New(s.T())

	serverURL, err := utils.GetOrcaURL(s.ControllerMachines[0])
	require.Nil(err)

	// Create the Test Team
	adminClient := getAdminHTTPClient(serverURL)
	teamID, err := utils.AddTeam(adminClient, serverURL, utils.GetAdminUser(), utils.GetAdminPassword(), TestTeamName)
	require.Nil(err)

	// Add an AccessList
	err = utils.AddAccessList(adminClient, serverURL, utils.GetAdminUser(), utils.GetAdminPassword(),
		TestTeamName, TestTeamLabel, TestTeamLabelRole)
	require.Nil(err)

	var tlsConfig *tls.Config
	for _, userTest := range TestMatrix {
		if userTest.Username == "admin" {
			// If the test is using the admin user, just return the tlsConfig
			tlsConfig, err = utils.GetUserTLSConfig(serverURL, utils.GetAdminUser(), utils.GetAdminPassword())
			require.Nil(err)
		} else {
			// Other users need to be created first
			err = utils.CreateNewUser(adminClient, serverURL,
				utils.GetAdminUser(), utils.GetAdminPassword(),
				userTest.Username, userTest.Password, false, userTest.Role)
			require.Nil(err)
			tlsConfig, err = utils.GetUserTLSConfig(serverURL, userTest.Username, userTest.Password)
			require.Nil(err)
		}
		// Add the user to the Test Team
		err = utils.AddTeamMember(adminClient, serverURL, utils.GetAdminUser(), utils.GetAdminPassword(),
			teamID, userTest.Username)
		require.Nil(err)

		// Create that user's HTTP client
		tlsConfig.InsecureSkipVerify = true

		tr := &http.Transport{
			TLSClientConfig:     tlsConfig,
			DisableKeepAlives:   true,
			MaxIdleConnsPerHost: 0,
		}

		userClient := &http.Client{
			Timeout:   TIMEOUT,
			Transport: tr,
		}

		// Test all the endpoints for that user
		for _, endpoint := range append(userTest.Endpoints, OpenAuthEndpoints...) {
			// Replace the {teamid} tags in the path using the gorilla/mux package
			route := mux.Route{}
			route.BuildOnly()
			route.Path(endpoint.Path)
			url, err := route.URLPath("teamid", teamID)

			req, err := http.NewRequest(endpoint.Method, serverURL+url.Path, bytes.NewBufferString(endpoint.Body))
			require.Nil(err)
			res, err := userClient.Do(req)
			require.Nil(err)
			// Require that no requests are forbidden
			log.Debug(url.Path)
			require.Equal(endpoint.Status, res.StatusCode)
			if endpoint.Resp != "" {
				// Perform assertions on the response body, if specified
				defer res.Body.Close()
				buf := new(bytes.Buffer)
				_, err = buf.ReadFrom(res.Body)
				require.Nil(err)
				require.Equal(endpoint.Resp, strings.TrimSpace(buf.String()))
			}
		}
	}
}
