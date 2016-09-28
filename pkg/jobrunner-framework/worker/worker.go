package worker

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/WatchBeam/clock"
	"github.com/docker/dhe-deploy/pkg/jobrunner-framework/schema"
	"github.com/docker/distribution/context"
	"github.com/robfig/cron"
)

// Worker is a simple interface for a simple on-demand and cron job worker
// system.
type Worker interface {
	// Run launches the lower-level methods: Init(), MonitorGoroutine(), ClaimJobsGoroutine(), SignalHandlerGoroutine() and finally does a Wait() for the worker to stop
	Run() error
	Init() error
	MonitorCronsGoroutine() error
	MonitorCancellationsGoroutine() error
	MonitorRecoverableJobsGoroutine() error
	ClaimJobsGoroutine() error
	SignalHandlerGoroutine() error
	DetectExpirationGoroutine() error
	Wait() error

	WaitTimeout(timeout time.Duration) error

	CancelJob(jobID string)
	DeleteJob(jobID string, status string)
	Shutdown()
}

type worker struct {
	ctx        context.Context
	schemaMgr  schema.JobrunnerManager
	id         string
	actionInfo map[string]ActionInfo
	// A function to be deferred in new goroutines
	// used to allow catching of exceptions during testing
	goroutineHack            func()
	expirationTicker         clock.Ticker
	expirationTickerChan     <-chan time.Time
	expirationTickerStopChan chan struct{}
	wgEndExpirationTicker    sync.WaitGroup
	expirationHeap           *SafeTimeHeap

	schedulers               map[string]*cronScheduler
	signalEndSchedulers      chan *sync.WaitGroup
	signalEndJobs            chan *sync.WaitGroup
	signalEndCancelJobs      chan *sync.WaitGroup
	signalEndRecoverableJobs chan *sync.WaitGroup
	wgEndSchedulers          sync.WaitGroup
	wgEndJobs                sync.WaitGroup
	wgEndCancelJobs          sync.WaitGroup
	wgEndRecoverableJobs     sync.WaitGroup
	clock                    clock.Clock

	jobsLock sync.Mutex
	jobs     map[string]*job

	newExecutor NewExecutor
}

// New returns a new worker ready to run.
func New(ctx context.Context, schemaMgr schema.JobrunnerManager, id string, actionInfo map[string]ActionInfo, c clock.Clock, goroutineHack func(), ne NewExecutor) Worker {
	if ne == nil {
		ne = DefaultNewExecutor
	}
	w := &worker{
		id:                       id,
		actionInfo:               actionInfo,
		ctx:                      ctx,
		schemaMgr:                schemaMgr,
		schedulers:               make(map[string]*cronScheduler),
		jobs:                     make(map[string]*job),
		signalEndSchedulers:      make(chan *sync.WaitGroup, 1),
		signalEndJobs:            make(chan *sync.WaitGroup, 1),
		signalEndCancelJobs:      make(chan *sync.WaitGroup, 1),
		signalEndRecoverableJobs: make(chan *sync.WaitGroup, 1),
		clock: c,
		expirationTickerStopChan: make(chan struct{}),
		goroutineHack:            goroutineHack,
		newExecutor:              ne,
		expirationHeap:           NewSafeTimeHeap(),
	}
	w.expirationTicker = w.clock.NewTicker(time.Second)
	w.expirationTickerChan = w.expirationTicker.Chan()
	w.wgEndExpirationTicker.Add(1)
	w.wgEndJobs.Add(1)
	w.wgEndCancelJobs.Add(1)
	w.wgEndRecoverableJobs.Add(1)
	w.wgEndSchedulers.Add(1)

	return w
}

// init initializes this worker by registering itself with the system.
func (w *worker) Init() error {
	if err := w.initializeJobState(); err != nil {
		return fmt.Errorf("unable to initialize job state: %s", err)
	}

	return nil
}

