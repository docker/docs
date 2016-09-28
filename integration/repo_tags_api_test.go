package integration

import (
	"net/http"
	"path"
	"testing"

	"github.com/docker/dhe-deploy/adminserver/api/common/errors"
	"github.com/docker/dhe-deploy/integration/framework"
	"github.com/docker/dhe-deploy/integration/util"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// RepoTagsAPITestSuite is an integration test suite for the DTR API.
type RepoTagsAPITestSuite struct {
	suite.Suite
	*framework.IntegrationFramework
	u *util.Util

	adminAccountId string
	imageName      string
	setOldAuth     func()
	//deletRepo  func()
}

// SetupSuite handles setting up the API test suite by initializing auth in
// managed mode with our initial admin user.
func (suite *RepoTagsAPITestSuite) SetupSuite() {
	suite.IntegrationFramework, suite.u = setupFramework(suite)
	suite.setOldAuth = suite.u.SwitchAuth()
	suite.imageName = "tianon/true"

	account, err := suite.API.GetAccount(suite.Config.AdminUsername)
	require.Nil(suite.T(), err, "Getting account returned error: %s", err)
	suite.adminAccountId = account.ID

	err = suite.Docker.PullImage(suite.imageName, nil)
	defer func() {
		if err != nil {
			suite.setOldAuth()
		}
	}()
	require.Nil(suite.T(), err)
}

// TearDownSuite handles tearing down the API test suite by restoring the
// original auth settings.
func (suite *RepoTagsAPITestSuite) TearDownSuite() {
	suite.setOldAuth()
}

func (suite *RepoTagsAPITestSuite) SetupTest() {
	util.WipeDTRIgnorableLoggedErrors()
	util.WipeDockerIgnorableLoggedErrors()
}

func (suite *RepoTagsAPITestSuite) TearDownTest() {
	suite.u.TestLogs()
}

type Catalog struct {
	Repositories []string `json:"repositories"`
}

// YOLO (this broke when we added etcd, TODO: figure out why and fix it before we release HA)
//func (suite *RepoTagsAPITestSuite) TestOldTagsMigration() {
//	var err error
//	// DTR doesn't let you push a repo with a non-namespaced name, so we spin up our own registry with no auth to do it
//
//	// make sure we have the distribution image
//	suite.u.Execute(fmt.Sprintf("sudo docker pull docker/trusted-registry-distribution:v2.2.0"), false)
//	suite.u.Execute(fmt.Sprintf("sudo docker pull quay.io/coreos/etcd:v2.2.4"), false)
//
//	// create new storage config with the same contents as the old one without the auth section and with tls enabled
//	newStorageYmlFilename := "/tmp/storage2.yml"
//	newStorageYmlMountFilename := "/config/storage2.yml"
//
//	certPemFilename := "/tmp/server.pem"
//	certPemMountFilename := "/ssl/server.pem"
//	session, err := suite.SSH.ConnectToRemoteHost()
//	require.Nil(suite.T(), err)
//	pemReader := strings.NewReader(serverPem)
//	err = scp.Copy(int64(len(serverPem)), os.FileMode(0666), path.Base(certPemFilename), pemReader, certPemFilename, session)
//	require.Nil(suite.T(), err)
//
//	// we can't use the actual hubconfig because we are not necessarily on the same network as dtr, so we don't have access to etcd
//	oldYml := suite.u.Execute(fmt.Sprintf("sudo docker run --net %s --entrypoint /etcdctl quay.io/coreos/etcd:v2.2.4 --endpoint %s get %s", deploy.NetworkMode, strings.Join(deploy.EtcdUrls, ","), path.Join(deploy.EtcdPath, deploy.RegistryConfigFilename)), false)
//	storageConfig := &configuration.Configuration{}
//	err = yaml.Unmarshal([]byte(oldYml), storageConfig)
//	require.Nil(suite.T(), err)
//	storageConfig.Auth = nil
//	storageConfig.HTTP.TLS.Certificate = "/ssl/server.pem"
//	storageConfig.HTTP.TLS.Key = "/ssl/server.pem"
//
//	// copy the yaml to where we'll mount it from
//	newYml, err := yaml.Marshal(storageConfig)
//	require.Nil(suite.T(), err)
//	session, err = suite.SSH.ConnectToRemoteHost()
//	require.Nil(suite.T(), err)
//	ymlReader := bytes.NewReader(newYml)
//	err = scp.Copy(int64(len(newYml)), os.FileMode(0666), path.Base(newStorageYmlFilename), ymlReader, newStorageYmlFilename, session)
//	require.Nil(suite.T(), err)
//
//	defer suite.u.Execute(fmt.Sprintf("sudo rm %s", newStorageYmlFilename), false)
//
//	// make sure the custom registry isn't running already
//	defer suite.u.Execute(fmt.Sprintf("sudo docker rm -f noauth_registry"), true)
//
//	// run the custom registry
//	// XXX: we assume that the registry config is the default filesystem config
//	// we reuse dtr's certs because it's easier than generating our own
//	opts := fmt.Sprintf("-v %s:%s -v %s:%s -v %s:%s -v %s:%s -e %s=%s -p 0.0.0.0:1337:5000", path.Join(deploy.ImageStorageRootPath, "/local"), "/storage", newStorageYmlFilename, newStorageYmlMountFilename, deploy.GeneratedConfigsPath, deploy.StorageContainerPEMDir, certPemFilename, certPemMountFilename, "REGISTRY_STORAGE_FILESYSTEM_ROOTDIRECTORY", "/storage")
//	dockerArgs := fmt.Sprintf("--net %s", deploy.NetworkMode)
//	suite.u.Execute(fmt.Sprintf("sudo docker run -d %s --name noauth_registry %s docker/trusted-registry-distribution:v2.2.0 %s", dockerArgs, opts, newStorageYmlMountFilename), false)
//
//	// defer delete temporary registry container
//	defer suite.u.Execute(fmt.Sprintf("sudo docker rm -f noauth_registry"), false)
//
//	// tag the image to push to the temporary registry without a namespace
//	err = suite.Docker.TagImage(suite.imageName, path.Join(suite.Config.DTRHost+":1337", "oldname"), "latest", true)
//	require.Nil(suite.T(), err)
//	defer func() {
//		_, err := suite.Docker.RemoveImage(path.Join(suite.Config.DTRHost+":1337", "oldname:latest"), true)
//		assert.Nil(suite.T(), err)
//	}()
//
//	err = deploy.Poll(time.Second, 10, func() error {
//		// push to the temporary registry to hack it into DTR
//		return suite.Docker.PushImage(path.Join(suite.Config.DTRHost+":1337", "oldname"), "latest", suite.Config.AdminAuthConfig)
//	})
//	require.Nil(suite.T(), err)
//
//	//TODO: defer a delete from the registry when we have a real registry client
//
//	// using the DTR registry, make sure the catalog api works and lists the "oldname" repo
//	// TODO: use a real registry client
//	err = deploy.Poll(time.Second, 10, func() error {
//		catalogStr := suite.u.Execute(fmt.Sprintf("token=$(curl -u %s:%s 'https://%s/auth/token?service=%s&scope=registry:catalog:*' -k | grep -o ':\"[^\"]*' | cut -c 3-) && curl -H \"Authorization: Bearer $token\" 'https://%s/v2/_catalog' -k", suite.Config.AdminUsername, suite.Config.AdminPassword, suite.Config.DTRHost, suite.Config.DTRHost, suite.Config.DTRHost), false)
//		var catalog Catalog
//		err := json.Unmarshal([]byte(catalogStr), &catalog)
//		if err != nil {
//			return err
//		}
//		errorChecker := util.NewErrorChecker()
//		assert.Contains(errorChecker, catalog.Repositories, "oldname")
//		return errorChecker.Errors()
//	})
//	require.Nil(suite.T(), err)
//
//	// using the DTR registry, make sure we can't push to the old tag
//	err = suite.Docker.TagImage(suite.imageName, path.Join(suite.Config.DTRHost, "oldname"), "latest", true)
//	require.Nil(suite.T(), err)
//	defer func() {
//		_, err = suite.Docker.RemoveImage(path.Join(suite.Config.DTRHost, "oldname:latest"), true)
//		assert.Nil(suite.T(), err)
//	}()
//
//	// make sure we can't push as admin
//	err = suite.Docker.PushImage(path.Join(suite.Config.DTRHost, "oldname"), "latest", suite.Config.AdminAuthConfig)
//	assert.NotNil(suite.T(), err)
//
//	// make sure we can pull the old tag
//	err = suite.Docker.PullImage(path.Join(suite.Config.DTRHost, "oldname:latest"), suite.Config.AdminAuthConfig)
//	assert.Nil(suite.T(), err)
//
//	// make sure we can't pull or push as non-admin, also can't user the catalog api as non-admin
//	user, rmuser := suite.u.CreateActivateRandomUser()
//	userAuthConfig := &dockerclient.AuthConfig{
//		Username: user.Name,
//		Password: user.Password,
//		Email:    "a@a.a",
//	}
//	defer rmuser()
//
//	catalogStr := suite.u.Execute(fmt.Sprintf("token=$(curl -u '%s:%s' 'https://%s/auth/token?service=%s&scope=registry:catalog:*' -k | grep -o ':\"[^\"]*' | cut -c 3-) && curl -H \"Authorization: Bearer $token\" 'https://%s/v2/_catalog' -k", user.Name, user.Password, suite.Config.DTRHost, suite.Config.DTRHost, suite.Config.DTRHost), false)
//	assert.Contains(suite.T(), catalogStr, "UNAUTHORIZED")
//	err = suite.Docker.PushImage(path.Join(suite.Config.DTRHost, "oldname"), "latest", userAuthConfig)
//	assert.NotNil(suite.T(), err)
//	err = suite.Docker.PullImage(path.Join(suite.Config.DTRHost, "oldname:latest"), userAuthConfig)
//	assert.NotNil(suite.T(), err)
//
//	util.AppendDockerIgnorableLoggedErrors([]string{
//		"no such id: noauth_registry",        // docker 1.9
//		"No such container: noauth_registry", // docker 1.10
//	})
//}

func (suite *RepoTagsAPITestSuite) TestListTags() {
	var err error
	testRepo := "testrepo"

	// create a repo under the admin namespace using the admin
	suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)
	defer suite.u.CreateRepoWithChecks(suite.Config.AdminUsername, testRepo, "blah", "blah blah", "public")()

	// make sure it has no tags
	response, err := suite.API.GetRepositoryTags(suite.Config.AdminUsername, testRepo)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 0, len(response))

	// make sure we get an error for non-existent repos
	_, err = suite.API.GetRepositoryTags(suite.Config.AdminUsername, "noooooot")
	suite.u.AssertErrorCodes(err, http.StatusNotFound, errors.ErrorCodeNoSuchRepository)

	// make sure we get an error for non-existent user
	_, err = suite.API.GetRepositoryTags("wat", "not")
	suite.u.AssertErrorCodes(err, http.StatusNotFound, errors.ErrorCodeNoSuchAccount)

	// push an image and make sure it fails if the repo doesn't exist
	tag := "1337"
	defer suite.u.TagImageWithChecks(suite.Config.DTRHost+"/"+suite.Config.AdminUsername, "nonexistent", tag, suite.imageName)()
	err = suite.Docker.PushImage(path.Join(suite.Config.DTRHost, suite.Config.AdminUsername, "nonexistent"), tag, suite.Config.AdminAuthConfig)
	require.NotNil(suite.T(), err)

	// push an image and see if it appears in the list
	defer suite.u.TagImageWithChecks(suite.Config.DTRHost+"/"+suite.Config.AdminUsername, testRepo, tag, suite.imageName)()
	func() {
		defer suite.u.PushImageWithChecks(suite.Config.DTRHost+"/"+suite.Config.AdminUsername, testRepo, tag)()

		// make sure it has a tag
		response, err = suite.API.GetRepositoryTags(suite.Config.AdminUsername, testRepo)
		assert.Nil(suite.T(), err)
		assert.Len(suite.T(), response, 1)
	}()

	// make sure it has no tags after we clean up
	response, err = suite.API.GetRepositoryTags(suite.Config.AdminUsername, testRepo)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 0, len(response))
}

