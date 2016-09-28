package ha_utils

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strings"

	log "github.com/Sirupsen/logrus"
	dc "github.com/docker/engine-api/client"
	"github.com/docker/orca"
	"github.com/samalba/dockerclient"
)

type BuildMachine struct {
	name       string
	dockerHost string
	tlsConfig  *tls.Config
	sshKeyPath string
	sshUser    string
}

type DockerMachineInspectDriver struct {
	IPAddress  string
	SSHUser    string
	SSHKeyPath string
}

type DockerMachineInspect struct {
	Driver DockerMachineInspectDriver
}

// Generate a new machine using docker-machine CLI
func NewBuildMachine(machineID int, prefix string, mcf MachineCreateFlags) (Machine, error) {
	// Some cloud providers can be a little flaky, so try a few times before we give up
	for i := 0; i < RetryCount; i++ {
		machine, err := buildMachineOnce(machineID, prefix, mcf)
		if err != nil {
			machine.Remove()
		} else {
			return machine, err
		}
	}
	return nil, fmt.Errorf("Failed to create machine after %d tries", RetryCount)
}

func buildMachineOnce(machineID int, prefix string, mcf MachineCreateFlags) (Machine, error) {
	_, err := exec.LookPath("docker-machine")
	if err != nil {
		return nil, fmt.Errorf("You must have docker-machine installed locally in your path!")
	}

	args := []string{
		"create",
		"--driver",
		mcf.MachineDriver,
	}
	args = append(args, mcf.CreateFlags...)

	if mcf.MachineDriver == "generic" {
		args = append(args, "--generic-ip-address", mcf.GenericMachineList[machineID])
	}

	// The first machine is the UCP controller while the rest are DTR nodes each of which ultimatley have one replica
	m := &BuildMachine{}
	m.name = fmt.Sprintf("%s-%04d", prefix, machineID)

	args = append(args, m.name)

	log.Infof("Creating new test VM: %s", m.name)
	cmd := exec.Command("docker-machine", args...)
	cmd.Stdout = LogFile
	cmd.Stderr = LogFile
	cmd.Run()
	if err != nil {
		return nil, err
	}
	log.Infof("Created new test VM: %s", m.name)

	cmd = exec.Command("docker-machine", "inspect", m.name)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Error(err)
		log.Error(string(out))
		return nil, err
	}
	machineInfo := DockerMachineInspect{}
	if err := json.Unmarshal([]byte(out), &machineInfo); err != nil {
		log.Error(err)
		return nil, err
	}

	if mcf.FixupCommand != "" {
		log.Infof("Fixing test VM by running: %s %s", m.name, mcf.FixupCommand)

		// Can't use `docker-machine ssh` because it doesn't allocate a tty
		// which is needed for sudo.
		cmd = exec.Command("ssh", "-tt", "-l", machineInfo.Driver.SSHUser, "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", "-i", machineInfo.Driver.SSHKeyPath, machineInfo.Driver.IPAddress, mcf.FixupCommand)
		out, err = cmd.CombinedOutput()
		if err != nil {
			log.Error(err)
			log.Error(string(out))
			// If something went wrong, make sure to clean up after ourselves
			return nil, err
		}
	}

	// Now get the env settings for this new VM
	cmd = exec.Command("docker-machine", formatDockerMachineEnv(m.name)...)
	out, err = cmd.CombinedOutput()
	if err != nil {
		log.Error(err)
		log.Error(string(out))
		// If something went wrong, make sure to clean up after ourselves
		return nil, err
	} else {
		for _, line := range strings.Split(string(out), "\n") {
			if strings.Contains(line, "export DOCKER_HOST=") {
				vals := strings.Split(line, "=")
				m.dockerHost = strings.Trim(vals[1], `"`)
			} else if strings.Contains(line, "export DOCKER_CERT_PATH=") {
				vals := strings.Split(line, "=")
				certDir := strings.Trim(vals[1], `"`)
				log.Infof("Loading certs from %s", certDir)
				cert, err := tls.LoadX509KeyPair(
					fmt.Sprintf("%s/cert.pem", certDir),
					fmt.Sprintf("%s/key.pem", certDir))
				if err != nil {
					// If something went wrong, make sure to clean up after ourselves
					return nil, err
				}
				caCert, err := ioutil.ReadFile(
					fmt.Sprintf("%s/ca.pem", certDir))
				if err != nil {
					// If something went wrong, make sure to clean up after ourselves
					return nil, err
				}
				caCertPool := x509.NewCertPool()
				caCertPool.AppendCertsFromPEM(caCert)
				m.tlsConfig = &tls.Config{
					Certificates: []tls.Certificate{cert},
					RootCAs:      caCertPool,
				}
			}
		}
	}

	log.Infof("Host: %s", m.dockerHost)
	return m, nil
}

// These are so that we can sort the machines after retrieving them in parallel
type MachinesByName []Machine

