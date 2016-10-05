<!--[metadata]>
+++
draft=true
aliases = ["/ucp/evaluation-install/"]
title = "Install DDC in a sandbox"
description = "Evaluation installation"
keywords = ["UCP, DTR, install, sandbox, evaluation, free trial"]
[menu.main]
parent="mn_ucp"
identifier="ucp_evaluate_sandbox"
weight=10
+++
<![end-metadata]-->


# Install DDC in a sandbox for evaluation

This page introduces Docker Datacenter: a combination of Docker Universal
Control Plane (UCP) and Docker Trusted Registry (DTR), and walks you through
installing it on a local (non-production) host or sandbox. Once you've installed, we'll also give you a guided tour so you can evaluate its features.

The instructions here are for a sandbox installation on Mac OS X or Windows
systems. If you're an experienced Linux user, or if you want more detailed
technical information, you might want to skip this evaluation and go directly to
[Plan a production installation](installation/plan-production-install.md) and
then to [Install UCP for production](installation/install-production.md).

> **Note**: This evaluation installs using the open source software version of
Docker Engine inside of a VirtualBox VM which runs the small-footprint
`boot2docker.iso` Linux. This configuration is **not** a production
configuration.

## Introduction: About this example

In this tutorial, we'll use Docker's provisioning tool - Docker Machine - to
create two virtual hosts. These two hosts are VirtualBox VMs running a small
footprint Linux image called `boot2docker.iso`, with the open source version of
Docker Engine installed.

![Docker Machine setup](images/explain.png)

A UCP installation consists of an UCP controller and one or more hosts. We'll
install UCP on one host, then join the second node to UCP as a swarm member. The
two VMs create a simple swarm cluster with one controller, which by default
secures the cluster via self-signed TLS certificates.

![Sandbox](images/sandbox.png)

 DDC's second component is DTR, which must be installed on a host that's a member of the UCP swarm. So next, we'll then install DTR on that second node.

Once you've installed UCP and DTR you'll work through a tutorial to deploy a
container through UCP, and explore the user interface.

>**Note**: The command examples in this page were tested for a Mac OS X environment. If you are in another, you may need to change the commands to use the correct ones for you environment.

## Prerequisites

This example requires that you have:

* [Docker Toolbox installed](/toolbox/overview.md) (contains Docker Machine and Docker Engine)
* A free Docker ID account

## Step 1: Provision hosts with Engine

In this step, you'll provision two VMs for your UCP sandbox installation. One
will run UCP and one will be used to run containers, so the host specifications
will be slightly different.

In a production environment you would use enterprise-grade hosts instead
of local VMs. These nodes could be on your company's private network or
in the cloud. You would also use the Commercially Supported (CS Engine) version of Docker Engine required by UCP.

Set up the nodes for your evaluation:

1. Open a terminal on your computer.

2. Use Docker Machine to list any VMs in VirtualBox.

    ```none
    $ docker-machine ls
    NAME         ACTIVE   DRIVER       STATE     URL                         SWARM
    default    *        virtualbox   Running   tcp://192.168.99.100:2376
    ```

3. Create a VM named `node1` using the following command.

    ```
    $ docker-machine create -d virtualbox \
    --virtualbox-memory "2000" \
    --virtualbox-disk-size "5000" node1
    ```

    When you create your virtual host you specify the memory and disk size
    options. UCP requires a minimum of 3.00 GB disk space and runs best with a
    minimum of 2 GB of memory.

4. Create a VM named `node2` using the command below.

    ```none
    $ docker-machine create -d virtualbox \
    --virtualbox-memory "2000" node2
    ```

5. Use the `docker-machine ls` command to list your hosts.

    ```
    $ docker-machine ls
    NAME      ACTIVE   DRIVER       STATE     URL                         SWARM   DOCKER    ERRORS
    default   -        virtualbox   Stopped                                       Unknown
    node1     -        virtualbox   Running   tcp://192.168.99.100:2376           v1.12.1
    node2     -        virtualbox   Running   tcp://192.168.99.101:2376           v1.12.1

    ```
    At this point, both nodes are in the `Running` state and ready for UCP installation.

## About the ucp tool

To install UCP, you'll use the Docker CLI to pull and run the `docker/ucp`
image, which contains a bootstrapper tool, which is designed to make UCP easier
to install than many enterprise-grade applications. The `ucp` tool runs `docker run` commands to `install` a UCP controller or `join` a node to a UCP
controller.

The general format of these commands are a `docker run --rm -it docker/ucp` with
one or more subcommands, and you'll find them later in this document. For the
tutorial purposes, we use the `-i` options for "interactive" install mode, but
you can run them unattended in production.

Regardless of how you use the `docker/ucp` tool, the default install supplies
default options for both data volumes and the certificate authority (CA). In a
production installation you can also optionally:

* use the high availability feature
* customize the port used by the UCP web application
* customize the port used by the Swarm manager
* create your own data volumes
* use your own TLS certificates

You can learn more about these when you <a
href="https://docs.docker.com/ucp/plan-production-install/" target="_blank">Plan
a production installation</a>.

