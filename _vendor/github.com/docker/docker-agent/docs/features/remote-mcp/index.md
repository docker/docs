---
title: "Remote MCP Servers"
description: "Connect docker-agent to cloud services via remote MCP servers with built-in OAuth authentication."
keywords: docker agent, ai agents, features, remote mcp servers
weight: 120
canonical: https://docs.docker.com/ai/docker-agent/features/remote-mcp/
---

_Connect docker-agent to cloud services via remote MCP servers with built-in OAuth authentication._

## Overview

Docker Agent supports connecting to remote MCP servers over **Streamable HTTP**, **SSE** (Server-Sent Events), and **Unix domain sockets**. Streamable HTTP is the current recommended transport for most hosted MCP servers. Many popular services offer MCP endpoints with OAuth — docker-agent handles the authentication flow automatically.

```yaml
toolsets:
  - type: mcp
    remote:
      url: "https://mcp.linear.app/mcp"
      transport_type: "streamable"
```

## Unix Domain Sockets

Use a `unix://` URL to connect to an MCP server listening on a local Unix socket. This is useful when running docker-agent inside a container and exposing an MCP server from the host via a bind-mounted socket:

```yaml
toolsets:
  - type: mcp
    remote:
      url: "unix:///tmp/mcp-notify.sock"
      transport_type: "streamable"
```

The path after `unix://` is the absolute path to the socket file. Configured `headers` are forwarded over the socket connection. OAuth discovery is not supported for Unix socket URLs.

> [!TIP]
> **OAuth flow**
>
> When you connect to a remote MCP server that requires OAuth, docker-agent opens your browser automatically for authentication. Tokens are cached for subsequent sessions.

> [!TIP]
> **Cancelling the authorization dialog**
>
> If you dismiss the OAuth authorization dialog, the request is cancelled cleanly — no repeated prompts appear. The agent will report that authorization was declined. To try again, simply re-enable the server or repeat the request that triggered the flow.

## Configuration

```yaml
toolsets:
  - type: mcp
    remote:
      url: "https://mcp.example.com/mcp"
      transport_type: "streamable" # or "sse" for legacy servers
      headers:
        Authorization: "Bearer token" # optional: static auth
    # Optional: use only for trusted internal/private MCP or OAuth endpoints.
    allow_private_ips: true
```

For full configuration details, see the [Tool Config](../../configuration/tools/index.md) page.

Set `allow_private_ips: true` on a remote MCP toolset only when the MCP server or its OAuth registration/token endpoints intentionally resolve to private, loopback, or link-local addresses. The default blocks those OAuth helper requests to reduce SSRF risk.

> [!NOTE]
> **Headers forwarded during OAuth discovery**
>
> Configured `headers` are forwarded to OAuth protected-resource-metadata discovery requests directed at the MCP server's own host — not to third-party authorization servers. This allows services like Grafana Cloud that require a routing header (e.g. `X-Grafana-URL`) on the discovery request to scope the OAuth flow correctly. Headers are never sent to a different host than the one in `remote.url`.

