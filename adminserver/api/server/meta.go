package server

import (
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/coreos/etcd/client"
	"github.com/docker/dhe-deploy"
	"github.com/docker/dhe-deploy/adminserver/api/common/errors"
	"github.com/docker/dhe-deploy/adminserver/api/common/forms"
	"github.com/docker/dhe-deploy/adminserver/api/common/responses"
	"github.com/docker/dhe-deploy/bootstrap"
	"github.com/docker/dhe-deploy/hubconfig"
	"github.com/docker/dhe-deploy/hubconfig/util"
	"github.com/docker/dhe-deploy/shared/containers"
	"github.com/docker/distribution/context"
	"github.com/emicklei/go-restful"
	rethink "gopkg.in/dancannon/gorethink.v2"
)

func (a *APIServer) getAdminSettingsHandler(ctx context.Context, r *restful.Request) responses.APIResponse {
	rd := newRequestData(a, ctx, r)
	if rd.addFilters(
		makeFilterGetAuthenticatedUser(true),
		ensureIsAdmin,
	).evaluateFilters(); rd.errResponse != nil {
		return rd.errResponse
	}

	hubConfig, err := a.settingsStore.UserHubConfig()
	if err != nil {
		log.WithField("error", err).Error("Failed to retrieve hub config")
		return responses.APIError(errors.InternalError(ctx, err))
	}

	haConfig, err := a.settingsStore.HAConfig()
	if err != nil {
		log.WithField("error", err).Error("Failed to retrieve ha config")
		return responses.APIError(errors.InternalError(ctx, err))
	}
	safeHAConfig := maskPasswordsInHAConfig(haConfig)

	replicaSettings := map[string]responses.ReplicaSettings{}
	for id, settings := range haConfig.ReplicaConfig {
		replicaSettings[id] = responses.ReplicaSettings{
			HTTPPort:  settings.HTTPPort,
			HTTPSPort: settings.HTTPSPort,
			Node:      settings.Node,
		}
	}

	settings := &responses.Settings{
		DTRHost:            hubConfig.DTRHost,
		AuthBypassCA:       hubConfig.AuthBypassCA,
		AuthBypassOU:       hubConfig.AuthBypassOU,
		DisableUpgrades:    hubConfig.DisableUpgrades,
		ReportAnalytics:    hubConfig.ReportAnalytics,
		AnonymizeAnalytics: hubConfig.AnonymizeAnalytics,
		ReleaseChannel:     &hubConfig.ReleaseChannel,
		WebTLSCert:         hubConfig.WebTLSCert,
		WebTLSCA:           hubConfig.WebTLSCA,
		GCMode:             hubConfig.GCMode,
		// immutable settings (changed from the bootstrapper)
		ReplicaSettings:       replicaSettings,
		HTTPProxy:             safeHAConfig.HTTPProxy,
		HTTPSProxy:            safeHAConfig.HTTPSProxy,
		NoProxy:               safeHAConfig.NoProxy,
		LogProtocol:           safeHAConfig.LogProtocol,
		LogHost:               safeHAConfig.LogHost,
		LogLevel:              safeHAConfig.LogLevel,
		EtcdHeartbeatInterval: safeHAConfig.EtcdHeartbeatInterval,
		EtcdElectionTimeout:   safeHAConfig.EtcdElectionTimeout,
		EtcdSnapshotCount:     safeHAConfig.EtcdSnapshotCount,
		ReplicaID:             os.Getenv(deploy.ReplicaIDEnvVar),
	}
	return responses.JSONResponse(http.StatusOK, nil, nil, settings)
}

