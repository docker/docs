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

var ret = passphrase.ConstantRetriever("pass")

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
func TestPrettyPrintRootAndSigningKeys(t *testing.T) {
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

	// starts with headers
	assert.True(t, reflect.DeepEqual(strings.Fields(lines[0]),
		[]string{"ROLE", "GUN", "KEY", "ID", "LOCATION"}))
	assert.Equal(t, "----", lines[1][:4])

	for i, line := range lines[2:] {
		// we are purposely not putting spaces in test data so easier to split
		splitted := strings.Fields(line)
		for j, v := range splitted {
			assert.Equal(t, expected[i][j], strings.TrimSpace(v))
		}
	}
}

// If there are no keys in any of the key stores, a message that there are no
// signing keys should be displayed.
func TestPrettyPrintZeroKeys(t *testing.T) {
	ret := passphrase.ConstantRetriever("pass")
	emptyKeyStore := trustmanager.NewKeyMemoryStore(ret)

	var b bytes.Buffer
	prettyPrintKeys([]trustmanager.KeyStore{emptyKeyStore}, &b)
	text, err := ioutil.ReadAll(&b)
	assert.NoError(t, err)

	lines := strings.Split(strings.TrimSpace(string(text)), "\n")
	assert.Len(t, lines, 1)
	assert.Equal(t, "No signing keys found.", lines[0])
}

// If there are no keys, removeKeyInteractively will just return an error about
// there not being any key
func TestRemoveIfNoKey(t *testing.T) {
	var buf bytes.Buffer
	stores := []trustmanager.KeyStore{trustmanager.NewKeyMemoryStore(nil)}
	err := removeKeyInteractively(stores, "12345", &buf, &buf)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "No key with ID")
}

// If there is one key, asking to remove it will ask for confirmation.  Passing
// anything other than 'yes'/'y'/'' response will abort the deletion and
// not delete the key.
func TestRemoveOneKeyAbort(t *testing.T) {
	nos := []string{"no", "NO", "AAAARGH", "   N    "}
	store := trustmanager.NewKeyMemoryStore(ret)

	key, err := trustmanager.GenerateED25519Key(rand.Reader)
	assert.NoError(t, err)
	err = store.AddKey(key.ID(), "root", key)
	assert.NoError(t, err)

	stores := []trustmanager.KeyStore{store}

	for _, noAnswer := range nos {
		var out bytes.Buffer
		in := bytes.NewBuffer([]byte(noAnswer + "\n"))

		err := removeKeyInteractively(stores, key.ID(), in, &out)
		assert.NoError(t, err)
		text, err := ioutil.ReadAll(&out)
		assert.NoError(t, err)

		output := string(text)
		assert.Contains(t, output, "Are you sure")
		assert.Contains(t, output, "Aborting action")
		assert.Len(t, store.ListKeys(), 1)
	}
}

// If there is one key, asking to remove it will ask for confirmation.  Passing
// 'yes'/'y'/'' response will continue the deletion.
func TestRemoveOneKeyConfirm(t *testing.T) {
	yesses := []string{"yes", " Y ", "yE", "   ", ""}

	for _, yesAnswer := range yesses {
		store := trustmanager.NewKeyMemoryStore(ret)

		key, err := trustmanager.GenerateED25519Key(rand.Reader)
		assert.NoError(t, err)
		err = store.AddKey(key.ID(), "root", key)
		assert.NoError(t, err)

		var out bytes.Buffer
		in := bytes.NewBuffer([]byte(yesAnswer + "\n"))

		err = removeKeyInteractively(
			[]trustmanager.KeyStore{store}, key.ID(), in, &out)
		assert.NoError(t, err)
		text, err := ioutil.ReadAll(&out)
		assert.NoError(t, err)

		output := string(text)
		assert.Contains(t, output, "Are you sure")
		assert.Contains(t, output, "Deleted "+key.ID())
		assert.Len(t, store.ListKeys(), 0)
	}
}

