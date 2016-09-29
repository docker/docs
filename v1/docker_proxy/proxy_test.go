package proxy

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"runtime"
	"sync"
	"testing"
	"time"
)

type rpc struct {
	name            string
	req             string // Sent to proxy
	inspect         string // Reply to optional inspect call
	expectReq       string // Modified request received by Docker
	expectApprovals []string
	resp            string // Response sent by Docker
	expectResp      string // Response received by client
}

var inspectResponse = "HTTP/1.1 200 OK\r\n" +
	"Content-Type: application/json\r\n" +
	"Server: Docker/1.10.0-rc4 (linux)\r\n" +
	"Date: Fri, 05 Feb 2016 10:27:25 GMT\r\n" +
	"Transfer-Encoding: chunked\r\n" +
	"\r\n" +
	"cb9\n" +
	"{\"Id\":\"876d34c85bb77f504fda9d6f9f6701537caec43b280d8d241a21a3a4083d455f\",\"Created\":\"2016-02-05T10:26:37.889618655Z\",\"Path\":\"/bin/bash\",\"Args\":[],\"State\":{\"Status\":\"exited\",\"Running\":false,\"Paused\":false,\"Restarting\":false,\"OOMKilled\":false,\"Dead\":false,\"Pid\":0,\"ExitCode\":0,\"Error\":\"\",\"StartedAt\":\"2016-02-05T10:26:38.136213354Z\",\"FinishedAt\":\"2016-02-05T10:26:38.262794417Z\"},\"Image\":\"sha256:3876b81b5a81678926c601cd842040a69eb0456d295cd395e7a895a416cf7d91\",\"ResolvConfPath\":\"/var/lib/docker/containers/876d34c85bb77f504fda9d6f9f6701537caec43b280d8d241a21a3a4083d455f/resolv.conf\",\"HostnamePath\":\"/var/lib/docker/containers/876d34c85bb77f504fda9d6f9f6701537caec43b280d8d241a21a3a4083d455f/hostname\",\"HostsPath\":\"/var/lib/docker/containers/876d34c85bb77f504fda9d6f9f6701537caec43b280d8d241a21a3a4083d455f/hosts\",\"LogPath\":\"/var/lib/docker/containers/876d34c85bb77f504fda9d6f9f6701537caec43b280d8d241a21a3a4083d455f/876d34c85bb77f504fda9d6f9f6701537caec43b280d8d241a21a3a4083d455f-json.log\",\"Name\":\"/foo\",\"RestartCount\":0,\"Driver\":\"aufs\",\"MountLabel\":\"\",\"ProcessLabel\":\"\",\"AppArmorProfile\":\"\",\"ExecIDs\":null,\"HostConfig\":{\"Binds\":[\"/Users/tal:/mnt\"],\"ContainerIDFile\":\"\",\"LogConfig\":{\"Type\":\"json-file\",\"Config\":null},\"NetworkMode\":\"default\",\"PortBindings\":{},\"RestartPolicy\":{\"Name\":\"no\",\"MaximumRetryCount\":0},\"VolumeDriver\":\"\",\"VolumesFrom\":null,\"CapAdd\":null,\"CapDrop\":null,\"Dns\":[],\"DnsOptions\":[],\"DnsSearch\":[],\"ExtraHosts\":null,\"GroupAdd\":null,\"IpcMode\":\"\",\"Links\":null,\"OomScoreAdj\":0,\"PidMode\":\"\",\"Privileged\":false,\"PublishAllPorts\":false,\"ReadonlyRootfs\":false,\"SecurityOpt\":null,\"UTSMode\":\"\",\"ShmSize\":67108864,\"ConsoleSize\":[0,0],\"Isolation\":\"\",\"CpuShares\":0,\"CgroupParent\":\"\",\"BlkioWeight\":0,\"BlkioWeightDevice\":null,\"BlkioDeviceReadBps\":null,\"BlkioDeviceWriteBps\":null,\"BlkioDeviceReadIOps\":null,\"BlkioDeviceWriteIOps\":null,\"CpuPeriod\":0,\"CpuQuota\":0,\"CpusetCpus\":\"\",\"CpusetMems\":\"\",\"Devices\":[],\"KernelMemory\":0,\"Memory\":0,\"MemoryReservation\":0,\"MemorySwap\":0,\"MemorySwappiness\":-1,\"OomKillDisable\":false,\"PidsLimit\":0,\"Ulimits\":null},\"GraphDriver\":{\"Name\":\"aufs\",\"Data\":null},\"Mounts\":[{\"Source\":\"/Users/tal\",\"Destination\":\"/mnt\",\"Mode\":\"\",\"RW\":true,\"Propagation\":\"rprivate\"}],\"Config\":{\"Hostname\":\"876d34c85bb7\",\"Domainname\":\"\",\"User\":\"\",\"AttachStdin\":false,\"AttachStdout\":false,\"AttachStderr\":false,\"Tty\":false,\"OpenStdin\":false,\"StdinOnce\":false,\"Env\":null,\"Cmd\":[\"/bin/bash\"],\"Image\":\"ubuntu\",\"Volumes\":null,\"WorkingDir\":\"\",\"Entrypoint\":null,\"OnBuild\":null,\"Labels\":{},\"StopSignal\":\"SIGTERM\"},\"NetworkSettings\":{\"Bridge\":\"\",\"SandboxID\":\"18f53c1b06f9a65f3c68e18b9081470a3c5e4d42d65fd0514bfea5cee5e31401\"," +
	"\"HairpinMode\":false,\"LinkLocalIPv6Address\":\"\",\"LinkLocalIPv6PrefixLen\":0,\"Ports\":null,\"SandboxKey\":\"/var/run/docker/netns/18f53c1b06f9\",\"SecondaryIPAddresses\":null,\"SecondaryIPv6Addresses\":null,\"EndpointID\":\"\",\"Gateway\":\"\",\"GlobalIPv6Address\":\"\",\"GlobalIPv6PrefixLen\":0,\"IPAddress\":\"\",\"IPPrefixLen\":0,\"IPv6Gateway\":\"\",\"MacAddress\":\"\",\"Networks\":{\"bridge\":{\"IPAMConfig\":null,\"Links\":null,\"Aliases\":null,\"NetworkID\":\"b09157ef71afcb91fc22ea1937e4f97865a6d6ae5191950f5bbfba625e71b4df\",\"EndpointID\":\"\",\"Gateway\":\"\",\"IPAddress\":\"\",\"IPPrefixLen\":0,\"IPv6Gateway\":\"\",\"GlobalIPv6Address\":\"\",\"GlobalIPv6PrefixLen\":0,\"MacAddress\":\"\"}}}}\n" +
	"\n" +
	"0\n"

