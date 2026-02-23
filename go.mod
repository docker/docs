module github.com/docker/docs

go 1.25.0

// This go.mod file is used by hugo to vendor documentation from upstream
// reposities. Use the "require" section to specify the version of the
// upstream repository.
//
// Make sure to add an entry in the "tools" section when adding a new repository.
require (
	github.com/docker/buildx v0.31.1
	github.com/docker/cli v29.2.0+incompatible
	github.com/docker/compose/v5 v5.0.2
	github.com/docker/model-runner/cmd/cli v1.0.3
	github.com/moby/buildkit v0.27.0
	github.com/moby/moby/api v1.53.0
)

tool (
	github.com/docker/buildx
	github.com/docker/cli
	github.com/docker/compose/v5
	github.com/docker/model-runner/cmd/cli
	github.com/docker/scout-cli
	github.com/moby/buildkit
	github.com/moby/moby/api
)
