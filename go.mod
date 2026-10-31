module github.com/docker/docs

go 1.27

// This go.mod file is used by hugo to vendor documentation from upstream
// repositories. Use the "require" section to specify the version of the
// upstream repository.
//
// Make sure to add an entry in the "tools" section when adding a new repository.
require (
	github.com/docker/buildx v0.38.1
	github.com/docker/cli v29.9.2+incompatible
	github.com/docker/compose/v5 v5.6.1
	github.com/docker/docker-agent v1.154.0
	github.com/docker/model-runner v1.1.36
	github.com/moby/buildkit v0.34.0
	github.com/moby/moby/api v1.56.1
)

tool (
	github.com/docker/buildx
	github.com/docker/cli
	github.com/docker/compose/v5
	github.com/docker/docker-agent
	github.com/docker/model-runner
	github.com/docker/scout-cli
	github.com/moby/buildkit
	github.com/moby/moby/api
)
