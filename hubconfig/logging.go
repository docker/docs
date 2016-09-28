package hubconfig

import (
	"fmt"
	"log/syslog"
	"strings"

	log "github.com/Sirupsen/logrus"
)

// SyslogLogrusLevels interprets the log levels for syslog and logrus
func (c HAConfig) SyslogLogrusLevels() (syslogLevel syslog.Priority, logrusLevel log.Level, err error) {
	level := strings.ToUpper(c.LogLevel)
	switch {
	case strings.Contains(level, "EMERG"):
		syslogLevel = syslog.LOG_EMERG
		logrusLevel = log.FatalLevel
	case strings.Contains(level, "ALERT"):
		syslogLevel = syslog.LOG_ALERT
		logrusLevel = log.FatalLevel
	case strings.Contains(level, "CRIT"):
		syslogLevel = syslog.LOG_CRIT
		logrusLevel = log.FatalLevel
	case strings.Contains(level, "ERR"):
		syslogLevel = syslog.LOG_ERR
		logrusLevel = log.ErrorLevel
	case strings.Contains(level, "WARN"):
		syslogLevel = syslog.LOG_WARNING
		logrusLevel = log.WarnLevel
	case strings.Contains(level, "NOTICE"):
		syslogLevel = syslog.LOG_NOTICE
		logrusLevel = log.InfoLevel
	case strings.Contains(level, "INFO"):
		syslogLevel = syslog.LOG_INFO
		logrusLevel = log.InfoLevel
	case strings.Contains(level, "DEBUG"):
		syslogLevel = syslog.LOG_DEBUG
		logrusLevel = log.DebugLevel
	default:
		log.Infof("Unrecognised log level %s - you must use standard syslog levels", level)
		// Use default log levels
		syslogLevel = syslog.LOG_INFO
		logrusLevel = log.InfoLevel
		err = fmt.Errorf("Unknown LogLevel: %s", c.LogLevel)
	}
	return
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

var validProtocols = []string{"tcp", "tcp+tls", "udp", "internal"}

func (c HAConfig) LoggingValid() error {
	logProtocol := c.LogProtocol
	if !stringInSlice(logProtocol, validProtocols) {
		return fmt.Errorf("LogProtocol should be one of: %s. Instead, it's '%s'.", strings.Join(validProtocols, ","), logProtocol)
	}

	if c.LogHost == "" && c.LogProtocol != "internal" {
		return fmt.Errorf("LogHost should be set if the LogProtocol is tcp, tcp+tls, or udp.")
	}
	return nil
}

// LogTest validates the logging config, and optionally sends a ping to the syslog server if there is one
func (c HAConfig) LogTest() error {
	if err := c.LoggingValid(); err != nil {
		return err
	}
	if _, _, err := c.SyslogLogrusLevels(); err != nil {
		return err
	}
	if c.LogProtocol == "tcp" || c.LogProtocol == "udp" {
		syslogWriter, err := syslog.Dial(c.LogProtocol, c.LogHost, syslog.LOG_INFO, "dtr-testing-ping")
		if err != nil {
			return err
		}
		return syslogWriter.Close()
	}
	return nil
}
