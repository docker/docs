+++
title = "Set up container networking with UCP"
description = "Docker Universal Control Plane"
[menu.main]
parent="mn_ucp"
+++

# Set up container networking with UCP

Beginning in release 1.9, the Docker Engine updated and expanded its networking
subsystem. Along with host and bridge networks, you can now create custom
networks that encompass multiple hosts running Docker Engine. This last feature
is known as multi-host networking.

#### About these installation instructions

These installation instructions were written using the Ubuntu 14.0.3 operating system. This means that the file paths and commands used in the instructions are specific to Ubuntu 14.0.3. If you are installing on another operating system, the steps are the same but your commands and paths may differ.

## Understand multi-host networking and UCP

You create a multi-host network using the Docker client or the UCP administration console.  Multi-host networks rely on the `overlay` network plugin driver. To create a network, using the CLI.

```
$ docker network create -d overlay my-custom-network
```

For the UCP beta, enabling multi-host networking is a manual process. You must
enable it after you install UCP on all your nodes. If you do not have
networking enabled, the Docker client returns the following error when you try
to create a network with the overlay driver:

```
$ docker network create -d overlay my-custom-network
Error response from daemon: failed to parse pool request for address space "GlobalDefault" pool "" subpool "": cannot find address space GlobalDefault (most likely the backing datastore is not configured)
```

If you attempt the same operation from UCP's web administration, you receive
the same error.

![Network error](images/network_gui_error.png)

This error returns because the networking features rely on a key-value store. In
a UCP environment, that key-value store is configured through UCP and protected
by the Swarm TLS certificate chain. To avoid this error, you need to manually
configure the Docker daemon to use UCP's key-value store in a secure manner.

This page explains how to configure the Docker Engine daemon startup options.
Once the daemon is configured and restarted, the `docker network` CLI and the
resources they create will use the Swarm TLS certificate chain managed by UCP.


## Enable multi-host networking

You'll do this procedure on all your UCP controller and nodes. Once your
configure and restart the Docker daemon, you'll have secure communication within
your cluster as you create custom multi-host networks on the controller or
nodes.


### Prerequisites

You must install UCP on your entire cluster (sever and nodes), before following
these instructions.  Make sure you have run on the `install` and `join` on each
node as appropriate. Then, enable mult-host networking on every node in your
cluster using these instructions.

UCP requires that all clients, including Docker Engine, use a Swarm TLS
certificate chain signed by the UCP Swarm Root CA. You configured these
certificates at bootstrap time either interactively or by passing the `--san`
(subject alternative sames) option to `install` or `join`.

To continue with this procedure, you need to know the SAN values you used on
each controller or node. Because you can pass a SAN either as an IP address or
fully-qualified hostname, make sure you know how to find these.

If you used a public IP address, log into the controller host and run these two
commands on the controller:

```bash
$ IP_ADDRESS=$(ip -o -4 route get 8.8.8.8 | cut -f8 -d' ')
$ echo ${IP_ADDRESS}
```

If your cluster is installed on a cloud provider, the public IP may not be the
same as the IP address returned by this command. Confirm through your cloud
provider's console or command line that this value is indeed the public IP. For
example, the AWS console shows these values to you:

![Open certs](images/ip_cloud_provider.png)

You can get also the SAN values of the controller by examining the certificate
through your browser. This would include a fully qualified hostname you used for
the controller. Each browser has a different way to view a website's certificate. To
do this on Chrome:

1. Open the browser to the UCP console.

2. In the address bar, click on the connection icon.

    ![Open certs](images/browser_cert_open.png)

    The browser opens displays the connection information. Depending on your Chrome version the dialog may be slightly different.

3. Click the **Certificate Information**.

4. Open the **Details** view and scroll down to the **Subject Alternative Name** section.

    ![SAN](images/browser_cert_san.png)

If you are using a private network, and don't know the fully-qualified DNS for a node, you can ask your network administrator.


### Configure and restart the daemon

