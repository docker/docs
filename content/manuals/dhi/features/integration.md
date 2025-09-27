---
title: Seamless integration
description: Learn how Docker Hardened Images integrate into your existing development and deployment workflows for enhanced security without compromising usability.
description_short: See how Docker Hardened Images integrate with CI/CD pipelines, vulnerability scanners, and container registries across your toolchain
keywords: ci cd containers, vulnerability scanning, slsa build level 3, signed sbom, oci compliant registry
---

Docker Hardened Images (DHI) are designed to integrate effortlessly into your
existing development and deployment workflows, ensuring that enhanced security
does not come at the cost of usability.

## Explore images in Docker Hub

After your organization [signs
up](https://www.docker.com/products/hardened-images/#getstarted), teams can
explore the full DHI catalog directly on Docker Hub. There, developers and
security teams can:

- Review available images and language/framework variants
- Understand supported distros
- Compare development vs. runtime variants

Each repository includes metadata like supported tags, base image
configurations, and image-specific documentation, helping you choose the right variant
for your project.

## Use DHIs in CI/CD workflows

You can use DHIs as the same base image in any CI/CD pipeline that is built
using a Dockerfile. They integrate easily into platforms like GitHub Actions,
GitLab CI/CD, Jenkins, CircleCI, and other automation systems your team already
uses.

## Built to fit your DevSecOps stack

Docker Hardened Images are designed to work seamlessly with your existing
DevSecOps toolchain. They integrate with scanning tools, registries, CI/CD
systems, and policy engines that teams already use.

Docker has partnered with a broad range of ecosystem providers in order to
ensure that DHIs work out of the box with your existing workflows and tools.
These partners help deliver enhanced scanning, metadata validation, and
compliance insights directly into your pipelines.

All DHIs include:

- Signed Software Bill of Materials (SBOMs)
- CVE data
- Vulnerability Exploitability eXchange (VEX) documents
- SLSA Build Level 3 provenance

Because the metadata is signed and structured, you can feed it into policy
engines and dashboards for auditing or compliance workflows.

## Distribute through your preferred registry

DHIs are mirrored to your organization's namespace on Docker Hub. From there,
you can optionally push them to any OCI-compliant registry, such as:

- Amazon ECR
- Google Artifact Registry
- GitHub Container Registry
- Azure Container Registry
- Harbor
- JFrog Artifactory
- Other OCI-compliant on-premises or cloud registries

Mirroring ensures teams can pull images from their preferred location without
breaking policies or build systems.

## Summary

Docker Hardened Images integrate with the tools you already use, from development
and CI to scanning and deployment. They:

- Work with standard Docker tooling and pipelines
- Support popular scanners and registries
- Include security metadata that plugs into your existing compliance systems

This means you can adopt stronger security controls without disrupting your
engineering workflows.
