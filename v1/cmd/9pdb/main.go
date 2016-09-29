package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"time"

	datakit "github.com/docker/datakit/api/go-datakit"
	"golang.org/x/net/context"
)

var dbPath string

func usage() {
	fmt.Fprintf(os.Stderr, "Flags of %s:\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\nAddional arguments describe a get or set command:\n")
	fmt.Fprintf(os.Stderr, "  %s get <key1> <key2> ... <keyN>: read a set of keys from a snapshot\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s set <key1>=<value1> <key2>=<value2> ... <keyN>=<valueN>: update keys in a transaction\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s set-file <key1>=<value1> <key2>=<value2> ... <keyN>=<valueN>: update keys in a transaction with values extracted from files\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s rm <key1> <key2> ... <keyN>: remove a set of keys in a transaction\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\nExamples:\n")
	fmt.Fprintf(os.Stderr, "  $ %s get com.docker.driver.amd64-linux/network\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  com.docker.driver.amd64-linux/network=hybrid\n")
	fmt.Fprintf(os.Stderr, "  $ %s set com.docker.driver.amd64-linux/test=hello\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  $ %s rm com.docker.driver.amd64-linux\n", os.Args[0])
	os.Exit(1)
}

func main() {
	flag.StringVar(&dbPath, "path", GetDefaultDBPath(), "path to the database Unix domain socket or Named pipe")
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Please supply a command: either get, set or set-file\n")
		usage()
	}

	ctx := context.Background()

	var conn net.Conn
	var err error
	for retries := 0; retries < 10; retries++ {
		conn, err = Dial(dbPath)
		if err == nil {
			break
		}

		duration := 100 * time.Millisecond
		log.Println("Retry dialing", dbPath, "in", duration)
		time.Sleep(duration)
	}
	if err != nil {
		log.Fatalf("Failed to contact the database on %s: %s", dbPath, err)
	}

	client, err := datakit.NewClient(ctx, conn)
	if err != nil {
		log.Fatalf("Failed to establish 9P session with the database on %s: %s", dbPath, err)
	}

	if args[0] == "get" {
		sha, err := datakit.Head(ctx, client, "master")
		if err != nil {
			log.Fatalf("Failed to discover the HEAD of master: %#v", err)
		}
		snap := datakit.NewSnapshot(ctx, client, datakit.COMMIT, sha)

		for _, key := range args[1:] {
			path := strings.Split(key, "/")
			value, err := snap.Read(ctx, path)
			if err != nil {
				log.Fatalf("Failed to read %s from snapshot: %#v", key, err)
			}
			fmt.Printf("%s=%s\n", key, value)
		}
		os.Exit(0)
	}

	if args[0] == "rm" {
		t, err := datakit.NewTransaction(ctx, client, "master", "9pdb")
		if err != nil {
			log.Fatalf("Failed to create database branch: %#v", err)
		}
		for _, key := range args[1:] {
			path := strings.Split(key, "/")
			if err = t.Remove(ctx, path); err != nil {
				log.Fatalf("Failed to remove %s: %#v", key, err)
			}
		}
		if err = t.Commit(ctx, "9pdb commit"); err != nil {
			log.Fatalf("Failed to commit transaction: %#v", err)
		}
		os.Exit(0)
	}

	if args[0] == "set" {
		t, err := datakit.NewTransaction(ctx, client, "master", "9pdb")
		if err != nil {
			log.Fatalf("Failed to create database branch: %#v", err)
		}
		for _, key := range args[1:] {
			parts := strings.SplitN(key, "=", 2)
			if len(parts) != 2 {
				log.Fatalf("Expected argument of the form <key>=<value>")
			}
			key := parts[0]
			value := parts[1]
			path := strings.Split(key, "/")
			if err = t.Write(ctx, path, value); err != nil {
				log.Fatalf("Failed to write %s = %v: %#v", key, value, err)
			}
		}
		if err = t.Commit(ctx, "9pdb commit"); err != nil {
			log.Fatalf("Failed to commit transaction: %#v", err)
		}
		os.Exit(0)
	}

	if args[0] != "set-file" {
		fmt.Fprintf(os.Stderr, "The command should be either 'get' or 'set' or 'rm'\n")
		usage()
	}

	// set-file
	t, err := datakit.NewTransaction(ctx, client, "master", "9pdb")
	if err != nil {
		log.Fatalf("Failed to create database branch: %#v", err)
	}
	for _, key := range args[1:] {
		parts := strings.SplitN(key, "=", 2)
		if len(parts) != 2 {
			log.Fatalf("Expected argument of the form <key>=<value>")
		}
		key := parts[0]
		fileName := parts[1]
		// eliminate double quotes for paths with spaces
		if fileName[:1] == "\"" && fileName[len(fileName)-1:] == "\"" {
			fileName = fileName[1 : len(fileName)-1]
		}
		value, err := ioutil.ReadFile(fileName)
		if err != nil {
			log.Fatalf("Failed to read input file: %#v", err)
		}
		path := strings.Split(key, "/")
		if err = t.Write(ctx, path, string(value)); err != nil {
			log.Fatalf("Failed to write %s = %v: %#v", key, parts[1], err)
		}
	}
	if err = t.Commit(ctx, "9pdb commit"); err != nil {
		log.Fatalf("Failed to commit transaction: %#v", err)
	}
	os.Exit(0)

}