// initializeJobState creates an in-memory record of the jobs which existed
// locally when the worker starts up.
func (w *worker) initializeJobState() error {
	w.jobsLock.Lock()
	defer w.jobsLock.Unlock()

	logger := context.GetLogger(w.ctx)

	// Note: use offset and limit of 0 to get all jobs.
	assignedJobs, err := w.schemaMgr.GetMostRecentlyScheduledJobsForWorker(w.id, 0, 0)
	if err != nil {
		return fmt.Errorf("unable to list assigned jobs: %s", err)
	}

	logger.Infof("found %d jobs assigned in database", len(assignedJobs))

	// Iterate through the jobs and set them to failed because we just came back up from the dead
	// TODO: research if it's possible for a job's binary to run away and keep
	// running even though the worker has restarted: probably not unless it double forks?
	for _, assignedJob := range assignedJobs {
		if assignedJob.Status == schema.JobStatusRunning {
			jobID := assignedJob.ID

			// TODO: maybe record a special log line explaining the worker died?
			if err := w.schemaMgr.UpdateJobStatus(jobID, schema.JobStatusWorkerDead); err != nil {
				logger.Errorf("unable to update job status: %s", err)
			}
		}
	}

	return nil
}

func (w *worker) MonitorCronsGoroutine() error {
	// TODO: error reporting
	context.GetLogger(w.ctx).Info("monitoring crons")

	for w.monitorCronsLoop() {
		// Lost connection to the cron changes stream.
		// Wait 5 seconds before trying again.
		w.clock.Sleep(5 * time.Second)
	}
	return nil
}

func (w *worker) MonitorRecoverableJobsGoroutine() error {
	context.GetLogger(w.ctx).Info("monitoring recoverable jobs")

	// decide which jobs are recoverable so we can construct this query optimally
	recoverableActions := []string{}
	for action, info := range w.actionInfo {
		if info.CleanupJob != nil {
			recoverableActions = append(recoverableActions, action)
		}
	}

	// TODO: error reporting
	for w.monitorRecoverableJobsLoop(recoverableActions) {
		// Lost connection to the unclaimed jobs changes stream.
		// Wait 5 seconds before trying again.
		time.Sleep(5 * time.Second)
	}
	return nil
}

func (w *worker) MonitorCancellationsGoroutine() error {
	context.GetLogger(w.ctx).Info("monitoring cancellations")

	// TODO: error reporting
	for w.monitorCancellationsLoop() {
		// Lost connection to the unclaimed jobs changes stream.
		// Wait 5 seconds before trying again.
		time.Sleep(5 * time.Second)
	}
	return nil
}

func (w *worker) ClaimJobsGoroutine() error {
	context.GetLogger(w.ctx).Info("monitoring new jobs")

	// TODO: error reporting
	for w.claimJobsLoop() {
		// Lost connection to the unclaimed jobs changes stream.
		// Wait 5 seconds before trying again.
		time.Sleep(5 * time.Second)
	}
	return nil
}

func (w *worker) SignalHandlerGoroutine() error {
	// Use the current goroutine to listen for SIGINT or SIGTERM.
	// The worker will cancel all cron schedulers and running jobs in order
	// to shutdown gracefully.
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-signals:
		context.GetLogger(w.ctx).Infof("responding to %s", sig)
	}

	w.Shutdown()
	return nil
}

// Wait for all cron schedulers and jobs to end.
func (w *worker) Wait() error {
	// TODO: report any errors that happened during shutdown
	w.wgEndSchedulers.Wait()
	w.wgEndJobs.Wait()
	w.wgEndCancelJobs.Wait()
	w.wgEndRecoverableJobs.Wait()
	return nil
}

// WaitTimeout is useful in case the goroutines panicked, such as in tests
func (w *worker) WaitTimeout(timeout time.Duration) error {
	done := make(chan struct{})
	go func() {
		// TODO: report any errors that happened during shutdown
		if w.goroutineHack != nil {
			defer w.goroutineHack()
		}
		w.wgEndSchedulers.Wait()
		w.wgEndJobs.Wait()
		w.wgEndCancelJobs.Wait()
		w.wgEndRecoverableJobs.Wait()
		close(done)
	}()

	timedOut := time.After(timeout)
	select {
	case <-done:
		return nil
	case <-timedOut:
		return fmt.Errorf("timed out waiting for goroutines to end after %s", timeout)
	}
	panic("impossible")
}

