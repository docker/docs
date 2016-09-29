/*
 * Copyright (C) 2016 Docker Inc
 * All rights reserved
 */

package osxTasks

import (
	"encoding/json"
	"github.com/Sirupsen/logrus"
	appleutil "github.com/docker/pinata/v1/apple/util"
	"github.com/mitchellh/go-ps"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

// TODO(djs55): merge this with "type child struct"
type process struct {
	Task Task
	Pid  int
}

// WritePidFile guarantees to kill the process should it somehow survive a
// full app shutdown.
func WritePidFile(task Task, pid int) {
	process := process{task, pid}

	filename := getPidFilename(task)
	js, err := json.Marshal(process)
	if err != nil {
		logrus.Fatalln("Failed to marshal json", err)
	}
	err = ioutil.WriteFile(filename, js, 0644)
	if err != nil {
		logrus.Fatalln("Failed to write", filename, err)
	}
}

func removePidFile(t Task) {
	filename := getPidFilename(t)
	err := os.Remove(filename)
	if err != nil && !(os.IsNotExist(err)) {
		logrus.Fatalln("Failed to remove pidfile", filename, err)
	}
}

func getPidFilename(t Task) string {
	dir := getPidDirectory()

	return filepath.Join(dir, t.Argv0)
}

// ShutdownLeakedProcesses will shutdown any leaked processes
func ShutdownLeakedProcesses() {
	procs := getProcesses()

	minLevel := ^uint(0)
	maxLevel := uint(0)
	for _, p := range procs {
		if p.Task.ShutdownOrder > maxLevel {
			maxLevel = p.Task.ShutdownOrder
		}
		if p.Task.ShutdownOrder < minLevel {
			minLevel = p.Task.ShutdownOrder
		}
	}
	for level := minLevel; level <= maxLevel; level++ {
		var toShutdown []process
		for _, p := range procs {
			if p.Task.ShutdownOrder == level {
				toShutdown = append(toShutdown, p)
			}
		}
		shutdownLeakedConcurrent(toShutdown)
	}
}

func shutdownLeakedConcurrent(all []process) {
	// Send a SIGTERM
	signal := syscall.SIGTERM
	for _, p := range all {
		if !isStillRunning(p) {
			logrus.Println("Removing leaked pidfile", p.Task.Argv0)
			removePidFile(p.Task)
			continue
		}
		logrus.Printf("Signal %s to %s (pid %d)", signal.String(), p.Task.Argv0, p.Pid)
		syscall.Kill(p.Pid, signal)
	}
	// Wait for up to 30s
	retries := 30
	finished := false
	for !finished && retries > 0 {
		retries = retries - 1
		finished = true
		for _, p := range all {
			if isStillRunning(p) {
				finished = false
			}
		}
		if !finished {
			time.Sleep(time.Second)
		}
	}
	// Send a SIGKILL
	signal = syscall.SIGKILL
	for _, p := range all {
		if isStillRunning(p) {
			logrus.Printf("Signal %s to %s (pid %d)", signal.String(), p.Task.Argv0, p.Pid)
			syscall.Kill(p.Pid, signal)
		}
		removePidFile(p.Task)
	}
}

func getProcesses() []process {
	dir := getPidDirectory()
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		logrus.Fatalln("readdir", dir, err)
	}
	var all []process
	for _, fi := range files {
		pidFile := filepath.Join(dir, fi.Name())
		contents, err := ioutil.ReadFile(pidFile)
		if err != nil {
			logrus.Println("Failed to read contents of pidfile", fi.Name(), err)
			continue
		}
		var process process
		err = json.Unmarshal(contents, &process)
		if err != nil {
			logrus.Println("Failed to read contents of pidfile", pidFile, err)
			continue
		}
		all = append(all, process)
	}
	return all
}

func isStillRunning(p process) bool {
	proc, err := ps.FindProcess(p.Pid)
	// if the pid still refers to an executable with the same name, then it's still
	// running.
	if err != nil {
		logrus.Fatalln("FindProcess failed for pid", p.Pid, err)
	}
	// Executable() might be truncated:
	// https://github.com/mitchellh/go-ps/issues/15
	// workaround by looking for com.docker. prefix
	if proc != nil && strings.HasPrefix(proc.Executable(), "com.docker.") {
		return true
	}
	return false
}

func getPidDirectory() string {
	dir := filepath.Join(appleutil.GetContainerPath(), "tasks")
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		logrus.Fatalln("mkdir", dir, err)
	}
	return dir
}
