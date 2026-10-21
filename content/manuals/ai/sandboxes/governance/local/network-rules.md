---
title: Manage local network rules
linkTitle: Manage network rules
weight: 20
description: Allow, deny, and remove local network rules for sandboxes on your machine with sbx policy.
keywords: docker sandboxes, local policy, network rules, sbx policy allow, sbx policy deny, HTTP rules, UDP
---

Local network rules allow or block destinations on top of your
[preset](_index.md#start-from-a-preset). You can add rules for all sandboxes or
for one sandbox, narrow a rule to HTTP methods and paths, and remove rules you
no longer need. For the host patterns a rule accepts, see
[Network and HTTP rule syntax](../reference/network-rules.md#destinations).

## Allow or block a destination

Use [`sbx policy allow`](/reference/cli/sbx/policy/allow/) and
[`sbx policy deny`](/reference/cli/sbx/policy/deny/) to add or restrict access
on top of the active preset. Changes take effect immediately. Rules apply to
all sandboxes by default:

```console
$ sbx policy allow network api.anthropic.com
$ sbx policy deny network ads.example.com
```

Specify multiple hosts in one command with a comma-separated list:

```console
$ sbx policy allow network "api.anthropic.com,*.npmjs.org,*.pypi.org"
```

## Allow or block a destination for one sandbox

Pass `--sandbox <name>` to scope a rule to one sandbox:

```console
$ sbx policy allow network --sandbox my-sandbox api.example.com
$ sbx policy deny network --sandbox my-sandbox ads.example.com
```

You can also add deny rules for a sandbox when you create it, with
`--deny-network` on `sbx create` or `sbx run`:

```console
$ sbx create --deny-network ads.example.com claude .
$ sbx run --deny-network ads.example.com claude
```

Pass the flag multiple times to deny more than one host. Rules added this way
appear in `sbx policy ls <name>`, like any other rule scoped to the sandbox.

To block a destination that a kit allows, add a deny rule for the sandbox that
uses the kit. A deny rule takes precedence over the kit's allow rule:

```console
$ sbx policy deny network --sandbox my-sandbox openrouter.ai
```

Removing a sandbox also removes the rules scoped to it.

## Allow or block HTTP methods and paths

The rules in the previous sections match a whole host. An HTTP rule matches
specific HTTP methods and paths on a host instead, so you can allow or block
specific API requests. For example, you can let an agent read from an API
without letting it write to it.

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
use a network rule without `--method` or `--path`.

Deny rules take the same flags, which is the usual way to carve a method or
path out of a broader allow:

```console
$ sbx policy allow network api.example.com
$ sbx policy deny network api.example.com --method POST --path '/admin/**'
```

For more details on how the two layers combine, see
[HTTP method and path](../reference/network-rules.md#http-method-and-path).

To list only HTTP rules, run `sbx policy ls --type http`. In
`sbx policy ls --wide`, the `METHOD` and `PATH` columns show `-` for rules that
match a whole host.

> [!NOTE]
> `sbx policy check network` and `sbx policy log` don't evaluate or display
> HTTP methods and paths. A check reports the decision for the host, which can
> differ from the decision for a specific method and path on that host.

## Remove a rule

Remove a rule by resource or by rule ID:

```console
$ sbx policy rm network --resource ads.example.com
$ sbx policy rm network --id 2d3c1f0e-4a73-4e05-bc9d-f2f9a4b50d67
```

To remove a rule scoped to a sandbox, pass `--sandbox <name>`:

```console
$ sbx policy rm network --sandbox my-sandbox --resource api.example.com
```

To remove an HTTP rule, name the same method and path you added it with, or
use its rule ID:

```console
$ sbx policy rm network --resource api.github.com --method GET --path '/repos/org/project/**'
```

To find a rule's ID, run `sbx policy ls --wide`. For more details, see
[Inspect policy status and rules](inspect.md).

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

## Set or reset the preset

You choose a preset the first time you run a sandbox. For more details, see
[Start from a preset](_index.md#start-from-a-preset). To set the preset without
the prompt, or to start over with a new one, use `sbx policy init` or `sbx
policy reset`.

### Set the preset without the prompt

In non-interactive environments such as CI pipelines or headless servers, the
interactive prompt can't be displayed. Use `sbx policy init` to set the
preset before running any other `sbx` commands:

```console
$ sbx policy init balanced
```

Available values are `allow-all`, `balanced`, and `deny-all`. `sbx policy init`
only works before a preset is set. To change the preset after that, reset local
policy.

### Reset local policy

To remove all custom rules and start fresh with a new preset, use
`sbx policy reset`:

```console
$ sbx policy reset
```

This deletes the local policy store, restarts the daemon, and prompts you to
choose a new preset. Running sandboxes stop when the daemon shuts down. Pass
`--force` to skip the confirmation prompt:

```console
$ sbx policy reset --force
```

> [!WARNING]
> `sbx policy reset` deletes all locally configured policy rules, including any
> destinations you approved. Those destinations are requested again the next
> time a sandbox reaches them. To pick up organization policy changes without
> deleting local rules, use `sbx policy refresh` instead. For more details, see
> [Organization policy changes aren't taking effect](../../troubleshooting.md#organization-policy-changes-arent-taking-effect).
