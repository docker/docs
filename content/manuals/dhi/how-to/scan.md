---
title: Scan Docker Hardened Images
linktitle: Scan an image
description: Learn how to scan Docker Hardened Images for known vulnerabilities using Docker Scout, Grype, or Trivy.
keywords: scan container image, docker scout cves, grype scanner, trivy container scanner, vex attestation
weight: 46
---

Docker Hardened Images (DHIs) are designed to be secure by default, but like any
container image, it's important to scan them regularly as part of your
vulnerability management process.

You can scan DHIs using the same tools you already use for standard images, such
as Docker Scout, Grype, and Trivy. DHIs follow the same formats and standards
for compatibility across your security tooling. Before you scan an image, the image must
be mirrored into your organization on Docker Hub.

> [!NOTE]
>
> When you have a Docker Hardened Images Enterprise subscription, [Docker
> Scout](/manuals/scout/_index.md) is automatically enabled at no additional
> cost for all mirrored Docker Hardened Image repositories on Docker Hub. You
> can view scan results directly in the Docker Hub UI under your organization's
> repository.

> [!IMPORTANT]
>
> You must authenticate to the Docker Hardened Images registry (`dhi.io`) to
> pull images. Use your Docker ID credentials (the same username and password
> you use for Docker Hub) when signing in. If you don't have a Docker account,
> [create one](../../accounts/create-account.md) for free.
>
> Run `docker login dhi.io` to authenticate.

## Docker Scout

Docker Scout is integrated into Docker Desktop and the Docker CLI. It provides
vulnerability insights, CVE summaries, and direct links to remediation guidance.

### Scan a DHI using Docker Scout

To scan a Docker Hardened Image using Docker Scout, run the following
command:

```console
$ docker scout cves dhi.io/<image>:<tag> --platform <platform>
```

Example output:

```plaintext
    v SBOM obtained from attestation, 101 packages found
    v Provenance obtained from attestation
    v VEX statements obtained from attestation
    v No vulnerable package detected
    ...
```

For more detailed filtering and JSON output, see [Docker Scout CLI reference](../../../reference/cli/docker/scout/_index.md).

### Automate DHI scanning in CI/CD with Docker Scout

Integrating Docker Scout into your CI/CD pipeline enables you to automatically
verify that images built from Docker Hardened Images remain free from known
vulnerabilities during the build process. This proactive approach ensures the
continued security integrity of your images throughout the development
lifecycle.

#### Example GitHub Actions workflow

The following is a sample GitHub Actions workflow that builds an image and scans
it using Docker Scout:

```yaml {collapse="true"}
name: DHI Vulnerability Scan

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ "**" ]

env:
  REGISTRY: docker.io
  IMAGE_NAME: ${{ github.repository }}
  SHA: ${{ github.event.pull_request.head.sha || github.event.after }}

jobs:
  scan:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write
      pull-requests: write

    steps:
      - name: Checkout repository
        uses: actions/checkout@v3

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Log in to Docker Hub
        uses: docker/login-action@v2
        with:
          username: ${{ secrets.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKER_PASSWORD }}

      - name: Build Docker image
        run: |
          docker build -t ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:${{ env.SHA }} .

      - name: Run Docker Scout CVE scan
        uses: docker/scout-action@v1
        with:
          command: cves
          image: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:${{ env.SHA }}
          only-severities: critical,high
          exit-code: true
```

The `exit-code: true` parameter ensures that the workflow fails if any critical or
high-severity vulnerabilities are detected, preventing the deployment of
insecure images.

For more details on using Docker Scout in CI, see [Integrating Docker
Scout with other systems](/manuals/scout/integrations/_index.md).

### Comparing Docker Scout results with other scanners

Some vulnerabilities reported by other scanners may not appear in Docker Scout results. This can happen for several
reasons:

- Hardware-specific vulnerabilities: Certain vulnerabilities may only affect specific hardware architectures (for
  example, Power10 processors) that are not relevant to Docker images, so they are not reported by Docker Scout.
- VEX statement filtering: Docker Scout automatically applies VEX statements to document and suppress vulnerabilities
  that do not apply to the image. If your scanner does not consume VEX statements, you may see more vulnerabilities
  reported than what appears in Docker Scout results.
- Temporary vulnerability identifiers: Temporary vulnerability identifiers (like `TEMP-xxxxxxx` from Debian) are not
  surfaced by Docker Scout, as they are not intended for external reference.

