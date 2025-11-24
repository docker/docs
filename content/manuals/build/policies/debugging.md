---
title: Debugging build policies
description: Debug policies during development with inspection and testing tools
weight: 40
---

Writing policies requires understanding what inputs your policy receives and how
your rules evaluate. This guide shows practical debugging techniques for policy
development.

## Quick reference

Essential debugging commands:

```console
# Explore input structure for different sources
$ docker buildx policy eval --print .
$ docker buildx policy eval --print https://github.com/org/repo.git
$ docker buildx policy eval --print docker-image://alpine:3.19

# Test if policy allows a source
$ docker buildx policy eval .

# See all inputs from your Dockerfile
$ BUILDX_POLICY_DEBUG=1 docker buildx build --progress=plain . 2>&1 | grep "checking policy"
```

## Inspect input data

Before writing policy logic, see what data you have available.

### View all inputs from your Dockerfile

> [!IMPORTANT]
> `policy eval --print` only shows the specific source you provide - it doesn't
> parse your Dockerfile. To see **all** inputs from your Dockerfile (images, Git
> repos, HTTP downloads), use debug mode with an actual build.

Use debug mode to see every input BuildKit evaluates:

```console
$ BUILDX_POLICY_DEBUG=1 docker buildx build --progress=plain . 2>&1 | grep "checking policy"
```

For a Dockerfile with an image:

```dockerfile
FROM alpine:3.19
RUN echo "hello"
```

You see all policy checks:

```text
#1 0.012 checking policy for source local://dockerfile
#1 0.026 checking policy for source docker-image://docker.io/library/alpine:3.19
#1 1.015 checking policy for source local://context
```

This shows:
- The Dockerfile itself (`local://dockerfile`)
- The alpine image (`docker-image://...`)
- The build context (`local://context`)

To see the full input data, add `--print` to the debug output (but note this is
verbose).

### Explore input structure for specific types

To understand what fields are available for different input types, use `policy
eval --print` with specific sources:

**Explore Git input structure:**

```console
$ docker buildx policy eval --print https://github.com/moby/buildkit.git#v0.26.1
```

Shows Git fields:

```json
{
  "git": {
    "schema": "https",
    "host": "github.com",
    "remote": "https://github.com/moby/buildkit.git",
    "tagName": "v0.26.1"
  }
}
```

**Explore HTTP input structure:**

```console
$ docker buildx policy eval --print https://releases.hashicorp.com/terraform/1.5.0/terraform.zip
```

Shows HTTP fields:

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

**Explore image input structure:**

```console
$ docker buildx policy eval --print docker-image://alpine:3.19
```

Shows image fields:

```json
{
  "image": {
    "ref": "docker.io/library/alpine:3.19",
    "host": "docker.io",
    "repo": "alpine",
    "tag": "3.19",
    "os": "linux",
    "arch": "arm64"
  }
}
```

> [!IMPORTANT]
> `policy eval --print` doesn't fetch sources, so many fields show as
> "unresolved". For images, metadata fields like `hasProvenance`, `signatures`,
> `labels`, and `env` are only populated during actual builds. Use debug mode
> with a build to see complete input data.

Note: For images, use the `docker-image://` prefix.

### Check field availability

Not all fields are always populated. Fields may be unavailable for two reasons:

**1. Using `policy eval` instead of builds:**
- `policy eval` doesn't fetch sources, leaving many fields unresolved
- Image fields like `hasProvenance`, `signatures`, `labels`, `env` require builds
- Git commit/tag details require builds
- HTTP checksums require builds

**2. Sources don't contain the data:**
- Images without attestations won't have `hasProvenance` or `signatures`
- Unsigned Git commits won't have `pgpSignature`
- Images without labels won't populate the `labels` field

To see which fields are actually populated for your inputs, use debug mode
during a build and examine the full input JSON, or add `print()` statements to
your policy:

```rego
allow if {
    input.image
    print("Full image input:", input.image)
    true
}
```

Run with debug mode to see the output:

```console
$ BUILDX_POLICY_DEBUG=1 docker buildx build --progress=plain .
```

## Debug with print statements

Add `print()` calls to see what your policy evaluates.

### Basic print usage

```rego
package docker

default allow := false

allow if input.local

allow if {
    input.image
    print("Checking image:", input.image.ref)
    print("Host:", input.image.host)
    print("Repo:", input.image.repo)
    input.image.repo == "alpine"
}

decision := {"allow": allow}
```

Run the build with debug mode:

```console
$ BUILDX_POLICY_DEBUG=1 docker buildx build --progress=plain .
```

Output shows your print statements:

```text
#1 [policy] loading policy: Dockerfile.rego
#1 [policy] evaluating policy for image
#1 Dockerfile.rego:8: Checking image: alpine:3.19
#1 Dockerfile.rego:9: Host:
#1 Dockerfile.rego:10: Repo: alpine
#1 [policy] decision: {"allow":true}
```

