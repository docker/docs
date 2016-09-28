package schema

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/docker/dhe-deploy/rethinkutil"
	"github.com/satori/go.uuid"
	rethink "gopkg.in/dancannon/gorethink.v2"
	rethinkencoding "gopkg.in/dancannon/gorethink.v2/encoding"
)

var (
	// ErrNoSuchJob conveys that a job with the given id does not exist.
	ErrNoSuchJob = errors.New("no such job")
	// ErrJobExists conveys that a job with the same PK already exists.
	ErrJobExists = errors.New("job already exists")
)

// A Job represents a job which is performed by a worker.
type Job struct {
	ID                  string            `gorethink:"id"`                  // Randomly generated uuid for foreign references.
	PK                  string            `gorethink:"pk"`                  // Hash of (cronID + scheduledTime + retriesLeft + recoveryFromID), primary key. Use a randomly generated ID if not a cron job.
	RecoveryFromID      string            `gorethink:"recoveryFromID"`      // This job is the recovery job from a failed job with this ID, empty if n/a
	WorkerID            string            `gorethink:"workerID"`            // ID of the worker assigned this job, empty if unassigned.
	CronID              string            `gorethink:"cronID"`              // ID of the cron this job was spawned from. Empty if the job was created on-demand.
	Status              string            `gorethink:"status"`              // One of "waiting", "running", "done", or "errored".
	ScheduledAt         time.Time         `gorethink:"scheduledAt"`         // The time at which this job was scheduled.
	LastUpdated         time.Time         `gorethink:"lastUpdated"`         // Last time the status of this job was updated.
	HeartbeatExpiration time.Time         `gorethink:"heartbeatExpiration"` // When this worker should be considered dead
	Action              string            `gorethink:"action"`              // The action this job performs.
	RetriesLeft         int               `gorethink:"retriesLeft"`         // Number of retries left
	Parameters          map[string]string `gorethink:"parameters"`          // Parameters to start the job with; they are sent as environment variables in all uppercase, prefixed with PARAM_
	Deadline            string            `gorethink:"deadline"`            // After this amount of time has passed, a SIGTERM will be sent
	StopTimeout         string            `gorethink:"stopTimeout"`         // This long after SIGTERM is sent, SIGKILL will be sent if the proccess is still alive
}

// Possible job statuses.
const (
	JobStatusWaiting = "waiting"
	JobStatusRunning = "running"
	JobStatusDone    = "done"
	JobStatusError   = "error"
	// another worker has noticed that this worker has died
	JobStatusWorkerDead = "worker_dead"
	// graceful shutdown of the worker
	JobStatusWorkerShutdown = "worker_shutdown"
	// This status is meant to be used for when a worker was proclaimed dead but
	// then came back and was unhappy with the proclamation. This can happen if
	// it was separated from the network
	JobStatusWorkerResurrection = "worker_resurrection"
	// This status is given to exclusive jobs when they are picked up by a worker, but
	// the worker finds out that another job is running at the same time. In this case
	// the job was never started
	JobStatusConflict      = "conflict"
	JobStatusCancelRequest = "cancel_request"
	JobStatusCanceled      = "canceled"
	JobStatusDeleted       = "deleted"
	JobStatusTimeout       = "timeout"
)

var FailedStatuses = []string{JobStatusError, JobStatusTimeout, JobStatusWorkerDead, JobStatusWorkerResurrection, JobStatusConflict}
var FinishedStatuses = append(FailedStatuses, JobStatusDone, JobStatusCanceled, JobStatusDeleted)

func StatusIsFailed(status string) bool {
	for _, s := range FailedStatuses {
		if status == s {
			return true
		}
	}
	return false
}
func StatusIsFinished(status string) bool {
	for _, s := range FinishedStatuses {
		if status == s {
			return true
		}
	}
	return false
}

func computeJobPK(job *Job) string {
	hash := sha256.New()

	if job.CronID == "" && job.RecoveryFromID == "" {
		// Not a cron job or recovery job, use a randomly generated ID.
		hash.Write([]byte(uuid.NewV4().String()))
		return hex.EncodeToString(hash.Sum(nil))
	}

	hash.Write([]byte(job.CronID))
	hash.Write([]byte(fmt.Sprintf("%d", job.RetriesLeft)))
	hash.Write([]byte(job.RecoveryFromID))
	hash.Write([]byte(job.ScheduledAt.UTC().String()))

	return hex.EncodeToString(hash.Sum(nil))
}

