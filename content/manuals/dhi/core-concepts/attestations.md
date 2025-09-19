---
title: Attestations
description: Review the full set of signed attestations included with each Docker Hardened Image, such as SBOMs, VEX, build provenance, and scan results.
keywords: container image attestations, signed sbom, build provenance, slsa compliance, vex document
---

Docker Hardened Images (DHIs) include comprehensive, signed security
attestations that verify the image's build process, contents, and security
posture. These attestations are a core part of secure software supply chain
practices and help users validate that an image is trustworthy and
policy-compliant.

## What is an attestation?

An attestation is a signed statement that provides verifiable information
about an image, such as how it was built, what's inside it, and what security
checks it has passed. Attestations are typically signed using Sigstore tooling
(such as Cosign), making them tamper-evident and cryptographically verifiable.

Attestations follow standardized formats (like [in-toto](https://in-toto.io/),
[CycloneDX](https://cyclonedx.org/), and [SLSA](https://slsa.dev/)) and are
attached to the image as OCI-compliant metadata. They can be generated
automatically during image builds or added manually to document extra tests,
scan results, or custom provenance.

## Why are attestations important?

Attestations provide critical visibility into the software supply chain by:

- Documenting *what* went into an image (e.g., SBOMs)
- Verifying *how* it was built (e.g., build provenance)
- Capturing *what security scans* it has passed or failed (e.g., CVE reports,
  secrets scans, test results)
- Helping organizations enforce compliance and security policies
- Supporting runtime trust decisions and CI/CD policy gates

They are essential for meeting industry standards such as SLSA,
and help teams reduce the risk of supply chain attacks by making build and
security data transparent and verifiable.

## How Docker Hardened Images use attestations

All DHIs are built using [SLSA Build Level
3](https://slsa.dev/spec/latest/levels) practices, and each image variant is
published with a full set of signed attestations. These attestations allow users
to:

- Verify that the image was built from trusted sources in a secure environment
- View SBOMs in multiple formats to understand component-level details
- Review scan results to check for vulnerabilities or embedded secrets
- Confirm the build and deployment history of each image

Attestations are automatically published and associated with each mirrored DHI
in your Docker Hub organization. They can be inspected using tools like [Docker
Scout](../how-to/verify.md) or
[Cosign](https://docs.sigstore.dev/cosign/overview), and are consumable by CI/CD
tooling or security platforms.

## Available attestations

While every DHI variant includes a set of attestations, the attestations may
vary based on the image variant. For example, some images may include a STIG
scan attestation. The following table is a comprehensive list of all
attestations that may be included with a DHI. To see which attestations are
available for a specific image variant, you can [view the image variant
details](../how-to/explore.md#view-image-variant-details) in Docker Hub.

| Attestation type           | Description                                                                                                                                                                                                                     | Predicate type URI                                 |
|----------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------------------------------------------------|
| CycloneDX SBOM             | A software bill of materials in [CycloneDX](https://cyclonedx.org/) format, listing components, libraries, and versions.                                                                                                      | `https://cyclonedx.org/bom/v1.5`                  |
| STIG scan                  | Results of a STIG scan, with output in HTML and XCCDF formats.                                                                                                                           | `https://docker.com/dhi/stig/v0.1`                |
| CVEs (In-Toto format)      | A list of known vulnerabilities (CVEs) affecting the image's components, based on package and distribution scanning.                                                                           | `https://in-toto.io/attestation/vulns/v0.1`       |
| VEX                        | A [Vulnerability Exploitability eXchange (VEX)](https://openvex.dev/) document that identifies vulnerabilities that do not apply to the image and explains why (e.g., not reachable or not present).                         | `https://openvex.dev/ns/v0.2.0`                   |
| Scout health score         | A signed attestation from Docker Scout that summarizes the overall security and quality posture of the image.                                                                           | `https://scout.docker.com/health/v0.1`            |
| Scout provenance           | Provenance metadata generated by Docker Scout, including the source Git commit, build parameters, and environment details.                                                               | `https://scout.docker.com/provenance/v0.1`        |
| Scout SBOM                 | An SBOM generated and signed by Docker Scout, including additional Docker-specific metadata.                                                                                             | `https://scout.docker.com/sbom/v0.1`              |
| Secrets scan               | Results of a scan for accidentally included secrets, such as credentials, tokens, or private keys.                                                                                       | `https://scout.docker.com/secrets/v0.1`           |
| Tests                      | A record of automated tests run against the image, such as functional checks or validation scripts.                                                                                      | `https://scout.docker.com/tests/v0.1`             |
| Virus scan                 | Results of antivirus scans performed on the image layers.                                                                                                                                 | `https://scout.docker.com/virus/v0.1`             |
| CVEs (Scout format)                      | A vulnerability report generated by Docker Scout, listing known CVEs and severity data.                                                                                                  | `https://scout.docker.com/vulnerabilities/v0.1`   |
| SLSA provenance            | A standard [SLSA](https://slsa.dev/) provenance statement describing how the image was built, including build tool, parameters, and source.                                               | `https://slsa.dev/provenance/v0.2`                |
| SLSA verification summary  | A summary attestation indicating the image's compliance with SLSA requirements.                                                                                                          | `https://slsa.dev/verification_summary/v1`        |
| SPDX SBOM                  | An SBOM in [SPDX](https://spdx.dev/) format, widely adopted in open-source ecosystems.                                                                                                   | `https://spdx.dev/Document`                       |
| FIPS compliance            | An attestation that verifies the image uses FIPS 140-validated cryptographic modules.                              | `https://docker.com/dhi/fips/v0.1`                |

## View and verify attestations

To view and verify attestations for an image, see [Verify a Docker Hardened
Image](../how-to/verify.md).

## Add your own attestations

In addition to the comprehensive attestations provided by Docker Hardened
Images, you can add your own signed attestations when building derivative
images. This is especially useful if you’re building new applications on top of
a DHI and want to maintain transparency, traceability, and trust in your
software supply chain.

By attaching attestations such as SBOMs, build provenance, or custom metadata,
you can meet compliance requirements, pass security audits, and support policy
evaluation tools like Docker Scout.

These attestations can then be verified downstream using tools
like Cosign or Docker Scout.

To learn how to attach custom attestations during the build process, see [Build
attestations](/manuals/build/metadata/attestations.md).
