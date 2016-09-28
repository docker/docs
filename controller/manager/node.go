package manager

import (
	"archive/tar"
	"bytes"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"path"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/swarm"
	kvstore "github.com/docker/libkv/store"
	"golang.org/x/net/context"

	"github.com/docker/orca"
	"github.com/docker/orca/config"
	"github.com/docker/orca/pkg/pki"
	orcatypes "github.com/docker/orca/types"
	"github.com/docker/orca/utils"
)

var (
	ErrNodeDoesNotExist = errors.New("node does not exist")
	KvSecretKey         = path.Join(datastoreVersion, "secret")
	KvSecretTTL         = 5 * time.Minute
)

func parseClusterNodes(systemStatus [][2]string) ([]*orca.Node, error) {
	nodes := []*orca.Node{}
	node := &orca.Node{}

	//First find where the nodes list starts
	startOffset := 0
	for startOffset < len(systemStatus) {
		if systemStatus[startOffset][0] == "Nodes" {
			break
		}
		startOffset++
	}
	if startOffset == 0 {
		log.Debug("No nodes detected in SystemStatus")
	}

	// Now loop through the remaining lines and build up Node objects
	for _, l := range systemStatus[startOffset+1:] {
		label := l[0]
		data := l[1]

		if strings.Index(label, "  └") == -1 {
			// Special case for first node
			if node.Name != "" {
				nodes = append(nodes, node)
				node = &orca.Node{}
			}
			node.Name = strings.TrimSpace(label)
			node.Addr = strings.TrimSpace(data)
			continue
		}

		// node info like "Containers"
		switch label {
		case "  └ ID":
			node.ID = strings.TrimSpace(data)
		case "  └ ServerVersion":
			node.ServerVersion = strings.TrimSpace(data)
		case "  └ Containers":
			node.Containers = strings.TrimSpace(data)
		case "  └ Reserved CPUs":
			node.ReservedCPUs = strings.TrimSpace(data)
		case "  └ Reserved Memory":
			node.ReservedMemory = strings.TrimSpace(data)
		case "  └ Labels":
			lbls := strings.Split(data, ",")
			node.Labels = lbls
		case "  └ Status":
			node.Status = data
		case "  └ Error":
			if data != "(none)" { // Special case so we keep the UI cleaner
				node.Error = strings.TrimSpace(data)
			}
		case "  └ UpdatedAt":
			parsedTime, err := time.Parse(time.RFC3339, data)
			if err != nil {
				log.Warnf("Failed to parse node UpdatedAt time: %s - %s", data, err)
			} else {
				node.UpdatedAt = parsedTime
			}
		default:
			// From time to time swarm adds new labels, so we should make sure this doesn't show
			// up in the logs
			log.Debug("unrecognized node label from swarm:%s,%s", label, data)
		}
	}
	if node.Name != "" {
		nodes = append(nodes, node)
	}

	return nodes, nil
}

func (m DefaultManager) Nodes() ([]*orca.Node, error) {
	info, err := m.client.Info(context.TODO())
	if err != nil {
		return nil, err
	}

	return parseClusterNodes(info.SystemStatus)
}

func (m DefaultManager) Node(name string) (*orca.Node, error) {
	nodes, err := m.Nodes()
	if err != nil {
		return nil, err
	}

	for _, node := range nodes {
		if node.Name == name {
			return node, nil
		}
	}

	return nil, nil
}

func (m DefaultManager) nodeIDFromHostname(hostname string) (string, error) {
	nodes, err := m.proxyClient.NodeList(context.TODO(), types.NodeListOptions{})
	if err != nil {
		return "", err
	}

	for _, node := range nodes {
		if node.Description.Hostname == hostname {
			return node.ID, nil
		}
	}
	return "", fmt.Errorf("No swarm-mode nodes detected with hostname %s", hostname)
}

