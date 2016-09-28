package jobrunner

import (
	"fmt"
	"io"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/WatchBeam/clock"
	"github.com/docker/dhe-deploy/pkg/jobrunner-framework/constants"
	"github.com/docker/dhe-deploy/pkg/jobrunner-framework/schema"
	"github.com/docker/dhe-deploy/pkg/jobrunner-framework/worker"
	"github.com/docker/distribution/context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vikstrous/mock/gomock"
)

var ActionSleep = "sleep"
var ActionExclusiveSleep = "exclusive_sleep"
var ActionExclusiveRecoverableSleep = "exclusive_recoverable_sleep"
var ActionRecoverableSleep = "recoverable_sleep"
var ActionEcho = "echo"
var ActionFail = "fail"

// TODO: test all the special flags
var RegisteredActions = map[string]worker.ActionInfo{
	ActionFail: {
		Command: "aaaaaaa",
		Args:    []string{""},
		CleanupJob: &schema.Job{
			Action: ActionEcho,
		},
	},
	ActionEcho: {
		Command: "echo",
		Args:    []string{"success"},
	},
	ActionSleep: {
		Command: "sleep",
		Args:    []string{"60"},
	},
	ActionRecoverableSleep: {
		Command: "sleep",
		Args:    []string{"60"},
		CleanupJob: &schema.Job{
			Action: ActionEcho,
		},
	},
	ActionExclusiveSleep: {
		Command:   "sleep",
		Exclusive: true,
		Args:      []string{"60"},
	},
	ActionExclusiveRecoverableSleep: {
		Command:   "sleep",
		Exclusive: true,
		Args:      []string{"60"},
		CleanupJob: &schema.Job{
			Action: ActionEcho,
		},
	},
}

type TestWorker struct {
	SchemaMgr             *schema.MockJobrunnerManager
	Worker                worker.Worker
	JobsCh                chan schema.JobChange
	JobsWaiter            *Waiter
	CancelJobsCh          chan schema.JobChange
	CancelJobsWaiter      *Waiter
	RecoverableJobsCh     chan schema.JobChange
	RecoverableJobsWaiter *Waiter
	CronsCh               chan schema.CronChange
	CronWaiter            *Waiter
	ID                    string
	Executors             []worker.Executor
	Clock                 *clock.MockClock
}

func newTestWorker(mockCtrl *gomock.Controller, id string) *TestWorker {
	tw := &TestWorker{}
	tw.SchemaMgr = schema.NewMockJobrunnerManager(mockCtrl)
	date, _ := time.Parse(time.UnixDate, "0")
	tw.Clock = clock.NewMockClock(date)
	tw.ID = id
	tw.Worker = worker.New(context.Background(), tw.SchemaMgr, tw.ID, RegisteredActions, tw.Clock, GinkgoRecover, func() worker.Executor {
		// this function just serves as a way for us to shim a mock executor into the worker
		// if necessary
		if len(tw.Executors) == 0 {
			return worker.DefaultNewExecutor()
		} else {
			executor := tw.Executors[0]
			tw.Executors = tw.Executors[1:]
			return executor
		}
	})
	tw.JobsCh = make(chan schema.JobChange, 1)
	tw.JobsWaiter = NewWaiter("jobs channel")
	tw.CancelJobsCh = make(chan schema.JobChange, 1)
	tw.CancelJobsWaiter = NewWaiter("cancel jobs channel")
	tw.RecoverableJobsCh = make(chan schema.JobChange, 1)
	tw.RecoverableJobsWaiter = NewWaiter("recoverable jobs channel")
	tw.CronsCh = make(chan schema.CronChange, 1)
	tw.CronWaiter = NewWaiter("crons channel")
	return tw
}
func (tw *TestWorker) setUpChannelsAndWaiters() {
	tw.SchemaMgr.EXPECT().GetNewJobChanges().Return(tw.JobsCh, io.Closer(tw.JobsWaiter), nil).AnyTimes()
	tw.SchemaMgr.EXPECT().GetOwnJobCancellations(tw.ID).Return(tw.CancelJobsCh, io.Closer(tw.CancelJobsWaiter), nil).AnyTimes()
	tw.SchemaMgr.EXPECT().GetRecoverableJobChanges(gomock.Any()).Return(tw.RecoverableJobsCh, io.Closer(tw.RecoverableJobsWaiter), nil).AnyTimes()
}

func (tw *TestWorker) waitForFinish() {
	WaitersWait(time.Second, tw.JobsWaiter, tw.RecoverableJobsWaiter, tw.CancelJobsWaiter, tw.CronWaiter)
	Expect(tw.Worker.WaitTimeout(time.Second)).To(BeNil())
}

func TestJobrunner(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Jobrunner Suite")
}

// Note 1:
// We always have to have a short sleep after a job succeeds, before calling shutdown.
// If we call shutdown before the job is removed from the list of active jobs for
// the worker, the worker will try to cancel it during shutdown, which will cause
// the test to fail. In practice this time period is very small, but this can
// also happen in real life. If it happens, the job will be marked as failed even
// though it succeeded.

