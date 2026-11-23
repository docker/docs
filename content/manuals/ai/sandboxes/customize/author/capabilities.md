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

Capability descriptions link to the versioned Docker Sandbox Kit
Specification. The table identifies Docker Sandboxes runtime support for
each capability separately.

Version numbers identify the first `sbx` release that supports the
capability, subject to the restrictions in the table. Limited support means
that Docker Sandboxes applies only part of the capability. Declaration only
means that Docker Sandboxes accepts the entry without using its commands.

All capability names in the table have the `com.docker.sandbox/` prefix.
Not verified means Docker Sandboxes support has not been confirmed for that
capability. Specification links are pinned to the catalog release.

| Capability | Local sandboxes | Cloud sandboxes | Restrictions |
| --- | --- | --- | --- |
| [`acp@1`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/acp@1.md) | Not verified | Not verified | Docker Sandboxes support has not been verified. |
| [`agent-context@1`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/agent-context@1.md) | v0.45.0 | v0.45.0 | Agent instructions are applied when creating the sandbox. |
| [`agent-context@2`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/agent-context@2.md) | v0.51.0 | Limited: v0.51.0 | Cloud sandboxes write one workload profile in `/home/agent`. Mixins use the workload's profile and cannot specify a `directory`. Unsupported required profiles prevent creation; optional ones are skipped. |
| [`agent-interactive-sessions@1`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/agent-interactive-sessions@1.md) | Unsupported | Unsupported | Docker Sandboxes does not use these session commands. |
| [`agent-sessions@1`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/agent-sessions@1.md) | Declaration only | Declaration only | Accepted since v0.45.0, but Docker Sandboxes does not execute the declared session commands. |
| [`agent-skill@1`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/agent-skill@1.md) | v0.48.0 | Unsupported | Requires a selected `agent-skills@1` destination. Bundled skills remain available with `--skills=off`. |
| [`agent-skills@1`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/agent-skills@1.md) | v0.45.0 | Unsupported | Local sandboxes mount the shared skills store. From v0.48.0, discovery paths also receive bundled skills. |
| [`auto-restart@1`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/auto-restart@1.md) | Not verified | Not verified | Docker Sandboxes support has not been verified. |
| [`credential@1`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/credential@1.md) | v0.45.0 | Limited: v0.45.0 | In cloud sandboxes, credentials scoped to installation remain available after setup. |
| [`docker-auth@1`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/docker-auth@1.md) | Not verified | Not verified | Docker Sandboxes support has not been verified. |
| [`dockerd@1`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/dockerd@1.md) | Not verified | Not verified | Docker Sandboxes support has not been verified. |
| [`git-identity@1`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/git-identity@1.md) | v0.48.0 | Unsupported | Requires `user.name` and `user.email` in the host Git config. |
| [`host-mount@1`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/host-mount@1.md) | v0.50.0 | Unsupported | Users can [select host directories](../use-kits.md#choose-directories-for-kit-host-mounts) at creation time. Cloud sandboxes reject required entries and skip optional ones. |
| [`host-port@1`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/host-port@1.md) | Not verified | Not verified | Docker Sandboxes support has not been verified. |
| [`kit-registry@1`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/kit-registry@1.md) | v0.45.0 | Unsupported | Restricted to approved OCI builder kits. Local and Git kit sources do not receive registry access. |
| [`lifecycle@1`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/lifecycle@1.md) | v0.45.0 | v0.45.0 | Local startup hooks run alongside the agent, rather than blocking its launch. |
| [`long-running@1`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/long-running@1.md) | v0.47.0 | Unsupported | Keeps a local sandbox running after its last session disconnects. |
| [`network-policy@1`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/network-policy@1.md) | v0.45.0 | Limited: v0.45.0 | Cloud sandboxes keep installation-phase network access after setup. |
| [`network-policy@2`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/network-policy@2.md) | v0.45.0 | Limited: v0.45.0 | Cloud sandboxes enforce host rules, but not HTTP method or path restrictions. Installation-phase access remains after setup. |
| [`open-url@1`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/open-url@1.md) | Not verified | Not verified | Docker Sandboxes support has not been verified. |
| [`port@1`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/port@1.md) | v0.45.0 | Limited: v0.45.0 | Cloud sandboxes expose IPv4 TCP ports only. See [Port address families](#port-address-families). |
| [`privileged@1`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/privileged@1.md) | v0.48.0 | Unsupported | Accepted by the local runtime, which runs Docker-in-Docker containers in privileged mode inside the sandbox VM. |
| [`resources@1`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/resources@1.md) | Limited: v0.45.0 | Limited: v0.45.0 | Workload CPU and memory defaults are applied. The `gpu` field does not select GPUs; cloud creation rejects a required GPU request. |
| [`sbx@1`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/sbx@1.md) | v0.45.0 | v0.45.0 | Docker Sandboxes launches the agent using the workload command. |
| [`ssh-agent@1`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/ssh-agent@1.md) | v0.48.0 | Unsupported | Requires SSH-agent forwarding to be enabled and a host SSH agent to be available. |
| [`usb-device@1`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/usb-device@1.md) | Unsupported | Unsupported | Kit requests for USB device access are not supported. |
| [`volume@1`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/volume@1.md) | v0.45.0 | Limited: v0.45.0 | Local sandboxes provide block storage or tmpfs and [retain block volumes across removal and recreation](../use-kits.md#keep-kit-data-when-recreating-a-sandbox). Cloud sandboxes use ordinary sandbox disk at the declared paths, including for tmpfs requests. |

<!-- capability-catalog:descriptions:start -->

## Capability descriptions

### [`acp@1`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/acp@1.md)

How to run an [Agent Client Protocol](https://agentclientprotocol.com/) (ACP) agent inside the sandbox. This capability describes the launch command, its environment, and the protocol major version.

### [`agent-context@1`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/agent-context@1.md)

Written instruction content the agent reads — the AGENTS.md family. A host with no such concept skips an optional declaration or refuses a required one.

### [`agent-context@2`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/agent-context@2.md)

Instruction content stored inline in the descriptor, with configurable runtime delivery. Mixin guidance defaults to progressive surfacing: each discovery profile indexes a per-Kit file that the agent reads on demand. An explicit `delivery: inline` places that contribution directly in every profile. Workload guidance always appears directly in the profiles.

### [`agent-interactive-sessions@1`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/agent-interactive-sessions@1.md)

The workload's interactive session control surface: the argv shapes a human-facing host uses to put the agent's terminal UI in front of a person — start a session, seed it with a prompt, reopen or pick a past one. Not a grant: the host consumes it to operate the agent, the way it consumes the image config's entrypoint.

### [`agent-sessions@1`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/agent-sessions@1.md)

How to run the workload's agent headlessly inside the sandbox. This capability describes commands to run a prompt, resume or continue a session, and list past sessions.

### [`agent-skill@1`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/agent-skill@1.md)

One bundled agent skill, supplied by the Kit's image and exposed in the composed agent's skill directories. A workload or a mixin can supply it.

### [`agent-skills@1`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/agent-skills@1.md)

A directory this Kit's agent scans for skills. It receives selected Kit-bundled skills and, when available and enabled, the host's shared skills store. Declare it on the Kit that supplies the agent, whether it is a workload or a mixin.

### [`auto-restart@1`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/auto-restart@1.md)

Resume a previously running sandbox after the host or runtime recovers. Service Kits can request this behavior so that other sandboxes can reach them again after a host reboot or a runtime upgrade.

### [`credential@1`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/credential@1.md)

One service the workload authenticates to, and how the runtime presents proof — never where the secret lives. The host's credential store is the sole source; a Kit declares the need, the user's bindings answer it.

### [`docker-auth@1`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/docker-auth@1.md)

Registry authentication through opaque, run-scoped placeholders. Real login credentials and issued tokens remain outside the sandbox.

### [`dockerd@1`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/dockerd@1.md)

A request for a Docker Engine inside the sandbox: elevation, a private persistent store for the engine, and by default the engine started for the Kit. The host may refuse.

### [`git-identity@1`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/git-identity@1.md)

The user's runtime-provided Git author identity, supplied as `user.name` and `user.email` defaults inside the sandbox. This requests attribution, not authentication, signing authority or evidence of user review. Authentication uses separate grants such as `credential@1`.

### [`host-mount@1`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/host-mount@1.md)

One host directory shared with the host and across sandboxes composing the same Kit. The runtime chooses its location and mount mechanism; Linux bind mounts and VM filesystem sharing can satisfy the same grant. Use [volume@1](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/volume@1.md) for private storage keyed on owner scope, sandbox name, and path. It survives sandbox deletion for reattachment under that name, without this host or cross-sandbox sharing contract.

### [`host-port@1`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/host-port@1.md)

One host loopback TCP or UDP service made available at an IPv4 loopback endpoint inside the sandbox. This reverses [port@1](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/port@1.md), which publishes a sandbox service to the host. It grants access to a specific host service, not general network egress.

### [`kit-registry@1`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/kit-registry@1.md)

Access to the runtime's Kit registry — the runtime-hosted registry endpoint over its own image store, where Kit builds push results and pull the Kit frontend. Declaring this capability is the only thing that makes the registry reachable: the typed request literally is the enforcement.

### [`lifecycle@1`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/lifecycle@1.md)

The Kit's setup and launch behavior, executed by the engine across the sandbox's life: install hooks once at create, startup hooks every boot, files written at start, and the interactive argv tail for TTY sessions. A capability rather than grammar fields so a host that cannot execute them refuses the Kit by name at preflight instead of silently never running its setup.

### [`long-running@1`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/long-running@1.md)

The workload runs independently of client sessions. A background startup hook or a published port alone does not request this behavior: on a runtime with session-based auto-stop, the sandbox can otherwise stop when its last client disconnects.

### [`network-policy@1`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/network-policy@1.md)

Phase-scoped egress policy: what the sandbox may reach while the Kit's install hooks run, and what the agent may reach in steady state.

### [`network-policy@2`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/network-policy@2.md)

Phase-scoped egress policy whose entries may bound a host to particular HTTP requests: what the sandbox may reach while the Kit's install hooks run, what the agent may reach in steady state, and what it may do there.

### [`open-url@1`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/open-url@1.md)

Request HTTP or HTTPS navigation in the runtime host's default browser. This is an explicit host interaction grant, not browser automation or sandbox network access. The host is the machine running the runtime, not necessarily the CLI client machine. A remote runtime does not silently substitute the client's browser.

### [`port@1`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/port@1.md)

One in-container port the runtime should publish to the host. Inbound service exposure — distinct from the outbound [network policy](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/network-policy@1.md).

### [`privileged@1`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/privileged@1.md)

A request for elevated privilege. Present or absent — there is nothing to configure, and the host may refuse.

### [`resources@1`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/resources@1.md)

The workload's compute: CPU, memory, GPU. A constraint on the Kit, not a grant to it.

### [`sbx@1`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/sbx@1.md)

A declaration that the workload targets the sandbox agent platform: the host launches the agent rather than letting the image entrypoint be PID 1, and it honors the identity the image config states instead of assuming one.

### [`ssh-agent@1`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/ssh-agent@1.md)

An SSH agent the workload can use on the user's behalf. The workload can list the agent's public keys and use them to authenticate to remote machines or sign data, while private keys stay outside the sandbox.

### [`usb-device@1`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/usb-device@1.md)

One USB passthrough request, matched by vendor/product identity or by device class. Hardware access is host-answerable and permission-gated like any other grant.

### [`volume@1`](https://github.com/docker/sandbox-kit-spec/blob/v3.0.0-m.13/docs/spec/capabilities/com.docker.sandbox/volume@1.md)

One sandbox-private storage destination. Block storage belongs to the sandbox name within its owner scope and outlives sandbox deletion; tmpfs is scratch storage.

<!-- capability-catalog:descriptions:end -->

## Port address families

For v3 `port@1` entries, `family` selects the address families used to publish
the port: `ipv4`, `ipv6`, or `dual` for both. If no kit specifies a family,
the port defaults to IPv4, even with an explicit `transport: tcp`. Use
`family: dual` when the service needs both IPv4 and IPv6 clients.

When kits declare the same port and transport, an entry without `family`
uses the family specified by another kit in the set. Kits that specify
conflicting families for that port cannot be used together.

Cloud sandboxes support IPv4 TCP ports only. They reject required `ipv6` and
`dual` entries and skip optional ones. For capability syntax and
examples, follow the `port@1` specification link in the table.

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
