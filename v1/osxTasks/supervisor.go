/*
 * Copyright (C) 2016 Docker Inc
 * All rights reserved
 */

package osxTasks

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"os"
	"path"
	"strings"
	"sync"
	"syscall"
	"time"
)

var mutex sync.Mutex
var cond = sync.NewCond(&mutex)

var shuttingDown bool

type child struct {
	task Task
	proc *os.Process
	quit *bool // updates signalled by cond
}

var children []child // protected by mutex

// lock and unlock do not account for Cond locking/unlocking

func lock(descr string) {
	mutex.Lock()
	//logrus.Printf("locked mutex by %s", descr)
}

func unlock(descr string) {
	mutex.Unlock()
	//logrus.Printf("unlocked mutex by %s", descr)
}

// Supervise keeps the listed tasks running until signaled to stop via the
// quit channel.
func Supervise(tasks []Task, quit chan struct{}) {
	StartAll(tasks, false)

	for restartOne(tasks, quit) {
	}
}

// restartOne performs one round of process restarts
func restartOne(tasks []Task, quit chan struct{}) bool {
	work := make(chan struct{})
	go waitAny(work)

	select {
	case <-quit:
		// orderly shutdown requested
		lock("restartOne quit")
		defer unlock("restartOne quit")
		shuttingDown = true
		shutdownOrdered(children)
		quit <- struct{}{}
		return false
	case <-work:
		// one or more processes have quit
		// Sleep for 1s to allow other dependent processes to fail. Otherwise
		// we might only restart fewer processes than we need, leading to more
		// failures.
		time.Sleep(time.Second)
		// Discover the greatest shutdown_order of any of the failed processes.
		lock("restartOne work")
		maxDeadLevel := uint(0)
		for _, c := range children {
			if *c.quit && c.task.ShutdownOrder > maxDeadLevel {
				maxDeadLevel = c.task.ShutdownOrder
			}
		}
		// Shutdown all processes of this order and lower
		toShutdown := []child{}
		for _, c := range children {
			if c.task.ShutdownOrder <= maxDeadLevel {
				toShutdown = append(toShutdown, c)
			}
		}
		// Restart all tasks of this order and lower
		shutdownOrdered(toShutdown)
		tasksNeeded := []Task{}
		for _, c := range toShutdown {
			tasksNeeded = append(tasksNeeded, c.task)
		}
		unlock("restartOne work")

		StartAll(tasksNeeded, false)
		return true
	}
}

// waitAny sends on the result channel when any of the children have quit
func waitAny(result chan struct{}) {
	lock("waitAny")
	defer unlock("waitAny")
	for {
		if shuttingDown {
			return
		}
		for _, c := range children {
			if *c.quit {
				result <- struct{}{}
				return
			}
		}
		cond.Wait()
	}
}

// StartAll starts processes for the given list of tasks
func StartAll(tasks []Task, debug bool) {
	lock("StartAll")
	defer unlock("StartAll")
	logrus.Printf("Starting %s", taskNames(tasks))

	for _, t := range tasks {
		files := []*os.File{os.Stdin, os.Stdout, os.Stderr}
		args := make([]string, len(t.Arguments))
		copy(args, t.Arguments)
		fd := 3

		for _, s := range t.sockets {
			files = append(files, s.File())
			args = s.Replace(args, fd)
			fd++
		}
		if debug && t.DebugArg != "" {
			args = append(args, t.DebugArg)
		}
		sys := syscall.SysProcAttr{
			Setsid: true,
		}
		attr := os.ProcAttr{
			Files: files,
			Sys:   &sys,
		}
		args = append([]string{t.Argv0}, args...)
		proc, err := os.StartProcess(t.Command, args, &attr)
		if err != nil {
			// shutdown everything already started
			logrus.Fatalf("Failed to start %s %s: %#v", t.Command, strings.Join(args, " "), err)
		}
		logrus.Printf("Start %s (pid %d)", path.Base(t.Command), proc.Pid)
		WritePidFile(t, proc.Pid)

		quit := false
		go func(t Task) {
			state, err := proc.Wait()
			if err != nil {
				logrus.Printf("Error waiting for pid %d task %s to exit: %#v", proc.Pid, t.Command, err)
				return
			}
			logrus.Printf("Reap %s (pid %d): %s", path.Base(t.Command), proc.Pid, state.String())
			removePidFile(t)
			lock("reap")
			defer unlock("reap")
			quit = true
			cond.Broadcast()
		}(t)
		addChild(t, proc, &quit)
	}
}