var _ = Describe("Multiple Workers", func() {
	var (
		mockCtrl *gomock.Controller
		tw1      *TestWorker
		tw2      *TestWorker
	)
	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		tw1 = newTestWorker(mockCtrl, "worker1")
		tw1.SchemaMgr.EXPECT().SafeGetActionConfig(gomock.Any()).Return(&schema.ActionConfig{HeartbeatTimeout: constants.DefaultHeartbeatTimeout.String()}, nil).AnyTimes()
		tw2 = newTestWorker(mockCtrl, "worker2")
		tw2.SchemaMgr.EXPECT().SafeGetActionConfig(gomock.Any()).Return(&schema.ActionConfig{HeartbeatTimeout: constants.DefaultHeartbeatTimeout.String()}, nil).AnyTimes()
	})
	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("running crons at the same time", func() {

		BeforeEach(func() {
			for _, tw := range []*TestWorker{tw1, tw2} {
				tw.setUpChannelsAndWaiters()
				tw.SchemaMgr.EXPECT().GetCronChanges().Return(tw.CronsCh, io.Closer(tw.CronWaiter), nil).AnyTimes()
				tw.SchemaMgr.EXPECT().GetMostRecentlyScheduledJobsForWorker(tw.ID, uint(0), uint(0)).Return([]schema.Job{}, nil).AnyTimes()
			}
		})

		It("should be able to schedule a job at the right time and execute it only once", func() {
			InitAndAsyncRunForTest(tw1.Worker)
			InitAndAsyncRunForTest(tw2.Worker)
			scheduledSecond := time.Duration(8)
			cron := &schema.Cron{
				Action:   ActionEcho,
				Schedule: fmt.Sprintf("%d * * * * *", scheduledSecond),
			}

			tw1.SchemaMgr.EXPECT().CreateJob(&schema.Job{
				WorkerID:    tw1.ID,
				CronID:      "",
				Status:      schema.JobStatusRunning,
				ScheduledAt: tw1.Clock.Now().Add(time.Second * scheduledSecond),
				LastUpdated: tw1.Clock.Now().Add(time.Second * scheduledSecond),
				Action:      ActionEcho,
			}).Return(nil).Times(1)

			tw2.SchemaMgr.EXPECT().CreateJob(&schema.Job{
				WorkerID:    tw2.ID,
				CronID:      "",
				Status:      schema.JobStatusRunning,
				ScheduledAt: tw2.Clock.Now().Add(time.Second * scheduledSecond),
				LastUpdated: tw2.Clock.Now().Add(time.Second * scheduledSecond),
				Action:      ActionEcho,
			}).Return(fmt.Errorf("tw1 already took the job")).Times(1)

			// wait for this function to be called
			jobDone := NewWaiter("job done")
			tw1.SchemaMgr.EXPECT().UpdateJobStatus("", schema.JobStatusDone).Do(func(_ string, _ string) error {
				jobDone.Close()
				return nil
			}).Times(1)
			tw1.SchemaMgr.EXPECT().InsertJobLog("", "success", 0).Return(nil).Times(1)

			tw1.CronsCh <- schema.CronChange{
				OldValue: nil,
				NewValue: cron,
			}
			tw2.CronsCh <- schema.CronChange{
				OldValue: nil,
				NewValue: cron,
			}
			// XXX: we need to make sure it picks up the cron and calculates the wait time before we increment the time
			// The only way to know this right now is to wait for it to call NewTimer, so we wait instead. Sorry.
			time.Sleep(time.Millisecond * 100)
			// by moving time for tw1 we make sure it's the one that picks up the cron first
			tw1.Clock.AddTime(time.Second * scheduledSecond)
			tw2.Clock.AddTime(time.Second * scheduledSecond)

			// TODO: consider how we can mock the executor to make sure this executes echo
			WaitersWait(time.Second, jobDone)
			// TODO: come up with a better way to make sure tw2 doesn't try to create the job too
			time.Sleep(time.Millisecond * 100)

			tw1.Worker.Shutdown()
			tw2.Worker.Shutdown()

			tw1.waitForFinish()
			tw2.waitForFinish()
		})
	})

	Describe("picking up jobs at the same time", func() {

		Describe("with a new job", func() {
			BeforeEach(func() {
				for _, tw := range []*TestWorker{tw1, tw2} {
					// no old jobs
					tw.SchemaMgr.EXPECT().GetMostRecentlyScheduledJobsForWorker(tw.ID, uint(0), uint(0)).Return([]schema.Job{}, nil).AnyTimes()
					// jobsch allows us to send out a job to the worker to make it try to take it
					tw.setUpChannelsAndWaiters()
					// no crons
					tw.SchemaMgr.EXPECT().GetCronChanges().Return(make(chan schema.CronChange), io.Closer(tw.CronWaiter), nil).AnyTimes()
				}
			})

			// TODO: test a worker's ability to notice that it was separated from the rest and cancel jobs (separate from noticing that it's proclaimed a zombie)
			It("should be able to run a job, go MIA, have another worker notice the stale job and fail the job", func() {
				InitAndAsyncRunForTest(tw1.Worker)
				InitAndAsyncRunForTest(tw2.Worker)

				job := schema.Job{
					ID:     "job1",
					Status: schema.JobStatusWaiting,
					Action: ActionRecoverableSleep,
				}
				heartbeatenJob := job
				heartbeatenJob.HeartbeatExpiration = heartbeatenJob.LastUpdated.Add(constants.HeartbeatInterval + constants.DefaultHeartbeatTimeout)
				heartbeatenJob.WorkerID = tw1.ID
				deadJob := heartbeatenJob
				deadJob.Status = schema.JobStatusWorkerDead
				resurrectionJob := deadJob
				resurrectionJob.Status = schema.JobStatusWorkerResurrection

				// # First we create a job and heart beat it once

				// both will claim it, one of them needs to be told to go away
				tw1.SchemaMgr.EXPECT().ClaimJob(job.ID, tw1.ID, gomock.Any()).Return(&job, nil).MaxTimes(1)
				tw2.SchemaMgr.EXPECT().ClaimJob(job.ID, tw2.ID, gomock.Any()).Return(nil, fmt.Errorf("u raff u ruse")).MaxTimes(1)

				// we set up the mock waiting for a heartbeat
				heartbeat := NewWaiter("heartbeat")
				tw1.SchemaMgr.EXPECT().HeartbeatJob(job.ID, tw1.Clock.Now().Add(constants.HeartbeatInterval+constants.DefaultHeartbeatTimeout)).Do(func(_ string, _ time.Time) error {
					defer heartbeat.Close()
					// we deliver the heartbeat notification to all workers
					tw1.RecoverableJobsCh <- schema.JobChange{
						OldValue: &job,
						NewValue: &heartbeatenJob,
					}
					tw2.RecoverableJobsCh <- schema.JobChange{
						OldValue: &job,
						NewValue: &heartbeatenJob,
					}
					return nil
				})

				// kick off the jobs
				tw1.JobsCh <- schema.JobChange{
					OldValue: nil,
					NewValue: &job,
				}
				tw2.JobsCh <- schema.JobChange{
					OldValue: nil,
					NewValue: &job,
				}
				// wait for the jobs to be actually picked up and created and for heartbeat timer to start
				time.Sleep(100 * time.Millisecond)

				// now that the first worker has claimed the job, we advance time to make sure they send a heartbeat
				tw1.Clock.AddTime(constants.HeartbeatInterval)
				// in a mock function above we tell the second worker about the heartbeat from the first worker, so we need to wait for that to happen
				WaitersWait(time.Second, heartbeat)
				// now the second worker notices the heartbeat and is happy
				tw2.Clock.AddTime(constants.HeartbeatInterval)

				// Wait here to reduce raciness of test in case the worker is proclaimed dead prematurely
				time.Sleep(10 * time.Millisecond)

				// # Next we make the job time out and watch a recovery job start

				// The second worker should now fail the job as soon as we let it run
				workerDead := NewWaiter("update job status to indicate worker dead")
				tw2.SchemaMgr.EXPECT().UpdateJobStatus(job.ID, schema.JobStatusWorkerDead).Do(func(_ string, _ string) error {
					defer workerDead.Close()

					// now that we've failed a job, tell both workers about it and see if they create the recovery job
					tw1.RecoverableJobsCh <- schema.JobChange{
						OldValue: &heartbeatenJob,
						NewValue: &deadJob,
					}
					tw2.RecoverableJobsCh <- schema.JobChange{
						OldValue: &heartbeatenJob,
						NewValue: &deadJob,
					}
					return nil
				}).Times(1)
				tw1.SchemaMgr.EXPECT().InsertJobLog(job.ID, worker.ForceKillMsg, gomock.Any()).Return(nil).Times(1)
				tw1.SchemaMgr.EXPECT().InsertJobLog(job.ID, worker.ExitStatusMsg(9), gomock.Any()).Return(nil).Times(1)

				// at this point a recovery job is launched
				scheduled := NewWaiter("recovery job scheduled")
				job2 := &schema.Job{
					WorkerID:       tw2.ID,
					RecoveryFromID: job.ID,
					Status:         schema.JobStatusRunning,
					ScheduledAt:    tw2.Clock.Now().Add(constants.DefaultHeartbeatTimeout),
					LastUpdated:    tw2.Clock.Now().Add(constants.DefaultHeartbeatTimeout),
					Action:         ActionEcho,
				}
				tw2.SchemaMgr.EXPECT().CreateJob(job2).Do(func(_ *schema.Job) error {
					scheduled.Close()
					return nil
				}).Times(1)
				// wait for this function to be called
				jobDone := NewWaiter("update job status")
				tw2.SchemaMgr.EXPECT().UpdateJobStatus(job2.ID, schema.JobStatusDone).Do(func(_ string, _ string) error {
					jobDone.Close()
					return nil
				}).Times(1)
				tw2.SchemaMgr.EXPECT().InsertJobLog(job2.ID, "success", 0).Return(nil).Times(1)

				// now that the first worker's job is marked as dead, we let it know that it's been buried and wait for it to note its resurrection and kill itself again, kinda like Jesus
				workerResurrection := NewWaiter("update job status to indicate worker resurrection")
				tw1.SchemaMgr.EXPECT().UpdateJobStatus(job.ID, schema.JobStatusWorkerResurrection).Do(func(_ string, _ string) error {
					workerResurrection.Close()

					// now that we've failed a job, tell both workers about it and see if they create the recovery job
					tw1.RecoverableJobsCh <- schema.JobChange{
						OldValue: &deadJob,
						NewValue: &resurrectionJob,
					}
					tw2.RecoverableJobsCh <- schema.JobChange{
						OldValue: &deadJob,
						NewValue: &resurrectionJob,
					}
					return nil
				}).Times(1)

				// Now we freeze the first worker in time and the second worker notices that the first worker hasn't updated its job, then does all of the things above
				tw2.Clock.AddTime(constants.DefaultHeartbeatTimeout)

				// wait for it to finish marking the first worker's job as dead and running the recovery job
				WaitersWait(time.Second, workerDead)
				WaitersWait(time.Second, scheduled)
				WaitersWait(time.Second, jobDone)

				// everyone should note the resurrection and feel bad
				WaitersWait(time.Second, workerResurrection)

				tw1.Worker.Shutdown()
				tw2.Worker.Shutdown()

				tw1.waitForFinish()
				tw2.waitForFinish()
			})

			It("only one of them should pick it up and run it and it should succeed", func() {
				// XXX: note that this doesn't test the rethinkdb part of the process
				InitAndAsyncRunForTest(tw1.Worker)
				InitAndAsyncRunForTest(tw2.Worker)

				job := &schema.Job{
					Status: schema.JobStatusWaiting,
					Action: ActionEcho,
				}

				// both will claim it, one of them needs to be told to go away
				tw1.SchemaMgr.EXPECT().ClaimJob("", tw1.ID, gomock.Any()).Return(job, nil).Times(1)
				tw2.SchemaMgr.EXPECT().ClaimJob("", tw2.ID, gomock.Any()).Return(nil, fmt.Errorf("u raff u ruse")).Times(1)

				// wait for this function to be called
				updated := NewWaiter("update job status")
				tw1.SchemaMgr.EXPECT().UpdateJobStatus("", schema.JobStatusDone).Do(func(_ string, _ string) error {
					updated.Close()
					return nil
				}).Times(1)
				tw1.SchemaMgr.EXPECT().InsertJobLog("", "success", 0).Return(nil).Times(1)

				tw1.JobsCh <- schema.JobChange{
					OldValue: nil,
					NewValue: job,
				}
				tw2.JobsCh <- schema.JobChange{
					OldValue: nil,
					NewValue: job,
				}

				// wait for it to succeed at running echo
				WaitersWait(time.Second, updated)

				// TODO: consider how we can mock the executor to make sure this executes echo

				tw1.Worker.Shutdown()
				tw2.Worker.Shutdown()

				tw1.waitForFinish()
				tw2.waitForFinish()
			})

			It("while one is running an exclusive job, the other should not be able to", func() {
				// XXX: note that this doesn't test the rethinkdb part of the process
				InitAndAsyncRunForTest(tw1.Worker)
				InitAndAsyncRunForTest(tw2.Worker)

				job := &schema.Job{
					ID:     "one",
					Status: schema.JobStatusWaiting,
					Action: ActionExclusiveSleep,
				}
				job2 := &schema.Job{
					ID:     "two",
					Status: schema.JobStatusWaiting,
					Action: ActionExclusiveSleep,
				}

				// both will claim it, one of them needs to be told to go away
				// these functions will be called twice, and we need to return different things both times, so we
				// need them to be implemented as functions
				tw1.SchemaMgr.EXPECT().ClaimJob("one", tw1.ID, gomock.Any()).Return(job, nil).Times(1)
				tw1.SchemaMgr.EXPECT().ClaimJob("two", tw1.ID, gomock.Any()).Return(nil, fmt.Errorf("nope")).Times(1)
				tw2.SchemaMgr.EXPECT().ClaimJob("one", tw2.ID, gomock.Any()).Return(nil, fmt.Errorf("nope2")).Times(1)
				tw2.SchemaMgr.EXPECT().ClaimJob("two", tw2.ID, gomock.Any()).Return(job2, nil).Times(1)

				// tw1 is in the clear, but tw2 shouldn't execute the job because it should already be running at that time
				tw1.SchemaMgr.EXPECT().CountJobsWithActionStatus(ActionExclusiveSleep, schema.JobStatusRunning).Return(uint(1), nil).Times(1)
				tw2.SchemaMgr.EXPECT().CountJobsWithActionStatus(ActionExclusiveSleep, schema.JobStatusRunning).Return(uint(2), nil).Times(1)

				// make sure we can wait for the first job to complete successfully
				updated := NewWaiter("update job status")
				tw1.SchemaMgr.EXPECT().UpdateJobStatus("one", schema.JobStatusDone).Do(func(_ string, _ string) error {
					updated.Close()
					return nil
				}).Times(1)

				// make sure we can control when the first job starts and ends
				ch1 := make(chan struct{})
				executor, ch2 := MockDoubleSyncExecutor(mockCtrl, ch1)
				tw1.Executors = []worker.Executor{executor}

				// send the first job to both; only the first one will pick it up
				tw1.JobsCh <- schema.JobChange{
					OldValue: nil,
					NewValue: job,
				}
				tw2.JobsCh <- schema.JobChange{
					OldValue: nil,
					NewValue: job,
				}

				// wait for it to start running
				ChWait("execution started", ch2, time.Second)

				// set up expectation for for the second job to fail
				rejected := NewWaiter("reject second job")
				tw2.SchemaMgr.EXPECT().UpdateJobStatus("two", schema.JobStatusConflict).Do(func(_ string, _ string) error {
					rejected.Close()
					return nil
				}).Times(1)

				// send the second job to both; only the second will pick it up
				tw1.JobsCh <- schema.JobChange{
					OldValue: nil,
					NewValue: job2,
				}
				tw2.JobsCh <- schema.JobChange{
					OldValue: nil,
					NewValue: job2,
				}

				// wait for the second job to fail due to the exclusivity check
				WaitersWait(time.Second, rejected)

				// tell the executor (first job) to finish running
				close(ch1)

				// wait for it to finish running
				WaitersWait(time.Second, updated)

				tw1.Worker.Shutdown()
				tw2.Worker.Shutdown()

				tw1.waitForFinish()
				tw2.waitForFinish()
			})

			// in this test we don't actually create the first worker, but we tell the second worker that
			// another worker is running the cleanup job
			It("while one is running an exclusive job's cleanup job, the other should not be able to start", func() {
				InitAndAsyncRunForTest(tw2.Worker)

				job2 := &schema.Job{
					ID:     "two",
					Status: schema.JobStatusWaiting,
					Action: ActionExclusiveRecoverableSleep,
				}

				tw2.SchemaMgr.EXPECT().ClaimJob("two", tw2.ID, gomock.Any()).Return(job2, nil).Times(1)

				// tw1 is in the clear, but tw2 shouldn't execute the job because it should already be running at that time
				tw2.SchemaMgr.EXPECT().CountJobsWithActionStatus(ActionExclusiveRecoverableSleep, schema.JobStatusRunning).Return(uint(1), nil).Times(1)
				tw2.SchemaMgr.EXPECT().CountJobsWithActionStatus(ActionEcho, schema.JobStatusRunning).Return(uint(2), nil).Times(1)

				// set up expectation for for the second job to fail
				rejected := NewWaiter("reject second job")
				tw2.SchemaMgr.EXPECT().UpdateJobStatus("two", schema.JobStatusConflict).Do(func(_ string, _ string) error {
					rejected.Close()
					return nil
				}).Times(1)

				// now try to run the job - it should be rejected because its cleanup job is still running
				tw2.JobsCh <- schema.JobChange{
					OldValue: nil,
					NewValue: job2,
				}

				// wait for the second job to fail due to the exclusivity check
				WaitersWait(time.Second, rejected)

				tw2.Worker.Shutdown()

				tw2.waitForFinish()
			})

		})
	})
})

