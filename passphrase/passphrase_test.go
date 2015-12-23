package passphrase

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/docker/notary/tuf/data"
	"github.com/stretchr/testify/assert"
)

func assertAskOnceForKey(t *testing.T, in, out *bytes.Buffer, retriever Retriever, password, role string) {
	_, err := in.WriteString(password + "\n")
	assert.NoError(t, err)

	pass, giveUp, err := retriever("repo/0123456789abcdef", role, false, 0)
	assert.NoError(t, err)
	assert.False(t, giveUp)
	assert.Equal(t, password, pass)

	text, err := ioutil.ReadAll(out)
	assert.NoError(t, err)
	assert.Equal(t, "Enter passphrase for "+role+" key with ID 0123456 (repo):",
		strings.TrimSpace(string(text)))
}

// PromptRetrieverWithInOut prompts for delegations passwords (non creation) if needed
func TestGetPassphraseForUsingDelegationKey(t *testing.T) {
	var in bytes.Buffer
	var out bytes.Buffer

	retriever := PromptRetrieverWithInOut(&in, &out, nil)

	for i := 0; i < 3; i++ {
		target := fmt.Sprintf("targets/level%d", i)
		password := fmt.Sprintf("randompass%d", i)

		assertAskOnceForKey(t, &in, &out, retriever, password, target)
	}
}

// PromptRetrieverWithInOut prompts for creating delegations passwords if needed
func TestGetPassphraseForCreatingDelegationKey(t *testing.T) {
	var in bytes.Buffer
	var out bytes.Buffer

	retriever := PromptRetrieverWithInOut(&in, &out, nil)

	_, err := in.WriteString("passphrase\npassphrase\n")
	assert.NoError(t, err)

	pass, giveUp, err := retriever("repo/0123456789abcdef", "targets/a", true, 0)
	assert.NoError(t, err)
	assert.False(t, giveUp)
	assert.Equal(t, "passphrase", pass)

	text, err := ioutil.ReadAll(&out)
	assert.NoError(t, err)
	lines := strings.Split(strings.TrimSpace(string(text)), "\n")

	expectedText := []string{
		`Enter passphrase for new targets/a key with ID 0123456 (repo): `,
		`Repeat passphrase for new targets/a key with ID 0123456 (repo):`,
	}

	assert.Equal(t, expectedText, lines)
}

// PromptRetrieverWithInOut, if asked for root, targets, delegation, and
// snapshot passphrases in that order will only prompt for root, targets, and
// delegation passphrases because it caches the targets password and uses it
// for snapshot.
func TestGetRootTargetsDelegation(t *testing.T) {
	var in bytes.Buffer
	var out bytes.Buffer

	retriever := PromptRetrieverWithInOut(&in, &out, nil)

	assertAskOnceForKey(t, &in, &out, retriever, "rootpassword", data.CanonicalRootRole)
	assertAskOnceForKey(t, &in, &out, retriever, "targetspassword", data.CanonicalTargetsRole)
	assertAskOnceForKey(t, &in, &out, retriever, "delegationpass", "targets/delegation")

	// now ask for snapshot password, but it should already be cached, it
	// won't ask and  no input necessary.
	pass, giveUp, err := retriever("repo/0123456789abcdef", data.CanonicalSnapshotRole, false, 0)
	assert.NoError(t, err)
	assert.False(t, giveUp)
	assert.Equal(t, "targetspassword", pass)

	text, err := ioutil.ReadAll(&out)
	assert.NoError(t, err)
	assert.Empty(t, text)
}
