package schema

import (
	"io"
	"time"

	"github.com/docker/dhe-deploy/rethinkutil"
	rethink "gopkg.in/dancannon/gorethink.v2"
)

// Manager exports CRUDy methods for Accounts and Teams.
type JobrunnerManager interface {
	rethinkutil.Manager

	// CreateJob creates a new job using the values supplied by the given
	// job. The ID field of the job is set to a random UUID. The PK field
	// of the job is set to a hash of the cronID and scheduledAt time from
	// the job. If a job with the same PK already exists, the returned
	// error will be ErrJobExists.
	CreateJob(job *Job) error
	// ClaimJob attempts to claim the job with the given jobID for the
	// worker with the given workerID by performing a conditional update on
	// the job. If the job is still unclaimed, the job's workerID will be
	// set to the given workerID and its status will be set to "running"
	// and it lastUpdated time will be set to the current UTC time. If
	// claiming the job was successful, the job will be returned. If it is
	// nil, then the job has already been claimed by another worker.
	ClaimJob(jobID, workerID string, heartbeatExpiration time.Time) (*Job, error)
	// HeartbeatJob updates the job with the given jobID by setting its
	// last updated field to current time.
	HeartbeatJob(jobID string, heartbeatExpiration time.Time) error
	// UpdateJobStatus updates the job with the given jobID by setting its
	// status to the given status value.
	UpdateJobStatus(jobID, status string) error
	// GetJob retrieves the job with the given jobID. If no such job exists
	// the returned error will be ErrNoSuchJob.
	GetJob(jobID string) (*Job, error)
	// GetMostRecentlyScheduledJobs retreives a slice of the most recently
	// scheduled jobs sorted from newest to oldest.
	// The length of the slice will be at most limit
	GetMostRecentlyScheduledJobs(offset, limit uint) (jobs []Job, err error)
	// GetMostRecentlyScheduledJobsForWorker returns a slice of the most
	// recently scheduled jobs which are claimed by the worker with the
	// given workerID. The length of the slice will be at most limit. If
	// limit is 0, all jobs are returned.
	GetMostRecentlyScheduledJobsForWorker(workerID string, offset, limit uint) (jobs []Job, err error)
	// GetMostRecentlyScheduledJobsWithAction retreives a slice of the most
	// recently scheduled jobs with the given action. The length of the
	// slice will be at most limit. If limit is 0, all jobs are returned.
	GetMostRecentlyScheduledJobsWithAction(action string, offset, limit uint) (jobs []Job, err error)
	// GetLeastRecentlyScheduledJobsForWorker returns a slice of the most
	// recently scheduled jobs which are claimed by the worker with the
	// given workerID. The length of the slice will be at most limit. If
	// limit is 0, all jobs are returned.
	GetLeastRecentlyScheduledJobsForWorkerWithAction(workerID, action string, offset, limit uint) (jobs []Job, err error)
	// CountJobsWithActionStatus returns the number jobs with the given
	// action and status. This is a convenience method used to determine if
	// there are multiple jobs running which are performing the same action.
	CountJobsWithActionStatus(action, status string) (count uint, err error)
	// GetNextUnclaimedJob returns a job or nil depending on if there's a job queued up for this action
	// This is used to try to claim another job. It might have to be run multiple times in a loop if there's
	// contention for claiming the next job in the queue
	GetNextUnclaimedJob(parallelAction string) (job *Job, err error)
	// GetRecoverableJobChanges is meant to track the status of all pending
	// recoverable jobs so that all workers can use internal timers to track
	// jobs for heartbeat timeouts
	// This scales poorly because we have to listen to heartbeats on all
	// workers for all jobs on other workers and for cancellation events
	// on our own jobs. It's at least O(n*m), maybe O(n*m)
	// where n is number of outstanding jobs and
	// m is number of replicas. Every outstanding job needs to heartbeat
	// out to all workers so they can monitor its liveness.
	// Luckily we use this only for GC, so n <= 1 and m is officially allowed to be <= 7
	// In rethink 2.4 we might be able to use time based triggers instead:
	// https://github.com/rethinkdb/rethinkdb/issues/5813 or something like this
	// https://github.com/rethinkdb/rethinkdb/issues/3583 if/when that makes it in
	GetRecoverableJobChanges(recoverableActions []string) (changeStream <-chan JobChange, streamCloser io.Closer, err error)
	// GetOwnJobCancellations tracks currently running jobs for the worker
	// it watches for the status being changed to a cancellation request or
	// for the job being deleted
	GetOwnJobCancellations(workerID string) (changeStream <-chan JobChange, streamCloser io.Closer, err error)
	// GetNewJobChanges notifies workers about new jobs that are unclaimed
	GetNewJobChanges() (changeStream <-chan JobChange, streamCloser io.Closer, err error)
	// CancelJob creates a cancellation request for the given job. The job is not
	// guaranteed to be canceled immediately or to have its status set to
	// canceled. This is a best effort sort of thing.
	CancelJob(jobID string) error
	// DeleteJob deletes the job with the given jobID.
	DeleteJob(jobID string) error

	// GetJobLogs returns the logs for a job given a jobID.
	GetJobLogs(jobID string, offset, limit uint) ([]JobLog, error)
	// InsertJobLog inserts a log line for the given job
	InsertJobLog(jobID string, log string, lineNum int) error

	// ListCrons lists all scheduled crons in no particular order
	ListCrons() ([]Cron, error)
	// UpdateCron updates a cron with the given action and sets its
	// properties to the given values. It returns the updated cron.
	UpdateCron(cron *Cron) (*Cron, error)
	// GetCron retrieves the cron for the given action. If no such cron
	// exists the returned error will be ErrNoSuchCron.
	GetCron(action string) (*Cron, error)
	// GetCronChanges begins listening for any changes to cron configs.
	// Returns a channel on which the caller may receive a stream of
	// CronChange objects and an io.Closer which performs necessary cleanup
	// to end the stream's underlying goroutine. After closing, the
	// changeStream should be checked for a possible remaining value.
	GetCronChanges() (changeStream <-chan CronChange, streamCloser io.Closer, err error)
	// DeleteCron deletes the cron with the given action.
	DeleteCron(action string) error

	// ListActionConfigs lists all scheduled crons in no particular order
	ListActionConfigs() ([]ActionConfig, error)
	// UpdateActionConfig updates an action config with the given action and sets its
	// properties to the given values. It returns the updated action config.
	UpdateActionConfig(actionConfig *ActionConfig) (*ActionConfig, error)
	// GetActionConfig retrieves the action config for the given action. If no such action config
	// exists the returned error will be ErrNoSuchActionConfig.
	GetActionConfig(action string) (*ActionConfig, error)
	// SafeGetActionConfig returns a default config even in case of error
	SafeGetActionConfig(action string) (*ActionConfig, error)
	// DeleteActionConfig deletes the action config with the given action.
	DeleteActionConfig(action string) error
}

type jobrunnerManager struct {
	rethinkutil.Manager
}

var _ JobrunnerManager = &jobrunnerManager{}

// NewRethinkDBManager returns a new schema manager which connects to a
// RethinkDB cluster for storing data.
func NewJobrunnerManager(db string, session *rethink.Session) JobrunnerManager {
	return &jobrunnerManager{
		rethinkutil.NewGenericManager(
			db,
			session,
			[]rethinkutil.Table{
				JobsTable,
				CronsTable,
				JobLogsTable,
				ActionConfigsTable,
			},
		),
	}
}
