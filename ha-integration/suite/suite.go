package suite

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	dc "github.com/docker/engine-api/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/docker/dhe-deploy/ha-integration/ha_utils"
	"github.com/docker/dhe-deploy/ha-integration/util"
	"github.com/docker/dhe-deploy/shared/dtrutil"
)

type Args struct {
	machinesPrefix                string
	numControllers                int
	numMachines                   int
	numPrometheusExporters        int
	machineParallelCreate         bool
	machineDriver                 string
	createFlags                   []string
	fixupCommand                  string
	genericMachineList            []string
	dtrLBConstraints              []string
	ucpLBConstraints              []string
	instanceTypes                 []string
	prometheusServerConstraints   []string
	prometheusExporterConstraints []string
	hubUsername                   string
	hubPassword                   string
	enableUCPLB                   bool
	enablePrometheus              bool
	forcePurge                    bool
	forcePurgePrometheus          bool
}

func parseArgs() Args {
	args := Args{}

	args.machinesPrefix = os.Getenv("MACHINE_PREFIX") + "-DTRTest"

	numControllers64, err := strconv.ParseInt(os.Getenv("UCP_CONTROLLERS"), 10, 64)
	Expect(err).To(BeNil())
	args.numControllers = int(numControllers64)

	numMachines64, err := strconv.ParseInt(os.Getenv("NUM_MACHINES"), 10, 64)
	Expect(err).To(BeNil())
	args.numMachines = int(numMachines64)

	numPrometheusExporters64, err := strconv.ParseInt(os.Getenv("NUM_PROMETHEUS_EXPORTERS"), 10, 64)
	Expect(err).To(BeNil())
	args.numPrometheusExporters = int(numPrometheusExporters64)

	args.machineDriver = os.Getenv("MACHINE_DRIVER")
	Expect(args.machineDriver).To(Not(BeEmpty()))
	createFlagsStr := os.Getenv("MACHINE_CREATE_FLAGS")
	args.createFlags = strings.Split(createFlagsStr, " ")
	args.fixupCommand = os.Getenv("MACHINE_FIXUP_COMMAND")
	genericMachineListStr := os.Getenv("GENERIC_MACHINE_LIST")
	args.genericMachineList = strings.Split(genericMachineListStr, ",")
	for i, machine := range args.genericMachineList {
		args.genericMachineList[i] = strings.Trim(machine, " ")
	}
	if len(args.genericMachineList) == 1 && args.genericMachineList[0] == "" {
		args.genericMachineList = []string{}
	}

	instanceTypesStr := os.Getenv("INSTANCE_TYPES")
	args.instanceTypes = strings.Split(instanceTypesStr, ",")
	for i, machine := range args.instanceTypes {
		args.instanceTypes[i] = strings.Trim(machine, " ")
	}
	if len(args.instanceTypes) == 1 && args.instanceTypes[0] == "" {
		args.instanceTypes = []string{}
	}

	args.dtrLBConstraints = strings.Split(os.Getenv("DTR_LB_CONSTRAINTS"), "|")
	if len(args.dtrLBConstraints) == 1 && args.dtrLBConstraints[0] == "" {
		args.dtrLBConstraints = []string{}
	}
	args.ucpLBConstraints = strings.Split(os.Getenv("UCP_LB_CONSTRAINTS"), "|")
	if len(args.ucpLBConstraints) == 1 && args.ucpLBConstraints[0] == "" {
		args.ucpLBConstraints = []string{}
	}
	args.prometheusServerConstraints = strings.Split(os.Getenv("PROMETHEUS_SERVER_CONSTRAINTS"), "|")
	if len(args.prometheusServerConstraints) == 1 && args.prometheusServerConstraints[0] == "" {
		args.prometheusServerConstraints = []string{}
	}
	args.prometheusExporterConstraints = strings.Split(os.Getenv("PROMETHEUS_EXPORTER_CONSTRAINTS"), "|")
	if len(args.prometheusExporterConstraints) == 1 && args.prometheusExporterConstraints[0] == "" {
		args.prometheusExporterConstraints = []string{}
	}

	args.machineParallelCreate = os.Getenv("MACHINE_PARALLEL_CREATE") != ""

	args.hubUsername = os.Getenv("REGISTRY_USERNAME")
	args.hubPassword = os.Getenv("REGISTRY_PASSWORD")

	args.enableUCPLB = os.Getenv("ENABLE_UCP_LB") != ""
	args.enablePrometheus = os.Getenv("ENABLE_PROMETHEUS") != ""
	args.forcePurge = os.Getenv("FORCE_PURGE") != ""
	args.forcePurgePrometheus = os.Getenv("FORCE_PURGE_PROMETHEUS") != ""

	return args
}

