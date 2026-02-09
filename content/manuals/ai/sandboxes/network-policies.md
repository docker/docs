---
title: Network policies
description: Configure network filtering policies to control outbound HTTP and HTTPS access from sandboxed agents.
weight: 55
---

{{< summary-bar feature_name="Docker Sandboxes" >}}

Network policies control what external resources a sandbox can access through
an HTTP/HTTPS filtering proxy. Use policies to prevent agents from accessing
internal networks, enforce compliance requirements, or restrict internet access
to specific services.

Each sandbox has a filtering proxy that enforces policies for outbound HTTP and
HTTPS traffic. Connection to external services over other protocols, such as
raw TCP and UDP connections, are blocked. All agent communication must go
through the HTTP proxy or remain contained within the sandbox.

The proxy runs on an ephemeral port on your host, but from the agent
container's perspective it is accessible at `host.docker.internal:3128`.

### Security considerations

Use network policies as one layer of security, not the only layer. The microVM
boundary provides the primary isolation. Network policies add an additional
control for HTTP/HTTPS traffic.

The network filtering restricts which domains processes can connect to, but
doesn't inspect the traffic content. When configuring policies:

- Allowing broad domains like `github.com` permits access to any content on
  that domain, including user-generated content. Agents could use these as
  channels for data exfiltration.
- Domain fronting techniques may bypass filters by routing traffic through
  allowed domains. This is inherent to HTTPS proxying.

Only allow domains you trust with your data. You're responsible for
understanding what your policies permit.

## How network filtering works

Each sandbox has an HTTP/HTTPS proxy running on your host. The proxy is
accessible from inside the sandbox at `host.docker.internal:3128`.

When an agent makes HTTP or HTTPS requests, the proxy:

1. Checks the policy rules against the host in the request. If the host is
   blocked, the request is stopped immediately
2. Connects to the server. If the host was not explicitly allowed, checks the
   server's IP address against BlockCIDR rules

For example, `localhost` is not in the default allow-list, but it's allowed by the
default "allow" policy. Because it's not explicitly allowed, the proxy then checks
the resolved IP addresses (`127.0.0.1` or `::1`) against the BlockCIDR rules.
Since `127.0.0.1/8` and `::1/128` are both blocked by default, this catches any
DNS aliases for localhost like `ip6-localhost`.

If an agent needs access to a service on localhost, include a port number in
the allow-list (e.g., `localhost:1234`).

HTTP requests to `host.docker.internal` are rewritten to `localhost`, so only
the name `localhost` will work in the allow-list.

## Default policy

New sandboxes use this default policy unless you configure a custom policy:

