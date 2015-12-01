package main

import (
	"bytes"
	"crypto/rand"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/docker/notary/passphrase"
	"github.com/docker/notary/trustmanager"
	"github.com/stretchr/testify/assert"
)

var ret = passphrase.ConstantRetriever("pass")

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
