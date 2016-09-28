package garant

import (
	"net/http"

	"github.com/bugsnag/bugsnag-go"
	"github.com/yvasiyarov/gorelic"
)

// configureReporting wraps the given HTTP handler with Bugsnag and Newrelic
// reporting handlers if they have been added to the app's configuration.
func (app *App) configureReporting(handler http.Handler) http.Handler {
	var (
		bugsnagOpts  = app.reportingOpts.Bugsnag
		newrelicOpts = app.reportingOpts.NewRelic
	)

	if bugsnagOpts.APIKey != "" {
		handler = bugsnag.Handler(handler, bugsnag.Configuration{
			APIKey:              bugsnagOpts.APIKey,
			ReleaseStage:        bugsnagOpts.ReleaseStage,
			Endpoint:            bugsnagOpts.Endpoint,
			NotifyReleaseStages: bugsnagOpts.NotifyReleaseStages,
			AppVersion:          bugsnagOpts.AppVersion,
			ProjectPackages:     bugsnagOpts.ProjectPackages,
		})
	}

	if newrelicOpts.LicenseKey != "" {
		agent := gorelic.NewAgent()
		agent.NewrelicLicense = newrelicOpts.LicenseKey
		agent.NewrelicName = newrelicOpts.Name
		agent.CollectHTTPStat = true
		agent.Run()

		handler = agent.WrapHTTPHandler(handler)
	}

	return suppressPanic(handler)
}

// suppressPanic wraps the given handler and suppresses any runtime panics that
// occur in the wrapped handler. The garant app instance will have already
// logged a stacktrace with context information and written a simple error
// response to the client, so don't do any logging, just supress the panic.
func suppressPanic(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			recover() // Explicitly ignore the panic value.
		}()

		handler.ServeHTTP(w, r)
	})
}
