---
title: Cloud sandboxes
description: Run Docker Sandboxes on Docker-managed cloud infrastructure and understand the cloud-specific command, storage, and billing model.
keywords: docker sandboxes, cloud sandboxes, sbx cloud, ai agents, agentic platform
weight: 15
---

Cloud sandboxes run AI agents on Docker-managed infrastructure instead of your
local machine. Use them when you need an isolated environment that doesn't
depend on the compute resources or virtualization support of your host.

Cloud sandboxes use the same `sbx` CLI as local sandboxes. Add the global
`--cloud` flag to send a supported command to the Cloud Sandboxes API:

```console
$ sbx --cloud ls
```

Cloud and local sandboxes have separate state and different capabilities. A
cloud sandbox can't mount a host workspace or use host hardware, and its
secrets, network policy, ports, and lifecycle are managed in the cloud. See
[Local and cloud differences](local-vs-cloud.md) before adapting a local
workflow.

## Prerequisites

To use cloud sandboxes, you need:

- The [`sbx` CLI](../install.md)
- A Docker account signed in through `sbx login`
- An active [Docker Agentic Platform plan](/manuals/subscription-billing/plans/docker-agentic-platform.md)

Cloud sandbox compute is metered through the Docker Agentic Platform
pay-as-you-go plan. Inference charges aren't included. Your model provider
charges for requests made with the API keys or OAuth credentials that you
configure.

## Get started

Credentials configured for local sandboxes aren't available to cloud
sandboxes. Configure a cloud credential for your agent before launching it.
For example, authenticate Claude Code with Anthropic OAuth:

```console
$ sbx --cloud secret set anthropic --oauth
```

Then create and attach to a cloud sandbox:

```console
$ sbx --cloud run claude --name cloud-project
```

The command creates the sandbox when the name doesn't exist, then opens the
agent session. Cloud sandboxes don't accept a local workspace path. Copy files
into the sandbox after creation or clone a repository from inside the sandbox.

Use `sbx --cloud` for later operations on the same sandbox:

```console
$ sbx --cloud ls
$ sbx --cloud exec cloud-project git status
$ sbx --cloud attach cloud-project
```

## Learn more

- [Local and cloud differences](local-vs-cloud.md) compares the two execution
  environments
- [Use cloud sandboxes](usage.md) covers creation, files, ports, and lifecycle
- [Authenticate cloud agents](credentials.md) covers cloud-specific secrets,
  OAuth, and safe sign-in workflows
- [Manage cloud network policy](network-policy.md) covers account-level and
  sandbox-level network access
- [Move a sandbox](move.md) explains filesystem transfers between local and
  cloud environments
- [`sbx` CLI reference](/reference/cli/sbx/) lists commands and options
