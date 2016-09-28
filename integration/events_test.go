package integration

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"testing"
	"time"

	"github.com/docker/dhe-deploy/events/types"
	"github.com/docker/dhe-deploy/integration/apiclient"
	"github.com/docker/dhe-deploy/integration/framework"
	"github.com/docker/dhe-deploy/integration/util"
	"github.com/docker/dhe-deploy/manager/schema"

	log "github.com/Sirupsen/logrus"
	"github.com/samalba/dockerclient"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type EventsTestSuite struct {
	suite.Suite
	*framework.IntegrationFramework
	u              *util.Util
	adminAccountId string
}

func (suite *EventsTestSuite) SetupSuite() {
	suite.IntegrationFramework, suite.u = setupFramework(suite)

	account, err := suite.API.GetAccount(suite.Config.AdminUsername)
	require.Nil(suite.T(), err, "Getting account returned error: %s", err)
	suite.adminAccountId = account.ID

	// Ensure we've ran migrations so that GC works in blob delete test
	job, err := suite.API.RunJobByAction("tagmigration")
	require.NotEmpty(suite.T(), job.ID, fmt.Sprintf("tagmiration not created: %#v", job))
	require.Nil(suite.T(), err)
	err = suite.u.WaitForJob(job.ID)
	require.Nil(suite.T(), err, "Failed to get complete tagmigration task in time: %v", err)
}

func (suite *EventsTestSuite) SetupTest() {
	util.WipeDTRIgnorableLoggedErrors()
	util.WipeDockerIgnorableLoggedErrors()
}

func (suite *EventsTestSuite) TearDownTest() {
	suite.u.TestLogs()
}

func (suite *EventsTestSuite) TestCreateRepoEvent() {
	reponame := "testcreaterepo"

	// create a repo under the admin namespace using the admin
	require.NoError(suite.T(), suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword))
	defer suite.u.CreateRepoWithChecks(suite.Config.AdminUsername, reponame, "desc", "longer desc", "public")()

	evts, _, err := suite.API.GetEvents()
	require.Nil(suite.T(), err)

	require.NotEmpty(suite.T(), evts.Events, "There should be at least one event returned but got %d", len(evts.Events))
	require.Nil(suite.T(), validateEvent(evts.Events[0], types.Create, suite.adminAccountId, types.Repository, suite.Config.AdminUsername+"/"+reponame))
}

func (suite *EventsTestSuite) TestDeleteRepoEvent() {
	reponame := "testdeleterepo"

	// create and delete a repo under the admin namespace using the admin
	require.NoError(suite.T(), suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword))
	suite.u.CreateRepoWithChecks(suite.Config.AdminUsername, reponame, "desc", "longer desc", "public")()

	evts, _, err := suite.API.GetEvents()
	require.Nil(suite.T(), err)

	require.NotEmpty(suite.T(), evts.Events, "There should be at least one event returned but got %d", len(evts.Events))
	require.Nil(suite.T(), validateEvent(evts.Events[0], types.Delete, suite.adminAccountId, types.Repository, suite.Config.AdminUsername+"/"+reponame))
}

func (suite *EventsTestSuite) TestQueryParamsPagination() {
	reponame1 := "testrepo1"
	reponame2 := "testrepo2"

	// create a repo under the admin namespace using the admin
	require.NoError(suite.T(), suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword))
	defer suite.u.CreateRepoWithChecks(suite.Config.AdminUsername, reponame1, "desc", "longer desc", "public")()
	defer suite.u.CreateRepoWithChecks(suite.Config.AdminUsername, reponame2, "desc", "longer desc", "public")()

	// test limit 1
	values := url.Values{}
	values.Add("limit", "1")
	evts, next, err := suite.API.GetEventsWithParams(values)
	require.Nil(suite.T(), err)

	require.Len(suite.T(), evts.Events, 1, "Only 1 event should be returned when limit is set to 1")
	require.Nil(suite.T(), validateEvent(evts.Events[0], types.Create, suite.adminAccountId, types.Repository, suite.Config.AdminUsername+"/"+reponame2))

	require.NotNil(suite.T(), next, "First request should return a header specifying the next page")
	evts, _, err = suite.API.GetEventsWithParams(*next)
	require.Nil(suite.T(), err)

	require.NotEmpty(suite.T(), evts.Events, "There should be at least one event returned but got %d", len(evts.Events))
	require.Nil(suite.T(), validateEvent(evts.Events[0], types.Create, suite.adminAccountId, types.Repository, suite.Config.AdminUsername+"/"+reponame1))
}

