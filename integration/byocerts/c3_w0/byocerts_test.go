package instance

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/docker/orca/integration/byocerts"
	"github.com/docker/orca/integration/utils"
)

type TestSuite struct {
	byocerts.BYOCertsSuite
}

func (s *TestSuite) GetNodeCounts() (controllerCount int, workerCount int) {
	controllerCount = 3
	workerCount = 0
	return
}

func TestAcceptance(t *testing.T) {
	utils.HandleTestArgs(t)
	s := &TestSuite{}
	s.Init(s)
	suite.Run(t, s)
}
