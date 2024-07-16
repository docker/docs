module github.com/docker/docs

go 1.21

toolchain go1.21.1

require (
	github.com/docker/buildx v0.15.1 // indirect
	github.com/docker/cli v27.0.3+incompatible // indirect
	github.com/docker/compose/v2 v2.28.1 // indirect
	github.com/docker/scout-cli v1.10.0 // indirect
	github.com/moby/buildkit v0.14.1 // indirect
	github.com/moby/moby v27.0.3+incompatible // indirect
)

replace (
	github.com/docker/buildx => github.com/docker/buildx v0.15.1
	github.com/docker/cli => github.com/docker/cli v27.0.3+incompatible
	github.com/docker/compose/v2 => github.com/docker/compose/v2 v2.28.1
	github.com/docker/scout-cli => github.com/docker/scout-cli v1.10.0
	github.com/moby/buildkit => github.com/moby/buildkit v0.14.0-rc2.0.20240611065153-eed17a45c62b
	github.com/moby/moby => github.com/moby/moby v27.0.3+incompatible
)
