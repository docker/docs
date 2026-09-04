---
title: Workflow patterns
linkTitle: Workflows
weight: 50
description: Choose a Docker Sandboxes workflow for agent skills, Git, local development, authenticated tools, or headless automation.
keywords: docker sandboxes, sbx, workflows, agent skills, shared skills, clone mode, git, local development, authentication, ci, headless
toc_max: 2
---

Choose a workflow based on how you want to develop, authenticate tools, or run
sandboxes in automation. For command syntax and lifecycle basics, see
[Usage](../usage.md).

## Choose how code moves

Your workspace strategy determines when an agent's changes appear on the host.
Direct mode edits the host working tree in place. Clone mode keeps changes in a
private clone until you fetch or push them. Host worktrees provide branch
isolation while keeping Git operations on the host. See [Git workflows](git.md)
to choose a strategy for single or parallel tasks.

## Develop in the sandbox

Each sandbox has a private Docker daemon and runtime for building images,
installing dependencies, and running tests. You can publish services from the
sandbox or connect to services on the host. See
[Develop and test locally](development.md).

Tools inside the sandbox can use credentials configured on the host without
copying secret values into the VM. See
[Authenticate command-line tools](authentication.md) for GitHub CLI, registry,
and external secret-provider workflows.

## Reuse and automate workflows

Sandbox environment files work like Compose files for sandboxes: they capture
project configuration in a versioned YAML file. Use `sbxenv.yaml` to define
the agent, workspaces, tools, resources, credentials, and ports so contributors
can start a consistent environment without reproducing CLI flags and setup
steps. See [Sandbox environment files](../configuration/environment-files.md).

You can also import skills from supported host agents into a persistent store
shared with new sandboxes. See [Share agent skills](agent-skills.md).

For unattended jobs, use headless authentication and manage the sandbox
lifecycle from scripts. See [Run sandboxes in CI](automation.md).