## Step 2. Install the UCP controller

In this step, you install the UCP controller on the `node1` you provisioned
earlier. A controller serves the UCP application and runs the processes that
manage an installation's Docker objects.

In a production installation, a system administrator can set up UCP's high
availability feature, which allows you to designate several nodes as controller
replicas. This way if one controller fails, a replica node is ready to take its
place.

For this sandbox installation, we don't need high availability, so a single
host for the controller works fine.

1. Open a terminal on your computer if you don't have one open already.

2. Connect the terminal environment to the `node1` you created.

    a. Use the `docker-machine env` command to get the settings.

      ```none
      $ docker-machine env node1
      export DOCKER_TLS_VERIFY="1"
      export DOCKER_HOST="tcp://192.168.99.100:2376"
      export DOCKER_CERT_PATH="/Users/ldr/.docker/machine/machines/node1"
      export DOCKER_MACHINE_NAME="node1"
      # Run this command to configure your shell:
      # eval $(docker-machine env node1)
      ```

    b. Run the `eval` command found in the final line to set your environment.

      ````
      $ eval $(docker-machine env node1)
      ````

      Running this `eval` command sends the `docker` commands in the following
      steps to the Docker Engine on on `node1`.

    c. Verify that `node1` is the active environment.

    You can do this by running `docker-machine ls` and checking that there is an `*` (asterisk) in the `ACTIVE` field next to `node1`.

3. Start the `ucp` tool to install interactively.

    ```none
    $ docker run --rm -it \
    -v /var/run/docker.sock:/var/run/docker.sock \
    --name ucp docker/ucp install -i \
    --swarm-port 3376 --host-address $(docker-machine ip node1)
    ```

    > **Note**: If you are on a Windows system, your shell won't be able to
    resolve the `$(docker-machine ip node1)` variable. Instead, edit the command
    supplied to replace it with the actual IP address.

    The first time you run the `ucp` tool, the `docker run` command pulls the
    UCP bootstrapper image from Docker Cloud. The tool downloads the packages it
    needs, and verifies that your system will support a UCP installation.

4. Enter a password for UCP when prompted, and then confirm it.

    The system prompts you for Subject alternative names (SANs). In this
    sandbox, you've already provided the IP address and the `ucp` tool
    discovered this for you and shows it in the controller list.

5. Press enter to proceed using the list the `ucp` tool provided.

     UCP requires
    that all clients, including the Docker Engine, use a Swarm TLS certificate
    chain signed by the UCP Swarm Root CA. You can provide the certificate
    system with subject alternative names or SANs, which allow you to set up
    individual "leaf certificates."

    When it completes, the `ucp` tool prompts you to log in into the UCP web
    interface and gives you its location. You'll do this in the next stepo so
    you can install a license.

## Step 3. License your installation

In this step, you'll get a license, log in to the UCP web interface and install the license. Docker allows you to run an evaluation version of UCP with a single controller and node for up to 30 days.

