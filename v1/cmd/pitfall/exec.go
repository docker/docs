package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

// DownloadResults extracts the results of the tests over SSH
func DownloadResults(client *ssh.Client, src string, dst string) error {
	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf

	cmd := fmt.Sprintf("cat %s", src)
	if err := session.Run(cmd); err != nil {
		return err
	}

	_, err = f.Write(stdoutBuf.Bytes())
	if err != nil {
		return (err)
	}
	return nil
}

// ExecCmd executes a command over SSH
func ExecCmd(client *ssh.Client, cmd string, stream bool) error {
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	if stream {
		stdout, err := session.StdoutPipe()
		if err != nil {
			return err
		}
		outScanner := bufio.NewScanner(stdout)
		go func() {
			for outScanner.Scan() {
				fmt.Printf("STDOUT| %s\n", outScanner.Text())
			}
		}()
		stderr, err := session.StderrPipe()
		if err != nil {
			return err
		}
		errScanner := bufio.NewScanner(stderr)
		go func() {
			for errScanner.Scan() {
				fmt.Printf("STDERR| %s\n", errScanner.Text())
			}
		}()

	} else {
		var stdoutBuf bytes.Buffer
		session.Stdout = &stdoutBuf

		var stderrBuf bytes.Buffer
		session.Stderr = &stderrBuf
	}

	if err := session.Run(cmd); err != nil {
		return fmt.Errorf("Ohnoes: %s\nStdout:\n%s\nStderr:%s\n", err, session.Stdout, session.Stderr)
	}
	return nil
}
