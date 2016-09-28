package worker

import (
	"bufio"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/WatchBeam/clock"
	"github.com/docker/dhe-deploy/pkg/jobrunner-framework/constants"
	"github.com/docker/dhe-deploy/pkg/jobrunner-framework/schema"
	"github.com/docker/distribution/context"
)

type ActionInfo struct {
	Command   string   // The command to execute. Can be relative to $PATH or absolute
	Args      []string // Command line arguments
	Exclusive bool     // Exclusive jobs refuse to start if another worker is already running another job with the same action
	// cleanup jobs cannot themselves have cleanup jobs
	// NOTE: cleanup jobs should not NOT be called regularly and they should not be chained because these behaviours are untested
	// if the cleanup job is called regularly and the recoverable job is exclusive, it will not be able to run while an instance
	// of the cleanup job is running either
	CleanupJob *schema.Job // Another job to start in case this job fails; useful for recovering the state of the cluster if it was modified by a job. Ex. returning registries to rw mode after GC. This is run only if there are no retries left
}

type job struct {
	ctx        context.Context
	worker     *worker
	actionInfo ActionInfo
	job        *schema.Job

	sync.Mutex
	sync.WaitGroup

	newExecutor   NewExecutor
	executor      Executor
	goroutineHack func()
	clock         clock.Clock

	isComplete bool
	logLineNum int
}

// init handles creating the log file for this job, and setting up the command
// to run.
func (j *job) init() error {
	j.executor = j.newExecutor()

	return nil
}

func (j *job) writeLogs(logs ...string) {
	// flush extra logs
	for _, logLine := range logs {
		err := j.worker.schemaMgr.InsertJobLog(j.job.ID, logLine, j.logLineNum)
		if err != nil {
			context.GetLogger(j.ctx).Warnf("failed to save log to db: %s", err)
		}
		j.logLineNum++
	}
}

func (j *job) exclusivityCheck(action string) error {
	// Don't run the job if there's already another running job
	// performing this aciton.
	numJobs, err := j.worker.schemaMgr.CountJobsWithActionStatus(action, schema.JobStatusRunning)
	if err != nil {
		// something's wrong with the db? we can't guarantee anything here, so we shouldn't proceed
		if err := j.worker.schemaMgr.UpdateJobStatus(j.job.ID, schema.JobStatusError); err != nil {
			return fmt.Errorf("unable to update job status to error: %s", err)
		}
		return fmt.Errorf("unable to check for other running jobs with the same action: %s", err)
	}

	// Exclusivity is checked after the job is run so there should be exactly
	// one running
	if numJobs > 1 {
		// There's at least one other job (besides this one)
		// which is performing this exclusive action already.
		// Refuse to run and set the status to "conflict".
		if err := j.worker.schemaMgr.UpdateJobStatus(j.job.ID, schema.JobStatusConflict); err != nil {
			return fmt.Errorf("unable to update job status to conflict: %s", err)
		}
		return fmt.Errorf("unable to perform exclusive action: %s - %d other jobs running", action, numJobs)
	}
	return nil
}

