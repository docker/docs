package auth

import (
	"math"
)

type Role int

const (
	// these must remain in order from least access to highest
	// the getAccess func uses comparison for access level
	None Role = iota
	View
	RestrictedControl
	FullControl
	Admin Role = math.MaxInt64 // used internally for admin only resources
)

type AccessList struct {
	Id     string `json:"id,omitempty"`
	Role   Role   `json:"role"`
	TeamId string `json:"teamId,omitempty"`
	Label  string `json:"label,omitempty"`
}
