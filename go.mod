module github.com/docker/docs

go 1.21

toolchain go1.21.1

require (
	github.com/docker/buildx v0.13.1-0.20240307093612-37b7ad1465d2 // indirect
	github.com/docker/cli v26.0.0-rc1+incompatible // indirect
	github.com/docker/compose/v2 v2.24.6 // indirect
	github.com/docker/scout-cli v1.4.1 // indirect
	github.com/moby/buildkit v0.13.0 // indirect
	github.com/moby/moby v25.0.3-0.20240203133757-341a7978a541+incompatible // indirect
)

replace (
	github.com/docker/buildx => github.com/docker/buildx v0.13.1-0.20240307093612-37b7ad1465d2
	github.com/docker/cli => github.com/docker/cli v25.0.4-0.20240221083216-f67e569a8fb9+incompatible
	github.com/docker/compose/v2 => github.com/docker/compose/v2 v2.24.6
	github.com/docker/scout-cli => github.com/docker/scout-cli v1.4.1
	github.com/moby/buildkit => github.com/moby/buildkit v0.13.0-beta3.0.20240201135300-d906167d0b34
	github.com/moby/moby => github.com/moby/moby v25.0.3-0.20240203133757-341a7978a541+incompatible
)
