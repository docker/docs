package adminserver

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"runtime"
	"strings"

	"github.com/docker/dhe-deploy/adminserver/util"
	hubconfigutil "github.com/docker/dhe-deploy/hubconfig/util"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	gorillacontext "github.com/gorilla/context"
)

func (a *AdminServer) recoveryMiddleware() negroni.Handler {
	return negroni.HandlerFunc(func(writer http.ResponseWriter, request *http.Request, next http.HandlerFunc) {
		defer func() {
			if err := recover(); err != nil {
				writeJSONError(writer, errors.New(http.StatusText(http.StatusInternalServerError)), http.StatusInternalServerError)
				stack := make([]byte, 1024*8)
				stack = stack[:runtime.Stack(stack, false)]

				log.WithField("error", err).Errorf("PANIC: %s", stack)
			}
		}()

		next(writer, request)
	})
}

// authMiddleware attempts to authenticate the client using any available
// authentication scheme (sesison cookie, basic auth, etc). If authenticaiton
// succeeds, a "user" value is set on the request using gorilla/context.
func (a *AdminServer) authMiddleware() negroni.Handler {
	return negroni.HandlerFunc(func(writer http.ResponseWriter, request *http.Request, next http.HandlerFunc) {
		// for the following endpoints don't authenticate at all:
		if request.URL.Path == "/ca" || request.URL.Path == "/health" {
			next(writer, request)
			return
		}

		if err := a.authenticateRequest(writer, request); err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Authentication failed due to internal error")
			writeJSONError(writer, errors.New(http.StatusText(http.StatusInternalServerError)), http.StatusInternalServerError)
			return
		}

		next(writer, request)
	})
}

// authRedirectMiddleware handles auth redirection for admin and non-admin users:
// for authed users this redirects to the homepage and for non-authed users this
// redirects to the login.
func (a *AdminServer) authRedirectMiddleware() negroni.Handler {
	return negroni.HandlerFunc(func(writer http.ResponseWriter, request *http.Request, next http.HandlerFunc) {
		// for the following endpoints don't authenticate at all:
		if request.URL.Path == "/ca" || request.URL.Path == "/health" {
			next(writer, request)
			return
		}
		user := util.GetAuthenticatedUser(request)

		switch {
		case user != nil && !user.IsAnonymous && (request.URL.Path == "/api/v0/openid/begin"):
			http.Redirect(writer, request, "/", http.StatusFound)
		case user != nil && !user.IsAnonymous && (request.URL.Path == "/logout"):
			haConfig, err := a.settingsStore.HAConfig()
			if err != nil {
				writeJSONError(writer, errors.New(http.StatusText(http.StatusInternalServerError)), http.StatusInternalServerError)
			}
			enziConfig := hubconfigutil.GetEnziConfig(haConfig)
			urlValues := url.Values{}
			urlValues.Add("next", fmt.Sprintf("https://%s", request.Host))

			// Override cookies setting expiration to zero time and maxage to -1
			// to delete cookie
			http.SetCookie(writer, &http.Cookie{
				Name:     "session",
				Value:    "",
				Path:     "/",
				Secure:   true,
				HttpOnly: true,
				MaxAge:   -1,
			})

			http.Redirect(writer, request, fmt.Sprintf("https://%s%s/logout?%s", enziConfig.Host, enziConfig.Prefix, urlValues.Encode()), http.StatusSeeOther)
		case (user == nil || user.IsAnonymous) && !strings.Contains(request.URL.Path, "openid"):
			http.Redirect(writer, request, "/api/v0/openid/begin", http.StatusFound)
		default:
			next(writer, request)
		}
	})
}

func (a *AdminServer) adminAuthMiddleware() negroni.Handler {
	return negroni.HandlerFunc(func(writer http.ResponseWriter, request *http.Request, next http.HandlerFunc) {
		user := util.GetAuthenticatedUser(request)

		// Allow the login page and POST requests to the login/logout
		// endpoints. This switch is easier than writing a long if.
		switch request.URL.Path {
		case "/login", "/admin/login", "/api/v0/openid/begin":
			if user != nil && !user.IsAnonymous && *user.Account.IsAdmin {
				// If they're already an admin, redirect to home.
				http.Redirect(writer, request, "/", http.StatusFound)
				return
			}

			fallthrough
		case "/admin/logout":
			// Unnecessary?
			next(writer, request)
		}

		if user != nil && !user.IsAnonymous && *user.Account.IsAdmin {
			next(writer, request)
		} else {
			writeJSONError(writer, errors.New(http.StatusText(http.StatusUnauthorized)), http.StatusUnauthorized)
		}
	})
}

func (a *AdminServer) adminAPIAuthMiddleware() negroni.Handler {
	return negroni.HandlerFunc(func(writer http.ResponseWriter, request *http.Request, next http.HandlerFunc) {
		user := util.GetAuthenticatedUser(request)
		if (user == nil || user.IsAnonymous) || !*user.Account.IsAdmin {
			writeJSONError(writer, errors.New(http.StatusText(http.StatusUnauthorized)), http.StatusUnauthorized)
			return
		}

		next(writer, request)
	})
}

func (a *AdminServer) noCacheMiddleware() negroni.Handler {
	return negroni.HandlerFunc(func(writer http.ResponseWriter, request *http.Request, next http.HandlerFunc) {
		writer.Header().Add("Cache-Control", "max-age=0, no-cache")
		next(writer, request)
	})
}

func (a *AdminServer) clearContext() negroni.Handler {
	return negroni.HandlerFunc(func(writer http.ResponseWriter, request *http.Request, next http.HandlerFunc) {
		defer func() {
			gorillacontext.Clear(request)
		}()

		next(writer, request)
	})
}

// upgradeMiddleware checks whether upgrades are enabled in hub config settings.  If not,
// the request flow is halted and route handlers are not called.
//
// Any upgrade handler should use this middleware.
func (a *AdminServer) upgradeMiddleware() negroni.Handler {
	return negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if settings, err := a.settingsStore.UserHubConfig(); err != nil {
			log.WithField("error", err).Error("Error loading user hub config settings")
			writeJSONError(w, errors.New(http.StatusText(http.StatusInternalServerError)), http.StatusInternalServerError)
			return
		} else if settings.DisableUpgrades {
			log.Info("Upgrades disabled; not calling route handler")
			writeJSONError(w, errors.New(http.StatusText(http.StatusForbidden)), http.StatusForbidden)
			return
		}
		next(w, r)
	})
}
