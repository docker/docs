---
title: Debugging build policies
linkTitle: Debugging
description: Debug policies during development with inspection and testing tools
keywords: build policies, debugging, policy troubleshooting, log-level, policy eval, rego debugging
weight: 70
---

When policies don't work as expected, use the tools available to inspect policy
evaluation and understand what's happening. This guide covers the debugging
techniques and common gotchas.

## Quick reference

Essential debugging commands:

```console
# See complete input data during builds (recommended)
$ docker buildx build --progress=plain --policy log-level=debug .

# See policy checks and decisions
$ docker buildx build --progress=plain .

# Explore input structure for different sources
$ docker buildx policy eval --print .
$ docker buildx policy eval --print https://github.com/org/repo.git
$ docker buildx policy eval --print docker-image://alpine:3.19

# Test if policy allows a source
$ docker buildx policy eval .
```

## Policy output with `--progress=plain`

To see policy evaluation during builds, use `--progress=plain`:

```console
$ docker buildx build --progress=plain .
```

This shows all policy checks, decisions, and `print()` output. Without
`--progress=plain`, policy evaluation is silent unless there's an error.

```plaintext
#1 loading policies Dockerfile.rego
#1 0.010 checking policy for source docker-image://alpine:3.19 (linux/arm64)
#1 0.011 Dockerfile.rego:8: image: {"ref":"alpine:3.19","repo":"alpine","tag":"3.19"}
#1 0.012 policy decision for source docker-image://alpine:3.19: ALLOW
```

If a policy denies a source, you'll see:

```text
#1 0.012 policy decision for source docker-image://nginx:latest: DENY
ERROR: source "docker-image://nginx:latest" not allowed by policy
```

## Debug logging

For detailed debugging, add `--policy log-level=debug` to see the full input
JSON, unresolved fields, and policy responses:

```console
$ docker buildx build --progress=plain --policy log-level=debug .
```

This shows significantly more information than the default level, including the
complete input structure for each source without needing `print()` statements
in your policy.

Complete input JSON:

```text
#1 0.007 policy input: {
#1 0.007   "env": {
#1 0.007     "filename": "."
#1 0.007   },
#1 0.007   "image": {
#1 0.007     "ref": "docker.io/library/alpine:3.19",
#1 0.007     "host": "docker.io",
#1 0.007     "repo": "alpine",
#1 0.007     "fullRepo": "docker.io/library/alpine",
#1 0.007     "tag": "3.19",
#1 0.007     "platform": "linux/arm64",
#1 0.007     "os": "linux",
#1 0.007     "arch": "arm64"
#1 0.007   }
#1 0.007 }
```

Unresolved fields:

```text
#1 0.007 unknowns for policy evaluation: [input.image.checksum input.image.labels input.image.user input.image.volumes input.image.workingDir input.image.env input.image.hasProvenance input.image.signatures]
```

Policy response:

```text
#1 0.008 policy response: map[allow:true]
```

This detailed output is invaluable for understanding exactly what data your
policy receives and which fields are not yet resolved. Use debug logging when
developing policies to avoid needing extensive `print()` statements.

## Conditional debugging with print()

While `--policy log-level=debug` shows all input data automatically, the
`print()` function is useful for debugging specific rule logic and conditional
flows:

```rego
allow if {
    input.image
    print("Checking image:", input.image.repo, "isCanonical:", input.image.isCanonical)
    input.image.repo == "alpine"
    input.image.isCanonical
}
```

Use `print()` to debug conditional logic within rules or track which rules are
evaluating. For general input inspection during development, use `--policy
log-level=debug` instead - it requires no policy modifications.

> [!NOTE]
> Print statements only execute when their containing rule evaluates. A rule
> like `allow if { input.image; print(...) }` only prints for image inputs,
> not for Git repos, HTTP downloads, or local files.

## Common issues

### Full repository path or repository name

Symptom: Policy checking repository names doesn't match as expected.

Cause: Docker Hub images use `input.image.repo` for the short name
(`"alpine"`) but `input.image.fullRepo` includes the full path
(`"docker.io/library/alpine"`).

Solution:

```rego
# Match just the repo name (works for Docker Hub and other registries)
allow if {
    input.image
    input.image.repo == "alpine"
}

# Or match the full repository path
allow if {
    input.image
    input.image.fullRepo == "docker.io/library/alpine"
}
```

### Policy evaluation happens multiple times

Symptom: Build output shows the same source evaluated multiple times.

Cause: BuildKit may evaluate policies at different stages (reference
resolution, actual pull) or for different platforms.

This is normal behavior. Policies should be idempotent (produce same result
each time for the same input).

### Fields missing with `policy eval --print`

Symptom: `docker buildx policy eval --print` doesn't show expected fields
like `hasProvenance`, `labels`, or `checksum`.

Cause: `--print` shows only reference information by default, without
fetching from registries.

Solution: Use `--fields` to fetch specific metadata fields:

```console
$ docker buildx policy eval --print --fields image.labels docker-image://alpine:3.19
```

See [Using build policies](./usage.md#testing-policies-with-policy-eval) for
details.

## Next steps

- See complete field reference: [Input reference](./inputs.md)
- Review example policies: [Examples](./examples.md)
- Learn policy usage patterns: [Using build policies](./usage.md)
