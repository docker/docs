package manager

import (
	"encoding/json"
	"fmt"
	"path"

	"github.com/docker/orca/dockerhub"
	"github.com/docker/orca/utils"
)

var (
	ksWebhookKeys = datastoreVersion + "/webhookkeys"
)

func getWebhookKeyHash(key string) string {
	return GetKeyHash(key)
}

func (m DefaultManager) WebhookKey(key string) (*dockerhub.WebhookKey, error) {
	keyHash := getWebhookKeyHash(key)
	k := path.Join(ksWebhookKeys, keyHash)

	exists, err := m.Datastore().Exists(k)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrWebhookKeyDoesNotExist
	}

	var webhookKey *dockerhub.WebhookKey

	kvPair, err := m.Datastore().Get(k)
	if err != nil {
		return nil, utils.MaybeWrapEtcdClusterErr(err)
	}

	if err := json.Unmarshal(kvPair.Value, &webhookKey); err != nil {
		return nil, err
	}

	return webhookKey, nil
}

func (m DefaultManager) WebhookKeys() ([]*dockerhub.WebhookKey, error) {
	keys := []*dockerhub.WebhookKey{}

	kvPairs, err := m.Datastore().List(ksWebhookKeys)
	if len(kvPairs) == 0 {
		// we cannot check if the key exists since it is a tree
		// therefore we simply return an empty slice
		return keys, nil
	}

	if err != nil {
		return nil, err
	}

	for _, kvPair := range kvPairs {
		var key *dockerhub.WebhookKey
		if err := json.Unmarshal(kvPair.Value, &key); err != nil {
			return nil, err
		}

		keys = append(keys, key)
	}

	return keys, nil
}

func (m DefaultManager) NewWebhookKey(image string) (*dockerhub.WebhookKey, error) {
	k := utils.GenerateID(16)
	key := &dockerhub.WebhookKey{
		Key:   k,
		Image: image,
	}

	if err := m.SaveWebhookKey(key); err != nil {
		return nil, err
	}

	return key, nil
}

func (m DefaultManager) SaveWebhookKey(key *dockerhub.WebhookKey) error {
	keyHash := getWebhookKeyHash(key.Key)
	k := path.Join(ksWebhookKeys, keyHash)

	data, err := json.Marshal(key)
	if err != nil {
		return err
	}

	if err := m.Datastore().Put(k, data, nil); err != nil {
		return utils.MaybeWrapEtcdClusterErr(err)
	}

	m.logEvent("add-webhook-key", fmt.Sprintf("image=%s", key.Image), []string{"webhook"})

	return nil
}

func (m DefaultManager) DeleteWebhookKey(key string) error {
	keyHash := getWebhookKeyHash(key)
	k := path.Join(ksWebhookKeys, keyHash)

	exists, err := m.Datastore().Exists(k)
	if err != nil {
		return err
	}

	if !exists {
		return ErrWebhookKeyDoesNotExist
	}

	if err := m.Datastore().Delete(k); err != nil {
		return err
	}

	m.logEvent("delete-webhook-key", fmt.Sprintf("key=%s", key), []string{"webhook"})

	return nil
}
