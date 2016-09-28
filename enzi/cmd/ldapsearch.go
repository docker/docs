package cmd

import (
	"fmt"
	"io/ioutil"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/docker/orca/enzi/api/forms"
	"github.com/docker/orca/enzi/authn/ldap"
	"github.com/docker/orca/enzi/authn/ldap/config"
	goldap "github.com/go-ldap/ldap"
)

// LDAPSearch is the command for performing searches against an LDAP server.
var LDAPSearch = cli.Command{
	Name:      "ldapsearch",
	Usage:     "LDAP Search Tool",
	ArgsUsage: "[filter [attributes...]]",
	Action:    runLDAPSearch,
}

var (
	ldapURI            string
	searchBaseDN       string
	searchScope        string
	bindDN             string
	password           string
	promptBindPassword bool
	startTLS           bool
	pageSize           int
	caFile             string
	insecureSkipVerify bool
)

func init() {
	LDAPSearch.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "host, H",
			Usage:       "URI of LDAP Host",
			Destination: &ldapURI,
		},
		cli.StringFlag{
			Name:        "base-dn, b",
			Usage:       "Search Base DN",
			Destination: &searchBaseDN,
		},
		cli.StringFlag{
			Name:        "scope, s",
			Value:       "one",
			Usage:       "Search Scope {base|one|sub}",
			Destination: &searchScope,
		},
		cli.StringFlag{
			Name:        "bind-dn, D",
			Usage:       "Bind DN (if omitted, no client authentication will be performed)",
			Destination: &bindDN,
		},
		cli.StringFlag{
			Name:        "bind-password, w",
			Usage:       "Password for simple bind authentication (use -W for interactive prompt)",
			Destination: &password,
		},
		cli.BoolFlag{
			Name:        "prompt-bind-password, W",
			Usage:       "Prompt for bind password",
			Destination: &promptBindPassword,
		},
		cli.BoolFlag{
			Name:        "start-tls, Z",
			Usage:       "Use StartTLS",
			Destination: &startTLS,
		},
		cli.IntFlag{
			Name:        "page-limit, p",
			Usage:       "Truncate results size (zero indicates no pagination control - result size will be the maximum allowed by the server)",
			Destination: &pageSize,
		},
		cli.StringFlag{
			Name:        "cafile",
			Usage:       "Path to RootCA certificate bundle to verify TLS connection",
			Destination: &caFile,
		},
		cli.BoolFlag{
			Name:        "insecure",
			Usage:       "Whether to skip server certificate verification",
			Destination: &insecureSkipVerify,
		},
	}
}

func runLDAPSearch(ctx *cli.Context) error {
	validationErr := forms.ValidateLDAPServerURL("host", ldapURI, startTLS)
	if validationErr != nil {
		log.Fatalf("Invalid LDAP URI: %s", validationErr.Detail)
	}

	var caBundle string
	if caFile != "" {
		pemCerts, err := ioutil.ReadFile(caFile)
		if err != nil {
			log.Fatalf("Unable to read CA file: %s", err)
		}
		caBundle = string(pemCerts)
	}

	// We only need to fill out TLS setting here.
	settings := &config.Settings{
		StartTLS:      startTLS,
		TLSSkipVerify: insecureSkipVerify,
		RootCerts:     caBundle,
	}

	conn, err := ldap.GetConn(ldapURI, settings)
	if err != nil {
		log.Fatalf("Unable to connect to LDAP server: %s", err)
	}
	defer conn.Close()

	if bindDN != "" {
		if promptBindPassword {
			prompt := fmt.Sprintf("Password for DN=%s", bindDN)
			password = promptPassword(false, prompt)
		}

		if err := conn.Bind(bindDN, password); err != nil {
			log.Fatalf("Unable to bind as %s: %s", bindDN, err)
		}
	}

	var scope int
	switch searchScope {
	case "base":
		scope = goldap.ScopeBaseObject
	case "one":
		scope = goldap.ScopeSingleLevel
	case "sub":
		scope = goldap.ScopeWholeSubtree
	default:
		log.Fatalf("Invalid search scope: %q", searchScope)
	}

	filter := ctx.Args().First()
	if filter == "" {
		filter = "(objectClass=*)"
	}

	attributes := ctx.Args().Tail()

	searchRequest := goldap.NewSearchRequest(
		searchBaseDN, scope, goldap.DerefAlways,
		0, 0, false, filter, attributes, nil,
	)

	if pageSize > 0 {
		searchRequest.Controls = []goldap.Control{
			goldap.NewControlPaging(uint32(pageSize)),
		}
	}

	result, err := conn.Search(searchRequest)
	if err != nil {
		log.Fatalf("Unable to perform search: %s", err)
	}

	if len(result.Entries) == 0 {
		fmt.Println("No Results")
	} else {
		fmt.Printf("NumResults: %d\nResults:\n", len(result.Entries))
		for _, entry := range result.Entries {
			fmt.Printf("  -\n    DN: %q\n", entry.DN)
			for _, attr := range entry.Attributes {
				for _, val := range attr.Values {
					fmt.Printf("    %s: %q\n", attr.Name, val)
				}
			}
		}
	}

	if len(result.Referrals) > 0 {
		fmt.Printf("NumReferrals: %d\nReferrals:\n", len(result.Referrals))
		for _, referral := range result.Referrals {
			fmt.Printf("  - %s\n", referral)
		}
	}

	return nil
}
