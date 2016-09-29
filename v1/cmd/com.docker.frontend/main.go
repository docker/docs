package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/docker/pinata/v1/apple"
	appleutil "github.com/docker/pinata/v1/apple/util"
	"github.com/docker/pinata/v1/reportError"
)

func main() {
	logrus.AddHook(apple.NewLogrusASLHook())
	reportError.Initialize()
	// logrus.Println("Starting HTTP/2 client...")

	// os.Args[0] // program's name
	requestBody := os.Args[1]
	// logrus.Println("Sending request:", requestBody)
	// TODO: gdevillele: test that requestBody is not empty

	// First, construct the path to the socket
	var socketPath = filepath.Join(appleutil.GetContainerPath(), "s20")

	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		logrus.Fatal(err)
	}

	clientConn := httputil.NewClientConn(conn, nil)
	request, err := http.NewRequest("GET", "/", strings.NewReader(requestBody))
	if err != nil {
		logrus.Fatal(err)
	}

	// setup timeout goroutine
	timedOut := false
	go func() {
		time.Sleep(time.Second * 10)
		timedOut = true
		clientConn.Close()
	}()

	// perform request
	res, err := clientConn.Do(request)
	if err != nil {
		if timedOut {
			os.Exit(2)
		} else {
			logrus.Fatal(err)
		}
	}
	defer res.Body.Close()

	// process response
	if res.StatusCode != 200 {
		logrus.Fatalf("Got response code %d from backend", res.StatusCode)
	}

	bodyData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		if timedOut {
			os.Exit(2)
		} else {
			logrus.Fatal(err)
		}
	}
	bodyString := string(bodyData)

	fmt.Printf(bodyString)
	// logrus.Println("RESPONSE:", res.StatusCode)
	// logrus.Println("RESPONSE:", string(body))
	// logrus.Println("Quitting HTTP/2 client")
}