This shows:
- Which rules evaluated
- Actual field values at evaluation time
- Whether the rule matched

### Print complex objects

Print entire objects to see their structure:

```rego
allow if {
    input.git
    print("Git object:", input.git)
    print("Commit details:", input.git.commit)
    true
}
```

Output:

```text
#1 Dockerfile.rego:5: Git object: {"host":"github.com","remote":"https://github.com/moby/buildkit.git","tagName":"v0.26.1","ref":"refs/tags/v0.26.1"}
#1 Dockerfile.rego:6: Commit details: {"author":{"name":"John Doe","email":"john@example.com","when":"2025-01-15T10:00:00Z"},"message":"Release v0.26.1"}
```

### Debug conditional logic

Print to understand which branches execute:

```rego
allow if {
    input.image
    print("Evaluating image rules")

    is_alpine := input.image.repo == "alpine"
    print("Is Alpine?", is_alpine)

    has_tag := input.image.tag != ""
    print("Has tag?", has_tag)

    is_alpine
    has_tag
}
```

Output shows the evaluation flow:

```text
#1 Dockerfile.rego:5: Evaluating image rules
#1 Dockerfile.rego:8: Is Alpine? true
#1 Dockerfile.rego:11: Has tag? true
#1 [policy] decision: {"allow":true}
```

This immediately shows which conditions pass or fail.

## Understand debug mode

Debug mode shows complete policy evaluation during builds.

### Enable debug mode

```console
$ BUILDX_POLICY_DEBUG=1 docker buildx build --progress=plain .
```

The `--progress=plain` flag is required. Debug output doesn't appear with other
progress modes.

### Debug output structure

Debug mode shows five key pieces of information:

#### 1. Policy loading

```text
#1 [policy] loading policy: Dockerfile.rego
```

Confirms which policy file loaded. If missing, buildx didn't find your policy
file.

#### 2. Input evaluation

```text
#1 [policy] evaluating policy for context
#1 [policy] evaluating policy for image
```

Shows each input being evaluated. Policies evaluate separately for:
- Local context (Dockerfile itself)
- Each image reference
- Each Git repository
- Each HTTP download

#### 3. Input data

```text
#1 [policy] input: {"image":{"ref":"alpine:3.19","repo":"alpine","tag":"3.19"}}
```

The complete input object your policy received. Use this to verify your rules
match actual data.

#### 4. Decision

```text
#1 [policy] decision: {"allow":true}
```

Whether the input passed (`true`) or failed (`false`).

#### 5. Print statements

```text
#1 Dockerfile.rego:8: Checking image: alpine:3.19
```

Any `print()` calls in your policy appear with their line numbers.

### Example debug session

Full debug output for a build:

```text
#1 [policy] loading policy: Dockerfile.rego
#1 [policy] evaluating policy for context
#1 [policy] input: {"local":{"name":"context"}}
#1 [policy] decision: {"allow":true}
#1 [policy] evaluating policy for image
#1 [policy] input: {"image":{"ref":"alpine:3.19","repo":"alpine","tag":"3.19","host":""}}
#1 Dockerfile.rego:8: Checking image: alpine:3.19
#1 [policy] decision: {"allow":true}
```

This shows:
1. Policy loaded successfully
2. Local context allowed (required)
3. Alpine image evaluated
4. Print statement executed
5. Image allowed
6. Build can proceed

When a policy fails:

```text
#1 [policy] evaluating policy for image
#1 [policy] input: {"image":{"ref":"nginx:latest","repo":"nginx","tag":"latest"}}
#1 [policy] decision: {"allow":false}
ERROR: policy violation: image nginx:latest denied
```

The input shows what was denied. Update your policy to handle this input or
confirm the denial is correct.

## Common debugging scenarios

### Policy denies everything

**Symptom**: Build fails immediately. Debug shows `{"allow":false}` for local
context.

**Solution**: Add the local input rule:

```rego
package docker

default allow := false

allow if input.local  # Required for Dockerfile access

# Your other rules...

decision := {"allow": allow}
```

Every policy needs to allow local access. Without it, buildx can't read your
Dockerfile.

### Rule doesn't match

**Symptom**: Policy should allow input but decision is `false`. Debug shows
input data looks correct.

**Cause**: Field values don't match expectations.

**Debug approach**:

1. Add print statements:

```rego
allow if {
    input.image
    print("Image host:", input.image.host)
    print("Expected host: docker.io")
    input.image.host == "docker.io"
}
```

1. Check debug output:

```text
#1 Dockerfile.rego:4: Image host:
#1 Dockerfile.rego:5: Expected host: docker.io
#1 [policy] decision: {"allow":false}
```

1. See the issue: Docker Hub uses empty host, not `"docker.io"`

1. Fix the policy:

```rego
allow if {
    input.image
    input.image.host == ""  # Docker Hub
    input.image.repo == "alpine"
}
```

### Field always empty

**Symptom**: Expected field like `hasProvenance` is always false.

**Debug approach**:

