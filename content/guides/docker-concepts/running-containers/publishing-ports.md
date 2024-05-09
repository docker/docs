---
title: Publishing and exposing ports
keywords: concepts, build, images, container, docker desktop
description: This concept page will teach you the significance of publishing and exposing ports in Docker 
---

{{< youtube-embed 9JnqOmJ96ds >}}

## Explanation

If you've been following the guides so far, you understand that containers provide isolated processes for each component of your application. Each component - a React frontend, a Python API, and a Postgres database - runs in its own sandbox environment, completely isolated from everything else on your host machine. This isolation is great for security and managing dependencies, but it also means you can’t access them directly. For example, you can’t access the web app in your browser.

That’s where port publishing comes in.

### Publishing ports

Publishing a port provides the ability to break through a little bit of networking isolation by setting up a forwarding rule. As an example, you can indicate that requests on your host’s port `8080` should be forwarded to the container’s port `80`. Publishing ports happens during container creation using the `-p` (or `--publish`) flag with `docker run`. The syntax is:

```console
docker run -p HOST_PORT:CONTAINER_PORT my-image
```

- `HOST_PORT`: The port number on your host machine where you want to receive traffic
- `CONTAINER_PORT`: The port number within the container that's listening for connections

For example, to publish the container's port `80` to host port `8080`:

```console
docker run -p 8080:80 my-image
```

Now, any traffic sent to port `8080` on your host machine will be forwarded to port `80` within the container.

> **Important**
>
> When a port is published, it's published to all network interfaces by default. This means any traffic that reaches your machine can access the published application. Be mindful of publishing databases or any sensitive information. [Learn more about published ports here](https://docs.docker.com/network/#published-ports).
{ .important }

### Publishing to ephemeral ports

At times, you may want to simply publish the port but don’t care which host port is used. In these cases, you can let Docker pick the port for you. To do so, simply omit the `HOST_PORT` configuration. 

For example, the following command will publish the container’s port `80` onto an ephemeral port on the host:

```console
docker run -p 80 docker/welcome-to-docker
```
 
Once the container is running, using `docker ps` will show you the port that was chosen:

```console
CONTAINER ID   IMAGE                      COMMAND                  CREATED         STATUS        PORTS                                     NAMES
5646fec21960   docker/welcome-to-docker   "/docker-entrypoint.…"   2 seconds ago   Up 1 second   0.0.0.0:32770->80/tcp, :::32770->80/tcp   busy_shaw
```

In this example, the app is exposed on the host at port `32770`.

### Publishing all ports

When creating a container image, the `EXPOSE` instruction is used to indicate the packaged application will use the specified port. These ports aren't published by default. 

With the `-P` or `--publish-all` flag, you can automatically publish all exposed ports to ephemeral ports. This is quite useful when you’re trying to avoid port conflicts in development or testing environments.

For example, the following command will publish all of the exposed ports configured by the image:

```console
docker run -P docker/welcome-to-docker
```

## Try it now

In this hands-on guide, you'll learn how to publish container ports using both the CLI and Docker Compose for deploying a web application.

### Use the Docker CLI

In this step, you will run a container and publish its port using the Docker CLI.

1. Download and install Docker Desktop

2. In a terminal, run the following command to start a new container:

    ```console
    $ docker run -d -p 8080:80 docker/welcome-to-docker
    ```

    The first `8080` refers to the host port. This is the port on your local machine that will be used to access the application running inside the container. The second `80` refers to the container port. This is the port that the application inside the container listens on for incoming connections. Hence, the command binds to port `8080` of the host to port `80` on the container system.

3. Verify the published port by going to the **Containers** view of the Docker Dashboard.

   ![A screenshot of Docker dashboard showing the published port](images/published-ports.webp?border=true)

4. Open the website by either selecting the link in the **Port(s)** column of your container or visiting [http://localhost:8080](http://localhost:8080) in your browser.

   ![A screenshot of the landing page of the Nginx web server running in a container](images/access-the-frontend.webp?border=true)

### Use Docker Compose

This example will launch the same application using Docker Compose:

1. Create a new directory and inside that directory, create a `compose.yaml` file with the following contents:

    ```yaml {hl_lines=[4,5]}
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

* [Docker Container Port](https://docs.docker.com/reference/cli/docker/container/port/)
* [Published Ports](https://docs.docker.com/network/#published-ports)

## Next steps

Now that you understand how to publish and expose ports, you're ready to learn how to override the container defaults using the `docker run` command.

{{< button text="Overriding container defaults" url="overriding-container-defaults" >}}

