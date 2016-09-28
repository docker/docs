package util

import (
	"github.com/docker/dhe-deploy/integration/framework"
	"github.com/stretchr/testify/suite"
)

func MakeUtil(suite suite.TestingSuite, fw *framework.IntegrationFramework) *Util {
	return &Util{
		TestingSuite:         suite,
		IntegrationFramework: fw,
		// This is the name of the image loaded in Util.LoadPackedImage()
		PackedImageName: "dtr/true",
	}
}

type Util struct {
	*framework.IntegrationFramework
	suite.TestingSuite

	PackedImageName string
}