1. Add print statements to see what's available:

   ```rego
   allow if {
       input.image
       print("Full image input:", input.image)
       true
   }
   ```

1. Run with debug mode:

   ```console
   $ BUILDX_POLICY_DEBUG=1 docker buildx build --progress=plain .
   ```

1. Check the output for your image:

   ```text
   #1 Dockerfile.rego:3: Full image input: {"ref":"alpine:3.19","repo":"alpine","tag":"3.19","host":"","isCanonical":false}
   ```

1. Notice `hasProvenance` is missing

1. Understand why: Image doesn't have provenance attestations, or BuildKit
   couldn't fetch them

1. Use available fields instead:

```rego
# Instead of requiring provenance
allow if {
    input.image.hasProvenance  # May not work
}

# Use digest references as alternative
allow if {
    input.image.isCanonical  # Requires @sha256:...
}
```

### Unexpected input structure

**Symptom**: Can't access nested fields. Policy errors or unexpected behavior.

**Debug approach**:

1. Print the object:

```rego
allow if {
    input.git.commit
    print("Commit object:", input.git.commit)
    print("Author:", input.git.commit.author)
    true
}
```

1. Check structure:

   ```text
   #1 Dockerfile.rego:4: Commit object: {"author":{"name":"John","email":"john@example.com","when":"2025-01-15T10:00:00Z"}}
   #1 Dockerfile.rego:5: Author: {"name":"John","email":"john@example.com","when":"2025-01-15T10:00:00Z"}
   ```

1. Notice timestamp field is `when` not `timestamp`

1. Fix field access:

```rego
allow if {
    input.git.commit
    input.git.commit.author.when  # Correct field name
}
```

### Policy works locally but fails in CI

**Symptom**: Same Dockerfile and policy pass locally but fail in CI.

**Debug approach**:

1. Enable debug in CI:

   ```yaml
   - name: Build
     run: BUILDX_POLICY_DEBUG=1 docker buildx build --progress=plain .
   ```

1. Compare debug output from local and CI

1. Look for differences in input data:
   - Different BuildKit versions
   - Different image registries or mirrors
   - Network access differences affecting metadata

1. Make policy resilient:

```rego
# Check builder info if needed
allow if {
    input.image
    # Don't rely on metadata that may be unavailable
    input.image.isCanonical  # Digest always available
}
```

## Debugging workflow

Follow this process when developing policies:

### 1. Understand input structure

Before writing rules, familiarize yourself with the available fields:

- Review the [Input reference](./inputs.md) documentation
- Use `policy eval --print` with specific sources to explore structure:
  ```console
  $ docker buildx policy eval --print https://github.com/org/repo.git
  $ docker buildx policy eval --print https://example.com/file.tar.gz
  ```

### 2. Write initial policy

Start with basic structure:

```rego
package docker

default allow := false

allow if input.local

# Add your first rule
allow if {
    input.image
    print("Image ref:", input.image.ref)
    true  # Allow everything initially
}

decision := {"allow": allow}
```

### 3. Test policy syntax

Test if your policy allows the local context:

```console
$ docker buildx policy eval .
```

If you get "ERROR: policy denied", either:
- Your policy correctly denies the local context (check if you have `allow if input.local`)
- Your policy has a syntax error (error message will indicate this)

No output means the local context is allowed.

### 4. Add debug prints

Instrument your policy:

```rego
allow if {
    input.image
    print("Checking image:", input.image.ref)
    print("Repo:", input.image.repo)
    print("Is canonical:", input.image.isCanonical)

    # Your actual logic here
    input.image.repo == "alpine"
}
```

### 5. Build with debug mode

See evaluation in action:

```console
$ BUILDX_POLICY_DEBUG=1 docker buildx build --progress=plain .
```

Watch your print statements and decisions.

### 6. Refine logic

Based on debug output, adjust rules to match actual inputs.

### 7. Remove debug prints

Once working, remove or comment out print statements:

```rego
allow if {
    input.image
    # print("Checking image:", input.image.ref)  # Debug
    input.image.repo == "alpine"
}
```

## Quick troubleshooting checklist

When a policy isn't working:

- [ ] Policy file in same directory as Dockerfile?
- [ ] Policy file named `Dockerfile.rego` (or `<name>.Dockerfile.rego`)?
- [ ] `allow if input.local` rule present?
- [ ] Used debug build to see which inputs are being evaluated?
- [ ] Added `print()` statements to debug logic?
- [ ] Enabled debug mode (`BUILDX_POLICY_DEBUG=1 --progress=plain`)?
- [ ] Checked debug output for which input failed?
- [ ] Compared input data against your rule conditions?
- [ ] Reviewed [Input reference](./inputs.md) for available fields?
- [ ] BuildKit version 0.26.0 or later?

## Next steps

- See complete field reference: [Input reference](./inputs.md)
- Review example policies: [Examples](./examples.md)
- Learn policy usage patterns: [Using build policies](./usage.md)