func (w *worker) Run() error {
	if err := w.Init(); err != nil {
		return fmt.Errorf("unable to initialize worker: %s", err)
	}

	context.GetLogger(w.ctx).Info("running worker")

	go func() {
		err := w.MonitorCronsGoroutine()
		if err != nil {
			context.GetLogger(w.ctx).Error("error from monitor crons goroutine: %s", err)
		}
	}()

	go func() {
		err := w.MonitorCancellationsGoroutine()
		if err != nil {
			context.GetLogger(w.ctx).Error("error from monitor cancellations goroutine: %s", err)
		}
	}()

	go func() {
		err := w.ClaimJobsGoroutine()
		if err != nil {
			context.GetLogger(w.ctx).Error("error from claim jobs goroutine: %s", err)
		}
	}()

	go func() {
		err := w.DetectExpirationGoroutine()
		if err != nil {
			context.GetLogger(w.ctx).Error("error from detect expiration goroutine: %s", err)
		}
	}()

	go func() {
		err := w.SignalHandlerGoroutine()
		if err != nil {
			context.GetLogger(w.ctx).Error("error from signal handler goroutine: %s", err)
		}
	}()

	return w.Wait()
}

func (w *worker) Shutdown() {
	// Note: end the schedulers first so that they don't try to claim any
	// more jobs.

	w.endSchedulers()
	w.wgEndSchedulers.Wait()
	w.endJobs()
}

func (w *worker) endSchedulers() {
	w.signalEndSchedulers <- &w.wgEndSchedulers
}

func (w *worker) endJobs() {
	w.expirationTicker.Stop()
	close(w.expirationTickerStopChan)
	w.wgEndExpirationTicker.Wait()
	w.signalEndJobs <- &w.wgEndJobs
	w.signalEndCancelJobs <- &w.wgEndCancelJobs
	w.signalEndRecoverableJobs <- &w.wgEndRecoverableJobs
}

// monitorCronsLoop consumes a stream of changes to all registered cron jobs.
// If the stream ends or there is a connection issue it returns true to signal
// the caller to wait and try again. Returns false after receiving a shutdown
// signal.
func (w *worker) monitorCronsLoop() bool {
	changeStream, streamCloser, err := w.schemaMgr.GetCronChanges()
	if err != nil {
		context.GetLogger(w.ctx).Errorf("unable to get stream of cron changes: %s", err)
		return true // Try again.
	}

	defer streamCloser.Close()

	for {
		select {
		case wg := <-w.signalEndSchedulers:
			defer wg.Done()

			// The main goroutine has received a signal to end the
			// process. Cancel all schedulers and use the given
			// waitgroup to signal that cleanup is complete. The
			// sender will still need to wait for each scheduler to
			// gracefully shutdown.
			w.cancelAllSchedulers()

			return false // Don't try again.

		case change, ok := <-changeStream:
			if !ok {
				context.GetLogger(w.ctx).Info("cron change stream unexpectedly closed")
				// Cancel all schedulers and try again.
				w.cancelAllSchedulers()
				return true
			}

			context.GetLogger(w.ctx).Debugf("old cron value: %#v", change.OldValue)
			context.GetLogger(w.ctx).Debugf("new cron value: %#v", change.NewValue)

			if change.OldValue != nil {
				// Cancel the scheduler with this cronID.
				w.cancelScheduler(change.OldValue.ID)
			}

			if change.NewValue != nil {
				subCtx := context.WithValue(w.ctx, "cronID", change.NewValue.ID)
				subCtx = context.WithLogger(subCtx, context.GetLogger(subCtx, "cronID"))

				w.addScheduler(subCtx, change.NewValue)
			}
		}
	}
}

// addScheduler creates and starts a new scheduler with the given cronID using
// the given cronSpec schedule to create jobs to run the given action. The
// scheduler is added to this worker's set of currently running scheduler. If
// there is an error parsing the given cronSpec then an error is logged and the
// scheduler does not run.
func (w *worker) addScheduler(ctx context.Context, c *schema.Cron) {
	schedule, err := cron.Parse(c.Schedule)
	if err != nil {
		// How did this get through validation?
		context.GetLogger(ctx).Errorf("invalid cron spec %q: %s", c.Schedule, err)
		return
	}

	scheduler := &cronScheduler{
		ctx:      ctx,
		worker:   w,
		cron:     c,
		schedule: schedule,
		// signalCancel has size of 1 so that when we signal a cancellation, we don't have to wait for it to be received
		// this is useful in case the goroutine has already crashed
		signalCancel: make(chan struct{}, 1),
	}

	w.schedulers[c.ID] = scheduler

	scheduler.Add(1)

	go func() {
		if w.goroutineHack != nil {
			defer w.goroutineHack()
		}
		scheduler.run()
	}()
}

