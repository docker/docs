---
title: Tools
description: Interfaces and tools for browsing, managing, and automating Docker Hardened Images.
weight: 25
params:
  grid_tools:
    - title: Use Docker Hub
      description: Browse the DHI catalog on Docker Hub to search repositories, inspect image metadata, and view SBOMs, CVEs, and attestations.
      icon: squares-2x2
      link: /dhi/tools/hub/
    - title: CLI
      description: Install and use the `docker dhi` command-line interface to browse the catalog, inspect images, and manage mirrors from your terminal.
      icon: command-line
      link: /dhi/tools/cli/
    - title: MCP server
      description: Connect an AI assistant to the DHI catalog to search repositories, inspect images, retrieve SBOMs, and check CVEs using plain language.
      icon: cpu-chip
      link: /dhi/tools/mcp/
    - title: Use the DHI Terraform provider
      description: Use the DHI Terraform provider to manage mirrors and automate DHI configuration as infrastructure as code.
      icon: wrench-screwdriver
      link: /dhi/tools/terraform/
---

Docker Hardened Images can be accessed and managed through several interfaces.
Choose the tool that fits your workflow.

{{< grid items="grid_tools" >}}
