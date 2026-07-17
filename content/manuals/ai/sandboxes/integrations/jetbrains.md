---
title: Connect a JetBrains IDE to a sandbox
linkTitle: JetBrains IDEs
weight: 25
description: Use JetBrains Remote Development to develop inside a Docker Sandbox over SSH.
keywords: docker sandboxes, jetbrains, remote development, remote ssh, gateway, sbx
---

{{< summary-bar feature_name="Docker Sandboxes SSH" >}}

JetBrains Remote Development runs the IDE backend inside the sandbox and opens
the project locally in JetBrains Client. Connect through JetBrains Gateway or
the Remote Development option in a supported JetBrains IDE.

## Prerequisites

- SSH access set up. See [Editor and app integrations](_index.md#enable-ssh-access).
- [JetBrains Gateway installed](https://www.jetbrains.com/help/idea/jetbrains-gateway.html),
  or a supported JetBrains IDE with the Remote Development Gateway plugin.

## Allow JetBrains network access

JetBrains Gateway downloads the IDE backend into the sandbox. The Balanced
network preset doesn't include all the required JetBrains endpoints. Add a
sandbox-scoped rule for them:

```console
$ sbx policy allow network --sandbox demo "*.jetbrains.com,data.services.jetbrains.com"
```

If your organization manages sandbox network policy, ask your administrator to
allow these endpoints instead. Organization policy overrides local rules.

## Connect

Confirm that you can connect to the sandbox from a terminal:

```console
$ ssh demo.sbx
```

1. Open JetBrains Gateway. Alternatively, select **Remote Development** from
   the welcome screen of a supported JetBrains IDE.
2. Under **SSH Connection**, select **New Connection**.
3. Create an SSH configuration with `demo.sbx` as the host and select
   **OpenSSH config and authentication agent** as the authentication type. The
   managed SSH config supplies the remaining connection settings.
4. Select **Check Connection and Continue**.
5. Choose the backend IDE version and the project folder in the sandbox, then
   connect. The first connection downloads and installs the IDE backend inside
   the sandbox.

For more connection options, see the JetBrains instructions to
[connect and work with JetBrains Gateway](https://www.jetbrains.com/help/idea/remote-development-a.html).

## Related

- [Editor and app integrations](_index.md) — how SSH access works and how to
  set it up
- [Local policy](../governance/local.md) — manage sandbox network access
