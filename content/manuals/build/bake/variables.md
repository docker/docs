---
title: Variables in Bake
linkTitle: Variables
weight: 40
description:
keywords: build, buildx, bake, buildkit, hcl, variables
---

You can define and use variables in a Bake file to set attribute values,
interpolate them into other values, and perform arithmetic operations.
Variables can be defined with default values, and can be overridden with
environment variables.

## Using variables as attribute values

Use the `variable` block to define a variable.

```hcl {title=docker-bake.hcl}
variable "TAG" {
  default = "docker.io/username/webapp:latest"
}
```

The following example shows how to use the `TAG` variable in a target.

```hcl {title=docker-bake.hcl}
target "webapp" {
  context = "."
  dockerfile = "Dockerfile"
  tags = [ TAG ]
}
```

## Interpolate variables into values

Bake supports string interpolation of variables into values. You can use the
`${}` syntax to interpolate a variable into a value. The following example
defines a `TAG` variable with a value of `latest`.

```hcl {title=docker-bake.hcl}
variable "TAG" {
  default = "latest"
}
```

To interpolate the `TAG` variable into the value of an attribute, use the
`${TAG}` syntax.

```hcl {title=docker-bake.hcl}
group "default" {
  targets = [ "webapp" ]
}

variable "TAG" {
  default = "latest"
}

target "webapp" {
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

## Validating variables

To verify that the value of a variable conforms to an expected type, value
range, or other condition, you can define custom validation rules using the
`validation` block.

In the following example, validation is used to enforce a numeric constraint on
a variable value; the `PORT` variable must be 1024 or greater.

```hcl {title=docker-bake.hcl}
# Define a variable `PORT` with a default value and a validation rule
variable "PORT" {
  default = 3000  # Default value assigned to `PORT`

  # Validation block to ensure `PORT` is a valid number within the acceptable range
  validation {
    condition = PORT >= 1024  # Ensure `PORT` is at least 1024
    error_message = "The variable 'PORT' must be 1024 or greater."  # Error message for invalid values
  }
}
```

If the `condition` expression evaluates to `false`, the variable value is
considered invalid, whereby the build invocation fails and `error_message` is
emitted. For example, if `PORT=443`, the condition evaluates to `false`, and
the error is raised.

Values are coerced into the expected type before the validation is set. This
ensures that any overrides set with environment variables work as expected.

### Validate multiple conditions

To evaluate more than one condition, define multiple `validation` blocks for
the variable. All conditions must be `true`.

Here’s an example:

```hcl {title=docker-bake.hcl}
# Define a variable `VAR` with multiple validation rules
variable "VAR" {
  # First validation block: Ensure the variable is not empty
  validation {
    condition = VAR != ""
    error_message = "The variable 'VAR' must not be empty."
  }

  # Second validation block: Ensure the value contains only alphanumeric characters
  validation {
    # VAR and the regex match must be identical:
    condition = VAR == regex("[a-zA-Z0-9]+", VAR)
    error_message = "The variable 'VAR' can only contain letters and numbers."
  }
}
```

This example enforces:

- The variable must not be empty.
- The variable must match a specific character set.

For invalid inputs like `VAR="hello@world"`, the validation would fail.

### Validating variable dependencies

You can reference other Bake variables in your condition expression, enabling
validations that enforce dependencies between variables. This ensures that
dependent variables are set correctly before proceeding.

Here’s an example:

```hcl {title=docker-bake.hcl}
# Define a variable `FOO`
variable "FOO" {}

# Define a variable `BAR` with a validation rule that references `FOO`
variable "BAR" {
  # Validation block to ensure `FOO` is set if `BAR` is used
  validation {
    condition = FOO != ""  # Check if `FOO` is not an empty string
    error_message = "The variable 'BAR' requires 'FOO' to be set."
  }
}
```

This configuration ensures that the `BAR` variable can only be used if `FOO`
has been assigned a non-empty value. Attempting to build without setting `FOO`
will trigger the validation error.

## Escape variable interpolation

If you want to bypass variable interpolation when parsing the Bake definition,
use double dollar signs (`$${VARIABLE}`).

```hcl {title=docker-bake.hcl}
target "webapp" {
  dockerfile-inline = <<EOF
  FROM alpine
  ARG TARGETARCH
  RUN echo "Building for $${TARGETARCH/amd64/x64}"
  EOF
  platforms = ["linux/amd64", "linux/arm64"]
}
```

```console
$ docker buildx bake --progress=plain
...
#8 [linux/arm64 2/2] RUN echo "Building for arm64"
#8 0.036 Building for arm64
#8 DONE 0.0s

#9 [linux/amd64 2/2] RUN echo "Building for x64"
#9 0.046 Building for x64
#9 DONE 0.1s
...
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

target "webapp" {
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
    "webapp": {
      "context": ".",
      "contexts": {
        "base": "docker.io/library/alpine:latest"
      },
      "dockerfile": "Dockerfile"
    }
  }
}
```

## Additional resources

Here are some additional resources that show how you can use variables in Bake:

- You can override `variable` values using environment variables. See
  [Overriding configurations](./overrides.md#environment-variables) for more
  information.
- You can refer to and use global variables in functions. See [HCL
  functions](./funcs.md#variables-in-functions)
- You can use variable values when evaluating expressions. See [Expression
  evaluation](./expressions.md#expressions-with-variables)