func (suite *RepoTagsAPITestSuite) TestGetTagTrust() {
	var err error
	testRepo := "testrepo"
	nonExistentTag := "tag"

	// create a repo under the admin namespace using the admin
	suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)
	defer suite.u.CreateRepoWithChecks(suite.Config.AdminUsername, testRepo, "blah", "blah blah", "public")()

	// make sure we get an error for non-existent tags
	_, err = suite.API.GetTagTrust(suite.Config.AdminUsername, testRepo, nonExistentTag)
	suite.u.AssertErrorCodes(err, http.StatusNotFound, errors.ErrorCodeNoSuchTag)

	// make sure we get an error for non-existent repos
	_, err = suite.API.GetTagTrust(suite.Config.AdminUsername, "noooooot", nonExistentTag)
	suite.u.AssertErrorCodes(err, http.StatusNotFound, errors.ErrorCodeNoSuchRepository)

	// make sure we get an error for non-existent user
	_, err = suite.API.GetTagTrust("wat", "not", "tag")
	suite.u.AssertErrorCodes(err, http.StatusNotFound, errors.ErrorCodeNoSuchAccount)

	// push and image and see if it appears in the list
	tag := "1337"
	defer suite.u.TagImageWithChecks(suite.Config.DTRHost+"/"+suite.Config.AdminUsername, testRepo, tag, suite.imageName)()
	func() {
		defer suite.u.PushImageWithChecks(suite.Config.DTRHost+"/"+suite.Config.AdminUsername, testRepo, tag)()

		// make sure we get the correct response when the tag exists but is not signed
		response, err := suite.API.GetTagTrust(suite.Config.AdminUsername, testRepo, tag)
		assert.Nil(suite.T(), err)
		assert.Equal(suite.T(), response.InNotary, false)
		assert.Equal(suite.T(), response.HashMismatch, false)
	}()

	// make sure that after we remove the image we are back to getting 404 errors
	suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)
	_, err = suite.API.GetTagTrust(suite.Config.AdminUsername, testRepo, tag)
	suite.u.AssertErrorCodes(err, http.StatusNotFound, errors.ErrorCodeNoSuchTag)

	// TODO: test signing with notary, etc.
}

