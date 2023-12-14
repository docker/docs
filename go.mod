module github.com/docker/docs

go 1.21

toolchain go1.21.1

require (
	github.com/compose-spec/compose-spec v0.0.0-20231121152139-478928e7c9f8 // indirect
	github.com/docker/buildx v0.12.1-0.20231214091505-b68ee824c673 // indirect
	github.com/docker/cli v25.0.0-beta.1+incompatible // indirect
	github.com/docker/compose/v2 v2.23.3 // indirect
	github.com/docker/scout-cli v1.2.0 // indirect
	github.com/moby/buildkit v0.13.0-beta1.0.20231214000015-a960fe501f00 // indirect
	github.com/moby/moby v24.0.5+incompatible // indirect
)

// buildkit depends on cli v25 beta1, pin to v24
replace github.com/docker/cli => github.com/docker/cli v24.0.8-0.20231213094340-0f82fd88610a+incompatible