func (machines MachinesByName) Len() int {
	return len(machines)
}

func (machines MachinesByName) Swap(i, j int) {
	machines[i], machines[j] = machines[j], machines[i]
}

func (machines MachinesByName) Less(i, j int) bool {
	return strings.Compare(machines[i].GetName(), machines[j].GetName()) == -1
}

// So we can sort the UCP machines after retrieving them from a controller
// and comparing them with the local list of sorted machines
type UCPMachines []orca.Node

func (machines UCPMachines) Len() int {
	return len(machines)
}

func (machines UCPMachines) Swap(i, j int) {
	machines[i], machines[j] = machines[j], machines[i]
}

func (machines UCPMachines) Less(i, j int) bool {
	return strings.Compare(machines[i].Name, machines[j].Name) == -1
}

func RetrieveClusterMachines(prefix string) (machines []Machine, err error) {
	cmd := exec.Command("docker-machine", "ls")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Error(err)
		log.Error(string(out))
		return nil, fmt.Errorf("failed to ls: %s", err)
	}

	type result struct {
		machine Machine
		err     error
	}

	resultChan := make(chan result)
	lines := strings.Split(string(out), "\n")
	numMachines := 0

	for _, line := range lines {
		vals := strings.Fields(line)
		if len(vals) == 0 || !strings.HasPrefix(vals[0], prefix) {
			continue
		}

		numMachines++
		go func(line string) {
			vals := strings.Fields(line)
			machineName := vals[0]
			machineDockerHost, err := GetDockerHost(machineName)
			if err != nil {
				resultChan <- result{nil, fmt.Errorf("failed to get docker host for %s: %s", machineName, err)}
				return
			}

			machineTLSConfig, err := GetTLSConfig(machineName)
			if err != nil {
				resultChan <- result{nil, err}
				return
			}

			resultChan <- result{&BuildMachine{
				name:       machineName,
				dockerHost: machineDockerHost,
				tlsConfig:  machineTLSConfig,
			}, nil}
		}(line)
	}

	errs := []error{}
	machines = []Machine{}
	for i := 0; i < numMachines; i++ {
		result := <-resultChan
		if result.err != nil {
			errs = append(errs, result.err)
		} else {
			machines = append(machines, result.machine)
		}
	}

	if len(errs) > 0 {
		return nil, fmt.Errorf("%s", errs)
	}

	sort.Sort(MachinesByName(machines))
	return machines, nil
}

func (m *BuildMachine) GetName() string {
	return m.name
}

func (m *BuildMachine) GetDockerHost() string {
	return m.dockerHost
}

func GetDockerHost(machineName string) (string, error) {
	cmd := exec.Command("docker-machine", formatDockerMachineEnv(machineName)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Error(err)
		log.Error(string(out))
		// If something went wrong, make sure to clean up after ourselves
		return "", err
	} else {
		for _, line := range strings.Split(string(out), "\n") {
			if strings.Contains(line, "export DOCKER_HOST=") {
				vals := strings.Split(line, "=")
				return strings.Trim(vals[1], `"`), nil
			}
		}
	}
	return "", fmt.Errorf("No docker host found in: %s", out)
}

// Return a docker client connected to the machine
func (m *BuildMachine) GetClient() (*dockerclient.DockerClient, error) {
	return dockerclient.NewDockerClientTimeout(m.dockerHost, m.tlsConfig, Timeout, nil)
}

func GetClient(machineName string) (*dockerclient.DockerClient, error) {
	dockerHost, err := GetDockerHost(machineName)
	if err != nil {
		return nil, err
	}
	tlsConfig, err := GetTLSConfig(machineName)
	if err != nil {
		return nil, err
	}
	return dockerclient.NewDockerClientTimeout(dockerHost, tlsConfig, Timeout, nil)
}

func (m *BuildMachine) GetNewClient() (*dc.Client, error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: m.tlsConfig,
		},
	}

	headers := map[string]string{}
	//if hubUsername != "" && hubPassword != "" {
	//	headers["X-Registry-Auth"] = MakeRegistryAuth(hubUsername, hubPassword)
	//}
	return dc.NewClient(m.dockerHost, "v1.22", client, headers)
}

func GetNewClient(machineName string) (*dc.Client, error) {
	dockerHost, err := GetDockerHost(machineName)
	if err != nil {
		return nil, err
	}
	tlsConfig, err := GetTLSConfig(machineName)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	headers := map[string]string{}
	//if hubUsername != "" && hubPassword != "" {
	//	headers["X-Registry-Auth"] = MakeRegistryAuth(hubUsername, hubPassword)
	//}
	return dc.NewClient(dockerHost, "v1.22", client, headers)
}

func formatDockerMachineEnv(name string) []string {
	if os.Getenv("NEW_DOCKER_MACHINE") != "" {
		return []string{"env", "--shell", "bash", name}
	} else {
		return []string{"env", name}
	}
}