1. Go to the [Docker Datacenter page](https://store.docker.com/bundles/docker-datacenter) in the Docker Store.

7. Click **Free 30-day evaluation** to select the free trial license type.

    If you're not logged in to the Docker Store, you can log in with an existing Docker ID, or create a new Docker ID from this page.
    Once you're logged in, continue to the next step.

8. Fill out the short form that appears. and click **Start your evaluation!**
    The screen refreshes to show your active subscription.

7. From the Subscription page, click **Subscription Details**, and select **Setup instructions** from the drop down menu.

8. The screen that appears contains installation instructions for when you are installing DDC on a production system. For now, you just need the trial license key.

9. Click `License Key` to download the `.lic` file to your local computer.

    Save the file to a safe location.

    ![](images/get-license.png)

10. Go back to your terminal window.
11. Copy the local IP address from the installer output.

    It will look something like https://192.168.99.100:443

12. Paste this IP address into your browser to view the UCP login screen.

    Your browser may warn you about the security of the connection. The warning
    appears because the UCP installer generated its own certificate which was
    issued by a built-in certificate authority (CA). The certificate's
    fingerprint is displayed during install and you can compare it to verify
    that it's the same one you expect.

2. Click the **Advanced** link and then the **Proceed to** link.

    The login screen appears.

    ![](images/login-ani.gif)

5. Enter the administrator username and password you provided during installation.

    Once you're logged in, the UCP dashboard appears and prompts for a license.

    ![](images/skip-this.png)

12. Click the **Upload License** button,

13. Locate and upload your `.lic` file.

    ![](images/license.png)

    Once you upload the file, the license message disappears from UCP.

You should now see the UCP Dashboard, showing one node connected.

![](images/dashboard.png)

## Step 4. Join a node

In this step, you join your `node2` to the controller using the `ucp join`
command. In a production installation, you'd do this step for each node
you want to add.

1. Open a terminal on your computer if you don't already have one open.

2. Connect the terminal environment to the `node2` you provisioned earlier.

    a. Use `docker-machine env` command to get the settings command for `node2`.

        ```none
        $$ docker-machine env node2
        export DOCKER_TLS_VERIFY="1"
        export DOCKER_HOST="tcp://192.168.99.101:2376"
        export DOCKER_CERT_PATH="/Users/ldr/.docker/machine/machines/node2"
        export DOCKER_MACHINE_NAME="node2"
        # Run this command to configure your shell:
        # eval $(docker-machine env node2)
        ```

    b. Run the `eval` command to set your environment.

        ```
        $ eval $(docker-machine env node2)
        ```

    Running this `eval` command sends the `docker` commands in the following steps to the Docker Engine on `node2`.

2. Run the `docker/ucp join` command.

    > **Note**: If you are on a Windows system, your shell won't be able to
    resolve the `$(docker-machine ip node2)` variable. Instead, edit the command
    supplied to replace it with the actual IP address.

    ```none
    $ docker run --rm -it \
    -v /var/run/docker.sock:/var/run/docker.sock \
    --name ucp docker/ucp join -i \
    --host-address $(docker-machine ip node2)
    ```

    The `join` command pulls several images, then prompts you for the URL of the UCP Server.

3. Enter the URL of the UCP server to continue.

4. Press `y` when prompted to continue and join the node to the swarm.

5. Enter the admin username and password for the UCP server when prompted.

    The installer continues and prompts you for SANs. In this sandbox, you've already provided the IP address and the `ucp` tool discovered this for you and shows it in the controller list.

5. Press `enter` to proceed without providing a SAN.

    The installation is complete when you see the message `Starting local swarm containers`.

4. Log in to UCP with your browser and confirm that the new node appears.

      ![](images/nodes.png)


## Step 5: Install Docker Trusted Registry

Next, we'll install Docker Trusted Registry (DTR). DTR provides a secure
location to store your organization's Docker images. Images are used by UCP to
run containers that make up a service. By providing a secure connection between
DTR and UCP, you can verify that your production services contain only signed
code produced by your own organization.

1. First, make sure you know the IP addresses of both your UCP and DTR nodes. You can find this easily by running `docker-machine ls`.

2. Open a terminal window, and enter the following command, replacing `$UCP_NODE_IP` with the IP address of your actual UCP instance.

    ```none
    $ curl -k https://$UCP_NODE_IP/ca > ucp-ca.pem
    ```

    This command downloads the ca certificate from your UCP installation, and saves it to a file. You'll use this in the next step. You may want to run `cat ucp-ca.pem` to make sure the file actually contains the certificate.

3. Next, use the following command to install DTR on `node2`.

    ```
    $ docker run -it --rm docker/dtr install --ucp-url $UCP_URL \
     --ucp-username $ADMIN_NAME --ucp-password $ADMIN_PASSWD --ucp-node node2 \
     --dtr-external-url $DTR_URL --ucp-ca "$(cat ucp-ca.pem)" \
    ```

    > **Tip**: You'll need to edit the command so it uses the correct IP addresses for your UCP and DTR nodes, and the correct administrator credentials. You might want to do this in a text editor. You may also omit the admin credentials from the command if you would prefer to be prompted for them during installation.

4. Verify that DTR is running by navigating your browser to the DTR server's IP.

5. Confirm that you can log in using your UCP administrator credentials.

## Step 6: Link UCP to your DTR instance

Now that you have your DTR instance up and running, we'll link it to your UCP instance. This allows you to use UCP to pull images from the DTR instance.

1. Navigate to the UCP web interface in your browser.
2. Log in using the administrator credentials.
3. Click **Settings** in the left menu, and then click the **DTR** tab.
4. Enter the URL of your DTR instance.
5. Make sure the **Insecure** checkbox is selected.

    In a production environment, you would upload a certificate instead. However, for this evaluation install we are using self-signed certificates which may not validate.

6. Click **Update Registry**, and click **Yes** in the confirmation dialog that appears.

**Congratulations!** You now have a working installation of Docker Datacenter running in your sandbox. You can explore on your own, or continue your evaluation by walking through our [guided tour](tutorial-sandbox.md).

### Further reading
* [UCP architecture](architecture.md)
* [UCP system requirements](installation/system-requirements.md)
* [Plan a production installation](installation/plan-production-install.md)
* [Install UCP for production](installation/install-production.md).



<!-- Wat.
Take a minute and explore UCP. At this point, you have a single controller
running. How many nodes is that? What makes a controller is the containers it
runs. Locate the Containers page and show the system containers on your
controller. You'll know you've succeeded if you see this list:

![](images/controller-containers.png)

The containers reflect the architecture of UCP.  The containers are running
Swarm, a key-value store process, and some containers with certificate volumes.
Explore the other resources. -->

<!--For this sandbox installation however, we're using self-signed certificates,
which will prevent you from being able to `docker pull` from the registry. We'll
work around that for this sandbox installation, using the steps below, or you
can read more about [Configuring security settings for DTR](https://docs.docker.com/docker-trusted-registry/configure/config-security/)
as you would do for a production deployment.-->
