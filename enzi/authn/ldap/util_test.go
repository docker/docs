package ldap

import (
	"fmt"
	"testing"

	"github.com/go-ldap/ldap"
	"github.com/stretchr/testify/require"
)

func TestGetRangedAttrValues(t *testing.T) {
	// We have 15 total members. We will test the getRangedAttrValues
	// using varying slices of this full set.
	memberDNset := []string{
		"uid=jane.doe,ou=people,dc=example,dc=com",
		"uid=john.doe,ou=people,dc=example,dc=com",
		"uid=archie.takis,ou=people,dc=example,dc=com",
		"uid=filipe,varuna,ou=people,dc=example,dc=com",
		"uid=mick.folke,ou=people,dc=example,dc=com",
		"uid=ira.eduard,ou=people,dc=example,dc=com",
		"uid=sehzade.garvan,ou=people,dc=example,dc=com",
		"uid=gereon.ethelred,ou=people,dc=example,dc=com",
		"uid=rodge.auxentius,ou=people,dc=example,dc=com",
		"uid=mario.amuary,ou=people,dc=example,dc=com",
		"uid=phoebus.andon,ou=people,dc=example,dc=com",
		"uid=niko.trafford,ou=people,dc=example,dc=com",
		"uid=edom.stef,ou=people,dc=example,dc=com",
		"uid=ilia.daniel,ou=people,dc=example,dc=com",
		"uid=temuri.vasya,ou=people,dc=example,dc=com",
	}
	groupDN := "cn=devs,ou=groups,dc=example,dc=com"

	// Each test case covers one or more attribute queries with their
	// corresponding results.
	testCases := [][]struct {
		attr   string
		result *ldap.Entry
	}{
		// The first test case has 1 query with all members in the
		// first result.
		{
			{
				attr: "member",
				result: &ldap.Entry{
					DN: groupDN,
					Attributes: []*ldap.EntryAttribute{
						{
							Name:   "member",
							Values: memberDNset[:],
						},
					},
				},
			},
		},
		// The second test case has 2 queries with 10 members each.
		{
			{
				attr: "member",
				result: &ldap.Entry{
					DN: groupDN,
					Attributes: []*ldap.EntryAttribute{
						{
							Name:   "member",
							Values: nil,
						},
						{
							Name:   "member;range=0-9",
							Values: memberDNset[:10],
						},
					},
				},
			},
			{
				attr: "member;range=10-*",
				result: &ldap.Entry{
					DN: groupDN,
					Attributes: []*ldap.EntryAttribute{
						{
							Name:   "member;range=10-*",
							Values: memberDNset[10:],
						},
					},
				},
			},
		},
		// The third test case has 3 queries with 5 members each.
		{
			{
				attr: "member",
				result: &ldap.Entry{
					DN: groupDN,
					Attributes: []*ldap.EntryAttribute{
						{
							Name:   "member",
							Values: nil,
						},
						{
							Name:   "member;range=0-4",
							Values: memberDNset[:5],
						},
					},
				},
			},
			{
				attr: "member;range=5-*",
				result: &ldap.Entry{
					DN: groupDN,
					Attributes: []*ldap.EntryAttribute{
						{
							Name:   "member;range=5-*",
							Values: nil,
						},
						{
							Name:   "member;range=5-9",
							Values: memberDNset[5:10],
						},
					},
				},
			},
			{
				attr: "member;range=10-*",
				result: &ldap.Entry{
					DN: groupDN,
					Attributes: []*ldap.EntryAttribute{
						{
							Name:   "member;range=10-*",
							Values: memberDNset[10:],
						},
					},
				},
			},
		},
	}

	for i, testCase := range testCases {
		var allVals []string
		expectedQueryAttr := "member"

		for j, query := range testCase {
			require.Equal(t, expectedQueryAttr, query.attr, "test case %d, query attr %d", i, j)

			vals, rangeEnd := getRangedAttrValues(query.result, "member")
			allVals = append(allVals, vals...)

			if rangeEnd == 0 {
				break
			}

			expectedQueryAttr = fmt.Sprintf("member;range=%d-*", rangeEnd+1)
		}

		// All of the values returned should be equal to our original
		// member DN set.
		require.Equal(t, len(memberDNset), len(allVals), "test case %d", i)
		for j := 0; j < len(memberDNset); j++ {
			require.Equal(t, memberDNset[j], allVals[j])
		}

		t.Logf("test case %d passed!", i)
	}
}
