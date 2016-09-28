package util

import (
	"net/http"
	"path"

	"github.com/docker/dhe-deploy/adminserver/api/common/errors"
	"github.com/docker/dhe-deploy/integration/apiclient"
	"github.com/samalba/dockerclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type UserPermissionsOnRepo struct {
	AccessLevel      string
	NSAccessLevel    string
	TeamsVisible     bool
	Push             bool
	Pull             bool
	View             bool
	DeleteTags       bool
	EditDescription  bool
	MakePublic       bool
	ManageTeamAccess bool
	ExplicitlyShared bool
}

func (u *Util) CheckUserPermissionsOnRepo(testDesc string, perms UserPermissionsOnRepo, namespace, repoName, localImg1, localImg2, teamName, team2Name, username, password string) {
	userAuthConfig := &dockerclient.AuthConfig{
		Username: username,
		Password: password,
		Email:    "a@a.a",
	}
	require.Nil(u.T(), u.API.Login(username, password))

	// push should work according to the requirements in the table
	imageTag := "latest"
	untag := u.TagImageWithChecks(path.Join(u.Config.DTRHost, namespace), repoName, imageTag, localImg1)
	err := u.Docker.PushImage(path.Join(u.Config.DTRHost, namespace, repoName), imageTag, userAuthConfig)
	if perms.Push {
		assert.Nil(u.T(), err, testDesc)
	} else {
		assert.NotNil(u.T(), err, testDesc)
	}

	// even if the previous push didn't succeed, we push as admin to make sure the tag exists for the next tests
	if err != nil {
		err = u.Docker.PushImage(path.Join(u.Config.DTRHost, namespace, repoName), imageTag, u.Config.AdminAuthConfig)
		require.Nil(u.T(), err)
	}

	// if we can push, we should also be able to re-push aka "update" a tag
	untag()
	untag2 := u.TagImageWithChecks(path.Join(u.Config.DTRHost, namespace), repoName, imageTag, localImg2)
	err = u.Docker.PushImage(path.Join(u.Config.DTRHost, namespace, repoName), imageTag, userAuthConfig)
	if perms.Push {
		assert.Nil(u.T(), err, testDesc)
	} else {
		assert.NotNil(u.T(), err, testDesc)
	}
	untag2()

	// check if we list the correct access level
	rua, err := u.API.GetUserRepositoryAccess(username, namespace, repoName)
	if perms.View {
		if assert.Nil(u.T(), err, testDesc) {
			if assert.NotNil(u.T(), rua, testDesc) {
				assert.Equal(u.T(), repoName, rua.Repository.Name, testDesc)
				assert.Equal(u.T(), perms.AccessLevel, rua.AccessLevel, testDesc)
			}
		}
	} else {
		u.AssertErrorCodesWithMsg(err, http.StatusNotFound, testDesc, errors.ErrorCodeNoSuchRepository)
	}
	// NOTE: this isn't a feature in 2.0
	// // check if an admin can see the user's access level
	// require.Nil(u.T(), u.API.Login(u.API.EnziSession()))
	// rua, err = u.API.GetUserRepositoryAccess(username, namespace, repoName)
	// assert.Nil(u.T(), err, testDesc)
	// assert.Equal(u.T(), repoName, rua.Repository.Name, testDesc)
	// assert.Equal(u.T(), perms.AccessLevel, rua.UserAccess.AccessLevel, testDesc)
	// u.API.Login(username, password)

	if perms.NSAccessLevel != "" {
		// // check if we list the correct namespace access level for the user
		// rnua, err := u.API.GetUserRepositoryNamespaceAccess(username, namespace)
		// assert.Nil(u.T(), err, testDesc)
		// assert.Equal(u.T(), namespace, rnua.Namespace.Name, testDesc)
		// assert.Equal(u.T(), perms.NSAccessLevel, rnua.AccessLevel, testDesc)

		// // check if an admin can see the user's namespace access level for the user
		// u.API.Login(u.Config.AdminUsername, u.Config.AdminPassword)
		// rnua, err = u.API.GetUserRepositoryNamespaceAccess(username, namespace)
		// assert.Nil(u.T(), err, testDesc)
		// assert.Equal(u.T(), namespace, rnua.Namespace.Name, testDesc)
		// assert.Equal(u.T(), perms.NSAccessLevel, rnua.AccessLevel, testDesc)
		// u.API.Login(username, password)

		// check if we list the correct namespace access levels list
		ns, teamAccess, err := u.API.ListRepositoryNamespaceTeamAccess(namespace)
		if !perms.TeamsVisible {
			assert.NotNil(u.T(), err, testDesc)
		} else {
			if assert.Nil(u.T(), err, testDesc) {
				assert.Equal(u.T(), namespace, string(*ns), testDesc)
				if perms.NSAccessLevel == "" {
					assert.Len(u.T(), teamAccess, 0, testDesc)
				} else {
					if assert.Len(u.T(), teamAccess, 1, testDesc) {
						assert.Equal(u.T(), perms.NSAccessLevel, teamAccess[0].AccessLevel, testDesc)
					}
				}
			}
		}

		// check if we list the correct namespace access levels for the team
		rnta, err := u.API.GetRepositoryNamespaceTeamAccess(namespace, teamName)
		if !perms.TeamsVisible {
			assert.NotNil(u.T(), err, testDesc)
		} else {
			if assert.Nil(u.T(), err, testDesc) {
				assert.Equal(u.T(), perms.AccessLevel, rnta.AccessLevel, testDesc)
				assert.Equal(u.T(), namespace, string(rnta.Namespace), testDesc)
			}
		}
	}

	// viewing should work according to the requirements in the table
	repo, err := u.API.GetRepository(namespace, repoName)
	if perms.View {
		assert.Nil(u.T(), err, testDesc)
		if assert.NotNil(u.T(), repo) {
			assert.Equal(u.T(), repoName, repo.Name, testDesc)
		}
	} else {
		u.AssertErrorCodesWithMsg(err, http.StatusNotFound, testDesc, errors.ErrorCodeNoSuchRepository)
	}

	// viewing also implies that listing all repos, "shared with me" repos or the org's repos should include/not include the repo
	repos, err := u.API.ListRepositories(namespace)
	assert.Nil(u.T(), err, testDesc)
	if perms.View {
		if assert.Len(u.T(), repos, 1, testDesc) {
			assert.Equal(u.T(), repoName, repos[0].Name, testDesc)
		}
	} else {
		assert.Len(u.T(), repos, 0, testDesc)
	}
	repos, err = u.API.ListAllRepositories()
	assert.Nil(u.T(), err, testDesc)
	if perms.View {
		if assert.Len(u.T(), repos, 1, testDesc) {
			assert.Equal(u.T(), repoName, repos[0].Name, testDesc)
		}
	} else {
		assert.Len(u.T(), repos, 0, testDesc)
	}
	// ListSharedRepositories is gone
	// repos, err = u.API.ListSharedRepositories(username)
	// assert.Nil(u.T(), err, testDesc)
	// if perms.View && perms.ExplicitlyShared {
	// 	if assert.Len(u.T(), repos, 1, testDesc) {
	// 		assert.Equal(u.T(), repoName, repos[0].Name, testDesc)
	// 	}
	// } else {
	// 	assert.Len(u.T(), repos, 0, testDesc)
	// }

	// pull should work according to the requirements in the table
	err = u.Docker.PullImage(path.Join(u.Config.DTRHost, namespace, repoName)+":"+imageTag, userAuthConfig)
	if perms.Pull {
		assert.Nil(u.T(), err, testDesc)
	} else {
		assert.NotNil(u.T(), err, testDesc)
	}

	newShort := "changed"
	newLong := "changed"
	public := "public"

	// editing the description should work according to the requirements in the table
	repo, err = u.API.UpdateRepository(namespace, repoName, apiclient.RepositoryUpdateForm{
		ShortDescription: &newShort,
		LongDescription:  &newLong,
	})
	if perms.EditDescription {
		if assert.Nil(u.T(), err, testDesc) {
			assert.Equal(u.T(), newShort, repo.ShortDescription, testDesc)
			assert.Equal(u.T(), newLong, repo.LongDescription, testDesc)
		}
	} else {
		if perms.View {
			u.AssertErrorCodesWithMsg(err, http.StatusForbidden, testDesc, errors.ErrorCodeNotAuthorized)
		} else {
			u.AssertErrorCodesWithMsg(err, http.StatusNotFound, testDesc, errors.ErrorCodeNoSuchRepository)
		}
	}

	// editing the visibility should work according to the requirements in the table
	repo, err = u.API.UpdateRepository(namespace, repoName, apiclient.RepositoryUpdateForm{
		Visibility: &public,
	})
	if perms.MakePublic {
		if assert.Nil(u.T(), err, testDesc) {
			assert.Equal(u.T(), public, repo.Visibility, testDesc)
		}
	} else {
		if perms.View {
			u.AssertErrorCodesWithMsg(err, http.StatusForbidden, testDesc, errors.ErrorCodeNotAuthorized)
		} else {
			u.AssertErrorCodesWithMsg(err, http.StatusNotFound, testDesc, errors.ErrorCodeNoSuchRepository)
		}
	}

	// deleting should work according to the requirements in the table
	err = u.API.DeleteTag(namespace, repoName, imageTag)
	if perms.DeleteTags {
		assert.Nil(u.T(), err, testDesc)
	} else {
		if perms.View {
			u.AssertErrorCodesWithMsg(err, http.StatusForbidden, testDesc, errors.ErrorCodeNotAuthorized)
		} else {
			u.AssertErrorCodesWithMsg(err, http.StatusNotFound, testDesc, errors.ErrorCodeNoSuchRepository)
		}
	}

	// giving another team permissions
	if team2Name != "" {
		_, err = u.API.SetRepositoryTeamAccess(namespace, repoName, team2Name, "read-only")
		if perms.ManageTeamAccess {
			assert.Nil(u.T(), err, testDesc)
		} else {
			if perms.View {
				u.AssertErrorCodesWithMsg(err, http.StatusForbidden, testDesc, errors.ErrorCodeNotAuthorized)
			} else {
				u.AssertErrorCodesWithMsg(err, http.StatusNotFound, testDesc, errors.ErrorCodeNoSuchRepository)
			}
		}
		u.API.RevokeRepositoryTeamAccess(namespace, repoName, team2Name)
	}
}
