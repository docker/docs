---
title: Overriding configurations
description: Learn how to override configurations in Bake files to build with different attributes.
keywords: build, buildx, bake, buildkit, hcl, json, overrides, configuration
aliases:
  - /build/bake/configuring-build/
---

Bake supports loading build definitions from files, but sometimes you need even
more flexibility to configure these definitions. For example, you might want to
override an attribute when building in a particular environment or for a
specific target.

To override target attributes or the variables used to configure them, use:

- [File overrides](#file-overrides)
- [CLI overrides](#command-line)
- [Variable overrides](#variables)
- [Environment variable overrides](#environment-variables)

## File overrides

You can load multiple Bake files that define build configurations for your
targets. This is useful when you want to separate configurations into different
files for better organization, or to conditionally override configurations
based on which files are loaded.

### Default file lookup

You can use the `--file` or `-f` flag to specify which files to load.
If you don't specify any files, Bake will use the following lookup order:

1. `compose.yaml`
2. `compose.yml`
3. `docker-compose.yml`
4. `docker-compose.yaml`
5. `docker-bake.json`
6. `docker-bake.hcl`
7. `docker-bake.override.json`
8. `docker-bake.override.hcl`

If more than one Bake file is found, all files are loaded and merged into a
single definition. Files are merged according to the lookup order.

```console
$ docker buildx bake --print
[+] Building 0.0s (1/1) FINISHED                                                                                                                                                                                            
 => [internal] load local bake definitions                                                                                                                                                                             0.0s
 => => reading compose.yaml 45B / 45B                                                                                                                                                                                  0.0s
 => => reading docker-bake.hcl 113B / 113B                                                                                                                                                                             0.0s
 => => reading docker-bake.override.hcl 65B / 65B
```

If merged files contain duplicate attribute definitions, those definitions are
either merged or overridden by the last occurrence, depending on the attribute.

Bake will attempt to load all of the files in the order they are found. If
multiple files define the same target, attributes are either merged or
overridden. In the case of overrides, the last one loaded takes precedence.

For example, given the following files:

```hcl {title=docker-bake.hcl}
variable "TAG" {
  default = "foo"
}

target "default" {
  tags = ["username/my-app:${TAG}"]
}
```

```hcl {title=docker-bake.override.hcl}
variable "TAG" {
  default = "bar"
}
```

Since `docker-bake.override.hcl` is loaded last in the default lookup order,
the `TAG` variable is overridden with the value `bar`.

```console
$ docker buildx bake --print
{
  "target": {
    "default": {
      "context": ".",
      "dockerfile": "Dockerfile",
      "tags": ["username/my-app:bar"]
    }
  }
}
```

### Manual file overrides

You can use the `--file` flag to explicitly specify which files to load,
and use this as a way to conditionally apply override files.

For example, you can create a file that defines a set of configurations for a
specific environment, and load it only when building for that environment. The
following example shows how to load an `override.hcl` file that sets the `TAG`
variable to `bar`. The `TAG` variable is then used in the `default` target.

```hcl {title=docker-bake.hcl}
variable "TAG" {
  default = "foo"
}

target "default" {
  tags = ["username/my-app:${TAG}"]
}
```

```hcl {title=overrides.hcl}
variable "TAG" {
  default = "bar"
}
```

Printing the build configuration without the `--file` flag shows the `TAG`
variable is set to the default value `foo`.

```console
$ docker buildx bake --print
{
  "target": {
    "default": {
      "context": ".",
      "dockerfile": "Dockerfile",
      "tags": [
        "username/my-app:foo"
      ]
    }
  }
}
```

Using the `--file` flag to load the `overrides.hcl` file overrides the `TAG`
variable with the value `bar`.

```console
$ docker buildx bake -f docker-bake.hcl -f overrides.hcl --print
{
  "target": {
    "default": {
      "context": ".",
      "dockerfile": "Dockerfile",
      "tags": [
        "username/my-app:bar"
      ]
    }
  }
}
```

## Command line

You can also override target configurations from the command line with the
[`--set` flag](/reference/cli/docker/buildx/bake/#set):

```hcl
# docker-bake.hcl
target "app" {
  args = {
    mybuildarg = "foo"
  }
}
```

```console
$ docker buildx bake --set app.arg.mybuildarg=bar \
    --set app.platform=linux/arm64 app --print
```

```json
{
  "group": {
    "default": {
      "targets": ["app"]
    }
  },
  "target": {
    "app": {
      "context": ".",
      "dockerfile": "Dockerfile",
      "args": {
        "mybuildarg": "bar"
      },
      "platforms": ["linux/arm64"]
    }
  }
}
```

> [!NOTE]
>
> `--set` is a repeatable flag. For array fields such as `tag`, repeat `--set`
> to provide multiple values or use `+=` to append without replacing.
> An empty assignment, such as `--set app.tag=`, clears the array. Subsequent
> assignments add values to the cleared array. Array literal syntax such as
> `--set app.tag=[a,b]` is not supported.

Override keys use singular names such as `arg`, `secret`, and `tag`. These
names apply to `--set` selectors; Bake file attributes keep their existing
names, such as `args` and `tags`. Older plural override spellings remain
supported as aliases.

`context` and `contexts` aren't aliases: `context` sets the main build
context, while `contexts.<name>` sets one named context. The alias `secrets`
applies to the secret list override only; `secrets.<id>` is invalid.

You can match multiple targets with a pattern. Quote the override so your
shell doesn't expand the pattern before Bake reads it:

```console
$ docker buildx bake --set 'foo*.arg.mybuildarg=value'
$ docker buildx bake --set '*.platform=linux/arm64'
$ docker buildx bake --set 'foo*.no-cache'
```

The first and third commands affect targets whose names start with `foo`.
The second sets the platform for all targets. For the full list of override
keys, aliases, pattern syntax, and fields supporting `+=`, see the
[`--set` reference](/reference/cli/docker/buildx/bake/#set).

### Override list attributes

Attributes such as `secret`, `ssh`, and `output` can use lists of objects in a
Bake file. In a `--set` value, use the inline syntax of the corresponding
build flag, rather than an HCL object or array. For example, `secret` uses the
same `id=...,src=...` syntax as
[`build --secret`](/reference/cli/docker/buildx/build/#secret).

Use `=` to replace the secret list, `+=` to add a secret, or an empty value
to clear the list:

```console
$ docker buildx bake --set app.secret=id=aws,src=./credentials app
$ docker buildx bake --set app.secret+=id=gcp,env=GCP_CREDENTIALS app
$ docker buildx bake --set app.secret= app
```

These operations affect the list as a whole. To change only the source of one
declared secret while preserving the others, use `secret.<id>`.

### Override a secret source

Use `--set <target>.secret.<id>=<source>` to change the source of an existing
secret without replacing the target's secret list. The secret ID must already
be declared by the target.

For example, declare an `aws` secret in the Bake file:

```hcl {title=docker-bake.hcl}
target "app" {
  secret = [{
    id  = "aws"
    src = "./credentials"
  }]
}
```

Supply its value from an environment variable or a different file:

```console
$ docker buildx bake --set app.secret.aws=env=AWS_CREDENTIALS app
$ docker buildx bake --set app.secret.aws=src=/path/to/credentials app
```

The build still receives the secret as `aws`. Other declared secrets keep
their sources. The override fails if the target doesn't declare `aws`, or if
an `id=` in the override value names a different secret. This form doesn't
support `+=`.

Named source overrides are applied after secret list overrides, so the ID
must exist in the resulting list. The value describes the complete source
using `env=...` or `src=...`; it doesn't address an object field such as
`secret.aws.env`.

### Override resource limits

Use `--set <target>.resource.<name>=<value>` to change a resource limit for
the containers that run build steps. For example, override memory and CPU
limits for an `app` target:

```console
$ docker buildx bake --set app.resource.memory=2g \
    --set app.resource.cpu-quota=50000 app
```

Supported names are `memory`, `memory-swap`, `cpu-shares`, `cpu-period`,
`cpu-quota`, `cpuset-cpus`, and `cpuset-mems`. Each override changes only the
named limit, preserving other limits from the Bake definition. Resource
overrides don't support `+=`.

These limits require a BuildKit daemon with per-step resource limit support
and take effect on Linux. They control resources within a build step, rather
than how many Bake targets run concurrently. See the
[`--resource` reference](/reference/cli/docker/buildx/build/#resource) for
units and supported values, and the
[`resources` attribute](./reference.md#targetresources) for file configuration.

## Variables

Use `--var NAME=value` to override a declared Bake variable for a single
invocation. For example, set `TAG` without exporting it in your shell:

```console
$ docker buildx bake --var TAG=dev webapp
```

The value takes precedence over an environment variable with the same name.
Unlike `--set webapp.tags=...`, this changes the variable wherever the Bake
definition uses it. See
[Override variable values](./variables.md#override-variable-values) for
precedence, types, and repeated flags.

## Environment variables

You can also use environment variables to override configurations.

Bake lets you use environment variables to override the value of a `variable`
block. Only `variable` blocks can be overridden with environment variables.
This means you need to define the variables in the bake file and then set the
environment variable with the same name to override it.

To disable these lookups, set `BUILDX_BAKE_DISABLE_VARS_ENV_LOOKUP=1`; see the
[variable reference](../building/variables.md#buildx_bake_disable_vars_env_lookup).
You can still provide variable values explicitly with `--var NAME=value`.

The following example shows how you can define a `TAG` variable with a default
value in the Bake file, and override it with an environment variable.

```hcl
variable "TAG" {
  default = "latest"
}

target "default" {
  context = "."
  dockerfile = "Dockerfile"
  tags = ["docker.io/username/webapp:${TAG}"]
}
```

```console
$ export TAG=$(git rev-parse --short HEAD)
$ docker buildx bake --print webapp
```

The `TAG` variable is overridden with the value of the environment variable,
which is the short commit hash generated by `git rev-parse --short HEAD`.

```json
{
  "group": {
    "default": {
      "targets": ["webapp"]
    }
  },
  "target": {
    "webapp": {
      "context": ".",
      "dockerfile": "Dockerfile",
      "tags": ["docker.io/username/webapp:985e9e9"]
    }
  }
}
```

### Type coercion

Overriding non-string variables with environment variables is supported. Values
passed as environment variables are coerced into suitable types first.

The following example defines a `PORT` variable. The `backend` target uses the
`PORT` variable as-is, and the `frontend` target uses the value of `PORT`
incremented by one.

```hcl
variable "PORT" {
  default = 3000
}

group "default" {
  targets = ["backend", "frontend"]
}

target "backend" {
  args = {
    PORT = PORT
  }
}

target "frontend" {
  args = {
    PORT = add(PORT, 1)
  }
}
```

Overriding `PORT` using an environment variable will first coerce the value
into the expected type, an integer, before the expression in the `frontend`
target runs.

```console
$ PORT=7070 docker buildx bake --print
```

```json
{
  "group": {
    "default": {
      "targets": [
        "backend",
        "frontend"
      ]
    }
  },
  "target": {
    "backend": {
      "context": ".",
      "dockerfile": "Dockerfile",
      "args": {
        "PORT": "7070"
      }
    },
    "frontend": {
      "context": ".",
      "dockerfile": "Dockerfile",
      "args": {
        "PORT": "7071"
      }
    }
  }
}
```
