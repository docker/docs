package utils

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/engine-api/client"
	"github.com/gorilla/websocket"
	"github.com/samalba/dockerclient"
)

var (
	NamePrefix         = os.Getenv("MACHINE_PREFIX") + "OrcaTest"
	Timeout            = 30 * time.Second
	BusyboxImage       = "busybox"
	PathDoesNotExist   = errors.New("The specified path does not exist on the host")
	ServerCertHostPath = "/var/lib/docker/volumes/ucp-controller-server-certs/_data/cert.pem"
	RetryCount         = 3
)

// Interface for test machine management
type Machine interface {
	GetName() string
	GetDockerHost() string
	GetClient() (*dockerclient.DockerClient, error)
	GetEngineAPI() (*client.Client, error)
	Remove() error
	Stop() error
	Start() error
	GetIP() (string, error)
	GetInternalIP() (string, error)
	CatHostFile(hostPath string) ([]byte, error)
	TarHostDir(hostPath string) ([]byte, error)
	MachineSSH(command string) (string, error)
}

// Use docker-machine to create a test engine which can then be used for integration tests (try RetryCount times)
func GetTestMachines(count int) ([]Machine, error) {
	if os.Getenv("MACHINE_LOCAL") != "" {
		if count == 1 {
			m, err := NewLocalMachine()
			return []Machine{m}, err
		} else {
			return nil, fmt.Errorf("MACHINE_LOCAL only works for single node tests")
		}
	} else if os.Getenv("MACHINE_INVENTORY") != "" {
		return nil, fmt.Errorf("MACHINE_INVENTORY not yet supported")
		//return GetMachinesFromInventory(count) // TODO
	} else {
		machines := []Machine{}
		for i := 0; i < count; i++ {
			m, err := NewBuildMachine()
			if err != nil {
				for _, m := range machines {
					m.Remove()
					//m.Finished(false) // TODO
				}
				return nil, err
			}
			machines = append(machines, m)
		}
		return machines, nil
	}
}

// Get rid of this
func CreateTestMachine() (Machine, error) {

	if os.Getenv("MACHINE_LOCAL") != "" {
		return NewLocalMachine()
	} else {
		return NewBuildMachine()
	}
	// TODO - implement other external machine cases here...
}

// Return a manifest of the files on the host in the directory (using find $hostpath)
func HostDirManifest(m Machine, hostPath string) (map[string]interface{}, error) {
	cmd := []string{"sh", "-c", "cd /theDir && find . -print"}
	binds := []string{fmt.Sprintf("%s:/theDir:ro", hostPath)}
	log.Debug("Running find")
	data, err := runCommand(m, BusyboxImage, cmd, binds, []string{})
	if err != nil {
		log.Debugf("failed: %s", err)
		return nil, err
	}
	res := make(map[string]interface{})
	for _, line := range strings.Split(string(data), "\n") {
		log.Debug(line)
		res[line] = struct{}{}
	}
	return res, nil
}

// Get the contents of a specific file on the engine
func catHostFile(m Machine, hostPath string) ([]byte, error) {
	cmd := []string{"cat", "/thefile"}
	binds := []string{fmt.Sprintf("%s:/thefile:ro", hostPath)}
	return runCommand(m, BusyboxImage, cmd, binds, []string{})
}

// Get the content of a directory as a tar file from the engine
func tarHostDir(m Machine, hostPath string) ([]byte, error) {
	// TODO - Might want to consider compression if we find we're transfering significant data during tests
	cmd := []string{"tar", "--directory", "/theDir", "-cf", "-", "."}
	binds := []string{fmt.Sprintf("%s:/theDir:ro", hostPath)}
	return runCommand(m, BusyboxImage, cmd, binds, []string{})
}

