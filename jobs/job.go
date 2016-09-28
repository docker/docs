package jobs

// Job defines a task which may be scheduled for multiple runs, typically by a
// JobRunner.
type Job interface {
	// IsReady returns whether the job is ready to be run at the current time
	// (used for scheduling at fixed intervals or times).
	IsReady() bool
	// Run synchronously runs the job, returning an error when finished or if it
	// fails to start.
	Run(bool) error
}
