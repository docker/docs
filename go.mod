module github.com/docker/docs

go 1.21

toolchain go1.21.1

require (
	github.com/docker/buildx v0.12.1 // indirect
	github.com/docker/cli v25.0.0-rc.1+incompatible // indirect
	github.com/docker/compose/v2 v2.24.0 // indirect
	github.com/docker/scout-cli v1.2.0 // indirect
	github.com/moby/buildkit v0.13.0-beta1.0.20231219135447-957cb50df991 // indirect
	github.com/moby/moby v24.0.5+incompatible // indirect
)

// buildkit depends on cli v25 beta1, pin to v24
replace github.com/docker/cli => github.com/docker/cli v24.0.8-0.20240103162225-b0c5946ba5d8+incompatible
