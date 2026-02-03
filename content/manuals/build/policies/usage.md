---
title: Using build policies
linkTitle: Usage
description: Apply policies to builds and develop policies iteratively
keywords: build policies, policy eval, docker buildx, policy development, debugging
weight: 20
---

Build policies validate inputs before builds execute. This guide covers how to
develop policies iteratively and apply them to real builds with `docker buildx
build` and `docker buildx bake`.

## Prerequisites

- Buildx 0.31.0 or later - Check your version: `docker buildx version`
- BuildKit 0.26.0 or later - Verify with: `docker buildx inspect
  --bootstrap`

If you're using Docker Desktop, ensure you're on a version that includes these
updates.

## Policy development workflow

Buildx automatically loads policies that match your Dockerfile name. When you
build with `Dockerfile`, Buildx looks for `Dockerfile.rego` in the same
directory. For a file named `app.Dockerfile`, it looks for
`app.Dockerfile.rego`. See the [Advanced: Policy configuration](#advanced-policy-configuration)
section for configuration options and manual policy loading.

Writing policies is an iterative process:

1. Start with a basic deny-all policy.
2. Build with debug logging to see what inputs your Dockerfile uses.
3. Add rules to allow specific sources based on the debug output.
4. Test and refine.

### Viewing inputs from your Dockerfile

To see the inputs that your Dockerfile references (images, Git repos, HTTP
downloads), build with debug logging:

```console
$ docker buildx build --progress=plain --policy log-level=debug .
```

Example output for an image source:

```text
#1 0.010 checking policy for source docker-image://alpine:3.19 (linux/arm64)
#1 0.011 policy input: {
#1 0.011   "env": {
#1 0.011     "filename": "."
#1 0.011   },
#1 0.011   "image": {
#1 0.011     "ref": "docker.io/library/alpine:3.19",
#1 0.011     "host": "docker.io",
#1 0.011     "repo": "alpine",
#1 0.011     "tag": "3.19",
#1 0.011     "platform": "linux/arm64"
#1 0.011   }
#1 0.011 }
#1 0.011 unknowns for policy evaluation: [input.image.checksum input.image.labels ...]
#1 0.012 policy decision for source docker-image://alpine:3.19: ALLOW
```

This shows the complete input structure, which fields are unresolved, and the
policy decision for each source. See [Input reference](./inputs.md) for all
available fields.

### Testing policies with policy eval

Use [`docker buildx policy eval`](/reference/cli/docker/buildx/policy/eval/) to
test whether your policy allows a specific source without running a full build.

Note: `docker buildx policy eval` tests the source specified as the argument.
It doesn't parse your Dockerfile to evaluate all inputs - for that, [build with
--progress=plain](#viewing-inputs-from-your-dockerfile).

Test if your policy allows the local context:

```console
$ docker buildx policy eval .
```

No output means the policy allowed the source. If denied, you see:

```console
ERROR: policy denied
```

Test other sources:

```console
$ docker buildx policy eval https://example.com              # Test HTTP
$ docker buildx policy eval https://github.com/org/repo.git  # Test Git
```

By default, `--print` shows reference information parsed from the source string
(like `repo`, `tag`, `host`) without fetching from registries. To inspect
metadata that requires fetching the source (like `labels`, `checksum`, or
`hasProvenance`), specify which fields to fetch with `--fields`:

```console
$ docker buildx policy eval --print --fields image.labels docker-image://alpine:3.19
```

Multiple fields can be specified as a comma-separated list.

### Iterative development example

Here's a practical workflow for developing policies:

1. Start with basic deny-all policy:

   ```rego {title="Dockerfile.rego"}
   package docker

   default allow := false

   allow if input.local

   decision := {"allow": allow}
   ```

2. Build with debug logging to see what inputs your Dockerfile uses:

   ```console
   $ docker buildx build --progress=plain --policy log-level=debug .
   ```

   The output shows the denied image and its input structure:

   ```text
   #1 0.026 checking policy for source docker-image://docker.io/library/alpine:3.19
   #1 0.027 policy input: {
   #1 0.027   "image": {
   #1 0.027     "repo": "alpine",
   #1 0.027     "tag": "3.19",
   #1 0.027     ...
   #1 0.027   }
   #1 0.027 }
   #1 0.028 policy decision for source docker-image://alpine:3.19: DENY
   #1 ERROR: source "docker-image://alpine:3.19" not allowed by policy
   ```

3. Add a rule allowing the alpine image:

   ```rego
   allow if {
       input.image.repo == "alpine"
   }
   ```

4. Build again to verify the policy works:

   ```console
   $ docker buildx build .
   ```

If it fails, see [Debugging](./debugging.md) for troubleshooting guidance.

## Using policies with `docker build`

Once you've developed and tested your policy, apply it to real builds.

### Basic usage

Create a policy alongside your Dockerfile:

```dockerfile {title="Dockerfile"}
FROM alpine:3.19
RUN echo "hello"
```

```rego {title="Dockerfile.rego"}
package docker

default allow := false

allow if input.local

allow if {
    input.image.repo == "alpine"
}

decision := {"allow": allow}
```

Build normally:

```console
$ docker buildx build .
```

Buildx loads the policy automatically and validates the `alpine:3.19` image
before building.

### Build with different Dockerfile names

Specify the Dockerfile with `-f`:

```console
$ docker buildx build -f app.Dockerfile .
```

Buildx looks for `app.Dockerfile.rego` in the same directory.

### Build with manual policy

Add an extra policy to the automatic one:

```console
$ docker buildx build --policy filename=extra-checks.rego .
```

Both `Dockerfile.rego` (automatic) and `extra-checks.rego` (manual) must pass.

### Build without automatic policy

Use only your specified policy:

```console
$ docker buildx build --policy reset=true,filename=strict.rego .
```

## Using policies with bake

[Bake](/build/bake/) supports automatic policy loading just like `docker buildx
build`. Place `Dockerfile.rego` alongside your Dockerfile and run:

```console
$ docker buildx bake
```

### Manual policy in bake files

Specify additional policies in your `docker-bake.hcl`:

```hcl {title="docker-bake.hcl"}
target "default" {
  dockerfile = "Dockerfile"
  policy = ["extra.rego"]
}
```

The `policy` attribute takes a list of policy files. Bake loads these in
addition to the automatic `Dockerfile.rego` (if it exists).

### Multiple policies in bake

```hcl {title="docker-bake.hcl"}
target "webapp" {
  dockerfile = "Dockerfile"
  policy = [
    "shared/base-policy.rego",
    "security/image-signing.rego"
  ]
}
```

All policies must pass for the target to build successfully.

### Different policies per target

Apply different validation rules to different targets:

```hcl {title="docker-bake.hcl"}
target "development" {
  dockerfile = "dev.Dockerfile"
  policy = ["policies/permissive.rego"]
}

target "production" {
  dockerfile = "prod.Dockerfile"
  policy = ["policies/strict.rego", "policies/signing-required.rego"]
}
```

Build with the appropriate target:

```console
$ docker buildx bake development  # Uses permissive policy
$ docker buildx bake production   # Uses strict policies
```

### Bake with policy options

Currently, bake doesn't support policy options (reset, strict, disabled) in the
HCL file. Use command-line flags instead:

```console
$ docker buildx bake --policy disabled=true production
```

## Testing in CI/CD

Validate policies in continuous integration by running builds with the `--policy` flag. For unit testing policies before running builds, see [Test build policies](./testing.md).

Test policies during CI builds:

```yaml {title=".github/workflows/test-policies.yml"}
name: Test Build Policies
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: docker/setup-buildx-action@v3
      - name: Test build with policy
        run: docker buildx build --policy strict=true .
```

This ensures policy changes don't break builds and that new rules work as
intended. The `strict=true` flag fails the build if policies aren't loaded (for
example, if the BuildKit instance used by the build is too old and doesn't
support policies).

