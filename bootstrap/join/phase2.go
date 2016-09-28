package join

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/coreos/etcd/client"
	"github.com/docker/dhe-deploy"
	"github.com/docker/dhe-deploy/bootstrap"
	"github.com/docker/dhe-deploy/bootstrap/dumpcerts"
	"github.com/docker/dhe-deploy/bootstrap/flags"
	"github.com/docker/dhe-deploy/bootstrap/install"
	"github.com/docker/dhe-deploy/bootstrap/reconfigure"
	"github.com/docker/dhe-deploy/hubconfig/etcd"
	"github.com/docker/dhe-deploy/hubconfig/sanitizers"
	"github.com/docker/dhe-deploy/hubconfig/settingsstore"
	jschema "github.com/docker/dhe-deploy/pkg/jobrunner-framework/schema"
	"github.com/docker/dhe-deploy/shared/containers"
	"github.com/docker/dhe-deploy/shared/dtrutil"
	"github.com/docker/engine-api/types"
)

func phase2(c *cli.Context) (int, error) {
	log.Debug("phase 2 starting...")
	var err error

	// this will not do any prompts if it's run correctly
	dropperOpts, err := bootstrap.ParseAndPromptDropperOpts(c)
	if err != nil {
		return 1, err
	}
	// this is set during dropping
	ourReplicaID := flags.ReplicaID
	masterReplicaID := flags.ExistingReplicaID

	log.Debugf("Collecting certs")
	log.Debugf("replica id = %s", ourReplicaID)
	log.Debugf("join replica id = %s", masterReplicaID)

	certsTar, err := dumpcerts.DumpCerts(c, masterReplicaID)
	if err != nil {
		return 1, err
	}

	bs, err := bootstrap.GetBootstrapClient(c, dropperOpts)
	if err != nil {
		return 1, err
	}

	bs.SetReplicaID(ourReplicaID)
	err = bootstrap.SetupReplicaAndNode(c, bs)
	if err != nil {
		return 1, err
	}

	err = bs.ConfirmNodeIsOkay()
	if err != nil {
		return 1, err
	}

	if err := writeCerts(certsTar); err != nil {
		return 1, err
	}

	if !flags.NoUCP && !flags.SkipNetworkTest {
		firstNode, err := bootstrap.GetNodeWithContainer(bs, containers.Etcd.ReplicaName(masterReplicaID))
		if err != nil {
			return 1, fmt.Errorf("Failed to get node for existing etcd: %s", err)
		}
		secondNode := bs.GetNodeName()
		err = overlayNetworkTest(bs, firstNode, secondNode)
		if err != nil {
			return 1, fmt.Errorf("Overlay networking test failed. Confirm that UCP's overlay networking is working correctly. Reason: %s", err)
		} else {
			log.Info("Overlay networking test passed.")
		}
	}

	err = install.SetupNetworks(c, bs)
	if err != nil {
		return 1, err
	}

	if !flags.NoUCP {
		err := bootstrap.UCPConnTest()
		if err != nil {
			return 1, err
		}
	}

	err = install.SetupVolumes(c, bs)
	if err != nil {
		return 1, err
	}

	if err = install.CreateClientCerts(bs); err != nil {
		log.Debugf("Couldn't create client certs: %s", err)
		return 1, err
	}

	// TODO: make this a function or constant somewhere
	clientURL := fmt.Sprintf("https://%s:%d", containers.Etcd.OverlayName(masterReplicaID), containers.EtcdClientPort1)
	cluster, err := AddEtcdNode(bs, clientURL)
	if err != nil {
		log.Errorf("Error joining cluster: %s", err)
		return 1, err
	}

	_, err = install.StartEtcd(bs, c.Command.Name, cluster)
	if err != nil {
		return 1, err
	}

	kvStore, err := etcd.NewKeyValueStore(containers.EtcdUrls(), deploy.EtcdPath)
	if err != nil {
		log.Errorf("Couldn't set up kvStore: %s", err)
		return 1, err
	}
	settingsStore := sanitizers.Wrap(settingsstore.New(kvStore))

	// Configure Enzi host, ports, etc. This will also start the containers.
	err = reconfigure.Reconfigure(bs, kvStore, settingsStore, c, nil, false, false, false, false)
	if err != nil {
		log.Errorf("Couldn't get ha config: %s", err)
		return 1, err
	}

	// Update the number of replicas we have and set rethink replication
	log.Info("Updating replication settings...")
	haConfig, err := settingsStore.HAConfig()
	if err != nil {
		return 1, fmt.Errorf("Couldn't get ha config: %s", err)
	}
	numReplicas := len(haConfig.ReplicaConfig)
	err = UpdateReplication(numReplicas, ourReplicaID)
	if err != nil {
		return 1, err
	}

	registryConfig, err := settingsStore.RegistryConfig()
	if err != nil {
		return 1, fmt.Errorf("Couldn't get registry config: %s", err)
	}
	if registryConfig.Storage.Type() == "filesystem" {
		log.Warnf("Warning! You are using DTR with the filesystem driver in an HA setup. This will not work with non-clustered file systems. See the documentation for how to configure HA image storage.")
	}

	log.Info("Join is complete")
	log.Infof("Replica ID is set to: %s", ourReplicaID)
	if numReplicas > 1 {
		log.Infof("There are currently %d replicas in your Docker Trusted Registry cluster", numReplicas)
		if (numReplicas % 2) == 0 {
			log.Info("You currently have an even number of replicas which can impact cluster availability")
			log.Info("It is recommended that you have 3, 5 or 7 replicas in your cluster")
		}
	}

	return 0, nil
}

