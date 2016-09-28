package util

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
	// "testing"

	"github.com/docker/dhe-deploy/bootstrap"
	"github.com/docker/dhe-deploy/bootstrap/flags"
	"github.com/docker/dhe-deploy/bootstrap/ucpclient"
	"github.com/docker/dhe-deploy/ha-integration/ha_utils"
	"github.com/docker/dhe-deploy/integration/apiclient"
	"github.com/docker/dhe-deploy/shared/dtrutil"
	dc "github.com/docker/engine-api/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	log "github.com/Sirupsen/logrus"
	// "github.com/codegangsta/cli"
	"github.com/samalba/dockerclient"
)

type DTRInstallation struct {
	Args             DTRArgs
	DTRLBNode        string
	PrimaryReplicaID string
	Replicas         map[string]string
	API              apiclient.APIClient
}

var DefaultInstallation DTRInstallation

func GetAvailableMachine(machines []ha_utils.Machine) ha_utils.Machine {
	for _, machine := range machines {
		if _, ok := DefaultInstallation.Replicas[machine.GetName()]; !ok {
			return machine
		}
	}

	return nil
}

func GetDockerClient(ucpAddr, username, password, hubUsername, hubPassword string) (*dc.Client, error) {
	// TODO: don't always use insecure tls?
	httpClient, err := dtrutil.HTTPClient(true)
	if err != nil {
		return nil, err
	}
	ucp := ucpclient.New(ucpAddr, httpClient)
	err = ucp.Login(username, password)
	if err != nil {
		return nil, err
	}
	client, err := bootstrap.DockerClientFromJWT(ucpAddr, ucp.JWT(), httpClient, hubUsername, hubPassword)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func SetDefaults(ucp, dtr, dtrLBNode string) error {
	DefaultInstallation.Args.DTRHost = dtr
	DefaultInstallation.DTRLBNode = dtrLBNode

	DefaultInstallation.API = MakeAPIClient(DefaultInstallation.Args.DTRHost)

	DefaultInstallation.Args.UCPHost = ucp

	return nil
}

func MakeAPIClient(DTRHost string) apiclient.APIClient {
	httpClient := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 10 * time.Second,
			}).Dial,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			TLSHandshakeTimeout:   5 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			MaxIdleConnsPerHost:   5,
		},
	}

	httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return apiclient.RedirectAttemptedError
	}

	return apiclient.New(DefaultInstallation.Args.DTRHost, 5, httpClient)
}

func GetUCPControllerIP(machines []ha_utils.Machine) (string, error) {
	return machines[0].GetIP()
}

type UCPArgs struct {
	AdminUsername string
	AdminPassword string
	UCPHost       string
	UCPCA         string
	UCPNode       string
	Constraints   string

	UCPInsecureTLS bool
	// NoUCP         bool

}

func (args *UCPArgs) AddDefaults() {
	if args.AdminUsername == "" {
		args.AdminUsername = DefaultInstallation.Args.AdminUsername
	}
	if args.AdminPassword == "" {
		args.AdminPassword = DefaultInstallation.Args.AdminPassword
	}
	if args.UCPHost == "" {
		args.UCPHost = DefaultInstallation.Args.UCPHost
	}
	if args.UCPCA == "" {
		args.UCPInsecureTLS = true
	}
}

func (args *UCPArgs) AsCommand() string {
	sshCommand := ""
	if args.AdminUsername != "" {
		sshCommand = fmt.Sprintf("%s--%s %s ", sshCommand, flags.UsernameFlag.Name, args.AdminUsername)
	}
	if args.AdminPassword != "" {
		sshCommand = fmt.Sprintf("%s--%s %s ", sshCommand, flags.PasswordFlag.Name, args.AdminPassword)
	}
	if args.UCPHost != "" {
		sshCommand = fmt.Sprintf("%s--%s %s ", sshCommand, flags.UCPHostFlag.Name, args.UCPHost)
	}
	if args.UCPNode != "" {
		sshCommand = fmt.Sprintf("%s--%s %s ", sshCommand, flags.UCPNodeFlag.Name, args.UCPNode)
	}

	// UCPCA

	if args.UCPInsecureTLS {
		sshCommand = fmt.Sprintf("%s--%s ", sshCommand, flags.UCPInsecureTLSFlag.Name)
	}

	return strings.TrimSuffix(sshCommand, " ")
}

type DTRReplicaArgs struct {
	ReplicaID         string
	ExistingReplicaID string
}

func (args *DTRReplicaArgs) AddDefaults() {
	if args.ExistingReplicaID == "" {
		args.ExistingReplicaID = DefaultInstallation.PrimaryReplicaID
	}
	// No ReplicaID default, it's either for specifying or identical in the case of remove.
}

