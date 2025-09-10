---
title: Security for developers
linkTitle: Security
description: Learn about developer-level security features like 2FA and access tokens
keywords: docker, docker hub, docker desktop, security, developer security, 2FA, access tokens
weight: 40
params:
  sidebar:
    group: Platform
grid_developers:
- title: Set up two-factor authentication
  description: Add an extra layer of authentication to your Docker account.
  link: /security/2fa/
  icon: phonelink_lock
- title: Manage access tokens
  description: Create personal access tokens as an alternative to your password.
  icon: password
  link: /security/access-tokens/
- title: Static vulnerability scanning
  description: Automatically run a point-in-time scan on your Docker images for vulnerabilities.
  icon: image_search
  link: /docker-hub/repos/manage/vulnerability-scanning/
- title: Docker Engine security
  description: Understand how to keep Docker Engine secure.
  icon: security
  link: /engine/security/
- title: Secrets in Docker Compose
  description: Learn how to use secrets in Docker Compose.
  icon: privacy_tip
  link: /compose/how-tos/use-secrets/
grid_resources:
- title: Security FAQs
  description: Explore common security FAQs.
  icon: help
  link: /faq/security/general/
- title: Security best practices
  description: Understand the steps you can take to improve the security of your container.
  icon: category
  link: /develop/security-best-practices/
- title: Suppress CVEs with VEX
  description: Learn how to suppress non-applicable or fixed vulnerabilities found in your images.
  icon: query_stats
  link: /scout/guides/vex/
- title: Docker Hardened Images
  description: Learn how to use Docker Hardened Images to enhance your software supply security.
  icon: encrypted_add_circle
  link: /dhi/
---

Docker helps you protect your local environments, infrastructure, and networks
with its developer-level security features.

Use tools like two-factor authentication (2FA), personal access tokens, and
Docker Scout to manage access and detect vulnerabilities early in your workflow.
You can also integrate secrets securely into your development stack using Docker Compose,
or enhance your software supply security with Docker Hardened Images.

Explore the following sections to learn more.

## For developers

{{< grid items="grid_developers" >}}

## More resources

{{< grid items="grid_resources" >}}
