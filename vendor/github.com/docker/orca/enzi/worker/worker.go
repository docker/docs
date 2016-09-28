package worker

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/docker/distribution/context"
	"github.com/docker/orca/enzi/schema"
	"github.com/satori/go.uuid"
)

// Worker is a simple interface for a simple on-demand and cron job worker
// system.
type Worker interface {
	Run()
	CancelJob(jobID string)
	DeleteJob(jobID string)
	GetJobLogs(jobID string) (io.ReadCloser, error)
	Shutdown()
}

type worker struct {
	ctx       context.Context
	schemaMgr schema.Manager
	id        string
	address   string
	workDir   string

	schedulers          map[string]*cronScheduler
	signalEndSchedulers chan *sync.WaitGroup

	jobsLock      sync.Mutex
	jobs          map[string]*job
	signalEndJobs chan *sync.WaitGroup

	signalEndStateReconciler chan struct{}
}

// New returns a new worker ready to run.
func New(ctx context.Context, schemaMgr schema.Manager, address, workDir string) (Worker, error) {
	w := &worker{
		ctx:                      ctx,
		schemaMgr:                schemaMgr,
		address:                  address,
		workDir:                  workDir,
		schedulers:               make(map[string]*cronScheduler),
		signalEndSchedulers:      make(chan *sync.WaitGroup),
		jobs:                     make(map[string]*job),
		signalEndJobs:            make(chan *sync.WaitGroup),
		signalEndStateReconciler: make(chan struct{}),
	}

	if err := w.init(); err != nil {
		return nil, fmt.Errorf("unable to initialize worker: %s", err)
	}

	return w, nil
}

// init initializes this worker by registering itself with the system.
func (w *worker) init() error {
	// Check if worker ID already exists.
	idFile := filepath.Join(w.workDir, "id")
	idBytes, err := ioutil.ReadFile(idFile)
	switch {
	case err == nil:
		// We must be a worker that died and has been restarted.
		w.id = string(idBytes)
	case os.IsNotExist(err):
		// We must be a new worker and need to create the id file.
		w.id = uuid.NewV4().String()
		if err := ioutil.WriteFile(idFile, []byte(w.id), os.FileMode(0644)); err != nil {
			return fmt.Errorf("unable to write worker id file: %s", err)
		}
	default:
		return fmt.Errorf("unable to read worker id file: %s", err)
	}

	// Register ourselves in the DB.
	if err := w.schemaMgr.CreateOrUpdateWorker(w.id, w.address); err != nil {
		return fmt.Errorf("unable to register worker with the jobs system: %s", err)
	}

	// Create the jobs folder if it does not already exist.
	if err := os.MkdirAll(w.jobsDir(), os.FileMode(0755)); err != nil {
		return fmt.Errorf("unable to create local jobs directory: %s", err)
	}

	if err := w.initializeJobState(); err != nil {
		return fmt.Errorf("unable to initialize job state: %s", err)
	}

	return nil
}

func (w *worker) jobsDir() string {
	return filepath.Join(w.workDir, "jobs")
}

func (w *worker) jobDirPath(jobID string) string {
	// This should be safe as all of our jobIDs should be UUIDs.
	// Splitting after the first 2 characters means there should be no more
	// than 256 unique subdirectories due to having hex-encoded jobIDs.
	return filepath.Join(w.jobsDir(), jobID[:2], jobID[2:])
}

