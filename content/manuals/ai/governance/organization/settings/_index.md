---
title: Sandboxes settings management
linkTitle: Sandboxes settings
weight: 35
description: How settings policies set defaults for Docker Sandboxes settings or lock them for your organization, teams, or members.
keywords: docker sandboxes, settings management, settings policy, sbx settings, AI governance, Docker Home, locked settings, organization settings
---

{{< summary-bar feature_name="AI Governance" >}}

Developers configure Docker Sandboxes on their machines with
[`sbx settings`](/manuals/ai/sandboxes/configuration/settings.md). Settings
policies let you, the administrator, manage some of those settings centrally
from Docker Home. You can give a setting a default value that developers can
still change, or lock it so developers can't change it. For example, you can
turn off experimental features, stop agents from reading images from the host
clipboard, or control how developers connect to sandboxes over SSH. For the
settings you can manage, see the
[Settings policy reference](../../reference/settings-policy.md). To create a
policy, see [Create and manage settings policies](create-manage.md).

You can apply a settings policy to your whole organization, to teams, or to
individual members.

## Setting states

A settings policy assigns a state to each setting it manages.

| State            | What developers get                                                                                                  |
| ---------------- | -------------------------------------------------------------------------------------------------------------------- |
| **Locked**       | The value or values you allow. Developers can't set anything outside them.                                           |
| **Default**      | Your value as their starting point. Developers can change it with `sbx settings set`.                                |
| **User defined** | Their own choice. The policy leaves the setting to developers.                                                       |
| **Not set**      | Whatever the next policy that applies sets. If no policy sets it, the developer's own value or the built-in default. |

What you can allow when you lock a setting depends on its type. For details,
see [Setting states by type](../../reference/settings-policy.md#setting-states-by-type).

## Policy scope and precedence

Each settings policy applies to one of the following:

- **All organization members**. Only one settings policy can apply to the
  whole organization at a time. To assign a different one, first unassign the
  current one.
- **Teams**. Teams you select.
- **Members**. Individual members you select.

When more than one policy applies to a developer, each setting is resolved
separately. For each setting, a member policy overrides a team policy, and a
team policy overrides the organization policy. If a developer is on more than
one team with a settings policy, and those policies set the same setting, the
policy highest in your precedence order wins for that setting. A policy that
leaves a setting as **Not set** doesn't override other policies for it, so
they can still supply that setting. To set the order, see
[Set the precedence of team policies](create-manage.md#set-the-precedence-of-team-policies).

## How policies combine with developer settings

On a developer's machine, a value set by your policy as **Default** replaces
the built-in default. A value the developer sets with `sbx settings set`, or
through the setting's environment variable if it has one, still takes
precedence over it.

A **Locked** setting applies on top of everything else. Developers can't use
`sbx settings set` or an environment variable to get around it.

## Developer experience

Developers keep using the same `sbx settings` commands. In `sbx settings ls`,
a setting your policy manages shows `remote` in the `SOURCE` column, marked
`(managed)` or `(locked)`:

```console
$ sbx settings ls
KEY                                VALUE   TYPE   SOURCE           RESTART   DESCRIPTION
platform.allowExperimentalFeatures false   bool   remote (locked)  no        Allow experimental features.
```

`sbx settings get <KEY> --json` includes a `management` object with the
policy state and whether the setting is locked.

If a developer tries to set a value that a locked setting doesn't allow, the
command fails:

```console
$ sbx settings set platform.allowExperimentalFeatures true
error: request failed: 403 Forbidden: value is restricted by administrator settings
```
