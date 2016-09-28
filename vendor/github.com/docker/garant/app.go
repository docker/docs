package garant

import (
	"fmt"
	"net/http"
	"os"

	"github.com/docker/distribution/context"
	"github.com/docker/garant/auth"
	"github.com/docker/garant/config"
	"github.com/docker/libtrust"
	"github.com/gorilla/mux"
)

// App is a Garant Token Server Application.
type App struct {
	config        *config.Configuration
	authorizer    auth.Authorizer
	signingKey    libtrust.PrivateKey
	reportingOpts config.Reporting
	baseContext   context.Context
}

// NewApp configures a new Garant Token Server Application.
func NewApp(configurationPath string) (*App, error) {
	fmt.Println("initializing token signing app")
	appConfig, err := loadConfiguration(configurationPath)
	if err != nil {
		return nil, fmt.Errorf("unable to load configuration: %s", err)
	}

	authorizer, err := auth.NewAuthorizer(appConfig.Auth.BackendName, appConfig.Auth.Parameters)
	if err != nil {
		return nil, err
	}

	signingKey, err := libtrust.LoadKeyFile(appConfig.SigningKey)
	if err != nil {
		return nil, err
	}

	ctx := context.WithValues(
		configureLoggingContext(appConfig.Logging),
		map[string]interface{}{
			"authBackend":  appConfig.Auth.BackendName,
			"signingKeyID": signingKey.KeyID(),
		},
	)

	logger := context.GetLogger(ctx, "authBackend", "signingKeyID")

	logger.Info("token signing app initialized")

	app := &App{
		config:        appConfig,
		authorizer:    authorizer,
		signingKey:    signingKey,
		reportingOpts: appConfig.Reporting,
		baseContext:   context.WithLogger(ctx, logger),
	}

	return app, nil
}

// loadConfiguration reads and parses a garant token server configuration file
// located at the given file path.
func loadConfiguration(configurationPath string) (*config.Configuration, error) {
	configFile, err := os.Open(configurationPath)
	if err != nil {
		return nil, err
	}
	defer configFile.Close()

	appConfig, err := config.Parse(configFile)
	if err != nil {
		return nil, fmt.Errorf("error parsing %s: %v", configurationPath, err)
	}

	return appConfig, nil
}

// ListenAndServe listens on the configured HTTP server to handle token and
// simple account info requests.
func (app *App) ListenAndServe() error {
	router := mux.NewRouter()

	app.registerHandlers(app.config.HTTP.Prefix, router)

	handler := app.configureReporting(router)

	addr := app.config.HTTP.Addr
	if addr == "" {
		// Use the default address.
		addr = "localhost:8080"
	}

	if app.config.HTTP.TLS != nil {
		certFile := app.config.HTTP.TLS.Certificate
		keyFile := app.config.HTTP.TLS.Key

		return http.ListenAndServeTLS(addr, certFile, keyFile, handler)
	}

	return http.ListenAndServe(addr, handler)
}
