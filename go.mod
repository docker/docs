module github.com/docker/docs

go 1.24.0

require (
	github.com/docker/buildx v0.24.0 // indirect
	github.com/docker/cli v28.2.1+incompatible // indirect
	github.com/docker/compose/v2 v2.36.2 // indirect
	github.com/docker/model-cli v0.1.25 // indirect
	github.com/docker/scout-cli v1.15.0 // indirect
	github.com/moby/buildkit v0.22.0 // indirect
	github.com/moby/moby v28.2.1+incompatible // indirect
)

replace (
	github.com/docker/buildx => github.com/docker/buildx v0.24.0
	github.com/docker/cli => github.com/docker/cli v28.2.1+incompatible
	github.com/docker/compose/v2 => github.com/docker/compose/v2 v2.36.2
	github.com/docker/model-cli => github.com/docker/model-cli v0.1.26-0.20250527144806-15d0078a3c01
	github.com/docker/scout-cli => github.com/docker/scout-cli v1.15.0
	github.com/moby/buildkit => github.com/moby/buildkit v0.22.0
	github.com/moby/moby => github.com/moby/moby v28.2.1+incompatible
)
