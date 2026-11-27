---
title: Settings policy reference
linkTitle: Settings policy
weight: 37
description: Reference for the Docker Sandboxes settings that settings policies can manage, and the setting states available for each type.
keywords: docker sandboxes, settings policy, sbx settings, managed settings, locked settings, AI Governance
---

This reference lists the Docker Sandboxes settings that
[settings policies](../organization/settings/_index.md) can manage, and the
setting states Docker Home offers for each type.

## Settings you can manage

| Setting in Docker Home             | `sbx settings` key                   | Type    |
| ---------------------------------- | ------------------------------------ | ------- |
| Clipboard image paste              | `clipboard.imagePaste`               | Boolean |
| Fail-closed credentials            | `credentials.failClosed`             | Boolean |
| Allow experimental features        | `platform.allowExperimentalFeatures` | Boolean |
| SSH accepted environment variables | `ssh.acceptEnv`                      | List    |
| SSH authorized keys                | `ssh.authorizedKeys`                 | List    |
| SSH auto-create sandbox            | `ssh.autoCreate`                     | Boolean |
| SSH default agent                  | `ssh.defaultAgent`                   | String  |
| SSH default template               | `ssh.defaultTemplate`                | String  |
| SSH managed key                    | `ssh.manageKey`                      | Boolean |
| SSH port                           | `ssh.port`                           | Integer |
| SSH workspace root                 | `ssh.workspaceRoot`                  | String  |

For what each setting does, see the
[settings reference](/manuals/ai/sandboxes/configuration/settings.md#settings-reference).
Settings that aren't in this table stay under each developer's control.

## Setting states by type

For a boolean setting, you pick the state and value together from one menu.
For what each state means, see
[Setting states](../organization/settings/_index.md#setting-states).

| Option                | State        | Value |
| --------------------- | ------------ | ----- |
| **Enabled (locked)**  | Locked       | On    |
| **Enabled**           | Default      | On    |
| **User defined**      | User defined | None  |
| **Disabled**          | Default      | Off   |
| **Disabled (locked)** | Locked       | Off   |
| **Not set**           | Not set      | None  |

For integer, string, and list settings, you pick **Locked**, **Default**,
**User defined**, or **Not set**. **Locked** and **Default** then ask for
values:

| Type    | **Locked** asks for                                                          | **Default** asks for    |
| ------- | ---------------------------------------------------------------------------- | ----------------------- |
| Integer | A minimum, a maximum, or both                                                | A number                |
| String  | The values developers can choose from. If you allow one, that value applies. | A value                 |
| List    | The entries developers can include in their list                             | The entries in the list |

For example, locking **SSH accepted environment variables** to `LANG` and
`TERM` stops developers from forwarding other environment variables into
sandboxes over SSH.
