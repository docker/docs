package legacy

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/docker/distribution/context"
	"github.com/docker/libkv"
	kvstore "github.com/docker/libkv/store"
	"github.com/docker/libkv/store/etcd"
	orcaauth "github.com/docker/orca/auth"
	"github.com/docker/orca/auth/builtin"
	enziadaptor "github.com/docker/orca/auth/enzi"
	"github.com/docker/orca/enzi/api/forms"
	ldapconfig "github.com/docker/orca/enzi/authn/ldap/config"
	enziconfig "github.com/docker/orca/enzi/config"
	"github.com/docker/orca/enzi/schema"
	enziutil "github.com/docker/orca/enzi/util"
	"github.com/docker/orca/utils"
)

const (
	defaultOrgName           = "docker-datacenter"
	defaultOrgFullName       = "Docker Datacenter"
	ucpServiceName           = "Docker Universal Control Plane"
	ucpServiceDescription    = "Docker Datacenter Container Orchestration"
	ucpAuthConfigKVKey       = "orca/v1/config/auth2"
	ucpLegacyAuthConfigKVKey = "orca/v1/config/auth"
)

func init() {
	etcd.Register()
}

// MigrateLegacyAuthData migrates user, team, membership, and auth config from
// storage in the UCP KV Store to storage in the eNZi auth service. The UCP KV
// store must be running at the given kvURL. The eNZi auth database must be
// running at the given authDBAddr. Both stores are connected to over TLS using
// the given tlsConfig. The given ucpServiceAddr must be set to the host:port
// of the UCP controller. The given authAPIAddr must be set to the host:port of
// the eNZi auth API server.
//
// First, a default org and service representing UCP is registered with the
// auth provider if they do not already exist. Data is only migrated if either
// forceMigrate is true or if we detect that a new UCP auth config has not yet
// been set in its KV store. Once the migration is complete with no errors, or
// if the migration was skipped, a new auth config is set in the UCP KV store
// so that a new version of the UCP controller may be started.
func MigrateLegacyAuthData(kvURL *url.URL, tlsConfig *tls.Config, ucpServiceAddr, serviceCAFilename, authAPIAddr, authDBAddr string, forceMigrate bool) error {
	log := context.GetLogger(context.Background())

	dbSession, err := enziutil.GetDBSession([]string{authDBAddr}, tlsConfig)
	if err != nil {
		return fmt.Errorf("unable to get Auth DB session: %s", err)
	}
	defer dbSession.Close()

	// Wait for all table replicas to be ready or timeout after 1 minute.
	waitForDB := func() error {
		return schema.WaitForReadyTables(dbSession, schema.AllReplicasReady)
	}

	if err := utils.RunWithTimeout(waitForDB, time.Minute); err != nil {
		return fmt.Errorf("unable to wait for auth store replicas: %s", err)
	}

	schemaMgr := schema.NewRethinkDBManager(dbSession)

	kvStore, err := getKVHandle(log, kvURL, tlsConfig)
	if err != nil {
		return fmt.Errorf("unable to connect to KV store: %s", err)
	}
	defer kvStore.Close()

	// This legacy authenticator cannot be used to verify session tokens
	// and can be used safely with a nil value given as auth context.
	legacyAuthenticator := builtin.NewWithStore(kvStore)

	// Create the default Docker Datacenter org if it does not exist.
	defaultOrg, err := getOrCreateDefaultOrg(log, schemaMgr)
	if err != nil {
		return err
	}

	// Create UCP service registration if it does not exist.
	ucpService, err := getOrCreateUCPService(log, schemaMgr, defaultOrg, authAPIAddr, ucpServiceAddr, serviceCAFilename)
	if err != nil {
		return err
	}

	// Ensure that our provider address in in the list of provider
	// identities for the service. It should be if we just registered it.
	updateServiceWithNewAuthProvider(schemaMgr, ucpService, authAPIAddr)

	// Try to get the new auth config (if it exists).
	var ucpAuthConfig *orcaauth.AuthenticatorConfiguration
	kvPair, err := kvStore.Get(ucpAuthConfigKVKey)
	if err == nil {
		// The new auth config already exists.
		ucpAuthConfig = new(orcaauth.AuthenticatorConfiguration)
		if err := json.Unmarshal(kvPair.Value, ucpAuthConfig); err != nil {
			return fmt.Errorf("unable to decode current UCP auth config: %s", err)
		}
	} else if err != kvstore.ErrKeyNotFound {
		err = utils.MaybeWrapEtcdClusterErr(err)
		return fmt.Errorf("unable to get UCP auth config from KV store: %s", err)
	}

	// We consider the data to be already migrated *iff* the new auth
	// config has been set in the KV store and we have successfully
	// configured at least one auth provider instance (which could only
	// have been done by a fresh install/join of a newer version or by this
	// migratino having been completed in the past).
	alreadyMigrated := ucpAuthConfig != nil && len(ucpAuthConfig.EnziConfig.ProviderAddrs) > 0

	if !alreadyMigrated || forceMigrate {
		// We have either detected that a migration has has not yet
		// ran and succeeded or we are forcing another migration.
		ucpAuthConfig, err = performMigration(log, kvStore, legacyAuthenticator, schemaMgr, defaultOrg)
		if err != nil {
			return fmt.Errorf("unable to perform migration: %s", err)
		}
	}

	// Update UCP Auth Config KV entry.
	updatedEnziConfig := orcaauth.EnziConfig{
		DefaultOrgID:  defaultOrg.ID,
		ServiceID:     ucpService.ID,
		ProviderAddrs: ucpService.ProviderIdentities,
		// Preserve existing user default role.
		UserDefaultRole: ucpAuthConfig.EnziConfig.UserDefaultRole,
	}

	ucpAuthConfig.AuthenticatorType = orcaauth.AuthenticatorEnzi
	ucpAuthConfig.EnziConfig = updatedEnziConfig

	ucpAuthConfigJSON, err := json.Marshal(ucpAuthConfig)
	if err != nil {
		return fmt.Errorf("unable to encode updated UCP auth config to JSON: %s", err)
	}

	if err := kvStore.Put(ucpAuthConfigKVKey, ucpAuthConfigJSON, nil); err != nil {
		err = utils.MaybeWrapEtcdClusterErr(err)
		return fmt.Errorf("unable to update UCP auth config in KV store: %s", err)
	}

	log.Info("Successfully Updated UCP Auth Config")

	// Run this function to fix a bug with the original 1.1.0 migrator.
	// See Github issue https://github.com/docker/orca/issues/1550
	if err := ensureOrgTeamMembersAreOrgMembers(schemaMgr, defaultOrg); err != nil {
		return fmt.Errorf("unable to ensure team members are org members: %s", err)
	}

	return nil
}

