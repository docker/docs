---
title: Quick install Docker Sandboxes
linkTitle: Docker Sandboxes
description: Install the sbx CLI, sign in, and prepare to run an AI coding agent in a Docker Sandbox.
keywords: Docker Sandboxes, sbx, quick install, AI agents
---

Install the `sbx` CLI to run coding agents in local or cloud sandboxes.
Docker Desktop and Docker Engine are not required.

## Install the CLI

{{< tabs >}}
{{< tab name="macOS" >}}

Requires macOS Sonoma 14 or later on Apple silicon, with Homebrew installed.

```console
$ brew trust docker/tap
$ brew install docker/tap/sbx
```

{{< /tab >}}
{{< tab name="Windows" >}}

Requires Windows 11 on a 64-bit Intel or AMD processor. Install for your user
account with Windows Package Manager:

```powershell
winget install -h Docker.sbx
```

Open a terminal after installation so it picks up the updated `PATH`.

{{< /tab >}}
{{< tab name="Ubuntu" >}}

Requires Ubuntu 24.04 or later on a 64-bit Intel, AMD, or Arm processor.
Add Docker's package repository, then install `sbx`:

```console
$ curl -fsSL https://get.docker.com | sudo REPO_ONLY=1 sh
$ sudo apt install docker-sbx
```

{{< /tab >}}
{{< /tabs >}}

Local sandboxes also require hardware virtualization. Before starting a local
sandbox, follow the [platform prerequisites](/manuals/ai/sandboxes/install.md#prerequisites)
to set up Windows Hypervisor Platform or KVM where needed. Cloud sandboxes do
not require local virtualization.

For manual downloads, machine-wide deployment, and troubleshooting, see the
[full Sandboxes installation instructions](/manuals/ai/sandboxes/install.md).

## Check the installation and sign in

```console
$ sbx version
$ sbx login
```

The login command opens your browser to sign in to Docker.
