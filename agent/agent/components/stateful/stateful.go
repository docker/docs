package stateful

import (
	"github.com/cloudflare/cfssl/log"
	"github.com/docker/engine-api/client"

	"github.com/docker/orca/agent/agent/components/stateful/authstore"
	"github.com/docker/orca/agent/agent/components/stateful/kv"
	"github.com/docker/orca/types"
)

type Stateful struct {
	kv.KV
	authstore.AuthStore
}

func (p *Stateful) BuildCurrentConfig(dclient *client.Client, currentCfg *types.NodeConfig) error {
	err := p.KV.BuildCurrentConfig(dclient, currentCfg)
	if err != nil {
		return err
	}
	err = p.AuthStore.BuildCurrentConfig(dclient, currentCfg)
	if err != nil {
		return err
	}
	return nil
}

func (p *Stateful) RequiresReconciliation(expectedCfg *types.NodeConfig, currentCfg *types.NodeConfig) (bool, error) {
	kvReconcile, err := p.KV.RequiresReconciliation(expectedCfg, currentCfg)
	if err != nil {
		return false, err
	}
	if kvReconcile {
		return true, nil
	}
	authReconcile, err := p.AuthStore.RequiresReconciliation(expectedCfg, currentCfg)
	if err != nil {
		return false, err
	}
	if authReconcile {
		return true, nil
	}
	return false, nil
}

func (p *Stateful) Reconcile(dclient *client.Client, expectedCfg *types.NodeConfig, currentCfg *types.NodeConfig) error {
	kvReconcile, err := p.KV.RequiresReconciliation(expectedCfg, currentCfg)
	if err != nil {
		return err
	}
	authReconcile, err := p.AuthStore.RequiresReconciliation(expectedCfg, currentCfg)
	if err != nil {
		return err
	}

	if !kvReconcile && !authReconcile {
		return nil
	}

	// Take the lock
	chDone := make(chan bool)
	// Defer a lock release
	defer func() {
		// TODO: This channel hangs forever, debug and rewrite how the lock is taken
		//	chDone <- true
		log.Debug("closing channel for kv lock")
		close(chDone)
		log.Info("Stateful component lock released")
	}()

	// Reconcile KV, if needed
	if kvReconcile {
		err = p.KV.Reconcile(dclient, expectedCfg, currentCfg, chDone)
		if err != nil {
			return err
		}
	}

	// Reconcile Auth Store, if needed
	if authReconcile {
		err = p.AuthStore.Reconcile(dclient, expectedCfg, currentCfg, chDone)
		if err != nil {
			return err
		}
	}
	log.Info("Stateful component reconciled succesfully")
	return nil
}

func NewStateful() *Stateful {
	return &Stateful{
		KV:        kv.KV{},
		AuthStore: authstore.AuthStore{},
	}
}