func getKVHandle(log context.Logger, kvURL *url.URL, tlsConfig *tls.Config) (kvstore.Store, error) {
	log.Debugf("connecting to legacy kv store: %s", kvURL)

	kvOpts := &kvstore.Config{
		ConnectionTimeout: time.Second * 10,
		TLS:               tlsConfig,
	}

	// Note: this call doesn't immediately contact the KV store.
	kvStore, err := libkv.NewStore(kvstore.ETCD, []string{kvURL.Host}, kvOpts)
	if err != nil {
		return nil, fmt.Errorf("unable to get kvStore handle: %s", err)
	}

	// Retry things during startup. There is an Orca Instance Key stored in
	// the KV Store. If we can successfully get this value then we may
	// assume that the kv store is ready to use.
	retryTimeout := 5 * time.Minute
	retryWaitDuration := 2 * time.Second
	startTime := time.Now()

	// Check for the existence of the legacy auth config key. It might not
	// exist if this system never ran with legacy auth configuration. We
	// don't really care though, we just want to make sure that the KV
	// store is ready.
	_, err = kvStore.Exists(ucpLegacyAuthConfigKVKey)
	for err != nil && time.Since(startTime) < retryTimeout {
		log.Debugf("unable to access kv store: %s", err)
		log.Debugf("KV store not yet ready. Trying again in %s ...", retryWaitDuration)

		time.Sleep(retryWaitDuration)
		_, err = kvStore.Exists(ucpLegacyAuthConfigKVKey)
	}

	if err != nil {
		return nil, fmt.Errorf("unable to wait for KV store to be ready: %s", err)
	}

	return kvStore, nil
}

