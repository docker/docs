package utils

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/swarm"
	"golang.org/x/net/context"
)

func JoinNode(primaryController, joiningNode Machine, manager bool) error {
	// Get the swarm URL
	primaryEngineClient, err := primaryController.GetEngineAPI()
	if err != nil {
		return err
	}
	joiningEngineClient, err := joiningNode.GetEngineAPI()
	if err != nil {
		return err
	}

	// Discover the manager join secret
	info, err := primaryEngineClient.SwarmInspect(context.TODO())
	if err != nil {
		return err
	}

	//Consider using filter to only look at managers?
	options := types.NodeListOptions{}
	nodes, err := primaryEngineClient.NodeList(context.TODO(), options)
	if err != nil {
		return err
	}
	// Just find the first available manager
	swarmManagerURL := ""
	for _, node := range nodes {
		if node.ManagerStatus != nil && node.ManagerStatus.Addr != "" {
			swarmManagerURL = node.ManagerStatus.Addr
			break
		}
	}
	if swarmManagerURL == "" {
		return fmt.Errorf("Something went wrong - couldn't find any existing manager nodes")
	}

	var token string
	if manager {
		token = info.JoinTokens.Manager
	} else {
		token = info.JoinTokens.Worker
	}

	// Optional: CACertHash
	joinReq := swarm.JoinRequest{
		RemoteAddrs: []string{swarmManagerURL},
		JoinToken:   token,
		ListenAddr:  "0.0.0.0:2377",
	}
	log.Infof("Joining %s to %s", joiningNode.GetName(), swarmManagerURL)
	err = joiningEngineClient.SwarmJoin(context.TODO(), joinReq)
	if err != nil {
		return err
	}
	log.Info("Join completed")

	return nil
}