**Policy mode:** `allow` (permit all traffic except what's explicitly blocked)

**Blocked CIDRs:**

- `10.0.0.0/8` - Private network (Class A)
- `127.0.0.0/8` - Loopback addresses
- `169.254.0.0/16` - Link-local addresses
- `172.16.0.0/12` - Private network (Class B)
- `192.168.0.0/16` - Private network (Class C)
- `::1/128` - IPv6 loopback
- `fc00::/7` - IPv6 unique local addresses
- `fe80::/10` - IPv6 link-local addresses

**Allowed hosts:**

- `*.anthropic.com` - Claude API and services
- `platform.claude.com:443` - Claude platform services

The default policy blocks access to private networks, localhost, and cloud
metadata services while allowing internet access. Explicitly allowed hosts
bypass CIDR checks for performance.

## Monitor network activity

View what your agent is accessing and whether requests are being blocked:

```console
$ docker sandbox network log
```

Network logs help you understand agent behavior and define policies.

## Applying policies

Use `docker sandbox network proxy` to configure network policies for your
sandboxes. The sandbox must be running when you apply policy changes. Changes
take effect immediately and persist across sandbox restarts.

### Example: Block internal networks

Prevent agents from accessing your local network, Docker networks, and cloud
metadata services. It prevents accidental access to internal services while
allowing agents to install packages and access public APIs.

```console
$ docker sandbox network proxy my-sandbox \
  --policy allow \
  --block-cidr 10.0.0.0/8 \
  --block-cidr 172.16.0.0/12 \
  --block-cidr 192.168.0.0/16 \
  --block-cidr 127.0.0.0/8
```

This policy:

- Allows internet access
- Blocks RFC 1918 private networks (10.x.x.x, 172.16-31.x.x, 192.168.x.x)
- Blocks localhost (127.x.x.x)

> [!NOTE]
> These CIDR blocks are already blocked by default. The example above shows how
> to explicitly configure them. See [Default policy](#default-policy) for the
> complete list.

### Example: Restrict to package managers only

For strict control, use a denylist policy that only allows package repositories:

```console
$ docker sandbox network proxy my-sandbox \
  --policy deny \
  --allow-host "*.npmjs.org" \
  --allow-host "*.pypi.org" \
  --allow-host "files.pythonhosted.org" \
  --allow-host "*.rubygems.org" \
  --allow-host github.com
```

> [!NOTE]
> This policy would block the backend of your AI coding agent (e.g., for Claude
> Code: `xyz.anthropic.com`). Make sure you also include hostnames that the
> agent requires.

The agent can install dependencies but can't access arbitrary internet
resources. This is useful for CI/CD environments or when running untrusted code.

### Example: Allow AI APIs and development tools

Combine AI service access with package managers and version control:

```console
$ docker sandbox network proxy my-sandbox \
  --policy deny \
  --allow-host api.anthropic.com \
  --allow-host "*.npmjs.org" \
  --allow-host "*.pypi.org" \
  --allow-host github.com \
  --allow-host "*.githubusercontent.com"
```

This allows agents to call AI APIs, install packages, and fetch code from
GitHub while blocking other internet access.

### Example: Allow specific APIs while blocking subdomains

Use port-specific rules and subdomain patterns for fine-grained control:

```console
$ docker sandbox network proxy my-sandbox \
  --policy deny \
  --allow-host api.example.com:443 \
  --allow-host cdn.example.com \
  --allow-host "*.storage.example.com:443"
```

This policy allows:

- Requests to `api.example.com` on port 443 (typically
  `https://api.example.com`)
- Requests to `cdn.example.com` on any port
- Requests to any subdomain of `storage.example.com` on port 443, like
  `us-west.storage.example.com:443` or `eu-central.storage.example.com:443`

Requests to `example.com` (any port), `www.example.com` (any port), or
`api.example.com:8080` would all be blocked because none match the specific
patterns.

To allow both a domain and all its subdomains, specify both patterns:

```console
$ docker sandbox network proxy my-sandbox \
  --policy deny \
  --allow-host example.com \
  --allow-host "*.example.com"
```

## Policy modes: allowlist versus denylist

Policies have two modes that determine default behavior.

### Allowlist mode

Default: Allow all traffic, block specific destinations.

```console
$ docker sandbox network proxy my-sandbox \
  --policy allow \
  --block-cidr 192.0.2.0/24 \
  --block-host example.com
```

Use allowlist mode when you want agents to access most resources but need to
block specific networks or domains. This is less restrictive and works well for
development environments where the agent needs broad internet access.

### Denylist mode

Default: Block all traffic, allow specific destinations.

```console
$ docker sandbox network proxy my-sandbox \
  --policy deny \
  --allow-host api.anthropic.com \
  --allow-host "*.npmjs.org"
```

Use denylist mode when you want tight control over external access. This
requires explicitly allowing each service the agent needs, making it more
secure but also more restrictive. Good for production deployments, CI/CD
pipelines, or untrusted code.

### Domain and CIDR matching

Domain patterns support exact matches, wildcards, and port specifications:

- `example.com` matches only that exact domain (any port)
- `example.com:443` matches requests to `example.com` on port 443 (the default
  HTTPS port)
- `*.example.com` matches all subdomains like `api.example.com` or
  `www.example.com`
- `*.example.com:443` matches requests to any subdomain on port 443

Important pattern behaviors:

- `example.com` does NOT match subdomains. A request to `api.example.com`
  won't match a rule for `example.com`.
- `*.example.com` does NOT match the root domain. A request to `example.com`
  won't match a rule for `*.example.com`.
- To allow both a domain and its subdomains, specify both patterns:
  `example.com` and `*.example.com`.

When multiple patterns could match a request, the most specific pattern wins:

1. Exact hostname and port: `api.example.com:443`
2. Exact hostname, any port: `api.example.com`
3. Wildcard patterns (longest match first): `*.api.example.com:443`,
   `*.api.example.com`, `*.example.com:443`, `*.example.com`
4. Catch-all wildcards: `*:443`, `*`
5. Default policy (allow or deny)

This specificity lets you block broad patterns while allowing specific
exceptions. For example, you can block `example.com` and `*.example.com` but
allow `api.example.com:443`.

CIDR blocks match IP addresses after DNS resolution:

- `192.0.2.0/24` blocks all 192.0.2.x addresses
- `198.51.100.0/24` blocks all 198.51.100.x addresses
- `203.0.113.0/24` blocks all 203.0.113.x addresses

When you block or allow a domain, the proxy resolves it to IP addresses and
checks those IPs against CIDR rules. This means blocking a CIDR range affects
any domain that resolves to an IP in that range.

## Bypass mode for HTTPS tunneling

By default, the proxy acts as a man-in-the-middle for HTTPS connections,
terminating TLS and re-encrypting traffic with its own certificate authority.
This allows the proxy to enforce policies and inject authentication credentials
for certain services. The sandbox container trusts the proxy's CA certificate.

Some applications use certificate pinning or other techniques that don't work
with MITM proxies. For these cases, use bypass mode to tunnel HTTPS traffic
directly without inspection:

```console
$ docker sandbox network proxy my-sandbox \
  --bypass-host api.service-with-pinning.com
```

You can also bypass traffic to specific IP ranges:

```console
$ docker sandbox network proxy my-sandbox \
  --bypass-cidr 203.0.113.0/24
```

When traffic is bypassed, the proxy:

- Acts as a simple TCP tunnel without inspecting content
- Cannot inject authentication credentials into requests
- Cannot detect domain fronting or other evasion techniques
- Still enforces domain and port matching based on the initial connection

Use bypass mode only when necessary. Bypassed traffic loses the visibility and
security benefits of MITM inspection. If you bypass broad domains like
`github.com`, the proxy has no visibility into what the agent accesses on that
domain.

## Policy persistence

Network policies are stored in JSON configuration files.

### Per-sandbox configuration

When you run `docker sandbox network proxy my-sandbox`, the command updates the
configuration for that specific sandbox only. The policy is persisted to
`~/.docker/sandboxes/vm/my-sandbox/proxy-config.json`.

The default policy (used when creating new sandboxes) remains unchanged unless
you modify it directly.

### Default configuration

The default network policy for new sandboxes is stored at
`~/.sandboxd/proxy-config.json`. This file is created automatically when the
first sandbox starts, but only if it doesn't already exist.

The current default policy is `allow`, which permits all outbound connections
except to blocked CIDR ranges (private networks, localhost, and cloud metadata
services).

You can modify the default policy:

1. Edit `~/.sandboxd/proxy-config.json` directly
2. Or copy a sandbox policy to the default location:

   ```console
   $ cp ~/.docker/sandboxes/vm/my-sandbox/proxy-config.json ~/.sandboxd/proxy-config.json
   ```

Review and customize the default policy to match your security requirements
before creating new sandboxes. Once created, the default policy file won't be
modified by upgrades, preserving your custom configuration.

## Next steps

- [Architecture](architecture.md)
- [Using sandboxes effectively](workflows.md)
- [Custom templates](templates.md)
