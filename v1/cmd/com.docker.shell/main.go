/*
 * Copyright (C) 2016 Docker Inc
 * All rights reserved
 */

package main

import (
	"flag"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/chzyer/readline"
	datakit "github.com/docker/datakit/api/go-datakit"
	"github.com/docker/pinata/v1/apple"
	"github.com/docker/pinata/v1/osxTasks"
	"github.com/docker/pinata/v1/pinataSockets"
	"github.com/docker/pinata/v1/reportError"
	"golang.org/x/net/context"
	"io"
	"strings"
	"time"
)

func main() {

	logrus.AddHook(apple.NewLogrusASLHook())
	reportError.Initialize()

	osxTasks.AcquireTaskManagerLock()

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

	tasks := osxTasks.ListTasks(bundle)

	osxTasks.IncreaseFdLimit()
	osxTasks.ShutdownLeakedProcesses()

	shell(tasks)
}

func shell(tasks []osxTasks.Task) {
	completer := readline.NewPrefixCompleter(
		readline.PcItem("start"),
		readline.PcItem("stop"),
	)

	rl, err := readline.NewEx(&readline.Config{
		HistoryFile:  ".history",
		AutoComplete: completer,
	})
	if err != nil {
		logrus.Fatalln(err)
	}
	running := false
	for {
		msg := "stopped"
		if running {
			msg = "running"
		}
		rl.SetPrompt(fmt.Sprintf("%s 🐳 > ", msg))

		line, err := rl.Readline()
		if err == io.EOF {
			logrus.Printf("Quitting (EOF)")
			osxTasks.StopAll()
			return
		} else if err != nil {
			logrus.Println("error: ", err)
			osxTasks.StopAll()
			logrus.Fatalln("error: ", err)
		}

		if line == "" {
			continue
		}
		line = strings.Trim(line, " \r\n\t")
		words := strings.Split(line, " ")
		if len(words) == 0 {
			logrus.Println("Please restate the question. You have twenty seconds to comply.")
			continue
		}
		switch words[0] {
		case "start":
			logrus.Printf("Starting tasks")
			osxTasks.StartAll(tasks, true)
			running = true
		case "stop":
			logrus.Printf("Stopping tasks")
			osxTasks.StopAll()
			running = false
		case "memory":
			if len(words) != 2 {
				logrus.Println("Usage: memory <GiB>")
				continue
			}
			setMemory(words[1])
		case "exit":
			logrus.Printf("Quitting")
			osxTasks.StopAll()
			return
		}
	}
}

func setMemory(memory string) {
	logrus.Printf("Setting memory to %s GiB", memory)
	ctx := context.TODO()
	client := newClient(ctx)
	defer client.Close(ctx)
	err := osxTasks.SetMemory(ctx, client, memory)
	if err != nil {
		logrus.Printf("Failed")
	} else {
		logrus.Printf("OK")
	}
}

func newClient(ctx context.Context) *datakit.Client {
	// share one database connection
	var client *datakit.Client
	var err error
	attemptsRemaining := 100
	for attemptsRemaining > 0 {

		client, err = datakit.Dial(ctx, "unix", pinataSockets.GetDBSocketPath())
		if err == nil {
			break
		}
		logrus.Printf("Failed to connect to db: %#v\n", err)
		time.Sleep(100 * time.Millisecond)
		attemptsRemaining = attemptsRemaining - 1
	}
	if client == nil {
		logrus.Fatalln("Failed to connect to the database after 10s")
	}
	return client
}
