package util

import (
	"fmt"
	"path"
	"strings"

	"github.com/docker/dhe-deploy"
	"github.com/docker/dhe-deploy/hubconfig"

	"github.com/docker/distribution/configuration"
	storagedriver "github.com/docker/distribution/registry/storage/driver"
	// import storage drivers
	_ "github.com/docker/distribution/registry/storage/driver/azure"
	"github.com/docker/distribution/registry/storage/driver/factory"
	_ "github.com/docker/distribution/registry/storage/driver/filesystem"
	_ "github.com/docker/distribution/registry/storage/driver/gcs"
	_ "github.com/docker/distribution/registry/storage/driver/inmemory"
	_ "github.com/docker/distribution/registry/storage/driver/middleware/cloudfront"
	_ "github.com/docker/distribution/registry/storage/driver/oss"
	_ "github.com/docker/distribution/registry/storage/driver/s3-aws"
	_ "github.com/docker/distribution/registry/storage/driver/swift"
)

func GetReadonlyMode(storage *configuration.Storage) bool {
	storageReadOnly := false
	if storageMaintenance, ok := (*storage)["maintenance"]; ok {
		if readOnly, ok := storageMaintenance["readonly"]; ok {
			// Depending on how we parsed the yaml, it might give us a map of interface
			// to interface. The YAML parser sucks. I don't know. In any case, this hack
			// is easier than trying to figure out why we can't parse the yaml into the
			// right types.
			if readOnlyMap, ok := readOnly.(map[string]interface{}); ok {
				if readOnlyEnabled, ok := readOnlyMap["enabled"]; ok {
					// will default to false if not a boolean
					storageReadOnly, _ = readOnlyEnabled.(bool)
				}
			} else if readOnlyMap, ok := readOnly.(map[interface{}]interface{}); ok {
				if readOnlyEnabled, ok := readOnlyMap["enabled"]; ok {
					// will default to false if not a boolean
					storageReadOnly, _ = readOnlyEnabled.(bool)
				}
			}

		}
	}
	return storageReadOnly
}

func SetReadonlyMode(storage *configuration.Storage, readonly bool) {
	maintenanceMap := (*storage)["maintenance"]
	if maintenanceMap == nil {
		maintenanceMap = configuration.Parameters{}
		(*storage)["maintenance"] = maintenanceMap
	}
	// XXX: In yaml, this needs to be a map of interface to interface. Otherwise the registry will freak out.
	maintenanceMap["readonly"] = map[interface{}]interface{}{"enabled": readonly}
}

func SetReadonlyModeJSON(storage *configuration.Storage, readonly bool) {
	maintenanceMap := (*storage)["maintenance"]
	if maintenanceMap == nil {
		maintenanceMap = configuration.Parameters{}
		(*storage)["maintenance"] = maintenanceMap
	}
	// XXX: In yaml, this needs to be a map of interface to interface. Otherwise the registry will freak out.
	maintenanceMap["readonly"] = map[string]interface{}{"enabled": readonly}
}

func RemoveUseragent(storage *configuration.Storage) {
	parameters := (*storage)[storage.Type()]
	delete(parameters, "useragent")
	(*storage)[storage.Type()] = parameters
}

func DriverFromStorage(storage configuration.Storage) (storagedriver.StorageDriver, error) {
	driverType := storage.Type()
	driverParams := storage.Parameters()
	// when using the filesystem driver the rootdirectory is hard-coded
	// to /storage.
	if driverType == "filesystem" {
		driverParams["rootdirectory"] = "/storage"
	}

	return factory.Create(driverType, driverParams)
}

func GetRegistryAuthConfig(domainName string) hubconfig.JSONGarantOptions {
	var realm string
	parts := strings.Split(domainName, ":")
	domainPart := parts[0]
	portPart := ""
	if len(parts) > 1 {
		portPart = parts[1]
	}

	// we omit the port if it's 443
	if portPart == "443" {
		realm = fmt.Sprintf("https://%s/%s/token", domainPart, deploy.GarantSubroute)
		domainName = domainPart
	} else {
		realm = fmt.Sprintf("https://%s/%s/token", domainName, deploy.GarantSubroute)
	}

	return hubconfig.JSONGarantOptions{
		Realm:      realm,
		Issuer:     domainName,
		Service:    domainName,
		CertBundle: path.Join(deploy.ConfigDirPath, deploy.GarantRootCertFilename),
	}
}
