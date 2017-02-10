---
description: Get started with Docker Compose
keywords: documentation, docs, docker, compose, orchestration, containers
title: Get started with Docker Compose
---

On this page you build a simple Python web application running on Docker
Compose. The application uses the Flask framework and maintains a hit counter in
Redis. While the sample uses Python, the concepts demonstrated here should be
understandable even if you're not familiar with it.

## Prerequisites

Make sure you have already [installed both Docker Engine and Docker
Compose](install.md). You don't need to install Python or Redis, as both are
provided by Docker images.

## Step 1: Setup

1.  Create a directory for the project:

        $ mkdir composetest
        $ cd composetest

2.  Create a file called `app.py` in your project directory and paste this in:

        from flask import Flask
        from redis import Redis

        app = Flask(__name__)
        redis = Redis(host='redis', port=6379)

        @app.route('/')
        def hello():
            count = redis.incr('hits')
            return 'Hello World! I have been seen {} times.\n'.format(count)

        if __name__ == "__main__":
            app.run(host="0.0.0.0", debug=True)

3.  Create another file called `requirements.txt` in your project directory and
    paste this in:

        flask
        redis

   These define the application's dependencies.


## Step 2: Create a Dockerfile

In this step, you write a Dockerfile that builds a Docker image. The image
contains all the dependencies the Python application requires, including Python
itself.

In your project directory, create a file named `Dockerfile` and paste the
following:

    FROM python:3.4-alpine
    ADD . /code
    WORKDIR /code
    RUN pip install -r requirements.txt
    CMD ["python", "app.py"]

This tells Docker to:

* Build an image starting with the Python 3.4 image.
* Add the current directory `.` into the path `/code` in the image.
* Set the working directory to `/code`.
* Install the Python dependencies.
* Set the default command for the container to `python app.py`

For more information on how to write Dockerfiles, see the [Docker user
guide](/engine/tutorials/dockerimages.md#building-an-image-from-a-dockerfile)
and the [Dockerfile reference](/engine/reference/builder.md).


## Step 3: Define services in a Compose file

Create a file called `docker-compose.yml` in your project directory and paste
the following:

    version: '2'
    services:
      web:
        build: .
        ports:
         - "5000:5000"
        volumes:
         - .:/code
      redis:
        image: "redis:alpine"

This Compose file defines two services, `web` and `redis`. The web service:

* Uses an image that's built from the `Dockerfile` in the current directory.
* Forwards the exposed port 5000 on the container to port 5000 on the host
  machine.
* Mounts the project directory on the host to `/code` inside the container,
  allowing you to modify the code without having to rebuild the image.

The `redis` service uses a public
[Redis](https://registry.hub.docker.com/_/redis/) image pulled from the Docker
Hub registry.

>**Tip:** If your project is outside of the `Users` directory (`cd ~`), then you
need to share the drive or location of the Dockerfile and volume you are using.
If you get runtime errors indicating an application file is not found, a volume
mount is denied, or a service cannot start, try enabling file or drive sharing.
Volume mounting requires shared drives for projects that live outside of
`C:\Users` (Windows) or `/Users` (Mac), and is required for _any_ project on
Docker for Windows that uses [Linux
containers](/docker-for-windows/index.md#switch-between-windows-and-linux-containers-beta-feature).
For more information, see [Shared
Drives](../docker-for-windows/index.md#shared-drives) on Docker for Windows,
[File sharing](../docker-for-mac/index.md#file-sharing) on Docker for Mac, and
the general examples on how to [Manage data in
containers](../engine/tutorials/dockervolumes.md).

## Step 4: Build and run your app with Compose

1. From your project directory, start up your application.

        $ docker-compose up
        Pulling image redis...
        Building web...
        Starting composetest_redis_1...
        Starting composetest_web_1...
        redis_1 | [8] 02 Jan 18:43:35.576 # Server started, Redis version 2.8.3
        web_1   |  * Running on http://0.0.0.0:5000/
        web_1   |  * Restarting with stat

   Compose pulls a Redis image, builds an image for your code, and start the
   services you defined.

2. Enter `http://0.0.0.0:5000/` in a browser to see the application running.

   If you're using Docker on Linux natively, then the web app should now be
   listening on port 5000 on your Docker daemon host. If `http://0.0.0.0:5000`
   doesn't resolve, you can also try `http://localhost:5000`.

   If you're using Docker Machine on a Mac, use `docker-machine ip MACHINE_VM` to get
   the IP address of your Docker host. Then, `open http://MACHINE_VM_IP:5000` in a
   browser.

   You should see a message in your browser saying:

   `Hello World! I have been seen 1 times.`

3. Refresh the page.

   The number should increment.

>**Tip:** You can list local images with `docker image ls` and inspect them with `docker inspect <tag or id>`. Listing images at this point should return `redis` and `web`.


## Step 5: Update the application

Because the application code is mounted into the container using a volume, you
can make changes to its code and see the changes instantly, without having to
rebuild the image.

1.  Change the greeting in `app.py` and save it. For example:

        return 'Hello from Docker! I have been seen {} times.\n'.format(count)

2.  Refresh the app in your browser. The greeting should be updated, and the
    counter should still be incrementing.


## Step 6: Experiment with some other commands

If you want to run your services in the background, you can pass the `-d` flag
(for "detached" mode) to `docker-compose up` and use `docker-compose ps` to
see what is currently running:

    $ docker-compose up -d
    Starting composetest_redis_1...
    Starting composetest_web_1...

    $ docker-compose ps
    Name                 Command            State       Ports
    -------------------------------------------------------------------
    composetest_redis_1   /usr/local/bin/run         Up
    composetest_web_1     /bin/sh -c python app.py   Up      5000->5000/tcp

The `docker-compose run` command allows you to run one-off commands for your
services. For example, to see what environment variables are available to the
`web` service:

    $ docker-compose run web env

See `docker-compose --help` to see other available commands. You can also install [command completion](completion.md) for the bash and zsh shell, which will also show you available commands.

If you started Compose with `docker-compose up -d`, you'll probably want to stop
your services once you've finished with them:

    $ docker-compose stop

You can bring everything down, removing the containers entirely, with the `down`
command. Pass `--volumes` to also remove the data volume used by the Redis
container:

    $ docker-compose down --volumes

At this point, you have seen the basics of how Compose works.


## Where to go next

- Next, try the quick start guide for [Django](django.md),
  [Rails](rails.md), or [WordPress](wordpress.md).
- [Explore the full list of Compose commands](./reference/index.md)
- [Compose configuration file reference](compose-file.md)
