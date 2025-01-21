---
title: Bake file reference
---

The Bake file is a file for defining workflows that you run using `docker buildx bake`.

## File format

You can define your Bake file in the following file formats:

- HashiCorp Configuration Language (HCL)
- JSON
- YAML (Compose file)

By default, Bake uses the following lookup order to find the configuration file:

1. `compose.yaml`
2. `compose.yml`
3. `docker-compose.yml`
4. `docker-compose.yaml`
5. `docker-bake.json`
6. `docker-bake.hcl`
7. `docker-bake.override.json`
8. `docker-bake.override.hcl`

You can specify the file location explicitly using the `--file` flag:

```console
$ docker buildx bake --file ../docker/bake.hcl --print
```

If you don't specify a file explicitly, Bake searches for the file in the
current working directory. If more than one Bake file is found, all files are
merged into a single definition. Files are merged according to the lookup
order. That means that if your project contains both a `compose.yaml` file and
a `docker-bake.hcl` file, Bake loads the `compose.yaml` file first, and then
the `docker-bake.hcl` file.

If merged files contain duplicate attribute definitions, those definitions are
either merged or overridden by the last occurrence, depending on the attribute.
The following attributes are overridden by the last occurrence:

- `target.cache-to`
- `target.dockerfile-inline`
- `target.dockerfile`
- `target.outputs`
- `target.platforms`
- `target.pull`
- `target.tags`
- `target.target`

For example, if `compose.yaml` and `docker-bake.hcl` both define the `tags`
attribute, the `docker-bake.hcl` is used.

```console
$ cat compose.yaml
services:
  webapp:
    build:
      context: .
      tags:
        - bar
$ cat docker-bake.hcl
target "webapp" {
  tags = ["foo"]
}
$ docker buildx bake --print webapp
{
  "group": {
    "default": {
      "targets": [
        "webapp"
      ]
    }
  },
  "target": {
    "webapp": {
      "context": ".",
      "dockerfile": "Dockerfile",
      "tags": [
        "foo"
      ]
    }
  }
}
```

All other attributes are merged. For example, if `compose.yaml` and
`docker-bake.hcl` both define unique entries for the `labels` attribute, all
entries are included. Duplicate entries for the same label are overridden.

```console
$ cat compose.yaml
services:
  webapp:
    build:
      context: .
      labels: 
        com.example.foo: "foo"
        com.example.name: "Alice"
$ cat docker-bake.hcl
target "webapp" {
  labels = {
    "com.example.bar" = "bar"
    "com.example.name" = "Bob"
  }
}
$ docker buildx bake --print webapp
{
  "group": {
    "default": {
      "targets": [
        "webapp"
      ]
    }
  },
  "target": {
    "webapp": {
      "context": ".",
      "dockerfile": "Dockerfile",
      "labels": {
        "com.example.foo": "foo",
        "com.example.bar": "bar",
        "com.example.name": "Bob"
      }
    }
  }
}
```

## Syntax

The Bake file supports the following property types:

- `target`: build targets
- `group`: collections of build targets
- `variable`: build arguments and variables
- `function`: custom Bake functions

You define properties as hierarchical blocks in the Bake file.
You can assign one or more attributes to a property.

The following snippet shows a JSON representation of a simple Bake file.
This Bake file defines three properties: a variable, a group, and a target.

```json
{
  "variable": {
    "TAG": {
      "default": "latest"
    }
  },
  "group": {
    "default": {
      "targets": ["webapp"]
    }
  },
  "target": {
    "webapp": {
      "dockerfile": "Dockerfile",
      "tags": ["docker.io/username/webapp:${TAG}"]
    }
  }
}
```

In the JSON representation of a Bake file, properties are objects,
and attributes are values assigned to those objects.

The following example shows the same Bake file in the HCL format:

```hcl
variable "TAG" {
  default = "latest"
}

group "default" {
  targets = ["webapp"]
}

target "webapp" {
  dockerfile = "Dockerfile"
  tags = ["docker.io/username/webapp:${TAG}"]
}
```

HCL is the preferred format for Bake files.
Aside from syntactic differences,
HCL lets you use features that the JSON and YAML formats don't support.

The examples in this document use the HCL format.

## Target

A target reflects a single `docker build` invocation.
Consider the following build command:

```console
$ docker build \
  --file=Dockerfile.webapp \
  --tag=docker.io/username/webapp:latest \
  https://github.com/username/webapp
```

