module github.com/docker/docs

go 1.24.0

require (
	github.com/docker/buildx v0.25.0 // indirect
	github.com/docker/cli v28.3.2+incompatible // indirect
	github.com/docker/compose/v2 v2.38.2 // indirect
	github.com/docker/model-cli v0.1.33-0.20250703103301-d4e4936a9eb2 // indirect
	github.com/docker/scout-cli v1.15.0 // indirect
	github.com/moby/buildkit v0.23.2 // indirect
	github.com/moby/moby v28.3.2+incompatible // indirect
)

replace (
	github.com/docker/buildx => github.com/docker/buildx v0.25.0
	github.com/docker/cli => github.com/docker/cli v28.3.2+incompatible
	github.com/docker/compose/v2 => github.com/docker/compose/v2 v2.38.2
	github.com/docker/model-cli => github.com/docker/model-cli v0.1.33-0.20250703103301-d4e4936a9eb2
	github.com/docker/scout-cli => github.com/docker/scout-cli v1.15.0
	github.com/moby/buildkit => github.com/moby/buildkit v0.23.2
	github.com/moby/moby => github.com/moby/moby v28.3.2+incompatible
)