func (suite *EventsTestSuite) TestQueryParamsBeforeAndAfter() {
	reponame1 := "pubbeforerepo"
	reponame2 := "pubafterrepo"

	// create a repo under the admin namespace using the admin
	require.NoError(suite.T(), suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword))
	defer suite.u.CreateRepoWithChecks(suite.Config.AdminUsername, reponame1, "desc", "longer desc", "public")()
	// sleep 2 seconds to ensure an isolated time between the two repos
	time.Sleep(2 * time.Second)
	defer suite.u.CreateRepoWithChecks(suite.Config.AdminUsername, reponame2, "desc", "longer desc", "public")()

	params := url.Values{}
	params.Add("limit", "2")
	evts, _, err := suite.API.GetEventsWithParams(params)
	require.Len(suite.T(), evts.Events, 2, "Limit is set to 2 and 2 repos have been created so 2 events should be returned")

	// since we waited 2 seconds, adding 1 second should guarantee time between events
	between := evts.Events[1].PublishedAt.Add(time.Second)

	// test publishedBefore
	values := url.Values{}
	layout := "2006-01-02T15:04:05.000Z"
	values.Add("publishedBefore", between.Format(layout))
	evts, _, err = suite.API.GetEventsWithParams(values)
	require.Nil(suite.T(), err)

	require.NotEmpty(suite.T(), evts.Events, "There should be an event published before the time provided")
	require.Nil(suite.T(), validateEvent(evts.Events[0], types.Create, suite.adminAccountId, types.Repository, suite.Config.AdminUsername+"/"+reponame1))

	// test publishedAfter
	values = url.Values{}
	values.Add("publishedAfter", between.Format(layout))
	evts, _, err = suite.API.GetEventsWithParams(values)
	require.Nil(suite.T(), err)

	require.NotEmpty(suite.T(), evts.Events, "There should be an event published after the time provided")
	require.Nil(suite.T(), validateEvent(evts.Events[0], types.Create, suite.adminAccountId, types.Repository, suite.Config.AdminUsername+"/"+reponame2))
}

func (suite *EventsTestSuite) TestByActor() {
	reponame := "actortest"

	require.NoError(suite.T(), suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword))
	defer suite.u.CreateRepoWithChecks(suite.Config.AdminUsername, reponame, "desc", "longer desc", "public")()

	values := url.Values{}
	values.Add("actorId", suite.adminAccountId)
	evts, _, err := suite.API.GetEventsWithParams(values)
	require.Nil(suite.T(), err)

	require.NotEmpty(suite.T(), evts.Events, "There should be an event by admin actor")
	require.Nil(suite.T(), validateEvent(evts.Events[0], types.Create, suite.adminAccountId, types.Repository, suite.Config.AdminUsername+"/"+reponame))

	values = url.Values{}
	values.Add("actorId", "randomaccountid")
	evts, _, err = suite.API.GetEventsWithParams(values)
	require.Nil(suite.T(), err)

	require.Empty(suite.T(), evts.Events, "There should be no events with actor randomaccountid")
}

func (suite *EventsTestSuite) TestByEventType() {
	reponame := "eventtest"
	var err error

	require.NoError(suite.T(), suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword))
	defer suite.u.CreateRepoWithChecks(suite.Config.AdminUsername, reponame, "desc", "longer desc", "public")()

	values := url.Values{}
	values.Add("eventType", types.Create)
	evts, _, err := suite.API.GetEventsWithParams(values)
	require.Nil(suite.T(), err)

	require.NotEmpty(suite.T(), evts.Events, "There should be an event with type create")
	require.Nil(suite.T(), validateEvent(evts.Events[0], types.Create, suite.adminAccountId, types.Repository, suite.Config.AdminUsername+"/"+reponame))

	values = url.Values{}
	values.Add("eventType", "invalideventtype")
	evts, _, err = suite.API.GetEventsWithParams(values)
	require.Nil(suite.T(), err)

	require.Empty(suite.T(), evts.Events, "There should be no events with type invalideventtype")
}