var JobsTable = rethinkutil.Table{
	Name:       "jobs",
	PrimaryKey: "pk", // Guarantees uniqueness of (cron_id, scheduled time). Quick lookups.
	SecondaryIndexes: map[string][]string{
		"id":                          nil,                                   // For quick lookups by ID.
		"scheduledAt":                 nil,                                   // For quickly listing jobs ordered by scheduledAt time.
		"workerID_scheduledAt":        {"workerID", "scheduledAt"},           // For quickly listing jobs with a particular workerID, ordered by scheduledAt time.
		"action_scheduledAt":          {"action", "scheduledAt"},             // For quickly listing jobs with a particilar action, ordered by scheduledAt time.
		"workerID_action_scheduledAt": {"workerID", "action", "scheduledAt"}, // For quickly listing jobs with a particular workerID and action, ordered by scheduledAt time.
		"action_status":               {"action", "status"},                  // For quickly listing jobs with a particular action and status.
	},
}

// CreateJob inserts a new job into the *jobs* table of the database using the
// values supplied by the given job. The ID field of the job is set to a random
// UUID. The PK field of the job is set to a hash of the cronID and scheduledAt
// time from the job. If a job with the same PK already exists, the returned
// error will be ErrJobExists.
func (m *jobrunnerManager) CreateJob(job *Job) error {
	job.ID = uuid.NewV4().String()
	job.PK = computeJobPK(job)

	if resp, err := JobsTable.Term(m.DB()).Insert(job).RunWrite(m.Session()); err != nil {
		if isDuplicatePrimaryKeyErr(resp) {
			return ErrJobExists
		}

		return fmt.Errorf("unable to insert job into db: %s", err)
	}

	return nil
}

// ClaimJob attempts to claim the job with the given jobID for the worker with
// the given workerID by performing a conditional update on the job. If the job
// is still unclaimed, the job's workerID will be set to the given workerID and
// its status will be set to "running" and it lastUpdated time will be set to
// the current UTC time. If claiming the job was successful, the job will be
// returned. If it is nil, then the job has already been claimed by another
// worker.
func (m *jobrunnerManager) ClaimJob(jobID, workerID string, heartbeatExpiration time.Time) (*Job, error) {
	response, err := JobsTable.Term(m.DB()).GetAllByIndex("id", jobID).Update(
		func(job rethink.Term) rethink.Term {
			return rethink.Branch(
				job.Field("workerID").Eq(""),
				map[string]interface{}{
					"workerID":            workerID,
					"status":              JobStatusRunning,
					"lastUpdated":         rethink.Now(),
					"heartbeatExpiration": heartbeatExpiration,
				},
				map[string]interface{}{},
			)
		}, rethink.UpdateOpts{ReturnChanges: true},
	).RunWrite(m.Session())
	if err != nil {
		return nil, fmt.Errorf("unable to execute update query: %s", err)
	}

	if response.Replaced == 0 {
		// The job was not claimed.
		return nil, nil
	}

	// The update succeeded in claiming the job. We need to convert
	// the "new_value" from a simple interface{} type (actually a
	// map[string]interface{}) to our Job struct type using the gorethink
	// decoder.
	encodedJob := response.Changes[0].NewValue

	var job Job
	if err := rethinkencoding.Decode(&job, encodedJob); err != nil {
		return nil, fmt.Errorf("unable to decode new job value: %s", err)
	}

	return &job, nil

}

// HeartbeatJob updates the job with the given jobID by setting its last
// updated time to the current time
func (m *jobrunnerManager) HeartbeatJob(jobID string, heartbeatExpiration time.Time) error {
	if _, err := JobsTable.Term(m.DB()).GetAllByIndex("id", jobID).Update(
		map[string]interface{}{
			"heartbeatExpiration": heartbeatExpiration,
		},
	).RunWrite(m.Session()); err != nil {
		return fmt.Errorf("unable to execute update query: %s", err)
	}

	return nil
}

// UpdateJobStatus updates the job with the given jobID by setting its status
// to the given status value.
func (m *jobrunnerManager) UpdateJobStatus(jobID, status string) error {
	if _, err := JobsTable.Term(m.DB()).GetAllByIndex("id", jobID).Update(
		map[string]interface{}{
			"status":      status,
			"lastUpdated": rethink.Now(),
		},
	).RunWrite(m.Session()); err != nil {
		return fmt.Errorf("unable to execute update query: %s", err)
	}

	return nil
}

