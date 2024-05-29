module github.com/docker/docs

go 1.21

toolchain go1.21.1

require (
	github.com/docker/buildx v0.14.1 // indirect
	github.com/docker/cli v26.1.3+incompatible // indirect
	github.com/docker/compose/v2 v2.27.0 // indirect
	github.com/docker/scout-cli v1.9.3 // indirect
	github.com/moby/buildkit v0.13.1 // indirect
	github.com/moby/moby v26.1.2+incompatible // indirect
)

replace (
	github.com/docker/buildx => github.com/docker/buildx v0.14.1
	github.com/docker/cli => github.com/docker/cli v26.1.3-0.20240513184838-60f2d38d5341+incompatible
	github.com/docker/compose/v2 => github.com/docker/compose/v2 v2.27.0
	github.com/docker/scout-cli => github.com/docker/scout-cli v1.9.3
	github.com/moby/buildkit => github.com/moby/buildkit v0.13.0-rc3.0.20240424175633-5fce077ed0e0
	github.com/moby/moby => github.com/moby/moby v26.1.2+incompatible
)
