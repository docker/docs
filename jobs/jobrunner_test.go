package jobs

import (
	"errors"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/stretchr/testify/assert"
)

type testJob struct {
	isReady  bool
	run      func(this *testJob) error
	timesRan int
}

func (j *testJob) IsReady() bool {
	return j.isReady
}

func (j *testJob) Run(manual bool) error {
	err := j.run(j)
	j.timesRan++
	return err
}

func TestJobRunnerNoJobsReady(t *testing.T) {
	c := clock.NewMock()
	jr := NewRunner(c)
	job := &testJob{
		isReady: false,
		run: func(this *testJob) error {
			t.Error("this job wasn't supposed to run")
			return errors.New("this job wasn't supposed to run")
		},
	}
	err := jr.AddJob("shouldntRun", job)
	assert.Nil(t, err)

	err = jr.Start()
	assert.Nil(t, err)

	c.Add(2 * time.Minute)
	assert.Equal(t, 0, job.timesRan)
}

func TestJobRunnerRunsJob(t *testing.T) {
	c := clock.NewMock()
	jr := NewRunner(c)
	job := &testJob{
		isReady: true,
		run: func(this *testJob) error {
			this.isReady = false
			return nil
		},
	}
	err := jr.AddJob("runOnce", job)
	assert.Nil(t, err)

	err = jr.Start()
	assert.Nil(t, err)

	c.Add(time.Minute + time.Second)

	assert.Equal(t, 1, job.timesRan)
}

func TestJobRunnerRunsJobMultipleTimes(t *testing.T) {
	c := clock.NewMock()
	jr := NewRunner(c)
	job := &testJob{
		isReady: true,
		run: func(this *testJob) error {
			if this.timesRan == 2 {
				this.isReady = false
			}
			return nil
		},
	}
	err := jr.AddJob("run3Times", job)
	assert.Nil(t, err)

	err = jr.Start()
	assert.Nil(t, err)

	c.Add(5 * time.Minute)

	assert.Equal(t, 3, job.timesRan)
}

func TestJobRunnerRunsSimultaneousJobs(t *testing.T) {
	c := clock.NewMock()
	jr := NewRunner(c)
	job1 := &testJob{
		isReady: true,
		run: func(this *testJob) error {
			this.isReady = false
			c.Sleep(2 * time.Minute)
			return nil
		},
	}
	job2 := &testJob{
		isReady: true,
		run: func(this *testJob) error {
			this.isReady = false
			c.Sleep(2 * time.Minute)
			return nil
		},
	}
	err := jr.AddJob("job1", job1)
	assert.Nil(t, err)

	err = jr.AddJob("job2", job2)
	assert.Nil(t, err)

	err = jr.Start()
	assert.Nil(t, err)

	c.Add(time.Minute + time.Second)

	assert.Equal(t, 0, job1.timesRan)
	assert.Equal(t, 0, job2.timesRan)

	statusMap, err := jr.Status()
	assert.Nil(t, err)

	assert.Equal(t, JobStatus{Status: StatusRunning}, statusMap["job1"])
	assert.Equal(t, JobStatus{Status: StatusRunning}, statusMap["job2"])

	c.Add(2*time.Minute + time.Second)

	assert.Equal(t, 1, job1.timesRan)
	assert.Equal(t, 1, job2.timesRan)

	statusMap, err = jr.Status()
	assert.Nil(t, err)

	assert.Equal(t, JobStatus{Status: StatusWaiting}, statusMap["job1"])
	assert.Equal(t, JobStatus{Status: StatusWaiting}, statusMap["job2"])
}

func TestJobRunnerRecordsFailure(t *testing.T) {
	c := clock.NewMock()
	jr := NewRunner(c)
	job := &testJob{
		isReady: true,
		run: func(this *testJob) error {
			return errors.New("expected failure")
		},
	}
	err := jr.AddJob("shouldFail", job)
	assert.Nil(t, err)

	err = jr.Start()
	assert.Nil(t, err)

	c.Add(time.Minute + time.Second)

	assert.Equal(t, 1, job.timesRan)

	statusMap, err := jr.Status()
	assert.Nil(t, err)

	assert.Equal(t, JobStatus{Status: StatusError, Details: "expected failure"}, statusMap["shouldFail"])
}

func TestRemoveJob(t *testing.T) {
	c := clock.NewMock()
	jr := NewRunner(c)
	job := &testJob{
		isReady: true,
		run: func(this *testJob) error {
			return nil
		},
	}
	err := jr.AddJob("toRemove", job)
	assert.Nil(t, err)

	err = jr.RemoveJob("toRemove")
	assert.Nil(t, err)

	err = jr.Start()
	assert.Nil(t, err)

	c.Add(time.Minute + time.Second)

	assert.Equal(t, 0, job.timesRan)

	statusMap, err := jr.Status()
	assert.Nil(t, err)

	assert.Equal(t, 0, len(statusMap))

}

