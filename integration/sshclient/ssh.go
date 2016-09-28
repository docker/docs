package sshclient

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"golang.org/x/crypto/ssh"
)

type NonzeroExitCodeError struct {
	ExitCode int
}

func (e NonzeroExitCodeError) Error() string {
	return fmt.Sprintf("exit: %d", e.ExitCode)
}

type SSHClient interface {
	ConnectToRemoteHost() (*ssh.Session, error)
	RunRemoteCommand(cmd string) (string, string, error)
}

type sshClient struct {
	Host     string
	Port     int
	Username string
	Password string
	KeyPath  string
}

func New(host, username, password, keypath string, port int) SSHClient {
	return &sshClient{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		KeyPath:  keypath,
	}
}

func (c *sshClient) ConnectToRemoteHost() (*ssh.Session, error) {
	// sign the private key string
	// XXX: Does not support loading from string
	var auth []ssh.AuthMethod
	if c.Password == "" {
		buffer, err := ioutil.ReadFile(c.KeyPath)
		if err != nil {
			return nil, fmt.Errorf("Could not find key at %s : %s", c.KeyPath, err)
		}

		signer, err := ssh.ParsePrivateKey(buffer)
		if err != nil {
			return nil, fmt.Errorf("ParsePrivateKey error: %s", err)
		}
		auth = []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		}
	} else {
		auth = []ssh.AuthMethod{
			ssh.Password(c.Password),
		}
	}

	// create ssh client config
	sshClientConfig := &ssh.ClientConfig{
		User: c.Username,
		Auth: auth,
	}
	// connect to remote host with config
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", c.Host, c.Port), sshClientConfig)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to ssh host: %s", err)
	}
	// est. an ssh client session
	session, err := client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("Fail to create ssh remote session: %s", err)
	}
	return session, nil
}

func (c *sshClient) RunRemoteCommand(cmd string) (string, string, error) {
	// connect to remote host and est. a session
	session, err := c.ConnectToRemoteHost()
	if err != nil {
		return "", "", err
	}
	defer session.Close()

	// create buffer to capture stdout
	var bOut bytes.Buffer
	var bErr bytes.Buffer
	session.Stdout = &bOut
	session.Stderr = &bErr

	// Warning: No timeout exists
	// Execute the command and capture err.
	// Ignore all 'ExitError', and capture exit code
	if err := session.Run(cmd); err != nil {
		if exiterr, ok := err.(*ssh.ExitError); ok {
			if exitcode := exiterr.ExitStatus(); exitcode != 0 {
				return bOut.String(), bErr.String(), NonzeroExitCodeError{exitcode}
			}
		}
		return "", "", fmt.Errorf("Failed to run %q command remotely: %s", cmd, err)
	}

	return bOut.String(), bErr.String(), nil
}
