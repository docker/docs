package worker

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/docker/distribution/context"
	"github.com/docker/orca/enzi/schema"
)

// The currently available actions to perform in a job.
const (
	ActionLdapSync  = "ldap-sync"
	ActionCleanupDB = "cleanup-db"
	ActionTick      = "tick"
)

type actionInfo struct {
	exclusive bool
}

// RegisteredActions is the set of valid actions available to run in a job and
// their corresponding configuration.
var RegisteredActions = map[string]actionInfo{
	ActionLdapSync: {
		exclusive: true,
	},
	ActionCleanupDB: {
		exclusive: true,
	},
	ActionTick: {
		exclusive: true,
	},
}

func (w *worker) claimJobs() {
	context.GetLogger(w.ctx).Info("claiming jobs")

	for w.claimJobsLoop() {
		// Lost connection to the unclaimed jobs changes stream.
		// Wait 5 seconds before trying again.
		time.Sleep(5 * time.Second)
	}
}

func (w *worker) claimJobsLoop() bool {
	changeStream, streamCloser, err := w.schemaMgr.GetUnclaimedJobChanges()
	if err != nil {
		context.GetLogger(w.ctx).Errorf("unable to get stream of unclaimed job changes: %s", err)
		return true // Try again.
	}

	defer streamCloser.Close()

	for {
		select {
		case wg := <-w.signalEndJobs:
			defer wg.Done()

			// The main goroutine has received a signal to end the
			// process. Cancel all jobs and use the given waitgroup
			// to signal that cleanup is complete. The sender will
			// still need to wait for each job to gracefully
			// shutdown.
			w.cancelAllJobs()

			return false // Don't try again.

		case change, ok := <-changeStream:
			if !ok {
				context.GetLogger(w.ctx).Info("unclaimed job change stream unexpectedly closed")
				return true // Try again.
			}

			context.GetLogger(w.ctx).Debugf("old job value: %#v", change.OldValue)
			context.GetLogger(w.ctx).Debugf("new job value: %#v", change.NewValue)

			if change.OldValue != nil {
				continue
			}

			newJob := change.NewValue
			if newJob.WorkerID == "" {
				w.claimJob(newJob)
			}
		}
	}
}

func (w *worker) claimJob(job *schema.Job) {
	w.jobsLock.Lock()
	defer w.jobsLock.Unlock()

	// Attempt to claim this job.
	claimedJob, err := w.schemaMgr.ClaimJob(job.ID, w.id)
	if err != nil {
		context.GetLogger(w.ctx).Errorf("unable to claim job: %s", err)
		return
	}

	if claimedJob == nil {
		context.GetLogger(w.ctx).Infof("job %s was already claimed", job.ID)
		return
	}

	subCtx := context.WithValue(w.ctx, "jobID", claimedJob.ID)
	subCtx = context.WithLogger(subCtx, context.GetLogger(subCtx, "jobID"))

	w.addJob(subCtx, claimedJob.ID, claimedJob.Action)
}

// DeleteJob attempts to delete the job with the given ID.
func (w *worker) DeleteJob(jobID string) {
	w.jobsLock.Lock()
	defer w.jobsLock.Unlock()

	w.deleteJobLocked(jobID)
}

func (w *worker) deleteJobLocked(jobID string) {
	job, ok := w.jobs[jobID]
	if !ok {
		return // No such job; must have already been deleted.
	}

	job.cancel()
	job.Wait()

	if err := os.RemoveAll(w.jobDirPath(jobID)); err != nil {
		context.GetLogger(w.ctx).Errorf("unable to remove log directory for job %s: %s", jobID, err)
	}

	delete(w.jobs, jobID)
}

// CancelJob signals the job with the given jobID to cancel itself and waits
// for it to gracefully shutdown.
func (w *worker) CancelJob(jobID string) {
	w.jobsLock.Lock()
	defer w.jobsLock.Unlock()

	job, ok := w.jobs[jobID]
	if !ok {
		return // No such job; may have already completed.
	}

	job.cancel()
	job.Wait()
}

// cancelAllJobs signals each currently running job to cancel itself and waits
// for all running jobs to gracefully shutdown.
func (w *worker) cancelAllJobs() {
	context.GetLogger(w.ctx).Info("canceling all running jobs")

	w.jobsLock.Lock()
	defer w.jobsLock.Unlock()

	for _, job := range w.jobs {
		job.cancel()
	}

	for _, job := range w.jobs {
		job.Wait()
	}
}

// addJob creates and starts a new job with the given jobID to run the given
// action. The job is added to this worker's set of currently running jobs. If
// there is an error initializing the job, the job status is updated to
// the errored status.
func (w *worker) addJob(ctx context.Context, jobID, action string) {
	job := &job{
		ctx:    ctx,
		worker: w,
		id:     jobID,
		action: action,
		status: schema.JobStatusRunning,
	}

	if err := job.init(); err != nil {
		logger := context.GetLogger(ctx)
		logger.Errorf("failed to initialize job: %s", err)

		if err := w.schemaMgr.UpdateJobStatus(jobID, schema.JobStatusError); err != nil {
			logger.Errorf("unable to update job status: %s", err)
		}

		return
	}

	w.jobs[jobID] = job

	job.Add(1)

	go job.run()
}

