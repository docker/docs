---
title: Compare local and cloud sandboxes
linkTitle: Local and cloud
description: Compare local and cloud Docker Sandboxes, including workspaces, host integrations, ports, secrets, storage, and lifecycle behavior.
keywords: docker sandboxes, local sandbox, cloud sandbox, sbx cloud, comparison
weight: 10
---

Local and cloud sandboxes provide isolated environments for AI agents, but
they run against different resources and stores. This comparison helps you
choose an environment and identify workflows that need cloud-specific setup.

| Capability | Local sandbox | Cloud sandbox |
| --- | --- | --- |
| Compute | Uses resources from the host | Uses Docker-managed cloud resources |
| Workspace | Mounts host paths or uses a private Git clone backed by the host repository | Has no access to host paths; transfer files or clone a repository inside the sandbox |
| Hardware | Can use supported host integrations, such as GPU, USB, display, and nested virtualization | Has no access to host hardware |
| Ports | Binds sandbox ports to host addresses and ports | Exposes a sandbox port through a public HTTPS URL |
| Secrets | Reads from the local `sbx` secret store | Reads from the separate `sbx --cloud` secret store |
| Network policy | Uses local and organization policy sources supported by the local runtime | Uses separate, network-only account and sandbox policy without centralized organization governance |
| Storage | Persists in the local sandbox and its attached host resources | Persists in the cloud sandbox and optional cloud volumes |
| Lifetime | Remains until you stop or remove it | Expires according to its time-to-live and timeout action |
| Billing | No metered sandbox compute charge | Metered through the Docker Agentic Platform plan |

## Host-dependent features

A cloud sandbox has no path back to the machine where you run `sbx`. The
following local features don't apply in cloud mode:

- Workspace paths, bind mounts, `--clone`, and Git worktrees created with
  `--branch`
- GPU, USB, display, and nested virtualization options
- Models served by a local model runtime
- Host port bindings
- Declarative `.sbxenv.yaml` workflows
- Host-backed agent skills and the full local MCP management workflow

The CLI rejects local-only flags used with `--cloud` instead of ignoring them.

## Separate resources

Adding `--cloud` changes the backend for the command. A sandbox shown by
`sbx ls` doesn't appear in `sbx --cloud ls`, and resources created for one
backend don't automatically become available to the other.

This separation applies to sandboxes, templates, secrets, volumes, and network
policy. Use [`sbx move`](move.md) when you need to copy a sandbox filesystem
between backends. Moving doesn't unify the resource stores or transfer
host-mounted files and secrets.
