module github.com/docker/docs

go 1.23.1

require (
	github.com/docker/buildx v0.20.1 // indirect
	github.com/docker/cli v27.5.1+incompatible // indirect
	github.com/docker/compose/v2 v2.32.4 // indirect
	github.com/docker/scout-cli v1.15.0 // indirect
	github.com/moby/buildkit v0.19.0 // indirect
	github.com/moby/moby v27.5.1+incompatible // indirect
)

replace (
	github.com/docker/buildx => github.com/docker/buildx v0.20.1
	github.com/docker/cli => github.com/docker/cli v27.5.1+incompatible
	github.com/docker/compose/v2 => github.com/docker/compose/v2 v2.32.4
	github.com/docker/scout-cli => github.com/docker/scout-cli v1.15.0
	github.com/moby/buildkit => github.com/moby/buildkit v0.19.0
	github.com/moby/moby => github.com/moby/moby v27.5.1+incompatible
)
