package manager

import (
	"crypto"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/orca"
	"github.com/docker/orca/auth"
	"github.com/docker/orca/auth/builtin"
	"github.com/docker/orca/auth/dtr"
	"github.com/docker/orca/auth/enzi"
	"github.com/docker/orca/auth/ldap"
)

var (
	singletonAuthKvKey = "auth2"
	MinimumSyncTime    = 60   // Minutes
	DefaultSyncTime    = 1440 // 1 day
)

var (
	ErrInvalidPublicKey = errors.New("invalid public key")
)

type AuthConfigSubsystem struct {
	ksKey string
	m     *DefaultManager
	cfg   *auth.AuthenticatorConfiguration
	// We need an extra layer of indirection since Authenticator is an interface
	a *struct{ authenticator auth.Authenticator }
}

func NewAuthConfigSubsystem(key, jsonConfig string, m *DefaultManager) (ConfigSubsystem, error) {
	// Auth doesn't support instances, just a single one
	if key != singletonAuthKvKey {
		log.Debugf("Malformed auth config key: %s", key)
		return nil, fmt.Errorf("Only one auth configuration supported")
	}

	authType := auth.AuthenticatorBuiltin
	var cfg auth.AuthenticatorConfiguration
	if jsonConfig == "" {
		cfg = auth.AuthenticatorConfiguration{
			AuthenticatorType: authType,
			// Set up some sensible defaults
			LDAPConfig: auth.LDAPSettings{
				StartTLS:          true,
				SyncInterval:      DefaultSyncTime,
				UserLoginAttrName: "sAMAccountName",
			},
		}
	}
	s := AuthConfigSubsystem{
		m:     m,
		ksKey: path.Join(KsConfigDir, singletonAuthKvKey),
		cfg:   &cfg,
		a:     &struct{ authenticator auth.Authenticator }{},
	}
	if jsonConfig != "" {
		cfgInt, err := s.ValidateConfig(jsonConfig, false)
		if err != nil {
			return nil, err
		}
		cfg, ok := cfgInt.(auth.AuthenticatorConfiguration)
		if !ok {
			// Should not happen
			return nil, fmt.Errorf("Malformed configuration type")
		}
		s.UpdateConfig(cfg)
	}
	m.configSubsystems[filepath.Base(s.ksKey)] = s
	return s, nil
}

func setupAuthenticator(m *DefaultManager) {
	// Use an empty string as the default jsonConfig. The contructor will
	// use a sane default in this case.
	setupSingletonConfigSubsystem(m, singletonAuthKvKey, "", NewAuthConfigSubsystem)

	// Spawn a goroutine to keep our enzi signing key from expiring in the
	// KV store.
	go m.persistEnziSigningKey()
}

func (s AuthConfigSubsystem) GetKvKey() string {
	return s.ksKey
}

