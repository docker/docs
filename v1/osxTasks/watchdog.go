/*
 * Copyright (C) 2016 Docker Inc
 * All rights reserved
 */

package osxTasks

import (
	"github.com/Sirupsen/logrus"
	"io"
	"os"
	"strconv"
	"syscall"
)

var expectedHandshake = "This is not a stable interface. Do not try to use it, it will change."

// StartWatchdog initialises the watchdog connection over a pipe to the parent
// process. Initially we read and verify the handshake message. The expected
// message is derived from the source repo and changes with every release,
// to avoid people relying on this as a stable interface. Next we block reading
// from the PIPE and when the parent goes away we get a read of 0.
func StartWatchdog(w string, c chan os.Signal) {
	if w == "" {
		logrus.Fatalln("Failed to initialise watchdog protocol")
	}
	if len(w) < 3 || w[0:3] != "fd:" {
		logrus.Fatalf("Failed to parse watchdog fd: %s", w)
	}
	w = w[3:]
	watchdogFd, err := strconv.Atoi(w)
	if err != nil {
		logrus.Fatalf("Failed to parse watchdog fd: %#v", err)
	}
	r := os.NewFile(uintptr(watchdogFd), "watchdog")
	bs := make([]byte, len(expectedHandshake))
	_, err = io.ReadAtLeast(r, bs, len(bs))
	if err != nil {
		logrus.Fatalln("Failed to read watchdog handshake")
	}
	if string(bs) != expectedHandshake {
		logrus.Fatalln("Watchdog handshake failed")
	}
	go func() {
		for {
			n, _ := r.Read(bs)
			if n == 0 {
				logrus.Println("EOF on watchdog pipe: shutting down")
				c <- syscall.SIGTERM
				return
			}
		}
	}()
}
