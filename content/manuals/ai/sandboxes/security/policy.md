---
title: Policies
weight: 30
description: Configure network access rules for sandboxes.
---

{{< summary-bar feature_name="Docker Sandboxes sbx" >}}

Sandboxes are [network-isolated](isolation.md) from your host and from each
other. A policy system controls what a sandbox can access over the network.

Use the `sbx policy` command to configure network access rules. Rules apply
to all sandboxes on the machine.

## Network policies

The only way traffic can leave a sandbox is through an HTTP/HTTPS proxy on
your host, which enforces access rules on every outbound request.

Non-HTTP TCP traffic, including SSH, can be allowed by adding a policy rule
for the destination IP address and port (for example,
`sbx policy allow network "10.1.2.3:22"`). UDP and ICMP traffic is blocked
at the network layer and can't be unblocked with policy rules.

### Initial policy selection

On first start, and after running `sbx policy reset`, the daemon prompts you to
choose a network policy:

```plaintext
Choose a default network policy:

     1. Open         — All network traffic allowed, no restrictions.
     2. Balanced     — Default deny, with common dev sites allowed.
     3. Locked Down  — All network traffic blocked unless you allow it.

  Use ↑/↓ to navigate, Enter to select, or press 1–3.
```

| Policy      | Description                                                                                                                                                                                    |
| ----------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Open        | All outbound traffic is allowed. No restrictions. Equivalent to adding a wildcard allow rule with `sbx policy allow network "**"`.                                                             |
| Balanced    | Default deny, with a baseline allowlist covering AI provider APIs, package managers, code hosts, container registries, and common cloud services. You can extend this with `sbx policy allow`. |
| Locked Down | All outbound traffic is blocked, including model provider APIs (for example, `api.anthropic.com`). You must explicitly allow everything you need.                                              |

You can change your effective policy at any time using `sbx policy allow` and
`sbx policy deny`, or start over by running `sbx policy reset`.

### Non-interactive environments

In non-interactive environments such as CI pipelines or headless servers, the
interactive prompt can't be displayed. Use `sbx policy set-default` to set the
default network policy before running any other `sbx` commands:

```console
$ sbx policy set-default balanced
```

Available values are `allow-all`, `balanced`, and `deny-all`. After setting the
default, you can customize further with `sbx policy allow` and
`sbx policy deny` as usual.

### Default policy

All outbound HTTP/HTTPS traffic is blocked by default unless an explicit rule
allows it. The **Balanced** policy ships with a baseline allowlist covering AI provider
APIs, package managers, code hosts, container registries, and common cloud
services. Run `sbx policy ls` to see the active rules for your installation.

### Managing rules

Use [`sbx policy allow`](/reference/cli/sbx/policy/allow/) and
[`sbx policy deny`](/reference/cli/sbx/policy/deny/) to add network access
rules. Changes take effect immediately and apply to all sandboxes:

```console
$ sbx policy allow network api.anthropic.com
$ sbx policy deny network ads.example.com
```

Specify multiple hosts in one command with a comma-separated list:

```console
$ sbx policy allow network "api.anthropic.com,*.npmjs.org,*.pypi.org"
```

List all active policy rules with `sbx policy ls`:

```console
$ sbx policy ls
ID                                     TYPE      DECISION   RESOURCES
a1b2c3d4-e5f6-7890-abcd-ef1234567890   network   allow      api.anthropic.com, *.npmjs.org
f9e8d7c6-b5a4-3210-fedc-ba0987654321   network   deny       ads.example.com
```

Use `--type network` to show only network policies.

Remove a policy by resource or by rule ID:

```console
$ sbx policy rm network --resource ads.example.com
$ sbx policy rm network --id 2d3c1f0e-4a73-4e05-bc9d-f2f9a4b50d67
```

### Resetting to defaults

To remove all custom policies and restore the default policy, use
`sbx policy reset`:

```console
$ sbx policy reset
```

This deletes the local policy store and stops the daemon. When the daemon
restarts on the next command, you are prompted to choose a new network policy.
If sandboxes are running, they stop when the daemon shuts down. You are prompted for
confirmation unless you pass `--force`:

```console
$ sbx policy reset --force
```

### Switching to allow-by-default

If you prefer a permissive policy where all outbound traffic is allowed, add
a wildcard allow rule:

```console
$ sbx policy allow network "**"
```

This lets agents install packages and call any external API without additional
configuration. You can still deny specific hosts with `sbx policy deny`.

### Wildcard syntax

Rules support exact domains (`example.com`), wildcard subdomains
(`*.example.com`), and optional port suffixes (`example.com:443`).

Note that `example.com` doesn't match subdomains, and `*.example.com` doesn't
match the root domain. Specify both to cover both.

### Common patterns

Allow access to package managers so agents can install dependencies:

```console
$ sbx policy allow network "*.npmjs.org,*.pypi.org,files.pythonhosted.org,github.com"
```

The **Balanced** policy already includes AI provider APIs, package managers,
code hosts, and container registries. You only need to add allow rules for
additional domains your workflow requires. If you chose **Locked Down**, you
must explicitly allow everything.

> [!WARNING]
> Allowing broad domains like `github.com` permits access to any content on
> that domain, including user-generated content. Only allow domains you trust
> with your data.

### Monitoring

Use `sbx policy log` to see which hosts your sandboxes have contacted:

```console
$ sbx policy log
Blocked requests:
SANDBOX      TYPE     HOST                   PROXY        RULE       LAST SEEN        COUNT
my-sandbox   network  blocked.example.com    transparent  policykit  10:15:25 29-Jan  1

Allowed requests:
SANDBOX      TYPE     HOST                   PROXY        RULE       LAST SEEN        COUNT
my-sandbox   network  api.anthropic.com      forward      policykit  10:15:23 29-Jan  42
my-sandbox   network  registry.npmjs.org     transparent  policykit  10:15:20 29-Jan  18
```

The **PROXY** column shows how the request left the sandbox:

| Value         | Description                                                                                                    |
| ------------- | -------------------------------------------------------------------------------------------------------------- |
| `forward`     | Routed through the forward proxy. Supports [credential injection](credentials.md).                             |
| `transparent` | Intercepted by the transparent proxy. Policy is enforced but credential injection is not available.            |
| `network`     | Non-HTTP traffic (raw TCP, UDP, ICMP). TCP can be allowed with a policy rule; UDP and ICMP are always blocked. |

Filter by sandbox name by passing it as an argument:

```console
$ sbx policy log my-sandbox
```

Use `--limit N` to show only the last `N` entries, `--json` for
machine-readable output, or `--type network` to filter by policy type.

## Precedence

All outbound traffic is blocked by default unless an explicit rule allows it.
If a domain matches both an allow and a deny rule, the deny rule wins
regardless of specificity.

To unblock a domain, find the deny rule with `sbx policy ls` and remove it
with `sbx policy rm`.

## Troubleshooting

### Policy changes not taking effect

If policy changes aren't taking effect, run `sbx policy reset` to wipe the
local policy store and restart the daemon. On next start, you are prompted to
choose a new network policy.
