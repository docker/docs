// Package passphrase is a utility function for managing passphrase
// for TUF and Notary keys.
package passphrase

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/pkg/term"
	"github.com/docker/notary"
)

const (
	idBytesToDisplay            = 7
	tufRootAlias                = "root"
	tufTargetsAlias             = "targets"
	tufSnapshotAlias            = "snapshot"
	tufRootKeyGenerationWarning = `You are about to create a new root signing key passphrase. This passphrase
will be used to protect the most sensitive key in your signing system. Please
choose a long, complex passphrase and be careful to keep the password and the
key file itself secure and backed up. It is highly recommended that you use a
password manager to generate the passphrase and keep it safe. There will be no
way to recover this key. You can find the key in your config directory.`
)

var (
	// ErrTooShort is returned if the passphrase entered for a new key is
	// below the minimum length
	ErrTooShort = errors.New("Passphrase too short")

	// ErrDontMatch is returned if the two entered passphrases don't match.
	// new key is below the minimum length
	ErrDontMatch = errors.New("The entered passphrases do not match")

	// ErrTooManyAttempts is returned if the maximum number of passphrase
	// entry attempts is reached.
	ErrTooManyAttempts = errors.New("Too many attempts")
)

// PromptRetriever returns a new Retriever which will provide a prompt on stdin
// and stdout to retrieve a passphrase. The passphrase will be cached such that
// subsequent prompts will produce the same passphrase.
func PromptRetriever() notary.PassRetriever {
	return PromptRetrieverWithInOut(os.Stdin, os.Stdout, nil)
}

type boundRetriever struct {
	in              io.Reader
	out             io.Writer
	aliasMap        map[string]string
	passphraseCache map[string]string
}

func (br *boundRetriever) getPassphrase(keyName, alias string, createNew bool, numAttempts int) (string, bool, error) {
	if numAttempts == 0 {
		if alias == tufRootAlias && createNew {
			fmt.Fprintln(br.out, tufRootKeyGenerationWarning)
		}

		if pass, ok := br.passphraseCache[alias]; ok {
			return pass, false, nil
		}
	} else if !createNew { // per `if`, numAttempts > 0 if we're at this `else`
		if numAttempts > 3 {
			return "", true, ErrTooManyAttempts
		}
		fmt.Fprintln(br.out, "Passphrase incorrect. Please retry.")
	}

	// passphrase not cached and we're not aborting, get passphrase from user!
	return br.requestPassphrase(keyName, alias, createNew, numAttempts)
}

func (br *boundRetriever) requestPassphrase(keyName, alias string, createNew bool, numAttempts int) (string, bool, error) {
	// Figure out if we should display a different string for this alias
	displayAlias := alias
	if val, ok := br.aliasMap[alias]; ok {
		displayAlias = val
	}

	// If typing on the terminal, we do not want the terminal to echo the
	// password that is typed (so it doesn't display)
	if term.IsTerminal(0) {
		state, err := term.SaveState(0)
		if err != nil {
			return "", false, err
		}
		term.DisableEcho(0, state)
		defer term.RestoreTerminal(0, state)
	}

	indexOfLastSeparator := strings.LastIndex(keyName, string(filepath.Separator))
	if indexOfLastSeparator == -1 {
		indexOfLastSeparator = 0
	}

	var shortName string
	if len(keyName) > indexOfLastSeparator+idBytesToDisplay {
		if indexOfLastSeparator > 0 {
			keyNamePrefix := keyName[:indexOfLastSeparator]
			keyNameID := keyName[indexOfLastSeparator+1 : indexOfLastSeparator+idBytesToDisplay+1]
			shortName = keyNameID + " (" + keyNamePrefix + ")"
		} else {
			shortName = keyName[indexOfLastSeparator : indexOfLastSeparator+idBytesToDisplay]
		}
	}

	withID := fmt.Sprintf(" with ID %s", shortName)
	if shortName == "" {
		withID = ""
	}

	switch {
	case createNew:
		fmt.Fprintf(br.out, "Enter passphrase for new %s key%s: ", displayAlias, withID)
	case displayAlias == "yubikey":
		fmt.Fprintf(br.out, "Enter the %s for the attached Yubikey: ", keyName)
	default:
		fmt.Fprintf(br.out, "Enter passphrase for %s key%s: ", displayAlias, withID)
	}

	stdin := bufio.NewReader(br.in)
	passphrase, err := stdin.ReadBytes('\n')
	fmt.Fprintln(br.out)
	if err != nil {
		return "", false, err
	}

	retPass := strings.TrimSpace(string(passphrase))

	if createNew {
		err = br.verifyAndConfirmPassword(stdin, retPass, displayAlias, withID)
		if err != nil {
			return "", false, err
		}
	}

	br.cachePassword(alias, retPass)

	return retPass, false, nil
}

func (br *boundRetriever) verifyAndConfirmPassword(stdin *bufio.Reader, retPass, displayAlias, withID string) error {
	if len(retPass) < 8 {
		fmt.Fprintln(br.out, "Passphrase is too short. Please use a password manager to generate and store a good random passphrase.")
		return ErrTooShort
	}

	fmt.Fprintf(br.out, "Repeat passphrase for new %s key%s: ", displayAlias, withID)
	confirmation, err := stdin.ReadBytes('\n')
	fmt.Fprintln(br.out)
	if err != nil {
		return err
	}
	confirmationStr := strings.TrimSpace(string(confirmation))

	if retPass != confirmationStr {
		fmt.Fprintln(br.out, "Passphrases do not match. Please retry.")
		return ErrDontMatch
	}
	return nil
}

func (br *boundRetriever) cachePassword(alias, retPass string) {
	br.passphraseCache[alias] = retPass
}

// PromptRetrieverWithInOut returns a new Retriever which will provide a
// prompt using the given in and out readers. The passphrase will be cached
// such that subsequent prompts will produce the same passphrase.
// aliasMap can be used to specify display names for TUF key aliases. If aliasMap
// is nil, a sensible default will be used.
func PromptRetrieverWithInOut(in io.Reader, out io.Writer, aliasMap map[string]string) notary.PassRetriever {
	bound := &boundRetriever{
		in:              in,
		out:             out,
		aliasMap:        aliasMap,
		passphraseCache: make(map[string]string),
	}

	return bound.getPassphrase
}

// ConstantRetriever returns a new Retriever which will return a constant string
// as a passphrase.
func ConstantRetriever(constantPassphrase string) notary.PassRetriever {
	return func(k, a string, c bool, n int) (string, bool, error) {
		return constantPassphrase, false, nil
	}
}
