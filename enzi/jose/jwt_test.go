package jose

import (
	"crypto/elliptic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestJWT(t *testing.T) {
	key := generateECDSAKey(t, elliptic.P256())

	now := time.Now()
	expiration := now.Add(time.Minute)

	testCases := []JWTClaims{
		{
			Issuer:          "me",
			Subject:         "her",
			AuthorizedParty: "you",
			Audience:        []string{"me", "you", "him"},
			IssuedAt:        now.Unix(),
			Expiration:      expiration.Unix(),
		},
		{
			Issuer:          "me",
			Subject:         "her",
			AuthorizedParty: "you",
			Audience:        []string{"him"},
			IssuedAt:        now.Unix(),
			Expiration:      expiration.Unix(),
		},
	}

	for _, testCase := range testCases {
		jwt, err := NewJWT(key, testCase)
		require.NoError(t, err)

		require.Equal(t, "JWT", jwt.Header.Type)
		require.Equal(t, key.ID, jwt.Header.KeyID)
		require.Equal(t, testCase.Issuer, jwt.Claims.Issuer)
		require.Equal(t, testCase.Subject, jwt.Claims.Subject)
		require.Equal(t, testCase.AuthorizedParty, jwt.Claims.AuthorizedParty)
		require.Equal(t, testCase.Audience, jwt.Claims.Audience)

		decodedJWT, err := DecodeJWT(jwt.String())
		require.NoError(t, err)

		require.Equal(t, jwt.Header.Type, decodedJWT.Header.Type)
		require.Equal(t, jwt.Header.KeyID, decodedJWT.Header.KeyID)
		require.Equal(t, jwt.Header.SigningAlg, decodedJWT.Header.SigningAlg)
		require.Equal(t, jwt.Claims.Issuer, decodedJWT.Claims.Issuer)
		require.Equal(t, jwt.Claims.Subject, decodedJWT.Claims.Subject)
		require.Equal(t, jwt.Claims.AuthorizedParty, decodedJWT.Claims.AuthorizedParty)
		require.Equal(t, jwt.Claims.Audience, decodedJWT.Claims.Audience)
		require.Equal(t, jwt.Claims.IssuedAt, decodedJWT.Claims.IssuedAt)
		require.Equal(t, jwt.Claims.Expiration, decodedJWT.Claims.Expiration)

		require.NoError(t, decodedJWT.Verify(&key.PublicKey, "him", testCase.Issuer))
	}
}
