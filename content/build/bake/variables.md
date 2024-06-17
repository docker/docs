---
title: Variables in Bake
description: 
keywords: build, buildx, bake, buildkit, hcl, variables
---

You can define and use variables in a Bake file to set attribute values,
interpolate them into other values, and perform arithmetic operations.
Variables can be defined with default values, and can be overridden with
environment variables.

## Using variables as attribute values

Use the `variable` block to define a variable.

```hcl
variable "TAG" {
  default = "docker.io/username/webapp:latest"
}
```

The following example shows how to use the `TAG` variable in a target.

```hcl
target "default" {
  context = "."
  dockerfile = "Dockerfile"
  tags = [ TAG ]
}
```

## Interpolate variables into values

Bake supports string interpolation of variables into values. You can use the
`${}` syntax to interpolate a variable into a value. The following example
defines a `TAG` variable with a value of `latest`.

```hcl
variable "TAG" {
  default = "latest"
}
```

To interpolate the `TAG` variable into the value of an attribute, use the
`${TAG}` syntax.

```hcl
target "default" {
  context = "."
  dockerfile = "Dockerfile"
  tags = ["docker.io/username/webapp:${TAG}"]
}
```

Printing the Bake file with the `--print` flag shows the interpolated value in
the resolved build configuration.

```console
$ docker buildx bake --print
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

## Using variables in variables across files

When multiple files are specified, one file can use variables defined in
another file. In the following example, the `vars.hcl` file defines a
`BASE_IMAGE` variable with a default value of `docker.io/library/alpine`.

```hcl {title=vars.hcl}
variable "BASE_IMAGE" {
  default = "docker.io/library/alpine"
}
```

The following `docker-bake.hcl` file defines a `BASE_LATEST` variable that
references the `BASE_IMAGE` variable.

```hcl {title=docker-bake.hcl}
variable "BASE_LATEST" {
  default = "${BASE_IMAGE}:latest"
}

target "default" {
  contexts = {
    base = BASE_LATEST
  }
}
```

When you print the resolved build configuration, using the `-f` flag to specify
the `vars.hcl` and `docker-bake.hcl` files, you see that the `BASE_LATEST`
variable is resolved to `docker.io/library/alpine:latest`.

```console
$ docker buildx bake -f vars.hcl -f docker-bake.hcl --print app
```

```json
{
  "target": {
    "default": {
      "context": ".",
      "contexts": {
        "base": "docker.io/library/alpine:latest"
      },
      "dockerfile": "Dockerfile"
    }
  }
}
```
