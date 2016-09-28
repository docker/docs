package adminserver

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/docker/dhe-deploy/adminserver/util"
	configutil "github.com/docker/dhe-deploy/hubconfig/util"
	"github.com/docker/dhe-deploy/shared/dtrutil"
)

func (a *AdminServer) enziProxyHandler(writer http.ResponseWriter, request *http.Request) {
	// XXX: maybe don't request the enzi url every time???
	haConfig, err := a.settingsStore.HAConfig()
	if err != nil {
		writeJSONError(writer, err, http.StatusInternalServerError)
		return
	}

	user := util.GetAuthenticatedUser(request)
	if user == nil || user.IsAnonymous {
		// this shouldn't be possible because we require auth in middleware
		writeJSONError(writer, fmt.Errorf("not authenticted???"), http.StatusInternalServerError)
		return
	}

	enziConfig := configutil.GetEnziConfig(haConfig)
	httpClient, err := dtrutil.HTTPClient(!enziConfig.VerifyCert, enziConfig.CA)
	enziURL, err := url.Parse(fmt.Sprintf("https://%s", enziConfig.Host))
	if err != nil {
		writeJSONError(writer, err, http.StatusInternalServerError)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(enziURL)
	proxy.Transport = httpClient.Transport

	// we use url.Parse so we can manipulate each part of the URL separately rather than
	// it being opaque
	request.URL, err = url.Parse(request.URL.String())
	if err != nil {
		writeJSONError(writer, err, http.StatusInternalServerError)
		return
	}
	// we remove the prefix on our end in all cases
	// we optionall add the prefix depending on whether we are running on ucp or not
	request.URL.Path = fmt.Sprintf("%s%s", enziConfig.Prefix, strings.TrimPrefix(request.URL.Path, "/enzi"))
	user.Token.AuthenticateRequest(request)
	proxy.ServeHTTP(writer, request)
}
