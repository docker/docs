package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/dancannon/gorethink"
	"github.com/docker/distribution/health"
	"github.com/docker/notary"
	"github.com/docker/notary/cryptoservice"
	"github.com/docker/notary/passphrase"
	"github.com/docker/notary/signer"
	"github.com/docker/notary/signer/keydbstore"
	"github.com/docker/notary/storage/rethinkdb"
	"github.com/docker/notary/trustmanager"
	"github.com/docker/notary/tuf/data"
	"github.com/docker/notary/utils"
	"github.com/spf13/viper"
)

func passphraseRetriever(keyName, alias string, createNew bool, attempts int) (passphrase string, giveup bool, err error) {
	passphrase = mainViper.GetString(strings.ToUpper(alias))

	if passphrase == "" {
		return "", false, errors.New("expected env variable to not be empty: " + alias)
	}

	return passphrase, false, nil
}

// Reads the configuration file for storage setup, and sets up the cryptoservice
// mapping
func setUpCryptoservices(configuration *viper.Viper, allowedBackends []string) (
	signer.CryptoServiceIndex, error) {
	backend := configuration.GetString("storage.backend")

	var keyStore trustmanager.KeyStore
	switch backend {
	case notary.MemoryBackend:
		keyStore = trustmanager.NewKeyMemoryStore(
			passphrase.ConstantRetriever("memory-db-ignore"))
	case notary.RethinkDBBackend:
		var sess *gorethink.Session
		storeConfig, err := utils.ParseRethinkDBStorage(configuration)
		if err != nil {
			return nil, err
		}
		defaultAlias, err := getDefaultAlias(configuration)
		if err != nil {
			return nil, err
		}
		sess, err = rethinkdb.Connection(storeConfig.CA, storeConfig.Source)
		if err != nil {
			return nil, err
		}
		s := keydbstore.NewRethinkDBKeyStore(passphraseRetriever, defaultAlias, sess)
		health.RegisterPeriodicFunc("DB operational", s.CheckHealth, time.Minute)
		keyStore = s
	case notary.MySQLBackend, notary.SQLiteBackend:
		storeConfig, err := utils.ParseSQLStorage(configuration)
		if err != nil {
			return nil, err
		}
		defaultAlias, err := getDefaultAlias(configuration)
		if err != nil {
			return nil, err
		}
		dbStore, err := keydbstore.NewKeyDBStore(
			passphraseRetriever, defaultAlias, storeConfig.Backend, storeConfig.Source)
		if err != nil {
			return nil, fmt.Errorf("failed to create a new keydbstore: %v", err)
		}

		health.RegisterPeriodicFunc(
			"DB operational", dbStore.HealthCheck, time.Minute)
		keyStore = dbStore
	}

	if doBootstrap {
		err := bootstrap(keyStore)
		if err != nil {
			logrus.Fatal(err.Error())
		}
		os.Exit(0)
	}

	cryptoService := cryptoservice.NewCryptoService(keyStore)
	cryptoServices := make(signer.CryptoServiceIndex)
	cryptoServices[data.ED25519Key] = cryptoService
	cryptoServices[data.ECDSAKey] = cryptoService
	return cryptoServices, nil
}

func getDefaultAlias(configuration *viper.Viper) (string, error) {
	defaultAlias := configuration.GetString("storage.default_alias")
	if defaultAlias == "" {
		// backwards compatibility - support this environment variable
		defaultAlias = configuration.GetString(defaultAliasEnv)
	}

	if defaultAlias == "" {
		return "", fmt.Errorf("must provide a default alias for the key DB")
	}
	logrus.Debug("Default Alias: ", defaultAlias)
	return defaultAlias, nil
}