// cancelScheduler signals the scheduler with the given cronID to cancel itself
// and waits for it to gracefully shutdown. The scheduler is then removed from
// this worker's set of schedulers.
func (w *worker) cancelScheduler(cronID string) {
	scheduler, ok := w.schedulers[cronID]
	if !ok {
		return // No such scheduler; may have already been canceled.
	}

	scheduler.cancel()
	scheduler.Wait()

	delete(w.schedulers, cronID)
}

// cancelAllSchedulers signals each currently running scheduler to cancel
// itself. Once each scheduler has gracefully shutdown, the worker's set of
// schedulers is reset to be empty.
func (w *worker) cancelAllSchedulers() {
	context.GetLogger(w.ctx).Info("canceling all cron schedulers")

	for _, scheduler := range w.schedulers {
		scheduler.cancel()
	}

	for _, scheduler := range w.schedulers {
		scheduler.Wait()
	}

	w.schedulers = make(map[string]*cronScheduler)
}

func (w *worker) DetectExpirationGoroutine() error {
	// we want to notify the caller when we've existed
	defer w.wgEndExpirationTicker.Done()
	// run every second, looking for expired timers
For:
	for {
		select {
		case <-w.expirationTickerStopChan:
			break For
		case <-w.expirationTickerChan:
			// note that this inner loop is just to make sure we process all events in one check
			// the outside loop is the one that does the periodic checking loop
			for {
				job := w.expirationHeap.PopIfExpired(w.clock.Now())
				// if an expired job was found, mark its worker as dead
				if job != nil {
					// TODO: consider races with others updating the status at the same time?
					err := w.schemaMgr.UpdateJobStatus(job.ID, schema.JobStatusWorkerDead)
					if err != nil {
						context.GetLogger(w.ctx).Errorf("error updating job status to mark worker as dead: %s", err)
					}
				} else {
					break
				}
			}
		}
	}
	return nil
}

func (w *worker) monitorCancellationsLoop() bool {
	// get the feed of job changes we care about
	changeStream, streamCloser, err := w.schemaMgr.GetOwnJobCancellations(w.id)
	if err != nil {
		context.GetLogger(w.ctx).Errorf("unable to get stream of unclaimed job changes: %s", err)
		return true // Try again.
	}

	defer streamCloser.Close()

	for {
		select {
		case wg := <-w.signalEndCancelJobs:
			defer wg.Done()
			return false // Don't try again.

		case change, ok := <-changeStream:
			if !ok {
				context.GetLogger(w.ctx).Info("cancellation change stream unexpectedly closed")
				return true // Try again.
			}

			context.GetLogger(w.ctx).Debugf("old job value: %#v", change.OldValue)
			context.GetLogger(w.ctx).Debugf("new job value: %#v", change.NewValue)

			// Find the job and cancel it. If it's not ours or already canceled, do nothing.
			if change.NewValue != nil && change.NewValue.Status == schema.JobStatusCancelRequest && change.NewValue.WorkerID == w.id {
				w.DeleteJob(change.NewValue.ID, schema.JobStatusCanceled)
			}
			// Detect job deletions and cancel the job if it's deleted
			if change.OldValue != nil && change.NewValue == nil {
				w.DeleteJob(change.OldValue.ID, schema.JobStatusDeleted)
			}
		}
	}
}

