package adminserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/dhe-deploy"
	"github.com/docker/dhe-deploy/hubconfig"
	"github.com/docker/dhe-deploy/hubconfig/defaultconfigs"
	"github.com/docker/dhe-deploy/hubconfig/util"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/dhe-license-server/tiers"
	"github.com/satori/go.uuid"
	"gopkg.in/yaml.v2"

	// register all storage and auth drivers
	"github.com/docker/distribution/configuration"
	"github.com/docker/distribution/context"
	"github.com/docker/distribution/registry"
	_ "github.com/docker/distribution/registry/auth/htpasswd"
	_ "github.com/docker/distribution/registry/auth/silly"
	_ "github.com/docker/distribution/registry/auth/token"
	_ "github.com/docker/distribution/registry/proxy"
	_ "github.com/docker/distribution/registry/storage/driver/azure"
	_ "github.com/docker/distribution/registry/storage/driver/filesystem"
	_ "github.com/docker/distribution/registry/storage/driver/gcs"
	_ "github.com/docker/distribution/registry/storage/driver/inmemory"
	_ "github.com/docker/distribution/registry/storage/driver/middleware/cloudfront"
	_ "github.com/docker/distribution/registry/storage/driver/oss"
	_ "github.com/docker/distribution/registry/storage/driver/s3-aws"
	_ "github.com/docker/distribution/registry/storage/driver/swift"
)

func (a *AdminServer) getRegistrySettingsHandler(writer http.ResponseWriter, request *http.Request) {
	settings, err := a.getRegistrySettings()
	if err != nil {
		writeJSONError(writer, err, http.StatusInternalServerError)
		return
	}

	writeJSON(writer, settings)
}

func (a *AdminServer) getRegistryConfigurationsHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, storageConfigurations)
}

func (a *AdminServer) getLicenseSettingsHandler(writer http.ResponseWriter, request *http.Request) {
	settings, err := a.getLicenseSettings()
	if err != nil {
		writeJSONError(writer, err, http.StatusInternalServerError)
		return
	}

	writeJSON(writer, settings)
}

func (a *AdminServer) getRegistrySettings() (*RegistrySettings, error) {
	registryConfig, err := a.settingsStore.RegistryConfig()
	if err != nil {
		log.WithField("error", err).Error("Failed to retrieve registry config")
		return nil, err
	}

	registryConfigBytes, err := yaml.Marshal(registryConfig)
	if err != nil {
		log.WithField("error", err).Error("Failed to marshal registry config")
		return nil, err
	}

	registryConfigString := string(registryConfigBytes)
	if registryConfig == nil {
		registryConfig = &configuration.Configuration{}
	}

	return &RegistrySettings{
		Config:  &registryConfigString,
		Storage: registryConfig.Storage,
	}, nil
}

func (a *AdminServer) getLicenseSettings() (*LicenseSettings, error) {
	if a.licenseChecker.LicenseType() == tiers.Hourly {
		return &LicenseSettings{
			IsValid:     a.licenseChecker.IsValid(),
			KeyID:       "",
			LicenseType: a.licenseChecker.LicenseType(),
			LicenseTier: a.licenseChecker.LicenseTier(),
		}, nil
	} else {
		licenseConfig, err := a.settingsStore.LicenseConfig()
		if err != nil {
			log.WithField("error", err).Error("Failed to retrieve license config")
			return nil, err
		}

		// Sanitize the license info!
		var sanitizedLicenseConfig *LicenseSettings
		if licenseConfig != nil {
			sanitizedLicenseConfig = &LicenseSettings{
				IsValid:     a.licenseChecker.IsValid(),
				KeyID:       licenseConfig.KeyID,
				AutoRefresh: licenseConfig.AutoRefresh,
				Expiration:  a.licenseChecker.Expiration(),
				LicenseType: a.licenseChecker.LicenseType(),
				LicenseTier: a.licenseChecker.LicenseTier(),
			}
		}

		return sanitizedLicenseConfig, nil
	}
}

func (a *AdminServer) TryRegistryConfig(storageConfig *configuration.Configuration) (err error) {
	// fix parsing bugs in json vs yaml
	util.SetReadonlyMode(&storageConfig.Storage, util.GetReadonlyMode(&storageConfig.Storage))

	cert, err := a.settingsStore.GarantRootCert()
	if err != nil {
		log.WithField("error", err).Error("Failed to write get garant root cert")
		return fmt.Errorf("Failed to write get garant root cert")
	}
	err = os.MkdirAll(deploy.ConfigDirPath, 0777)
	if err != nil {
		log.WithField("error", err).Error("Failed to create dir for garant root cert")
		return fmt.Errorf("Failed to create dir for garant root cert")
	}
	err = ioutil.WriteFile(filepath.Join(deploy.ConfigDirPath, deploy.GarantRootCertFilename), []byte(cert), 0666)
	if err != nil {
		log.WithField("error", err).Error("Failed to write write out garant root cert")
		return fmt.Errorf("Failed to write write out garant root cert")
	}

	defer func() {
		if recoverErr := recover(); recoverErr != nil {
			// this sets the outside error so the function can return an error
			err = fmt.Errorf("invalid storage settings: %s", recoverErr)
			log.WithField("error", err).Error("failed to run registry with new settings")
			return
		}
	}()

	if _, err = registry.NewRegistry(context.Background(), storageConfig); err != nil {
		return err
	}

	// test out the new storage option
	driver, err := util.DriverFromStorage(storageConfig.Storage)
	if err != nil {
		log.Errorf("Error creating driver: %s", err.Error())
		return err
	}

	filename := "/" + uuid.NewV4().String()

	// test write
	err = driver.PutContent(context.Background(), filename, []byte("testfile"))
	if err != nil {
		return err
	}

	// make sure file got written
	_, err = driver.Stat(context.Background(), filename)
	if err != nil {
		return err
	}

	// finally, make sure the file can be deleted
	return driver.Delete(context.Background(), filename)
}

