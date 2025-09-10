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
      - AI
      - Products
      - Platform
      - Enterprise
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
    icon: /icons/Compose.svg
    link: /compose/
  - title: Testcontainers
    description: Run containers programmatically in your preferred programming language.
    icon: /icons/Testcontainers.svg
    link: /testcontainers/
  - title: MCP Gateway
    description: Manage and secure your AI tools with a single gateway.
    icon: /icons/toolkit.svg
    link: /ai/mcp-gateway/
  - title: Cagent
    description: The open-source multi-agent solution to assist you in your tasks.
    icon: /icons/cagent.svg
    link: /ai/cagent
  ai:
  - title: Ask Gordon
    description: Streamline your workflow and get the most out of the Docker ecosystem with your personal AI assistant.
    icon: note_add
    link: /ai/gordon/
  - title: Docker Model Runner
    description: View and manage your local models.
    icon: /icons/models.svg
    link: /ai/model-runner/
  - title: MCP Catalog and Toolkit
    description: Augment your AI workflow with MCP servers.
    icon: /icons/toolkit.svg
    link: /ai/mcp-catalog-and-toolkit/
  products:
  - title: Docker Desktop
    description: Your command center for container development.
    icon: /icons/Whale.svg
    link: /desktop/
  - title: Docker Hardened Images
    description: Secure, minimal images for trusted software delivery.
    icon: /icons/dhi.svg
    link: /dhi/
  - title: Docker Offload
    description: Build and run containers in the cloud.
    icon: cloud
    link: /offload/
  - title: Build Cloud
    description: Build your images faster in the cloud.
    icon: /icons/logo-build-cloud.svg
    link: /build-cloud/
  - title: Docker Hub
    description: Discover, share, and integrate container images.
    icon: hub
    link: /docker-hub/
  - title: Docker Scout
    description: Image analysis and policy evaluation.
    icon: /icons/Scout.svg
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
  enterprise:
  - title: Deploy Docker Desktop
    description: Deploy Docker Desktop at scale within your company
    icon: download
    link: /enterprise/enterprise-deployment/
---

This section contains user guides on how to install, set up, configure, and use
Docker products.

## Open source

Open source development and containerization technologies.

{{< grid items=open-source >}}

## AI

All the Docker AI tools in one easy-to-access location.

{{< grid items=ai >}}

## Products

End-to-end developer solutions for innovative teams.

{{< grid items=products >}}

## Platform

Documentation related to the Docker platform, such as administration and
subscription management.

{{< grid items=platform >}}

## Enterprise

Targeted at IT administrators with help on deploying Docker Desktop at scale with configuration guidance on security related features.

{{< grid items=enterprise >}}