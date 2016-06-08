package keydbstore

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRDBTUFFileMarshalling(t *testing.T) {
	created := time.Now().AddDate(-1, -1, -1)
	updated := time.Now().AddDate(0, -5, 0)
	deleted := time.Time{}

	createdMarshalled, err := json.Marshal(created)
	require.NoError(t, err)
	updatedMarshalled, err := json.Marshal(updated)
	require.NoError(t, err)
	deletedMarshalled, err := json.Marshal(deleted)
	require.NoError(t, err)

	jsonBytes := []byte(fmt.Sprintf(`
	{
		"created_at": %s,
		"updated_at": %s,
		"deleted_at": %s,
		"key_id": "56ee4a23129fc22c6cb4b4ba5f78d730c91ab6def514e80d807c947bb21f0d63",
		"encryption_alg": "A256GCM",
		"keywrap_alg": "PBES2-HS256+A128KW",
		"algorithm": "ecdsa",
		"passphrase_alias": "timestamp_1",
		"public": "Hello world public",
		"private": "Hello world private"
	}
	`, createdMarshalled, updatedMarshalled, deletedMarshalled))
	fmt.Println(string(jsonBytes))

	unmarshalledAnon, err := PrivateKeysRethinkTable.JSONUnmarshaller(jsonBytes)
	require.NoError(t, err)
	unmarshalled, ok := unmarshalledAnon.(RDBPrivateKey)
	require.True(t, ok)

	// There is some weirdness with comparing time.Time due to a location pointer,
	// so let's use time.Time's equal function to compare times, and then re-assign
	// the timing struct to compare the rest of the RDBTUFFile struct
	require.True(t, created.Equal(unmarshalled.CreatedAt))
	require.True(t, updated.Equal(unmarshalled.UpdatedAt))
	require.True(t, deleted.Equal(unmarshalled.DeletedAt))

	expected := RDBPrivateKey{
		Timing:          unmarshalled.Timing,
		KeyID:           "56ee4a23129fc22c6cb4b4ba5f78d730c91ab6def514e80d807c947bb21f0d63",
		EncryptionAlg:   "A256GCM",
		KeywrapAlg:      "PBES2-HS256+A128KW",
		Algorithm:       "ecdsa",
		PassphraseAlias: "timestamp_1",
		Public:          "Hello world public",
		Private:         "Hello world private",
	}
	require.Equal(t, expected, unmarshalled)
}
