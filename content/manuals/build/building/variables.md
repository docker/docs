---
title: Build variables
linkTitle: Variables
weight: 20
description: Using build arguments and environment variables to configure builds
keywords: build, args, variables, parameters, env, environment variables, config
aliases:
- /build/buildkit/color-output-controls/
- /build/building/env-vars/
- /build/guide/build-args/
---

In Docker Build, build arguments (`ARG`) and environment variables (`ENV`)
both serve as a means to pass information into the build process.
You can use them to parameterize the build, allowing for more flexible and configurable builds.

> [!WARNING]
>
> Build arguments and environment variables are inappropriate for passing secrets
> to your build, because they're exposed in the final image. Instead, use
> secret mounts or SSH mounts, which expose secrets to your builds securely.
>
> See [Build secrets](./secrets.md) for more information.

## Similarities and differences

Build arguments and environment variables are similar.
They're both declared in the Dockerfile and can be set using flags for the `docker build` command.
Both can be used to parameterize the build.
But they each serve a distinct purpose.

### Build arguments

Build arguments are variables for the Dockerfile itself.
Use them to parameterize values of Dockerfile instructions.
For example, you might use a build argument to specify the version of a dependency to install.

Build arguments have no effect on the build unless it's used in an instruction.
They're not accessible or present in containers instantiated from the image
unless explicitly passed through from the Dockerfile into the image filesystem or configuration.
They may persist in the image metadata, as provenance attestations and in the image history,
which is why they're not suitable for holding secrets.

They make Dockerfiles more flexible, and easier to maintain.

