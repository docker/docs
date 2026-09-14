---
title: Manuals
description: Learn how to install, set up, configure, and use Docker products with this collection of user guides
keywords: docker, docs, manuals, products, user guides, how-to
# hard-code the URL of this page
url: /manuals/
layout: wide
params:
  icon: document-text
  sidebar:
    groups:
      - AI and agents
      - Application development
      - Supply chain security
      - Accounts and admin
      - Enterprise
  notoc: true
  ai-and-agents:
  - title: Docker Sandboxes
    description: Run AI coding agents in isolated environments.
    icon: command-line
    link: /ai/sandboxes/
  - title: MCP Catalog and Toolkit
    description: Augment your AI workflow with MCP servers.
    icon: /icons/toolkit.svg
    link: /ai/mcp-catalog-and-toolkit/
  - title: Gordon
    description: Streamline your workflow and get the most out of the Docker ecosystem with your personal AI assistant.
    icon: document-plus
    link: /ai/gordon/
  - title: Docker Model Runner
    description: View and manage your local models.
    icon: /icons/models.svg
    link: /ai/model-runner/
  - title: Docker Agent
    description: The open-source multi-agent solution to assist you in your tasks.
    icon: /icons/cagent.svg
    link: /ai/docker-agent
  application-development:
  - title: Docker Desktop
    description: Your command center for container development.
    icon: /icons/Whale.svg
    link: /desktop/
  - title: Docker Offload
    description: Build and run containers in the cloud.
    icon: cloud
    link: /offload/
  - title: Docker Build Cloud
    description: Build your images faster in the cloud.
    icon: /icons/logo-build-cloud.svg
    link: /build-cloud/
  - title: Testcontainers
    description: Run containers programmatically in your preferred programming language.
    icon: /icons/Testcontainers.svg
    link: /testcontainers/
  - title: Docker Build
    description: Build and ship any application anywhere.
    icon: wrench-screwdriver
    link: /build/
  - title: Docker Engine
    description: The industry-leading container runtime.
    icon: cpu-chip
    link: /engine/
  - title: Docker Compose
    description: Define and run multi-container applications.
    icon: /icons/Compose.svg
    link: /compose/
  supply-chain-security:
  - title: Docker Hub
    description: Discover, share, and integrate container images.
    icon: globe-alt
    link: /docker-hub/
  - title: Docker Hardened Images
    description: Secure, minimal images for trusted software delivery.
    icon: /icons/dhi.svg
    link: /dhi/
  - title: Docker Scout
    description: Image analysis and policy evaluation.
    icon: /icons/Scout.svg
    link: /scout/
  platform:
  - title: Accounts
    description: Manage Docker individual and organization accounts.
    icon: user-circle
    link: /accounts/
  - title: Subscription and billing
    description: Manage Docker subscriptions, plans, billing, and payments.
    icon: credit-card
    link: /subscription-billing/
  - title: Security
    description: Security guardrails for both administrators and developers.
    icon: lock-closed
    link: /security/
  - title: FAQs
    description: Frequently asked questions about Docker accounts, organizations, companies, subscriptions, billing, and security.
    icon: question-mark-circle
    link: /faqs/
  - title: Support
    description: Support options for paid subscriptions and community resources.
    icon: chat-bubble-left
    link: /support/
  - title: Release notes
    description: Features, bug fixes, and breaking changes for Docker Home, billing, security, and subscriptions.
    icon: document-plus
    link: /platform-release-notes/
  enterprise:
  - title: Deploy Docker Desktop
    description: Deploy Docker Desktop at scale within your company
    icon: arrow-down-tray
    link: /enterprise/enterprise-deployment/
  - title: Hardened Docker Desktop
    description: Security features that strengthen developer environments.
    icon: shield-check
    link: /enterprise/security/hardened-desktop/
---

This section contains user guides on how to install, set up, configure, and use
Docker products.

## AI and agents

All the Docker AI tools in one easy-to-access location.

{{< grid items=ai-and-agents >}}

## Application development

End-to-end developer solutions for innovative teams.

{{< grid items=application-development >}}

## Supply chain security

Security guardrails and image analysis for your software supply chain.

{{< grid items=supply-chain-security >}}

## Accounts and admin

Manage Docker accounts, administration, subscriptions, billing, and security.

{{< grid items=platform >}}

## Enterprise

Targeted at IT administrators with help on deploying Docker Desktop at scale with configuration guidance on security related features.

{{< grid items=enterprise >}}
