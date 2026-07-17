---
title: Resources and feedback
linktitle: Resources and feedback
description: Additional resources, community links, GitHub repositories, and how to give feedback for Docker Hardened Images.
keywords: docker hardened images resources, dhi feedback, dhi github, dhi community, report issue, dhi support
weight: 999
aliases:
  - /dhi/about/feedback/
  - /dhi/explore/feedback/
---

This page provides links to additional resources, community channels, and ways
to give feedback on Docker Hardened Images (DHI).

For product information and feature comparison, visit the [Docker Hardened
Images product page](https://www.docker.com/products/hardened-images/).

## Guides

For guides that demonstrate how to use Docker Hardened Images in various
scenarios, see the [guides section filtered by DHI](/guides/?tags=dhi).

## Docker Hub

Docker Hardened Images are available on Docker Hub:

- [Docker Hardened Images Catalog](https://dhi.io): Browse and pull Docker
  Hardened Images from the official catalog
- [Docker Hub MCP Server](https://hub.docker.com/mcp/server/dockerhub/overview):
  MCP server to list Docker Hardened Images (DHIs) available in your
  organizations

## GitHub repositories and resources

Docker Hardened Images repositories are available in the
[docker-hardened-images](https://github.com/docker-hardened-images) GitHub
organization:

- [Catalog](https://github.com/docker-hardened-images/catalog): DHI definition
  files and catalog metadata
- [Advisories](https://github.com/docker-hardened-images/advisories): CVE
  advisories for OSS packages distributed with DHIs
  - [Scanner vendor integration guide](https://github.com/docker-hardened-images/advisories/tree/main/integration):
    Reference for scanner vendors integrating DHI VEX support
- [Keyring](https://github.com/docker-hardened-images/keyring): Public signing
  keys and verification tools
- [Log](https://github.com/docker-hardened-images/log): Log of references (tag >
  digest) for Docker Hardened Images
- [dhictl](https://github.com/docker-hardened-images/dhictl): Command-line
  interface for managing and interacting with Docker Hardened Images
- [Terraform Provider](https://github.com/docker-hardened-images/terraform-provider-dhi):
  Terraform provider for managing DHI resources
  ([Terraform Registry](https://registry.terraform.io/providers/docker-hardened-images/dhi/latest/docs))

## Additional resources

- [Start a free trial](https://hub.docker.com/hardened-images/start-free-trial):
  Explore DHI Select and Enterprise features including FIPS/STIG variants, customization,
  and SLA-backed support
- [Support Service Level Agreement](https://docs.docker.com/go/dhi-sla/):
  Review the SLA commitments for DHI Select and Enterprise subscriptions
- [Request a demo](https://www.docker.com/products/hardened-images/#getstarted): Get a
  personalized demo and information about DHI Select and Enterprise subscriptions
- [Request an image](https://github.com/docker-hardened-images/catalog/issues):
  Submit a request for a specific Docker Hardened Image
- [Debian package index](https://dhi.io/deb/debian/main/index.html):
  Browse hardened Debian packages in Docker's public repository
- [Alpine package index](https://dhi.io/apk/alpine/v3.24/main/index.html):
  Browse hardened Alpine packages in Docker's public repository
- <a href="https://www.docker.com/pricing/contact-sales/" id="dkr_docs_cs_dhi_resources" class="link" rel="noopener">Contact Sales</a>: Connect with
  Docker sales team for enterprise inquiries
- [Docker Support](https://www.docker.com/support/): Access support resources
  for DHI Select and Enterprise customers

## Feedback and community

Use the [GitHub Discussions board](https://github.com/orgs/docker-hardened-images/discussions)
to engage with the DHI team for general questions, best practices, security
tips, and community announcements.

To report a bug, request a feature, or suggest a documentation improvement,
[open an issue](https://github.com/docker-hardened-images/catalog/issues) in
the catalog repository.

## Security disclosure

Do not post details of vulnerabilities before coordinated disclosure and
resolution. If you discover a security vulnerability, report it responsibly
by following Docker's [security disclosure
policy](https://www.docker.com/trust/vulnerability-disclosure-policy/).
