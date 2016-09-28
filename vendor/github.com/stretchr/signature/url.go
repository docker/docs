package signature

import (
	stewstrings "github.com/stretchr/stew/strings"
	"net/url"
	"sort"
)

// OrderParams gets the parameters ordered by key, then by values as a URL string.
func OrderParams(values url.Values) string {

	// get the keys
	var keys []string
	for k, _ := range values {
		keys = append(keys, k)
	}

	// sort the keys
	sort.Strings(keys)

	// ordered items
	var ordered []string

	// sort the values
	for _, key := range keys {
		sort.Strings(values[key])
		for _, val := range values[key] {
			ordered = append(ordered, stewstrings.MergeStrings(key, "=", val))
		}
	}

	joined := stewstrings.JoinStrings("&", ordered...)
	return joined

}
