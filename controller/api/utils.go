package api

import (
	"strings"
)

func paramBoolValue(val string) bool {
	v := strings.ToLower(strings.TrimSpace(val))
	return (v == "1" || v == "true")
}
