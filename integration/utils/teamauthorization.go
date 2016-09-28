package utils

import (
	"fmt"
	"strings"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/swarm"
	"github.com/samalba/dockerclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"

	"github.com/docker/orca"
	"github.com/docker/orca/auth"
)

// Test data tables

type TestTeamScenario struct {
	// Short names to keep the table from getting unwieldy
	N string
	L *TestLabel
	T *TestTeam
	O func(serverURL string, team *TestTeam, label *TestLabel) string
	E string
}

type TestUser struct {
	Username string
	Client   *dockerclient.DockerClient
	Teams    []*TestTeam
}

type TestLabel struct {
	Name       string
	Containers []string
	Services   []string
	Networks   []string
}

type TestTeam struct {
	Name   string
	TeamID string // Loaded once initialized
	Users  []*TestUser
	Label  string // Label is the Label with the highest permission level for a given team
}

type TestAccess struct {
	Team  *TestTeam
	Label *TestLabel
	Role  auth.Role
}

// Worker routines

func createTeamCnt(serverURL string, team *TestTeam, label *TestLabel) string {
	for _, user := range team.Users {
		log.Debugf("About to create container for %s with user %s", team.Name, user.Username)
		if user.Client == nil {
			client, err := GetUserDockerClient(serverURL, user.Username, Password)
			if err != nil {
				return err.Error()
			}
			user.Client = client
		}
		containerName := fmt.Sprintf("%s_container_%s", team.Name, user.Username)
		cfg := &dockerclient.ContainerConfig{
			Image: Image,
			Cmd:   []string{"sh", "-c", fmt.Sprintf("echo %s; sleep 1h", ContainerOutput)},
			Labels: map[string]string{
				orca.UCPAccessLabel: label.Name,
			},
		}
		containerId, err := user.Client.CreateContainer(cfg, containerName, nil)
		if err != nil {
			return err.Error()
		}
		err = user.Client.StartContainer(containerId, nil)
		if err != nil {
			return err.Error()
		}
		log.Debugf("Created and started %s with user %s", containerName, user.Username)
		// If we managed to start it, keep track so we can use it later for poking around
		label.Containers = append(label.Containers, containerName)
	}

	return "OK"
}

func createTeamService(serverURL string, team *TestTeam, label *TestLabel) string {
	for _, user := range team.Users {
		log.Debugf("About to create service for %s with user %s", team.Name, user.Username)
		if user.Client == nil {
			client, err := GetUserDockerClient(serverURL, user.Username, Password)
			if err != nil {
				return err.Error()
			}
			user.Client = client
		}

		engineClient, err := ConvertToEngineAPI(user.Client)
		if err != nil {
			return err.Error()
		}

		serviceName := user.Username + team.Name
		serviceResp, err := engineClient.ServiceCreate(context.TODO(), swarm.ServiceSpec{
			Annotations: swarm.Annotations{
				Name: serviceName,
				Labels: map[string]string{
					orca.UCPAccessLabel: label.Name,
				},
			},
			TaskTemplate: swarm.TaskSpec{
				ContainerSpec: swarm.ContainerSpec{
					Image:   Image,
					Command: []string{"ping", "google.com"},
				},
			},
		}, types.ServiceCreateOptions{})
		if err != nil {
			return err.Error()
		}

		// If we managed to start it, keep track so we can use it later for poking around
		label.Services = append(label.Services, serviceResp.ID)
	}

	return "OK"
}

func inspectTeamService(serverURL string, team *TestTeam, label *TestLabel) string {
	for _, user := range team.Users {
		log.Debugf("About to inspect services for team %s with user %s", team.Name, user.Username)
		if user.Client == nil {
			client, err := GetUserDockerClient(serverURL, user.Username, Password)
			if err != nil {
				return err.Error()
			}
			user.Client = client
		}

		engineClient, err := ConvertToEngineAPI(user.Client)
		if err != nil {
			return err.Error()
		}

		for _, serviceID := range label.Services {
			// Inspect the current service
			service, _, err := engineClient.ServiceInspectWithRaw(context.TODO(), serviceID)
			if err != nil {
				return err.Error()
			}

			// Require service ID equality
			if service.ID != serviceID {
				return fmt.Sprintf("Inspected service ID %s is different than expected: %s ", service.ID, serviceID)
			}
		}
	}

	return "OK"
}

