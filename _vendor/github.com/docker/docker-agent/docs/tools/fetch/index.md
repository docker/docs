---
title: "Fetch Tool"
description: "Read content from HTTP/HTTPS URLs."
keywords: docker agent, ai agents, tools, toolsets, fetch tool
linkTitle: "Fetch"
weight: 50
canonical: https://docs.docker.com/ai/docker-agent/tools/fetch/
---

_Read content from HTTP/HTTPS URLs._

## Overview

The fetch tool lets agents retrieve content from one or more HTTP/HTTPS URLs. It is **read-only** — only `GET` requests are supported. The tool respects `robots.txt`, limits response size (1 MB per URL), and can return content as plain text, Markdown (converted from HTML), or raw HTML.

> [!NOTE]
> **GET only**
>
> The fetch tool does **not** support `POST`, `PUT`, `DELETE` or other methods, and does not expose request bodies or per-call custom headers (the toolset can still attach static [credential headers](#custom-headers) to every request). To call REST endpoints with other verbs, use the [API tool](../api/index.md) or an [OpenAPI toolset](../openapi/index.md).

## Configuration

```yaml
toolsets:
  - type: fetch
```

### Options

| Property            | Type          | Default | Description                                                                                                                                                                                                                                                                                                      |
| ------------------- | ------------- | ------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `timeout`           | int           | `30`    | Default request timeout in seconds (overridable per tool call).                                                                                                                                                                                                                                                  |
| `allowed_domains`   | array[string] | _none_  | Allow-list of hosts the tool may fetch. When set, every URL whose host is **not** in the list is rejected before any network call is made. Mutually exclusive with `blocked_domains`.                                                                                                                            |
| `blocked_domains`   | array[string] | _none_  | Deny-list of hosts the tool must not fetch. URLs whose host matches one of these patterns are rejected before any network call (including `robots.txt`) is made. Mutually exclusive with `allowed_domains`.                                                                                                      |
| `allow_private_ips` | boolean       | `false` | Opt in to dialling **non-public** IP addresses (loopback, RFC1918, link-local — including the cloud-metadata endpoint at `169.254.169.254` — multicast, and the unspecified address). Required to reach `localhost` / internal services. See [SSRF protection](#ssrf-protection-and-reaching-localhost) below. |
| `headers`           | map[string]string | _none_ | Static HTTP headers attached to **every** request the toolset issues (including `robots.txt`). Values support `${env.VAR}` for secrets. Caller-supplied entries override the default `User-Agent` and the format-driven `Accept` header. Headers are stripped on cross-host redirects so credentials never leak to a third-party host. See [Custom headers](#custom-headers) below. |

### Domain matching

Domain patterns in `allowed_domains` and `blocked_domains` use the following rules (case-insensitive):

- **Bare domain** — `example.com` matches the host `example.com` _and_ any subdomain such as `docs.example.com`. It does **not** match unrelated hosts that share a suffix (e.g. `badexample.com`).
- **Leading dot** — `.example.com` matches **only** strict subdomains (`docs.example.com`, `a.b.example.com`), not the apex `example.com`.
- **Wildcard glob** — `*.example.com` is an alias for the leading-dot form; the apex is excluded. The `*` is only valid as a leading `*.` token (entries like `foo.*`, `*.*.example.com`, or a bare `*` are rejected at config-load time).
- **IP literal** — IP addresses are matched exactly (`169.254.169.254`).
- **CIDR range** — `169.254.0.0/16`, `10.0.0.0/8`, `::1/128`, `fc00::/7`. Matches when the URL's host parses as an IP inside the network. Hostname hosts never match a CIDR pattern. Malformed CIDRs are rejected at config-load time.
- **Trailing dots** in FQDN-form URLs (`http://example.com./`) are stripped before matching, so they cannot bypass a deny-list entry.

The lists are mutually exclusive: a single fetch toolset may set either `allowed_domains` or `blocked_domains`, but not both.

When a list is configured, every redirect target is re-checked against the same list. A request to an allowed origin that redirects to a forbidden host is rejected before any data is read from the redirect.

> [!WARNING]
> **Limitations**
>
> Matching is purely string-based on the URL host. It does **not** perform DNS resolution and does **not** normalise alternative IP encodings (decimal `2852039166`, hex `0xa9.0xfe.0xa9.0xfe`, octal, etc. IPv4-mapped IPv6 addresses ARE normalized to their IPv4 form). If you need to deny access to a specific IP, also list its alternative encodings, or block at the network layer.

### Custom Timeout

```yaml
toolsets:
  - type: fetch
    timeout: 60
```

### Custom headers

Attach static headers — typically credentials — to every request. Values support `${env.VAR}` interpolation so secrets stay out of YAML, and headers are dropped on cross-host redirects so a redirect chain cannot leak them to a third-party host:

```yaml
toolsets:
  - type: fetch
    allowed_domains:
      - docs.internal.example.com
    headers:
      Authorization: "Bearer ${env.INTERNAL_DOCS_TOKEN}"
      X-Internal-Client: "docker-agent"
```

> [!WARNING]
> **Pair credential headers with an allow-list**
>
> When `headers` carries credentials (e.g. `Authorization`), set `allowed_domains` to the specific hosts that should receive them. Stdlib already strips a small allow-list (`Authorization`, `Cookie`, `WWW-Authenticate`) on cross-domain redirects, and the fetch tool additionally strips every operator-supplied header on cross-host redirects — but an allow-list is the strongest guarantee against accidental exfiltration.

### Restrict to specific domains

```yaml
toolsets:
  - type: fetch
    allowed_domains:
      - docker.com          # docker.com and *.docker.com
      - github.com          # github.com and *.github.com
      - .githubusercontent.com  # only subdomains, e.g. raw.githubusercontent.com
```

### Block sensitive hosts

```yaml
toolsets:
  - type: fetch
    blocked_domains:
      - 169.254.169.254       # cloud metadata endpoint (literal IP)
      - 169.254.0.0/16        # entire link-local range (CIDR)
      - 10.0.0.0/8            # RFC1918 private range
      - "*.internal.example.com"  # any subdomain (wildcard)
      - internal.example.com  # internal corporate hostname
```

> [!NOTE]
> **Already blocked by default**
>
> You do **not** need to add loopback, RFC1918, link-local (incl. `169.254.169.254`), multicast or the unspecified address to `blocked_domains` to be safe — the fetch tool already refuses connections to those ranges at dial time, after DNS resolution. The example above is only useful if you also want to reject those hosts _before_ any network call (and to surface a clearer error message to the agent), or if you have set `allow_private_ips: true` and want to deny a specific subset.

### SSRF protection and reaching localhost

By default, the fetch tool refuses connections to **non-public IP addresses** — even when DNS for an otherwise-public host resolves to one of them (so DNS rebinding is also blocked). The check happens at dial time, after DNS resolution, and rejects:

- **Loopback** — `127.0.0.0/8`, `::1` (this is what blocks `http://localhost/...` and `http://127.0.0.1/...`)
- **RFC1918 private ranges** — `10.0.0.0/8`, `172.16.0.0/12`, `192.168.0.0/16`
- **Link-local** — `169.254.0.0/16` (IPv4, including the cloud-metadata endpoint `169.254.169.254`) and `fe80::/10` (IPv6)
- **Multicast** and the **unspecified** address (`0.0.0.0`, `::`)
- **IPv4-mapped IPv6** — addresses like `::ffff:127.0.0.1` or `::ffff:169.254.169.254` are normalized to their IPv4 form and blocked accordingly

This is the default because LLM-driven fetches are a classic Server-Side Request Forgery (SSRF) vector: a prompt-injected URL can otherwise reach internal services, cloud metadata, or admin interfaces on the host running the agent.

If an agent legitimately needs to call **localhost** or an **internal service**, opt in with `allow_private_ips: true`:

```yaml
toolsets:
  - type: fetch
    allow_private_ips: true
    allowed_domains:
      - localhost
      - 127.0.0.1
      - 10.0.0.0/8            # internal corporate range
```

> [!WARNING]
> **Pair with an allow-list**
>
> Setting `allow_private_ips: true` alone re-exposes the SSRF surface. We strongly recommend combining it with an `allowed_domains` entry that restricts the tool to the specific internal hosts or CIDRs the agent actually needs (e.g. `localhost`, `127.0.0.1`, or your internal CIDR).
>
> **Note:** `allowed_domains` is checked _before_ DNS resolution (string-based on hostname), while the SSRF check happens _after_ DNS resolution (on the resolved IP). This means `allowed_domains` and `blocked_domains` are evaluated independently of `allow_private_ips` and continue to apply. A public hostname in `allowed_domains` that resolves to a private IP will still be blocked unless `allow_private_ips: true` is set.

## Tool Interface

The toolset exposes a single tool, `fetch`, with the following parameters:

| Parameter | Type           | Required | Description                                                                                                 |
| --------- | -------------- | -------- | ----------------------------------------------------------------------------------------------------------- |
| `urls`    | array[string]  | ✓        | One or more HTTP/HTTPS URLs to fetch (all via `GET`).                                                       |
| `format`  | string         | ✓        | Output format: `text`, `markdown`, or `html`. HTML responses are converted to text/markdown when requested. |
| `timeout` | integer        | ✗        | Per-call request timeout in seconds. Overrides the toolset default. Valid range: `1`–`300`.                 |

Responses are capped at **1 MB** per URL. Hosts that disallow the agent's user-agent via `robots.txt` are skipped with a clear error.

> [!TIP]
> **Fetch vs. API Tool**
>
> Use `fetch` when the agent needs to read arbitrary public URLs at runtime. Use the [API tool](../api/index.md) to expose specific, structured HTTP endpoints (including non-`GET` verbs) as named tools.

## Domain Filtering

The `allowed_domains`, `blocked_domains`, and `allow_private_ips` options let you control which hosts the fetch tool may reach. The complete reference is in the [Options](#options) table and [Domain matching](#domain-matching) section above.

**Key points:**

- `allowed_domains` — allow-list; only listed hosts (and their subdomains for bare-domain entries) are reachable
- `blocked_domains` — deny-list; mutually exclusive with `allowed_domains` (a config error is thrown if both are set)
- `allow_private_ips` — defaults to `false`; set to `true` to reach loopback / RFC-1918 / link-local addresses
- The same `allow_private_ips` flag is also supported on `api`, `openapi`, `a2a`, and remote `mcp` toolsets

See [`examples/fetch_domain_filtering.yaml`](https://github.com/docker/docker-agent/blob/main/examples/fetch_domain_filtering.yaml) for a complete filtering example, and [`examples/remote_mcp_allow_private_ips.yaml`](https://github.com/docker/docker-agent/blob/main/examples/remote_mcp_allow_private_ips.yaml) for the equivalent pattern on remote MCP toolsets.
