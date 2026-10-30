---
title: Develop remotely in a cloud sandbox
linkTitle: Remote development
description: Connect VS Code or Cursor to a cloud sandbox over SSH, edit and preview a sample, and save your work.
keywords: docker agentic platform, cloud sandboxes, workflows, remote development
weight: 20
---

Keep project files and development tools in a cloud sandbox while using VS Code
or Cursor on your computer. Provisional duration: 25 minutes; this estimate needs
measurement during review.

## Before you begin

You need an active [Docker Agentic Platform subscription](/manuals/agentic-platform/signup.md#activate-cloud-access),
the [Docker Sandboxes CLI](/manuals/ai/sandboxes/install.md) version 0.45.0 or
later, an SSH client, and VS Code or Cursor. The public sample needs no GitHub
token or model API key. Use the same Docker account in the CLI and Console.

Your editor window runs locally. Project files, terminal commands, and the
editor's remote server run inside the sandbox. Files don't sync to your computer.
Docker bills compute while the sandbox runs, even with the editor disconnected.

## 1. Create the remote project

1. Open the [Console](https://agentic-platform.docker.com/) and select **+ New sandbox**.
2. Choose **shell**, deselect saved credentials and MCP tools, and select
   **Open** so the editor can download its remote server and dependencies.
   Open permits unrestricted outbound network access; use it here with the
   public sample and no saved credentials.
3. Review the compute price and timer. Select **Stop** under **When the time is up**
   to retain files, then select **Run**.
4. In the sandbox's browser terminal, clone the sample:

   ```console
   $ git clone https://github.com/docker/welcome-to-docker.git project
   ```

Expect `/home/agent/workspace/project` to contain the sample source and Dockerfile.
Local workspaces aren't mounted in the cloud sandbox.

## 2. Configure and test SSH

In a terminal on your computer, run:

```console
$ sbx login
$ sbx --cloud setup ssh
$ sbx --cloud ls
```

Copy the sandbox's `sbx_`-prefixed ID from the list, then test the connection:

```console
$ ssh <SANDBOX_ID>@sbx_cloud
$ ls /home/agent/workspace/project
$ exit
```

Expect a shell inside the sandbox and the project files. `exit` returns you to
your computer. Use the ID, rather than the sandbox's display name, for SSH.
See [Connect with SSH](/manuals/ai/sandboxes/cloud/usage.md#connect-with-ssh)
for proxy troubleshooting.

## 3. Connect your editor

Choose the instructions for your editor:

- In VS Code, install the
  [Remote - SSH extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-ssh).
  Open the Command Palette and run **Remote-SSH: Connect to Host**.
- In Cursor, install its Remote - SSH support. Open the Command Palette and run
  **Remote-SSH: Connect to Host**.

Enter `<SANDBOX_ID>@sbx_cloud` using the ID you tested. In the remote window,
open `/home/agent/workspace/project`. The first connection installs the editor's
server inside the sandbox. Expect the editor to show a remote connection and
the same project files you listed over SSH.

For local sandbox integration, see the separate
[VS Code](/manuals/ai/sandboxes/integrations/vscode.md) and
[Cursor](/manuals/ai/sandboxes/integrations/cursor.md) pages.

## 4. Edit and preview remotely

In the remote editor, find a visible heading in the application's source,
change its text, and save the file. In the editor's remote terminal, build and run:

```console
$ cd /home/agent/workspace/project
$ docker build -t welcome-to-docker .
$ docker run -d --name welcome-to-docker -p 3000:3000 welcome-to-docker
```

In the Console, select **Connect**, enter `3000` under **Connect via a Port**,
and select **Open Port**. Open the generated HTTPS URL and expect to see your
changed heading. This URL is public: Docker account authentication doesn't
protect it. Use non-sensitive data and application authentication for private
content. Close the port in **Connect** after checking the preview and reload
the URL to verify the application is no longer reachable.

## 5. Reconnect and keep your work

Close the remote editor window, then connect again using the same host entry
and open the same directory. Expect your saved change to remain. Disconnecting
doesn't pause the sandbox. If the timer stopped it, resume it in the Console
before reconnecting; check the fresh timer and restart the workload if needed.

From a terminal on your computer, copy the project out:

```console
$ sbx --cloud cp <SANDBOX_ID>:/home/agent/workspace/project ./cloud-project
```

Inspect the copied source and confirm that your edit is present. Alternatively,
commit and push to a repository you control. Configure a
[GitHub credential](/manuals/agentic-platform/secrets.md#github-credential)
with the required repository permissions before launching a sandbox for that workflow.

## 6. Clean up

In **Sandboxes**, pause the sandbox and confirm its paused status to stop compute
charges and retain files. A resumed sandbox gets a fresh timer with the original
duration. Once you've verified your saved copy, delete the sandbox if you no
longer need it. Deletion removes the remote files and editor server.
See [Lifecycle actions](/manuals/agentic-platform/sandboxes.md#manage-the-lifecycle).
