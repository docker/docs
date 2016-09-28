package schema

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/satori/go.uuid"
	rethink "gopkg.in/dancannon/gorethink.v2"
	rethinkencoding "gopkg.in/dancannon/gorethink.v2/encoding"
)

var (
	// ErrNoSuchWorker conveys that a worker with the given id does not
	// exist.
	ErrNoSuchWorker = errors.New("no such worker")
	// ErrNoSuchJob conveys that a job with the given id does not exist.
	ErrNoSuchJob = errors.New("no such job")
	// ErrJobExists conveys that a job with the same PK already exists.
	ErrJobExists = errors.New("job already exists")
	// ErrNoSuchCron conveys that a cron with the given name or id does not
	// exist.
	ErrNoSuchCron = errors.New("no such cron")
)

// A Worker represents a worker node in the jobs cluster.
type Worker struct {
	ID      string `gorethink:"id"`      // Randomly generated uuid for foreign references.
	Address string `gorethink:"address"` // Address at which an API server can contact the worker.
}

var workersTable = table{
	db:               dbName,
	name:             "workers",
	primaryKey:       "id", // Guarantees uniqueness. Quick lookups by id.
	secondaryIndexes: nil,  // No secondary indexes required.
}

// CreateOrUpdateWorker creates a worker with the given ID and address. If a
// worker with the given ID already exists, its address is set to the given
// address value.
func (m *manager) CreateOrUpdateWorker(id, address string) error {
	_, err := workersTable.Term().Get(id).Replace(Worker{
		ID:      id,
		Address: address,
	}).RunWrite(m.session)

	return err
}

// GetWorker retrieves the worker object with the given ID.
func (m *manager) GetWorker(id string) (*Worker, error) {
	cursor, err := workersTable.Term().Get(id).Run(m.session)
	if err != nil {
		return nil, fmt.Errorf("unable to query db: %s", err)
	}

	var worker Worker
	if err := cursor.One(&worker); err != nil {
		if err == rethink.ErrEmptyResult {
			return nil, ErrNoSuchWorker
		}

		return nil, fmt.Errorf("unable to get query result: %s", err)
	}

	return &worker, nil
}

// ListWorkers returns a slice of all registered workers.
func (m *manager) ListWorkers() ([]Worker, error) {
	cursor, err := workersTable.Term().Run(m.session)
	if err != nil {
		return nil, fmt.Errorf("unable to query db: %s", err)
	}

	workers := []Worker{}
	if err := cursor.All(&workers); err != nil {
		return nil, fmt.Errorf("unable to scan query results: %s", err)
	}

	return workers, nil
}

// DeleteWorker deletes the worker with the given ID.
func (m *manager) DeleteWorker(id string) error {
	_, err := workersTable.Term().Get(id).Delete().RunWrite(m.session)

	return err
}

// A Job represents a job which is performed by a worker.
type Job struct {
	ID          string    `gorethink:"id"`          // Randomly generated uuid for foreign references.
	PK          string    `gorethink:"pk"`          // Hash of (cronID + scheduledTime), primary key. Use a randomly generated ID if not a cron job.
	WorkerID    string    `gorethink:"workerID"`    // ID of the worker assigned this job, empty if unassigned.
	CronID      string    `gorethink:"cronID"`      // ID of the cron this job was spawned from. Empty if the job was created on-demand.
	Status      string    `gorethink:"status"`      // One of "waiting", "running", "done", or "errored".
	ScheduledAt time.Time `gorethink:"scheduledAt"` // The time at which this job was scheduled.
	LastUpdated time.Time `gorethink:"lastUpdated"` // Last time the status of this job was updated.
	Action      string    `gorethink:"action"`      // The action this job performs.
}

// Possible job statuses.
const (
	JobStatusWaiting  = "waiting"
	JobStatusRunning  = "running"
	JobStatusDone     = "done"
	JobStatusError    = "error"
	JobStatusCanceled = "canceled"
)

func computeJobPK(cronID string, scheduledTime time.Time) string {
	hash := sha256.New()

	if cronID == "" {
		// Not a cron job, use a randomly generated ID.
		cronID = uuid.NewV4().String()
	}

	hash.Write([]byte(cronID))
	hash.Write([]byte(scheduledTime.UTC().String()))

	return hex.EncodeToString(hash.Sum(nil))
}