func getOrCreateDefaultOrg(log context.Logger, schemaMgr schema.Manager) (*schema.Account, error) {
	defaultOrg, err := schemaMgr.GetOrgByName(defaultOrgName)
	if err == nil {
		log.Debugf("%s Organization Already Exists", defaultOrgName)

		return defaultOrg, nil
	}

	if err != schema.ErrNoSuchAccount {
		return nil, fmt.Errorf("unable to get %q org: %s", defaultOrgName, err)
	}

	// The default org does not yet exist. We need to create it.
	defaultOrg = &schema.Account{
		Name:     defaultOrgName,
		FullName: defaultOrgFullName,
		IsOrg:    true,
	}

	// This will have set the ID for the org.
	if err := schemaMgr.CreateAccount(defaultOrg); err != nil {
		return nil, fmt.Errorf("unable to create %q org: %s", defaultOrgName, err)
	}

	log.Debugf("Successfully Created %s Organization", defaultOrgName)

	return defaultOrg, nil
}

func getOrCreateUCPService(log context.Logger, schemaMgr schema.Manager, defaultOrg *schema.Account, providerAddr, serviceAddr, caBundleFilename string) (*schema.Service, error) {
	ucpService, err := schemaMgr.GetServiceByName(defaultOrg.ID, ucpServiceName)
	if err == nil {
		log.Debugf("%s Service Already Exists", ucpServiceName)

		return ucpService, nil
	}

	if err != schema.ErrNoSuchService {
		return nil, fmt.Errorf("unable to get %q service: %s", ucpServiceName, err)
	}

	// The UCP service does not yet exist. We need to create it.
	CAPEMBundle, err := ioutil.ReadFile(caBundleFilename)
	if err != nil {
		return nil, fmt.Errorf("unable to read Root CA Certificate file: %s", err)
	}

	ucpService = &schema.Service{
		OwnerID:            defaultOrg.ID,
		Name:               ucpServiceName,
		Description:        ucpServiceDescription,
		URL:                fmt.Sprintf("https://%s", serviceAddr),
		Privileged:         true,
		RedirectURIs:       []string{fmt.Sprintf("https://%s/openid_callback", serviceAddr)},
		JWKsURIs:           []string{fmt.Sprintf("https://%s/openid_keys", serviceAddr)},
		GrantTypes:         []string{"authorization_code", "refresh_token", "service_session", "password", "root_session"},
		ResponseTypes:      []string{"code"},
		ProviderIdentities: []string{providerAddr},
		CABundle:           string(CAPEMBundle),
	}

	// This will have set the ID for the service.
	if err := schemaMgr.CreateService(ucpService); err != nil {
		return nil, fmt.Errorf("unable to create %q service: %s", ucpServiceName, err)
	}

	log.Debugf("Successfully Created %s Service", ucpServiceName)

	return ucpService, nil
}

func updateServiceWithNewAuthProvider(schemaMgr schema.Manager, ucpService *schema.Service, providerAddr string) error {
	providerIDs := make(map[string]struct{}, len(ucpService.ProviderIdentities))
	for _, providerID := range ucpService.ProviderIdentities {
		providerIDs[providerID] = struct{}{}
	}

	if _, exists := providerIDs[providerAddr]; exists {
		return nil // No need to update.
	}

	// Need to update the service in the same way we would during
	// bootstrap/join.
	ucpService.ProviderIdentities = append(ucpService.ProviderIdentities, providerAddr)
	updateFields := schema.ServiceUpdateFields{
		ProviderIdentities: &ucpService.ProviderIdentities,
	}

	if err := schemaMgr.UpdateService(ucpService.ID, updateFields); err != nil {
		return fmt.Errorf("unable to update %q service: %s", ucpServiceName, err)
	}

	return nil
}

