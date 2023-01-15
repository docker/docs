---
title: "Bake file definition"
keywords: build, buildx, bake, buildkit, hcl, json, compose
aliases:
  - /build/customize/bake/file-definition/
---

`buildx bake` supports HCL, JSON and Compose file format for defining build
[groups](#group), [targets](#target) as well as [variables](#variable) and
[functions](#functions). It looks for build definition files in the current
directory in the following order:

- `docker-compose.yml`
- `docker-compose.yaml`
- `docker-bake.json`
- `docker-bake.override.json`
- `docker-bake.hcl`
- `docker-bake.override.hcl`

## Specification

Inside a bake file you can declare group, target and variable blocks to define
project specific reusable build flows.

### Target

A target reflects a single `docker build` invocation with the same options that
you would specify for `docker build`:

```hcl
# docker-bake.hcl
target "webapp-dev" {
  dockerfile = "Dockerfile.webapp"
  tags = ["docker.io/username/webapp:latest"]
}
```

```console
$ docker buildx bake webapp-dev
```

> **Note**
>
> In the case of compose files, each service corresponds to a target.
> If compose service name contains a dot it will be replaced with an underscore.

Complete list of valid target fields available for [HCL](#hcl-definition) and
[JSON](#json-definition) definitions:

| Name                | Type   | Description                                                                                                                                                                 |
| ------------------- | ------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `args`              | Map    | Set build-time variables (same as [`--build-arg` flag](../../engine/reference/commandline/buildx_build.md))                                                                 |
| `attest`            | List   | Define attestations that should be applied to the image, see [SBOM attestations](../attestations/sbom.md) and [Provenance attestations](../attestations/slsa-provenance.md) |
| `cache-from`        | List   | External cache sources (same as [`--cache-from` flag](../../engine/reference/commandline/buildx_build.md))                                                                  |
| `cache-to`          | List   | Cache export destinations (same as [`--cache-to` flag](../../engine/reference/commandline/buildx_build.md))                                                                 |
| `context`           | String | Set of files located in the specified path or URL                                                                                                                           |
| `contexts`          | Map    | Additional build contexts (same as [`--build-context` flag](../../engine/reference/commandline/buildx_build.md))                                                            |
| `dockerfile-inline` | String | Inline Dockerfile content                                                                                                                                                   |
| `dockerfile`        | String | Name of the Dockerfile (same as [`--file` flag](../../engine/reference/commandline/buildx_build.md))                                                                        |
| `inherits`          | List   | [Inherit build options](#merging-and-inheritance) from other targets                                                                                                        |
| `labels`            | Map    | Set metadata for an image (same as [`--label` flag](../../engine/reference/commandline/buildx_build.md))                                                                    |
| `no-cache-filter`   | List   | Do not cache specified stages (same as [`--no-cache-filter` flag](../../engine/reference/commandline/buildx_build.md))                                                      |
| `no-cache`          | Bool   | Do not use cache when building the image (same as [`--no-cache` flag](../../engine/reference/commandline/buildx_build.md))                                                  |
| `output`            | List   | Output destination (same as [`--output` flag](../../engine/reference/commandline/buildx_build.md))                                                                          |
| `platforms`         | List   | Set target platforms for build (same as [`--platform` flag](../../engine/reference/commandline/buildx_build.md))                                                            |
| `pull`              | Bool   | Always attempt to pull all referenced images (same as [`--pull` flag](../../engine/reference/commandline/buildx_build.md))                                                  |
| `secret`            | List   | Secret to expose to the build (same as [`--secret` flag](../../engine/reference/commandline/buildx_build.md))                                                               |
| `ssh`               | List   | SSH agent socket or keys to expose to the build (same as [`--ssh` flag](../../engine/reference/commandline/buildx_build.md))                                                |
| `tags`              | List   | Name and optionally a tag in the format `name:tag` (same as [`--tag` flag](../../engine/reference/commandline/buildx_build.md))                                             |
| `target`            | String | Set the target build stage to build (same as [`--target` flag](../../engine/reference/commandline/buildx_build.md))                                                         |

### Group

A group is a grouping of targets:

```hcl
# docker-bake.hcl
group "build" {
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

```console
$ docker buildx bake build
```

### Variable

Similar to how Terraform provides a way to [define variables](https://www.terraform.io/docs/configuration/variables.html#declaring-an-input-variable){:target="blank" rel="noopener" class=""},
the HCL file format also supports variable block definitions. These can be used
to define variables with values provided by the current environment, or a
default value when unset:

```hcl
# docker-bake.hcl
variable "TAG" {
  default = "latest"
}

target "webapp-dev" {
  dockerfile = "Dockerfile.webapp"
  tags = ["docker.io/username/webapp:${TAG}"]
}
```

```console
$ docker buildx bake webapp-dev          # will use the default value "latest"
$ TAG=dev docker buildx bake webapp-dev  # will use the TAG environment variable value
```

> **Tip**
>
> See also the [Configuring builds](configuring-build.md) page for advanced usage.

### Null values

Null values for `args` and `labels` are supported, so default sets in your
Dockerfile will be used:

```hcl
# docker-bake.hcl
variable "GO_VERSION" {
  default = null
}
target "default" {
  args = {
    GO_VERSION = GO_VERSION
  }
}
```

```dockerfile
ARG GO_VERSION="1.18"
FROM golang:${GO_VERSION}
```

```console
$ docker buildx bake --print
```

```json
{
  "target": {
    "default": {
      "context": ".",
      "dockerfile": "Dockerfile"
    }
  }
}
```

```console
$ GO_VERSION=1.19 docker buildx bake --print
```

```json
{
  "target": {
    "default": {
      "context": ".",
      "dockerfile": "Dockerfile",
      "args": {
        "GO_VERSION": "1.19"
      }
    }
  }
}
```

### Functions

A [set of generally useful functions](https://github.com/docker/buildx/blob/master/bake/hclparser/stdlib.go){:target="blank" rel="noopener" class=""}
provided by [go-cty](https://github.com/zclconf/go-cty/tree/main/cty/function/stdlib){:target="blank" rel="noopener" class=""}
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

In addition, [user defined functions](https://github.com/hashicorp/hcl/tree/main/ext/userfunc){:target="blank" rel="noopener" class=""}
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

> **Note**
>
> See [User defined HCL functions](hcl-funcs.md) page for more details.

## Built-in variables

- `BAKE_CMD_CONTEXT` can be used to access the main `context` for bake command
  from a bake file that has been [imported remotely](file-definition.md#remote-definition).
- `BAKE_LOCAL_PLATFORM` returns the current platform's default platform
  specification (e.g. `linux/amd64`).

## Merging and inheritance

Multiple files can include the same target and final build options will be
determined by merging them together:

```hcl
# docker-bake.hcl
target "webapp-dev" {
  dockerfile = "Dockerfile.webapp"
  tags = ["docker.io/username/webapp:latest"]
}
```

```hcl
# docker-bake2.hcl
target "webapp-dev" {
  tags = ["docker.io/username/webapp:dev"]
}
```

```console
$ docker buildx bake -f docker-bake.hcl -f docker-bake2.hcl webapp-dev
```

A group can specify its list of targets with the `targets` option. A target can
inherit build options by setting the `inherits` option to the list of targets or
groups to inherit from:

```hcl
# docker-bake.hcl
target "webapp-dev" {
  dockerfile = "Dockerfile.webapp"
  tags = ["docker.io/username/webapp:${TAG}"]
}

target "webapp-release" {
  inherits = ["webapp-dev"]
  platforms = ["linux/amd64", "linux/arm64"]
}
```

## `default` target/group

When you invoke `bake` you specify what targets/groups you want to build. If no
arguments is specified, the group/target named `default` will be built:

```hcl
# docker-bake.hcl
target "default" {
  dockerfile = "Dockerfile.webapp"
  tags = ["docker.io/username/webapp:latest"]
}
```

```console
$ docker buildx bake
```

## Definitions

### HCL definition

HCL definition file is recommended as its experience is more aligned with buildx UX
and also allows better code reuse, different target groups and extended features.

```hcl
# docker-bake.hcl
variable "TAG" {
  default = "latest"
}

group "default" {
  targets = ["db", "webapp-dev"]
}

target "webapp-dev" {
  dockerfile = "Dockerfile.webapp"
  tags = ["docker.io/username/webapp:${TAG}"]
}

target "webapp-release" {
  inherits = ["webapp-dev"]
  platforms = ["linux/amd64", "linux/arm64"]
}

target "db" {
  dockerfile = "Dockerfile.db"
  tags = ["docker.io/username/db"]
}
```

### JSON definition

```json
{
  "variable": {
    "TAG": {
      "default": "latest"
    }
  },
  "group": {
    "default": {
      "targets": ["db", "webapp-dev"]
    }
  },
  "target": {
    "webapp-dev": {
      "dockerfile": "Dockerfile.webapp",
      "tags": ["docker.io/username/webapp:${TAG}"]
    },
    "webapp-release": {
      "inherits": ["webapp-dev"],
      "platforms": ["linux/amd64", "linux/arm64"]
    },
    "db": {
      "dockerfile": "Dockerfile.db",
      "tags": ["docker.io/username/db"]
    }
  }
}
```

### Compose file

```yaml
# docker-compose.yml
services:
  webapp:
    image: docker.io/username/webapp:latest
    build:
      dockerfile: Dockerfile.webapp

  db:
    image: docker.io/username/db
    build:
      dockerfile: Dockerfile.db
```

> **Note**
>
> See [Building from Compose file](compose-file.md) page for more details.

## Remote definition

You can also build bake files directly from a remote Git repository or HTTPS URL:

```console
$ docker buildx bake "https://github.com/docker/cli.git#v20.10.11" --print
#1 [internal] load git source https://github.com/docker/cli.git#v20.10.11
#1 0.745 e8f1871b077b64bcb4a13334b7146492773769f7       refs/tags/v20.10.11
#1 2.022 From https://github.com/docker/cli
#1 2.022  * [new tag]         v20.10.11  -> v20.10.11
#1 DONE 2.9s
```

```json
{
  "group": {
    "default": {
      "targets": [
        "binary"
      ]
    }
  },
  "target": {
    "binary": {
      "context": "https://github.com/docker/cli.git#v20.10.11",
      "dockerfile": "Dockerfile",
      "args": {
        "BASE_VARIANT": "alpine",
        "GO_STRIP": "",
        "VERSION": ""
      },
      "target": "binary",
      "platforms": [
        "local"
      ],
      "output": [
        "build"
      ]
    }
  }
}
```

As you can see the context is fixed to `https://github.com/docker/cli.git` even if
[no context is actually defined](https://github.com/docker/cli/blob/2776a6d694f988c0c1df61cad4bfac0f54e481c8/docker-bake.hcl#L17-L26){:target="blank" rel="noopener" class=""}
in the definition.

If you want to access the main context for bake command from a bake file
that has been imported remotely, you can use the [`BAKE_CMD_CONTEXT` built-in var](#built-in-variables).

```console
$ cat https://raw.githubusercontent.com/tonistiigi/buildx/remote-test/docker-bake.hcl
```

```hcl
target "default" {
  context = BAKE_CMD_CONTEXT
  dockerfile-inline = <<EOT
FROM alpine
WORKDIR /src
COPY . .
RUN ls -l && stop
EOT
}
```

```console
$ docker buildx bake "https://github.com/tonistiigi/buildx.git#remote-test" --print
```

```json
{
  "target": {
    "default": {
      "context": ".",
      "dockerfile": "Dockerfile",
      "dockerfile-inline": "FROM alpine\nWORKDIR /src\nCOPY . .\nRUN ls -l \u0026\u0026 stop\n"
    }
  }
}
```

```console
$ touch foo bar
$ docker buildx bake "https://github.com/tonistiigi/buildx.git#remote-test"
```

```text
...
 > [4/4] RUN ls -l && stop:
#8 0.101 total 0
#8 0.102 -rw-r--r--    1 root     root             0 Jul 27 18:47 bar
#8 0.102 -rw-r--r--    1 root     root             0 Jul 27 18:47 foo
#8 0.102 /bin/sh: stop: not found
```

```console
$ docker buildx bake "https://github.com/tonistiigi/buildx.git#remote-test" "https://github.com/docker/cli.git#v20.10.11" --print
#1 [internal] load git source https://github.com/tonistiigi/buildx.git#remote-test
#1 0.429 577303add004dd7efeb13434d69ea030d35f7888       refs/heads/remote-test
#1 CACHED
```

```json
{
  "target": {
    "default": {
      "context": "https://github.com/docker/cli.git#v20.10.11",
      "dockerfile": "Dockerfile",
      "dockerfile-inline": "FROM alpine\nWORKDIR /src\nCOPY . .\nRUN ls -l \u0026\u0026 stop\n"
    }
  }
}
```

```console
$ docker buildx bake "https://github.com/tonistiigi/buildx.git#remote-test" "https://github.com/docker/cli.git#v20.10.11"
```

```text
...
 > [4/4] RUN ls -l && stop:
#8 0.136 drwxrwxrwx    5 root     root          4096 Jul 27 18:31 kubernetes
#8 0.136 drwxrwxrwx    3 root     root          4096 Jul 27 18:31 man
#8 0.136 drwxrwxrwx    2 root     root          4096 Jul 27 18:31 opts
#8 0.136 -rw-rw-rw-    1 root     root          1893 Jul 27 18:31 poule.yml
#8 0.136 drwxrwxrwx    7 root     root          4096 Jul 27 18:31 scripts
#8 0.136 drwxrwxrwx    3 root     root          4096 Jul 27 18:31 service
#8 0.136 drwxrwxrwx    2 root     root          4096 Jul 27 18:31 templates
#8 0.136 drwxrwxrwx   10 root     root          4096 Jul 27 18:31 vendor
#8 0.136 -rwxrwxrwx    1 root     root          9620 Jul 27 18:31 vendor.conf
#8 0.136 /bin/sh: stop: not found
```