func (args *DTRReplicaArgs) AsCommand() string {
	sshCommand := ""
	if args.ReplicaID != "" {
		sshCommand = fmt.Sprintf("%s--%s %s ", sshCommand, flags.ReplicaIDFlag.Name, args.ReplicaID)
	}
	if args.ExistingReplicaID != "" {
		sshCommand = fmt.Sprintf("%s--%s %s ", sshCommand, flags.ExistingReplicaIDFlag.Name, args.ExistingReplicaID)
	}

	return strings.TrimSuffix(sshCommand, " ")
}

type DTRImageArgs struct {
	DTRRepo           string
	DTRTag            string
	DTRBootstrapImage string
	DTRImages         []string

	SkipImages bool
	PullImages bool
}

func (args *DTRImageArgs) AddDefaults() {
	if args.DTRRepo == "" {
		args.DTRRepo = DefaultInstallation.Args.DTRRepo
	}
	if args.DTRTag == "" {
		args.DTRTag = DefaultInstallation.Args.DTRTag
	}
	if args.DTRBootstrapImage == "" {
		args.DTRBootstrapImage = DefaultInstallation.Args.DTRBootstrapImage
	}
	if len(args.DTRImages) == 0 {
		args.DTRImages = DefaultInstallation.Args.DTRImages
	}

	if DefaultInstallation.Args.SkipImages {
		args.SkipImages = true
	}
	if DefaultInstallation.Args.PullImages {
		args.PullImages = true
	}
}

type HubArgs struct {
	HubUsername string
	HubPassword string
}

type HTTPArgs struct {
	HTTPProxy        string
	HTTPSProxy       string
	NoProxy          bool
	ReplicaHTTPPort  string
	ReplicaHTTPSPort string
	DTRHost          string
}

func (args *HTTPArgs) AddDefaults() {
	if args.DTRHost == "" {
		args.DTRHost = DefaultInstallation.Args.DTRHost
	}
}

func (args *HTTPArgs) AsCommand() string {
	sshCommand := ""
	if args.DTRHost != "" {
		sshCommand = fmt.Sprintf("%s--%s %s ", sshCommand, flags.DTRHostFlag.Name, args.DTRHost)
	}
	if args.ReplicaHTTPSPort != "" {
		sshCommand = fmt.Sprintf("%s--%s %s ", sshCommand, flags.ReplicaHTTPSPortFlag.Name, args.ReplicaHTTPSPort)
	}

	return strings.TrimSuffix(sshCommand, " ")
}

type EnziArgs struct {
	EnziInsecureTLS bool
	EnziCA          string
}

type ForceArgs struct {
	Force bool
}

func (args *ForceArgs) AsCommand() string {
	sshCommand := ""
	if args.Force {
		sshCommand = fmt.Sprintf("%s--%s ", sshCommand, flags.ForceRemoveFlag.Name)
	}

	return strings.TrimSuffix(sshCommand, " ")
}

type ConstraintArgs struct {
	Constraints string
}

func (args *ConstraintArgs) AsCommand() string {
	sshCommand := ""
	if args.Constraints != "" {
		sshCommand = fmt.Sprintf("%s--%s '%s'", sshCommand, flags.ExtraEnvsFlag.Name, args.Constraints)
	} else {
		sshCommand = fmt.Sprintf("%s--%s $(hostname) ", sshCommand, flags.UCPNodeFlag.Name)
	}

	return strings.TrimSuffix(sshCommand, " ")
}

type PProfArgs struct {
	EnablePProf bool
}

func (args *PProfArgs) AsCommand() string {
	sshCommand := ""
	if args.EnablePProf {
		sshCommand = fmt.Sprintf("%s--%s ", sshCommand, flags.EnablePProfFlag.Name)
	}

	return strings.TrimSuffix(sshCommand, " ")
}

type DTRArgs struct {
	UCPArgs
	DTRReplicaArgs
	DTRImageArgs
	HubArgs
	HTTPArgs
	EnziArgs
}

type InstallArgs struct {
	UCPArgs
	DTRReplicaArgs
	// DTRImageArgs
	// HubArgs
	HTTPArgs
	EnziArgs
	ConstraintArgs
	PProfArgs
}

type JoinArgs struct {
	UCPArgs
	DTRReplicaArgs
	// DTRImageArgs
	// HubArgs
	HTTPArgs
	EnziArgs
	ConstraintArgs
	PProfArgs
}

type RemoveArgs struct {
	UCPArgs
	DTRReplicaArgs
	DTRImageArgs
	ForceArgs
}

type UpgradeArgs struct {
	UCPArgs
	DTRImageArgs
	DTRReplicaArgs
}

