---
title: MCP gateway
description: Register MCP servers, authorize OAuth-backed servers, and connect MCP tools to Docker Sandboxes.
keywords: docker sandboxes, sbx, MCP gateway, Model Context Protocol, MCP servers, sbx mcp, static MCP, OAuth
weight: 30
---

Docker Sandboxes includes an MCP gateway for connecting agents to Model Context
Protocol servers. The gateway gives the agent inside the sandbox one MCP
endpoint, while `sbx` manages the registered servers, OAuth credentials, and
sandbox lifecycle on the host.

This is different from configuring an MCP server directly in an agent such as
Claude Code. Direct MCP setup configures that agent's own MCP client. With
Docker Sandboxes, you register MCP servers once on the host, and the sandbox
gateway exposes them to supported agents inside isolated sandboxes. That
host-managed gateway provides a single path for credentials, dynamic or fixed
server sets, live server loading, and organization governance.

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
- For local containerized MCP servers, use a host with Docker installed and
  running.

## Quick start

Start by registering one MCP server on the host:

```console
$ sbx mcp add notion --url https://mcp.notion.com/mcp
```

If the server requires OAuth, `sbx` opens an authorization flow before it stores
the registration. After registration, verify that the server is in your local
MCP catalog:

```console
$ sbx mcp ls
NAME                 TYPE     URL/COMMAND
notion               remote   https://mcp.notion.com/mcp
```

Then start a sandbox:

```console
$ sbx run claude --name mcp-demo
```

The sandbox starts with an MCP gateway. In the default dynamic mode, the agent
can discover registered MCP servers through gateway tools and use them during
the session. Because you registered the server before starting the sandbox, it
is part of that sandbox's initial MCP catalog. The registration remains on the
host and can be reused by other sandboxes.

## Register an MCP server

`sbx mcp add` registers an MCP server by name. The registration records the
server definition on the host. It doesn't attach the server to a running sandbox
by itself. A server registered before sandbox creation is available when that
sandbox's gateway starts. To add a server to a sandbox that's already running,
use [`sbx mcp load`](#add-a-server-to-a-running-sandbox).

Server names can contain letters, numbers, dots, hyphens, and underscores.

### Remote endpoint URL

For a remote MCP endpoint, pass the server URL:

```console
$ sbx mcp add notion --url https://mcp.notion.com/mcp
$ sbx mcp add linear --url https://mcp.linear.app/mcp
```

### Local stdio server

Some MCP servers communicate over stdio instead of exposing a remote HTTP
endpoint. Docker Sandboxes can register local stdio servers in two ways: by
resolving a package from registry or manifest metadata, or by running an
explicit command.

#### From registry or manifest metadata

Use `--local` with an MCP community registry URL or a URL that returns a
`server.json` or `server.yaml` document. The registry entry or manifest must
describe an OCI package that uses stdio transport. `sbx` resolves the image from
the metadata and starts it on the host with Docker.

```console
$ sbx mcp add fetch --local \
  --url https://registry.modelcontextprotocol.io/v0/servers/fetch-mcp/versions/latest
```

The `--local` flag tells `sbx` to run the registry or manifest entry as a
host-side stdio server. If the entry doesn't publish a stdio package, `sbx`
rejects the registration instead of starting it locally.

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
  --args run --args -i --args --rm --args your/image
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

When a server isn't authorized, the gateway exposes a helper tool named
`<server>-authorize`, such as `notion-authorize`. The agent can call that tool
from inside the sandbox when it needs the server. The gateway also exposes
`<server>-revoke-auth` for OAuth-backed servers.

You can manage OAuth credentials from the host:

```console
$ sbx mcp auth status notion
$ sbx mcp auth notion
$ sbx mcp auth rm notion
```

Use `--all` to apply `auth`, `auth status`, or `auth rm` to all registered
OAuth-backed servers. Use `--format=json` for machine-readable output.

## Expose servers to a sandbox

Every sandbox starts an MCP gateway. When the sandbox starts, supported agent
integrations read the gateway URL and register it with the agent.

A sandbox can use MCP in dynamic mode or static mode. The mode is selected when
the sandbox is created.

| Mode    | How to use it             | Behavior                                                                  |
| ------- | ------------------------- | ------------------------------------------------------------------------- |
| Dynamic | Omit `--static-mcp`       | The registered catalog is searchable from the agent through gateway tools |
| Static  | Pass `--static-mcp` names | Only the named servers are pre-loaded, and discovery tools are disabled   |

### Dynamic mode

Dynamic mode is the default. Start a sandbox without `--static-mcp`:

```console
$ sbx run claude --name my-session
```

In dynamic mode, the agent can search the registered MCP catalog and connect
servers during the session. The gateway exposes discovery tools such as
`mcp-find` and `mcp-add`.

### Static mode

Use static mode when you want a fixed MCP server set for the sandbox:

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
static set is fixed at sandbox creation time. To use a different static set,
create a new sandbox.

## Add a server to a running sandbox

To attach an already-registered server to a running sandbox outside the initial
dynamic or static setup, use `sbx mcp load`:

```console
$ sbx mcp add linear --url https://mcp.linear.app/mcp
$ sbx mcp load linear --sandbox my-session
MCP server "linear" loaded into sandbox "my-session" (live)
```

Connected agent sessions receive a tool-list update, so the agent can discover
the added tools without reconnecting.

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
