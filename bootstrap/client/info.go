package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/context"

	"github.com/docker/orca/bootstrap/config"
	orcaconfig "github.com/docker/orca/config"
)

func dedup(hostnames []string) []string {
	ret := []string{}
	for _, newname := range hostnames {
		found := false
		for _, oldname := range ret {
			if newname == oldname {
				found = true
				break
			}
		}
		if !found {
			ret = append(ret, newname)
		}
	}
	return ret
}

func (c *EngineClient) GatherHostnames(interactive bool) error {
	localName, err := c.GetLocalID()
	if err != nil {
		log.Error("Failed to retrieve local name")
		return err
	}
	config.OrcaLocalName = localName
	log.Debugf("Local Name: %s", config.OrcaLocalName)

	hostnames, err := c.GetHostnames()
	if err != nil {
		log.Error("Failed to retrieve local hostnames")
		return err
	}

	hostAddress, err := c.GetHostAddress()
	if err != nil {
		log.Error("Failed to determine host address")
		return err
	}
	config.OrcaHostAddress = hostAddress
	log.Debugf("Host Address: %s", config.OrcaHostAddress)

	hostnames = append(hostnames, hostAddress)

	if len(config.OrcaSANs) > 0 {
		hostnames = append(hostnames, config.OrcaSANs...)
	}

	config.OrcaHostnames = dedup(c.FQDNCheck(hostnames, interactive))
	log.Debugf("Hostnames: %v", config.OrcaHostnames)
	// Update the env var with the full set we've got now
	os.Setenv("UCP_HOSTNAMES", strings.Join(config.OrcaHostnames, ","))

	return nil
}

func (c *EngineClient) SetHostnamesAsLabels() error {
	// Update the present Swarm-mode node's labels to include all SANs
	info, err := c.GetClient().Info(context.TODO())
	if err != nil {
		return err
	}

	node, _, err := c.GetClient().NodeInspectWithRaw(context.TODO(), info.Swarm.NodeID)
	if err != nil {
		return err
	}

	if node.Spec.Labels == nil {
		node.Spec.Labels = make(map[string]string)
	}

	// Append SANs to this node's SAN label
	node.Spec.Labels[orcaconfig.SANNodeLabel] = strings.Join(config.OrcaHostnames, ",")

	return c.GetClient().NodeUpdate(context.TODO(), node.ID, node.Version, node.Spec)
}

func (c *EngineClient) GetLocalID() (string, error) {
	name := os.Getenv("UCP_LOCAL_NAME")
	if name != "" {
		return name, nil
	}
	info, err := c.client.Info()
	if err != nil {
		return "", err
	}

	// XXX Is this the best choice?
	name = info.Name

	os.Setenv("UCP_LOCAL_NAME", name)
	return name, nil
}

func (c *EngineClient) GetHostnames() ([]string, error) {
	names := os.Getenv("UCP_HOSTNAMES")
	if names != "" {
		return strings.Split(names, ","), nil
	}

	info, err := c.client.Info()
	if err != nil {
		return nil, err
	}

	hostnames := []string{info.Name}
	// XXX Are there better sources of information?
	hostnames = append(hostnames, "127.0.0.1", c.bootstrapper.NetworkSettings.Gateway)
	return hostnames, nil
}

// Display a warning if no FQDNs are detected, and in interactive mode, prompt for additional hostnames
func (c *EngineClient) FQDNCheck(hostnames []string, interactive bool) []string {
	fqdnFound := false
	for _, hostname := range hostnames {
		if net.ParseIP(hostname) == nil {
			if strings.Contains(hostname, ".") {
				// It might not be fully qualified, but at least has a domain component
				fqdnFound = true
				// TODO - we might want to attempt a lookup with a short timeout and warn if we can't resolve it...
			}
		}
	}

	if fqdnFound && !interactive {
		return hostnames
	}

	if !fqdnFound && !config.InPhase2 {
		log.Warnf("None of the hostnames we'll be using in the UCP certificates %v contain a domain component.  Your generated certs may fail TLS validation unless you only use one of these shortnames or IPs to connect.  You can use the --san flag to add more aliases", hostnames)
	} else if interactive {
		fmt.Printf("We detected the following hostnames/IP addresses for this system %v\n", hostnames)
	}

	if interactive {
		fmt.Print("\nYou may enter additional aliases (SANs) now or press enter to proceed with the above list.\nAdditional aliases: ")

		reader := bufio.NewReader(os.Stdin)
		value, err := reader.ReadString('\n')
		if err != nil {
			log.Warnf("Failed to read input: %s", err)
			return hostnames
		}
		log.Debugf("User entered: %s", value)
		// Split by comma or space
		if strings.Contains(value, ",") {
			for _, newName := range strings.Split(value, ",") {
				cleaned := strings.TrimSpace(newName)
				if cleaned != "" {
					hostnames = append(hostnames, cleaned)
				}
			}

		} else if strings.Contains(value, " ") {
			for _, newName := range strings.Split(value, " ") {
				cleaned := strings.TrimSpace(newName)
				if cleaned != "" {
					hostnames = append(hostnames, cleaned)
				}
			}
		} else {
			// A single hostname/ip
			cleaned := strings.TrimSpace(value)
			if cleaned != "" {
				hostnames = append(hostnames, cleaned)
			}
		}
	}
	return hostnames
}