var jobsTable = table{
	db:         dbName,
	name:       "jobs",
	primaryKey: "pk", // Guarantees uniqueness of (cron_id, scheduled time). Quick lookups.
	secondaryIndexes: map[string][]string{
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
func (m *manager) CreateJob(job *Job) error {
	job.ID = uuid.NewV4().String()
	job.PK = computeJobPK(job.CronID, job.ScheduledAt)

	if resp, err := jobsTable.Term().Insert(job).RunWrite(m.session); err != nil {
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
func (m *manager) ClaimJob(jobID, workerID string) (*Job, error) {
	response, err := jobsTable.Term().GetAllByIndex("id", jobID).Update(
		func(job rethink.Term) rethink.Term {
			return rethink.Branch(
				job.Field("workerID").Eq(""),
				map[string]interface{}{
					"workerID":    workerID,
					"status":      JobStatusRunning,
					"lastUpdated": rethink.Now(),
				},
				map[string]interface{}{},
			)
		}, rethink.UpdateOpts{ReturnChanges: true},
	).RunWrite(m.session)
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

// UpdateJobStatus updates the job with the given jobID by setting its status
// to the given status value.
func (m *manager) UpdateJobStatus(jobID, status string) error {
	if _, err := jobsTable.Term().GetAllByIndex("id", jobID).Update(
		map[string]interface{}{
			"status":      status,
			"lastUpdated": rethink.Now(),
		},
	).RunWrite(m.session); err != nil {
		return fmt.Errorf("unable to execute update query: %s", err)
	}

	return nil
}

// GetJob retrieves the job with the given jobID. If no such job exists the
// returned error will be ErrNoSuchJob.
func (m *manager) GetJob(jobID string) (*Job, error) {
	cursor, err := jobsTable.Term().GetAllByIndex("id", jobID).Run(m.session)
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

// GetMostRecentlyScheduledJobs retreives a slice of the most recently
// scheduled jobs. The length of the slice will be at most limit. If limit is
// 0, all jobs are returned starting from the given offset.
func (m *manager) GetMostRecentlyScheduledJobs(offset, limit uint) (jobs []Job, err error) {
	query := jobsTable.Term().OrderBy(
		rethink.OrderByOpts{Index: rethink.Desc("scheduledAt")},
	)

	return paginateJobsQuery(m.session, query, offset, limit)
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
func (m *manager) GetMostRecentlyScheduledJobsWithAction(action string, offset, limit uint) (jobs []Job, err error) {
	query := jobsTable.Term().OrderBy(
		rethink.OrderByOpts{Index: rethink.Desc("action_scheduledAt")},
	).Between(
		[]interface{}{action, rethink.MinVal},
		[]interface{}{action, rethink.MaxVal},
	)

	return paginateJobsQuery(m.session, query, offset, limit)
}

// GetMostRecentlyScheduledJobsForWorker returns a slice of the most recently
// scheduled jobs which are claimed by the worker with the given workerID. The
// length of the slice will be at most limit. If limit is 0, all jobs
// are returned starting from the given offset.
func (m *manager) GetMostRecentlyScheduledJobsForWorker(workerID string, offset, limit uint) (jobs []Job, err error) {
	query := jobsTable.Term().OrderBy(
		rethink.OrderByOpts{Index: rethink.Desc("workerID_scheduledAt")},
	).Between(
		[]interface{}{workerID, rethink.MinVal},
		[]interface{}{workerID, rethink.MaxVal},
	)

	return paginateJobsQuery(m.session, query, offset, limit)
}

// GetMostRecentlyScheduledJobsForWorkerWithAction returns a slice of the most
// recently scheduled jobs which are claimed by the worker with the given
// workerID and perform the given action. The length of the slice will be at
// most limit. If limit is 0, all jobs are returned starting from the given
// offset.
func (m *manager) GetMostRecentlyScheduledJobsForWorkerWithAction(workerID, action string, offset, limit uint) (jobs []Job, err error) {
	query := jobsTable.Term().OrderBy(
		rethink.OrderByOpts{Index: rethink.Desc("workerID_action_scheduledAt")},
	).Between(
		[]interface{}{workerID, action, rethink.MinVal},
		[]interface{}{workerID, action, rethink.MaxVal},
	)

	return paginateJobsQuery(m.session, query, offset, limit)
}

// CountJobsWithActionStatus returns the number jobs with the given action and
// status. This is a convenience method used to determine if there are multiple
// jobs running which are performing the same action.
func (m *manager) CountJobsWithActionStatus(action, status string) (count uint, err error) {
	cursor, err := jobsTable.Term().GetAllByIndex(
		"action_status", []interface{}{action, status},
	).Count().Run(m.session)
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

// GetUnclaimedJobChanges begins listening for any changes to jobs with an
// empty workerID field. Returns a channel on which the caller may receive a
// stream of JobChange objects and an io.Closer which performs necessary
// cleanup to end the stream's underlying goroutine. After closing, the
// changeStream should be checked for a possible remaining value.
func (m *manager) GetUnclaimedJobChanges() (changeStream <-chan JobChange, streamCloser io.Closer, err error) {
	cursor, err := jobsTable.Term().Between(
		// Doesn't matter what the scheduledAt time is, we just want
		// jobs where workerID is empty.
		[]interface{}{"", rethink.MinVal},
		[]interface{}{"", rethink.MaxVal},
		rethink.BetweenOpts{Index: "workerID_scheduledAt"},
	).Changes(
		rethink.ChangesOpts{IncludeInitial: true},
	).Run(m.session)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to query db: %s", err)
	}

	changes := make(chan JobChange)
	cursor.Listen(changes)

	return changes, cursor, nil
}

// DeleteJob deletes the job with the given jobID.
func (m *manager) DeleteJob(jobID string) error {
	if _, err := jobsTable.Term().GetAllByIndex("id", jobID).Delete().RunWrite(m.session); err != nil {
		return fmt.Errorf("unable to delete job from database: %s", err)
	}

	return nil
}