func PushDTRImages(machines []ha_utils.Machine, args DTRImageArgs) (err error) {
	(&args).AddDefaults()

	if args.SkipImages {
		return nil
	}

	// If Pulling the DTR images, use the provided DTR_REPO and DTR_TAG vars
	if args.PullImages {
		log.Info("Pulling DTR Images on all nodes")
		var errs []error
		wg := sync.WaitGroup{}
		wg.Add(len(machines))
		for _, machine := range machines {
			go func(machine ha_utils.Machine) {
				defer GinkgoRecover()
				log.Infof("pulling on %s", machine.GetName())
				remoteClient, _ := machine.GetClient()

				err := ha_utils.PullImages(remoteClient, args.DTRImages)
				if err != nil {
					errs = append(errs, err)
				}
				wg.Done()
			}(machine)
		}
		wg.Wait()
		Expect(errs).To(BeEmpty())
	} else {
		log.Info("Transferring DTR Images to all nodes")
		log.Info(args.DTRImages)

		localClient, err := dockerclient.NewDockerClient("unix://var/run/docker.sock", nil)
		Expect(err).To(BeNil())

		for _, machine := range machines {
			remoteClient, _ := machine.GetClient()
			err := ha_utils.TransferImages(localClient, remoteClient, args.DTRImages)
			Expect(err).To(BeNil())
		}
	}
	return
}

// XXX: This is aweful. TODO: Remove it.
func GetReplicaID(machine ha_utils.Machine) (string, error) {
	return machine.MachineSSH("sudo docker ps | grep 'dtr-etcd' | awk '{print $(NF)}' | cut -d '-' -f 3")
}

func InstallDTRNode(machine ha_utils.Machine, args InstallArgs) {
	(&args.UCPArgs).AddDefaults()
	(&args.HTTPArgs).AddDefaults()
	Expect(DefaultInstallation.PrimaryReplicaID).To(BeEmpty())

	sshArgument := fmt.Sprintf("sudo docker run --rm %s install --debug %s %s %s %s %s", DefaultInstallation.Args.DTRBootstrapImage, (&args.UCPArgs).AsCommand(), (&args.HTTPArgs).AsCommand(), (&args.ConstraintArgs).AsCommand(), (&args.PProfArgs).AsCommand(), (&args.DTRReplicaArgs).AsCommand())
	log.Info(machine.GetName())
	log.Info(sshArgument)
	out, err := machine.MachineSSH(sshArgument)
	log.Debug(out)
	Expect(err).To(BeNil())

	replicaID := args.DTRReplicaArgs.ReplicaID
	if replicaID == "" {
		replicaID, err = GetReplicaID(machine)
		Expect(err).To(BeNil())
	}

	Expect(replicaID).ToNot(BeEmpty(), "Got back an empty replica ID!")

	DefaultInstallation.PrimaryReplicaID = replicaID
	DefaultInstallation.Replicas[machine.GetName()] = replicaID
}

func JoinDTRNode(machine ha_utils.Machine, args JoinArgs) {
	(&args.UCPArgs).AddDefaults()
	(&args.DTRReplicaArgs).AddDefaults()

	sshArgument := fmt.Sprintf("sudo docker run --rm %s join --debug %s %s %s %s", DefaultInstallation.Args.DTRBootstrapImage, (&args.UCPArgs).AsCommand(), (&args.DTRReplicaArgs).AsCommand(), (&args.ConstraintArgs).AsCommand(), (&args.DTRReplicaArgs).AsCommand())
	log.Info(machine.GetName())
	log.Info(sshArgument)
	out, err := machine.MachineSSH(sshArgument)
	log.Debug(out)
	Expect(err).To(BeNil())

	replicaID := args.DTRReplicaArgs.ReplicaID
	if replicaID == "" {
		replicaID, err = GetReplicaID(machine)
		Expect(err).To(BeNil())
	}

	DefaultInstallation.Replicas[machine.GetName()] = replicaID
}

func RemoveDTRNode(machine ha_utils.Machine, args RemoveArgs) {
	(&args.UCPArgs).AddDefaults()
	(&args.DTRImageArgs).AddDefaults()

	sshArgument := fmt.Sprintf("sudo docker run --rm %s remove --debug %s %s %s", args.DTRImageArgs.DTRBootstrapImage, (&args.ForceArgs).AsCommand(), (&args.UCPArgs).AsCommand(), (&args.DTRReplicaArgs).AsCommand())
	log.Info(machine.GetName())
	log.Info(sshArgument)
	out, err := machine.MachineSSH(sshArgument)
	if err != nil {
		log.Debug(out)
	}
	Expect(err).To(BeNil())

	replicaID := DefaultInstallation.Replicas[machine.GetName()]
	delete(DefaultInstallation.Replicas, machine.GetName())
	if replicaID == DefaultInstallation.PrimaryReplicaID {
		if len(DefaultInstallation.Replicas) > 0 {
			for key := range DefaultInstallation.Replicas {
				DefaultInstallation.PrimaryReplicaID = DefaultInstallation.Replicas[key]
				break
			}
		} else {
			DefaultInstallation.PrimaryReplicaID = ""
		}
	}
}

