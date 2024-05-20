---
title: Configure builds with Bake
description: Learn how to create a flexible build configuration with Bake
keywords: build, buildx, bake, buildkit, hcl, json
aliases:
  - /build/customize/bake/configuring-build/
---

Bake supports loading build definitions from files, but sometimes you need even
more flexibility to configure these definitions.

For this use case, you can define variables inside the bake files that can be
set by the user with environment variables or by [attribute definitions](#global-scope-attributes)
in other bake files. If you wish to change a specific value for a single
invocation you can use the `--set` flag [from the command line](#from-command-line).

## Global scope attributes

You can define global scope attributes in HCL/JSON and use them for code reuse
and setting values for variables. This means you can do a "data-only" HCL file
with the values you want to set/override and use it in the list of regular
output files.

```hcl
# docker-bake.hcl
variable "FOO" {
  default = "abc"
}

target "app" {
  args = {
    v1 = "pre-${FOO}"
  }
}
```

You can use this file directly:

```console
$ docker buildx bake --print app
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
        "v1": "pre-abc"
      }
    }
  }
}
```

Or create an override configuration file:

```hcl
# env.hcl
WHOAMI="myuser"
FOO="def-${WHOAMI}"
```

And invoke bake together with both of the files:

```console
$ docker buildx bake -f docker-bake.hcl -f env.hcl --print app
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
        "v1": "pre-def-myuser"
      }
    }
  }
}
```

## Resource interpolation

You can also refer to attributes defined as part of other targets, to help
reduce duplication between targets.

```hcl
# docker-bake.hcl
target "foo" {
  dockerfile = "${target.foo.name}.Dockerfile"
  tags       = [target.foo.name]
}
target "bar" {
  dockerfile = "${target.foo.name}.Dockerfile"
  tags       = [target.bar.name]
}
```

You can use this file directly:

```console
$ docker buildx bake --print foo bar
```

```json
{
  "group": {
    "default": {
      "targets": ["foo", "bar"]
    }
  },
  "target": {
    "foo": {
      "context": ".",
      "dockerfile": "foo.Dockerfile",
      "tags": ["foo"]
    },
    "bar": {
      "context": ".",
      "dockerfile": "foo.Dockerfile",
      "tags": ["bar"]
    }
  }
}
```

## From command line

You can also override target configurations from the command line with the
[`--set` flag](../../reference/cli/docker/buildx/bake.md#set):

```hcl
# docker-bake.hcl
target "app" {
  args = {
    mybuildarg = "foo"
  }
}
```

```console
$ docker buildx bake --set app.args.mybuildarg=bar --set app.platform=linux/arm64 app --print
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

Pattern matching syntax defined in [https://golang.org/pkg/path/#Match](https://golang.org/pkg/path/#Match)
is also supported:

```console
$ docker buildx bake --set foo*.args.mybuildarg=value  # overrides build arg for all targets starting with "foo"
$ docker buildx bake --set *.platform=linux/arm64      # overrides platform for all targets
$ docker buildx bake --set foo*.no-cache               # bypass caching only for targets starting with "foo"
```

Complete list of overridable fields:

- `args`
- `cache-from`
- `cache-to`
- `context`
- `dockerfile`
- `labels`
- `no-cache`
- `output`
- `platform`
- `pull`
- `secrets`
- `ssh`
- `tags`
- `target`

## Using variables in variables across files

When multiple files are specified, one file can use variables defined in
another file.

```hcl
# docker-bake1.hcl
variable "FOO" {
  default = upper("${BASE}def")
}

variable "BAR" {
  default = "-${FOO}-"
}

target "app" {
  args = {
    v1 = "pre-${BAR}"
  }
}
```

```hcl
# docker-bake2.hcl
variable "BASE" {
  default = "abc"
}

target "app" {
  args = {
    v2 = "${FOO}-post"
  }
}
```

```console
$ docker buildx bake -f docker-bake1.hcl -f docker-bake2.hcl --print app
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
        "v1": "pre--ABCDEF-",
        "v2": "ABCDEF-post"
      }
    }
  }
}
```
