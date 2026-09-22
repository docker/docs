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

Cloud creation uses cloud account policy, network rules passed to the command,
and network access declared by the agent or kit. Both `sbx --cloud create` and
`sbx --cloud run` leave local network and organization policies on the host.
Moving a local sandbox to the cloud also uses cloud policy.

Define the intended policy in the cloud store. After creation, inspect the
configured rules and [verify connection decisions](#inspect-network-policy).

## Initialize account policy

Set the account policy to `allow-all`, `balanced`, or `deny-all`:

```console
$ sbx --cloud policy init deny-all
```

The default applies to your cloud sandboxes in the Docker account. You can add
rules after initialization or specify initial rules when creating a sandbox.
A deny-all default still permits destinations allowed by applicable sandbox
or agent-kit rules.

The `balanced` preset sets deny-all as the default and adds allow rules for
common development services. To set a default for one sandbox, use
`sbx --cloud policy init deny-all --sandbox cloud-project`.

You can run `init` again to change the default. Existing allow and deny rules
remain in place. Use `reset` to remove account rules before choosing another
preset.

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

Cloud rules match network destinations. HTTP method and path restrictions,
`--protocol`, and local governance profiles aren't supported.

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

Remove a rule by its pattern, from whichever allow or deny list contains it:

```console
$ sbx --cloud policy rm network --resource api.github.com:443
$ sbx --cloud policy rm network --sandbox cloud-project --resource api.anthropic.com:443
```

Remove the custom account policy and return to the platform default:

```console
$ sbx --cloud policy reset
```

The command asks for confirmation and prints the default stored by the server
after the reset. Use `--force` in scripts. The local daemon is unaffected.
Sandbox-specific policies remain in place and can still grant access. Inspect
those policies separately.
