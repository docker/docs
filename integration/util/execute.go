package util

import (
	"fmt"
	"strings"
	"testing"

	"github.com/docker/dhe-deploy/integration/sshclient"
	"github.com/stretchr/testify/require"
)

func (u *Util) Execute(cmd string, ignoreFail bool) string {
	output, stderr, exitcode := RawExecute(u.T(), u.SSH, cmd)
	if exitcode != 0 && !ignoreFail {
		require.FailNow(u.T(), fmt.Sprintf("Command exited with non-zero exit code : %s returned %v:\n%s\n%s", cmd, exitcode, output, stderr))
	}
	// trim trailing newline of execute output
	output = strings.TrimSuffix(output, "\n")
	return output
}

func (u *Util) RawExecute(cmd string) (string, string, int) {
	return RawExecute(u.T(), u.SSH, cmd)
}

func Execute(t *testing.T, ssh sshclient.SSHClient, cmd string, ignoreFail bool) string {
	output, stderr, exitcode := RawExecute(t, ssh, cmd)
	if exitcode != 0 && !ignoreFail {
		t.Fatalf("Command exited with non-zero exit code : %s returned %v:\n%s\n%s", cmd, exitcode, output, stderr)
	}
	// trim trailing newline of execute output
	output = strings.TrimSuffix(output, "\n")
	return output
}

func RawExecute(t *testing.T, ssh sshclient.SSHClient, cmd string) (string, string, int) {
	output, stderr, err := ssh.RunRemoteCommand(cmd)
	exitCode := 0
	if nonzeroExitCodeErr, ok := err.(sshclient.NonzeroExitCodeError); ok {
		exitCode = nonzeroExitCodeErr.ExitCode
	} else if err != nil {
		t.Fatalf("Failed to run command on remote host: %s", err)
	}
	output = strings.TrimSuffix(output, "\n")
	return output, stderr, exitCode
}
