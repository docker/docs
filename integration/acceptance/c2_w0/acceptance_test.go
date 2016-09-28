package instance

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/docker/orca/integration/acceptance"
	"github.com/docker/orca/integration/utils"
)

type AcceptanceSuite struct {
	acceptance.AcceptanceSuite
}

func (s *AcceptanceSuite) GetNodeCounts() (controllerCount int, workerCount int) {
	controllerCount = 2
	workerCount = 0
	return
}

func TestAcceptance(t *testing.T) {
	utils.HandleTestArgs(t)
	s := &AcceptanceSuite{}
	s.Init(s)
	suite.Run(t, s)
}
