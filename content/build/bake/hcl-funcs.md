---
title: "User defined HCL functions"
keywords: build, buildx, bake, buildkit, hcl
aliases:
- /build/customize/bake/hcl-funcs/
---

## Using interpolation to tag an image with the git sha

As shown in the [File definition](file-definition.md#variable) page, `bake`
supports variable blocks which are assigned to matching environment variables
or default values:

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

alternatively, in json format:

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

```console
$ docker buildx bake --print webapp
```

```json
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
        "docker.io/username/webapp:latest"
      ]
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
        "docker.io/username/webapp:985e9e9"
      ]
    }
  }
}
```

## Using the `add` function

You can use [`go-cty` stdlib functions](https://github.com/zclconf/go-cty/tree/main/cty/function/stdlib){:target="blank" rel="noopener" class=""}.
Here we are using the `add` function.

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
      "targets": [
        "webapp"
      ]
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

## Defining an `increment` function

It also supports [user defined functions](https://github.com/hashicorp/hcl/tree/main/ext/userfunc){:target="blank" rel="noopener" class=""}.
The following example defines a simple an `increment` function.

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
      "targets": [
        "webapp"
      ]
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

## Only adding tags if a variable is not empty using an `notequal`

Here we are using the conditional `notequal` function which is just for
symmetry with the `equal` one.

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
        "my-image:latest"
      ]
    }
  }
}
```

## Using variables in functions

You can refer variables to other variables like the target blocks can. Stdlib
functions can also be called but user functions can't at the moment.

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
        "user/repo:v1"
      ]
    }
  }
}
```

## Using typed variables

Non-string variables are also accepted. The value passed with env is parsed
into suitable type first.

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
      "targets": [
        "app"
      ]
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
