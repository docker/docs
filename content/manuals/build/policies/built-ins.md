---
title: Built-in functions
description: Buildx includes built-in helper functions to make writing policies easier
weight: 60
---

<!-- TODO: The "built-in functions" discussed on this page are built-ins
provided by buildx. NOT TO BE CONFUSED WITH built-ins that exist in rego. This
must be clarified. -->

Buildx provides built-in functions for common security checks like signature
verification, attestation validation, and version comparison. These functions
simplify policy writing by handling complex cryptographic operations.

## Image verification functions

### `trusted_github_builder(repo)`

Verifies that an image was built by GitHub Actions using GitHub's trusted
builder infrastructure.

**Parameters:**

- `repo` (string) - GitHub repository in `owner/repo` format

**Returns:** Boolean - `true` if the image has valid GitHub Actions build
attestations

**Example:**

```rego
# Require images are built by trusted GitHub Actions
allow if {
    input.image.repo == "tonistiigi/xx"
    trusted_github_builder("tonistiigi/xx")
}
```

This function checks that:

- The image has provenance attestations
- The attestations were created by GitHub Actions
- The build came from the specified repository's workflows

Use this to ensure images come from official GitHub Actions builds rather than
being pushed manually or built elsewhere.

### `trusted_github_builder_ref(repo, ref)`

Like `trusted_github_builder()`, but also validates the image was built from a
specific Git reference.

**Parameters:**

- `repo` (string) - GitHub repository in `owner/repo` format
- `ref` (string) - Git reference (branch, tag, or SHA)

**Returns:** Boolean - `true` if the image was built by GitHub Actions from the
specified ref

**Example:**

```rego
# Only allow images built from the main branch
allow if {
    input.image.repo == "myorg/myapp"
    trusted_github_builder_ref("myorg/myapp", "refs/heads/main")
}
```

Use this to restrict images to specific branches or tags, ensuring production
builds only come from release branches.

### `sigstore_self_signed(image, repo)`

Verifies an image has a valid Sigstore keyless signature.

**Parameters:**

- `image` (object) - The image input object (`input.image`)
- `repo` (string) - Expected repository name

**Returns:** Boolean - `true` if the image has a valid Sigstore signature

**Example:**

```rego
# Require Sigstore signatures
allow if {
    input.image.repo == "crazymax/diun"
    sigstore_self_signed(input.image, "crazymax/diun")
}
```

Sigstore provides keyless signing using OIDC identity. This function validates
the signature without requiring you to manage public keys.

You can inspect additional signature metadata in the policy:

```rego
allow if {
    input.image.repo == "myorg/myapp"
    sigstore_self_signed(input.image, "myorg/myapp")
    # Verify the build workflow
    input.image.signature.buildSignerURI == "https://github.com/myorg/ci/.github/workflows/build.yml@refs/heads/main"
    # Verify it ran on GitHub's infrastructure
    input.image.signature.runnerEnvironment == "github_hosted"
}
```

### `docker_hardened_image(image, name)`

Verifies an image is from Docker's hardened image program.

**Parameters:**

- `image` (object) - The image input object (`input.image`)
- `name` (string) - Expected hardened image name

**Returns:** Boolean - `true` if the image is a valid Docker hardened image

**Example:**

```rego
# Require hardened Go images
allow if {
    input.image.repo == "docker/dhi-golang"
    docker_hardened_image(input.image, "docker/dhi-golang")
}
```

Docker hardened images receive additional security testing, vulnerability
scanning, and timely security updates. Use this function to enforce hardened
base images for production builds.

## Git verification functions

### `git_signed(git, pubkey)`

Verifies that a Git commit or tag is signed with a specific PGP public key.

**Parameters:**

- `git` (object) - The Git input object (`input.git`)
- `pubkey` (string) - Path to PGP public key file (relative to policy
  directory)

**Returns:** Boolean - `true` if the Git ref is signed with the specified key

**Example:**

```rego
# Require signed tags from maintainer key
allow if {
    input.git.host == "github.com"
    input.git.remote == "myorg/myapp"
    input.git.tagName != ""
    git_signed(input.git, "keys/maintainer.asc")
}
```

Place the public key file in your policy directory:

```text
policy/
├── git.rego
└── keys/
    └── maintainer.asc
```

### `git_signed_github_user(git, user)`

Verifies that a Git commit or tag is signed by a specific GitHub user's
verified signing key.

**Parameters:**

- `git` (object) - The Git input object (`input.git`)
- `user` (string) - GitHub username