var _ = Describe("One Worker", func() {
	var (
		mockCtrl *gomock.Controller
		tw1      *TestWorker
	)
	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		tw1 = newTestWorker(mockCtrl, "worker1")
		tw1.SchemaMgr.EXPECT().SafeGetActionConfig(gomock.Any()).Return(&schema.ActionConfig{HeartbeatTimeout: constants.DefaultHeartbeatTimeout.String()}, nil).AnyTimes()
	})
	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("managing crons and starting scheduled jobs", func() {

		Context("with some new crons and no manual jobs", func() {

			BeforeEach(func() {
				tw1.setUpChannelsAndWaiters()
				tw1.SchemaMgr.EXPECT().GetCronChanges().Return(tw1.CronsCh, io.Closer(tw1.CronWaiter), nil).AnyTimes()
				tw1.SchemaMgr.EXPECT().GetMostRecentlyScheduledJobsForWorker(tw1.ID, uint(0), uint(0)).Return([]schema.Job{}, nil).AnyTimes()
			})

			It("should be able to schedule a job at the right time and execute it (twice)", func() {
				InitAndAsyncRunForTest(tw1.Worker)
				cron := &schema.Cron{
					Action:   ActionEcho,
					Schedule: fmt.Sprintf("* * * * * *"),
				}

				tw1.SchemaMgr.EXPECT().CreateJob(&schema.Job{
					WorkerID:    tw1.ID,
					CronID:      "",
					Status:      schema.JobStatusRunning,
					ScheduledAt: tw1.Clock.Now().Add(time.Second),
					LastUpdated: tw1.Clock.Now().Add(time.Second),
					Action:      ActionEcho,
				}).Return(nil).Times(1)

				tw1.SchemaMgr.EXPECT().CreateJob(&schema.Job{
					WorkerID:    tw1.ID,
					CronID:      "",
					Status:      schema.JobStatusRunning,
					ScheduledAt: tw1.Clock.Now().Add(time.Second * 2),
					LastUpdated: tw1.Clock.Now().Add(time.Second * 2),
					Action:      ActionEcho,
				}).Return(nil).Times(1)

				// wait for this function to be called
				jobDone := NewWaiter("job done")
				job2Done := NewWaiter("job 2 done")
				first := true
				tw1.SchemaMgr.EXPECT().UpdateJobStatus("", schema.JobStatusDone).Do(func(_ string, _ string) error {
					if first {
						jobDone.Close()
						first = false
					} else {
						job2Done.Close()
					}
					return nil
				}).Times(2)
				tw1.SchemaMgr.EXPECT().InsertJobLog("", "success", 0).Return(nil).Times(2)

				exec1 := MockEchoExecutor(mockCtrl, nil)
				exec2 := MockEchoExecutor(mockCtrl, nil)
				tw1.Executors = []worker.Executor{exec1, exec2}

				tw1.CronsCh <- schema.CronChange{
					OldValue: nil,
					NewValue: cron,
				}
				// XXX: we need to make sure it picks up the cron and calculates the wait time before we increment the time
				// The only way to know this right now is to wait for it to call NewTimer, so we wait instead. Sorry.
				time.Sleep(time.Millisecond * 100)
				tw1.Clock.AddTime(time.Second)
				// wait for it to execute the job and be ready for another second to tick
				time.Sleep(time.Millisecond * 100)
				tw1.Clock.AddTime(time.Second)

				WaitersWait(time.Second, jobDone, job2Done)

				tw1.Worker.Shutdown()

				tw1.waitForFinish()
			})

			It("should be able to handle a schedule cancellation", func() {
				InitAndAsyncRunForTest(tw1.Worker)
				cron := &schema.Cron{
					Action:   ActionEcho,
					Schedule: "* * * * * *",
				}

				tw1.CronsCh <- schema.CronChange{
					OldValue: nil,
					NewValue: cron,
				}
				tw1.CronsCh <- schema.CronChange{
					OldValue: cron,
					NewValue: nil,
				}
				// XXX: we need to make sure it picks up the cron and calculates the wait time before we increment the time
				// The only way to know this right now is to wait for it to call NewTimer, so we wait instead. Sorry.
				time.Sleep(time.Millisecond * 100)
				tw1.Clock.AddTime(time.Second)
				// we need to wait to make sure it creates the cron before we increment time again
				time.Sleep(time.Millisecond * 100)
				// it should delete the timer now
				tw1.Clock.AddTime(time.Second)
				// wait for things to settle
				time.Sleep(time.Millisecond * 100)
				// make sure it doesn't try to run any jobs
				tw1.Clock.AddTime(time.Second)

				tw1.Worker.Shutdown()

				tw1.waitForFinish()
			})
		})
	})

	Describe("starting and/or picking up jobs", func() {

		Describe("with no new jobs", func() {

			BeforeEach(func() {
				tw1.setUpChannelsAndWaiters()
				tw1.SchemaMgr.EXPECT().GetCronChanges().Return(make(chan schema.CronChange), io.Closer(tw1.CronWaiter), nil).AnyTimes()
			})

			Context("with no jobs", func() {
				BeforeEach(func() {
					tw1.SchemaMgr.EXPECT().GetMostRecentlyScheduledJobsForWorker(tw1.ID, uint(0), uint(0)).Return([]schema.Job{}, nil).Times(1)
				})

				It("should be able to start and stop", func() {
					InitAndAsyncRunForTest(tw1.Worker)

					tw1.Worker.Shutdown()

					tw1.waitForFinish()
				})
			})

			Context("with an unfinished job", func() {
				BeforeEach(func() {
					tw1.SchemaMgr.EXPECT().GetMostRecentlyScheduledJobsForWorker(tw1.ID, uint(0), uint(0)).Return([]schema.Job{
						{
							WorkerID: tw1.ID,
							Status:   schema.JobStatusRunning,
							Action:   ActionSleep,
						},
					}, nil).Times(1)
				})

				It("should mark it as failed", func() {
					tw1.SchemaMgr.EXPECT().UpdateJobStatus("", schema.JobStatusWorkerDead).Return(nil).AnyTimes()

					err := tw1.Worker.Init()
					Expect(err).To(BeNil())
				})

			})
		})
		Describe("with a new job", func() {
			BeforeEach(func() {
				tw1.SchemaMgr.EXPECT().GetMostRecentlyScheduledJobsForWorker(tw1.ID, uint(0), uint(0)).Return([]schema.Job{}, nil).Times(1)
				tw1.setUpChannelsAndWaiters()
				tw1.SchemaMgr.EXPECT().GetCronChanges().Return(make(chan schema.CronChange), io.Closer(tw1.CronWaiter), nil).Times(1)
			})

			It("should be able to pick it up and run it, then cancel it through rethinkdb", func() {
				InitAndAsyncRunForTest(tw1.Worker)
				job := &schema.Job{
					Status: schema.JobStatusWaiting,
					Action: ActionSleep,
				}
				jobCancelRequest := &schema.Job{
					WorkerID: tw1.ID,
					Status:   schema.JobStatusCancelRequest,
					Action:   ActionSleep,
				}

				tw1.SchemaMgr.EXPECT().ClaimJob("", tw1.ID, gomock.Any()).Return(job, nil).Times(1)
				// wait for this function to be called
				jobCancelled := NewWaiter("update job status")
				tw1.SchemaMgr.EXPECT().UpdateJobStatus("", schema.JobStatusCanceled).Do(func(_ string, _ string) error {
					jobCancelled.Close()
					return nil
				}).Times(1)
				tw1.SchemaMgr.EXPECT().InsertJobLog("", worker.ForceKillMsg, gomock.Any()).Return(nil).Times(1)
				tw1.SchemaMgr.EXPECT().InsertJobLog("", worker.ExitStatusMsg(9), gomock.Any()).Return(nil).Times(1)

				tw1.JobsCh <- schema.JobChange{
					OldValue: nil,
					NewValue: job,
				}

				// TODO: use a mock executor and actually wait for execution to start
				time.Sleep(time.Millisecond * 100)

				tw1.CancelJobsCh <- schema.JobChange{
					OldValue: job,
					NewValue: jobCancelRequest,
				}

				// wait for it to cancel the job
				WaitersWait(time.Second, jobCancelled)

				tw1.Worker.Shutdown()

				tw1.waitForFinish()
			})

			It("should be able to pick it up and run it, then cancel it by deleting it", func() {
				InitAndAsyncRunForTest(tw1.Worker)
				job := &schema.Job{
					Status: schema.JobStatusWaiting,
					Action: ActionSleep,
				}

				tw1.SchemaMgr.EXPECT().ClaimJob("", tw1.ID, gomock.Any()).Return(job, nil).Times(1)
				// wait for this function to be called
				jobCancelled := NewWaiter("update job status")
				tw1.SchemaMgr.EXPECT().UpdateJobStatus("", schema.JobStatusDeleted).Do(func(_ string, _ string) error {
					jobCancelled.Close()
					return nil
				}).Times(1)
				tw1.SchemaMgr.EXPECT().InsertJobLog("", worker.ForceKillMsg, gomock.Any()).Return(nil).Times(1)
				tw1.SchemaMgr.EXPECT().InsertJobLog("", worker.ExitStatusMsg(9), gomock.Any()).Return(nil).Times(1)

				tw1.JobsCh <- schema.JobChange{
					OldValue: nil,
					NewValue: job,
				}

				// TODO: use a mock executor and actually wait for execution to start
				time.Sleep(time.Millisecond * 100)

				tw1.CancelJobsCh <- schema.JobChange{
					OldValue: job,
					NewValue: nil,
				}

				// wait for it to cancel the job
				WaitersWait(time.Second, jobCancelled)

				tw1.Worker.Shutdown()

				tw1.waitForFinish()
			})

			// XXX: note that there is no actual killing of a process happening in this test because of mocks
			It("should be able to pick it up and run it, then cancel it using the deadline with SIGKILL", func() {
				InitAndAsyncRunForTest(tw1.Worker)
				job := &schema.Job{
					Status:   schema.JobStatusWaiting,
					Action:   ActionSleep,
					Deadline: "1s",
				}

				ch1 := make(chan struct{})
				executor, ch2 := MockDoubleSyncExecutor(mockCtrl, ch1)
				tw1.Executors = []worker.Executor{executor}

				tw1.SchemaMgr.EXPECT().ClaimJob("", tw1.ID, gomock.Any()).Return(job, nil).Times(1)
				// wait for this function to be called
				jobFinished := NewWaiter("update job status")
				tw1.SchemaMgr.EXPECT().UpdateJobStatus("", schema.JobStatusTimeout).Do(func(_ string, _ string) error {
					jobFinished.Close()
					return nil
				}).Times(1)

				tw1.JobsCh <- schema.JobChange{
					OldValue: nil,
					NewValue: job,
				}

				ChWait("execution started", ch2, time.Second)

				// XXX: unfortunately, we have to wait for the timer to start before we increment time
				// we can probably make this better
				time.Sleep(100 * time.Millisecond)
				// advance the time to trigger the killing of the process
				tw1.Clock.AddTime(time.Second)

				// XXX: we have to wait here for the timeout to take effect before we allow
				// the job to complete. Otherwise it won't be canceled in time
				time.Sleep(100 * time.Millisecond)

				// release the executor to let it finish
				close(ch1)

				// wait for it to cancel the job
				WaitersWait(time.Second, jobFinished)

				tw1.Worker.Shutdown()

				tw1.waitForFinish()
			})

			It("should be able to pick it up and run it, then cancel it using the deadline with SIGTERM", func() {
				InitAndAsyncRunForTest(tw1.Worker)
				job := &schema.Job{
					Status:      schema.JobStatusWaiting,
					Action:      ActionSleep,
					Deadline:    "100ms",
					StopTimeout: "2s",
				}

				tw1.SchemaMgr.EXPECT().ClaimJob("", tw1.ID, gomock.Any()).Return(job, nil).Times(1)
				// wait for this function to be called
				jobFinished := NewWaiter("update job status")
				tw1.SchemaMgr.EXPECT().UpdateJobStatus("", schema.JobStatusTimeout).Do(func(_ string, _ string) error {
					jobFinished.Close()
					return nil
				}).Times(1)
				updateLogs := NewWaiter("update logs")
				tw1.SchemaMgr.EXPECT().InsertJobLog("", worker.ExitStatusMsg(15), gomock.Any()).Do(func(_ string, _ string, _ int) error {
					updateLogs.Close()
					return nil
				})

				tw1.JobsCh <- schema.JobChange{
					OldValue: nil,
					NewValue: job,
				}
				// XXX: We need to make sure the job has started to execute at this time
				// This test can be improved to make it less racy
				time.Sleep(100 * time.Millisecond)
				tw1.Clock.AddTime(100 * time.Millisecond)
				// wait for the job to be killed
				time.Sleep(200 * time.Millisecond)

				// wait for it to cancel the job
				WaitersWait(time.Second, jobFinished, updateLogs)
				//WaitersWait(time.Second, updateLogs)

				tw1.Worker.Shutdown()

				tw1.waitForFinish()
			})

			It("should be able to pick it up and run it, then cancel it by shutting down", func() {
				InitAndAsyncRunForTest(tw1.Worker)
				job := &schema.Job{
					Status: schema.JobStatusWaiting,
					Action: ActionSleep,
				}

				tw1.SchemaMgr.EXPECT().ClaimJob("", tw1.ID, gomock.Any()).Return(job, nil).Times(1)
				tw1.SchemaMgr.EXPECT().UpdateJobStatus("", schema.JobStatusWorkerShutdown).Return(nil).Times(1)
				tw1.SchemaMgr.EXPECT().InsertJobLog("", worker.ForceKillMsg, gomock.Any()).Return(nil).Times(1)
				tw1.SchemaMgr.EXPECT().InsertJobLog("", worker.ExitStatusMsg(9), gomock.Any()).Return(nil).Times(1)

				tw1.JobsCh <- schema.JobChange{
					OldValue: nil,
					NewValue: job,
				}

				// XXX: wait for it to claim the job better
				time.Sleep(100 * time.Millisecond)

				tw1.Worker.Shutdown()

				tw1.waitForFinish()
			})

			It("should be able to pick it up and run it and succeed, with params just because", func() {
				InitAndAsyncRunForTest(tw1.Worker)
				params := map[string]string{"abc": "def"}
				job := &schema.Job{
					Status:     schema.JobStatusWaiting,
					Action:     ActionEcho,
					Parameters: params,
				}

				tw1.SchemaMgr.EXPECT().ClaimJob("", tw1.ID, gomock.Any()).Return(job, nil).Times(1)
				// wait for this function to be called
				jobDone := NewWaiter("update job status")
				tw1.SchemaMgr.EXPECT().UpdateJobStatus("", schema.JobStatusDone).Do(func(_ string, _ string) error {
					jobDone.Close()
					return nil
				}).Times(1)
				tw1.SchemaMgr.EXPECT().InsertJobLog("", "success", 0).Return(nil).Times(1)

				exec1 := MockEchoExecutor(mockCtrl, params)
				tw1.Executors = []worker.Executor{exec1}

				tw1.JobsCh <- schema.JobChange{
					OldValue: nil,
					NewValue: job,
				}

				// wait for it to succeed at running echo
				WaitersWait(time.Second, jobDone)

				// See Note 1
				time.Sleep(10 * time.Millisecond)

				tw1.Worker.Shutdown()

				tw1.waitForFinish()
			})

			It("once a job is failed, make sure the recovery job is created", func() {
				InitAndAsyncRunForTest(tw1.Worker)

				// don't bother actually running the job, just tell a worker that it failed
				job := schema.Job{
					ID:     "job1",
					Status: schema.JobStatusWaiting,
					Action: ActionFail,
				}
				jobFailed := job
				// There are other status codes that represent errors, but we test only this one
				jobFailed.Status = schema.JobStatusError

				// wait for this function to be called
				scheduled := NewWaiter("recovery job scheduled")
				tw1.SchemaMgr.EXPECT().CreateJob(&schema.Job{
					WorkerID:       tw1.ID,
					RecoveryFromID: "job1",
					Status:         schema.JobStatusRunning,
					ScheduledAt:    tw1.Clock.Now(),
					LastUpdated:    tw1.Clock.Now(),
					Action:         ActionEcho,
				}).Do(func(_ *schema.Job) error {
					scheduled.Close()
					return nil
				}).Times(1)
				// wait for this function to be called
				jobDone := NewWaiter("update job status")
				tw1.SchemaMgr.EXPECT().UpdateJobStatus("", schema.JobStatusDone).Do(func(_ string, _ string) error {
					jobDone.Close()
					return nil
				}).Times(1)
				tw1.SchemaMgr.EXPECT().InsertJobLog("", "success", 0).Return(nil).Times(1)

				tw1.RecoverableJobsCh <- schema.JobChange{
					OldValue: &job,
					NewValue: &jobFailed,
				}

				// wait for it to finish scheduling a new job and running it
				WaitersWait(time.Second, scheduled)
				WaitersWait(time.Second, jobDone)

				tw1.Worker.Shutdown()

				tw1.waitForFinish()
			})

		})
	})
})

