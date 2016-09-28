package tmpresponses

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/dhe-deploy/pkg/jobrunner-framework/schema"
	cronUtilLib "github.com/robfig/cron"
)

// This package is temporary. It should be merged with enzi's responses package when the jobrunner is moved back into enzi

type JobLog struct {
	Data string `json:"data",description:"Contents of the log line"`
}

type JobLogs struct {
	JobLogs []JobLog `json:"jobLogs"`
}

// A Job object contains fields for a worker job.
type Job struct {
	ID          string            `json:"id"          description:"The ID of the job"`
	WorkerID    string            `json:"workerID"    description:"The ID of the worker which performed the job, unclaimed by a worker if empty"`
	Status      string            `json:"status"      description:"The current status of the job" enum:"waiting|running|done|canceled|errored"`
	ScheduledAt time.Time         `json:"scheduledAt" description:"The time at which this job was scheduled"`
	LastUpdated time.Time         `json:"lastUpdated" description:"The last time at which the status of this job was updated"`
	Action      string            `json:"action"      description:"The action this job performs"`
	RetriesLeft int               `json:"retriesLeft",description:"The number of times to retry the job if it fails"`
	Parameters  map[string]string `json:"parameters",description:"Extra parameters to pass to the job. The available parameters depend on the job."`
	Deadline    string            `json:"deadline",desrciption:"After this amount of time has passed, a SIGTERM will be sent"`
	StopTimeout string            `json:"stopTimeout",description:"This long after SIGTERM is sent, SIGKILL will be sent if the proccess is still alive"`
}

func MakeJob(job *schema.Job) Job {
	return Job{
		ID:          job.ID,
		WorkerID:    job.WorkerID,
		Status:      job.Status,
		ScheduledAt: job.ScheduledAt,
		LastUpdated: job.LastUpdated,
		Action:      job.Action,
		RetriesLeft: job.RetriesLeft,
		Parameters:  job.Parameters,
		Deadline:    job.Deadline,
		StopTimeout: job.StopTimeout,
	}
}

// A Jobs object is used as a response for a list of jobs.
type Jobs struct {
	Jobs []Job `json:"jobs"`
}

func MakeJobs(jobs []schema.Job) Jobs {
	jobObjects := make([]Job, len(jobs))
	for i := range jobs {
		jobObjects[i] = MakeJob(&jobs[i])
	}

	return Jobs{
		Jobs: jobObjects,
	}
}

// A Cron is used to run a job on a regular schedule.
type Cron struct {
	ID          string            `json:"id",description:"Randomly generated UUID for foreign references."`
	Action      string            `json:"action",description:"The action to be performed by jobs spawned from this cron."`
	Schedule    string            `json:"schedule",description:"The schedule for this cron as a cronspec string: (seconds) (minutes) (hours) (day of month) (month) (day of week) or @hourly, @weekly, etc."`
	Retries     int               `json:"retries",description:"The number of times to retry the job if it fails"`
	Parameters  map[string]string `json:"parameters",description:"Extra parameters to pass to the job. The available parameters depend on the job."`
	Deadline    string            `json:"deadline",desrciption:"After this amount of time has passed, a SIGTERM will be sent"`
	StopTimeout string            `json:"stopTimeout",description:"This long after SIGTERM is sent, SIGKILL will be sent if the proccess is still alive"`
	NextRun     *time.Time        `json:"nextRun",description:"The next time the job will run."`
}

// MakeCron returns a Cron response object for the given worker cron.
func MakeCron(cron *schema.Cron) Cron {
	// set the next run field
	var nextRun *time.Time

	// if there is no schedule, there is no next run time
	if cron.Schedule != "" {
		schedule, err := cronUtilLib.Parse(cron.Schedule)
		if err != nil {
			// dont want to error because the user put in invalid data
			// just skip the NextRun field
			log.Errorf("unable to parse cron schedule: %s, error: %s", cron.Schedule, err)
		} else {
			n := schedule.Next(time.Now())
			nextRun = &n
		}
	}

	return Cron{
		ID:          cron.ID,
		Action:      cron.Action,
		Schedule:    cron.Schedule,
		Retries:     cron.Retries,
		Parameters:  cron.Parameters,
		Deadline:    cron.Deadline,
		StopTimeout: cron.StopTimeout,
		NextRun:     nextRun,
	}
}

// A Crons object is used as a response for a list of crons.
type Crons struct {
	Crons []Cron `json:"crons"`
}

func MakeCrons(crons []schema.Cron) Crons {
	cronObjects := make([]Cron, len(crons))
	for i := range crons {
		cronObjects[i] = MakeCron(&crons[i])
	}

	return Crons{
		Crons: cronObjects,
	}
}

// ActionConfig represents action-specific configs
type ActionConfig struct {
	ID               string            `json:"id",description:"Randomly generated UUID for foreign references."`
	Action           string            `json:"action",description:"The action this config refers to."`
	MaxJobsPerWorker int               `json:"maxJobsPerWorker",description:"The maximum number of jobs to run on the same worker at the same time; 0 means unlimited"`
	HeartbeatTimeout string            `json:"heartbeatTimeout",description:"The amount of time to wait before declaring a job with this action to be abandoned by its worker"`
	Parameters       map[string]string `json:"parameters",description:"Extra parameters to pass to the job. The available parameters depend on the job."`
}

// MakeActionConfig returns a ActionConfig response object for the given worker actionConfig.
func MakeActionConfig(actionConfig *schema.ActionConfig) ActionConfig {
	return ActionConfig{
		ID:               actionConfig.ID,
		Action:           actionConfig.Action,
		MaxJobsPerWorker: actionConfig.MaxJobsPerWorker,
		HeartbeatTimeout: actionConfig.HeartbeatTimeout,
		Parameters:       actionConfig.Parameters,
	}
}

// A ActionConfigs object is used as a response for a list of actionConfigs.
type ActionConfigs struct {
	ActionConfigs []ActionConfig `json:"actionConfigs"`
}

func MakeActionConfigs(actionConfigs []schema.ActionConfig) ActionConfigs {
	actionConfigObjects := make([]ActionConfig, len(actionConfigs))
	for i := range actionConfigs {
		actionConfigObjects[i] = MakeActionConfig(&actionConfigs[i])
	}

	return ActionConfigs{
		ActionConfigs: actionConfigObjects,
	}
}
