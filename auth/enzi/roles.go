package enzi

import (
	"encoding/json"
	"fmt"
	"path"

	libkv "github.com/docker/libkv/store"
	"github.com/docker/orca/auth"
	"github.com/docker/orca/utils"
)

// Users' Default Roles are stored in the following directory
// structure:
//     orca/v1/auth/
//         user_roles/
//             {username} -> JSON({role})
//             {username} -> JSON({role})
//             ...
const kvStoreUserRolesPrefix = "orca/v1/auth/user_roles"

func kvStoreUserRolePath(username string) string {
	return path.Join(kvStoreUserRolesPrefix, username)
}

type userRoleData struct {
	Role auth.Role `json:"role"`
}

func (a *Authenticator) SetUserRole(username string, role auth.Role) error {
	roleData, err := json.Marshal(userRoleData{Role: role})
	if err != nil {
		return fmt.Errorf("unable to encode user default role to JSON: %s", err)
	}

	if err := a.kvStore.Put(kvStoreUserRolePath(username), roleData, nil); err != nil {
		return utils.MaybeWrapEtcdClusterErr(err)
	}

	return nil
}

func (a *Authenticator) getUserRole(username string) (auth.Role, error) {
	key := kvStoreUserRolePath(username)

	kvPair, err := a.kvStore.Get(key)
	if err != nil {
		if err == libkv.ErrKeyNotFound {
			// A user-specific role is not specified so fall back
			// to the global default from our config.
			return a.userDefaultRole, nil
		}

		// FIXME: Wrap this error? If it's an etcd client error then
		// we will lose details if we wrap with fmt.Errorf().
		return auth.None, utils.MaybeWrapEtcdClusterErr(err)
	}

	var roleData userRoleData
	if err := json.Unmarshal(kvPair.Value, &roleData); err != nil {
		return auth.None, fmt.Errorf("unable to decode default role JSON data: %s", err)
	}

	return roleData.Role, nil
}
