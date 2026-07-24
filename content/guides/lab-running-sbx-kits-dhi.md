---
title: "Lab: AI Agents in Docker Sandboxes with Kits and Hardened Images"
linkTitle: "Lab: Sandboxes, Kits, and DHI"
description: |
  Run AI coding agents inside isolated Docker Sandboxes and progressively
  harden what they produce using sbx kits and Docker Hardened Images in this
  hands-on interactive lab.
summary: |
  Hands-on lab: Run AI coding agents in isolated Docker Sandboxes, then use `sbx`
  kits and Docker Hardened Images to turn their output into secure,
  production-ready container images.
keywords: AI, Docker, Docker Sandboxes, `sbx`, kits, Docker Hardened Images, DHI, Docker Scout, container security, lab, labspace
params:
  tags: [ai, labs]
  time: 30 minutes
  resource_links:
    - title: Docker Sandboxes documentation
      url: https://docs.docker.com/ai/sandboxes/
    - title: Docker Hardened Images documentation
      url: https://docs.docker.com/dhi/
    - title: Docker Scout
      url: https://docs.docker.com/scout/
    - title: Labspace repository
      url: https://github.com/dockersamples/labspace-demo-sbx-kits-dhi
---

This lab shows you how to run AI coding agents inside isolated Docker Sandboxes
(`sbx`) and progressively harden what the agent produces using sbx **kits** and
**Docker Hardened Images (DHI)**. You'll start with a plain sandbox, attach a
container best-practices kit to change how the agent writes Dockerfiles, then add
a DHI kit so the agent builds and runs on hardened base images. Along the way
you'll compare baseline and hardened images on size, packages, vulnerabilities,
and attestations.

## Launch the lab

{{< labspace-launch image="dockersamples/labspace-sbx-kits-dhi" >}}

## What you'll learn

- Run an AI coding agent (Claude) in an isolated Docker Sandbox microVM with its own daemon, filesystem, and network
- Apply sandbox network policy that allows approved development endpoints and denies everything else
- Define an `sbx` kit: declarative, shareable agent configuration (tools, credentials, network rules, files, startup commands, and guidance) in a single `spec.yaml`
- Compare the output of a plain sandbox against one guided by a container best-practices kit
- Use a DHI kit to direct the agent to build and run on Docker Hardened Images
- Keep registry credentials on the host with sbx custom secrets, so your Docker PAT never enters the sandbox VM
- Push baseline vs. DHI image tags and compare size, package count, vulnerabilities, and attestations (SBOM + provenance) in Docker Hub and Docker Scout

## Modules

| # | Module | Description |
|---|--------|-------------|
| 0 | Prerequisites | Set up Docker Desktop, the `sbx` CLI, and a Docker Personal Access Token for pulling Docker Hardened Images |
| 1 | Start with a Plain Sandbox | Run an AI coding agent in an isolated Docker Sandbox microVM and review its default container output |
| 2 | Add the Best Practices Kit | Attach a container best-practices kit and compare how the agent's Dockerfile changes |
| 3 | Add the DHI Kit | Direct the agent to use Docker Hardened Images, then compare image size, packages, vulnerabilities, and attestations |
