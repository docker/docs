package garant

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/docker/distribution/context"
	"github.com/docker/garant/auth"
	"github.com/docker/garant/auth/common"
	"github.com/gorilla/mux"
)

var clientIDRegexp = regexp.MustCompile("[ -~]+")

type getTokenResponse struct {
	Token        string    `json:"token"`
	AccessToken  string    `json:"access_token"`
	ExpiresIn    int       `json:"expires_in"`
	IssuedAt     time.Time `json:"issued_at"`
	RefreshToken string    `json:"refresh_token,omitempty"`
}

type postTokenResponse struct {
	AccessToken  string    `json:"access_token"`
	Scope        string    `json:"scope"`
	ExpiresIn    int       `json:"expires_in"`
	IssuedAt     time.Time `json:"issued_at"`
	RefreshToken string    `json:"refresh_token,omitempty"`
}

// registerHandlers registers HTTP handler functions on the given router for
// getting basic account info and getting a signed JSON Web Token. Registered
// URL paths will begin with the given prefix.
func (app *App) registerHandlers(prefix string, router *mux.Router) {
	for _, route := range []struct {
		path    string
		method  string
		handler http.Handler
	}{
		{
			// TODO: remove this once no one refers to it anymore.
			// Use `/<prefix>/account_info` instead.
			path:    "/v2/",
			method:  "GET",
			handler: app.handlerWithContext(app.getAccountInfo),
		},
		{
			// TODO: remove this once no one refers to it anymore.
			// Use `/<prefix>/token` instead.
			path:    "/v2/token/",
			method:  "GET",
			handler: app.handlerWithContext(app.getToken),
		},
		{
			path:    path.Join("/", prefix, "account_info"),
			method:  "GET",
			handler: app.handlerWithContext(app.getAccountInfo),
		},
		{
			path:    path.Join("/", prefix, "token"),
			method:  "GET",
			handler: app.handlerWithContext(app.getToken),
		},
		{
			path:    path.Join("/", prefix, "token"),
			method:  "POST",
			handler: app.handlerWithContext(app.postToken),
		},
		{
			path:    path.Join("/", prefix, "__cause_panic"),
			method:  "GET",
			handler: app.handlerWithContext(causePanic),
		},
		{
			path:    path.Join("/", prefix, "__health_check"),
			method:  "GET",
			handler: app.handlerWithContext(app.healthCheck),
		},
	} {
		router.Path(route.path).Methods(route.method).Handler(route.handler)
	}
}

// handlerWithContext wraps the given context-aware handler by setting up the
// request context from this app's baseContext. It also handles logging panics
// and re-panicing so that a wrapping bugsnag reporting notifier may notify of
// the panic as well.
func (app *App) handlerWithContext(handler func(context.Context, http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithRequest(app.baseContext, r)
		logger := context.GetRequestLogger(ctx)
		ctx = context.WithLogger(ctx, logger)

		defer func() {
			err := recover()
			if err == nil {
				// Not currently panicking.
				return
			}

			var stack []byte
			if stacker, ok := err.(common.Stacker); ok {
				stack = stacker.Stack()
			} else {
				stack = debug.Stack()
			}

			// Push the stacktrace onto the context.
			ctx = context.WithValue(ctx, "stackTrace", string(stack))
			ctx = context.WithLogger(ctx, context.GetLogger(ctx, "stackTrace"))

			// Write a simple error response to the client.
			ctx, w = context.WithResponseWriter(ctx, w)

			body := map[string]string{
				"details": "internal error",
			}

			if requestID, ok := ctx.Value("http.request.id").(string); ok {
				body["request_id"] = requestID
			}

			if app.reportingOpts.Bugsnag.APIKey != "" {
				body["message"] = "an administrator has been notified of the issue"
			} else {
				body["message"] = "please contact an administrator to resolve this issue"
			}

			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(body)

			context.GetResponseLogger(ctx).Errorf("runtime panic: %s", err)

			// Re-panic so that the bugsnag handler (if configured) can notify
			// about the panic as well. The panic will finally be suppressed by
			// the outer-most handler.
			panic(err)
		}()

		handler(ctx, w, r)
	})
}

