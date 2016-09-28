package schema

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/docker/dhe-deploy"
	"github.com/docker/garant/auth"
	"github.com/pborman/uuid"
	rethink "gopkg.in/dancannon/gorethink.v2"
)

var (
	ErrClientTokenExists = fmt.Errorf("client token already exists")
	ErrNoSuchClientToken = fmt.Errorf("no such client token")
)

type ClientToken struct {
	// auth.RefreshToken properties
	Token    string `gorethink:"-"`
	ClientID string `gorethink:"clientID"`
	// custom properties

	// AccountID refers to the account ID within enzi. This can be used as the
	// `refresh_token` parameter within a refresh_token grant, allowing us continued
	// auth access.
	AccountID   string    `gorethink:"accountID"`
	HashedToken string    `gorethink:"token"`
	LastFour    string    `gorethink:"lastFour"`
	CreatorIP   string    `gorethink:"creatorIP"`
	CreatorUA   string    `gorethink:"creatorUA"`
	CreatedAt   time.Time `gorethink:"createdAt"`
}

// RefreshToken returns an auth.RefreshToken which is a struct
// containing a subset of the ClientToken fields used in Garant.
func (ct *ClientToken) RefreshToken() auth.RefreshToken {
	return auth.RefreshToken{
		Token:    ct.Token,
		ClientID: ct.ClientID,
	}
}

// NewClientToken generates a new ClientToken with Token set to a random UUID
// and HashedToken set to its hash
func NewClientToken() *ClientToken {
	token := new(ClientToken)
	token.Token = uuid.New()
	token.HashedToken = hashToken(token.Token)
	token.LastFour = token.Token[len(token.Token)-4:]
	return token
}

var clientTokenTable = table{
	db:         deploy.DTRDBName,
	name:       "client_tokens",
	primaryKey: "token", // Guarantees uniqueness of token
	secondaryIndexes: map[string][]string{
		"accountID_clientID": {"accountID", "clientID"}, // For looking up all tokens for a specific user/client
	},
}

type ClientTokenManager struct {
	session *rethink.Session
}

func NewClientTokenManager(s *rethink.Session) *ClientTokenManager {
	return &ClientTokenManager{s}
}

// CreateClientToken stores a new client token in the client token table.
// Note that the ClientToken must have all struct fields present.
func (m *ClientTokenManager) CreateClientToken(token *ClientToken) error {
	if resp, err := clientTokenTable.Term().Insert(token).RunWrite(m.session); err != nil {
		if isDuplicatePrimaryKeyErr(resp) {
			return ErrClientTokenExists
		}
		return fmt.Errorf("unable to create client token in database: %s", err)
	}
	return nil
}

// GetClientToken returns a complete ClientToken with the given token.
// Note that the token must be passed unhashed.
func (m *ClientTokenManager) GetClientToken(token string) (ct *ClientToken, err error) {
	ct = &ClientToken{}
	err = clientTokenTable.getRowByIndexVal(m.session, "token", hashToken(token), ct, ErrNoSuchClientToken)
	return
}

// ListAccountTokens retrieves all stored client tokens for the given account ID
func (m *ClientTokenManager) ListAccountTokens(accountID string) (tokens []*ClientToken, err error) {
	cursor, err := clientTokenTable.Term().Between(
		[]interface{}{accountID, rethink.MinVal},
		[]interface{}{accountID, rethink.MaxVal},
		rethink.BetweenOpts{Index: "accountID_clientID"},
	).Run(m.session)
	if err != nil {
		return nil, fmt.Errorf("unable to query for account tokens: %s", err)
	}
	if err = cursor.All(&tokens); err != nil {
		return nil, fmt.Errorf("unable to scan query results: %s", err)
	}
	return tokens, nil
}

// RevokeToken revokes a token by the unhashed token.
func (m *ClientTokenManager) RevokeToken(unhashedToken string) error {
	return m.RevokeHashedToken(hashToken(unhashedToken))
}

// RevokeHashedToken revokes a token by the token hash
func (m *ClientTokenManager) RevokeHashedToken(hashedToken string) error {
	if _, err := clientTokenTable.Term().Get(hashedToken).Delete().Run(m.session); err != nil {
		return fmt.Errorf("unable to delete client token: %s", err)
	}
	return nil
}

// RevokeByClientID revokes all of a single account's tokens for a specific client.
// This can be used to log out of all docker daemons which are logged in by oauth
func (m *ClientTokenManager) RevokeByClientID(accountID, clientID string) error {
	err := clientTokenTable.Term().GetAllByIndex(
		"accountID_clientID",
		[]interface{}{accountID, clientID},
	).Delete().Exec(m.session)

	if err != nil {
		return fmt.Errorf("unable to revoke tokens by client id: %s", err)
	}
	return nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