func runCommand(m Machine, image string, cmd, binds, entrypoint []string) ([]byte, error) {
	log.Debugf("Running - image:%s entrypoint:%s cmd:%s binds:%s", image, entrypoint, cmd, binds)
	c, err := m.GetClient()
	if err != nil {
		return nil, err
	}
	if _, err := c.InspectImage(image); err != nil {
		log.Infof("Pulling %s", image)
		if err := c.PullImage(image, nil); err != nil {
			return nil, fmt.Errorf("Failed to pull %s: %s", image, err)
		}
	}
	cfg := &dockerclient.ContainerConfig{
		Image:        image,
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          cmd,
	}
	if len(entrypoint) > 0 {
		cfg.Entrypoint = entrypoint
	}
	containerId, err := c.CreateContainer(cfg, "", nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to create cat container%s", err)
	}

	defer c.RemoveContainer(containerId, true, false)

	hostConfig := &dockerclient.HostConfig{Binds: binds}
	// Start the container
	if err := c.StartContainer(containerId, hostConfig); err != nil {
		log.Debugf("Failed to launch inspection container: %s", err)
		return nil, err
	}

	defer c.StopContainer(containerId, 5)

	reader, err := c.ContainerLogs(containerId, &dockerclient.LogOptions{
		Follow:     true,
		Stdout:     true,
		Stderr:     true,
		Timestamps: false,
	})
	if err != nil {
		return nil, err
	}

	stdoutBuffer := new(bytes.Buffer)
	stderrBuffer := new(bytes.Buffer)

	// stdCopy is really chatty in debug mode
	oldLevel := log.GetLevel()
	log.SetLevel(log.InfoLevel)
	defer log.SetLevel(oldLevel)
	if _, err = stdcopy.StdCopy(stdoutBuffer, stderrBuffer, reader); err != nil {
		log.Info("cannot read logs from logs reader")
		return nil, err
	}
	stderr := stderrBuffer.String()
	if strings.Contains(strings.ToLower(stderr), "no such file") {
		// XXX This doesn't seem to hit...
		log.Info("Got a no such file on stderr")
		log.Info(stderr)
		return nil, PathDoesNotExist
	}

	info, err := c.InspectContainer(containerId)
	if err != nil {
		return nil, fmt.Errorf("Failed to inspect container after completion: %s", err)
	}
	if info.State == nil {
		return nil, fmt.Errorf("Container didn't finish!")
	}

	if info.State.ExitCode != 0 {
		// return nil, fmt.Errorf("Container exited with %d", info.State.ExitCode)
		// XXX We'll assume an error is the path didn't exist, not some other random glitch.
		//     Not ideal, but the log output doesn't seem to contain the "no such file"
		//     as expected.
		log.Info("Non zero exit code: %d", info.State.ExitCode)
		return nil, PathDoesNotExist
	}

	// Looks like it worked OK
	return stdoutBuffer.Bytes(), nil
}

// ~Copied from bootstrapper...
func VolumeExists(client *dockerclient.DockerClient, name string) bool {
	// TODO - dockerclient really needs a "GetVolume" call for large deployments - this may be a perf problem
	volumes, err := client.ListVolumes()
	if err != nil {
		return false
	}
	for _, volume := range volumes {
		if volume.Name == name {
			return true
		}
	}
	return false
}

func LoadFileInVolume(client *dockerclient.DockerClient, volname, filename, contents string) error {
	if !VolumeExists(client, volname) {
		if _, err := client.CreateVolume(&dockerclient.VolumeCreateRequest{Name: volname}); err != nil {
			return err
		}
	}

	cfg := &dockerclient.ContainerConfig{
		OpenStdin: true,
		Image:     "busybox",
		Cmd: []string{
			"sh", "-c",
			fmt.Sprintf("mkdir -p $(dirname /data/%s); cat - > /data/%s", filename, filename),
		},
	}
	hostConfig := &dockerclient.HostConfig{
		Binds: []string{
			fmt.Sprintf("%s:/data", volname),
		},
	}

	containerId, err := client.CreateContainer(cfg, "", nil)
	if err != nil {
		log.Fatal(err)
	}

	defer client.RemoveContainer(containerId, false, false)

	if err := client.StartContainer(containerId, hostConfig); err != nil {
		return fmt.Errorf("Failed to launch container to load data to volume: %s\n", err)
	}
	defer client.StopContainer(containerId, 5)

	attachURL := fmt.Sprintf("wss://%s/containers/%s/attach/ws?logs=0&stream=1&stdin=1", client.URL.Host, containerId)
	websocket.DefaultDialer.TLSClientConfig = client.TLSConfig
	ws, _, err := websocket.DefaultDialer.Dial(attachURL, nil)
	if err != nil {
		return fmt.Errorf("dial:", err)
	}
	defer ws.Close()

	err = ws.WriteMessage(websocket.BinaryMessage, []byte(contents))
	if err != nil {
		return fmt.Errorf("write:", err)
	}
	log.Debugf("Created file %s on volume %s", filename, volname)
	return nil
}

