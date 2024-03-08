module github.com/docker/docs

go 1.21

toolchain go1.21.1

require (
	github.com/docker/buildx v0.12.0-rc2.0.20231219140829-617f538cb315 // indirect
	github.com/docker/cli v26.0.0-rc1+incompatible // indirect
	github.com/docker/compose/v2 v2.0.0-00010101000000-000000000000 // indirect
	github.com/docker/scout-cli v1.4.1 // indirect
	github.com/moby/buildkit v0.13.0 // indirect
	github.com/moby/moby v25.0.4+incompatible // indirect
)

replace (
	github.com/docker/buildx => github.com/docker/buildx v0.13.1-0.20240307093612-37b7ad1465d2
	github.com/docker/cli => github.com/docker/cli v25.0.4+incompatible
	github.com/docker/compose/v2 => github.com/docker/compose/v2 v2.24.7
	github.com/docker/scout-cli => github.com/docker/scout-cli v1.4.1
	github.com/moby/buildkit => github.com/moby/buildkit v0.13.0-rc3.0.20240308080452-a38011b9f57d
	github.com/moby/moby => github.com/moby/moby v25.0.4+incompatible
)