// UpgradeDTR upgrades the entire DTR cluster
func UpgradeDTR(machines []ha_utils.Machine, args UpgradeArgs) {
	if len(DefaultInstallation.Replicas) == 0 {
		return
	}

	(&args.UCPArgs).AddDefaults()
	args.ExistingReplicaID = DefaultInstallation.PrimaryReplicaID
	machine := machines[0]

	sshArgument := fmt.Sprintf("sudo docker run --rm %s upgrade --debug %s %s", args.DTRBootstrapImage, (&args.UCPArgs).AsCommand(), (&args.DTRReplicaArgs).AsCommand())
	log.Info(machine.GetName())
	log.Info(sshArgument)
	out, err := machine.MachineSSH(sshArgument)
	if err != nil {
		log.Debug(out)
	}
	Expect(err).To(BeNil())
}

func RemoveDTR(machines []ha_utils.Machine, imageArgs DTRImageArgs) {
	(&imageArgs).AddDefaults()
	for _, machine := range machines {
		if _, ok := DefaultInstallation.Replicas[machine.GetName()]; ok {
			RemoveDTRNode(machine, RemoveArgs{
				DTRReplicaArgs: DTRReplicaArgs{
					ReplicaID:         DefaultInstallation.Replicas[machine.GetName()],
					ExistingReplicaID: DefaultInstallation.Replicas[machine.GetName()],
				},
				UCPArgs:      DefaultInstallation.Args.UCPArgs,
				DTRImageArgs: imageArgs,
				ForceArgs: ForceArgs{
					Force: true,
				},
			})
		}
	}
}

// deployDTR deploys an HA installation of DTR on a cluster
func DeployDTR(machines []ha_utils.Machine, args InstallArgs, replicas int) {
	installed := 0
	for _, machine := range machines {
		if installed == replicas {
			break
		}
		args := args
		args.ReplicaID = fmt.Sprintf("%012x", installed)
		if DefaultInstallation.PrimaryReplicaID == "" {
			InstallDTRNode(machine, args)
		} else {
			JoinDTRNode(machine, JoinArgs(args))
		}
		installed++
	}
	Expect(installed).To(Equal(replicas))

	log.Debugf("Deployed DTR on %d nodes", installed)
}

func PurgeDTR(machines []ha_utils.Machine) {
	log.Info("purging DTR on all machines")
	DefaultInstallation.PrimaryReplicaID = ""
	for _, machine := range machines {
		delete(DefaultInstallation.Replicas, machine.GetName())
		machine.MachineSSH("sudo docker ps -a | grep dtr- | grep -v enzi | awk '{print $1}' | xargs sudo docker rm -f; sudo docker volume ls | grep dtr | awk '{print $2}' | xargs sudo docker volume rm")
	}
}

// We consider a UCP controller node healthy for our setup if
// UCP itself can confirm node is healthy and part of the cluster
// Poll any UCP controller and check against list of healthy controller nodes
func UCPHealthCheck(ucpIP string, machines []ha_utils.Machine, numControllers int) error {
	retryInterval := time.Second * time.Duration(numControllers)
	retryAttempts := 300
	errorMsg := ""

	httpClient, err := dtrutil.HTTPClient(true)
	if err != nil {
		log.Errorf("Can't get http client: %s", err)
		return err
	}

	ucp := ucpclient.New(ucpIP, httpClient)
	err = dtrutil.Poll(retryInterval, retryAttempts, func() error {
		err = ucp.Login("admin", "orca")
		if err != nil {
			log.Debugf("Retrying login - Error: %s", err)
			return err
		}
		return nil
	})
	log.Infof("Num of controllers: %d", numControllers)
	if err != nil {
		return err
	}

	err = dtrutil.Poll(retryInterval, retryAttempts, func() error {
		log.Info("Pinging UCP to check controllers' health")
		nodes, err := ucp.GetNodes()
		if err != nil {
			log.Debugf("Can't get UCP's nodes: %s", err)
			return err
		}
		if len(nodes) != len(machines) {
			errorMsg = fmt.Sprintf("Cluster isn't of intended size yet")
			log.Debugf("%s", errorMsg)
			return fmt.Errorf("%s", errorMsg)
		}

		sort.Sort(ha_utils.UCPMachines(nodes))
		healthyCount := 0
		for i, node := range nodes {
			if i >= numControllers {
				break
			}
			if node.Status == "Healthy" {
				healthyCount++
				if healthyCount == numControllers {
					return nil
				}
			}
		}
		errorMsg = fmt.Sprintf("Not all controllers have joined yet")
		log.Debug("%s", errorMsg)
		return fmt.Errorf("%s", errorMsg)
	})
	return err
}