func performMigration(log context.Logger, kvStore kvstore.Store, legacyAuthenticator orcaauth.Authenticator, schemaMgr schema.Manager, defaultOrg *schema.Account) (ucpAuthConfig *orcaauth.AuthenticatorConfiguration, err error) {
	// Start by migrating users.
	usernameIDs, userMigrationErrors, err := migrateUserData(log, kvStore, legacyAuthenticator, schemaMgr)
	if err != nil {
		return nil, fmt.Errorf("unable to migrate user data: %s", err)
	}

	// Next, migrate teams and members.
	teamMigrationErrors, err := migrateTeamData(log, kvStore, legacyAuthenticator, schemaMgr, defaultOrg, usernameIDs)
	if err != nil {
		return nil, fmt.Errorf("unable to migrate team data: %s", err)
	}

	// Now is the time to check for errors.
	numErrs := 0
	if len(userMigrationErrors) > 0 {
		numErrs += len(userMigrationErrors)
		for username, err := range userMigrationErrors {
			log.Warnf("error migrating user %s: %s", username, err)
		}
	}
	if len(teamMigrationErrors) > 0 {
		numErrs += len(teamMigrationErrors)
		for teamName, err := range teamMigrationErrors {
			log.Warnf("error migrating team %s: %s", teamName, err)
		}
	}
	if numErrs > 0 {
		return nil, fmt.Errorf("encountered %d errors during data migration", numErrs)
	}

	// Finally, migrate auth config.
	ucpAuthConfig, err = migrateAuthConfig(kvStore, schemaMgr)
	if err != nil {
		return nil, fmt.Errorf("unable to migrate legacy auth config: %s", err)
	}

	log.Info("Successfully Completed Data Migration")

	return ucpAuthConfig, nil
}

func migrateUserData(log context.Logger, kvStore kvstore.Store, legacyAuthenticator orcaauth.Authenticator, schemaMgr schema.Manager) (usernameIDs map[string]string, userMigrationErrors map[string]error, err error) {
	// Migrate users.
	legacyUsers, err := legacyAuthenticator.ListUsers(nil)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to list users from legacy auth store: %s", err)
	}

	log.Infof("Migrating %d users from legacy storage", len(legacyUsers))

	legacyAdaptor := enziadaptor.NewLegacyAdaptor(kvStore)

	usernameIDs = make(map[string]string, len(legacyUsers))
	userMigrationErrors = map[string]error{}

	progress := &enziutil.ProgressLogger{
		Logger:   log,
		Message:  "Migrating User Accounts",
		StepSize: 1000,
		Length:   len(legacyUsers),
	}

	for progress.Next() {
		legacyUser := legacyUsers[progress.Index()]

		// Normalize the username.
		username := forms.NormalizeAccountName(legacyUser.Username)
		// Legacy accounts have separate first and last names.
		fullName := strings.TrimSpace(fmt.Sprintf("%s %s", legacyUser.FirstName, legacyUser.LastName))
		// Legacy accounts use only bycrypt password hashing while new
		// accounts can use either bycrypt or pbkdf2_sha256.
		var passwordHash string
		if legacyUser.Password != "" {
			passwordHash = fmt.Sprintf("bcrypt$%s", legacyUser.Password)
		}

		newUser := &schema.Account{
			Name:         username,
			FullName:     fullName,
			IsOrg:        false,
			IsAdmin:      legacyUser.Admin,
			IsActive:     !legacyUser.Disabled,
			PasswordHash: passwordHash,
			LdapDN:       legacyUser.LdapDN,
		}

		// Creating the account will also generate and set an ID.
		err := schemaMgr.CreateAccount(newUser)
		if err == schema.ErrAccountExists {
			newUser, err = schemaMgr.GetUserByName(username)
		}
		if err != nil {
			// Save the error and continue to the next user.
			userMigrationErrors[username] = fmt.Errorf("unable to migrate user: %s", err)
			continue
		}

		usernameIDs[username] = newUser.ID

		// Set these values for later.
		legacyUser.Username = username
		legacyUser.ID = newUser.ID

		// Migrate the user's UCP role.
		if err := legacyAdaptor.SetUserRole(username, legacyUser.Role); err != nil {
			// Save the error and continue to the next user.
			userMigrationErrors[username] = fmt.Errorf("unable to set user role: %s", err)
			continue
		}

		// Migrate the user's Public Keys.
		for _, accountKey := range legacyUser.PublicKeys {
			pubKeyBlock, _ := pem.Decode([]byte(accountKey.PublicKey))
			if pubKeyBlock == nil {
				// Ignore invalid keys.
				log.Debugf("found invalid key for user %s", username)
				continue
			}

			pubKey, err := x509.ParsePKIXPublicKey(pubKeyBlock.Bytes)
			if err != nil {
				// Ignore invalid keys.
				log.Debugf("found invalid key for user %s", username)
				continue
			}

			if err := legacyAdaptor.AddUserPublicKey(legacyUser, accountKey.Label, pubKey); err != nil {
				// Unable to add to KV store?
				userMigrationErrors[username] = fmt.Errorf("unable to add user public key: %s", err)
				// This just skips to the next key. If there's
				// another error then only the last error is
				// saved. This is Okay.
				continue
			}
		}
	}

	return usernameIDs, userMigrationErrors, nil
}

