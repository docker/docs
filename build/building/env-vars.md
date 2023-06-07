---
title: Environment variables for Docker Build
description: Use environment variables to configure builds
keywords: docker build, buildx, buildkit, env, environment variables, config
redirect_from:
  - /build/buildkit/color-output-controls/
---

You can set the following environment variables to enable, disable, or change
the behavior of features related to building:

| Variable                                                          | Type              | Description                                          |
| ----------------------------------------------------------------- | ----------------- | ---------------------------------------------------- |
| [BUILDKIT_COLORS](#buildkit_colors)                               | String            | Configure text color for the terminal output.        |
| [BUILDKIT_HOST](#buildkit_host)                                   | String            | Specify host to use for remote builders.             |
| [BUILDKIT_PROGRESS](#buildkit_progress)                           | String            | Configure type of progress output.                   |
| [BUILDX_BUILDER](#buildx_builder)                                 | String            | Specify the builder instance to use.                 |
| [BUILDX_CONFIG](#buildx_config)                                   | String            | Specify location for configuration, state, and logs. |
| [BUILDX_EXPERIMENTAL](#buildx_experimental)                       | Boolean           | Turn on experimental features.                       |
| [BUILDX_GIT_CHECK_DIRTY](#buildx_git_check_dirty)                 | Boolean           | Enable dirty Git checkout detection.                 |
| [BUILDX_GIT_INFO](#buildx_git_info)                               | Boolean           | Remove Git information in provenance attestations.   |
| [BUILDX_GIT_LABELS](#buildx_git_labels)                           | String \| Boolean | Add Git provenance labels to images.                 |
| [BUILDX_NO_DEFAULT_ATTESTATIONS](#buildx_no_default_attestations) | Boolean           | Turn off default provenance attestations.            |
| [BUILDX_NO_DEFAULT_LOAD](#buildx_no_default_load)                 | Boolean           | Turn off loading images to image store by default.   |

See also
[BuildKit built-in build args](../../engine/reference/builder.md#buildkit-built-in-build-args).

You can express Boolean values for environment variables in different ways. For
example, `true`, `1`, and `T` all evaluate to true. Evaluation is done using the
`strconv.ParseBool` function in the Go standard library. See the
[reference documentation](https://pkg.go.dev/strconv#ParseBool){: target="_blank" rel="noopener" class="_" }
for details.

## BUILDKIT_COLORS

Changes the colors of the terminal output. Set `BUILDKIT_COLORS` to a CSV string
in the following format:

```console
$ export BUILDKIT_COLORS="run=123,20,245:error=yellow:cancel=blue:warning=white"
```

Color values can be any valid RGB hex code, or one of the
[BuildKit predefined colors](https://github.com/moby/buildkit/blob/master/util/progress/progressui/colors.go){: target="_blank" rel="noopener" class="_" }.

Setting `NO_COLOR` to anything turns off colorized output, as recommended by
[no-color.org](https://no-color.org/){: target="_blank" rel="noopener" class="_" }.

## BUILDKIT_HOST

_Introduced in [Buildx v0.9.0](../release-notes.md#090)_

You use the `BUILDKIT_HOST` to specify the address of a BuildKit daemon to use
as a remote builder. This is the same as specifying the address as a positional
argument to `docker buildx create`.

Usage:

```console
$ export BUILDKIT_HOST=tcp://localhost:1234
$ docker buildx create --name=remote --driver=remote
```

If you specify both the `BUILDKIT_HOST` environment variable and a positional
argument, the argument takes priority.

## BUILDKIT_PROGRESS

Sets the type of the BuildKit progress output. Valid values are:

- `auto` (default)
- `plain`
- `tty`

Usage:

```console
$ export BUILDKIT_PROGRESS=plain
```

## BUILDX_BUILDER

Overrides the configured builder instance. Same as the `docker buildx --builder`
CLI flag.

Usage:

```console
$ export BUILDX_BUILDER=my-builder
```

## BUILDX_CONFIG

You can use `BUILDX_CONFIG` to specify the directory to use for build
configuration, state, and logs. The lookup order for this directory is as
follows:

- `$BUILDX_CONFIG`
- `$DOCKER_CONFIG/buildx`
- `~/.docker/buildx` (default)

Usage:

```console
$ export BUILDX_CONFIG=/usr/local/etc
```

## BUILDX_EXPERIMENTAL

Enables experimental build features.

Usage:

```console
$ export BUILDX_EXPERIMENTAL=1
```

## BUILDX_GIT_CHECK_DIRTY

_Introduced in [Buildx v0.10.4](../release-notes.md#0104)_

When set to true, checks for dirty state in source control information for
[provenance attestations](../attestations/slsa-provenance.md).

Usage:

```console
$ export BUILDX_GIT_CHECK_DIRTY=1
```

## BUILDX_GIT_INFO

_Introduced in [Buildx v0.10.0](../release-notes.md#0100)_

When set to false, removes source control information from
[provenance attestations](../attestations/slsa-provenance.md).

Usage:

```console
$ export BUILDX_GIT_INFO=0
```

## BUILDX_GIT_LABELS

_Introduced in [Buildx v0.10.0](../release-notes.md#0100)_

Adds provenance labels, based on Git information, to images that you build. The
labels are:

- `com.docker.image.source.entrypoint`: Location of the Dockerfile relative to
  the project root
- `org.opencontainers.image.revision`: Git commit revision
- `org.opencontainers.image.source`: SSH or HTTPS address of the repository

Example:

```json
  "Labels": {
    "com.docker.image.source.entrypoint": "Dockerfile",
    "org.opencontainers.image.revision": "5734329c6af43c2ae295010778cd308866b95d9b",
    "org.opencontainers.image.source": "git@github.com:foo/bar.git"
  }
```

Usage:

- Set `BUILDX_GIT_LABELS=1` to include the `entrypoint` and `revision` labels.
- Set `BUILDX_GIT_LABELS=full` to include all labels.

If the repository is in a dirty state, the `revision` gets a `-dirty` suffix.

## BUILDX_NO_DEFAULT_ATTESTATIONS

_Introduced in [Buildx v0.10.4](../release-notes.md#0104)_

By default, BuildKit v0.11 and later adds
[provenance attestations](../attestations/slsa-provenance.md) to images you
build. Set `BUILDX_NO_DEFAULT_ATTESTATIONS=1` to disable the default provenance
attestations.

Usage:

```console
$ export BUILDX_NO_DEFAULT_ATTESTATIONS=1
```

## BUILDX_NO_DEFAULT_LOAD

When you build an image using the `docker` driver, the image is automatically
loaded to the image store when the build finishes. Set `BUILDX_NO_DEFAULT_LOAD`
to disable automatic loading of images to the local container store.

Usage:

```console
$ export BUILDX_NO_DEFAULT_LOAD=1
```
