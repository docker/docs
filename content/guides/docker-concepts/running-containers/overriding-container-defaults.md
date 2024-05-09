---
title: Overriding container defaults
keywords: concepts, build, images, container, docker desktop
description: This concept page will teach you how to override the container defaults using the `docker run` command.
---

{{< youtube-embed seY7D7Jx690 >}}

## Explanation

When a Docker container starts, it executes an application or command. The container gets this executable (script or file) from its image’s configuration. Containers come with default settings that usually work well, but you can change them if needed. These adjustments help the container's program run exactly how you want it to.

For example, if you have an existing database container that listens on the standard port and you want to run a new instance of the same database container, then you might want to change the port settings the new container listens on so that it doesn’t conflict with the existing container. Sometimes you might want to increase the memory available to the container if the program needs more resources to handle a heavy workload or set the environment variables to provide specific configuration details the program needs to function properly.

Now you might ask, "How can I override these container defaults?" The answer is `docker run` command.

The `docker run` command offers a powerful way to override these defaults and tailor the container's behavior to your liking. The command offers several flags that let you to customize container behavior on the fly.

Here's a few ways you can achieve this.

### Overriding the network ports

Sometimes you might want to use separate database instances for development and testing purposes. Running these database instances on the same port might conflict. You can use the `-p` option in `docker run` to map container ports to host ports, allowing you to run the multiple instances of the container without any conflict.

```console
$ docker run -d -p HOST_PORT:CONTAINER_PORT my_image
```

### Setting environment variables

This option sets an environment variable `foo` inside the container with the value `bar`.

```console
$ docker run -e foo=bar my_image env
```

You will see output like the following:

```console
HOSTNAME=2042f2e6ebe4
foo=bar
```

> **Tip**
>
> The `.env` file acts as a convenient way to set environment variables for your Docker containers without cluttering your command line with numerous `-e` flags. To use a `.env` file, you can pass `--env-file` option with the `docker run` command.
{ .tip }

### Restricting the container to consume the resources

You can use the `--memory` and `--cpus` flags with the `docker run` command to restrict how much CPU and memory a container can use. For example, you can set a memory limit for the Python API container, preventing it from consuming excessive resources on your host. Here's the command:

```console
$ docker run --memory="512m" --cpus="0.5" my_image
 ```

This command limits container memory usage to 512 MB and defines the CPU quota of 0.5 for half a core.

> **Monitor the real-time resource usage**
>
> You can use the `docker stats` command to monitor the real-time resource usage of running containers. This helps you understand whether the allocated resources are sufficient or need adjustment.
{ .tip }

By effectively using these `docker run` flags, you can tailor your containerized application's behavior to fit your specific requirements.

## Try it now

In this hands-on guide, you'll see how to use the `docker run` command to override the container defaults.

1. [Download and install](/get-docker/) Docker Desktop.

### Run the multiple instance of Redis database

1.  Start a container using the [Redis image](https://hub.docker.com/_/redis) with the following command:
    
    ```console
    $ docker run -d -p 6379:6379 redis
    ```

    This will start the Redis database in the background, listening on the standard container port `6379` and mapped to port `6379` on the host machine.

2. Start a new Redis container mapped to the different port

    ```console
    $ docker run -d -p 6380:6379 redis
    ```

    This will start a new Redis container in the background, listening on the standard container port `6379` but mapped to port `6380` on the host machine. You override the host port just to ensure that this new container doesn't conflict with the existing running container.

3. Verify if both containers are running via the Docker Dashboard.

    ![A screenshot of Docker Dashboard showing the running instances of Redis containers](images/running-redis-containers.webp?border=true)

### Run Redis container in a controlled network

By default, containers automatically connect to a special network called a bridge network when you run them. This bridge network acts like a virtual bridge, allowing containers on the same host to communicate with each other while keeping them isolated from the outside world and other hosts. It's a convenient starting point for most container interactions. However, for specific scenarios, you might want more control over the network configuration.

Here's where the `--network` flag with the `docker run` command comes in.

Follow the steps to see how to connect a Redis container to a custom network.

1. Create a network by using the following command:

    ```console
    $ docker network create mynetwork
    ```

2. Verify the network by running the following command:

    ```console
    $ docker network ls
    ```

3. Connect Redis to the existing network by using the following command:

    ```console
    $ docker run -d -p 6381:6379 –network mynetwork redis
    ```

    This will start Redis container in the background, mapped to the host port 6381 and attached to the `mynetwork` network. You passed the `--network` parameter to override the container default by connecting the container to custom Docker network for better isolation and communication with other containers.

    > **Inspecting the container network**
    >
    > You can use `docker network inspect` command to see if the container is tied to this new bridge network.
    { .tip }

### Manage the resources

By default, containers are not limited in their resource usage. However, on shared systems, it's crucial to manage resources effectively. It's important not to allow a running container to consume too much of the host machine's memory.

This is where the `docker run` command shines again. It offers flags like `--memory` and `--cpus` to restrict how much CPU and memory a container can use.

```console
$ docker run -d --memory=”512m” --cpus=”.5” redis
```

The `--cpus` flag specifies the CPU quota for the container. Here, it's set to half a CPU core (0.5) whereas the `--memory` flag specifies the memory limit for the container. In this case, it's set to 512 megabytes (512MB).

### Override the default `CMD` and `ENTRYPOINT` in Docker Compose

Sometimes, you might need to override the default commands (CMD) or entry points (ENTRYPOINT) defined in a Docker image, especially when using Docker Compose.

1. Create a `docker-compose.yaml` file with the following content:

    ```yaml
    services:
        redis:
         image: redis
         entrypoint: ["docker-entrypoint.sh"]
         command: ["redis-server", "--requirepass", "redispassword"]
    ```

    The YAML file defines a service named `redis` that uses the official Redis image, sets an entrypoint script, and starts the container with password authentication.

2. Bring up the service by running the following command:

    ```console
    $ docker compose up -d
    ```

    This command starts the Redis service defined in the Docker Compose file.

3. Verify the authentication with Docker Dashboard

    Open the **Docker Dashboard**, select the **Redis** container and select **Exec** to enter into the container shell.

    ![A screenshot of the Docker Dashboard selecting the Redis container and entering into its shell using EXEC button](images/exec-into-redis-container.webp?border=true) 

## Additional resources

* [Build Variables](https://docs.docker.com/build/building/variables/)
* [Ways to set environment variables with Compose](https://docs.docker.com/compose/environment-variables/set-environment-variables/)
* [What is a container](/guides/docker-concepts/the-basics/what-is-a-container/)

## Next steps

Now that you have learned about overriding container defaults, it's time to learn how to persist container data.

{{< button text="Persisting container data" url="persisting-container-data" >}}