func (suite *EventsTestSuite) TestByEventTypeAndActor() {
	reponame := "eventactortest"
	var err error

	require.NoError(suite.T(), suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword))
	defer suite.u.CreateRepoWithChecks(suite.Config.AdminUsername, reponame, "desc", "longer desc", "public")()

	values := url.Values{}
	values.Add("actorId", suite.adminAccountId)
	values.Add("eventType", types.Create)
	evts, _, err := suite.API.GetEventsWithParams(values)
	require.Nil(suite.T(), err)

	require.NotEmpty(suite.T(), evts.Events, "There should be an event with type create by admin actor")
	require.Nil(suite.T(), validateEvent(evts.Events[0], types.Create, suite.adminAccountId, types.Repository, suite.Config.AdminUsername+"/"+reponame))

	values = url.Values{}
	values.Add("actorId", "randomaccountid")
	values.Add("eventType", types.Create)
	evts, _, err = suite.API.GetEventsWithParams(values)
	require.Nil(suite.T(), err)

	require.Empty(suite.T(), evts.Events, "There should be no events with actor randomaccountid, regardless of event type")
}

func (suite *EventsTestSuite) TestByActorInvalidPermissions() {
	reponame1 := "actortestinvalidpermissions1"
	reponame2 := "actortestinvalidpermissions2"
	var err error

	// login as admin
	require.NoError(suite.T(), suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword))

	// create user
	normalUser := "eventuser"
	normalPassword := util.GenerateRandomPassword(10)
	defer suite.u.CreateUserWithChecks(normalUser, normalPassword)()
	if err := suite.u.ActivateUser(normalUser); err != nil {
		suite.T().Error(err)
	}

	// create repo
	defer suite.u.CreateRepoWithChecks(suite.Config.AdminUsername, reponame1, "desc", "longer desc", "public")()

	// logout
	require.Nil(suite.T(), suite.API.Logout())

	// login as normal user and try to get events
	require.NoError(suite.T(), suite.API.Login(normalUser, normalPassword))
	values := url.Values{}
	values.Add("actorId", suite.adminAccountId)
	evts, _, err := suite.API.GetEventsWithParams(values)

	// should get an unauthorized error
	apiErr, ok := err.(*apiclient.APIError)
	if !ok {
		suite.T().Errorf("Requesting another user's data should give back an api error")
	}
	require.Equal(suite.T(), http.StatusForbidden, apiErr.HTTPStatusCode)
	require.Len(suite.T(), apiErr.Errors, 1)
	require.Equal(suite.T(), "The client is not authorized.", apiErr.Errors[0].Message)
	require.Equal(suite.T(), "Users can only query for events created by themselves unless they have admin privileges", apiErr.Errors[0].Detail)

	// create private repo as normal user
	defer suite.u.CreateRepoWithChecks(normalUser, reponame2, "desc", "longer desc", "private")()

	// get normal user's id
	normalUserAccount, err := suite.API.GetAccount(normalUser)
	require.Nil(suite.T(), err, "Getting account returned error: %s", err)
	normalUserId := normalUserAccount.ID

	// get events as normal user
	values = url.Values{}
	values.Add("actorId", normalUserId)
	evts, _, err = suite.API.GetEventsWithParams(values)

	// normal user should be able to see his own private repo
	require.NotEmpty(suite.T(), evts.Events, "Normal user should get events for his own private repo")
	require.Nil(suite.T(), validateEvent(evts.Events[0], types.Create, normalUserId, types.Repository, normalUser+"/"+reponame2))

	// logout
	require.Nil(suite.T(), suite.API.Logout())

	// login as admin and get events
	require.NoError(suite.T(), suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword))
	values = url.Values{}
	values.Add("actorId", normalUserId)
	evts, _, err = suite.API.GetEventsWithParams(values)

	// admin should be able to see normal user's private repo
	require.NotEmpty(suite.T(), evts.Events, "Admin should be able to see all private events")
	require.Nil(suite.T(), validateEvent(evts.Events[0], types.Create, normalUserId, types.Repository, normalUser+"/"+reponame2))
}