// initializeJobState creates an in-memory record of the jobs which existed
// locally when the worker starts up.
func (w *worker) initializeJobState() error {
	w.jobsLock.Lock()
	defer w.jobsLock.Unlock()

	logger := context.GetLogger(w.ctx)

	localJobIDs, err := w.getLocalJobIDs()
	if err != nil {
		return fmt.Errorf("unable to get local job IDs: %s", err)
	}

	logger.Infof("found %d local jobs", len(localJobIDs))

	for jobID := range localJobIDs {
		jobCtx := context.WithValue(w.ctx, "jobID", jobID)
		jobCtx = context.WithLogger(jobCtx, context.GetLogger(jobCtx, "jobID"))

		// All we know now is what the ID is and that the job must be
		// complete. We do not know what action or what status the job
		// was in. We will update it later with the action and status
		// from the database.
		w.jobs[jobID] = &job{
			ctx:        jobCtx,
			worker:     w,
			id:         jobID,
			isComplete: true,
		}
	}

	// Note: use offset and limit of 0 to get all jobs.
	assignedJobs, err := w.schemaMgr.GetMostRecentlyScheduledJobsForWorker(w.id, 0, 0)
	if err != nil {
		return fmt.Errorf("unable to list assigned jobs: %s", err)
	}

	logger.Infof("found %d jobs assigned in database", len(assignedJobs))

	/*
	   Let A be the set of assigned jobs, according to the database.
	   Let B be the set of local jobs on this worker.

	   A ∩ B is the set of jobs for which there is both a record in the
	   database and a local record on this worker. The status of each of
	   these jobs should match the status recorded in the DB. The job
	   should be errored if they do not match.

	   A - B is the set of jobs assigned to this worker in the DB for
	   which there is no local record. For these jobs, a local record
	   should be created and the status should be set to errored in the
	   database if it was supposed to be running.

	   B - A is the set of jobs with a local record on this worker but
	   which have no record in the DB. These jobs should simply be removed
	   from the local record.
	*/

	// Iterate through set A (jobs assigned to this worker in the DB).
	for _, assignedJob := range assignedJobs {
		jobID := assignedJob.ID

		// Reduce set B to B - A.
		delete(localJobIDs, jobID)

		// Determine whether this job is in the set A ∩ B or A - B.
		localJob, isLocal := w.jobs[jobID]

		if !isLocal {
			logger.Warnf("no local record of job %s", jobID)

			jobCtx := context.WithValue(w.ctx, "jobID", jobID)
			jobCtx = context.WithLogger(jobCtx, context.GetLogger(jobCtx, "jobID"))

			// Create a local record of the job.
			localJob = &job{
				ctx:        jobCtx,
				worker:     w,
				id:         jobID,
				isComplete: true,
			}

			w.jobs[jobID] = localJob
		}

		// Set the action and status of the local job.
		localJob.action = assignedJob.Action
		localJob.status = assignedJob.Status
		localJob.isCanceled = assignedJob.Status == schema.JobStatusCanceled

		// During initialization, no local jobs should be running. If
		// the DB record says the job is running, we should make an
		// attempt to update the job status to "errored".
		if assignedJob.Status == schema.JobStatusRunning {
			localJob.status = schema.JobStatusError

			logger.Warnf("updating stray job %s status to %s", jobID, localJob.status)

			if err := w.schemaMgr.UpdateJobStatus(jobID, localJob.status); err != nil {
				logger.Errorf("unable to update job status: %s", err)
			}
		}
	}

	// localJobIDs contains only those local jobIDs which have no record in
	// the DB. Simply delete these local jobs.
	for jobID := range localJobIDs {
		logger.Warnf("deleting local job %s with no database record", jobID)
		w.deleteJobLocked(jobID)
	}

	return nil
}

func (w *worker) Run() {
	context.GetLogger(w.ctx).Info("running worker")

	go w.monitorCrons()
	go w.claimJobs()
	go w.periodicallyReconcileJobState()

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

	os.Exit(0)
}

func (w *worker) Shutdown() {
	// Wait for all cron schedulers and jobs to end.
	// Note: end the schedulers first so that they don't try to claim any
	// more jobs.
	w.endStateReconciler()
	w.endSchedulers()
	w.endJobs()
}

func (w *worker) endStateReconciler() {
	w.signalEndStateReconciler <- struct{}{}
}

func (w *worker) endSchedulers() {
	var wg sync.WaitGroup

	wg.Add(1)
	w.signalEndSchedulers <- &wg
	wg.Wait()
}

func (w *worker) endJobs() {
	var wg sync.WaitGroup

	wg.Add(1)
	w.signalEndJobs <- &wg
	wg.Wait()
}

// getLocalJobIDs returns a set of jobIDs for the jobs this worker has a local
// record of. Each job is located in 2 subdirs: hex[:2]/hex[2:] We Loop through
// the outer dirs (max 256) and concat the contents of those directories
func (w *worker) getLocalJobIDs() (jobIDs map[string]struct{}, err error) {
	jobIDs = make(map[string]struct{}, 100)

	err = filepath.Walk(w.jobsDir(), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			return nil
		}

		// The given path will contain the jobs dir. We want the next
		// two parts with the jobID prefix and suffix. Strip the prefix
		// and then split on the system's path separator.
		relPath := strings.TrimPrefix(path, w.jobsDir()+string(filepath.Separator))

		parts := strings.Split(relPath, string(filepath.Separator))
		if len(parts) == 2 {
			jobID := parts[0] + parts[1]
			jobIDs[jobID] = struct{}{}

			// Do not descend any further.
			return filepath.SkipDir
		}

		return nil
	})

	return jobIDs, err
}

