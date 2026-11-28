module github.com/docker/docs

go 1.27.4

// This go.mod file is used by hugo to vendor documentation from upstream
// repositories. Use the "require" section to specify the version of the
// upstream repository.
//
// Make sure to add an entry in the "tools" section when adding a new repository.
require (
	github.com/docker/buildx v0.39.1
	github.com/docker/cli v29.10.0+incompatible
	github.com/docker/compose/v5 v5.8.0
	github.com/docker/docker-agent v1.167.0
	github.com/docker/model-runner v1.1.36
	github.com/moby/buildkit v0.35.1
	github.com/moby/moby/api v1.57.0
	github.com/moby/moby/v2 v2.0.0-beta.31
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
	github.com/moby/moby/v2
)

// Use the upstream glossary-link fix until a tagged release includes it.
replace github.com/moby/buildkit => github.com/moby/buildkit v0.35.0-rc1.0.20261118113413-bfe7dac5cde4
