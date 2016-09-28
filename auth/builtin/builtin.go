package builtin

import (
	"crypto"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"path"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	jwt "github.com/dgrijalva/jwt-go"
	kvstore "github.com/docker/libkv/store"
	"github.com/docker/orca"
	"github.com/docker/orca/auth"
	"github.com/docker/orca/utils"
	"golang.org/x/crypto/bcrypt"
)

var (
	TokenExp                  = time.Hour * 72
	datastoreVersion          = "orca/v1"
	datastoreCerts            = datastoreVersion + "/controller_pub_certs"
	datastoreTokens           = datastoreVersion + "/token_ids/active"
	datastoreBlacklist        = datastoreVersion + "/token_ids/expired"
	datastoreAccounts         = datastoreVersion + "/accounts"
	datastoreAccountTeams     = datastoreVersion + "/accountteams"
	datastoreTeams            = datastoreVersion + "/teams"
	datastoreTeamMembers      = datastoreVersion + "/teammembers"
	datastoreLastSyncMessages = datastoreVersion + "/authlastsyncmessages"
)

type (
	BuiltinAuthenticator struct {
		orcaID    string
		tlsCert   []byte
		tlsKey    []byte
		store     kvstore.Store
		tokenType string
		certStore map[string][]byte
	}
)

var _ auth.Authenticator = (*BuiltinAuthenticator)(nil)

// NewWithStore returns a new builtin authenticator which is incapable of
// generating or verifying session tokens. It is exported for use with other
// packages that do not need that functionality.
func NewWithStore(store kvstore.Store) *BuiltinAuthenticator {
	return &BuiltinAuthenticator{
		store: store,
	}
}

func NewAuthenticator(store *kvstore.Store, orcaID string, tlsCert, tlsKey []byte) *BuiltinAuthenticator {
	var tokenType string
	var cert *x509.Certificate
	var err error

	if len(tlsCert) > 0 && len(tlsKey) > 0 {
		log.Debug("Using RSA style tokens")
		tokenType = "rsa"

		cert, err = parseCert(tlsCert)
		if err != nil {
			log.Errorf("Couldn't parse cert: %s", err)
		}

		// use the serial number of the public cert so that we can use it to verify incoming JWTs
		orcaID = cert.SerialNumber.String()
	} else {
		log.Info("Couldn't find RSA certs, using HMAC style tokens")
		tokenType = "hmac"
		secret, err := auth.GenerateRandID()
		if err != nil {
			log.Errorf("Couldn't generate a random token; unable to set up authenticator: %s", err)
			return nil
		}
		tlsCert = []byte(secret)
		tlsKey = []byte(secret)
	}

	// Store our external public cert in the kv store so that other controllers can
	// use it based on the token issuer.  The certStore will work as a cache so we
	// only have to do a lookup once.
	kv := *store
	p := path.Join(datastoreCerts, orcaID)
	if err := kv.Put(p, tlsCert, nil); err != nil {
		err = utils.MaybeWrapEtcdClusterErr(err)
		log.Error(err)
		return nil
	}

	a := &BuiltinAuthenticator{
		orcaID:    orcaID,
		tlsCert:   tlsCert,
		tlsKey:    tlsKey,
		store:     *store,
		tokenType: tokenType,
		certStore: map[string][]byte{orcaID: tlsCert},
	}

	a.updateOldAccounts()
	return a
}

func (a *BuiltinAuthenticator) Store() kvstore.Store {
	return a.store
}

func parseCert(pemCert []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(pemCert)
	if block == nil {
		return nil, fmt.Errorf("Invalid cert")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	return cert, nil
}

func pubKeyPem(pubKey crypto.PublicKey) (pemOut string, err error) {
	derBytes, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return "", err
	}

	pemBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derBytes,
	})

	return strings.TrimSpace(string(pemBytes)), nil
}

func (a *BuiltinAuthenticator) AddUserPublicKey(user *auth.Account, label string, publicKey crypto.PublicKey) error {
	keyPem, err := pubKeyPem(publicKey)
	if err != nil {
		return err
	}

	user.PublicKeys = append(user.PublicKeys, auth.AccountKey{
		Label:     label,
		PublicKey: keyPem,
	})

	// reset the password field so it is not updated
	user.Password = ""
	if _, err := a.SaveUser(nil, user); err != nil {
		return err
	}

	return nil
}

