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
$ sbx --cloud exec cloud-project git status
```

Use `attach` to connect to the agent session, or configure SSH access and
connect with your SSH client:

```console
$ sbx --cloud setup ssh
$ ssh sbx_01abc123@sbx_cloud
```

Replace the example with `ssh <sandbox-id>@sbx_cloud`, using the
`sbx_`-prefixed ID from `sbx --cloud ls`.

## Transfer files

Use `sbx --cloud cp` to copy files or directories between the client machine
and a cloud sandbox:

```console
$ sbx --cloud cp ./src cloud-project:/workspace/src
$ sbx --cloud cp cloud-project:/workspace/result.json ./result.json
```

Copying creates a point-in-time transfer. It doesn't mount or synchronize the
local path. For source control workflows, you can also clone a remote
repository from inside the sandbox and push changes to the remote.

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

The default timeout action is `delete`. The `stop` action preserves the
sandbox so it can be started again, but it requires an eligible account. The
cloud service enforces a maximum lifetime of 24 hours from creation.

Inspect or extend the expiration, subject to that ceiling:

```console
$ sbx --cloud ttl cloud-project
$ sbx --cloud ttl +30m cloud-project
```

## Stop or remove a sandbox

Stop a cloud sandbox without deleting its filesystem:

```console
$ sbx --cloud stop cloud-project
```

Stopping and restarting cloud sandboxes requires the corresponding account
entitlement. Compute isn't billed while a sandbox is stopped. To resume it, run
the agent again and select the stopped sandbox when prompted:

```console
$ sbx --cloud run claude
```

The `attach` command doesn't resume a stopped sandbox. A detached run creates a
new sandbox instead of resuming the stopped one.

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

## Load an MCP server

Load an MCP server from Docker Agentic Platform into a running cloud sandbox:

```console
$ sbx --cloud mcp load SERVER --sandbox cloud-project
```

The server name is resolved by the MCP gateway associated with your Docker
Agentic Platform account. Cloud sandboxes don't use servers registered in the
local MCP store with `sbx mcp add`.
