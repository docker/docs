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

	"github.com/docker/distribution/context"
	"github.com/docker/orca/enzi/authn/ldap/config"
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

func SearchWithReferrals(ctx context.Context, settings *config.Settings, searchOpts SearchOpts) *SearchResult {
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
		func(ctx context.Context) {
			currentLdapURL = ldapURLList.Remove(listNode).(LDAPSearchURL)
			if _, visited := visitedSet[currentLdapURL]; visited {
				return
			}
			visitedSet[currentLdapURL] = struct{}{}

			ctx = context.WithValue(ctx, "LDAPSearchURL", currentLdapURL)
			ctx = context.WithLogger(ctx, context.GetLogger(ctx, "LDAPSearchURL"))

			searchReq := ldap.NewSearchRequest(currentLdapURL.BaseDN, searchOpts.Scope, ldap.DerefAlways, 1, 5, false, searchOpts.Filter, searchOpts.Attributes, nil)
			if currentLdapURL == searchOpts.RootURL {
				context.GetLogger(ctx).Debug("beginning search")
			} else {
				context.GetLogger(ctx).Debug("following referral")
			}

			ldapConn, err := GetConn(currentLdapURL.HostURLString, settings)
			if err != nil {
				ctx = context.WithValues(ctx, map[string]interface{}{
					"error":    err.Error(),
					"hostURL":  currentLdapURL.HostURLString,
					"startTLS": fmt.Sprint(settings.StartTLS),
				})
				context.GetLogger(ctx, "error", "hostURL", "startTLS").Error("failed to connect to LDAP host")
				return
			}
			defer ldapConn.Close()

			if err := ldapConn.Bind(settings.ReaderDN, settings.ReaderPassword); err != nil {
				context.GetLogger(ctx).Errorf("error binding reader: %s", err)
				return
			}

			searchRes, err := ldapConn.Search(searchReq)
			if err != nil {
				context.GetLogger(ctx).Errorf("error making search: %q: %s", searchReq.Filter, err)
				return
			} else if len(searchRes.Entries) > 0 {
				ctx = context.WithValue(ctx, "searchEntries", printEntries(searchRes.Entries))
				context.GetLogger(ctx, "searchEntries").Debug("found entry")

				entry = searchRes.Entries[0]
				return
			} else {
				ctx = context.WithValue(ctx, "referrals", searchRes.Referrals)
				context.GetLogger(ctx, "referrals").Debug("no entry results but got some referrals")

				for _, referralString := range searchRes.Referrals {
					newReferral, err := ParseReferral(referralString)
					if err != nil {
						ctx = context.WithValues(ctx, map[string]interface{}{
							"error":          err.Error(),
							"referralString": referralString,
						})
						context.GetLogger(ctx, "error", "referralString").Error("failed to parse referral")
						continue
					}

					ldapURLList.PushBack(newReferral)
				}
				return
			}
		}(ctx)

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
func GetGroupMembers(ctx context.Context, settings *config.Settings, groupDN, memberAttrName string) (memberDNSet DNSet) {
	memberDNSet = make(DNSet)

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
		context.GetLogger(ctx).Infof("Searching group: %s for attr: %s", groupDN, memberAttr)

		searchOpts.Attributes = []string{memberAttr}

		result := SearchWithReferrals(ctx, settings, searchOpts)
		if result == nil {
			return memberDNSet
		}

		vals, rangeEnd := getRangedAttrValues(result.Entry, memberAttrName)

		for _, memberDN := range vals {
			memberDNSet.Add(memberDN)
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
		if len(attr.Values) == 0 {
			// Skip an attribute with no values. This can occur
			// when a query for an attribute returns both the
			// originally queried attribute (with no values) and
			// the original queried attribute with a range option
			// containing the next slice of the result values. In
			// this case, we will skip it and continue to next
			// attribute which does contain values.
			continue
		}

		if attr.Name == attrName {
			// If we have an exact attribute match (with values),
			// simply return those values.
			return attr.Values, 0
		}

		// A semicolon after the attribute name indicates options.
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
