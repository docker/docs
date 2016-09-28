package api

import (
	log "github.com/Sirupsen/logrus"
	"github.com/docker/orca/controller/mock_test"
)

func getTestApi() (*Api, error) {
	log.SetLevel(log.ErrorLevel)
	m := mock_test.MockManager{}

	config := ApiConfig{
		ListenAddr:         "",
		Manager:            m,
		AuthWhiteListCIDRs: nil,
		AllowInsecure:      true,
		ControllerCertPEM:  nil,
		ControllerKeyPEM:   nil,
	}

	return NewApi(config)
}
