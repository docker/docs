package idauth

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/docker/distribution/context"
	"github.com/docker/orca/enzi/api/client"
	"github.com/docker/orca/enzi/jose"
	"github.com/stretchr/testify/require"
)

/*
This file tests the IDToken Authenticator to make sure that it can successfully
verify tokens from the issuer.
*/

type testKeyServer struct {
	keys  jose.JWKSet
	calls int
}

// The client should only ever need to list signing keys so that is all that
// this server does.
func (s *testKeyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.calls++
	json.NewEncoder(w).Encode(s.keys)
}

func (s *testKeyServer) addKey(key *jose.PublicKey) {
	s.keys.Keys = append(s.keys.Keys, key)
}

func testClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // Because test.
			},
		},
	}
}

func testServerAddr(t *testing.T, s *httptest.Server) string {
	parsed, err := url.Parse(s.URL)
	require.NoError(t, err)

	return parsed.Host
}

func generateKey(t *testing.T) *jose.PrivateKey {
	cryptoPrivateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	key, err := jose.NewPrivateKey(cryptoPrivateKey)
	require.NoError(t, err)

	return key
}

func makeToken(t *testing.T, key *jose.PrivateKey, claims jose.JWTClaims) string {
	token, err := jose.NewJWT(key, claims)
	require.NoError(t, err)

	return token.String()
}

func prepTest(t *testing.T) (*testKeyServer, *httptest.Server, *Authenticator) {
	keyServer := &testKeyServer{
		keys: jose.JWKSet{},
	}

	testServer := httptest.NewTLSServer(keyServer)

	enziClient := client.New(testClient(), testServerAddr(t, testServer), "", nil)

	keyCache := NewInMemoryKeyCache()
	keyCache.Start()

	authenticator := &Authenticator{
		Context:          context.Background(),
		APIClient:        enziClient,
		IssuerIdentifier: "https://auth.example.com",
		ServiceID:        "testService",
		SigningKeyCache:  keyCache,
	}

	return keyServer, testServer, authenticator
}

func TestIDTokenAuthenticatorSuccess(t *testing.T) {
	keyServer, testServer, authenticator := prepTest(t)
	defer testServer.Close()
	defer authenticator.SigningKeyCache.Stop()

	// Generate a couple of keys and use them to sign some identity tokens.
	keyA := generateKey(t)
	keyB := generateKey(t)

	claims := jose.JWTClaims{
		Issuer:          authenticator.IssuerIdentifier,
		Subject:         "testSubject",
		AuthorizedParty: "someone",
		Audience:        []string{authenticator.ServiceID},
		IssuedAt:        time.Now().Unix(),
		Expiration:      time.Now().Add(time.Minute).Unix(),
	}

	tokenA := makeToken(t, keyA, claims)
	tokenB := makeToken(t, keyB, claims)

	// Add the keys to the key server so that the authenticator can look
	// them up and cache them.
	keyServer.addKey(&keyA.PublicKey)
	keyServer.addKey(&keyB.PublicKey)

	for _, token := range []string{tokenA, tokenB} {
		verifiedClaims, err := authenticator.AuthenticateIdentityToken(token)
		require.NoError(t, err)
		require.Equal(t, claims.Subject, verifiedClaims.Subject)
	}

	// The authenticator should have only called the key server once.
	require.Exactly(t, 1, keyServer.calls)
}

func TestIDTokenAuthenticatorInvalid(t *testing.T) {
	keyServer, testServer, authenticator := prepTest(t)
	defer testServer.Close()
	defer authenticator.SigningKeyCache.Stop()

	// Generate a couple of keys and use them to sign some identity tokens.
	keyA := generateKey(t)
	keyB := generateKey(t)

	claims := jose.JWTClaims{
		Issuer:          "notTrusted", // Not a trusted issuer.
		Subject:         "testSubject",
		AuthorizedParty: "someone",
		Audience:        []string{"notMe"}, // Wrong audience.
		IssuedAt:        time.Now().Unix(),
		Expiration:      time.Now().Add(-5 * time.Minute).Unix(), // Already expired.
	}

	// Only add keyA to the key server. The authenticator should fail to
	// verify a token signed by keyA since it wasn't on the token server.
	keyServer.addKey(&keyA.PublicKey)

	_, err := authenticator.AuthenticateIdentityToken(makeToken(t, keyB, claims))
	require.Error(t, err)
	require.Contains(t, err.Error(), "unable to get verifying key")

	// Adding the other key should fix that error.
	keyServer.addKey(&keyB.PublicKey)

	_, err = authenticator.AuthenticateIdentityToken(makeToken(t, keyB, claims))
	require.Error(t, err)
	require.Contains(t, err.Error(), "JWT intended for a different audience")

	// Changing the claims audience should fix that error.
	claims.Audience = []string{authenticator.ServiceID}

	_, err = authenticator.AuthenticateIdentityToken(makeToken(t, keyB, claims))
	require.Error(t, err)
	require.Contains(t, err.Error(), "JWT issuer untrusted")

	// Changing the claims issuer should fix that error.
	claims.Issuer = authenticator.IssuerIdentifier

	_, err = authenticator.AuthenticateIdentityToken(makeToken(t, keyB, claims))
	require.Error(t, err)
	require.Contains(t, err.Error(), "JWT expired")

	// Making it expire later should fix that error.
	claims.Expiration = time.Now().Add(time.Minute).Unix()

	_, err = authenticator.AuthenticateIdentityToken(makeToken(t, keyB, claims))
	require.NoError(t, err)

	// The authenticator should have only called the key server twice:
	// once for the original token and again after the other signing key
	// was added.
	require.Exactly(t, 2, keyServer.calls)
}
