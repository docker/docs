package registry

import (
	"encoding/base64"
	"encoding/json"
	"os"

	"github.com/docker/engine-api/types"
)

func RequestPrivilegeFunc() (string, error) {
	auth := types.AuthConfig{
		Username: os.Getenv("REGISTRY_USERNAME"),
		Password: os.Getenv("REGISTRY_PASSWORD"),
		Email:    os.Getenv("REGISTRY_EMAIL"),
	}
	buf, err := json.Marshal(auth)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(buf), nil
}