func (suite *RepoTagsAPITestSuite) TestDeleteTag() {
	var err error
	testRepo := "testrepo"

	// create a repo under the admin namespace using the admin
	suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)
	defer suite.u.CreateRepoWithChecks(suite.Config.AdminUsername, testRepo, "blah", "blah blah", "public")()

	// make sure we don't get an error for non-existent tags
	err = suite.API.DeleteTag(suite.Config.AdminUsername, testRepo, "manifestortag")
	assert.Nil(suite.T(), err)

	// make sure we get an error for non-existent repos
	err = suite.API.DeleteTag(suite.Config.AdminUsername, "noooooot", "manifestortag")
	suite.u.AssertErrorCodes(err, http.StatusNotFound, errors.ErrorCodeNoSuchRepository)

	// make sure we get an error for non-existent user
	err = suite.API.DeleteTag("wat", "not", "manifestortag")
	suite.u.AssertErrorCodes(err, http.StatusNotFound, errors.ErrorCodeNoSuchAccount)

	// push an image so we can try deleting it
	tag := "1337"
	// also tag it with a second tag and make sure that one isn't deleted when the first one is deleted
	tag2 := "1338"
	func() {
		defer suite.u.TagImageWithChecks(suite.Config.DTRHost+"/"+suite.Config.AdminUsername, testRepo, tag, suite.imageName)()
		defer suite.u.TagImageWithChecks(suite.Config.DTRHost+"/"+suite.Config.AdminUsername, testRepo, tag2, suite.imageName)()
		deleteFirst := suite.u.PushImageWithChecks(suite.Config.DTRHost+"/"+suite.Config.AdminUsername, testRepo, tag)
		defer suite.u.PushImageWithChecks(suite.Config.DTRHost+"/"+suite.Config.AdminUsername, testRepo, tag2)()
		// make sure we can pull it
		err = suite.Docker.PullImage(path.Join(suite.Config.DTRHost, suite.Config.AdminUsername, testRepo)+":"+tag, suite.Config.AdminAuthConfig)
		assert.Nil(suite.T(), err)

		// make sure we can list both
		response, err := suite.API.GetRepositoryTags(suite.Config.AdminUsername, testRepo)
		assert.Nil(suite.T(), err)
		assert.Equal(suite.T(), 2, len(response))

		deleteFirst()

		// make sure we can list the one tag left
		response, err = suite.API.GetRepositoryTags(suite.Config.AdminUsername, testRepo)
		assert.Nil(suite.T(), err)
		assert.Equal(suite.T(), 1, len(response))

	}()
	// at this point the created repos have been deleted by going out of scope
	err = suite.Docker.PullImage(path.Join(suite.Config.DTRHost, suite.Config.AdminUsername, testRepo)+":"+tag, suite.Config.AdminAuthConfig)
	// make sure we can't pull it
	assert.NotNil(suite.T(), err)
	// make sure we can't list it
	response, err := suite.API.GetRepositoryTags(suite.Config.AdminUsername, testRepo)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 0, len(response))

	// TODO: try to delete by digest; how do we get the hash???

	util.AppendDockerIgnorableLoggedErrors([]string{"Error streaming logs: unexpected EOF"})
}

