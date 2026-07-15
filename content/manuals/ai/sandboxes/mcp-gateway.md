---
title: MCP gateway
description: Register MCP servers, authorize OAuth-backed servers, and connect MCP tools to Docker Sandboxes.
keywords: docker sandboxes, sbx, MCP gateway, Model Context Protocol, MCP servers, sbx mcp, static MCP, OAuth
weight: 30
---

Docker Sandboxes includes an MCP gateway for connecting agents to Model Context
Protocol servers. The gateway runs with the sandbox and presents one MCP
endpoint to the agent. Use `sbx mcp` on the host to register servers, manage
OAuth credentials, and attach registered servers to running sandboxes.

This page focuses on local gateway mode, where MCP traffic is handled by a
gateway running on your machine. Use `sbx mcp status` to check the effective
gateway for your account.

## Prerequisites

- Sign in with `sbx login`.
- Use an agent integration that configures MCP at startup: Claude Code, Codex,
  Gemini, Kiro, or OpenCode.
- For remote servers that require OAuth, use MCP servers that support OAuth
  Dynamic Client Registration.

The Docker Desktop MCP Toolkit isn't required for sandbox MCP support.

## Check the gateway mode

Run `sbx mcp status` to see the effective MCP gateway:

```console
$ sbx mcp status
Gateway mode: local
Gateway URL:  (none)
Source:       mcp-saas gateway-mode API (resolved by daemon)
```

The output can also include `Decision` and `Reason` fields. These fields show
the raw gateway decision and why that decision produced the effective gateway.

The daemon resolves the gateway mode and caches the result for its lifetime. If
you change sign-in state or an MCP gateway setting, restart the daemon before
checking again:

```console
$ sbx daemon stop
$ sbx mcp status
```

## Register an MCP server

Register servers on the host before you connect them to a sandbox. Server names
can contain letters, numbers, dots, hyphens, and underscores.

For a remote MCP endpoint, pass the server URL:

```console
$ sbx mcp add notion --url https://mcp.notion.com/mcp
$ sbx mcp add linear --url https://mcp.linear.app/mcp
```

You can also register a server from an MCP registry URL, a `server.json` or
`server.yaml` manifest URL, or a Docker Hardened Image reference that carries
an MCP server manifest:

```console
$ sbx mcp add fetch \
  --url https://registry.modelcontextprotocol.io/v0/servers/fetch-mcp/versions/latest

$ sbx mcp add opine --url https://example.com/mcp/opine/server.yaml

$ sbx mcp add fetch --url dhi.io/fetch-mcp:latest
```

To run a registry server locally as a container, add `--local`. This is useful
for stdio-based servers packaged as containers:

```console
$ sbx mcp add fetch --local \
  --url https://registry.modelcontextprotocol.io/v0/servers/fetch-mcp/versions/latest
```

You can also register a local stdio command:

```console
$ sbx mcp add playwright --command npx --args @playwright/mcp@latest
$ sbx mcp add local-image-server --command docker \
  --args run --args -i --args your/image
```

> [!WARNING]
> Local stdio commands run on the host, outside the sandbox. They run with your
> host user's permissions and can access host files, host network resources, and
> credentials available to that user. Use `--command` only with executables you
> trust.

`sbx mcp add` registers the server. It doesn't attach the server to a running
sandbox by itself. To connect registered servers to a sandbox, use dynamic mode,
static mode, or `sbx mcp load`.

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

## Start a sandbox with MCP

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

To attach an already-registered server to a running sandbox, use `sbx mcp
load`:

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