type job struct {
	ctx    context.Context
	worker *worker
	id     string
	action string

	sync.Mutex
	sync.WaitGroup

	cmd *exec.Cmd

	isComplete bool
	isCanceled bool
	status     string
}

func (j *job) logFilename() string {
	return filepath.Join(j.worker.jobDirPath(j.id), "log")
}

// GetJobLogs returns an io.ReadCloser for the logs of the job with the given
// ID. If no such job exists, (nil, nil) is returned.
func (w *worker) GetJobLogs(jobID string) (io.ReadCloser, error) {
	w.jobsLock.Lock()
	defer w.jobsLock.Unlock()

	job, exists := w.jobs[jobID]
	if !exists {
		return nil, nil
	}

	return os.Open(job.logFilename())
}

// init handles creating the log file for this job, and setting up the command
// to run.
func (j *job) init() error {
	if err := os.MkdirAll(j.worker.jobDirPath(j.id), os.FileMode(0755)); err != nil {
		return fmt.Errorf("unable to make job directory: %s", err)
	}

	logFile, err := os.OpenFile(j.logFilename(), os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.FileMode(0644))
	if err != nil {
		return fmt.Errorf("unable to create log file: %s", err)
	}

	actionInfo, ok := RegisteredActions[j.action]
	if !ok {
		logFile.Close()

		return fmt.Errorf("unregistered action: %s", j.action)
	}

	if actionInfo.exclusive {
		// Don't run the job if there's already another running job
		// performing this aciton.
		numJobs, err := j.worker.schemaMgr.CountJobsWithActionStatus(j.action, schema.JobStatusRunning)
		if err != nil {
			return fmt.Errorf("unable to check for other running jobs with the same action: %s", err)
		}

		if numJobs > 1 {
			// There's at least one other job (besides this one)
			// which is performing this exclusive action already.
			// Refuse to run.
			return fmt.Errorf("unable to perform exclusive action: %s - %d other jobs running", j.action, numJobs)
		}
	}

	// Jobs are ran with the same arguments as this worker process, but
	// with the name of the job as an aditional argument. This ensures that
	// the job is started with the same options as this worker process.
	// TODO: Using enviroment variables would be more elegant.
	args := make([]string, len(os.Args), len(os.Args)+1)
	copy(args, os.Args)
	args = append(args, j.action)

	j.cmd = exec.Command(args[0], args[1:]...)
	j.cmd.Stdout = logFile
	j.cmd.Stderr = logFile

	return nil
}

// run starts this job's subprocess and waits for this job to be completed. It
// blocks until the job subprocess has either completed successfully, errored,
// or been canceled. This function should only be called once.
func (j *job) run() {
	defer j.Done() // Signals to anyone waiting for this job to complete.

	// Close the log file when the job completes.
	defer j.cmd.Stdout.(*os.File).Close()

	// Run blocks until the subprocess exits.
	err := j.cmd.Run()

	j.Lock()
	defer j.Unlock()

	context.GetLogger(j.ctx).Info("job process ended")

	j.status = schema.JobStatusDone
	if err != nil {
		// The process was either canceled or errored on its own.
		context.GetLogger(j.ctx).Errorf("job exited with error: %s", err)

		j.status = schema.JobStatusError
		if j.isCanceled {
			j.status = schema.JobStatusCanceled
		}
	}

	// Mark the job as complete if it is not already.
	j.isComplete = true

	if err := j.worker.schemaMgr.UpdateJobStatus(j.id, j.status); err != nil {
		context.GetLogger(j.ctx).Errorf("unable to update job status: %s", err)
	}

	context.GetLogger(j.ctx).Info("job completed")
}

// cancel ends this job by killing the job's subprocess. The caller should call
// Wait() on this job to wait for the job to shutdown gracefully.
func (j *job) cancel() {
	j.Lock()
	defer j.Unlock()

	j.cancelLocked()
}

func (j *job) cancelLocked() {
	if j.isComplete || j.isCanceled {
		return
	}

	context.GetLogger(j.ctx).Info("job canceled")

	j.isCanceled = true

	// Killing the process will unblock the run() goroutine.
	j.cmd.Process.Kill()
}

func (j *job) reconcileState(assignedJob schema.Job) {
	j.Lock()
	defer j.Unlock()

	if assignedJob.Status == j.status {
		return // Everything is consistent.
	}

	// If the job status in the DB is not "running" then cancel the local
	// job. It may have already completed. Once it is complete, its local
	// status will also not be "running".
	if assignedJob.Status != schema.JobStatusRunning {
		j.cancelLocked()
		return
	}

	// Here we know that the status in the DB is "running" but the local
	// state is not running which means the job has completed (only the
	// run goroutine could have updated the local job status). However, the
	// given assigned job status might be slightly stale if the job
	// completed and updated the status after it was retreived from the
	// database. Either way, if we update the record in the DB to the
	// current local status it should just be a nop if there was a race.

	context.GetLogger(j.ctx).Warnf("updating stray job status to %s", j.status)

	if err := j.worker.schemaMgr.UpdateJobStatus(j.id, j.status); err != nil {
		context.GetLogger(j.ctx).Errorf("unable to update job status: %s", err)
	}
}