func (a *BuiltinAuthenticator) AuthenticatePublicKey(pubKey crypto.PublicKey) (*auth.Context, error) {
	keyPem, err := pubKeyPem(pubKey)
	if err != nil {
		return nil, err
	}

	// TODO: this is not very efficient.  we may have to
	// change the data model to improve
	accounts, err := a.ListUsers(nil)
	if err != nil {
		log.Warnf("Failed to lookup accounts: %s", err)
		return nil, err
	}

	for _, account := range accounts {
		for _, publicKey := range account.PublicKeys {
			if publicKey.PublicKey != keyPem {
				continue // Not a match.
			}

			// Found an account with the same public key.
			if account.Disabled {
				return nil, auth.ErrAccountDisabled
			}

			// Note: We make a call to GetUser() here so
			// that the user's teams field is populated.
			// It's okay to pass a nil context since it is
			// ignored.
			user, err := a.GetUser(nil, account.Username)
			if err != nil {
				return nil, err
			}

			return &auth.Context{
				User: user,
			}, nil
		}
	}

	return nil, auth.ErrInvalidPublicKey
}

func (a *BuiltinAuthenticator) CanChangePassword(ctx *auth.Context) bool {
	return true
}

func (a *BuiltinAuthenticator) Name() string {
	return "builtin"
}

func (a *BuiltinAuthenticator) AuthenticateUsernamePassword(username, password string) (*auth.Context, error) {
	acct, err := a.GetUser(nil, username)
	if err != nil {
		return nil, err
	}

	if acct.Disabled {
		return nil, auth.ErrAccountDisabled
	}

	passHash := acct.Password

	if err := bcrypt.CompareHashAndPassword([]byte(passHash), []byte(password)); err != nil {
		return nil, auth.ErrInvalidPassword
	}

	tokenStr, err := a.GenerateToken(username, nil)
	if err != nil {
		return nil, err
	}

	authCtx := &auth.Context{
		User:         acct,
		SessionToken: tokenStr,
	}

	return authCtx, nil
}

func (a *BuiltinAuthenticator) parseToken(tokenStr string) (*jwt.Token, error) {
	// Tokens don't include the public cert, so we look up the cert of the issuer
	// and then return that to jwt.Parse() which will validate the signature.
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if a.tokenType == "rsa" {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
		} else if a.tokenType == "hmac" {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
		}
		issuer, _ := token.Claims["iss"].(string)
		return a.lookupCert(issuer)
	})

	if err != nil {
		return nil, err
	}
	return token, err
}

func (a *BuiltinAuthenticator) lookupCert(certID string) ([]byte, error) {
	var cert []byte

	// Attempt to get the cert from our cache.  If it's not there, look it up in the kvstore.
	cert, ok := a.certStore[certID]
	if !ok {
		p := path.Join(datastoreCerts, certID)
		kv, err := a.store.Get(p)
		if err != nil {
			return nil, utils.MaybeWrapEtcdClusterErr(err)
		}
		cert = kv.Value
		a.certStore[certID] = cert
	}
	return cert, nil
}

func (a *BuiltinAuthenticator) AuthenticateSessionToken(tokenStr string) (*auth.Context, error) {
	token, err := a.parseToken(tokenStr)
	if err != nil {
		// Invalid Token.
		log.Errorf("unable to parse token as JWT: %s", err)
		return nil, auth.ErrInvalidSessionToken
	}

	// TODO: in case the KV store is busted, we can still allow a person to
	//       be authenticated here, since all the info we need is passed in on
	//       the claims of the jwt

	// check to see if the token was in our blacklist
	blackListPath := path.Join(datastoreBlacklist, token.Claims["jti"].(string))
	isBlacklisted, err := a.store.Exists(blackListPath)
	if err != nil {
		return nil, err
	}
	if isBlacklisted {
		log.Errorf("JWT %s was blacklisted", token.Claims["jti"].(string))
		return nil, auth.ErrInvalidSessionToken
	}

	sub, ok := token.Claims["sub"]
	if !ok {
		log.Error("JWT has no 'sub' field")
		return nil, auth.ErrInvalidSessionToken
	}
	username := sub.(string)

	// Get the account object
	acct, err := a.GetUser(nil, username)
	if err != nil {
		return nil, err
	}

	if acct.Disabled {
		return nil, auth.ErrAccountDisabled
	}

	authCtx := &auth.Context{
		User:         acct,
		SessionToken: tokenStr,
	}

	return authCtx, nil
}