// addChild adds a child to the global list
func addChild(task Task, proc *os.Process, quit *bool) {
	children = append(children, child{task: task, proc: proc, quit: quit})
}

// delChild removes a child from the global list
func delChild(c2 child) {
	children2 := []child{}
	for _, c := range children {
		if c.proc.Pid != c2.proc.Pid {
			children2 = append(children2, c)
		}
	}
	children = make([]child, len(children2))
	copy(children, children2)
}

func kill(child child, signal syscall.Signal) {
	logrus.Printf("Signal %s to %s (pid %d)", signal.String(), path.Base(child.task.Argv0), child.proc.Pid)
	syscall.Kill(child.proc.Pid, signal)
}

// shutdownConcurrent performs a concurrent shutdown of the given children, by first sending
// SIGTERM and eventually sending SIGKILL if they fail to respond.
func shutdownConcurrent(these []child) {
	// remove from the global list
	for _, c := range these {
		delChild(c)
	}
	// send them all a SIGTERM to encourage a graceful shutdown
	for _, c := range these {
		if !*c.quit {
			kill(c, syscall.SIGTERM)
		}
	}
	// wait for up to 60s for them to exit
	timer := time.NewTimer(60 * time.Second)
Loop:
	for {
		remaining := []child{}
		for _, c := range these {
			if !*c.quit {
				remaining = append(remaining, c)
			}
		}
		these := make([]child, len(remaining))
		copy(these, remaining)

		if len(these) == 0 {
			timer.Stop()
			break Loop
		}
		unlock("shutdownConcurrent loop")
		select {
		case <-timer.C:
			logrus.Warn("Timed-out waiting for graceful shutdown")
			lock("shutdownConcurrent timeout")
			break Loop
		default:
			time.Sleep(100 * time.Millisecond)
			lock("shutdownConcurrent loop")
		}
	}
	// any children which haven't shutdown should be sent a SIGKILL
	for _, c := range these {
		if !*c.quit {
			logrus.Warnf("Sending a SIGKILL to %s", c.task.Command)
			kill(c, syscall.SIGKILL)
		}
	}
	for _, c := range these {
		for !*c.quit {
			cond.Wait()
		}
	}
}

// shutdownOrdered gracefully shuts down all running processes, where the set of
// processes with the same shutdown_order are shutdown concurrently.
func shutdownOrdered(these []child) {
	minLevel := ^uint(0)
	maxLevel := uint(0)
	for _, c := range these {
		if c.task.ShutdownOrder > maxLevel {
			maxLevel = c.task.ShutdownOrder
		}
		if c.task.ShutdownOrder < minLevel {
			minLevel = c.task.ShutdownOrder
		}
	}
	for level := minLevel; level <= maxLevel; level++ {
		var toShutdown []child
		for _, c := range these {
			if c.task.ShutdownOrder == level {
				toShutdown = append(toShutdown, c)
			}
		}
		logrus.Printf("Stop %d children with order %d: %s", len(toShutdown), level, childNames(toShutdown))
		shutdownConcurrent(toShutdown)
	}
}

// StopAll gracefully shuts down all running processes.
func StopAll() {
	lock("StopAll")
	defer unlock("StopAll")
	shuttingDown = true
	shutdownOrdered(children)
}

func taskNames(tasks []Task) string {
	names := []string{}
	for _, t := range tasks {
		names = append(names, path.Base(t.Command))
	}
	return strings.Join(names, ", ")
}

func childNames(children []child) string {
	names := []string{}
	for _, c := range children {
		names = append(names, fmt.Sprintf("%s (pid %d)", path.Base(c.task.Argv0), c.proc.Pid))
	}
	return strings.Join(names, ", ")
}
