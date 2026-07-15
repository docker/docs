---
title: MCP gateway
description: Register MCP servers, authorize OAuth-backed servers, and connect MCP tools to Docker Sandboxes.
keywords: docker sandboxes, sbx, MCP gateway, Model Context Protocol, MCP servers, sbx mcp, static MCP, OAuth
weight: 50
---

Docker Sandboxes includes an MCP gateway for connecting agents to Model Context
Protocol servers. The gateway gives the agent inside the sandbox one MCP
endpoint, while `sbx` manages the registered servers, OAuth credentials, and
sandbox lifecycle on the host.

This is different from configuring an MCP server directly in an agent such as
Claude Code. Direct MCP setup configures that agent's own MCP client. With
Docker Sandboxes, you register MCP servers once on the host, and the sandbox
gateway exposes them to supported agents inside isolated sandboxes. That
host-managed gateway provides a single path for credentials, explicit server
loading, live updates, and organization governance.

> [!NOTE]
> The Docker Sandboxes MCP gateway is separate from the Docker Desktop MCP
> Toolkit. You don't need the Docker Desktop MCP Toolkit to use `sbx mcp`, and
> MCP Toolkit server settings aren't shared with Docker Sandboxes.

## Prerequisites

- Sign in with `sbx login`.
- Use an agent integration that configures MCP at startup: Claude Code, Codex,
  Gemini, Kiro, or OpenCode.
- For remote servers that require OAuth, use MCP servers that support OAuth
  Dynamic Client Registration.
- For `--local --url` registrations that resolve to OCI packages, use a host
  with Docker installed and running. Docker is also required for explicit
  `--command docker ...` registrations.

## Quick start

Start by registering one MCP server on the host:

```console
$ sbx mcp add notion --url https://mcp.notion.com/mcp
```

If the server requires OAuth, `sbx` opens an authorization flow before it stores
the registration. After registration, verify that the server is registered:

```console
$ sbx mcp ls
NAME                 TYPE     URL/COMMAND
notion               remote   https://mcp.notion.com/mcp
```

Then start a sandbox and expose the registered server:

```console
$ sbx run claude --name mcp-demo --static-mcp notion
```

The sandbox starts with an MCP gateway and pre-loads the `notion` server. The
registration remains on the host and can be reused by other sandboxes.

## Register an MCP server