func writeCerts(buf *[]byte) error {
	log.Debugf("Reading certs")

	r := bytes.NewReader(*buf)
	tr := tar.NewReader(r)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			// all done
			break
		}
		if err != nil {
			log.Errorf("Error reading CA certificate tar ball: %s", err)
			return err
		}
		log.Debugf("file: %s", hdr.Name)
		if err := os.MkdirAll(filepath.Dir(hdr.Name), 0755); err != nil {
			return err
		}
		newBuf, err := ioutil.ReadAll(tr)
		if err != nil {
			return err
		}
		log.Debugf("cert: %s", string(newBuf))
		if err := ioutil.WriteFile(hdr.Name, newBuf, 0644); err != nil {
			log.Errorf("Error writing CA certificate: %s", err)
			return err
		}
	}
	log.Debugf("Finished writing certs")
	return nil
}

func pickTwoNodes(bs bootstrap.Bootstrap) (string, string, error) {
	log.Debug("Picking two nodes for overlay network test")
	dc := bs.GetDockerClient()
	info, err := dc.Info(context.Background())
	if err != nil {
		return "", "", err
	}
	startedNodesList := false
	skippingNodeData := false
	firstNode := ""
	secondNode := ""
	// example output that we need to parse
	//Filters: health, port, dependency, affinity, constraint
	//Nodes: 1
	// compooter: 172.17.0.1:12376
	//  └ Status: Healthy
	//  └ Containers: 15
	//  └ Reserved CPUs: 0 / 4
	//  └ Reserved Memory: 0 B / 8.094 GiB
	//  └ Labels: executiondriver=, kernelversion=4.4.5-1-docker-aufs, operatingsystem=Arch Linux, storagedriver=aufs
	//  └ Error: (none)
	//  └ UpdatedAt: 2016-04-27T03:50:41Z
	//  └ ServerVersion: 1.11.0-rc5
	//Cluster Managers: 1
	// 172.17.0.1: Healthy
	for _, line := range info.SystemStatus {
		log.Debug("Parsing UCP status line: %v", line)
		key := line[0]
		value := line[1]
		if skippingNodeData {
			if len(key) > 1 && key[0] == ' ' && key[1] == ' ' {
				continue
			} else {
				secondNode = strings.TrimPrefix(key, " ")
				break
			}
		}
		if startedNodesList && len(key) > 1 && key[0] == ' ' && key[1] != ' ' {
			firstNode = strings.TrimPrefix(key, " ")
			skippingNodeData = true
			continue
		}
		if key == "Nodes" {
			nodes, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return "", "", err
			}
			if nodes < 2 {
				log.Warn("Skipping overlay networking test because only one UCP node was found.")
				return "", "", nil
			}
			startedNodesList = true
			continue
		}
	}
	log.Debug("Picked the first two UCP nodes for the test: %s %s", firstNode, secondNode)
	return firstNode, secondNode, nil
}
func overlayNetworkTest(bs bootstrap.Bootstrap, firstNode, secondNode string) error {
	log.Info("Starting overlay network test")
	// We do the equivalent of:
	// docker run -it --rm --net dtr-ol --name overlay-test1 --entrypoint sh dockerhubenterprise/dtr-dev
	// docker run -it --rm --net dtr-ol --name overlay-test2 --entrypoint ping dockerhubenterprise/dtr-dev -c 3 overlay-test1

	container1 := containers.ContainerConfig{
		Name:     deploy.OverlayTestContainer1Name,
		Image:    deploy.BootstrapRepo.TaggedName(),
		Networks: []containers.NetworkConfig{{Name: deploy.OverlayNetworkName}},
		Node:     firstNode,
		Entrypoint: []string{
			"sleep", "600",
		},
	}
	container2 := containers.ContainerConfig{
		Name:     deploy.OverlayTestContainer2Name,
		Image:    deploy.BootstrapRepo.TaggedName(),
		Networks: []containers.NetworkConfig{{Name: deploy.OverlayNetworkName}},
		Node:     secondNode,
		Entrypoint: []string{
			"ping", "-c", "3", deploy.OverlayTestContainer1Name,
		},
	}

	// now we:
	// 1. start container 1
	// 2. wait for it to start
	// 3. start container 2
	// 4. wait for it to stop
	// 5. check the return code of container 2
	// 6. delete both containers

	// we don't need to set node or replica ids on the bootstrap because they don't affect the particular features that we are using here
	// 1.
	resp1, err := bs.ContainerCreateFromContainerConfig(container1)
	if err != nil {
		return fmt.Errorf("Failed to create container1: %s", err)
	}
	err = bs.ContainerStart(resp1.ID)
	if err != nil {
		return nil
	}

	// 6.
	defer func() {
		err := bs.ContainerRemove(resp1.ID, types.ContainerRemoveOptions{
			Force: true,
		})
		if err != nil {
			log.Warnf("Failed to remove container %s with id %s", deploy.OverlayTestContainer1Name, resp1.ID)
		}
	}()

	// 2.
	err = dtrutil.Poll(time.Second, 30, func() error {
		info, err := bs.ContainerInspect(resp1.ID)
		if err != nil {
			return err
		}
		if info.State == nil || !info.State.Running {
			return fmt.Errorf("Container not running.")
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Failed to wait for container1 to start: %s", err)
	}

	// 3.
	resp2, err := bs.ContainerCreateFromContainerConfig(container2)
	if err != nil {
		return fmt.Errorf("Failed to create container2: %s", err)
	}
	err = bs.ContainerStart(resp2.ID)
	if err != nil {
		return nil
	}

	// 6.
	defer func() {
		err := bs.ContainerRemove(resp2.ID, types.ContainerRemoveOptions{
			Force: true,
		})
		if err != nil {
			log.Warnf("Failed to remove container %s with id %s", deploy.OverlayTestContainer2Name, resp2.ID)
		}
	}()

	// 4.
	result, err := bs.ContainerWait(resp2.ID)
	if err != nil {
		return fmt.Errorf("Failed to wait for container2 to finish: %s", err)
	}

	// 5.
	if result != 0 {
		return fmt.Errorf("Failed to ping from container2 to container1.")
	}

	return nil
}

func AddEtcdNode(bs bootstrap.Bootstrap, clientURL string) ([]string, error) {
	// XXX - the way the etcd client works here is different than the way it works in
	//       hubconfig/etcd/keyvaluestore.go.  We should refactor the other code to work
	//       similarly to this, plus also have a generic interface for getting a
	//       reusable etcd client.
	etcdName := containers.Etcd.ReplicaName(bs.GetReplicaID())
	// TODO: make this a constant or function somewhere. It's too magic.
	peerURL := fmt.Sprintf("https://%s:%d", containers.Etcd.OverlayName(bs.GetReplicaID()), containers.EtcdPeerPort1)

	etcdClient, err := bootstrap.GetEtcdConn(clientURL)
	if err != nil {
		log.Errorf("Couldn't get etcd connection: %s", err)
		return nil, err
	}

	m := client.NewMembersAPI(etcdClient)
	log.Debugf("client = %q", m)

	ctx := context.Background()
	member, err := m.Add(ctx, peerURL)
	if err != nil {
		if err == context.Canceled {
			log.Error("canceled")
			return nil, err
		} else if err == context.DeadlineExceeded {
			log.Error("deadline exceeded")
			return nil, err
		} else if cerr, ok := err.(*client.ClusterError); ok {
			log.Errorf("cluster error:  %q", cerr.Errors)
			return nil, err
		} else {
			log.Errorf("Couldn't join cluster: %s", err)
			return nil, err
		}
		return nil, err
	}

	cluster := []string{}
	members, err := m.List(ctx)
	if err != nil {
		log.Errorf("Error getting cluster member list: %s", err)
		return nil, err
	}

	for _, memb := range members {
		for _, u := range memb.PeerURLs {
			mName := memb.Name
			log.Debugf("found '%s' -> '%s'", mName, member.ID)
			if memb.ID == member.ID {
				// skip adding ourself
				mName = etcdName
				continue
			}
			cluster = append(cluster, fmt.Sprintf("%s=%s", mName, u))
		}
	}

	log.Debugf("cluster = %s", cluster)
	return cluster, nil
}

func UpdateReplication(numReplicas int, replicaID string) error {
	dbSession, err := dtrutil.GetRethinkSession(replicaID)
	if err != nil {
		return fmt.Errorf("Failed to create db session: %s", err)
	}
	schemaMgr := jschema.NewJobrunnerManager(deploy.JobrunnerDBName, dbSession)
	err = schemaMgr.ScaleDB(uint(numReplicas), false)
	if err != nil {
		return fmt.Errorf("Failed to set up jobrunner tables: %s", err)
	}

	schemaManager := bootstrap.CreateManager(replicaID)
	if err = schemaManager.SetReplication(numReplicas); err != nil {
		return fmt.Errorf("Error changing rethink replication factor: %s", err)
	}
	return nil
}
