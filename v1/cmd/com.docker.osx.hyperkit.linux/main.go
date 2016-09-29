/*
 * Copyright (C) 2016 Docker Inc
 * All rights reserved
 */

package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/reexec"
	"github.com/docker/pinata/v1/apple"
	"github.com/docker/pinata/v1/osxTasks"
	"github.com/docker/pinata/v1/reportError"
)

func main() {

	logrus.AddHook(apple.NewLogrusASLHook())
	reportError.Initialize()
	reexec.Register(osxTasks.APIServerCommandName, osxTasks.APIServer)
	if reexec.Init() {
		return
	}

	// NOTE(aduermael) logrus should be used like this for logs to use ASL API
	// logrus.Debugln("A debug message")
	// logrus.Println("A info message")
	// logrus.Warnln("A warning message")
	// logrus.Errorln("An error message")
	// logrus.Fatalln("A fatal error message")

	osxTasks.AcquireTaskManagerLock()

	// NOTE(aduermael): we don't check for pinataSockets.GetVsockDirPath() existence
	// here. We could do it, but in any case we should not create it manually.
	// This directory is the app's container directory, automatically created somehow
	// by the system for our application. If we follow sandboxing guidelines, this
	// directory is the only one available to read and write data without asking
	// for user permissions. So I guess it may be created with some metadata and we
	// shoulnd't get in the way.

	osxTasks.IncreaseFdLimit()

	bundle := ""
	debug := false
	watchdog := ""

	flag.StringVar(&bundle, "bundle", "", "path to the root of the application bundle (\"\" means autodetect)")
	flag.StringVar(&watchdog, "watchdog", "", "fd of watchdog pipe")
	flag.BoolVar(&debug, "debug", false, "run in debug mode")
	flag.Parse()

	if bundle == "" {
		bundle = apple.FindBundle()
	}

	osxTasks.ShutdownLeakedProcesses()

	tasks := osxTasks.ListTasks(bundle)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGINT)

	osxTasks.StartWatchdog(watchdog, c)

	quit := make(chan struct{})
	go osxTasks.Supervise(tasks, quit)

	<-c
	logrus.Println("Received SIGTERM, shutting down")
	quit <- struct{}{}
	<-quit
}
