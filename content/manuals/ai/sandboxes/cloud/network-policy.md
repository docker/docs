---
title: Manage cloud network policy
linkTitle: Network policy
description: Control outbound connections from Docker cloud sandboxes with account-level and sandbox-level allow and deny network rules.
keywords: docker sandboxes, cloud network policy, sbx cloud, allow network, deny network
weight: 40
---

Cloud network policy controls outbound connections from cloud sandboxes. It is
a separate, network-only policy store with Docker account and individual
sandbox scopes.

> [!IMPORTANT]
>
> Docker Sandboxes organization policies aren't a central enforcement layer
> for cloud sandboxes. Organization and team network policy changes aren't
> continuously applied to cloud sandboxes. Configure and verify cloud rules with
> `sbx --cloud policy` rather than relying on organization governance.

## Understand policy scope

The cloud policy service doesn't provide organization or team scopes. An
account policy supplies the default for cloud sandboxes in the Docker account,
and a sandbox policy targets one cloud sandbox.

When a local sandbox daemon is available, `sbx --cloud create` reads active
network rules from its local policy store and copies them into the cloud
sandbox's initial policy. This point-in-time copy can include organization
rules that were synchronized to that machine. It doesn't make the cloud
sandbox organization-governed, and later organization policy changes don't
update the copied cloud rules.

This copy applies only to `sbx --cloud create`. A sandbox created by
`sbx --cloud run` doesn't read the local policy store. It uses cloud account
policy, network rules passed to the command, and network access declared by the
agent or kit.

Define the intended policy in the cloud store and use `sbx --cloud policy ls`
to verify it after creating the sandbox. Don't treat copied local rules as a
central enforcement boundary.

## Initialize account policy

Initialize the account policy with an allow-all or deny-all default:

```console
$ sbx --cloud policy init deny-all
```

The default applies to cloud sandboxes in the Docker account. You can add
rules after initialization or specify initial rules when creating a sandbox.

## Add network rules

Add an account-level exception:

```console
$ sbx --cloud policy allow network api.github.com:443
```

Scope a rule to one sandbox:

```console
$ sbx --cloud policy allow network api.anthropic.com:443 \
    --sandbox cloud-project
```

Deny rules take precedence when the same destination matches both an allow
rule and a deny rule.

You can also add initial rules while creating a sandbox:

```console
$ sbx --cloud create --name cloud-project \
    --allow-network api.github.com:443 \
    --deny-network example.com claude
```

## Inspect network policy

List account-level rules or the rules that apply to one sandbox:

```console
$ sbx --cloud policy ls
$ sbx --cloud policy ls cloud-project
```

Review connection decisions for a sandbox:

```console
$ sbx --cloud policy log cloud-project
```

These records describe cloud network policy decisions. They aren't Docker AI
Governance organization audit logs.

## Remove network rules

Use `sbx --cloud policy rm network` to remove an individual rule. Run the
command with `--help` to see the resource and scope options:

```console
$ sbx --cloud policy rm network --help
```

Reset the account's cloud network policy when you need to clear its default
and rules:

```console
$ sbx --cloud policy reset
```
