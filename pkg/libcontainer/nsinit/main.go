package main

import (
	"encoding/json"
	"errors"
	"flag"
	"github.com/dotcloud/docker/pkg/libcontainer"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

var (
	ErrUnsupported    = errors.New("Unsupported method")
	ErrWrongArguments = errors.New("Wrong argument count")
)

func main() {
	var (
		console = flag.String("console", "", "Console (pty slave) name")
		tty     = flag.Bool("tty", false, "Create a tty")
		pipeFd  = flag.Int("pipe", 0, "sync pipe fd")
	)
	flag.Parse()

	container, err := loadContainer()
	if err != nil {
		log.Fatal(err)
	}

	if flag.NArg() < 1 {
		log.Fatal(ErrWrongArguments)
	}
	switch flag.Arg(0) {
	case "exec": // this is executed outside of the namespace in the cwd
		var exitCode int
		nspid, err := readPid()
		if err != nil {
			if !os.IsNotExist(err) {
				log.Fatal(err)
			}
		}
		if nspid > 0 {
			exitCode, err = execinCommand(container, nspid, flag.Args()[1:])
		} else {
			exitCode, err = execCommand(container, *tty, flag.Args()[1:])
		}
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(exitCode)
	case "init": // this is executed inside of the namespace to setup the container
		if flag.NArg() < 2 {
			log.Fatal(ErrWrongArguments)
		}
		if err := initCommand(container, *console, os.NewFile(uintptr(*pipeFd), "pipe"), flag.Args()[1:]); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("command not supported for nsinit %s", flag.Arg(0))
	}
}

func loadContainer() (*libcontainer.Container, error) {
	f, err := os.Open("container.json")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var container *libcontainer.Container
	if err := json.NewDecoder(f).Decode(&container); err != nil {
		return nil, err
	}
	return container, nil
}

func readPid() (int, error) {
	data, err := ioutil.ReadFile(".nspid")
	if err != nil {
		return -1, err
	}
	pid, err := strconv.Atoi(string(data))
	if err != nil {
		return -1, err
	}
	return pid, nil
}