## Advanced: Policy configuration

This section covers advanced policy loading mechanisms and configuration
options.

### Automatic policy loading

Buildx automatically loads policies that match your Dockerfile name. When you
build with `Dockerfile`, Buildx looks for `Dockerfile.rego` in the same
directory. For a file named `app.Dockerfile`, it looks for
`app.Dockerfile.rego`.

```text
project/
├── Dockerfile
├── Dockerfile.rego          # Loaded automatically for Dockerfile
├── app.Dockerfile
├── app.Dockerfile.rego      # Loaded automatically for app.Dockerfile
└── src/
```

This automatic loading means you don't need command-line flags in most cases.
Create the policy file alongside your Dockerfile and build:

```console
$ docker buildx build .
```

Buildx detects `Dockerfile.rego` and evaluates it before running the build.

> [!NOTE]
> Policy files must be in the same directory as the Dockerfile they validate.
> Buildx doesn't search parent directories or subdirectories.

### When policies don't load

If buildx can't find a matching `.rego` file, the build proceeds without policy
evaluation. To require policies and fail if none are found, use strict mode:

```console
$ docker buildx build --policy strict=true .
```

This fails the build if no policy loads or if the BuildKit daemon doesn't
support policies.

### Manual policy configuration

The `--policy` flag lets you specify additional policies, override automatic
loading, or control policy behavior.