var rpcs = []rpc{
	{
		name: "version",
		req: "GET /v1.22/version HTTP/1.1\r\n" +
			"Host: \r\n" +
			"User-Agent: Docker-Client/1.10.0-rc1 (linux)\r\n" +
			"\r\n",
		resp: "HTTP/1.1 200 OK\r\n" +
			"Content-Type: application/json\r\n" +
			"Server: Docker/1.10.0-rc1 (linux)\r\n" +
			"Date: Wed, 27 Jan 2016 11:07:15 GMT\r\n" +
			"Content-Length: 213\r\n" +
			"\r\n" +
			"{\"Version\":\"1.10.0-rc1\",\"ApiVersion\":\"1.22\",\"GitCommit\":\"677c593\",\"GoVersion\":\"go1.5.3\",\"Os\":\"linux\",\"Arch\":\"amd64\",\"KernelVersion\":\"4.1.13-8.pvops.qubes.x86_64\",\"BuildTime\":\"2016-01-15T18:21:13.402088518+00:00\"}\n",
		expectResp: "{\"Version\":\"1.10.0-rc1\",\"ApiVersion\":\"1.22\",\"GitCommit\":\"677c593\",\"GoVersion\":\"go1.5.3\",\"Os\":\"linux\",\"Arch\":\"amd64\",\"KernelVersion\":\"4.1.13-8.pvops.qubes.x86_64\",\"BuildTime\":\"2016-01-15T18:21:13.402088518+00:00\"}\n",
	},

	{
		name: "docker crash",
		req: "GET /v1.22/version HTTP/1.1\r\n" +
			"Host: \r\n" +
			"User-Agent: Docker-Client/1.10.0-rc1 (linux)\r\n" +
			"\r\n",
		resp:       "",
		expectResp: "Bad response from Docker engine\n",
	},

	{
		name:       "attach",
		req:        "POST /v1.22/containers/cb3bc087a0167f39f6f44c43359adc243344e569de1bf41ceed54907e400d306/attach?stderr=1&stdout=1&stream=1 HTTP/1.1\r\nHost: \r\nUser-Agent: Docker-Client/1.10.0-rc1 (darwin)\r\nContent-Length: 0\r\nConnection: Upgrade\r\nContent-Type: text/plain\r\nUpgrade: tcp\r\n\r\nHello",
		resp:       "HTTP/1.1 101 UPGRADED\r\nContent-Type: application/vnd.docker.raw-stream\r\nConnection: Upgrade\r\nUpgrade: tcp\r\n\r\nHello\n",
		expectResp: "Hello",
	},

	{
		name: "start",
		req: "POST /v1.22/containers/foo/start HTTP/1.1\r\n" +
			"Host: \r\n" +
			"User-Agent: Docker-Client/1.10.0-rc4 (darwin)\r\n" +
			"Content-Length: 0\r\n" +
			"Content-Type: text/plain\r\n" +
			"\r\n",
		inspect:         inspectResponse,
		expectApprovals: []string{"/Users/tal"},
		resp: "HTTP/1.1 204 No Content\r\n" +
			"Server: Docker/1.10.0-rc4 (linux)\r\n" +
			"Date: Fri, 05 Feb 2016 10:14:52 GMT\r\n" +
			"\r\n",
		expectResp: "",
	},

	{
		name: "start-compose",
		req: "POST /v1.21/containers/foo/start HTTP/1.1\r\n" +
			"Host: localhost\r\n" +
			"User-Agent: python-requests/2.7.0 CPython/2.7.9 Darwin/15.3.0\r\n" +
			"Content-Length: 2\r\n" +
			"Accept: */*\r\n" +
			"Accept-Encoding: gzip, deflate\r\n" +
			"Connection: keep-alive\r\n" +
			"Content-Type: application/json\r\n" +
			"\r\n" +
			"{}",
		inspect:         inspectResponse,
		expectApprovals: []string{"/Users/tal"},
		resp: "HTTP/1.1 204 No Content\r\n" +
			"Server: Docker/1.10.0-rc4 (linux)\r\n" +
			"Date: Fri, 05 Feb 2016 10:14:52 GMT\r\n" +
			"\r\n",
		expectResp: "",
	},
}