func (m DefaultManager) PromoteNode(remoteAddr string, replicateServerCerts bool) error {
	var localNodeName string
	var remoteNodeName string

	nodes, err := m.Nodes()
	if err != nil {
		return err
	}

	// Get Hostname from IP
	for _, node := range nodes {
		nodeAddr := strings.Split(node.Addr, ":")[0]

		if strings.Contains(nodeAddr, m.hostAddr) {
			localNodeName = node.Name
		} else if strings.Contains(nodeAddr, remoteAddr) {
			remoteNodeName = node.Name
		}
		if localNodeName != "" && remoteNodeName != "" {
			break
		}
	}

	// Ensure both nodes were found
	if localNodeName == "" {
		err = fmt.Errorf("Unable to find a classic swarm node for the present controller at %s", m.hostAddr)
		log.Error(err)
		return err
	}
	if remoteNodeName == "" {
		err = fmt.Errorf("Unable to find a classic swarm node for the requested node at %s", remoteAddr)
		log.Error(err)
		return err
	}

	// Get NodeID of the target node and validate whether it's a swarm-mode manager
	remoteNodeID, err := m.nodeIDFromHostname(remoteNodeName)
	if err != nil {
		return err
	}
	node, _, err := m.proxyClient.NodeInspectWithRaw(context.TODO(), remoteNodeID)
	if err != nil {
		return err
	}
	if node.Spec.Role != swarm.NodeRoleManager {
		return fmt.Errorf("Remote node %s attempted to promote without being a swarm-mode manager", remoteNodeID)
	}

	log.Infof("Copying root key material to node %s", remoteNodeID)

	// Copy over the cert.pem and key.pem files for both the Client CA and the Cluster CA
	reader, _, err := m.client.CopyFromContainer(context.TODO(), localNodeName+"/ucp-cluster-root-ca", "/etc/cfssl/cert.pem")
	if err != nil {
		return err
	}
	err = m.client.CopyToContainer(context.TODO(), fmt.Sprintf("%s/%s", remoteNodeName, config.OrcaReconcileContainerName),
		"/var/lib/docker/ucp/ucp-cluster-root-ca", reader, types.CopyToContainerOptions{})
	if err != nil {
		return err
	}
	reader, _, err = m.client.CopyFromContainer(context.TODO(), localNodeName+"/ucp-cluster-root-ca", "/etc/cfssl/key.pem")
	if err != nil {
		return err
	}
	err = m.client.CopyToContainer(context.TODO(), fmt.Sprintf("%s/%s", remoteNodeName, config.OrcaReconcileContainerName),
		"/var/lib/docker/ucp/ucp-cluster-root-ca", reader, types.CopyToContainerOptions{})
	if err != nil {
		return err
	}
	reader, _, err = m.client.CopyFromContainer(context.TODO(), localNodeName+"/ucp-client-root-ca", "/etc/cfssl/cert.pem")
	if err != nil {
		return err
	}
	err = m.client.CopyToContainer(context.TODO(), fmt.Sprintf("%s/%s", remoteNodeName, config.OrcaReconcileContainerName),
		"/var/lib/docker/ucp/ucp-client-root-ca", reader, types.CopyToContainerOptions{})
	if err != nil {
		return err
	}
	reader, _, err = m.client.CopyFromContainer(context.TODO(), localNodeName+"/ucp-client-root-ca", "/etc/cfssl/key.pem")
	if err != nil {
		return err
	}
	err = m.client.CopyToContainer(context.TODO(), fmt.Sprintf("%s/%s", remoteNodeName, config.OrcaReconcileContainerName),
		"/var/lib/docker/ucp/ucp-client-root-ca", reader, types.CopyToContainerOptions{})
	if err != nil {
		return err
	}

	// The controller was requested to replicate its existing server certs to the node under promotion
	if replicateServerCerts {
		log.Info("Attempting to replicate server certs or SANs")
		// Determine if this controller's cert is signed by the UCP Client CA and, if so, don't replicate
		der, _ := pem.Decode(m.controllerCertPEM)
		if der == nil {
			return fmt.Errorf("Error decoding existing controller's certificate")

		}
		cert, err := x509.ParseCertificate(der.Bytes)
		if err != nil {
			return err
		}
		if strings.Contains(cert.Issuer.CommonName, "UCP Client Root CA") {
			// The present node's server cert is not externally signed.
			// Append this node's SAN labels to the remote node
			info, err := m.proxyClient.Info(context.TODO())
			if err != nil {
				return err
			}

			log.Infof("Updating SAN Labels to node %s", remoteNodeID)
			return m.updateSANLabels(info.Swarm.NodeID, remoteNodeID)
		}

		// Certificate is externally signed, proceed with replication
		serverCerts := orcatypes.ServerCertRequest{
			CA:   m.clientCAChain,
			Cert: string(m.controllerCertPEM),
			Key:  string(m.controllerKeyPEM),
		}
		log.Infof("Replicating this controller's server cert to node %s", remoteNodeID)
		return m.copyCertsToRemoteNodes([]string{remoteNodeName}, config.OrcaReconcileContainerName,
			"/var/lib/docker/ucp/ucp-controller-server-certs", serverCerts)
	}

	return nil
}