While Docker Scout handles this filtering automatically, you can manually configure similar filtering with other
scanners using [Grype ignore rules](https://github.com/anchore/grype#specifying-matches-to-ignore) in its configuration
file (`~/.grype.yaml`) or [Trivy policy exceptions](https://trivy.dev/v0.19.2/misconfiguration/policy/exceptions/) using
REGO rules to filter out specific vulnerabilities by CVE ID, package name, fix state, or other criteria. You can also
use VEX statements with other scanners as described in [Use VEX to filter known non-exploitable
CVEs](#use-vex-to-filter-known-non-exploitable-cves).

## Grype

[Grype](https://github.com/anchore/grype) is an open-source scanner that checks
container images against vulnerability databases like the NVD and distro
advisories.

### Scan a DHI using Grype

After installing Grype, you can scan a Docker Hardened Image by pulling
the image and running the scan command:

```console
$ docker pull dhi.io/<image>:<tag>
$ grype dhi.io/<image>:<tag>
```

Example output:

```plaintext
NAME               INSTALLED              FIXED-IN     TYPE  VULNERABILITY     SEVERITY    EPSS%  RISK
libperl5.36        5.36.0-7+deb12u2       (won't fix)  deb   CVE-2023-31484    High        79.45    1.1
perl               5.36.0-7+deb12u2       (won't fix)  deb   CVE-2023-31484    High        79.45    1.1
perl-base          5.36.0-7+deb12u2       (won't fix)  deb   CVE-2023-31484    High        79.45    1.1
...
```

You should include the `--vex` flag to apply VEX statements during the scan,
which filter out known non-exploitable CVEs. For more information, see the [VEX
section](#use-vex-to-filter-known-non-exploitable-cves).

## Trivy

[Trivy](https://github.com/aquasecurity/trivy) is an open-source vulnerability
scanner for containers and other artifacts. It detects vulnerabilities in OS
packages and application dependencies.

### Scan a DHI using Trivy

After installing Trivy, you can scan a Docker Hardened Image by pulling
the image and running the scan command:

```console
$ docker pull dhi.io/<image>:<tag>
$ trivy image --scanners vuln dhi.io/<image>:<tag>
```

To filter vulnerabilities using VEX statements, Trivy supports multiple
approaches. Docker recommends using VEX Hub, which provides a seamless workflow
for automatically downloading and applying VEX statements from configured
repositories.

#### Using VEX Hub (recommended)

Configure Trivy to download the Docker Hardened Images advisories repository
from VEX Hub. Run the following commands to set up the VEX repository:

```console
$ trivy vex repo init
$ cat << REPO > ~/.trivy/vex/repository.yaml
repositories:
  - name: default
    url: https://github.com/aquasecurity/vexhub
    enabled: true
    username: ""
    password: ""
    token: ""
  - name: dhi-vex
    url: https://github.com/docker-hardened-images/advisories
    enabled: true
REPO
$ trivy vex repo list
$ trivy vex repo download
```

After setting up VEX Hub, you can scan a Docker Hardened Image with VEX filtering:

```console
$ docker pull dhi.io/<image>:<tag>
$ trivy image --scanners vuln --vex repo dhi.io/<image>:<tag>
```

For example, scanning the `dhi.io/python:3.13` image:

```console
$ trivy image --scanners vuln --vex repo dhi.io/python:3.13
```

Example output:

```plaintext
Report Summary

┌─────────────────────────────────────────────────────────────────────────────┬────────────┬─────────────────┐
│                                   Target                                    │    Type    │ Vulnerabilities │
├─────────────────────────────────────────────────────────────────────────────┼────────────┼─────────────────┤
│ dhi.io/python:3.13 (debian 13.2)                                            │   debian   │        0        │
├─────────────────────────────────────────────────────────────────────────────┼────────────┼─────────────────┤
│ opt/python-3.13.11/lib/python3.13/site-packages/pip-25.3.dist-info/METADATA │ python-pkg │        0        │
└─────────────────────────────────────────────────────────────────────────────┴────────────┴─────────────────┘
Legend:
- '-': Not scanned
- '0': Clean (no security findings detected)
```

The `--vex repo` flag applies VEX statements from the configured repository during the scan,
which filters out known non-exploitable CVEs.

#### Using local VEX files

In addition to VEX Hub, Trivy also supports the use of local VEX files for
vulnerability filtering. You can download the VEX attestation that Docker
Hardened Images provide and use it directly with Trivy.

First, download the VEX attestation for your image:

```console
$ docker scout vex get dhi.io/<image>:<tag> --output vex.json
```

Then scan the image with the local VEX file:

```console
$ trivy image --scanners vuln --vex vex.json dhi.io/<image>:<tag>
```

## Use VEX to filter known non-exploitable CVEs

Docker Hardened Images include signed VEX (Vulnerability Exploitability
eXchange) attestations that identify vulnerabilities not relevant to the image’s
runtime behavior.

When using Docker Scout, these VEX statements are automatically applied and no
manual configuration needed.

> [!NOTE]
>
> By default, VEX attestations are fetched from `registry.scout.docker.com`. Ensure that you can access this registry
> if your network has outbound restrictions. You can also mirror the attestations to an alternate registry. For more
> details, see [Mirror to a third-party registry](mirror.md#mirror-to-a-third-party-registry).

To manually create a JSON file of VEX attestations for tools that support it:

```console
$ docker scout vex get dhi.io/<image>:<tag> --output vex.json
```

> [!NOTE]
>
> The `docker scout vex get` command requires [Docker Scout
> CLI](https://github.com/docker/scout-cli/) version 1.18.3 or later.
>
> If the image exists locally on your device, you must prefix the image name with `registry://`. For example, use
> `registry://docs/dhi-python:3.13` instead of `docs/dhi-python:3.13`.

For example:

```console
$ docker scout vex get dhi.io/python:3.13 --output vex.json
```

This creates a `vex.json` file containing the VEX statements for the specified
image. You can then use this file with tools that support VEX to filter out
known non-exploitable CVEs.

For example, with Grype you can use the `--vex` flag to apply the VEX
statements during the scan:

```console
$ grype dhi.io/python:3.13 --vex vex.json
```
