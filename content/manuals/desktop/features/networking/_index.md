---
description: Understand how Docker Desktop handles networking, firewalls, file access, proxies, and endpoint visibility.
keywords: docker desktop, networking, architecture, firewall, proxies, crowdstrike, vpn
title: Networking on Docker Desktop
linkTitle: Networking
weight: 30
--- 

This page explains how Docker Desktop routes network traffic and file I/O between containers, the VM, and the host, and how this behavior is visible to firewalls and endpoint protection tools.

## Overview

Docker Desktop runs the Docker Engine inside a lightweight Linux virtual machine (VM). Docker Desktop routes all network and file operations between the Docker VM and the host through the `com.docker.backend` process. 

Th backend acts as:

- A network proxy, translating traffic between the host and Linux VM.
- A file server, using gRPC FUSE which handles file access from containers to the host filesystem.
- A control plane, handling Docker API calls, port forwarding, and proxy settings.

## How containers connect to the internet

Each Linux container in Docker Desktop runs inside a small virtual network managed by Docker and every container is attached to a Docker-managed network and receives its own internal IP address. You can view and these networks with `docker network ls`, `docker network create`, and `docker network inspect`. They are managed by the [`daemon.json`](/manuals/engine/daemon/_index.md).

When a container initiates a network request, for example with `apt-get update` or `docker pull`:

- The container’s `eth0` interface connects to a virtual bridge (`docker0`) inside the VM.
- Outbound traffic from the container is sent through Network Address Translation (NAT) using a virtual adapter (typically with an internal IP such as `192.168.65.3`). You can view or change this with the [Docker Desktop settings](/manuals/desktop/settings-and-maintenance/settings.md#network).
- The traffic is transferred to the host system over a shared-memory channel rather than through a traditional virtual network interface. This approach ensures reliable communication and avoids conflicts with host-level network adapters or firewall configurations.
- On the host, Docker Desktop’s backend process receives the traffic and creates standard TCP/IP connections using the same networking APIs as other applications.

All outbound container network traffic originates from the `com.docker.backend` process. Firewalls, VPNs, and security tools, like Crowdstrike, see traffic coming from this process — not from a VM or unknown source so firewall and endpoint security software can apply rules directly to `com.docker.backend`.

## How exposed ports work

When you publish a container port using the `-p` or `--publish` flag, Docker Desktop makes that container port accessible from your host system or local network.

For example, with `docker run -p 80:80 nginx`:

- Docker Desktop's backend process listens on the specified host port, in this case, port `80`.
- When an application such as a web browser connects to that port, Docker Desktop forwards the connection into the Linux VM where the container is running over a shared-memory channel.
- Inside the VM, the connection is routed to the container’s internal IP address and port, for example `172.17.0.2:80`.
- The container responds through the same path, so you can access it from your host just like any other local service.

By default, `docker run -p` listens on all network interfaces (`0.0.0.0`), but you can restrict it to a specific address, such as `127.0.0.1` (localhost) or a particular network adapter.

Host firewalls can permit or deny inbound connections by filtering on `com.docker.backend`.

## Using Docker Desktop with a proxy

Docker Desktop can use your system’s default proxy settings or custom settings that you configure with [Docker Desktop's proxy setting](/manuals/desktop/settings-and-maintenance/settings.md#proxies). All proxy traffic passes through `com.docker.backend.exe`.

When a proxy is enabled:

- The backend process forwards the network requests, for example `docker pull`, through an internal proxy at `http.docker.internal:3128`. 
- The internal proxy then connects either directly to the internet or through your upstream proxy, depending on your configuration and adding authentication if necessary. 
- Docker Desktop then downloads the requested images or data through the proxy as usual.

Note that: 
- The proxy honors system or manual proxy configuration.
- On Windows, Basic, NTLM, and Kerberos authentication is supported.
- For Mac, NTLM/Kerberos is not supported natively. Run a local proxy on `localhost` as a workaround.
- CLI plugins and other tools that use the Docker API directly must be configured separately with the `HTTP_PROXY`, `HTTPS_PROXY`, and `NO_PROXY` environment variables.
  
## Firewalls and endpoint visibility 

Docker Desktop doesn't include a built-in firewall. 

To restrict VM or container networking apply rules to `com.docker.backend.exe` (Windows) `com.docker.backend` (Mac) or `qemu` (Linux) as all VM networking is funneled through these processes. 

Use Windows Defender Firewall or enterprise endpoint firewalls for control. This enables traffic inspection and restriction at the host level without modifying the Docker Engine. 

CrowdStrike and similar tools can observe all traffic and file access that passes through the backend process. To monitor in-VM operations, install the agent inside the Docker VM.

| Action | Visible to host EDR? | Reason | 
|---------|----------------------|---------| 
| Container reads host files | Yes | Access handled by `com.docker.backend` | 
| Container writes host files | Yes | Same process performs the write | 
| Container accesses its own filesystem layers | No | Exists only inside the VM |
| Endpoint agent inside VM | Yes | Full visibility | 
