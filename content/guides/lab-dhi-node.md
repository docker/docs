---
title: "Lab: Migrating a Node App to Docker Hardened Images"
linkTitle: "Lab: Migrating to DHI (Node)"
description: |
  Migrate a Node.js application from a standard base image to Docker Hardened
  Images. Use Docker Scout to analyze CVEs, compare images, and inspect
  supply chain attestations.
summary: |
  Hands-on lab: Replace a Node.js base image with a Docker Hardened Image.
  Analyze CVEs with Docker Scout, rewrite the Dockerfile to use multi-stage
  builds with DHI, and explore SBOMs, VEX, and compliance attestations.
keywords: Docker, Hardened Images, DHI, Node.js, Docker Scout, CVE, security, SBOM, lab, labspace
params:
  tags: [labs, dhi]
  time: 30 minutes
  resource_links:
    - title: Docker Hardened Images
      url: /dhi/
    - title: Docker Scout docs
      url: /scout/
    - title: Build attestations
      url: /build/metadata/attestations/
    - title: Labspace repository
      url: https://github.com/dockersamples/labspace-dhi-node
---

Migrate a Node.js application from a standard `node:24-trixie-slim` base image
to a Docker Hardened Image. You'll measure the before-and-after impact on CVE
count, image size, and policy compliance using Docker Scout, then explore the
supply chain attestations DHI ships with every image.

## Launch the lab

{{< labspace-launch image="dockersamples/labspace-dhi-node" >}}

## What you'll learn

By the end of this Labspace, you will have completed the following:

- Analyze a Node.js container image with Docker Scout to identify CVE and policy failures
- Rewrite a Dockerfile to use a multi-stage build with DHI dev and runtime variants
- Compare image size and vulnerability counts before and after the migration
- Inspect supply chain attestations included with Docker Hardened Images (SBOMs, SLSA, VEX)
- Export VEX documents for integration with external scanners such as Grype or Trivy

## Modules

| #   | Module                                   | Description                                                                     |
| --- | ---------------------------------------- | ------------------------------------------------------------------------------- |
| 1   | Introduction                             | Overview of Docker Hardened Images and their security benefits                  |
| 2   | Setup                                    | Perform setup tasks required for the lab.                                       |
| 3   | Analyzing the Starting Image             | Build the app, scan it with Docker Scout, and review failing policies           |
| 4   | Migrating to DHI                         | Rewrite the Dockerfile with multi-stage DHI build and compare results           |
| 5   | DHI Attestations and Scanner Integration | Inspect SBOMs, FIPS attestations, STIG scans, and export VEX for external tools |
