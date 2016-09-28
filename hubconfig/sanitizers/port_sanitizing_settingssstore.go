package sanitizers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/docker/dhe-deploy"
	"github.com/docker/dhe-deploy/hubconfig"
)

type PortSanitizingSettingsStore struct {
	hubconfig.SettingsStore
}

// list of restricted ports from http://goo.gl/3hydqK
var PORT_BLACKLIST = []uint16{
	0, 1, 7, 9, 11, 13, 15, 17, 19, 20, 21, 22, 23, 25, 37, 42, 43, 53, 77,
	79, 87, 95, 101, 102, 103, 104, 109, 110, 111, 113, 115, 117, 119, 123,
	135, 139, 143, 179, 389, 465, 512, 513, 514, 515, 526, 530, 531, 532, 540,
	556, 563, 587, 601, 636, 993, 995, 2049, 4045, 6000,
}

var PORT_BLACKLIST_SET map[uint16]struct{}

func init() {
	PORT_BLACKLIST_SET = make(map[uint16]struct{})
	marker := struct{}{}
	for _, p := range PORT_BLACKLIST {
		PORT_BLACKLIST_SET[p] = marker
	}
}

func (s PortSanitizingSettingsStore) SetHubConfig(hubConfig *hubconfig.UserHubConfig) error {
	host := hubConfig.DTRHost
	parts := strings.Split(host, ":")
	if len(parts) == 0 {
		return s.SettingsStore.SetUserHubConfig(hubConfig)
	}

	port, err := strconv.ParseInt(parts[1], 10, 16)
	if err != nil {
		return fmt.Errorf("Failed to validate port in DTR host: %s", err)
	}

	err = checkValidPort(uint16(port))
	if err != nil {
		return fmt.Errorf("Invalid HTTPS port: %s", err)
	}

	return s.SettingsStore.SetUserHubConfig(hubConfig)
}

func (s PortSanitizingSettingsStore) SetHAConfig(haConfig *hubconfig.HAConfig) error {
	for id, replicaConfig := range haConfig.ReplicaConfig {
		if replicaConfig.HTTPPort == 0 {
			replicaConfig.HTTPPort = deploy.AdminPort
		}
		err := checkValidPort(replicaConfig.HTTPPort)
		if err != nil {
			return fmt.Errorf("Invalid HTTP port: %v", err)
		}
		if replicaConfig.HTTPSPort == 0 {
			replicaConfig.HTTPSPort = deploy.AdminTlsPort
		}
		err = checkValidPort(replicaConfig.HTTPSPort)
		if err != nil {
			return fmt.Errorf("Invalid HTTPS port: %v", err)
		}
		if replicaConfig.HTTPPort == replicaConfig.HTTPSPort {
			return fmt.Errorf("Can not use the same port for HTTP and HTTPS")
		}
		haConfig.ReplicaConfig[id] = replicaConfig
	}
	return s.SettingsStore.SetHAConfig(haConfig)
}

func checkValidPort(port uint16) error {
	if _, blacklisted := PORT_BLACKLIST_SET[port]; blacklisted {
		return fmt.Errorf("%d is not a valid port", port)
	}
	return nil
}
