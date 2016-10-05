<!--[metadata]>
+++
title = "Evaluate DDC in a sandbox"
description = "Evaluation installation"
keywords = ["tbd, tbd"]
[menu.main]
parent="mn_ucp"
identifier="ucp_tutorial_sandbox"
weight=10
+++
<![end-metadata]-->


# Evaluate DDC in a sandbox deployment

This tutorial assumes that you have [installed and configured](install-sandbox.md) a two-node Docker Datacenter installation including both UCP and DTR using the instructions [here](install-sandbox.md). If you haven't done this, we can't promise that this tutorial workflow will work exactly the same.

In the second half of this tutorial, we'll walk you through a typical deployment workflow using your sandbox installation of DDC as if it was a production instance installed on your organization's network.


Over the course of this tutorial, we will:
- Create a repository in DTR
- Set up certificates or set insecure flag
- Pull a Docker image, tag it and push it to your DTR repo.
- Edit the Docker image you just pulled, tag it and push it to your DTR repo.
- Set up your shell so you can interact with Docker objects in UCP using the command line
- Use UCP to deploy your edited image to a node


## Set --insecure registry or set up DTR trust and login

Next, we'll set up a security exception that allows a specific Docker-machine to push images to DTR even though the DTR instance has a self-signed certificate. For a production deployment, you would [set up certificate trust](https://docs.docker.com/ucp/configuration/dtr-integration/) between UCP and DTR, and [between DTR and your Docker Engine](https://docs.docker.com/docker-trusted-registry/repos-and-images/), but for our sandbox deployment we can skip this.

> **Warning**: These steps produce an insecure DTR connection. Do not use these configuration steps for a production deployment.

To allow your shell to connect to DTR despite it having a self-signed certificate, we'll specify that there is one insecure registry that we'll allow the Docker instance on `node1` (the UCP node) to connect to.  We'll add this exception by editing the configuration file where docker-machine stores `node1`'s configuration details.

1. Edit the file found at `~/.docker/machine/machines/node1/config.json` using your preferred text editor.

    For example `$ vi ~/.docker/machine/machines/node1/config.json`

2. Locate `InsecureRegistry` key in `EngineOptions` section, and add your DTR instance's IP between the brackets, enclosed in quotes.

    For example, `"InsecureRegistry": ["192.168.99.100"],`

3. Save your changes to the file and exit.

4. Run the command `docker-machine provision node1` to update `node1`'s configuration with the new `InsecureRegistry` setting.

This allows you to push docker images to, and pull docker images from, the registry running on `node2`.


## Step 1: Create an image repository in DTR

In this step, we'll create an image repository in DTR that you will be able to push Docker images to. Remember a Docker image is a combination of code and filesystem used as a template to create a container.

1. In your web browser, go to the DTR web UI.

    If you need help finding the URL for this host, you can use `docker-machine ls` to find the IP for `node2` where you installed it.

2. Log in to DTR using your administrator credentials.

3. Navigate to the Repositories screen and click New Repository.

4. In the repository name field, enter `my-nginx`.

5. Click **Save**.

## Step 2: Pull an image, tag and push to DTR

1. In your terminal, make sure `node1` is active using `docker-machine ls`.

    This is the node that you configured the security exception for, and if you are connecting to a Docker Engine without this exception you won't be able to push to your DTR instance.

    If necessary, use `docker-machine env` to make `node1` active.

    ```none
    $ eval "$(docker-machine env node1)"
    ```

2. Pull the latest Nginx image by running `docker pull ngnix:latest`

    Because you aren't specifying a registry as part of the `pull` command, Docker Engine locates and downloads the latest `ngnix` image from Docker Cloud's registry.

2. Log in to your DTR instance on `node2` using the `docker login` command.

    ```none
    docker login $DTR_IP
    ```

    Enter your administrator username and password when prompted.

3. Tag the `nginx` image you downloaded.
    Use the IP of your DTR instance to specify the repository path, and the .

    `docker tag nginx:latest [$DTR]/admin/my-nginx:official`

4. Push the tagged image to your DTR instance.

    `docker push my-nginx:official $DTR_IP/admin/my-nginx:official`

You now have a copy of the official Nginx Docker image available on your sandbox DTR instance.

