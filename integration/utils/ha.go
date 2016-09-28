package utils

import (
	"fmt"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/samalba/dockerclient"
)

// Return an error if the controller count is incorrect (after retrying a little
func VerifyControllerCount(client *dockerclient.DockerClient, expectedCount, retryCount int) error {
	log.Infof("Verifying controller count matches expcted value: %d", expectedCount)
	lastError := ""
	for i := 0; i < retryCount; i++ {
		info, err := client.Info()
		if err != nil {
			return err
		}
		log.Debugf("Info Driver Status: %v", info.DriverStatus)
		for _, driver := range info.DriverStatus {
			if strings.Contains(driver[0], "Cluster Managers") {
				log.Debugf("Cluster Managers: %v", driver)
				if strings.TrimSpace(driver[1]) == fmt.Sprintf("%d", expectedCount) {
					return nil
				} else {
					lastError = fmt.Sprintf("Incorrect number of controllers.  Expected %d but got %s", expectedCount, strings.TrimSpace(driver[1]))
					break
				}
			}
		}
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf(lastError)
}
