---
title: FIPS
description: Learn how Docker Hardened Images support FIPS 140 through validated cryptographic modules to help organizations meet compliance requirements.
keywords: docker fips, fips 140 images, fips docker images, docker compliance, secure container images
---

## What is FIPS 140?

[FIPS 140](https://csrc.nist.gov/publications/detail/fips/140/3/final) is a U.S.
government standard that defines security requirements for cryptographic modules
that protect sensitive information. It is widely used in regulated environments
such as government, healthcare, and financial services.

FIPS certification is managed by the [NIST Cryptographic Module Validation
Program
(CMVP)](https://csrc.nist.gov/projects/cryptographic-module-validation-program),
which ensures cryptographic modules meet rigorous security standards.

## Why FIPS compliance matters

FIPS 140 compliance is required or strongly recommended in many regulated
environments where sensitive data must be protected, such as government,
healthcare, finance, and defense. These standards ensure that cryptographic
operations are performed using vetted, trusted algorithms implemented in secure
modules.

Using software components that rely on validated cryptographic modules can help organizations:

- Satisfy federal and industry mandates, such as FedRAMP, which require or
  strongly recommend FIPS 140-validated cryptography.
- Demonstrate audit readiness, with verifiable evidence of secure,
  standards-based cryptographic practices.
- Reduce security risk, by blocking unapproved or unsafe algorithms (e.g., MD5)
  and ensuring consistent behavior across environments.

## How Docker Hardened Images support FIPS compliance

Docker Hardened Images (DHIs) include variants that use cryptographic modules
validated under FIPS 140. These images are intended to help organizations meet
compliance requirements by incorporating components that meet the standard.

- FIPS image variants use cryptographic modules that are already validated under
  FIPS 140.
- These variants are built and maintained by Docker to support environments with
  regulatory or compliance needs.
- Docker provides signed test attestations that document the use of validated
  cryptographic modules. These attestations can support internal audits and
  compliance reporting.

> [!NOTE]
>
> Using a FIPS image variant helps meet compliance requirements but does not
> make an application or system fully compliant. Compliance depends on how the
> image is integrated and used within the broader system.

## Identify images that support FIPS

Docker Hardened Images that support FIPS are marked as **FIPS** compliant
in the Docker Hardened Images catalog.

To find DHI repositories with FIPS image variants, [explore images](../how-to/explore.md) and:

- Use the **FIPS** filter on the catalog page
- Look for **FIPS** compliant on individual image listings

These indicators help you quickly locate repositories that support FIPS-based
compliance needs. Image variants that include FIPS support will have a tag
ending with `-fips`, such as `3.13-fips`.

## Validate FIPS-related tests using attestations

Docker Hardened Images include a signed [test
attestation](../core-concepts/attestations.md) that documents the results of
automated image validation. For FIPS variants, this includes test cases that
verify whether the image uses FIPS-validated cryptographic modules.

You can retrieve and inspect this attestation using the Docker Scout CLI:

```console
$ docker scout attest get \
  --predicate-type https://scout.docker.com/tests/v0.1 \
  --predicate \
  <your-namespace>/dhi-<image>:<tag> --platform <platform>
```

For example:

```console
$ docker scout attest get \
  --predicate-type https://scout.docker.com/tests/v0.1 \
  --predicate \
  docs/dhi-python:3.13-fips --platform linux/amd64
```

The output is a structured JSON report. Individual test outputs are
base64-encoded under fields like `stdout`. You can decode them to review the raw
test output.

To decode and view test results:

```console
$ docker scout attest get \
  --predicate-type https://scout.docker.com/tests/v0.1 \
  --predicate \
  docs/dhi-python:3.13-fips --platform linux/amd64 \
  | jq -r '.results.tests[].extra.stdout' \
  | base64 -d
```