func migrateTeamData(log context.Logger, kvStore kvstore.Store, legacyAuthenticator orcaauth.Authenticator, schemaMgr schema.Manager, defaultOrg *schema.Account, usernameIDs map[string]string) (teamMigrationErrors map[string]error, err error) {
	// Migrate Teams and Members.
	legacyTeams, err := legacyAuthenticator.ListTeams(nil)
	if err != nil {
		return nil, fmt.Errorf("unable to list teams from legacy auth store: %s", err)
	}

	log.Infof("Migrating %d teams from legacy storage", len(legacyTeams))

	teamMigrationErrors = map[string]error{}

	progress := &enziutil.ProgressLogger{
		Logger:   log,
		Message:  "Migrating Teams",
		StepSize: 1000,
		Length:   len(legacyTeams),
	}

	for progress.Next() {
		legacyTeam := legacyTeams[progress.Index()]

		// Normalize the team name.
		teamName := forms.NormalizeFullName(legacyTeam.Name)

		newTeam := &schema.Team{
			OrgID: defaultOrg.ID,
			Name:  teamName,
			// Keep the legacy ID so that UCP team access lists
			// don't break. Legacy Team IDs are the first 8 bytes
			// of a hash of the time when the team was created.
			ID: legacyTeam.Id,
		}

		if legacyTeam.LdapDN != "" {
			newTeam.MemberSyncConfig = schema.MemberSyncOpts{
				EnableSync:         true,
				SelectGroupMembers: true,
				GroupDN:            legacyTeam.LdapDN,
				GroupMemberAttr:    legacyTeam.LdapMemberAttr,
			}
		}

		// Creating the team will also generate and set the PK.
		err := schemaMgr.CreateTeam(newTeam)
		if err == schema.ErrTeamExists {
			newTeam, err = schemaMgr.GetTeamByName(defaultOrg.ID, teamName)
		}
		if err != nil {
			// Save the error and continue to the next team.
			teamMigrationErrors[teamName] = fmt.Errorf("unable to migrate team: %s", err)
			continue
		}

		// Now migrate the team's current membership. UCP stores
		// managed and ldap-synced members separately.
		managedTeamMembersKVDir := path.Join("orca/v1/teammembers", legacyTeam.Id)
		managedKVPairs, err := kvStore.List(managedTeamMembersKVDir)
		if err != nil && err != kvstore.ErrKeyNotFound {
			teamMigrationErrors[teamName] = fmt.Errorf("unable to list managed team members from legacy store: %s", err)
			continue
		}
		ldapSyncedTeamMembersKVDir := path.Join("orca/v1/ldapteammembers", legacyTeam.Id)
		ldapSyncedKVPairs, err := kvStore.List(ldapSyncedTeamMembersKVDir)
		if err != nil && err != kvstore.ErrKeyNotFound {
			teamMigrationErrors[teamName] = fmt.Errorf("unable to list ldap-synced team members from legacy store: %s", err)
			continue
		}

		kvPairs := append(managedKVPairs, ldapSyncedKVPairs...)

		// The values are usernames of team members.
		for _, kvPair := range kvPairs {
			username := forms.NormalizeAccountName(string(kvPair.Value))
			userID, ok := usernameIDs[username]
			if !ok {
				// This user must have been deleted but there's
				// still a record of their team membership?
				log.Debugf("found team %s membership for non-existant user: %s", teamName, username)
				continue
			}

			// AddTeamMembership also adds an org membership entry.
			if err := schemaMgr.AddTeamMembership(defaultOrg.ID, legacyTeam.Id, userID, nil, nil); err != nil {
				teamMigrationErrors[teamName] = fmt.Errorf("unable to add user %s to team: %s", username, err)
				// This just skips to the next key. If there's
				// another error then only the last error is
				// saved. This is Okay.
				continue
			}
		}
	}

	return teamMigrationErrors, nil
}

