---
title: Default security posture
linkTitle: Defaults
weight: 15
description: What a sandbox permits and blocks before you change any settings.
keywords: docker sandboxes, security defaults, network policy, credentials, sbx
---

{{< summary-bar feature_name="Docker Sandboxes sbx" >}}

A sandbox created with `sbx run` and no additional flags has the following
security posture.

## Network defaults

All outbound HTTP and HTTPS traffic is blocked unless an explicit rule allows
it (deny-by-default). All non-HTTP protocols (raw TCP, UDP including DNS, and
ICMP) are blocked at the network layer. Traffic to private IP ranges, loopback
addresses, and link-local addresses is also blocked.

Run `sbx policy ls` to see the active network rules for your installation. To
customize network access, see [Policies](policy.md). If your organization
manages sandbox policies centrally, those rules apply on top of the defaults
described here. See [Organization governance](governance.md).

## Workspace defaults

Sandboxes use a direct mount by default. The agent sees and modifies your
working tree directly, and changes appear on your host immediately.

The agent can read, write, and delete any file within the workspace directory,
including hidden files, configuration files, build scripts, and Git hooks.
See [Workspace trust](workspace.md) for what to review after an agent session.

## Credential defaults

No credentials are available to the sandbox unless you provide them using
`sbx secret` or environment variables. When credentials are provided, the
host-side proxy injects them into outbound HTTP headers. The agent cannot
read the raw credential values.

See [Credentials](credentials.md) for setup instructions.

## Agent capabilities inside the sandbox

The agent runs with full control inside the sandbox VM:

- `sudo` access (the agent runs as a non-root user with sudo privileges)
- A private Docker Engine for building images and running containers
- Package installation through `apt`, `pip`, `npm`, and other package managers
- Full read and write access to the VM filesystem

Everything the agent installs or creates inside the VM, including packages,
Docker images, and configuration changes, persists across stop and restart
cycles. When you remove the sandbox with `sbx rm`, the VM and its contents
are deleted. Only workspace files remain on the host.

## What is blocked by default

The following are blocked for all sandboxes and cannot be changed through
policy configuration:

- Host filesystem access outside the workspace directory
- Host Docker daemon
- Host network and localhost
- Communication between sandboxes
- Raw TCP, UDP, and ICMP connections
- Traffic to private IP ranges and link-local addresses

Outbound HTTP/HTTPS to domains not in the allow list is also blocked by
default, but you can add allow rules with `sbx policy allow`.
