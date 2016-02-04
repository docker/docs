package main

import (
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/docker/notary/client"
	"github.com/docker/notary/trustmanager"
	"github.com/docker/notary/tuf/data"
	"github.com/olekukonko/tablewriter"
)

// returns a tablewriter
func getTable(headers []string, writer io.Writer) *tablewriter.Table {
	table := tablewriter.NewWriter(writer)
	table.SetBorder(false)
	table.SetColumnSeparator(" ")
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("-")
	table.SetAutoWrapText(false)
	table.SetHeader(headers)
	return table
}

// --- pretty printing certs ---

func truncateWithEllipsis(str string, maxWidth int, leftTruncate bool) string {
	if len(str) <= maxWidth {
		return str
	}
	if leftTruncate {
		return fmt.Sprintf("...%s", str[len(str)-(maxWidth-3):])
	}
	return fmt.Sprintf("%s...", str[:maxWidth-3])
}

const (
	maxGUNWidth = 25
	maxLocWidth = 40
)

type keyInfo struct {
	gun      string // assumption that this is "" if role is root
	role     string
	keyID    string
	location string
}

// We want to sort by gun, then by role, then by keyID, then by location
// In the case of a root role, then there is no GUN, and a root role comes
// first.
type keyInfoSorter []keyInfo

func (k keyInfoSorter) Len() int      { return len(k) }
func (k keyInfoSorter) Swap(i, j int) { k[i], k[j] = k[j], k[i] }
func (k keyInfoSorter) Less(i, j int) bool {
	// special-case role
	if k[i].role != k[j].role {
		if k[i].role == data.CanonicalRootRole {
			return true
		}
		if k[j].role == data.CanonicalRootRole {
			return false
		}
		// otherwise, neither of them are root, they're just different, so
		// go with the traditional sort order.
	}

	// sort order is GUN, role, keyID, location.
	orderedI := []string{k[i].gun, k[i].role, k[i].keyID, k[i].location}
	orderedJ := []string{k[j].gun, k[j].role, k[j].keyID, k[j].location}

	for x := 0; x < 4; x++ {
		switch {
		case orderedI[x] < orderedJ[x]:
			return true
		case orderedI[x] > orderedJ[x]:
			return false
		}
		// continue on and evalulate the next item
	}
	// this shouldn't happen - that means two values are exactly equal
	return false
}

// Given a list of KeyStores in order of listing preference, pretty-prints the
// root keys and then the signing keys.
func prettyPrintKeys(keyStores []trustmanager.KeyStore, writer io.Writer) {
	var info []keyInfo

	for _, store := range keyStores {
		for keyPath, role := range store.ListKeys() {
			gun := ""
			if role != data.CanonicalRootRole {
				dirPath := filepath.Dir(keyPath)
				if dirPath != "." { // no gun
					gun = dirPath
				}
			}
			info = append(info, keyInfo{
				role:     role,
				location: store.Name(),
				gun:      gun,
				keyID:    filepath.Base(keyPath),
			})
		}
	}

	if len(info) == 0 {
		writer.Write([]byte("No signing keys found.\n"))
		return
	}

	sort.Stable(keyInfoSorter(info))

	table := getTable([]string{"ROLE", "GUN", "KEY ID", "LOCATION"}, writer)

	for _, oneKeyInfo := range info {
		table.Append([]string{
			oneKeyInfo.role,
			truncateWithEllipsis(oneKeyInfo.gun, maxGUNWidth, true),
			oneKeyInfo.keyID,
			truncateWithEllipsis(oneKeyInfo.location, maxLocWidth, true),
		})
	}
	table.Render()
}

// --- pretty printing targets ---

type targetsSorter []*client.TargetWithRole

func (t targetsSorter) Len() int      { return len(t) }
func (t targetsSorter) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
func (t targetsSorter) Less(i, j int) bool {
	return t[i].Name < t[j].Name
}

// --- pretty printing roles ---

type roleSorter []*data.Role

func (r roleSorter) Len() int      { return len(r) }
func (r roleSorter) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r roleSorter) Less(i, j int) bool {
	return r[i].Name < r[j].Name
}

// Pretty-prints the sorted list of TargetWithRoles.
func prettyPrintTargets(ts []*client.TargetWithRole, writer io.Writer) {
	if len(ts) == 0 {
		writer.Write([]byte("\nNo targets present in this repository.\n\n"))
		return
	}

	sort.Stable(targetsSorter(ts))

	table := getTable([]string{"Name", "Digest", "Size (bytes)", "Role"}, writer)

	for _, t := range ts {
		table.Append([]string{
			t.Name,
			hex.EncodeToString(t.Hashes["sha256"]),
			fmt.Sprintf("%d", t.Length),
			t.Role,
		})
	}
	table.Render()
}

// Pretty-prints the list of provided Roles
func prettyPrintRoles(rs []*data.Role, writer io.Writer, roleType string) {
	if len(rs) == 0 {
		writer.Write([]byte(fmt.Sprintf("\nNo %s present in this repository.\n\n", roleType)))
		return
	}

	// this sorter works for Role types
	sort.Stable(roleSorter(rs))

	table := getTable([]string{"Role", "Paths", "Key IDs", "Threshold"}, writer)

	for _, r := range rs {
		table.Append([]string{
			r.Name,
			prettyPrintPaths(r.Paths),
			strings.Join(r.KeyIDs, ","),
			fmt.Sprintf("%v", r.Threshold),
		})
	}
	table.Render()
}

// Pretty-prints a list of delegation paths, and ensures the empty string is printed as "" in the console
func prettyPrintPaths(paths []string) string {
	// sort paths first
	sort.Strings(paths)
	prettyPaths := []string{}
	for _, path := range paths {
		// manually escape "" and designate that it is all paths with an extra print <all paths>
		if path == "" {
			path = "\"\" <all paths>"
		}
		prettyPaths = append(prettyPaths, path)
	}
	return strings.Join(prettyPaths, ",")
}

// --- pretty printing certs ---

// cert by repo name then expiry time.  Don't bother sorting by fingerprint.
type certSorter []*x509.Certificate

func (t certSorter) Len() int      { return len(t) }
func (t certSorter) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
func (t certSorter) Less(i, j int) bool {
	if t[i].Subject.CommonName < t[j].Subject.CommonName {
		return true
	} else if t[i].Subject.CommonName > t[j].Subject.CommonName {
		return false
	}

	return t[i].NotAfter.Before(t[j].NotAfter)
}

// Given a list of Ceritifcates in order of listing preference, pretty-prints
// the cert common name, fingerprint, and expiry
func prettyPrintCerts(certs []*x509.Certificate, writer io.Writer) {
	if len(certs) == 0 {
		writer.Write([]byte("\nNo trusted root certificates present.\n\n"))
		return
	}

	sort.Stable(certSorter(certs))

	table := getTable([]string{
		"GUN", "Fingerprint of Trusted Root Certificate", "Expires In"}, writer)

	for _, c := range certs {
		days := math.Floor(c.NotAfter.Sub(time.Now()).Hours() / 24)
		expiryString := "< 1 day"
		if days == 1 {
			expiryString = "1 day"
		} else if days > 1 {
			expiryString = fmt.Sprintf("%d days", int(days))
		}

		certID, err := trustmanager.FingerprintCert(c)
		if err != nil {
			fatalf("Could not fingerprint certificate: %v", err)
		}

		table.Append([]string{c.Subject.CommonName, certID, expiryString})
	}
	table.Render()
}