## Step 3. Deploy a container from the UCP web interface

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

6. Click **Network** to expand the networking settings.

    A Docker container is isolated from other processes on your network and has its own internal network configuration. To access the service inside a container, you need to expose the container's port, which maps to a port on the node. The node is hosting an instance of Docker Engine, so its port is called the **Host Port**.

7. Enter `443` in the **Port** field and enter `4443` the **Host Port** field.

    We're mapping port 443 in the container to a different port on the host because your UCP instance is already serving the web interface on port 443.

8. Click the plus sign to add another **Port**.

9. For this port, enter `80` in the **Port** field, and enter `8080` in the **Host Port** field.

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

2. Click the nginx container.

    ![Edit](images/container_edit.png)

    The system displays the container's details and some operations you can run on the container.

3. Scroll down to the ports section.

    You'll see an IP address with port `8080` for the server.

4. Copy the IP address to your browser and paste the information you copied.

    You should see the welcome message for nginx.

    ![Port 80](images/welcome_nginx.png)


## Step 4: Edit your image, tag, and push

In this step, we'll edit the Nginx image so that it shows a customized webpage when you run it in a container. Then we'll tag the edited image, and push it to DTR so you can deploy it using UCP.

2. Change to your user `$HOME` directory.

    ```none
    $ cd $HOME
    ```

2. Make a `site` directory and open it.

    ```none
    $ mkdir site
    $ cd site
    ```

4. Create an `index.html` file.

    ```none
    $ echo "my new site" > index.html
    ```
5. Copy the `site` directory to the `node1` VM's file system.

    This allows you to copy files from the file system into a container.

    `docker-machine scp $HOME/site/ node1:~/site`

5. Start a new `nginx` container and replace the `html` folder with your `site` directory.

    ```none
    $ docker run -d -P -v $HOME/site:/usr/share/nginx/html --name mysite nginx
    ```

    This command runs an `nginx` image in a container called `mysite`. The `-P` tells the Engine to expose all the ports on the container.

6. Stop the container, and see that it's still around and called `mysite`

2. docker tag
    `docker tag nginx:mysite [$DTR]/admin/my-nginx:mysite`
7. docker push (edited) `[$DTR]/admin/my-nginx:mysite`


## Step 1. Download a client bundle

In this step, you download the *client bundle*, which contains your user
certificates, and a script to configure a shell environment with them so you can
connect to UCP.

Both nodes in your UCP cluster are running an instance of Engine. A UCP operator
can use the command line Engine client instead of the UCP web interface to
interact with the Docker objects and resources UCP manages.

However, to issue commands to a UCP node, your local shell environment must be
configured with the same security certificates as the UCP application itself.
Think of this as logging in to UCP with your user credentials, but from your
shell.

First, we'll download the bundle and configure your environment.

1. If you haven't already done so, log into UCP.

2. From the username menu at the top right, select **Profile**.

3. Scroll down and click **Create Client Bundle**.

    The browser downloads a `ucp-bundle-$USERNAME.zip` file.

4. Open a new shell on your local machine.

5. Make sure this new shell is not connected to a Docker Machine host by running `docker-machine ls`. No Docker host is connective if none of the hosts in the list show an * (asterisk) in the `ACTIVE` column.

<!-- Not needed here - needed elsewhere?
4. Create a directory to hold the deploy information.

    ```
    $ mkdir deploy-app
    ```
-->

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

5. From your shell, change into the directory created when the bundle was
unzipped.

6. Execute the `env.sh` script.

    This sets the appropriate environment variables for your UCP deployment.

    ```
    $ source env.sh
    ```

    If you are on Windows, you may need to open the `env.sh` file and review it so you can set the environment variables manually.

7. Run `docker info` to examine the UCP deployment.

    Your output should reflect a connection to UCP, instead of just a single local Docker instance. If you completed the previous installation tutorial, you should see `Nodes: 2` in the output.

    ```none
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






## Docker pull to put images on ucp

1. docker pull nginx - from docker hub/cloud
```
$ docker pull nginx
Using default tag: latest
node1: Pulling nginx:latest... : downloaded
node2: Pulling nginx:latest... : downloaded
```
 see that it shows up on both nodes




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
