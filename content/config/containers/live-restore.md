---
description: Learn how to keep containers running when the daemon isn't available
keywords: docker, upgrade, daemon, dockerd, live-restore, daemonless container
title: Live restore
aliases:
  - /engine/admin/live-restore/
---

By default, when the Docker daemon terminates, it shuts down running containers.
You can configure the daemon so that containers remain running if the daemon
becomes unavailable. This functionality is called _live restore_. The live restore
option helps reduce container downtime due to daemon crashes, planned outages,
or upgrades.

> **Note**
>
> Live restore isn't supported for Windows containers, but it does work for
> Linux containers running on Docker Desktop for Windows.

## Enable live restore

There are two ways to enable the live restore setting to keep containers alive
when the daemon becomes unavailable. **Only do one of the following**.

- Add the configuration to the daemon configuration file. On Linux, this
  defaults to `/etc/docker/daemon.json`. On Docker Desktop for Mac or Docker
  Desktop for Windows, select the Docker icon from the task bar, then click
  **Settings** -> **Docker Engine**.

  - Use the following JSON to enable `live-restore`.

    ```json
    {
      "live-restore": true
    }
    ```

  - Restart the Docker daemon. On Linux, you can avoid a restart (and avoid any
    downtime for your containers) by reloading the Docker daemon. If you use
    `systemd`, then use the command `systemctl reload docker`. Otherwise, send a
    `SIGHUP` signal to the `dockerd` process.

- If you prefer, you can start the `dockerd` process manually with the
  `--live-restore` flag. This approach isn't recommended because it doesn't
  set up the environment that `systemd` or another process manager would use
  when starting the Docker process. This can cause unexpected behavior.

## Live restore during upgrades

Live restore allows you to keep containers running across Docker daemon updates,
but is only supported when installing patch releases (`YY.MM.x`), not for
major (`YY.MM`) daemon upgrades.

If you skip releases during an upgrade, the daemon may not restore its
connection to the containers. If the daemon can't restore the connection, it
can't manage the running containers and you must stop them manually.

## Live restore upon restart

The live restore option only works to restore containers if the daemon options,
such as bridge IP addresses and graph driver, didn't change. If any of these
daemon-level configuration options have changed, the live restore may not work
and you may need to manually stop the containers.

## Impact of live restore on running containers

If the daemon is down for a long time, running containers may fill up the FIFO
log the daemon normally reads. A full log blocks containers from logging more
data. The default buffer size is 64K. If the buffers fill, you must restart
the Docker daemon to flush them.

On Linux, you can modify the kernel's buffer size by changing
`/proc/sys/fs/pipe-max-size`. You can't modify the buffer size on Docker Desktop for
Mac or Docker Desktop for Windows.

## Live restore and Swarm mode

The live restore option only pertains to standalone containers, and not to Swarm
services. Swarm services are managed by Swarm managers. If Swarm managers are
not available, Swarm services continue to run on worker nodes but can't be
managed until enough Swarm managers are available to maintain a quorum.
