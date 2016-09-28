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
	ksConsoleSessions = datastoreVersion + "/consolesessions"
)

func (m DefaultManager) CreateConsoleSession(c *orca.ConsoleSession) error {
	k := path.Join(ksConsoleSessions, c.Token)
	data, err := json.Marshal(c)
	if err != nil {
		return nil
	}

	if err := m.Datastore().Put(k, data, &kvstore.WriteOptions{TTL: 5 * time.Second}); err != nil {
		return utils.MaybeWrapEtcdClusterErr(err)
	}

	m.logEvent("create-console-session", fmt.Sprintf("container=%s", c.ContainerID), []string{"console"})

	return nil
}

func (m DefaultManager) RemoveConsoleSession(c *orca.ConsoleSession) error {
	k := path.Join(ksConsoleSessions, c.Token)
	exists, err := m.Datastore().Exists(k)
	if err != nil {
		return err
	}

	if !exists {
		return ErrConsoleSessionDoesNotExist
	}

	if err := m.Datastore().Delete(k); err != nil {
		return err
	}

	return nil
}

func (m DefaultManager) ConsoleSession(token string) (*orca.ConsoleSession, error) {
	k := path.Join(ksConsoleSessions, token)
	exists, err := m.Datastore().Exists(k)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrConsoleSessionDoesNotExist
	}

	kvPair, err := m.Datastore().Get(k)
	if err != nil {
		return nil, utils.MaybeWrapEtcdClusterErr(err)
	}

	var c *orca.ConsoleSession
	if err := json.Unmarshal(kvPair.Value, &c); err != nil {
		return nil, err
	}

	return c, nil
}

func (m DefaultManager) ValidateConsoleSessionToken(containerId string, token string) bool {
	cs, err := m.ConsoleSession(token)
	if err != nil {
		log.Errorf("error validating console session token: %s", err)
		return false
	}

	if cs == nil || cs.ContainerID != containerId {
		log.Warnf("unauthorized token request: %s", token)
		return false
	}

	if err := m.RemoveConsoleSession(cs); err != nil {
		log.Error(err)
		return false
	}

	return true
}
