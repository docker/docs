<!--[metadata]>
+++
title = "Evaluate DDC in a sandbox "
description = "Evaluation installation"
keywords = ["tbd, tbd"]
[menu.main]
parent="mn_ucp"
identifier="ucp_tutorial_sandbox"
weight=10
+++
<![end-metadata]-->


# Evaluate DDC in a sandbox deployment

This tutorial assumes that you have [installed and configured](install-sandbox.md) a two-node DDC installation of UCP and DTR using the instructions [here](install-sandbox.md). If you haven't done this, we can't promise that this tutorial workflow will work exactly the same.

In the second half of this tutorial, we'll walk you through a normal deployment workflow using your sandbox installation of DDC as if it was a production instance installed on your organization's network.

Over the course of this tutorial, we will:
- Set up your shell so you can interact with Docker obejects in UCP using the command line .
- pull a Docker image, edit it and tag it,
- Create a repository in DTR, and push your edited image to that repo.
-




## Step 10. Download a client bundle

In this step, you download the *client bundle*, which contains the certificates
so your shell can connect to UCP, and a script to configure a shell environment
with them.

Both nodes in your UCP cluster are running an instance of Engine. A UCP operator
can use the command line Engine client instead of the UCP web interface to
interact with the Docker objects and resources UCP manages.

However, to issue commands to a UCP node, your local shell environment must be
configured with the same security certificates as the UCP application itself.
Think of this as logging in to UCP with your shell.

Download the bundle and configure your environment.

1. If you haven't already done so, log into UCP.

2. From the username menu at the top right, select **Profile**.

3. Scroll down and click **Create Client Bundle**.

    The browser downloads the `ucp-bundle-$USERNAME.zip` file.

4. Open a new shell on your local machine.

5. Make sure your shell is does not have an active Docker Machine host.

    ```none
    $ docker-machine ls
    NAME      ACTIVE   DRIVER       STATE     URL                         SWARM   DOCKER    ERRORS
    default   -        virtualbox   Stopped                                       Unknown
    node1     -        virtualbox   Running   tcp://192.168.99.100:2376           v1.12.1
    node2     -        virtualbox   Running   tcp://192.168.99.101:2376           v1.12.1
    ```

    While Machine has a stopped and running host, neither is active in the shell. You know this because neither host shows an * (asterisk) indicating the shell is configured.

4. Create a directory to hold the deploy information.

    ```
    $ mkdir deploy-app
    ```

4. Navigate to where the bundle was downloaded, and unzip the client bundle

    ```none
    $ unzip bundle.zip
    Archive:  bundle.zip
    	extracting: ca.pem
    extracting: cert.pem
    extracting: key.pem
    	extracting: cert.pub
    extracting: env.sh
    ```

5. Change into the directory that was created when the bundle was unzipped

6. Execute the `env.sh` script to set the appropriate environment variables for your UCP deployment.

    ```
    $ source env.sh
    ```

    If you are on Windows, you may need to open the `env.sh` file and review it so you can set the environment variables manually.