func (a *APIServer) updateAdminSettingsHandler(ctx context.Context, r *restful.Request) responses.APIResponse {
	rd := newRequestData(a, ctx, r)
	if rd.addFilters(
		makeFilterGetAuthenticatedUser(true),
		ensureIsAdmin,
	).evaluateFilters(); rd.errResponse != nil {
		return rd.errResponse
	}

	settings := forms.Settings{}
	if err := json.NewDecoder(r.Request.Body).Decode(&settings); err != nil {
		return responses.APIError(errors.MakeError(errors.ErrorCodeInvalidJSON, err))
	}

	hubConfig, err := a.settingsStore.UserHubConfig()
	if err != nil {
		log.WithField("error", err).Error("creating new hub config")
	}
	if hubConfig == nil {
		hubConfig = new(hubconfig.UserHubConfig)
	}

	needToReregister := false

	// these must be taken into account before the domain name is saved in case they don't match the domain name
	// the web certs will be validated when we try to save the hubConfig
	if settings.WebTLSCert != nil {
		hubConfig.WebTLSCert = *settings.WebTLSCert
	}
	if settings.WebTLSKey != nil {
		hubConfig.WebTLSKey = *settings.WebTLSKey
	}
	if settings.WebTLSCA != nil && *settings.WebTLSCA != "" && *settings.WebTLSCA != hubConfig.WebTLSCA {
		needToReregister = true
		hubConfig.WebTLSCA = *settings.WebTLSCA
	}

	if settings.DTRHost != nil && *settings.DTRHost != "" && *settings.DTRHost != hubConfig.DTRHost {
		needToReregister = true
		hubConfig.DTRHost = *settings.DTRHost
		domainName := strings.Split(hubConfig.DTRHost, ":")[0]
		// if you change the domain name, we may need a new cert
		err := util.HubConfigTLSDomainConsistent(hubConfig)
		// if it's no longer consistent for whatever reason, generate a new cert
		if err != nil {
			cert, err := util.GenTLSCert(domainName)
			if err != nil {
				return responses.APIError(errors.InternalError(ctx, err))
			}

			err = util.SetTLSCertificateInHubConfig(hubConfig, cert, cert)
			if err != nil {
				return responses.APIError(errors.InternalError(ctx, err))
			}
		}
	}
	if settings.AuthBypassCA != nil {
		hubConfig.AuthBypassCA = *settings.AuthBypassCA
	}
	if settings.AuthBypassOU != nil {
		hubConfig.AuthBypassOU = *settings.AuthBypassOU
	}
	if settings.DisableUpgrades != nil {
		hubConfig.DisableUpgrades = *settings.DisableUpgrades
	}
	if settings.ReportAnalytics != nil {
		hubConfig.ReportAnalytics = *settings.ReportAnalytics
	}
	if settings.AnonymizeAnalytics != nil {
		hubConfig.AnonymizeAnalytics = *settings.AnonymizeAnalytics
	}
	if settings.ReleaseChannel != nil {
		hubConfig.ReleaseChannel = *settings.ReleaseChannel
	}
	if settings.GCMode != nil {
		hubConfig.GCMode = *settings.GCMode
	}

	// Validate that the cert in AuthBypassCA is parsable if there is one
	if hubConfig.AuthBypassCA != "" {
		certBytes := []byte(hubConfig.AuthBypassCA)
		cp := x509.NewCertPool()
		hasCerts := cp.AppendCertsFromPEM(certBytes)
		if !hasCerts {
			err := fmt.Errorf("Failed to parse Auth Bypass CA.")
			return responses.APIError(errors.InvalidSettings(err))
		}
	}

	if err := a.settingsStore.SetUserHubConfig(hubConfig); err != nil {
		log.WithField("error", err).Error("Failed to update hub config")
		// we don't report this as an internal error because there could be a validation error
		return responses.APIError(errors.InvalidSettings(err))
	}

	// after setting the hubconfig we need to update our registration with enzi if we messed with the certs or domain name (also https port, handled in reconfigure.go)
	if needToReregister {
		err := util.RegisterAuth(rd.user.EnziSession, a.settingsStore)
		if err != nil {
			err := fmt.Errorf("Failed to re-register auth: %s", err)
			return responses.APIError(errors.InvalidSettings(err))
		}
	}

	return responses.JSONResponse(http.StatusAccepted, nil, nil, nil)
}

