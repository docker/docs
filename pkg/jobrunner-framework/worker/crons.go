package worker

import (
	"sync"
	"time"

	"github.com/WatchBeam/clock"
	"github.com/docker/dhe-deploy/pkg/jobrunner-framework/schema"
	"github.com/docker/distribution/context"
	"github.com/robfig/cron"
)

type cronScheduler struct {
	ctx          context.Context
	worker       *worker
	cron         *schema.Cron
	schedule     cron.Schedule
	signalCancel chan struct{}

	sync.Mutex     // Used to protect the isCanceled field.
	sync.WaitGroup // Used to signal that the scheduler is done.
	isCanceled     bool
}

func (s *cronScheduler) run() {
	// Signal to any goroutines waiting for the scheduler to shutdown.
	defer func() {
		s.Done()
	}()

	for {
		now := s.worker.clock.Now()
		scheduledAt := s.schedule.Next(now)

		var (
			timer clock.Timer
			wait  <-chan time.Time
		)

		if scheduledAt.IsZero() {
			// This means that the schedule could not find a valid
			// time in the next 5 years, so just block forever in
			// the hope that an administrator fixes it.
			// Note: receive from a nil channel blocks forever.
			timer = nil
			wait = nil
		} else {
			timer = s.worker.clock.NewTimer(scheduledAt.Sub(now))
			wait = timer.Chan()
		}

		select {
		case <-wait:
			context.GetLogger(s.ctx).Infof("timer fired to run job %s scheduled at %s", s.cron.Action, scheduledAt)

			s.claimJob(scheduledAt)

		case <-s.signalCancel:
			context.GetLogger(s.ctx).Info("cron canceled")

			// Stop the timer if it has been initialized.
			if timer != nil {
				timer.Stop()
			}

			return
		}
	}
}

func (s *cronScheduler) claimJob(scheduledAt time.Time) {
	s.worker.jobsLock.Lock()
	defer s.worker.jobsLock.Unlock()

	job := &schema.Job{
		WorkerID:    s.worker.id,
		CronID:      s.cron.ID,
		Status:      schema.JobStatusRunning,
		ScheduledAt: scheduledAt,
		LastUpdated: scheduledAt,
		Action:      s.cron.Action,
		Parameters:  s.cron.Parameters,
		Deadline:    s.cron.Deadline,
		StopTimeout: s.cron.StopTimeout,
	}

	if err := s.worker.schemaMgr.CreateJob(job); err != nil {
		if err == schema.ErrJobExists {
			context.GetLogger(s.ctx).Info("job already claimed by another worker")
		} else {
			context.GetLogger(s.ctx).Errorf("unable to create job: %s", err)
		}

		return
	}

	subCtx := context.WithValue(s.ctx, "jobID", job.ID)
	subCtx = context.WithLogger(subCtx, context.GetLogger(subCtx, "jobID"))

	s.worker.addJob(subCtx, job)
}

// cancel ends this scheduler by signaling it to end. The caller should call
// Wait() on this scheduler to wait for it to shutdown gracefully.
func (s *cronScheduler) cancel() {
	s.Lock()
	defer s.Unlock()

	if s.isCanceled {
		return
	}

	s.isCanceled = true
	s.signalCancel <- struct{}{}
}
