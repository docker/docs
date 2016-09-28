package util

import (
	"fmt"

	"github.com/docker/dhe-deploy/hubconfig"
	"github.com/docker/dhe-deploy/shared/dtrutil"

	log "github.com/Sirupsen/logrus"
	enziclient "github.com/docker/orca/enzi/api/client"
	enzierrors "github.com/docker/orca/enzi/api/errors"
	enziforms "github.com/docker/orca/enzi/api/forms"
)

type EnziConfig struct {
	Host       string
	Prefix     string
	CA         string
	VerifyCert bool
}

const (
	accountName = "docker-datacenter"
	fullName    = "Docker Datacenter"

	serviceName        = "Docker Trusted Registry"
	serviceDescription = "Docker Datacenter container storage"
)

func GetEnziConfig(haConfig *hubconfig.HAConfig) EnziConfig {
	var config EnziConfig

	if haConfig.EnziHost == "" {
		config.Host = haConfig.UCPHost
		config.Prefix = "/enzi"
		config.CA = haConfig.UCPCA
		config.VerifyCert = haConfig.UCPVerifyCert
	} else {
		config.Host = haConfig.EnziHost
		config.Prefix = ""
		config.CA = haConfig.EnziCA
		config.VerifyCert = haConfig.EnziVerifyCert
	}

	return config
}

func getFullHostAddress(userHubConfig *hubconfig.UserHubConfig) (string, error) {
	return fmt.Sprintf("https://%s", userHubConfig.DTRHost), nil
}

func RegisterAuth(enziSession *enziclient.Session, settingsStore hubconfig.SettingsStore) error {
	haConfig, err := settingsStore.HAConfig()
	if err != nil {
		return err
	}

	userHubConfig, err := settingsStore.UserHubConfig()
	if err != nil {
		return err
	}

	// NOTE (remove? This will be done by UCP when dtr is installed on top of UCP)
	createOrgForm := enziforms.CreateAccount{
		Name:     accountName,
		FullName: fullName,
		IsOrg:    true,
	}

	_, err = enziSession.CreateAccount(createOrgForm)
	if err != nil {
		apiErrs, ok := err.(*enzierrors.APIErrors)
		if !ok {
			return fmt.Errorf("unable to cast api errors when creating account: %s", err)
		} else {
			for _, apiErr := range apiErrs.Errors {
				// we get an internal error if our cert changed
				if apiErr.Code != "ACCOUNT_EXISTS" && apiErr.Code != "INTERNAL_ERROR" {
					return apiErr
				}
			}
		}
	}

	hostAddress, err := getFullHostAddress(userHubConfig)
	if err != nil {
		log.Errorln("unable to retreive DTR host address")
	}

	RedirectURIs := []string{fmt.Sprintf("%s/api/v0/openid/callback", hostAddress)}
	JWKsURIs := []string{fmt.Sprintf("%s/api/v0/openid/keys", hostAddress)}
	enziConfig := GetEnziConfig(haConfig)
	ProviderIdentities := []string{enziConfig.Host}
	URL := fmt.Sprintf("%s/", hostAddress)

	// Try to Update
	updateForm := enziforms.UpdateService{
		RedirectURIs:       &RedirectURIs,
		JWKsURIs:           &JWKsURIs,
		ProviderIdentities: &ProviderIdentities,
		CABundle:           &userHubConfig.WebTLSCA,
		URL:                &URL,
	}

	dtrService, err := enziSession.UpdateService(accountName, serviceName, updateForm)
	if err != nil {
		log.WithField("error", err).Debug("failed to update DTR service config in Auth Service")
	} else {
		if dtrService.ID == "" {
			log.Errorln("service registration update response did not have a service ID")
			return fmt.Errorf("Service registration update response did not have a service ID")
		}
		if dtrService.OwnerID == "" {
			log.Errorln("service registration update response did not have an owner ID")
			return fmt.Errorf("Service registration update response did not have an owner ID")
		}
		return settingsStore.SetEnziService(dtrService)
	}

	createServiceForm := enziforms.CreateService{
		Name:               serviceName,
		Description:        serviceDescription,
		URL:                URL,
		Privileged:         true, // VERY IMPORTANT.
		RedirectURIs:       RedirectURIs,
		JWKsURIs:           JWKsURIs,
		ProviderIdentities: ProviderIdentities,
		CABundle:           userHubConfig.WebTLSCA,
	}

	dtrService, err = enziSession.CreateService(createOrgForm.Name, createServiceForm)
	if err != nil {
		apiErrs, ok := err.(*enzierrors.APIErrors)
		if !ok {
			return fmt.Errorf("unable to cast api errors when creating account: %s", err)
		}

		for _, apiErr := range apiErrs.Errors {
			log.Errorf("%#v", apiErr)
		}

		return err
	}

	if dtrService.ID == "" {
		log.Errorln("service registration response did not have a service ID")
		return fmt.Errorf("Service registration response did not have a service ID")
	}
	if dtrService.OwnerID == "" {
		log.Errorln("service registration response did not have an owner ID")
		return fmt.Errorf("Service registration response did not have an owner ID")
	}

	return settingsStore.SetEnziService(dtrService)
}

func GetAuthAPISession(username, password string, settingsStore hubconfig.SettingsStore) (*enziclient.Session, error) {
	haConfig, err := settingsStore.HAConfig()
	if err != nil {
		return nil, err
	}

	enziConfig := GetEnziConfig(haConfig)
	client, err := dtrutil.HTTPClient(!enziConfig.VerifyCert, enziConfig.CA)
	if err != nil {
		return nil, err
	}

	basicAuthenticator := &enziclient.BasicAuthenticator{
		Username: username,
		Password: password,
	}

	return enziclient.New(client, enziConfig.Host, enziConfig.Prefix, basicAuthenticator), nil
}