func (w *worker) monitorRecoverableJobsLoop(recoverableActions []string) bool {
	// get the feed of job changes we care about
	changeStream, streamCloser, err := w.schemaMgr.GetRecoverableJobChanges(recoverableActions)
	if err != nil {
		context.GetLogger(w.ctx).Errorf("unable to get stream of unclaimed job changes: %s", err)
		return true // Try again.
	}

	defer streamCloser.Close()

	for {
	Select:
		select {
		case wg := <-w.signalEndRecoverableJobs:
			defer wg.Done()

			return false // Don't try again.

		case change, ok := <-changeStream:
			if !ok {
				context.GetLogger(w.ctx).Info("unclaimed job change stream unexpectedly closed")
				return true // Try again.
			}

			context.GetLogger(w.ctx).Debugf("old job value: %#v", change.OldValue)
			context.GetLogger(w.ctx).Debugf("new job value: %#v", change.NewValue)

			// update our knowledge of the state of everyone's recoverable jobs
			for _, ra := range recoverableActions {
				if change.OldValue != nil && change.OldValue.Action == ra || change.NewValue != nil && change.NewValue.Action == ra {
					if change.OldValue != nil && change.OldValue.WorkerID == w.id || change.NewValue != nil && change.NewValue.WorkerID == w.id {
						// watch out for workers proclaiming us dead and cancel the jobs proclaimed dead
						if change.NewValue != nil && change.NewValue.Status == schema.JobStatusWorkerDead {
							w.DeleteJob(change.NewValue.ID, schema.JobStatusWorkerResurrection)
							// we are done processing this event. we shouldn't react to this even further by running the recovery job
							break Select
						}
					} else {
						// we don't track our own jobs for liveness because that's silly
						w.expirationHeap.Update(&change)
					}
					break
				}
			}

			// This job was failed by a worker, we should try to create its recovery job if there is one (whether we or someone else failed it)
			if change.NewValue != nil && schema.StatusIsFinished(change.NewValue.Status) && change.OldValue != nil && !schema.StatusIsFinished(change.OldValue.Status) {
				ai := w.actionInfo[change.NewValue.Action]
				if ai.CleanupJob != nil {
					timeNow := w.clock.Now()
					job := &schema.Job{
						WorkerID:       w.id,
						Status:         schema.JobStatusRunning,
						RecoveryFromID: change.NewValue.ID,
						ScheduledAt:    timeNow,
						LastUpdated:    timeNow,
						Action:         ai.CleanupJob.Action,
						Parameters:     ai.CleanupJob.Parameters,
						Deadline:       ai.CleanupJob.Deadline,
						StopTimeout:    ai.CleanupJob.StopTimeout,
					}

					if err := w.schemaMgr.CreateJob(job); err != nil {
						if err == schema.ErrJobExists {
							context.GetLogger(w.ctx).Info("cleanup job already claimed by another worker")
						} else {
							context.GetLogger(w.ctx).Errorf("unable to create cleanup job: %s", err)
						}
					} else {
						// XXX: maybe it's not ideal that we are relying on having the ID modified in an argument rather than having
						// a return value. We should change CreateJob to return the job it created
						subCtx := context.WithValue(w.ctx, "jobID", job.ID)
						subCtx = context.WithLogger(subCtx, context.GetLogger(subCtx, "jobID"))
						func() {
							w.jobsLock.Lock()
							defer w.jobsLock.Unlock()
							w.addJob(subCtx, job)
						}()
					}
				}
			}
		}
	}
}

func (w *worker) claimJobsLoop() bool {
	// get the feed of job changes we care about
	changeStream, streamCloser, err := w.schemaMgr.GetNewJobChanges()
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
			w.cancelAllJobs(schema.JobStatusWorkerShutdown)

			return false // Don't try again.

		case change, ok := <-changeStream:
			if !ok {
				context.GetLogger(w.ctx).Info("unclaimed job change stream unexpectedly closed")
				return true // Try again.
			}

			context.GetLogger(w.ctx).Debugf("old job value: %#v", change.OldValue)
			context.GetLogger(w.ctx).Debugf("new job value: %#v", change.NewValue)

			// new unassigned job
			if change.OldValue == nil {
				newJob := change.NewValue
				if newJob.WorkerID == "" {
					// Check the parallelism level for this job before trying to claim it
					// if we are overloaded, ignore the new job. It will be picked up when we
					// or someone else finishes the next job

					// merge parameters from the action config for this job
					ac, err := w.schemaMgr.SafeGetActionConfig(newJob.Action)
					if err != nil {
						context.GetLogger(w.ctx).Errorf("job failed to get action config for job, using a default config: %s", err)
					}
					// count how many of these jobs we are currently executing
					count := 0
					w.jobsLock.Lock()
					for _, j := range w.jobs {
						if j.job.Action == newJob.Action {
							count += 1
						}
					}
					w.jobsLock.Unlock()
					if ac != nil && ac.MaxJobsPerWorker > 0 && count > ac.MaxJobsPerWorker {
						context.GetLogger(w.ctx).Infof("queuing job %s", newJob.ID)
						// we have too many of these jobs, ignore it, we'll try to pick it up later
						// TODO: make sure there's no race where jobs are dropped on the floor
					} else {
						context.GetLogger(w.ctx).Infof("claiming job %s", newJob.ID)
						w.claimJob(newJob)
					}
				}
			}
		}
	}
}

