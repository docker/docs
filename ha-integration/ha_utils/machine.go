package ha_utils

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/stdcopy"
	dc "github.com/docker/engine-api/client"
	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/samalba/dockerclient"
	"github.com/stretchr/testify/require"
)

var (
	Timeout            = 90 * time.Second
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
	GetNewClient() (*dc.Client, error)
	Remove() error
	Stop() error
	Start() error
	GetIP() (string, error)
	GetInternalIP() (string, error)
	CatHostFile(hostPath string) ([]byte, error)
	TarHostDir(hostPath string) ([]byte, error)
	MachineSSH(command string) (string, error)
}

// MachineSCP copies files from one Machine to another using docker-machine scp
func MachineSCP(source, target Machine, sourcePath, targetPath string) (string, error) {
	cmd := exec.Command("docker-machine", "scp", "-r",
		source.GetName()+":"+sourcePath, target.GetName()+":"+targetPath)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

// LocalMachineSCP copies files from the local host to another Machine using docker-maching scp
func LocalMachineSCP(target Machine, sourcePath, targetPath string) (string, error) {
	cmd := exec.Command("docker-machine", "scp", "-r",
		sourcePath, target.GetName()+":"+targetPath)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
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

	if _, err = stdcopy.StdCopy(stdoutBuffer, stderrBuffer, reader); err != nil {
		log.Debug("cannot read logs from logs reader")
		return nil, err
	}
	stderr := stderrBuffer.String()
	if strings.Contains(strings.ToLower(stderr), "no such file") {
		// XXX This doesn't seem to hit...
		log.Debug("Got a no such file on stderr")
		log.Debug(stderr)
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
		log.Debug("Non zero exit code: %d", info.State.ExitCode)
		return nil, PathDoesNotExist
	}

	// Looks like it worked OK
	return stdoutBuffer.Bytes(), nil
}

type MachineCreateFlags struct {
	MachineDriver      string
	CreateFlags        []string
	FixupCommand       string
	GenericMachineList []string
}

// DeployMachines deploys a number of docker-machine machines and installs UCP them
func CreateMachines(t require.TestingT, parallelCreate bool, numMachines int, prefix string, mcf MachineCreateFlags, offset int, instanceTypes []string) []Machine {
	if len(instanceTypes) > 0 {
		Expect(len(instanceTypes)).To(BeNumerically(">=", numMachines))
	}
	// this tells docker machine to create a bunch of machines in parallel
	machines := make([]Machine, numMachines)
	var wg sync.WaitGroup
	require := require.New(t)

	if parallelCreate {
		wg.Add(numMachines)
	}

	for i := 0; i < numMachines; i++ {
		if parallelCreate {
			go func(i int) {
				defer ginkgo.GinkgoRecover()
				defer wg.Done()
				machines[i] = createMachine(require, i, prefix, mcf, offset, instanceTypes)
			}(i)
		} else {
			machines[i] = createMachine(require, i, prefix, mcf, offset, instanceTypes)
		}
	}

	if parallelCreate {
		wg.Wait()
	}

	return machines
}

func createMachine(require *require.Assertions, i int, prefix string, mcf MachineCreateFlags, offset int, instanceTypes []string) Machine {
	mcfCopy := mcf
	if len(instanceTypes) > 0 {
		// this removes the instance type parameter if found, then adds the right instance type
		newFlags := []string{}
		skip := false
		for _, flag := range mcf.CreateFlags {
			if skip {
				skip = false
			} else if flag == "--amazonec2-instance-type" {
				skip = true
			} else {
				newFlags = append(newFlags, flag)
			}
		}
		newFlags = append(newFlags, "--amazonec2-instance-type", instanceTypes[i])
		mcfCopy.CreateFlags = newFlags
	}
	machine, err := NewBuildMachine(offset+i, prefix, mcfCopy)
	require.Nil(err)
	client, err := machine.GetClient()
	require.Nil(err)
	err = LoadAllLocalOrcaImages(client)
	require.Nil(err, "Failed to load images: %s", err)
	// writing to a fixed size array should be safe, right?

	return machine
}

// given a list of controller nodes, regenerate their certs
func RegenCerts(machines []Machine, san string) {
	log.Infof("regnerating certs for %d machines if necessary", len(machines))
	// this is a massively parallel action
	wg := sync.WaitGroup{}
	wg.Add(len(machines))
	errs := []string{}
	for _, machine := range machines {
		go func(machine Machine) {
			defer ginkgo.GinkgoRecover()
			defer wg.Done()

			// first check if the certs are already good because this is an expensive operation
			// and we can't just yolo do it every time
			IP, err := machine.GetIP()
			if err != nil {
				errs = append(errs, err.Error())
				return
			}
			conn, err := tls.Dial("tcp", fmt.Sprintf("%s:444", IP), &tls.Config{
				InsecureSkipVerify: true,
			})
			if err != nil {
				errs = append(errs, err.Error())
				return
			}
			certs := conn.ConnectionState().PeerCertificates
			conn.Close()
			cert := certs[0]
			err = cert.VerifyHostname(san)
			if err == nil {
				// if this cert validates for the new name, our work here is already done
				return
			} else {
				log.Debugf("regenerating cert because name verification failed with %s", err)
			}

			client, err := machine.GetClient()
			if err != nil {
				errs = append(errs, err.Error())
				return
			}

			id, err := GetID(machine)
			if err != nil {
				errs = append(errs, err.Error())
				return
			}

			// regenerate certs for this machine
			args := []string{"regen-certs", "-D", "--san", san, "--id", id}
			err = RunUCPBootstrapper(client, args, []string{})
			if err != nil {
				errs = append(errs, err.Error())
				return
			}

			_, err = machine.MachineSSH("sudo service docker restart")
			if err != nil {
				errs = append(errs, err.Error())
				return
			}
		}(machine)
	}
	wg.Wait()
	Expect(errs).To(BeEmpty())
}

func GetID(machine Machine) (string, error) {
	client, err := machine.GetClient()
	if err != nil {
		return "", err
	}

	args := []string{"id"}
	output, err := RunUCPBootstrapperWithOutput(client, args, []string{})
	if err != nil {
		return "", err
	}
	lines := strings.Split(output, "\n")
	secondLastLine := lines[len(lines)-2]
	return secondLastLine, nil
}

func BackupCA(machine Machine) (string, error) {
	client, err := machine.GetNewClient()
	if err != nil {
		return "", err
	}
	id, err := GetID(machine)
	if err != nil {
		return "", err
	}

	// regenerate certs for this machine
	args := []string{"backup", "-D", "--root-ca-only", "--id", id}
	out, stderr, err := RunUCPBootstrapperWithInput(client, args, []string{}, "")
	if err != nil {
		log.Debugf("output: %s; %s", out, stderr)
	}
	return out, err
}

func RestoreCA(machine Machine, str string) error {
	client, err := machine.GetNewClient()
	if err != nil {
		return err
	}
	id, err := GetID(machine)
	if err != nil {
		return err
	}

	// regenerate certs for this machine
	args := []string{"restore", "-D", "--root-ca-only", "--id", id}
	_, _, err = RunUCPBootstrapperWithInput(client, args, []string{}, str)
	return err
}

func DeployUCPMachines(t require.TestingT, machines []Machine, numControllers int, replicateCAs bool) {
	var err error
	require := require.New(t)

	// Handle the creation of the very first machine
	firstControllerMachine := machines[0]

	firstControllerIP, err := firstControllerMachine.GetIP()
	require.Nil(err)
	client, err := firstControllerMachine.GetClient()
	require.Nil(err)
	// Install UCP controller on port 444
	args := []string{"install", "--disable-tracking", "--disable-usage", "--controller-port", "444",
		"-D", "--swarm-port", "3376", "--fresh-install", "--san", firstControllerIP}
	require.Nil(RunUCPBootstrapper(client, args, []string{}))

	firstControllerURL, err := GetOrcaURL(firstControllerMachine)
	require.Nil(err)

	// Join requires a valid license
	AddValidLicense(t, firstControllerURL, GetAdminUser(), GetAdminPassword())
	log.Debug("First machine ready and running UCP")

	backup := ""
	if replicateCAs {
		backup, err = BackupCA(firstControllerMachine)
		if err != nil {
			log.Debug(backup)
		}
		require.Nil(err)
		require.NotEmpty(backup)
	}

	JoinUCPMachines(t, machines[0], machines[1:], numControllers, backup)

	log.Debugf("UCP cluster of %d nodes initialized and configured", len(machines))
}

func JoinUCPMachines(t require.TestingT, controller Machine, machines []Machine, numControllers int, backup string) {
	require := require.New(t)
	firstControllerURL, err := GetOrcaURL(controller)
	require.Nil(err)
	fingerprint, err := GetOrcaFingerprint(controller)
	require.Nil(err)

	// Join the controllers sequentially
	for i, machine := range machines {
		// Currently only the first numControllers machines in order are going to be controllers.
		if i >= (numControllers - 1) {
			break
		}

		IP, err := machine.GetIP()
		require.Nil(err)

		args := []string{"join", "-D", "--controller-port", "444", "--fingerprint", fingerprint, "--admin-username", "admin", "--admin-password", "orca", "--url", firstControllerURL, "--swarm-port", "3376", "--fresh-install", "--san", IP}
		args = append(args, "--replica")

		client, err := machine.GetClient()
		require.Nil(err)

		log.Debugf("Adding host %d to: %s %s", i, firstControllerURL, fingerprint)
		require.Nil(RunUCPBootstrapper(client, args, []string{"UCP_ADMIN_USER=admin", "UCP_ADMIN_PASSWORD=orca"}))

		if backup != "" {
			err = RestoreCA(machine, backup)
			require.Nil(err)
		}
	}

	// Join the non-controllers in parallel
	errChan := make(chan error)
	for i, machine := range machines {
		if i < (numControllers - 1) {
			continue
		}

		IP, err := machine.GetIP()
		require.Nil(err)

		args := []string{"join", "-D", "--controller-port", "444", "--fingerprint", fingerprint, "--admin-username", "admin", "--admin-password", "orca", "--url", firstControllerURL, "--swarm-port", "3376", "--fresh-install", "--san", IP}

		go func(machine Machine, args []string) {
			client, err := machine.GetClient()
			if err != nil {
				errChan <- err
				return
			}

			log.Debugf("Adding host %d to: %s %s", i, firstControllerURL, fingerprint)
			errChan <- RunUCPBootstrapper(client, args, []string{"UCP_ADMIN_USER=admin", "UCP_ADMIN_PASSWORD=orca"})
		}(machine, args)
	}

	for i := 0; i < len(machines)-numControllers; i++ {
		require.Nil(<-errChan)
	}
}