// TestTagstoreSavesMetadata loads dtr/true from the embedded .tar.gz below,
// pushes this to the tagstore and asserts that the tagstore saves metadata
// correctly.
func (suite *RepoTagsAPITestSuite) TestTagstoreSavesMetadata() {
	repoName := "tagstore"
	tag := "le-metadata"

	defer suite.u.CreateRepoWithChecks(suite.Config.AdminUsername, repoName, "short", "long", "public")()
	defer suite.u.LoadPackedImage()()
	defer suite.u.TagImageWithChecks(suite.Config.AdminNamespace(), repoName, tag, suite.u.PackedImageName)()
	deleteTag := suite.u.PushImageWithChecks(suite.Config.AdminNamespace(), repoName, tag)

	response, err := suite.API.GetRepositoryTags(suite.Config.AdminUsername, repoName)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 1, len(response))
	if len(response) != 1 {
		return
	}

	tr := response[0]
	// NOTE: The docker bug which pushes invalid manifests prior to 1.12 (see
	// https://github.com/docker/docker/pull/21949) means that the manifest's
	// hash can either be
	// '9c5e68903b1a8a4ef96dc3188183df2f5e52cb992cfd86c6dee3511e995b660b' (technically invalid)
	// or
	// '2d8ba47db0b607ae16709130751024840386d74aa682c3e4fcace438c82040ae' (valid)
	if tr.Digest != "sha256:2d8ba47db0b607ae16709130751024840386d74aa682c3e4fcace438c82040ae" && tr.Digest != "sha256:9c5e68903b1a8a4ef96dc3188183df2f5e52cb992cfd86c6dee3511e995b660b" {
		suite.T().Fatalf("unexpected manifest hash within tagstore: %s", tr.Digest)
	}
	// The tag and manifest API use the correct manifest digest of sha256:... - they omit that we store them as `namespace/repo@sha256:...`
	assert.Equal(suite.T(), tr.Digest, tr.Manifest.Digest)
	assert.Equal(suite.T(), int64(146), tr.Manifest.Size)
	assert.Equal(suite.T(), "linux", tr.Manifest.OS)
	assert.Equal(suite.T(), "amd64", tr.Manifest.Architecture)
	assert.Equal(suite.T(), []string{"sha256:464a2e7411dd2596f41c76cf85053f6cfb4353f1128a947a6a16f8439d86ec84"}, tr.Manifest.Layers)

	// Delete teh tag and assert that we no longer see it in the tagstore
	deleteTag()
	response, err = suite.API.GetRepositoryTags(suite.Config.AdminUsername, repoName)
	assert.Equal(suite.T(), 0, len(response))

	// Repush the tag and assert it's back in, ensuring our soft deletes undelete.
	defer suite.u.PushImageWithChecks(suite.Config.AdminNamespace(), repoName, tag)()
	response, err = suite.API.GetRepositoryTags(suite.Config.AdminUsername, repoName)
	assert.Equal(suite.T(), 1, len(response))
}

