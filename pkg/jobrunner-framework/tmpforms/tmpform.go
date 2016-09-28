package tmpforms

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/docker/dhe-deploy/pkg/jobrunner-framework/worker"
	"github.com/docker/orca/enzi/api/errors"
)

// This package is temporary. It should be merged with enzi's responses package when the jobrunner is moved back into enzi

// CronCreate is a form for creating a new cron
type CronCreate struct {
	Action      string            `json:"action" description:"The action which the cron will perform"`
	Schedule    string            `json:"schedule" description:"The for the cron as a cronspec string: (seconds) (minutes) (hours) (day of month) (month) (day of week) or @hourly, @weekly, etc."`
	Retries     int               `json:"retries" description"The number of times to retry a job if it fails"`
	Parameters  map[string]string `json:"parameters",description:"Extra parameters to pass to the job. The available parameters depend on the job."`
	Deadline    string            `json:"deadline",description:"After this amount of time has passed, a SIGTERM will be sent"`
	StopTimeout string            `json:"stopTimeout",description:"This long after SIGTERM is sent, SIGKILL will be sent if the proccess is still alive"`
	ActionInfo  map[string]worker.ActionInfo
}

// ValidateJSON unmarshals this form from the given JSON and then validates
// the form, returning a list of any decoding or form validation errors.
func (form *CronCreate) ValidateJSON(r io.Reader) []*errors.APIError {
	return validateJSONForm(r, form)
}

// Validate validates this CronCreate form.
func (form *CronCreate) Validate() (apiErrs []*errors.APIError) {
	if form.Deadline != "" {
		_, err := time.ParseDuration(form.Deadline)
		if err != nil {
			return []*errors.APIError{errors.InvalidFormField("deadline", fmt.Sprintf("not a valid duration: %s", err))}
		}
	}
	if form.StopTimeout != "" {
		_, err := time.ParseDuration(form.StopTimeout)
		if err != nil {
			return []*errors.APIError{errors.InvalidFormField("stopTimeout", fmt.Sprintf("not a valid duration: %s", err))}
		}
	}
	if _, ok := form.ActionInfo[form.Action]; !ok {
		return []*errors.APIError{errors.InvalidFormField("action", "not a registered action")}
	}

	return nil
}

// JobSubmission is a form for scheduling an on-demand job to run.
type JobSubmission struct {
	Action string `json:"action" description:"The action which the job will perform"`
	// TODO: go-restful doesn't handle this type properly, we need to fix the docs
	Parameters map[string]string `json:"parameters" description:"Parameters to start the job with"`
	// RetriesLeft int               `json:"retriesLeft" description"The number of times to retry a job if it fails"`
	Deadline    string `json:"deadline",description:"After this amount of time has passed, a SIGTERM will be sent"`
	StopTimeout string `json:"stopTimeout",description:"This long after SIGTERM is sent, SIGKILL will be sent if the proccess is still alive"`
	ActionInfo  map[string]worker.ActionInfo
}

// ValidateJSON unmarshals this form from the given JSON and then validates
// the form, returning a list of any decoding or form validation errors.
func (form *JobSubmission) ValidateJSON(r io.Reader) []*errors.APIError {
	return validateJSONForm(r, form)
}

// Validate validates this JobSubmission form.
func (form *JobSubmission) Validate() (apiErrs []*errors.APIError) {
	if form.Deadline != "" {
		_, err := time.ParseDuration(form.Deadline)
		if err != nil {
			return []*errors.APIError{errors.InvalidFormField("deadline", fmt.Sprintf("not a valid duration: %s", err))}
		}
	}
	if form.StopTimeout != "" {
		_, err := time.ParseDuration(form.StopTimeout)
		if err != nil {
			return []*errors.APIError{errors.InvalidFormField("stopTimeout", fmt.Sprintf("not a valid duration: %s", err))}
		}
	}
	if _, ok := form.ActionInfo[form.Action]; !ok {
		return []*errors.APIError{errors.InvalidFormField("action", "not a registered action")}
	}

	return nil
}

func validateJSONForm(r io.Reader, form validator) []*errors.APIError {
	if err := json.NewDecoder(r).Decode(form); err != nil {
		return []*errors.APIError{errors.InvalidJSON(err)}
	}

	return form.Validate()
}

type validator interface {
	Validate() []*errors.APIError
}

// ActionConfigCreate is a form for creating a new cron
type ActionConfigCreate struct {
	Action           string            `json:"action" description:"The action to modify the config for"`
	Parameters       map[string]string `json:"parameters",description:"Extra parameters to pass to the job. The available parameters depend on the job. These are overwritten by any corresponding parameters set in the job itself."`
	MaxJobsPerWorker int               `json:"maxJobsPerWorker" description:"The maximum number of jobs to run on the same worker at the same time"`
	HeartbeatTimeout string            `json:"heartbeatTimeout" description:"The amount of time to wait before declaring a job with this action to be abandoned by its worker"`
	ActionInfo       map[string]worker.ActionInfo
}

// ValidateJSON unmarshals this form from the given JSON and then validates
// the form, returning a list of any decoding or form validation errors.
func (form *ActionConfigCreate) ValidateJSON(r io.Reader) []*errors.APIError {
	return validateJSONForm(r, form)
}

// Validate validates this ActionConfigCreate form.
func (form *ActionConfigCreate) Validate() (apiErrs []*errors.APIError) {
	if form.HeartbeatTimeout != "" {
		_, err := time.ParseDuration(form.HeartbeatTimeout)
		if err != nil {
			return []*errors.APIError{errors.InvalidFormField("heartbeatTimeout", fmt.Sprintf("not a valid duration: %s", err))}
		}
	}
	if _, ok := form.ActionInfo[form.Action]; !ok {
		return []*errors.APIError{errors.InvalidFormField("action", "not a registered action")}
	}

	return nil
}
