---
title: Security for developers
linkTitle: Security 
description: Learn about developer-level security features Docker has to offer and explore best practices
keywords: docker, docker hub, docker desktop, security
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
---

Docker provides security guardrails for both administrators and developers.

If you're an administrator, you can enforce sign-in across Docker products for your developers, and
scale, manage, and secure your instances of Docker Desktop with DevOps security controls like Enhanced Container Isolation and Registry Access Management.

For both administrators and developers, Docker provides security-specific products such as Docker Scout, for securing your software supply chain with proactive image vulnerability monitoring and remediation strategies.

## For developers

See how you can protect your local environments, infrastructure, and networks without impeding productivity.

{{< grid items="grid_developers" >}}

## Further resources

{{< grid items="grid_resources" >}}
