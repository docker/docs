package ldap

import (
	"strings"

	"github.com/docker/orca/enzi/schema"
)

// DNSet represents a set of case-insensitive DN strings.
type DNSet map[string]struct{}

// Add adds the given rawDN to the set.
func (s DNSet) Add(rawDN string) {
	// Normalize the DN to be all lowercase.
	s[strings.ToLower(rawDN)] = struct{}{}
}

// Exists checks for existince of the given rawDN.
func (s DNSet) Exists(rawDN string) bool {
	_, exists := s[strings.ToLower(rawDN)]
	return exists
}

// Remove deletes the given rawDN from this set.
func (s DNSet) Remove(rawDN string) {
	delete(s, strings.ToLower(rawDN))
}

// DNAccountMap represents a mapping from case-insensitive DN strings to
// accounts.
type DNAccountMap map[string]*schema.Account

// Put creates a mapping from the given rawDN to the given acct.
func (s DNAccountMap) Put(rawDN string, acct *schema.Account) {
	s[strings.ToLower(rawDN)] = acct
}

// Get retrieves the account for the given rawDN.
func (s DNAccountMap) Get(rawDN string) (acct *schema.Account, exists bool) {
	acct, exists = s[strings.ToLower(rawDN)]
	return acct, exists
}

// Delete deletes the given rawDN from this map.
func (s DNAccountMap) Delete(rawDN string) {
	delete(s, strings.ToLower(rawDN))
}
