package util

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/docker/dhe-deploy"
	"github.com/docker/dhe-deploy/shared/dtrutil"

	"github.com/docker/orca/enzi/api/forms"
	enzischema "github.com/docker/orca/enzi/schema"
	enziworker "github.com/docker/orca/enzi/worker"
	goldap "github.com/go-ldap/ldap"
	"github.com/stretchr/testify/require"
)

const (
	ldapImage                 = "brianbland/ldap-base"
	containerName             = "ldap_1"
	ldapPort                  = "3000"
	ldapAdminDN               = "cn=Manager,dc=example,dc=com"
	ldapAdminPassword         = "welcome0"
	adminLDAPGroupDN          = "cn=dtradmins,dc=example,dc=com"
	adminGroupMemberAttribute = "member"
	searchRootDN              = "dc=example,dc=com"
	userSearchFilter          = "objectClass=person"
	userLoginAttrName         = "cn"
)

func (u *Util) Docker0IP() string {
	return u.Execute(fmt.Sprintf("sudo docker exec ldap_1 bash -c \"route -n | tail -n 2 | head -n 1 | awk '{print \\$2}'\""), false)
}

func (u *Util) GetLDAPAuthSettings() forms.LDAPSettings {
	return forms.LDAPSettings{
		RecoveryAdminUsername: u.IntegrationFramework.Config.AdminUsername,
		RecoveryAdminPassword: &[]string{u.IntegrationFramework.Config.AdminPassword}[0],
		ServerURL:             "ldap://" + u.Docker0IP() + ":" + ldapPort,
		StartTLS:              false,
		RootCerts:             "",   // TODO use root certs!
		TLSSkipVerify:         true, // TODO stoppit??
		ReaderDN:              ldapAdminDN,
		ReaderPassword:        ldapAdminPassword,
		UserSearchConfigs: []forms.UserSearchOpts{{
			BaseDN:       searchRootDN,
			ScopeSubtree: false,
			UsernameAttr: userLoginAttrName,
			Filter:       userSearchFilter,
		}},
		AdminSyncOpts: forms.MemberSyncOpts{
			SelectGroupMembers: true,
			GroupDN:            adminLDAPGroupDN,
			GroupMemberAttr:    adminGroupMemberAttribute,
		},
		SyncSchedule: "@hourly",
	}
}
func (u *Util) GetLDAPExternalUrl() string {
	return "ldap://" + strings.Split(u.Config.DTRHost, ":")[0] + ":" + ldapPort
}

func (u *Util) RestartLDAPContainer() {
	AppendDockerIgnorableLoggedErrors([]string{
		"Handler for DELETE /v1.21/containers/ldap_1 return",
		"no such id: ldap_1",
		"No such container: ldap_1",
	})

	u.Execute(fmt.Sprintf("sudo docker pull %s", ldapImage), false)

	// make sure we have clean state of the ldap server
	u.Execute(fmt.Sprintf("sudo docker rm -f %s", containerName), true)

	dockerArgs := fmt.Sprintf("--net %s", deploy.BridgeNetworkName)
	// TODO: restrict the IP to the IP of the admin server container
	u.Execute(fmt.Sprintf("sudo docker run -d %s -p 0.0.0.0:%s:389 --name %s %s", dockerArgs, ldapPort, containerName, ldapImage), false)

	require.Nil(u.T(), dtrutil.Poll(time.Second, u.Config.RetryAttempts, func() error {
		// try to connect to the ldap port, once we succeed, end the polling
		conn, err := GetLDAPConn(u.GetLDAPExternalUrl(), false)
		if err != nil {
			return err
		}
		defer conn.Close()

		err = conn.Bind(ldapAdminDN, ldapAdminPassword)
		if err != nil {
			return err
		}
		return nil
	}))
}

func (u *Util) GetBoundLDAPConn() *goldap.Conn {
	conn, err := GetLDAPConn(u.GetLDAPExternalUrl(), false)
	require.Nil(u.T(), err)
	err = conn.Bind(ldapAdminDN, ldapAdminPassword)
	require.Nil(u.T(), err)
	return conn
}

// TODO I have no idea how to trigger an LDAP sync
func (u *Util) LDAPSync() {
	u.API.Login(u.IntegrationFramework.Config.AdminUsername, u.IntegrationFramework.Config.AdminPassword)

	jobResp, err := u.API.EnziSession().CreateJob(forms.JobSubmission{Action: enziworker.ActionLdapSync})
	require.Nil(u.T(), err)
	require.NotEmpty(u.T(), jobResp.ID)

	// Wait for the sync job to finish.
	for {
		// Give it a little more time to run. If the job is still in
		// "waiting" status, this should give it enough time to get
		// picked up by a worker.
		time.Sleep(5 * time.Second)

		jobResp, err = u.API.EnziSession().GetJob(jobResp.ID)
		require.Nil(u.T(), err)

		if jobResp.Status == enzischema.JobStatusDone {
			break
		}

		jobLogsReadCloser, err := u.API.EnziSession().GetJobLogs(jobResp.ID)
		require.Nil(u.T(), err)

		jobLogs, err := ioutil.ReadAll(jobLogsReadCloser)
		require.Nil(u.T(), err)

		require.Equal(u.T(), enzischema.JobStatusRunning, jobResp.Status, string(jobLogs))
		jobLogsReadCloser.Close()
	}
}

func (u *Util) CreateUserInLDAPServer(username, password string) error {
	ldapConn := u.GetBoundLDAPConn()
	defer ldapConn.Close()

	// Make adminUpper and adminuser admins according to the ldap server
	addRequest := goldap.NewAddRequest(u.GetDN(username))
	addRequest.Attribute("objectclass", []string{"person"})
	addRequest.Attribute("cn", []string{username})
	addRequest.Attribute("sn", []string{username})
	addRequest.Attribute("userPassword", []string{password})
	return ldapConn.Add(addRequest)
}

func (u *Util) CreateTeamInLDAPServer(teamname string) error {
	ldapConn := u.GetBoundLDAPConn()
	defer ldapConn.Close()

	// Make adminUpper and adminuser admins according to the ldap server
	addRequest := goldap.NewAddRequest(u.GetDN(teamname))
	addRequest.Attribute("objectclass", []string{"groupofnames"})
	addRequest.Attribute("cn", []string{teamname})
	return ldapConn.Add(addRequest)
}

func (u *Util) ChangeUserLDAPPassword(username, oldPassword, newPassword string) error {
	ldapConn := u.GetBoundLDAPConn()
	defer ldapConn.Close()

	modifyRequest := goldap.NewPasswordModifyRequest(u.GetDN(username), oldPassword, newPassword)
	_, err := ldapConn.PasswordModify(modifyRequest)
	return err
}

func (u *Util) DeleteLDAPUser(username string) {
	u.DeleteLDAPEntry(u.GetDN(username))
}

func (u *Util) GetDN(username string) string {
	return fmt.Sprintf("cn=%s,dc=example,dc=com", strings.ToLower(username))
}

func (u *Util) DeleteLDAPEntry(dn string) {
	ldapConn := u.GetBoundLDAPConn()
	defer ldapConn.Close()

	deleteRequest := goldap.NewDelRequest(dn, []goldap.Control{})
	ldapConn.Del(deleteRequest)
}
