/*
 * Copyright (C) 2016 Docker Inc
 * All rights reserved
 */

package osxTasks

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"net"
	"os"
	"strings"
)

// Task encapsulates an OSX process
type Task struct {
	Command       string
	Argv0         string
	Arguments     []string
	ShutdownOrder uint
	sockets       []ListeningSocket
	DebugArg      string // argument to append in debug mode
}

// ListeningSocket represents a listening unix soxket
type ListeningSocket struct {
	file  *os.File
	index int
}

// String returns a string represetation of a ListeningSocket
func (l *ListeningSocket) String() string {
	return fmt.Sprintf("{->%d<-}", l.index)
}

// Replace when provided with args, returns these args with any occurances of the
// ListeningSocket.String replaced with the provided fd
func (l *ListeningSocket) Replace(args []string, fd int) []string {
	var replaced []string
	for _, a := range args {
		replaced = append(replaced, strings.Replace(a, l.String(), fmt.Sprintf("%d", fd), -1))
	}
	return replaced
}

// File returns the file of a ListeningSocket
func (l *ListeningSocket) File() *os.File {
	return l.file
}

var nextListeningSocket int

// ListenUnix returns a new ListeningSocket
func ListenUnix(path string) ListeningSocket {
	err := os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		logrus.Fatalf("Failed to remove existing socket on %s: %s", path, err.Error())
	}
	l, err := net.ListenUnix("unix", &net.UnixAddr{Name: path, Net: "unix"})
	if err != nil {
		logrus.Fatalf("Failed to listen on %s: %#v", path, err.Error())
	}
	file, err := l.File()
	if err != nil {
		logrus.Fatalf("Failed to find fd of Unix socket listener: %s", err.Error())
	}
	index := nextListeningSocket
	nextListeningSocket++
	return ListeningSocket{file, index}
}

// NewTask creates a new Task
func NewTask(Command string, Argv0 string, Arguments []string, ShutdownOrder uint, sockets []ListeningSocket, DebugArg string) Task {
	return Task{Command, Argv0, Arguments, ShutdownOrder, sockets, DebugArg}
}
