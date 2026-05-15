---
title: "Lab: Container Image Attestations"
linkTitle: "Lab: Container Image Attestations"
description: |
  Learn to attach SBOMs, build provenance, image signatures, and VEX
  statements to container images for a verifiable software supply chain.
summary: |
  Hands-on lab: Add supply chain metadata to a container image. Generate
  SBOMs and SLSA provenance with BuildKit, sign images with Cosign, and
  attach OpenVEX statements to declare vulnerability exploitability status.
keywords: Docker, supply chain, SBOM, provenance, SLSA, Cosign, VEX, attestations, security, lab, labspace
params:
  tags: [labs]
  time: 45 minutes
  resource_links:
    - title: Build attestations
      url: /build/metadata/attestations/
    - title: SBOM attestations
      url: /build/metadata/attestations/sbom/
    - title: Provenance attestations
      url: /build/metadata/attestations/slsa-provenance/
    - title: Labspace repository
      url: https://github.com/dockersamples/labspace-attestation-basics
---

Prove where your container images came from and that they haven't been
tampered with. This lab walks through generating SBOMs and SLSA build
provenance with BuildKit, signing images with Cosign, and writing VEX
statements to declare which CVEs affect your image — the techniques used
to meet supply chain security requirements like NIST SSDF and EO 14028.

## Launch the lab

{{< labspace-launch image="dockersamples/labspace-attestation-basics" >}}

## What you'll learn

By the end of this Labspace, you will have completed the following:

- Generate and inspect an SPDX SBOM attached to a container image with `--sbom=true`
- Generate SLSA build provenance with `--provenance=mode=max` and understand how multi-stage builds are fully recorded
- Install Cosign and use key-based signing to sign and verify a container image
- Write an OpenVEX statement to declare CVE exploitability status and attach it as a signed attestation
- Understand how SBOMs, provenance, signatures, and VEX complement each other in a complete supply chain story

## Modules

| #   | Module                            | Description                                                                          |
| --- | --------------------------------- | ------------------------------------------------------------------------------------ |
| 1   | Introduction                      | Overview of supply chain attestations and the sample Go app                          |
| 2   | Software Bill of Materials (SBOM) | Build with `--sbom=true`, inspect SPDX contents, and understand scanner integration  |
| 3   | Build Provenance                  | Generate SLSA provenance and explore how multi-stage builds are recorded             |
| 4   | Signing Images with Cosign        | Generate a key pair, sign the image, verify the signature, and learn keyless signing |
| 5   | VEX Statements                    | Scan for CVEs, write an OpenVEX document, and attach it as a signed attestation      |
| 6   | Bringing It All Together          | Run the complete build-sign-attest workflow and see the full supply chain picture    |
| 7   | Recap                             | Summary of skills and next steps for policy enforcement and higher SLSA levels       |
