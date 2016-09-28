package apiclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
)

const scheduleSubroute = "/api/v0/admin/settings/registry/garbageCollection/schedule"

type ScheduleContainer struct {
	Schedule string `json:"schedule"`
}

func (c *apiClient) SetGarbageCollectionSchedule(schedule string) error {
	response, err := c.makeRequest("POST", url.URL{Path: path.Join(scheduleSubroute)}, &ScheduleContainer{Schedule: schedule})
	defer response.Body.Close()

	if err != nil {
		return err
	}

	return validateStatusCode(response, http.StatusOK)
}

func (c *apiClient) DeleteGarbageCollectionSchedule() error {
	response, err := c.makeRequest("DELETE", url.URL{Path: path.Join(scheduleSubroute)}, nil)
	defer response.Body.Close()

	if err != nil {
		return err
	}

	return validateStatusCode(response, http.StatusOK)
}

func (c *apiClient) GetGarbageCollectionSchedule() (string, error) {
	response, err := c.makeRequest("GET", url.URL{Path: path.Join(scheduleSubroute)}, nil)
	defer func() {
		if response != nil {
			response.Body.Close()
		}
	}()

	if err != nil {
		return "", fmt.Errorf("Failed to retrieve response: %v", err)
	}

	if err = validateStatusCode(response, http.StatusOK); err != nil {
		return "", err
	}

	var scheduleContainer ScheduleContainer
	if err := json.NewDecoder(response.Body).Decode(&scheduleContainer); err != nil {
		return "", fmt.Errorf("Failed to decode response: %v", err)
	}

	return scheduleContainer.Schedule, nil
}