If you followed the prerequisites, you should have a list of the SAN values you used with the UCP boostrapper to `install` the controller and `join` each node. With these values in hand, do the following:

1. Log into the host running the UCP controller.

2. Leave UCP processes running.

3. Determine the Docker daemon's startup configuration file.

    Each Linux distribution have different approaches for configuring daemon
    startup (init) options. On Centos/RedHat systems that rely systemd,
    the Docker daemon startup options are stored in the
    `/lib/systemd/system/docker.service` file. Ubuntu 14.04 stores these in the `/etc/init/docker.conf` file.

4. Open the configuration file with your favorite editor.

  **Ubuntu**:
        $ sudo vi /etc/init/docker.conf

  **Centos/Rehat**:
        $ sudo vi /lib/systemd/system/docker.service

5. Uncomment the `DOCKER_OPTS` line and add the following options.

        --cluster-advertise eth0:12376
        --cluster-store etcd://CONTROLLER_PUBLIC_IP_OR_DOMAIN:12379
        --cluster-store-opt kv.cacertfile=/var/lib/docker/discovery_certs/ca.pem
        --cluster-store-opt kv.certfile=/var/lib/docker/discovery_certs/cert.pem
        --cluster-store-opt kv.keyfile=/var/lib/docker/discovery_certs/key.pem

  Replace `CONTROLLER_PUBLIC_IP_OR_DOMAIN` with the IP address of the UCP
  controller. Use `ifconfig` to ensure the host you are installing on is
  accessible over `eth0`, on your host the systems with multiple ether
  interfaces this value might differ. When you are done, the line should look
  similar to the following:

        DOCKER_OPTS="--dns 8.8.8.8 --dns 8.8.4.4 --cluster-advertise eth0:12376 --cluster-store etcd://52.70.188.239:12379 --cluster-store-opt kv.cacertfile=/var/lib/docker/discovery_certs/ca.pem --cluster-store-opt kv.certfile=/var/lib/docker/discovery_certs/cert.pem --cluster-store-opt kv.keyfile=/var/lib/docker/discovery_certs/key.pem"

6. Save and close the `/etc/init/docker.conf` file.

6. Restart the Docker daemon.

  **Ubuntu**:
        $ sudo docker restart

  **Centos/RedHat**:
        $ sudo systemdctl restart docker.service

7. Review the Docker logs to check the restart.

  **Ubuntu**:
        $ sudo tail -f /var/log/upstart/docker.log

  **Centos/RedHat**:
        $ sudo sudo journalctl -fu docker.service

8. Verify that you can create and delete a custom network.

        $ docker network create -d overlay my-custom-network
        $ docker network ls
        $ docker network rm my-custom-network

9. Repeat steps 1-8 on the remaining nodes in your cluster.


## Troubleshooting the daemon configuration

If you have trouble, try these troubleshooting measures:

* Review the daemon logs to ensure the daemon was started.
* Add the `-D` (debug) to the Docker daemon start options.
* Check your daemon configuration to ensure that `--cluster-advertise eth0:12376` is set properly.  
* Check your daemon configuration `--cluster-store` options is point to the
key-store `etcd://CONTROLLER_PUBLIC_IP_OR_DOMAIN:12379` on the UCP controller.
* Make sure the controller is accessible over the network, for example `ping CONTROLLER_PUBLIC_IP_OR_DOMAIN`.
A ping requires that inbound ICMP requests are allowed on the controller.
* Stop the daemon and start it manually from the command line.

      sudo /usr/bin/docker daemon -D --cluster-advertise eth0:12376 --cluster-store etcd://CONTROLLER_PUBLIC_IP_OR_DOMAIN:12379 --cluster-store-opt kv.cacertfile=/var/lib/docker/discovery_certs/ca.pem --cluster-store-opt kv.certfile=/var/lib/docker/discovery_certs/cert.pem --cluster-store-opt kv.keyfile=/var/lib/docker/discovery_certs/key.pem

Remember, you'll need to restart the daemon each time you change the start options.
