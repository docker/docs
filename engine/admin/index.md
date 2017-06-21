---
description: Configuring and troubleshooting the Docker daemon
keywords: docker, daemon, configuration, troubleshooting
redirect_from:
- /engine/articles/configuring/
- /engine/admin/configuring/
title: Configure and troubleshoot the Docker daemon
---

After successfully installing Docker and starting Docker, the `dockerd` daemon
runs with its default configuration. This topic shows how to customize
the configuration, start the daemon manually, and troubleshoot and debug the
daemon if you run into issues.

## Start the daemon using operating system utilities

The command to start Docker depends on your operating system. Check the correct
page under [Install Docker](/engine/installation/index.md). To configure Docker
to start automatically at system boot, see
[Configure Docker to start on boot](/engine/installation/linux/linux-postinstall.md#configure-docker-to-start-on-boot)

## Start the daemon manually

Typically, you start Docker using operating system utilities. For debugging
purposes, you can start Docker manually using the `dockerd` command. You
may need to use `sudo`, depending on your operating system configuration. When
you start Docker this way, it runs in the foreground and sends its logs directly
to your terminal.

```bash
$ dockerd

INFO[0000] +job init_networkdriver()
INFO[0000] +job serveapi(unix:///var/run/docker.sock)
INFO[0000] Listening for HTTP on unix (/var/run/docker.sock)
...
...
```

To stop Docker when you have started it manually, issue a `Ctrl+C` in your
terminal.

## Configure the Docker daemon

The daemon includes many configuration options, which you can pass as flags
when starting Docker manually, or set in the `daemon.json` configuration file.
The second method is recommended because those configuration changes persist
when you restart Docker.

See [dockerd](/engine/reference/commandline/dockerd.md) for a full list of
configuration options.

Here is an example of starting the Docker daemon manually with some configuration
options:

```bash
$ dockerd -D --tls=true --tlscert=/var/docker/server.pem --tlskey=/var/docker/serverkey.pem -H tcp://192.168.59.3:2376
```

This command enables debugging (`-D`), enables TLS (`-tls`), specifies the server
certificate and key (`--tlscert` and `--tlskey`), and specifies the network
interface where the daemon listens for connections (`-H`).

A better approach is to put these options into the `daemon.json` file and
restart Docker. This method works for every Docker platform. The following
`daemon.json` example sets all the same options as the above command:

```json
{
  "debug": true,
  "tls": true,
  "tlscert": "/var/docker/server.pem",
  "tlskey": "/var/docker/serverkey.pem",
  "hosts": ["tcp://192.168.59.3:2376"]
}
```

## Isolate containers within a user namespace

Linux namespaces provide isolation for running processes, limiting their access
to system resources without the running process being aware of the limitations.
For more information on Linux namespaces, see
[Linux namespaces](https://www.linux.com/news/understanding-and-securing-linux-namespaces){: target="_blank" class="_" };

The best practice for preventing privilege-escalation attacks from within a
container is to configure your container's applications to run as unprivileged
users. For containers whose processes must run as the `root` user within the
container, you can re-map this user to a less-privileged user on the Docker
host.

This re-mapping is relatively transparent to the container, but introduces some
configuration complexity when the container needs access to privileged resources
on the Docker host, such as bind mounts into areas of the filesystem that the
system user cannot write to.

### Prerequisites

1.  The username or user ID (and optionally group name or group ID) must exist
    already. Typically, this means that the relevant entries need to be in
    `/etc/password` and `/etc/group`, but if you are using a different
    authentication back-end, this requirement may translate differently.

    To verify this, use the `id` command:

    ```bash
    $ id testuser
    
    uid=1001(testuser) gid=1001(testuser) groups=1001(testuser)
    ```

2.  If there are any locations on the Docker host where the unprivileged
    user needs to write, adjust the permissions of those locations
    accordingly.

3.  The way the namespace remapping is handled on the host is using two files,
    `/etc/subuid` and `/etc/subgid`. These files are typically managed
    automatically when you add or remove users or groups, but on a few
    distributions such as RHEL and CentOS 7.3, you may need to manage these
    files manually.

    Each file contains three fields: the username or ID of the user, followed by
    a beginning UID or GID (which is treated as UID or GID 0 within the namespace)
    and a maxumum number of UIDs or GIDs available to the user. For instance,
    given the following entry:

    ```none
    testuser:231072:65536
    ```
   
    This means that user-namespaced processes started by `testuser` will be owne
    by host UID `231072` (which will look like UID `0` inside the namespace)
    through 296608 (231072 + 65536). These ranges should not overlap, to ensure
    that namespaced processes cannot access each other's namespaces.

    After adding your user, check `/etc/subuid` and `/etc/subgid` to see if your
    user has an entry in each. If not, you need to add it, being careful to
    avoid overlap.

4.  Enabling `userns-remap` will effectively mask existing image and container
    layers, as well as other Docker objects within `/var/lib/docker/`. This is
    because Docker needs to adjust the ownership of these resources and actually
    stores them in a subdirectory within `/var/lib/docker/`. It is best to enable
    this feature on a new Docker installation rather than an existing one.

    Along the same lines, if you disable `userns-remap` you will not see any
    of the resources created while it was enabled.

#### Enable userns-remap

1.  Edit `/etc/docker/daemon.json`. Assuming the file was previously empty, the
    following entry will enable `userns-remap` using user and group called
    `testuser`. You can address the user and group by ID or name. You only need to
    specify the group name or ID if it is different from the user name or ID.

    ```json
    {
      "userns-remap": "testuser"
    }
    ```

    Save the file and restart Docker.

2.  Verify that previous images are not available using the `docker image ls`
    command. The output should be empty.

3.  Start a container from the `hello-world` image.

    ```bash
    $ docker run hello-world
    ```

4.  Verify that a namespaced directory exists within `/var/lib/docker/` named
    with the UID and GID of the namespaced user, owned by that UID and GID,
    and not group-or-world-readable. Some of the subdirectories are still
    owned by `root` and have different permissions.

    ```bash
    $ sudo ls -ld /var/lib/docker/231072.231072/

    drwx------ 11 231072 231072 11 Jun 21 21:19 /var/lib/docker/231072.231072/

    $ sudo ls -l /var/lib/docker/231072.231072/

    total 14
    drwx------ 5 231072 231072 5 Jun 21 21:19 aufs
    drwx------ 3 231072 231072 3 Jun 21 21:21 containers
    drwx------ 3 root   root   3 Jun 21 21:19 image
    drwxr-x--- 3 root   root   3 Jun 21 21:19 network
    drwx------ 4 root   root   4 Jun 21 21:19 plugins
    drwx------ 2 root   root   2 Jun 21 21:19 swarm
    drwx------ 2 231072 231072 2 Jun 21 21:21 tmp
    drwx------ 2 root   root   2 Jun 21 21:19 trust
    drwx------ 2 231072 231072 3 Jun 21 21:19 volumes
    ```

    Your directory listing may have some differences, especially if you
    user a different container storage driver than `aufs`.

    The directories which are owned by the remapped user are used instead
    of the same directories directly beneath `/var/lib/docker/` and the
    unused versions (such as `/var/lib/docker/tmp/` in the example here)
    can be removed. Docker will not use them while `userns-remap` is
    enabled.

## Troubleshoot the daemon

You can enable debugging on the daemon to learn about the runtime activity of
the daemon and to aid in troubleshooting. If the daemon is completely
non-responsive, you can also
[force a full stack trace](#force-a-full-stack-trace-to-be-logged) of all
threads to be added to the daemon log by sending the `SIGUSR` signal to the
Docker daemon.

### Out Of Memory Exceptions (OOME)

If your containers attempt to use more memory than the system has available,
you may experience an Out Of Memory Exception (OOME) and a container, or the
Docker daemon, might be killed by the kernel OOM killer. To prevent this from
happening, ensure that your application runs on hosts with adequate memory and
see
[Understand the risks of running out of memory](/engine/admin/resource_constraints.md#understand-the-risks-of-running-out-of-memory).

### Read the logs

The daemon logs may help you diagnose problems. The logs may be saved in one of
a few locations, depending on the operating system configuration and the logging
subsystem used:

| Operating system | Location |
|------------------|----------|
| RHEL, Oracle Linux | `/var/log/messages` |
| Debian           | `/var/log/daemon.log` |
| Ubuntu 16.04+, CentOS | Use the command `journalctl -u docker.service` |
| Ubuntu 14.10-    | `/var/log/upstart/docker.log` |
| macOS            | `~/Library/Containers/com.docker.docker/Data/com.docker.driver.amd64-linux/console-ring` |
| Windows          | `AppData\Local` |


### Enable debugging

There are two ways to enable debugging. The recommended approach is to set the
`debug` key to `true` in the `daemon.json` file. This method works for every
Docker platform.

1.  Edit the `daemon.json` file, which is usually located in `/etc/docker/`.
    You may need to create this file, if it does not yet exist. On macOS or
    Windows, do not edit the file directly. Instead, go to
    **Preferences** / **Daemon** / **Advanced**.

2.  If the file is empty, add the following:

    ```json
    {
      "debug": true
    }
    ```

    If the file already contains JSON, just add the key `"debug": true`, being
    careful to add a comma to the end of the line if it is not the last line
    before the closing bracket. Also verify that if the `log-level` key is set,
    it is set to either `info` or `debug`. `info` is the default, and possible
    values are `debug`, `info`, `warn`, `error`, `fatal`.

3.  Send a `HUP` signal to the daemon to cause it to reload its configuration.
    On Linux hosts, use the following command.

    ```bash
    $ sudo kill -SIGHUP $(pidof dockerd)
    ```

    On Windows hosts, restart Docker.

Instead of following this procedure, you can also stop the Docker daemon and
restart it manually with the `-D` flag. However, this may result in Docker
restarting with a different environment than the one the hosts's startup scripts
will create, and this may make debugging more difficult.

### Force a stack trace to be logged

If the daemon is unresponsive, you can force a full stack trace to be logged
by sending a `SIGUSR1` signal to the daemon.

- **Linux**:

  ```bash
  $ sudo kill -SIGUSR1 $(pidof dockerd)
  ```

- **Windows Server**:

  Download [docker-signal](https://github.com/jhowardmsft/docker-signal).

  Run the executable with the flag `--pid=<PID of daemon>`.

This will force a stack trace to be logged but will not stop the daemon.

The daemon will continue operating after handling the `SIGUSR1` signal and
dumping the stack traces to the log. The stack traces can be used to determine
the state of all goroutines and threads within the daemon.

## Check whether Docker is running

The operating-system independent way to check whether Docker is running is to
ask Docker, using the `docker info` command.

You can also use operating system utilities, such as
`sudo systemctl is-active docker` or `sudo status docker` or
`sudo service docker status`, or checking the service status using Windows
utilities.

Finally, you can check in the process list for the `dockerd` process, using
commands like `ps` or `top`.