func MockEchoExecutor(mockCtrl *gomock.Controller, params map[string]string) *worker.MockExecutor {
	if params == nil {
		params = map[string]string{}
	}
	executor := worker.NewMockExecutor(mockCtrl)
	executor.EXPECT().Start(gomock.Any(), "echo", []string{"success"}, params, time.Duration(0)).Return(strings.NewReader("success\n"), nil)
	executor.EXPECT().Wait().Return(nil)
	return executor
}

func MockSyncExecutor(mockCtrl *gomock.Controller) (*worker.MockExecutor, <-chan struct{}) {
	ch := make(chan struct{})
	executor := worker.NewMockExecutor(mockCtrl)
	executor.EXPECT().Start(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), time.Duration(0)).Do(func(ctx context.Context, command string, args []string, parameters map[string]string, stopTimeout time.Duration) (io.Reader, error) {
		close(ch)
		return strings.NewReader("success\n"), nil
	})
	executor.EXPECT().Wait().Return(nil)
	return executor, ch
}

func MockDoubleSyncExecutor(mockCtrl *gomock.Controller, ch1 chan struct{}) (*worker.MockExecutor, <-chan struct{}) {
	ch2 := make(chan struct{})
	executor := worker.NewMockExecutor(mockCtrl)
	reader, writer := io.Pipe()
	executor.EXPECT().Start(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), time.Duration(0)).Do(func(ctx context.Context, command string, args []string, parameters map[string]string, stopTimeout time.Duration) (io.Reader, error) {
		close(ch2)
		return reader, nil
	})
	statusLock := sync.Mutex{}
	status := schema.JobStatusDone
	executor.EXPECT().Wait().Do(func() error {
		ChWait("unblock execute wait", ch1, time.Second)
		writer.Close()
		statusLock.Lock()
		defer statusLock.Unlock()
		if status == schema.JobStatusDone {
			return nil
		} else {
			return &worker.ExecutorError{
				Status: status,
				Err:    fmt.Errorf("blah"),
			}
		}
	})
	executor.EXPECT().Cancel(gomock.Any()).Do(func(s string) {
		statusLock.Lock()
		defer statusLock.Unlock()
		status = s
	}).AnyTimes()
	return executor, ch2
}

func InitAndAsyncRunForTest(w worker.Worker) {
	err := w.Init()
	Expect(err).To(BeNil())

	go func() {
		defer GinkgoRecover()
		err := w.MonitorCronsGoroutine()
		Expect(err).To(BeNil())
	}()

	go func() {
		defer GinkgoRecover()
		err := w.MonitorCancellationsGoroutine()
		Expect(err).To(BeNil())
	}()

	go func() {
		defer GinkgoRecover()
		err := w.MonitorRecoverableJobsGoroutine()
		Expect(err).To(BeNil())
	}()

	go func() {
		defer GinkgoRecover()
		err := w.ClaimJobsGoroutine()
		Expect(err).To(BeNil())
	}()

	go func() {
		defer GinkgoRecover()
		err := w.DetectExpirationGoroutine()
		Expect(err).To(BeNil())
	}()
}

func WaitersWait(timeout time.Duration, waiters ...*Waiter) {
	for _, waiter := range waiters {
		err := waiter.Wait(timeout)
		Expect(err).To(BeNil())
	}
}