func createTeamNetwork(serverURL string, team *TestTeam, label *TestLabel) string {
	for _, user := range team.Users {
		log.Debugf("About to create network for %s with user %s", team.Name, user.Username)
		if user.Client == nil {
			client, err := GetUserDockerClient(serverURL, user.Username, Password)
			if err != nil {
				return err.Error()
			}
			user.Client = client
		}

		engineClient, err := ConvertToEngineAPI(user.Client)
		if err != nil {
			return err.Error()
		}

		networkName := user.Username + team.Name
		networkResp, err := engineClient.NetworkCreate(context.TODO(), networkName, types.NetworkCreate{
			Driver: "bridge", // TODO: convert to overlay for HA acceptance tests when backwards compatibility is in
			Labels: map[string]string{
				orca.UCPAccessLabel: label.Name,
			},
		})

		if err != nil {
			return err.Error()
		}

		// If we managed to create the network, keep track so we can use it later for poking around
		label.Networks = append(label.Networks, networkResp.ID)
	}

	return "OK"
}

func inspectTeamNetwork(serverURL string, team *TestTeam, label *TestLabel) string {
	for _, user := range team.Users {
		log.Debugf("About to inspect networks for team %s with user %s", team.Name, user.Username)
		if user.Client == nil {
			client, err := GetUserDockerClient(serverURL, user.Username, Password)
			if err != nil {
				return err.Error()
			}
			user.Client = client
		}

		engineClient, err := ConvertToEngineAPI(user.Client)
		if err != nil {
			return err.Error()
		}
		for _, networkID := range label.Networks {
			// Inspect the current service
			network, err := engineClient.NetworkInspect(context.TODO(), networkID)
			if err != nil {
				return err.Error()
			}

			// Require network ID equality
			if network.ID != networkID {
				return fmt.Sprintf("Inspected network ID %s is different than expected: %s ", network.ID, networkID)
			}
		}
	}

	return "OK"
}

func createNetworkContainer(serverURL string, team *TestTeam, label *TestLabel) string {
	for _, user := range team.Users {
		log.Debugf("About to create network container for %s with user %s", team.Name, user.Username)
		cli, err := NewCLIHandle(serverURL, user.Username, Password)
		if err != nil {
			return err.Error()
		}

		res, err := cli.RunCLICommand([]string{

			"docker", "run", "-d", "--net", label.Networks[0], "--label", orca.UCPAccessLabel + "=" + team.Label, Image, "sh", "-c", fmt.Sprintf("echo %s; sleep 1h", ContainerOutput),
		})
		if err != nil {
			return res + " " + err.Error()
		}

		log.Debugf("Created and started networked container with user %s", user.Username)
	}

	return "OK"
}

// connnectDisconnectContainer tests whether it's possible to connect and disconnect a container
// to a network of the same target label
func connectDisconnectContainer(serverURL string, team *TestTeam, label *TestLabel) string {
	for _, user := range team.Users {
		log.Debugf("About to connect/disconnect networks for team %s with user %s", team.Name, user.Username)

		if user.Client == nil {
			client, err := GetUserDockerClient(serverURL, user.Username, Password)
			if err != nil {
				return err.Error()
			}
			user.Client = client
		}

		engineClient, err := ConvertToEngineAPI(user.Client)
		if err != nil {
			return err.Error()
		}
		// Pick one network with the target label
		networkID := label.Networks[0]
		// Pick one container with the target label
		containerName := label.Containers[0]

		// Perform a network connect
		err = engineClient.NetworkConnect(context.TODO(), networkID, containerName, nil)
		if err != nil {
			return err.Error()
		}
		// Perform a network disconnect - don't --force
		err = engineClient.NetworkDisconnect(context.TODO(), networkID, containerName, false)
		if err != nil {
			return err.Error()
		}
	}
	return "OK"
}

func createTeamCntPriv(serverURL string, team *TestTeam, label *TestLabel) string {
	for _, user := range team.Users {
		log.Debugf("About to create privileged container for %s with user %s", team.Name, user.Username)
		if user.Client == nil {
			client, err := GetUserDockerClient(serverURL, user.Username, Password)
			if err != nil {
				return err.Error()
			}
			user.Client = client
		}
		containerName := fmt.Sprintf("%s_priv_container_%s", team.Name, user.Username)
		cfg := &dockerclient.ContainerConfig{
			Image: Image,
			Cmd:   []string{"sh", "-c", fmt.Sprintf("echo %s; sleep 1h", ContainerOutput)},
			Labels: map[string]string{
				orca.UCPAccessLabel: label.Name,
			},
			HostConfig: dockerclient.HostConfig{
				Privileged: true,
			},
		}
		containerId, err := user.Client.CreateContainer(cfg, containerName, nil)
		if err != nil {
			return err.Error()
		}
		err = user.Client.StartContainer(containerId, nil)
		if err != nil {
			return err.Error()
		}
		log.Debugf("Created and started %s with user %s", containerName, user.Username)
		// If we managed to start it, keep track so we can use it later for poking around
		label.Containers = append(label.Containers, containerName)
	}

	return "OK"
}

func listContainers(serverURL string, team *TestTeam, label *TestLabel) string {
	for _, user := range team.Users {
		if _, err := user.Client.ListContainers(true, false, ""); err != nil {
			return err.Error()
		}
	}
	return "OK"
}

