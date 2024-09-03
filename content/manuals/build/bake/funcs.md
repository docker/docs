---
title: Functions
weight: 60
description: Learn about built-in and user-defined HCL functions with Bake
keywords: build, buildx, bake, buildkit, hcl, functions, user-defined, built-in, custom, gocty
aliases:
  - /build/customize/bake/hcl-funcs/
  - /build/bake/hcl-funcs/
---

HCL functions are great for when you need to manipulate values in your build
configuration in more complex ways than just concatenation or interpolation.

## Standard library

Bake ships with built-in support for the [`go-cty` standard library functions](https://github.com/zclconf/go-cty/tree/main/cty/function/stdlib).
The following example shows the `add` function.

```hcl {title=docker-bake.hcl}
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

```hcl {title=docker-bake.hcl}
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

## Variables in functions

You can make references to [variables](./variables) and standard library
functions inside your functions.

You can't reference user-defined functions from other functions.

The following example uses a global variable (`REPO`) in a custom function.

```hcl {title=docker-bake.hcl}
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

Printing the Bake file with the `--print` flag shows that the `tag` function
uses the value of `REPO` to set the prefix of the tag.

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
