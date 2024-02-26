---
title: Run your Rust image as a container
keywords: rust, run, image, container,
description: Learn how to run your Rust image as a container.
---

## Prerequisite

You have completed  [Build your Rust image](build-images.md) and you have built an image.

## Overview

A container is a normal operating system process except that Docker isolates this process so that it has its own file system, its own networking, and its own isolated process tree separate from the host.

To run an image inside of a container, you use the `docker run` command. The `docker run` command requires one parameter which is the name of the image.

## Run an image

Use `docker run` to run the image you built in [Build your Rust image](build-images.md).

```console
$ docker run docker-rust-image
```

After running this command, you’ll notice that you weren't returned to the command prompt. This is because your application is a server that runs in a loop waiting for incoming requests without returning control back to the OS until you stop the container.

Open a new terminal then make a request to the server using the `curl` command.

```console
$ curl http://localhost:8000
```

You should see output like the following.

```console
curl: (7) Failed to connect to localhost port 8000 after 2236 ms: Couldn't connect to server
```

As you can see, your `curl` command failed. This means you weren't able to connect to the localhost on port 8000. This is normal because your container is running in isolation which includes networking. Stop the container and restart with port 8000 published on your local network.

To stop the container, press ctrl-c. This will return you to the terminal prompt.

To publish a port for your container, you’ll use the `--publish` flag (`-p` for short) on the `docker run` command. The format of the `--publish` command is `[host port]:[container port]`. So, if you wanted to expose port 8000 inside the container to port 3001 outside the container, you would pass `3001:8000` to the `--publish` flag.

You didn't specify a port when running the application in the container and the default is 8000. If you want your previous request going to port 8000 to work, you can map the host's port 3001 to the container's port 8000:

```console
$ docker run --publish 3001:8000 docker-rust-image
```

Now, rerun the curl command. Remember to open a new terminal.

```console
$ curl http://localhost:3001
```

You should see output like the following.

```console
Hello, Docker!
```

Success! You were able to connect to the application running inside of your container on port 8000. Switch back to the terminal where your container is running and stop it.

Press ctrl-c to stop the container.

## Run in detached mode

This is great so far, but your sample application is a web server and you don't have to be connected to the container. Docker can run your container in detached mode or in the background. To do this, you can use the `--detach` or `-d` for short. Docker starts your container the same as before but this time will "detach" from the container and return you to the terminal prompt.

```console
$ docker run -d -p 3001:8000 docker-rust-image
ce02b3179f0f10085db9edfccd731101868f58631bdf918ca490ff6fd223a93b
```

Docker started your container in the background and printed the Container ID on the terminal.

Again, make sure that your container is running properly. Run the curl command again.

```console
$ curl http://localhost:3001
```

You should see output like the following.

```console
Hello, Docker!
```

## List containers

Since you ran your container in the background, how do you know if your container is running or what other containers are running on your machine? Well, to see a list of containers running on your machine, run `docker ps`. This is similar to how you use the ps command in Linux to see a list of processes.

You should see output like the following.

```console
CONTAINER ID   IMAGE                   COMMAND         CREATED         STATUS         PORTS                    NAMES
3074745e412c   docker-rust-image       "/bin/server"   8 seconds ago   Up 7 seconds   0.0.0.0:3001->8000/tcp   wonderful_kalam
```

The `docker ps` command provides a bunch of information about your running containers. You can see the container ID, the image running inside the container, the command that was used to start the container, when it was created, the status, ports that were exposed, and the name of the container.

You are probably wondering where the name of your container is coming from. Since you didn’t provide a name for the container when you started it, Docker generated a random name. You’ll fix this in a minute, but first you need to stop the container. To stop the container, run the `docker stop` command which does just that, stops the container. You need to pass the name of the container or you can use the container ID.

```console
$ docker stop wonderful_kalam
wonderful_kalam
```

Now, rerun the `docker ps` command to see a list of running containers.

```console
$ docker ps
CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS              PORTS               NAMES
```

