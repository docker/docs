+++
title = "Quickstart Guide"
description = "Docker Universal Control Plane"
[menu.ucp]
weight="-98"
+++


# Docker UCP Quickstart Guide

These instructions explain how to install Docker Universal Control Plane (UCP). A UCP installation consists of an UCP controller and one or more nodes. The same machine can serve as both the controller and the node. These instructions show you how to install both a host and a node. It contains the following sections:

- [Plan your installation](#plan-your-installation)
- [Step 1: Verify you have the prerequisites](#step-1-verify-you-have-the-prerequisites)
- [Step 2: Configure your network for UCP](#step-2-configure-your-network-for-ucp)
- [Step 3: Install Docker CS Engine 1.9](#step-3-install-docker-cs-engine-1-9)
- [Step 4: (optional) Create user-named volumes](#step-4-optional-create-user-named-volumes)
- [Step 5: Install the UCP controller](#step-5-install-the-ucp-controller)
- [Step 6: (optional) Add a controller replica to the UCP cluster](#step-6-optional-add-a-controller-replica-to-the-ucp-cluster)
- [Step 7: Add an Engine node to the UCP cluster](#step-7-add-an-engine-node-to-the-ucp-cluster)
- [Step 8: Set up certs for the Docker CLI](#step-8-set-up-certs-for-the-docker-cli)
- [Uninstall](#uninstall)
- [Block Mixpanel analytics](#block-mixpanel-analytics)
- [Installing with your own certificates](#installing-with-your-own-certificates)
- [Where to go next](#where-to-go-next)

## Plan your installation

The UCP installation consists of running the `ucp bootstrapper` image using the
Docker Engine CLI. The bootstrapper has subcommands to `install` a controller or `join` a node to a UCP controller.  During an install, the bootstrapper:

* configures data volumes on the controller
* generates a certificate chain on the controller
* creates a Swarm cluster
* loads and launches the appropriate images

You can use the bootstrapper as an interactive script or by passing command-line
options. Regardless of how you use the bootstrapper, the installer supplies some
quick default options for both data volumes and the certificate authority (CA).
### Default versus the custom installation options

The first time you install, you should build a sandbox environment using
`bootstrapper` defaults. After installing and using this sandbox environment,
you can uninstall it and try a custom installation. In a custom installation you
can:

* use the high availability feature
* customize the port used by the Swarm manager
* create your own data volumes
* use your own certs

This install documentation describes the default installation and the
customization steps.  Customize steps are identified with the keyword
(optional). Make sure you skip these steps when doing the default installation
in your sandbox.

### IP addresses and fully-qualified domain names

When you bootstrap a controller or node, you must supply a host address either
interactively or using the `--host-address` option. The host address can be an accessible IP address and/or fully-qualified domain name.  

If you are using a cloud provider such as AWS or Digital Ocean, you may need to
allocate a private network for your UCP installation. You can use this network
as long as the controller and nodes can communicate via
their private IPs. If the private IPs do not support communication among the
UCP cluster, using public IPs or full-qualified domain names are required. For
more information about what ports and protocols are required see [Step 2: Configure your network for UCP](#step-2-configure-your-network-for-ucp).

### Subject alternative names (SANs)

Further, UCP requires that all clients, including the Docker Engine, use a Swarm
TLS certificate chain signed by the UCP Swarm Root CA. You build the certificate
by passing the `--san` (subject alternative names or SANs) values to the
boostrapper's `install` or `join`. A SAN value can be the pubic IP address
and/or fully-qualified domain name.

For the controller and each node, you must specify at least one SAN; you can
specify more.

If you are using a cloud provider and specified private IPs for the host address
values, consider whether you need to access your cluster through a public
network as well as the private network space. If the answer is yes, your SAN
values should contain both the public IPs or full-qualified hostnames and the private network IPs.

### Mixpanel analytics

The UCP BETA program makes use of Mixpanel to collect analytics. This feature
collects data on your usage of UCP and returns it to Docker. The information is
entirely anonymous and does not identify your Company or users. Currently, you
cannot turn the collection off but you can block the outgoing messaging. Later
in this documentation [Block Mixpanel analytics](#block-mixpanel-analytics)
explains how.

### About these installation instructions

These installation instructions were written using the Ubuntu 14.0.3 operating
system. The file paths and commands used in the instructions are specific to
Ubuntu 14.0.3. If you are installing on another operating system, the steps are
the same but your commands and paths may differ.

The Ubuntu system used to write these instructions was an AWS EC2 instance. The
instance was running in a in an AWS virtual private network. Elastic IPs were
configured for all the hosts.

## Step 1: Verify you have the prerequisites

You can install UCP on your network or on a cloud provider such AWS or Digital
Ocean. To install, the controller and the nodes require a minimum of 1.50 GB of memory. You can run any of these supported operating systems:

* RHEL 7.0, 7.1
* Ubuntu 14.04 LTS
* CentOS 7.1

Your system must have a 3.16.0 kernel or higher. If you don't have the proper kernel installed, the UCP bootstrapper returns this error.

```
INFO[0000] Verifying your system is compatible with UCP
FATA[0000] Your kernel version 3.13.0 is too old.  UCP requires at least version 3.16.0 for all features to work.  To proceed with an old kernel use the '--old-kernel' flag
```

If you proceed with the '--old-kernel' flag, you cannot use the Docker custom networking features.

Installing Docker UCP requires that you first install the CS Docker Engine v1.9
on both the controller and the nodes. The CS Engine can be installed manually or
from an image if your cloud provider support its. These instructions assume you
are installing both UCP and Engine manually.

Finally, installing UCP requires you to pull an image from the Docker Hub. If
you don't already have a Docker Hub account with access to the UCP Beta. Make
sure you [create an account first](https://hub.docker.com/). Once you have a Hub
account, [request access to the UCP BETA
program](https://goto.docker.com/try-universal-control-plane.html).

## Step 2: Configure your network for UCP

UCP includes Docker Swarm as part of its installation. So, you don't need to
install Docker Swarm. You do need to ensure that the UCP controller and nodes
can communicate across your network. Configure your network making sure to open
the following ports:

| Port             | Description     |
|------------------|-----------------|
| `443`            | UCP controller     |
| `2376`           | Swarm manager   |
| `12376`          | Engine proxy    |
| `12379`, `12380` | Key Value store |
| `12381`          | Swarm CA service|
| `12382`          | UCP CA service |

All the communication among the controller, nodes, and key value store is
protected by mutual TLS. The UCP installation of Swarm provides and configures
TLS for you automatically.

Finally, you can specify a different port for the Swarm manager if you need to.
Using a different  port is a customization. These instructions assume you are using the default `2376` port.

## Step 3: Install Docker CS Engine 1.9

The BETA program requires that you install the Docker CS Engine 1.9.0 or above.
Follow the instructions for your particular operating system and ensure
you are pointing at the proper repo.

Install the Docker CS Engine on both the controller node and each member node.

### RHEL 7.0, 7.1 and CentOS 7.1

Use the detailed [Red Hat Linux installation
instructions](https://docs.docker.com/docker-trusted-registry/install/install-csengine/#centos-7-1-rhel-7-0-7-1-yum-based-systems)

### Ubuntu 14.04 LTS

Use the [detailed Ubuntu installation
instructions](https://docs.docker.com/docker-trusted-registry/install/install-csengine/#install-on-ubuntu-14-04-lts)


## Step 4: (optional) Create user-named volumes

UCP uses named volumes for persistence of user data.  By default, the `ucp`
bootstrapper creates for you. It uses the default volume driver and flags. The
first time you install, you should skip this step and try it later. Later, try
an install where your take the option to use custom volume driver and create
your own volumes.

If you choose this option, create your volumes prior to installing UCP. The volumes UCP requires are:

| Volume name             | Data                                                                                 |
|-------------------------|--------------------------------------------------------------------------------------|
| `ucp-root-ca`          | The certificate and key for the UCP root CA. Do not create this volume if you are using your own certificates.                                      |
| `ucp-swarm-root-ca`    | The certificate and key for the Swarm root CA.                                       |
| `ucp-server-certs`     | The controller certificates for the UCP controllers web server.                                     |
| `ucp-swarm-node-certs` | The Swarm certificates for the current node (repeated on every node in the cluster). |
| `ucp-swarm-kv-certs`   | The Swarm KV client certificates for the current node (repeated on every node in the cluster). |
| `ucp-swarm-controller-certs` | The UCP Controller Swarm client certificates for the current node. |
| `ucp-kv`               | Key value store persistence.                                                         |

## Step 5: Install the UCP controller

In this step you install the UCP controller. The controller includes a running Swarm manager and node as well. Use the following command to pull the bootstrapper image and review the `install` options:

```bash
docker run --rm -it docker/ucp install --help
```


When you install, the script prompts you for the following information:

* a password to use for the UCP `admin` account
* your Docker Hub username/password/email
*  at lease one SAN value which is the accessible IP address or fully-qualified domain name for the controller node

When you have the information you'll be prompted for, do the following to
install:

1. Log into the system where you mean to install the UCP controller.

    If you are installing on a cloud provider such as AWS, make sure the instance has a public IP or hostname.

    ![Open certs](../images/ip_cloud_provider.png)

2. Run the `ucp` bootstrapper.

        $ docker run --rm -it \
          -v /var/run/docker.sock:/var/run/docker.sock \
          --name ucp \
          docker/ucp \
          install -i

    The bootstrapper pulls several images and prompts you for the installation values it needs. When it completes, the bootstrapper prompts you to login into the UCP GUI.

        INFO[0053] Login to UCP at https://10.0.0.32:443

3. Enter the address into your browser to view the UCP login screen.

    Your browser may warn you about the connection. The warning appears because
    the UCP certificate was issued by a built-in certificate authority. Your
    actions with the install actually created the certificate. If you are
    concerned, the certificate's fingerprint is displayed during install and you
    can compare it.  

4. Use the Advanced link to proceed to UCP.

    The login screen displays.

    ![](../images/login.png)

5. Enter `admin` for the username along with the password you provided to the bootstrapper.

    If you didn't enter an admin password, the default password is `orca` After
    you enter the correct credentials, the UCP dashboard displays.

    ![](../images/dashboard.png)

    The dashboard shows a single node, your controller node.

## Step 6: (optional) Add a controller replica to the UCP cluster

In this optional step, you configure support for UCP's high-availability
feature. You do this by adding one or more UCP *replicas* using the
bootstrapper's `ucp join` subcommand.  The first time you install, you should
skip this optional step and try it later. Later, try an install where you
configure high-availability.

When adding nodes to your cluster, you decide which nodes you to use as
*replicas* and which nodes are simply for extra capacity.  A
replica is a node in your cluster that can act as an additional UCP controller.
Should the primary controller fail, a replica can take over the controller role
for the cluster.  If you are trying out the optional HA deployment:

* Configure the controller and replicas before adding additional Engine nodes.
* Configure a minimum of two replicas in addition to the controller.

Repeat the install for each node you want to add. Use the following command to pull the bootstrapper image and review the `join` options:


```bash
docker run --rm -it docker/ucp join --help
```

The bootstrapper prompts you for the following information:

* the URL of the UCP controller, for example `https://52.70.188.239`
* the username/password of an UCP administrator account
* your Docker Hub username/password/email
* at least one SAN value which is an accessible IP address or fully-qualified domain name for node

When you have the information you'll be prompted for, do the following to install:

1. Log into the host where you mean to install the node.

2. Run the `ucp` bootstrapper.

        $ docker run --rm -it -v /var/run/docker.sock:/var/run/docker.sock --name ucp docker/ucp join --replica -i

    The bootstrapper pulls several images and prompts you for the installation values it needs. When it completes, the bootstrapper notifies you that it is starting swarm.

        INFO[0005] Verifying your system is compatible with UCP
        INFO[0011] Sending add host request to UCP server      
        INFO[0011] Starting local swarm containers  

3. Repeat steps 1 thru 3 on the other replicas.

    You should configure a minimum of 3 controllers configured, a primary and
    two replicas. Never run a cluster with only the primary controller and a
    single replica.  

4. Login into UCP with your browser and navigate to the **NODES** page.

    Simply clicking on the nodes from the Dashboard takes you to the page. The page should display your new nodes.

      ![](../images/nodes.png)


## Step 7: Add an Engine node to the UCP cluster

In this step, you install one or more UCP nodes using the `ucp join` subcommand. Repeat the install for each node you want to add. Use the following command to pull the bootstrapper image and review the `join` options:

```bash
docker run --rm -it docker/ucp join --help
```

The bootstrapper prompts you for the following information:

* the URL of the UCP controller, for example `https://52.70.188.239`
* the username/password of an UCP administrator account
* your Docker Hub username/password/email
* at least one SAN value which is the actual external, publically-accessible IP address or fully-qualified domain name for node

When you have the information you'll be prompted for, do the following to install:

1. Log into the system where you mean to install the node.

2. Run the `ucp` bootstrapper.

        $ docker run --rm -it -v /var/run/docker.sock:/var/run/docker.sock --name ucp docker/ucp join -i

    The bootstrapper pulls several images and prompts you for the installation values it needs. When it completes, the bootstrapper notifies you that it is starting swarm.

        INFO[0005] Verifying your system is compatible with UCP
        INFO[0011] Sending add host request to UCP server      
        INFO[0011] Starting local swarm containers  

3. Repeat steps 1 thru 2 on the other nodes.

4. Login into UCP with your browser and navigate to the **NODES** page.

    Simply clicking on the nodes from the Dashboard takes you to the page. The page should display your new nodes.

      ![](../images/nodes.png)

## Step 8: Set up certs for the Docker CLI

Once you install UCP on a machine, it is a good idea to download a client bundle.  The bundle contains the certificates a user needs to run the Docker Engine command line client (`docker`) against the UCP controller and nodes.

You can download the bundle from the UCP interface or using `curl` command. The
bundle is in a `.zip` package.  You need `zip` or similar to unzip the file. If you plan to use `curl` you also need JQuery. This is used to pass the `curl`
command an authorization token. Of course, you need to have `curl` installed as well.

### Download the bundle from the UCP interface

1. If you haven't already done so, log into UCP.

2. Choose **ADMIN > Profile** from the right-hand menu.

    Any user can download their certificates. So, if you were logged in under a user name such as `davey` the path to download bundle is **davey > Profile**. Since you are logged ins as `ADMIN`, the path is `ADMIN`.

3. Click **Create Client Bundle**.

    The browser downloads the `ucp-bundle-admin.zip` file.

### Download the bundle with curl

1. Log into a machine with network access to the UCP controller.

    You might log into the controller itself. You could also log into any arbitrary machine able to `ping` the controller.

2. Install the prerequisite `curl`, `zip`, `jq` (JQuery) packages if you haven't already.

    On Ubuntu, the installation looks like this:

          $ sudo apt-get install zip curl jq
          Reading package lists... Done
          Building dependency tree       
          Reading state information... Done
          The following extra packages will be installed:
            libcurl3
          The following NEW packages will be installed:
            jq zip
          The following packages will be upgraded:
            curl libcurl3
            ----output snipped----

    To curl the bundle, you must export your user security token from the UCP controller. You do this in the next step.

3. Create an environment variable to hold your user security token.

		AUTHTOKEN=$(curl -sk -d '{"username":"admin","password":"<password>"}' https://<ducp-0 IP>/auth/login | jq -r .auth_token)

4. Curl the client bundle down to your node.

		    $ curl -k -H "X-Access-Token:admin:$AUTHTOKEN" https://<ducp-0 IP>/api/clientbundle -o bundle.zip

    The browser downloads a `bundle.zip` file.

### Install the certificate bundle

Once you download the bundle, you can unzip and use it.

1. Make sure you have `zip` installed.

        $ which unzip
        /usr/bin/unzip

    If you don't, install it before continuing.

2. Open the folder containing the bundle file.

4. Unzip the file to reveal its contents.

        ucp-bundle
        ├── ca.pem
        ├── cert.pem
        ├── cert.pub
        ├── env.sh
        └── key.pem

5.  Set up your environment by sourcing the `env.sh` file.

        $ source env.sh

6.  Use the `docker info` command to get the location of the Swarm managers and engines.

        $ docker info
        Containers: 9
        Images: 9
        Role: primary
        Strategy: spread
        Filters: health, port, dependency, affinity, constraint
        Nodes: 1
         node1: 192.168.122.7:12376
          └ Containers: 9
          └ Reserved CPUs: 0 / 1
          └ Reserved Memory: 0 B / 2.054 GiB
          └ Labels: executiondriver=native-0.2, kernelversion=4.0.9-boot2docker, operatingsystem=Boot2Docker 1.8.1 (TCL 6.3); master : eb5571f - Thu Sep  3 22:18:54 UTC 2015, provider=kvm, storagedriver=aufs
        Cluster Managers: 1
         192.168.122.7: Healthy
          └ Orca Controller: https://192.168.122.7
          └ Swarm Manager: tcp://192.168.122.7:3376
          └ KV: etcd://192.168.122.7:12379
        CPUs: 1
        Total Memory: 2.054 GiB
        Name: node1
        ID: PNLT:MFCO:DDWL:MSLF:YVHU:35Z3:66KM:DFZM:OPBK:D4BQ:EKNT:6DXA
        Labels:
         com.docker.ucp.license_key=unlicensed
         com.docker.ucp.license_max_engines=0
         com.docker.ucp.license_expires=EXPIRED

### Client Bundles on Externally Managed CA Configuration                                           

If UCP is configured with an external CA, it will be unable to sign client bundles for non-admin users automatically. It is still possible to manually issue certificates signed by the CA that UCP users can use to interact with UCP via the CLI.

Generate an 2048-bit RSA private key.

```
openssl genrsa -out key.pem 2048
```

Generate a Certificate Signing Request (CSR).  The output `cert.csr` should be provided to your organization's CA owner to be signed, with a minimum of client authentication usage.

```
openssl req -new -sha256 -key key.pem -out cert.csr
```

Your CA owner will sign the CSR, and provide `cert.pem` and `ca.pem` files.

Extract the public key from the signed certificate:

```
openssl x509 -pubkey -noout -in cert.pem  > cert.pub
```

The contents of cert.pub will then need to be added to your profile.  You can add this in the UI by clicking the User Menu in the top right corner, and select profile.

Once you are on the User Profile screen, click the "Add an Existing Public Key" button and provide the contents of cert.pub, along with a memorable label for this bundle.

Now that you have linked the public key to you account, the next step is to configure your CLI. To configure your CLI to use the certificate bundle that you have generated, you will need to export the following environment variables:

```
export DOCKER_TLS_VERIFY=1
export DOCKER_CERT_PATH=$(pwd)
export DOCKER_HOST=tcp://<ucp-hostname>:443
```

## Uninstall

The bootstrapper can also uninstall UCP from the controller and the nodes. The uninstall process will not remove any other containers that are running, except those recognized to be part of UCP. To see the uninstall options before you uninstall, use the following:

```bash
docker run --rm -it docker/ucp uninstall --help
```

To uninstall, do the following:

1. Log into the node you want to remove UCP from.

2. Enter the following command to uninstall:

        $ docker run --rm -it -v /var/run/docker.sock:/var/run/docker.sock --name ucp docker/ucp uninstall

3. Repeat the uninstall on each node making sure to save the controller till last.

## Block Mixpanel analytics

To block the outflow of Mixplanel analytic data to Docker, do the following:

1. Log into the system running the UCP controller.

2. Add a rule to drop the forward to port 80.

        $ iptables -I FORWARD -p tcp --dport 80 -j DROP

Reboots unset this iptables chain, so it is a good idea to add this command to the controller's startup script.

## Installing with your own certificates

UCP uses two separate root CAs for access control - one for Swarm, and one for
the UCP controller itself.  The dual root certificates supply differentiation
between the Docker remote API access to UCP vs. Swarm.  Unlike Docker Engine or
Docker Swarm, UCP implements ACL and audit logging on a per-user basis.  Swarm
and the Engine proxies trust only the Swarm Root CA, while the UCP controller
trusts both Root CAs.  Admins can access UCP, Swarm and the engines while
normal users are only granted access to UCP.

UCP v1.0 supports user provided externally signed certificates
for the UCP controller.  This cert is used by UCP's main management web UI
and the Docker remote API. The remote API is visible to the Docker CLI. In this release, the Swarm Root CA is always managed by UCP.

The external UCP Root CA model supports customers managing their own CA, or
purchasing certs from a commercial CA.  When operating in this mode, UCP can
not generate regular user certificates, as those must be managed and signed
externally, however admin account certs can be generated as they are signed by
the internal Swarm Root CA.  Normal user accounts should be signed by the same
external Root CA (or a trusted intermediary), and the public keys manually added
through the UI.

The first time you install, we recommend you skip user-supplied certs and use
the default certificates instead. The default TLS certificate files are placed
on the host filesystem of each Docker Engine in
`/var/lib/docker/discovery_certs/`. Later, do a second install and try the
option to use your own certs.


### Configure user-supplied Certificates

To install UCP with your own external root CA, you create a named volume called
**ucp-server-certs** on the same system where you plan to install the UCP
controller.

1. Log into the machine where you intend to install UCP.

2. If you haven't already done so, create a named volume called **ucp-server-certs**.

3. Ensure the volume's top-level directory contains these files:

    <table>
    <tr>
      <th>File</th>
      <th>Description</th>
    </tr>
    <tr>
      <td><code>ca.pem</code></td>
      <td>Your Root CA Certificate chain (including any intermediaries).</td>
    </tr>
    <tr>
      <td><code>cert.pem</code></td>
      <td>Your signed UCP controller cert.</td>
    </tr>
    <tr>
      <td><code>key.pem</code></td>
      <td>Your UCP controller private key.</td>
    </tr>
    </table>

4. Follow "Step 5" above to install UCP but pass in an additional `--external-ucp-ca` option to the bootstrapper, for example:

        docker run --rm -it \
          -v /var/run/docker.sock:/var/run/docker.sock \
          ...snip...
          install -i --external-ucp-ca


## Where to go next

To learn more you can also investigate:

* Read more [about the product](https://www.docker.com/universal-control-plane).
* Visit the UCP forum to [ask questions and get answers](https://forums.docker.com/c/commercial-products/ucpbeta).
* Try the [UCP hands-on lab](https://github.com/docker/ucp_lab).
* [How to use the Docker Client](http://docs.docker.com/reference/commandline/cli/)
* [An overview of Docker Swarm](http://docs.docker.com/swarm/)