For an example on how you can use build arguments,
see [`ARG` usage example](#arg-usage-example).

### Environment variables

Environment variables are passed through to the build execution environment,
and persist in containers instantiated from the image.

Environment variables are primarily used to:

- Configure the execution environment for builds
- Set default environment variables for containers

Environment variables, if set, can directly influence the execution of your build,
and the behavior or configuration of the application.

You can't override or set an environment variable at build-time.
Values for environment variables must be declared in the Dockerfile.
You can combine environment variables and build arguments to allow
environment variables to be configured at build-time.

For an example on how to use environment variables for configuring builds,
see [`ENV` usage example](#env-usage-example).

## `ARG` usage example

Build arguments are commonly used to specify versions of components,
such as image variants or package versions, used in a build.

Specifying versions as build arguments lets you build with different versions
without having to manually update the Dockerfile.
It also makes it easier to maintain the Dockerfile,
since it lets you declare versions at the top of the file.

Build arguments can also be a way to reuse a value in multiple places.
For example, if you use multiple flavors of `alpine` in your build,
you can ensure you're using the same version of `alpine` everywhere:

- `golang:1.22-alpine${ALPINE_VERSION}`
- `python:3.12-alpine${ALPINE_VERSION}`
- `nginx:1-alpine${ALPINE_VERSION}`

The following example defines the version of `node` and `alpine` using build arguments.

```dockerfile
# syntax=docker/dockerfile:1

ARG NODE_VERSION="{{% param example_node_version %}}"
ARG ALPINE_VERSION="{{% param example_alpine_version %}}"

FROM node:${NODE_VERSION}-alpine${ALPINE_VERSION} AS base
WORKDIR /src

FROM base AS build
COPY package*.json ./
RUN npm ci
RUN npm run build

FROM base AS production
COPY package*.json ./
RUN npm ci --omit=dev && npm cache clean --force
COPY --from=build /src/dist/ .
CMD ["node", "app.js"]
```

In this case, the build arguments have default values.
Specifying their values when you invoke a build is optional.
To override the defaults, you would use the `--build-arg` CLI flag:

```console
$ docker build --build-arg NODE_VERSION=current .
```

For more information on how to use build arguments, refer to:

- [`ARG` Dockerfile reference](/reference/dockerfile.md#arg)
- [`docker build --build-arg` reference](/reference/cli/docker/buildx/build.md#build-arg)

## `ENV` usage example

Declaring an environment variable with `ENV` makes the variable
available to all subsequent instructions in the build stage.
The following example shows an example setting `NODE_ENV` to `production`
before installing JavaScript dependencies with `npm`.
Setting the variable makes `npm` omits packages needed only for local development.

```dockerfile
# syntax=docker/dockerfile:1

FROM node:20
WORKDIR /app
COPY package*.json ./
ENV NODE_ENV=production
RUN npm ci && npm cache clean --force
COPY . .
CMD ["node", "app.js"]
```

Environment variables aren't configurable at build-time by default.
If you want to change the value of an `ENV` at build-time,
you can combine environment variables and build arguments:

```dockerfile
# syntax=docker/dockerfile:1

FROM node:20
ARG NODE_ENV=production
ENV NODE_ENV=$NODE_ENV
WORKDIR /app
COPY package*.json ./
RUN npm ci && npm cache clean --force
COPY . .
CMD ["node", "app.js"]
```

With this Dockerfile, you can use `--build-arg` to override the default value of `NODE_ENV`:

```console
$ docker build --build-arg NODE_ENV=development .
```

Note that, because the environment variables you set persist in containers,
using them can lead to unintended side-effects for the application's runtime.

For more information on how to use environment variables in builds, refer to:

- [`ENV` Dockerfile reference](/reference/dockerfile.md#env)

## Scoping

Build arguments declared in the global scope of a Dockerfile
aren't automatically inherited into the build stages.
They're only accessible in the global scope.

```dockerfile
# syntax=docker/dockerfile:1

# The following build argument is declared in the global scope:
ARG NAME="joe"

FROM alpine
# The following instruction doesn't have access to the $NAME build argument
# because the argument was defined in the global scope, not for this stage.
RUN echo "hello ${NAME}!"
```

The `echo` command in this example evaluates to `hello !`
because the value of the `NAME` build argument is out of scope.
To inherit global build arguments into a stage, you must consume them:

```dockerfile
# syntax=docker/dockerfile:1

# Declare the build argument in the global scope
ARG NAME="joe"

FROM alpine
# Consume the build argument in the build stage
ARG NAME
RUN echo $NAME
```

Once a build argument is declared or consumed in a stage,
it's automatically inherited by child stages.

```dockerfile
# syntax=docker/dockerfile:1
FROM alpine AS base
# Declare the build argument in the build stage
ARG NAME="joe"

# Create a new stage based on "base"
FROM base AS build
# The NAME build argument is available here
# since it's declared in a parent stage
RUN echo "hello $NAME!"
```

The following diagram further exemplifies how build argument
and environment variable inheritance works for multi-stage builds.

{{< figure src="../../images/build-variables.svg" class="invertible" >}}

## Pre-defined build arguments

This section describes pre-defined build arguments available to all builds by default.

### Multi-platform build arguments

Multi-platform build arguments describe the build and target platforms for the build.

The build platform is the operating system, architecture, and platform variant
of the host system where the builder (the BuildKit daemon) is running.

- `BUILDPLATFORM`
- `BUILDOS`
- `BUILDARCH`
- `BUILDVARIANT`

The target platform arguments hold the same values for the target platforms for the build,
specified using the `--platform` flag for the `docker build` command.

- `TARGETPLATFORM`
- `TARGETOS`
- `TARGETARCH`
- `TARGETVARIANT`

These arguments are useful for doing cross-compilation in multi-platform builds.
They're available in the global scope of the Dockerfile,
but they aren't automatically inherited by build stages.
To use them inside stage, you must declare them:

```dockerfile
# syntax=docker/dockerfile:1

# Pre-defined build arguments are available in the global scope
FROM --platform=$BUILDPLATFORM golang
# To inherit them to a stage, declare them with ARG
ARG TARGETOS
RUN GOOS=$TARGETOS go build -o ./exe .
```

For more information about multi-platform build arguments, refer to
[Multi-platform arguments](/reference/dockerfile.md#automatic-platform-args-in-the-global-scope)

### Proxy arguments

Proxy build arguments let you specify proxies to use for your build.
You don't need to declare or reference these arguments in the Dockerfile.
Specifying a proxy with `--build-arg` is enough to make your build use the proxy.

Proxy arguments are automatically excluded from the build cache
and the output of `docker history` by default.
If you do reference the arguments in your Dockerfile,
the proxy configuration ends up in the build cache.

The builder respects the following proxy build arguments.
The variables are case insensitive.

- `HTTP_PROXY`
- `HTTPS_PROXY`
- `FTP_PROXY`
- `NO_PROXY`
- `ALL_PROXY`

To configure a proxy for your build:

```console
$ docker build --build-arg HTTP_PROXY=https://my-proxy.example.com .
```

For more information about proxy build arguments, refer to
[Proxy arguments](/reference/dockerfile.md#predefined-args).

## Build tool configuration variables

The following environment variables enable, disable, or change the behavior of Buildx and BuildKit.
Note that these variables aren't used to configure the build container;
they aren't available inside the build and they have no relation to the `ENV` instruction.
They're used to configure the Buildx client, or the BuildKit daemon.

| Variable                                                                    | Type              | Description                                                      |
|-----------------------------------------------------------------------------|-------------------|------------------------------------------------------------------|
| [BUILDKIT_COLORS](#buildkit_colors)                                         | String            | Configure text color for the terminal output.                    |
| [BUILDKIT_HOST](#buildkit_host)                                             | String            | Specify host to use for remote builders.                         |
| [BUILDKIT_PROGRESS](#buildkit_progress)                                     | String            | Configure type of progress output.                               |
| [BUILDKIT_TTY_LOG_LINES](#buildkit_tty_log_lines)                           | String            | Number of log lines (for active steps in TTY mode).              |
| [BUILDX_BAKE_FILE](#buildx_bake_file)                                       | String            | Specify the build definition file(s) for `docker buildx bake`.   |
| [BUILDX_BAKE_FILE_SEPARATOR](#buildx_bake_file_separator)                   | String            | Specify the file-path separator for `BUILDX_BAKE_FILE`.          |
| [BUILDX_BAKE_GIT_AUTH_HEADER](#buildx_bake_git_auth_header)                 | String            | HTTP authentication scheme for remote Bake files.                |
| [BUILDX_BAKE_GIT_AUTH_TOKEN](#buildx_bake_git_auth_token)                   | String            | HTTP authentication token for remote Bake files.                 |
| [BUILDX_BAKE_GIT_SSH](#buildx_bake_git_ssh)                                 | String            | SSH authentication for remote Bake files.                        |
| [BUILDX_BUILDER](#buildx_builder)                                           | String            | Specify the builder instance to use.                             |
| [BUILDX_CONFIG](#buildx_config)                                             | String            | Specify location for configuration, state, and logs.             |
| [BUILDX_CPU_PROFILE](#buildx_cpu_profile)                                   | String            | Generate a `pprof` CPU profile at the specified location.        |
| [BUILDX_EXPERIMENTAL](#buildx_experimental)                                 | Boolean           | Turn on experimental features.                                   |
| [BUILDX_GIT_CHECK_DIRTY](#buildx_git_check_dirty)                           | Boolean           | Enable dirty Git checkout detection.                             |
| [BUILDX_GIT_INFO](#buildx_git_info)                                         | Boolean           | Remove Git information in provenance attestations.               |
| [BUILDX_GIT_LABELS](#buildx_git_labels)                                     | String \| Boolean | Add Git provenance labels to images.                             |
| [BUILDX_MEM_PROFILE](#buildx_mem_profile)                                   | String            | Generate a `pprof` memory profile at the specified location.     |
| [BUILDX_METADATA_PROVENANCE](#buildx_metadata_provenance)                   | String \| Boolean | Customize provenance information included in the metadata file.  |
| [BUILDX_METADATA_WARNINGS](#buildx_metadata_warnings)                       | String            | Include build warnings in the metadata file.                     |
| [BUILDX_NO_DEFAULT_ATTESTATIONS](#buildx_no_default_attestations)           | Boolean           | Turn off default provenance attestations.                        |
| [BUILDX_NO_DEFAULT_LOAD](#buildx_no_default_load)                           | Boolean           | Turn off loading images to image store by default.               |
| [EXPERIMENTAL_BUILDKIT_SOURCE_POLICY](#experimental_buildkit_source_policy) | String            | Specify a BuildKit source policy file.                           |

BuildKit also supports a few additional configuration parameters. Refer to
[BuildKit built-in build args](/reference/dockerfile.md#buildkit-built-in-build-args).

You can express Boolean values for environment variables in different ways.
For example, `true`, `1`, and `T` all evaluate to true.
Evaluation is done using the `strconv.ParseBool` function in the Go standard library.
See the [reference documentation](https://pkg.go.dev/strconv#ParseBool) for details.

<!-- vale Docker.HeadingSentenceCase = NO -->

### BUILDKIT_COLORS

Changes the colors of the terminal output. Set `BUILDKIT_COLORS` to a CSV string
in the following format:

```console
$ export BUILDKIT_COLORS="run=123,20,245:error=yellow:cancel=blue:warning=white"
```

Color values can be any valid RGB hex code, or one of the
[BuildKit predefined colors](https://github.com/moby/buildkit/blob/master/util/progress/progressui/colors.go).

Setting `NO_COLOR` to anything turns off colorized output, as recommended by
[no-color.org](https://no-color.org/).

### BUILDKIT_HOST

{{< summary-bar feature_name="Buildkit host" >}}

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

### BUILDKIT_PROGRESS

Sets the type of the BuildKit progress output. Valid values are:

- `auto` (default)
- `plain`
- `tty`
- `quiet`
- `rawjson`

Usage:

```console
$ export BUILDKIT_PROGRESS=plain
```

### BUILDKIT_TTY_LOG_LINES

You can change how many log lines are visible for active steps in TTY mode by
setting `BUILDKIT_TTY_LOG_LINES` to a number (default to `6`).

```console
$ export BUILDKIT_TTY_LOG_LINES=8
```

### EXPERIMENTAL_BUILDKIT_SOURCE_POLICY

Lets you specify a
[BuildKit source policy](https://github.com/moby/buildkit/blob/master/docs/build-repro.md#reproducing-the-pinned-dependencies)
file for creating reproducible builds with pinned dependencies.

```console
$ export EXPERIMENTAL_BUILDKIT_SOURCE_POLICY=./policy.json
```

Example:

```json
{
  "rules": [
    {
      "action": "CONVERT",
      "selector": {
        "identifier": "docker-image://docker.io/library/alpine:latest"
      },
      "updates": {
        "identifier": "docker-image://docker.io/library/alpine:latest@sha256:4edbd2beb5f78b1014028f4fbb99f3237d9561100b6881aabbf5acce2c4f9454"
      }
    },
    {
      "action": "CONVERT",
      "selector": {
        "identifier": "https://raw.githubusercontent.com/moby/buildkit/v0.10.1/README.md"
      },
      "updates": {
        "attrs": {"http.checksum": "sha256:6e4b94fc270e708e1068be28bd3551dc6917a4fc5a61293d51bb36e6b75c4b53"}
      }
    },
    {
      "action": "DENY",
      "selector": {
        "identifier": "docker-image://docker.io/library/golang*"
      }
    }
  ]
}
```

### BUILDX_BAKE_FILE

{{< summary-bar feature_name="Buildx bake file" >}}

Specify one or more build definition files for `docker buildx bake`. 

This environment variable provides an alternative to the `-f` / `--file` command-line flag.

Multiple files can be specified by separating them with the system path separator (":" on Linux/macOS, ";" on Windows):

```console
export BUILDX_BAKE_FILE=file1.hcl:file2.hcl
```

Or with a custom separator defined by the [BUILDX_BAKE_FILE_SEPARATOR](#buildx_bake_file_separator) variable:

```console
export BUILDX_BAKE_FILE_SEPARATOR=@
export BUILDX_BAKE_FILE=file1.hcl@file2.hcl
```

If both `BUILDX_BAKE_FILE` and the `-f` flag are set, only the files provided via `-f` are used. 

If a listed file does not exist or is invalid, bake returns an error.

### BUILDX_BAKE_FILE_SEPARATOR

{{< summary-bar feature_name="Buildx bake file separator" >}}

Controls the separator used between file paths in the `BUILDX_BAKE_FILE` environment variable. 

This is useful if your file paths contain the default separator character or if you want to standardize separators across different platforms.

```console
export BUILDX_BAKE_PATH_SEPARATOR=@
export BUILDX_BAKE_FILE=file1.hcl@file2.hcl
```

### BUILDX_BAKE_GIT_AUTH_HEADER

{{< summary-bar feature_name="Buildx bake Git auth token" >}}

Sets the HTTP authentication scheme when using a remote Bake definition in a private Git repository.
This is equivalent to the [`GIT_AUTH_HEADER` secret](./secrets#http-authentication-scheme),
but facilitates the pre-flight authentication in Bake when loading the remote Bake file.
Supported values are `bearer` (default) and `basic`.

Usage:

```console
$ export BUILDX_BAKE_GIT_AUTH_HEADER=basic
```

### BUILDX_BAKE_GIT_AUTH_TOKEN

{{< summary-bar feature_name="Buildx bake Git auth token" >}}

Sets the HTTP authentication token when using a remote Bake definition in a private Git repository.
This is equivalent to the [`GIT_AUTH_TOKEN` secret](./secrets#git-authentication-for-remote-contexts),
but facilitates the pre-flight authentication in Bake when loading the remote Bake file.

Usage:

```console
$ export BUILDX_BAKE_GIT_AUTH_TOKEN=$(cat git-token.txt)
```

### BUILDX_BAKE_GIT_SSH

{{< summary-bar feature_name="Buildx bake Git SSH" >}}

Lets you specify a list of SSH agent socket filepaths to forward to Bake
for authenticating to a Git server when using a remote Bake definition in a private repository.
This is similar to SSH mounts for builds, but facilitates the pre-flight authentication in Bake when resolving the build definition.

Setting this environment is typically not necessary, because Bake will use the `SSH_AUTH_SOCK` agent socket by default.
You only need to specify this variable if you want to use a socket with a different filepath.
This variable can take multiple paths using a comma-separated string.

Usage:

```console
$ export BUILDX_BAKE_GIT_SSH=/run/foo/listener.sock,~/.creds/ssh.sock
```

### BUILDX_BUILDER

Overrides the configured builder instance. Same as the `docker buildx --builder`
CLI flag.

Usage:

```console
$ export BUILDX_BUILDER=my-builder
```

### BUILDX_CONFIG

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

### BUILDX_CPU_PROFILE

{{< summary-bar feature_name="Buildx CPU profile" >}}

If specified, Buildx generates a `pprof` CPU profile at the specified location.

> [!NOTE]
> This property is only useful for when developing Buildx. The profiling data
> is not relevant for analyzing a build's performance.

Usage:

```console
$ export BUILDX_CPU_PROFILE=buildx_cpu.prof
```

### BUILDX_EXPERIMENTAL

Enables experimental build features.

Usage:

```console
$ export BUILDX_EXPERIMENTAL=1
```

### BUILDX_GIT_CHECK_DIRTY

{{< summary-bar feature_name="Buildx Git check dirty" >}}

When set to true, checks for dirty state in source control information for
[provenance attestations](/manuals/build/metadata/attestations/slsa-provenance.md).

Usage:

```console
$ export BUILDX_GIT_CHECK_DIRTY=1
```

### BUILDX_GIT_INFO

{{< summary-bar feature_name="Buildx Git info" >}}

When set to false, removes source control information from
[provenance attestations](/manuals/build/metadata/attestations/slsa-provenance.md).

Usage:

```console
$ export BUILDX_GIT_INFO=0
```

### BUILDX_GIT_LABELS

{{< summary-bar feature_name="Buildx Git labels" >}}

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

### BUILDX_MEM_PROFILE

{{< summary-bar feature_name="Buildx mem profile" >}}

If specified, Buildx generates a `pprof` memory profile at the specified
location.

> [!NOTE]
> This property is only useful for when developing Buildx. The profiling data
> is not relevant for analyzing a build's performance.

Usage:

```console
$ export BUILDX_MEM_PROFILE=buildx_mem.prof
```

### BUILDX_METADATA_PROVENANCE

{{< summary-bar feature_name="Buildx metadata provenance" >}}

By default, Buildx includes minimal provenance information in the metadata file
through [`--metadata-file` flag](/reference/cli/docker/buildx/build/#metadata-file).
This environment variable allows you to customize the provenance information
included in the metadata file:
* `min` sets minimal provenance (default).
* `max` sets full provenance.
* `disabled`, `false` or `0` does not set any provenance.

### BUILDX_METADATA_WARNINGS

{{< summary-bar feature_name="Buildx metadata warnings" >}}

By default, Buildx does not include build warnings in the metadata file through
[`--metadata-file` flag](/reference/cli/docker/buildx/build/#metadata-file).
You can set this environment variable to `1` or `true` to include them.

### BUILDX_NO_DEFAULT_ATTESTATIONS

{{< summary-bar feature_name="Buildx no default" >}}

By default, BuildKit v0.11 and later adds
[provenance attestations](/manuals/build/metadata/attestations/slsa-provenance.md) to images you
build. Set `BUILDX_NO_DEFAULT_ATTESTATIONS=1` to disable the default provenance
attestations.

Usage:

```console
$ export BUILDX_NO_DEFAULT_ATTESTATIONS=1
```

### BUILDX_NO_DEFAULT_LOAD

When you build an image using the `docker` driver, the image is automatically
loaded to the image store when the build finishes. Set `BUILDX_NO_DEFAULT_LOAD`
to disable automatic loading of images to the local container store.

Usage:

```console
$ export BUILDX_NO_DEFAULT_LOAD=1
```

<!-- vale Docker.HeadingSentenceCase = YES -->
