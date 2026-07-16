---
title: Validating Git repositories
linkTitle: Git validation
description: Write policies to validate Git repositories used in your builds
keywords: build policies, git validation, git signatures, gpg, signed commits, signed tags
weight: 40
---

Git repositories often appear in Docker builds as source code inputs. The `ADD`
instruction can clone repositories, and build contexts can reference Git URLs.
Validating these inputs ensures you're building from trusted sources with
verified versions.

This guide teaches you to write policies that validate Git inputs, from basic
version pinning to verifying signed commits and tags.

## Prerequisites

You should understand the policy basics from the [Introduction](./intro.md):
creating policy files, basic Rego syntax, and how policies evaluate during
builds.

## What are Git inputs?

Git inputs come from `ADD` instructions that reference Git repositories:

```dockerfile
# Clone a specific tag
ADD https://github.com/moby/buildkit.git#v0.26.1 /buildkit

# Clone a branch
ADD https://github.com/user/repo.git#main /src

# Clone a commit
ADD https://github.com/user/repo.git#abcde123 /src
```

The build context can also be a Git repository when you build with:

```console
$ docker build https://github.com/user/repo.git#main
```

Each Git reference triggers a policy evaluation. Your policy can inspect
repository URLs, validate versions, check commit metadata, and verify
signatures.

## Match specific repositories

The simplest Git policy restricts which repositories can be used:

```rego {title="Dockerfile.rego"}
package docker

default allow := false

allow if input.local

allow if {
  input.git.host == "github.com"
  input.git.remote == "https://github.com/moby/buildkit.git"
}

decision := {"allow": allow}
```

This policy:

- Denies all inputs by default
- Allows local build context
- Allows only the BuildKit repository from GitHub

The `host` field contains the Git server hostname, and `remote` contains the
full repository URL. Test it:

```dockerfile {title="Dockerfile"}
FROM scratch
ADD https://github.com/moby/buildkit.git#v0.26.1 /
```

```console
$ docker build .
```

The build succeeds. Try a different repository and it fails.

You can match multiple repositories with additional rules:

```rego
allow if {
  input.git.host == "github.com"
  input.git.remote == "https://github.com/moby/buildkit.git"
}

allow if {
  input.git.host == "github.com"
  input.git.remote == "https://github.com/docker/cli.git"
}

decision := {"allow": allow}
```

## Pin to specific versions

Tags and branches can change over time. Pin to specific versions to ensure
reproducible builds:

```rego
package docker

default allow := false

allow if input.local

allow if {
  input.git.remote == "https://github.com/moby/buildkit.git"
  input.git.tagName == "v0.26.1"
}

decision := {"allow": allow}
```

The `tagName` field contains the tag name when the Git reference points to a
tag. Use `branch` for branches:

```rego
allow if {
  input.git.remote == "https://github.com/user/repo.git"
  input.git.branch == "main"
}
```

Or use `ref` for any type of reference (branch, tag, or commit SHA):

```rego
allow if {
  input.git.ref == "v0.26.1"
}
```

## Use version allowlists

For repositories you trust but want to control versions, maintain an allowlist:

```rego
package docker

default allow := false

allowed_versions = [
    {"tag": "v0.26.1", "annotated": true, "sha": "abc123"},
]

is_buildkit if {
    input.git.remote == "https://github.com/moby/buildkit.git"
}

allow if {
    not is_buildkit
}

allow if {
    is_buildkit
    some version in allowed_versions
    input.git.tagName == version.tag
    input.git.isAnnotatedTag == version.annotated
    startswith(input.git.commitChecksum, version.sha)
}

decision := {"allow": allow}
```

This policy:

- Defines an allowlist of approved versions with metadata
- Uses a helper rule (`is_buildkit`) for readability
- Allows all non-BuildKit inputs
- For BuildKit, checks the tag name, whether it's an annotated tag, and the commit SHA against the allowlist

The helper rule makes complex policies more maintainable. You can expand the
allowlist as new versions are approved:

```rego
allowed_versions = [
    {"tag": "v0.26.1", "annotated": true, "sha": "abc123"},
    {"tag": "v0.27.0", "annotated": true, "sha": "def456"},
    {"tag": "v0.27.1", "annotated": true, "sha": "789abc"},
]
```

## Validate with regex patterns

Use pattern matching for semantic versioning:

```rego
package docker

default allow := false

allow if input.local

allow if {
  input.git.remote == "https://github.com/moby/buildkit.git"
  regex.match(`^v[0-9]+\.[0-9]+\.[0-9]+$`, input.git.tagName)
}

decision := {"allow": allow}
```

