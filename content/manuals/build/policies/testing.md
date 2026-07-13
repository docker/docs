---
title: Test build policies
linkTitle: Testing
description: Write and run unit tests for build policies, similar to the opa test command
keywords: build policies, opa, rego, testing, unit tests, policy validation
weight: 60
---

The [`docker buildx policy test`](/reference/cli/docker/buildx/policy/test/)
command runs unit tests for build policies using OPA's [standard test
framework](https://www.openpolicyagent.org/docs/policy-testing).

```console
$ docker buildx policy test <path>
```

This validates policy logic with mocked inputs.

For testing against real sources (actual image metadata, Git repositories), use
[`docker buildx policy eval`](/reference/cli/docker/buildx/policy/eval/)
instead. You can use the `eval --print` option to resolve input for a specific
source for writing a test case.

## Basic example

Start with a simple policy that only allows `alpine` images:

```rego {title="Dockerfile.rego"}
package docker

default allow = false

allow if {
    input.image.repo == "alpine"
}

decision := {"allow": allow}
```

Create a test file with the `*_test.rego` suffix. Test functions must start
with `test_`:

```rego {title="Dockerfile_test.rego"}
package docker

test_alpine_allowed if {
    decision.allow with input as {"image": {"repo": "alpine"}}
}

test_ubuntu_denied if {
    not decision.allow with input as {"image": {"repo": "ubuntu"}}
}
```

Run the tests:

```console
$ docker buildx policy test .
test_alpine_allowed: PASS (allow=true)
test_ubuntu_denied: PASS (allow=false)
```

`PASS` indicates that the tests defined in `Dockerfile_test.rego` executed
successfully and all assertions were satisfied.

## Command options

Filter tests by name with `--run`:

```console
$ docker buildx policy test --run alpine .
test_alpine_allowed: PASS (allow=true)
```

Test policies with non-default filenames using `--filename`:

```console
$ docker buildx policy test --filename app.Dockerfile .
```

This loads `app.Dockerfile.rego` and runs `*_test.rego` files against it.

## Test output

Passed tests show the allow status and any deny messages:

```console
test_alpine_allowed: PASS (allow=true)
test_ubuntu_denied: PASS (allow=false, deny_msg=only alpine images are allowed)
```

Failed tests show input, decision output, and missing fields:

```console
test_invalid: FAIL (allow=false)
input:
  {
    "image": {}
  }
decision:
  {
    "allow": false,
    "deny_msg": [
      "only alpine images are allowed"
    ]
  }
missing_input: input.image.repo
```

## Test deny messages

To test custom error messages, capture the full decision result and assert on
the `deny_msg` field.

For a policy with deny messages:

```rego {title="Dockerfile.rego"}
package docker

default allow = false

allow if {
    input.image.repo == "alpine"
}

deny_msg contains msg if {
    not allow
    msg := "only alpine images are allowed"
}

decision := {"allow": allow, "deny_msg": deny_msg}
```

Test the deny message:

```rego {title="Dockerfile_test.rego"}
test_deny_message if {
    result := decision with input as {"image": {"repo": "ubuntu"}}
    not result.allow
    "only alpine images are allowed" in result.deny_msg
}
```

## Test patterns

**Test environment-specific rules:**

```rego
test_production_requires_digest if {
    decision.allow with input as {
        "env": {"target": "production"},
        "image": {"isCanonical": true}
    }
}

test_development_allows_tags if {
    decision.allow with input as {
        "env": {"target": "development"},
        "image": {"isCanonical": false}
    }
}
```

**Test multiple registries:**

```rego
test_dockerhub_allowed if {
    decision.allow with input as {
        "image": {
            "ref": "docker.io/library/alpine",
            "host": "docker.io",
            "repo": "alpine"
        }
    }
}

test_ghcr_allowed if {
    decision.allow with input as {
        "image": {
            "ref": "ghcr.io/myorg/myapp",
            "host": "ghcr.io",
            "repo": "myorg/myapp"
        }
    }
}
```

For available input fields, see the [Input reference](./inputs.md).

## Organize test files

The test runner discovers all `*_test.rego` files recursively:

```plaintext
build-policies/
├── Dockerfile.rego
├── Dockerfile_test.rego
└── tests/
    ├── registries_test.rego
    ├── signatures_test.rego
    └── environments_test.rego
```

Run all tests:

```console
$ docker buildx policy test .
```

Or test specific files:

```console
$ docker buildx policy test tests/registries_test.rego
```