// updateSANLabels concatenates the comma-separated SANs of a source and target node
// and updates the target node's labels with the result, excluding duplicates
func (m DefaultManager) updateSANLabels(sourceNodeID, targetNodeID string) error {
	sourceSANs := make(map[string]struct{})
	targetSANs := make(map[string]struct{})

	// Append the source node's SANs to the final list of SANs
	sourceNode, _, err := m.proxyClient.NodeInspectWithRaw(context.TODO(), sourceNodeID)
	if err != nil {
		return err
	}
	if sourceNode.Spec.Annotations.Labels != nil {
		existingSANs, ok := sourceNode.Spec.Annotations.Labels[config.SANNodeLabel]
		if ok {
			for _, key := range strings.Split(existingSANs, ",") {
				sourceSANs[key] = struct{}{}
			}
		}
	}

	// Append the target node's SANs to the final list of SANs, if not a duplicate
	targetNode, _, err := m.proxyClient.NodeInspectWithRaw(context.TODO(), targetNodeID)
	if err != nil {
		return err
	}
	if targetNode.Spec.Annotations.Labels == nil {
		targetNode.Spec.Annotations.Labels = make(map[string]string)
	}
	targetSANList, ok := targetNode.Spec.Annotations.Labels[config.SANNodeLabel]
	if ok {
		for _, key := range strings.Split(targetSANList, ",") {
			targetSANs[key] = struct{}{}
		}
	}

	// Append the sourceSANs to targetSANs, excluding duplicates
	for key, _ := range sourceSANs {
		targetSANs[key] = struct{}{}
	}

	// TODO: augment targetSANs with any standard expected SANs

	// flatten the final set of SANs in a list
	finalSANList := make([]string, len(targetSANs))
	i := 0
	for san, _ := range targetSANs {
		finalSANList[i] = san
		i++
	}

	targetNode.Spec.Annotations.Labels[config.SANNodeLabel] = strings.Join(finalSANList, ",")

	// Update the target node labels with the SANs
	return m.proxyClient.NodeUpdate(context.TODO(), targetNodeID, targetNode.Meta.Version, targetNode.Spec)
}

func (m DefaultManager) controllers() ([]orcatypes.Controller, error) {
	// Obtain the controller blobs from the KV Store
	var res []orcatypes.Controller
	controllersData, err := m.Datastore().List(path.Join(datastoreVersion, "controllers"))
	if err != nil {
		return res, err
	}

	// Unmarshal into controller types
	var controller orcatypes.Controller
	for _, controllerData := range controllersData {
		err = json.Unmarshal(controllerData.Value, &controller)
		if err != nil {
			return res, err
		}
		res = append(res, controller)
	}
	return res, nil
}

func (m DefaultManager) controllerNames() ([]string, error) {
	var controllerNames []string
	// get all controllers
	controllers, err := m.controllers()
	if err != nil {
		return controllerNames, err
	}

	// get all controller IPs
	controllerIPs := orcatypes.GetIPsFromControllers(controllers)

	// get all Nodes
	nodes, err := m.Nodes()
	if err != nil {
		return controllerNames, err
	}

	// iterate over all nodes until all the controllers are found
	addLast := ""
	for _, node := range nodes {
		for _, controllerIP := range controllerIPs {
			if strings.Contains(node.Addr, controllerIP) {
				if strings.Contains(m.engineProxyURL.String(), controllerIP) {
					// Make sure that this controller is the last one in the list
					addLast = node.Name
				} else {
					controllerNames = append(controllerNames, node.Name)
				}
			}
		}
	}

	if addLast != "" {
		controllerNames = append(controllerNames, addLast)
	}

	if len(controllerNames) != len(controllers) {
		return controllerNames, fmt.Errorf("internal error: unable to detect hostnames for all controllers")
	}
	return controllerNames, nil
}