This allows any BuildKit tag matching the pattern `vX.Y.Z` where X, Y, and Z
are numbers. The regex ensures you're using release versions, not pre-release
tags like `v0.26.0-rc1`.

Match major versions:

```rego
# Only allow v0.x releases
allow if {
  input.git.remote == "https://github.com/moby/buildkit.git"
  regex.match(`^v0\.[0-9]+\.[0-9]+$`, input.git.tagName)
}
```

## Inspect commit metadata

The `commit` object provides detailed information about commits:

```rego
package docker

default allow := false

allow if input.local

# Check commit author
allow if {
  input.git.remote == "https://github.com/user/repo.git"
  input.git.commit.author.email == "trusted@example.com"
}

decision := {"allow": allow}
```

The `commit` object includes:

- `author.name`: Author's name
- `author.email`: Author's email
- `author.when`: When the commit was authored
- `committer.name`: Committer's name
- `committer.email`: Committer's email
- `committer.when`: When the commit was committed
- `message`: Commit message

Validate commit messages:

```rego
allow if {
  input.git.commit
  contains(input.git.commit.message, "Signed-off-by:")
}
```

Pin to specific commit SHA:

```rego
allow if {
  input.git.commitChecksum == "abc123def456..."
}
```

## Require signed commits

GPG-signed commits prove authenticity. Check for commit signatures:

```rego
package docker

default allow := false

allow if input.local

allow if {
  input.git.remote == "https://github.com/moby/buildkit.git"
  input.git.commit.pgpSignature != null
}

decision := {"allow": allow}
```

The `pgpSignature` field is `null` for unsigned commits. For signed commits, it
contains signature details.

SSH signatures work similarly:

```rego
allow if {
  input.git.commit.sshSignature != null
}
```

## Require signed tags

Annotated tags can be signed, providing a cryptographic guarantee of the
release:

```rego
package docker

default allow := false

allow if input.local

allow if {
  input.git.remote == "https://github.com/moby/buildkit.git"
  input.git.tag.pgpSignature != null
}

decision := {"allow": allow}
```

The `tag` object is only available for annotated tags. It includes:

- `tagger.name`: Who created the tag
- `tagger.email`: Tagger's email
- `tagger.when`: When the tag was created
- `message`: Tag message
- `pgpSignature`: GPP signature (if signed)
- `sshSignature`: SSH signature (if signed)

Lightweight tags don't have a `tag` object, so this policy effectively requires
annotated, signed tags.

## Verify signatures with public keys

Use the `verify_git_signature()` function to cryptographically verify Git
signatures against trusted public keys:

```rego
package docker

default allow := false

allow if input.local

allow if {
  input.git.remote == "https://github.com/moby/buildkit.git"
  input.git.tagName != ""
  verify_git_signature(input.git.tag, "keys.asc")
}

decision := {"allow": allow}
```

This verifies that Git tags are signed by keys in the `keys.asc` public
key file. To set this up:

1. Export maintainer public keys:
   ```console
   $ curl https://github.com/user.gpg > keys.asc
   ```
2. Place `keys.asc` alongside your policy file

The function verifies PGP signatures on commits or tags. See [Built-in
functions](./built-ins.md) for more details.

## Apply conditional rules

Use different rules for different contexts. Allow unsigned refs during
development but require signing for production:

```rego
package docker

default allow := false

allow if input.local

is_buildkit if {
    input.git.remote == "https://github.com/moby/buildkit.git"
}

is_version_tag if {
    is_buildkit
    regex.match(`^v[0-9]+\.[0-9]+\.[0-9]+$`, input.git.tagName)
}

# Version tags must be signed
allow if {
    is_version_tag
    input.git.tagName != ""
    verify_git_signature(input.git.tag, "keys.asc")
}

# Non-version refs allowed in development
allow if {
    is_buildkit
    not is_version_tag
    input.env.target != "release"
}

decision := {"allow": allow}
```

This policy:

- Defines helper rules for readability
- Requires signed version tags from maintainers
- Allows unsigned refs (branches, commits) unless building the release target
- Uses `input.env.target` to detect the build target

Build a development target without signatures:

```console
$ docker buildx build --target=dev .
```

Build the release target, and signing is enforced:

```console
$ docker buildx build --target=release .
```

## Next steps

You now understand how to validate Git repositories in build policies. To
continue learning:

- Browse [Example policies](./examples.md) for complete policy patterns
- Read [Built-in functions](./built-ins.md) for Git signature verification
  functions
- Check the [Input reference](./inputs.md) for all available Git fields
