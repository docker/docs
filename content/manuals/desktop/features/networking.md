---
description: Understand how networking works on Docker Desktop and see the known limitations
keywords: networking, docker desktop, proxy, vpn, Linux, Mac, Windows
title: Explore networking features on Docker Desktop
linkTitle: Networking
aliases:
- /desktop/linux/networking/
- /docker-for-mac/networking/
- /mackit/networking/
- /desktop/mac/networking/
- /docker-for-win/networking/
- /docker-for-windows/networking/
- /desktop/windows/networking/
- /desktop/networking/
weight: 30
---

Docker Desktop includes built-in networking capabilities to help you connect containers with services on your host, across containers, or through proxies and VPNs.

## Networking features for all platforms

### VPN Passthrough

Docker Desktop networking can work when attached to a VPN. To do this,
Docker Desktop intercepts traffic from the containers and injects it into
the host as if it originated from the Docker application.

### Port mapping

When you run a container with the `-p` argument, for example:

```console
$ docker run -p 80:80 -d nginx
```

Docker Desktop makes whatever is running on port `80` in the container, in
this case, `nginx`, available on port `80` of `localhost`. In this example, the
host and container ports are the same. 

To avoid conflicts with services already using port `80` on the host:

```console
$ docker run -p 8000:80 -d nginx
```

Now connections to `localhost:8000` are sent to port `80` in the container. 

> [!TIP]
>
> The syntax for `-p` is `HOST_PORT:CLIENT_PORT`.

### HTTP/HTTPS Proxy support

See [Proxies](/manuals/desktop/settings-and-maintenance/settings.md#proxies)

### SOCKS5 proxy support

{{< summary-bar feature_name="SOCKS5 proxy support" >}}

SOCKS (Socket Secure) is a protocol that facilitates the routing of network packets between a client and a server through a proxy server. It provides a way to enhance privacy, security, and network performance for users and applications. 

You can enable SOCKS proxy support to allow outgoing requests, such as pulling images, and access Linux container backend IPs from the host. 

To enable and set up SOCKS proxy support:

1. Navigate to the **Resources** tab in **Settings**. 
2. From the dropdown menu select **Proxies**.
3. Switch on the **Manual proxy configuration** toggle. 
4. In the **Secure Web Server HTTPS** box, paste your `socks5://host:port` URL.

## Networking mode and DNS behaviour for Mac and Windows

With Docker Desktop version 4.42 and later, you can customize how Docker handles container networking and DNS resolution to better support a range of environments — from IPv4-only to dual-stack and IPv6-only systems. These settings help prevent timeouts and connectivity issues caused by incompatible or misconfigured host networks.

> [!NOTE]
>
> These settings can be overridden on a per-network basis using CLI flags or Compose file options.

### Default networking mode

Choose the default IP protocol used when Docker creates new networks. This allows you to align Docker with your host’s network capabilities or organizational requirements, such as enforcing IPv6-only access.

The options available are:

- **Dual IPv4/IPv6** (Default): Supports both IPv4 and IPv6. Most flexible and ideal for environments with dual-stack networking.
- **IPv4 only**: Only IPv4 addresses are used. Use this if your host or network does not support IPv6.
- **IPv6 only**: Only IPv6 addresses are used. Best for environments transitioning to or enforcing IPv6-only connectivity.

> [!NOTE]
>
> This setting can be overridden on a per-network basis using CLI flags or Compose file options.

### DNS resolution behavior 

Control how Docker filters DNS records returned to containers, improving reliability in environments where only IPv4 or IPv6 is supported. This setting is especially useful for preventing apps from trying to connect using IP families that aren't actually available, which can cause avoidable delays or failures.

Depending on your selected network mode, the options available are:

- **Auto (recommended)**: Docker detects your host's network stack and automatically filters out unsupported DNS record types (A for IPv4, AAAA for IPv6).
- **Filter IPv4 (A records)**: Prevents containers from resolving IPv4 addresses. Only available in dual-stack mode.
- **Filter IPv6 (AAAA records)**: Prevents containers from resolving IPv6 addresses. Only available in dual-stack mode.
- **No filtering**: Docker returns all DNS records (A and AAAA), regardless of host support.

> [!IMPORTANT]
>
> Switching the default networking mode resets the DNS filter to Auto.

### Using Settings Management

If you're an administrator, you can use [Settings Management](/manuals/enterprise/security/hardened-desktop/settings-management/configure-json-file.md#networking) to enforce this Docker Desktop setting across your developer's machines. Choose from the following code snippets and at it to your `admin-settings.json` file,
or configure this setting using the [Admin Console](/manuals/enterprise/security/hardened-desktop/settings-management/configure-admin-console.md)

{{< tabs >}}
{{< tab name="Networking mode" >}}

Dual IPv4/IPv6:

```json
{
  "defaultNetworkingMode": {
    "locked": true
    "value": "dual-stack"
  }
}
```

IPv4 only:

```json
{
  "defaultNetworkingMode": {
    "locked": true
    "value": "ipv4only"
  }
}
```

IPv6 only:

```json
{
  "defaultNetworkingMode": {
    "locked": true
    "value": "ipv6only"
  }
}
```

{{< /tab >}}
{{< tab name="DNS resolution" >}}

Auto filter:

```json
{
  "dnsInhibition": {
    "locked": true
    "value": "auto"
  }
}
```

Filter IPv4:

```json
{
  "dnsInhibition": {
    "locked": true
    "value": "ipv4"
  }
}
```

Filter IPv6:

```json
{
  "dnsInhibition": {
    "locked": true
    "value": "ipv6"
  }
}
```

No filter:

```json
{
  "dnsInhibition": {
    "locked": true
    "value": "none"
  }
}
```

{{< /tab >}}
{{< /tabs >}}

## Networking features for Mac and Linux

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

## Use cases and workarounds 

### I want to connect from a container to a service on the host

The host has a changing IP address, or none if you have no network access.
Docker recommends you connect to the special DNS name `host.docker.internal`,
which resolves to the internal IP address used by the host.

You can also reach the gateway using `gateway.docker.internal`.

If you have installed Python on your machine, use the following instructions as an example to connect from a container to a service on the host:

1. Run the following command to start a simple HTTP server on port 8000.

    `python -m http.server 8000`

    If you have installed Python 2.x, run `python -m SimpleHTTPServer 8000`.

2. Now, run a container, install `curl`, and try to connect to the host using the following commands:

    ```console
    $ docker run --rm -it alpine sh
    # apk add curl
    # curl http://host.docker.internal:8000
    # exit
    ```

### I want to connect to a container from the host

Port forwarding works for `localhost`. `--publish`, `-p`, or `-P` all work.
Ports exposed from Linux are forwarded to the host.

Docker recommends you publish a port, or to connect from another
container. This is what you need to do even on Linux if the container is on an
overlay network, not a bridge network, as these are not routed.

For example, to run an `nginx` webserver:

```console
$ docker run -d -p 80:80 --name webserver nginx
```

To clarify the syntax, the following two commands both publish container's port `80` to host's port `8000`:

```console
$ docker run --publish 8000:80 --name webserver nginx

$ docker run -p 8000:80 --name webserver nginx
```

To publish all ports, use the `-P` flag. For example, the following command
starts a container (in detached mode) and the `-P` flag publishes all exposed ports of the
container to random ports on the host.

```console
$ docker run -d -P --name webserver nginx
```

Alternatively, you can also use [host networking](/manuals/engine/network/drivers/host.md#docker-desktop)
to give the container direct access to the network stack of the host.

See the [run command](/reference/cli/docker/container/run.md) for more details on
publish options used with `docker run`.