func (a *BuiltinAuthenticator) Logout(ctx *auth.Context) error {
	if ctx.SessionToken == "" {
		// Nothing to do.
		return nil
	}

	token, err := a.parseToken(ctx.SessionToken)
	if err != nil {
		return err
	}

	return a.expireToken(token)
}

func (a *BuiltinAuthenticator) GenerateToken(username string, extraClaims map[string]string) (string, error) {
	var token *jwt.Token

	if username == "" {
		return "", fmt.Errorf("No username for generating a token")
	}

	ID, err := auth.GenerateRandID()
	if err != nil {
		return "", err
	}

	exp := time.Now().Add(TokenExp)
	if a.tokenType == "rsa" {
		token = jwt.New(jwt.SigningMethodRS256)
	} else {
		token = jwt.New(jwt.SigningMethodHS256)
	}
	token.Claims["exp"] = exp.Unix()
	token.Claims["sub"] = username
	token.Claims["iss"] = a.orcaID
	token.Claims["jti"] = ID

	if extraClaims != nil {
		for k, v := range extraClaims {
			token.Claims[k] = v
		}
	}

	tokenStr, err := token.SignedString(a.tlsKey)
	if err != nil {
		return "", err
	}

	k := path.Join(datastoreTokens, username, ID)
	authToken := &auth.AuthToken{
		ID: ID,
	}

	data, err := json.Marshal(authToken)
	if err != nil {
		return "", err
	}

	if err := a.store.Put(k, data, &kvstore.WriteOptions{TTL: TokenExp}); err != nil {
		err = utils.MaybeWrapEtcdClusterErr(err)
		log.Errorf("Couldn't store token: %s", err)
		return "", err
	}

	return tokenStr, nil
}

// TODO: there needs to be a way of expiring all tokens for a user when they do things like
//       reset their password.

func (a *BuiltinAuthenticator) expireToken(token *jwt.Token) error {
	username, ok := token.Claims["sub"].(string)
	if !ok {
		return fmt.Errorf("No subject found in token")
	}
	ID, ok := token.Claims["jti"].(string)
	if !ok {
		return fmt.Errorf("No id found in token")
	}

	k := path.Join(datastoreTokens, username, ID)
	if err := a.store.Delete(k); err != nil {
		log.Errorf("Couldn't delete active token: %s", err)
		return err
	}

	// store the token in the blacklist
	p := path.Join(datastoreBlacklist, ID)
	authToken := &auth.AuthToken{
		ID: ID,
	}

	data, err := json.Marshal(authToken)
	if err != nil {
		return err
	}
	if err = a.store.Put(p, data, &kvstore.WriteOptions{TTL: TokenExp}); err != nil {
		err = utils.MaybeWrapEtcdClusterErr(err)
		log.Errorf("Couldn't expire token: %s", err)
		return err
	}

	return nil
}

func (a *BuiltinAuthenticator) expireTokensForUser(username string) error {
	k := path.Join(datastoreTokens, username)

	kvPairs, err := a.store.List(k)
	if err != nil && err != kvstore.ErrKeyNotFound {
		log.Debugf("couldn't get tokens: %s", err)
		return err
	}

	for _, kvPair := range kvPairs {
		var t *auth.AuthToken
		if err := json.Unmarshal(kvPair.Value, &t); err != nil {
			log.Errorf("couldn't unmarshal token: %s", err)
			return err
		}
		// insert the token ID into the blacklist
		p := path.Join(datastoreBlacklist, t.ID)
		if err := a.store.Put(p, kvPair.Value, &kvstore.WriteOptions{TTL: TokenExp}); err != nil {
			return utils.MaybeWrapEtcdClusterErr(err)
		}
	}

	// delete old user tokens
	if err := a.store.Delete(k); err != nil && err != kvstore.ErrKeyNotFound {
		return err
	}
	return nil
}

func (a *BuiltinAuthenticator) GetUser(ctx *auth.Context, username string) (*auth.Account, error) {
	k := path.Join(datastoreAccounts, username)
	exists, err := a.store.Exists(k)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, auth.ErrAccountDoesNotExist
	}

	kvPair, err := a.store.Get(k)
	if err != nil {
		return nil, utils.MaybeWrapEtcdClusterErr(err)
	}

	data := kvPair.Value

	var account *auth.Account

	if err := json.Unmarshal(data, &account); err != nil {
		return nil, err
	}

	teams, err := a.getAccountTeams(account.Username)
	if err != nil {
		return nil, err
	}

	account.ManagedTeams = teams

	return account, nil
}