func (s AuthConfigSubsystem) ValidateConfig(jsonConfig string, userInitiated bool) (interface{}, error) {
	var cfg auth.AuthenticatorConfiguration
	if jsonConfig == "" {
		cfg = auth.AuthenticatorConfiguration{
			AuthenticatorType: auth.AuthenticatorBuiltin,
		}
	} else {
		if err := json.Unmarshal([]byte(jsonConfig), &cfg); err != nil {
			return nil, fmt.Errorf("Malformed authenticator configuration: %s", err)
		}
	}
	if cfg.AuthenticatorType == auth.AuthenticatorLDAP {
		// Perform some very basic validation of the configuration - deeper validation comes during sync
		if userInitiated {
			if cfg.LDAPConfig.AdminPassword == "" {
				return nil, fmt.Errorf("You must specify an admin username and password to verify the configuration")
			}
		}
		if cfg.LDAPConfig.AdminUsername == "" {
			return nil, fmt.Errorf("You must specify an admin username and password to verify the configuration")
		}
		if cfg.LDAPConfig.ServerURL == "" {
			return nil, fmt.Errorf("You must specify an LDAP server URL (e.g., ldap://server.acme.com)")
		} // TODO - might consider parsing it as a URL to catch garbage, but the sync will fail anyway...
		if cfg.LDAPConfig.ReaderDN == "" {
			return nil, fmt.Errorf("You must specify a LDAP account used for performing searches")
		} // Theoretically the search account might be unauthenticated... so allow blank passwords
		if cfg.LDAPConfig.UserBaseDN == "" {
			return nil, fmt.Errorf("You must specify a base DN for users (e.g., DC=example,DC=com)")
		}
		if cfg.LDAPConfig.UserLoginAttrName == "" {
			return nil, fmt.Errorf("You must specify an LDAP attribute used for account names (e.g., sAMAccountName)")
		}
		if cfg.LDAPConfig.UserSearchFilter == "" {
			return nil, fmt.Errorf("You must specify an LDAP search to select the user accounts - e.g., (memberOf=OU=Groups,DC=example,DC=com)")
		}
		// Force a minimum sync time.
		if cfg.LDAPConfig.SyncInterval < MinimumSyncTime {
			cfg.LDAPConfig.SyncInterval = MinimumSyncTime
		}
		// LDAP config looks plausible
	}
	return cfg, nil
}
func (s AuthConfigSubsystem) UpdateConfig(cfgInt interface{}) (err error) {
	priorAuthenticator := s.cfg.AuthenticatorType
	cfg, ok := cfgInt.(auth.AuthenticatorConfiguration)
	if !ok {
		return fmt.Errorf("Incorrect configuration type: %t", cfgInt)
	}

	// Make sure we don't consider the LDAP admin password in comparisons
	// or write it out to any back-end configuration
	ldapAdminPassword := cfg.LDAPConfig.AdminPassword
	cfg.LDAPConfig.AdminPassword = ""

	if cfg.Equals(s.cfg) && s.a.authenticator != nil {
		log.Debug("Auth config unchanged")
		return nil
	}

	kv := s.m.Datastore()

	var authenticator auth.Authenticator
	switch cfg.AuthenticatorType {
	case auth.AuthenticatorEnzi:
		log.Info("Using eNZi authenticator")
		authenticator, err = enzi.NewAuthenticator(cfg.EnziConfig, s.m.hostAddr, kv, s.m.enziTokenSigningKey, s.m.clusterCAChain, s.m.enziTokenCertChain...)
		if err != nil {
			log.Errorf("unable to create eNZi authenticator backend: %s", err)
			return err
		}
	case auth.AuthenticatorDTR:
		log.Info("Using DTR authenticator")
		authenticator = dtr.NewAuthenticator(
			&kv,
			s.m.trustKey.PublicKey().KeyID(),
			s.m.controllerCertPEM,
			s.m.controllerKeyPEM,
			cfg.DTRConfig.URL,
			cfg.DTRConfig.Insecure)
		log.Infof("Setting admin user to '%s'", cfg.DTRConfig.AdminUser)
		s.m.createAdmin(cfg.DTRConfig.AdminUser, false)
	case auth.AuthenticatorLDAP:

		log.Info("Using LDAP authenticator")
		ldapAuth := ldap.NewAuthenticator(
			&kv,
			s.m.trustKey.PublicKey().KeyID(),
			s.m.controllerCertPEM,
			s.m.controllerKeyPEM,
			&cfg.LDAPConfig)
		if ldapAdminPassword != "" {
			// For user initiated, force the sync, and only do admin for quicker initial sync to verify settings
			err := ldapAuth.Sync(nil, true, true)
			if err != nil {
				return fmt.Errorf("Failed to perform initial sync - please verify your configuration: %s", err)
			}

			// Now make sure the admin password works...
			log.Infof("Verifying login works with LDAP account %s", cfg.LDAPConfig.AdminUsername)
			_, err = ldapAuth.AuthenticateUsernamePassword(cfg.LDAPConfig.AdminUsername, ldapAdminPassword)
			if err != nil {
				return fmt.Errorf("Unable to authenticate with provided admin credentials: %s", err)
			}
			// If we got this far, it looks good, so start a real full sync in the background
			go ldapAuth.Sync(nil, true, false)
		} else {
			// For non-user initiated events, don't force the sync (another HA node probably did it)
			err := ldapAuth.Sync(nil, false, false)
			if err != nil {
				return fmt.Errorf("Failed to perform initial sync - please verify your configuration: %s", err)
			}
		}

		authenticator = ldapAuth
	default:
		cfg.AuthenticatorType = auth.AuthenticatorBuiltin // force it in case we get bad data
		log.Info("Using builtin authenticator")
		authenticator = builtin.NewAuthenticator(
			&kv,
			s.m.trustKey.PublicKey().KeyID(),
			s.m.controllerCertPEM,
			s.m.controllerKeyPEM)
	}

	s.a.authenticator = authenticator
	*s.cfg = cfg

	// If we just switched from LDAP to builtin, make sure we have an admin account
	if priorAuthenticator == auth.AuthenticatorLDAP && cfg.AuthenticatorType == auth.AuthenticatorBuiltin {
		account, err := s.a.authenticator.GetUser(nil, "admin")
		if err != nil {
			s.m.createAdmin("admin", true)
		} else if account.Disabled {
			log.Info(`Re-enabling the "admin" account`)
			account.Disabled = false
			account.Password = "" // Don't clobber the password
			s.a.authenticator.SaveUser(nil, account)
		}
	}
	return nil
}

