module github.com/docker/docs

go 1.21.0

toolchain go1.22.5

require (
	github.com/docker/buildx v0.17.0 // indirect
	github.com/docker/cli v27.2.2-0.20240909090509-65decb573126+incompatible // indirect
	github.com/docker/compose/v2 v2.29.6 // indirect
	github.com/docker/scout-cli v1.13.0 // indirect
	github.com/moby/buildkit v0.16.0 // indirect
	github.com/moby/moby v27.2.1+incompatible // indirect
)

replace (
	github.com/docker/buildx => github.com/docker/buildx v0.17.0
	github.com/docker/cli => github.com/docker/cli v27.2.2-0.20240909090509-65decb573126+incompatible
	github.com/docker/compose/v2 => github.com/docker/compose/v2 v2.29.2
	github.com/docker/scout-cli => github.com/docker/scout-cli v1.13.0
	github.com/moby/buildkit => github.com/moby/buildkit v0.16.0
	github.com/moby/moby => github.com/moby/moby v27.2.1+incompatible
)
