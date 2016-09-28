package ui

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/docker/orca/integration/utils"
	"github.com/docker/orca/integration/utils/ui"
)

func (s *TestSuite) TestWebTests() {
	serverURL := "https://" + s.OrcaClientTestSuite.Client.URL.Host
	ui.TestUI(s.T(), serverURL)
}

type TestSuite struct {
	utils.OrcaClientTestSuite
}

func TestStress(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
