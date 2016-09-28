package worker

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/docker/distribution/context"
	ldapconfig "github.com/docker/orca/enzi/authn/ldap/config"
	"github.com/docker/orca/enzi/config"
	"github.com/docker/orca/enzi/schema"
	"github.com/robfig/cron"
)

const (
	cronIDCleanupDB = "cleanup-db"
	cronIDLDAPSync  = "ldap-sync"
)

var cronAffectingProperties = []string{
	config.AuthConfigPropertyKey,
	ldapconfig.LDAPConfigPropertyKey,
}

func (w *worker) scheduleCleanupDB() {
	cronID := cronIDCleanupDB
	action := ActionCleanupDB
	// Run cleanup on the 42nd minute of every hour.
	cronSpec := "0 42 * * *"

	subCtx := context.WithValue(w.ctx, "cronID", cronID)
	subCtx = context.WithLogger(subCtx, context.GetLogger(subCtx, "cronID"))

	w.addScheduler(subCtx, cronID, action, cronSpec)
}

func (w *worker) monitorCrons() {
	context.GetLogger(w.ctx).Info("monitoring crons")

	for w.monitorCronsLoop() {
		// Lost connection to the cron changes stream.
		// Wait 5 seconds before trying again.
		time.Sleep(5 * time.Second)
	}
}

// monitorCronsLoop consumes a stream of changes to configuration properties
// which may affect cron job schedules. If the stream ends or there is a
// connection issue it returns true to signal the caller to wait and try again.
// Returns false after receiving a shutdown signal.
func (w *worker) monitorCronsLoop() bool {
	changeStream, streamCloser, err := w.schemaMgr.GetPropertyChanges(cronAffectingProperties...)
	if err != nil {
		context.GetLogger(w.ctx).Errorf("unable to get stream of property changes: %s", err)
		return true // Try again.
	}

	defer streamCloser.Close()

	// Always start the cleanup-db cron scheduler.
	w.scheduleCleanupDB()

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
				context.GetLogger(w.ctx).Info("property change stream unexpectedly closed")
				// Cancel all schedulers and try again.
				w.cancelAllSchedulers()
				return true
			}

			context.GetLogger(w.ctx).Debugf("old cron value: %#v", change.OldValue)
			context.GetLogger(w.ctx).Debugf("new cron value: %#v", change.NewValue)

			w.handlePropertyChange(change)
		}
	}
}

func (w *worker) handlePropertyChange(change schema.PropertyChange) {
	// Determine the property key.
	var changedPropertyKey string
	if change.OldValue != nil {
		changedPropertyKey = change.OldValue.Key
	} else if change.NewValue != nil {
		changedPropertyKey = change.NewValue.Key
	} else {
		context.GetLogger(w.ctx).Warn("got a property change with both old and new values nil")
		return // Nothing to do.
	}

	// Check which property is being changed and perform the appropriate
	// action.
	switch changedPropertyKey {
	case config.AuthConfigPropertyKey:
		w.handleChangedAuthBackend(change)
	case ldapconfig.LDAPConfigPropertyKey:
		w.handleChangedLDAPSettings(change)
	default:
		context.GetLogger(w.ctx).Warnf("not sure how to handle changed property key: %q", changedPropertyKey)
	}
}

func (w *worker) handleChangedAuthBackend(change schema.PropertyChange) {
	// Always cancel our current scheduler.
	w.cancelScheduler(cronIDLDAPSync)

	if change.NewValue == nil {
		// The config was deleted? That's weird, but the LDAP sync cron
		// has been canceled anyway.
		return
	}

	var authConfig config.Auth
	if err := json.Unmarshal([]byte(change.NewValue.JSONValue), &authConfig); err != nil {
		context.GetLogger(w.ctx).Warnf("unable to decode new auth config property value: %s", err)
		// We can't understand the new value but the LDAP sync cron has
		// been canceled anyway.
		return
	}

	if authConfig.Backend != config.AuthBackendLDAP {
		// We're not using the LDAP backend in the new config so don't
		// bother starting the cron job.
		return
	}

	// Get the current LDAP settings.
	ldapSettings, err := ldapconfig.GetLDAPConfig(w.schemaMgr)
	if err != nil {
		context.GetLogger(w.ctx).Errorf("unable to get system LDAP config: %s", err)
		return
	}

	cronID := cronIDLDAPSync
	action := ActionLdapSync
	cronSpec := ldapSettings.SyncSchedule

	subCtx := context.WithValue(w.ctx, "cronID", cronID)
	subCtx = context.WithLogger(subCtx, context.GetLogger(subCtx, "cronID"))

	w.addScheduler(subCtx, cronID, action, cronSpec)
}

