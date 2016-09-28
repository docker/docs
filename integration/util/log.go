package util

import (
	"fmt"
	"strings"

	"github.com/docker/dhe-deploy/shared/containers"
)

const (
	allDTR    = "allDTR"
	allDocker = "allDocker"
)

// TODO: move all of these variables into the framework instead of making them global
var KnownDockerRestartLogErrors = []string{"write unix @: broken pipe", "write: broken pipe", "Error streaming logs: unexpected EOF", "Couldn't run auplink before unmount: exit status 22", "network with name dtr already exists"}
var KnownDTRRestartLogErrors = []string{"Pinging the database failed.", "ERROR unexpected end of JSON input", "host not found in upstream", "No such file or directory)\\nnginx", "exit status 1", "No such file or directory:fopen('/nginx-ssl/server.pem", "invalid PID number"}

var KnownDockerInstallLogErrors = []string{"Cannot kill container",
	"notrunning: Container",
	"write unix @: broken pipe",
	"write: broken pipe",
	"No such image:",
	"no such id:",
	"No such container:",
	"remove /var/run/docker.pid: no such file or directory",
	"Error streaming logs: unexpected EOF",
	"network with name dtr already exists",
	"external resolution failed",
	"Couldn't run auplink before unmount: exit status 22",
	"Unexpected watch error",
	"client: etcd cluster is unavailable or misconfigured",
	"Error unmounting container",
	"Failed to load container",
	"Container already stopped",
	"discovery error: 102: Not a file (/docker/nodes)",
	"Force shutdown daemon",
	"is already stopped",
	"device or resource busy",
}
var KnownDTRInstallLogErrors = append([]string{"ERROR:  duplicate key value violates unique constraint \"accounts_name_key\"", "error reading container stats; no longer getting stats for container", "ERROR:  relation \"properties\" does not exist at character 19", "Pinging the database failed."}, KnownDTRRestartLogErrors...)

var KnownDockerLogErrors = []string{"external resolution failed"}
var KnownDTRLogErrors = []string{}

var logLineTracker = make(map[string]int)
var logIgnorablesTracker = make(map[string][]string)

func (u *Util) TestLogs() {
	u.testForDockerLogErrors()
	u.testForDTRLogErrors()
}

func AppendDTRIgnorableLoggedErrors(ignorableErrorLogs []string) {
	logIgnorablesTracker[allDTR] = append(logIgnorablesTracker[allDTR], ignorableErrorLogs...)
}

func AppendDockerIgnorableLoggedErrors(ignorableErrorLogs []string) {
	logIgnorablesTracker[allDocker] = append(logIgnorablesTracker[allDocker], ignorableErrorLogs...)
}

func WipeDTRIgnorableLoggedErrors() {
	logIgnorablesTracker[allDTR] = KnownDTRLogErrors
}

func WipeDockerIgnorableLoggedErrors() {
	logIgnorablesTracker[allDocker] = KnownDockerLogErrors
}

// TestForDTRLogErrors calls ErrorsInDTRLogFile to retrieve all the errors in the dtr logs and parses out any known errors that we want to ignore while printing the rest.
// NOTE: if there are any new bugs that cause errors that you want to ignore, add them to the string of if statements.
func (u *Util) testForDTRLogErrors() {
	errorsToIgnore := logIgnorablesTracker[allDTR]
	// add in stupid non-error registry errors...
	errorsToIgnore = append(errorsToIgnore, `msg="response completed with error"`)

	for _, component := range containers.AllContainers {
		logFile := component.Name
		output := strings.Split(strings.TrimRight(u.errorsInDTRLogFile(logFile), "\n"), "\n")
		linesRead := len(output)
		start := logLineTracker[logFile]
		logLineTracker[logFile] = linesRead
		errors := 0
		for _, v := range output[start:] {
			ignorable := false

			for _, ignorableSubstring := range errorsToIgnore {
				if strings.Contains(v, ignorableSubstring) {
					ignorable = true
				}
			}

			if !ignorable && v != "" {
				fmt.Printf("FAILURE: found unexpected error: %s\n", v)
				errors++
			}
		}

		//assert.Equal(u.T(), 0, errors, fmt.Sprintf("%d errors were found in the DTR logs for %s", errors, logFile))
	}
}

// TestForDockerLogErrors calls ErrorsInDockerLogFile to retrieve all the errors in the docker logs and parses out any known errors that we want to ignore while printing the rest.
// NOTE: if there are any new bugs that cause errors that you want to ignore, add them to the string of if statements.
func (u *Util) testForDockerLogErrors() {
	output := strings.Split(strings.TrimRight(u.errorsInDockerLogFile(), "\n"), "\n")
	linesRead := len(output)
	start := logLineTracker[allDocker]
	logLineTracker[allDocker] = linesRead
	errorsToIgnore := append(logIgnorablesTracker[allDocker], "/stop returned error: Cannot stop container")
	errors := 0

	for _, v := range output[start:] {
		ignorable := false

		for _, ignorableSubstring := range errorsToIgnore {
			if strings.Contains(v, ignorableSubstring) {
				ignorable = true
			}
		}

		if !ignorable && v != "" {
			fmt.Printf("FAILURE: found unexpected error: %s\n", v)
			errors++
		}
	}
	//assert.Equal(u.T(), 0, errors, fmt.Sprintf("%d errors were found in the Docker logs", errors))
}

// ErrorsInDockerLogFile uses DetectHostOS to get the host os so that the location of the docker logs can be found. The logs are then parsed for errors and the errors are returned.
func (u *Util) errorsInDockerLogFile() string {
	hostOS := DetectHostOS(u.T(), u.SSH)
	if hostOS.IsSystemd {
		return Execute(u.T(), u.SSH, "sudo journalctl SYSLOG_IDENTIFIER=docker --boot | grep level=error", true)
	}

	return Execute(u.T(), u.SSH, fmt.Sprintf("sudo egrep -Hn %s %s", hostOS.SyslogRegex, hostOS.SyslogPath), true)
}

// ErrorsInDTRLogFile parses the dtr logs for any errors and returns those errors.
func (u *Util) errorsInDTRLogFile(filename string) string {
	return Execute(u.T(), u.SSH, fmt.Sprintf(`sudo sh -c 'grep -Hn "ERROR\|level=error" /usr/local/etc/dtr/logs/%s'`, filename), true)
}