func (s AuthConfigSubsystem) GetConfiguration() (string, error) {
	data, err := json.Marshal(s.cfg)
	return string(data), err
}

func (m DefaultManager) GetAuthenticator() auth.Authenticator {
	s := m.configSubsystems[singletonAuthKvKey].(AuthConfigSubsystem)
	return s.a.authenticator
}

func (m DefaultManager) AuthenticateUsernamePassword(username, password, remoteAddr string) (*auth.Context, error) {
	authenticator := m.GetAuthenticator()
	ctx, err := authenticator.AuthenticateUsernamePassword(username, password)
	if err != nil {
		m.SaveEvent(&orca.Event{
			Type:       "auth fail",
			Time:       time.Now(),
			Username:   username,
			RemoteAddr: remoteAddr,
			Message:    authenticator.Name() + ":Password based auth failure",
		})
	} else {
		m.SaveEvent(&orca.Event{
			Type:       "auth ok",
			Time:       time.Now(),
			Username:   username,
			RemoteAddr: remoteAddr,
			Message:    authenticator.Name() + ":Password based auth suceeded",
		})
	}
	return ctx, err
}

// WARNING: the returned auth.Account object does not contain team membership information
func (m DefaultManager) AuthenticatePublicKey(pubKey crypto.PublicKey, remoteAddr string) (*auth.Context, error) {
	ctx, err := m.GetAuthenticator().AuthenticatePublicKey(pubKey)
	if err == auth.ErrInvalidPublicKey {
		m.SaveEvent(&orca.Event{
			Type:       "auth fail",
			Time:       time.Now(),
			RemoteAddr: remoteAddr,
			Message:    "Public key based auth failure",
		})

		return nil, ErrInvalidPublicKey
	}
	if err == auth.ErrAccountDisabled {
		m.SaveEvent(&orca.Event{
			Type:       "auth fail",
			Time:       time.Now(),
			RemoteAddr: remoteAddr,
			Message:    "Public key based auth failure by disabled account",
		})

		return nil, ErrInvalidPublicKey
	}

	return ctx, err
}

func (m DefaultManager) periodicAuthSync() {
	if err := m.GetAuthenticator().Sync(nil, false, false); err != nil {
		log.Error(err)
	}
}
func (m DefaultManager) AuthSyncMessages(ctx *auth.Context) string {
	return m.GetAuthenticator().LastSyncStatus(ctx)
}
func (m DefaultManager) AuthSync(ctx *auth.Context) string {
	go m.GetAuthenticator().Sync(ctx, true, false)
	return "Sync started in the background"
}

// persistEnziSigningKey should run in a goroutine. It periodically wakes up to
// extend the TTL of this server's signing key for eNZi service tokens.
func (m DefaultManager) persistEnziSigningKey() {
	ticks := time.Tick(time.Hour)
	for range ticks {
		authenticator := m.GetAuthenticator()
		enziAuthenticator, ok := authenticator.(*enzi.Authenticator)
		if !ok {
			// We're not currently using the eNZi authenticator.
			// Just wait for the next clock tick.
			continue
		}

		err := enziAuthenticator.SaveSigningKey()
		for err != nil {
			log.Errorf("unable to extend signing key expiration (trying again in 30s): %s", err)

			time.Sleep(time.Second * 30)
			err = enziAuthenticator.SaveSigningKey()
		}

		log.Debugf("successfully extended expiration time for signing key %s", m.enziTokenSigningKey.ID)
	}
}