func inspByLbl(serverURL string, team *TestTeam, label *TestLabel) string {
	for _, user := range team.Users {
		for _, containerName := range label.Containers {
			_, err := user.Client.InspectContainer(containerName)
			if err != nil {
				return fmt.Sprintf("User %s failed to lookup container %s: %s", user.Username, containerName, err)
			}
		}
	}

	return "OK"
}

func copyByLbl(serverURL string, team *TestTeam, label *TestLabel) string {
	for _, user := range team.Users {
		for _, containerName := range label.Containers {
			err := CopyToFromContainer(serverURL, user.Username, Password, containerName)
			if err != nil {
				return "access denied"
			}
		}
	}
	return "OK"
}

// Test cases follow

func TestTeamOwnedContainers(t *testing.T, serverURL string) {
	teamMap := map[string]*TestTeam{
		"security": {Name: "security", Label: "sec"},
		"dba":      {Name: "dba", Label: "db"},
		"dev":      {Name: "dev", Label: "app"},
		"it":       {Name: "it", Label: "lb"},
	}

	log.Info("Building up teams")
	for name, team := range teamMap {
		teamID, err := AddTeam(nil, serverURL, GetAdminUser(), GetAdminPassword(), team.Name)
		require.Nil(t, err, "Failed to create team %s", name)
		team.TeamID = teamID
	}

	// Labels we'll stuff on the containers
	labelMap := map[string]*TestLabel{
		"sec": {Name: "sec"},
		"db":  {Name: "db"},
		"app": {Name: "app"},
		"lb":  {Name: "lb"},
	}

	// Map the labels to access rights in a fashion that might be plausible in a real deployment
	accessList := []TestAccess{
		{Team: teamMap["security"], Label: labelMap["sec"], Role: auth.FullControl},
		{Team: teamMap["security"], Label: labelMap["db"], Role: auth.View},
		{Team: teamMap["security"], Label: labelMap["app"], Role: auth.View},
		{Team: teamMap["security"], Label: labelMap["lb"], Role: auth.View},

		{Team: teamMap["dba"], Label: labelMap["sec"], Role: auth.None},
		{Team: teamMap["dba"], Label: labelMap["db"], Role: auth.RestrictedControl},
		{Team: teamMap["dba"], Label: labelMap["app"], Role: auth.View},
		{Team: teamMap["dba"], Label: labelMap["lb"], Role: auth.None},

		{Team: teamMap["dev"], Label: labelMap["sec"], Role: auth.None},
		{Team: teamMap["dev"], Label: labelMap["db"], Role: auth.View},
		{Team: teamMap["dev"], Label: labelMap["app"], Role: auth.RestrictedControl},
		{Team: teamMap["dev"], Label: labelMap["lb"], Role: auth.None},

		{Team: teamMap["it"], Label: labelMap["sec"], Role: auth.View},
		{Team: teamMap["it"], Label: labelMap["db"], Role: auth.View},
		{Team: teamMap["it"], Label: labelMap["app"], Role: auth.View},
		{Team: teamMap["it"], Label: labelMap["lb"], Role: auth.RestrictedControl},
	}

	log.Info("Building up access lists")
	for _, a := range accessList {
		err := AddAccessList(nil, serverURL, GetAdminUser(), GetAdminPassword(), a.Team.TeamID, a.Label.Name, a.Role)
		require.Nil(t, err)
	}

	log.Info("Building up users for each team")
	// Note: This builds upon the assumption that user default permission/role isn't relevant, and will verify that during test
	for teamName, team := range teamMap {
		for userRole := auth.None; userRole <= auth.FullControl; userRole++ {
			user := &TestUser{
				Username: fmt.Sprintf("%s_user_%d", teamName, userRole),
				Teams:    []*TestTeam{team},
			}
			err := CreateNewUser(nil, serverURL, GetAdminUser(), GetAdminPassword(), user.Username, Password, false, userRole)
			require.Nil(t, err, "Failed to create user %s", user.Username)
			log.Debugf("created user: username=%s role=%d", user.Username, userRole)
			err = AddTeamMember(nil, serverURL, GetAdminUser(), GetAdminPassword(), team.TeamID, user.Username)
			require.Nil(t, err, "Failed to add user %s to team %s", user.Username, team.Name)
			team.Users = append(team.Users, user)
		}
	}

	// Now build up some scenarios and try to create containers with them
	scenarios := []TestTeamScenario{
		// Create containers across all the teams/labels
		{T: teamMap["security"], L: labelMap["sec"], O: createTeamCnt, N: "createTeamCnt", E: "OK"},
		{T: teamMap["security"], L: labelMap["db"], O: createTeamCnt, N: "createTeamCnt", E: "access denied"},
		{T: teamMap["security"], L: labelMap["app"], O: createTeamCnt, N: "createTeamCnt", E: "access denied"},
		{T: teamMap["security"], L: labelMap["lb"], O: createTeamCnt, N: "createTeamCnt", E: "access denied"},

		{T: teamMap["dba"], L: labelMap["sec"], O: createTeamCnt, N: "createTeamCnt", E: "access denied"},
		{T: teamMap["dba"], L: labelMap["db"], O: createTeamCnt, N: "createTeamCnt", E: "OK"},
		{T: teamMap["dba"], L: labelMap["app"], O: createTeamCnt, N: "createTeamCnt", E: "access denied"},
		{T: teamMap["dba"], L: labelMap["lb"], O: createTeamCnt, N: "createTeamCnt", E: "access denied"},

		{T: teamMap["dev"], L: labelMap["sec"], O: createTeamCnt, N: "createTeamCnt", E: "access denied"},
		{T: teamMap["dev"], L: labelMap["db"], O: createTeamCnt, N: "createTeamCnt", E: "access denied"},
		{T: teamMap["dev"], L: labelMap["app"], O: createTeamCnt, N: "createTeamCnt", E: "OK"},
		{T: teamMap["dev"], L: labelMap["lb"], O: createTeamCnt, N: "createTeamCnt", E: "access denied"},

		{T: teamMap["it"], L: labelMap["sec"], O: createTeamCnt, N: "createTeamCnt", E: "access denied"},
		{T: teamMap["it"], L: labelMap["db"], O: createTeamCnt, N: "createTeamCnt", E: "access denied"},
		{T: teamMap["it"], L: labelMap["app"], O: createTeamCnt, N: "createTeamCnt", E: "access denied"},
		{T: teamMap["it"], L: labelMap["lb"], O: createTeamCnt, N: "createTeamCnt", E: "OK"},

		// Try to create a privileged container (only one set will be allowed)
		{T: teamMap["security"], L: labelMap["sec"], O: createTeamCntPriv, N: "createTeamCntPriv", E: "OK"},
		{T: teamMap["security"], L: labelMap["db"], O: createTeamCntPriv, N: "createTeamCntPriv", E: "access denied"},
		{T: teamMap["security"], L: labelMap["app"], O: createTeamCntPriv, N: "createTeamCntPriv", E: "access denied"},
		{T: teamMap["security"], L: labelMap["lb"], O: createTeamCntPriv, N: "createTeamCntPriv", E: "access denied"},

		{T: teamMap["dba"], L: labelMap["sec"], O: createTeamCntPriv, N: "createTeamCntPriv", E: "access denied"},
		{T: teamMap["dba"], L: labelMap["db"], O: createTeamCntPriv, N: "createTeamCntPriv", E: "privileged not allowed"},
		{T: teamMap["dba"], L: labelMap["app"], O: createTeamCntPriv, N: "createTeamCntPriv", E: "access denied"},
		{T: teamMap["dba"], L: labelMap["lb"], O: createTeamCntPriv, N: "createTeamCntPriv", E: "access denied"},

		{T: teamMap["dev"], L: labelMap["sec"], O: createTeamCntPriv, N: "createTeamCntPriv", E: "access denied"},
		{T: teamMap["dev"], L: labelMap["db"], O: createTeamCntPriv, N: "createTeamCntPriv", E: "access denied"},
		{T: teamMap["dev"], L: labelMap["app"], O: createTeamCntPriv, N: "createTeamCntPriv", E: "privileged not allowed"},
		{T: teamMap["dev"], L: labelMap["lb"], O: createTeamCntPriv, N: "createTeamCntPriv", E: "access denied"},

		{T: teamMap["it"], L: labelMap["sec"], O: createTeamCntPriv, N: "createTeamCntPriv", E: "access denied"},
		{T: teamMap["it"], L: labelMap["db"], O: createTeamCntPriv, N: "createTeamCntPriv", E: "access denied"},
		{T: teamMap["it"], L: labelMap["app"], O: createTeamCntPriv, N: "createTeamCntPriv", E: "access denied"},
		{T: teamMap["it"], L: labelMap["lb"], O: createTeamCntPriv, N: "createTeamCntPriv", E: "privileged not allowed"},

		// Inspect containers across all teams/labels
		{T: teamMap["security"], L: labelMap["sec"], O: inspByLbl, N: "inspByLbl", E: "OK"},
		{T: teamMap["security"], L: labelMap["db"], O: inspByLbl, N: "inspByLbl", E: "OK"},
		{T: teamMap["security"], L: labelMap["app"], O: inspByLbl, N: "inspByLbl", E: "OK"},
		{T: teamMap["security"], L: labelMap["lb"], O: inspByLbl, N: "inspByLbl", E: "OK"},

		{T: teamMap["dba"], L: labelMap["sec"], O: inspByLbl, N: "inspByLbl", E: "access denied"},
		{T: teamMap["dba"], L: labelMap["db"], O: inspByLbl, N: "inspByLbl", E: "OK"},
		{T: teamMap["dba"], L: labelMap["app"], O: inspByLbl, N: "inspByLbl", E: "OK"},
		{T: teamMap["dba"], L: labelMap["lb"], O: inspByLbl, N: "inspByLbl", E: "access denied"},

		{T: teamMap["dev"], L: labelMap["sec"], O: inspByLbl, N: "inspByLbl", E: "access denied"},
		{T: teamMap["dev"], L: labelMap["db"], O: inspByLbl, N: "inspByLbl", E: "OK"},
		{T: teamMap["dev"], L: labelMap["app"], O: inspByLbl, N: "inspByLbl", E: "OK"},
		{T: teamMap["dev"], L: labelMap["lb"], O: inspByLbl, N: "inspByLbl", E: "access denied"},

		{T: teamMap["it"], L: labelMap["sec"], O: inspByLbl, N: "inspByLbl", E: "OK"},
		{T: teamMap["it"], L: labelMap["db"], O: inspByLbl, N: "inspByLbl", E: "OK"},
		{T: teamMap["it"], L: labelMap["app"], O: inspByLbl, N: "inspByLbl", E: "OK"},
		{T: teamMap["it"], L: labelMap["lb"], O: inspByLbl, N: "inspByLbl", E: "OK"},

		// Copy files to/from containers across all teams/labels
		// Only possible with Full Control over a label
		{T: teamMap["security"], L: labelMap["sec"], O: copyByLbl, N: "copyByLbl", E: "OK"},
		{T: teamMap["security"], L: labelMap["db"], O: copyByLbl, N: "copyByLbl", E: "access denied"},
		{T: teamMap["security"], L: labelMap["app"], O: copyByLbl, N: "copyByLbl", E: "access denied"},
		{T: teamMap["security"], L: labelMap["lb"], O: copyByLbl, N: "copyByLbl", E: "access denied"},

		{T: teamMap["dba"], L: labelMap["sec"], O: copyByLbl, N: "copyByLbl", E: "access denied"},
		{T: teamMap["dba"], L: labelMap["db"], O: copyByLbl, N: "copyByLbl", E: "access denied"},
		{T: teamMap["dba"], L: labelMap["app"], O: copyByLbl, N: "copyByLbl", E: "access denied"},
		{T: teamMap["dba"], L: labelMap["lb"], O: copyByLbl, N: "copyByLbl", E: "access denied"},

		{T: teamMap["dev"], L: labelMap["sec"], O: copyByLbl, N: "copyByLbl", E: "access denied"},
		{T: teamMap["dev"], L: labelMap["db"], O: copyByLbl, N: "copyByLbl", E: "access denied"},
		{T: teamMap["dev"], L: labelMap["app"], O: copyByLbl, N: "copyByLbl", E: "access denied"},
		{T: teamMap["dev"], L: labelMap["lb"], O: copyByLbl, N: "copyByLbl", E: "access denied"},

		{T: teamMap["it"], L: labelMap["sec"], O: copyByLbl, N: "copyByLbl", E: "access denied"},
		{T: teamMap["it"], L: labelMap["db"], O: copyByLbl, N: "copyByLbl", E: "access denied"},
		{T: teamMap["it"], L: labelMap["app"], O: copyByLbl, N: "copyByLbl", E: "access denied"},
		{T: teamMap["it"], L: labelMap["lb"], O: copyByLbl, N: "copyByLbl", E: "access denied"},

		// List containers across teams/labels
		{T: teamMap["security"], L: labelMap["sec"], O: listContainers, N: "listContainers", E: "OK"},
		{T: teamMap["security"], L: labelMap["db"], O: listContainers, N: "listContainers", E: "OK"},
		{T: teamMap["security"], L: labelMap["app"], O: listContainers, N: "listContainers", E: "OK"},
		{T: teamMap["security"], L: labelMap["lb"], O: listContainers, N: "listContainers", E: "OK"},

		{T: teamMap["dba"], L: labelMap["sec"], O: listContainers, N: "listContainers", E: "OK"},
		{T: teamMap["dba"], L: labelMap["db"], O: listContainers, N: "listContainers", E: "OK"},
		{T: teamMap["dba"], L: labelMap["app"], O: listContainers, N: "listContainers", E: "OK"},
		{T: teamMap["dba"], L: labelMap["lb"], O: listContainers, N: "listContainers", E: "OK"},

		{T: teamMap["dev"], L: labelMap["sec"], O: listContainers, N: "listContainers", E: "OK"},
		{T: teamMap["dev"], L: labelMap["db"], O: listContainers, N: "listContainers", E: "OK"},
		{T: teamMap["dev"], L: labelMap["app"], O: listContainers, N: "listContainers", E: "OK"},
		{T: teamMap["dev"], L: labelMap["lb"], O: listContainers, N: "listContainers", E: "OK"},

		{T: teamMap["it"], L: labelMap["sec"], O: listContainers, N: "listContainers", E: "OK"},
		{T: teamMap["it"], L: labelMap["db"], O: listContainers, N: "listContainers", E: "OK"},
		{T: teamMap["it"], L: labelMap["app"], O: listContainers, N: "listContainers", E: "OK"},
		{T: teamMap["it"], L: labelMap["lb"], O: listContainers, N: "listContainers", E: "OK"},

		// Create services across all the teams/labels
		{T: teamMap["security"], L: labelMap["sec"], O: createTeamService, N: "createTeamService", E: "OK"},
		{T: teamMap["security"], L: labelMap["db"], O: createTeamService, N: "createTeamService", E: "access denied"},
		{T: teamMap["security"], L: labelMap["app"], O: createTeamService, N: "createTeamService", E: "access denied"},
		{T: teamMap["security"], L: labelMap["lb"], O: createTeamService, N: "createTeamService", E: "access denied"},

		{T: teamMap["dba"], L: labelMap["sec"], O: createTeamService, N: "createTeamService", E: "access denied"},
		{T: teamMap["dba"], L: labelMap["db"], O: createTeamService, N: "createTeamService", E: "OK"},
		{T: teamMap["dba"], L: labelMap["app"], O: createTeamService, N: "createTeamService", E: "access denied"},
		{T: teamMap["dba"], L: labelMap["lb"], O: createTeamService, N: "createTeamService", E: "access denied"},

		{T: teamMap["dev"], L: labelMap["sec"], O: createTeamService, N: "createTeamService", E: "access denied"},
		{T: teamMap["dev"], L: labelMap["db"], O: createTeamService, N: "createTeamService", E: "access denied"},
		{T: teamMap["dev"], L: labelMap["app"], O: createTeamService, N: "createTeamService", E: "OK"},
		{T: teamMap["dev"], L: labelMap["lb"], O: createTeamService, N: "createTeamService", E: "access denied"},

		{T: teamMap["it"], L: labelMap["sec"], O: createTeamService, N: "createTeamService", E: "access denied"},
		{T: teamMap["it"], L: labelMap["db"], O: createTeamService, N: "createTeamService", E: "access denied"},
		{T: teamMap["it"], L: labelMap["app"], O: createTeamService, N: "createTeamService", E: "access denied"},
		{T: teamMap["it"], L: labelMap["lb"], O: createTeamService, N: "createTeamService", E: "OK"},

		// Inspect services across all teams/labels
		{T: teamMap["security"], L: labelMap["sec"], O: inspectTeamService, N: "inspectTeamService", E: "OK"},
		{T: teamMap["security"], L: labelMap["db"], O: inspectTeamService, N: "inspectTeamService", E: "OK"},
		{T: teamMap["security"], L: labelMap["app"], O: inspectTeamService, N: "inspectTeamService", E: "OK"},
		{T: teamMap["security"], L: labelMap["lb"], O: inspectTeamService, N: "inspectTeamService", E: "OK"},

		{T: teamMap["dba"], L: labelMap["sec"], O: inspectTeamService, N: "inspectTeamService", E: "access denied"},
		{T: teamMap["dba"], L: labelMap["db"], O: inspectTeamService, N: "inspectTeamService", E: "OK"},
		{T: teamMap["dba"], L: labelMap["app"], O: inspectTeamService, N: "inspectTeamService", E: "OK"},
		{T: teamMap["dba"], L: labelMap["lb"], O: inspectTeamService, N: "inspectTeamService", E: "access denied"},

		{T: teamMap["dev"], L: labelMap["sec"], O: inspectTeamService, N: "inspectTeamService", E: "access denied"},
		{T: teamMap["dev"], L: labelMap["db"], O: inspectTeamService, N: "inspectTeamService", E: "OK"},
		{T: teamMap["dev"], L: labelMap["app"], O: inspectTeamService, N: "inspectTeamService", E: "OK"},
		{T: teamMap["dev"], L: labelMap["lb"], O: inspectTeamService, N: "inspectTeamService", E: "access denied"},

		{T: teamMap["it"], L: labelMap["sec"], O: inspectTeamService, N: "inspectTeamService", E: "OK"},
		{T: teamMap["it"], L: labelMap["db"], O: inspectTeamService, N: "inspectTeamService", E: "OK"},
		{T: teamMap["it"], L: labelMap["app"], O: inspectTeamService, N: "inspectTeamService", E: "OK"},
		{T: teamMap["it"], L: labelMap["lb"], O: inspectTeamService, N: "inspectTeamService", E: "OK"},

		// Create networks across all the teams/labels
		{T: teamMap["security"], L: labelMap["sec"], O: createTeamNetwork, N: "createTeamNetwork", E: "OK"},
		{T: teamMap["security"], L: labelMap["db"], O: createTeamNetwork, N: "createTeamNetwork", E: "access denied"},
		{T: teamMap["security"], L: labelMap["app"], O: createTeamNetwork, N: "createTeamNetwork", E: "access denied"},
		{T: teamMap["security"], L: labelMap["lb"], O: createTeamNetwork, N: "createTeamNetwork", E: "access denied"},

		{T: teamMap["dba"], L: labelMap["sec"], O: createTeamNetwork, N: "createTeamNetwork", E: "access denied"},
		{T: teamMap["dba"], L: labelMap["db"], O: createTeamNetwork, N: "createTeamNetwork", E: "OK"},
		{T: teamMap["dba"], L: labelMap["app"], O: createTeamNetwork, N: "createTeamNetwork", E: "access denied"},
		{T: teamMap["dba"], L: labelMap["lb"], O: createTeamNetwork, N: "createTeamNetwork", E: "access denied"},

		{T: teamMap["dev"], L: labelMap["sec"], O: createTeamNetwork, N: "createTeamNetwork", E: "access denied"},
		{T: teamMap["dev"], L: labelMap["db"], O: createTeamNetwork, N: "createTeamNetwork", E: "access denied"},
		{T: teamMap["dev"], L: labelMap["app"], O: createTeamNetwork, N: "createTeamNetwork", E: "OK"},
		{T: teamMap["dev"], L: labelMap["lb"], O: createTeamNetwork, N: "createTeamNetwork", E: "access denied"},

		{T: teamMap["it"], L: labelMap["sec"], O: createTeamNetwork, N: "createTeamNetwork", E: "access denied"},
		{T: teamMap["it"], L: labelMap["db"], O: createTeamNetwork, N: "createTeamNetwork", E: "access denied"},
		{T: teamMap["it"], L: labelMap["app"], O: createTeamNetwork, N: "createTeamNetwork", E: "access denied"},
		{T: teamMap["it"], L: labelMap["lb"], O: createTeamNetwork, N: "createTeamNetwork", E: "OK"},

		// Inspect networks across all teams/labels
		{T: teamMap["security"], L: labelMap["sec"], O: inspectTeamNetwork, N: "inspectTeamNetwork", E: "OK"},
		{T: teamMap["security"], L: labelMap["db"], O: inspectTeamNetwork, N: "inspectTeamNetwork", E: "OK"},
		{T: teamMap["security"], L: labelMap["app"], O: inspectTeamNetwork, N: "inspectTeamNetwork", E: "OK"},
		{T: teamMap["security"], L: labelMap["lb"], O: inspectTeamNetwork, N: "inspectTeamNetwork", E: "OK"},

		{T: teamMap["dba"], L: labelMap["sec"], O: inspectTeamNetwork, N: "inspectTeamNetwork", E: "access denied"},
		{T: teamMap["dba"], L: labelMap["db"], O: inspectTeamNetwork, N: "inspectTeamNetwork", E: "OK"},
		{T: teamMap["dba"], L: labelMap["app"], O: inspectTeamNetwork, N: "inspectTeamNetwork", E: "OK"},
		{T: teamMap["dba"], L: labelMap["lb"], O: inspectTeamNetwork, N: "inspectTeamNetwork", E: "access denied"},

		{T: teamMap["dev"], L: labelMap["sec"], O: inspectTeamNetwork, N: "inspectTeamNetwork", E: "access denied"},
		{T: teamMap["dev"], L: labelMap["db"], O: inspectTeamNetwork, N: "inspectTeamNetwork", E: "OK"},
		{T: teamMap["dev"], L: labelMap["app"], O: inspectTeamNetwork, N: "inspectTeamNetwork", E: "OK"},
		{T: teamMap["dev"], L: labelMap["lb"], O: inspectTeamNetwork, N: "inspectTeamNetwork", E: "access denied"},

		{T: teamMap["it"], L: labelMap["sec"], O: inspectTeamNetwork, N: "inspectTeamNetwork", E: "OK"},
		{T: teamMap["it"], L: labelMap["db"], O: inspectTeamNetwork, N: "inspectTeamNetwork", E: "OK"},
		{T: teamMap["it"], L: labelMap["app"], O: inspectTeamNetwork, N: "inspectTeamNetwork", E: "OK"},
		{T: teamMap["it"], L: labelMap["lb"], O: inspectTeamNetwork, N: "inspectTeamNetwork", E: "OK"},

		// Perform a connect and disconnect across all teams/labels
		{T: teamMap["security"], L: labelMap["sec"], O: connectDisconnectContainer, N: "connectDisconnectContainer", E: "OK"},
		{T: teamMap["security"], L: labelMap["db"], O: connectDisconnectContainer, N: "connectDisconnectContainer", E: "access denied"},
		{T: teamMap["security"], L: labelMap["app"], O: connectDisconnectContainer, N: "connectDisconnectContainer", E: "access denied"},
		{T: teamMap["security"], L: labelMap["lb"], O: connectDisconnectContainer, N: "connectDisconnectContainer", E: "access denied"},

		{T: teamMap["dba"], L: labelMap["sec"], O: connectDisconnectContainer, N: "connectDisconnectContainer", E: "access denied"},
		{T: teamMap["dba"], L: labelMap["db"], O: connectDisconnectContainer, N: "connectDisconnectContainer", E: "OK"},
		{T: teamMap["dba"], L: labelMap["app"], O: connectDisconnectContainer, N: "connectDisconnectContainer", E: "access denied"},
		{T: teamMap["dba"], L: labelMap["lb"], O: connectDisconnectContainer, N: "connectDisconnectContainer", E: "access denied"},

		{T: teamMap["dev"], L: labelMap["sec"], O: connectDisconnectContainer, N: "connectDisconnectContainer", E: "access denied"},
		{T: teamMap["dev"], L: labelMap["db"], O: connectDisconnectContainer, N: "connectDisconnectContainer", E: "access denied"},
		{T: teamMap["dev"], L: labelMap["app"], O: connectDisconnectContainer, N: "connectDisconnectContainer", E: "OK"},
		{T: teamMap["dev"], L: labelMap["lb"], O: connectDisconnectContainer, N: "connectDisconnectContainer", E: "access denied"},

		{T: teamMap["it"], L: labelMap["sec"], O: connectDisconnectContainer, N: "connectDisconnectContainer", E: "access denied"},
		{T: teamMap["it"], L: labelMap["db"], O: connectDisconnectContainer, N: "connectDisconnectContainer", E: "access denied"},
		{T: teamMap["it"], L: labelMap["app"], O: connectDisconnectContainer, N: "connectDisconnectContainer", E: "access denied"},
		{T: teamMap["it"], L: labelMap["lb"], O: connectDisconnectContainer, N: "connectDisconnectContainer", E: "OK"},

		// Create a container with the team's highest label and attach it at creation time to a target network
		{T: teamMap["security"], L: labelMap["sec"], O: createNetworkContainer, N: "createNetworkContainer", E: "OK"},
		{T: teamMap["security"], L: labelMap["db"], O: createNetworkContainer, N: "createNetworkContainer", E: "access denied"},
		{T: teamMap["security"], L: labelMap["app"], O: createNetworkContainer, N: "createNetworkContainer", E: "access denied"},
		{T: teamMap["security"], L: labelMap["lb"], O: createNetworkContainer, N: "createNetworkContainer", E: "access denied"},

		{T: teamMap["dba"], L: labelMap["sec"], O: createNetworkContainer, N: "createNetworkContainer", E: "access denied"},
		{T: teamMap["dba"], L: labelMap["db"], O: createNetworkContainer, N: "createNetworkContainer", E: "OK"},
		{T: teamMap["dba"], L: labelMap["app"], O: createNetworkContainer, N: "createNetworkContainer", E: "access denied"},
		{T: teamMap["dba"], L: labelMap["lb"], O: createNetworkContainer, N: "createNetworkContainer", E: "access denied"},

		{T: teamMap["dev"], L: labelMap["sec"], O: createNetworkContainer, N: "createNetworkContainer", E: "access denied"},
		{T: teamMap["dev"], L: labelMap["db"], O: createNetworkContainer, N: "createNetworkContainer", E: "access denied"},
		{T: teamMap["dev"], L: labelMap["app"], O: createNetworkContainer, N: "createNetworkContainer", E: "OK"},
		{T: teamMap["dev"], L: labelMap["lb"], O: createNetworkContainer, N: "createNetworkContainer", E: "access denied"},

		{T: teamMap["it"], L: labelMap["sec"], O: createNetworkContainer, N: "createNetworkContainer", E: "access denied"},
		{T: teamMap["it"], L: labelMap["db"], O: createNetworkContainer, N: "createNetworkContainer", E: "access denied"},
		{T: teamMap["it"], L: labelMap["app"], O: createNetworkContainer, N: "createNetworkContainer", E: "access denied"},
		{T: teamMap["it"], L: labelMap["lb"], O: createNetworkContainer, N: "createNetworkContainer", E: "OK"},
	}

	for _, s := range scenarios {
		res := s.O(serverURL, s.T, s.L)
		if strings.Contains(res, s.E) {
			log.Infof("PASS %s:%s:%s found expected match in output: %s", s.T.Name, s.L.Name, s.N, res)
		} else {
			assert.Fail(t, fmt.Sprintf("Failed %s:%s:%s - expected %s found %s", s.T.Name, s.L.Name, s.N, s.E, res))
		}
	}
}
