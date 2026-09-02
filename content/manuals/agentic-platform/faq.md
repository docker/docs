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

Docker Agentic Platform provides predefined sandbox types for Claude Code,
Codex, OpenCode, Copilot, and Gemini CLI. Each runs in an isolated, Docker-hosted
sandbox with a live terminal.

## How does Docker Agentic Platform differ from Docker Sandboxes?

Docker Agentic Platform runs sandboxes on Docker-managed cloud infrastructure
through a web Console. Docker Sandboxes runs sandboxes on your development
machine through the `sbx` command. Docker Agentic Platform manages the compute,
MCP connections, secrets, and network policies used by its hosted sandboxes.

## Can I move a sandbox between my machine and Docker Agentic Platform?

No. Local and hosted sandboxes are separate in the initial release. You cannot
move a running sandbox or its local bind mounts into Docker Agentic Platform.

## Can I share sandboxes and configuration with a team?

The initial self-service experience is single-user. You manage your own
sandboxes, MCP connections, secrets, and network policies. Shared workspaces
and collaborative ownership are not part of the initial release.

## How does a sandbox access external services?

By default, every new sandbox uses the **Open** user policy, regardless of
sandbox type. **Open** allows access to all outbound destinations. To restrict
egress, replace **Open** with **Balanced**, a custom policy, or no user policy.

Network policies control the destinations a sandbox can reach. The sandbox
type's read-only kit policy applies automatically, and you can select zero or
more user policies when you create the sandbox. If you select no user policies,
only the kit policy applies and destinations that it does not allow are blocked.
A deny rule takes precedence over an allow rule. Docker stores configured
secret values outside the sandbox and applies them to matching requests through
the sandbox proxy.

## How are sandbox usage and model inference billed?

Docker bills sandbox compute per second while the sandbox runs. The Console
also shows the equivalent hourly rate.

Model inference is billed separately. The sandbox uses your credential for an
external model provider, which meters and bills inference under that provider
account. See [Docker Billing](/billing/) for account, usage, and payment
information.

## How long are logs, telemetry, and snapshots retained?

Docker retains Docker Agentic Platform service logs for 31 days and raw
telemetry for 12 months.

All snapshots, including the most recent snapshot created when you pause a
sandbox, are automatically deleted after seven days of non-use.