// migrateAuthConfig migrates the legacy auth config from the UCP kv store to
// enzi. A new UCP auth config is prepared and returned, but not yet set. On
// success, the caller should finalize the new config in the KV store with the
// service's ID and default OrgID.
func migrateAuthConfig(kvStore kvstore.Store, schemaMgr schema.Manager) (ucpAuthConfig *orcaauth.AuthenticatorConfiguration, err error) {
	// Get the legacy auth config, it might not exist if this system never
	// ran UCP that used the legacy auth config.
	ucpAuthConfig = new(orcaauth.AuthenticatorConfiguration)
	kvPair, err := kvStore.Get(ucpLegacyAuthConfigKVKey)
	if err == nil {
		// Legacy auth config exists.
		if err := json.Unmarshal(kvPair.Value, ucpAuthConfig); err != nil {
			return nil, fmt.Errorf("unable to decode UCP legacy auth config: %s", err)
		}
	} else if err != kvstore.ErrKeyNotFound {
		err = utils.MaybeWrapEtcdClusterErr(err)
		return nil, fmt.Errorf("unable to get UCP legacy auth config from KV store: %s", err)
	}

	legacyLDAPConfig := ucpAuthConfig.LDAPConfig

	newLDAPConfig := &ldapconfig.Settings{
		RecoveryAdminUsername: legacyLDAPConfig.AdminUsername,
		ServerURL:             legacyLDAPConfig.ServerURL,
		StartTLS:              legacyLDAPConfig.StartTLS,
		RootCerts:             legacyLDAPConfig.RootCerts,
		TLSSkipVerify:         legacyLDAPConfig.TLSSkipVerify,
		ReaderDN:              legacyLDAPConfig.ReaderDN,
		ReaderPassword:        legacyLDAPConfig.ReaderPassword,
		UserSearchConfigs: []ldapconfig.UserSearchOpts{{
			BaseDN:       legacyLDAPConfig.UserBaseDN,
			ScopeSubtree: true, // Legacy config always did this.
			UsernameAttr: legacyLDAPConfig.UserLoginAttrName,
			Filter:       legacyLDAPConfig.UserSearchFilter,
		}},
		// Disabled by default, but we are explicit here.
		AdminSyncOpts: ldapconfig.MemberSyncOpts{
			EnableSync: false,
		},
		// The legacy config syncs at regular intervals in units of
		// minutes. It's pretty easy to convert that to a CronSpec
		// which is used by the new system.
		SyncSchedule: enziutil.MinutesToCron(legacyLDAPConfig.SyncInterval),
	}

	// The user default role has been moved to a new location, separate
	// from LDAP sync configuration.
	ucpAuthConfig.EnziConfig.UserDefaultRole = legacyLDAPConfig.UserDefaultRole

	if err := ldapconfig.SetLDAPConfig(schemaMgr, newLDAPConfig); err != nil {
		return nil, fmt.Errorf("unable to set new LDAP config: %s", err)
	}

	newAuthConfig := &enziconfig.Auth{
		Backend: enziconfig.AuthBackendManaged,
	}

	if ucpAuthConfig.AuthenticatorType == orcaauth.AuthenticatorLDAP {
		newAuthConfig.Backend = enziconfig.AuthBackendLDAP
	}

	if err := enziconfig.SetAuthConfig(schemaMgr, newAuthConfig); err != nil {
		return nil, fmt.Errorf("unable to set new auth config: %s", err)
	}

	return ucpAuthConfig, nil
}

// ensureTeamMembersAreOrgMembers loops through all of the distinct team
// members in the org and ensures that all members of teams in that org are
// also members of the organization. This is done as a consistency repair after
// a bug in early legacy migration code. After a long enough while, it can
// probably be safely removed, when we think there will be minimal impact, i.e.
// barely anyone is running pre-v1.1.1.
func ensureOrgTeamMembersAreOrgMembers(schemaMgr schema.Manager, org *schema.Account) error {
	userIDs, err := schemaMgr.ListTeamMembersInOrg(org.ID)
	if err != nil {
		return fmt.Errorf("unable to list distinct team members for org %s: %s", org.Name, err)
	}

	for _, userID := range userIDs {
		if err := schemaMgr.AddOrgMembership(org.ID, userID, nil, nil); err != nil {
			return fmt.Errorf("unable to add org membership: %s", err)
		}
	}

	return nil
}
