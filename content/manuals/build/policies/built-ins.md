---
title: Built-in functions
linkTitle: Built-in functions
description: Buildx includes built-in helper functions to make writing policies easier
keywords: build policies, built-in functions, rego functions, signature verification, policy helpers
weight: 90
---

Buildx provides built-in functions, in addition to the [Rego
built-ins](#rego-built-in-functions), to extend Rego policies with
Docker-specific operations like loading local files, verifying Git signatures,
and pinning image digests.

## Rego built-in functions

The functions [documented on this page](#buildx-built-in-functions) are
Buildx-specific functions, distinct from [Rego's standard built-in
functions](https://www.openpolicyagent.org/docs/policy-language#built-in-functions)

Buildx also supports standard Rego built-in functions, but only a subset. To
see the exact list of supported functions, refer to the Buildx [source
code](https://github.com/docker/buildx/blob/master/policy/builtins.go).

## Buildx built-in functions

Buildx provides the following custom built-in functions for policy development:

- [`print`](#print)
- [`load_json`](#load_json)
- [`verify_git_signature`](#verify_git_signature)
- [`pin_image`](#pin_image)

### `print`

Outputs debug information during policy evaluation.

Parameters:

- Any number of values to print

Returns: The values (pass-through)

Example:

```rego
allow if {
    input.image.repo == "alpine"
    print("Allowing alpine image:", input.image.tag)
}
```

Debug output appears when building with `--progress=plain`.

### `load_json`

Loads and parses JSON data from local files in the build context.

Parameters:

- `filename` (string) - Path to JSON file relative to policy directory

Returns: Parsed JSON data as Rego value

Example:

```rego
# Load approved versions from external file
approved_versions = load_json("versions.json")

allow if {
    input.image.repo == "alpine"
    some version in approved_versions.alpine
    input.image.tag == version
}
```

File structure:

```text
project/
├── Dockerfile
├── Dockerfile.rego
└── versions.json
```

versions.json:

```json
{
  "alpine": ["3.19", "3.20"],
  "golang": ["1.21", "1.22"]
}
```

The JSON file must be in the same directory as the policy or in a
subdirectory accessible from the policy location.

### `verify_git_signature`

Verifies PGP signatures on Git commits or tags.

Parameters:

- `git_object` (object) - Either `input.git.commit` or `input.git.tag`
- `keyfile` (string) - Path to PGP public key file (relative to policy
  directory)

Returns: Boolean - `true` if signature is valid, `false` otherwise

Example:

```rego
# Require signed Git tags
allow if {
    input.git.tagName != ""
    verify_git_signature(input.git.tag, "maintainer.asc")
}

# Require signed commits
allow if {
    input.git.commit
    verify_git_signature(input.git.commit, "keys/team.asc")
}
```

Directory structure:

```text
project/
├── Dockerfile.rego
└── maintainer.asc          # PGP public key
```

Or with subdirectory:

```text
project/
├── Dockerfile.rego
└── keys/
    ├── maintainer.asc
    └── team.asc
```

Obtaining public keys:

```console
$ gpg --export --armor user@example.com > maintainer.asc
```

### `pin_image`

Pins an image to a specific digest, overriding the tag-based reference. Use
this to force builds to use specific image versions.

Parameters:

- `image_object` (object) - Must be `input.image` (the current image being
  evaluated)
- `digest` (string) - Target digest in format `sha256:...`

Returns: Boolean - `true` if pinning succeeds

Example:

```rego
# Pin alpine 3.19 to specific digest
alpine_3_19_digest = "sha256:4b7ce07002c69e8f3d704a9c5d6fd3053be500b7f1c69fc0d80990c2ad8dd412"

allow if {
    input.image.repo == "alpine"
    input.image.tag == "3.19"
    pin_image(input.image, alpine_3_19_digest)
}
```

Automatic digest replacement:

```rego
# Replace old digests with patched versions
replace_map = {
  "3.22.0": "3.22.2",
  "3.22.1": "3.22.2",
}

alpine_digests = {
  "3.22.0": "sha256:8a1f59ffb675680d47db6337b49d22281a139e9d709335b492be023728e11715",
  "3.22.2": "sha256:4b7ce07002c69e8f3d704a9c5d6fd3053be500b7f1c69fc0d80990c2ad8dd412",
}

allow if {
    input.image.repo == "alpine"
    some old_version, new_version in replace_map
    input.image.checksum == alpine_digests[old_version]
    print("Replacing", old_version, "with", new_version)
    pin_image(input.image, alpine_digests[new_version])
}
```

This pattern automatically upgrades old image versions to patched releases.

## Next steps

- Browse complete examples: [Example policies](./examples.md)
- Learn policy development workflow: [Using build policies](./usage.md)
- Reference input fields: [Input reference](./inputs.md)
