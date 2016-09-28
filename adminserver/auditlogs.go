package adminserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/docker/dhe-deploy"
	"github.com/docker/dhe-deploy/adminserver/util"

	"github.com/codegangsta/negroni"
	"github.com/docker/distribution/manifest/schema1"
	"github.com/docker/distribution/notifications"
	"github.com/gorilla/mux"
)

type SyslogWriter interface {
	Info(...interface{})
	Debug(...interface{})
	Warning(...interface{})
	Error(...interface{})
}

type auditLogsHandler struct {
	Writer SyslogWriter
}

func (a *auditLogsHandler) logAPICall() negroni.Handler {
	return negroni.HandlerFunc(func(writer http.ResponseWriter, request *http.Request, next http.HandlerFunc) {
		user := util.GetAuthenticatedUser(request)
		var username string
		if user != nil {
			username = user.Account.Name
		}
		logLine := fmt.Sprintf("Call of type %s to %s from %s by %s", request.Method, request.RequestURI, request.RemoteAddr, username)
		a.Writer.Info(logLine)
		next(writer, request)
	})
}

func (a *auditLogsHandler) getEventsHandler(writer http.ResponseWriter, request *http.Request) {
	// This header is filtered out by the nginx reverse proxy
	if request.Header.Get(deploy.RegistryEventsHeaderName) != "true" {
		writeJSONError(writer, errors.New(http.StatusText(http.StatusForbidden)), http.StatusForbidden)
		return
	}
	eventsEnvelope := new(notifications.Envelope)
	if err := json.NewDecoder(request.Body).Decode(eventsEnvelope); err != nil {
		writeJSONError(writer, err, http.StatusBadRequest)
		return
	}

	for _, event := range eventsEnvelope.Events {
		if event.Target.Descriptor.MediaType != schema1.MediaTypeManifest {
			continue
		}

		target := event.Target.Repository

		lastSlashIndex := strings.LastIndex(event.Target.URL, "/")
		if lastSlashIndex >= 0 {
			target += ":" + event.Target.URL[lastSlashIndex+1:]
		}

		logLine := fmt.Sprintf("%s %s %s from %s", event.Actor.Name, actionString(event.Action), target, event.Request.Addr)

		fmt.Fprintln(writer, logLine)
		a.Writer.Info(logLine)
	}
}

func actionString(action string) string {
	action = strings.ToLower(action)
	switch action {
	case "pull":
		return "pulled from"
	case "push":
		return "pushed to"
	default:
		if action[len(action)-1] == 'e' {
			return action + "d"
		}
		return action + "ed"
	}
}

func (a auditLogsHandler) WireSubroutes(router *mux.Router) {
	router.HandleFunc("/", a.getEventsHandler).Methods("POST")
}
