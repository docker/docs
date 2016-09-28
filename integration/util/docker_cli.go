package util

import (
	"fmt"
	"strings"
)

func (u *Util) DockerLogin(username, password string) error {
	loginCmd := fmt.Sprintf("sudo docker login -u %s -p %s -e 'a@a.a' %s", username, password, u.Config.DTRHost)
	if stdout, stderr, exitCode := u.RawExecute(loginCmd); exitCode != 0 {
		return fmt.Errorf("cannot do docker login: stdout=%q, stderr=%q, exitCode=%d", stdout, stderr, exitCode)
	}
	return nil
}

func (u *Util) DockerLogout() error {
	logoutCmd := fmt.Sprintf("sudo docker logout %s", u.Config.DTRHost)
	if stdout, stderr, exitCode := u.RawExecute(logoutCmd); exitCode != 0 {
		return fmt.Errorf("cannot do docker logout: stdout=%q, stderr=%q, exitCode=%d", stdout, stderr, exitCode)
	}
	return nil
}

func (u *Util) DockerSearch(query string) ([]string, error) {
	searchCmd := fmt.Sprintf("sudo docker search %s/%s", u.Config.DTRHost, "foobar")
	stdout, stderr, exitcode := u.RawExecute(searchCmd)
	if exitcode != 0 {
		return nil, fmt.Errorf("Exited with code %d, stdout=%q, stderr=%q", exitcode, stdout, stderr)
	}
	var searchResults []string
	lines := strings.Split(stdout, "\n")
	// skip the column names line:
	// NAME      DESCRIPTION   STARS     OFFICIAL   AUTOMATED
	// and the trailing newline
	for _, line := range lines[1 : len(lines)-1] {
		searchResults = append(searchResults, strings.Fields(line)[0])
	}
	return searchResults, nil
}
