---
description: Permission requirements for Docker Desktop for Mac
keywords: Docker Desktop, mac, security, install
title: Docker Desktop permission requirements for Mac
redirect_from:
- /docker-for-mac/privileged-helper/
---

This page contains information about the permission requirements for running and installing Docker Desktop on Mac, the functionality of the privileged helper process `com.docker.vmnetd` and the reasoning behind this approach, as well as clarification about running containers as `root` as opposed to having `root` access on the host.

## Permission requirements

In the default set up flow, Docker Desktop for Mac does not require root privileges for installation but does require `root` access to be granted on the first run. The first time that Docker Desktop is launched the user gets an admin prompt to grant permissions for a privileged helper service to be installed. For subsequent runs, no `root` privileges are required.

The reason for this is that Docker Desktop needs to perform a limited set of privileged operations using the privileged helper process `com.docker.vmnetd`. This approach allows, following the principle of least privilege, `root` access to be used only for the operations for which it is absolutely necessary, while still being able to use Docker Desktop as an unprivileged user.

From version 4.11 of Docker Desktop for Mac it will be possible to avoid running the privileged helper service in the background by using the `--user` flag on the [install command](../install/mac-install.md#install-from-the-command-line). This will result in `com.docker.vmnetd` being used for set up during installation and then disabled at runtime. In this case, the user will not be prompted to grant root privileges on the first run of Docker Desktop. Specifically, the `--user` flag:
- Uninstalls the previous `com.docker.vmnetd` if present
- Sets up symlinks for the user
- Ensures that `localhost` and `kubernetes.docker.internal` are present in `/etc/hosts`

This approach has the following limitations:
- Docker Desktop can only be run by one user account per machine, namely the one specified in the `-–user` flag.
- Ports 1-79 are blocked. the containers will run but the port won’t be exposed on the host.
- Spindump diagnostics for fine grained CPU utilization are not gathered.

## Privileged Helper

The privileged helper is started by `launchd` and runs in the background unless it is disabled at runtime as previously described. The Docker Desktop backend communicates with it over the UNIX domain socket `/var/run/com.docker.vmnetd.sock`. The functionalities it performs are: 
- Installing and uninstalling symlinks in `/usr/local/bin`. This ensures the `docker` CLI is on the user’s PATH without having to reconfigure shells, log out then log back in for example.
- Binding privileged ports that are less than 1024. The so-called "privileged ports" have not generally been used as a security boundary, however OSes still prevent unprivileged processes from binding them which breaks commands like `docker run -p 80:80 nginx`
- Ensuring `localhost` and `kubernetes.docker.internal` are defined in `/etc/hosts`. Some old macOS installs did not have `localhost` in `/etc/hosts`, which caused Docker to fail. Defining the DNS name `kubernetes.docker.internal` allows us to share Kubernetes contexts with containers.
- Securely caching the Registry Access Management policy which is read-only for the developer.
- Performing some diagnostic actions, in particular gathering a performance trace of Docker itself.
- Uninstalling the privileged helper.

## Containers running as root within the Linux VM

The Docker daemon and containers run in a lightweight Linux VM managed by Docker. This means that although containers run by default as `root`, this does not grant `root` access to the Mac host machine. The Linux VM serves as a security boundary and limits what resources can be accessed from the host. Any directories from the host bind mounted into Docker containers still retain their original permissions.


