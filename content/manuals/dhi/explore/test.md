---
title: How Docker Hardened Images are tested
linktitle: Image testing
description: See how Docker Hardened Images are automatically tested for standards compliance, functionality, and security.
keywords: docker scout, test attestation, cosign verify, image testing, vulnerability scan
weight: 45
---

Docker Hardened Images (DHIs) are designed to be secure, minimal, and
production-ready. To ensure their reliability and security, Docker employs a
comprehensive testing strategy, which you can independently verify using signed
attestations and open tooling.

Every image is tested for standards compliance, functionality, and security. The
results of this testing are embedded as signed attestations, which can be
[inspected and verified](#view-and-verify-the-test-attestation) programmatically
using the Docker Scout CLI.

## Testing strategy overview

The testing process for DHIs focuses on two main areas:

- Image standards compliance: Ensuring that each image adheres to strict size,
  security, and compatibility standards.
- Application functionality: Verifying that applications within the images
  function correctly.

## Image standards compliance

Each DHI undergoes rigorous checks to meet the following standards:

- Minimal attack surface: Images are built to be as small as possible, removing
  unnecessary components to reduce potential vulnerabilities.
- Near-zero known CVEs: Images are scanned using tools like Docker Scout to
  ensure they are free from known Common Vulnerabilities and Exposures (CVEs).
- Multi-architecture support: DHIs are built for multiple architectures
  (`linux/amd64` and `linux/arm64`) to ensure broad compatibility.
- Kubernetes compatibility: Images are tested to run seamlessly within
  Kubernetes clusters, ensuring they meet the requirements for container
  orchestration environments.

## Application functionality testing

Docker tests Docker Hardened Images to ensure they behave as expected in typical
usage scenarios. This includes verifying that:

- Applications start and run successfully in containerized environments.
- Runtime behavior aligns with upstream expectations.
- Build variants (like `-dev` images) support common development and build tasks.

The goal is to ensure that DHIs work out of the box for the most common use
cases while maintaining the hardened, minimal design.

## Automated testing and CI/CD integration

Docker integrates automated testing into its Continuous Integration/Continuous
Deployment (CI/CD) pipelines:

- Automated scans: Each image build triggers automated scans for vulnerabilities
  and compliance checks.
- Reproducible builds: Build processes are designed to be reproducible, ensuring
  consistency across different environments.
- Continuous monitoring: Docker continuously monitors for new vulnerabilities
  and updates images accordingly to maintain security standards.

## Testing attestation

Docker provides a test attestation that details the testing and validation
processes each DHI has undergone.

### View and verify the test attestation

You can view and verify this attestation using the Docker Scout CLI.

1. Use the `docker scout attest get` command with the test predicate type:

   ```console
   $ docker scout attest get \
     --predicate-type https://scout.docker.com/tests/v0.1 \
     --predicate \
     dhi.io/<image>:<tag>
   ```

   > [!NOTE]
   >
   > If the image exists locally on your device, you must prefix the image name with `registry://`. For example, use
   > `registry://dhi.io/python` instead of `dhi.io/python`.

   For example:

   ```console
   $ docker scout attest get \
     --predicate-type https://scout.docker.com/tests/v0.1 \
     --predicate \
     dhi.io/python:3.13
   ```

   This contains a list of tests and their results.

   Example output:

    ```console
        v SBOM obtained from attestation, 101 packages found
        v Provenance obtained from attestation
        {
          "reportFormat": "CTRF",
          "results": {
            "summary": {
              "failed": 0,
              "passed": 1,
              "skipped": 0,
              "start": 1749216533,
              "stop": 1749216574,
              "tests": 1
            },
            "tests": [
              {
                ...
   ```

2. Verify the test attestation signature. To ensure the attestation is authentic
   and signed by Docker, run:

   ```console
   docker scout attest get \
     --predicate-type https://scout.docker.com/tests/v0.1 \
     --verify \
     dhi.io/<image>:<tag> --platform <platform>
   ```

   Example output:
   
   ```console
    v SBOM obtained from attestation, 101 packages found
    v Provenance obtained from attestation
    v cosign verify registry.scout.docker.com/docker/dhi-python@sha256:70c8299c4d3cb4d5432734773c45ae58d8acc2f2f07803435c65515f662136d5 \
        --key https://registry.scout.docker.com/keyring/dhi/latest.pub --experimental-oci11

      Verification for registry.scout.docker.com/docker/dhi-python@sha256:70c8299c4d3cb4d5432734773c45ae58d8acc2f2f07803435c65515f662136d5 --
      The following checks were performed on each of these signatures:
        - The cosign claims were validated
        - Existence of the claims in the transparency log was verified offline
        - The signatures were verified against the specified public key

    i Signature payload
    ...
    ```

If the attestation is valid, Docker Scout will confirm the signature and show
the matching `cosign verify` command.

To view other attestations, such as SBOMs or vulnerability reports, see [Verify
an image](../how-to/verify.md).
