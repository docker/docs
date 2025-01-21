---
title: Expression evaluation in Bake
linkTitle: Expressions
weight: 50
description: Learn about advanced Bake features, like user-defined functions
keywords: build, buildx, bake, buildkit, hcl, expressions, evaluation, math, arithmetic, conditionals
aliases:
  - /build/bake/advanced/
---

Bake files in the HCL format support expression evaluation, which lets you
perform arithmetic operations, conditionally set values, and more.

## Arithmetic operations

You can perform arithmetic operations in expressions. The following example
shows how to multiply two numbers.

```hcl {title=docker-bake.hcl}
sum = 7*6

target "default" {
  args = {
    answer = sum
  }
}
```

Printing the Bake file with the `--print` flag shows the evaluated value for
the `answer` build argument.

```console
$ docker buildx bake --print
```

```json
{
  "target": {
    "default": {
      "context": ".",
      "dockerfile": "Dockerfile",
      "args": {
        "answer": "42"
      }
    }
  }
}
```

## Ternary operators

You can use ternary operators to conditionally register a value.

The following example adds a tag only when a variable is not empty, using the
built-in `notequal` [function](./funcs.md).

```hcl {title=docker-bake.hcl}
variable "TAG" {}

target "default" {
  context="."
  dockerfile="Dockerfile"
  tags = [
    "my-image:latest",
    notequal("",TAG) ? "my-image:${TAG}": "",
  ]
}
```

In this case, `TAG` is an empty string, so the resulting build configuration
only contains the hard-coded `my-image:latest` tag.

```console
$ docker buildx bake --print
```

```json
{
  "target": {
    "default": {
      "context": ".",
      "dockerfile": "Dockerfile",
      "tags": ["my-image:latest"]
    }
  }
}
```

## Expressions with variables

You can use expressions with [variables](./variables.md) to conditionally set
values, or to perform arithmetic operations.

The following example uses expressions to set values based on the value of
variables. The `v1` build argument is set to "higher" if the variable `FOO` is
greater than 5, otherwise it is set to "lower". The `v2` build argument is set
to "yes" if the `IS_FOO` variable is true, otherwise it is set to "no".

```hcl {title=docker-bake.hcl}
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

Printing the Bake file with the `--print` flag shows the evaluated values for
the `v1` and `v2` build arguments.

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
