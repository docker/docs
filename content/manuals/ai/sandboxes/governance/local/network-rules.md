---
title: Manage local network rules
linkTitle: Manage network rules
weight: 20
description: Allow, deny, and remove local network rules for sandboxes on your machine with sbx policy.
keywords: docker sandboxes, local policy, network rules, sbx policy allow, sbx policy deny, HTTP rules, UDP
---

A rule covers a destination host. It can also name HTTP methods and paths to
narrow the match to part of that host.

## Network rules

Use [`sbx policy allow`](/reference/cli/sbx/policy/allow/) and
[`sbx policy deny`](/reference/cli/sbx/policy/deny/) to add or restrict access
on top of the active preset. Changes take effect immediately. Rules apply to
all sandboxes by default:

```console
$ sbx policy allow network api.anthropic.com
$ sbx policy deny network ads.example.com
```

Pass `--sandbox <name>` to scope a rule to one sandbox:

```console
$ sbx policy allow network --sandbox my-sandbox api.example.com
$ sbx policy deny network --sandbox my-sandbox ads.example.com
```

As of v0.38.0, you can also set per-sandbox deny rules at creation time with
`--deny-network` on `sbx create` or `sbx run`, instead of adding them after the
fact:

```console
$ sbx create --deny-network ads.example.com claude .
$ sbx run --deny-network ads.example.com claude
```

Pass the flag multiple times to deny more than one host. Rules added this way
appear in `sbx policy ls <name>` and can be removed with
`sbx policy rm network --sandbox <name> --resource <host>`.

Specify multiple hosts in one command with a comma-separated list:

```console
$ sbx policy allow network "api.anthropic.com,*.npmjs.org,*.pypi.org"
```

Remove a rule by resource or by rule ID:

```console
$ sbx policy rm network --resource ads.example.com
$ sbx policy rm network --id 2d3c1f0e-4a73-4e05-bc9d-f2f9a4b50d67
```

To remove a sandbox-scoped rule, pass `--sandbox <name>`:

```console
$ sbx policy rm network --sandbox my-sandbox --resource api.example.com
```

## HTTP method and path rules

Add `--method` to an allow or deny rule to match specific HTTP methods on a
host, and `--path` to restrict it to part of the host's URL space:

```console
$ sbx policy allow network api.github.com --method GET --path '/repos/org/project/**'
```

Quote the path so your shell doesn't expand the wildcard. Pass several methods
as a comma-separated list:

```console
$ sbx policy allow network api.github.com --method GET,HEAD
```

`--method ANY` matches every HTTP method, and `--path` defaults to `/**` when
you omit it:

```console
$ sbx policy allow network api.github.com --method ANY
```

`ANY` can't be combined with specific methods, and a path without a method is
rejected. Pass a method, or use `ANY` when you mean every method.

Method names are case-insensitive. The accepted values are `GET`, `HEAD`,
`POST`, `PUT`, `PATCH`, `DELETE`, `OPTIONS`, `CONNECT`, and `TRACE`.

A path must start with `/` and be canonical. It can't contain a query string, a
fragment, percent-encoding, control characters, surrounding whitespace,
repeated or trailing slashes, or dot segments such as `.` and `..`. Each rule
takes one path. Repeating `--path` keeps only the last value, and a comma is
read as part of the path, so add a separate rule for each path:

```console
$ sbx policy allow network api.github.com --method GET --path '/repos/**'
$ sbx policy allow network api.github.com --method GET --path '/users/**'
```

Hosts follow the same patterns as network rules and can include a port. Write
the host on its own, without a scheme, so an HTTP rule takes `api.example.com`
rather than `https://api.example.com`.

A local HTTP rule takes a hostname. To match an IP address or a CIDR range,
add a plain network rule for that destination instead.

Deny rules take the same flags, which is the usual way to carve a method or
path out of a broader allow:

```console
$ sbx policy allow network api.example.com
$ sbx policy deny network api.example.com --method POST --path '/admin/**'
```

For how the two layers combine, see
[HTTP rules](../reference/network-rules.md#http-method-and-path).

Remove an HTTP rule by naming the same qualifiers you added it with, or by
rule ID:

```console
$ sbx policy rm network --resource api.github.com --method GET --path '/repos/org/project/**'
$ sbx policy rm network --id 7f3a1c2e-4a73-4e05-bc9d-f2f9a4b50d67
```

List HTTP rules with `--type http`, or see them alongside network rules in a
wide listing, where the `METHOD` and `PATH` columns are empty for rules that
match a whole host:

```console
$ sbx policy ls --wide
TYPE      METHOD   PATH
network   -        -
http      GET      /repos/org/project/**
```

> [!NOTE]
> `sbx policy check network` and `sbx policy log` don't evaluate or display
> HTTP methods and paths. A check reports the decision for the host, which can
> differ from the decision for a specific method and path on that host.

## Allow outbound UDP

Outbound UDP is experimental and disabled by default. Turn on experimental
features and UDP egress before adding UDP allow rules:

```console
$ sbx settings set platform.allowExperimentalFeatures true
$ sbx settings set feature.udp-egress true
$ sbx policy allow network --protocol udp api.example.com:443
```

Local allow rules apply to TCP by default. Use `--protocol udp` for UDP or
`--protocol tcp,udp` for both. Deny rules apply to both protocols by default.
Use `--protocol` to restrict a deny rule to one protocol.

UDP follows the same organization and local policy precedence as TCP. It is
refused when the destination requires an HTTP, SOCKS5, system, or PAC-selected
proxy, because those proxies can't carry UDP. ICMP remains blocked.

UDP to local-network destinations is blocked before any policy is evaluated,
so no allow rule can permit it, including one that allows all hosts. This
covers multicast, link-local, and unspecified addresses, `255.255.255.255`,
and the broadcast addresses of the host's networks. As a result, discovery
protocols that use multicast, such as mDNS and SSDP, don't work from a
sandbox. UDP to `host.docker.internal` isn't affected. The network log records
these connections under the rule `<udp local network boundary>`.

Inspect UDP rules or check a destination:

```console
$ sbx policy ls --protocol udp
$ sbx policy check network --protocol udp api.example.com:443
```

The CLI warns if you save a UDP rule while UDP egress is disabled.

## Troubleshooting

### An HTTP method or path is blocked on an allowed host

A host that a network rule allows can still have individual methods or paths
denied by an HTTP rule. Run `sbx policy ls --type http` to see which HTTP rules
apply. `sbx policy check network` reports the decision for the host only, so it
shows a host as allowed even when the specific request is denied. See
[HTTP method and path rules](#http-method-and-path-rules).