`sbx mcp add` registers an MCP server by name. The registration records the
server definition on the host. It doesn't attach the server to a sandbox by
itself. To expose a registered server to a sandbox, pass it with
[`--static-mcp`](#expose-servers-at-sandbox-creation) when you create the
sandbox, or use [`sbx mcp load`](#add-a-server-to-a-running-sandbox) for a
sandbox that's already running.

Server names can contain letters, numbers, dots, hyphens, and underscores.

The `--url` flag can point to different kinds of input. The execution location
depends on what you register:

- **Remote endpoint URL**: The URL is the running MCP server endpoint. The
  server runs remotely, and the sandbox gateway connects to it.
- **Metadata URL with `--local`**: The URL returns a registry entry,
  `server.json`, or `server.yaml` that describes an OCI-packaged stdio server.
  `sbx` resolves the image and runs it on the host with Docker.
- **Explicit command**: `sbx` runs the command on the host as a stdio MCP
  server.

Local stdio servers run on the host, not inside the sandbox. The agent inside
the sandbox connects only to the MCP gateway.

### Remote endpoint URL

For a remote MCP endpoint, pass the server URL:

```console
$ sbx mcp add notion --url https://mcp.notion.com/mcp
$ sbx mcp add linear --url https://mcp.linear.app/mcp
```

### Local stdio server

Some MCP servers communicate over stdio instead of exposing a remote HTTP
endpoint. Use a local stdio server when `sbx` should launch the MCP server on
the host. You can provide a metadata URL or an explicit command.

#### From registry or manifest metadata

Use `--local --url` when you have an MCP community registry URL or a URL that
returns a `server.json` or `server.yaml` document. The registry entry or
manifest must describe an OCI package that uses stdio transport. `sbx` doesn't
launch non-OCI package types, such as `npm`, from metadata. To use those
servers, register an explicit command.

This path resolves the image from the metadata and starts it on the host with
Docker, so Docker must be installed and running on the host.

```console
$ sbx mcp add fetch --local \
  --url https://registry.modelcontextprotocol.io/v0/servers/fetch-mcp/versions/latest
```

If the entry doesn't publish an OCI stdio package, `sbx` rejects the
registration instead of starting it locally.

A server manifest describes the MCP server package and how to start it. It can
be hosted on a GitHub raw URL, internal HTTP server, or CDN.

```console
$ sbx mcp add opine --local --url https://example.com/mcp/opine/server.yaml
```

#### From an explicit command

Use `--command` when you already know the executable and arguments, or when you
need custom Docker flags. The command can be a package runner such as `npx` or
a Docker container command:

```console
$ sbx mcp add playwright --command npx --args @playwright/mcp@latest
$ sbx mcp add local-image-server --command docker \
  --args "run,-i,--rm,your/image"
```

Use registry or manifest metadata when you have a published server definition
and don't need to customize `docker run`. Use `--command` for local
development, private servers, or custom container flags.

> [!WARNING]
> Local stdio servers run on the host, outside sandbox isolation. If the command
> starts a Docker container, that container uses host Docker isolation, not
> sandbox isolation. The process or container can access host files, host
> network resources, and credentials made available to it. Use trusted commands
> and images, and avoid mounting host paths or passing credentials unless the
> server needs them.

## Authorize OAuth-backed servers

If a registered remote server requires OAuth, `sbx mcp add` starts the
authorization flow by default:

```console
$ sbx mcp add notion --url https://mcp.notion.com/mcp
Resolving MCP server "notion"...
Open this URL to authorize MCP server "notion":
https://api.notion.com/v1/oauth/authorize?...
MCP server "notion" authorized
MCP server "notion" registered (type: remote)
```

OAuth credentials stay on the host. In local gateway mode, `sbx` stores tokens
in the host operating system's credential store.

To register an OAuth-backed server without authorizing it, pass `--skip_auth`:

```console
$ sbx mcp add notion --url https://mcp.notion.com/mcp --skip_auth
```

When an unauthorized OAuth-backed server is exposed to a sandbox, the gateway
exposes a helper tool named `<server>-authorize`, such as `notion-authorize`.
The agent can call that tool from inside the sandbox when it needs the server.

You can manage OAuth credentials from the host:

```console
$ sbx mcp auth status notion
$ sbx mcp auth notion
$ sbx mcp auth rm notion
```

Use `--all` to apply `auth`, `auth status`, or `auth rm` to all registered
OAuth-backed servers. Use `--format=json` for machine-readable output.

## Expose servers at sandbox creation

Every sandbox starts an MCP gateway. When the sandbox starts, supported agent
integrations read the gateway URL and register it with the agent.

Use `--static-mcp` when you want to pre-load registered MCP servers into a
sandbox:

```console
$ sbx mcp add notion --url https://mcp.notion.com/mcp
$ sbx mcp add linear --url https://mcp.linear.app/mcp
$ sbx run claude --name my-session --static-mcp notion,linear
```

You can pass `--static-mcp` as a comma-separated list or repeat the flag:

```console
$ sbx run claude --name my-session \
  --static-mcp notion --static-mcp linear
```

Every name in the static set must already be registered with `sbx mcp add`. The
static set is fixed at sandbox creation time. If you omit `--static-mcp`, the
sandbox starts with an MCP gateway but no registered MCP servers. To use a
different static set, create a new sandbox.

## Add a server to a running sandbox

To attach an already-registered server to a running sandbox, use
`sbx mcp load`:

```console
$ sbx mcp add linear --url https://mcp.linear.app/mcp
$ sbx mcp load linear --sandbox my-session
MCP server "linear" loaded into sandbox "my-session" (live)
```

Connected agent sessions receive a tool-list update, so the added tools become
visible without reconnecting.

## Built-in gateway tools

The local MCP gateway exposes a small set of built-in tools. These tools belong
to the gateway itself, not to a registered MCP server.

| Tool                 | Description                                                                                    |
| -------------------- | ---------------------------------------------------------------------------------------------- |
| `mcp-exec`           | Executes a tool by name in the current gateway session.                                        |
| `code-mode`          | Creates a session-scoped JavaScript tool that can call selected tools through the MCP gateway. |
| `<server>-authorize` | Starts OAuth authorization for an exposed server that requires authorization.                  |

The gateway exposes `<server>-authorize` only when an OAuth-backed server needs
authorization. If `code-mode` creates a generated tool, that generated tool is
available only in the current session.

In MCP access policies, built-in gateway tools are `MCP::Primordial` resources
and use the `invokePrimordial` action. Tools from registered MCP servers are
`MCP::Tool` resources and use the `invokeTool` action. For details, see the
[MCP policy reference](governance/reference/mcp-policy.md).

## Manage registrations

List registered servers:

```console
$ sbx mcp ls
```

Inspect a registered server:

```console
$ sbx mcp inspect notion
```

Remove a registered server:

```console
$ sbx mcp rm notion
```

For OAuth-backed servers, `sbx mcp rm` removes the OAuth credential before it
removes the server registration. To remove only the OAuth credential, use
`sbx mcp auth rm`.

## Register a bundle

An MCP bundle is a JSON array of server definitions fetched from a URL. Each
entry maps to one `sbx mcp add` registration.

```json
[
  { "name": "notion", "url": "https://mcp.notion.com/mcp" },
  { "name": "linear", "url": "https://mcp.linear.app/mcp" },
  {
    "name": "github",
    "command": "npx",
    "args": ["@modelcontextprotocol/server-github"]
  }
]
```

Fetch and register a bundle:

```console
$ sbx mcp bundle add core --url https://example.com/mcp-bundle.json
```

List registered bundles:

```console
$ sbx mcp bundle ls
```

Remove a bundle and the servers it registered:

```console
$ sbx mcp bundle rm core
```

## Governance

Organizations with AI Governance can use
[MCP access policies](governance/access-controls/mcp.md) to control MCP server
registration, tool calls, gateway meta-tools, resources, prompts, and approval
requirements. MCP access policies are organization policies written in Cedar.
