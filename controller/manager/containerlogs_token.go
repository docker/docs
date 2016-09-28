package manager

import (
	"encoding/json"
	"fmt"
	"path"
	"time"

	log "github.com/Sirupsen/logrus"
	kvstore "github.com/docker/libkv/store"
	"github.com/docker/orca"
	"github.com/docker/orca/utils"
)

var (
	ksContainerLogsTokens = datastoreVersion + "/containerlogs_tokens"
)

func (m DefaultManager) CreateContainerLogsToken(c *orca.ContainerLogsToken) error {
	k := path.Join(ksContainerLogsTokens, c.Token)
	data, err := json.Marshal(c)
	if err != nil {
		return nil
	}

	if err := m.Datastore().Put(k, data, &kvstore.WriteOptions{TTL: 5 * time.Second}); err != nil {
		return utils.MaybeWrapEtcdClusterErr(err)
	}

	m.logEvent("create-containerlogs-token", fmt.Sprintf("container=%s", c.ContainerID), []string{"security"})

	return nil
}

func (m DefaultManager) RemoveContainerLogsToken(c *orca.ContainerLogsToken) error {
	k := path.Join(ksContainerLogsTokens, c.Token)
	exists, err := m.Datastore().Exists(k)
	if err != nil {
		return err
	}

	if !exists {
		return ErrContainerLogsTokenDoesNotExist
	}

	if err := m.Datastore().Delete(k); err != nil {
		return err
	}

	return nil
}

func (m DefaultManager) ContainerLogsToken(token string) (*orca.ContainerLogsToken, error) {
	k := path.Join(ksContainerLogsTokens, token)
	exists, err := m.Datastore().Exists(k)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrContainerLogsTokenDoesNotExist
	}

	kvPair, err := m.Datastore().Get(k)
	if err != nil {
		return nil, utils.MaybeWrapEtcdClusterErr(err)
	}

	var c *orca.ContainerLogsToken
	if err := json.Unmarshal(kvPair.Value, &c); err != nil {
		return nil, err
	}

	return c, nil
}

func (m DefaultManager) ValidateContainerLogsToken(containerId string, token string) bool {
	cs, err := m.ContainerLogsToken(token)
	if err != nil {
		log.Errorf("error validating containerlogs token: %s", err)
		return false
	}

	if cs == nil || cs.ContainerID != containerId {
		log.Warnf("unauthorized token request: %s", token)
		return false
	}

	if err := m.RemoveContainerLogsToken(cs); err != nil {
		log.Error(err)
		return false
	}

	return true
}
