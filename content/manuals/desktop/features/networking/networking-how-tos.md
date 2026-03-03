---
description: Learn how to connect containers to the host, across containers, or through proxies and VPNs in Docker Desktop.
keywords: docker desktop, networking, vpn, proxy, port mapping, dns
title: Explore networking how-tos on Docker Desktop
linkTitle: How-tos
aliases:
- /desktop/linux/networking/
- /docker-for-mac/networking/
- /mackit/networking/
- /desktop/mac/networking/
- /docker-for-win/networking/
- /docker-for-windows/networking/
- /desktop/windows/networking/
- /desktop/networking/
---

This page explains how to configure and use networking features, connect containers to host services, work behind proxies or VPNs, and troubleshoot common issues.

For details on how Docker Desktop routes network traffic and file I/O between containers, the VM, and the host, see [Network overview](/manuals/desktop/features/networking/index.md#overview).

## Core networking how-tos

### Connect a container to a service on the host

The host has a changing IP address, or none if you have no network access. To connect to services running on your host, use the special DNS name:

| Name                      | Description                                      |
| ------------------------- | ------------------------------------------------ |
| `host.docker.internal`    | Resolves to the internal IP address of your host |
| `gateway.docker.internal` | Resolves to the gateway IP of the Docker VM      |


#### Example

Run a simple HTTP server on port `8000`:

```console
$ python -m http.server 8000
```

Then run a container, install `curl`, and try to connect to the host using the following commands:

```console
$ docker run --rm -it alpine sh
# apk add curl
# curl http://host.docker.internal:8000
# exit
```

### Connect to a container from the host

To access containerized services from your host or local network, publish ports with the `-p` or `--publish` flag. For example:

```console
$ docker run -d -p 80:80 --name webserver nginx
```

Docker Desktop makes whatever is running on port `80` in the container, in
this case, `nginx`, available on port `80` of `localhost`. 

> [!TIP]
>
> The syntax for `-p` is `HOST_PORT:CLIENT_PORT`.

To publish all ports, use the `-P` flag. For example, the following command
starts a container (in detached mode) and the `-P` flag publishes all exposed ports of the
container to random ports on the host.

```console
$ docker run -d -P --name webserver nginx
``` 

Alternatively, you can also use [host networking](/manuals/engine/network/drivers/host.md#docker-desktop)
to give the container direct access to the network stack of the host.

See the [run command](/reference/cli/docker/container/run/) for more details on
publish options used with `docker run`.

All inbound connections pass through the Docker Desktop backend process (`com.docker.backend` (Mac), `com.docker.backend` (Windows), or `qemu` (Linux), which handles port forwarding into the VM.
For more details, see [How exposed ports work](/manuals/desktop/features/networking/index.md#how-exposed-ports-work)

### Working with VPNs

Docker Desktop networking can work when attached to a VPN. 

To do this, Docker Desktop intercepts traffic from the containers and injects it into
the host as if it originated from the Docker application.

For details about how this traffic appears to host firewalls and endpoint detection systems, see [Firewalls and endpoint visibility](/manuals/desktop/features/networking/index.md#firewalls-and-endpoint-visibility).

### Working with proxies

Docker Desktop can use your system proxy or a manual configuration.
To configure proxies:

1. Navigate to the **Resources** tab in **Settings**. 
2. From the dropdown menu select **Proxies**.
3. Switch on the **Manual proxy configuration** toggle.
4. Enter your HTTP, HTTPS or SOCKS5 proxy URLS.

For more details on proxies and proxy configurations, see the [Proxy settings documentation](/manuals/desktop/settings-and-maintenance/settings.md#proxies).

## Network how-tos for Mac and Windows

You can control how Docker handles container networking and DNS resolution to better support a range of environments — from IPv4-only to dual-stack and IPv6-only systems. These settings help prevent timeouts and connectivity issues caused by incompatible or misconfigured host networks.

You can set the following settings on the **Network** tab in the Docker Desktop Dashboard settings, or if you're an admin, with Settings Management via the [`admin-settings.json` file](/manuals/enterprise/security/hardened-desktop/settings-management/configure-json-file.md#networking), or the [Admin Console](/manuals/enterprise/security/hardened-desktop/settings-management/configure-admin-console.md)

> [!NOTE]
>
> These settings can be overridden on a per-network basis using CLI flags or Compose file options.

### Default networking mode

Choose the default IP protocol used when Docker creates new networks. This allows you to align Docker with your host’s network capabilities or organizational requirements, such as enforcing IPv6-only access.

| Mode                         | Description                                 |
| ---------------------------- | ------------------------------------------- |
| **Dual IPv4/IPv6 (default)** | Supports both IPv4 and IPv6. Most flexible. |
| **IPv4 only**                | Uses only IPv4 addressing.                  |
| **IPv6 only**                | Uses only IPv6 addressing.                  |

### DNS resolution behavior 

Control how Docker filters DNS records returned to containers, improving reliability in environments where only IPv4 or IPv6 is supported. This setting is especially useful for preventing apps from trying to connect using IP families that aren't actually available, which can cause avoidable delays or failures.

| Option                         | Description                                                                 |
| ------------------------------ | --------------------------------------------------------------------------- |
| **Auto (recommended)**         | Automatically filters unsupported record types. (A for IPv4, AAAA for IPv6) |
| **Filter IPv4 (A records)**    | Blocks IPv4 lookups. Only available in dual-stack mode.                     |
| **Filter IPv6 (AAAA records)** | Blocks IPv6 lookups. Only available in dual-stack mode.                     |
| **No filtering**               | Returns both A and AAAA records.                                            |

> [!IMPORTANT]
>
> Switching the default networking mode resets the DNS filter to Auto.

## Network how-tos for Mac and Linux

### SSH agent forwarding

Docker Desktop for Mac and Linux lets you use the host’s SSH agent inside a container. To do this:

1. Bind mount the SSH agent socket by adding the following parameter to your `docker run` command:

   ```console
   $--mount type=bind,src=/run/host-services/ssh-auth.sock,target=/run/host-services/ssh-auth.sock
   ```

2. Add the `SSH_AUTH_SOCK` environment variable in your container:

    ```console
    $ -e SSH_AUTH_SOCK="/run/host-services/ssh-auth.sock"
    ```

To enable the SSH agent in Docker Compose, add the following flags to your service:

 ```yaml
services:
  web:
    image: nginx:alpine
    volumes:
      - type: bind
        source: /run/host-services/ssh-auth.sock
        target: /run/host-services/ssh-auth.sock
    environment:
      - SSH_AUTH_SOCK=/run/host-services/ssh-auth.sock
 ```

## Known limitations

### Changing internal IP addresses

The internal IP addresses used by Docker can be changed from **Settings**. After changing IPs, you need to reset the Kubernetes cluster and to leave any active Swarm.

### There is no `docker0` bridge on the host

Because of the way networking is implemented in Docker Desktop, you cannot
see a `docker0` interface on the host. This interface is actually within the
virtual machine.

### I cannot ping my containers

Docker Desktop can't route traffic to Linux containers. However if you're a Windows user, you can
ping the Windows containers.

### Per-container IP addressing is not possible

This is because the Docker `bridge` network is not reachable from the host.
However if you are a Windows user, per-container IP addressing is possible with Windows containers.