// If there is more than one key, removeKeyInteractively will ask which key to
// delete and will do so over and over until the user quits if the answer is
// invalid.
func TestRemoveMultikeysInvalidInput(t *testing.T) {
	in := bytes.NewBuffer([]byte("nota number\n9999\n-3\n0"))

	key, err := trustmanager.GenerateED25519Key(rand.Reader)
	assert.NoError(t, err)

	stores := []trustmanager.KeyStore{
		trustmanager.NewKeyMemoryStore(ret),
		trustmanager.NewKeyMemoryStore(ret),
	}

	err = stores[0].AddKey(key.ID(), "root", key)
	assert.NoError(t, err)

	err = stores[1].AddKey("gun/"+key.ID(), "target", key)
	assert.NoError(t, err)

	var out bytes.Buffer

	err = removeKeyInteractively(stores, key.ID(), in, &out)
	assert.Error(t, err)
	text, err := ioutil.ReadAll(&out)
	assert.NoError(t, err)

	assert.Len(t, stores[0].ListKeys(), 1)
	assert.Len(t, stores[1].ListKeys(), 1)

	// It should have listed the keys over and over, asking which key the user
	// wanted to delete
	output := string(text)
	assert.Contains(t, output, "Found the following matching keys")
	var rootCount, targetCount int
	for _, line := range strings.Split(output, "\n") {
		if strings.Contains(line, key.ID()) {
			if strings.Contains(line, "target") {
				targetCount++
			} else {
				rootCount++
			}
		}
	}
	assert.Equal(t, rootCount, targetCount)
	assert.Equal(t, 4, rootCount) // for each of the 4 invalid inputs
}

// If there is more than one key, removeKeyInteractively will ask which key to
// delete.  Then it will confirm whether they want to delete, and the user can
// abort at that confirmation.
func TestRemoveMultikeysAbortChoice(t *testing.T) {
	in := bytes.NewBuffer([]byte("1\nn\n"))

	key, err := trustmanager.GenerateED25519Key(rand.Reader)
	assert.NoError(t, err)

	stores := []trustmanager.KeyStore{
		trustmanager.NewKeyMemoryStore(ret),
		trustmanager.NewKeyMemoryStore(ret),
	}

	err = stores[0].AddKey(key.ID(), "root", key)
	assert.NoError(t, err)

	err = stores[1].AddKey("gun/"+key.ID(), "target", key)
	assert.NoError(t, err)

	var out bytes.Buffer

	err = removeKeyInteractively(stores, key.ID(), in, &out)
	assert.NoError(t, err) // no error to abort deleting
	text, err := ioutil.ReadAll(&out)
	assert.NoError(t, err)

	assert.Len(t, stores[0].ListKeys(), 1)
	assert.Len(t, stores[1].ListKeys(), 1)

	// It should have listed the keys, asked whether the user really wanted to
	// delete, and then aborted.
	output := string(text)
	assert.Contains(t, output, "Found the following matching keys")
	assert.Contains(t, output, "Are you sure")
	assert.Contains(t, output, "Aborting action")
}

// If there is more than one key, removeKeyInteractively will ask which key to
// delete.  Then it will confirm whether they want to delete, and if the user
// confirms, will remove it from the correct key store.
func TestRemoveMultikeysRemoveOnlyChosenKey(t *testing.T) {
	in := bytes.NewBuffer([]byte("1\ny\n"))

	key, err := trustmanager.GenerateED25519Key(rand.Reader)
	assert.NoError(t, err)

	stores := []trustmanager.KeyStore{
		trustmanager.NewKeyMemoryStore(ret),
		trustmanager.NewKeyMemoryStore(ret),
	}

	err = stores[0].AddKey(key.ID(), "root", key)
	assert.NoError(t, err)

	err = stores[1].AddKey("gun/"+key.ID(), "target", key)
	assert.NoError(t, err)

	var out bytes.Buffer

	err = removeKeyInteractively(stores, key.ID(), in, &out)
	assert.NoError(t, err)
	text, err := ioutil.ReadAll(&out)
	assert.NoError(t, err)

	// It should have listed the keys, asked whether the user really wanted to
	// delete, and then deleted.
	output := string(text)
	assert.Contains(t, output, "Found the following matching keys")
	assert.Contains(t, output, "Are you sure")
	assert.Contains(t, output, "Deleted "+key.ID())

	// figure out which one we picked to delete, and assert it was deleted
	for _, line := range strings.Split(output, "\n") {
		if strings.HasPrefix(line, "\t1.") { // we picked the first item
			if strings.Contains(line, "root") { // first key store
				assert.Len(t, stores[0].ListKeys(), 0)
				assert.Len(t, stores[1].ListKeys(), 1)
			} else {
				assert.Len(t, stores[0].ListKeys(), 1)
				assert.Len(t, stores[1].ListKeys(), 0)
			}
		}
	}
}