7. Run `docker info` to examine the UCP deployment.

    Your output should show that you are managing UCP vs. a single node.

    ```none
    $ docker info
    Containers: 12
     Running: 0
     Paused: 0
     Stopped: 0
    Images: 17
    Role: primary
    Strategy: spread
    Filters: health, port, dependency, affinity, constraint
    Nodes: 2
     node1: 192.168.99.106:12376
      └ Status: Healthy
      └ Containers: 9
      └ Reserved CPUs: 0 / 1
      └ Reserved Memory: 0 B / 3.01 GiB
      └ Labels: executiondriver=native-0.2, kernelversion=4.1.17-boot2docker, operatingsystem=Boot2Docker 1.10.0 (TCL 6.4.1); master : b09ed60 - Thu Feb  4 20:16:08 UTC 2016, provider=virtualbox, storagedriver=aufs
      └ Error: (none)
      └ UpdatedAt: 2016-02-09T12:03:16Z
     node2: 192.168.99.107:12376
      └ Status: Healthy
      └ Containers: 3
      └ Reserved CPUs: 0 / 1
      └ Reserved Memory: 0 B / 4.956 GiB
      └ Labels: executiondriver=native-0.2, kernelversion=4.1.17-boot2docker, operatingsystem=Boot2Docker 1.10.0 (TCL 6.4.1); master : b09ed60 - Thu Feb  4 20:16:08 UTC 2016, provider=virtualbox, storagedriver=aufs
      └ Error: (none)
      └ UpdatedAt: 2016-02-09T12:03:11Z
    Cluster Managers: 1
     192.168.99.106: Healthy
      └ Orca Controller: https://192.168.99.106:443
      └ Swarm Manager: tcp://192.168.99.106:3376
      └ KV: etcd://192.168.99.106:12379
    Plugins:
     Volume:
     Network:
    CPUs: 2
    Total Memory: 7.966 GiB
    Name: ucp-controller-node1
    ID: P5QI:ZFCX:ELZ6:RX2F:ADCT:SJ7X:LAMQ:AA4L:ZWGR:IA5V:CXDE:FTT2
    WARNING: No oom kill disable support
    WARNING: No cpu cfs quota support
    WARNING: No cpu cfs period support
    WARNING: No cpu shares support
    WARNING: No cpuset support
    Labels:
     com.docker.ucp.license_key=p3vPAznHhbitGG_KM36NvCWDiDDEU7aP_Y9z4i7V4DNb
     com.docker.ucp.license_max_engines=1
     com.docker.ucp.license_expires=2016-11-11 00:53:53 +0000 UTC
     ```

     ```
    $ docker info
    Containers: 18
    Images: 28
    Server Version: swarm/1.2.5
    Role: primary
    Strategy: spread
    Filters: health, port, containerslots, dependency, affinity, constraint
    Nodes: 2
     node1: 192.168.99.100:12376
      └ ID: ILMX:BKG4:NW6K:YP3E:EERA:WAHA:V6KW:TC22:TND2:DUTM:NALD:OFAI
      └ Status: Healthy
      └ Containers: 11 (10 Running, 0 Paused, 1 Stopped)
      └ Reserved CPUs: 0 / 1
      └ Reserved Memory: 0 B / 2.004 GiB
      └ Labels: kernelversion=4.4.17-boot2docker, operatingsystem=Boot2Docker 1.12.1 (TCL 7.2); HEAD : ef7d0b4 - Thu Aug 18 21:18:06 UTC 2016, provider=virtualbox, storagedriver=aufs
      └ UpdatedAt: 2016-09-30T00:18:00Z
      └ ServerVersion: 1.12.1
     node2: 192.168.99.101:12376
      └ ID: DUHM:NIOB:WRIP:IDRL:RNNC:HQQ4:HL24:KWEZ:BHQ2:JFDX:DZCK:YV7Q
      └ Status: Healthy
      └ Containers: 7 (7 Running, 0 Paused, 0 Stopped)
      └ Reserved CPUs: 0 / 1
      └ Reserved Memory: 0 B / 2.004 GiB
      └ Labels: kernelversion=4.4.17-boot2docker, operatingsystem=Boot2Docker 1.12.1 (TCL 7.2); HEAD : ef7d0b4 - Thu Aug 18 21:18:06 UTC 2016, provider=virtualbox, storagedriver=aufs
      └ UpdatedAt: 2016-09-30T00:18:18Z
      └ ServerVersion: 1.12.1
    Cluster Managers: 1
     192.168.99.100: Healthy
      └ Orca Controller: https://192.168.99.100:443
      └ Swarm Manager: tcp://192.168.99.100:3376
      └ KV: etcd://192.168.99.100:12379
    Kernel Version: 4.4.17-boot2docker
    Operating System: linux
    CPUs: 2
    Total Memory: 4.008 GiB
    Name: ucp-controller-node1
    ID: H6IE:2NQ6:JWMJ:UM2B:GKOM:N6CI:Q3WX:RDRW:3RO4:HKKD:YMDU:SEFS
    Labels:
     com.docker.ucp.license_key=cbikfA44-5gAJ3iOgPp_6AMl_V_uRxFiITvyVvESdFWx
     com.docker.ucp.license_max_engines=10
     com.docker.ucp.license_expires=2016-10-29 00:01:45 +0000 UTC
     ```

<!-- ## Create repo

log in to dtr.
create a repo called hi-there
docker pull hello-world
docker tag hello-world:latest as [$DTR]/admin/hi-there
docker push ohai [$DTR]/admin/hi-there


Next, make sure node1 is active
Docker login

`docker tag alpine:latest <path to dtr>/admin/foobar`
 then
 `docker push <path to dtr>/admin/foobar`

-->

## Step 8. Deploy a container

UCP allows you to deploy and manage "Dockerized" applications in production. An
application is built using Docker objects, such as images and containers, and
Docker resources, such as volumes and networks.

