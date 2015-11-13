package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/docker/notary/passphrase"
	"github.com/docker/notary/trustmanager"
	"github.com/docker/notary/tuf/data"
	"github.com/stretchr/testify/assert"
)

func TestTruncateWithEllipsis(t *testing.T) {
	digits := "1234567890"
	// do not truncate
	assert.Equal(t, truncateWithEllipsis(digits, 10, true), digits)
	assert.Equal(t, truncateWithEllipsis(digits, 10, false), digits)
	assert.Equal(t, truncateWithEllipsis(digits, 11, true), digits)
	assert.Equal(t, truncateWithEllipsis(digits, 11, false), digits)

	// left and right truncate
	assert.Equal(t, truncateWithEllipsis(digits, 8, true), "...67890")
	assert.Equal(t, truncateWithEllipsis(digits, 8, false), "12345...")
}

func TestKeyInfoSorter(t *testing.T) {
	expected := []keyInfo{
		{role: data.CanonicalRootRole, gun: "", keyID: "a", location: "i"},
		{role: data.CanonicalRootRole, gun: "", keyID: "a", location: "j"},
		{role: data.CanonicalRootRole, gun: "", keyID: "z", location: "z"},
		{role: "a", gun: "a", keyID: "a", location: "y"},
		{role: "b", gun: "a", keyID: "a", location: "y"},
		{role: "b", gun: "a", keyID: "b", location: "y"},
		{role: "b", gun: "a", keyID: "b", location: "z"},
		{role: "a", gun: "b", keyID: "a", location: "z"},
	}
	jumbled := make([]keyInfo, len(expected))
	// randomish indices
	for j, e := range []int{3, 6, 1, 4, 0, 7, 5, 2} {
		jumbled[j] = expected[e]
	}

	sort.Sort(keyInfoSorter(jumbled))
	assert.True(t, reflect.DeepEqual(expected, jumbled),
		fmt.Sprintf("Expected %v, Got %v", expected, jumbled))
}

type otherMemoryStore struct {
	trustmanager.KeyMemoryStore
}

func (l *otherMemoryStore) Name() string {
	return strings.Repeat("z", 70)
}

// Given a list of key stores, the keys should be pretty-printed with their
// roles, locations, IDs, and guns first in sorted order in the key store
func TestPrettyPrintKeys(t *testing.T) {
	ret := passphrase.ConstantRetriever("pass")
	keyStores := []trustmanager.KeyStore{
		trustmanager.NewKeyMemoryStore(ret),
		&otherMemoryStore{KeyMemoryStore: *trustmanager.NewKeyMemoryStore(ret)},
	}

	longNameShortened := "..." + strings.Repeat("z", 37)

	// just use the same key for testing
	key, err := trustmanager.GenerateED25519Key(rand.Reader)
	assert.NoError(t, err)

	root := data.CanonicalRootRole

	// add keys to the key stores
	err = keyStores[0].AddKey(key.ID(), root, key)
	assert.NoError(t, err)

	err = keyStores[1].AddKey(key.ID(), root, key)
	assert.NoError(t, err)

	err = keyStores[0].AddKey(strings.Repeat("a/", 30)+key.ID(), "targets", key)
	assert.NoError(t, err)

	err = keyStores[1].AddKey("short/gun/"+key.ID(), "snapshot", key)
	assert.NoError(t, err)

	expected := [][]string{
		{root, key.ID(), keyStores[0].Name()},
		{root, key.ID(), longNameShortened},
		{"targets", "..." + strings.Repeat("/a", 11), key.ID(), keyStores[0].Name()},
		{"snapshot", "short/gun", key.ID(), longNameShortened},
	}

	var b bytes.Buffer
	prettyPrintKeys(keyStores, &b)
	text, err := ioutil.ReadAll(&b)
	assert.NoError(t, err)

	lines := strings.Split(strings.TrimSpace(string(text)), "\n")
	assert.Len(t, lines, len(expected)+2)
	for i, line := range lines[2:] {
		// we are purposely not putting spaces in test data so easier to split
		splitted := strings.Fields(line)
		for j, v := range splitted {
			assert.Equal(t, expected[i][j], strings.TrimSpace(v))
		}
	}
}
