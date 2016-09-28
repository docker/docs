package adminserver

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/docker/dhe-deploy"
	"github.com/docker/dhe-deploy/adminserver/util"
	"github.com/docker/dhe-deploy/shared/containers"
)

func (a *AdminServer) rethinkProxyHandler(writer http.ResponseWriter, request *http.Request) {
	user := util.GetAuthenticatedUser(request)
	if user == nil || user.IsAnonymous {
		// this shouldn't be possible because we require auth in middleware
		writeJSONError(writer, fmt.Errorf("not authenticted???"), http.StatusInternalServerError)
		return
	}
	if !*user.Account.IsAdmin {
		writeJSONError(writer, fmt.Errorf("Not admin"), http.StatusForbidden)
		return
	}

	rethinkURL, err := url.Parse(fmt.Sprintf("http://%s:%d", containers.Rethinkdb.BridgeNameLocalReplica(), deploy.RethinkAdminPort))
	if err != nil {
		writeJSONError(writer, err, http.StatusInternalServerError)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(rethinkURL)

	// we use url.Parse so we can manipulate each part of the URL separately rather than
	// it being opaque
	request.URL, err = url.Parse(request.URL.String())
	if err != nil {
		writeJSONError(writer, err, http.StatusInternalServerError)
		return
	}
	// we remove the prefix on our end in all cases
	request.URL.Path = strings.TrimPrefix(request.URL.Path, "/db")
	user.Token.AuthenticateRequest(request)
	proxy.ServeHTTP(writer, request)
}
