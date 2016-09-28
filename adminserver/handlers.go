package adminserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/docker/dhe-deploy"
	"github.com/docker/dhe-deploy/adminserver/healthcheck"
	"github.com/docker/dhe-deploy/adminserver/util"
	"github.com/docker/dhe-deploy/manager/versions"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"github.com/samalba/dockerclient"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (a *AdminServer) serveHTML(writer http.ResponseWriter, request *http.Request) {
	var err error
	var alerts = []alert{}

	user := util.GetAuthenticatedUser(request)

	if user != nil && !user.IsAnonymous {
		alerts = a.globalAlerts(request)
	}

	var userJS []byte
	if user != nil {
		userJS, err = json.Marshal(user.Account)
	}
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to marshal user information")
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	templateArgs := struct {
		CurrentVersion string
		ContainerID    string
		Alerts         []alert
		User           template.JS
		IsAdmin        bool
	}{
		CurrentVersion: deploy.Version,
		Alerts:         alerts,
		User:           template.JS(userJS),
		IsAdmin:        user != nil && *user.Account.IsAdmin,
	}

	if err := getTemplate().Execute(writer, templateArgs); err != nil {
		log.WithFields(log.Fields{
			"error":        err,
			"templateArgs": templateArgs,
		}).Error("Failed to render template")
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}

func (a *AdminServer) eventWsHandler(writer http.ResponseWriter, request *http.Request) {
	user := util.GetAuthenticatedUser(request)

	conn, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Errorf("Error upgrading websocket: %+v", err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	a.eventWebSocketManager.NewWebSocketClient(conn, a.eventWebSocketManager, user)
}

func (a *AdminServer) caHandler(writer http.ResponseWriter, request *http.Request) {
	hubConfig, err := a.settingsStore.UserHubConfig()
	if err != nil {
		log.WithField("error", err).Error("Failed to retrieve hub config")
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(writer, hubConfig.WebTLSCA)
}

func (a *AdminServer) healthHandler(writer http.ResponseWriter, request *http.Request) {
	out := struct {
		Error   string
		Healthy bool
	}{
		Error:   "",
		Healthy: true,
	}
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	// by default we do a health check when we are asked to instead of getting it from
	// the cache thing. We should remove the cache thing in the future somehow
	if request.URL.Query().Get("cached") != "" {
		verdict, err := a.propertyManager.Get(fmt.Sprintf("%s%s", deploy.ReplicaHealthPropertyPrefix, os.Getenv(deploy.ReplicaIDEnvVar)))
		if err != nil {
			verdict = fmt.Sprintf("Failed to get health check result: %s", err)
		}

		if verdict != "OK" {
			out.Healthy = false
			out.Error = verdict
		}
	} else {
		err := healthcheck.HealthCheckAll(a.etcdClient, a.notaryClient, a.rethinkSession)

		if err != nil {
			out.Healthy = false
			out.Error = err.Error()
		}
	}
	if !out.Healthy {
		writer.WriteHeader(http.StatusServiceUnavailable)
	}
	output, _ := json.Marshal(out)
	writer.Write([]byte(output))
}

func (a *AdminServer) infoHandler(writer http.ResponseWriter, request *http.Request) {
	// after DTR 2.0 the only thing we are keeping here for backwards compatibililty is the DTR version info
	type DTR struct {
		Version string
		GitSHA  string
	}
	infoStruct := struct {
		DTR DTR
	}{
		DTR: DTR{
			Version: deploy.Version,
			GitSHA:  deploy.GitSHA,
		},
	}

	dheVersionBytes, err := json.MarshalIndent(infoStruct, "", "  ")
	if err != nil {
		log.WithField("error", err).Error("Failed to marshal info struct")
		writeJSONError(writer, err, http.StatusInternalServerError)
	}
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.Write(dheVersionBytes)
}

func (a *AdminServer) getUpgradeHandler(writer http.ResponseWriter, request *http.Request) {
	var response struct {
		UpgradeAvailable bool     `json:"upgradeAvailable,omitempty"`
		NewVersion       string   `json:"newVersion,omitempty"`
		NewPatch         string   `json:"newPatch,omitempty"`
		VersionsList     []string `json:"versions,omitempty"`
	}
	versionsList, err := a.versionChecker.VersionList(nil)
	if err != nil {
		log.WithField("error", err).Warn("Failed to determine newest manager version")
		writeJSONError(writer, err, http.StatusBadGateway)
		return
	}

	newestVersion := versionsList[len(versionsList)-1]
	response.UpgradeAvailable = versions.Less(deploy.Version, newestVersion)
	response.NewVersion = newestVersion
	response.VersionsList = versionsList
	response.NewPatch = versionsList.PatchFor(deploy.Version)

	writeJSON(writer, response)
}

func (a *AdminServer) dockerHubLoginHandler(writer http.ResponseWriter, request *http.Request) {
	var authConfig dockerclient.AuthConfig
	defer request.Body.Close()
	if err := json.NewDecoder(request.Body).Decode(&authConfig); err != nil {
		log.WithField("error", err).Warn("Failed to parse login request body")
		writeJSONError(writer, err, http.StatusBadRequest)
		return
	}
	authConfig.Email = deploy.DummyHubUserEmail
	if _, err := a.versionChecker.NewestVersion(&authConfig); err != nil {
		log.WithField("error", err).Warn("Failed to get newest manager version with new login credentials")
		writeJSONError(writer, errors.New("Username/password combination invalid"), http.StatusBadRequest)
		return
	}

	if err := a.settingsStore.SetHubCredentials(&authConfig); err != nil {
		log.WithField("error", err).Error("Failed to save new docker credentials")
		writeJSONError(writer, errors.New("Failed to save new docker credentials"), http.StatusInternalServerError)
		return
	}

	writeJSON(writer, nil)
}

func (a *AdminServer) getVersionHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprint(writer, deploy.Version)
}

func (a *AdminServer) clientAnalyticsHandler(writer http.ResponseWriter, request *http.Request) {
	useragent := request.UserAgent()
	sourceIP := request.RemoteAddr
	a.sendClientAnalytics(useragent, sourceIP)
}
