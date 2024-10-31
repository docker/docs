module github.com/docker/docs

go 1.23.1

require (
	github.com/docker/buildx v0.17.1 // indirect
	github.com/docker/cli v27.3.2-0.20241008150905-cb3048fbebb1+incompatible // indirect
	github.com/docker/compose/v2 v2.30.1 // indirect
	github.com/docker/scout-cli v1.13.0 // indirect
	github.com/moby/buildkit v0.17.0 // indirect
	github.com/moby/moby v27.3.1+incompatible // indirect
)

replace (
	github.com/docker/buildx => github.com/docker/buildx v0.17.1
	github.com/docker/cli => github.com/docker/cli v27.3.1+incompatible
	github.com/docker/compose/v2 => github.com/docker/compose/v2 v2.30.1
	github.com/docker/scout-cli => github.com/docker/scout-cli v1.13.0
	github.com/moby/buildkit => github.com/moby/buildkit v0.17.0
	github.com/moby/moby => github.com/moby/moby v27.3.1+incompatible
)