func (a *AdminServer) updateRegistrySettingsHandler(writer http.ResponseWriter, request *http.Request) {
	var settings RegistrySettings
	if err := json.NewDecoder(request.Body).Decode(&settings); err != nil {
		log.WithField("error", err).Warn("Failed to decode new registry settings")
		writeJSONError(writer, err, http.StatusBadRequest)
		return
	}
	defer request.Body.Close()

	storageConfig, err := configuration.Parse(strings.NewReader(*settings.Config))
	if err != nil {
		// TODO: log a sanitized yml? maybe this is impossible because we failed to parse
		log.WithField("error", err).Info("Invalid storage config file")
		writeJSONError(writer, errors.New("invalid storage configuration"), http.StatusBadRequest)
		return
	}

	err = a.TryRegistryConfig(storageConfig)
	if err != nil {
		log.WithField("error", err).Error("Failed to update registry config")
		writeJSONError(writer, err, http.StatusBadRequest)
		return
	}

	log.Info("Updating registry config")

	if err := a.settingsStore.SetRegistryConfig(storageConfig); err != nil {
		log.WithField("error", err).Error("Failed to update registry config")
		writeJSONError(writer, errors.New(http.StatusText(http.StatusInternalServerError)), http.StatusInternalServerError)
		return
	}

	writeJSONStatus(writer, nil, http.StatusAccepted)
}

func (a *AdminServer) updateRegistrySettingsViaFormHandler(w http.ResponseWriter, r *http.Request) {
	config, err := a.settingsStore.RegistryConfig()
	if err != nil {
		log.WithField("error", err).Warn("Failed to load registry storage settings")
		writeJSONError(w, err, http.StatusInternalServerError)
		return
	}

	// If somehow no registry configuration exists settingsStore.RegistryConfig() returns
	// nil, nil, so start with our default configuration.
	if config == nil {
		config = new(configuration.Configuration)
		*config = defaultconfigs.DefaultRegistryConfig
	}

	// This endpoint overwrites the existing storage configuration.  Without this, if a user
	// selects "S3" when "filesystem" is already configured there will be multiple drivers selected.
	config.Storage = configuration.Storage{}

	if err := json.NewDecoder(r.Body).Decode(&config.Storage); err != nil {
		log.WithField("error", err).Warn("Failed to decode JSON input for registry settings via form")
		writeJSONError(w, err, http.StatusInternalServerError)
		return
	}

	err = a.TryRegistryConfig(config)
	if err != nil {
		log.WithField("error", err).Error("Failed to update registry config")
		writeJSONError(w, err, http.StatusBadRequest)
		return
	}

	log.Info("Updating registry config")

	if err := a.settingsStore.SetRegistryConfig(config); err != nil {
		log.WithField("error", err).Error("Failed to update registry config")
		writeJSONError(w, err, http.StatusInternalServerError)
		return
	}

	// fix parsing bugs in json vs yaml
	util.SetReadonlyModeJSON(&config.Storage, util.GetReadonlyMode(&config.Storage))
	util.RemoveUseragent(&config.Storage)

	writeJSONStatus(w, config.Storage, http.StatusAccepted)
}

func (a *AdminServer) updateLicenseSettingsHandler(writer http.ResponseWriter, request *http.Request) {
	var settings *hubconfig.LicenseConfig
	if err := json.NewDecoder(request.Body).Decode(&settings); err != nil {
		log.WithField("error", err).Warn("Failed to decode new license settings")
		writeJSONError(writer, err, http.StatusBadRequest)
		return
	}
	defer request.Body.Close()

	log.Info("Updating license key")

	err := a.licenseChecker.LoadLicenseFromConfig(settings, true)
	if err != nil {
		log.Error("Failed to update license - invalid license")
		writeJSONError(writer, errors.New("failed to update license - invalid license"), http.StatusBadRequest)
		return
	}

	newLicense, err := a.getLicenseSettings()
	if err != nil {
		log.WithField("error", err).Error("failed to check license")
		writeJSONError(writer, errors.New("failed to check license"), http.StatusInternalServerError)
		return
	}

	writeJSONStatus(writer, newLicense, http.StatusAccepted)
}

func (a *AdminServer) toggleLicenseAutoRefreshHandler(writer http.ResponseWriter, request *http.Request) {
	var jsonBody json.RawMessage
	if err := json.NewDecoder(request.Body).Decode(&jsonBody); err != nil {
		log.WithField("error", err).Warn("Failed to decode license auto refresh toggle request")
		writeJSONError(writer, err, http.StatusBadRequest)
		return
	}

	toggleRequest := struct {
		Setting bool `json:"auto_refresh"`
	}{}

	if err := json.Unmarshal(jsonBody, &toggleRequest); err != nil {
		log.WithField("error", err).Warn("Failed to parse json")
		writeJSONError(writer, err, http.StatusInternalServerError)
		return
	}

	err := a.licenseChecker.ToggleAutoRefresh(toggleRequest.Setting)
	if err != nil {
		log.WithField("error", err).Error("failed to toggle license auto refresh")
		writeJSONError(writer, errors.New("failed to toggle license auto refresh"), http.StatusInternalServerError)
		return
	}

	writeJSONStatus(writer, nil, http.StatusAccepted)
}
