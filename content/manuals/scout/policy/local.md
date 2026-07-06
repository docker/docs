---
title: Evaluate policies
description: Evaluate Docker Scout policies using the CLI, with built-in and custom Rego policies
keywords: scout, policy, rego, opa, cli, custom policies, policy bundle
---

{{< summary-bar feature_name="Evaluate policies" >}}

`docker scout policy` lets you evaluate images against a configurable policy
set using the CLI. You can use the built-in defaults, adjust thresholds to
match your requirements, or write custom policies in
[Rego](https://www.openpolicyagent.org/docs/latest/policy-language/).

## How it works

When you run `docker scout policy`, the CLI indexes the image into an SBOM,
enriches it with CVE and VEX data, then evaluates each configured policy
in-process. No data is sent to the Scout service, and an
organization is not required for most use cases.

Policies come from three sources, which can be combined:

- Built-in defaults: a curated set embedded in the CLI, used when no other
  source is given.
- OCI policy bundles: Rego packaged as an OCI artifact and pulled from a
  registry with `--policy-bundle`.
- Local `.rego` files: for authoring and iterating on custom policies with
  `--policy-file` or `--policy-dir`.

## Migrate from Policy Evaluation in the Dashboard

If you used the Policies page in the Docker Scout Dashboard, `docker scout
policy` provides the same capability from the CLI. The built-in policies are
the same set. To evaluate an image:

```console
$ docker scout policy <image>
```

If you had customized policies in the dashboard, such as adjusted severity
thresholds or disabled policies, you can replicate those settings with a
`--policy-config` file. See [Configure built-in policies](#configure-built-in-policies).

### Use in CI

Use the [Docker Scout GitHub Action](https://github.com/marketplace/actions/docker-scout)
to evaluate policies as part of your workflow:

```yaml
- name: Evaluate policies
  uses: docker/scout-action@v1
  with:
    command: policy
    image: ${{ env.IMAGE_NAME }}
    organization: <ORG>
```

For other CI platforms, install the
[Docker Scout CLI plugin](/manuals/scout/install.md) on your runner and run
`docker scout policy <image> --exit-code`.

### Migrate the GitHub Action from dashboard-based policy evaluation

The Docker Scout GitHub Action now supports the same local policy configuration
flags as `docker scout policy`. If you used `compare --exit-on policy` with
dashboard-managed policy settings, replicate those settings locally with
`--policy-config`:

```yaml
- uses: docker/scout-action@v1.23.0
  with:
    command: compare
    image: ${{ env.IMAGE_NAME }}
    to-env: production
    exit-on: policy
    policy-config: policies.json
    organization: <ORG>
```

See [Configure built-in policies](#configure-built-in-policies) for the
`policy-config` file format.

## Examples

Evaluate an image against the built-in policy set:

```console
$ docker scout policy myorg/app:latest
```

Use `--exit-code` to fail the pipeline when any policy is not met:

```console
$ docker scout policy myorg/app:latest --exit-code
```

To customize which policies run and their thresholds, pass a policy-config
file:

```console
$ docker scout policy myorg/app:latest --policy-config policies.json
```

Policy results are also surfaced by `docker scout quickview` and
`docker scout compare`, which accept the same `--policy-file`,
`--policy-dir`, `--policy-bundle`, and `--policy-config` flags.

Other useful flags:

```console
# Evaluate a specific platform for a multi-platform image
$ docker scout policy myorg/app:latest --platform linux/arm64

# Show results for a specific policy only
$ docker scout policy myorg/app:latest --only-policy "No copyleft licenses"

# Write the report to a file
$ docker scout policy myorg/app:latest --output report.txt
```

## Built-in policies

The following policies are available by default:

| Policy | What it checks |
| --- | --- |
| No fixable critical or high vulnerabilities | Critical/high CVEs that have a fix available |
| No high-profile vulnerabilities | Curated list of well-known CVEs (Log4Shell, XZ backdoor, and others) |
| No copyleft licenses | Packages under AGPL, GPL, LGPL, MPL, and similar licenses |
| No outdated base images | Base image is behind the latest digest of its tag |
| Supply chain attestations | Provenance and SBOM attestations are attached |
| Default non-root user | Image is configured to run as a non-root user |
| No unapproved base images | Base image matches a configurable allowlist |

## Configure built-in policies

A JSON policy-config file controls which policies run and their thresholds.
Pass it with `--policy-config`.

```json
{
  "policies": [
    {
      "name": "fixable-vulnerabilities",
      "config": {
        "severities": ["CRITICAL"],
        "grace_period_days": 14
      }
    },
    {
      "name": "no-stale-base-images",
      "enabled": false
    }
  ]
}
```

- `policies[].name`: the policy's stable ID (see the following table).
- `policies[].enabled`: set to `false` to skip the policy. Policies not listed are enabled by default.
- `policies[].config`: an object passed to the policy as `data.config`.

### Configuration reference

The following table lists the configurable keys for each built-in policy.

| Policy (stable ID) | config key | Default | Description |
| --- | --- | --- | --- |
| `fixable-vulnerabilities` | `severities` | `["CRITICAL","HIGH"]` | Severity levels that count as a violation |
| `fixable-vulnerabilities` | `fixable_only` | `true` | When `true`, only vulnerabilities with a known fix count |
| `fixable-vulnerabilities` | `package_types` | `[]` | Allowlist of PURL package types to consider; empty means all |
| `fixable-vulnerabilities` | `grace_period_days` | `0` | Days a newly disclosed CVE is exempt |
| `high-profile-vulnerabilities` | `cves` | curated list | CVE IDs considered high-profile |
| `high-profile-vulnerabilities` | `ignored_cves` | `[]` | CVE IDs excluded from causing a failure |
| `high-profile-vulnerabilities` | `include_cisa_kev` | `true` | Also flag vulnerabilities in the CISA KEV catalog |
| `copyleft-license` | `licenses` | AGPL/GPL/LGPL/MPL/… | SPDX license IDs treated as copyleft |
| `copyleft-license` | `ignored_packages` | `[]` | Package URLs exempted from the check |
| `approved-base-images` | `allowed_base_images` | `["*"]` | Glob patterns of allowed base image references |
| `approved-base-images` | `allowed_distros_only` | `true` | When enabled, base image must use an allowed OS distribution |
| `approved-base-images` | `allowed_distros` | curated list | OS distributions considered allowed |
| `supply-chain-attestations` | `required_attestations` | provenance and SBOM predicate types | Attestation predicate types that must be present |

## Write custom policies

Policies are Rego modules in the `docker.scout` package. A policy declares a
boolean `pass` rule and a `violation` set. Use `--policy-file` for a single
file or `--policy-dir` to load a directory recursively.

```rego
# METADATA
# title: No packages from internal registry
# description: Flags packages sourced from registry.internal.example.com.
# custom:
#   name: no-internal-registry
#   result_type: generic
#   weight: 5
#   not_compliant_title: Packages from internal registry found
#   details_order:
#   - purl
#   - reason
package docker.scout

import rego.v1

default pass := false

pass if {
    count(violation) == 0
}

violation contains v if {
    att := oci.referrer("https://scout.docker.com/sbom/v0.1")
    some pkg in att.statement.predicate.artifacts
    contains(pkg.purl, "registry.internal.example.com")
    v := {
        "message": sprintf("Package %s sourced from internal registry", [pkg.purl]),
        "detail": {
            "purl": pkg.purl,
            "reason": "matches registry.internal.example.com",
        },
    }
}
```

```console
# Single file
$ docker scout policy myorg/app:latest --policy-file ./no-internal-registry.rego

# Directory of policies
$ docker scout policy myorg/app:latest --policy-dir ./rego

# Custom policies combined with a config file
$ docker scout policy myorg/app:latest \
  --policy-dir ./rego \
  --policy-config ./policies.json
```

Both `--policy-file` and `--policy-dir` are repeatable. When either is
provided, the built-in defaults are not loaded automatically. To run both
built-in and custom policies together, publish the built-in set to a registry
with `docker scout policy publish`, then pass both bundles using
`--policy-bundle`:

```console
$ docker scout policy myorg/app:latest \
  --policy-bundle registry.example.com/default-policies:latest \
  --policy-bundle registry.example.com/dhi-policies:latest \
  --policy-file ./custom.rego
```

`--policy-bundle` is repeatable, so you can combine as many bundles as needed
alongside local policy files.

### Metadata annotations

The CLI reads [OPA metadata annotations](https://www.openpolicyagent.org/docs/latest/annotations/)
to render results. Place them in a `# METADATA` block immediately above the
`package` declaration.

| Annotation | Purpose |
| --- | --- |
| `title` | Human-readable policy name shown in the report |
| `description` | Longer explanation |
| `custom.name` | Stable ID used to match `--policy-config` entries. Defaults to the package path if omitted |
| `custom.result_type` | How violations are rendered: `vulnerability`, `license`, `boolean`, or `generic` (default) |
| `custom.weight` | Higher weights sort first in the report |
| `custom.not_compliant_title` | Status label shown when the policy fails |
| `custom.details_order` | Ordered list of `detail` keys to display as columns |

### Output contract

- `pass`: `true` when the policy is met. The standard form is `pass if { count(violation) == 0 }`.
- `violation`: a set of objects. Each object should have a `message` string and
  a `detail` object whose keys match `custom.details_order`. An optional
  `remediation` string is shown as remediation guidance.

### Input and built-in functions

The evaluation input is the enriched SBOM. Key entry points:

- `input.source.image`: image metadata, including `input.source.image.config.config.User`,
  `input.source.image.name`, and `input.source.image.digest`.
- `data.config`: the per-policy `config` object from the policy-config file.

The following built-in functions are available in policy Rego:

| Function | Description |
| --- | --- |
| `oci.referrer(predicateType)` | Retrieve an attestation by predicate type |
| `oci.canonical_name(ref)` | Normalize an image reference, for example `"node:25"` to `"docker.io/library/node"` |
| `oci.image_digest(ref)` | Resolve the current digest of a tag from its registry |
| `oci.index(ref, digest)` | Get the OCI image index for an image |
| `oci.referrer_index(ref, digest)` | Get the OCI referrer index for an image |
| `oci.referrer_by_digest(ref, digest)` | Get a referrer for an image by digest |
| `scout.parse_purl(purl)` | Parse a PURL into its components |
| `scout.package_provenance(purl)` | Provenance of a package from the SBOM |
| `scout.vulnerabilities(purls)` | Vulnerabilities for the given PURLs |
| `scout.package_recommendation(purl)` | Recommended (fixed) version for a package |
| `scout.base_image()` | Base image matches recorded in the SBOM |
| `cosign.verify_dsse(envelope, opts)` | Verify an in-toto DSSE envelope's cosign signature |
| `cosign.verify_image(ref, opts)` | Verify an image's cosign signature |
| `gpg.verify_commit(...)` | Verify a detached commit signature |

Common predicate types for `oci.referrer`:

| Predicate type | Content |
| --- | --- |
| `https://scout.docker.com/vulnerabilities/v0.1` | CVEs per package |
| `https://scout.docker.com/sbom/v0.1` | SBOM artifacts, with `purl` and `licenses` |
| `https://scout.docker.com/provenance/v0.1` | Build provenance, including `base_image` |
| `https://openvex.dev/ns/v0.2.0` | VEX statements |

## Debug policies

Add `print()` statements to your Rego and enable debug output by setting
`"debug": true` in the policy-config file:

```rego
violation contains v if {
    att := oci.referrer("https://scout.docker.com/sbom/v0.1")
    some pkg in att.statement.predicate.artifacts
    print("checking", pkg.purl)
    contains(pkg.purl, blocked)
    # ...
}
```

```json
{
  "debug": true,
  "policies": [
    { "name": "no-internal-registry" }
  ]
}
```

Output is prefixed with the policy name and source line:

```text
no-internal-registry#28: checking pkg:deb/debian/curl@7.88.1
```

> [!NOTE]
>
> The `debug` field in the policy-config controls `print()` output from your
> Rego. The global `--debug` flag is separate: it enables CLI-level debug
> logging for bundle loading, registry resolution, and similar internals.

### Inspect raw evaluation results

Use `--result-file` to write the full evaluation result for every policy to a
JSON file. This is useful when iterating on a custom policy to inspect
intermediate values.

```console
$ docker scout policy myorg/app:latest \
  --policy-file ./no-internal-registry.rego \
  --result-file result.json
```

Each entry in the output contains:

- `pass`: the policy's boolean outcome.
- `violations`: the `detail` object of each reported violation.
- `bindings`: the raw `data.docker.scout` document, including the `violation`
  set and any other complete rules, for inspecting intermediate values.
- `metrics`: OPA evaluation metrics (timers and counters).

```json
{
  "no-internal-registry": {
    "pass": false,
    "violations": [
      { "purl": "pkg:deb/debian/curl@7.88.1", "reason": "matches \"registry.internal.example.com\"" }
    ],
    "bindings": {
      "blocked": "registry.internal.example.com",
      "violation": [
        {
          "message": "Package pkg:deb/debian/curl@7.88.1 sourced from internal registry",
          "detail": {
            "purl": "pkg:deb/debian/curl@7.88.1",
            "reason": "matches \"registry.internal.example.com\""
          }
        }
      ]
    },
    "metrics": {
      "timer_rego_query_eval_ns": 1234567
    }
  }
}
```

## Share policies as OCI bundles

Package `.rego` files as an OCI artifact and distribute them through any
registry.

### Publish a bundle

```console
# Publish a directory of policies
$ docker scout policy publish \
  --policy-dir ./rego \
  registry.example.com/my-policies:latest

# Publish specific files
$ docker scout policy publish \
  --policy-file fixable.rego \
  --policy-file licenses.rego \
  registry.example.com/my-policies:latest

# Publish the built-in default set
$ docker scout policy publish registry.example.com/my-policies:latest
```

Each module's metadata is validated before publishing. The command prints the
resulting digest and the list of bundled policies.

### Use a bundle

```console
$ docker scout policy myorg/app:latest \
  --policy-bundle registry.example.com/my-policies:latest

# Combine a bundle with local files and a config
$ docker scout policy myorg/app:latest \
  --policy-bundle registry.example.com/my-policies:latest \
  --policy-file ./extra.rego \
  --policy-config ./policies.json
```

`--policy-bundle` is repeatable. Authentication uses your existing Docker
registry credentials. Bundles are cached by digest, so re-running against the
same bundle does not re-download it. A new digest (for example, after
re-publishing `:latest`) is fetched automatically.
