---
title: Kit capability support
linkTitle: Capability support
description: Check which v3 kit capabilities Docker Sandboxes supports locally and in the cloud, their minimum sbx versions, and restrictions that affect kit authors.
keywords: sandboxes, sbx, kits, capabilities, compatibility, cloud
weight: 37
---

{{< summary-bar feature_name="Docker Sandboxes sbx" >}}

Use this reference to check whether a v3 kit can run on your version of Docker
Sandboxes and whether it needs a local sandbox. Follow each capability link
for its settings and examples in the Docker Sandbox Kit Specification.

## Support by capability

Version numbers identify the first `sbx` release that supports the
capability, subject to the restrictions in the table. Limited support means
that Docker Sandboxes applies only part of the capability. Declaration only
means that Docker Sandboxes accepts the entry without using its commands.
The table covers behavior through v0.48.0.

All capability names in the table have the `com.docker.sandbox/` prefix.

| Capability | Local sandboxes | Cloud sandboxes | Restrictions |
| --- | --- | --- | --- |
| [`agent-context@1`](https://github.com/docker/sandbox-kit-spec/blob/main/docs/spec/capabilities/com.docker.sandbox/agent-context@1.md) | v0.45.0 | v0.45.0 | Agent instructions are applied when creating the sandbox. |
| [`agent-interactive-sessions@1`](https://github.com/docker/sandbox-kit-spec/blob/main/docs/spec/capabilities/com.docker.sandbox/agent-interactive-sessions@1.md) | Unsupported | Unsupported | Docker Sandboxes does not use these session commands. |
| [`agent-sessions@1`](https://github.com/docker/sandbox-kit-spec/blob/main/docs/spec/capabilities/com.docker.sandbox/agent-sessions@1.md) | Declaration only | Declaration only | Accepted since v0.45.0, but Docker Sandboxes does not execute the declared session commands. |
| [`agent-skill@1`](https://github.com/docker/sandbox-kit-spec/blob/main/docs/spec/capabilities/com.docker.sandbox/agent-skill@1.md) | v0.48.0 | Unsupported | Requires a selected `agent-skills@1` destination. Bundled skills remain available with `--skills=off`. |
| [`agent-skills@1`](https://github.com/docker/sandbox-kit-spec/blob/main/docs/spec/capabilities/com.docker.sandbox/agent-skills@1.md) | v0.45.0 | Unsupported | Local sandboxes mount the shared skills store. From v0.48.0, discovery paths also receive bundled skills. |
| [`credential@1`](https://github.com/docker/sandbox-kit-spec/blob/main/docs/spec/capabilities/com.docker.sandbox/credential@1.md) | v0.45.0 | Limited: v0.45.0 | In cloud sandboxes, credentials scoped to installation remain available after setup. |
| [`git-identity@1`](https://github.com/docker/sandbox-kit-spec/blob/main/docs/spec/capabilities/com.docker.sandbox/git-identity@1.md) | v0.48.0 | Unsupported | Requires `user.name` and `user.email` in the host Git config. |
| [`host-mount@1`](https://github.com/docker/sandbox-kit-spec/blob/main/docs/spec/capabilities/com.docker.sandbox/host-mount@1.md) | Unsupported | Unsupported | Kit requests for shared host directories are not supported. |
| [`kit-registry@1`](https://github.com/docker/sandbox-kit-spec/blob/main/docs/spec/capabilities/com.docker.sandbox/kit-registry@1.md) | v0.45.0 | Unsupported | Restricted to approved OCI builder kits. Local and Git kit sources do not receive registry access. |
| [`lifecycle@1`](https://github.com/docker/sandbox-kit-spec/blob/main/docs/spec/capabilities/com.docker.sandbox/lifecycle@1.md) | v0.45.0 | v0.45.0 | Local startup hooks run alongside the agent, rather than blocking its launch. |
| [`long-running@1`](https://github.com/docker/sandbox-kit-spec/blob/main/docs/spec/capabilities/com.docker.sandbox/long-running@1.md) | v0.47.0 | Unsupported | Keeps a local sandbox running after its last session disconnects. |
| [`network-policy@1`](https://github.com/docker/sandbox-kit-spec/blob/main/docs/spec/capabilities/com.docker.sandbox/network-policy@1.md) | v0.45.0 | Limited: v0.45.0 | Cloud sandboxes keep installation-phase network access after setup. |
| [`network-policy@2`](https://github.com/docker/sandbox-kit-spec/blob/main/docs/spec/capabilities/com.docker.sandbox/network-policy@2.md) | v0.45.0 | Limited: v0.45.0 | Cloud sandboxes enforce host rules, but not HTTP method or path restrictions. Installation-phase access remains after setup. |
| [`port@1`](https://github.com/docker/sandbox-kit-spec/blob/main/docs/spec/capabilities/com.docker.sandbox/port@1.md) | v0.45.0 | Limited: v0.45.0 | Cloud sandboxes expose TCP ports only. |
| [`privileged@1`](https://github.com/docker/sandbox-kit-spec/blob/main/docs/spec/capabilities/com.docker.sandbox/privileged@1.md) | v0.48.0 | Unsupported | Accepted by the local runtime, which runs Docker-in-Docker containers in privileged mode inside the sandbox VM. |
| [`resources@1`](https://github.com/docker/sandbox-kit-spec/blob/main/docs/spec/capabilities/com.docker.sandbox/resources@1.md) | Limited: v0.45.0 | Limited: v0.45.0 | Workload CPU and memory defaults are applied. The `gpu` field does not select GPUs; cloud creation rejects a required GPU request. |
| [`sbx@1`](https://github.com/docker/sandbox-kit-spec/blob/main/docs/spec/capabilities/com.docker.sandbox/sbx@1.md) | v0.45.0 | v0.45.0 | Docker Sandboxes launches the agent using the workload command. |
| [`ssh-agent@1`](https://github.com/docker/sandbox-kit-spec/blob/main/docs/spec/capabilities/com.docker.sandbox/ssh-agent@1.md) | v0.48.0 | Unsupported | Requires SSH-agent forwarding to be enabled and a host SSH agent to be available. |
| [`usb-device@1`](https://github.com/docker/sandbox-kit-spec/blob/main/docs/spec/capabilities/com.docker.sandbox/usb-device@1.md) | Unsupported | Unsupported | Kit requests for USB device access are not supported. |
| [`volume@1`](https://github.com/docker/sandbox-kit-spec/blob/main/docs/spec/capabilities/com.docker.sandbox/volume@1.md) | v0.45.0 | Limited: v0.45.0 | Local sandboxes provide block storage or tmpfs. Cloud sandboxes use ordinary sandbox disk at the declared paths, including for tmpfs requests. |

## Capability groups

Docker Sandboxes supports capability groups in `sbx` v0.48.0 and later.
A group is selected only when the runtime can provide every capability in it.
If a capability is unavailable, an optional group is skipped in full, including
its network rules, setup commands, generated files, and agent instructions.
An unavailable capability in a required group prevents sandbox creation.

Availability also depends on the host: Git identity requires configured
`user.name` and `user.email`, and SSH-agent access requires an available
agent with forwarding enabled. See the upstream
[capability group definition](https://github.com/docker/sandbox-kit-spec/blob/main/docs/spec/SPEC-v3.md#711-capability-groups)
for syntax and composition rules.

## Capabilities that require sandbox creation

Include kits that need the following settings when you create the sandbox.
`sbx kit add` cannot apply these settings to an existing sandbox:

| Capability | Settings that require creation |
| --- | --- |
| `credential@1` | Credential routing and authentication |
| `kit-registry@1` | Access to the kit registry |
| `long-running@1` | Keeping the sandbox running after sessions disconnect |
| `port@1` | Published ports |
| `resources@1` | CPU and memory allocation |
| `volume@1` | Block storage and tmpfs mounts |
| `lifecycle@1` | Generated `files` require sandbox creation. Install and startup hooks can be added with `sbx kit add`. |
| `network-policy@1` and `network-policy@2` | Host-level deny rules require sandbox creation. Allow rules can be added with `sbx kit add`, as can HTTP method and path rules for `network-policy@2`. |

`sbx kit add` rejects kits that declare credentials, published ports,
resource settings, volumes, generated lifecycle files, or host-level deny
rules. To use those kits, remove and recreate the sandbox with the kits
included. Save any work you need from the sandbox before removing it.

Skill discovery paths (`agent-skills@1`) and bundled skills (`agent-skill@1`)
can be added with `sbx kit add`: the command replaces the sandbox container
and applies the updated skill configuration.
