package integration

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"runtime"
	"testing"

	"github.com/docker/dhe-deploy/integration/framework"
	"github.com/docker/dhe-deploy/integration/util"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// used to throw to the channel variable in the function that is called multiple times in parallel
type APIReturn struct {
	Resp *http.Response
	Err  error
}

type RequestTestSuite struct {
	suite.Suite
	*framework.IntegrationFramework
	u      *util.Util
	status chan *APIReturn
}

func (suite *RequestTestSuite) SetupSuite() {
	suite.IntegrationFramework, suite.u = setupFramework(suite)
}

func (suite *RequestTestSuite) SetupTest() {
	util.WipeDTRIgnorableLoggedErrors()
	util.WipeDockerIgnorableLoggedErrors()
}

func (suite *RequestTestSuite) TearDownTest() {
	suite.u.TestLogs()
}

//test that ports are up
func (suite *RequestTestSuite) TestPort80Up() {
	_, err := tls.Dial("tcp", fmt.Sprintf("%s:80", suite.Config.DTRHost), &tls.Config{
		InsecureSkipVerify: true,
	})
	assert.NotNil(suite.T(), err, fmt.Sprintf("port 80 not up error: %v", err))
}

func (suite *RequestTestSuite) TestPort443Up() {
	_, err := tls.Dial("tcp", fmt.Sprintf("http://%s:443", suite.Config.DTRHost), &tls.Config{
		InsecureSkipVerify: true,
	})
	assert.NotNil(suite.T(), err, fmt.Sprintf("port 443 not up error: %v", err))
}

// main function to run the parrallel api calls from
// NOTE: still need to add to this and to the apicall function library
func (suite *RequestTestSuite) TestParallelSpam() {
	previousThreads := runtime.GOMAXPROCS(runtime.NumCPU()) // allow for parallelization
	defer runtime.GOMAXPROCS(previousThreads)
	arr := suite.throwParallelRequests(1000, 50, http.StatusOK, suite.testLoadBalancer)
	assert.NotEqual(suite.T(), arr, nil, "Error array not initialized")
}

// A function to run in parallel specifically for checking the load balancer endpoint
// NOTE: need to add a function like this that send a value of type APIReturn through the channel for every api endpoint that you want to test
func (suite *RequestTestSuite) testLoadBalancer() *APIReturn {
	resp, err := suite.API.LoadBalancerStatusResponse()
	return &APIReturn{
		Resp: resp,
		Err:  err,
	}
}

// runs given function call in parallel numReqs number of times. Set up now so that it returns an array of all the errors and Reponses and expects caller to handle checking
// however you could add to it so that it handles the error checking and returns just one individual error of some type
func (suite *RequestTestSuite) throwParallelRequests(numReqs, connectionPoolSize, expectedStatus int, callToThrow func() *APIReturn) []*APIReturn {
	responses := make([]*APIReturn, numReqs)
	poolChannel := make(chan struct{}, connectionPoolSize)

	statusChannel := make(chan *APIReturn)
	numReqsFailed := 0

	for i := 0; i < numReqs; i++ {
		go func() {
			poolChannel <- struct{}{}
			defer func() { <-poolChannel }()

			statusChannel <- callToThrow()
		}()
	}

	for i := 0; i < numReqs; i++ {
		apiReturn := <-statusChannel
		if apiReturn.Err != nil || apiReturn.Resp.StatusCode != expectedStatus { // if status channel returns an unexpected status code or an error that is not nil, do something
			numReqsFailed++ // change this to whatever you want in this case. error/unexpected status codes should be handled by whatever is calling this function so that all requests can finish
		}
		responses[i] = apiReturn
	}
	return responses
}

func TestRequestSuite(t *testing.T) {
	suite.Run(t, new(RequestTestSuite))
}
