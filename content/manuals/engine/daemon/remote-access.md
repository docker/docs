---
description:
  Configuring remote access allows Docker to accept requests from remote
  hosts by configuring it to listen on an IP address and port as well as the Unix
  socket
keywords: configuration, daemon, remote access, engine
title: Configure remote access for Docker daemon
aliases:
  - /config/daemon/remote-access/
---

By default, the Docker daemon listens for connections on a Unix socket to accept
requests from local clients. You can configure Docker to accept requests
from remote clients by configuring it to listen on an IP address and port as well
as the Unix socket.

<!-- prettier-ignore -->
> [!WARNING]
>
> Configuring Docker to accept connections from remote clients can leave you
> vulnerable to unauthorized access to the host and other attacks.
>
> It's critically important that you understand the security implications of opening Docker to the network.
> If steps aren't taken to secure the connection, it's possible for remote non-root users to gain root access on the host.
>
> Remote access without TLS is **not recommended**, and will require explicit opt-in in a future release.
> For more information on how to use TLS certificates to secure this connection, see
> [Protect the Docker daemon socket](/manuals/engine/security/protect-access.md).

## Enable remote access

You can enable remote access to the daemon either using a `docker.service` systemd unit file for Linux distributions using systemd.
Or you can use the `daemon.json` file, if your distribution doesn't use systemd.

Configuring Docker to listen for connections using both the systemd unit file
and the `daemon.json` file causes a conflict that prevents Docker from starting.

### Configuring remote access with systemd unit file

1. Use the command `sudo systemctl edit docker.service` to open an override file
   for `docker.service` in a text editor.

2. Add or modify the following lines, substituting your own values.

   ```systemd
   [Service]
   ExecStart=
   ExecStart=/usr/bin/dockerd -H fd:// -H tcp://127.0.0.1:2375
