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

A cloud sandbox expires after one hour by default and is deleted on expiration.
See [Configure expiration](#configure-expiration) to choose another timeout or
action before creating it.

> [!IMPORTANT]
>
> Credentials saved for local sandboxes aren't available in cloud sandboxes.
> [Configure a cloud credential](credentials.md) before launching an agent.
> Signing in through the agent can store credentials in the sandbox filesystem.

Create a sandbox and attach to its agent:

```console
$ sbx --cloud run claude --name cloud-project
```

To create the sandbox without opening an agent session, use `create`:

```console
$ sbx --cloud create --name cloud-project claude
```

The command prints the cloud sandbox ID. Attach by ID or name:

```console
$ sbx --cloud attach cloud-project
```

Without resource flags, a cloud sandbox starts with 2 CPUs and 4 GiB of
memory. Use `--cpus` and `--memory` to select another supported configuration.

Cloud creation doesn't accept workspace paths. For example,
`sbx --cloud run claude .` returns an error because `.` refers to the local
filesystem.

## Run without attaching

For scripts or terminals without interactive input, create and start an agent
without attaching:

```console
$ sbx --cloud run --detached claude --name cloud-task
```

A detached run creates a fresh sandbox. It does not resume a stopped sandbox.
Use `sbx --cloud exec` to run commands and `sbx --cloud rm` to clean up.

To detach from an interactive `run` or `attach` session while leaving the
agent running, press `Ctrl+\`. Reconnect with `sbx --cloud attach <sandbox-name>`.

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
Cloud mode accepts only the sandbox port number. Host IP addresses, host port
numbers, and protocol suffixes don't apply.

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

The default time-to-live is one hour, and the default timeout action is
`delete`. The `stop` action preserves the sandbox so it can be started again.
Volume-backed sandboxes require the `delete` action.

Inspect or extend the expiration. Extensions cannot move expiration beyond
24 hours from creation:

```console
$ sbx --cloud ttl cloud-project
$ sbx --cloud ttl +30m cloud-project
```

## Stop or remove a sandbox

Stop a cloud sandbox without deleting its filesystem:

```console
$ sbx --cloud stop cloud-project
```

Compute isn't billed while a sandbox is stopped. To resume it, run the agent
again and select the stopped sandbox when prompted:

```console
$ sbx --cloud run claude
```

Resuming keeps the sandbox ID and state and sets expiration to one hour after
resume. Inspect it with `sbx --cloud ttl cloud-project`. The `attach` command
doesn't resume a stopped sandbox. A detached run creates a fresh sandbox instead
of resuming the stopped one.

If stop or resume reports that an existing sandbox was not found, the operation
may be disabled for your account.

Volume-backed sandboxes can't be stopped. Remove a volume-backed sandbox to end
it and save the volume snapshot.

Remove a sandbox when you no longer need its state:

```console
$ sbx --cloud rm cloud-project
```

Removal deletes the cloud sandbox and can't be undone.

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

Cloud sandboxes also support sandbox kits and `--kit` mixins. Kit support does
not make host mounts, shared host skills, or arbitrary credential services
available in the cloud. Check [Local and cloud differences](local-vs-cloud.md)
and configure [cloud credentials](credentials.md) before adapting a local kit.

## Load an MCP server

First, [connect and authorize an MCP server](/manuals/agentic-platform/mcp.md)
in Docker Agentic Platform. Then load that server into a running cloud sandbox:

```console
$ sbx --cloud mcp load <server-name> --sandbox cloud-project
```

The server name is resolved by the MCP gateway associated with your Docker
Agentic Platform account. Cloud sandboxes don't use servers registered in the
local MCP store with `sbx mcp add`.
