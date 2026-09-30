---
title: Docker Agentic Platform FAQ
linkTitle: FAQ
description: Find answers about Docker Agentic Platform sandboxes, access, credentials, billing, and data retention.
keywords: docker agentic platform, faq, sandboxes, agents, docker sandboxes, usage, billing, inference, logs, telemetry, snapshots, data retention
weight: 80
aliases:
  - /agentic-platform/concepts/local-cloud/
---

## What can I run in Docker Agentic Platform?

You can run agents such as Claude Code, Codex, and Hermes in isolated cloud
sandboxes with live terminals. Browse [Kits](kits.md) for curated and community
agents, or choose Shell to work without an agent. You can also enter a public
kit reference in the launcher.

## Can I use private kits?

You can run private Docker Hub kits using the Docker Sandboxes CLI. Launching
a kit from the Console requires both the kit and its base image to be public.
The **Kits** page provides the CLI command to copy.

To access private GitHub repositories from inside a sandbox, use a
[GitHub credential](secrets.md#github-credential).

## How does Docker Agentic Platform differ from Docker Sandboxes?

Docker Agentic Platform runs sandboxes on Docker-managed cloud infrastructure
through a web Console. The `sbx` CLI runs local sandboxes on your development
machine and cloud sandboxes with `sbx --cloud`. The Console and CLI have
different workflows and secret names. See
[Cloud sandboxes](/manuals/ai/sandboxes/cloud/_index.md) for the CLI experience.

## Can I move a sandbox between my machine and Docker Agentic Platform?

The `sbx move` command copies a sandbox filesystem between local and cloud
environments. It does not transfer running processes, host bind mounts, or
managed secrets, and it leaves the source sandbox in place. See
[Move a sandbox](/manuals/ai/sandboxes/cloud/move.md) for the CLI workflow and
its limitations.

## Can I share sandboxes and configuration with a team?

The initial release is for individual use. You manage your own sandboxes, MCP
connections, secrets, and network policies. Shared workspaces and team
ownership aren't supported.

## How does a sandbox access external services?

The **Open** user policy allows outbound access to any host. **Balanced**
allows a curated set of hosts and services. Allow rules from all applicable
policies are combined: selecting both **Open** and **Balanced** permits all
outbound destinations except those blocked by explicit deny rules. Balanced's
allow list doesn't restrict Open's access. To restrict access with an allow
list, deselect **Open** if it is selected and can be removed.

The launcher remembers your policy selections from the previous launch in the
same browser. Without saved selections, it selects the account's policies
marked **Always applied**. These policies can't be deselected in the launcher.

Your kit's network rules also apply. If an allow rule and a deny rule match
the same destination, the deny rule wins. See [Network policies](policies.md)
for details.

For services that need authentication, save your credentials under
[Secrets](secrets.md). The sandbox proxy adds them to matching requests without
exposing their values to the agent.

## How are sandbox usage and model inference billed?

You pay for compute by the second while your sandbox runs. The Console also
shows the equivalent hourly rate.

Your model provider bills inference separately, under the account associated
with your API key. See [Docker Billing](/subscription-billing/) for account,
usage, and payment information.

## How long are logs, telemetry, and snapshots retained?

Service logs are kept for 31 days and raw telemetry for 12 months.

All snapshots, including the most recent snapshot created when you pause a
sandbox, are automatically deleted after seven days of non-use.
