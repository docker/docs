---
title: Apply Docker Hardened Image policies to your images
linktitle: Apply image policies
description: Learn how to hold your own images to Docker Hardened Image security and compliance standards using the Docker Scout CLI.
weight: 50
keywords: docker scout policies, image security policy, container compliance, dhi policies, vulnerability policy check
---

Docker publishes the set of security and compliance policies that Docker
Hardened Images (DHIs) are built to meet, so you can hold your own images to the
same standards. You evaluate images against these policies with the
[`docker scout policy`](../../scout/policy/local.md) command.

These policies encode requirements such as running as a non-root user, being
free of fixable critical and high vulnerabilities, containing no embedded
malware or secrets, and shipping signed supply chain attestations. They don't
verify whether an image is a DHI or built on a DHI base image; they check
whether an image meets the same bar that DHIs are held to.

Unlike the built-in Docker Scout policies, the DHI policies aren't embedded in
the CLI. They're maintained as Rego source in the
[`docker-hardened-images/policies`](https://github.com/docker-hardened-images/policies)
repository and published as an OCI policy bundle at
[`dhi/policies`](https://hub.docker.com/repository/docker/dhi/policies/general).
You pull the bundle at evaluation time with the `--policy-bundle` flag, so you
can apply DHI standards locally, in CI, or both, without sending any data to the
Docker Scout service.

## Policies in the DHI bundle

The `dhi/policies` bundle includes the following policies:

| Policy | What it checks |
| --- | --- |
| No default root user for non-dev images | The image is configured to run as a non-root user. |
| No fixable vulnerabilities past their remediation SLA | No fixable CVEs remain unaddressed past their remediation SLA (7 days for critical and high, 30 days for others). |
| No high-profile vulnerabilities | The image is free of a curated list of well-known CVEs, optionally including the CISA KEV catalog. |
| No embedded malware | A malware scan attestation is present and passing. |
| No embedded secrets | A secret scan attestation is present and passing. |
| No failing tests | A test attestation is present and passing. |
| Signed supply chain attestations | SBOM and provenance attestations are attached and signed. |
| Unintentional shell or package manager | No undeclared shell or package manager is present in the image. |
| STIG scan | For FIPS-compliant images, the STIG scan meets the required score. |

For the authoritative list and the Rego source for each policy, see the
[`docker-hardened-images/policies`](https://github.com/docker-hardened-images/policies)
repository.

## Prerequisites

- The [Docker Scout CLI plugin](/manuals/scout/install.md). It's included with
  Docker Desktop.
- Access to pull the `dhi/policies` bundle from Docker Hub. Authentication uses
  your existing Docker registry credentials, so sign in first:

  ```console
  $ docker login
  ```

## Evaluate an image against the DHI policies

To evaluate an image against the DHI policy bundle, pass the bundle reference to
`docker scout policy` with the `--policy-bundle` flag:

```console
$ docker scout policy <image> --policy-bundle dhi/policies:latest
```

The CLI pulls the bundle, indexes the image into an SBOM, enriches it with CVE
and VEX data, and evaluates each policy in the bundle against that data. Bundles
are cached by digest, so re-running against the same bundle doesn't re-download
it.

### Example: Build and evaluate a DHI-based image

The following example builds an image from a DHI base image and evaluates it
against the DHI policy bundle.

#### Step 1: Use a DHI base image in your Dockerfile

Create a Dockerfile that uses a Docker Hardened Image from the DHI catalog as
the base. For example:

```dockerfile
# Dockerfile
FROM dhi.io/python:3.13

ENTRYPOINT ["python", "-c", "print('Hello from a DHI-based image')"]
```

#### Step 2: Build the image

Open a terminal and navigate to the directory containing your Dockerfile. Then,
build the image and load it into your local image store:

```console
$ docker build --load -t my-dhi-app:v1 .
```

#### Step 3: Evaluate the image against the DHI policies

Sign in and evaluate the local image against the DHI policy bundle:

```console
$ docker login
$ docker scout policy my-dhi-app:v1 \
  --policy-bundle dhi/policies:latest
```

The command prints a compliance result for each policy in the bundle, along with
the details of any violations.

## Customize the DHI policies

You can tune which policies run and their thresholds with a `--policy-config`
file. The
[`docker-hardened-images/policies`](https://github.com/docker-hardened-images/policies)
repository includes an example `config.json` you can start from.

```console
$ docker scout policy my-dhi-app:v1 \
  --policy-bundle dhi/policies:latest \
  --policy-config ./config.json
```

The config file matches policies by their stable name and lets you disable
individual policies or adjust their settings, such as severity levels and grace
periods. For the config file format, see
[Configure built-in policies](../../scout/policy/local.md#configure-built-in-policies).

You can also combine the DHI bundle with the built-in Docker Scout policies,
additional bundles, or your own custom Rego files. `--policy-bundle`,
`--policy-file`, and `--policy-dir` are all repeatable:

```console
$ docker scout policy my-dhi-app:v1 \
  --policy-bundle dhi/policies:latest \
  --policy-file ./custom.rego
```

For more on authoring custom policies and combining policy sources, see
[Evaluate policies](../../scout/policy/local.md).

## Enforce policy compliance in CI

Use the [Docker Scout GitHub Action](https://github.com/docker/scout-action) to
evaluate the DHI policies on every push and fail the workflow when an image
doesn't meet them. The following workflow builds the image, then evaluates it
against the DHI policy bundle:

```yaml
name: DHI policy check

on:
  push:

env:
  IMAGE_NAME: my-dhi-app:${{ github.sha }}

jobs:
  policy:
    runs-on: ubuntu-latest
    steps:
      - name: Check out the repository
        uses: actions/checkout@v4

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Log in to Docker Hub
        uses: docker/login-action@v3
        with:
          username: ${{ secrets.DOCKER_USER }}
          password: ${{ secrets.DOCKER_PAT }}

      - name: Build the image
        uses: docker/build-push-action@v6
        with:
          context: .
          load: true
          tags: ${{ env.IMAGE_NAME }}

      - name: Evaluate DHI policies
        uses: docker/scout-action@v1.23.1
        with:
          command: policy
          image: ${{ env.IMAGE_NAME }}
          policy-bundle: dhi/policies:latest
          exit-code: true
```

The `docker/login-action` step authenticates with Docker Hub so the runner can
pull the DHI base image and the `dhi/policies` bundle. Store your Docker Hub
username and a [personal access token](/manuals/security/access-tokens.md) as the
`DOCKER_USER` and `DOCKER_PAT` repository secrets.

Set `exit-code: true` to fail the step when any policy isn't met. The
`policy-bundle` input accepts a comma-separated list of bundles, and you can
combine it with the `policy-file`, `policy-dir`, and `policy-config` inputs, the
same as the CLI flags.

For more on running policy evaluation in CI, see
[Evaluate policies](../../scout/policy/local.md#use-in-ci).
