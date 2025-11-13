---
title: Continuous patching and secure maintenance
linktitle: Continuous patching
description: Learn how Docker Hardened Images are automatically rebuilt, tested, and updated to stay in sync with upstream security patches.
keywords: docker hardened images, secure base image, automatic patching, CVE updates, compatibility, dev containers, runtime containers, image maintenance
---

Docker Hardened Images (DHI) offer a secure and enterprise-ready foundation for
containerized applications, backed by a robust, automated patching process that
helps maintain compliance and reduce vulnerability exposure.

## Secure base images with strong compatibility

DHI includes a curated set of minimal base images designed to work across a
broad range of environments and language ecosystems. These images provide secure
building blocks with high compatibility, making it easier to integrate into your
existing infrastructure and development workflows without sacrificing security.

## Development and runtime variants

To support different stages of the software lifecycle, DHI provides two key
variants:

- Development images: Include essential tools and libraries required to build
  and test applications securely.
- Runtime images: Contain only the core components needed to run applications,
  offering a smaller attack surface and improved runtime efficiency.

This variant structure supports multi-stage builds, enabling developers to
compile code in secure development containers and deploy with lean runtime
images in production.

## Automated patching and secure updates

Docker monitors upstream open-source packages and security advisories for
vulnerabilities (CVEs) and other updates. When changes are detected, affected
Docker Hardened Images are automatically rebuilt and tested.

Updated images are published with cryptographic provenance attestations to
support verification and compliance workflows. This automated process reduces
the operational burden of manual patching and helps teams stay aligned with
secure software development practices.

## Automatic patching for customized images

When you [customize a Docker Hardened Image](../how-to/customize.md), your
customized images also benefit from automatic patching. When the base Docker
Hardened Image receives a security update, Docker automatically rebuilds your
customized images in the background, ensuring they stay current with the latest
security patches without requiring manual intervention.

This means your customizations maintain continuous compliance and protection by
default, with no additional operational overhead.