package ldap

import (
	"bytes"
	"container/list"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/orca/auth"
	"github.com/go-ldap/ldap"
)

// ErrGroupNotFound indicates that an LDAP group could not be found.
var ErrGroupNotFound = errors.New("group not found")

// AndSearchFilter handles organizing the given arguments into a combined
// logical-AND LDAP search filter, handling wrapping arguments in parens if
// they are not already.
func AndSearchFilter(args ...string) string {
	if len(args) == 0 {
		return ""
	}

	for i, arg := range args {
		if !strings.HasPrefix(arg, "(") {
			args[i] = fmt.Sprintf("(%s)", arg)
		}
	}

	if len(args) == 1 {
		return args[0]
	}

	return fmt.Sprintf("(&%s)", strings.Join(args, ""))
}

type LDAPSearchURL struct {
	HostURLString string
	BaseDN        string
	// TODO: add scope and extra filters to LDAPSearchURL
}

type SearchResult struct {
	Entry         *ldap.Entry
	HostURLString string
}

// SearchOpts are the options used in an LDAP search for a user or group.
type SearchOpts struct {
	RootURL    LDAPSearchURL
	Scope      int
	Filter     string
	Attributes []string
}

func SearchWithReferrals(settings auth.LDAPSettings, searchOpts SearchOpts) *SearchResult {
	var (
		entry          *ldap.Entry
		currentLdapURL LDAPSearchURL
		ldapURLList    = list.New()
		visitedSet     = make(map[LDAPSearchURL]struct{})
	)

	listNode := ldapURLList.PushFront(searchOpts.RootURL)
	for listNode != nil && entry == nil {
		// use an anonymous function so that we can use defer to ensure that
		// connections are released after each iteration
		// TODO maybe build a list of errors to return?
		func() {
			currentLdapURL = ldapURLList.Remove(listNode).(LDAPSearchURL)
			if _, visited := visitedSet[currentLdapURL]; visited {
				return
			}
			visitedSet[currentLdapURL] = struct{}{}

			searchReq := ldap.NewSearchRequest(currentLdapURL.BaseDN, searchOpts.Scope, ldap.DerefAlways, 1, 5, false, searchOpts.Filter, searchOpts.Attributes, nil)
			if currentLdapURL == searchOpts.RootURL {
				log.Info("beginning search")
			} else {
				log.Info("following referral")
			}

			ldapConn, err := GetConn(currentLdapURL.HostURLString, settings)
			if err != nil {
				log.Errorf("failed to connect to LDAP host: %s", err)
				return
			}
			defer ldapConn.Close()

			if err := ldapConn.Bind(settings.ReaderDN, settings.ReaderPassword); err != nil {
				log.Errorf("error binding reader: %s", err)
				return
			}

			searchRes, err := ldapConn.Search(searchReq)
			if err != nil {
				log.Errorf("error making search: %q: %s", searchReq.Filter, err)
				return
			} else if len(searchRes.Entries) > 0 {
				log.Info("found entry")

				entry = searchRes.Entries[0]
				return
			} else {
				log.Info("no entry results but got some referrals")

				for _, referralString := range searchRes.Referrals {
					newReferral, err := ParseReferral(referralString)
					if err != nil {
						log.Errorf("failed to parse referral: %s", err)
						continue
					}

					ldapURLList.PushBack(newReferral)
				}
				return
			}
		}()

		listNode = ldapURLList.Front()
	}

	if entry == nil {
		return nil
	}

	return &SearchResult{
		Entry:         entry,
		HostURLString: currentLdapURL.HostURLString,
	}
}

// ParseReferral parses a given LDAP URL into its host and path components.
// These typcially look like: ldap://ldap.example.com/ou=us,dc=example,dc=com
// and are used for partitioning data in the directory tree by country, for
// example.
func ParseReferral(referral string) (LDAPSearchURL, error) {
	referralURL, err := url.Parse(referral)
	if err != nil {
		return LDAPSearchURL{}, err
	}
	referralDN := ""
	if len(referralURL.Path) > 0 {
		referralDN = referralURL.Path[1:]
	}
	return LDAPSearchURL{
		HostURLString: referral,
		BaseDN:        referralDN,
	}, nil
}

func printEntries(entries []*ldap.Entry) string {
	entriesStrBuf := new(bytes.Buffer)
	fmt.Fprint(entriesStrBuf, "[")
	for i, entry := range entries {
		if i > 0 {
			fmt.Fprint(entriesStrBuf, " ")
		}
		fmt.Fprintf(entriesStrBuf, `{DN: %s, attributes: [`, entry.DN)
		for j, attr := range entry.Attributes {
			if j > 0 {
				fmt.Fprint(entriesStrBuf, " ")
			}
			fmt.Fprintf(entriesStrBuf, "%s: %s", attr.Name, attr.Values)
		}
		fmt.Fprint(entriesStrBuf, "]}")
	}
	fmt.Fprint(entriesStrBuf, "]")
	return entriesStrBuf.String()
}

// GetGroupMembers iterates through the members of a group and returns a set of
// all member DNs.
func GetGroupMembers(settings auth.LDAPSettings, groupDN, memberAttrName string) (memberDNSet map[string]struct{}) {
	memberDNSet = make(map[string]struct{})

	memberAttr := memberAttrName
	searchOpts := SearchOpts{
		RootURL: LDAPSearchURL{
			HostURLString: settings.ServerURL,
			BaseDN:        groupDN,
		},
		Scope:  ldap.ScopeBaseObject,
		Filter: "(objectClass=*)",
	}

	for {
		log.Infof("Searching group: %s for attr: %s", groupDN, memberAttr)

		searchOpts.Attributes = []string{memberAttr}

		result := SearchWithReferrals(settings, searchOpts)
		if result == nil {
			return memberDNSet
		}

		vals, rangeEnd := getRangedAttrValues(result.Entry, memberAttrName)

		for _, memberDN := range vals {
			memberDNSet[memberDN] = struct{}{}
		}

		if rangeEnd <= 0 {
			break
		}

		// update the queryable attribute, continue search
		memberAttr = fmt.Sprintf("%s;range=%d-*", memberAttrName, rangeEnd+1)
	}

	return memberDNSet
}

var attrRangeOptionPatten = regexp.MustCompile(`;range=[\d*]+-([\d*]+);?`)

// getRangedAttrValues searches the given entry for an attribute with the given
// attrName and returns the corresponding values. If the attribute contains a
// range option, the end value of that range will be returned. A rangeEnd value
// of 0 indicates the end of the range or that there was no range option.
func getRangedAttrValues(entry *ldap.Entry, attrName string) (vals []string, rangeEnd int) {
	for _, attr := range entry.Attributes {
		if attr.Name == attrName {
			return attr.Values, 0
		}

		// A semicolon after the attribute name indic
		if !strings.HasPrefix(attr.Name, attrName+";") {
			continue // Is not the attribute with options.
		}

		// Trim off the attribute name to get only the options string.
		attrOpts := attr.Name[len(attrName):]

		// Search the attribute options for a "range=start-end" option.
		// Note: this is a special Microsoft Active Directory option.
		matches := attrRangeOptionPatten.FindStringSubmatch(attrOpts)
		if len(matches) == 2 {
			rangeEnd, _ = strconv.Atoi(matches[1])
		}

		return attr.Values, rangeEnd
	}

	return nil, 0
}