**Returns:** Boolean - `true` if the Git ref is signed by the user's GitHub key

**Example:**

```rego
# Require signatures from trusted maintainers
maintainers := ["tonistiigi", "crazy-max", "jsternberg"]

allow if {
    input.git.host == "github.com"
    input.git.remote == "moby/buildkit"
    some i
    git_signed_github_user(input.git, maintainers[i])
}
```

This function queries GitHub's API to verify the signature matches a key
associated with the user's account. It works for both PGP and SSH signatures.

Use this to validate releases are signed by known maintainers without manually
managing public keys.

## HTTP verification functions

### `gpg_signed(http, pubkey)`

Verifies that an HTTP download has a detached GPG signature file.

**Parameters:**

- `http` (object) - The HTTP input object (`input.http`)
- `pubkey` (string) - Path to PGP public key file (relative to policy
  directory)

**Returns:** Boolean - `true` if a valid `.sig` or `.asc` signature file exists

**Example:**

```rego
# Require GPG signatures for downloads
allow if {
    input.http.url
    startswith(input.http.url, "https://nginx.org/download/")
    gpg_signed(input.http, "keys/nginx-signing-key.asc")
}
```

This function looks for a signature file with the same name as the download but
with a `.sig` or `.asc` extension. For example, if downloading
`nginx-1.24.0.tar.gz`, it looks for `nginx-1.24.0.tar.gz.sig` and verifies the
signature against the provided public key.

Place public keys in your policy directory:

```text
policy/
├── downloads.rego
└── keys/
    ├── nginx-signing-key.asc
    └── vendor.asc
```

### `github_attested(repo)`

Verifies that a GitHub release artifact has GitHub's official build
attestations.

**Parameters:**

- `repo` (string) - GitHub repository in `owner/repo` format

**Returns:** Boolean - `true` if the artifact has valid GitHub attestations

**Example:**

```rego
# Require attestations for GitHub releases
allow if {
    input.http.url
    startswith(input.http.url, "https://github.com/cli/cli/releases/download/")
    github_attested("cli/cli")
}
```

GitHub generates attestations for release artifacts uploaded through GitHub
Actions. This function verifies the artifact was built and published by the
repository's official workflows, not manually uploaded by someone with write
access.

Use this to ensure you're downloading official releases rather than potentially
compromised files.

## Version comparison functions

### `version_gt(str)`

Compares version strings and returns true if the input version is strictly
greater than the specified minimum.

**Parameters:**

- `str` (string) - Minimum version to compare against

**Returns:** Boolean - `true` if input version is greater than the minimum

**Example:**

```rego
# Require golangci-lint version greater than v2.0.4
allow if {
    input.image.repo == "golangci/golangci-lint"
    version_gt("v2.0.4")
}
```

This function uses semantic version comparison, understanding version formats
like `v1.2.3`, `1.2.3`, `v1.2.3-beta`, etc.

### `version_gte(str)`

Compares version strings and returns true if the input version is greater than
or equal to the specified minimum.

**Parameters:**

- `str` (string) - Minimum version to compare against

**Returns:** Boolean - `true` if input version is >= the minimum

**Example:**

```rego
# Require at least Python 3.11
allow if {
    input.image.repo == "library/python"
    startswith(input.image.tag, "3.")
    version_gte("3.11")
}
```

Use version comparison functions to enforce minimum versions for tools, base
images, or dependencies with known security vulnerabilities in older releases.

## Using built-in functions

Built-in functions integrate with standard Rego syntax. Combine them with other
conditions to build comprehensive policies:

```rego
package docker

default allow := false

# Production images need multiple checks
allow if {
    input.image.repo == "myorg/backend"
    input.image.isCanonical
    trusted_github_builder_ref("myorg/backend", "refs/heads/main")
    input.image.hasProvenance
    input.image.hasSBOM
}

# Dependencies need signatures
allow if {
    input.http.url
    startswith(input.http.url, "https://releases.myvendor.com/")
    gpg_signed(input.http, "keys/vendor.asc")
    input.http.checksum != ""
}

decision := {"allow": allow}
```

All built-in functions return boolean values, so they work naturally in `allow`
rules alongside other conditions.

## Next steps

- Browse [Example policies](./examples.md) to see built-in functions in context
- Read the [Input reference](./inputs.md) for available input fields
- Learn about [Rego](https://www.openpolicyagent.org/docs/latest/policy-language/)
  for advanced policy patterns
