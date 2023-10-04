---
title: Run your image as a container
keywords: Python, run, image, container,
description: Learn how to run the image as a container.
---

## Prerequisites

Work through the steps to build a Python image in [Build your Python image](build-images.md).

## Overview

In the previous module, you got the sample application and then created a Dockerfile that you used to produce an image. You created your image using the Docker command `docker build`. Now that you have an image, you can run that image and see if your application is running correctly.

A container is a normal operating system process except that this process is isolated in that it has its own file system, its own networking, and its own isolated process tree separate from the host.

To run an image inside of a container, use the `docker run` command. The `docker run` command requires one parameter which is the name of the image. Start your image and make sure it is running correctly. Run the following command in your terminal.

```console
$ docker run python-docker
```

After running this command, you’ll notice that you were not returned to the command prompt. This is because your application is a REST server and runs in a loop waiting for incoming requests without returning control back to the OS until you stop the container.

Open a new terminal then make a `GET` request to the server using the `curl` command.

```console
$ curl localhost:5000
curl: (7) Failed to connect to localhost port 5000: Connection refused
```

As you can see, your `curl` command failed because the connection to your server was refused. This means you were not able to connect to localhost on port 5000. This is expected because your container is running in isolation which includes networking. Stop the container and restart with port 5000 published on your local network.

To stop the container, press ctrl-c. This will return you to the terminal prompt.

To publish a port for your container, use the `--publish` flag (`-p` for short) on the `docker run` command. The format of the `--publish` command is `[host port]:[container port]`. So, if you wanted to expose port 5000 inside the container to port 3000 outside the container, you would pass `3000:5000` to the `--publish` flag.

You did not specify a port when running the flask application in the container and the default is 5000. If you want your previous request going to port 5000 to work, you can map the host's port 8000 to the container's port 5000:

```console
$ docker run --publish 8000:5000 python-docker
```

Now, open a new terminal window and re-run the `curl` command.

```console
$ curl localhost:8000
Hello, Docker!
```

Success! You were able to connect to the application running inside of your container on port 8000. Switch back to the terminal where your container is running and you should see the `GET` request logged to the console.

```shell
[31/Jan/2021 23:39:31] "GET / HTTP/1.1" 200 -
```

Press ctrl-c to stop the container.

## Run in detached mode

This is great so far, but your sample application is a web server and you don't have to be connected to the container. Docker can run your container in detached mode or in the background. To do this, you can use the `--detach` or `-d` for short. Docker starts your container the same as before but this time will “detach” from the container and return you to the terminal prompt.

```console
$ docker run -d -p 8000:5000 python-docker
ce02b3179f0f10085db9edfccd731101868f58631bdf918ca490ff6fd223a93b
```

Docker started your container in the background and printed the container ID on the terminal.

Again, make sure that your container is running properly. Run the same curl command from above.

```console
$ curl localhost:8000
Hello, Docker!
```

## List containers

Since you ran your container in the background, how do you know if your container is running or what other containers are running on your machine? Well, to see a list of containers running on your machine, run `docker ps`. This is similar to how the ps command is used to see a list of processes on a Linux machine.

```console
$ docker ps
CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS              PORTS                    NAMES
ce02b3179f0f        python-docker         "python3 -m flask ru…"   6 minutes ago       Up 6 minutes        0.0.0.0:8000->5000/tcp   wonderful_kalam
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

You can start, stop, and restart Docker containers. When you stop a container, it is not removed, but the status is changed to stopped and the process inside the container is stopped. When you ran the `docker ps` command in the previous module, the default output only shows running containers. When you pass the `--all` or `-a` for short, you see all containers on your machine, irrespective of their start or stop status.

```console
$ docker ps -a
CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS                      PORTS               NAMES
ce02b3179f0f        python-docker         "python3 -m flask ru…"   16 minutes ago      Exited (0) 5 minutes ago                        wonderful_kalam
ec45285c456d        python-docker         "python3 -m flask ru…"   28 minutes ago      Exited (0) 20 minutes ago                       agitated_moser
fb7a41809e5d        python-docker         "python3 -m flask ru…"   37 minutes ago      Exited (0) 36 minutes ago                       goofy_khayyam
```

You should now see several containers listed. These are containers that you started and stopped but have not been removed.

Restart the container that you just stopped. To do this, locate the name of the container you just stopped and replace the name of the container below in the restart command.

```console
$ docker restart wonderful_kalam
```

Now list all the containers again using the `docker ps` command.

```console
$ docker ps --all
CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS                      PORTS                    NAMES
ce02b3179f0f        python-docker         "python3 -m flask ru…"   19 minutes ago      Up 8 seconds                0.0.0.0:8000->5000/tcp   wonderful_kalam
ec45285c456d        python-docker         "python3 -m flask ru…"   31 minutes ago      Exited (0) 23 minutes ago                            agitated_moser
fb7a41809e5d        python-docker         "python3 -m flask ru…"   40 minutes ago      Exited (0) 39 minutes ago                            goofy_khayyam
```

Notice that the container you just restarted has been started in detached mode and has port 8000 exposed. Also, observe the status of the container is “Up X seconds”. When you restart a container, it starts with the same flags or commands that it was originally started with.

Now, stop and remove all of your containers and take a look at fixing the random naming issue. Stop the container you just started. Find the name of your running container and replace the name in the command below with the name of the container on your system.

```console
$ docker stop wonderful_kalam
wonderful_kalam
```

Now that all of your containers are stopped, remove them. When you remove a container, it is no longer running, nor it is in the stopped status, but the process inside the container has been stopped and the metadata for the container has been removed.

```console
$ docker ps --all
CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS                      PORTS               NAMES
ce02b3179f0f        python-docker         "python3 -m flask ru…"   19 minutes ago      Exited (0) 5 seconds ago                        wonderful_kalam
ec45285c456d        python-docker         "python3 -m flask ru…"   31 minutes ago      Exited (0) 23 minutes ago                       agitated_moser
fb7a41809e5d        python-docker         "python3 -m flask ru…"   40 minutes ago      Exited (0) 39 minutes ago                       goofy_khayyam
```

To remove a container, run the `docker rm` command with the container name. You can pass multiple container names to the command using a single command. Again, replace the container names in the following command with the container names from your system.

```console
$ docker rm wonderful_kalam agitated_moser goofy_khayyam
wonderful_kalam
agitated_moser
goofy_khayyam
```

Run the `docker ps --all` command again to see that all containers are removed.

Now, address the random naming issue. Standard practice is to name your containers for the simple reason that it is easier to identify what is running in the container and what application or service it is associated with.

To name a container, you just need to pass the `--name` flag to the `docker run` command.

```console
$ docker run -d -p 8000:5000 --name rest-server python-docker
1aa5d46418a68705c81782a58456a4ccdb56a309cb5e6bd399478d01eaa5cdda
$ docker ps
CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS              PORTS                    NAMES
1aa5d46418a6        python-docker         "python3 -m flask ru…"   3 seconds ago       Up 3 seconds        0.0.0.0:8000->5000/tcp   rest-server
```

That’s better! You can now easily identify your container based on the name.

## Summary

In this section, you took a look at running containers. You also took a look at managing containers by starting, stopping, and restarting them. And finally, you looked at naming your containers so they are more easily identifiable.

Related information:
 - [docker run CLI reference](../../engine/reference/commandline/run.md)

## Next steps

In the next section, you’ll learn how to run a database in a container and connect it to a Python application.

{{< button text="How to develop your application" url="develop.md" >}}
