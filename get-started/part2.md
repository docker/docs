---
title: "Get Started, Part 2: Containers"
keywords: containers, python, code, coding, build, push, run
description: Learn how to write, build, and run a simple app -- the Docker way.
---

{% include_relative nav.html selected="2" %}

## Prerequisites

Part 2 requires that you:

- Complete [Part 1](index.md).
- Sign up for a Docker account (to get a Docker ID) at [Docker Hub](https://hub.docker.com/){: target="_blank" class="_"}.

## Introduction to part 2

This tutorial teaches you how to build a simple Python web application (with
[Flask](http://flask.pocoo.org/){: target="_blank" class="_"}) that prints
**"Hello, Moby Dock!"** on a blue canvas and counts the number of visitors to
the site.

Here, in part 2, we start by building a Docker application image (that we call
`hellomoby`) and run it as a single container on a single node (localhost). We
also tag and push the new image to [Docker Hub](https://hub.docker.com/){: target="_blank" class="_"},
the Docker public reigstry.

Later, in [Part 3](part3.md), we run our image as a "Docker Service" so we can
scale the number of containers and do other operations. In [Part 4](part4.md),
we employ container orchestration and distribute the multiple containers running
our application across multiple nodes. In [Part 5](part5.md), we add a service
for a [Redis](https://hub.docker.com/_/redis/){: target="_blank" class="_"}
database (to store our visitor count) and rebuild our `hellomoby` image as a
Docker stack.

- Stack
- Docker Service
- **Container**  <span class="badge badge-danger">You are here</span>

## Build application image

Typically, before you can create a Python application, you need to install a
Python runtime on your dev environment and configure it _exactly_ like your
production server for the app to run as expected. With Docker, you simply pull
a portable Python runtime image from [Docker Hub](https://hub.docker.com/_/python/){: target="_blank" class="_"}
(or [Docker Store](https://store.docker.com/images/python){: target="_blank" class="_"})
and base your application image on that runtime image.

Docker images are built with
[Dockerfiles](https://docs.docker.com/engine/userguide/eng-image/dockerfile_best-practices/){:
target="_blank" class="_"}. A Dockerfile lets you define the environment
inside a container. Access to resources, such as networks and disk drives, is
virtualized  inside this environment. Because the container is isolated from the
rest of your system, you must map ports to the outside world, and be specific
about what files you want to copy into that environment. Let's start.

1.  Create an empty directory and change directories _into_ this new directory:

    ```shell
    $ mkdir ~/hellomoby && cd ~/hellomoby
    ```

2.  Create a file named `Dockerfile` (no extension) with the following content
    and save:

    ```dockerfile
    ## Dockerfile
    ## Use an official Python runtime as a parent image
    FROM python:2.7-slim

    ## Set the working directory to /app
    WORKDIR /app

    ## Copy the current directory contents into the container at /app
    ADD . /app

    ## Install any needed packages specified in requirements.txt
    RUN pip install -r requirements.txt

    ## Make port 80 available to the world outside this container
    EXPOSE 80

    ## Define environment variable
    ENV NAME Moby Dock

    ## Run app.py when the container launches
    CMD ["python", "app.py"]
    ```

    > Are you behind a proxy server?
    >
    > Proxy servers can block connections to your web app. Add the following lines
    > to `Dockerfile` before `RUN pip ...`:
    >
    > ```conf
    > # Replace host:port with values for your servers
    > ENV http_proxy host:port
    > ENV https_proxy host:port
    > ```

    A Dockerfile defines _how_ to build the Docker image; but we also need to code
    the application itself.

3.  Create two more files, `requirements.txt` and `app.py` (referred to in
    `Dockerfile`), and put them in the same directory:

    ```
    ## requirements.txt
    flask
    redis
    ```

    ```python
    ## app.py
    from flask import Flask
    from redis import Redis, RedisError
    import os
    import socket

    ## Connect to Redis
    redis = Redis(host="redis", db=0, socket_connect_timeout=2, socket_timeout=2)

    app = Flask(__name__)

    @app.route("/")
    def hello():
        try:
            visits = redis.incr("counter")
        except RedisError:
            visits = "<i>Cannot connect to Redis, counter disabled.</i>"

        html = "<h3>Hello, {name}!</h3>" \
               "<body bgcolor={bgcolor}>" \
               "<b>Hostname:</b> {hostname}<br/>" \
               "<b>Visits:</b> {visits}"
        return html.format(name=os.getenv("NAME", "whale"), bgcolor='#72AAFD', hostname=socket.gethostname(), visits=visits)

    if __name__ == "__main__":
        app.run(host='0.0.0.0', port=80)
    ```

4.  List your files to ensure you are still at the top level of your new directory:

    ```shell
    $ ls
    app.py  Dockerfile  requirements.txt
    ```

5.  Build the image with the files in _this_ directory (`.`) and tag it with a
    meaningful name such as `hellomoby`:

    ```shell
    $ docker build --tag hellomoby .
    ```

6.  List the new `hellomoby` image (and the python image it uses) in your local
    Docker registry:

    ```shell
    $ docker image ls
    REPOSITORY      TAG          IMAGE ID           CREATED              SIZE
    hellomoby       latest       3cd4cae1bbeb       2 minutes ago        149MB
    python          2.7-slim     4fd30fc83117       7 weeks ago          138MB
    hello-world     latest       f2a91732366c       2 months ago         1.85kB
    ...
    ```

## Run application image

1.  Run the app with the publish flag (`--publish` or `-p`) to map port 80
    published by the container to port 4000 on your machine:

    ```shell
    $ docker run --publish 4000:80 hellomoby
    ```

    Remapping here demonstrates the difference between what you `EXPOSE` within
    the `Dockerfile`, and what you publish. Python serves your app on port 80
    but that message comes from _inside_ the container, so Python does not know
    you remapped the port to 4000.

2.  Point a browser to [http://localhost:4000](http://localhost:4000).

    > On Windows 7 (with Docker Toolbox), use the Docker Machine IP instead of
    > `localhost`, for example: http://192.168.99.100:4000. To find the IP
    > address, run `docker-machine ip`.

    You should see:

    - **Hello, Moby Dock!**  (on a blue canvas)
    - **Hostname:** `CONTAINER ID`
    - **Visits:** _Cannot connect to Redis, counter disabled._

    ![Hello, Moby Dock in browser](images/hellomoby-web.png){:width="500px"}

    Scroll up and take another look at `Dockerfile`. The `ADD` command bundles
    `app.py` and `requirements.txt` into the build. The `EXPOSE` command makes
    the output from `app.py` accessible over HTTP. The `RUN` command installs
    the Python libraries for Flask and Redis. The environment variable `NAME` is
    printed by `app.py`, as is the output of the call to `socket.gethostname()`,
    which is the container ID.

    > Expected Redis error
    >
    > Redis fails because it is not actually running. We installed the Python
    > library, but not Redis itself. We fix this in [Part 5](part5.md).

3.  Press `CTRL+C` in your terminal to quit.

    > On Windows, press `CTRL+C` and manually stop the container: `docker container stop <CONTAINER ID>`

4.  Run the app again in _detached_ mode (`--detach` or `-d`). The container
    runs in the background and the full CONTAINER ID displays:

    ```shell
    $ docker run --detach --publish 4000:80 hellomoby
    9fc657da2cc76fd130ce141e006717bfdb58783727dd87a98611e63db5ae3a05
    ```

5.  List the running container:

    ```shell
    $ docker container ls
    CONTAINER ID    IMAGE        COMMAND            CREATED           PORTS                   NAMES
    9fc657da2cc7    hellomoby    "python app.py"    29 seconds ago    0.0.0.0:4000->80/tcp    eager_wiles
    ```

6.  List the container again with the `ps` command:

    ```shell
    $ docker ps
    ```

    > Container NAMES
    >
    > The concatenated container NAMES are automatically generated by the moby
    > [names-generator](https://github.com/moby/moby/blob/master/pkg/namesgenerator/names-generator.go).
    > You can set a custom name with `docker run --name <custom_name>`.

7.  Refresh the browser at [http://localhost:4000](http://localhost:4000) and
    notice that "Hostname" reflects the new `CONTAINER ID`.

8.  Stop the container with `CONTAINER ID` or `NAMES` (which in this case, is
    "eager_wiles"--yours will probably differ), for example:

    ```shell
    $ docker container stop eager_wiles
    ```

    > - To stop all running containers: `docker container stop $(docker ps -q)`
    > - To remove _all_ containers: `docker container rm $(docker ps -aq)`

## Share application image

To demonstrate the portability of your application, we upload the `hellomoby`
image to [Docker Hub](https://hub.docker.com/){: target="_blank" class="_"}, delete the
local image, then rerun the application. The Docker CLI automatically pulls the
image from the Docker public registry. If you don't have a Docker account, sign up for one at [Docker Hub](https://hub.docker.com/){: target="_blank" class="_"}.

1.  Log in to the Docker Hub on your local machine.

    ```shell
    $ docker login
    ```

    The notation for associating a local image with a repository on a registry is
    `username/repository:tag`. The tag is optional, but recommended, as it is
    the mechanism that registries use to give Docker images a version.

2.  Create a repository name and tag (and remember to use _your_ Docker account
    username):

    ```shell
    $ docker tag hellomoby username/hellomoby:v1
    ```

    > Use meaningful names for the repository and tag. Here, the image is stored in
    > the `hellomoby` repository (in your Docker account) and tagged as v1 (for
    > version 1).

3.  List your newly tagged image and notice the duplicate IMAGE ID:

    ```shell
    $ docker image ls
    REPOSITORY          TAG         IMAGE ID        CREATED           SIZE
    gordon/hellomoby    v1          3cd4cae1bbeb    18 minutes ago    149MB
    hellomoby           latest      3cd4cae1bbeb    18 minutes ago    149MB
    python              2.7-slim    4fd30fc83117    7 weeks ago       138MB
    hello-world         latest      f2a91732366c    2 months ago      1.85kB
    ...
    ```

4.  Upload your tagged image to Docker Hub: `docker push username/repository:tag`,
    for example:

    ```shell
    $ docker push gordon/hellomoby:v1
    The push refers to repository [docker.io/gordon/hellomoby]
    70f38f9f0633: Pushed
    532baf2e84df: Pushed
    8374fff6617b: Pushed
    94b0b6f67798: Mounted from library/python
    e0c374004259: Mounted from library/python
    56ee7573ea0f: Mounted from library/python
    cfce7a8ae632: Mounted from library/python
    v1: digest: sha256:6d9b791bb8e32fb44b6a00a5a2e65c0beac6442f45de8ece233f310168d16310 size: 1787
    ```

    Once complete, the results of this upload are publicly available.

5.  In a browser, log in to [Docker Hub](https://hub.docker.com/) to see your
    image and its pull command.

## Run image from remote repository

With your application image in the Docker registry, you can run your app on any
machine with `username/repository:tag`. If the image isn't available locally,
Docker pulls it from the repository. Let's test it.

1.  Remove all local `hellomoby` images with the shared IMAGE ID:

    ```shell
    docker image rm -f <IMAGE ID>
    ```

    > To remove all images: `docker image rm -f $(docker image ls -q)`

2.  Run the `hellomoby` image from your remote repository, for example:

    ```shell
    $ docker run -p 4000:80 gordon/hellomoby:v1
    Unable to find image 'gordon/hellomoby:v1' locally
    v1: Pulling from gordon/hellomoby
    c4bb02b17bb4: Already exists
    c5c896dce5ee: Already exists
    cf210b898cc6: Already exists
    5117cef49bdb: Already exists
    16ccfde469ab: Already exists
    d998e18b4542: Already exists
    33cf4efe6bb2: Already exists
    Digest: sha256:6d9b791bb8e32fb44b6a00a5a2e65c0beac6442f45de8ece233f310168d16310
    Status: Downloaded newer image for gordon/hellomoby:v1
     * Running on http://0.0.0.0:80/ (Press CTRL+C to quit)
    ```

3.  Press `CTRL+C` to quit and list your images again:

    ```shell
    $ docker image ls
    ```

## Recap and cheat sheet

```shell
## Create working directory and files:
mkdir ~/hellomoby && cd ~/hellomoby
[create Dockerfile, requirements.txt, app.py]

## Build `hellomoby` image and list:
docker build --tag hellomoby .
docker image ls

## Run `hellomoby` image as container and map ports:
docker run --publish 4000:80 hellomoby
CTRL+C

## Run again in detached mode then stop container (by ID or NAMES):
docker run --detach --publish 4000:80 hellomoby
docker container ls
docker ps
docker container stop NAMES
#docker container stop $(docker ps -q)
#docker container rm $(docker ps -aq)

## Login into Docker Hub, then tag and push `hellomboby` (with your username):
docker login
docker tag hellomoby USERNAME/hellomoby:v1
docker push USERNAME/hellomoby:v1

## Remove local copies of the `hellomoby` image then run from Docker Hub:
docker image rm -f IMAGE ID
docker run -p 4000:80 USERNAME/hellomoby:v1
docker image ls
```

## Conclusion of part 2

No matter where `docker run` executes, it pulls your image, with all its
dependencies, from `requirements.txt` and runs your code. They travel together
in a neat little package and the host machine does not have to install anything
except Docker to run it.

In the next section, we learn how to scale our application by running the
container in a **Docker service**.

[Continue to Part 3 >>](part3.md){: class="button outline-btn"}
