package manager

import (
	"encoding/json"
	"fmt"
	"path"
	"time"

	log "github.com/Sirupsen/logrus"
	kvstore "github.com/docker/libkv/store"
	"github.com/docker/orca/auth"
	"github.com/docker/orca/utils"
)

var (
	ksAccessLists           = datastoreVersion + "/accesslists"
	ksAccessListsTeams      = datastoreVersion + "/accessliststeams"
	datastoreAccessListLock = datastoreVersion + "/accesslistupdatelock"
)

func getAccessListHash(l *auth.AccessList) string {
	return GetKeyHash(fmt.Sprintf("%d.%s.%s", l.Role, l.TeamId, l.Label))
}

func (m DefaultManager) SaveAccessList(l *auth.AccessList) (string, error) {
	chDone := make(chan bool)

	if err := utils.GetKVLock(m.Datastore(), datastoreAccessListLock, &kvstore.LockOptions{TTL: time.Duration(10) * time.Second}, chDone); err != nil {
		log.Errorf("failed to get lock from KV store for access list update: %s", err)
		return "", err
	}

	defer func() {
		chDone <- true
		log.Debug("closing channel for kv lock")
		close(chDone)
	}()

	listHash := getAccessListHash(l)
	k := path.Join(ksAccessLists, listHash)

	tlKey := path.Join(ksAccessListsTeams, l.TeamId, listHash)

	l.Id = listHash

	data, err := json.Marshal(l)
	if err != nil {
		return "", err
	}

	if err := m.Datastore().Put(k, data, nil); err != nil {
		return "", utils.MaybeWrapEtcdClusterErr(err)
	}

	if err := m.Datastore().Put(tlKey, data, nil); err != nil {
		return "", utils.MaybeWrapEtcdClusterErr(err)
	}

	m.logEvent("add-access-list", fmt.Sprintf("role=%d teamId=%s label=%s", l.Role, l.TeamId, l.Label), []string{"security"})

	return "add-access-list", nil
}

func (m DefaultManager) RemoveAccessList(teamId, id string) error {
	chDone := make(chan bool)

	if err := utils.GetKVLock(m.Datastore(), datastoreAccessListLock, &kvstore.LockOptions{TTL: time.Duration(10) * time.Second}, chDone); err != nil {
		log.Errorf("failed to get lock from KV store for access list update: %s", err)
		return err
	}

	defer func() {
		chDone <- true
		log.Debug("closing channel for kv lock")
		close(chDone)
	}()

	k := path.Join(ksAccessLists, id)

	tlKey := path.Join(ksAccessListsTeams, teamId, id)

	exists, err := m.Datastore().Exists(k)
	if err != nil {
		return err
	}

	if !exists {
		return ErrAccessListDoesNotExist
	}

	tlExists, err := m.Datastore().Exists(tlKey)
	if err != nil {
		return err
	}

	if !tlExists {
		return ErrAccessListDoesNotExist
	}

	if err := m.Datastore().Delete(k); err != nil {
		return err
	}

	if err := m.Datastore().Delete(tlKey); err != nil {
		return err
	}

	m.logEvent("delete-access-list", fmt.Sprintf("id=%s", id), []string{"security"})

	return nil
}

func (m DefaultManager) AccessList(teamId, id string) (*auth.AccessList, error) {
	k := path.Join(ksAccessListsTeams, teamId, id)

	exists, err := m.Datastore().Exists(k)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrAccessListDoesNotExist
	}

	var list *auth.AccessList

	kvPair, err := m.Datastore().Get(k)
	if err != nil {
		return nil, utils.MaybeWrapEtcdClusterErr(err)
	}

	if err := json.Unmarshal(kvPair.Value, &list); err != nil {
		return nil, err
	}

	return list, nil
}

func (m DefaultManager) AccessLists() ([]*auth.AccessList, error) {
	lists := []*auth.AccessList{}

	kvPairs, err := m.Datastore().List(ksAccessLists)
	if len(kvPairs) == 0 {
		// we cannot check if the key exists since it is a tree
		// therefore we simply return an empty slice
		return lists, nil
	}

	if err != nil {
		return nil, err
	}

	for _, kvPair := range kvPairs {
		var list *auth.AccessList
		if err := json.Unmarshal(kvPair.Value, &list); err != nil {
			return nil, err
		}

		lists = append(lists, list)
	}

	return lists, nil
}

func (m DefaultManager) AccessListsForTeam(teamId string) ([]*auth.AccessList, error) {
	lists := []*auth.AccessList{}

	k := path.Join(ksAccessListsTeams, teamId)

	kvPairs, err := m.Datastore().List(k)
	if len(kvPairs) == 0 {
		// we cannot check if the key exists since it is a tree
		// therefore we simply return an empty slice
		return lists, nil
	}

	if err != nil {
		return nil, err
	}

	for _, kvPair := range kvPairs {
		var list *auth.AccessList
		if err := json.Unmarshal(kvPair.Value, &list); err != nil {
			return nil, err
		}

		lists = append(lists, list)
	}

	return lists, nil
}

func (m DefaultManager) AccessListsForAccount(ctx *auth.Context, username string) ([]*auth.AccessList, error) {
	lists := []*auth.AccessList{}

	teams, err := m.GetAuthenticator().ListUserTeams(ctx, username)
	if err != nil {
		return nil, err
	}

	found := map[string]bool{}

	for _, team := range teams {
		teamLists, err := m.AccessListsForTeam(team.Id)
		if err != nil {
			return nil, err
		}

		for _, list := range teamLists {
			if _, ok := found[list.Id]; !ok {
				lists = append(lists, list)
				found[list.Id] = true
			}
		}
	}

	log.Debugf("access lists for account: user=%s lists=%v", username, lists)

	return lists, nil
}