func GetTLSConfig(machineName string) (*tls.Config, error) {
	cmd := exec.Command("docker-machine", formatDockerMachineEnv(machineName)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Error(err)
		log.Error(string(out))
		return nil, err
	} else {
		for _, line := range strings.Split(string(out), "\n") {
			if strings.Contains(line, "export DOCKER_CERT_PATH=") {
				vals := strings.Split(line, "=")
				certDir := strings.Trim(vals[1], `"`)
				log.Infof("Loading certs from %s", certDir)
				cert, err := tls.LoadX509KeyPair(
					fmt.Sprintf("%s/cert.pem", certDir),
					fmt.Sprintf("%s/key.pem", certDir))
				if err != nil {
					return nil, err
				}
				caCert, err := ioutil.ReadFile(
					fmt.Sprintf("%s/ca.pem", certDir))
				if err != nil {
					return nil, err
				}
				caCertPool := x509.NewCertPool()
				caCertPool.AppendCertsFromPEM(caCert)
				return &tls.Config{
					Certificates: []tls.Certificate{cert},
					RootCAs:      caCertPool,
				}, nil
			}
		}
	}
	return nil, fmt.Errorf("Couldn't get TLSConfig")
}

// Remove the machine after the tests have completed
func (m *BuildMachine) Remove() error {
	if os.Getenv("PRESERVE_TEST_MACHINE") != "" {
		log.Info("Skipping removal of machine with PRESERVE_TEST_MACHINE set")
		return nil
	}
	if m.name == "" {
		log.Info("Machine already deleted")
		return nil
	}
	cmd := exec.Command("docker-machine", "rm", "-f", m.name)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Error(string(out))
		return err
	}
	log.Infof("Machine %s deleted", m.name)
	m.name = ""
	return nil
}

func Remove(machineName string) error {
	if os.Getenv("PRESERVE_TEST_MACHINE") != "" {
		log.Info("Skipping removal of machine with PRESERVE_TEST_MACHINE set")
		return nil
	}
	if machineName == "" {
		log.Info("Machine already deleted")
		return nil
	}
	cmd := exec.Command("docker-machine", "rm", "-f", machineName)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Error(string(out))
		return err
	}
	log.Infof("Machine %s deleted", machineName)
	return nil
}

func (m *BuildMachine) Stop() error {
	cmd := exec.Command("docker-machine", "stop", m.name)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Error(string(out))
		return err
	}
	return nil
}

func Stop(machineName string) error {
	cmd := exec.Command("docker-machine", "stop", machineName)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Error(string(out))
		return err
	}
	return nil
}

func (m *BuildMachine) Start() error {
	cmd := exec.Command("docker-machine", "start", m.name)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Error(string(out))
		return err
	}
	return nil
}

func Start(machineName string) error {
	cmd := exec.Command("docker-machine", "start", machineName)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Error(string(out))
		return err
	}
	return nil
}

// Return the public IP of the machine
func (m *BuildMachine) GetIP() (string, error) {
	if os.Getenv("USE_PRIVATE_IP") == "" {
		cmd := exec.Command("docker-machine", "ip", m.name)
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Error(string(out))
			return "", err
		} else {
			return strings.TrimSpace(string(out)), nil
		}
	} else {
		return m.GetInternalIP()
	}
}

func GetIP(machineName string) (string, error) {
	if os.Getenv("USE_PRIVATE_IP") == "" {
		cmd := exec.Command("docker-machine", "ip", machineName)
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Error(string(out))
			return "", err
		} else {
			return strings.TrimSpace(string(out)), nil
		}
	} else {
		return GetInternalIP(machineName)
	}
}

// Get the internal IP (useful for join operations)
func (m *BuildMachine) GetInternalIP() (string, error) {
	cmd := exec.Command("docker-machine", "ssh", m.name, "ip route get 8.8.8.8| cut -d' ' -f8|head -1")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Error(string(out))
		return "", err
	} else {
		return strings.TrimSpace(string(out)), nil
	}
}

func GetInternalIP(machineName string) (string, error) {
	cmd := exec.Command("docker-machine", "ssh", machineName, "ip route get 8.8.8.8| cut -d' ' -f8|head -1")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Error(string(out))
		return "", err
	} else {
		return strings.TrimSpace(string(out)), nil
	}
}

func (m *BuildMachine) MachineSSH(command string) (string, error) {
	cmd := exec.Command("docker-machine", "ssh", m.name, command)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func MachineSSH(machineName string, command string) (string, error) {
	cmd := exec.Command("docker-machine", "ssh", machineName, command)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

// Get the contents of a specific file on the engine
func (m *BuildMachine) CatHostFile(hostPath string) ([]byte, error) {
	return catHostFile(m, hostPath)
}

// Get the content of a directory as a tar file from the engine
func (m *BuildMachine) TarHostDir(hostPath string) ([]byte, error) {
	return tarHostDir(m, hostPath)
}