func Prepare() []ha_utils.Machine {
	// TODO: get rid of the t
	t := GinkgoT()

	args := parseArgs()

	machinesPrefix := args.machinesPrefix
	numControllers := args.numControllers
	numMachines := args.numMachines
	numPrometheusExporters := args.numPrometheusExporters
	machineParallelCreate := args.machineParallelCreate
	machineDriver := args.machineDriver
	createFlags := args.createFlags
	fixupCommand := args.fixupCommand
	genericMachineList := args.genericMachineList
	dtrLBConstraints := args.dtrLBConstraints
	ucpLBConstraints := args.ucpLBConstraints
	instanceTypes := args.instanceTypes
	prometheusServerConstraints := args.prometheusServerConstraints
	prometheusExporterConstraints := args.prometheusExporterConstraints
	hubUsername := args.hubUsername
	hubPassword := args.hubPassword
	enableUCPLB := args.enableUCPLB
	enablePrometheus := args.enablePrometheus
	forcePurge := args.forcePurge
	forcePurgePrometheus := args.forcePurgePrometheus

	// this discovery existing docker machine machines to use
	machines, err := ha_utils.RetrieveClusterMachines(machinesPrefix)
	Expect(err).To(BeNil())

	if len(machines) == 0 {
		log.Info("no machines found, creating some.")
		machines = ha_utils.CreateMachines(t, machineParallelCreate, numMachines, machinesPrefix, ha_utils.MachineCreateFlags{
			MachineDriver:      machineDriver,
			CreateFlags:        createFlags,
			FixupCommand:       fixupCommand,
			GenericMachineList: genericMachineList,
		}, 0, instanceTypes)
		ha_utils.DeployUCPMachines(t, machines, numControllers, enableUCPLB)
	} else {
		// Here we assume that machines don't randomly disappear and the only reason
		// we might have fewer machines than expected is because we want to add more
		// at the end
		// In the future we can expand this to fill in gaps in the numbers as well
		if int(numMachines) > len(machines) {
			difference := int(numMachines) - len(machines)
			flags := ha_utils.MachineCreateFlags{
				MachineDriver: machineDriver,
				CreateFlags:   createFlags,
				FixupCommand:  fixupCommand,
			}
			if machineDriver == "generic" {
				Expect(len(genericMachineList) > difference).To(BeTrue())
			}
			if len(genericMachineList) > difference {
				flags.GenericMachineList = genericMachineList[difference:]
			}
			types := []string{}
			if len(instanceTypes) >= len(machines) {
				types = instanceTypes[len(machines):]
			}
			newMachines := ha_utils.CreateMachines(t, machineParallelCreate, difference, machinesPrefix, flags, len(machines), types)

			// the first X machines are controllers, so we calculate how many controllers are left to deploy:
			// XXX: in the future we should consider a more dynamic way to specify which subset should be controllers
			newControllers := numControllers - len(machines)
			if newControllers < 0 {
				newControllers = 0
			}

			backup := ""
			// we replicate the root ca only if ucp is using a load balancer
			if enableUCPLB {
				if newControllers > 0 {
					backup, err = ha_utils.BackupCA(machines[0])
					if err != nil {
						log.Debug(backup)
					}
					Expect(err).To(BeNil())
				}
			}
			// install ucp on the new machines to add them to the cluster
			ha_utils.JoinUCPMachines(t, machines[0], newMachines, newControllers, backup)

			machines = append(machines, newMachines...)
		}

		Expect(len(machines)).To(Equal(int(numMachines)))
		log.Info("machines found, using them.")

		if forcePurge {
			wg := sync.WaitGroup{}
			wg.Add(len(machines))
			errs := []error{}
			for i := range machines {
				go func(i int) {
					log.Infof("purging machine %d", i)
					_, err := machines[i].MachineSSH("sudo docker ps -aq | xargs sudo docker rm -f; sudo docker volume ls -q | xargs sudo docker volume rm; sudo rm /etc/docker/daemon.json; sudo service docker restart")
					if err != nil {
						errs = append(errs, err)
					}
					wg.Done()
					log.Infof("purged machine %d", i)
				}(i)
			}
			Expect(errs).To(BeEmpty(), "errors purging")
			wg.Wait()
			ha_utils.DeployUCPMachines(t, machines, numControllers, enableUCPLB)
		}
	}
	log.Infof("Using %d docker machines: %v", len(machines), machines)

	ucpIP := ""
	ucp := ""
	if enableUCPLB {
		// to bootstrap the ucp load balancer we have to use ucp without a load balancer
		IP, err := machines[0].GetIP()
		Expect(err).To(BeNil())

		client, err := util.GetDockerClient(fmt.Sprintf("%s:444", IP), "admin", "orca", hubUsername, hubPassword)
		Expect(err).To(BeNil())

		// It's not mandatory to use the ucp load balancer, and in fact, the regular tests don't use it
		ucpLBNode := util.SetupDefaultLoadBalancer(machines, client, ucpLBConstraints, 446, 444, util.UCPLoadBalancerContainerName)

		ucpIP, err = ha_utils.GetIP(ucpLBNode)
		Expect(err).To(BeNil())

		ucp = fmt.Sprintf("%s:446", ucpIP)
	} else {
		ucpIP, err = machines[0].GetIP()
		Expect(err).To(BeNil())
		ucp = fmt.Sprintf("%s:444", ucpIP)
	}

	ha_utils.RegenCerts(machines[:numControllers], ucpIP)
	log.Infof("The ucp load balancer is at: %s", ucp)

	// Perform health checks to make sure UCP controllers are all up and joined to
	// form the cluster of intended size. Blocks until checks for all controllers are done.
	err = util.UCPHealthCheck(ucp, machines, numControllers)
	Expect(err).To(BeNil())

	// from now on, use the ucp load balancer to do ucp things (if one was created)
	var client *dc.Client
	if !enableUCPLB {
		// if no load balancer is being used, try to log in right away
		client, err = util.GetDockerClient(ucp, "admin", "orca", hubUsername, hubPassword)
		Expect(err).To(BeNil())
	} else {
		attempts := 0
		// otherwise, wait forever for ucp to come back up. Thanks, Obama.
		err = dtrutil.Poll(time.Second, 300, func() error {
			client, err = util.GetDockerClient(ucp, "admin", "orca", hubUsername, hubPassword)
			attempts++
			if err != nil {
				return err
			}
			return nil
		})
		Expect(err).To(BeNil())
		log.Infof("ucp came up after %d attempts", attempts)
	}

	purged, prometheusNode := util.MaybePurgePrometheus(client, forcePurgePrometheus, numPrometheusExporters)
	if enablePrometheus && purged {
		// if no constraints are given, put prometheus on the first node
		log.Infof("deploying prometheus with %d constraints: %v", len(prometheusServerConstraints), prometheusServerConstraints)
		if len(prometheusServerConstraints) == 0 {
			nodeName := machines[0].GetName()
			log.Infof("No constraints given. Putting prometheus server on %s.", nodeName)
			prometheusServerConstraints = []string{fmt.Sprintf("constraint:node==%s", nodeName)}
		}

		prometheusNode = util.DeployPrometheus(client, prometheusServerConstraints, prometheusExporterConstraints, numPrometheusExporters)

	}
	if prometheusNode != "" {
		prometheusIP, err := ha_utils.GetIP(prometheusNode)
		Expect(err).To(BeNil())
		log.Infof("Prometheus is at: %s:666", prometheusIP)
	}

	// Note: empty constraints means put it on the first machine
	dtrLBNode := util.SetupDefaultLoadBalancer(
		machines, client, dtrLBConstraints, util.DTRLoadBalancerPort, util.DefaultDTRNodePort, util.DTRLoadBalancerContainerName)

	dtrIP, err := ha_utils.GetIP(dtrLBNode)
	Expect(err).To(BeNil())
	log.Infof("The dtr load balancer is at: %s:%d", dtrIP, util.DTRLoadBalancerPort)

	util.SetDefaults(ucp, fmt.Sprintf("%s:%d", dtrIP, util.DTRLoadBalancerPort), dtrLBNode)
	return machines
}