func writeStringToTar(tw *tar.Writer, path, input string) error {
	data := []byte(input)
	hdr := &tar.Header{
		Name: path,
		Mode: 0600,
		Size: int64(len(data)),
	}
	err := tw.WriteHeader(hdr)
	if err != nil {
		return err
	}
	_, err = tw.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (m DefaultManager) copyCertsToRemoteNodes(targets []string, containerName, path string, req orcatypes.ServerCertRequest) error {

	// Create a buffer to hold an archive of the cert material
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)

	err := writeStringToTar(tw, "./cert.pem", req.Cert)
	if err != nil {
		return err
	}
	err = writeStringToTar(tw, "./ca.pem", req.CA)
	if err != nil {
		return err
	}
	err = writeStringToTar(tw, "./key.pem", req.Key)
	if err != nil {
		return err
	}

	archive := buf.Bytes()
	for _, target := range targets {
		err = m.client.CopyToContainer(context.TODO(), fmt.Sprintf("%s/%s", target, containerName), path, bytes.NewBuffer(archive), types.CopyToContainerOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

func (m DefaultManager) SetControllerServerCerts(nodeName string, req orcatypes.ServerCertRequest) error {
	// TODO: parse req and make sure it's a valid set of certificates that will not break the cluster

	// Determine whether this operation will affect one or all the controllers
	var targets []string
	var err error
	if nodeName == "" {
		targets, err = m.controllerNames()
		if err != nil {
			return err
		}
	} else {
		targets = append(targets, nodeName)
	}

	// Copy the certs to the target controllers
	err = m.copyCertsToRemoteNodes(targets, "ucp-controller", "/etc/docker/ssl/orca", req)

	// Restart the target controller(s) after sending a response back to the API caller
	go func() {
		time.Sleep(time.Second)
		for _, target := range targets {
			timeout := 10 * time.Second
			_ = m.client.ContainerRestart(context.TODO(), target+"/ucp-controller", &timeout)
		}
	}()

	return nil
}

// TODO Refine this
func (m DefaultManager) getSwarmArgs() ([]string, error) {
	// find swarm container
	containers, err := m.client.ContainerList(context.TODO(), types.ContainerListOptions{All: true, Size: false})
	if err != nil {

		log.Debugf("Swarm discovery failed %s", err)
		return nil, err
	}
	swarmArgs := []string{}
	for _, cnt := range containers {
		cInfo, err := m.client.ContainerInspect(context.TODO(), cnt.ID)
		if err != nil {
			log.Debugf("Failed to lookup container %s, continuing...", cnt.ID)
			continue
		}

		if _, node := cInfo.Config.Labels["com.docker.ucp.node"]; node {
			log.Debugf("using swarm container for node args: id=%s", cnt.ID)
			swarmArgs = cInfo.Config.Cmd
			break
		}
	}
	// TODO Consider returning error if len(swarmArgs) == 0
	return swarmArgs, nil
}

func (m DefaultManager) refreshSecret() error {
	log.Info("Refreshing the UCP Join secret")
	// Create a new Secret
	newSecret, err := m.genRandomUUID()
	if err != nil {
		return err
	}
	// Put it in the KV store
	err = m.Datastore().Put(KvSecretKey, []byte(newSecret), &kvstore.WriteOptions{
		IsDir: false,
		TTL:   KvSecretTTL,
	})
	if err != nil {
		return utils.MaybeWrapEtcdClusterErr(err)
	}

	// Inspect the existing agent service
	service, _, err := m.proxyClient.ServiceInspectWithRaw(context.TODO(), "ucp-agent")
	if err != nil {
		return err
	}

	// Replace the secret in the service definition
	env := service.Spec.TaskTemplate.ContainerSpec.Env
	newSpec := service.Spec
	newEnv := []string{}

	// Attempt to replace an existing secret in the service env vars
	replaced := false
	for _, entry := range env {
		if strings.Contains(entry, "SECRET") {
			newEnv = append(newEnv, "SECRET="+newSecret)
			replaced = true
			continue
		}
		newEnv = append(newEnv, entry)
	}

	// If the secret was not replaced, add it as a new env var
	if !replaced {
		newEnv = append(newEnv, "SECRET="+newSecret)
	}

	newSpec.TaskTemplate.ContainerSpec.Env = newEnv

	// Update the service
	return m.proxyClient.ServiceUpdate(context.Background(), service.ID, service.Version, newSpec, types.ServiceUpdateOptions{})
}

func (m DefaultManager) AuthorizeNodeRequest(r *orca.NodeRequest) (*orca.NodeConfiguration, error) {
	swarmArgs, err := m.getSwarmArgs()
	if err != nil {
		return nil, err
	}
	config := orca.NodeConfiguration{
		ClusterCertificates: map[string]string{},
		UserCertificates:    map[string]string{},
		OrcaID:              m.ID(),
		KvStore:             []string{m.datastoreAddr},
		SwarmArgs:           swarmArgs,
	}

	// Someday we may support the "intermedate" profile
	profile := "node"
	if r.Replica {
		log.Debugf("Joining replica node")
	} else {
		log.Debugf("Joining non-replica node")
	}

	for label, request := range r.ClusterCertificateRequests {
		csr := &pki.CertificateSigningRequest{
			CertificateRequest: request,
			Profile:            profile,
		}
		// TODO - refactor to loop
		certResponse, err := m.ClusterSignCSR(csr)
		if err != nil {
			log.Debugf("Cluster CSR call failed %s", err)
			return nil, err
		}
		config.ClusterCertificates[label] = certResponse.Certificate

		// WARNING:
		// Only put in the root from the cluster - this might not always be the right answer
		// but for now, clusterCAChain contains all the roots, and we don't want cluster
		// only components on the joined node to trust the user certs
		config.ClusterCertificateChain = certResponse.CertificateChain
	}

	for label, request := range r.UserCertificateRequests {
		csr := &pki.CertificateSigningRequest{
			CertificateRequest: request,
			Profile:            profile,
		}
		// TODO - refactor to loop
		certResponse, err := m.ClientSignCSR(csr)
		if err != nil {
			log.Debugf("User CSR call failed %s", err)
			return nil, err
		}
		config.UserCertificates[label] = certResponse.Certificate
		// XXX do we need the join, or is the chain already good-to-go?
		config.UserCertificateChain = utils.JoinCerts(m.clientCAChain, certResponse.CertificateChain)
	}
	if config.UserCertificateChain == "" {
		config.UserCertificateChain = utils.JoinCerts(m.clusterCAChain, m.clientCAChain)
	}

	return &config, nil
}

// injectUCPStatus enhances a swarm.Node with the status of the UCP components on that node
func (m DefaultManager) injectUCPStatus(node swarm.Node) swarm.Node {
	if node.Status.State != swarm.NodeStateReady {
		// The swarm-mode node is not ready yet, do not check for the status of the UCP node
		return node
	}

	// Check if the node's hostname is visible in the classic swarm cluster
	swarmNode, err := m.Node(node.Description.Hostname)
	if err != nil || swarmNode == nil {
		node.Status.State = swarm.NodeStateDown
		node.Status.Message = "Classic Swarm unavailable"
		return node
	}
	node.Status.Message = "Classic Swarm available"

	if node.Spec.Role == swarm.NodeRoleManager {
		// Check if the UCP controller is running on the target node
		// NOTE: The beachhead will not launch the controller container if another component failed to start
		// TODO(alexmavr): more fine-grained checks on the status of KV, Rethinkdb and CAs
		_, err = m.client.ContainerInspect(context.TODO(), node.Description.Hostname+"/ucp-controller")
		if err != nil {
			node.Status.State = swarm.NodeStateDown
			node.Status.Message = "Classic Swarm available, UCP controller not started"
			return node
		}
		node.Status.Message = "Classic Swarm available, UCP controller available"
	}

	return node
}

func (m DefaultManager) InspectNode(nodeID string) (swarm.Node, error) {
	node, _, err := m.proxyClient.NodeInspectWithRaw(context.TODO(), nodeID)
	if err != nil {
		return node, err
	}
	return m.injectUCPStatus(node), nil
}

// ListNodes obtains a list of all swarm-mode nodes
func (m DefaultManager) ListNodes() ([]swarm.Node, error) {
	allNodes, err := m.proxyClient.NodeList(context.TODO(), types.NodeListOptions{})
	if err != nil {
		return allNodes, err
	}

	// Enhance the status of all the nodes in parallel
	var wg sync.WaitGroup
	finalNodes := make([]swarm.Node, len(allNodes))
	wg.Add(len(allNodes))
	for i, node := range allNodes {
		go func(c int, n swarm.Node) {
			finalNodes[c] = m.injectUCPStatus(n)
			wg.Done()
		}(i, node)
	}
	wg.Wait()
	return finalNodes, nil
}
