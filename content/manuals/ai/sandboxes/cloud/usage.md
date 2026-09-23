---
title: Use cloud sandboxes
description: Create and manage Docker cloud sandboxes with the sbx CLI, including file transfer, commands, ports, storage, and lifecycle controls.
keywords: docker sandboxes, cloud sandbox, sbx cloud, cloud ports, sandbox ttl
weight: 20
---

Use the `--cloud` flag with supported `sbx` commands to create and manage
sandboxes on Docker-managed infrastructure. Cloud operations use cloud IDs,
names, resources, and lifecycle controls rather than the local sandbox daemon.

## Create a sandbox

A cloud sandbox expires after one hour by default. On expiration, the service
stops sandboxes that can be resumed and deletes the rest. See
[Configure expiration](#configure-expiration) to choose the timeout and action
before creating it.

Credentials saved for local sandboxes aren't available in cloud sandboxes.
[Configure a cloud credential](credentials.md) before launching an agent.

Create a sandbox and attach to its agent, or reuse the named sandbox if it
already exists:

```console
$ sbx --cloud run claude --name cloud-project
```

Without `--name`, an interactive run offers existing sandboxes for that agent
and an option to create another. Pass `--new` to create a fresh sandbox.
Launches that bake a kit template also create a fresh sandbox.

Reusing a sandbox keeps its creation settings. Flags such as `--cpus`,
`--memory`, `--platform`, `--ttl`, `--env`, and `--allow-network` are rejected
when resuming. Use `--new` to create a sandbox with different settings.

To create the sandbox without opening an agent session, use `create`:

```console
$ sbx --cloud create --name cloud-project claude
```

The command prints the cloud sandbox ID. Attach by ID or name:

```console
$ sbx --cloud attach cloud-project
```

Cloud sandbox names must have at least two characters, start with a letter or
number, and contain only letters, numbers, and hyphens. The name `default` is
reserved. Periods accepted in local sandbox names aren't accepted in the cloud.

Cloud creation doesn't accept workspace paths. For example,
`sbx --cloud run claude .` returns an error because `.` refers to the local
filesystem.

### Choose resources and platform

Without resource flags, a cloud sandbox starts with 2 CPUs and 4 GiB of
memory. Use `--cpus` and `--memory` to select one of these configurations:

| Size | CPUs | Memory |
| --- | --- | --- |
| micro | 1 | 2 GiB |
| small | 2 | 4 GiB |
| medium | 4 | 8 GiB |
| large | 8 | 16 GiB |
| xl | 16 | 32 GiB |

For example:

```console
$ sbx --cloud create --name cloud-project --cpus 4 --memory 8g claude
```

If you specify only CPU or memory, the CLI selects the matching value for the
other resource. Unsupported combinations are rejected.

Use `--platform linux/amd64` or `--platform linux/arm64` to select an
architecture supported by your account. This matters when you plan to
[move the sandbox to your machine](move.md): the architectures must match.

## Run without attaching

For scripts or terminals without interactive input, start the sandbox without
opening an agent session:

```console
$ sbx --cloud run --detached claude --name cloud-task
```

A detached run with `--name` reuses the named sandbox and starts it if stopped.
If the sandbox doesn't exist, or you omit `--name`, it creates a sandbox.
Use `sbx --cloud exec` to run commands and `sbx --cloud rm --force` for
unattended cleanup. An interactive agent session requires `run` or `attach`
from a terminal.

To detach from an interactive `run` or `attach` session while leaving the
agent running, press `Ctrl+\`. Reconnect with `sbx --cloud attach <sandbox-name>`.
Reconnecting joins the existing agent session. Use `--detach-keys` with `run`
or `attach` to change the detach gesture, for example `--detach-keys ctrl-x,ctrl-d`.

## List and inspect sandboxes

List cloud sandboxes separately from local sandboxes:

```console
$ sbx --cloud ls
```

Most cloud commands accept either the sandbox name or the `sbx_`-prefixed ID
shown in the output.

## Run commands

Run a command inside a cloud sandbox:

```console
$ sbx --cloud exec cloud-project pwd
```

## Connect with SSH

Configure cloud SSH access and connect with your SSH client:

```console
$ sbx --cloud setup ssh
$ ssh sbx_01abc123@sbx_cloud
```

Replace the example with `ssh <sandbox-id>@sbx_cloud`, using the
`sbx_`-prefixed ID from `sbx --cloud ls`.

## Transfer files

Use `sbx --cloud cp` to copy files or directories between the client machine
and a cloud sandbox. Use absolute sandbox paths:

```console
$ sbx --cloud cp ./src cloud-project:/home/agent/workspace/src
$ sbx --cloud cp cloud-project:/home/agent/workspace/result.json ./result.json
```

Copying creates a point-in-time transfer. It doesn't mount or synchronize the
local path. For source control workflows, you can also clone a remote
repository from inside the sandbox and push changes to the remote. Configure
[cloud credentials](credentials.md) before creating a sandbox that needs access
to a private repository. The [cloud walkthrough](_index.md#get-started) shows a
public repository example.

## Expose a port

Expose a TCP service by specifying its sandbox port:

```console
$ sbx --cloud ports cloud-project --publish 8080
```

The command returns a public HTTPS URL assigned by the cloud control plane.
Cloud mode accepts a sandbox port number with an optional `/tcp` suffix, such
as `8080/tcp`. Host IP addresses, host port bindings, and other protocols are
rejected.

List or remove exposed ports:

```console
$ sbx --cloud ports cloud-project
$ sbx --cloud ports cloud-project --unpublish 8080
```

Treat an exposed URL as a public endpoint. Apply authentication in the service
and remove the exposure when you no longer need it.

## Configure expiration

Set the time-to-live and the action taken when it lapses during creation:

```console
$ sbx --cloud create --name cloud-project --ttl 2h --on-timeout delete claude
```

The default time-to-live is one hour. If you omit `--on-timeout`, the server
stops sandboxes that can be resumed and deletes the rest. Choose an action
explicitly when you need a particular outcome:

- `stop` preserves the sandbox so it can be started again. This requires
  support for stopping the sandbox.
- `restart` stops and immediately starts the sandbox. If you also specify
  `--ttl`, it must be at least one hour.
- `delete` removes the sandbox.

Volume-backed sandboxes require the `delete` action. Omitting `--ttl` uses
the server default. Setting `--ttl 0` is an error, not a way to disable
expiration.

Inspect or extend the expiration. Extensions cannot move expiration beyond
24 hours from creation:

```console
$ sbx --cloud ttl cloud-project
$ sbx --cloud ttl +30m cloud-project
```

## Stop or remove a sandbox

Stop a cloud sandbox while preserving its memory and filesystem:

```console
$ sbx --cloud stop cloud-project
```

The command returns when the stop request is accepted. Check `sbx --cloud ls`
to confirm that the sandbox has stopped. Compute isn't billed while it is
stopped.

To resume the sandbox and connect to its agent:

```console
$ sbx --cloud attach cloud-project
```

You can also use `sbx --cloud run claude --name cloud-project`, or run the agent
without `--name` and select the sandbox when prompted. Add `--detached` to a
named run to resume without attaching.

Resuming keeps the sandbox ID and state. Check its expiration with
`sbx --cloud ttl cloud-project` after resuming.

If stop or resume reports that an existing sandbox was not found, the operation
may be disabled for your account.

Volume-backed sandboxes can't be stopped. Remove a volume-backed sandbox to end
it and save the volume snapshot.

Remove a sandbox when you no longer need its state:

```console
$ sbx --cloud rm cloud-project
```

Removal asks for confirmation, deletes the cloud sandbox, and can't be undone.
Use `--force` to skip the prompt in scripts.

## Use persistent volumes

Cloud volumes are experimental and preserve data independently of a sandbox.
Create a volume, then attach it at sandbox creation:

```console
$ sbx --cloud volume create dependency-cache
$ sbx --cloud create --name cloud-project \
    --volume dependency-cache:/workspace/cache claude
```

The root directory of a newly created volume is owned by `root`. Change its
ownership after attaching it so the agent can write to it:

```console
$ sbx --cloud exec cloud-project \
    sudo chown agent:agent /workspace/cache
```

Volume data is saved as a snapshot when a sandbox exits, not continuously. If
multiple sandboxes mount the same volume at the same time, the last sandbox to
exit overwrites the stored snapshot.

## Customize a cloud sandbox

Cloud templates have their own store. A local template is not available to
`sbx --cloud` until you transfer it. To capture a running cloud sandbox and
create another sandbox from that template:

```console
$ sbx --cloud template save cloud-project cloud-template
$ sbx --cloud create --name cloud-copy --template cloud-template
```

The template supplies its CPU and memory configuration. Do not combine
`--template` with an agent name, `--cpus`, or `--memory`. To launch an OCI image
directly instead, use `--image-ref` with explicit CPU and memory values.

Snapshots include credentials written to the sandbox filesystem. Remove those
credentials before saving a template. Managed cloud secrets stay in the secret
store. See [Authenticate cloud agents](credentials.md).

Cloud sandboxes also support sandbox kits and `--kit` mixins. See
[Kits](../customize/kits.md) for customization and
[Local and cloud differences](local-vs-cloud.md) for host-dependent features.
Configure [cloud credentials](credentials.md) before adapting a local kit.

To declare reusable cloud configuration in a file, see
[Use a cloud environment](../configuration/environment-files.md#use-a-cloud-environment).

## Load an MCP server

[Connect an MCP server in Docker Agentic Platform](/manuals/agentic-platform/mcp.md)
before loading it into a cloud sandbox. Use the same Docker account you use
with `sbx`.

Load the connected server into a running cloud sandbox, using its name from
the console:

```console
$ sbx --cloud mcp load <server-name> --sandbox cloud-project
```

The server name is resolved by the MCP gateway associated with your Docker
Agentic Platform account. Cloud sandboxes don't use servers registered in the
local MCP store with `sbx mcp add`.

List servers reported by existing cloud sandbox gateways, or inspect one
sandbox's gateway:

```console
$ sbx --cloud mcp ls
$ sbx --cloud mcp ls cloud-project
```

The account listing shows servers reported by gateways, with the sandboxes
that use them. It omits unused server configurations and servers skipped by
a gateway. The sandbox view includes skipped servers.

## Diagnose cloud access

Check the CLI, Docker sign-in, cloud API connectivity, and account access:

```console
$ sbx --cloud diagnose
```

These checks don't require a local sandbox daemon. For local diagnostics,
see [Troubleshooting](../troubleshooting.md).
