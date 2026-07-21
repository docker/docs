---
title: Configure an upstream proxy
linkTitle: Upstream proxy
description: Route sandbox and daemon traffic through a corporate or upstream proxy, including PAC files, SOCKS5, and your OS system proxy.
keywords: docker sandboxes, sbx, upstream proxy, corporate proxy, pac, socks5, system proxy, no_proxy, egress
weight: 50
---

An upstream proxy is the corporate or network proxy that Docker Sandboxes
forwards outbound traffic through on its way to the internet. This is separate
from the [network policy](governance/), which decides _which_ destinations are
allowed. The upstream proxy decides _how_ allowed traffic reaches them.

Docker Sandboxes sends two kinds of outbound traffic, and you can proxy them
independently:

- Sandbox traffic — network access from inside your sandboxes.
- Daemon traffic — the `sbx` daemon's own access: image pulls, telemetry,
  sign-in, and feature flags.

## Default behavior

By default, both kinds of traffic use your operating system's proxy settings,
including any PAC URL configured there. You don't need to configure anything. On
macOS and Windows, `sbx` tracks the OS proxy setting while it runs, so a change
to your network, VPN, or PAC configuration is picked up without a restart. If
your OS has no proxy configured, traffic goes direct.

## Set a proxy manually

Use `sbx settings set` to override the default for one or both kinds of traffic:

```console
$ sbx settings set proxy http://proxy.corp:3128          # both kinds of traffic
$ sbx settings set proxy.sandbox socks5://proxy.corp:1080 # sandbox traffic only
$ sbx settings set proxy.daemon direct                    # daemon traffic only
```

A proxy value can be any of the following:

| Value                                                                                  | Meaning                                                                                          |
| -------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------ |
| _(unset)_                                                                              | Fall back to the wider scope, then environment variables, then the OS system proxy (the default) |
| `http://host:port` or `https://host:port`                                              | An HTTP or HTTPS proxy                                                                           |
| `socks5://host:port` or `socks5h://host:port`                                          | A SOCKS5 proxy                                                                                   |
| `pac+http://host/proxy.pac`, `pac+https://host/proxy.pac`, or `file:///path/proxy.pac` | A PAC (proxy auto-config) file                                                                   |
| `system`                                                                               | Force the use of the OS system proxy                                                             |
| `direct`                                                                               | Force a direct connection with no proxy                                                          |

With `socks5://`, DNS is resolved locally before the connection is handed to the
proxy. With `socks5h://`, DNS resolution is delegated to the proxy.

### Exclude destinations from the proxy

Exclusion lists mirror the same scopes. Each takes a comma-separated list of
hosts, domain suffixes, IP addresses, or CIDR ranges, or `*` to bypass the
proxy entirely:

```console
$ sbx settings set no_proxy "*.internal.corp,10.0.0.0/8"    # both kinds of traffic
$ sbx settings set no_proxy.sandbox "*.svc.cluster.local"   # sandbox traffic only
$ sbx settings set no_proxy.daemon "registry.internal"      # daemon traffic only
```

## Environment variables

Because `sbx` runs from your shell, it also honors the standard and legacy proxy
environment variables, so existing setups keep working without migration:

- `HTTP_PROXY`, `HTTPS_PROXY`, and `NO_PROXY` (and their lowercase forms) — the
  standard variables. They apply to both kinds of traffic when no `proxy` or
  `no_proxy` setting is configured.
- `DOCKER_SANDBOXES_PROXY` and `DOCKER_SANDBOXES_NO_PROXY` — the environment
  form of `proxy.sandbox` and `no_proxy.sandbox`. They apply to sandbox traffic
  only and never affect daemon traffic.

The daemon reads these variables when it starts, so set them before your first
`sbx` command, or restart the daemon for a change to take effect.

## Precedence

For each kind of traffic, the first match wins:

1. The scope-specific setting (`proxy.sandbox` or `proxy.daemon`; sandbox traffic
   can also use the `DOCKER_SANDBOXES_PROXY` environment variable)
2. The `proxy` setting
3. `HTTP_PROXY` or `HTTPS_PROXY` from the shell
4. The OS system proxy (the default)
5. Direct

The matching exclusion list (`no_proxy.<scope>`, then `no_proxy`) applies to the
chosen proxy, and the standard `NO_PROXY` variable still applies on the
environment path.

> [!NOTE]
> The most specific value you set wins. Set nothing, and traffic uses the OS
> system proxy. A shell `HTTP_PROXY` counts as being set and takes precedence
> over the OS setting, the same way `curl` behaves.

## When changes take effect

The two kinds of traffic are resolved at different times:

- Sandbox scope (`proxy.sandbox`, `no_proxy.sandbox`, and the sandbox side of
  `proxy` and `no_proxy`) is re-resolved every time a sandbox is created or
  restarted. A change takes effect on the next sandbox you create or restart;
  already-running sandboxes keep the proxy they were created with.
- Daemon scope (`proxy.daemon`, `no_proxy.daemon`, the daemon side of `proxy` and
  `no_proxy`, and the `DOCKER_SANDBOXES_*` environment variables) is resolved
  once when the daemon starts. A change requires a daemon restart.

When a `system` or PAC proxy is in use, `sbx` still tracks OS-level proxy changes
(such as switching networks, connecting a VPN, or updated PAC contents) live.

## Authentication

Authenticating to the upstream proxy itself is supported only through explicit
credentials in the proxy URL, for example `http://user:pass@host:port`.
Integrated NTLM, Kerberos, and system single sign-on aren't yet supported.

## Related pages

- [Network isolation](security/isolation.md) — how traffic leaves a sandbox and
  the network policy it passes through
- [Troubleshooting: API calls fail with a certificate error](troubleshooting.md#api-calls-fail-with-a-certificate-error)
  — installing an internal root CA when your proxy inspects HTTPS traffic