func contextWithAccountSubject(ctx context.Context, account auth.Account) context.Context {
	var acctSubject string
	if account != nil {
		acctSubject = account.Subject()
	}

	ctx = context.WithValue(ctx, "acctSubject", acctSubject)
	ctx = context.WithLogger(ctx, context.GetLogger(ctx, "acctSubject"))

	return ctx
}

func getAuthChallengeOrPanic(err error) auth.Challenge {
	if err != nil {
		if challenge, ok := err.(auth.Challenge); ok {
			return challenge
		}

		panic(err) // Let the app's panic/error handler/wrapper handle it.
	}

	return nil
}

// getAccountInfo handles rendering basic info about the authenticated account.
func (app *App) getAccountInfo(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	context.GetLogger(ctx).Info("getAccountInfo")

	account, err := app.authorizer.Authenticate(ctx, r)
	challenge := getAuthChallengeOrPanic(err)

	// Get response context.
	ctx, w = context.WithResponseWriter(ctx, w)

	if challenge != nil {
		challenge.ServeHTTP(w, r)

		context.GetResponseLogger(ctx).Info("authentication challenged")

		return
	}

	ctx = contextWithAccountSubject(ctx, account)
	context.GetLogger(ctx).Info("authenticated client")

	var info map[string]interface{}
	if account != nil {
		info = account.Info()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"accountInfo": info})

	context.GetResponseLogger(ctx).Info("getAccountInfo complete")
}

