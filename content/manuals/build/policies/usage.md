---
title: Using build policies
description: Apply policies to builds and develop policies iteratively
weight: 20
---

Build policies validate inputs before builds execute. This guide covers how to
develop policies iteratively and apply them to real builds with `docker buildx
build` and `docker buildx bake`.

## Prerequisites

- **Buildx 0.31.0 or later** - Check your version: `docker buildx version`
- **BuildKit 0.26.0 or later** - Verify with: `docker buildx inspect
  --bootstrap`

If you're using Docker Desktop, ensure you're on a version that includes these
updates.

## Policy development workflow

Buildx automatically loads policies that match your Dockerfile name. When you
build with `Dockerfile`, buildx looks for `Dockerfile.rego` in the same
directory. For a file named `app.Dockerfile`, it looks for
`app.Dockerfile.rego`. See the [Advanced: Policy configuration](#advanced-policy-configuration)
section for configuration options and manual policy loading.

Writing policies is an iterative process that includes inspecting the sources
you use, the input properties associated with those sources, and writing policy
checks to verify that the properties meet your expectations.

1. Inspect your Dockerfile inputs with debug builds
2. Understand available fields with `policy eval --print`
3. Write initial policy rules based on what you found
4. Test specific sources with `policy eval`
5. Refine rules based on actual input data
6. Apply to builds with confidence

This workflow prevents surprises when applying policies to real builds.

### Viewing inputs from your Dockerfile

To see the inputs that your Dockerfile references (images, Git repos, HTTP
downloads), use debug mode with an actual build:

```console
$ BUILDX_POLICY_DEBUG=1 docker buildx build --progress=plain . 2>&1 | grep "checking policy"
```

Example output:

```text
#1 0.012 checking policy for source local://dockerfile
#1 0.026 checking policy for source docker-image://docker.io/library/alpine:3.19
#1 1.015 checking policy for source local://context
#1 1.034 checking policy for source https://github.com/gohugoio/hugo/releases/download/v0.141.0/hugo.tar.gz
```

This shows every input BuildKit evaluates against your policy. Use this to
understand what rules you need to write.

### Viewing input structure

`docker policy eval --print` shows the input JSON structure for a specified
source. This is useful for understanding what fields are available.

Explore input structure for different source types:

#### View Git repository structure

```console
$ docker buildx policy eval --print https://github.com/moby/buildkit.git
```

Output shows Git fields:

```json
{
  "git": {
    "schema": "https",
    "host": "github.com",
    "remote": "https://github.com/moby/buildkit.git"
  }
}
```

#### View HTTP download structure

```console
$ docker buildx policy eval --print https://releases.hashicorp.com/terraform/1.5.0/terraform.zip
```

Output shows HTTP fields:

```json
{
  "http": {
    "url": "https://releases.hashicorp.com/terraform/1.5.0/terraform.zip",
    "schema": "https",
    "host": "releases.hashicorp.com",
    "path": "/terraform/1.5.0/terraform.zip"
  }
}
```

#### View image structure

```console
$ docker buildx policy eval --print docker-image://alpine:{{% param example_alpine_version %}}
```

> [!NOTE]
> The `docker-image://` prefix is required for image references.

Output shows image fields:

```json
{
  "image": {
    "ref": "docker.io/library/alpine:3.19",
    "host": "docker.io",
    "repo": "alpine",
    "tag": "{{% param example_alpine_version %}}",
    "platform": "linux/arm64",
    "os": "linux",
    "arch": "arm64"
  }
}
```

For complete field documentation, see the [Input reference](./inputs.md).

### Testing policies with policy eval

The `docker buildx policy eval` command tests whether your policy allows a
specific source (local directory, Git URL, or HTTP URL). It's useful for
testing policy logic without running a full build.

> [!IMPORTANT]
> `policy eval` doesn't fetch sources, so many fields remain unresolved. For
> images, only basic reference and platform fields are available - metadata
> like `hasProvenance`, `signatures`, `labels`, and `env` are not populated.
> To test policies that check these fields, run an actual build with debug
> mode instead.

Note: `docker buildx policy eval` tests one source at a time. It doesn't parse
your Dockerfile to evaluate all inputs - for that, [run a debug
build](#viewing-inputs-from-your-dockerfile).

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

**Available fields during `policy eval` for images:**
- Reference info: `ref`, `host`, `repo`, `fullRepo`, `tag`, `isCanonical`
- Platform info: `platform`, `os`, `arch`, `variant`

**Unresolved fields (only available during builds):**
- Metadata: `checksum`, `labels`, `user`, `volumes`, `workingDir`, `env`, `createdTime`
- Attestations: `hasProvenance`, `signatures`

### Iterative development example

Here's a practical workflow for developing policies:

**Step 1:** Start with basic deny-all

```rego {title="Dockerfile.rego"}
package docker

default allow := false

allow if input.local

decision := {"allow": allow}
```

**Step 2:** See what inputs your Dockerfile uses

Use debug mode with a build attempt:

```console
$ BUILDX_POLICY_DEBUG=1 docker buildx build --progress=plain . 2>&1 | grep "checking policy"
```

Output shows:

```text
#1 0.012 checking policy for source local://dockerfile
#1 0.026 checking policy for source docker-image://docker.io/library/alpine:3.19
#1 ERROR: policy denied
```

Now you know you need an image rule for alpine.

**Step 3:** Add rule for the image

```rego
allow if {
    input.image.repo == "alpine"
}
```

**Step 4:** Build again

```console
$ docker buildx build .
```

If successful, your policy is working. If it fails, see
[Debugging](./debugging.md) for troubleshooting guidance.

## Using policies with docker build

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

Bake supports automatic policy loading just like `docker buildx build`. Place
`Dockerfile.rego` alongside your Dockerfile and run:

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
  dockerfile = "Dockerfile.dev"
  policy = ["policies/permissive.rego"]
}

target "production" {
  dockerfile = "Dockerfile.prod"
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

> [!NOTE]
> Since `policy test` is not yet implemented, test your policies in CI/CD by
> running actual builds with the `--policy` flag. This ensures policies work
> correctly in your build pipeline.

Validate policies in continuous integration:

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
intended. The `strict=true` flag fails the build if policies aren't loaded.

## Advanced: Policy configuration

This section covers advanced policy loading mechanisms and configuration
options.

### Automatic policy loading

Buildx automatically loads policies that match your Dockerfile name. When you
build with `Dockerfile`, buildx looks for `Dockerfile.rego` in the same
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
| `log-level=<level>` | Control policy logging (error, warn, info, debug, none) | `log-level=debug`             |

Combine options with commas:

```console
$ docker buildx build --policy filename=extra.rego,strict=true .
```

### Testing sources with specific policy files

When using `policy eval`, the `--filename` flag specifies which policy file to
load by providing the base Dockerfile name (without the `.rego` extension).
This is useful for testing sources against policies associated with different
Dockerfiles.

For example, to test a source against the policy for `app.Dockerfile`:

```console
$ docker buildx policy eval --filename app.Dockerfile .
```

This loads `app.Dockerfile.rego` and tests whether it allows the source `.`
(the local directory). The flag defaults to `Dockerfile` if not specified.

Use cases:

- Testing if a policy allows a specific source before building
- Validating multiple policies against the same source
- Debugging policy behavior for specific Dockerfiles in multi-Dockerfile projects

Example with different sources:

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

- Debug policy failures: [Debugging](./debugging.md)
- Browse working examples: [Example policies](./examples.md)
- Reference all input fields: [Input reference](./inputs.md)
