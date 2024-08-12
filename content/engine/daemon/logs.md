---
title: Read the daemon logs
description: How to read the event logs for the Docker daemon
keywords: docker, daemon, configuration, troubleshooting, logging
aliases:
  - /config/daemon/logs/
---

The daemon logs may help you diagnose problems. The logs may be saved in one of
a few locations, depending on the operating system configuration and the logging
subsystem used:

| Operating system                   | Location                                                                                                                                 |
| :--------------------------------- | :--------------------------------------------------------------------------------------------------------------------------------------- |
| Linux                              | Use the command `journalctl -xu docker.service` (or read `/var/log/syslog` or `/var/log/messages`, depending on your Linux Distribution) |
| macOS (`dockerd` logs)             | `~/Library/Containers/com.docker.docker/Data/log/vm/dockerd.log`                                                                         |
| macOS (`containerd` logs)          | `~/Library/Containers/com.docker.docker/Data/log/vm/containerd.log`                                                                      |
| Windows (WSL2) (`dockerd` logs)    | `%LOCALAPPDATA%\Docker\log\vm\dockerd.log`                                                                                               |
| Windows (WSL2) (`containerd` logs) | `%LOCALAPPDATA%\Docker\log\vm\containerd.log`                                                                                            |
| Windows (Windows containers)       | Logs are in the Windows Event Log                                                                                                        |

To view the `dockerd` logs on macOS, open a terminal Window, and use the `tail`
command with the `-f` flag to "follow" the logs. Logs will be printed until you
terminate the command using `CTRL+c`:

```console
$ tail -f ~/Library/Containers/com.docker.docker/Data/log/vm/dockerd.log
2021-07-28T10:21:21Z dockerd time="2021-07-28T10:21:21.497642089Z" level=debug msg="attach: stdout: begin"
2021-07-28T10:21:21Z dockerd time="2021-07-28T10:21:21.497714291Z" level=debug msg="attach: stderr: begin"
2021-07-28T10:21:21Z dockerd time="2021-07-28T10:21:21.499798390Z" level=debug msg="Calling POST /v1.41/containers/35fc5ec0ffe1ad492d0a4fbf51fd6286a087b89d4dd66367fa3b7aec70b46a40/wait?condition=removed"
2021-07-28T10:21:21Z dockerd time="2021-07-28T10:21:21.518403686Z" level=debug msg="Calling GET /v1.41/containers/35fc5ec0ffe1ad492d0a4fbf51fd6286a087b89d4dd66367fa3b7aec70b46a40/json"
2021-07-28T10:21:21Z dockerd time="2021-07-28T10:21:21.527074928Z" level=debug msg="Calling POST /v1.41/containers/35fc5ec0ffe1ad492d0a4fbf51fd6286a087b89d4dd66367fa3b7aec70b46a40/start"
2021-07-28T10:21:21Z dockerd time="2021-07-28T10:21:21.528203579Z" level=debug msg="container mounted via layerStore: &{/var/lib/docker/overlay2/6e76ffecede030507fcaa576404e141e5f87fc4d7e1760e9ce5b52acb24
...
^C
```

## Enable debugging

There are two ways to enable debugging. The recommended approach is to set the
`debug` key to `true` in the `daemon.json` file. This method works for every
Docker platform.

1.  Edit the `daemon.json` file, which is usually located in `/etc/docker/`. You
    may need to create this file, if it doesn't yet exist. On macOS or Windows,
    don't edit the file directly. Instead, edit the file through the Docker Desktop settings.

2.  If the file is empty, add the following:

    ```json
    {
      "debug": true
    }
    ```

    If the file already contains JSON, just add the key `"debug": true`, being
    careful to add a comma to the end of the line if it's not the last line
    before the closing bracket. Also verify that if the `log-level` key is set,
    it's set to either `info` or `debug`. `info` is the default, and possible
    values are `debug`, `info`, `warn`, `error`, `fatal`.

3.  Send a `HUP` signal to the daemon to cause it to reload its configuration.
    On Linux hosts, use the following command.

    ```console
    $ sudo kill -SIGHUP $(pidof dockerd)
    ```

    On Windows hosts, restart Docker.

Instead of following this procedure, you can also stop the Docker daemon and
restart it manually with the debug flag `-D`. However, this may result in Docker
restarting with a different environment than the one the hosts' startup scripts
create, and this may make debugging more difficult.

## Force a stack trace to be logged

If the daemon is unresponsive, you can force a full stack trace to be logged by
sending a `SIGUSR1` signal to the daemon.

- **Linux**:

  ```console
  $ sudo kill -SIGUSR1 $(pidof dockerd)
  ```

- **Windows Server**:

  Download [docker-signal](https://github.com/moby/docker-signal).

  Get the process ID of dockerd `Get-Process dockerd`.

  Run the executable with the flag `--pid=<PID of daemon>`.

This forces a stack trace to be logged but doesn't stop the daemon. Daemon logs
show the stack trace or the path to a file containing the stack trace if it was
logged to a file.

The daemon continues operating after handling the `SIGUSR1` signal and dumping
the stack traces to the log. The stack traces can be used to determine the state
of all goroutines and threads within the daemon.

## View stack traces

The Docker daemon log can be viewed by using one of the following methods:

- By running `journalctl -u docker.service` on Linux systems using `systemctl`
- `/var/log/messages`, `/var/log/daemon.log`, or `/var/log/docker.log` on older
  Linux systems

> **Note**
>
> It isn't possible to manually generate a stack trace on Docker Desktop for
> Mac or Docker Desktop for Windows. However, you can click the Docker taskbar
> icon and choose **Troubleshoot** to send information to Docker if you run into
> issues.

Look in the Docker logs for a message like the following:

```none
...goroutine stacks written to /var/run/docker/goroutine-stacks-2017-06-02T193336z.log
```

The locations where Docker saves these stack traces and dumps depends on your
operating system and configuration. You can sometimes get useful diagnostic
information straight from the stack traces and dumps. Otherwise, you can provide
this information to Docker for help diagnosing the problem.