Basic syntax:

```console
$ docker buildx build --policy filename=custom.rego .
```

This loads `custom.rego` in addition to the automatic `Dockerfile.rego` (if it
exists).

Multiple policies:

```console
$ docker buildx build --policy filename=policy1.rego --policy filename=policy2.rego .
```

All policies must pass for the build to succeed. Use this to enforce layered
requirements (base policy + project-specific rules).

Available options:

| Option              | Description                                             | Example                       |
| ------------------- | ------------------------------------------------------- | ----------------------------- |
| `filename=<path>`   | Load policy from specified file                         | `filename=custom.rego`        |
| `reset=true`        | Ignore automatic policies, use only specified ones      | `reset=true`                  |
| `disabled=true`     | Disable all policy evaluation                           | `disabled=true`               |
| `strict=true`       | Fail if BuildKit doesn't support policies               | `strict=true`                 |
| `log-level=<level>` | Control policy logging (error, warn, info, debug, none). Use `debug` to see complete input JSON and unresolved fields | `log-level=debug`             |

Combine options with commas:

```console
$ docker buildx build --policy filename=extra.rego,strict=true .
```

### Exploring sources with policy eval

The `docker buildx policy eval` command lets you quickly explore and test
sources without running a build.

#### Inspect input structure with --print

Use `--print` to see the input structure for any source without running policy
evaluation:

```console
$ docker buildx policy eval --print https://github.com/moby/buildkit.git
```

```json
{
  "git": {
    "schema": "https",
    "host": "github.com",
    "remote": "https://github.com/moby/buildkit.git"
  }
}
```

Test different source types:

```console
# HTTP downloads
$ docker buildx policy eval --print https://releases.hashicorp.com/terraform/1.5.0/terraform.zip

# Images (requires docker-image:// prefix)
$ docker buildx policy eval --print docker-image://alpine:3.19

# Local context
$ docker buildx policy eval --print .
```

Shows information parsed from the source without fetching. Use `--fields` to
fetch specific metadata (see [above](#testing-policies-with-policy-eval)).

#### Test with specific policy files

The `--filename` flag specifies which policy file to load by providing the base
Dockerfile name (without the `.rego` extension). This is useful for testing
sources against policies associated with different Dockerfiles.

For example, to test a source against the policy for `app.Dockerfile`:

```console
$ docker buildx policy eval --filename app.Dockerfile .
```

This loads `app.Dockerfile.rego` and tests whether it allows the source `.`
(the local directory). The flag defaults to `Dockerfile` if not specified.

Test different sources against your policy:

```console
$ docker buildx policy eval --filename app.Dockerfile https://github.com/org/repo.git
$ docker buildx policy eval --filename app.Dockerfile docker-image://alpine:3.19
```

### Reset automatic loading

To use only your specified policies and ignore automatic `.rego` files:

```console
$ docker buildx build --policy reset=true,filename=custom.rego .
```

This skips `Dockerfile.rego` and loads only `custom.rego`.

### Disable policies temporarily

Disable policy evaluation for testing or emergencies:

```console
$ docker buildx build --policy disabled=true .
```

The build proceeds without any policy checks. Use this carefully - you're
bypassing security controls.

## Next steps

- Write unit tests for your policies: [Test build policies](./testing.md)
- Debug policy failures: [Debugging](./debugging.md)
- Browse working examples: [Example policies](./examples.md)
- Reference all input fields: [Input reference](./inputs.md)
