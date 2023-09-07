---
title: Advanced Bake patterns and functions
description: Learn about advanced Bake features, like user-defined functions
keywords: build, buildx, bake, buildkit, hcl
aliases:
  - /build/customize/bake/hcl-funcs/
  - /build/bake/hcl-funcs/
---

HCL functions are great for when you need to manipulate values in more complex ways than just concatenating or appending values.

The following sections contain some examples on custom functions and other
advanced use cases:

- [Interpolate environment variables](#interpolate-environment-variables)
- [Built-in functions](#built-in-functions)
- [User-defined functions](#user-defined-functions)
- [Ternary operators](#ternary-operators)
- [Variables in functions](#variables-in-functions)
- [Typed variables](#typed-variables)

## Interpolate environment variables

As shown in the [Bake file reference](reference.md#variable) page, Bake
supports variable blocks which are assigned to matching environment variables
or default values.

The following example shows how you can interpolate a `TAG` environment
variable to populate a variable in the Bake configuration.

{{< tabs >}}
{{< tab name="HCL" >}}

```hcl
# docker-bake.hcl
variable "TAG" {
  default = "latest"
}

group "default" {
  targets = ["webapp"]
}

target "webapp" {
  tags = ["docker.io/username/webapp:${TAG}"]
}
```

{{< /tab >}}
{{< tab name="JSON" >}}

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
      "tags": ["docker.io/username/webapp:${TAG}"]
    }
  }
}
```

{{< /tab >}}
{{< /tabs >}}

```console
$ docker buildx bake --print webapp
```

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
      "tags": ["docker.io/username/webapp:latest"]
    }
  }
}
```

```console
$ TAG=$(git rev-parse --short HEAD) docker buildx bake --print webapp
```

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

## Built-in functions

You can use [`go-cty` standard library functions](https://github.com/zclconf/go-cty/tree/main/cty/function/stdlib).
The following example shows the `add` function.

```hcl
# docker-bake.hcl
variable "TAG" {
  default = "latest"
}

group "default" {
  targets = ["webapp"]
}

target "webapp" {
  args = {
    buildno = "${add(123, 1)}"
  }
}
```

```console
$ docker buildx bake --print webapp
```

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
      "args": {
        "buildno": "124"
      }
    }
  }
}
```

## User-defined functions

You can create [user-defined functions](https://github.com/hashicorp/hcl/tree/main/ext/userfunc)
that do just what you want, if the built-in standard library functions don't
meet your needs.

The following example defines an `increment` function.

```hcl
# docker-bake.hcl
function "increment" {
  params = [number]
  result = number + 1
}

group "default" {
  targets = ["webapp"]
}

target "webapp" {
  args = {
    buildno = "${increment(123)}"
  }
}
```

```console
$ docker buildx bake --print webapp
```

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
      "args": {
        "buildno": "124"
      }
    }
  }
}
```

## Ternary operators

You can use ternary operators to conditionally register a value.

The following example adds a tag only when a variable is not empty, using the
`notequal` function.

```hcl
# docker-bake.hcl
variable "TAG" {default="" }

group "default" {
  targets = [
    "webapp",
  ]
}

target "webapp" {
  context="."
  dockerfile="Dockerfile"
  tags = [
    "my-image:latest",
    notequal("",TAG) ? "my-image:${TAG}": "",
  ]
}
```

```console
$ docker buildx bake --print webapp
```

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
      "tags": ["my-image:latest"]
    }
  }
}
```

## Variables in functions

You can make references to variables and standard library functions inside your
functions.

You can't reference user-defined functions from other functions.

```hcl
# docker-bake.hcl
variable "REPO" {
  default = "user/repo"
}

function "tag" {
  params = [tag]
  result = ["${REPO}:${tag}"]
}

target "webapp" {
  tags = tag("v1")
}
```

```console
$ docker buildx bake --print webapp
```

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
      "tags": ["user/repo:v1"]
    }
  }
}
```

## Typed variables

Non-string variables are supported. Values passed as environment variables are
coerced into suitable types first.

```hcl
# docker-bake.hcl
variable "FOO" {
  default = 3
}

variable "IS_FOO" {
  default = true
}

target "app" {
  args = {
    v1 = FOO > 5 ? "higher" : "lower"
    v2 = IS_FOO ? "yes" : "no"
  }
}
```

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
        "v1": "lower",
        "v2": "yes"
      }
    }
  }
}
```
