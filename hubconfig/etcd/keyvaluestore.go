package etcd

import (
	"bytes"
	"path"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/coreos/etcd/client"
	"github.com/coreos/etcd/pkg/transport"
	"github.com/docker/dhe-deploy/hubconfig"
	"github.com/docker/dhe-deploy/shared/containers"
	"github.com/docker/dhe-deploy/shared/dtrutil"
	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/etcd"
)

func init() {
	etcd.Register()
}

type keyValueStore struct {
	etcdUrls []string
	etcdPath string
	store    store.Store
}

func NewKeyValueStore(etcdUrls []string, etcdPath string) (hubconfig.KeyValueStore, error) {
	log.Debug("Connecting kvstore at %s", etcdUrls)
	tlsInfo := transport.TLSInfo{
		CertFile: containers.EtcdCertStore.CertPath(),
		KeyFile:  containers.EtcdCertStore.KeyPath(),
		CAFile:   containers.EtcdCACertStore.CertPath(),
	}
	t, err := transport.NewTransport(tlsInfo, time.Minute)
	if err != nil {
		return nil, err
	}

	config := &store.Config{
		ConnectionTimeout: time.Minute,
		TLS:               t.TLSClientConfig,
	}
	kvStore, err := libkv.NewStore(
		store.ETCD,
		etcdUrls,
		config,
	)
	if err != nil {
		if err, ok := err.(*client.ClusterError); ok {
			log.Errorf(err.Detail())
		}
		return nil, err
	}
	kv := keyValueStore{
		etcdUrls: etcdUrls,
		etcdPath: etcdPath,
		store:    kvStore,
	}
	// Wait for kv store to become avaiable
	err = kv.Ping()
	return &kv, err
}

// If the key to get is not found, we return an empty byte slice
func (s keyValueStore) Get(filePath string) ([]byte, error) {
	kvPair, err := s.store.Get(path.Join(s.etcdPath, filePath))
	if err != nil {
		if err, ok := err.(*client.ClusterError); ok {
			log.Errorf(err.Detail())
		}
		if err == store.ErrKeyNotFound {
			return []byte{}, nil
		}
		return nil, err
	}
	return kvPair.Value, nil
}

func (s keyValueStore) List(filePath string) ([]string, error) {
	kvPairs, err := s.store.List(path.Join(s.etcdPath, filePath))
	if err != nil {
		if err, ok := err.(*client.ClusterError); ok {
			log.Errorf(err.Detail())
		}
		if err == store.ErrKeyNotFound {
			return []string{}, nil
		}
		return nil, err
	}
	keys := []string{}
	for _, pair := range kvPairs {
		keys = append(keys, pair.Key)
	}
	return keys, nil
}

func (s keyValueStore) Put(filePath string, content []byte) error {
	err := s.store.Put(path.Join(s.etcdPath, filePath), content, nil)
	if err != nil {
		if err, ok := err.(*client.ClusterError); ok {
			log.Errorf(err.Detail())
		}
	}
	return err
}

func (s keyValueStore) Delete(filePath string) error {
	err := s.store.Delete(path.Join(s.etcdPath, filePath))
	if err != nil {
		if err, ok := err.(*client.ClusterError); ok {
			log.Errorf(err.Detail())
		}
	}
	return err
}

func (s keyValueStore) SemSignal(filePath string) error {
	return s.Put(filePath, []byte("signal"))
}

// Semaphore returns true if signaled, false if canceled; this is more of a barrier than a semaphore actually...
func (s keyValueStore) SemWait(filePath string, stopCh <-chan struct{}) (bool, error) {
	key := path.Join(s.etcdPath, filePath)
	err := s.store.Put(key, []byte("waiting"), nil)
	if err != nil {
		return false, err
	}

	events, err := s.store.Watch(key, stopCh)
	if err != nil {
		return false, err
	}

	// wait for either the first change or a cancellation
	for {
		select {
		case event := <-events:
			if bytes.Equal(event.Value, []byte("signal")) {
				return true, nil
			}
			// else keep waiting
		case <-stopCh:
			return false, nil
		}
	}

	// this case will never be triggered
	return false, nil
}

func (s keyValueStore) Lock(filePath string, value []byte, ttl time.Duration) error {
	lock, err := s.store.NewLock(path.Join(s.etcdPath, filePath), &store.LockOptions{Value: value, TTL: ttl})
	if err != nil {
		if err, ok := err.(*client.ClusterError); ok {
			log.Errorf(err.Detail())
		}
	}
	_, err = lock.Lock(nil)
	return err
}

func (s keyValueStore) Ping() error {
	return dtrutil.Poll(500*time.Millisecond, 40, func() error {
		log.Infof("Waiting for etcd...")
		_, err := s.store.Exists("/")
		if err, ok := err.(*client.ClusterError); ok {
			log.Debugf("Failed to connect to etcd: %s", strings.Replace(err.Detail(), "\n", ";", -1))
		}
		return err
	})
}