func (a *BuiltinAuthenticator) SaveUser(ctx *auth.Context, account *auth.Account) (string, error) {
	var (
		hash      string
		eventType string
	)

	if account.Password != "" {
		h, err := auth.Hash(account.Password)
		if err != nil {
			return "", err
		}

		hash = h
	}

	// check if exists; if so, update
	acct, err := a.GetUser(ctx, account.Username)
	if err != nil && err != auth.ErrAccountDoesNotExist {
		return "", err
	}

	// update
	if acct != nil {
		// if new password specified; use the hash
		if account.Password != "" {
			account.Password = hash
		} else { // keep the same hashed password
			account.Password = acct.Password
		}

		eventType = "update-account"
	} else {
		account.Password = hash
		eventType = "add-account"
	}

	// Make sure there isn't any membership info written out
	account.ManagedTeams = []string{}

	// kv
	data, err := json.Marshal(account)
	if err != nil {
		return "", err
	}

	k := path.Join(datastoreAccounts, account.Username)
	if err := a.store.Put(k, data, nil); err != nil {
		return "", utils.MaybeWrapEtcdClusterErr(err)
	}

	// If the account is disabled, make sure all tokens have been invalidated for it
	if account.Disabled {
		a.expireTokensForUser(account.Username)
	}

	return eventType, nil
}

// Used only when you know the account is new and pre-validated for more efficient bulk insertions
func (a *BuiltinAuthenticator) UnvalidatedAddAccount(account *auth.Account) error {
	var (
		hash string
	)

	if account.Password != "" {
		h, err := auth.Hash(account.Password)
		if err != nil {
			return err
		}

		hash = h
	}

	account.Password = hash

	// Make sure there isn't any membership info written out
	account.ManagedTeams = []string{}

	// kv
	data, err := json.Marshal(account)
	if err != nil {
		return err
	}

	k := path.Join(datastoreAccounts, account.Username)
	if err := a.store.Put(k, data, nil); err != nil {
		return utils.MaybeWrapEtcdClusterErr(err)
	}

	return nil
}

func (a *BuiltinAuthenticator) ListUsers(ctx *auth.Context) ([]*auth.Account, error) {
	kvPairs, err := a.store.List(datastoreAccounts)
	if err != nil && err == kvstore.ErrKeyNotFound {
		if err = a.store.Put(datastoreAccounts, nil, &kvstore.WriteOptions{IsDir: true}); err != nil {
			return nil, utils.MaybeWrapEtcdClusterErr(err)
		}
		// Corner case - make sure the directory exists
	} else if err != nil {
		return nil, err
	}

	accounts := []*auth.Account{}

	for _, kvPair := range kvPairs {
		var account *auth.Account
		if err := json.Unmarshal(kvPair.Value, &account); err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (a *BuiltinAuthenticator) ListTeamMembers(ctx *auth.Context, teamID string) ([]*auth.Account, error) {
	members, err := a.getTeamMembers(teamID)
	if err != nil {
		return nil, err
	}

	accounts := []*auth.Account{}

	for _, m := range members {
		acct, err := a.GetUser(ctx, m)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, acct)
	}

	return accounts, nil
}

func (a *BuiltinAuthenticator) DeleteUser(ctx *auth.Context, account *auth.Account) error {
	// invalidate any tokens
	if err := a.expireTokensForUser(account.Username); err != nil {
		return err
	}

	// delete the account
	k := path.Join(datastoreAccounts, account.Username)
	if err := a.store.Delete(k); err != nil {
		return err
	}

	return nil
}

// TODO: make certain all of the JWTs get invalidated
func (a *BuiltinAuthenticator) ChangePassword(ctx *auth.Context, username, oldPassword, newPassword string) error {
	acct, err := a.GetUser(ctx, username)
	if err != nil {
		return err
	}

	// Verify old password.
	if _, err := a.AuthenticateUsernamePassword(username, oldPassword); err != nil {
		return err
	}

	acct.Password = newPassword
	if _, err := a.SaveUser(ctx, acct); err != nil {
		return err
	}

	// user will have to sign in again
	if err = a.expireTokensForUser(username); err != nil {
		return err
	}

	return nil
}

func (a *BuiltinAuthenticator) GetTeam(ctx *auth.Context, id string) (*auth.Team, error) {
	if id == "" {
		return nil, auth.ErrTeamDoesNotExist
	}
	k := path.Join(datastoreTeams, id)
	exists, err := a.store.Exists(k)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, auth.ErrTeamDoesNotExist
	}

	kvPair, err := a.store.Get(k)
	if err != nil {
		return nil, utils.MaybeWrapEtcdClusterErr(err)
	}

	data := kvPair.Value

	var team *auth.Team

	if err := json.Unmarshal(data, &team); err != nil {
		return nil, err
	}

	members, err := a.getTeamMembers(team.Id)
	if err != nil {
		return nil, err
	}

	team.ManagedMembers = members

	return team, nil
}

