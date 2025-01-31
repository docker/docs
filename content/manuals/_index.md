---
title: Manuals
description: Learn how to install, set up, configure, and use Docker products with this collection of user guides
keywords: docker, docs, manuals, products, user guides, how-to
# hard-code the URL of this page
url: /manuals/
layout: wide
params:
  icon: description
  sidebar:
    groups:
      - Open source
      - Products
      - Platform
  notoc: true
  open-source:
  - title: Docker Build
    description: Build and ship any application anywhere.
    icon: build
    link: /build/
  - title: Docker Engine
    description: The industry-leading container runtime.
    icon: developer_board
    link: /engine/
  - title: Docker Compose
    description: Define and run multi-container applications.
    icon: /assets/icons/Compose.svg
    link: /compose/
  products:
  - title: Docker Desktop
    description: Your command center for container development.
    icon: /assets/icons/Whale.svg
    link: /desktop/
  - title: Build Cloud
    description: Build your images faster in the cloud.
    icon: /assets/images/logo-build-cloud.svg
    link: /build-cloud/
  - title: Docker Hub
    description: Discover, share, and integrate container images.
    icon: hub
    link: /docker-hub/
  - title: Docker Scout
    description: Image analysis and policy evaluation.
    icon: /assets/icons/Scout.svg
    link: /scout/
  - title: Docker for GitHub Copilot
    description: Integrate Docker's capabilities with GitHub Copilot.
    icon: chat
    link: /copilot/
  - title: Docker Extensions
    description: Customize your Docker Desktop workflow.
    icon: extension
    link: /extensions/
  - title: Testcontainers Cloud
    description: Run integration tests, with real dependencies, in the cloud.
    icon: package_2
    link: https://testcontainers.com/cloud/docs/
  platform:
  - title: Administration
    description: Centralized observability for companies and organizations.
    icon: admin_panel_settings
    link: /admin/
  - title: Billing
    description: Manage billing and payment methods.
    icon: payments
    link: /billing/
  - title: Accounts
    description: Manage your Docker account.
    icon: account_circle
    link: /accounts/
  - title: Security
    description: Security guardrails for both administrators and developers.
    icon: lock
    link: /security/
  - title: Subscription
    description: Commercial use licenses for Docker products.
    icon: card_membership
    link: /subscription/
---

This section contains user guides on how to install, set up, configure, and use
Docker products.

## Open source

Open source development and containerization technologies.

{{< grid items=open-source >}}

## Products

End-to-end developer solutions for innovative teams.

{{< grid items=products >}}

## Platform

Documentation related to the Docker platform, such as administration and
subscription management for organizations.

{{< grid items=platform >}}
