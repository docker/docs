package util

import (
	"fmt"
	"time"

	"github.com/docker/dhe-deploy/pkg/jobrunner-framework/schema"
	"github.com/docker/dhe-deploy/shared/dtrutil"
)

func (u *Util) PollAvailable() error {
	time.Sleep(time.Second)
	return dtrutil.Poll(time.Second, u.Config.RetryAttempts, func() error {
		nginxLoadBalancerStatus, err := u.API.LoadBalancerStatus()
		if err != nil {
			return err
		}

		for _, server := range nginxLoadBalancerStatus.NginxServers.NginxServer {
			if server.Status != "up" {
				return fmt.Errorf("Failed to wait for DTR to come up. Server is not up %v", server)
			}
		}

		return nil
	})
}

func (u *Util) WaitForJob(id string) error {
	// Wait for the job to finish
	return dtrutil.Poll(time.Second, u.Config.RetryAttempts, func() error {
		status, err := u.API.GetJobStatus(id)
		if err != nil {
			return fmt.Errorf("Failed to get job status: %v", err)
		}

		if status == schema.JobStatusError {
			return fmt.Errorf("job failed")
		}

		if status == schema.JobStatusDone {
			return nil
		}

		return fmt.Errorf("job status is not yet done: %s", status)
	})
}