func TestRemoveRunningJob(t *testing.T) {
	c := clock.NewMock()
	jr := NewRunner(c)
	job := &testJob{
		isReady: true,
		run: func(this *testJob) error {
			this.isReady = false
			c.Sleep(2 * time.Minute)
			return nil
		},
	}
	err := jr.AddJob("toRemove", job)
	assert.Nil(t, err)

	err = jr.Start()
	assert.Nil(t, err)

	c.Add(time.Minute + time.Second)

	assert.Equal(t, 0, job.timesRan)

	statusMap, err := jr.Status()
	assert.Nil(t, err)

	assert.Equal(t, JobStatus{Status: StatusRunning}, statusMap["toRemove"])

	err = jr.RemoveJob("toRemove")
	assert.Nil(t, err)

	c.Add(2*time.Minute + time.Second)

	assert.Equal(t, 1, job.timesRan)

	statusMap, err = jr.Status()
	assert.Nil(t, err)

	assert.Equal(t, 0, len(statusMap))
}

func TestAddDuplicateJob(t *testing.T) {
	c := clock.NewMock()
	jr := NewRunner(c)
	job := &testJob{
		isReady: true,
		run: func(this *testJob) error {
			return nil
		},
	}
	err := jr.AddJob("duplicate", job)
	assert.Nil(t, err)

	err = jr.AddJob("duplicate", job)
	assert.Equal(t, ErrJobAlreadyExists, err)
}

func TestRemoveNonexistentJob(t *testing.T) {
	c := clock.NewMock()
	jr := NewRunner(c)

	err := jr.RemoveJob("nonexistent")
	assert.Equal(t, ErrJobNotFound, err)
}

func TestStopJobRunner(t *testing.T) {
	c := clock.NewMock()
	jr := NewRunner(c)

	err := jr.Stop()
	assert.Equal(t, ErrAlreadyStopped, err)

	job := &testJob{
		isReady: true,
		run: func(this *testJob) error {
			return nil
		},
	}
	err = jr.AddJob("job", job)
	assert.Nil(t, err)

	err = jr.Start()
	assert.Nil(t, err)

	c.Add(time.Minute + time.Second)

	assert.Equal(t, 1, job.timesRan)

	err = jr.Stop()
	assert.Nil(t, err)

	c.Add(time.Minute + time.Second)

	assert.Equal(t, 1, job.timesRan)

	err = jr.Stop()
	assert.Equal(t, ErrAlreadyStopped, err)

	err = jr.Start()
	assert.Nil(t, err)

	c.Add(time.Minute + time.Second)

	assert.Equal(t, 2, job.timesRan)
}

func TestMultiStartJobRunner(t *testing.T) {
	c := clock.NewMock()
	jr := NewRunner(c)

	err := jr.Start()
	assert.Nil(t, err)

	err = jr.Start()
	assert.Equal(t, ErrAlreadyStarted, err)
}

func TestRunNowAlreadyRunning(t *testing.T) {
	c := clock.NewMock()
	jr := NewRunner(c)
	job := &testJob{
		isReady: true,
		run: func(this *testJob) error {
			c.Sleep(time.Minute)
			return nil
		},
	}
	err := jr.AddJob("job", job)
	assert.Nil(t, err)

	err = jr.Start()
	assert.Nil(t, err)

	c.Add(time.Minute + time.Second)

	assert.Equal(t, 0, job.timesRan)

	statusMap, err := jr.Status()
	assert.Nil(t, err)

	assert.Equal(t, JobStatus{Status: StatusRunning}, statusMap["job"])

	err = jr.RunNow("job")
	assert.Equal(t, ErrJobAlreadyRunning, err)

	c.Add(time.Minute + time.Second)

	assert.Equal(t, 1, job.timesRan)
}

func TestRunNowUnscheduled(t *testing.T) {
	c := clock.NewMock()
	jr := NewRunner(c)
	job := &testJob{
		isReady: false,
		run: func(this *testJob) error {
			c.Sleep(time.Minute)
			return nil
		},
	}
	err := jr.AddJob("job", job)
	assert.Nil(t, err)

	err = jr.Start()
	assert.Nil(t, err)

	c.Add(time.Minute + time.Second)

	assert.Equal(t, 0, job.timesRan)

	statusMap, err := jr.Status()
	assert.Nil(t, err)

	assert.Equal(t, JobStatus{Status: StatusWaiting}, statusMap["job"])

	err = jr.RunNow("job")
	assert.Nil(t, err)

	err = jr.RunNow("job")
	assert.Equal(t, ErrJobAlreadyRunning, err)

	c.Add(2 * time.Minute)

	assert.Equal(t, 1, job.timesRan)
}

func TestRunNowInvalidJob(t *testing.T) {
	c := clock.NewMock()
	jr := NewRunner(c)

	err := jr.RunNow("invalid")
	assert.Equal(t, ErrJobNotFound, err)
}