> [!NOTE]
> **Automatic reconnection after idle timeouts**
>
> Remote MCP connections (Streamable HTTP / SSE) automatically reconnect after the server closes an idle connection — no configuration needed. Services like Notion and Linear close idle connections periodically; docker-agent detects the clean close and reconnects with exponential backoff. To tune reconnect behaviour or disable reconnection entirely, use the [`lifecycle` block](../../configuration/tools/index.md#toolset-lifecycle).

> [!NOTE]
> **Automatic recovery from revoked or rotated OAuth tokens**
>
> If a remote MCP server rejects the cached token with a `401 invalid_token` error (for example, because the token was revoked or rotated server-side), docker-agent handles the failure automatically:
>
> - **Silent refresh:** when a refresh token is available, docker-agent silently exchanges it for a new access token and replays the request — no user interaction required.
> - **Re-authentication prompt:** when the refresh token is absent or has also expired, the toolset transitions to a "needs re-auth" state and surfaces an OAuth prompt on your next message (exactly like the first-time flow).
>
> Either way, the agent never burns 5 reconnect attempts on an auth failure — it fails fast and either refreshes silently or defers to interactive re-auth. If you want to trigger re-auth immediately without waiting for the next message, run `/toolset-restart <name>` from the TUI.

### OAuth for servers without Dynamic Client Registration

Most remote MCP servers that require OAuth support [Dynamic Client Registration (RFC 7591)](https://datatracker.ietf.org/doc/html/rfc7591) — no configuration is needed, docker-agent handles the flow for you.

For servers that do **not** support DCR, provide explicit OAuth credentials with the `oauth:` block:

```yaml
toolsets:
  - type: mcp
    remote:
      url: "https://mcp.example.com/mcp"
      transport_type: "streamable"
      oauth:
        clientId: "my-app-client-id"
        clientSecret: "my-app-client-secret" # optional (public clients may omit)
        callbackPort: 8765                   # optional; picks a free port otherwise
        scopes:                              # optional; server-specific
          - read
          - write
```

| Field          | Type            | Required | Description                                                                                      |
| -------------- | --------------- | -------- | ------------------------------------------------------------------------------------------------ |
| `clientId`     | string          | ✓        | OAuth client ID registered with the remote MCP server.                                           |
| `clientSecret` | string          | ✗        | OAuth client secret. Omit for public clients using PKCE.                                         |
| `callbackPort` | integer         | ✗        | Local port to receive the OAuth redirect. If omitted, docker-agent picks a random free port.    |
| `scopes`       | array[string]   | ✗        | Scopes to request during the authorization step. Values are server-specific.                     |
| `callbackRedirectURL` | string   | ✗        | Custom OAuth redirect URI. Useful when the auth server requires HTTPS or a pre-registered URL. The literal placeholder `${callbackPort}` is replaced with the actual local callback port. See below.            |

Secrets should be stored in a credential helper or environment variable rather than committed — see [Secrets](../../guides/secrets/index.md) for interpolation patterns.

### Custom redirect URI (`callbackRedirectURL`)

Some authorization servers require the OAuth `redirect_uri` to be HTTPS or to match a URL that was pre-registered during app creation — neither of which plays nicely with a locally-bound loopback address such as `http://127.0.0.1:8765/callback`.

To work around this, set `callbackRedirectURL` to a public URL that redirects back to the local callback server. The literal placeholder `${callbackPort}` is substituted with the actual port the local callback server is listening on (either `callbackPort` when set, or the randomly-assigned port otherwise).

```yaml
toolsets:
  - type: mcp
    remote:
      url: "https://mcp.example.com/mcp"
      transport_type: "streamable"
      oauth:
        clientId: "my-app-client-id"
        callbackPort: 8765
        # Advertise this URL to the authorization server. The external
        # service at redirect.example.com is expected to 302-redirect the
        # browser to http://127.0.0.1:8765/callback preserving the query
        # string (code, state, …).
        callbackRedirectURL: "https://redirect.example.com/cb?port=${callbackPort}"
```

The local callback server still listens on the loopback interface on `callbackPort`; only the `redirect_uri` advertised to the authorization server changes.

**Validation rules:**

- The URL must be absolute (scheme + host) once `${callbackPort}` has been substituted.
- Only `http` and `https` schemes are accepted.
- `http` is only allowed when the host is a loopback address (`127.0.0.1`, `::1`, `localhost`); any other host must use `https` to avoid exposing the authorization `code` on the wire (RFC 8252 §7.3).

### Unmanaged OAuth flow (server mode)

When running `docker-agent serve api` (no local browser, no callback server), the runtime delegates the OAuth dance to the connected client via an MCP elicitation. There are two sub-behaviors, selected by the `--mcp-oauth-redirect-uri` flag:

- **`--mcp-oauth-redirect-uri=<URL>` set** (recommended for hosts like Docker Desktop): the runtime generates `state` + PKCE + (optional) Dynamic Client Registration in-process, builds the full authorize URL, and emits an elicitation whose `Meta` includes:

  | Key                          | Value                                                            |
  | ---------------------------- | ---------------------------------------------------------------- |
  | `docker-agent/type`                | `"oauth_flow"`                                                   |
  | `docker-agent/server_url`          | The MCP server URL (for display / favicon)                       |
  | `docker-agent/authorize_url`       | The full URL the client should open in the user's browser        |
  | `docker-agent/state`               | The `state` value the client must echo back when replying        |
  | `auth_server`                | Issuer of the authorization server                               |
  | `auth_server_metadata`       | RFC 8414 authorization-server metadata document                  |
  | `resource_metadata`          | RFC 9728 protected-resource metadata document                    |

  The client opens the browser at the URL provided in the `docker-agent/authorize_url` meta field, receives the OAuth callback at whatever endpoint the configured `redirect_uri` resolves to (typically a host-controlled bouncer that 302s into a deeplink), and replies to the elicitation with `accept` and `Content = {"code": "...", "state": "..."}`. The runtime verifies the `state`, exchanges the `code` at the token endpoint (using the same `redirect_uri` for RFC 6749 §4.1.3 binding), stores the token, and replays the original MCP request with `Authorization: Bearer ...`.

- **Flag not set** (client-driven): the runtime emits the elicitation meta below and expects the client to drive the OAuth flow itself (PKCE, DCR, token exchange) and reply with `Content = {"access_token": "...", "refresh_token": "...", ...}`:

  | Key                          | Value                                                            |
  | ---------------------------- | ---------------------------------------------------------------- |
  | `docker-agent/type`          | `"oauth_flow"`                                                   |
  | `docker-agent/server_url`    | The MCP server URL (for display / favicon)                       |
  | `auth_server`                | Issuer of the authorization server                               |
  | `auth_server_metadata`       | RFC 8414 authorization-server metadata document                  |
  | `resource_metadata`          | RFC 9728 protected-resource metadata document                    |

The client-driven `{access_token, ...}` reply shape is still accepted on the `--mcp-oauth-redirect-uri` path too: a client that prefers to do the exchange itself can ignore the `docker-agent/authorize_url`/`docker-agent/state` keys.

A per-toolset `callbackRedirectURL` (in the YAML) overrides the runtime-wide `--mcp-oauth-redirect-uri` for that toolset.

> [!WARNING]
> **Security note**
>
> The `POST /api/mcp-oauth/callback` route is open by default (no auth required) when `--auth-token` is unset. State values are 128-bit opaque tokens, so brute-force is infeasible, but a state value that leaks (e.g. via debug logs or a compromised host) could be exploited by an attacker to inject a code. Set `--auth-token` when `docker agent serve api` listens on a network-reachable interface. When set, `--auth-token` enforces Bearer-token authentication on all API routes including this callback endpoint.

## Project Management &amp; Collaboration

| Service    | URL                                | Transport | Description                           |
| ---------- | ---------------------------------- | --------- | ------------------------------------- |
| Asana      | `https://mcp.asana.com/sse`        | sse       | Task and project management           |
| Atlassian  | `https://mcp.atlassian.com/v1/mcp/authv2` | streamable | Jira, Confluence integration          |
| Linear     | `https://mcp.linear.app/mcp`       | streamable | Issue tracking and project management |
| Monday.com | `https://mcp.monday.com/sse`       | sse       | Work management platform              |
| Intercom   | `https://mcp.intercom.com/sse`     | sse       | Customer communication platform       |

## Development &amp; Infrastructure

| Service                  | URL                                            | Transport  | Description                       |
| ------------------------ | ---------------------------------------------- | ---------- | --------------------------------- |
| GitHub                   | `https://api.githubcopilot.com/mcp`            | sse        | Version control and collaboration |
| Buildkite                | `https://mcp.buildkite.com/mcp`                | streamable | CI/CD platform                    |
| Netlify                  | `https://netlify-mcp.netlify.app/mcp`          | streamable | Web hosting and deployment        |
| Vercel                   | `https://mcp.vercel.com/`                      | sse        | Web deployment platform           |
| Cloudflare Bindings      | `https://bindings.mcp.cloudflare.com/sse`      | sse        | Edge computing resources          |
| Cloudflare Observability | `https://observability.mcp.cloudflare.com/sse` | sse        | Monitoring and analytics          |
| Grafbase                 | `https://api.grafbase.com/mcp`                 | streamable | GraphQL backend platform          |
| Neon                     | `https://mcp.neon.tech/sse`                    | sse        | Serverless Postgres database      |
| Prisma                   | `https://mcp.prisma.io/mcp`                    | streamable | Database ORM and toolkit          |
| Sentry                   | `https://mcp.sentry.dev/sse`                   | sse        | Error tracking and monitoring     |

## Content &amp; Media

| Service    | URL                                               | Transport  | Description                       |
| ---------- | ------------------------------------------------- | ---------- | --------------------------------- |
| Canva      | `https://mcp.canva.com/mcp`                       | streamable | Design and graphics platform      |
| Miro       | `https://mcp.miro.com/`                           | streamable | Collaborative whiteboard platform (Enterprise plan required; see [official docs](https://developers.miro.com/docs/miro-mcp)) |
| Cloudinary | `https://asset-management.mcp.cloudinary.com/sse` | sse        | Media management and optimization |
| InVideo    | `https://mcp.invideo.io/sse`                      | sse        | Video creation platform           |
| Webflow    | `https://mcp.webflow.com/sse`                     | sse        | Website builder and CMS           |
| Wix        | `https://mcp.wix.com/sse`                         | sse        | Website builder platform          |
| Notion     | `https://mcp.notion.com/mcp`                      | streamable | Documentation and knowledge base  |

## Communication &amp; Voice

| Service     | URL                                 | Transport  | Description                 |
| ----------- | ----------------------------------- | ---------- | --------------------------- |
| Fireflies   | `https://api.fireflies.ai/mcp`      | streamable | Meeting transcription       |
| Listenetic  | `https://mcp.listenetic.com/v1/mcp` | streamable | Audio intelligence platform |
| Carbonvoice | `https://mcp.carbonvoice.app`       | sse        | Voice communication tools   |
| Telnyx      | `https://api.telnyx.com/v2/mcp`     | streamable | Communications platform     |
| Dialer      | `https://getdialer.app/sse`         | sse        | Phone communication tools   |

## Storage &amp; File Management

| Service | URL                                 | Transport | Description              |
| ------- | ----------------------------------- | --------- | ------------------------ |
| Box     | `https://mcp.box.com`               | sse       | Cloud content management |
| Egnyte  | `https://mcp-server.egnyte.com/sse` | sse       | Enterprise file sharing  |

## Business &amp; Finance

| Service       | URL                                       | Transport  | Description                |
| ------------- | ----------------------------------------- | ---------- | -------------------------- |
| PayPal        | `https://mcp.paypal.com/sse`              | sse        | Payment processing         |
| Plaid         | `https://api.dashboard.plaid.com/mcp/sse` | sse        | Financial data integration |
| Square        | `https://mcp.squareup.com/sse`            | sse        | Payment processing         |
| Close         | `https://mcp.close.com/mcp`               | streamable | CRM platform               |
| Dodo Payments | `https://mcp.dodopayments.com/sse`        | sse        | Payment processing         |

## Analytics &amp; Data

| Service     | URL                                     | Transport  | Description                    |
| ----------- | --------------------------------------- | ---------- | ------------------------------ |
| ThoughtSpot | `https://agent.thoughtspot.app/mcp`     | streamable | Analytics and BI platform      |
| Meta Ads    | `https://mcp.pipeboard.co/meta-ads-mcp` | streamable | Facebook advertising analytics |

## Utilities &amp; Tools

| Service       | URL                                | Transport  | Description                     |
| ------------- | ---------------------------------- | ---------- | ------------------------------- |
| Apify         | `https://mcp.apify.com`            | sse        | Web scraping and automation     |
| SimpleScraper | `https://mcp.simplescraper.io/mcp` | streamable | Web scraping tool               |
| GlobalPing    | `https://mcp.globalping.dev/sse`   | sse        | Network diagnostics             |
| Jam           | `https://mcp.jam.dev/mcp`          | streamable | Bug reporting and collaboration |

## Example: Multi-Service Agent

Combine multiple remote MCP servers in a single agent:

```yaml
agents:
  root:
    model: anthropic/claude-sonnet-4-5
    instruction: |
      You help manage projects and deployments.
    toolsets:
      - type: mcp
        remote:
          url: "https://mcp.linear.app/mcp"
          transport_type: "streamable"
        instruction: Use Linear for issue tracking.
      - type: mcp
        remote:
          url: "https://api.githubcopilot.com/mcp"
          transport_type: "sse"
        instruction: Use GitHub for code and PRs.
      - type: mcp
        remote:
          url: "https://mcp.vercel.com/"
          transport_type: "sse"
        instruction: Use Vercel for deployments.
```

> [!NOTE]
> **Growing list**
>
> This list is updated as more services add MCP support. If a service you use isn't listed, check their documentation — many providers are adding MCP endpoints regularly.