func (suite *EventsTestSuite) TestEventWS() {
	reponame := "testrepo"
	var err error

	// login as admin
	suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)

	// create normal user
	normalUser := "wseventuser"
	normalPassword := util.GenerateRandomPassword(10)
	defer suite.u.CreateUserWithChecks(normalUser, normalPassword)()
	if err := suite.u.ActivateUser(normalUser); err != nil {
		suite.T().Error(err)
	}

	// open a WS
	adminConn, err := suite.API.EventsWebsocketAuthd()
	require.Nil(suite.T(), err)
	defer adminConn.Close()

	// logout as admin
	require.Nil(suite.T(), suite.API.Logout())

	// login as normal user and get websocket
	suite.API.Login(normalUser, normalPassword)

	userConn, err := suite.API.EventsWebsocketAuthd()
	require.Nil(suite.T(), err)
	defer userConn.Close()

	// make a channel to receive the event on
	adminEvtChan := make(chan schema.Event)
	defer close(adminEvtChan)
	userEvtChan := make(chan schema.Event)
	defer close(userEvtChan)

	// Because of the order of exiting tests, this gofunc is structured weirdly.
	// When `require` fails, first the event channel gets closed. Then conn.ReadJSON
	// returns an error. Then, the test actually fails.
	// If we had an error channel coming from the gofunc, the channel would be closed before
	// we could write the error through. Instead we can log the error and return, causing
	// the test to timeout and leaving a trace above to show the error.
	go func() {
		evt := schema.Event{}
		for i := 0; i < 2; i++ {
			// readjson is a blocking call
			if err := adminConn.ReadJSON(&evt); err != nil {
				log.Errorf("Reading JSON from websocket failed. This is likely because the test timed out waiting for an event to be read. Error: %s", err.Error())
				return
			}
			adminEvtChan <- evt
		}
	}()

	// create another listener, but this one should be anonymous
	go func() {
		evt := schema.Event{}
		for i := 0; i < 2; i++ {
			// readjson is a blocking call
			if err := userConn.ReadJSON(&evt); err != nil {
				// dont catch the error, because it will likely be that the socket is closed
				// if something does get pushed through the channel, we'll catch it later
				return
			}
			userEvtChan <- evt
		}
	}()

	// create repo
	defer suite.u.CreateRepoWithChecks(normalUser, reponame, "desc", "longer desc", "private")()

	normalUserAccount, err := suite.API.GetAccount(normalUser)
	require.Nil(suite.T(), err, "Getting normal user account returned error: %s", err)

	// make sure both WS get the event, since normal user created it and admin is admin
	select {
	case evt := <-adminEvtChan:
		require.Nil(suite.T(), validateEvent(evt, types.Create, normalUserAccount.ID, types.Repository, normalUser+"/"+reponame))
	case <-time.After(5 * time.Second):
		require.Nil(suite.T(), fmt.Errorf("Admin websocket timed out without receiving event"))
	}
	select {
	case evt := <-userEvtChan:
		require.Nil(suite.T(), validateEvent(evt, types.Create, normalUserAccount.ID, types.Repository, normalUser+"/"+reponame))
	case <-time.After(5 * time.Second):
		require.Nil(suite.T(), fmt.Errorf("User websocket timed out without receiving event"))
	}

	// wait a full two minutes, making sure that the ws ping/pong works
	time.Sleep(2 * time.Minute)

	require.Nil(suite.T(), suite.API.Logout())
	suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword)

	reponame = "otherrepo"

	// create repo
	defer suite.u.CreateRepoWithChecks(suite.Config.AdminUsername, reponame, "desc", "longer desc", "public")()

	// make sure that only the admin ws got the event, since normal user did not create the event
	select {
	case evt := <-adminEvtChan:
		require.Nil(suite.T(), validateEvent(evt, types.Create, suite.adminAccountId, types.Repository, suite.Config.AdminUsername+"/"+reponame))
	case <-time.After(5 * time.Second):
		require.Nil(suite.T(), fmt.Errorf("Authd websocket timed out without receiving event"))
	}
	select {
	case evt := <-userEvtChan:
		require.Nil(suite.T(), fmt.Errorf("Unauthed websocket should not get notified about private events, but got: %+v", evt))
	case <-time.After(5 * time.Second):
		// should time out
	}
}

