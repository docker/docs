package utils

import (
	"fmt"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/libkv"
	kvstore "github.com/docker/libkv/store"
	"github.com/docker/libkv/store/etcd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/docker/orca/config"
)

func init() {
	etcd.Register()
}

func TestKVCAConfigExists(t *testing.T, primaryServer Machine, machineIPs []string) {
	log.Infof("Checking for config of CA servers in KV store for machines: %s", machineIPs)

	serverURL, err := GetOrcaURL(primaryServer)
	require.Nil(t, err)

	serverIP, err := primaryServer.GetIP()
	require.Nil(t, err)

	// This test relies on the fact that a TLS client bundle for an admin
	// user has a cert that is signed by the cluster CA which is trusted
	// by the KV store.
	tlsConfig, err := GetUserTLSConfig(serverURL, GetAdminUser(), GetAdminPassword())
	require.Nil(t, err, "unable to get TLS config for admin user")

	var expectedKeys []string
	for _, machineIP := range machineIPs {
		expectedKeys = append(expectedKeys,
			fmt.Sprintf("/orca/v1/config/clusterca_%s:%d", machineIP, config.SwarmCAPort),
			fmt.Sprintf("/orca/v1/config/clientca_%s:%d", machineIP, config.OrcaCAPort),
		)
	}

	kv, err := libkv.NewStore(
		kvstore.ETCD,
		[]string{fmt.Sprintf("%s:%d", serverIP, config.KvPort)},
		&kvstore.Config{
			ConnectionTimeout: time.Second * 10,
			TLS:               tlsConfig,
		},
	)
	require.Nil(t, err, "unable to get KV Store client")

	kvPairs, err := kv.List("orca/v1/config")
	require.Nil(t, err, "unable to list config entries from kv store")

	// Make a set of all of the config keys.
	configKeys := make(map[string]struct{}, len(kvPairs))
	for _, kvPair := range kvPairs {
		configKeys[kvPair.Key] = struct{}{}
	}

	// Ensure that all of our expected cluster or client ca config keys
	// are in the set of existing config keys.
	for _, expectedKey := range expectedKeys {
		_, ok := configKeys[expectedKey]
		assert.True(t, ok, "config entry for %q not found in kv store", expectedKey)
	}
}
