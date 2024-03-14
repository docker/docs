---
description:
  Configuring remote access allows Docker to accept requests from remote
  hosts by configuring it to listen on an IP address and port as well as the Unix
  socket
keywords: configuration, daemon, remote access, engine
title: Configure remote access for Docker daemon
---

By default, the Docker daemon listens for connections on a Unix socket to accept
requests from local clients. It's possible to allow Docker to accept requests
from remote hosts by configuring it to listen on an IP address and port as well
as the Unix socket. For more detailed information on this configuration option,
refer to the
[dockerd CLI reference](/reference/cli/dockerd/#bind-docker-to-another-hostport-or-a-unix-socket).

<!-- prettier-ignore -->
> **Warning**
>
> Before configuring Docker to accept connections from remote hosts it's
> critically important that you understand the security implications of opening
> Docker to the network. If steps aren't taken to secure the connection, it's
> possible for remote non-root users to gain root access on the host. For more
> information on how to use TLS certificates to secure this connection, check
> [Protect the Docker daemon socket](../../engine/security/protect-access.md).
{ .warning }

You can configure Docker to accept remote connections. This can be done using
the `docker.service` systemd unit file for Linux distributions using systemd. Or
you can use the `daemon.json` file, if your distribution doesn't use systemd.

> systemd vs `daemon.json`
>
> Configuring Docker to listen for connections using both the systemd unit file
> and the `daemon.json` file causes a conflict that prevents Docker from
> starting.

### Configuring remote access with systemd unit file

1. Use the command `sudo systemctl edit docker.service` to open an override file
   for `docker.service` in a text editor.

2. Add or modify the following lines, substituting your own values.

   ```systemd
   [Service]
   ExecStart=
   ExecStart=/usr/bin/dockerd -H fd:// -H tcp://127.0.0.1:2375
   ```

3. Save the file.

4. Reload the `systemctl` configuration.

   ```console
   $ sudo systemctl daemon-reload
   ```

5. Restart Docker.

   ```console
   $ sudo systemctl restart docker.service
   ```

6. Verify that the change has gone through.

   ```console
   $ sudo netstat -lntp | grep dockerd
   tcp        0      0 127.0.0.1:2375          0.0.0.0:*               LISTEN      3758/dockerd
   ```

### Configuring remote access with `daemon.json`

1. Set the `hosts` array in the `/etc/docker/daemon.json` to connect to the Unix
   socket and an IP address, as follows:

   ```json
   {
     "hosts": ["unix:///var/run/docker.sock", "tcp://127.0.0.1:2375"]
   }
   ```

2. Restart Docker.

3. Verify that the change has gone through.

   ```console
   $ sudo netstat -lntp | grep dockerd
   tcp        0      0 127.0.0.1:2375          0.0.0.0:*               LISTEN      3758/dockerd
   ```
