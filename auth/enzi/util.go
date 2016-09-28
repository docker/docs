package enzi

import (
	"github.com/docker/orca/auth"
	"github.com/docker/orca/enzi/api/client"
	"github.com/docker/orca/enzi/api/responses"
)

func (a *Authenticator) getSession(creds client.RequestAuthenticator) *client.Session {
	return client.New(a.httpClient, a.providerAddr, "enzi", creds)
}

func (a *Authenticator) populateUserFields(creds client.RequestAuthenticator, userResp *responses.Account) (*auth.Account, error) {
	user := &auth.Account{
		ID:        userResp.ID,
		FirstName: userResp.FullName, // FIXME
		LastName:  "",                // FIXME
		Username:  userResp.Name,
	}

	if userResp.IsAdmin != nil {
		user.Admin = *userResp.IsAdmin
	}
	if userResp.IsActive != nil {
		user.Disabled = !*userResp.IsActive
	}

	// First, get the teams.
	teams, err := a.listUserTeams(creds, user.Username)
	if err != nil {
		return nil, err
	}

	user.ManagedTeams = make([]string, len(teams))
	for i, team := range teams {
		user.ManagedTeams[i] = team.Id
	}

	// Then get the user's role.
	user.Role, err = a.getUserRole(user.Username)
	if err != nil {
		return nil, err
	}

	// Finally, get the user's public keys.
	user.PublicKeys, err = a.listUserPublicKeys(user.Username)
	if err != nil {
		return nil, err
	}

	return user, nil
}