UCP deploys and manages these objects and resources using remote API calls to
the Engine daemons running on the nodes. For example, the `run` action may
deploy an image in a Docker container. That image might define a service such as
an Nginx web server or a database like Postgres.

A UCP administrator initiates Engine actions using the UCP dashboard or the
Docker Engine CLI. In this step, you deploy a container from the UCP dashboard.
The container runs an Nginx server, so you'll need to launch the `nginx` image
inside of it.

1. Log in to the UCP **Dashboard**.

2. Click **Containers**.

    The system displays the **Containers** page.

    > **Tip**: UCP runs some containers that support its own operations called
    "system" containers. These containers are hidden by default.

3. Click **+ Deploy Container**.

    UCP provides a dialog for you to enter configuration options for the container. For this example, we'll deploy a simple `nginx` container using specific values for each field. If you already know what you're doing, feel free to explore once you've completed this example.

4. Enter `nginx` for the image name.

    An image is a specific build of software you want to run. The software might
    be a stand-alone application, or component software necessary to support a
    complex service.

5. Enter `nginx_server` for the container name.

    This name just identifies the container on your network.

6. Click **Publish Ports** from the **Overview** menu.

    A Docker container is isolated from other processes on your network and has its own internal network configuration. To access the service inside a container, you need to expose the container's port, which maps to a port on the node. The node is hosting an instance of Docker Engine, so its port is called the **Host Port**.

7. Enter `443` in both the **Port** field and the **Host Port** field.

8. Click the plus sign to add another **Port**.

9. For this port, enter `80` in both the **Port** and **Host Port** fields.

    When you are done, your dialog should look like this:

    ![Port configuration](images/port_config.png)

10. Click **Run Container** to deploy the container.

    ![Deployed](images/display_container.png)

## Step 9. View a running service

At this point, you have deployed a container and you should see the container
status is `running`. Recall that you deployed an Nginx web server. That server
comes with a default page that you can view to validate that the server is
running. In this step, you open the running server.

1. Navigate back to the **Containers** page.

2. Click the edit icon on the container.

    ![Edit](images/container_edit.png)

    The system displays the container's details and some operations you can run on the container.

3. Scroll down to the ports section.

    You'll see an IP address with port `80` for the server.

4. Copy the IP address to your browser and paste the information you copied.

    You should see the welcome message for nginx.

    ![Port 80](images/welcome_nginx.png)

## Step 11. Deploy with the CLI

In this exercise, you'll launch another Nginx container. Only this time, you'll use the Engine CLI. Then, you'll look at the result in the UCP dashboard.

1. Connect the terminal environment to the `node2`.

    ```none
    $ eval "$(docker-machine env node2)"
    ```

2. Change to your user `$HOME` directory.

    ```none
    $ cd $HOME
    ```

2. Make a `site` directory.

    ```none
    $ mkdir site
    ```

3. Change into the `site` directory.

    ```none
    $ cd site
    ```

4. Create an `index.html` file.

    ```none
    $ echo "my new site" > index.html
    ```

5. Start a new `nginx` container and replace the `html` folder with your `site` directory.

    ```none
    $ docker run -d -P -v $HOME/site:/usr/share/nginx/html --name mysite nginx
    ```

    This command runs an `nginx` image in a container called `mysite`. The `-P` tells the Engine to expose all the ports on the container.

6. Open the UCP dashboard in your browser.

7. Navigate to the **Containers** page and locate your `mysite` container.

    ![mysite](images/second_node.png)

8. Scroll down to the ports section.

    You'll see an IP address with port `80/tcp` for the server. This time,
    you'll find that the port mapped on this container than the one created
    yourself. That's because the command didn't explicitly map a port, so the
    Engine chose mapped the default Nginx port `80` inside the container to an
    arbitrary port on the node.

4. Copy the IP address to your browser and paste the information you copied.

    You should see your `index.html` file display instead of the standard Nginx welcome.

    ![mysite](images/second_node.png)

## Explore UCP

At this point, you've completed the guided tour of a UCP installation. You've
learned how to create a UCP installation by creating two nodes and designating
one of them as a controller. You've created a container running a simple web
server both using UCP and directly on the command line.  You used UCP to get
information about what you created.

In a real UCP production installation, UCP admins and operators are expected to
do similar work every day. While the applications they launch will be more
complicated, the interaction channels a user can take, the GUI or the
certificate bundle plus a command line, remain the same.

Take some time to explore UCP some more. Investigate the documentation for other
activities you can perform with UCP.


| `docker run --rm -it docker/ucp` | `uninstall --help`      | Uninstalls UCP  |
