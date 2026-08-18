---
title: Install Docker Sandboxes
linkTitle: Install
weight: 5
description: Install the sbx CLI on macOS, Windows, or Linux and sign in to Docker Sandboxes.
keywords: sandbox, sbx, install, macOS, Windows, Linux, Ubuntu
---

Install the `sbx` CLI to run AI coding agents in isolated microVMs. You don't
need Docker Desktop or Docker Engine to use `sbx`.

## Prerequisites

### macOS

- macOS Sonoma version 14 or later
- Apple silicon

### Windows

- Windows 11
- A 64-bit Intel or AMD processor
- Windows Hypervisor Platform

To turn on Windows Hypervisor Platform, open an elevated PowerShell prompt and
run:

```powershell
Enable-WindowsOptionalFeature -Online -FeatureName HypervisorPlatform -All
```

### Linux

- Ubuntu 24.04 or later
- A 64-bit Intel or AMD processor, or a 64-bit Arm processor
- KVM hardware virtualization supported and turned on by the CPU
- Your user account in the `kvm` group

If you're running inside a virtual machine or virtual desktop infrastructure
environment, the environment must support nested virtualization.

Verify that KVM is available:

```console
$ lsmod | grep kvm
```

A working setup shows `kvm_intel`, `kvm_amd`, `kvm_arm64`, or `kvm` in the
output. If the output is empty, run `kvm-ok` for diagnostics. `sbx` requires
KVM to start.

Add your user to the `kvm` group:

```console
$ sudo usermod -aG kvm $USER
```

Sign out and back in, or run `newgrp kvm`, for the group change to take effect.

## Install on macOS

Install `sbx` using Homebrew:

```console
$ brew trust docker/tap
$ brew install docker/tap/sbx
```

## Install on Windows

Install `sbx` using Windows Package Manager:

```powershell
winget install -h Docker.sbx
```

## Install on Linux

You can install `sbx` with Docker Engine or install only the `sbx` package.

### Install Docker Engine and SBX

Run Docker's convenience script with `SBX=1` to install Docker Engine and the
`docker-sbx` package together:

```console
$ curl -fsSL https://get.docker.com | sudo SBX=1 sh
```

### Install SBX only

To install `sbx` without Docker Engine on the host, add Docker's `apt`
repository and install the `docker-sbx` package:

```console
$ curl -fsSL https://get.docker.com | sudo REPO_ONLY=1 sh
$ sudo apt-get install docker-sbx
```

## Install manually

To install `sbx` without a package manager, download a binary from the
[sbx-releases repository](https://github.com/docker/sbx-releases/releases).

## Sign in

Sign in to Docker:

```console
$ sbx login
```

The command opens a browser for Docker OAuth. See the [FAQ](faq.md) for why
sign-in is required and how Docker handles your data.

After signing in, [run your first sandbox](get-started.md).