// getToken handles authenticating the request and authorizing access to the
// requested scopes.
func (app *App) getToken(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	context.GetLogger(ctx).Info("getToken")

	account, err := app.authorizer.Authenticate(ctx, r)
	challenge := getAuthChallengeOrPanic(err)

	if challenge != nil {
		// Get response context.
		ctx, w = context.WithResponseWriter(ctx, w)

		challenge.ServeHTTP(w, r)

		context.GetResponseLogger(ctx).Info("authentication challenged")

		return
	}

	ctx = contextWithAccountSubject(ctx, account)
	context.GetLogger(ctx).Info("authenticated client")

	params := r.URL.Query()
	service := params.Get("service")
	clientID := params.Get("client_id")
	scopeSpecifiers := params["scope"]

	requestedAccessList := resolveScopeSpecifiers(scopeSpecifiers)

	ctx = context.WithValue(ctx, "requestedAccess", requestedAccessList)
	ctx = context.WithLogger(ctx, context.GetLogger(ctx, "requestedAccess"))

	grantedAccessList, err := app.authorizer.Authorize(ctx, account, service, requestedAccessList...)
	if err != nil {
		panic(err) // Let the app's panic/error handler/wrapper handle it.
	}

	ctx = context.WithValue(ctx, "grantedAccess", grantedAccessList)
	ctx = context.WithLogger(ctx, context.GetLogger(ctx, "grantedAccess"))

	expiresIn := 5 * time.Minute
	token := app.CreateJWT(account, service, grantedAccessList, expiresIn)

	context.GetLogger(ctx).Info("authorized client")

	// Get response context.
	ctx, w = context.WithResponseWriter(ctx, w)

	response := getTokenResponse{
		Token:       token,
		AccessToken: token,
		ExpiresIn:   int(expiresIn.Seconds()),
		IssuedAt:    time.Now(),
	}

	if offlineToken, _ := strconv.ParseBool(params.Get("offline_token")); offlineToken {
		if tokenAuthenticator, ok := app.authorizer.(auth.TokenAuthorizer); ok {
			refreshToken, err := tokenAuthenticator.GetToken(ctx, account, clientID)
			if err != nil {
				panic(err)
			}
			response.RefreshToken = refreshToken.Token
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	context.GetResponseLogger(ctx).Info("getToken complete")
}

// postToken handles authenticating the request and authorizing access to the
// reqested scopes for an oauth-compatible flow.
func (app *App) postToken(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	tokenAuthenticator, ok := app.authorizer.(auth.TokenAuthorizer)
	if !ok {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	context.GetLogger(ctx).Info("postToken")

	grantType := r.FormValue("grant_type")
	service := r.FormValue("service")
	clientID := r.FormValue("client_id")
	accessType := r.FormValue("access_type")
	scopeSpecifiers := strings.Split(r.FormValue("scope"), " ")

	if !clientIDRegexp.MatchString(clientID) {
		common.ErrInvalidClientID.ServeHTTP(w, r)
		return
	}

	var account auth.Account
	var err error

	var refreshToken string

	switch grantType {
	case "refresh_token":
		refreshToken = r.FormValue("refresh_token")
		account, err = tokenAuthenticator.AuthenticateWithToken(ctx, refreshToken)
	case "password":
		username := r.FormValue("username")
		password := r.FormValue("password")
		account, err = tokenAuthenticator.AuthenticateWithPassword(ctx, username, password)
	default:
		err = common.ErrInvalidGrantType
	}

	challenge := getAuthChallengeOrPanic(err)

	if challenge != nil {
		// Get response context.
		ctx, w = context.WithResponseWriter(ctx, w)

		challenge.ServeHTTP(w, r)

		context.GetResponseLogger(ctx).Info("authentication challenged")

		return
	}

	requestedAccessList := resolveScopeSpecifiers(scopeSpecifiers)
	grantedAccessList, err := app.authorizer.Authorize(ctx, account, service, requestedAccessList...)
	if err != nil {
		panic(err) // Let the app's panic/error handler/wrapper handle it.
	}

	ctx = context.WithValue(ctx, "requestedAccess", requestedAccessList)
	ctx = context.WithLogger(ctx, context.GetLogger(ctx, "requestedAccess"))

	ctx = context.WithValue(ctx, "grantedAccess", grantedAccessList)
	ctx = context.WithLogger(ctx, context.GetLogger(ctx, "grantedAccess"))

	expiresIn := 5 * time.Minute
	accessToken := app.CreateJWT(account, service, grantedAccessList, expiresIn)

	context.GetLogger(ctx).Info("authorized client")

	accessStrings := make([]string, len(grantedAccessList))
	for i, access := range grantedAccessList {
		accessStrings[i] = fmt.Sprintf("%s:%s:%s", access.Type, access.Name, access.Action)
	}
	scope := strings.Join(accessStrings, " ")

	// Get response context.
	ctx, w = context.WithResponseWriter(ctx, w)

	response := postTokenResponse{
		AccessToken: accessToken,
		Scope:       scope,
		ExpiresIn:   int(expiresIn.Seconds()),
		IssuedAt:    time.Now(),
	}

	if accessType == "offline" {
		if refreshToken == "" {
			token, err := tokenAuthenticator.GetToken(ctx, account, clientID)
			if err != nil {
				panic(err)
			}
			refreshToken = token.Token
		}
		response.RefreshToken = refreshToken
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	context.GetResponseLogger(ctx).Info("postToken complete")
}

func causePanic(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	message := params.Get("msg")
	if message == "" {
		message = "intentional panic"
	}

	panic(message)
}

// healthCheck handles reporting that this token signing server is running "OK"
// and also includes the expiration date of the key's leaf certificate (if
// there is one).
func (app *App) healthCheck(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	status := "OK"
	certExpires := "never"

	cert, err := app.getSigningKeyCert()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"details": fmt.Sprintf("unable to get signing key certificate: %s", err),
		})
		return
	}

	if cert != nil {
		certExpires = cert.NotAfter.Format(time.RFC3339)
	}

	json.NewEncoder(w).Encode(map[string]string{
		"status":       status,
		"authBackend":  app.config.Auth.BackendName,
		"signingKeyID": app.signingKey.KeyID(),
		"certExpires":  certExpires,
	})
}

func (app *App) getSigningKeyCert() (cert *x509.Certificate, err error) {
	x5cVal := app.signingKey.GetExtendedField("x5c")
	if x5cVal == nil {
		// There is no certificate list.
		return nil, nil
	}

	x5cValSlice, ok := x5cVal.([]interface{})
	if !ok {
		return nil, fmt.Errorf("x5c header value is not a list: %#v", x5cVal)
	}

	if len(x5cValSlice) == 0 {
		return nil, fmt.Errorf("x5c header list is empty: %#v", x5cValSlice)
	}

	certString, ok := x5cValSlice[0].(string)
	if !ok {
		return nil, fmt.Errorf("x5c header list entry is not a string: %#v", x5cValSlice[0])
	}

	certBytes, err := base64.StdEncoding.DecodeString(certString)
	if err != nil {
		return nil, fmt.Errorf("unable to decode x509 certificate: %s", err)
	}

	cert, err = x509.ParseCertificate(certBytes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse x509 certificate: %s", err)
	}

	return cert, nil
}