func TestRepoTagsAPISuite(t *testing.T) {
	suite.Run(t, new(RepoTagsAPITestSuite))
}

var serverPem = `
-----BEGIN CERTIFICATE-----
MIIFpDCCA4ygAwIBAgIRAPTE03PNFe0g0+bMF649chMwDQYJKoZIhvcNAQELBQAw
WzELMAkGA1UEBhMCVVMxDzANBgNVBAoTBkRvY2tlcjEPMA0GA1UECxMGRG9ja2Vy
MRYwFAYDVQQHEw1TYW4gRnJhbmNpc2NvMRIwEAYDVQQDEwlsb2NhbGhvc3QwHhcN
MTYwMTI2MTkwMzAxWhcNMTcwMTI1MTkwMzAxWjBbMQswCQYDVQQGEwJVUzEPMA0G
A1UEChMGRG9ja2VyMQ8wDQYDVQQLEwZEb2NrZXIxFjAUBgNVBAcTDVNhbiBGcmFu
Y2lzY28xEjAQBgNVBAMTCWxvY2FsaG9zdDCCAiIwDQYJKoZIhvcNAQEBBQADggIP
ADCCAgoCggIBAKpXwF1M9alCeiiHH7CqunAsKqROREqHKbPCrHsLsp5fRNU7n1Ix
1fIGCdOq9V2HKoL/jNXDQr9kWBC3yWdWmB2C7ue6xaWC3gaZZ5Cgt/rQa8JUBXtW
8qUeT2gQOQhr4acljIEPrMHeRnifthgTfCZi1TXoBTRLdUgA27UNu9UEwPivDt58
wcMofhULaMUVAzGpOI8h9I9QkmdB/Y2EMzGMandxvRV0FvkLcead6OmSyY3BPbFv
sLKYOm0JoKDbedwjYIiU4EuRze0nfnUNL2gnrU9/paWhSQoHDnDRvWeo9a6ashyL
U3jfFAy/8gqaaEOq2yktk0F5sZqv4TqiLc3OPvidyGQppNFEaN39BC8OSnjVrdni
j0yoKijoWHrsfMpyt7O4QNazQU09khvwc+SwBHWkqGBQe77CsHPrimM3D/ircKLY
2P/MZa4ryykre65Q8plxH7zVhPAnuE6MB9/gxT9wyC8xR8Kr3zozENMqWR1gzgDF
z0bxPsiH9fcwFA1vnpjzud/LNbXhxbTQFxOVylLX1BucjHTlnvDSVpSJrLX9FBZx
s1x+uXpiEI2ZhsL6i1+KqKOek2ZQVP2KTOIN4MS5Nvx990l1J9pEBwczCmVwU6fg
IX29i6CJqqmxXpgGdWdxJeSGp9JtE+snZLezNQI+aptM3BmM0ByPvrEFAgMBAAGj
YzBhMA4GA1UdDwEB/wQEAwICpDAPBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBSs
OkNxwwPtJS+Xl7wkZJT4W9z0rjAfBgNVHSMEGDAWgBSsOkNxwwPtJS+Xl7wkZJT4
W9z0rjANBgkqhkiG9w0BAQsFAAOCAgEAk9BByHfqEwnIFhV1Ue7v0GhE0nflzLAj
2m2QEcSvXO0eEWEpNqcRE+mkfQiADK+514x4RcXcjihS0IIlvq+C87GiC6QUxejZ
tZAa6mdh7KDdBOHcWDAIE4pfCh9bGUm9Z3+IdNdXQnq0Fc4FNIZo5VE65X1ipeXP
Sf0wo7/wC49cPxFSnJdghRShB2nVTAoIE1mCK2Qq6nDswhOV6uythOK+xe/N5yPN
ngHZW+6IHVhus/H9O3x6Jjeqp+HKtnYta/dA7I6bU5Tx2xKLH4yK0Ni4oMtMOOuD
c4qScYtiFzJWZTILH4UnJVIhx9Non2g4uMfNHLSZjNmgEbkOp7gSrOiwm1uG3FUC
hA6h0p2LweMTi4EP0iyuK2hm/uVHa9Lgx1QxgmBss8TmLM4QT9qh/VwMtES9Ococ
Zu7TxfNfEmKQzBMN3ng0iJ9929759ECpjvV0vYFd6c9TcpoJ/jhyzrCuUnqOiFAS
BqthGeOXPczaF4QXI+/uYpqVsnL7ygC7cOwT/ZhyXU952a3vnfZGUTuqEybHqzv+
1+7FllIdbzPNczswVZ+fibO2UMFb4B9AGqWuniHoRBZmn922Hbyy8FVMt0Th//h/
Ntx5MSVJDuCWsKsY5PYfrNyVShaj02sRJDqxjpfoSgE3uevgJIyNmABNpoFbNg0x
VzxsYLYIitY=
-----END CERTIFICATE-----
-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----
`
