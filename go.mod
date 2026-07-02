module github.com/docker/docs

go 1.26.4

// This go.mod file is used by hugo to vendor documentation from upstream
// repositories. Use the "require" section to specify the version of the
// upstream repository.
//
// Make sure to add an entry in the "tools" section when adding a new repository.
require (
	github.com/docker/buildx v0.35.0
	github.com/docker/cli v29.6.1+incompatible
	github.com/docker/compose/v5 v5.3.0
	github.com/docker/docker-agent v1.96.0
	github.com/docker/model-runner v1.1.36
	github.com/moby/buildkit v0.31.0
	github.com/moby/moby/api v1.55.0
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