func (a *BuiltinAuthenticator) ListTeams(ctx *auth.Context) ([]*auth.Team, error) {
	teams := []*auth.Team{}

	kvPairs, err := a.store.List(datastoreTeams)
	if err != nil && err == kvstore.ErrKeyNotFound {
		return teams, nil
	} else if err != nil {
		return nil, err
	}

	for _, kvPair := range kvPairs {
		var team *auth.Team
		if err := json.Unmarshal(kvPair.Value, &team); err != nil {
			return nil, err
		}

		teams = append(teams, team)
	}

	return teams, nil
}

func (a *BuiltinAuthenticator) SaveTeam(ctx *auth.Context, team *auth.Team) (string, error) {
	var (
		eventType string
	)

	// check if exists; if so, update
	grp, err := a.GetTeam(ctx, team.Id)
	if err != nil && err != auth.ErrTeamDoesNotExist {
		return "", err
	}

	// update
	if grp != nil {
		eventType = "update-team"
	} else {
		eventType = "add-team"
		id := utils.GenerateID(16)
		team.Id = id
	}

	// update default orgId if missing
	if team.OrgId == "" {
		team.OrgId = orca.UCPDefaultOrg
	}

	// Make sure there isn't any membership info written out
	team.ManagedMembers = []string{}

	// kv
	data, err := json.Marshal(team)
	if err != nil {
		return "", err
	}

	k := path.Join(datastoreTeams, team.Id)
	if err := a.store.Put(k, data, nil); err != nil {
		return "", utils.MaybeWrapEtcdClusterErr(err)
	}

	return eventType, nil
}

func (a *BuiltinAuthenticator) DeleteTeam(ctx *auth.Context, team *auth.Team) error {
	// Figure out which members we've got and update them so they no longer have this team listed
	members, _ := a.getTeamMembers(team.Id)
	for _, memberName := range members {
		k := path.Join(datastoreAccountTeams, memberName, team.Id)
		if err := a.store.Delete(k); err != nil && err != kvstore.ErrKeyNotFound {
			// Should we allow partial failures to fall through here so you can delete a busted team?
			// My concern is that may lead to more ominous corruption going undetected (stale memberships)
			// so we'll fail fast and force the user to take other corrective action to get the deletion done
			return fmt.Errorf("Failed to delelete membership for account %s from team %s: %s", memberName, team.Name, err)
		}
	}
	k := path.Join(datastoreTeamMembers, team.Id)
	if err := a.store.Delete(k); err != nil && err != kvstore.ErrKeyNotFound {
		return fmt.Errorf("Failed to delelete membership for team %s: %s", team.Name, err)
	}

	// delete the team
	k = path.Join(datastoreTeams, team.Id)
	if err := a.store.Delete(k); err != nil {
		return err
	}

	return nil
}

func (a *BuiltinAuthenticator) AddTeamMember(ctx *auth.Context, teamID, username string) error {
	// verify user exists
	if _, err := a.GetUser(ctx, username); err != nil {
		return err
	}

	k := path.Join(datastoreTeamMembers, teamID, username)
	if err := a.store.Put(k, []byte(username), nil); err != nil {
		return utils.MaybeWrapEtcdClusterErr(err)
	}

	k = path.Join(datastoreAccountTeams, username, teamID)
	if err := a.store.Put(k, []byte(teamID), nil); err != nil {
		return utils.MaybeWrapEtcdClusterErr(err)
	}

	return nil
}

func (a *BuiltinAuthenticator) DeleteTeamMember(ctx *auth.Context, teamID, username string) error {
	k := path.Join(datastoreTeamMembers, teamID, username)

	if err := a.store.Delete(k); err != nil {
		return err
	}

	k = path.Join(datastoreAccountTeams, username, teamID)

	if err := a.store.Delete(k); err != nil {
		return err
	}

	return nil
}

