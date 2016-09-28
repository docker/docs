package components

import (
	"github.com/docker/engine-api/client"

	"github.com/docker/orca/agent/agent/components/authapi"
	"github.com/docker/orca/agent/agent/components/authworker"
	"github.com/docker/orca/agent/agent/components/certs"
	"github.com/docker/orca/agent/agent/components/clientca"
	"github.com/docker/orca/agent/agent/components/clusterca"
	"github.com/docker/orca/agent/agent/components/controller"
	"github.com/docker/orca/agent/agent/components/csr"
	"github.com/docker/orca/agent/agent/components/proxy"
	"github.com/docker/orca/agent/agent/components/stateful"
	"github.com/docker/orca/agent/agent/components/swarmjoin"
	"github.com/docker/orca/agent/agent/components/swarmmanager"
	"github.com/docker/orca/types"
)

type Component interface {
	BuildCurrentConfig(dclient *client.Client, currentCfg *types.NodeConfig) error
	RequiresReconciliation(expectedCfg *types.NodeConfig, currentCfg *types.NodeConfig) (bool, error)
	Reconcile(dclient *client.Client, expectedCfg *types.NodeConfig, currentCfg *types.NodeConfig) error
}

// ComponentList is an ordered list of the UCP system components
var ComponentList []Component = []Component{
	&csr.CSR{},
	&proxy.Proxy{},
	&swarmjoin.SwarmJoin{},
	&certs.Certs{},
	stateful.NewStateful(),
	&clusterca.ClusterCA{},
	&clientca.ClientCA{},
	&swarmmanager.SwarmManager{},
	&authapi.AuthAPI{},
	&authworker.AuthWorker{},
	&controller.Controller{},
}
