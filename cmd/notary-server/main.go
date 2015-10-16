package main

import (
	_ "expvar"
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/bugsnag/bugsnag-go"
	"github.com/docker/distribution/health"
	_ "github.com/docker/distribution/registry/auth/htpasswd"
	_ "github.com/docker/distribution/registry/auth/token"
	"github.com/endophage/gotuf/signed"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/net/context"

	bugsnag_hook "github.com/Sirupsen/logrus/hooks/bugsnag"
	"github.com/docker/notary/server"
	"github.com/docker/notary/server/storage"
	"github.com/docker/notary/signer"
	"github.com/docker/notary/version"
	"github.com/spf13/viper"
)

// DebugAddress is the debug server address to listen on
const DebugAddress = "localhost:8080"

var debug bool
var configFile string

func init() {
	// set default log level to Error
	viper.SetDefault("logging", map[string]interface{}{"level": 2})

	// Setup flags
	flag.StringVar(&configFile, "config", "", "Path to configuration file")
	flag.BoolVar(&debug, "debug", false, "Enable the debugging server on localhost:8080")
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if debug {
		go debugServer(DebugAddress)
	}

	// when the server starts print the version for debugging and issue logs later
	logrus.Infof("Version: %s, Git commit: %s", version.NotaryVersion, version.GitCommit)

	ctx := context.Background()

	filename := filepath.Base(configFile)
	ext := filepath.Ext(configFile)
	configPath := filepath.Dir(configFile)

	viper.SetConfigType(strings.TrimPrefix(ext, "."))
	viper.SetConfigName(strings.TrimSuffix(filename, ext))
	viper.AddConfigPath(configPath)

	// Automatically accept configuration options from the environment
	viper.SetEnvPrefix("NOTARY_SERVER")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		logrus.Error("Viper Error: ", err.Error())
		logrus.Error("Could not read config at ", configFile)
		os.Exit(1)
	}
	lvl, err := logrus.ParseLevel(viper.GetString("logging.level"))
	if err != nil {
		lvl = logrus.ErrorLevel
		logrus.Error("Could not parse log level from config. Defaulting to ErrorLevel")
	}
	logrus.SetLevel(lvl)

	// set up bugsnag and attach to logrus
	bugs := viper.GetString("reporting.bugsnag")
	if bugs != "" {
		apiKey := viper.GetString("reporting.bugsnag_api_key")
		releaseStage := viper.GetString("reporting.bugsnag_release_stage")
		bugsnag.Configure(bugsnag.Configuration{
			APIKey:       apiKey,
			ReleaseStage: releaseStage,
		})
		hook, err := bugsnag_hook.NewBugsnagHook()
		if err != nil {
			logrus.Error("Could not attach bugsnag to logrus: ", err.Error())
		} else {
			logrus.AddHook(hook)
		}
	}
	keyAlgo := viper.GetString("trust_service.key_algorithm")
	if keyAlgo == "" {
		logrus.Fatal("no key algorithm configured.")
		os.Exit(1)
	}
	ctx = context.WithValue(ctx, "keyAlgorithm", keyAlgo)

	var trust signed.CryptoService
	if viper.GetString("trust_service.type") == "remote" {
		logrus.Info("Using remote signing service")
		trust = signer.NewNotarySigner(
			viper.GetString("trust_service.hostname"),
			viper.GetString("trust_service.port"),
			viper.GetString("trust_service.tls_ca_file"),
		)
		minute := 1 * time.Minute
		health.RegisterPeriodicFunc(
			"Trust operational",
			// If the trust service fails, the server is degraded but not
			// exactly unheatlthy, so always return healthy and just log an
			// error.
			func() error {
				err := trust.(*signer.NotarySigner).CheckHealth(minute)
				if err != nil {
					logrus.Error("Trust not fully operational: ", err.Error())
				}
				return nil
			},
			minute)
	} else {
		logrus.Info("Using local signing service")
		trust = signed.NewEd25519()
	}

	if viper.GetString("storage.backend") == "mysql" {
		logrus.Info("Using mysql backend")
		dbURL := viper.GetString("storage.db_url")
		store, err := storage.NewSQLStorage("mysql", dbURL)
		if err != nil {
			logrus.Fatal("Error starting DB driver: ", err.Error())
			return // not strictly needed but let's be explicit
		}
		health.RegisterPeriodicFunc(
			"DB operational", store.CheckHealth, time.Second*60)
		ctx = context.WithValue(ctx, "metaStore", store)
	} else {
		logrus.Debug("Using memory backend")
		ctx = context.WithValue(ctx, "metaStore", storage.NewMemStorage())
	}
	logrus.Info("Starting Server")
	err = server.Run(
		ctx,
		viper.GetString("server.addr"),
		viper.GetString("server.tls_cert_file"),
		viper.GetString("server.tls_key_file"),
		trust,
		viper.GetString("auth.type"),
		viper.Get("auth.options"),
	)

	logrus.Error(err.Error())
	return
}

func usage() {
	fmt.Println("usage:", os.Args[0])
	flag.PrintDefaults()
}

// debugServer starts the debug server with pprof, expvar among other
// endpoints. The addr should not be exposed externally. For most of these to
// work, tls cannot be enabled on the endpoint, so it is generally separate.
func debugServer(addr string) {
	logrus.Info("Debug server listening on", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		logrus.Fatal("error listening on debug interface: ", err)
	}
}