// GetJob retrieves the job with the given jobID. If no such job exists the
// returned error will be ErrNoSuchJob.
func (m *jobrunnerManager) GetJob(jobID string) (*Job, error) {
	cursor, err := JobsTable.Term(m.DB()).GetAllByIndex("id", jobID).Run(m.Session())
	if err != nil {
		return nil, fmt.Errorf("unable to query db: %s", err)
	}

	var job Job
	if err := cursor.One(&job); err != nil {
		if err == rethink.ErrEmptyResult {
			return nil, ErrNoSuchJob
		}

		return nil, fmt.Errorf("unable to get query result: %s", err)
	}

	return &job, nil
}

// GetMostRecentlyScheduledJobs retrieves a slice of the most recently
// scheduled jobs. The length of the slice will be at most limit. If limit is
// 0, all jobs are returned starting from the given offset.
// TODO: do we want to return only jobs that are in the running state?
func (m *jobrunnerManager) GetMostRecentlyScheduledJobs(offset, limit uint) (jobs []Job, err error) {
	query := JobsTable.Term(m.DB()).OrderBy(
		rethink.OrderByOpts{Index: rethink.Desc("scheduledAt")},
	)

	return paginateJobsQuery(m.Session(), query, offset, limit)
}

// paginateJobsQuery runs a paginated version of the given query using the
// given offset and limit. If limit is zero, all results are returned after
// skipping the first offset results.
func paginateJobsQuery(session *rethink.Session, query rethink.Term, offset, limit uint) (jobs []Job, err error) {
	if limit > 0 {
		query = query.Slice(offset, offset+limit)
	} else {
		query = query.Skip(offset)
	}

	cursor, err := query.Run(session)
	if err != nil {
		return nil, fmt.Errorf("unable to query db: %s", err)
	}

	jobs = []Job{}
	if err := cursor.All(&jobs); err != nil {
		return nil, fmt.Errorf("unable to scan query results: %s", err)
	}

	return jobs, nil
}

// GetMostRecentlyScheduledJobsWithAction retreives a slice of the most
// recently scheduled jobs with the given action. The length of the slice will
// be at most limit. If limit is 0, all jobs are returned starting from the
// given offset.
func (m *jobrunnerManager) GetMostRecentlyScheduledJobsWithAction(action string, offset, limit uint) (jobs []Job, err error) {
	query := JobsTable.Term(m.DB()).OrderBy(
		rethink.OrderByOpts{Index: rethink.Desc("action_scheduledAt")},
	).Between(
		[]interface{}{action, rethink.MinVal},
		[]interface{}{action, rethink.MaxVal},
	)

	return paginateJobsQuery(m.Session(), query, offset, limit)
}

// GetMostRecentlyScheduledJobsForWorker returns a slice of the most recently
// scheduled jobs which are claimed by the worker with the given workerID. The
// length of the slice will be at most limit. If limit is 0, all jobs
// are returned starting from the given offset.
func (m *jobrunnerManager) GetMostRecentlyScheduledJobsForWorker(workerID string, offset, limit uint) (jobs []Job, err error) {
	query := JobsTable.Term(m.DB()).OrderBy(
		rethink.OrderByOpts{Index: rethink.Desc("workerID_scheduledAt")},
	).Between(
		[]interface{}{workerID, rethink.MinVal},
		[]interface{}{workerID, rethink.MaxVal},
	)

	return paginateJobsQuery(m.Session(), query, offset, limit)
}

// GetMostRecentlyScheduledJobsForWorkerWithAction returns a slice of the most
// recently scheduled jobs which are claimed by the worker with the given
// workerID and perform the given action. The length of the slice will be at
// most limit. If limit is 0, all jobs are returned starting from the given
// offset.
func (m *jobrunnerManager) GetLeastRecentlyScheduledJobsForWorkerWithAction(workerID, action string, offset, limit uint) (jobs []Job, err error) {
	query := JobsTable.Term(m.DB()).OrderBy(
		rethink.OrderByOpts{Index: rethink.Asc("workerID_action_scheduledAt")},
	).Between(
		[]interface{}{workerID, action, rethink.MinVal},
		[]interface{}{workerID, action, rethink.MaxVal},
	)

	return paginateJobsQuery(m.Session(), query, offset, limit)
}

// CountJobsWithActionStatus returns the number jobs with the given action and
// status. This is a convenience method used to determine if there are multiple
// jobs running which are performing the same action.
func (m *jobrunnerManager) CountJobsWithActionStatus(action, status string) (count uint, err error) {
	cursor, err := JobsTable.Term(m.DB()).GetAllByIndex(
		"action_status", []interface{}{action, status},
	).Count().Run(m.Session())
	if err != nil {
		return 0, fmt.Errorf("unable to query db: %s", err)
	}

	if err := cursor.One(&count); err != nil {
		return 0, fmt.Errorf("unable to scan query results: %s", err)
	}

	return count, nil
}