// Set up a volume on the target machine for bring-your-own server cert scenarios
func BYOServerCertInit(m Machine) error {
	volname := "ucp-controller-server-certs"
	ip, _ := m.GetIP()
	internalIP, _ := m.GetInternalIP()

	configContents := `{
    "signing": {
        "default": {
            "expiry": "8760h"
        },
        "profiles": {
            "client": {
                    "usages": [
                            "signing",
                            "key encipherment",
                            "client auth"
                    ],
                    "expiry": "8760h"
            },
            "node": {
                    "usages": [
                            "signing",
                            "key encipherment",
                            "server auth",
                            "client auth"
                    ],
                    "expiry": "8760h"
            },
            "intermediate": {
                    "usages": [
                            "signing",
                            "key encipherment",
                            "cert sign",
                            "crl sign"
                    ],
                    "is_ca": true,
                    "expiry": "8760h"
            }
        }
    }
}
`
	caContents := `{
    "key": {
        "algo": "rsa",
        "size": 4096
    },
    "CN": "My Test CA"
}
`
	serverContents := fmt.Sprintf(`
{
    "hosts": [
        "127.0.0.1",
        "%s",
        "%s"
    ],
    "key": {
        "algo": "rsa",
        "size": 4096
    },
    "CN": "controller"
}`, ip, internalIP)

	c, err := m.GetClient()
	if err != nil {
		return err
	}

	// Now get down to work...  Load up all the config files
	if err := LoadFileInVolume(c, volname, "config.json", configContents); err != nil {
		return err
	}
	if err := LoadFileInVolume(c, volname, "ca.json", caContents); err != nil {
		return err
	}
	if err := LoadFileInVolume(c, volname, "server.json", serverContents); err != nil {
		return err
	}

	// And run the CFSSL commands to generate the bits and pieces
	// Normally you wouldn't want the ca to be in the same place, but it makes test setup much easier
	image := BuildImageString("ucp-cfssl")
	entrypoint := []string{"sh"}
	binds := []string{fmt.Sprintf("%s:/etc/cfssl", volname)}

	// Note: we have to make a copy of the ca without the chain so later when BYO install comes
	// along and glues the chain in there, cfssl doesn't barf when trying to create secondary certs
	// with the same "external root"
	cmd := []string{"-c", "cfssl genkey -initca ca.json | cfssljson -bare ca; cp ca.pem ca_no_chain.pem"}
	if out, err := runCommand(m, image, cmd, binds, entrypoint); err != nil {
		return err
	} else {
		fmt.Println(out)
	}

	cmd = []string{"-c", "cfssl gencert -config config.json -ca ca.pem -ca-key ca-key.pem -profile node server.json | cfssljson -bare cert; mv cert-key.pem key.pem"}
	if out, err := runCommand(m, image, cmd, binds, entrypoint); err != nil {
		return err
	} else {
		fmt.Println(out)
	}

	return nil
}

// Used for server cert regeneration (assumes the config is all wired up, just regen server
func BYOServerCertInitWithSameCA(m Machine) error {
	volname := "ucp-controller-server-certs"
	image := BuildImageString("ucp-cfssl")
	entrypoint := []string{"sh"}
	binds := []string{fmt.Sprintf("%s:/etc/cfssl", volname)}

	cmd := []string{"-c", "cfssl gencert -config config.json -ca ca_no_chain.pem -ca-key ca-key.pem -profile node server.json | cfssljson -bare cert; mv cert-key.pem key.pem"}
	if out, err := runCommand(m, image, cmd, binds, entrypoint); err != nil {
		return err
	} else {
		fmt.Println(out)
	}
	return nil
}

// Called to generate certificate material for secondary nodes signed by the first
// You must call BYOServerCertInit on the same machine first
func BYOServerCertGenSecondary(primary, secondary Machine) error {
	volname := "ucp-controller-server-certs"
	image := BuildImageString("ucp-cfssl")
	entrypoint := []string{"sh"}
	binds := []string{fmt.Sprintf("%s:/etc/cfssl", volname)}

	ip, _ := secondary.GetIP()
	serverContents := fmt.Sprintf(`
{
    "hosts": [
        "127.0.0.1",
        "%s"
    ],
    "key": {
        "algo": "rsa",
        "size": 4096
    },
    "CN": "controller"
}`, ip)
	primaryClient, err := primary.GetClient()
	if err != nil {
		return err
	}
	if err := LoadFileInVolume(primaryClient, volname, "secondary.json", serverContents); err != nil {
		return err
	}

	// Note: use the non-chained ca file so cfssl doesn't barf
	cmd := []string{"-c", "cfssl gencert -config config.json -ca ca_no_chain.pem -ca-key ca-key.pem -profile node secondary.json | cfssljson -bare secondary; cat ca.pem ; echo SPLIT; cat secondary.pem ; echo SPLIT ; cat secondary-key.pem"}
	out, err := runCommand(primary, image, cmd, binds, entrypoint)
	if err != nil {
		return err
	}
	//fmt.Println(out)
	results := strings.Split(string(out), "SPLIT")
	if len(results) != 3 {
		return fmt.Errorf("Malformed cert generated for secondary node")
	}
	ca := results[0]
	cert := results[1]
	key := results[2]

	secondaryClient, err := secondary.GetClient()
	if err != nil {
		return err
	}
	if err := LoadFileInVolume(secondaryClient, volname, "ca.pem", ca); err != nil {
		return err
	}
	if err := LoadFileInVolume(secondaryClient, volname, "cert.pem", cert); err != nil {
		return err
	}
	if err := LoadFileInVolume(secondaryClient, volname, "key.pem", key); err != nil {
		return err
	}
	return nil
}

func GetCAPEM(serverURL string) (string, error) {
	u, err := url.Parse(serverURL)
	if err != nil {
		return "", err
	}
	u.Path = "/ca"
	// It's just test code, so punt on the chain of trust...
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(u.String())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	caPEM, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(caPEM), nil
}

/*
BYO cert scenarios

1) Generate the machine
2) GetIP
3) Create the ucp-controller-server-certs volume
4) Use busybox to load ca.json into the volume
5) Use busybox to load server.json into the volume
6) Use busybox to load config.json into the volume



6) Spin up a dockerorca/cfssl container with the mount and run multiple commands
cfssl genkey -initca ca.json | cfssljson -bare ca
cfssl gencert -config=config.json -ca ca.pem -ca-key ca-key.pem -profile=node server.json | cfssljson -bare server



*/
