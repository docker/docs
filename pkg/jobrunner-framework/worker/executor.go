package worker

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/docker/dhe-deploy/pkg/jobrunner-framework/schema"
	"github.com/docker/distribution/context"
)

var logPrefix = "JOBRUNNER"
var ForceKillMsg = fmt.Sprintf("%s: process forcefully killed", logPrefix)
var EarlyCancelMsg = fmt.Sprintf("%s: job canceled before starting", logPrefix)

func ExitStatusMsg(status int) string {
	return fmt.Sprintf("%s: job exited with exit status: %d", logPrefix, status)
}

type executor struct {
	cmd          *exec.Cmd
	ctx          context.Context
	stopTimeout  time.Duration
	started      bool
	cancelStatus string
	cancelLock   sync.Mutex
	logLock      sync.Mutex
	extraLogs    []string
	logsWriter   io.WriteCloser
}

// An executor takes care of running and killing a specific instance of a job by running a binary
// This should be used in the following way:
// 1. create the executor
// 2. call Start() to launch a job
// 3. check for error and use the returned reader to read the combined stdout/stderr output
// 4. in another goroutine call wait to wait for the execution to end
// 5. call cancel at any time before, during or after execution and it will cause wait to unblock with error or start to fail
// Errors returned by the executor will always have the type ExecutorError
type Executor interface {
	Start(ctx context.Context, command string, args []string, parameters map[string]string, stopTimeout time.Duration) (io.Reader, error)
	Wait() error
	Cancel(status string)
}

func (e *executor) log(line string) {
	e.logLock.Lock()
	defer e.logLock.Unlock()
	e.extraLogs = append(e.extraLogs, line)
}

type ExecutorError struct {
	ExtraLogs []string
	Status    string
	Err       error
}

func (e *ExecutorError) Error() string {
	return fmt.Sprintf("Executor error: Return status: %s Extra logs: %s, err: %s", e.Status, strings.Join(e.ExtraLogs, ";"), e.Error())
}

func (e *executor) Err(status string, err error) error {
	e.logLock.Lock()
	defer e.logLock.Unlock()
	return &ExecutorError{
		ExtraLogs: e.extraLogs,
		Status:    status,
		Err:       err,
	}
}

func (e *executor) Start(ctx context.Context, command string, args []string, parameters map[string]string, stopTimeout time.Duration) (io.Reader, error) {
	e.cancelLock.Lock()
	defer e.cancelLock.Unlock()
	if e.cancelStatus != "" {
		return nil, e.Err(e.cancelStatus, fmt.Errorf("job canceled before starting"))
	}
	e.started = true

	e.ctx = ctx
	e.stopTimeout = stopTimeout
	// output of the binary is captured into here
	var logsReader io.Reader
	logsReader, e.logsWriter = io.Pipe()

	e.cmd = exec.Command(command, args...)
	// TODO: should we sanitize anything?
	e.cmd.Env = os.Environ()
	for k, v := range parameters {
		env := fmt.Sprintf("PARAM_%s=%s", strings.ToUpper(k), v)
		e.cmd.Env = append(e.cmd.Env, env)
	}
	e.cmd.Stdout = e.logsWriter
	e.cmd.Stderr = e.logsWriter

	err := e.cmd.Start()
	if err != nil {
		return nil, e.Err(schema.JobStatusError, err)
	} else {
		return logsReader, nil
	}
}

func (e *executor) Wait() error {
	defer func() {
		err := e.logsWriter.Close()
		if err != nil {
			context.GetLogger(e.ctx).Errorf("job failed to close logs pipe: %s", err)
		}
	}()
	err := e.cmd.Wait()
	// parse exit status
	// https://stackoverflow.com/questions/10385551/get-exit-code-go
	if exiterr, ok := err.(*exec.ExitError); ok {
		// This works on both Unix and Windows. Although package
		// syscall is generally platform dependent, WaitStatus is
		// defined for both Unix and Windows and in both cases has
		// an ExitStatus() method with the same signature.
		if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
			e.log(ExitStatusMsg(int(status)))
		}
	}

	e.cancelLock.Lock()
	if e.cancelStatus != "" {
		if err != nil {
			return e.Err(e.cancelStatus, err)
		} else {
			return e.Err(e.cancelStatus, fmt.Errorf("job canceled"))
		}
	}
	e.cancelLock.Unlock()

	if err != nil {
		return e.Err(schema.JobStatusError, err)
	} else {
		return nil
	}
}

// Cancel tries to stop the process, giving it stopTimeout time to stop gracefully
func (e *executor) Cancel(status string) {
	e.cancelLock.Lock()
	defer e.cancelLock.Unlock()
	e.cancelStatus = status
	if !e.started {
		e.log(EarlyCancelMsg)
		return
	}

	if e.cmd != nil && e.cmd.Process != nil {
		if e.stopTimeout > 0 {
			err := e.cmd.Process.Signal(syscall.SIGTERM)
			if err != nil {
				context.GetLogger(e.ctx).Errorf("failed to send SIGTERM: %s", err)
			}

			//timer := time.AfterFunc(e.stopTimeout, func() {
			time.AfterFunc(e.stopTimeout, func() {
				// XXX: we don't have a good way to cancel this timer, we should come up with
				// one eventually because pid reuse can technically cause it to kill
				// the wrong job
				// We can't use Wait inside Cancel because calling Wait() twice on a process is
				// bad and leads to unexpected behaviour
				if e.cmd != nil && e.cmd.Process != nil {
					// by putting this log line before Kill() we guarantee that it will be logged first before the exit code
					e.log(ForceKillMsg)
					err = e.cmd.Process.Kill()
					if err != nil {
						context.GetLogger(e.ctx).Errorf("failed force kill: %s", err)
					}
				}
			})
		} else {
			// by putting this log line before Kill() we guarantee that it will be logged first before the exit code
			e.log(ForceKillMsg)
			e.cmd.Process.Kill()
		}
	} else {
		context.GetLogger(e.ctx).Errorf("Tried to cancel job after start was called, but before the process was created. This should be impossible.")
	}
}

type NewExecutor func() Executor

func DefaultNewExecutor() Executor {
	return &executor{}
}
