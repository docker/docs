---
title: Publishing and exposing ports
keywords: concepts, build, images, container, docker desktop
description: This concept page will teach you the significance of publishing and exposing ports in Docker 
weight: 1
aliases: 
 - /guides/docker-concepts/running-containers/publishing-ports/
---

{{< youtube-embed 9JnqOmJ96ds >}}

## Explanation

If you've been following the guides so far, you understand that containers provide isolated processes for each component of your application. Each component - a React frontend, a Python API, and a Postgres database - runs in its own sandbox environment, completely isolated from everything else on your host machine. This isolation is great for security and managing dependencies, but it also means you can’t access them directly. For example, you can’t access the web app in your browser.

That’s where port publishing comes in.

### Publishing ports

Publishing a port provides the ability to break through a little bit of networking isolation by setting up a forwarding rule. As an example, you can indicate that requests on your host’s port `8080` should be forwarded to the container’s port `80`. Publishing ports happens during container creation using the `-p` (or `--publish`) flag with `docker run`. The syntax is:

```console
$ docker run -d -p HOST_PORT:CONTAINER_PORT nginx
```

- `HOST_PORT`: The port number on your host machine where you want to receive traffic
- `CONTAINER_PORT`: The port number within the container that's listening for connections

For example, to publish the container's port `80` to host port `8080`:

```console
$ docker run -d -p 8080:80 nginx
```

Now, any traffic sent to port `8080` on your host machine will be forwarded to port `80` within the container.

> [!IMPORTANT]
>
> When a port is published, it's published to all network interfaces by default. This means any traffic that reaches your machine can access the published application. Be mindful of publishing databases or any sensitive information. See [Binding to specific network interfaces](#binding-to-specific-network-interfaces) below, and [learn more about published ports here](/engine/network/#published-ports).

### Binding to specific network interfaces

By default, when you publish a port, Docker binds to all network interfaces (`0.0.0.0`). However, there are scenarios where you might want to restrict access to a specific network interface or IP address. You can do this by using the extended port syntax:

```console
$ docker run -d -p IP:HOST_PORT:CONTAINER_PORT nginx
```

- `IP`: The specific IP address or network interface on your host to bind to
- `HOST_PORT`: The port number on your host machine
- `CONTAINER_PORT`: The port number within the container

For example, if you want to make your container accessible only on localhost (127.0.0.1), you can use:

```console
$ docker run -d -p 127.0.0.1:8080:80 nginx
```

This restricts access to the container's port 80 to only local connections on your host's port 8080. External machines on your network would not be able to access this service. You can verify the binding with `docker ps`:

```console
$ docker ps
CONTAINER ID   IMAGE     COMMAND                  CREATED          STATUS          PORTS                          NAMES
a527355c9c53   nginx     "/docker-entrypoint.…"   4 seconds ago    Up 3 seconds    127.0.0.1:8080->80/tcp         elegant_newton
```

#### Common use cases for IP binding

- **Security**: Binding to localhost (127.0.0.1) to prevent external access to services like databases
- **Multi-homed hosts**: Binding to specific network interfaces on servers with multiple IP addresses
- **Service isolation**: Limiting which network segments can access certain containers

#### In Docker Compose
You can also use the IP binding syntax in your Docker Compose files:

```yaml
services:
  webapp:
    image: nginx
    ports:
      - "127.0.0.1:8080:80"
  
  database:
    image: postgres
    ports:
      - "127.0.0.1:5432:5432"
```

This configuration ensures your database is only accessible from the host machine itself, adding a layer of security to your application stack. You can also set a default bind address for all containers using the `host_binding_ipv4` configuration parameter. This allows you to change the default binding from `0.0.0.0` to a specific IP address. For more information, see [Setting the default bind address for containers](/engine/network/packet-filtering-firewalls/#setting-the-default-bind-address-for-containers).

### Publishing to ephemeral ports

At times, you may want to simply publish the port but don’t care which host port is used. In these cases, you can let Docker pick the port for you. To do so, simply omit the `HOST_PORT` configuration. 

For example, the following command will publish the container’s port `80` onto an ephemeral port on the host:

```console
$ docker run -p 80 nginx
```
 
Once the container is running, using `docker ps` will show you the port that was chosen:

```console
docker ps
CONTAINER ID   IMAGE         COMMAND                  CREATED          STATUS          PORTS                    NAMES
a527355c9c53   nginx         "/docker-entrypoint.…"   4 seconds ago    Up 3 seconds    0.0.0.0:54772->80/tcp    romantic_williamson
```

In this example, the app is exposed on the host at port `54772`.

### Publishing all ports

When creating a container image, the `EXPOSE` instruction is used to indicate the packaged application will use the specified port. These ports aren't published by default. 

With the `-P` or `--publish-all` flag, you can automatically publish all exposed ports to ephemeral ports. This is quite useful when you’re trying to avoid port conflicts in development or testing environments.

For example, the following command will publish all of the exposed ports configured by the image:

```console
$ docker run -P nginx
```

## Try it out

In this hands-on guide, you'll learn how to publish container ports using both the CLI and Docker Compose for deploying a web application.

### Use the Docker CLI

In this step, you will run a container and publish its port using the Docker CLI.

1. [Download and install](/get-started/get-docker/) Docker Desktop.

2. In a terminal, run the following command to start a new container:

    ```console
    $ docker run -d -p 8080:80 docker/welcome-to-docker
    ```

    The first `8080` refers to the host port. This is the port on your local machine that will be used to access the application running inside the container. The second `80` refers to the container port. This is the port that the application inside the container listens on for incoming connections. Hence, the command binds to port `8080` of the host to port `80` on the container system.

3. Verify the published port by going to the **Containers** view of the Docker Desktop Dashboard.

   ![A screenshot of Docker Desktop Dashboard showing the published port](images/published-ports.webp?w=5000&border=true)

4. Open the website by either selecting the link in the **Port(s)** column of your container or visiting [http://localhost:8080](http://localhost:8080) in your browser.

   ![A screenshot of the landing page of the Nginx web server running in a container](/get-started/docker-concepts/the-basics/images/access-the-frontend.webp?border=true)


### Use Docker Compose

This example will launch the same application using Docker Compose:

1. Create a new directory and inside that directory, create a `compose.yaml` file with the following contents:

    ```yaml
    services:
      app:
        image: docker/welcome-to-docker
        ports:
          - 8080:80
    ```

    The `ports` configuration accepts a few different forms of syntax for the port definition. In this case, you’re using the same `HOST_PORT:CONTAINER_PORT` used in the `docker run` command.

2. Open a terminal and navigate to the directory you created in the previous step.

3. Use the `docker compose up` command to start the application. 

4. Open your browser to [http://localhost:8080](http://localhost:8080).

## Additional resources

If you’d like to dive in deeper on this topic, be sure to check out the following resources:

* [`docker container port` CLI reference](/reference/cli/docker/container/port/)
* [Published ports](/engine/network/#published-ports)
* [Network drivers](/engine/network/drivers/) - Learn more about Docker's networking capabilities
* [Docker network command](/engine/reference/commandline/network/) - Reference for the `docker network` command

## Next steps

Now that you understand how to publish and expose ports, you're ready to learn how to override the container defaults using the `docker run` command.

{{< button text="Overriding container defaults" url="overriding-container-defaults" >}}

