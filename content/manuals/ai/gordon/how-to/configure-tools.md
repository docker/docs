---
title: Configure Gordon's tools
linkTitle: Configure tools
description: Enable and disable Gordon's built-in tools based on your needs
weight: 40
---

{{< summary-bar feature_name="Gordon" >}}

Gordon includes built-in tools that extend its capabilities. You can configure
which tools Gordon has access to based on your security requirements and
workflow needs.

Tool configuration provides an additional layer of control:

- Enabled tools: Gordon can propose actions using these tools (subject to
  your approval)
- Disabled tools: Gordon cannot use these tools, and will not request
  permission to use them

## Accessing tool settings

To configure Gordon's tools:

1. Open Docker Desktop.
2. Select **Ask Gordon** in the sidebar.
3. Select the settings icon at the bottom of the text input area.

   ![Session settings icon](../images/perm_settings.avif?border=true)

The tool settings dialog opens with two tabs: **Basic** and **Advanced**.

## Basic tool settings

In the **Basic** tab, you can enable or disable individual tools globally.

To disable a tool:

1. Find the tool you want to disable in the list.
2. Toggle it off.
3. Select **Save**.

Disabled tools cannot be used by Gordon, even with your approval.

## Advanced tool settings

The **Advanced** tab lets you create fine-grained allow-lists and deny-lists
for specific commands or patterns.

Allow-lists:
Gordon can use allow-listed commands even when the main tool is disabled. For
example, disable the shell tool but allow `cat`, `grep`, and `ls`.

Deny-lists:
Block specific commands while keeping the tool enabled. For example, allow the
shell tool but deny `chown` and `chmod`.

To configure:

1. Switch to the **Advanced** tab.
2. Add commands to **Allow rules** or **Deny rules**.
3. Select **Save**.

![Advanced tool configuration](../images/gordon_advanced_tool_config.avif?w=500px&border=true)

Gordon still requests approval before running allow-listed tools, unless YOLO
mode (auto-approve mode that bypasses permission checks) is enabled.

## Organizational controls

For Business subscriptions, administrators can control tool access for the
entire organization using Settings Management.

Administrators can:

- Disable specific tools for all users
- Lock tool configuration to prevent users from changing it
- Set organization-wide tool policies

See [Settings Management](/enterprise/security/hardened-desktop/settings-management/)
for details.
