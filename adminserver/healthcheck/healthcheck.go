package healthcheck

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/docker/dhe-deploy"
	"github.com/docker/dhe-deploy/hubconfig"
	"github.com/docker/dhe-deploy/manager/schema/interfaces"
	"github.com/docker/dhe-deploy/shared/containers"
	"github.com/docker/dhe-deploy/shared/dtrutil"
	rethink "gopkg.in/dancannon/gorethink.v2"
)

type HealthCheckJob struct {
	settingsStore   hubconfig.SettingsStore
	kvStore         hubconfig.KeyValueStore
	propertyManager interfaces.PropertyManager
	etcdClient      *http.Client
	notaryClient    *http.Client
	logger          *logrus.Logger
	rethinkSession  *rethink.Session
}

// the health check job reports only the status of the current pod
func NewJob(settingsStore hubconfig.SettingsStore, kvStore hubconfig.KeyValueStore, etcdClient, notaryClient *http.Client, propertyManager interfaces.PropertyManager, rethinkSession *rethink.Session, logger *logrus.Logger) (*HealthCheckJob, error) {
	job := &HealthCheckJob{
		settingsStore:   settingsStore,
		kvStore:         kvStore,
		propertyManager: propertyManager,
		etcdClient:      etcdClient,
		notaryClient:    notaryClient,
		logger:          logger,
		rethinkSession:  rethinkSession,
	}
	return job, nil
}

func RethinkdbHealth(session *rethink.Session) error {
	cursor, err := rethink.DB("rethinkdb").Table("current_issues").Run(session)
	if err != nil {
		return fmt.Errorf("There was an error reading rethinkdb current_issues table: %s", err)
	}
	info := []interface{}{}
	if err := cursor.All(&info); err != nil {
		return fmt.Errorf("There was an error reading rethinkdb current_issues table: %s", err)
	}
	if len(info) > 0 {
		data, err := json.Marshal(info)
		if err != nil {
			return fmt.Errorf("Failed to serialize error: %s; The data was: %v", err, info)
		}
		return fmt.Errorf("Rethinkdb has issues: %s", string(data))
	}
	return nil
}

func LocalEtcdHealth(etcdClient *http.Client) error {
	replicaID := os.Getenv(deploy.ReplicaIDEnvVar)
	etcdHealthURL := fmt.Sprintf("https://%s:%d/health", containers.Etcd.BridgeName(replicaID), containers.EtcdClientPort1)

	req, err := http.NewRequest("GET", etcdHealthURL, nil)
	if err != nil {
		return err
	}

	res, err := dtrutil.DoRequestWithClient(req, etcdClient)
	if err != nil {
		return fmt.Errorf("Failed to health check etcd: %s", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected status code from etcd: %s", res.Status)
	}

	var health struct {
		Health string `json:"health"`
	}

	if err := json.NewDecoder(res.Body).Decode(&health); err != nil {
		return fmt.Errorf("Error trying to decode health check response from etcd: %s", err)
	}

	if health.Health != "true" {
		logrus.WithField("health", health).Errorf("Etcd health check failed: %s", health.Health)
		return fmt.Errorf("Etcd unhealthy: %s", health.Health)
	}
	return nil
}

// LocalNotaryHealth checks the notary health endpoint and parses out the reasons Notary is unhealthy
// if the response is anything other than a 200
func LocalNotaryHealth(notaryClient *http.Client) error {
	replicaID := os.Getenv(deploy.ReplicaIDEnvVar)
	notaryHealthURL := fmt.Sprintf("https://%s:%d/_notary_server/health", containers.NotaryServer.BridgeName(replicaID), deploy.NotaryServerHTTPPort)

	req, err := http.NewRequest("GET", notaryHealthURL, nil)
	if err != nil {
		return err
	}

	res, err := dtrutil.DoRequestWithClient(req, notaryClient)
	if err != nil {
		return fmt.Errorf("Failed to health check notary: %s", err)
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusServiceUnavailable:
		reasons := make(map[string]string)
		if err := json.NewDecoder(res.Body).Decode(&reasons); err != nil {
			return fmt.Errorf("Notary unhealthy, but error trying to decode unhealthy response from notary: %s", err)
		}

		reasonStrings := make([]string, 0, len(reasons))
		for _, reasonStr := range reasons {
			reasonStrings = append(reasonStrings, reasonStr)
		}

		return fmt.Errorf("Notary unhealthy: %s", strings.Join(reasonStrings, ";"))
	default:
		return fmt.Errorf("Unexpected status code from notary: %s", res.Status)
	}
}

func HealthCheckLocalReplica() error {
	replicaID := os.Getenv(deploy.ReplicaIDEnvVar)

	// the admin server doesn't check its own endpoint because that would be silly
	endpoints := []struct {
		host string
		port uint16
	}{
		{
			host: containers.Registry.BridgeName(replicaID),
			port: deploy.StorageContainerPort,
		},
		{
			host: containers.Registry.BridgeName(replicaID),
			port: deploy.GarantPort,
		},
		{
			host: containers.Rethinkdb.BridgeName(replicaID),
			port: deploy.RethinkdbPort,
		},
		{
			host: containers.Etcd.BridgeName(replicaID),
			port: containers.EtcdClientPort1,
		},
		{
			host: containers.Etcd.BridgeName(replicaID),
			port: containers.EtcdClientPort2,
		},
		{
			host: containers.NotaryServer.BridgeName(replicaID),
			port: deploy.NotaryServerHTTPPort,
		},
		{
			host: containers.NotarySigner.BridgeName(replicaID),
			port: deploy.NotarySignerGRPCPort,
		},
	}

	for _, endpoint := range endpoints {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", endpoint.host, endpoint.port))
		if err != nil {
			return fmt.Errorf("Failed to connect to container at %s:%d: %s", endpoint.host, endpoint.port, err)
		}
		conn.Close()
	}
	return nil
}

func HealthCheckAll(etcdClient, notaryClient *http.Client, rethinkSession *rethink.Session) error {
	err := HealthCheckLocalReplica()
	if err != nil {
		return err
	}
	err = LocalEtcdHealth(etcdClient)
	if err != nil {
		return err
	}
	err = RethinkdbHealth(rethinkSession)
	if err != nil {
		return err
	}
	err = LocalNotaryHealth(notaryClient)
	if err != nil {
		return err
	}
	return nil
}

// run every minute
func (j *HealthCheckJob) IsReady() bool {
	return true
}

func (j *HealthCheckJob) Run(manual bool) error {
	verdict := "OK"
	err := HealthCheckAll(j.etcdClient, j.notaryClient, j.rethinkSession)
	if err != nil {
		verdict = err.Error()
	}

	err = j.propertyManager.Set(fmt.Sprintf("%s%s", deploy.ReplicaHealthPropertyPrefix, os.Getenv(deploy.ReplicaIDEnvVar)), verdict)
	if err != nil {
		return err
	}
	err = j.propertyManager.Set(fmt.Sprintf("%s%s", deploy.ReplicaHealthTimestampPropertyPrefix, os.Getenv(deploy.ReplicaIDEnvVar)), time.Now().Format(time.RFC3339))
	if err != nil {
		return err
	}

	return nil
}
