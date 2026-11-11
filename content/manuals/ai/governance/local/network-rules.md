---
title: Manage local network rules
linkTitle: Manage network rules
weight: 20
description: Allow, deny, and remove local network rules for sandboxes on your machine with sbx policy.
keywords: docker sandboxes, local policy, network rules, sbx policy allow, sbx policy deny, HTTP rules, UDP
---

Local network rules allow or block destinations on top of your
[preset](_index.md#start-from-a-preset). This page shows common tasks. For
every flag, see the [`sbx policy`](/reference/cli/sbx/policy/) CLI reference,
and for the patterns a rule accepts, see
[Network and HTTP rule syntax](../reference/network-rules.md).

## Set or reset the preset

The CLI asks you to choose a preset the first time you run a sandbox in an
interactive terminal. To set one without the prompt, or to start over, use
`sbx policy init` or `sbx policy reset`.

### Set the preset without the prompt

In CI pipelines and other non-interactive environments, commands such as
`sbx run` fail until a preset is set. Set it first with `sbx policy init`:

```console
$ sbx policy init balanced
```

`sbx policy init` only works before a preset is set. To change the preset after
that, reset local policy.

### Reset local policy

To remove all custom rules and start fresh with a new preset, use
`sbx policy reset`:

```console
$ sbx policy reset
```

This deletes the local policy store, restarts the daemon, and prompts you to
choose a new preset. Running sandboxes stop when the daemon shuts down. Pass
`--force` to skip the confirmation prompt.

> [!WARNING]
> `sbx policy reset` deletes all locally configured policy rules, including any
> destinations you approved. Those destinations are requested again the next
> time a sandbox reaches them. To pick up organization policy changes without
> deleting local rules, use `sbx policy refresh` instead. For more details, see
> [Organization policy changes aren't taking effect](../troubleshooting.md#organization-policy-changes-arent-taking-effect).

## Allow or block a destination

Use [`sbx policy allow`](/reference/cli/sbx/policy/allow/) and
[`sbx policy deny`](/reference/cli/sbx/policy/deny/) to add rules. Rules apply
to all sandboxes unless you pass `--sandbox`:

```console
$ sbx policy allow network "api.anthropic.com,*.npmjs.org"
$ sbx policy deny network --sandbox my-sandbox ads.example.com
```

To add a deny rule when you create a sandbox, pass `--deny-network` to
`sbx create` or `sbx run`:

```console
$ sbx run --deny-network ads.example.com claude
```

Rules scoped to a sandbox appear in `sbx policy ls <name>`, and removing the
sandbox removes them.

## Block a destination that a kit allows

A deny rule takes precedence over a kit's allow rule. To block a destination a
kit allows, add a deny rule for the sandbox that uses the kit:

```console
$ sbx policy deny network --sandbox my-sandbox openrouter.ai
```

## Let an agent read from an API without writing to it

An HTTP rule matches specific methods and paths on a host. To let an agent read
from an API but not write to it, allow only read methods on the paths it needs:

```console
$ sbx policy allow network api.github.com --method GET,HEAD --path '/repos/**'
```

Write the host without a scheme, such as `api.github.com` rather than
`https://api.github.com`. Quote the path so your shell doesn't expand the
wildcard. A rule takes one path, so add a separate rule for each path. To also
require query parameters, add `--query`, such as `--query 'team=platform'`.

For how methods, paths, and query parameters match, see
[HTTP rules](../reference/network-rules.md#http-rules).

## Block one endpoint on an allowed host

To keep a host available but block part of it, add an HTTP deny rule on top of
an allow:

```console
$ sbx policy allow network api.example.com
$ sbx policy deny network api.example.com --method POST --path '/admin/**'
```

`POST` requests under `/admin` are blocked. Every other request to
`api.example.com` is allowed.

## Remove a rule

Remove a rule by its resource, or by its ID from `sbx policy ls --wide`:

```console
$ sbx policy rm network --resource ads.example.com
$ sbx policy rm network --id 2d3c1f0e-4a73-4e05-bc9d-f2f9a4b50d67
```

Pass `--sandbox <name>` for a rule scoped to a sandbox. To remove an HTTP rule
by resource, also pass the same `--method`, `--path`, and `--query` flags you
added it with. For more details, see
[Inspect policy status and rules](inspect.md).

## Allow outbound UDP

Outbound UDP is experimental and disabled by default. Turn on experimental
features and UDP egress, then add a UDP rule. The CLI warns you if you save a
UDP rule while UDP egress is off:

```console
$ sbx settings set platform.allowExperimentalFeatures true
$ sbx settings set feature.udp-egress true
$ sbx policy allow network --protocol udp api.example.com:443
```

Allow rules apply to TCP unless you pass `--protocol`, while deny rules apply
to both TCP and UDP. UDP rules follow the same organization and local
precedence as TCP rules.

Some UDP traffic is blocked whatever your rules allow:

- Traffic to a destination that needs an HTTP, SOCKS5, system, or PAC-selected
  proxy, because those proxies can't carry UDP.
- Traffic to local-network destinations. This block applies before any policy
  is evaluated, so even a rule that allows all hosts doesn't permit it. It
  covers:
  - Multicast, link-local, and unspecified addresses
  - `255.255.255.255` and the broadcast addresses of the host's networks

  Discovery protocols such as mDNS and SSDP don't work from a sandbox. UDP to
  `host.docker.internal` isn't affected. The network log records these
  connections under the rule `<udp local network boundary>`.

To list UDP rules or check a destination:

```console
$ sbx policy ls --protocol udp
$ sbx policy check network --protocol udp api.example.com:443
```