// run starts this job's subprocess and waits for this job to be completed. It
// blocks until the job subprocess has either completed successfully, errored,
// or been canceled. This function should only be called once.
// At the end of run() the job status must be updated to SOMETHING one way or another
func (j *job) run() {
	defer j.Done() // Signals to anyone waiting for this job to complete.
	// Mark the job as complete when run is over
	defer func() {
		j.Lock()
		defer j.Unlock()
		j.isComplete = true
	}()

	if j.actionInfo.Exclusive {
		if err := j.exclusivityCheck(j.job.Action); err != nil {
			context.GetLogger(j.ctx).Errorf("exclusivity check failure: %s", err)
			return
		}
		if j.actionInfo.CleanupJob != nil {
			if err := j.exclusivityCheck(j.actionInfo.CleanupJob.Action); err != nil {
				context.GetLogger(j.ctx).Errorf("exclusivity check failure: %s", err)
				return
			}
		}
	}

	// merge parameters from the action config for this job
	ac, err := j.worker.schemaMgr.SafeGetActionConfig(j.job.Action)
	if err != nil {
		context.GetLogger(j.ctx).Errorf("job failed to get action config for job; executing anyway: %s", err)
	}
	params := map[string]string{}
	if ac.Parameters != nil {
		params = ac.Parameters
	}
	for k, v := range j.job.Parameters {
		params[k] = v
	}

	timeoutDuration, _ := time.ParseDuration(j.job.StopTimeout)
	var reader io.Reader
	// start the job; we need to lock to make sure it doesn't get canceled here
	// This is kind of ugly, but it's used to make sure we unlock in all cases and return if we get an error
	if func() bool {
		j.Lock()
		defer j.Unlock()
		var err error
		var startStatus string
		reader, err = j.executor.Start(j.ctx, j.actionInfo.Command, j.actionInfo.Args, params, timeoutDuration)
		if err != nil {
			ee := err.(*ExecutorError)
			context.GetLogger(j.ctx).Errorf("job failed to start with error: %s", ee.Err)
			j.writeLogs(ee.ExtraLogs...)

			if err := j.worker.schemaMgr.UpdateJobStatus(j.job.ID, startStatus); err != nil {
				context.GetLogger(j.ctx).Errorf("unable to update job status: %s", err)
			}
			return true
		}
		return false
	}() {
		return
	}

	heartbeatTimeout, err := time.ParseDuration(ac.HeartbeatTimeout)
	if err != nil {
		context.GetLogger(j.ctx).Errorf("failed to parse heartbeat timeout: %s", err)
	}
	heartbeatTicker := j.clock.NewTicker(constants.HeartbeatInterval)
	heartbeatStop := make(chan struct{})
	// we stop the heartbeat ticker only after we are truly done with the job
	defer func() {
		heartbeatTicker.Stop()
		close(heartbeatStop)
	}()
	// We heart beat as long as the job is running
	go func() {
		if j.goroutineHack != nil {
			defer j.goroutineHack()
		}
		for {
			select {
			case <-heartbeatStop:
				return
			case <-heartbeatTicker.Chan():
				err := j.worker.schemaMgr.HeartbeatJob(j.job.ID, j.clock.Now().Add(heartbeatTimeout))
				if err != nil {
					context.GetLogger(j.ctx).Errorf("job failed to heart beat job %s: %s", j.job.ID, err)
				}
			}
		}
	}()

	// We need to asynchronously close the logsWriter because cmd doesn't close it for us
	waitRet := make(chan error)
	go func() {
		if j.goroutineHack != nil {
			defer j.goroutineHack()
		}
		err := j.executor.Wait()
		waitRet <- err
	}()

	// schedule timeouts
	var timer clock.Timer
	deadline, _ := time.ParseDuration(j.job.Deadline)
	if deadline > 0 {
		timer = j.clock.AfterFunc(deadline, func() {
			if j.goroutineHack != nil {
				defer j.goroutineHack()
			}
			j.executor.Cancel(schema.JobStatusTimeout)
		})
	}

	scanner := bufio.NewScanner(reader)
	// read out the logs while it runs
	for scanner.Scan() {
		line := scanner.Text()
		j.writeLogs(line)
	}

	// stop the cancelation timer because the job is done
	if timer != nil {
		timer.Stop()
	}

	j.Lock()
	defer j.Unlock()

	context.GetLogger(j.ctx).Info("job process ended")

	err = <-waitRet
	status := schema.JobStatusDone
	if err != nil {
		ee := err.(*ExecutorError)
		status = ee.Status
		// The process was either canceled or errored on its own.
		context.GetLogger(j.ctx).Errorf("job exited with error: %s", ee.Err)

		j.writeLogs(ee.ExtraLogs...)
	}

	if err := j.worker.schemaMgr.UpdateJobStatus(j.job.ID, status); err != nil {
		context.GetLogger(j.ctx).Errorf("unable to update job status: %s", err)
	}

	context.GetLogger(j.ctx).Info("job completed")
}

// cancel ends this job by killing the job's subprocess. The caller should call
// Wait() on this job to wait for the job to shutdown gracefully.
func (j *job) cancel(status string) {
	j.Lock()
	defer j.Unlock()

	j.cancelLocked(status)
}

func (j *job) cancelLocked(status string) {
	if j.isComplete {
		return
	}

	context.GetLogger(j.ctx).Info("job canceled")

	// Killing the process will unblock the run() goroutine.
	j.executor.Cancel(status)
}
