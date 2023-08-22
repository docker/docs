# Bake file reference

The Bake file is a file for defining workflows that you run using `docker buildx bake`.

## File format

You can define your Bake file in the following file formats:

- HashiCorp Configuration Language (HCL)
- JSON
- YAML (Compose file)

By default, Bake uses the following lookup order to find the configuration file:

1. `docker-bake.override.hcl`
2. `docker-bake.hcl`
3. `docker-bake.override.json`
4. `docker-bake.json`
5. `docker-compose.yaml`
6. `docker-compose.yml`

Bake searches for the file in the current working directory.
You can specify the file location explicitly using the `--file` flag:

```console
$ docker buildx bake --file=../docker/bake.hcl --print
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
| ----------------------------------------------- | ------- | -------------------------------------------------------------------- |
| [`args`](#targetargs)                           | Map     | Build arguments                                                      |
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
| [`ssh`](#targetssh)                             | List    | SSH agent sockets or keys to expose to the build                     |
| [`tags`](#targettags)                           | List    | Image names and tags                                                 |
| [`target`](#targettarget)                       | String  | Target build stage                                                   |

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

> **Note**
>
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
  pull = "always"
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
RUN --mount=type=secret,id=KUBECONFIG \
    KUBECONFIG=$(cat /run/secrets/KUBECONFIG) helm upgrade --install
```

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

> **Note**
>
> See [User defined HCL functions][hcl-funcs] page for more details.

<!-- external links -->

[attestations]: https://docs.docker.com/build/attestations/
[bake_stdlib]: https://github.com/docker/buildx/blob/master/bake/hclparser/stdlib.go
[build-arg]: https://docs.docker.com/engine/reference/commandline/build/#build-arg
[build-context]: https://docs.docker.com/engine/reference/commandline/buildx_build/#build-context
[cache-backends]: https://docs.docker.com/build/cache/backends/
[cache-from]: https://docs.docker.com/engine/reference/commandline/buildx_build/#cache-from
[cache-to]: https://docs.docker.com/engine/reference/commandline/buildx_build/#cache-to
[context]: https://docs.docker.com/engine/reference/commandline/buildx_build/#build-context
[file]: https://docs.docker.com/engine/reference/commandline/build/#file
[go-cty]: https://github.com/zclconf/go-cty/tree/main/cty/function/stdlib
[hcl-funcs]: https://docs.docker.com/build/bake/hcl-funcs/
[output]: https://docs.docker.com/engine/reference/commandline/buildx_build/#output
[platform]: https://docs.docker.com/engine/reference/commandline/buildx_build/#platform
[run_mount_secret]: https://docs.docker.com/engine/reference/builder/#run---mounttypesecret
[secret]: https://docs.docker.com/engine/reference/commandline/buildx_build/#secret
[ssh]: https://docs.docker.com/engine/reference/commandline/buildx_build/#ssh
[tag]: https://docs.docker.com/engine/reference/commandline/build/#tag
[target]: https://docs.docker.com/engine/reference/commandline/build/#target
[userfunc]: https://github.com/hashicorp/hcl/tree/main/ext/userfunc
