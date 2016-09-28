package ha_utils

import (
	"os"
)

func GetAdminPassword() string {
	password := os.Getenv("UCP_ADMIN_PASSWORD")
	if password == "" {
		password = "orca"
	}
	return password
}

func GetAdminUser() string {
	username := os.Getenv("UCP_ADMIN")
	if username == "" {
		username = "admin"
	}
	return username
}