func (suite *EventsTestSuite) TestCreateAndDeleteTag() {
	testRepo := "deletetag"
	tag := "who"

	require.NoError(suite.T(), suite.API.Login(suite.Config.AdminUsername, suite.Config.AdminPassword))
	defer suite.u.CreateRepoWithChecks(suite.Config.AdminUsername, testRepo, "desc", "longer desc", "public")()

	// prepare for pushing
	smallImage := "tianon/true:latest"
	err := suite.Docker.PullImage(smallImage, nil)
	require.Nil(suite.T(), err)

	defer suite.u.TagImageWithChecks(suite.Config.DTRHost+"/"+suite.Config.AdminUsername, testRepo, tag, smallImage)()

	// make sure we can push before we delete
	authConfig := dockerclient.AuthConfig{
		Username: suite.Config.AdminUsername,
		Password: suite.Config.AdminPassword,
		Email:    "a@a.a",
	}

	// Test pushing events
	err = suite.Docker.PushImage(path.Join(suite.Config.DTRHost, suite.Config.AdminUsername, testRepo), tag, &authConfig)
	require.Nil(suite.T(), err)

	evts, _, err := suite.API.GetEvents()
	require.Nil(suite.T(), err)

	require.NotEmpty(suite.T(), evts.Events, "There should be a tag creation event last")
	require.Nil(suite.T(), validateEvent(evts.Events[0], types.Update, suite.adminAccountId, types.Tag, suite.Config.AdminUsername+"/"+testRepo+":"+tag))

	// Test deleting the tag

	err = suite.API.DeleteTag(suite.Config.AdminUsername, testRepo, tag)
	require.Nil(suite.T(), err)

	evts, _, err = suite.API.GetEvents()
	require.Nil(suite.T(), err)

	require.NotEmpty(suite.T(), evts.Events, "There should be a tag deletion event last")
	require.Nil(suite.T(), validateEvent(evts.Events[0], types.Delete, suite.adminAccountId, types.Tag, suite.Config.AdminUsername+"/"+testRepo+":"+tag))
}

func (suite *EventsTestSuite) TestDeleteBlob() {
	imageTag := "rick_sanchez" + strconv.FormatInt(util.Prng.Int63(), 16)
	imageName := "tianon/true"
	repoName := "mytrue"
	err := suite.Docker.PullImage(imageName, nil)
	require.Nil(suite.T(), err, "%s", err)

	deleteRepo := suite.u.CreateRepoWithChecks(suite.Config.AdminUsername, repoName, "test", "long test", "public")
	suite.u.TagImageWithChecks(suite.Config.DTRHost+"/"+suite.Config.AdminUsername, repoName, imageTag, imageName)
	suite.u.PushImageWithChecks(suite.Config.DTRHost+"/"+suite.Config.AdminUsername, repoName, imageTag)

	response, err := suite.API.GetRepositoryTags(suite.Config.AdminUsername, repoName)
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), 1, len(response))
	var layers []string
	for _, layer := range response[0].Manifest.Layers {
		layers = append(layers, layer)
	}
	require.NotEmpty(suite.T(), layers, "There should be layers because tianon/true was pushed")

	deleteRepo()
	suite.RunGC()

	evts, _, err := suite.API.GetEvents()
	require.Nil(suite.T(), err)

	require.NotEmpty(suite.T(), evts.Events, "There should be at least one blob deletion event")
	for _, layer := range layers {
		success := false
		for _, evt := range evts.Events {
			err = validateEvent(evt, types.Delete, types.GarbageCollector, types.Blob, layer)
			if err == nil {
				success = true
				break
			}
		}

		if !success {
			require.NoError(suite.T(), fmt.Errorf("Layer %s was part of the tag push but was not in any blob deletion events"))
		}
	}
}

func (suite *EventsTestSuite) RunGC() {
	job, err := suite.API.RunJobByAction("gc")
	require.Nil(suite.T(), err)
	err = suite.u.WaitForJob(job.ID)

	require.Nil(suite.T(), err, "Failed to get complete gc task in time: %v", err)
}

func validateEvent(evt schema.Event, eventType, actorId, objectType, objectId string) error {
	if evt.Type != eventType {
		return fmt.Errorf("Event type should be %s, but is actually %s", eventType, evt.Type)
	}
	if evt.Actor != actorId {
		return fmt.Errorf("Event actor should be %s, but is actually %s", actorId, evt.Actor)
	}
	if evt.PublishedAt == (time.Time{}) {
		return fmt.Errorf("Event published should not be empty")
	}

	if evt.Object.Type != objectType {
		return fmt.Errorf("Event object type should be %s, but is actually %s", objectType, evt.Object.Type)
	}
	if evt.Object.ID != objectId {
		return fmt.Errorf("Event object id should be %s, but is actually %s", objectId, evt.Object.ID)
	}
	return nil
}

func TestEventsSuite(t *testing.T) {
	suite.Run(t, new(EventsTestSuite))
}
