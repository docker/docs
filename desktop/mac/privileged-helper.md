---
description: Docker Desktop Privileged Helper on Mac
keywords: Docker Desktop, mac, security, install
title: Docker Desktop Privileged Helper on Mac
redirect_from:
- /docker-for-mac/privileged-helper/
---

## Permission requirements

While Docker Desktop does not generally require the user running or installing it to have `root` privileges, in the default use case it needs root access to be granted on the first run. The first time Docker Desktop is launched the user gets an admin prompt to grant permissions for a privileged helper service to be installed. For subsequent runs, no `root` privileges are required. 

The reason for this is that Docker Desktop needs to perform a limited set of privileged operations using the privileged helper process `com.docker.vmnetd`. This approach allows, following the principle of least privilege, `root` access to be used only for the operations for which it is absolutely necessary, while still being able to use Docker Desktop as an unprivileged user.

From version 4.11, it will be possible to avoid running a privileged service in the background by using `com.docker.vmnetd` for setup during installation and disabling it at runtime. In this case the user will not be prompted on the first run. Administrators would be able to do that by using the `–user` flag on the [install command](install.md#install-from-the-command-line) which would:
- Uninstall the previous `com.docker.vmnetd` if present
- Set up symlinks for the user
- Ensure that `localhost` and `kubernetes.docker.internal` are present in `/etc/hosts`

This approach will have certain limitations:
- Docker Desktop would only be able to be run by one user account per machine, namely the one specified in the `–user` flag.
- Ports 1-79 would be blocked - the containers would run but the port won’t be exposed on the host.
- Spindump diagnostics for fine grained CPU utilization would not be gathered.

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


