module github.com/docker/docs

go 1.24.9

require (
	github.com/docker/buildx v0.30.1 // indirect
	github.com/docker/cli v29.1.2+incompatible // indirect; see "replace" rule at the bottom for actual version
	github.com/docker/compose/v2 v2.40.3 // indirect
	github.com/docker/mcp-gateway v0.22.0 // indirect
	github.com/docker/model-runner/cmd/cli v0.1.44 // indirect
	github.com/docker/scout-cli v1.18.4 // indirect
	github.com/moby/buildkit v0.26.1 // indirect
	github.com/moby/moby/api v1.52.0 // indirect; see "replace" rule at the bottom for actual version
)

replace (
	github.com/docker/buildx => github.com/docker/buildx v0.30.1
	github.com/docker/cli => github.com/docker/cli v29.1.2+incompatible
	github.com/docker/compose/v2 => github.com/docker/compose/v2 v2.40.3
	github.com/docker/mcp-gateway => github.com/docker/mcp-gateway v0.22.0
	github.com/docker/model-runner/cmd/cli => github.com/docker/model-runner/cmd/cli v0.1.44
	github.com/docker/scout-cli => github.com/docker/scout-cli v1.18.4
	github.com/moby/buildkit => github.com/moby/buildkit v0.26.0
	github.com/moby/moby/api => github.com/moby/moby/api v1.52.0
)
