---
title: Default security posture
linkTitle: Defaults
weight: 20
description: What a sandbox permits and blocks before you change any settings.
keywords: docker sandboxes, security defaults, network policy, credentials, shared skills, sbx
---

A sandbox created with `sbx run claude` and no additional flags has the
following security posture.

## Network defaults

All outbound TCP traffic, including HTTP, HTTPS, and SSH, is blocked unless an
explicit rule allows the destination. Direct external UDP and ICMP traffic is
blocked at the network layer. DNS queries use the sandbox's internal resolver,
which enforces network policy.

Run `sbx policy ls` to see the active network rules for your installation.
Rules can be customized per machine with the `sbx policy` CLI, or managed
centrally across your organization. Org-level rules take precedence over local
rules. See
[Network access policies](../governance/access-controls/network.md).

## Workspace defaults

`sbx run` mounts the current directory when you don't pass a workspace path.
The agent can read, write, and delete any file within that directory, including
hidden files, configuration files, build scripts, and Git hooks.

When you omit the workspace path from `sbx create`, the sandbox doesn't mount a
host workspace. The agent works in the template image's configured `WORKDIR`.
The Docker-provided images for built-in agents set `WORKDIR` to
`/home/agent/workspace`. If the daemon can't resolve a usable absolute
`WORKDIR` from the image config, it falls back to that path. Files in this
directory persist across stops and restarts and are deleted when you remove the
sandbox. See
[Workspace isolation](isolation.md#workspace-isolation) for the available
workspace modes and what to review after a direct-mount session.

## Shared skills defaults

Sandboxes for supported agents mount a persistent shared skills store
read-write by default. Every sandbox that uses the store can change skills that
other participating sandboxes may load. Use `--no-share-skills` when creating a
sandbox to keep it outside this shared trust boundary. See
[Share agent skills](../workflows/agent-skills.md).

## Credential defaults

No credentials are available to the sandbox unless you provide them using
`sbx secret` or environment variables. When credentials are provided, the
host-side proxy injects them into outbound HTTP headers. The agent cannot
read the raw credential values.

See [Credentials](../configuration/credentials.md) for setup instructions.

## Agent capabilities inside the sandbox

The agent runs with full control inside the sandbox VM:

- `sudo` access (the agent runs as a non-root user with sudo privileges)
- A private Docker Engine for building images and running containers
- Package installation through `apt`, `pip`, `npm`, and other package managers
- Full read and write access to the VM filesystem

Everything the agent installs or creates inside the VM, including packages,
Docker images, mountless workspace files, and configuration changes, persists
across stop and restart cycles. When you remove the sandbox with `sbx rm`, the
VM and its contents are deleted. Direct-mounted workspace files and the shared
skills store remain on the host, as do repositories used as clone sources.

## What is blocked by default

The following are blocked for all sandboxes and cannot be changed through
policy configuration:

- Host filesystem access outside explicitly mounted workspaces and the shared
  skills store
- Host Docker daemon
- Direct network communication between sandboxes
- Direct external UDP and ICMP connections

Outbound TCP to destinations not in the allow list is also blocked by default,
but you can add allow rules with `sbx policy allow`.
