package util

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func (u *Util) CreateRepoWithChecks(namespace, reponame, shortDesc, longDesc, visibility string) func() {
	retFunc := func() {}
	if repo, err := u.API.CreateRepository(namespace, reponame, shortDesc, longDesc, visibility); err != nil {
		u.T().Fatalf("Failed to create repo %s/%s as %v: %s", namespace, reponame, err, u.API.Username())
	} else {
		retFunc = func() {
			u.API.Login(u.Config.AdminUsername, u.Config.AdminPassword)
			assert.Nil(u.T(), u.API.DeleteRepository(namespace, reponame))
		}
		require.NotNil(u.T(), repo)
		assert.NotEmpty(u.T(), repo.ID)
		assert.Equal(u.T(), namespace, repo.Namespace)
		assert.Equal(u.T(), reponame, repo.Name)
		assert.Equal(u.T(), shortDesc, repo.ShortDescription)
		assert.Equal(u.T(), longDesc, repo.LongDescription)
		assert.Equal(u.T(), visibility, repo.Visibility)
	}
	return retFunc
}