// JobChange is used to deliver old and new values of a job as part of a
// changes stream.
type JobChange struct {
	OldValue *Job `gorethink:"old_val"`
	NewValue *Job `gorethink:"new_val"`
}

func (m *jobrunnerManager) GetRecoverableJobChanges(recoverableActions []string) (<-chan JobChange, io.Closer, error) {
	// there's a bit of pointer shuffling here just because we want `ors` to be nilable
	var ors *rethink.Term
	if len(recoverableActions) > 0 {
		tmp := rethink.Row.Field("action").Eq(recoverableActions[0])
		ors = &tmp
		if len(recoverableActions) > 1 {
			// all events for jobs with recoverable actions
			for _, action := range recoverableActions[1:] {
				tmp2 := rethink.Or(
					*ors,
					rethink.Row.Field("action").Eq(action),
				)
				ors = &tmp2
			}
		}
	}

	term := JobsTable.Term(m.DB())
	if ors != nil {
		term = term.Filter(*ors)
	}

	cursor, err := term.Changes(
		// TODO: determine the value of this?
		rethink.ChangesOpts{IncludeInitial: false},
	).Filter(
		term,
	).Run(m.Session())
	if err != nil {
		return nil, nil, fmt.Errorf("unable to query db: %s", err)
	}

	changes := make(chan JobChange)
	cursor.Listen(changes)

	return changes, cursor, nil
}

func (m *jobrunnerManager) GetNextUnclaimedJob(parallelAction string) (*Job, error) {
	jobs, err := m.GetLeastRecentlyScheduledJobsForWorkerWithAction("", parallelAction, 0, 1)
	if err != nil {
		return nil, err
	}
	if len(jobs) > 1 {
		return &jobs[0], nil
	}
	return nil, nil
}

func (m *jobrunnerManager) GetOwnJobCancellations(workerID string) (<-chan JobChange, io.Closer, error) {
	// XXX: we ask for all changes about our own jobs because it's awkward to write
	// a filter for just the events we want. If someone feels the urge to do it,
	// be my guest, but it's not a huge difference.
	// The changes we really care about are:
	// 1. when a job's status is set to cancel request
	// 2. when a job is deleted
	cursor, err := JobsTable.Term(m.DB()).Filter(map[string]interface{}{
		"workerID": workerID,
	}).Changes(
		// TODO: determine the value of this?
		rethink.ChangesOpts{IncludeInitial: false},
	).Run(m.Session())
	if err != nil {
		return nil, nil, fmt.Errorf("unable to query db: %s", err)
	}

	changes := make(chan JobChange)
	cursor.Listen(changes)

	return changes, cursor, nil
}

func (m *jobrunnerManager) GetNewJobChanges() (<-chan JobChange, io.Closer, error) {
	// TODO: will this trigger only for inserts or inserts and claims?
	cursor, err := JobsTable.Term(m.DB()).Filter(map[string]interface{}{
		"workerID": "",
	}).Changes(
		rethink.ChangesOpts{IncludeInitial: true},
	).Run(m.Session())
	if err != nil {
		return nil, nil, fmt.Errorf("unable to query db: %s", err)
	}

	changes := make(chan JobChange)
	cursor.Listen(changes)

	return changes, cursor, nil
}

// CancelJob cancels the job with the given jobID.
func (m *jobrunnerManager) CancelJob(jobID string) error {
	// We make sure the job is in a cancellable state before trying to cancel it. The only such state right
	// now is "running". This needs to be an atomic update to make sure the job doesn't change status while being canceled.
	if _, err := JobsTable.Term(m.DB()).GetAllByIndex("id", jobID).Update(func(job rethink.Term) rethink.Term {
		return rethink.Branch(
			job.Field("status").Eq(JobStatusRunning),
			map[string]interface{}{
				"status": JobStatusCanceled,
			},
			map[string]interface{}{},
		)
	}).RunWrite(m.Session()); err != nil {
		return fmt.Errorf("unable to cancel job: %s", err)
	}

	return nil
}

// DeleteJob deletes the job with the given jobID.
func (m *jobrunnerManager) DeleteJob(jobID string) error {
	if _, err := JobsTable.Term(m.DB()).GetAllByIndex("id", jobID).Delete().RunWrite(m.Session()); err != nil {
		return fmt.Errorf("unable to delete job from database: %s", err)
	}

	return nil
}

func isDuplicatePrimaryKeyErr(resp rethink.WriteResponse) bool {
	return strings.HasPrefix(resp.FirstError, "Duplicate primary key")
}
