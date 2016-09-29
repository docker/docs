package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/crypto/ssh"
)

const (
	// WriteDirectory is the SCP Command to write a directory
	WriteDirectory = "D"
	// WriteFile is the SCP Command to write a file
	WriteFile = "C"
	// EndWriteFile is the SCP Command to stop writing a file
	EndWriteFile = "\x00"
	// EndWriteDirectory ends writing a directory
	EndWriteDirectory = "E"
)

// UploadDir uploads a directory to a server
func UploadDir(client *ssh.Client, src string, dst string) (<-chan error, error) {
	errc := make(chan error, 1)
	session, err := client.NewSession()
	if err != nil {
		fatal(err)
	}
	defer session.Close()

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf

	var stderrBuf bytes.Buffer
	session.Stderr = &stderrBuf
	go func() {
		w, _ := session.StdinPipe()
		defer w.Close()
		fmt.Fprintln(w, "D0755", 0, filepath.Base(dst))
		err := processDir(w, src)
		fmt.Fprintln(w, EndWriteDirectory)
		errc <- err
	}()
	if err := session.Run("/usr/bin/scp -qtr " + dst); err != nil {
		return errc, fmt.Errorf("Ohnoes: %s\nStdout:\n%s\nStderr:%s\n", err, session.Stdout, session.Stderr)
	}
	return errc, nil
}

func unixPermissions(f os.FileInfo) (string, error) {
	return fmt.Sprintf("%04o", f.Mode().Perm()), nil
}

func processDir(w io.WriteCloser, dir string) error {
	fl, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, f := range fl {
		name := filepath.Join(dir, f.Name())
		if f.IsDir() {
			dir, err := os.Open(name)
			if err != nil {
				return err
			}
			defer dir.Close()
			perms, err := unixPermissions(f)
			if err != nil {
				return err
			}
			mode := fmt.Sprintf("%s%s", WriteDirectory, perms)
			fmt.Fprintln(w, mode, 0, f.Name())
			processDir(w, name)
			fmt.Fprintln(w, EndWriteDirectory)
		} else {
			file, err := os.Open(name)
			if err != nil {
				return err
			}
			defer file.Close()
			fs, err := file.Stat()
			if err != nil {
				return err
			}
			perms, err := unixPermissions(f)
			mode := fmt.Sprintf("%s%s", WriteFile, perms)
			fmt.Fprintln(w, mode, fs.Size(), filepath.Base(name))
			io.Copy(w, file)
			fmt.Fprint(w, EndWriteFile)
		}
	}
	return nil
}