func (a *APIServer) getClusterStatusHandler(ctx context.Context, r *restful.Request) responses.APIResponse {
	rd := newRequestData(a, ctx, r)
	if rd.addFilters(
		makeFilterGetAuthenticatedUser(true),
		ensureIsAdmin,
	).evaluateFilters(); rd.errResponse != nil {
		return rd.errResponse
	}

	response := responses.ClusterStatus{
		RethinkSystemTables: map[string]interface{}{},
		EtcdStatus:          map[string]interface{}{},
		ReplicaHealth:       map[string]string{},
		ReplicaTimestamp:    map[string]string{},
		ReplicaRORegistry:   map[string]bool{},
	}

	systemTables := []string{"stats", "current_issues", "cluster_config", "db_config", "server_status", "table_config", "table_status"} // logs is not here because it can be big
	for _, table := range systemTables {
		cursor, err := rethink.DB("rethinkdb").Table(table).Run(a.rethinkSession)
		if err != nil {
			err := fmt.Sprintf("There was an error reading rethinkdb table %s: %s", table, err)
			return responses.APIError(errors.MakeError(errors.ErrorCodeStatusCheckFailed, err))
		}

		info := []interface{}{}
		if err := cursor.All(&info); err != nil {
			err := fmt.Sprintf("There was an error reading rethinkdb table %s: %s", table, err)
			return responses.APIError(errors.MakeError(errors.ErrorCodeStatusCheckFailed, err))
		}
		response.RethinkSystemTables[table] = info
	}

	// XXX: we shouldn't be doing this here but we don't have direct access to etcd from the kvstore, also coupling to the bootstrap utils is aweful
	clientURL := fmt.Sprintf("https://%s:%d", containers.Etcd.BridgeName(os.Getenv(deploy.ReplicaIDEnvVar)), containers.EtcdClientPort1)
	etcdClient, err := bootstrap.GetEtcdConn(clientURL)
	if err != nil {
		err := fmt.Sprintf("Couldn't set up etcd client: %s", err)
		return responses.APIError(errors.MakeError(errors.ErrorCodeStatusCheckFailed, err))
	}
	membersAPI := client.NewMembersAPI(etcdClient)
	members, err := membersAPI.List(context.Background())
	if err != nil {
		err := fmt.Sprintf("Couldn't get etcd members list: %s", err)
		return responses.APIError(errors.MakeError(errors.ErrorCodeStatusCheckFailed, err))
	}

	response.EtcdStatus["members"] = members

	haConfig, err := a.settingsStore.HAConfig()
	if err != nil {
		log.WithField("error", err).Error("Failed to retrieve ha config")
		return responses.APIError(errors.InternalError(ctx, err))
	}

	gcLocker, _ := a.kvStore.Get(deploy.GCRunLockPath)
	response.GCLockHolder = string(gcLocker)

	keys, _ := a.kvStore.List(deploy.RegistryROStatePath)

	roRegistriesMap := map[string]struct{}{}
	for _, keys := range keys {
		split := strings.Split(keys, "/")
		replicaID := split[len(split)-1]
		roRegistriesMap[replicaID] = struct{}{}
	}

	for replicaID := range haConfig.ReplicaConfig {
		verdict, err := a.propertyMgr.Get(fmt.Sprintf("%s%s", deploy.ReplicaHealthPropertyPrefix, replicaID))
		if err != nil {
			response.ReplicaHealth[replicaID] = fmt.Sprintf("Failed to get status of replica %s: %s", replicaID, err)
		} else {
			response.ReplicaHealth[replicaID] = verdict
		}

		timestamp, err := a.propertyMgr.Get(fmt.Sprintf("%s%s", deploy.ReplicaHealthTimestampPropertyPrefix, replicaID))
		if err != nil {
			response.ReplicaTimestamp[replicaID] = fmt.Sprintf("Failed to get timestamp of replica %s: %s", replicaID, err)
		} else {
			response.ReplicaTimestamp[replicaID] = timestamp
		}

		if _, ok := roRegistriesMap[replicaID]; ok {
			response.ReplicaRORegistry[replicaID] = true
		} else {
			response.ReplicaRORegistry[replicaID] = false
		}
	}

	return responses.JSONResponse(http.StatusOK, nil, nil, response)
}
