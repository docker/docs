---
title: Use environment files in the cloud
linkTitle: Environment files
description: Create a repeatable cloud sandbox with sbx env, copy in your project, work with an agent, and retrieve your changes.
keywords: docker sandboxes, cloud sandbox, sbx env, sbxenv, environment file
weight: 25
---

Use an environment file to share a cloud sandbox's agent, tools, environment
variables, and resource limits with your project contributors. Run
`sbx --cloud env` to apply that configuration on Docker-managed infrastructure.

A cloud environment has its own filesystem. It doesn't mount or synchronize
directories from your machine, so leave `workspace` out of the file. Copy your
project into the sandbox or clone it there, and copy the results back or push
your commits before removing the sandbox. Credentials and network policy are
also separate from your local sandbox configuration. Host lifecycle commands
still run on your machine, with your privileges.

> [!NOTE]
> `sbx env` is experimental. The command interface and file format may change.

## Prepare your environment

Complete the [cloud prerequisites](_index.md#prerequisites), including activating
cloud access and signing in with `sbx login`. Cloud compute charges apply.

This walkthrough uses Claude Code to work on a local project named `web-app`.
Configure an Anthropic API key in the cloud secret store before creating the
sandbox:

```console
$ sbx --cloud secret set anthropic
```

The command prompts for the key. For other agents and authentication options,
see [Authenticate cloud agents](credentials.md). A credential configured for a
local sandbox isn't available in the cloud.

Keep the environment file beside your project, outside the directory you'll
copy into the sandbox:

```text
web-app-env/
├── cloud.sbxenv.yaml
└── web-app/
```

Save the following as `web-app-env/cloud.sbxenv.yaml`:

```yaml
schemaVersion: "1"
name: cloud-web-app
agent: claude

env:
  NODE_ENV: test

sandboxOptions:
  cpus: 2
  memory: 4g
```

This selects Claude Code, sets `NODE_ENV` for the sandbox, and requests two
CPUs and 4 GiB of memory. Resource limits must match a
[cloud size](usage.md#choose-resources-and-platform). You can also add
[`kits`](../configuration/environment-files.md#kits) to install project tools.

## Create the sandbox and add your project

From `web-app-env`, review the plan and create the sandbox without attaching:

```console
$ sbx --cloud env plan ./cloud.sbxenv.yaml
$ sbx --cloud env create --ttl 2h --on-timeout delete ./cloud.sbxenv.yaml
```

The plan shows the configuration and inherited cloud credentials. Creation
asks you to approve the plan before applying it. This example explicitly
deletes the sandbox after two hours, so retrieve your work before then.
Without timeout flags, a cloud sandbox expires after one hour, and the service
stops sandboxes that can be resumed and deletes the rest. See
[Configure expiration](usage.md#configure-expiration) for other options.

Copy your project into the sandbox's workspace:

```console
$ sbx --cloud cp ./web-app cloud-web-app:/home/agent/workspace/web-app
```

Use a destination that doesn't already exist for this first copy. Copying a
directory to an existing directory nests it under its source name. Review the
files you're uploading, including hidden files, and exclude credentials and
other files you don't want in the cloud. Copying is a one-time transfer:
later edits on either side aren't synchronized. Directory uploads skip symbolic
links and empty directories and report what was skipped. Install dependencies
inside the sandbox rather than copying a local dependency directory.

For a Git workflow, you can instead clone a repository inside the sandbox and
push commits to your remote. Configure [cloud credentials](credentials.md) for
private repository access and [cloud network policy](network-policy.md) for
the destinations your project needs.

## Work with the agent

Attach to Claude Code using the same environment file:

```console
$ sbx --cloud env run ./cloud.sbxenv.yaml
```

For example, ask the agent:

```text
Work in /home/agent/workspace/web-app. Find the test and lint commands in
package.json, install the dependencies, and run those checks. Fix any failures
and explain your changes. Ask before changing dependencies.
```

Use `run` again to reconnect to the existing sandbox. It creates a sandbox if
one doesn't exist. Pass `--ttl` and `--on-timeout` only when creating a sandbox:
`run` rejects those flags when attaching to an existing one.

Press `Ctrl+\` to detach while leaving the agent running. You can also run a
command without an agent session, for example to inspect the working tree of
a copied Git project:

```console
$ sbx --cloud env exec ./cloud.sbxenv.yaml -- \
    git -C /home/agent/workspace/web-app diff --stat
```

For subsequent environment commands, keep the same machine, cloud endpoint,
Docker account, and ordered environment file paths. The CLI records the cloud
sandbox's identity locally to manage it across invocations.

## Retrieve your work and clean up

After the agent finishes, copy the project back to a separate local directory
so you can review it before updating your original files:

```console
$ sbx --cloud cp cloud-web-app:/home/agent/workspace/web-app ./web-app-result
```

Use a destination that doesn't already exist. If you cloned the project inside
the sandbox instead, push the commits you want to keep to your remote.

When you've saved your work, remove the environment:

```console
$ sbx --cloud env rm ./cloud.sbxenv.yaml
```

Removal asks for approval and deletes the sandbox and secrets provisioned by
this environment. Inherited credentials, such as the Anthropic key configured
earlier, remain. Global bindings remain unless you pass `--prune-bindings`.
Files stored only in the sandbox are lost on deletion.

For unattended commands, use `--force` with `create`, `run`, or `rm` to approve
the plan for that invocation. Use `run --detach` when no interactive terminal
is available.

## Adapt an existing environment file

Cloud environments support agents, kits, environment variables, CPU and memory
limits, credentials for supported providers, and host lifecycle commands.
They support a subset of the [environment file fields](../configuration/environment-files.md#file-reference).
Before reusing a local file, remove:

- `workspace`, `additionalWorkspaces`, and clone options
- `ports`, `registries`, and `mcp.servers`
- `sandboxOptions.template`, `profile`, `pullPolicy`, `skills`, `display`,
  `gpu`, `usb`, and `writableEnvFiles`

Unsupported settings are rejected before host commands or provisioning run.
Pass the cloud file path explicitly to avoid loading local settings from
`~/.sbxenv.yaml`. For services and integrations, use the cloud commands to
[expose ports](usage.md#expose-a-port) and [load MCP servers](usage.md#load-an-mcp-server).

Lifecycle commands run on your machine even for a cloud environment. A command
that installs dependencies locally won't install them in the cloud sandbox.
Use a kit or run the installation inside the sandbox instead.

The plan shows inherited cloud credentials. Sandbox-scoped credentials override
account defaults, and credentials declared in `secrets` override both. Use
[`snapshot: true`](../configuration/environment-files.md#secrets) to resolve a
host command or vault reference once and upload the value. Snapshots don't
refresh. Dynamic secret sources and custom credential providers aren't
supported. Credential bindings require cloud support for kit credentials and
explicit approval of the kit's injection domains.

Changes to secrets, bindings, kits, or resource limits require recreating the
sandbox. Updated `env` values apply to subsequent sessions. Rejoining a running
agent keeps that process's environment.

If creation is interrupted, retry the same command and unchanged declaration
within 23 hours. Unresolved requests prevent removal. Follow the recovery
message and retain its journal until the original requests are resolved.
