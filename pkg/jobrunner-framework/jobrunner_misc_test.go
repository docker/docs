package jobrunner

import (
	"bufio"
	"io"
	"time"

	"golang.org/x/net/context"

	"github.com/docker/dhe-deploy/pkg/jobrunner-framework/constants"
	"github.com/docker/dhe-deploy/pkg/jobrunner-framework/schema"
	"github.com/docker/dhe-deploy/pkg/jobrunner-framework/worker"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Misc", func() {
	Describe("time heap", func() {
		// TODO: way more tests
		It("test heap basics", func() {
			expirationHeap := worker.NewSafeTimeHeap()
			jobs := []*schema.Job{
				{ID: "first", HeartbeatExpiration: time.Now().Add(-constants.DefaultHeartbeatTimeout)},
				{ID: "second", HeartbeatExpiration: time.Now().Add(-constants.DefaultHeartbeatTimeout)},
				{ID: "third", HeartbeatExpiration: time.Now().Add(-constants.DefaultHeartbeatTimeout)},
				{ID: "notexpired", HeartbeatExpiration: time.Now().Add(time.Minute)},
			}

			for _, job := range jobs {
				expirationHeap.Update(&schema.JobChange{nil, job})
			}

			for i, job := range jobs[:3] {
				Expect(expirationHeap.Len()).To(Equal(4 - i))
				Expect(expirationHeap.Peek().ID).To(Equal(job.ID))
				Expect(expirationHeap.PopIfExpired(time.Now()).ID).To(Equal(job.ID))
			}

			Expect(expirationHeap.Len()).To(Equal(1))
			Expect(expirationHeap.Peek().ID).To(Equal(jobs[3].ID))
			Expect(expirationHeap.PopIfExpired(time.Now())).To(BeNil())
			Expect(expirationHeap.Len()).To(Equal(1))
		})
	})
	Describe("executor", func() {
		var executor worker.Executor

		BeforeEach(func() {
			executor = worker.DefaultNewExecutor()
		})

		It("executes echo", func() {
			reader, err := executor.Start(context.Background(), "echo", []string{"one", "two"}, map[string]string{}, 0)
			Expect(err).To(BeNil())

			ch := asyncWaitAndExpectStatus(executor, nil, schema.JobStatusDone)

			lines := readAll(reader)
			Expect(lines).To(Equal([]string{"one two"}))
			ChWait("async wait and expect status", ch, time.Second)
		})

		It("can cancel before starting", func() {
			executor.Cancel(schema.JobStatusWorkerDead)

			_, err := executor.Start(context.Background(), "echo", []string{"one", "two"}, map[string]string{}, 0)
			Expect(err).To(Not(BeNil()))
			ee := err.(*worker.ExecutorError)
			Expect(ee.Status).To(Equal(schema.JobStatusWorkerDead))
		})

		It("can cancel after starting (with SIGKILL)", func() {
			reader, err := executor.Start(context.Background(), "sleep", []string{"10"}, map[string]string{}, 0)
			Expect(err).To(BeNil())

			executor.Cancel(schema.JobStatusCanceled)

			ch := asyncWaitAndExpectStatus(executor, []string{worker.ForceKillMsg, worker.ExitStatusMsg(9)}, schema.JobStatusCanceled)

			lines := readAll(reader)
			Expect(lines).To(Equal([]string{}))
			ChWait("async wait and expect status", ch, time.Second)
		})

		It("can cancel after running for a bit (with SIGKILL)", func() {
			reader, err := executor.Start(context.Background(), "sleep", []string{"10"}, map[string]string{}, 0)
			Expect(err).To(BeNil())

			time.AfterFunc(time.Millisecond*100, func() {
				executor.Cancel(schema.JobStatusCanceled)
			})

			ch := asyncWaitAndExpectStatus(executor, []string{worker.ForceKillMsg, worker.ExitStatusMsg(9)}, schema.JobStatusCanceled)

			lines := readAll(reader)
			Expect(lines).To(Equal([]string{}))
			ChWait("async wait and expect status", ch, time.Second)
		})

		It("can cancel with a SIGTERM", func() {
			reader, err := executor.Start(context.Background(), "sleep", []string{"10"}, map[string]string{}, time.Second)
			Expect(err).To(BeNil())

			time.AfterFunc(time.Millisecond*100, func() {
				executor.Cancel(schema.JobStatusCanceled)
			})

			ch := asyncWaitAndExpectStatus(executor, []string{worker.ExitStatusMsg(15)}, schema.JobStatusCanceled)

			lines := readAll(reader)
			Expect(lines).To(Equal([]string{}))
			ChWait("async wait and expect status", ch, time.Second)
		})

		It("can cancel before calling wait", func() {
			reader, err := executor.Start(context.Background(), "sleep", []string{"10"}, map[string]string{}, time.Second)
			Expect(err).To(BeNil())

			time.AfterFunc(time.Millisecond*10, func() {
				executor.Cancel(schema.JobStatusCanceled)
			})

			// wait for cancel before launching the waiter
			time.Sleep(time.Millisecond * 100)

			// XXX: in this case wait gets called after cancel, so there's no chance for us to get the exit status of the process
			ch := asyncWaitAndExpectStatus(executor, []string{worker.ExitStatusMsg(15)}, schema.JobStatusCanceled)
			//ch := asyncWaitAndExpectStatus(executor, []string{}, schema.JobStatusCanceled)

			lines := readAll(reader)
			Expect(lines).To(Equal([]string{}))
			ChWait("async wait and expect status", ch, time.Second)
		})
	})
})

func asyncWaitAndExpectStatus(executor worker.Executor, expectExtraLogs []string, expectStatus string) chan struct{} {
	ch := make(chan struct{})
	go func() {
		defer GinkgoRecover()
		err := executor.Wait()
		if expectStatus == schema.JobStatusDone {
			Expect(err).To(BeNil())
		} else {
			Expect(err).To(Not(BeNil()))
			ee := err.(*worker.ExecutorError)
			Expect(ee.ExtraLogs).To(Equal(expectExtraLogs))
			Expect(ee.Status).To(Equal(expectStatus))
		}
		ch <- struct{}{}
	}()
	return ch
}
func readAll(reader io.Reader) []string {
	out := []string{}
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		out = append(out, scanner.Text())
	}
	return out
}
