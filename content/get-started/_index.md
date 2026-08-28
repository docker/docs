---
title: Get started
keywords: Docker, get started, containers, AI agents, sandboxes, administration
description: Choose a Docker journey to build and run applications, work with AI agents safely, or set up Docker for your organization.
layout: get-started
params:
  notoc: true
  sidebar:
    hidePage: true
    include:
      - Run an application with Docker
      - Run an AI agent in a sandbox
  journeys:
    - persona: New to Docker
      label: Container journey
      title: Build and run applications
      description: Start with a container, package an application as an image, and run a complete development stack.
      icon: cube-transparent
      length: 3 tutorials
      guides:
        - title: Run a container
          description: Start software from an existing image.
        - title: Build an image
          description: Package an application and its dependencies.
        - title: Run an application
          description: Start and change a multi-service development stack.
          link: /get-started/run-an-app/
          meta: 10 min
    - persona: Building with AI
      label: Agent journey
      title: Work with AI agents safely
      description: Give an agent an isolated place to change and test code without handing over your machine.
      icon: shield-check
      length: 1 tutorial
      guides:
        - title: Run an agent in a sandbox
          description: Let an agent work in a disposable microVM and choose which changes to keep.
          link: /get-started/run-an-agent/
          meta: 10 min
    - persona: Managing Docker
      label: Admin journey
      title: Set up Docker for your organization
      description: Prepare a controlled rollout for developers, then define how AI agents can access company resources.
      icon: building-office-2
      length: 2 tutorials
      guides:
        - title: Plan an organization rollout
          description: Prepare settings, sign-in, and a pilot deployment.
          link: /guides/admin-set-up/
          meta: 20 min
        - title: Govern AI agent access
          description: Apply network, filesystem, and MCP policies across an organization.
          link: /ai/sandboxes/governance/access-controls/organization/
aliases:
  - /engine/get-started/
  - /engine/tutorials/usingdocker/
  - /guides/getting-started/get-docker-desktop/
---