You can express this command in a Bake file as follows:

```hcl
target "webapp" {
  dockerfile = "Dockerfile.webapp"
  tags = ["docker.io/username/webapp:latest"]
  context = "https://github.com/username/webapp"
}
```

The following table shows the complete list of attributes that you can assign to a target:

| Name                                            | Type    | Description                                                          |
|-------------------------------------------------|---------|----------------------------------------------------------------------|
| [`args`](#targetargs)                           | Map     | Build arguments                                                      |
| [`annotations`](#targetannotations)             | List    | Exporter annotations                                                 |
| [`attest`](#targetattest)                       | List    | Build attestations                                                   |
| [`cache-from`](#targetcache-from)               | List    | External cache sources                                               |
| [`cache-to`](#targetcache-to)                   | List    | External cache destinations                                          |
| [`context`](#targetcontext)                     | String  | Set of files located in the specified path or URL                    |
| [`contexts`](#targetcontexts)                   | Map     | Additional build contexts                                            |
| [`dockerfile-inline`](#targetdockerfile-inline) | String  | Inline Dockerfile string                                             |
| [`dockerfile`](#targetdockerfile)               | String  | Dockerfile location                                                  |
| [`inherits`](#targetinherits)                   | List    | Inherit attributes from other targets                                |
| [`labels`](#targetlabels)                       | Map     | Metadata for images                                                  |
| [`matrix`](#targetmatrix)                       | Map     | Define a set of variables that forks a target into multiple targets. |
| [`name`](#targetname)                           | String  | Override the target name when using a matrix.                        |
| [`no-cache-filter`](#targetno-cache-filter)     | List    | Disable build cache for specific stages                              |
| [`no-cache`](#targetno-cache)                   | Boolean | Disable build cache completely                                       |
| [`output`](#targetoutput)                       | List    | Output destinations                                                  |
| [`platforms`](#targetplatforms)                 | List    | Target platforms                                                     |
| [`pull`](#targetpull)                           | Boolean | Always pull images                                                   |
| [`secret`](#targetsecret)                       | List    | Secrets to expose to the build                                       |
| [`shm-size`](#targetshm-size)                   | List    | Size of `/dev/shm`                                                   |
| [`ssh`](#targetssh)                             | List    | SSH agent sockets or keys to expose to the build                     |
| [`tags`](#targettags)                           | List    | Image names and tags                                                 |
| [`target`](#targettarget)                       | String  | Target build stage                                                   |
| [`ulimits`](#targetulimits)                     | List    | Ulimit options                                                       |

### `target.args`

Use the `args` attribute to define build arguments for the target.
This has the same effect as passing a [`--build-arg`][build-arg] flag to the build command.

```hcl
target "default" {
  args = {
    VERSION = "0.0.0+unknown"
  }
}
```

You can set `args` attributes to use `null` values.
Doing so forces the `target` to use the `ARG` value specified in the Dockerfile.

```hcl
variable "GO_VERSION" {
  default = "1.20.3"
}

target "webapp" {
  dockerfile = "webapp.Dockerfile"
  tags = ["docker.io/username/webapp"]
}

target "db" {
  args = {
    GO_VERSION = null
  }
  dockerfile = "db.Dockerfile"
  tags = ["docker.io/username/db"]
}
```

### `target.annotations`

The `annotations` attribute lets you add annotations to images built with bake.
The key takes a list of annotations, in the format of `KEY=VALUE`.

```hcl
target "default" {
  output = ["type=image,name=foo"]
  annotations = ["org.opencontainers.image.authors=dvdksn"]
}
```

is the same as

```hcl
target "default" {
  output = ["type=image,name=foo,annotation.org.opencontainers.image.authors=dvdksn"]
}
```

By default, the annotation is added to image manifests. You can configure the
level of the annotations by adding a prefix to the annotation, containing a
comma-separated list of all the levels that you want to annotate. The following
example adds annotations to both the image index and manifests.

```hcl
target "default" {
  output = ["type=image,name=foo"]
  annotations = ["index,manifest:org.opencontainers.image.authors=dvdksn"]
}
```

Read about the supported levels in
[Specifying annotation levels](https://docs.docker.com/build/building/annotations/#specifying-annotation-levels).

### `target.attest`

The `attest` attribute lets you apply [build attestations][attestations] to the target.
This attribute accepts the long-form CSV version of attestation parameters.

```hcl
target "default" {
  attest = [
    "type=provenance,mode=min",
    "type=sbom"
  ]
}
```

### `target.cache-from`

Build cache sources.
The builder imports cache from the locations you specify.
It uses the [Buildx cache storage backends][cache-backends],
and it works the same way as the [`--cache-from`][cache-from] flag.
This takes a list value, so you can specify multiple cache sources.

```hcl
target "app" {
  cache-from = [
    "type=s3,region=eu-west-1,bucket=mybucket",
    "user/repo:cache",
  ]
}
```

### `target.cache-to`

Build cache export destinations.
The builder exports its build cache to the locations you specify.
It uses the [Buildx cache storage backends][cache-backends],
and it works the same way as the [`--cache-to` flag][cache-to].
This takes a list value, so you can specify multiple cache export targets.

```hcl
target "app" {
  cache-to = [
    "type=s3,region=eu-west-1,bucket=mybucket",
    "type=inline"
  ]
}
```

### `target.call`

Specifies the frontend method to use. Frontend methods let you, for example,
execute build checks only, instead of running a build. This is the same as the
`--call` flag.

```hcl
target "app" {
  call = "check"
}
```

For more information about frontend methods, refer to the CLI reference for
[`docker buildx build --call`](https://docs.docker.com/reference/cli/docker/buildx/build/#call).

### `target.context`

Specifies the location of the build context to use for this target.
Accepts a URL or a directory path.
This is the same as the [build context][context] positional argument
that you pass to the build command.

```hcl
target "app" {
  context = "./src/www"
}
```

This resolves to the current working directory (`"."`) by default.

```console
$ docker buildx bake --print -f - <<< 'target "default" {}'
[+] Building 0.0s (0/0)
{
  "target": {
    "default": {
      "context": ".",
      "dockerfile": "Dockerfile"
    }
  }
}
```

### `target.contexts`

Additional build contexts.
This is the same as the [`--build-context` flag][build-context].
This attribute takes a map, where keys result in named contexts that you can
reference in your builds.

You can specify different types of contexts, such local directories, Git URLs,
and even other Bake targets. Bake automatically determines the type of
a context based on the pattern of the context value.

| Context type    | Example                                   |
| --------------- | ----------------------------------------- |
| Container image | `docker-image://alpine@sha256:0123456789` |
| Git URL         | `https://github.com/user/proj.git`        |
| HTTP URL        | `https://example.com/files`               |
| Local directory | `../path/to/src`                          |
| Bake target     | `target:base`                             |

#### Pin an image version

```hcl
# docker-bake.hcl
target "app" {
    contexts = {
        alpine = "docker-image://alpine:3.13"
    }
}
```

```Dockerfile
# Dockerfile
FROM alpine
RUN echo "Hello world"
```

#### Use a local directory

```hcl
# docker-bake.hcl
target "app" {
    contexts = {
        src = "../path/to/source"
    }
}
```

```Dockerfile
# Dockerfile
FROM scratch AS src
FROM golang
COPY --from=src . .
```

#### Use another target as base

> [!NOTE]
> You should prefer to use regular multi-stage builds over this option. You can
> Use this feature when you have multiple Dockerfiles that can't be easily
> merged into one.

```hcl
# docker-bake.hcl
target "base" {
    dockerfile = "baseapp.Dockerfile"
}
target "app" {
    contexts = {
        baseapp = "target:base"
    }
}
```

```Dockerfile
# Dockerfile
FROM baseapp
RUN echo "Hello world"
```

### `target.dockerfile-inline`

Uses the string value as an inline Dockerfile for the build target.

```hcl
target "default" {
  dockerfile-inline = "FROM alpine\nENTRYPOINT [\"echo\", \"hello\"]"
}
```

The `dockerfile-inline` takes precedence over the `dockerfile` attribute.
If you specify both, Bake uses the inline version.

### `target.dockerfile`

Name of the Dockerfile to use for the build.
This is the same as the [`--file` flag][file] for the `docker build` command.

```hcl
target "default" {
  dockerfile = "./src/www/Dockerfile"
}
```

Resolves to `"Dockerfile"` by default.

```console
$ docker buildx bake --print -f - <<< 'target "default" {}'
[+] Building 0.0s (0/0)
{
  "target": {
    "default": {
      "context": ".",
      "dockerfile": "Dockerfile"
    }
  }
}
```

### `target.entitlements`

Entitlements are permissions that the build process requires to run.

Currently supported entitlements are:

- `network.host`: Allows the build to use commands that access the host network. In Dockerfile, use [`RUN --network=host`](https://docs.docker.com/reference/dockerfile/#run---networkhost) to run a command with host network enabled.

- `security.insecure`: Allows the build to run commands in privileged containers that are not limited by the default security sandbox. Such container may potentially access and modify system resources. In Dockerfile, use [`RUN --security=insecure`](https://docs.docker.com/reference/dockerfile/#run---security) to run a command in a privileged container.

```hcl
target "integration-tests" {
  # this target requires privileged containers to run nested containers
  entitlements = ["security.insecure"]
}
```

Entitlements are enabled with a two-step process. First, a target must declare the entitlements it requires. Secondly, when invoking the `bake` command, the user must grant the entitlements by passing the `--allow` flag or confirming the entitlements when prompted in an interactive terminal. This is to ensure that the user is aware of the possibly insecure permissions they are granting to the build process.

### `target.inherits`

A target can inherit attributes from other targets.
Use `inherits` to reference from one target to another.

In the following example,
the `app-dev` target specifies an image name and tag.
The `app-release` target uses `inherits` to reuse the tag name.

```hcl
variable "TAG" {
  default = "latest"
}

target "app-dev" {
  tags = ["docker.io/username/myapp:${TAG}"]
}

target "app-release" {
  inherits = ["app-dev"]
  platforms = ["linux/amd64", "linux/arm64"]
}
```

The `inherits` attribute is a list,
meaning you can reuse attributes from multiple other targets.
In the following example, the `app-release` target reuses attributes
from both the `app-dev` and `_release` targets.

```hcl
target "app-dev" {
  args = {
    GO_VERSION = "1.20"
    BUILDX_EXPERIMENTAL = 1
  }
  tags = ["docker.io/username/myapp"]
  dockerfile = "app.Dockerfile"
  labels = {
    "org.opencontainers.image.source" = "https://github.com/username/myapp"
  }
}

target "_release" {
  args = {
    BUILDKIT_CONTEXT_KEEP_GIT_DIR = 1
    BUILDX_EXPERIMENTAL = 0
  }
}

target "app-release" {
  inherits = ["app-dev", "_release"]
  platforms = ["linux/amd64", "linux/arm64"]
}
```

When inheriting attributes from multiple targets and there's a conflict,
the target that appears last in the `inherits` list takes precedence.
The previous example defines the `BUILDX_EXPERIMENTAL` argument twice for the `app-release` target.
It resolves to `0` because the `_release` target appears last in the inheritance chain:

```console
$ docker buildx bake --print app-release
[+] Building 0.0s (0/0)
{
  "group": {
    "default": {
      "targets": [
        "app-release"
      ]
    }
  },
  "target": {
    "app-release": {
      "context": ".",
      "dockerfile": "app.Dockerfile",
      "args": {
        "BUILDKIT_CONTEXT_KEEP_GIT_DIR": "1",
        "BUILDX_EXPERIMENTAL": "0",
        "GO_VERSION": "1.20"
      },
      "labels": {
        "org.opencontainers.image.source": "https://github.com/username/myapp"
      },
      "tags": [
        "docker.io/username/myapp"
      ],
      "platforms": [
        "linux/amd64",
        "linux/arm64"
      ]
    }
  }
}
```

### `target.labels`

Assigns image labels to the build.
This is the same as the `--label` flag for `docker build`.

```hcl
target "default" {
  labels = {
    "org.opencontainers.image.source" = "https://github.com/username/myapp"
    "com.docker.image.source.entrypoint" = "Dockerfile"
  }
}
```

It's possible to use a `null` value for labels.
If you do, the builder uses the label value specified in the Dockerfile.

### `target.matrix`

A matrix strategy lets you fork a single target into multiple different
variants, based on parameters that you specify.
This works in a similar way to [Matrix strategies for GitHub Actions].
You can use this to reduce duplication in your bake definition.

The `matrix` attribute is a map of parameter names to lists of values.
Bake builds each possible combination of values as a separate target.

Each generated target **must** have a unique name.
To specify how target names should resolve, use the `name` attribute.

The following example resolves the `app` target to `app-foo` and `app-bar`.
It also uses the matrix value to define the [target build stage](#targettarget).

```hcl
target "app" {
  name = "app-${tgt}"
  matrix = {
    tgt = ["foo", "bar"]
  }
  target = tgt
}
```

```console
$ docker buildx bake --print app
[+] Building 0.0s (0/0)
{
  "group": {
    "app": {
      "targets": [
        "app-foo",
        "app-bar"
      ]
    },
    "default": {
      "targets": [
        "app"
      ]
    }
  },
  "target": {
    "app-bar": {
      "context": ".",
      "dockerfile": "Dockerfile",
      "target": "bar"
    },
    "app-foo": {
      "context": ".",
      "dockerfile": "Dockerfile",
      "target": "foo"
    }
  }
}
```

#### Multiple axes

You can specify multiple keys in your matrix to fork a target on multiple axes.
When using multiple matrix keys, Bake builds every possible variant.

The following example builds four targets:

- `app-foo-1-0`
- `app-foo-2-0`
- `app-bar-1-0`
- `app-bar-2-0`

```hcl
target "app" {
  name = "app-${tgt}-${replace(version, ".", "-")}"
  matrix = {
    tgt = ["foo", "bar"]
    version = ["1.0", "2.0"]
  }
  target = tgt
  args = {
    VERSION = version
  }
}
```

#### Multiple values per matrix target

If you want to differentiate the matrix on more than just a single value,
you can use maps as matrix values. Bake creates a target for each map,
and you can access the nested values using dot notation.

The following example builds two targets:

- `app-foo-1-0`
- `app-bar-2-0`

```hcl
target "app" {
  name = "app-${item.tgt}-${replace(item.version, ".", "-")}"
  matrix = {
    item = [
      {
        tgt = "foo"
        version = "1.0"
      },
      {
        tgt = "bar"
        version = "2.0"
      }
    ]
  }
  target = item.tgt
  args = {
    VERSION = item.version
  }
}
```

### `target.name`

Specify name resolution for targets that use a matrix strategy.
The following example resolves the `app` target to `app-foo` and `app-bar`.

```hcl
target "app" {
  name = "app-${tgt}"
  matrix = {
    tgt = ["foo", "bar"]
  }
  target = tgt
}
```

### `target.network`

Specify the network mode for the whole build request. This will override the default network mode
for all the `RUN` instructions in the Dockerfile. Accepted values are `default`, `host`, and `none`.

Usually, a better approach to set the network mode for your build steps is to instead use `RUN --network=<value>`
in your Dockerfile. This way, you can set the network mode for individual build steps and everyone building
the Dockerfile gets consistent behavior without needing to pass additional flags to the build command.

If you set network mode to `host` in your Bake file, you must also grant `network.host` entitlement when
invoking the `bake` command. This is because `host` network mode requires elevated privileges and can be a security risk.
You can pass `--allow=network.host` to the `docker buildx bake` command to grant the entitlement, or you can
confirm the entitlement when prompted if you are using an interactive terminal.

```hcl
target "app" {
  # make sure this build does not access internet
  network = "none"
}
```

### `target.no-cache-filter`

Don't use build cache for the specified stages.
This is the same as the `--no-cache-filter` flag for `docker build`.
The following example avoids build cache for the `foo` build stage.

```hcl
target "default" {
  no-cache-filter = ["foo"]
}
```

### `target.no-cache`

Don't use cache when building the image.
This is the same as the `--no-cache` flag for `docker build`.

```hcl
target "default" {
  no-cache = 1
}
```

### `target.output`

Configuration for exporting the build output.
This is the same as the [`--output` flag][output].
The following example configures the target to use a cache-only output,

```hcl
target "default" {
  output = ["type=cacheonly"]
}
```

### `target.platforms`

Set target platforms for the build target.
This is the same as the [`--platform` flag][platform].
The following example creates a multi-platform build for three architectures.

```hcl
target "default" {
  platforms = ["linux/amd64", "linux/arm64", "linux/arm/v7"]
}
```

### `target.pull`

Configures whether the builder should attempt to pull images when building the target.
This is the same as the `--pull` flag for `docker build`.
The following example forces the builder to always pull all images referenced in the build target.

```hcl
target "default" {
  pull = true
}
```

### `target.secret`

Defines secrets to expose to the build target.
This is the same as the [`--secret` flag][secret].

```hcl
variable "HOME" {
  default = null
}

target "default" {
  secret = [
    "type=env,id=KUBECONFIG",
    "type=file,id=aws,src=${HOME}/.aws/credentials"
  ]
}
```

This lets you [mount the secret][run_mount_secret] in your Dockerfile.

```dockerfile
RUN --mount=type=secret,id=aws,target=/root/.aws/credentials \
    aws cloudfront create-invalidation ...
RUN --mount=type=secret,id=KUBECONFIG,env=KUBECONFIG \
    helm upgrade --install
```

### `target.shm-size`

Sets the size of the shared memory allocated for build containers when using
`RUN` instructions.

The format is `<number><unit>`. `number` must be greater than `0`. Unit is
optional and can be `b` (bytes), `k` (kilobytes), `m` (megabytes), or `g`
(gigabytes). If you omit the unit, the system uses bytes.

This is the same as the `--shm-size` flag for `docker build`.

```hcl
target "default" {
  shm-size = "128m"
}
```

> [!NOTE]
> In most cases, it is recommended to let the builder automatically determine
> the appropriate configurations. Manual adjustments should only be considered
> when specific performance tuning is required for complex build scenarios.

### `target.ssh`

Defines SSH agent sockets or keys to expose to the build.
This is the same as the [`--ssh` flag][ssh].
This can be useful if you need to access private repositories during a build.

```hcl
target "default" {
  ssh = ["default"]
}
```

```dockerfile
FROM alpine
RUN --mount=type=ssh \
    apk add git openssh-client \
    && install -m 0700 -d ~/.ssh \
    && ssh-keyscan github.com >> ~/.ssh/known_hosts \
    && git clone git@github.com:user/my-private-repo.git
```

### `target.tags`

Image names and tags to use for the build target.
This is the same as the [`--tag` flag][tag].

```hcl
target "default" {
  tags = [
    "org/repo:latest",
    "myregistry.azurecr.io/team/image:v1"
  ]
}
```

### `target.target`

Set the target build stage to build.
This is the same as the [`--target` flag][target].

```hcl
target "default" {
  target = "binaries"
}
```

### `target.ulimits`

Ulimits overrides the default ulimits of build's containers when using `RUN`
instructions and are specified with a soft and hard limit as such:
`<type>=<soft limit>[:<hard limit>]`, for example:

```hcl
target "app" {
  ulimits = [
    "nofile=1024:1024"
  ]
}
```

> [!NOTE]
> If you do not provide a `hard limit`, the `soft limit` is used
> for both values. If no `ulimits` are set, they are inherited from
> the default `ulimits` set on the daemon.

> [!NOTE]
> In most cases, it is recommended to let the builder automatically determine
> the appropriate configurations. Manual adjustments should only be considered
> when specific performance tuning is required for complex build scenarios.

## Group

Groups allow you to invoke multiple builds (targets) at once.

```hcl
group "default" {
  targets = ["db", "webapp-dev"]
}

target "webapp-dev" {
  dockerfile = "Dockerfile.webapp"
  tags = ["docker.io/username/webapp:latest"]
}

target "db" {
  dockerfile = "Dockerfile.db"
  tags = ["docker.io/username/db"]
}
```

Groups take precedence over targets, if both exist with the same name.
The following bake file builds the `default` group.
Bake ignores the `default` target.

```hcl
target "default" {
  dockerfile-inline = "FROM ubuntu"
}

group "default" {
  targets = ["alpine", "debian"]
}
target "alpine" {
  dockerfile-inline = "FROM alpine"
}
target "debian" {
  dockerfile-inline = "FROM debian"
}
```

## Variable

The HCL file format supports variable block definitions.
You can use variables as build arguments in your Dockerfile,
or interpolate them in attribute values in your Bake file.

```hcl
variable "TAG" {
  default = "latest"
}

target "webapp-dev" {
  dockerfile = "Dockerfile.webapp"
  tags = ["docker.io/username/webapp:${TAG}"]
}
```

You can assign a default value for a variable in the Bake file,
or assign a `null` value to it. If you assign a `null` value,
Buildx uses the default value from the Dockerfile instead.

You can override variable defaults set in the Bake file using environment variables.
The following example sets the `TAG` variable to `dev`,
overriding the default `latest` value shown in the previous example.

```console
$ TAG=dev docker buildx bake webapp-dev
```

### Built-in variables

The following variables are built-ins that you can use with Bake without having
to define them.

| Variable              | Description                                                                         |
| --------------------- | ----------------------------------------------------------------------------------- |
| `BAKE_CMD_CONTEXT`    | Holds the main context when building using a remote Bake file.                      |
| `BAKE_LOCAL_PLATFORM` | Returns the current platform’s default platform specification (e.g. `linux/amd64`). |

### Use environment variable as default

You can set a Bake variable to use the value of an environment variable as a default value:

```hcl
variable "HOME" {
  default = "$HOME"
}
```

### Interpolate variables into attributes

To interpolate a variable into an attribute string value,
you must use curly brackets.
The following doesn't work:

```hcl
variable "HOME" {
  default = "$HOME"
}

target "default" {
  ssh = ["default=$HOME/.ssh/id_rsa"]
}
```

Wrap the variable in curly brackets where you want to insert it:

```diff
  variable "HOME" {
    default = "$HOME"
  }

  target "default" {
-   ssh = ["default=$HOME/.ssh/id_rsa"]
+   ssh = ["default=${HOME}/.ssh/id_rsa"]
  }
```

Before you can interpolate a variable into an attribute,
first you must declare it in the bake file,
as demonstrated in the following example.

```console
$ cat docker-bake.hcl
target "default" {
  dockerfile-inline = "FROM ${BASE_IMAGE}"
}
$ docker buildx bake
[+] Building 0.0s (0/0)
docker-bake.hcl:2
--------------------
   1 |     target "default" {
   2 | >>>   dockerfile-inline = "FROM ${BASE_IMAGE}"
   3 |     }
   4 |
--------------------
ERROR: docker-bake.hcl:2,31-41: Unknown variable; There is no variable named "BASE_IMAGE"., and 1 other diagnostic(s)
$ cat >> docker-bake.hcl

variable "BASE_IMAGE" {
  default = "alpine"
}

$ docker buildx bake
[+] Building 0.6s (5/5) FINISHED
```

## Function

A [set of general-purpose functions][bake_stdlib]
provided by [go-cty][go-cty]
are available for use in HCL files:

```hcl
# docker-bake.hcl
target "webapp-dev" {
  dockerfile = "Dockerfile.webapp"
  tags = ["docker.io/username/webapp:latest"]
  args = {
    buildno = "${add(123, 1)}"
  }
}
```

In addition, [user defined functions][userfunc]
are also supported:

```hcl
# docker-bake.hcl
function "increment" {
  params = [number]
  result = number + 1
}

target "webapp-dev" {
  dockerfile = "Dockerfile.webapp"
  tags = ["docker.io/username/webapp:latest"]
  args = {
    buildno = "${increment(123)}"
  }
}
```

> [!NOTE]
> See [User defined HCL functions][hcl-funcs] page for more details.

<!-- external links -->

[attestations]: https://docs.docker.com/build/attestations/
[bake_stdlib]: https://github.com/docker/buildx/blob/master/bake/hclparser/stdlib.go
[build-arg]: https://docs.docker.com/reference/cli/docker/image/build/#build-arg
[build-context]: https://docs.docker.com/reference/cli/docker/buildx/build/#build-context
[cache-backends]: https://docs.docker.com/build/cache/backends/
[cache-from]: https://docs.docker.com/reference/cli/docker/buildx/build/#cache-from
[cache-to]: https://docs.docker.com/reference/cli/docker/buildx/build/#cache-to
[context]: https://docs.docker.com/reference/cli/docker/buildx/build/#build-context
[file]: https://docs.docker.com/reference/cli/docker/image/build/#file
[go-cty]: https://github.com/zclconf/go-cty/tree/main/cty/function/stdlib
[hcl-funcs]: https://docs.docker.com/build/bake/hcl-funcs/
[output]: https://docs.docker.com/reference/cli/docker/buildx/build/#output
[platform]: https://docs.docker.com/reference/cli/docker/buildx/build/#platform
[run_mount_secret]: https://docs.docker.com/reference/dockerfile/#run---mounttypesecret
[secret]: https://docs.docker.com/reference/cli/docker/buildx/build/#secret
[ssh]: https://docs.docker.com/reference/cli/docker/buildx/build/#ssh
[tag]: https://docs.docker.com/reference/cli/docker/image/build/#tag
[target]: https://docs.docker.com/reference/cli/docker/image/build/#target
[userfunc]: https://github.com/hashicorp/hcl/tree/main/ext/userfunc