func readTimeout(r io.Reader, n int) (string, error) {
	buf := make([]byte, n)
	type resp struct {
		got int
		err error
	}
	ch := make(chan resp)
	go (func() {
		got, err := io.ReadFull(r, buf)
		ch <- resp{got: got, err: err}
	})()
	select {
	case resp := <-ch:
		return string(buf[0:resp.got]), resp.err
	case <-time.After(1 * time.Second):
		return "", errors.New("Timeout")
	}

}

func handleInspect(dockerSocket net.Listener, resp string, t *testing.T) {
	dockerS, err := dockerSocket.Accept()
	if err != nil {
		panic(err)
	}
	defer dockerS.Close()

	req, err := http.ReadRequest(bufio.NewReader(dockerS))
	if err != nil {
		panic(err)
	}
	expectedURL := "/containers/foo/json"
	if req.URL.Path != expectedURL {
		t.Errorf("Got URL: %s (expected %s)", req.URL, expectedURL)
	}

	_, err = dockerS.Write([]byte(resp))
	if err != nil {
		panic(err)
	}
}

//var approvals []string

func sliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, aa := range a {
		if aa != b[i] {
			return false
		}
	}
	return true
}

func runTest(r *rpc, proxyPath string, dockerSocket net.Listener, approver *testApprover, t *testing.T) {
	fmt.Printf("==== Running test '%s'\n", r.name)

	proxySocketConn, err := net.Dial("unix", proxyPath)
	if err != nil {
		panic(err)
	}
	defer proxySocketConn.Close()

	go func() {
		_, err := proxySocketConn.Write([]byte(r.req))
		if err != nil {
			panic(err)
		}
	}()

	if r.inspect != "" {
		handleInspect(dockerSocket, r.inspect, t)
	}

	dockerS, err := dockerSocket.Accept()
	if err != nil {
		panic(err)
	}
	defer dockerS.Close()

	expectReq := r.expectReq
	actualReq, err := readTimeout(dockerS, len(expectReq))
	if err != nil {
		panic(err)
	}

	if actualReq != expectReq {
		t.Errorf("Expected:\n%v\nGot:\n%s", expectReq, actualReq)
	} else {
		t.Logf("Request OK")
	}

	approved := approver.Approved()
	if !sliceEqual(approved, r.expectApprovals) {
		t.Errorf("Expected approvals: %v; got: %v", r.expectApprovals, approved)
	}

	done := make(chan error)
	go func() {
		if r.resp != "" {
			_, err := dockerS.Write([]byte(r.resp))
			dockerS.Close()
			done <- err
		} else {
			dockerS.Close()
			done <- nil
		}
	}()

	rawBodyReader := bufio.NewReader(proxySocketConn)
	resp, err := http.ReadResponse(rawBodyReader, nil)
	if err != nil {
		t.Fatalf("Failed to read HTTP response from proxy: %s", err)
	}

	body := bytes.NewBuffer(nil)
	_, err = io.Copy(body, resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body from proxy: %s", err)
	}

	if isRawStreamUpgrade(resp) {
		_, err = io.Copy(body, rawBodyReader)
		if err != nil {
			t.Fatalf("Error copying raw stream body from bufio: %s", err)
		}
	}

	actualResp := string(body.Bytes())
	expectResp := r.expectResp
	if actualResp != expectResp {
		t.Errorf("Expected:\n%v\nGot:\n%s", expectResp, actualResp)
	} else {
		t.Logf("Response OK")
	}

	if err := <-done; err != nil {
		panic(err)
	}
}

type testApprover struct {
	lock      sync.Mutex
	approvals []string
}

func (t *testApprover) Approve(containerID string, paths []string) error {
	if containerID != "foo" {
		panic(fmt.Errorf("Approve for wrong container %s", containerID))
	}

	t.lock.Lock()
	t.approvals = paths
	t.lock.Unlock()

	return nil
}

func (t *testApprover) Approved() []string {
	t.lock.Lock()
	approved := t.approvals
	t.approvals = make([]string, 0)
	t.lock.Unlock()

	return approved
}

func TestServe(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip()
	}

	tmpdir, err := ioutil.TempDir("", "proxy_test")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpdir)

	proxyPath := tmpdir + "/proxy.sock"
	dockerPath := tmpdir + "/docker.sock"

	proxySocket, err := net.Listen("unix", proxyPath)
	if err != nil {
		panic(err)
	}

	dockerSocket, err := net.Listen("unix", dockerPath)
	if err != nil {
		panic(err)
	}

	SetVerboseLogging(false)
	approver := &testApprover{}
	go Serve(proxySocket, dockerPath, approver, nil, nil)

	for _, rpc := range rpcs {
		runTest(&rpc, proxyPath, dockerSocket, approver, t)
	}
}
