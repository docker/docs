module github.com/docker/docs

go 1.26.0

// This go.mod file is used by hugo to vendor documentation from upstream
// reposities. Use the "require" section to specify the version of the
// upstream repository.
//
// Make sure to add an entry in the "tools" section when adding a new repository.
require (
	github.com/docker/buildx v0.31.1
	github.com/docker/cli v29.2.1+incompatible
	github.com/docker/compose/v5 v5.0.2
	github.com/docker/model-runner v1.1.9-0.20260303081710-59280ed7abd5
	github.com/moby/buildkit v0.27.0
	github.com/moby/moby/api v1.53.0
)

tool (
	github.com/docker/buildx
	github.com/docker/cli
	github.com/docker/compose/v5
	github.com/docker/model-runner
	github.com/docker/scout-cli
	github.com/moby/buildkit
	github.com/moby/moby/api
)
