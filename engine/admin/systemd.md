---
description: Controlling and configuring Docker using systemd
keywords: docker, daemon, systemd,  configuration
redirect_from:
- /engine/articles/systemd/
title: Control and configure Docker with systemd
---

Many Linux distributions use systemd to start the Docker daemon. This document
shows a few examples of how to customize Docker's settings.

## Starting the Docker daemon

Once Docker is installed, you will need to start the Docker daemon.

```bash
$ sudo systemctl start docker
# or on older distributions, you may need to use
$ sudo service docker start
```

If you want Docker to start at boot, you should also:

```bash
$ sudo systemctl enable docker
# or on older distributions, you may need to use
$ sudo chkconfig docker on
```

## Custom Docker daemon options

There are a number of ways to configure the daemon flags and environment variables
for your Docker daemon.

The recommended way is to use a systemd drop-in file (as described in the <a
target="_blank"
href="https://www.freedesktop.org/software/systemd/man/systemd.unit.html">systemd.unit</a>
documentation). These are local files named `<something>.conf` in the
`/etc/systemd/system/docker.service.d` directory.

However, if you had previously used a package which had an `EnvironmentFile`
(often pointing to `/etc/sysconfig/docker`) then for backwards compatibility,
you drop a file with a `.conf` extension into the
`/etc/systemd/system/docker.service.d` directory including the following:

```conf
[Service]
EnvironmentFile=-/etc/sysconfig/docker
EnvironmentFile=-/etc/sysconfig/docker-storage
EnvironmentFile=-/etc/sysconfig/docker-network
ExecStart=
ExecStart=/usr/bin/dockerd $OPTIONS \
          $DOCKER_STORAGE_OPTIONS \
          $DOCKER_NETWORK_OPTIONS \
          $BLOCK_REGISTRY \
          $INSECURE_REGISTRY
```

To check if the `docker.service` uses an `EnvironmentFile`:

```bash
$ systemctl show docker | grep EnvironmentFile

EnvironmentFile=-/etc/sysconfig/docker (ignore_errors=yes)
```

Alternatively, find out where the service file is located:

```bash
$ systemctl show --property=FragmentPath docker

FragmentPath=/usr/lib/systemd/system/docker.service

$ grep EnvironmentFile /usr/lib/systemd/system/docker.service

EnvironmentFile=-/etc/sysconfig/docker
```

You can customize the Docker daemon options using override files as explained in
the [HTTP Proxy example](systemd.md#http-proxy) below. The files located in
`/usr/lib/systemd/system` or `/lib/systemd/system` contain the default options
and should not be edited.

### Runtime directory and storage driver

You may want to control the disk space used for Docker images, containers
and volumes by moving it to a separate partition.

In this example, we'll assume that your `docker.service` file looks something
like:

```conf
[Unit]
Description=Docker Application Container Engine
Documentation=https://docs.docker.com
After=network.target

[Service]
Type=notify
# the default is not to use systemd for cgroups because the delegate issues still
# exists and systemd currently does not support the cgroup feature set required
# for containers run by docker
ExecStart=/usr/bin/dockerd
ExecReload=/bin/kill -s HUP $MAINPID
# Having non-zero Limit*s causes performance problems due to accounting overhead
# in the kernel. We recommend using cgroups to do container-local accounting.
LimitNOFILE=infinity
LimitNPROC=infinity
LimitCORE=infinity
# Uncomment TasksMax if your systemd version supports it.
# Only systemd 226 and above support this version.
#TasksMax=infinity
TimeoutStartSec=0
# set delegate yes so that systemd does not reset the cgroups of docker containers
Delegate=yes
# kill only the docker process, not all processes in the cgroup
KillMode=process

[Install]
WantedBy=multi-user.target
```

This will allow us to add extra flags via a drop-in file (mentioned above) by
placing a file containing the following in the `/etc/systemd/system/docker.service.d`
directory:

```conf
[Service]
ExecStart=
ExecStart=/usr/bin/dockerd --graph="/mnt/docker-data" --storage-driver=overlay
```

You can also set other environment variables in this file, for example, the
`HTTP_PROXY` environment variables described below.

To modify the ExecStart configuration, specify an empty configuration followed
by a new configuration as follows:

```conf
[Service]
ExecStart=
ExecStart=/usr/bin/dockerd --bip=172.17.42.1/16
```

If you fail to specify an empty configuration, Docker reports an error such as:

```conf
docker.service has more than one ExecStart= setting, which is only allowed for Type=oneshot services. Refusing.
```

### HTTP proxy

This example overrides the default `docker.service` file.

If you are behind an HTTP proxy server, for example in corporate settings,
you will need to add this configuration in the Docker systemd service file.

1.  Create a systemd drop-in directory for the docker service:

    ```bash
    $ mkdir -p /etc/systemd/system/docker.service.d
    ```

2.  Create a file called `/etc/systemd/system/docker.service.d/http-proxy.conf`
    that adds the `HTTP_PROXY` environment variable:

    ```conf
    [Service]
    Environment="HTTP_PROXY=http://proxy.example.com:80/"
    ```

3.  If you have internal Docker registries that you need to contact without
    proxying you can specify them via the `NO_PROXY` environment variable:

    ```conf
    Environment="HTTP_PROXY=http://proxy.example.com:80/" "NO_PROXY=localhost,127.0.0.1,docker-registry.somecorporation.com"
    ```

4.  Flush changes:

    ```bash
    $ sudo systemctl daemon-reload
    ```

5.  Verify that the configuration has been loaded:

    ```bash
    $ systemctl show --property=Environment docker
    Environment=HTTP_PROXY=http://proxy.example.com:80/
    ```
6.  Restart Docker:

    ```bash
    $ sudo systemctl restart docker
    ```

## Manually creating the systemd unit files

When installing the binary without a package, you may want
to integrate Docker with systemd. For this, simply install the two unit files
(service and socket) from [the github
repository](https://github.com/docker/docker/tree/master/contrib/init/systemd)
to `/etc/systemd/system`.