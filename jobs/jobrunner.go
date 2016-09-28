package jobs

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/benbjohnson/clock"
)

// Runner defines a scheduler which may periodically run various Jobs. A
// Runner will only begin running jobs after a call to Start() is successful.
type Runner interface {
	// AddJob adds a Job to the runnable list by name. Returns ErrJobAlreadyExists
	// if the name is already found in the Runner.
	AddJob(name string, job Job) error
	// RemoveJob removes a Job by name from the runnable list. Returns
	// ErrJobNotFound if the name is not found in the Runner.
	RemoveJob(name string) error
	// Start begins running the scheduled Jobs until the next call to Stop.
	// Returns ErrAlreadyStarted if already started.
	Start() error
	// Stop unschedules all Jobs until the next call to Start. Does not cancel
	// Jobs in progress. Returns ErrAlreadyStopped if already stopped.
	Stop() error
	// Status returns a map from Job name to Job status. Status may be one of the
	// following: running, waiting, error
	Status() (map[string]JobStatus, error)
	// RunNow immediately begins running the named Job. Returns
	// ErrJobAlreadyRunning if the Job is already running or ErrJobNotFound if the
	// name is not found in the Runner.
	RunNow(name string) error
}

const (
	StatusRunning = "running"
	StatusWaiting = "waiting"
	StatusError   = "error"
)

var (
	ErrJobAlreadyExists  = errors.New("job already exists")
	ErrJobNotFound       = errors.New("job not found")
	ErrAlreadyStarted    = errors.New("already started")
	ErrAlreadyStopped    = errors.New("already stopped")
	ErrJobAlreadyRunning = errors.New("job already running")
)

// jobInfo is a simple container of both a Job and its status, along with any
// details
type jobInfo struct {
	job    Job
	status JobStatus
}

type JobStatus struct {
	Status  string      `json:"status"`
	Details interface{} `json:"details,omitempty"`
}

type jobRunner struct {
	clock  clock.Clock
	rwLock *sync.RWMutex
	jobs   map[string]jobInfo
	ticker *clock.Ticker
}

func NewRunner(c clock.Clock) Runner {
	return &jobRunner{
		clock:  c,
		rwLock: new(sync.RWMutex),
		jobs:   make(map[string]jobInfo),
	}
}

func (jr *jobRunner) AddJob(name string, job Job) error {
	jr.rwLock.Lock()
	defer jr.rwLock.Unlock()

	if _, ok := jr.jobs[name]; ok {
		return ErrJobAlreadyExists
	}

	jr.jobs[name] = jobInfo{
		job: job,
		status: JobStatus{
			Status: StatusWaiting,
		},
	}

	return nil
}

func (jr *jobRunner) RemoveJob(name string) error {
	jr.rwLock.Lock()
	defer jr.rwLock.Unlock()

	if _, ok := jr.jobs[name]; !ok {
		return ErrJobNotFound
	}

	delete(jr.jobs, name)
	return nil
}

func (jr *jobRunner) Start() error {
	jr.rwLock.Lock()
	defer jr.rwLock.Unlock()

	if jr.ticker != nil {
		return ErrAlreadyStarted
	}

	jr.ticker = jr.clock.Ticker(time.Minute)

	go func() {
		for range jr.ticker.C {
			jr.runReadyJobs()
		}
	}()

	return nil
}

func (jr *jobRunner) Stop() error {
	jr.rwLock.Lock()
	defer jr.rwLock.Unlock()

	if jr.ticker == nil {
		return ErrAlreadyStopped
	}

	jr.ticker.Stop()
	jr.ticker = nil

	return nil
}

func (jr *jobRunner) Status() (map[string]JobStatus, error) {
	jr.rwLock.RLock()
	defer jr.rwLock.RUnlock()

	statusMap := make(map[string]JobStatus)
	for name, ji := range jr.jobs {
		statusMap[name] = ji.status
	}

	return statusMap, nil
}

func (jr *jobRunner) RunNow(name string) error {
	jr.rwLock.Lock()
	defer jr.rwLock.Unlock()

	return jr.runJob(name, true)
}

func (jr *jobRunner) runReadyJobs() {
	jr.rwLock.RLock()
	defer jr.rwLock.RUnlock()

	for name, ji := range jr.jobs {
		if ji.status.Status != StatusRunning && ji.job.IsReady() {
			jr.runJob(name, false)
		}
	}
}

func (jr *jobRunner) runJob(name string, manual bool) error {
	ji, ok := jr.jobs[name]
	if !ok {
		return ErrJobNotFound
	}

	if ji.status.Status == StatusRunning {
		return ErrJobAlreadyRunning
	}

	ji.status = JobStatus{
		Status: StatusRunning,
	}
	jr.jobs[name] = ji

	go func() {
		err := ji.job.Run(manual)

		jr.rwLock.Lock()
		defer jr.rwLock.Unlock()

		// Don't save status if job is removed
		if _, ok := jr.jobs[name]; !ok {
			return
		}

		if err != nil {
			logrus.WithField("error", err).WithField("job", name).Error("Job failed")
			ji.status = JobStatus{
				Status:  StatusError,
				Details: fmt.Sprintf("%v", err),
			}
		} else {
			ji.status = JobStatus{
				Status: StatusWaiting,
			}
		}

		jr.jobs[name] = ji
	}()

	return nil
}
