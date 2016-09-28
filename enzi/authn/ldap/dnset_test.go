package ldap

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDNSet(t *testing.T) {
	rawDNs := []string{
		// These 5 should normalize to the same value for Jane Doe.
		"uid=jane.doe,ou=people,dc=example,dc=com",
		"UID=jane.doe,OU=people,DC=example,DC=com",
		"uid=Jane.Doe,ou=People,dc=Example,dc=Com",
		"UID=Jane.Doe,OU=People,DC=Example,DC=Com",
		"UID=JANE.DOE,OU=PEOPLE,DC=EXAMPLE,DC=COM",
		// These 5 should normalize to the same value for John Doe.
		"uid=john.doe,ou=people,dc=example,dc=com",
		"UID=john.doe,OU=people,DC=example,DC=com",
		"uid=John.Doe,ou=People,dc=Example,dc=Com",
		"UID=John.Doe,OU=People,DC=Example,DC=Com",
		"UID=JOHN.DOE,OU=PEOPLE,DC=EXAMPLE,DC=COM",
	}

	dnSet := DNSet{}

	// Test adding values.
	for _, rawDN := range rawDNs {
		dnSet.Add(rawDN)
	}
	require.Len(t, dnSet, 2)

	// Test getting values from the set.
	for _, rawDN := range rawDNs {
		require.True(t, dnSet.Exists(rawDN))
	}

	// Test deleting values.
	for _, rawDN := range rawDNs {
		dnSet.Remove(rawDN)
	}
	require.Empty(t, dnSet)
}