func (a *BuiltinAuthenticator) getTeamMembers(teamID string) ([]string, error) {
	members := []string{}

	k := path.Join(datastoreTeamMembers, teamID)
	kvPairs, err := a.store.List(k)
	if err != nil && err == kvstore.ErrKeyNotFound {
		return members, nil
	} else if err != nil {
		return nil, err
	}

	for _, kvPair := range kvPairs {
		members = append(members, string(kvPair.Value))
	}

	log.Debugf("members: teamId=%s members=%v", teamID, members)

	return members, nil
}

func (a *BuiltinAuthenticator) getAccountTeams(username string) ([]string, error) {
	teams := []string{}

	k := path.Join(datastoreAccountTeams, username)
	kvPairs, err := a.store.List(k)
	if err != nil && err == kvstore.ErrKeyNotFound {
		return teams, nil
	} else if err != nil {
		return nil, err
	}

	for _, kvPair := range kvPairs {
		teams = append(teams, string(kvPair.Value))
	}

	log.Debugf("account teams: acct=%s teams=%v", username, teams)

	return teams, nil
}

func (a *BuiltinAuthenticator) ListUserTeams(ctx *auth.Context, username string) ([]*auth.Team, error) {
	teams := []*auth.Team{}

	teamIds, err := a.getAccountTeams(username)
	if err != nil {
		return nil, err
	}

	for _, id := range teamIds {
		t, err := a.GetTeam(ctx, id)
		if err != nil {
			return nil, err
		}

		teams = append(teams, t)
	}

	return teams, nil
}

func (a *BuiltinAuthenticator) Sync(ctx *auth.Context, force, adminOnly bool) error {
	// nothing to see here, move along
	return nil
}

func (a *BuiltinAuthenticator) LastSyncStatus(ctx *auth.Context) string {
	// Minor hack - when switching from builtin to ldap
	// during the transition, the API call for sync status
	// will show up here, so we'll return the sync data
	kvpair, err := a.store.Get(datastoreLastSyncMessages)
	if err != nil {
		return ""
	}
	return string(kvpair.Value)
}

// Iterate through all the accounts and make sure they're valid
func (a *BuiltinAuthenticator) updateOldAccounts() {
	log.Info("Checking for any accounts that need schema updates")
	// Note we use custom structures to preserve the old data model
	type (
		Version7 struct { // Pre 0.8
			FirstName string   `json:"first_name,omitempty"`
			LastName  string   `json:"last_name,omitempty"`
			Username  string   `json:"username,omitempty"`
			Password  string   `json:"password,omitempty"`
			Roles     []string `json:"roles,omitempty"`
			PublicKey string   `json:"public_key,omitempty"`
		}
		Version8 struct {
			Version7
			Admin           bool              `json:"admin,omitempty"`
			PublicKeys      []auth.AccountKey `json:"public_keys,omitempty"`
			ManagedTeams    []string          `json:"managed_teams,omitempty"`
			DiscoveredTeams []string          `json:"discovered_teams,omitempty"`
			LdapDN          string            `json:"ldap_dn,omitempty"`
			Disabled        bool              `json:"disabled,omitempty"` // Only used for Discovered accounts
		}
	)

	kvPairs, err := a.store.List(datastoreAccounts)
	if err != nil && err == kvstore.ErrKeyNotFound {
		// No accounts (corner case)
		return
	} else if err != nil {
		log.Errorf("Failed to get account list %s", err)
		return
	}

	for _, kvPair := range kvPairs {
		var account Version8
		if err := json.Unmarshal(kvPair.Value, &account); err != nil {
			log.Warnf("Failed to process old account for schema update: %s - %s", string(kvPair.Key), err)
			continue
		}

		// Data update algorithm:
		// 1. Roles[] -> Admin conversion
		//
		//(add more checks here as we need them)
		for _, role := range account.Version7.Roles {
			if role == "admin" && !account.Admin {
				log.Infof("Detected old admin account %s, updating schema", account.Username)
				account.Admin = true
				data, err := json.Marshal(account)
				if err != nil {
					log.Warnf("Failed to unmarshal updated account: %s", err)
					continue
				}

				k := path.Join(datastoreAccounts, account.Username)
				if err := a.store.Put(k, data, nil); err != nil {
					err = utils.MaybeWrapEtcdClusterErr(err)
					log.Warnf("Failed to save updated account: %s", err)
					continue
				}
			}
		}
	}

	log.Info("Completed account schema update check")
}
