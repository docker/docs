package idauth

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"
	"time"

	"github.com/docker/orca/enzi/jose"
	"github.com/stretchr/testify/require"
)

func TestKeyCache(t *testing.T) {
	cache := NewInMemoryKeyCache().(*keyCache)

	cache.Start()
	defer cache.Stop()

	// Generate 10 entries expiring over the next five seconds.
	keys := make([]*jose.PublicKey, 10)

	fullDuration := 5 * time.Second
	expirationInterval := fullDuration / 10

	// Use a separate loop to generate the keys since it is slow.
	for i := range keys {
		key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		require.NoError(t, err)

		keys[i], err = jose.NewPublicKey(key.Public())
		require.NoError(t, err)
	}

	for i := range keys {
		expiration := time.Duration(i+1) * expirationInterval
		cache.Set(keys[i], expiration)
	}

	// There should now be 10 keys in the cache.
	require.Len(t, cache.cachedKeys, 10)
	require.Len(t, cache.lookup, 10)

	for _, key := range keys {
		cached := cache.Get(key.ID)
		require.NotNil(t, cached)
		require.Equal(t, key.ID, cached.ID)
	}

	// Offset by half of an interval.
	time.Sleep(expirationInterval / 2)

	for i := 0; i < 5; i++ {
		time.Sleep(expirationInterval)

		// keys[i] should no longer be in the cache.
		require.Len(t, cache.cachedKeys, 9-i)
		require.Len(t, cache.lookup, 9-i)
		_, exists := cache.lookup[keys[i].ID]
		require.False(t, exists)
	}

	// Add those 5 keys back, expiring in reverse order over half a second
	// starting a seccond from now.
	for i := 0; i < 5; i++ {
		expiration := fullDuration - time.Duration(i)*expirationInterval
		cache.Set(keys[i], expiration)
	}

	// There should again be 10 keys in the cache.
	require.Len(t, cache.cachedKeys, 10)
	require.Len(t, cache.lookup, 10)

	for _, key := range keys {
		cached := cache.Get(key.ID)
		require.NotNil(t, cached)
		require.Equal(t, key.ID, cached.ID)
	}

	// keys 5-9 should expire over the next half second.
	for i := 0; i < 5; i++ {
		time.Sleep(expirationInterval)

		// keys[i+5] should no longer be in the cache.
		require.Len(t, cache.cachedKeys, 9-i)
		require.Len(t, cache.lookup, 9-i)
		_, exists := cache.lookup[keys[i+5].ID]
		require.False(t, exists)
	}

	// There should now be 5 keys in the cache.
	require.Len(t, cache.cachedKeys, 5)
	require.Len(t, cache.lookup, 5)

	for i, key := range keys {
		cached := cache.Get(key.ID)
		if i < 5 {
			require.NotNil(t, cached)
			require.Equal(t, key.ID, cached.ID)
		} else {
			require.Nil(t, cached)
		}
	}

	// keys 0-4 should expire again in reverse order.
	for i := 0; i < 5; i++ {
		time.Sleep(expirationInterval)

		// keys[4-i] should no longer be in the cache.
		require.Len(t, cache.cachedKeys, 4-i)
		require.Len(t, cache.lookup, 4-i)
		_, exists := cache.lookup[keys[4-i].ID]
		require.False(t, exists)
	}
}
