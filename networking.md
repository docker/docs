# Cross-Host Networking

For the Orca beta, enabling cross-host networking requires a manual step.
This document explains how to enable this feature after Orca has been
installed.


## Background

Libnetwork and cross-host networking is
officially supported as of Docker 1.9 (it is no longer an
[experimental](https://github.com/docker/docker/tree/master/experimental)
feature).  However, in 1.9, enabling this feature requires modifying
command line arguments for the Docker daemon.

The Key/Value store used by Orca, Swarm, and libnetwork is protected
by the Swarm TLS certificate chain managed by Orca.  This will prevent
unathorized access to the clusters configuration, and requires all
clients use a certificate signed by the Orca Swarm Root CA.  Admin users
certificate bundles, and the internal systems are signed by this CA.

During the install of Orca, TLS certificate files are placed on the
host filesystem of each engine in `/var/lib/docker/discovery_certs/`
to aid in the manual setup steps outlined below.


## Instructions

### 0) Install Orca

Install your controller, and join additional nodes as desired.
You can add additional nodes after enbabling cross-host networking,
however the manual steps described below must be performed **after**
the orca-bootstrap container has run `install` or `join` on the node.
The steps below must be performed on **every** node in your cluster.

### 1) Determine Orca IP

Before configuring each engine, determine the public facing IP address
(or hostname) of the primary orca system.  One example approach to
determine this address is by running the following on the host system
where the Orca controller is running:

```bash
ORCA_PUBLIC_IP=$(ip -o -4 route get 8.8.8.8 | cut -f8 -d' ')
echo ${ORCA_PUBLIC_IP}
```

**Note:  If you use a fully qualified hostnames instead of IP address, you
must make sure the TLS certificates includes that fully qualified name.
During install/join, use the `--san` flag, or `--interactive` mode.**

### 2) Determine configuration file

Different linux distributions and init sytems use different approaches
for configuring the Docker daemon.  The following lists a few examples:

* boot2docker based on tinycorelinux - Used by many docker-machine drivers
    * `/var/lib/boot2docker/profile` in the `EXTRA_ARGS` setting
* Ubuntu/Centos/RHEL + systemd
    * `/lib/systemd/system/docker.service` in the `ExecStart` setting

### 3) Add cluster flags to Docker Daemon

For each docker engine in your Orca/Swarm cluster **including the system running orca-controller**, you will have to update the docker daemon command line flags.

First determine the local engine's public IP:

```bash
LOCAL_ENGINE_IP=$(ip -o -4 route get 8.8.8.8 | cut -f8 -d' ')
echo ${LOCAL_ENGINE_IP}
```

Then edit the applicable configuration file (see step #2 above) with
the following settings.  **NOTE: You must replace the variables with
the actual values**

```bash
--cluster-advertise $LOCAL_ENGINE_IP:12376
--cluster-store etcd://$ORCA_PUBLIC_IP:12379
--cluster-store-opt kv.cacertfile=/var/lib/docker/discovery_certs/ca.pem
--cluster-store-opt kv.certfile=/var/lib/docker/discovery_certs/cert.pem
--cluster-store-opt kv.keyfile=/var/lib/docker/discovery_certs/key.pem
```

### 4) Restart the daemon

* boot2docker based on tinycorelinux - Used by many docker-machine drivers
    * **Restarting:** `sudo /etc/init.d/docker restart`
    * **Logs**: `tail -f /var/log/docker.log`
* Ubuntu + systemd
    * **Restarting:** `sudo systemdctl restart docker.service`
    * **Logs:** `sudo journalctl -fu docker.service`
