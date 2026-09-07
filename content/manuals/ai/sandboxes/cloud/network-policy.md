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
> Organization governance is not available for cloud sandboxes in this
> release. Configure cloud rules with `sbx --cloud policy` and verify network
> access using connection checks and policy logs.

## Understand policy scope

The cloud CLI supports account and sandbox policy scopes. Your account policy
supplies the default for cloud sandboxes you create in the Docker account. A
sandbox policy adds rules for one cloud sandbox. Matching deny rules take
precedence over allow rules across the applicable policies.

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

Define the intended policy in the cloud store. After creation, inspect the
configured rules and [verify connection decisions](#inspect-network-policy).
Don't treat copied local rules as a central enforcement boundary.

## Initialize account policy

Initialize the account policy with an allow-all or deny-all default:

```console
$ sbx --cloud policy init deny-all
```

The default applies to your cloud sandboxes in the Docker account. You can add
rules after initialization or specify initial rules when creating a sandbox.
A deny-all default still permits destinations allowed by applicable sandbox
or agent-kit rules.

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

Inspect your account policy or a sandbox's configured policy:

```console
$ sbx --cloud policy ls
$ sbx --cloud policy ls cloud-project
```

The sandbox view shows its policy document, or your account default when the
sandbox has no policy document. It does not show the complete combination of
applicable rules.

To verify enforcement, attempt the connection from the sandbox, then review
the connection decisions:

```console
$ sbx --cloud exec cloud-project curl -I https://api.github.com
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

Reset your account policy to deny-all and remove its rules:

```console
$ sbx --cloud policy reset
```

Sandbox-specific policies remain in place and can still grant access. Inspect
those policies separately. In `sbx` version 0.42.0, the reset command prints
"reset to allow-all", but the service resets the account policy to deny-all.
