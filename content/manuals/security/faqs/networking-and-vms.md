---
title: Network and VM FAQs
linkTitle: Network and VM
description: Frequently asked questions about Docker Desktop networking and virtualization security
keywords: docker desktop networking, virtualization, hyper-v, wsl2, network security, firewall
weight: 30
tags: [FAQ]
aliases:
- /faq/security/networking-and-vms/
---

## How can I limit container internet access to prevent data exfiltration?

Docker Desktop doesn't have a built-in mechanism for this, but you can use process-level firewalls on the host. Apply rules to the `com.docker.vpnkit` user-space process to control where it can connect (DNS allowlists, packet filters) and which ports/protocols it can use.

For enterprise environments, consider [Air-gapped containers](/manuals/enterprise/security/hardened-desktop/air-gapped-containers.md) which provide network access controls for containers.

## Can I apply firewall rules to container network traffic?

Yes. Docker Desktop uses a user-space process (`com.docker.vpnkit`) for network connectivity, which inherits constraints like firewall rules, VPN settings, and HTTP proxy properties from the user that launched it.

## Does Docker Desktop for Windows with Hyper-V allow users to create other VMs?

No. The `DockerDesktopVM` name is hard-coded in the service, so you cannot use Docker Desktop to create or manipulate other virtual machines.

## How does Docker Desktop achieve network isolation with Hyper-V and WSL 2?

Docker Desktop uses the same VM processes for both WSL 2 (in the `docker-desktop` distribution) and Hyper-V (in `DockerDesktopVM`). Host/VM communication uses `AF_VSOCK` hypervisor sockets (shared memory) rather than network switches or interfaces. All host networking is performed using standard TCP/IP sockets from the `com.docker.vpnkit.exe` and `com.docker.backend.exe` processes.

For more information, see [How Docker Desktop networking works under the hood](https://www.docker.com/blog/how-docker-desktop-networking-works-under-the-hood/).