const (
	maxNumJobs             = 512
	reconcileStateInterval = time.Minute * 5
)

func (w *worker) periodicallyReconcileJobState() {
	for {
		select {
		case <-w.signalEndStateReconciler:
			return
		case <-time.After(reconcileStateInterval):
			w.reconcileJobState()
		}
	}
}

func (w *worker) reconcileJobState() {
	logger := context.GetLogger(w.ctx)
	logger.Info("beginning state reconciliation")
	defer logger.Info("state reconciliation complete")

	w.jobsLock.Lock()
	defer w.jobsLock.Unlock()

	// Note: use offset and limit of 0 to get all jobs.
	assignedJobs, err := w.schemaMgr.GetMostRecentlyScheduledJobsForWorker(w.id, 0, 0)
	if err != nil {
		logger.Errorf("unable to list assinged jobs from database: %s", err)
		return
	}

	logger.Infof("found %d jobs assigned in database", len(assignedJobs))

	if len(assignedJobs) > maxNumJobs {
		logger.Info("to many assigned jobs")

		// We've got some cleanup to do. Sort the assigned jobs ordered
		// by most recent 'scheduledAt' time. We only want to keep the
		// latest maxNumJobs.
		sort.Sort(mostRecentlyScheduled(assignedJobs))

		toDelete := assignedJobs[maxNumJobs:]

		for _, job := range toDelete {
			logger.Infof("deleting job %s", job.ID)

			// Delete any local state.
			w.deleteJobLocked(job.ID)

			// Delete from DB.
			if err := w.schemaMgr.DeleteJob(job.ID); err != nil {
				logger.Errorf("unable to delete job %s from database: %s", job.ID, err)
			}
		}
	}

	// Create a set of local job IDs.
	localJobIDs := make(map[string]struct{}, len(w.jobs))
	for jobID := range w.jobs {
		localJobIDs[jobID] = struct{}{}
	}

	logger.Infof("have %d local jobs", len(localJobIDs))

	/*
	   Let A be the set of assigned jobs, according to the database.
	   Let B be the set of local jobs on this worker.

	   A ∩ B is the set of jobs for which there is both a record in the
	   database and a local record on this worker. The status of each of
	   these jobs should match the status recorded in the DB. The job
	   should be errored if they do not match.

	   A - B is the set of jobs assigned to this worker in the DB for
	   which there is no local record. For these jobs, a local record
	   should be created and the status should be set to errored in the
	   database if it was supposed to be running.

	   B - A is the set of jobs with a local record on this worker but
	   which have no record in the DB. These jobs should simply be removed
	   from the local record.
	*/

	for _, assignedJob := range assignedJobs {
		jobID := assignedJob.ID

		// Reduce set B to B - A.
		delete(localJobIDs, jobID)

		// Determine whether this job is in the set A ∩ B or A - B.
		localJob, isLocal := w.jobs[jobID]

		if !isLocal {
			logger.Warnf("no local record of job %s", jobID)

			jobCtx := context.WithValue(w.ctx, "jobID", jobID)
			jobCtx = context.WithLogger(jobCtx, context.GetLogger(jobCtx, "jobID"))

			// Create a local record of the job.
			localJob = &job{
				ctx:        jobCtx,
				worker:     w,
				id:         jobID,
				isComplete: true,
				action:     assignedJob.Action,
				status:     assignedJob.Status,
				isCanceled: assignedJob.Status == schema.JobStatusCanceled,
			}

			if err := os.MkdirAll(w.jobDirPath(jobID), os.FileMode(0755)); err != nil {
				logger.Errorf("unable to create job dir: %s", err)
			} else if err := ioutil.WriteFile(
				localJob.logFilename(),
				[]byte("record of this job's logs was lost\n"),
				os.FileMode(0644),
			); err != nil {
				logger.Errorf("unable to write job log file: %s", err)
			}

			w.jobs[jobID] = localJob

			// If the status is running, mark the job as errored
			// because it's clearly not running locally.
			if localJob.status == schema.JobStatusRunning {
				localJob.status = schema.JobStatusError
			}
		}

		localJob.reconcileState(assignedJob)
	}

	// localJobIDs contains only those local jobIDs which have no record in
	// the DB. Simply delete these local jobs.
	for jobID := range localJobIDs {
		logger.Warnf("deleting local job %s with no database record", jobID)
		w.deleteJobLocked(jobID)
	}
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
