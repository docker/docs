package main

import (
	"bufio"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/docker/pinata/v1/pinataSockets"
	"github.com/jeromer/syslogparser"
	"github.com/jeromer/syslogparser/rfc3164"
	"gopkg.in/natefinch/lumberjack.v2"
	"net"
	"os"
	"strconv"
	"time"
)

const (
	//VsyslogPort is the port used for syslog over vsock
	VsyslogPort = 514
)

// https://www.ietf.org/rfc/rfc3164.txt 4.1.1
func facilityString(f interface{}) string {
	v, ok := f.(int)
	if !ok {
		return "???"
	}

	fmap := [...]string{
		/* 0 */ "kern",
		/* 1 */ "user",
		/* 2 */ "mail",
		/* 3 */ "daemon",
		/* 4 */ "auth",
		/* 5 */ "internal",
		/* 6 */ "lpr",
		/* 7 */ "news",
		/* 8 */ "uucp",
		/* 9 */ "clock",
		/* 10 */ "auth",
		/* 11 */ "ftp",
		/* 12 */ "ntp",
		/* 13 */ "audit",
		/* 14 */ "alert",
		/* 15 */ "clock",
		/* 16 */ "local0",
		/* 17 */ "local1",
		/* 18 */ "local2",
		/* 19 */ "local3",
		/* 20 */ "local4",
		/* 21 */ "local5",
		/* 22 */ "local6",
		/* 23 */ "local7",
	}

	if v > len(fmap) {
		return fmt.Sprintf("%d", v)
	}
	return fmap[v]

}

func severityString(s interface{}) string {
	v, ok := s.(int)
	if !ok {
		return "???"
	}
	smap := [...]string{
		/* 0 */ "emerg",
		/* 1 */ "alert",
		/* 2 */ "crit",
		/* 3 */ "error",
		/* 4 */ "warn",
		/* 5 */ "notice",
		/* 6 */ "info",
		/* 7 */ "debug",
	}

	if v > len(smap) {
		return fmt.Sprintf("%d", v)
	}
	return smap[v]
}

func vsyslogOutput(logfile string, ch chan syslogparser.LogParts) {
	logger := lumberjack.Logger{
		Filename:   logfile,
		MaxSize:    50, // megabytes
		MaxBackups: 5,  // 50*5 == 250 megabytes of logs in total at one time
	}

	timeFormat := "2006-01-02 15:04:05 -0700 MST"

	logger.Write([]byte(fmt.Sprintf("%s --- STARTING ---\n",
		time.Now().UTC().Format(timeFormat))))

	for {
		lp := <-ch

		msg := fmt.Sprintf("%s %s.%s %s: %s",
			lp["timestamp"].(time.Time).Format(timeFormat),
			facilityString(lp["facility"]),
			severityString(lp["severity"]),
			lp["tag"], lp["content"])

		logger.Write([]byte(fmt.Sprintf("%s\n", msg)))

		// Log to console too if debug, which will hit ASL as
		// well, which we don't want for normal users
		if debug {
			msg = fmt.Sprintf("VM: %s", msg)
			sev, ok := lp["severity"].(int)
			if !ok {
				sev = 4 // Warn seems appropriate if severity was not an integer...
			}
			if sev <= 3 { // error or worse anything worse
				// than logrus.Error API will kill the
				// process which isn't what we want.
				logrus.Error(msg)
			} else if sev <= 5 { // warn or notice
				logrus.Warn(msg)
			} else if sev <= 6 { // info
				logrus.Info(msg)
			} else { // debug
				logrus.Debug(msg)
			}
		}
	}
}

func handleVSyslog() {
	logfile := driverDir + "/syslog"

	ch := make(chan syslogparser.LogParts)
	go vsyslogOutput(logfile, ch)

	sock := pinataSockets.GetVsockSocketPath(vsockHostCID, VsyslogPort)

	logrus.Printf("Syslog socket is %s", sock)
	logrus.Printf("Logfile is %s", logfile)

	err := os.Remove(sock)
	if err != nil && !os.IsNotExist(err) {
		logrus.Fatalf("Error removing %s: %s", sock, err)
	}

	l, err := net.ListenUnix("unix", &net.UnixAddr{Name: sock, Net: "unix"})
	if err != nil {
		logrus.Fatalf("Unable to listen to vsyslog sock %s: %s", sock, err)
	}

	for {
		s, err := l.AcceptUnix()
		if err != nil {
			logrus.Fatalf("Unable to accept on syslog vsock: %s", err)
			continue
		}

		go func(s *net.UnixConn) {
			defer s.Close()

			s.CloseWrite()

			r := bufio.NewReader(s)

			for {
				/* RFC5425 like scheme, section 4.3 */

				lenstr, err := r.ReadString(' ')
				if err != nil {
					logrus.Printf("Error reading vsyslog length: %s", err)
					return
				}

				lenstr = lenstr[:len(lenstr)-1] // Trim the ' '

				len, err := strconv.ParseUint(lenstr, 10, 32)
				if err != nil {
					logrus.Printf("Error parsing vsyslog length: %s", err)
					return
				}

				/* RFC5425 says we SHOULD support up to 0x2000, support twice that */
				if len >= 0x4000 {
					logrus.Printf("vsyslog length too large: %d > %d", len, 0x4000)
					return
				}

				buf := make([]byte, len)
				nr, err := r.Read(buf)
				if err != nil {
					logrus.Printf("Error reading vsyslog msg: read %d/%d (%#x/%#x) bytes: %s",
						nr, len, nr, len, string(buf[:nr]))
					return
				}

				if uint64(nr) != len {
					logrus.Printf("Incorrect read length: read %d/%d (%#x/%#x) bytes: %s",
						nr, len, nr, len, string(buf[:nr]))
					return
				}

				if buf[nr-1] == '\n' { // Strip trailing CR
					nr--
					buf = buf[:nr]
				}
				p := rfc3164.NewParser(buf)
				p.Hostname("docker") // Busybox doesn't include the hostname, have rfc3164 DTRT
				err = p.Parse()
				if err != nil {
					logrus.Printf("Failed to parse syslog: %s", string(buf[:nr]))
					/* We are now possibly out of
					/* sync with the stream, quit
					/* and allow the other end to
					/* reconnect */

					return
				}
				lp := p.Dump()
				ch <- lp
			}
		}(s)
	}
}