## Stop, start, and name containers

You can start, stop, and restart Docker containers. When you stop a container, it's not removed, but the status is changed to stopped and the process inside the container is stopped. When you ran the `docker ps` command in the previous module, the default output only shows running containers. When you pass the `--all` or `-a` for short, you see all containers on your machine, irrespective of their start or stop status.

```console
$ docker ps -a
CONTAINER ID   IMAGE                   COMMAND                  CREATED          STATUS                      PORTS                       
     NAMES
3074745e412c   docker-rust-image       "/bin/server"            3 minutes ago    Exited (0) 6 seconds ago                                
     wonderful_kalam
6cfa26e2e3c9   docker-rust-image       "/bin/server"            14 minutes ago   Exited (0) 5 minutes ago                                
     friendly_montalcini
4cbe94b2ea0e   docker-rust-image       "/bin/server"            15 minutes ago   Exited (0) 14 minutes ago                               
     tender_bose
```

You should now see several containers listed. These are containers that you started and stopped but you haven't removed.

Restart the container that you just stopped. Locate the name of the container you just stopped and replace the name of the container in following restart command.

```console
$ docker restart wonderful_kalam
```

Now list all the containers again using the `docker ps` command.

```console
$ docker ps --all
CONTAINER ID   IMAGE                   COMMAND                  CREATED          STATUS                      PORTS                       
     NAMES
3074745e412c   docker-rust-image       "/bin/server"            6 minutes ago    Up 4 seconds                0.0.0.0:3001->8000/tcp           wonderful_kalam
6cfa26e2e3c9   docker-rust-image       "/bin/server"            16 minutes ago   Exited (0) 7 minutes ago                                
     friendly_montalcini
4cbe94b2ea0e   docker-rust-image       "/bin/server"            18 minutes ago   Exited (0) 17 minutes ago                               
     tender_bose
```

Notice that the container you just restarted has been started in detached mode. Also, observe the status of the container is "Up X seconds". When you restart a container, it starts with the same flags or commands that it was originally started with.

Now, stop and remove all of your containers and take a look at fixing the random naming issue. Stop the container you just started. Find the name of your running container and replace the name in the following command with the name of the container on your system.

```console
$ docker stop wonderful_kalam
wonderful_kalam
```

Now that you have stopped all of your containers, remove them. When you remove a container, it's no longer running, nor is it in the stopped status, but the process inside the container has been stopped and the metadata for the container has been removed.

To remove a container, run the `docker rm` command with the container name. You can pass multiple container names to the command using a single command. Again, replace the container names in the following command with the container names from your system.

```console
$ docker rm wonderful_kalam friendly_montalcini tender_bose
wonderful_kalam
friendly_montalcini
tender_bose
```

Run the `docker ps --all` command again to see that Docker removed all containers.

Now, it's time to address the random naming issue. Standard practice is to name your containers for the simple reason that it's easier to identify what's running in the container and what application or service it's associated with.

To name a container, you just need to pass the `--name` flag to the `docker run` command.

```console
$ docker run -d -p 3001:8000 --name docker-rust-container docker-rust-image
1aa5d46418a68705c81782a58456a4ccdb56a309cb5e6bd399478d01eaa5cdda
$ docker ps
CONTAINER ID   IMAGE                   COMMAND         CREATED         STATUS         PORTS                    NAMES
c68fa18de1f6   docker-rust-image       "/bin/server"   7 seconds ago   Up 6 seconds   0.0.0.0:3001->8000/tcp   docker-rust-container
```

That’s better! You can now easily identify your container based on the name.

## Summary

In this section, you took a look at running containers. You also took a look at managing containers by starting, stopping, and restarting them. And finally, you looked at naming your containers so they are more easily identifiable.

Related information:
 - [docker run CLI reference](../../reference/cli/docker/container/run.md)

## Next steps

In the next section, you’ll learn how to run a database in a container and connect it to a Rust application.

{{< button text="How to develop your application" url="develop.md" >}}