func (w *worker) claimJob(job *schema.Job) {
	w.jobsLock.Lock()
	defer w.jobsLock.Unlock()
	ac, err := w.schemaMgr.SafeGetActionConfig(job.Action)
	if err != nil {
		context.GetLogger(w.ctx).Errorf("job failed to get action config for job, using a default config: %s", err)
	}

	duration, err := time.ParseDuration(ac.HeartbeatTimeout)
	if err != nil {
		context.GetLogger(w.ctx).Errorf("failed to parse heartbeat timeout: %s", err)
	}
	// Attempt to claim this job.
	claimedJob, err := w.schemaMgr.ClaimJob(job.ID, w.id, w.clock.Now().Add(duration))
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

	w.addJob(subCtx, claimedJob)
}

// DeleteJob attempts to delete the job with the given ID.
func (w *worker) DeleteJob(jobID string, status string) {
	w.jobsLock.Lock()
	defer w.jobsLock.Unlock()

	w.deleteJobLocked(jobID, status)
}

func (w *worker) deleteJobLocked(jobID string, status string) {
	job, ok := w.jobs[jobID]
	if !ok {
		return // No such job; must have already been deleted.
	}

	job.cancel(status)
	job.Wait()

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

	job.cancel(schema.JobStatusCanceled)
	job.Wait()
}

// cancelAllJobs signals each currently running job to cancel itself and waits
// for all running jobs to gracefully shutdown.
func (w *worker) cancelAllJobs(status string) {
	context.GetLogger(w.ctx).Info("canceling all running jobs")

	w.jobsLock.Lock()
	defer w.jobsLock.Unlock()

	for _, job := range w.jobs {
		job.cancel(status)
	}

	for _, job := range w.jobs {
		job.Wait()
	}
}

// addJob creates and starts a new job with the given jobID to run the given
// action. The job is added to this worker's set of currently running jobs. If
// there is an error initializing the job, the job status is updated to
// the errored status.
func (w *worker) addJob(ctx context.Context, dbjob *schema.Job) {
	// XXX: what do we do if the action is not defined? put it back in the queue? fail it?
	job := &job{
		ctx:           ctx,
		worker:        w,
		job:           dbjob,
		actionInfo:    w.actionInfo[dbjob.Action],
		newExecutor:   w.newExecutor,
		goroutineHack: w.goroutineHack,
		clock:         w.clock,
	}

	if err := job.init(); err != nil {
		logger := context.GetLogger(ctx)
		logger.Errorf("failed to initialize job: %s", err)

		if err := w.schemaMgr.UpdateJobStatus(job.job.ID, schema.JobStatusError); err != nil {
			logger.Errorf("unable to update job status: %s", err)
		}

		return
	}

	w.jobs[job.job.ID] = job

	job.Add(1)

	go func() {
		if w.goroutineHack != nil {
			defer w.goroutineHack()
		}
		// after a job completes, try to schedule another one if we have one in the queue
		defer func() {
			// merge parameters from the action config for this job
			ac, err := w.schemaMgr.SafeGetActionConfig(job.job.Action)
			if err != nil {
				context.GetLogger(w.ctx).Errorf("job failed to get action config for job: %s", err)
				return
			}
			if ac.MaxJobsPerWorker > 0 {
				queuedJob, err := w.schemaMgr.GetNextUnclaimedJob(job.job.Action)
				if err != nil {
					context.GetLogger(ctx).Errorf("unable to get next unclaimed job for action %s: %s", job.job.Action, err)
					return
				}
				if queuedJob != nil {
					w.claimJob(queuedJob)
				}
			}
		}()
		job.run()
	}()
}

type mostRecentlyScheduled []schema.Job

func (jobs mostRecentlyScheduled) Len() int {
	return len(jobs)
}

func (jobs mostRecentlyScheduled) Less(i, j int) bool {
	return jobs[i].ScheduledAt.After(jobs[j].ScheduledAt)
}

func (jobs mostRecentlyScheduled) Swap(i, j int) {
	jobs[i], jobs[j] = jobs[j], jobs[i]
}
