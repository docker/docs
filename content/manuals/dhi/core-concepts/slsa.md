---
title: Supply-chain Levels for Software Artifacts (SLSA)
linktitle: SLSA
description: Learn how Docker Hardened Images comply with SLSA Build Level 3 and how to verify provenance for secure, tamper-resistant builds.
keywords: slsa docker compliance, slsa build level 3, supply chain security, verified build provenance, secure container build
---

## What is SLSA?

Supply-chain Levels for Software Artifacts (SLSA) is a security framework
designed to enhance the integrity and security of software supply chains.
Developed by Google and maintained by the Open Source Security Foundation
(OpenSSF), SLSA provides a set of guidelines and best practices to prevent
tampering, improve integrity, and secure packages and infrastructure in software
projects.

SLSA defines [four build levels (0–3)](https://slsa.dev/spec/latest/build-track-basics) of
increasing security rigor, focusing on areas such as build provenance, source
integrity, and build environment security. Each level builds upon the previous
one, offering a structured approach to achieving higher levels of software
supply chain security.

## Why is SLSA important?

SLSA is crucial for modern software development due to the increasing complexity
and interconnectedness of software supply chains. Supply chain attacks, such as
the SolarWinds breach, have highlighted the vulnerabilities in software
development processes. By implementing SLSA, organizations can:

- Ensure artifact integrity: Verify that software artifacts have not been
  tampered with during the build and deployment processes.

- Enhance build provenance: Maintain verifiable records of how and when software
  artifacts were produced, providing transparency and accountability.

- Secure build environments: Implement controls to protect build systems from
  unauthorized access and modifications.

- Mitigate supply chain risks: Reduce the risk of introducing vulnerabilities or
  malicious code into the software supply chain.

## What is SLSA Build Level 3?

SLSA Build Level 3, Hardened Builds, is the highest of four progressive levels in
the SLSA framework. It introduces strict requirements to ensure that software
artifacts are built securely and traceably. To meet Level 3, a build must:

- Be fully automated and scripted to prevent manual tampering
- Use a trusted build service that enforces source and builder authentication
- Generate a signed, tamper-resistant provenance record describing how the artifact was built
- Capture metadata about the build environment, source repository, and build steps

This level provides strong guarantees that the software was built from the
expected source in a controlled, auditable environment, which significantly
reduces the risk of supply chain attacks.

## Docker Hardened Images and SLSA

Docker Hardened Images (DHIs) are secure-by-default container images
purpose-built for modern production environments. Each DHI is cryptographically
signed and complies with the [SLSA Build Level 3
standard](https://slsa.dev/spec/latest/build-track-basics#build-l3), ensuring
verifiable build provenance and integrity.

By integrating SLSA-compliant DHIs into your development and deployment processes, you can:

- Achieve higher security levels: Utilize images that meet stringent security
  standards, reducing the risk of vulnerabilities and attacks.

- Simplify compliance: Leverage built-in features like signed Software Bills of
  Materials (SBOMs) and vulnerability exception (VEX) statements to facilitate
  compliance with regulations such as FedRAMP.

- Enhance transparency: Access detailed information about the components and
  build process of each image, promoting transparency and trust.

- Streamline audits: Utilize verifiable build records and signatures to simplify
  security audits and assessments.

## Get and verify SLSA provenance for Docker Hardened Images

Each Docker Hardened Image (DHI) is cryptographically signed and includes
attestations. These attestations provide verifiable build provenance and
demonstrate adherence to SLSA Build Level 3 standards.

To get and verify SLSA provenance for a DHI, you can use Docker Scout.

```console
$ docker scout attest get dhi.io/<image>:<tag> \
  --predicate-type https://slsa.dev/provenance/v0.2 \
  --verify
```

For example:

```console
$ docker scout attest get dhi.io/node:20.19-debian12 \
  --predicate-type https://slsa.dev/provenance/v0.2 \
  --verify
```

## Resources

For more details about SLSA definitions and Docker Build, see [SLSA
definitions](/build/metadata/attestations/slsa-definitions/).