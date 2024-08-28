module github.com/docker/docs

go 1.21.0

toolchain go1.22.5

require (
	github.com/docker/buildx v0.16.2 // indirect
	github.com/docker/cli v27.1.1+incompatible // indirect
	github.com/docker/compose/v2 v2.29.2 // indirect
	github.com/docker/scout-cli v1.13.0 // indirect
	github.com/moby/buildkit v0.15.1 // indirect
	github.com/moby/moby v27.2.0+incompatible // indirect
)

replace (
	github.com/docker/buildx => github.com/docker/buildx v0.16.2
	github.com/docker/cli => github.com/docker/cli v27.1.1+incompatible
	github.com/docker/compose/v2 => github.com/docker/compose/v2 v2.29.2
	github.com/docker/scout-cli => github.com/docker/scout-cli v1.13.0
	github.com/moby/buildkit => github.com/moby/buildkit v0.15.1
	github.com/moby/moby => github.com/moby/moby v27.2.0+incompatible
)