func (w *worker) handleChangedLDAPSettings(change schema.PropertyChange) {
	// Always cancel our current scheduler.
	w.cancelScheduler(cronIDLDAPSync)

	if change.NewValue == nil {
		// The config was deleted? That's weird, but the LDAP sync cron
		// has been canceled anyway.
		return
	}

	var ldapSettings ldapconfig.Settings
	if err := json.Unmarshal([]byte(change.NewValue.JSONValue), &ldapSettings); err != nil {
		context.GetLogger(w.ctx).Warnf("unable to decode new LDAP settings property value: %s", err)
		// We can't understand the new value but the LDAP sync cron has
		// been canceled anyway.
		return
	}

	// Get the current auth config.
	authConfig, err := config.GetAuthConfig(w.schemaMgr)
	if err != nil {
		context.GetLogger(w.ctx).Errorf("unable to get system LDAP config: %s", err)
		return
	}

	if authConfig.Backend != config.AuthBackendLDAP {
		// The current backend is not the LDAP backend so we do not
		// need to schedule sync jobs.
		return
	}

	cronID := cronIDLDAPSync
	action := ActionLdapSync
	cronSpec := ldapSettings.SyncSchedule

	subCtx := context.WithValue(w.ctx, "cronID", cronID)
	subCtx = context.WithLogger(subCtx, context.GetLogger(subCtx, "cronID"))

	w.addScheduler(subCtx, cronID, action, cronSpec)
}

// addScheduler creates and starts a new scheduler with the given cronID using
// the given cronSpec schedule to create jobs to run the given action. The
// scheduler is added to this worker's set of currently running scheduler. If
// there is an error parsing the given cronSpec then an error is logged and the
// scheduler does not run.
func (w *worker) addScheduler(ctx context.Context, cronID, action, cronSpec string) {
	schedule, err := cron.Parse(cronSpec)
	if err != nil {
		// How did this get through validation?
		context.GetLogger(ctx).Errorf("invalid cron spec %q: %s", cronSpec, err)
		return
	}

	scheduler := &cronScheduler{
		ctx:          ctx,
		worker:       w,
		cronID:       cronID,
		schedule:     schedule,
		action:       action,
		signalCancel: make(chan struct{}),
	}

	w.schedulers[cronID] = scheduler

	scheduler.Add(1)

	context.GetLogger(ctx).Info("cron scheduled")

	go scheduler.run()
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

type cronScheduler struct {
	ctx          context.Context
	worker       *worker
	cronID       string
	schedule     cron.Schedule
	action       string
	signalCancel chan struct{}

	sync.Mutex     // Used to protect the isCanceled field.
	sync.WaitGroup // Used to signal that the scheduler is done.
	isCanceled     bool
}

func (s *cronScheduler) run() {
	// Signal to any goroutines waiting for the scheduler to shutdown.
	defer s.Done()

	for {
		now := time.Now()
		scheduledAt := s.schedule.Next(now)

		var (
			timer *time.Timer
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
			timer = time.NewTimer(scheduledAt.Sub(now))
			wait = timer.C
		}

		select {
		case <-wait:
			context.GetLogger(s.ctx).Infof("timer fired to run job %s scheduled at %s", s.action, scheduledAt)

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
		CronID:      s.cronID,
		Status:      schema.JobStatusRunning,
		ScheduledAt: scheduledAt,
		LastUpdated: scheduledAt,
		Action:      s.action,
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

	s.worker.addJob(subCtx, job.ID, job.Action)
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
