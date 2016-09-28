package dtrutil

import (
	"github.com/docker/distribution"
	"github.com/docker/distribution/configuration"
	"github.com/docker/distribution/context"
	"github.com/docker/distribution/registry/storage"
	"github.com/docker/distribution/registry/storage/driver/factory"
	"github.com/docker/libtrust"

	// register all storage and auth drivers
	_ "github.com/docker/distribution/registry/auth/htpasswd"
	_ "github.com/docker/distribution/registry/auth/silly"
	_ "github.com/docker/distribution/registry/auth/token"
	_ "github.com/docker/distribution/registry/proxy"
	_ "github.com/docker/distribution/registry/storage/driver/azure"
	_ "github.com/docker/distribution/registry/storage/driver/filesystem"
	_ "github.com/docker/distribution/registry/storage/driver/gcs"
	_ "github.com/docker/distribution/registry/storage/driver/inmemory"
	_ "github.com/docker/distribution/registry/storage/driver/middleware/cloudfront"
	_ "github.com/docker/distribution/registry/storage/driver/oss"
	_ "github.com/docker/distribution/registry/storage/driver/s3-aws"
	_ "github.com/docker/distribution/registry/storage/driver/swift"

	"github.com/palantir/stacktrace"
)

// NewRegistry constructs a new distribution.Namespace based off of the stored
// distribution config.
func NewRegistry(ctx context.Context, config *configuration.Configuration) (distribution.Namespace, error) {
	driver, err := factory.Create(config.Storage.Type(), config.Storage.Parameters())
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to construct %s driver", config.Storage.Type())
	}

	k, err := libtrust.GenerateECP256PrivateKey()
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to generate schema1 signature key")
	}

	registry, err := storage.NewRegistry(ctx, driver, storage.DisableSchema1Signatures, storage.Schema1SigningKey(k))
	if err != nil {
		return nil, stacktrace.Propagate(err, "failed to construct registry")
	}

	return registry, nil
}
