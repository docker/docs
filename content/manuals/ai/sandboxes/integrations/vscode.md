---
title: Connect VS Code to a sandbox
linkTitle: VS Code
weight: 10
description: Use VS Code Remote - SSH to develop inside a Docker Sandbox.
keywords: docker sandboxes, vs code, remote ssh, remote development, sbx
---

{{< summary-bar feature_name="Docker Sandboxes SSH" >}}

Use the Remote - SSH extension to open a VS Code window that runs inside a
sandbox. Your editor stays on your host while files, terminals, and extensions
run in the isolated sandbox.

## Prerequisites

- SSH access set up. See [Editor and app integrations](_index.md#enable-ssh-access).
- The [Remote - SSH](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-ssh)
  extension (`ms-vscode-remote.remote-ssh`) installed in VS Code.

## Connect

Confirm that you can connect to the sandbox from a terminal:

```console
$ ssh demo.sbx
```

In VS Code, open the Command Palette and run **Remote-SSH: Connect to Host...**.
Enter the sandbox hostname, such as `demo.sbx`, manually. After VS Code
connects, use the remote folder picker to
[select the mounted workspace](_index.md#select-the-workspace-folder).

For more connection options, see the VS Code instructions to
[connect to a remote host](https://code.visualstudio.com/docs/remote/ssh#_connect-to-a-remote-host).

## Notes

- The first connection installs the VS Code server inside the sandbox, so it
  can take a moment. Later connections are faster.

### Reconnect loop on macOS

Affected versions of VS Code can enter an infinite reconnect loop on macOS. If
this happens, set `remote.SSH.useLocalServer` to `false` in your VS Code user
settings:

```json
{
  "remote.SSH.useLocalServer": false
}
```

For details, see
[microsoft/vscode-remote-release#11672](https://github.com/microsoft/vscode-remote-release/issues/11672).

### SSH host key verification fails

VS Code can leave duplicate or malformed `Host *.sbx` blocks in your SSH
config after you add an SSH host. If a VS Code connection reports a
`KnownHostsCommand` error or `Host key verification failed`, remove every SSH
config block marked `docker sandboxes (managed)`, including the marker
comments:

```diff
-# >>> docker sandboxes (managed) >>>
-Host *.sbx
-    User _default_user_
-    ProxyCommand "sbx" ssh proxy %n
-    ...
-    UserKnownHostsFile "~/.ssh/sbx_known_hosts"
-    KnownHostsCommand "sbx" ssh known-hosts %H
-    StrictHostKeyChecking yes
-# <<< docker sandboxes (managed) <<<
```

Then regenerate the managed block:

```console
$ sbx setup ssh
```

Reconnect to the sandbox from VS Code.

## Related

- [Editor and app integrations](_index.md) — how SSH access works and how to
  set it up
